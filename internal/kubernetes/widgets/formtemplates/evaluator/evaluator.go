package evaluator

import (
	"context"
	"fmt"

	"github.com/krateoplatformops/krateo-bff/apis/ui/formtemplates/v1alpha1"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/dynamic"
	"github.com/krateoplatformops/krateo-bff/internal/kubernetes/formdefinitions"
	formdefinitionsutil "github.com/krateoplatformops/krateo-bff/internal/kubernetes/formdefinitions/util"
	"github.com/krateoplatformops/krateo-bff/internal/tmpl"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
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

	formdefinitionsClient, err := formdefinitions.NewClient(opts.RESTConfig)
	if err != nil {
		return err
	}

	ref, err := formdefinitionsClient.Namespace(in.Spec.DefinitionRef.Namespace).
		Get(ctx, in.Spec.DefinitionRef.Name)
	if err != nil {
		return err
	}

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

	gr := formdefinitionsutil.InferGroupResource(ref)

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

	return nil
}
