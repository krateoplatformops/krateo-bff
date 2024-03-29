package util

type PolicyRule interface {
	// Verbs is a list of Verbs that apply to ALL the ResourceKinds contained in this rule. '*' represents all verbs.
	Verbs() []string
	// APIGroups is the name of the APIGroup that contains the resources.  If multiple API groups are specified, any action requested against one of
	// the enumerated resources in any API group will be allowed. "" represents the core API group and "*" represents all API groups.
	APIGroups() []string
	// Resources is a list of resources this rule applies to. '*' represents all resources.
	Resources() []string
	// ResourceNames is an optional white list of names that the rule applies to.  An empty set means that everything is allowed.
	ResourceNames() []string
}

var _ PolicyRule = (*policyRule)(nil)

type policyRule struct {
	verbs         []string
	apiGroups     []string
	resources     []string
	resourceNames []string
}

func (pr *policyRule) Verbs() []string {
	return pr.verbs
}

func (pr *policyRule) APIGroups() []string {
	return pr.apiGroups
}

func (pr *policyRule) Resources() []string {
	return pr.resources
}

func (pr *policyRule) ResourceNames() []string {
	return pr.resourceNames
}
