package evaluator

import (
	"context"
	"fmt"
	"strings"

	"github.com/krateoplatformops/krateo-bff/apis/core"
	"github.com/krateoplatformops/krateo-bff/apis/ui/formtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/api"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/crdutil"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/endpoints"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/formdefinitions"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	rbacutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/rbac/util"
	"github.com/krateoplatformops/krateo-bff/internal/tmpl"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
)

const (
	filter = `.spec.versions[] | select(.name="%s") | .schema.openAPIV3Schema`
)

type EvalOptions struct {
	RESTConfig *rest.Config
	AuthnNS    string
	Subject    string
	Groups     []string
}

func Eval(ctx context.Context, in *v1alpha1.FormTemplate, opts EvalOptions) error {
	tpl, err := tmpl.New("${", "}")
	if err != nil {
		return err
	}

	ds, err := callAPIs(ctx, callAPIsOptions{
		restConfig: opts.RESTConfig,
		authnNS:    opts.AuthnNS,
		subject:    opts.Subject,
		tpl:        tpl,
		apiList:    in.Spec.APIList,
	})
	_ = ds // TODO

	// TODO(@lucasepe): resolve `definitionRef` and get - G=schema.group, V=schema.version, K=schema.kind
	formdefinitionsClient, err := formdefinitions.NewClient(opts.RESTConfig)
	if err != nil {
		return err
	}

	ref, err := formdefinitionsClient.Namespace(in.Spec.DefinitionRef.Namespace).
		Get(ctx, in.Spec.DefinitionRef.Name)
	if err != nil {
		return err
	}

	// TODO(@lucasepe): with GVK find GR and resolve `resourceRef`
	dyn, err := dynamic.NewGetter(opts.RESTConfig)
	if err != nil {
		return err
	}

	src, err := dyn.Get(ctx, dynamic.GetOptions{
		GVK: schema.GroupVersionKind{
			Group:   ref.Spec.Schema.Group,
			Version: ref.Spec.Schema.Version,
			Kind:    ref.Spec.Schema.Kind,
		},
		Namespace: ref.Namespace,
		Name:      ref.Name,
	})
	if err != nil {
		return err
	}

	gr := crdutil.InferGroupResource(schema.GroupKind{
		Group: ref.Spec.Schema.Group,
		Kind:  ref.Spec.Schema.Kind,
	})

	crd, err := dyn.Get(ctx, dynamic.GetOptions{
		GVK: schema.GroupVersionKind{
			Group:   "apiextensions.k8s.io",
			Version: "v1",
			Kind:    "CustomResourceDefinition",
		},
		Name: gr.String(),
	})
	if err != nil {
		return err
	}

	sch, err := dynamic.Extract(ctx, crd, fmt.Sprintf(filter, ref.Spec.Schema.Version))
	if err != nil {
		return err
	}

	in.Status.Content = &v1alpha1.FormTemplateStatusContent{
		Instance: &runtime.RawExtension{Object: src},
		Schema: &runtime.RawExtension{Object: &unstructured.Unstructured{
			Object: sch.(map[string]any),
		}},
	}

	// TODO(@lucasepe): eventually evaluate jq templates in `data`
	// TODO(@lucasepe): merge `resourceRef` specs with evaluated `data`
	// TODO(@lucasepe): apply to k8s the new evaluated `resource`
	// TODO(@lucasepe): in.Status.Schema = &runtime.RawExtension{Object: CRD} (store the CRD)
	// TODO(@lucasepe): in.Status.Instance = &runtime.RawExtension{Object: ""} (store the evaluated resource)

	return nil
	// return injectAllowedVerbs(in, allowedVerbsInjectorOptions{
	// 	restConfig: opts.RESTConfig,
	// 	subject:    opts.Subject,
	// 	groups:     opts.Groups,
	// })
}

type callAPIsOptions struct {
	restConfig *rest.Config
	authnNS    string
	subject    string
	tpl        tmpl.JQTemplate
	apiList    []*core.API
}

func callAPIs(ctx context.Context, opts callAPIsOptions) (map[string]any, error) {
	apiMap := map[string]*core.API{}
	for _, x := range opts.apiList {
		apiMap[x.Name] = x
	}

	sorted := core.SortApiByDeps(opts.apiList)

	ds := map[string]any{}
	for _, key := range sorted {
		x, ok := apiMap[key]
		if !ok {
			return nil, fmt.Errorf("API '%s' not found in apiMap", key)
		}

		ref := x.EndpointRef
		if ptr.Deref(x.KrateoGateway, false) {
			ref = &core.Reference{
				Name:      fmt.Sprintf("%s-clientconfig", opts.subject),
				Namespace: opts.authnNS,
			}
		}

		ep, err := endpoints.Resolve(context.TODO(), opts.restConfig, ref)
		if err != nil {
			return nil, err
		}

		hc, err := api.HTTPClientForEndpoint(ep)
		if err != nil {
			return nil, err
		}

		rt, err := api.Call(ctx, hc, api.CallOptions{
			API:      x,
			Endpoint: ep,
			Tpl:      opts.tpl,
			DS:       ds,
		})
		if err != nil {
			return nil, err
		}

		ds[x.Name] = rt
	}

	return ds, nil
}

const (
	allowedVerbsAnnotationKey = "krateo.io/allowed-verbs"
	resource                  = "formtemplates"
)

type allowedVerbsInjectorOptions struct {
	restConfig *rest.Config
	subject    string
	groups     []string
}

func injectAllowedVerbs(in *v1alpha1.FormTemplate, opts allowedVerbsInjectorOptions) error {
	verbs, err := rbacutil.GetAllowedVerbs(context.TODO(), opts.restConfig, util.ResourceInfo{
		Subject: opts.subject,
		Groups:  opts.groups,
		GroupResource: v1alpha1.FormTemplateGroupVersionKind.GroupVersion().
			WithResource(resource).
			GroupResource(),
		ResourceName: in.GetName(),
		Namespace:    in.GetNamespace(),
	})
	if err != nil {
		return err
	}

	m := in.GetAnnotations()
	if len(m) == 0 {
		m = map[string]string{}
	}
	m[allowedVerbsAnnotationKey] = strings.Join(verbs, ",")
	in.SetAnnotations(m)

	return nil
}
