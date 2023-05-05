package networkingv1

import (
	corev1 "go.f110.dev/kubeproto/go/apis/corev1"
	metav1 "go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilintstr "k8s.io/apimachinery/pkg/util/intstr"
)

const GroupName = "networking.k8s.io"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&Ingress{},
		&IngressClass{},
		&IngressClassList{},
		&IngressList{},
		&NetworkPolicy{},
		&NetworkPolicyList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type NetworkPolicyConditionReason string

const (
	NetworkPolicyConditionReasonFeatureNotSupported NetworkPolicyConditionReason = "FeatureNotSupported"
)

type NetworkPolicyConditionType string

const (
	NetworkPolicyConditionTypeAccepted       NetworkPolicyConditionType = "Accepted"
	NetworkPolicyConditionTypePartialFailure NetworkPolicyConditionType = "PartialFailure"
	NetworkPolicyConditionTypeFailure        NetworkPolicyConditionType = "Failure"
)

type PolicyType string

const (
	PolicyTypeIngress PolicyType = "Ingress"
	PolicyTypeEgress  PolicyType = "Egress"
)

type Ingress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec is the desired state of the Ingress.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec *IngressSpec `json:"spec,omitempty"`
	// Status is the current state of the Ingress.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Status *IngressStatus `json:"status,omitempty"`
}

func (in *Ingress) DeepCopyInto(out *Ingress) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(IngressSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(IngressStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *Ingress) DeepCopy() *Ingress {
	if in == nil {
		return nil
	}
	out := new(Ingress)
	in.DeepCopyInto(out)
	return out
}

func (in *Ingress) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type IngressClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec is the desired state of the IngressClass.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec *IngressClassSpec `json:"spec,omitempty"`
}

func (in *IngressClass) DeepCopyInto(out *IngressClass) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(IngressClassSpec)
		(*in).DeepCopyInto(*out)
	}
}

func (in *IngressClass) DeepCopy() *IngressClass {
	if in == nil {
		return nil
	}
	out := new(IngressClass)
	in.DeepCopyInto(out)
	return out
}

func (in *IngressClass) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type IngressClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []IngressClass `json:"items"`
}

func (in *IngressClassList) DeepCopyInto(out *IngressClassList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]IngressClass, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *IngressClassList) DeepCopy() *IngressClassList {
	if in == nil {
		return nil
	}
	out := new(IngressClassList)
	in.DeepCopyInto(out)
	return out
}

func (in *IngressClassList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type IngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Ingress `json:"items"`
}

func (in *IngressList) DeepCopyInto(out *IngressList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]Ingress, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *IngressList) DeepCopy() *IngressList {
	if in == nil {
		return nil
	}
	out := new(IngressList)
	in.DeepCopyInto(out)
	return out
}

func (in *IngressList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type NetworkPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Specification of the desired behavior for this NetworkPolicy.
	Spec *NetworkPolicySpec `json:"spec,omitempty"`
	// Status is the current state of the NetworkPolicy.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Status *NetworkPolicyStatus `json:"status,omitempty"`
}

