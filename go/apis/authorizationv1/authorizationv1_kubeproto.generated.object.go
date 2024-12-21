package authorizationv1

import (
	"go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "authorization.k8s.io"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&LocalSubjectAccessReview{},
		&LocalSubjectAccessReviewList{},
		&SelfSubjectAccessReview{},
		&SelfSubjectAccessReviewList{},
		&SelfSubjectRulesReview{},
		&SelfSubjectRulesReviewList{},
		&SubjectAccessReview{},
		&SubjectAccessReviewList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type LocalSubjectAccessReview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec holds information about the request being evaluated.  spec.namespace must be equal to the namespace
	// you made the request against.  If empty, it is defaulted.
	Spec SubjectAccessReviewSpec `json:"spec"`
	// Status is filled in by the server and indicates whether the request is allowed or not
	Status *SubjectAccessReviewStatus `json:"status,omitempty"`
}

func (in *LocalSubjectAccessReview) DeepCopyInto(out *LocalSubjectAccessReview) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(SubjectAccessReviewStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *LocalSubjectAccessReview) DeepCopy() *LocalSubjectAccessReview {
	if in == nil {
		return nil
	}
	out := new(LocalSubjectAccessReview)
	in.DeepCopyInto(out)
	return out
}

func (in *LocalSubjectAccessReview) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type LocalSubjectAccessReviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []LocalSubjectAccessReview `json:"items"`
}

func (in *LocalSubjectAccessReviewList) DeepCopyInto(out *LocalSubjectAccessReviewList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]LocalSubjectAccessReview, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *LocalSubjectAccessReviewList) DeepCopy() *LocalSubjectAccessReviewList {
	if in == nil {
		return nil
	}
	out := new(LocalSubjectAccessReviewList)
	in.DeepCopyInto(out)
	return out
}

func (in *LocalSubjectAccessReviewList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type SelfSubjectAccessReview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec holds information about the request being evaluated.  user and groups must be empty
	Spec SelfSubjectAccessReviewSpec `json:"spec"`
	// Status is filled in by the server and indicates whether the request is allowed or not
	Status *SubjectAccessReviewStatus `json:"status,omitempty"`
}

func (in *SelfSubjectAccessReview) DeepCopyInto(out *SelfSubjectAccessReview) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(SubjectAccessReviewStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *SelfSubjectAccessReview) DeepCopy() *SelfSubjectAccessReview {
	if in == nil {
		return nil
	}
	out := new(SelfSubjectAccessReview)
	in.DeepCopyInto(out)
	return out
}

func (in *SelfSubjectAccessReview) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type SelfSubjectAccessReviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []SelfSubjectAccessReview `json:"items"`
}

func (in *SelfSubjectAccessReviewList) DeepCopyInto(out *SelfSubjectAccessReviewList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]SelfSubjectAccessReview, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *SelfSubjectAccessReviewList) DeepCopy() *SelfSubjectAccessReviewList {
	if in == nil {
		return nil
	}
	out := new(SelfSubjectAccessReviewList)
	in.DeepCopyInto(out)
	return out
}

func (in *SelfSubjectAccessReviewList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type SelfSubjectRulesReview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec holds information about the request being evaluated.
	Spec SelfSubjectRulesReviewSpec `json:"spec"`
	// Status is filled in by the server and indicates the set of actions a user can perform.
	Status *SubjectRulesReviewStatus `json:"status,omitempty"`
}

func (in *SelfSubjectRulesReview) DeepCopyInto(out *SelfSubjectRulesReview) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(SubjectRulesReviewStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *SelfSubjectRulesReview) DeepCopy() *SelfSubjectRulesReview {
	if in == nil {
		return nil
	}
	out := new(SelfSubjectRulesReview)
	in.DeepCopyInto(out)
	return out
}

func (in *SelfSubjectRulesReview) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type SelfSubjectRulesReviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []SelfSubjectRulesReview `json:"items"`
}

func (in *SelfSubjectRulesReviewList) DeepCopyInto(out *SelfSubjectRulesReviewList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]SelfSubjectRulesReview, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *SelfSubjectRulesReviewList) DeepCopy() *SelfSubjectRulesReviewList {
	if in == nil {
		return nil
	}
	out := new(SelfSubjectRulesReviewList)
	in.DeepCopyInto(out)
	return out
}

