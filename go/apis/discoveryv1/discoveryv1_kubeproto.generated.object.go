package discoveryv1

import (
	"go.f110.dev/kubeproto/go/apis/corev1"
	"go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "discovery.k8s.io"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&EndpointSlice{},
		&EndpointSliceList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type EndpointSlice struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// addressType specifies the type of address carried by this EndpointSlice.
	// All addresses in this slice must be the same type. This field is
	// immutable after creation. The following address types are currently
	// supported:
	// * IPv4: Represents an IPv4 Address.
	// * IPv6: Represents an IPv6 Address.
	// * FQDN: Represents a Fully Qualified Domain Name.
	AddressType string `json:"addressType"`
	// endpoints is a list of unique endpoints in this slice. Each slice may
	// include a maximum of 1000 endpoints.
	Endpoints []Endpoint `json:"endpoints"`
	// ports specifies the list of network ports exposed by each endpoint in
	// this slice. Each port must have a unique name. When ports is empty, it
	// indicates that there are no defined ports. When a port is defined with a
	// nil port value, it indicates "all ports". Each slice may include a
	// maximum of 100 ports.
	Ports []EndpointPort `json:"ports"`
}

func (in *EndpointSlice) DeepCopyInto(out *EndpointSlice) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Endpoints != nil {
		l := make([]Endpoint, len(in.Endpoints))
		for i := range in.Endpoints {
			in.Endpoints[i].DeepCopyInto(&l[i])
		}
		out.Endpoints = l
	}
	if in.Ports != nil {
		l := make([]EndpointPort, len(in.Ports))
		for i := range in.Ports {
			in.Ports[i].DeepCopyInto(&l[i])
		}
		out.Ports = l
	}
}

func (in *EndpointSlice) DeepCopy() *EndpointSlice {
	if in == nil {
		return nil
	}
	out := new(EndpointSlice)
	in.DeepCopyInto(out)
	return out
}

func (in *EndpointSlice) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type EndpointSliceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []EndpointSlice `json:"items"`
}

func (in *EndpointSliceList) DeepCopyInto(out *EndpointSliceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]EndpointSlice, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *EndpointSliceList) DeepCopy() *EndpointSliceList {
	if in == nil {
		return nil
	}
	out := new(EndpointSliceList)
	in.DeepCopyInto(out)
	return out
}