func (in *NetworkPolicy) DeepCopyInto(out *NetworkPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(NetworkPolicySpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(NetworkPolicyStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *NetworkPolicy) DeepCopy() *NetworkPolicy {
	if in == nil {
		return nil
	}
	out := new(NetworkPolicy)
	in.DeepCopyInto(out)
	return out
}

func (in *NetworkPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type NetworkPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []NetworkPolicy `json:"items"`
}

func (in *NetworkPolicyList) DeepCopyInto(out *NetworkPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]NetworkPolicy, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *NetworkPolicyList) DeepCopy() *NetworkPolicyList {
	if in == nil {
		return nil
	}
	out := new(NetworkPolicyList)
	in.DeepCopyInto(out)
	return out
}

func (in *NetworkPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type IngressSpec struct {
	// IngressClassName is the name of the IngressClass cluster resource. The
	// associated IngressClass defines which controller will implement the
	// resource. This replaces the deprecated `kubernetes.io/ingress.class`
	// annotation. For backwards compatibility, when that annotation is set, it
	// must be given precedence over this field. The controller may emit a
	// warning if the field and annotation have different values.
	// Implementations of this API should ignore Ingresses without a class
	// specified. An IngressClass resource may be marked as default, which can
	// be used to set a default value for this field. For more information,
	// refer to the IngressClass documentation.
	IngressClassName string `json:"ingressClassName,omitempty"`
	// DefaultBackend is the backend that should handle requests that don't
	// match any rule. If Rules are not specified, DefaultBackend must be specified.
	// If DefaultBackend is not set, the handling of requests that do not match any
	// of the rules will be up to the Ingress controller.
	DefaultBackend *IngressBackend `json:"defaultBackend,omitempty"`
	// TLS configuration. Currently the Ingress only supports a single TLS
	// port, 443. If multiple members of this list specify different hosts, they
	// will be multiplexed on the same port according to the hostname specified
	// through the SNI TLS extension, if the ingress controller fulfilling the
	// ingress supports SNI.
	TLS []IngressTLS `json:"tls"`
	// A list of host rules used to configure the Ingress. If unspecified, or
	// no rule matches, all traffic is sent to the default backend.
	Rules []IngressRule `json:"rules"`
}

func (in *IngressSpec) DeepCopyInto(out *IngressSpec) {
	*out = *in
	if in.DefaultBackend != nil {
		in, out := &in.DefaultBackend, &out.DefaultBackend
		*out = new(IngressBackend)
		(*in).DeepCopyInto(*out)
	}
	if in.TLS != nil {
		l := make([]IngressTLS, len(in.TLS))
		for i := range in.TLS {
			in.TLS[i].DeepCopyInto(&l[i])
		}
		out.TLS = l
	}
	if in.Rules != nil {
		l := make([]IngressRule, len(in.Rules))
		for i := range in.Rules {
			in.Rules[i].DeepCopyInto(&l[i])
		}
		out.Rules = l
	}
}

func (in *IngressSpec) DeepCopy() *IngressSpec {
	if in == nil {
		return nil
	}
	out := new(IngressSpec)
	in.DeepCopyInto(out)
	return out
}

type IngressStatus struct {
	// LoadBalancer contains the current status of the load-balancer.
	LoadBalancer *corev1.LoadBalancerStatus `json:"loadBalancer,omitempty"`
}

func (in *IngressStatus) DeepCopyInto(out *IngressStatus) {
	*out = *in
	if in.LoadBalancer != nil {
		in, out := &in.LoadBalancer, &out.LoadBalancer
		*out = new(corev1.LoadBalancerStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *IngressStatus) DeepCopy() *IngressStatus {
	if in == nil {
		return nil
	}
	out := new(IngressStatus)
	in.DeepCopyInto(out)
	return out
}

type IngressClassSpec struct {
	// Controller refers to the name of the controller that should handle this
	// class. This allows for different "flavors" that are controlled by the
	// same controller. For example, you may have different Parameters for the
	// same implementing controller. This should be specified as a
	// domain-prefixed path no more than 250 characters in length, e.g.
	// "acme.io/ingress-controller". This field is immutable.
	Controller string `json:"controller,omitempty"`
	// Parameters is a link to a custom resource containing additional
	// configuration for the controller. This is optional if the controller does
	// not require extra parameters.
	Parameters *IngressClassParametersReference `json:"parameters,omitempty"`
}

func (in *IngressClassSpec) DeepCopyInto(out *IngressClassSpec) {
	*out = *in
	if in.Parameters != nil {
		in, out := &in.Parameters, &out.Parameters
		*out = new(IngressClassParametersReference)
		(*in).DeepCopyInto(*out)
	}
}

func (in *IngressClassSpec) DeepCopy() *IngressClassSpec {
	if in == nil {
		return nil
	}
	out := new(IngressClassSpec)
	in.DeepCopyInto(out)
	return out
}

type NetworkPolicySpec struct {
	// Selects the pods to which this NetworkPolicy object applies. The array of
	// ingress rules is applied to any pods selected by this field. Multiple network
	// policies can select the same set of pods. In this case, the ingress rules for
	// each are combined additively. This field is NOT optional and follows standard
	// label selector semantics. An empty podSelector matches all pods in this
	// namespace.
	PodSelector metav1.LabelSelector `json:"podSelector"`
	// List of ingress rules to be applied to the selected pods. Traffic is allowed to
	// a pod if there are no NetworkPolicies selecting the pod
	// (and cluster policy otherwise allows the traffic), OR if the traffic source is
	// the pod's local node, OR if the traffic matches at least one ingress rule
	// across all of the NetworkPolicy objects whose podSelector matches the pod. If
	// this field is empty then this NetworkPolicy does not allow any traffic (and serves
	// solely to ensure that the pods it selects are isolated by default)
	Ingress []NetworkPolicyIngressRule `json:"ingress"`
	// List of egress rules to be applied to the selected pods. Outgoing traffic is
	// allowed if there are no NetworkPolicies selecting the pod (and cluster policy
	// otherwise allows the traffic), OR if the traffic matches at least one egress rule
	// across all of the NetworkPolicy objects whose podSelector matches the pod. If
	// this field is empty then this NetworkPolicy limits all outgoing traffic (and serves
	// solely to ensure that the pods it selects are isolated by default).
	// This field is beta-level in 1.8
	Egress []NetworkPolicyEgressRule `json:"egress"`
	// List of rule types that the NetworkPolicy relates to.
	// Valid options are ["Ingress"], ["Egress"], or ["Ingress", "Egress"].
	// If this field is not specified, it will default based on the existence of Ingress or Egress rules;
	// policies that contain an Egress section are assumed to affect Egress, and all policies
	// (whether or not they contain an Ingress section) are assumed to affect Ingress.
	// If you want to write an egress-only policy, you must explicitly specify policyTypes [ "Egress" ].
	// Likewise, if you want to write a policy that specifies that no egress is allowed,
	// you must specify a policyTypes value that include "Egress" (since such a policy would not include
	// an Egress section and would otherwise default to just [ "Ingress" ]).
	// This field is beta-level in 1.8
	PolicyTypes []PolicyType `json:"policyTypes"`
}

func (in *NetworkPolicySpec) DeepCopyInto(out *NetworkPolicySpec) {
	*out = *in
	in.PodSelector.DeepCopyInto(&out.PodSelector)
	if in.Ingress != nil {
		l := make([]NetworkPolicyIngressRule, len(in.Ingress))
		for i := range in.Ingress {
			in.Ingress[i].DeepCopyInto(&l[i])
		}
		out.Ingress = l
	}
	if in.Egress != nil {
		l := make([]NetworkPolicyEgressRule, len(in.Egress))
		for i := range in.Egress {
			in.Egress[i].DeepCopyInto(&l[i])
		}
		out.Egress = l
	}
	if in.PolicyTypes != nil {
		t := make([]PolicyType, len(in.PolicyTypes))
		copy(t, in.PolicyTypes)
		out.PolicyTypes = t
	}
}

func (in *NetworkPolicySpec) DeepCopy() *NetworkPolicySpec {
	if in == nil {
		return nil
	}
	out := new(NetworkPolicySpec)
	in.DeepCopyInto(out)
	return out
}

type NetworkPolicyStatus struct {
	// Conditions holds an array of metav1.Condition that describe the state of the NetworkPolicy.
	// Current service state
	Conditions []metav1.Condition `json:"conditions"`
}

func (in *NetworkPolicyStatus) DeepCopyInto(out *NetworkPolicyStatus) {
	*out = *in
	if in.Conditions != nil {
		l := make([]metav1.Condition, len(in.Conditions))
		for i := range in.Conditions {
			in.Conditions[i].DeepCopyInto(&l[i])
		}
		out.Conditions = l
	}
}

func (in *NetworkPolicyStatus) DeepCopy() *NetworkPolicyStatus {
	if in == nil {
		return nil
	}
	out := new(NetworkPolicyStatus)
	in.DeepCopyInto(out)
	return out
}

type IngressBackend struct {
	// Service references a Service as a Backend.
	// This is a mutually exclusive setting with "Resource".
	Service *IngressServiceBackend `json:"service,omitempty"`
	// Resource is an ObjectRef to another Kubernetes resource in the namespace
	// of the Ingress object. If resource is specified, a service.Name and
	// service.Port must not be specified.
	// This is a mutually exclusive setting with "Service".
	Resource *corev1.TypedLocalObjectReference `json:"resource,omitempty"`
}

func (in *IngressBackend) DeepCopyInto(out *IngressBackend) {
	*out = *in
	if in.Service != nil {
		in, out := &in.Service, &out.Service
		*out = new(IngressServiceBackend)
		(*in).DeepCopyInto(*out)
	}
	if in.Resource != nil {
		in, out := &in.Resource, &out.Resource
		*out = new(corev1.TypedLocalObjectReference)
		(*in).DeepCopyInto(*out)
	}
}

func (in *IngressBackend) DeepCopy() *IngressBackend {
	if in == nil {
		return nil
	}
	out := new(IngressBackend)
	in.DeepCopyInto(out)
	return out
}

type IngressTLS struct {
	// Hosts are a list of hosts included in the TLS certificate. The values in
	// this list must match the name/s used in the tlsSecret. Defaults to the
	// wildcard host setting for the loadbalancer controller fulfilling this
	// Ingress, if left unspecified.
	Hosts []string `json:"hosts"`
	// SecretName is the name of the secret used to terminate TLS traffic on
	// port 443. Field is left optional to allow TLS routing based on SNI
	// hostname alone. If the SNI host in a listener conflicts with the "Host"
	// header field used by an IngressRule, the SNI host is used for termination
	// and value of the Host header is used for routing.
	SecretName string `json:"secretName,omitempty"`
}

func (in *IngressTLS) DeepCopyInto(out *IngressTLS) {
	*out = *in
	if in.Hosts != nil {
		t := make([]string, len(in.Hosts))
		copy(t, in.Hosts)
		out.Hosts = t
	}
}

func (in *IngressTLS) DeepCopy() *IngressTLS {
	if in == nil {
		return nil
	}
	out := new(IngressTLS)
	in.DeepCopyInto(out)
	return out
}

type IngressRule struct {
	// Host is the fully qualified domain name of a network host, as defined by RFC 3986.
	// Note the following deviations from the "host" part of the
	// URI as defined in RFC 3986:
	// 1. IPs are not allowed. Currently an IngressRuleValue can only apply to
	// the IP in the Spec of the parent Ingress.
	// 2. The `:` delimiter is not respected because ports are not allowed.
	// Currently the port of an Ingress is implicitly :80 for http and
	// :443 for https.
	// Both these may change in the future.
	// Incoming requests are matched against the host before the
	// IngressRuleValue. If the host is unspecified, the Ingress routes all
	// traffic based on the specified IngressRuleValue.
	// Host can be "precise" which is a domain name without the terminating dot of
	// a network host (e.g. "foo.bar.com") or "wildcard", which is a domain name
	// prefixed with a single wildcard label (e.g. "*.foo.com").
	// The wildcard character '*' must appear by itself as the first DNS label and
	// matches only a single label. You cannot have a wildcard label by itself (e.g. Host == "*").
	// Requests will be matched against the Host field in the following way:
	// 1. If Host is precise, the request matches this rule if the http host header is equal to Host.
	// 2. If Host is a wildcard, then the request matches this rule if the http host header
	// is to equal to the suffix (removing the first label) of the wildcard rule.
	Host string `json:"host,omitempty"`
	// IngressRuleValue represents a rule to route requests for this IngressRule.
	// If unspecified, the rule defaults to a http catch-all. Whether that sends
	// just traffic matching the host to the default backend or all traffic to the
	// default backend, is left to the controller fulfilling the Ingress. Http is
	// currently the only supported IngressRuleValue.
	IngressRuleValue IngressRuleValue `json:"ingressRuleValue"`
}

func (in *IngressRule) DeepCopyInto(out *IngressRule) {
	*out = *in
	in.IngressRuleValue.DeepCopyInto(&out.IngressRuleValue)
}

func (in *IngressRule) DeepCopy() *IngressRule {
	if in == nil {
		return nil
	}
	out := new(IngressRule)
	in.DeepCopyInto(out)
	return out
}

type IngressClassParametersReference struct {
	// APIGroup is the group for the resource being referenced. If APIGroup is
	// not specified, the specified Kind must be in the core API group. For any
	// other third-party types, APIGroup is required.
	APIGroup string `json:"apiGroup,omitempty"`
	// Kind is the type of resource being referenced.
	Kind string `json:"kind"`
	// Name is the name of resource being referenced.
	Name string `json:"name"`
	// Scope represents if this refers to a cluster or namespace scoped resource.
	// This may be set to "Cluster" (default) or "Namespace".
	Scope string `json:"scope,omitempty"`
	// Namespace is the namespace of the resource being referenced. This field is
	// required when scope is set to "Namespace" and must be unset when scope is set to
	// "Cluster".
	Namespace string `json:"namespace,omitempty"`
}

func (in *IngressClassParametersReference) DeepCopyInto(out *IngressClassParametersReference) {
	*out = *in
}

func (in *IngressClassParametersReference) DeepCopy() *IngressClassParametersReference {
	if in == nil {
		return nil
	}
	out := new(IngressClassParametersReference)
	in.DeepCopyInto(out)
	return out
}

type NetworkPolicyIngressRule struct {
	// List of ports which should be made accessible on the pods selected for this
	// rule. Each item in this list is combined using a logical OR. If this field is
	// empty or missing, this rule matches all ports (traffic not restricted by port).
	// If this field is present and contains at least one item, then this rule allows
	// traffic only if the traffic matches at least one port in the list.
	Ports []NetworkPolicyPort `json:"ports"`
	// List of sources which should be able to access the pods selected for this rule.
	// Items in this list are combined using a logical OR operation. If this field is
	// empty or missing, this rule matches all sources (traffic not restricted by
	// source). If this field is present and contains at least one item, this rule
	// allows traffic only if the traffic matches at least one item in the from list.
	From []NetworkPolicyPeer `json:"from"`
}

func (in *NetworkPolicyIngressRule) DeepCopyInto(out *NetworkPolicyIngressRule) {
	*out = *in
	if in.Ports != nil {
		l := make([]NetworkPolicyPort, len(in.Ports))
		for i := range in.Ports {
			in.Ports[i].DeepCopyInto(&l[i])
		}
		out.Ports = l
	}
	if in.From != nil {
		l := make([]NetworkPolicyPeer, len(in.From))
		for i := range in.From {
			in.From[i].DeepCopyInto(&l[i])
		}
		out.From = l
	}
}

func (in *NetworkPolicyIngressRule) DeepCopy() *NetworkPolicyIngressRule {
	if in == nil {
		return nil
	}
	out := new(NetworkPolicyIngressRule)
	in.DeepCopyInto(out)
	return out
}

type NetworkPolicyEgressRule struct {
	// List of destination ports for outgoing traffic.
	// Each item in this list is combined using a logical OR. If this field is
	// empty or missing, this rule matches all ports (traffic not restricted by port).
	// If this field is present and contains at least one item, then this rule allows
	// traffic only if the traffic matches at least one port in the list.
	Ports []NetworkPolicyPort `json:"ports"`
	// List of destinations for outgoing traffic of pods selected for this rule.
	// Items in this list are combined using a logical OR operation. If this field is
	// empty or missing, this rule matches all destinations (traffic not restricted by
	// destination). If this field is present and contains at least one item, this rule
	// allows traffic only if the traffic matches at least one item in the to list.
	To []NetworkPolicyPeer `json:"to"`
}

func (in *NetworkPolicyEgressRule) DeepCopyInto(out *NetworkPolicyEgressRule) {
	*out = *in
	if in.Ports != nil {
		l := make([]NetworkPolicyPort, len(in.Ports))
		for i := range in.Ports {
			in.Ports[i].DeepCopyInto(&l[i])
		}
		out.Ports = l
	}
	if in.To != nil {
		l := make([]NetworkPolicyPeer, len(in.To))
		for i := range in.To {
			in.To[i].DeepCopyInto(&l[i])
		}
		out.To = l
	}
}

func (in *NetworkPolicyEgressRule) DeepCopy() *NetworkPolicyEgressRule {
	if in == nil {
		return nil
	}
	out := new(NetworkPolicyEgressRule)
	in.DeepCopyInto(out)
	return out
}

type IngressServiceBackend struct {
	// Name is the referenced service. The service must exist in
	// the same namespace as the Ingress object.
	Name string `json:"name"`
	// Port of the referenced service. A port name or port number
	// is required for a IngressServiceBackend.
	Port *ServiceBackendPort `json:"port,omitempty"`
}

func (in *IngressServiceBackend) DeepCopyInto(out *IngressServiceBackend) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(ServiceBackendPort)
		(*in).DeepCopyInto(*out)
	}
}

func (in *IngressServiceBackend) DeepCopy() *IngressServiceBackend {
	if in == nil {
		return nil
	}
	out := new(IngressServiceBackend)
	in.DeepCopyInto(out)
	return out
}

type IngressRuleValue struct {
	HTTP *HTTPIngressRuleValue `json:"http,omitempty"`
}

func (in *IngressRuleValue) DeepCopyInto(out *IngressRuleValue) {
	*out = *in
	if in.HTTP != nil {
		in, out := &in.HTTP, &out.HTTP
		*out = new(HTTPIngressRuleValue)
		(*in).DeepCopyInto(*out)
	}
}

func (in *IngressRuleValue) DeepCopy() *IngressRuleValue {
	if in == nil {
		return nil
	}
	out := new(IngressRuleValue)
	in.DeepCopyInto(out)
	return out
}

type NetworkPolicyPort struct {
	// The protocol (TCP, UDP, or SCTP) which traffic must match. If not specified, this
	// field defaults to TCP.
	Protocol corev1.Protocol `json:"protocol,omitempty"`
	// The port on the given protocol. This can either be a numerical or named
	// port on a pod. If this field is not provided, this matches all port names and
	// numbers.
	// If present, only traffic on the specified protocol AND port will be matched.
	Port *utilintstr.IntOrString `json:"port,omitempty"`
	// If set, indicates that the range of ports from port to endPort, inclusive,
	// should be allowed by the policy. This field cannot be defined if the port field
	// is not defined or if the port field is defined as a named (string) port.
	// The endPort must be equal or greater than port.
	// This feature is in Beta state and is enabled by default.
	// It can be disabled using the Feature Gate "NetworkPolicyEndPort".
	EndPort int `json:"endPort,omitempty"`
}

func (in *NetworkPolicyPort) DeepCopyInto(out *NetworkPolicyPort) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(utilintstr.IntOrString)
		*out = *in
	}
}