func (in *SelfSubjectRulesReviewList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type SubjectAccessReview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec holds information about the request being evaluated
	Spec SubjectAccessReviewSpec `json:"spec"`
	// Status is filled in by the server and indicates whether the request is allowed or not
	Status *SubjectAccessReviewStatus `json:"status,omitempty"`
}

func (in *SubjectAccessReview) DeepCopyInto(out *SubjectAccessReview) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(SubjectAccessReviewStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *SubjectAccessReview) DeepCopy() *SubjectAccessReview {
	if in == nil {
		return nil
	}
	out := new(SubjectAccessReview)
	in.DeepCopyInto(out)
	return out
}

func (in *SubjectAccessReview) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type SubjectAccessReviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []SubjectAccessReview `json:"items"`
}

func (in *SubjectAccessReviewList) DeepCopyInto(out *SubjectAccessReviewList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]SubjectAccessReview, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *SubjectAccessReviewList) DeepCopy() *SubjectAccessReviewList {
	if in == nil {
		return nil
	}
	out := new(SubjectAccessReviewList)
	in.DeepCopyInto(out)
	return out
}

func (in *SubjectAccessReviewList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type SubjectAccessReviewSpec struct {
	// ResourceAuthorizationAttributes describes information for a resource access request
	ResourceAttributes *ResourceAttributes `json:"resourceAttributes,omitempty"`
	// NonResourceAttributes describes information for a non-resource access request
	NonResourceAttributes *NonResourceAttributes `json:"nonResourceAttributes,omitempty"`
	// User is the user you're testing for.
	// If you specify "User" but not "Groups", then is it interpreted as "What if User were not a member of any groups
	User string `json:"user,omitempty"`
	// Groups is the groups you're testing for.
	Groups []string `json:"groups"`
	// Extra corresponds to the user.Info.GetExtra() method from the authenticator.  Since that is input to the authorizer
	// it needs a reflection here.
	Extra map[string]ExtraValue `json:"extra,omitempty"`
	// UID information about the requesting user.
	UID string `json:"uid,omitempty"`
}

func (in *SubjectAccessReviewSpec) DeepCopyInto(out *SubjectAccessReviewSpec) {
	*out = *in
	if in.ResourceAttributes != nil {
		in, out := &in.ResourceAttributes, &out.ResourceAttributes
		*out = new(ResourceAttributes)
		(*in).DeepCopyInto(*out)
	}
	if in.NonResourceAttributes != nil {
		in, out := &in.NonResourceAttributes, &out.NonResourceAttributes
		*out = new(NonResourceAttributes)
		(*in).DeepCopyInto(*out)
	}
	if in.Groups != nil {
		t := make([]string, len(in.Groups))
		copy(t, in.Groups)
		out.Groups = t
	}
	if in.Extra != nil {
		in, out := &in.Extra, &out.Extra
		*out = make(map[string]ExtraValue, len(*in))
		for k, v := range *in {
			(*out)[k] = v
		}
	}
}

func (in *SubjectAccessReviewSpec) DeepCopy() *SubjectAccessReviewSpec {
	if in == nil {
		return nil
	}
	out := new(SubjectAccessReviewSpec)
	in.DeepCopyInto(out)
	return out
}

type SubjectAccessReviewStatus struct {
	// Allowed is required. True if the action would be allowed, false otherwise.
	Allowed bool `json:"allowed"`
	// Denied is optional. True if the action would be denied, otherwise
	// false. If both allowed is false and denied is false, then the
	// authorizer has no opinion on whether to authorize the action. Denied
	// may not be true if Allowed is true.
	Denied bool `json:"denied,omitempty"`
	// Reason is optional.  It indicates why a request was allowed or denied.
	Reason string `json:"reason,omitempty"`
	// EvaluationError is an indication that some error occurred during the authorization check.
	// It is entirely possible to get an error and be able to continue determine authorization status in spite of it.
	// For instance, RBAC can be missing a role, but enough roles are still present and bound to reason about the request.
	EvaluationError string `json:"evaluationError,omitempty"`
}

func (in *SubjectAccessReviewStatus) DeepCopyInto(out *SubjectAccessReviewStatus) {
	*out = *in
}

func (in *SubjectAccessReviewStatus) DeepCopy() *SubjectAccessReviewStatus {
	if in == nil {
		return nil
	}
	out := new(SubjectAccessReviewStatus)
	in.DeepCopyInto(out)
	return out
}

type SelfSubjectAccessReviewSpec struct {
	// ResourceAuthorizationAttributes describes information for a resource access request
	ResourceAttributes *ResourceAttributes `json:"resourceAttributes,omitempty"`
	// NonResourceAttributes describes information for a non-resource access request
	NonResourceAttributes *NonResourceAttributes `json:"nonResourceAttributes,omitempty"`
}

func (in *SelfSubjectAccessReviewSpec) DeepCopyInto(out *SelfSubjectAccessReviewSpec) {
	*out = *in
	if in.ResourceAttributes != nil {
		in, out := &in.ResourceAttributes, &out.ResourceAttributes
		*out = new(ResourceAttributes)
		(*in).DeepCopyInto(*out)
	}
	if in.NonResourceAttributes != nil {
		in, out := &in.NonResourceAttributes, &out.NonResourceAttributes
		*out = new(NonResourceAttributes)
		(*in).DeepCopyInto(*out)
	}
}

func (in *SelfSubjectAccessReviewSpec) DeepCopy() *SelfSubjectAccessReviewSpec {
	if in == nil {
		return nil
	}
	out := new(SelfSubjectAccessReviewSpec)
	in.DeepCopyInto(out)
	return out
}

type SelfSubjectRulesReviewSpec struct {
	// Namespace to evaluate rules for. Required.
	Namespace string `json:"namespace,omitempty"`
}

func (in *SelfSubjectRulesReviewSpec) DeepCopyInto(out *SelfSubjectRulesReviewSpec) {
	*out = *in
}

func (in *SelfSubjectRulesReviewSpec) DeepCopy() *SelfSubjectRulesReviewSpec {
	if in == nil {
		return nil
	}
	out := new(SelfSubjectRulesReviewSpec)
	in.DeepCopyInto(out)
	return out
}

type SubjectRulesReviewStatus struct {
	// ResourceRules is the list of actions the subject is allowed to perform on resources.
	// The list ordering isn't significant, may contain duplicates, and possibly be incomplete.
	ResourceRules []ResourceRule `json:"resourceRules"`
	// NonResourceRules is the list of actions the subject is allowed to perform on non-resources.
	// The list ordering isn't significant, may contain duplicates, and possibly be incomplete.
	NonResourceRules []NonResourceRule `json:"nonResourceRules"`
	// Incomplete is true when the rules returned by this call are incomplete. This is most commonly
	// encountered when an authorizer, such as an external authorizer, doesn't support rules evaluation.
	Incomplete bool `json:"incomplete"`
	// EvaluationError can appear in combination with Rules. It indicates an error occurred during
	// rule evaluation, such as an authorizer that doesn't support rule evaluation, and that
	// ResourceRules and/or NonResourceRules may be incomplete.
	EvaluationError string `json:"evaluationError,omitempty"`
}

func (in *SubjectRulesReviewStatus) DeepCopyInto(out *SubjectRulesReviewStatus) {
	*out = *in
	if in.ResourceRules != nil {
		l := make([]ResourceRule, len(in.ResourceRules))
		for i := range in.ResourceRules {
			in.ResourceRules[i].DeepCopyInto(&l[i])
		}
		out.ResourceRules = l
	}
	if in.NonResourceRules != nil {
		l := make([]NonResourceRule, len(in.NonResourceRules))
		for i := range in.NonResourceRules {
			in.NonResourceRules[i].DeepCopyInto(&l[i])
		}
		out.NonResourceRules = l
	}
}

func (in *SubjectRulesReviewStatus) DeepCopy() *SubjectRulesReviewStatus {
	if in == nil {
		return nil
	}
	out := new(SubjectRulesReviewStatus)
	in.DeepCopyInto(out)
	return out
}

type ResourceAttributes struct {
	// Namespace is the namespace of the action being requested.  Currently, there is no distinction between no namespace and all namespaces
	// "" (empty) is defaulted for LocalSubjectAccessReviews
	// "" (empty) is empty for cluster-scoped resources
	// "" (empty) means "all" for namespace scoped resources from a SubjectAccessReview or SelfSubjectAccessReview
	Namespace string `json:"namespace,omitempty"`
	// Verb is a kubernetes resource API verb, like: get, list, watch, create, update, delete, proxy.  "*" means all.
	Verb string `json:"verb,omitempty"`
	// Group is the API Group of the Resource.  "*" means all.
	Group string `json:"group,omitempty"`
	// Version is the API Version of the Resource.  "*" means all.
	Version string `json:"version,omitempty"`
	// Resource is one of the existing resource types.  "*" means all.
	Resource string `json:"resource,omitempty"`
	// Subresource is one of the existing resource types.  "" means none.
	Subresource string `json:"subresource,omitempty"`
	// Name is the name of the resource being requested for a "get" or deleted for a "delete". "" (empty) means all.
	Name string `json:"name,omitempty"`
	// fieldSelector describes the limitation on access based on field.  It can only limit access, not broaden it.
	// This field  is alpha-level. To use this field, you must enable the
	// `AuthorizeWithSelectors` feature gate (disabled by default).
	FieldSelector *FieldSelectorAttributes `json:"fieldSelector,omitempty"`
	// labelSelector describes the limitation on access based on labels.  It can only limit access, not broaden it.
	// This field  is alpha-level. To use this field, you must enable the
	// `AuthorizeWithSelectors` feature gate (disabled by default).
	LabelSelector *LabelSelectorAttributes `json:"labelSelector,omitempty"`
}

func (in *ResourceAttributes) DeepCopyInto(out *ResourceAttributes) {
	*out = *in
	if in.FieldSelector != nil {
		in, out := &in.FieldSelector, &out.FieldSelector
		*out = new(FieldSelectorAttributes)
		(*in).DeepCopyInto(*out)
	}
	if in.LabelSelector != nil {
		in, out := &in.LabelSelector, &out.LabelSelector
		*out = new(LabelSelectorAttributes)
		(*in).DeepCopyInto(*out)
	}
}

func (in *ResourceAttributes) DeepCopy() *ResourceAttributes {
	if in == nil {
		return nil
	}
	out := new(ResourceAttributes)
	in.DeepCopyInto(out)
	return out
}

type NonResourceAttributes struct {
	// Path is the URL path of the request
	Path string `json:"path,omitempty"`
	// Verb is the standard HTTP verb
	Verb string `json:"verb,omitempty"`
}

func (in *NonResourceAttributes) DeepCopyInto(out *NonResourceAttributes) {
	*out = *in
}

func (in *NonResourceAttributes) DeepCopy() *NonResourceAttributes {
	if in == nil {
		return nil
	}
	out := new(NonResourceAttributes)
	in.DeepCopyInto(out)
	return out
}

type ExtraValue []string

type ResourceRule struct {
	// Verb is a list of kubernetes resource API verbs, like: get, list, watch, create, update, delete, proxy.  "*" means all.
	Verbs []string `json:"verbs"`
	// APIGroups is the name of the APIGroup that contains the resources.  If multiple API groups are specified, any action requested against one of
	// the enumerated resources in any API group will be allowed.  "*" means all.
	APIGroups []string `json:"apiGroups"`
	// Resources is a list of resources this rule applies to.  "*" means all in the specified apiGroups.
	// "*/foo" represents the subresource 'foo' for all resources in the specified apiGroups.
	Resources []string `json:"resources"`
	// ResourceNames is an optional white list of names that the rule applies to.  An empty set means that everything is allowed.  "*" means all.
	ResourceNames []string `json:"resourceNames"`
}

func (in *ResourceRule) DeepCopyInto(out *ResourceRule) {
	*out = *in
	if in.Verbs != nil {
		t := make([]string, len(in.Verbs))
		copy(t, in.Verbs)
		out.Verbs = t
	}
	if in.APIGroups != nil {
		t := make([]string, len(in.APIGroups))
		copy(t, in.APIGroups)
		out.APIGroups = t
	}
	if in.Resources != nil {
		t := make([]string, len(in.Resources))
		copy(t, in.Resources)
		out.Resources = t
	}
	if in.ResourceNames != nil {
		t := make([]string, len(in.ResourceNames))
		copy(t, in.ResourceNames)
		out.ResourceNames = t
	}
}

func (in *ResourceRule) DeepCopy() *ResourceRule {
	if in == nil {
		return nil
	}
	out := new(ResourceRule)
	in.DeepCopyInto(out)
	return out
}

type NonResourceRule struct {
	// Verb is a list of kubernetes non-resource API verbs, like: get, post, put, delete, patch, head, options.  "*" means all.
	Verbs []string `json:"verbs"`
	// NonResourceURLs is a set of partial urls that a user should have access to.  *s are allowed, but only as the full,
	// final step in the path.  "*" means all.
	NonResourceURLs []string `json:"nonResourceURLs"`
}

func (in *NonResourceRule) DeepCopyInto(out *NonResourceRule) {
	*out = *in
	if in.Verbs != nil {
		t := make([]string, len(in.Verbs))
		copy(t, in.Verbs)
		out.Verbs = t
	}
	if in.NonResourceURLs != nil {
		t := make([]string, len(in.NonResourceURLs))
		copy(t, in.NonResourceURLs)
		out.NonResourceURLs = t
	}
}

func (in *NonResourceRule) DeepCopy() *NonResourceRule {
	if in == nil {
		return nil
	}
	out := new(NonResourceRule)
	in.DeepCopyInto(out)
	return out
}

type FieldSelectorAttributes struct {
	// rawSelector is the serialization of a field selector that would be included in a query parameter.
	// Webhook implementations are encouraged to ignore rawSelector.
	// The kube-apiserver's *SubjectAccessReview will parse the rawSelector as long as the requirements are not present.
	RawSelector string `json:"rawSelector,omitempty"`
	// requirements is the parsed interpretation of a field selector.
	// All requirements must be met for a resource instance to match the selector.
	// Webhook implementations should handle requirements, but how to handle them is up to the webhook.
	// Since requirements can only limit the request, it is safe to authorize as unlimited request if the requirements
	// are not understood.
	Requirements []metav1.FieldSelectorRequirement `json:"requirements"`
}

func (in *FieldSelectorAttributes) DeepCopyInto(out *FieldSelectorAttributes) {
	*out = *in
	if in.Requirements != nil {
		l := make([]metav1.FieldSelectorRequirement, len(in.Requirements))
		for i := range in.Requirements {
			in.Requirements[i].DeepCopyInto(&l[i])
		}
		out.Requirements = l
	}
}

func (in *FieldSelectorAttributes) DeepCopy() *FieldSelectorAttributes {
	if in == nil {
		return nil
	}
	out := new(FieldSelectorAttributes)
	in.DeepCopyInto(out)
	return out
}

type LabelSelectorAttributes struct {
	// rawSelector is the serialization of a field selector that would be included in a query parameter.
	// Webhook implementations are encouraged to ignore rawSelector.
	// The kube-apiserver's *SubjectAccessReview will parse the rawSelector as long as the requirements are not present.
	RawSelector string `json:"rawSelector,omitempty"`
	// requirements is the parsed interpretation of a label selector.
	// All requirements must be met for a resource instance to match the selector.
	// Webhook implementations should handle requirements, but how to handle them is up to the webhook.
	// Since requirements can only limit the request, it is safe to authorize as unlimited request if the requirements
	// are not understood.
	Requirements []metav1.LabelSelectorRequirement `json:"requirements"`
}

func (in *LabelSelectorAttributes) DeepCopyInto(out *LabelSelectorAttributes) {
	*out = *in
	if in.Requirements != nil {
		l := make([]metav1.LabelSelectorRequirement, len(in.Requirements))
		for i := range in.Requirements {
			in.Requirements[i].DeepCopyInto(&l[i])
		}
		out.Requirements = l
	}
}

func (in *LabelSelectorAttributes) DeepCopy() *LabelSelectorAttributes {
	if in == nil {
		return nil
	}
	out := new(LabelSelectorAttributes)
	in.DeepCopyInto(out)
	return out
}