func (in *EndpointSliceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type Endpoint struct {
	// addresses of this endpoint. The contents of this field are interpreted
	// according to the corresponding EndpointSlice addressType field. Consumers
	// must handle different types of addresses in the context of their own
	// capabilities. This must contain at least one address but no more than
	// 100. These are all assumed to be fungible and clients may choose to only
	// use the first element. Refer to: https://issue.k8s.io/106267
	Addresses []string `json:"addresses"`
	// conditions contains information about the current status of the endpoint.
	Conditions *EndpointConditions `json:"conditions,omitempty"`
	// hostname of this endpoint. This field may be used by consumers of
	// endpoints to distinguish endpoints from each other (e.g. in DNS names).
	// Multiple endpoints which use the same hostname should be considered
	// fungible (e.g. multiple A values in DNS). Must be lowercase and pass DNS
	// Label (RFC 1123) validation.
	Hostname string `json:"hostname,omitempty"`
	// targetRef is a reference to a Kubernetes object that represents this
	// endpoint.
	TargetRef *corev1.ObjectReference `json:"targetRef,omitempty"`
	// deprecatedTopology contains topology information part of the v1beta1
	// API. This field is deprecated, and will be removed when the v1beta1
	// API is removed (no sooner than kubernetes v1.24).  While this field can
	// hold values, it is not writable through the v1 API, and any attempts to
	// write to it will be silently ignored. Topology information can be found
	// in the zone and nodeName fields instead.
	DeprecatedTopology map[string]string `json:"deprecatedTopology,omitempty"`
	// nodeName represents the name of the Node hosting this endpoint. This can
	// be used to determine endpoints local to a Node.
	NodeName string `json:"nodeName,omitempty"`
	// zone is the name of the Zone this endpoint exists in.
	Zone string `json:"zone,omitempty"`
	// hints contains information associated with how an endpoint should be
	// consumed.
	Hints *EndpointHints `json:"hints,omitempty"`
}

func (in *Endpoint) DeepCopyInto(out *Endpoint) {
	*out = *in
	if in.Addresses != nil {
		t := make([]string, len(in.Addresses))
		copy(t, in.Addresses)
		out.Addresses = t
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = new(EndpointConditions)
		(*in).DeepCopyInto(*out)
	}
	if in.TargetRef != nil {
		in, out := &in.TargetRef, &out.TargetRef
		*out = new(corev1.ObjectReference)
		(*in).DeepCopyInto(*out)
	}
	if in.DeprecatedTopology != nil {
		in, out := &in.DeprecatedTopology, &out.DeprecatedTopology
		*out = make(map[string]string, len(*in))
		for k, v := range *in {
			(*out)[k] = v
		}
	}
	if in.Hints != nil {
		in, out := &in.Hints, &out.Hints
		*out = new(EndpointHints)
		(*in).DeepCopyInto(*out)
	}
}

func (in *Endpoint) DeepCopy() *Endpoint {
	if in == nil {
		return nil
	}
	out := new(Endpoint)
	in.DeepCopyInto(out)
	return out
}

type EndpointPort struct {
	// name represents the name of this port. All ports in an EndpointSlice must have a unique name.
	// If the EndpointSlice is dervied from a Kubernetes service, this corresponds to the Service.ports[].name.
	// Name must either be an empty string or pass DNS_LABEL validation:
	// * must be no more than 63 characters long.
	// * must consist of lower case alphanumeric characters or '-'.
	// * must start and end with an alphanumeric character.
	// Default is empty string.
	Name string `json:"name,omitempty"`
	// protocol represents the IP protocol for this port.
	// Must be UDP, TCP, or SCTP.
	// Default is TCP.
	Protocol corev1.Protocol `json:"protocol,omitempty"`
	// port represents the port number of the endpoint.
	// If this is not specified, ports are not restricted and must be
	// interpreted in the context of the specific consumer.
	Port int `json:"port,omitempty"`
	// The application protocol for this port.
	// This is used as a hint for implementations to offer richer behavior for protocols that they understand.
	// This field follows standard Kubernetes label syntax.
	// Valid values are either:
	// * Un-prefixed protocol names - reserved for IANA standard service names (as per
	// RFC-6335 and https://www.iana.org/assignments/service-names).
	// * Kubernetes-defined prefixed names:
	// * 'kubernetes.io/h2c' - HTTP/2 over cleartext as described in https://www.rfc-editor.org/rfc/rfc7540
	// * Other protocols should use implementation-defined prefixed names such as
	// mycompany.com/my-custom-protocol.
	AppProtocol string `json:"appProtocol,omitempty"`
}

func (in *EndpointPort) DeepCopyInto(out *EndpointPort) {
	*out = *in
}

func (in *EndpointPort) DeepCopy() *EndpointPort {
	if in == nil {
		return nil
	}
	out := new(EndpointPort)
	in.DeepCopyInto(out)
	return out
}

type EndpointConditions struct {
	// ready indicates that this endpoint is prepared to receive traffic,
	// according to whatever system is managing the endpoint. A nil value
	// indicates an unknown state. In most cases consumers should interpret this
	// unknown state as ready. For compatibility reasons, ready should never be
	// "true" for terminating endpoints, except when the normal readiness
	// behavior is being explicitly overridden, for example when the associated
	// Service has set the publishNotReadyAddresses flag.
	Ready bool `json:"ready,omitempty"`
	// serving is identical to ready except that it is set regardless of the
	// terminating state of endpoints. This condition should be set to true for
	// a ready endpoint that is terminating. If nil, consumers should defer to
	// the ready condition.
	Serving bool `json:"serving,omitempty"`
	// terminating indicates that this endpoint is terminating. A nil value
	// indicates an unknown state. Consumers should interpret this unknown state
	// to mean that the endpoint is not terminating.
	Terminating bool `json:"terminating,omitempty"`
}

func (in *EndpointConditions) DeepCopyInto(out *EndpointConditions) {
	*out = *in
}

func (in *EndpointConditions) DeepCopy() *EndpointConditions {
	if in == nil {
		return nil
	}
	out := new(EndpointConditions)
	in.DeepCopyInto(out)
	return out
}

type EndpointHints struct {
	// forZones indicates the zone(s) this endpoint should be consumed by to
	// enable topology aware routing.
	ForZones []ForZone `json:"forZones"`
}

func (in *EndpointHints) DeepCopyInto(out *EndpointHints) {
	*out = *in
	if in.ForZones != nil {
		l := make([]ForZone, len(in.ForZones))
		for i := range in.ForZones {
			in.ForZones[i].DeepCopyInto(&l[i])
		}
		out.ForZones = l
	}
}

func (in *EndpointHints) DeepCopy() *EndpointHints {
	if in == nil {
		return nil
	}
	out := new(EndpointHints)
	in.DeepCopyInto(out)
	return out
}

type ForZone struct {
	// name represents the name of the zone.
	Name string `json:"name"`
}

func (in *ForZone) DeepCopyInto(out *ForZone) {
	*out = *in
}

func (in *ForZone) DeepCopy() *ForZone {
	if in == nil {
		return nil
	}
	out := new(ForZone)
	in.DeepCopyInto(out)
	return out
}