func (in *NetworkPolicyPort) DeepCopy() *NetworkPolicyPort {
	if in == nil {
		return nil
	}
	out := new(NetworkPolicyPort)
	in.DeepCopyInto(out)
	return out
}

type NetworkPolicyPeer struct {
	// This is a label selector which selects Pods. This field follows standard label
	// selector semantics; if present but empty, it selects all pods.
	// If NamespaceSelector is also set, then the NetworkPolicyPeer as a whole selects
	// the Pods matching PodSelector in the Namespaces selected by NamespaceSelector.
	// Otherwise it selects the Pods matching PodSelector in the policy's own Namespace.
	PodSelector *metav1.LabelSelector `json:"podSelector,omitempty"`
	// Selects Namespaces using cluster-scoped labels. This field follows standard label
	// selector semantics; if present but empty, it selects all namespaces.
	// If PodSelector is also set, then the NetworkPolicyPeer as a whole selects
	// the Pods matching PodSelector in the Namespaces selected by NamespaceSelector.
	// Otherwise it selects all Pods in the Namespaces selected by NamespaceSelector.
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty"`
	// IPBlock defines policy on a particular IPBlock. If this field is set then
	// neither of the other fields can be.
	IPBlock *IPBlock `json:"ipBlock,omitempty"`
}

func (in *NetworkPolicyPeer) DeepCopyInto(out *NetworkPolicyPeer) {
	*out = *in
	if in.PodSelector != nil {
		in, out := &in.PodSelector, &out.PodSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.NamespaceSelector != nil {
		in, out := &in.NamespaceSelector, &out.NamespaceSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.IPBlock != nil {
		in, out := &in.IPBlock, &out.IPBlock
		*out = new(IPBlock)
		(*in).DeepCopyInto(*out)
	}
}

func (in *NetworkPolicyPeer) DeepCopy() *NetworkPolicyPeer {
	if in == nil {
		return nil
	}
	out := new(NetworkPolicyPeer)
	in.DeepCopyInto(out)
	return out
}

type ServiceBackendPort struct {
	// Name is the name of the port on the Service.
	// This is a mutually exclusive setting with "Number".
	Name string `json:"name,omitempty"`
	// Number is the numerical port number (e.g. 80) on the Service.
	// This is a mutually exclusive setting with "Name".
	Number int `json:"number,omitempty"`
}

func (in *ServiceBackendPort) DeepCopyInto(out *ServiceBackendPort) {
	*out = *in
}

func (in *ServiceBackendPort) DeepCopy() *ServiceBackendPort {
	if in == nil {
		return nil
	}
	out := new(ServiceBackendPort)
	in.DeepCopyInto(out)
	return out
}

type HTTPIngressRuleValue struct {
	// A collection of paths that map requests to backends.
	Paths []HTTPIngressPath `json:"paths"`
}

func (in *HTTPIngressRuleValue) DeepCopyInto(out *HTTPIngressRuleValue) {
	*out = *in
	if in.Paths != nil {
		l := make([]HTTPIngressPath, len(in.Paths))
		for i := range in.Paths {
			in.Paths[i].DeepCopyInto(&l[i])
		}
		out.Paths = l
	}
}

func (in *HTTPIngressRuleValue) DeepCopy() *HTTPIngressRuleValue {
	if in == nil {
		return nil
	}
	out := new(HTTPIngressRuleValue)
	in.DeepCopyInto(out)
	return out
}

type IPBlock struct {
	// CIDR is a string representing the IP Block
	// Valid examples are "192.168.1.1/24" or "2001:db9::/64"
	CIDR string `json:"cidr"`
	// Except is a slice of CIDRs that should not be included within an IP Block
	// Valid examples are "192.168.1.1/24" or "2001:db9::/64"
	// Except values will be rejected if they are outside the CIDR range
	Except []string `json:"except"`
}

func (in *IPBlock) DeepCopyInto(out *IPBlock) {
	*out = *in
	if in.Except != nil {
		t := make([]string, len(in.Except))
		copy(t, in.Except)
		out.Except = t
	}
}

func (in *IPBlock) DeepCopy() *IPBlock {
	if in == nil {
		return nil
	}
	out := new(IPBlock)
	in.DeepCopyInto(out)
	return out
}

type HTTPIngressPath struct {
	// Path is matched against the path of an incoming request. Currently it can
	// contain characters disallowed from the conventional "path" part of a URL
	// as defined by RFC 3986. Paths must begin with a '/' and must be present
	// when using PathType with value "Exact" or "Prefix".
	Path string `json:"path,omitempty"`
	// PathType determines the interpretation of the Path matching. PathType can
	// be one of the following values:
	// * Exact: Matches the URL path exactly.
	// * Prefix: Matches based on a URL path prefix split by '/'. Matching is
	// done on a path element by element basis. A path element refers is the
	// list of labels in the path split by the '/' separator. A request is a
	// match for path p if every p is an element-wise prefix of p of the
	// request path. Note that if the last element of the path is a substring
	// of the last element in request path, it is not a match (e.g. /foo/bar
	// matches /foo/bar/baz, but does not match /foo/barbaz).
	// * ImplementationSpecific: Interpretation of the Path matching is up to
	// the IngressClass. Implementations can treat this as a separate PathType
	// or treat it identically to Prefix or Exact path types.
	// Implementations are required to support all path types.
	PathType string `json:"pathType"`
	// Backend defines the referenced service endpoint to which the traffic
	// will be forwarded to.
	Backend IngressBackend `json:"backend"`
}

func (in *HTTPIngressPath) DeepCopyInto(out *HTTPIngressPath) {
	*out = *in
	in.Backend.DeepCopyInto(&out.Backend)
}

func (in *HTTPIngressPath) DeepCopy() *HTTPIngressPath {
	if in == nil {
		return nil
	}
	out := new(HTTPIngressPath)
	in.DeepCopyInto(out)
	return out
}
