package gatewayv1

import (
	"go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"maps"
)

const GroupName = "gateway.networking.k8s.io"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&BackendTLSPolicy{},
		&BackendTLSPolicyList{},
		&GRPCRoute{},
		&GRPCRouteList{},
		&Gateway{},
		&GatewayClass{},
		&GatewayClassList{},
		&GatewayList{},
		&HTTPRoute{},
		&HTTPRouteList{},
		&ListenerSet{},
		&ListenerSetList{},
		&ReferenceGrant{},
		&ReferenceGrantList{},
		&TLSRoute{},
		&TLSRouteList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type AddressType string

const (
	AddressTypeIPAddress    AddressType = "IPAddress"
	AddressTypeHostname     AddressType = "Hostname"
	AddressTypeNamedAddress AddressType = "NamedAddress"
)

type CookieLifetimeType string

const (
	CookieLifetimeTypeSession   CookieLifetimeType = "Session"
	CookieLifetimeTypePermanent CookieLifetimeType = "Permanent"
)

type FromNamespaces string

const (
	FromNamespacesAll      FromNamespaces = "All"
	FromNamespacesSelector FromNamespaces = "Selector"
	FromNamespacesSame     FromNamespaces = "Same"
	FromNamespacesNone     FromNamespaces = "None"
)

type FrontendValidationModeType string

const (
	FrontendValidationModeTypeAllowValidOnly        FrontendValidationModeType = "AllowValidOnly"
	FrontendValidationModeTypeAllowInsecureFallback FrontendValidationModeType = "AllowInsecureFallback"
)

type GRPCHeaderMatchType string

const (
	GRPCHeaderMatchTypeExact             GRPCHeaderMatchType = "Exact"
	GRPCHeaderMatchTypeRegularExpression GRPCHeaderMatchType = "RegularExpression"
)

type GRPCMethodMatchType string

const (
	GRPCMethodMatchTypeExact             GRPCMethodMatchType = "Exact"
	GRPCMethodMatchTypeRegularExpression GRPCMethodMatchType = "RegularExpression"
)

type GRPCRouteFilterType string

const (
	GRPCRouteFilterTypeRequestHeaderModifier  GRPCRouteFilterType = "RequestHeaderModifier"
	GRPCRouteFilterTypeResponseHeaderModifier GRPCRouteFilterType = "ResponseHeaderModifier"
	GRPCRouteFilterTypeRequestMirror          GRPCRouteFilterType = "RequestMirror"
	GRPCRouteFilterTypeExtensionRef           GRPCRouteFilterType = "ExtensionRef"
)

type GatewayClassConditionReason string

const (
	GatewayClassConditionReasonAccepted           GatewayClassConditionReason = "Accepted"
	GatewayClassConditionReasonInvalidParameters  GatewayClassConditionReason = "InvalidParameters"
	GatewayClassConditionReasonPending            GatewayClassConditionReason = "Pending"
	GatewayClassConditionReasonUnsupported        GatewayClassConditionReason = "Unsupported"
	GatewayClassConditionReasonWaiting            GatewayClassConditionReason = "Waiting"
	GatewayClassConditionReasonSupportedVersion   GatewayClassConditionReason = "SupportedVersion"
	GatewayClassConditionReasonUnsupportedVersion GatewayClassConditionReason = "UnsupportedVersion"
)

type GatewayClassConditionType string

const (
	GatewayClassConditionTypeAccepted         GatewayClassConditionType = "Accepted"
	GatewayClassConditionTypeSupportedVersion GatewayClassConditionType = "SupportedVersion"
)

type GatewayConditionReason string

const (
	GatewayConditionReasonProgrammed                  GatewayConditionReason = "Programmed"
	GatewayConditionReasonInvalid                     GatewayConditionReason = "Invalid"
	GatewayConditionReasonNoResources                 GatewayConditionReason = "NoResources"
	GatewayConditionReasonAddressNotAssigned          GatewayConditionReason = "AddressNotAssigned"
	GatewayConditionReasonAddressNotUsable            GatewayConditionReason = "AddressNotUsable"
	GatewayConditionReasonConfigurationChanged        GatewayConditionReason = "ConfigurationChanged"
	GatewayConditionReasonAccepted                    GatewayConditionReason = "Accepted"
	GatewayConditionReasonListenersNotValid           GatewayConditionReason = "ListenersNotValid"
	GatewayConditionReasonPending                     GatewayConditionReason = "Pending"
	GatewayConditionReasonUnsupportedAddress          GatewayConditionReason = "UnsupportedAddress"
	GatewayConditionReasonInvalidParameters           GatewayConditionReason = "InvalidParameters"
	GatewayConditionReasonScheduled                   GatewayConditionReason = "Scheduled"
	GatewayConditionReasonNotReconciled               GatewayConditionReason = "NotReconciled"
	GatewayConditionReasonResolvedRefs                GatewayConditionReason = "ResolvedRefs"
	GatewayConditionReasonInvalidClientCertificateRef GatewayConditionReason = "InvalidClientCertificateRef"
	GatewayConditionReasonRefNotPermitted             GatewayConditionReason = "RefNotPermitted"
	GatewayConditionReasonListenersNotResolved        GatewayConditionReason = "ListenersNotResolved"
	GatewayConditionReasonReady                       GatewayConditionReason = "Ready"
	GatewayConditionReasonListenersNotReady           GatewayConditionReason = "ListenersNotReady"
)

type GatewayConditionType string

const (
	GatewayConditionTypeProgrammed                     GatewayConditionType = "Programmed"
	GatewayConditionTypeInsecureFrontendValidationMode GatewayConditionType = "InsecureFrontendValidationMode"
	GatewayConditionTypeAccepted                       GatewayConditionType = "Accepted"
	GatewayConditionTypeScheduled                      GatewayConditionType = "Scheduled"
	GatewayConditionTypeResolvedRefs                   GatewayConditionType = "ResolvedRefs"
	GatewayConditionTypeReady                          GatewayConditionType = "Ready"
)

type GatewayDefaultScope string

const (
	GatewayDefaultScopeAll  GatewayDefaultScope = "All"
	GatewayDefaultScopeNone GatewayDefaultScope = "None"
)

type HTTPMethod string

const (
	HTTPMethodGET     HTTPMethod = "GET"
	HTTPMethodHEAD    HTTPMethod = "HEAD"
	HTTPMethodPOST    HTTPMethod = "POST"
	HTTPMethodPUT     HTTPMethod = "PUT"
	HTTPMethodDELETE  HTTPMethod = "DELETE"
	HTTPMethodCONNECT HTTPMethod = "CONNECT"
	HTTPMethodOPTIONS HTTPMethod = "OPTIONS"
	HTTPMethodTRACE   HTTPMethod = "TRACE"
	HTTPMethodPATCH   HTTPMethod = "PATCH"
)

type HTTPPathModifierType string

const (
	HTTPPathModifierTypeReplaceFullPath    HTTPPathModifierType = "ReplaceFullPath"
	HTTPPathModifierTypeReplacePrefixMatch HTTPPathModifierType = "ReplacePrefixMatch"
)

type HTTPRouteExternalAuthProtocol string

const (
	HTTPRouteExternalAuthProtocolGRPC HTTPRouteExternalAuthProtocol = "GRPC"
	HTTPRouteExternalAuthProtocolHTTP HTTPRouteExternalAuthProtocol = "HTTP"
)

type HTTPRouteFilterType string

const (
	HTTPRouteFilterTypeRequestHeaderModifier  HTTPRouteFilterType = "RequestHeaderModifier"
	HTTPRouteFilterTypeResponseHeaderModifier HTTPRouteFilterType = "ResponseHeaderModifier"
	HTTPRouteFilterTypeRequestRedirect        HTTPRouteFilterType = "RequestRedirect"
	HTTPRouteFilterTypeURLRewrite             HTTPRouteFilterType = "URLRewrite"
	HTTPRouteFilterTypeRequestMirror          HTTPRouteFilterType = "RequestMirror"
	HTTPRouteFilterTypeCors                   HTTPRouteFilterType = "CORS"
	HTTPRouteFilterTypeExternalAuth           HTTPRouteFilterType = "ExternalAuth"
	HTTPRouteFilterTypeExtensionRef           HTTPRouteFilterType = "ExtensionRef"
)

type HeaderMatchType string

const (
	HeaderMatchTypeExact             HeaderMatchType = "Exact"
	HeaderMatchTypeRegularExpression HeaderMatchType = "RegularExpression"
)

type ListenerConditionReason string

const (
	ListenerConditionReasonHostnameConflict         ListenerConditionReason = "HostnameConflict"
	ListenerConditionReasonProtocolConflict         ListenerConditionReason = "ProtocolConflict"
	ListenerConditionReasonNoConflicts              ListenerConditionReason = "NoConflicts"
	ListenerConditionReasonAccepted                 ListenerConditionReason = "Accepted"
	ListenerConditionReasonAttached                 ListenerConditionReason = "Attached"
	ListenerConditionReasonPortUnavailable          ListenerConditionReason = "PortUnavailable"
	ListenerConditionReasonUnsupportedProtocol      ListenerConditionReason = "UnsupportedProtocol"
	ListenerConditionReasonNoValidCACertificate     ListenerConditionReason = "NoValidCACertificate"
	ListenerConditionReasonUnsupportedValue         ListenerConditionReason = "UnsupportedValue"
	ListenerConditionReasonResolvedRefs             ListenerConditionReason = "ResolvedRefs"
	ListenerConditionReasonInvalidCertificateRef    ListenerConditionReason = "InvalidCertificateRef"
	ListenerConditionReasonInvalidRouteKinds        ListenerConditionReason = "InvalidRouteKinds"
	ListenerConditionReasonRefNotPermitted          ListenerConditionReason = "RefNotPermitted"
	ListenerConditionReasonInvalidCACertificateRef  ListenerConditionReason = "InvalidCACertificateRef"
	ListenerConditionReasonInvalidCACertificateKind ListenerConditionReason = "InvalidCACertificateKind"
	ListenerConditionReasonProgrammed               ListenerConditionReason = "Programmed"
	ListenerConditionReasonInvalid                  ListenerConditionReason = "Invalid"
	ListenerConditionReasonPending                  ListenerConditionReason = "Pending"
	ListenerConditionReasonOverlappingHostnames     ListenerConditionReason = "OverlappingHostnames"
	ListenerConditionReasonOverlappingCertificates  ListenerConditionReason = "OverlappingCertificates"
	ListenerConditionReasonReady                    ListenerConditionReason = "Ready"
)

type ListenerConditionType string

const (
	ListenerConditionTypeConflicted           ListenerConditionType = "Conflicted"
	ListenerConditionTypeAccepted             ListenerConditionType = "Accepted"
	ListenerConditionTypeDetached             ListenerConditionType = "Detached"
	ListenerConditionTypeResolvedRefs         ListenerConditionType = "ResolvedRefs"
	ListenerConditionTypeProgrammed           ListenerConditionType = "Programmed"
	ListenerConditionTypeOverlappingTLSConfig ListenerConditionType = "OverlappingTLSConfig"
	ListenerConditionTypeReady                ListenerConditionType = "Ready"
)

type ListenerEntryConditionReason string

const (
	ListenerEntryConditionReasonHostnameConflict      ListenerEntryConditionReason = "HostnameConflict"
	ListenerEntryConditionReasonProtocolConflict      ListenerEntryConditionReason = "ProtocolConflict"
	ListenerEntryConditionReasonListenerConflict      ListenerEntryConditionReason = "ListenerConflict"
	ListenerEntryConditionReasonAccepted              ListenerEntryConditionReason = "Accepted"
	ListenerEntryConditionReasonUnsupportedProtocol   ListenerEntryConditionReason = "UnsupportedProtocol"
	ListenerEntryConditionReasonTooManyListeners      ListenerEntryConditionReason = "TooManyListeners"
	ListenerEntryConditionReasonResolvedRefs          ListenerEntryConditionReason = "ResolvedRefs"
	ListenerEntryConditionReasonInvalidCertificateRef ListenerEntryConditionReason = "InvalidCertificateRef"
	ListenerEntryConditionReasonInvalidRouteKinds     ListenerEntryConditionReason = "InvalidRouteKinds"
	ListenerEntryConditionReasonRefNotPermitted       ListenerEntryConditionReason = "RefNotPermitted"
	ListenerEntryConditionReasonProgrammed            ListenerEntryConditionReason = "Programmed"
	ListenerEntryConditionReasonPortUnavailable       ListenerEntryConditionReason = "PortUnavailable"
	ListenerEntryConditionReasonPending               ListenerEntryConditionReason = "Pending"
	ListenerEntryConditionReasonInvalid               ListenerEntryConditionReason = "Invalid"
	ListenerEntryConditionReasonReady                 ListenerEntryConditionReason = "Ready"
)

type ListenerEntryConditionType string

const (
	ListenerEntryConditionTypeConflicted   ListenerEntryConditionType = "Conflicted"
	ListenerEntryConditionTypeAccepted     ListenerEntryConditionType = "Accepted"
	ListenerEntryConditionTypeResolvedRefs ListenerEntryConditionType = "ResolvedRefs"
	ListenerEntryConditionTypeProgrammed   ListenerEntryConditionType = "Programmed"
	ListenerEntryConditionTypeReady        ListenerEntryConditionType = "Ready"
)

type ListenerSetConditionReason string

const (
	ListenerSetConditionReasonProgrammed        ListenerSetConditionReason = "Programmed"
	ListenerSetConditionReasonAccepted          ListenerSetConditionReason = "Accepted"
	ListenerSetConditionReasonNotAllowed        ListenerSetConditionReason = "NotAllowed"
	ListenerSetConditionReasonParentNotAccepted ListenerSetConditionReason = "ParentNotAccepted"
	ListenerSetConditionReasonListenersNotValid ListenerSetConditionReason = "ListenersNotValid"
	ListenerSetConditionReasonInvalid           ListenerSetConditionReason = "Invalid"
	ListenerSetConditionReasonPending           ListenerSetConditionReason = "Pending"
)

type ListenerSetConditionType string

const (
	ListenerSetConditionTypeProgrammed ListenerSetConditionType = "Programmed"
	ListenerSetConditionTypeAccepted   ListenerSetConditionType = "Accepted"
)

type PathMatchType string

const (
	PathMatchTypeExact             PathMatchType = "Exact"
	PathMatchTypePathPrefix        PathMatchType = "PathPrefix"
	PathMatchTypeRegularExpression PathMatchType = "RegularExpression"
)

type PolicyConditionReason string

const (
	PolicyConditionReasonNoValidCACertificate    PolicyConditionReason = "NoValidCACertificate"
	PolicyConditionReasonResolvedRefs            PolicyConditionReason = "ResolvedRefs"
	PolicyConditionReasonInvalidCACertificateRef PolicyConditionReason = "InvalidCACertificateRef"
	PolicyConditionReasonInvalidKind             PolicyConditionReason = "InvalidKind"
	PolicyConditionReasonAccepted                PolicyConditionReason = "Accepted"
	PolicyConditionReasonConflicted              PolicyConditionReason = "Conflicted"
	PolicyConditionReasonInvalid                 PolicyConditionReason = "Invalid"
	PolicyConditionReasonTargetNotFound          PolicyConditionReason = "TargetNotFound"
)

type PolicyConditionType string

const (
	PolicyConditionTypeResolvedRefs PolicyConditionType = "ResolvedRefs"
	PolicyConditionTypeAccepted     PolicyConditionType = "Accepted"
)

type ProtocolType string

const (
	ProtocolTypeHTTP  ProtocolType = "HTTP"
	ProtocolTypeHTTPS ProtocolType = "HTTPS"
	ProtocolTypeTLS   ProtocolType = "TLS"
	ProtocolTypeTCP   ProtocolType = "TCP"
	ProtocolTypeUDP   ProtocolType = "UDP"
)

type QueryParamMatchType string

const (
	QueryParamMatchTypeExact             QueryParamMatchType = "Exact"
	QueryParamMatchTypeRegularExpression QueryParamMatchType = "RegularExpression"
)

type RouteConditionReason string

const (
	RouteConditionReasonAccepted                   RouteConditionReason = "Accepted"
	RouteConditionReasonNotAllowedByListeners      RouteConditionReason = "NotAllowedByListeners"
	RouteConditionReasonNoMatchingListenerHostname RouteConditionReason = "NoMatchingListenerHostname"
	RouteConditionReasonNoMatchingParent           RouteConditionReason = "NoMatchingParent"
	RouteConditionReasonUnsupportedValue           RouteConditionReason = "UnsupportedValue"
	RouteConditionReasonPending                    RouteConditionReason = "Pending"
	RouteConditionReasonIncompatibleFilters        RouteConditionReason = "IncompatibleFilters"
	RouteConditionReasonResolvedRefs               RouteConditionReason = "ResolvedRefs"
	RouteConditionReasonRefNotPermitted            RouteConditionReason = "RefNotPermitted"
	RouteConditionReasonInvalidKind                RouteConditionReason = "InvalidKind"
	RouteConditionReasonBackendNotFound            RouteConditionReason = "BackendNotFound"
	RouteConditionReasonUnsupportedProtocol        RouteConditionReason = "UnsupportedProtocol"
)

type RouteConditionType string

const (
	RouteConditionTypeAccepted         RouteConditionType = "Accepted"
	RouteConditionTypeResolvedRefs     RouteConditionType = "ResolvedRefs"
	RouteConditionTypePartiallyInvalid RouteConditionType = "PartiallyInvalid"
)

type SessionPersistenceType string

const (
	SessionPersistenceTypeCookie SessionPersistenceType = "Cookie"
	SessionPersistenceTypeHeader SessionPersistenceType = "Header"
)

type SubjectAltNameType string

const (
	SubjectAltNameTypeHostname SubjectAltNameType = "Hostname"
	SubjectAltNameTypeURI      SubjectAltNameType = "URI"
)

type TLSModeType string

const (
	TLSModeTypeTerminate   TLSModeType = "Terminate"
	TLSModeTypePassthrough TLSModeType = "Passthrough"
)

type WellKnownCACertificatesType string

const (
	WellKnownCACertificatesTypeSystem WellKnownCACertificatesType = "System"
)

type BackendTLSPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec defines the desired state of BackendTLSPolicy.
	Spec BackendTLSPolicySpec `json:"spec"`
	// Status defines the current state of BackendTLSPolicy.
	Status *PolicyStatus `json:"status,omitempty"`
}

func (in *BackendTLSPolicy) DeepCopyInto(out *BackendTLSPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(PolicyStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *BackendTLSPolicy) DeepCopy() *BackendTLSPolicy {
	if in == nil {
		return nil
	}
	out := new(BackendTLSPolicy)
	in.DeepCopyInto(out)
	return out
}

func (in *BackendTLSPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *BackendTLSPolicy) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.BackendTLSPolicy"
}

type BackendTLSPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []BackendTLSPolicy `json:"items,omitempty"`
}

func (in *BackendTLSPolicyList) DeepCopyInto(out *BackendTLSPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]BackendTLSPolicy, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *BackendTLSPolicyList) DeepCopy() *BackendTLSPolicyList {
	if in == nil {
		return nil
	}
	out := new(BackendTLSPolicyList)
	in.DeepCopyInto(out)
	return out
}

func (in *BackendTLSPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *BackendTLSPolicyList) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.BackendTLSPolicyList"
}

type GRPCRoute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec defines the desired state of GRPCRoute.
	Spec *GRPCRouteSpec `json:"spec,omitempty"`
	// Status defines the current state of GRPCRoute.
	Status *GRPCRouteStatus `json:"status,omitempty"`
}

func (in *GRPCRoute) DeepCopyInto(out *GRPCRoute) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(GRPCRouteSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(GRPCRouteStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *GRPCRoute) DeepCopy() *GRPCRoute {
	if in == nil {
		return nil
	}
	out := new(GRPCRoute)
	in.DeepCopyInto(out)
	return out
}

func (in *GRPCRoute) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *GRPCRoute) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.GRPCRoute"
}

type GRPCRouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []GRPCRoute `json:"items,omitempty"`
}

func (in *GRPCRouteList) DeepCopyInto(out *GRPCRouteList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]GRPCRoute, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *GRPCRouteList) DeepCopy() *GRPCRouteList {
	if in == nil {
		return nil
	}
	out := new(GRPCRouteList)
	in.DeepCopyInto(out)
	return out
}

func (in *GRPCRouteList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *GRPCRouteList) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.GRPCRouteList"
}

type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec defines the desired state of Gateway.
	Spec GatewaySpec `json:"spec"`
	// Status defines the current state of Gateway.
	Status *GatewayStatus `json:"status,omitempty"`
}

func (in *Gateway) DeepCopyInto(out *Gateway) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(GatewayStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *Gateway) DeepCopy() *Gateway {
	if in == nil {
		return nil
	}
	out := new(Gateway)
	in.DeepCopyInto(out)
	return out
}

func (in *Gateway) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *Gateway) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.Gateway"
}

type GatewayClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec defines the desired state of GatewayClass.
	Spec GatewayClassSpec `json:"spec"`
	// Status defines the current state of GatewayClass.
	// Implementations MUST populate status on all GatewayClass resources which
	// specify their controller name.
	Status *GatewayClassStatus `json:"status,omitempty"`
}

func (in *GatewayClass) DeepCopyInto(out *GatewayClass) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(GatewayClassStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *GatewayClass) DeepCopy() *GatewayClass {
	if in == nil {
		return nil
	}
	out := new(GatewayClass)
	in.DeepCopyInto(out)
	return out
}

func (in *GatewayClass) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *GatewayClass) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.GatewayClass"
}

type GatewayClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []GatewayClass `json:"items,omitempty"`
}

func (in *GatewayClassList) DeepCopyInto(out *GatewayClassList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]GatewayClass, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *GatewayClassList) DeepCopy() *GatewayClassList {
	if in == nil {
		return nil
	}
	out := new(GatewayClassList)
	in.DeepCopyInto(out)
	return out
}

func (in *GatewayClassList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *GatewayClassList) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.GatewayClassList"
}

type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Gateway `json:"items,omitempty"`
}

func (in *GatewayList) DeepCopyInto(out *GatewayList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]Gateway, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *GatewayList) DeepCopy() *GatewayList {
	if in == nil {
		return nil
	}
	out := new(GatewayList)
	in.DeepCopyInto(out)
	return out
}

func (in *GatewayList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *GatewayList) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.GatewayList"
}

type HTTPRoute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec defines the desired state of HTTPRoute.
	Spec HTTPRouteSpec `json:"spec"`
	// Status defines the current state of HTTPRoute.
	Status *HTTPRouteStatus `json:"status,omitempty"`
}

func (in *HTTPRoute) DeepCopyInto(out *HTTPRoute) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(HTTPRouteStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *HTTPRoute) DeepCopy() *HTTPRoute {
	if in == nil {
		return nil
	}
	out := new(HTTPRoute)
	in.DeepCopyInto(out)
	return out
}

func (in *HTTPRoute) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *HTTPRoute) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.HTTPRoute"
}

type HTTPRouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []HTTPRoute `json:"items,omitempty"`
}

func (in *HTTPRouteList) DeepCopyInto(out *HTTPRouteList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]HTTPRoute, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *HTTPRouteList) DeepCopy() *HTTPRouteList {
	if in == nil {
		return nil
	}
	out := new(HTTPRouteList)
	in.DeepCopyInto(out)
	return out
}

func (in *HTTPRouteList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *HTTPRouteList) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.HTTPRouteList"
}

type ListenerSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec defines the desired state of ListenerSet.
	Spec ListenerSetSpec `json:"spec"`
	// Status defines the current state of ListenerSet.
	Status *ListenerSetStatus `json:"status,omitempty"`
}

func (in *ListenerSet) DeepCopyInto(out *ListenerSet) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(ListenerSetStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *ListenerSet) DeepCopy() *ListenerSet {
	if in == nil {
		return nil
	}
	out := new(ListenerSet)
	in.DeepCopyInto(out)
	return out
}

func (in *ListenerSet) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *ListenerSet) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.ListenerSet"
}

type ListenerSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ListenerSet `json:"items,omitempty"`
}

func (in *ListenerSetList) DeepCopyInto(out *ListenerSetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]ListenerSet, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *ListenerSetList) DeepCopy() *ListenerSetList {
	if in == nil {
		return nil
	}
	out := new(ListenerSetList)
	in.DeepCopyInto(out)
	return out
}

func (in *ListenerSetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *ListenerSetList) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.ListenerSetList"
}

type ReferenceGrant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec defines the desired state of ReferenceGrant.
	Spec *ReferenceGrantSpec `json:"spec,omitempty"`
}

func (in *ReferenceGrant) DeepCopyInto(out *ReferenceGrant) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(ReferenceGrantSpec)
		(*in).DeepCopyInto(*out)
	}
}

func (in *ReferenceGrant) DeepCopy() *ReferenceGrant {
	if in == nil {
		return nil
	}
	out := new(ReferenceGrant)
	in.DeepCopyInto(out)
	return out
}

func (in *ReferenceGrant) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *ReferenceGrant) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.ReferenceGrant"
}

type ReferenceGrantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ReferenceGrant `json:"items,omitempty"`
}

func (in *ReferenceGrantList) DeepCopyInto(out *ReferenceGrantList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]ReferenceGrant, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *ReferenceGrantList) DeepCopy() *ReferenceGrantList {
	if in == nil {
		return nil
	}
	out := new(ReferenceGrantList)
	in.DeepCopyInto(out)
	return out
}

func (in *ReferenceGrantList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *ReferenceGrantList) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.ReferenceGrantList"
}

type TLSRoute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec defines the desired state of TLSRoute.
	Spec TLSRouteSpec `json:"spec"`
	// Status defines the current state of TLSRoute.
	Status *TLSRouteStatus `json:"status,omitempty"`
}

func (in *TLSRoute) DeepCopyInto(out *TLSRoute) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(TLSRouteStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *TLSRoute) DeepCopy() *TLSRoute {
	if in == nil {
		return nil
	}
	out := new(TLSRoute)
	in.DeepCopyInto(out)
	return out
}

func (in *TLSRoute) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *TLSRoute) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.TLSRoute"
}

type TLSRouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []TLSRoute `json:"items,omitempty"`
}

func (in *TLSRouteList) DeepCopyInto(out *TLSRouteList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]TLSRoute, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *TLSRouteList) DeepCopy() *TLSRouteList {
	if in == nil {
		return nil
	}
	out := new(TLSRouteList)
	in.DeepCopyInto(out)
	return out
}

func (in *TLSRouteList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *TLSRouteList) OpenAPIModelName() string {
	return "sigs.k8s.io.gateway_api.apis.v1.TLSRouteList"
}

type BackendTLSPolicySpec struct {
	// TargetRefs identifies an API object to apply the policy to.
	// Note that this config applies to the entire referenced resource
	// by default, but this default may change in the future to provide
	// a more granular application of the policy.
	// TargetRefs must be _distinct_. This means either that:
	// * They select different targets. If this is the case, then targetRef
	// entries are distinct. In terms of fields, this means that the
	// multi-part key defined by `group`, `kind`, and `name` must
	// be unique across all targetRef entries in the BackendTLSPolicy.
	// * They select different sectionNames in the same target.
	// When more than one BackendTLSPolicy selects the same target and
	// sectionName, implementations MUST determine precedence using the
	// following criteria, continuing on ties:
	// * The older policy by creation timestamp takes precedence. For
	// example, a policy with a creation timestamp of "2021-07-15
	// 01:02:03" MUST be given precedence over a policy with a
	// creation timestamp of "2021-07-15 01:02:04".
	// * The policy appearing first in alphabetical order by {namespace}/{name}.
	// For example, a policy named `foo/bar` is given precedence over a
	// policy named `foo/baz`.
	// For any BackendTLSPolicy that does not take precedence, the
	// implementation MUST ensure the `Accepted` Condition is set to
	// `status: False`, with Reason `Conflicted`.
	// Implementations SHOULD NOT support more than one targetRef at this
	// time. Although the API technically allows for this, the current guidance
	// for conflict resolution and status handling is lacking. Until that can be
	// clarified in a future release, the safest approach is to support a single
	// targetRef.
	// Support Levels:
	// * Extended: Kubernetes Service referenced by HTTPRoute backendRefs.
	// * Implementation-Specific: Services not connected via HTTPRoute, and any
	// other kind of backend. Implementations MAY use BackendTLSPolicy for:
	// - Services not referenced by any Route (e.g., infrastructure services)
	// - Gateway feature backends (e.g., ExternalAuth, rate-limiting services)
	// - Service mesh workload-to-service communication
	// - Other resource types beyond Service
	// Implementations SHOULD aim to ensure that BackendTLSPolicy behavior is consistent,
	// even outside of the extended HTTPRoute -(backendRef) -> Service path.
	// They SHOULD clearly document how BackendTLSPolicy is interpreted in these
	// scenarios, including:
	// - Which resources beyond Service are supported
	// - How the policy is discovered and applied
	// - Any implementation-specific semantics or restrictions
	// Note that this config applies to the entire referenced resource
	// by default, but this default may change in the future to provide
	// a more granular application of the policy.
	TargetRefs []LocalPolicyTargetReferenceWithSectionName `json:"targetRefs,omitempty"`
	// Validation contains backend TLS validation configuration.
	Validation BackendTLSPolicyValidation `json:"validation"`
	// Options are a list of key/value pairs to enable extended TLS
	// configuration for each implementation. For example, configuring the
	// minimum TLS version or supported cipher suites.
	// A set of common keys MAY be defined by the API in the future. To avoid
	// any ambiguity, implementation-specific definitions MUST use
	// domain-prefixed names, such as `example.com/my-custom-option`.
	// Un-prefixed names are reserved for key names defined by Gateway API.
	// Support: Implementation-specific
	Options map[string]string `json:"options,omitempty"`
}

func (in *BackendTLSPolicySpec) DeepCopyInto(out *BackendTLSPolicySpec) {
	*out = *in
	if in.TargetRefs != nil {
		l := make([]LocalPolicyTargetReferenceWithSectionName, len(in.TargetRefs))
		for i := range in.TargetRefs {
			in.TargetRefs[i].DeepCopyInto(&l[i])
		}
		out.TargetRefs = l
	}
	in.Validation.DeepCopyInto(&out.Validation)
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = make(map[string]string, len(*in))
		maps.Copy(*out, *in)
	}
}

func (in *BackendTLSPolicySpec) DeepCopy() *BackendTLSPolicySpec {
	if in == nil {
		return nil
	}
	out := new(BackendTLSPolicySpec)
	in.DeepCopyInto(out)
	return out
}

type PolicyStatus struct {
	// Ancestors is a list of ancestor resources (usually Gateways) that are
	// associated with the policy, and the status of the policy with respect to
	// each ancestor. When this policy attaches to a parent, the controller that
	// manages the parent and the ancestors MUST add an entry to this list when
	// the controller first sees the policy and SHOULD update the entry as
	// appropriate when the relevant ancestor is modified.
	// Note that choosing the relevant ancestor is left to the Policy designers;
	// an important part of Policy design is designing the right object level at
	// which to namespace this status.
	// Note also that implementations MUST ONLY populate ancestor status for
	// the Ancestor resources they are responsible for. Implementations MUST
	// use the ControllerName field to uniquely identify the entries in this list
	// that they are responsible for.
	// Note that to achieve this, the list of PolicyAncestorStatus structs
	// MUST be treated as a map with a composite key, made up of the AncestorRef
	// and ControllerName fields combined.
	// A maximum of 16 ancestors will be represented in this list. An empty list
	// means the Policy is not relevant for any ancestors.
	// If this slice is full, implementations MUST NOT add further entries.
	// Instead they MUST consider the policy unimplementable and signal that
	// on any related resources such as the ancestor that would be referenced
	// here. For example, if this list was full on BackendTLSPolicy, no
	// additional Gateways would be able to reference the Service targeted by
	// the BackendTLSPolicy.
	Ancestors []PolicyAncestorStatus `json:"ancestors,omitempty"`
}

func (in *PolicyStatus) DeepCopyInto(out *PolicyStatus) {
	*out = *in
	if in.Ancestors != nil {
		l := make([]PolicyAncestorStatus, len(in.Ancestors))
		for i := range in.Ancestors {
			in.Ancestors[i].DeepCopyInto(&l[i])
		}
		out.Ancestors = l
	}
}

func (in *PolicyStatus) DeepCopy() *PolicyStatus {
	if in == nil {
		return nil
	}
	out := new(PolicyStatus)
	in.DeepCopyInto(out)
	return out
}

type GRPCRouteSpec struct {
	CommonRouteSpec `json:",inline"`
	// Hostnames defines a set of hostnames to match against the GRPC
	// Host header to select a GRPCRoute to process the request. This matches
	// the RFC 1123 definition of a hostname with 2 notable exceptions:
	// 1. IPs are not allowed.
	// 2. A hostname may be prefixed with a wildcard label (`*.`). The wildcard
	// label MUST appear by itself as the first label.
	// If a hostname is specified by both the Listener and GRPCRoute, there
	// MUST be at least one intersecting hostname for the GRPCRoute to be
	// attached to the Listener. For example:
	// * A Listener with `test.example.com` as the hostname matches GRPCRoutes
	// that have either not specified any hostnames, or have specified at
	// least one of `test.example.com` or `*.example.com`.
	// * A Listener with `*.example.com` as the hostname matches GRPCRoutes
	// that have either not specified any hostnames or have specified at least
	// one hostname that matches the Listener hostname. For example,
	// `test.example.com` and `*.example.com` would both match. On the other
	// hand, `example.com` and `test.example.net` would not match.
	// Hostnames that are prefixed with a wildcard label (`*.`) are interpreted
	// as a suffix match. That means that a match for `*.example.com` would match
	// both `test.example.com`, and `foo.test.example.com`, but not `example.com`.
	// If both the Listener and GRPCRoute have specified hostnames, any
	// GRPCRoute hostnames that do not match the Listener hostname MUST be
	// ignored. For example, if a Listener specified `*.example.com`, and the
	// GRPCRoute specified `test.example.com` and `test.example.net`,
	// `test.example.net` MUST NOT be considered for a match.
	// If both the Listener and GRPCRoute have specified hostnames, and none
	// match with the criteria above, then the GRPCRoute MUST NOT be accepted by
	// the implementation. The implementation MUST raise an 'Accepted' Condition
	// with a status of `False` in the corresponding RouteParentStatus.
	// If a Route (A) of type HTTPRoute or GRPCRoute is attached to a
	// Listener and that listener already has another Route (B) of the other
	// type attached and the intersection of the hostnames of A and B is
	// non-empty, then the implementation MUST accept exactly one of these two
	// routes, determined by the following criteria, in order:
	// * The oldest Route based on creation timestamp.
	// * The Route appearing first in alphabetical order by
	// "{namespace}/{name}".
	// The rejected Route MUST raise an 'Accepted' condition with a status of
	// 'False' in the corresponding RouteParentStatus.
	// Support: Core
	Hostnames []string `json:"hostnames,omitempty"`
	// Rules are a list of GRPC matchers, filters and actions.
	// <gateway:experimental:validation:XValidation:message="Rule name must be unique within the route",rule="self.all(l1, !has(l1.name) || self.exists_one(l2, has(l2.name) && l1.name == l2.name))">
	Rules []GRPCRouteRule `json:"rules,omitempty"`
}

func (in *GRPCRouteSpec) DeepCopyInto(out *GRPCRouteSpec) {
	*out = *in
	out.CommonRouteSpec = in.CommonRouteSpec
	if in.Hostnames != nil {
		t := make([]string, len(in.Hostnames))
		copy(t, in.Hostnames)
		out.Hostnames = t
	}
	if in.Rules != nil {
		l := make([]GRPCRouteRule, len(in.Rules))
		for i := range in.Rules {
			in.Rules[i].DeepCopyInto(&l[i])
		}
		out.Rules = l
	}
}

func (in *GRPCRouteSpec) DeepCopy() *GRPCRouteSpec {
	if in == nil {
		return nil
	}
	out := new(GRPCRouteSpec)
	in.DeepCopyInto(out)
	return out
}

type GRPCRouteStatus struct {
	RouteStatus `json:",inline"`
}

func (in *GRPCRouteStatus) DeepCopyInto(out *GRPCRouteStatus) {
	*out = *in
	out.RouteStatus = in.RouteStatus
}

func (in *GRPCRouteStatus) DeepCopy() *GRPCRouteStatus {
	if in == nil {
		return nil
	}
	out := new(GRPCRouteStatus)
	in.DeepCopyInto(out)
	return out
}

type GatewaySpec struct {
	// GatewayClassName used for this Gateway. This is the name of a
	// GatewayClass resource.
	GatewayClassName string `json:"gatewayClassName"`
	// Listeners associated with this Gateway. Listeners define
	// logical endpoints that are bound on this Gateway's addresses.
	// At least one Listener MUST be specified.
	// ## Distinct Listeners
	// Each Listener in a set of Listeners (for example, in a single Gateway)
	// MUST be _distinct_, in that a traffic flow MUST be able to be assigned to
	// exactly one listener. (This section uses "set of Listeners" rather than
	// "Listeners in a single Gateway" because implementations MAY merge configuration
	// from multiple Gateways onto a single data plane, and these rules _also_
	// apply in that case).
	// Practically, this means that each listener in a set MUST have a unique
	// combination of Port, Protocol, and, if supported by the protocol, Hostname.
	// Some combinations of port, protocol, and TLS settings are considered
	// Core support and MUST be supported by implementations based on the objects
	// they support:
	// HTTPRoute
	// 1. HTTPRoute, Port: 80, Protocol: HTTP
	// 2. HTTPRoute, Port: 443, Protocol: HTTPS, TLS Mode: Terminate, TLS keypair provided
	// TLSRoute
	// 1. TLSRoute, Port: 443, Protocol: TLS, TLS Mode: Passthrough
	// "Distinct" Listeners have the following property:
	// **The implementation can match inbound requests to a single distinct
	// Listener**.
	// When multiple Listeners share values for fields (for
	// example, two Listeners with the same Port value), the implementation
	// can match requests to only one of the Listeners using other
	// Listener fields.
	// When multiple listeners have the same value for the Protocol field, then
	// each of the Listeners with matching Protocol values MUST have different
	// values for other fields.
	// The set of fields that MUST be different for a Listener differs per protocol.
	// The following rules define the rules for what fields MUST be considered for
	// Listeners to be distinct with each protocol currently defined in the
	// Gateway API spec.
	// The set of listeners that all share a protocol value MUST have _different_
	// values for _at least one_ of these fields to be distinct:
	// * **HTTP, HTTPS, TLS**: Port, Hostname
	// * **TCP, UDP**: Port
	// One **very** important rule to call out involves what happens when an
	// implementation:
	// * Supports TCP protocol Listeners, as well as HTTP, HTTPS, or TLS protocol
	// Listeners, and
	// * sees HTTP, HTTPS, or TLS protocols with the same `port` as one with TCP
	// Protocol.
	// In this case all the Listeners that share a port with the
	// TCP Listener are not distinct and so MUST NOT be accepted.
	// If an implementation does not support TCP Protocol Listeners, then the
	// previous rule does not apply, and the TCP Listeners SHOULD NOT be
	// accepted.
	// Note that the `tls` field is not used for determining if a listener is distinct, because
	// Listeners that _only_ differ on TLS config will still conflict in all cases.
	// ### Listeners that are distinct only by Hostname
	// When the Listeners are distinct based only on Hostname, inbound request
	// hostnames MUST match from the most specific to least specific Hostname
	// values to choose the correct Listener and its associated set of Routes.
	// Exact matches MUST be processed before wildcard matches, and wildcard
	// matches MUST be processed before fallback (empty Hostname value)
	// matches. For example, `"foo.example.com"` takes precedence over
	// `"*.example.com"`, and `"*.example.com"` takes precedence over `""`.
	// Additionally, if there are multiple wildcard entries, more specific
	// wildcard entries must be processed before less specific wildcard entries.
	// For example, `"*.foo.example.com"` takes precedence over `"*.example.com"`.
	// The precise definition here is that the higher the number of dots in the
	// hostname to the right of the wildcard character, the higher the precedence.
	// The wildcard character will match any number of characters _and dots_ to
	// the left, however, so `"*.example.com"` will match both
	// `"foo.bar.example.com"` _and_ `"bar.example.com"`.
	// ## Handling indistinct Listeners
	// If a set of Listeners contains Listeners that are not distinct, then those
	// Listeners are _Conflicted_, and the implementation MUST set the "Conflicted"
	// condition in the Listener Status to "True".
	// The words "indistinct" and "conflicted" are considered equivalent for the
	// purpose of this documentation.
	// Implementations MAY choose to accept a Gateway with some Conflicted
	// Listeners only if they only accept the partial Listener set that contains
	// no Conflicted Listeners.
	// Specifically, an implementation MAY accept a partial Listener set subject to
	// the following rules:
	// * The implementation MUST NOT pick one conflicting Listener as the winner.
	// ALL indistinct Listeners must not be accepted for processing.
	// * At least one distinct Listener MUST be present, or else the Gateway effectively
	// contains _no_ Listeners, and must be rejected from processing as a whole.
	// The implementation MUST set a "ListenersNotValid" condition on the
	// Gateway Status when the Gateway contains Conflicted Listeners whether or
	// not they accept the Gateway. That Condition SHOULD clearly
	// indicate in the Message which Listeners are conflicted, and which are
	// Accepted. Additionally, the Listener status for those listeners SHOULD
	// indicate which Listeners are conflicted and not Accepted.
	// ## General Listener behavior
	// Note that, for all distinct Listeners, requests SHOULD match at most one Listener.
	// For example, if Listeners are defined for "foo.example.com" and "*.example.com", a
	// request to "foo.example.com" SHOULD only be routed using routes attached
	// to the "foo.example.com" Listener (and not the "*.example.com" Listener).
	// This concept is known as "Listener Isolation", and it is an Extended feature
	// of Gateway API. Implementations that do not support Listener Isolation MUST
	// clearly document this, and MUST NOT claim support for the
	// `GatewayHTTPListenerIsolation` feature.
	// Implementations that _do_ support Listener Isolation SHOULD claim support
	// for the Extended `GatewayHTTPListenerIsolation` feature and pass the associated
	// conformance tests.
	// ## Compatible Listeners
	// A Gateway's Listeners are considered _compatible_ if:
	// 1. They are distinct.
	// 2. The implementation can serve them in compliance with the Addresses
	// requirement that all Listeners are available on all assigned
	// addresses.
	// Compatible combinations in Extended support are expected to vary across
	// implementations. A combination that is compatible for one implementation
	// may not be compatible for another.
	// For example, an implementation that cannot serve both TCP and UDP listeners
	// on the same address, or cannot mix HTTPS and generic TLS listens on the same port
	// would not consider those cases compatible, even though they are distinct.
	// Implementations MAY merge separate Gateways onto a single set of
	// Addresses if all Listeners across all Gateways are compatible.
	// In a future release the MinItems=1 requirement MAY be dropped.
	// Support: Core
	Listeners []Listener `json:"listeners,omitempty"`
	// Addresses requested for this Gateway. This is optional and behavior can
	// depend on the implementation. If a value is set in the spec and the
	// requested address is invalid or unavailable, the implementation MUST
	// indicate this in an associated entry in GatewayStatus.Conditions.
	// The Addresses field represents a request for the address(es) on the
	// "outside of the Gateway", that traffic bound for this Gateway will use.
	// This could be the IP address or hostname of an external load balancer or
	// other networking infrastructure, or some other address that traffic will
	// be sent to.
	// If no Addresses are specified, the implementation MAY schedule the
	// Gateway in an implementation-specific manner, assigning an appropriate
	// set of Addresses.
	// The implementation MUST bind all Listeners to every GatewayAddress that
	// it assigns to the Gateway and add a corresponding entry in
	// GatewayStatus.Addresses.
	// Support: Extended
	// <gateway:validateIPAddress>
	Addresses []GatewaySpecAddress `json:"addresses,omitempty"`
	// Infrastructure defines infrastructure level attributes about this Gateway instance.
	// Support: Extended
	Infrastructure *GatewayInfrastructure `json:"infrastructure,omitempty"`
	// AllowedListeners defines which ListenerSets can be attached to this Gateway.
	// The default value is to allow no ListenerSets.
	AllowedListeners *AllowedListeners `json:"allowedListeners,omitempty"`
	// TLS specifies frontend and backend tls configuration for entire gateway.
	// Support: Extended
	TLS *GatewayTLSConfig `json:"tls,omitempty"`
	// DefaultScope, when set, configures the Gateway as a default Gateway,
	// meaning it will dynamically and implicitly have Routes (e.g. HTTPRoute)
	// attached to it, according to the scope configured here.
	// If unset (the default) or set to None, the Gateway will not act as a
	// default Gateway; if set, the Gateway will claim any Route with a
	// matching scope set in its UseDefaultGateway field, subject to the usual
	// rules about which routes the Gateway can attach to.
	// Think carefully before using this functionality! While the normal rules
	// about which Route can apply are still enforced, it is simply easier for
	// the wrong Route to be accidentally attached to this Gateway in this
	// configuration. If the Gateway operator is not also the operator in
	// control of the scope (e.g. namespace) with tight controls and checks on
	// what kind of workloads and Routes get added in that scope, we strongly
	// recommend not using this just because it seems convenient, and instead
	// stick to direct Route attachment.
	// <gateway:experimental>
	DefaultScope GatewayDefaultScope `json:"defaultScope,omitempty"`
}

func (in *GatewaySpec) DeepCopyInto(out *GatewaySpec) {
	*out = *in
	if in.Listeners != nil {
		l := make([]Listener, len(in.Listeners))
		for i := range in.Listeners {
			in.Listeners[i].DeepCopyInto(&l[i])
		}
		out.Listeners = l
	}
	if in.Addresses != nil {
		l := make([]GatewaySpecAddress, len(in.Addresses))
		for i := range in.Addresses {
			in.Addresses[i].DeepCopyInto(&l[i])
		}
		out.Addresses = l
	}
	if in.Infrastructure != nil {
		in, out := &in.Infrastructure, &out.Infrastructure
		*out = new(GatewayInfrastructure)
		(*in).DeepCopyInto(*out)
	}
	if in.AllowedListeners != nil {
		in, out := &in.AllowedListeners, &out.AllowedListeners
		*out = new(AllowedListeners)
		(*in).DeepCopyInto(*out)
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(GatewayTLSConfig)
		(*in).DeepCopyInto(*out)
	}
}

func (in *GatewaySpec) DeepCopy() *GatewaySpec {
	if in == nil {
		return nil
	}
	out := new(GatewaySpec)
	in.DeepCopyInto(out)
	return out
}

type GatewayStatus struct {
	// Addresses lists the network addresses that have been bound to the
	// Gateway.
	// This list may differ from the addresses provided in the spec under some
	// conditions:
	// * no addresses are specified, all addresses are dynamically assigned
	// * a combination of specified and dynamic addresses are assigned
	// * a specified address was unusable (e.g. already in use)
	// <gateway:validateIPAddress>
	Addresses []GatewayStatusAddress `json:"addresses,omitempty"`
	// Conditions describe the current conditions of the Gateway.
	// Implementations should prefer to express Gateway conditions
	// using the `GatewayConditionType` and `GatewayConditionReason`
	// constants so that operators and tools can converge on a common
	// vocabulary to describe Gateway state.
	// Known condition types are:
	// * "Accepted"
	// * "Programmed"
	// * "Ready"
	// <gateway:util:excludeFromCRD>
	// Notes for implementors:
	// Conditions are a listType `map`, which means that they function like a
	// map with a key of the `type` field _in the k8s apiserver_.
	// This means that implementations must obey some rules when updating this
	// section.
	// * Implementations MUST perform a read-modify-write cycle on this field
	// before modifying it. That is, when modifying this field, implementations
	// must be confident they have fetched the most recent version of this field,
	// and ensure that changes they make are on that recent version.
	// * Implementations MUST NOT remove or reorder Conditions that they are not
	// directly responsible for. For example, if an implementation sees a Condition
	// with type `special.io/SomeField`, it MUST NOT remove, change or update that
	// Condition.
	// * Implementations MUST always _merge_ changes into Conditions of the same Type,
	// rather than creating more than one Condition of the same Type.
	// * Implementations MUST always update the `observedGeneration` field of the
	// Condition to the `metadata.generation` of the Gateway at the time of update creation.
	// * If the `observedGeneration` of a Condition is _greater than_ the value the
	// implementation knows about, then it MUST NOT perform the update on that Condition,
	// but must wait for a future reconciliation and status update. (The assumption is that
	// the implementation's copy of the object is stale and an update will be re-triggered
	// if relevant.)
	// </gateway:util:excludeFromCRD>
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// Listeners provide status for each unique listener port defined in the Spec.
	Listeners []ListenerStatus `json:"listeners,omitempty"`
	// AttachedListenerSets represents the total number of ListenerSets that have been
	// successfully attached to this Gateway.
	// A ListenerSet is successfully attached to a Gateway when all the following conditions are met:
	// - The ListenerSet is selected by the Gateway's AllowedListeners field
	// - The ListenerSet has a valid ParentRef selecting the Gateway
	// - The ListenerSet's status has the condition "Accepted: true"
	// Uses for this field include troubleshooting AttachedListenerSets attachment and
	// measuring blast radius/impact of changes to a Gateway.
	AttachedListenerSets int `json:"attachedListenerSets,omitempty"`
}

func (in *GatewayStatus) DeepCopyInto(out *GatewayStatus) {
	*out = *in
	if in.Addresses != nil {
		l := make([]GatewayStatusAddress, len(in.Addresses))
		for i := range in.Addresses {
			in.Addresses[i].DeepCopyInto(&l[i])
		}
		out.Addresses = l
	}
	if in.Conditions != nil {
		l := make([]metav1.Condition, len(in.Conditions))
		for i := range in.Conditions {
			in.Conditions[i].DeepCopyInto(&l[i])
		}
		out.Conditions = l
	}
	if in.Listeners != nil {
		l := make([]ListenerStatus, len(in.Listeners))
		for i := range in.Listeners {
			in.Listeners[i].DeepCopyInto(&l[i])
		}
		out.Listeners = l
	}
}

func (in *GatewayStatus) DeepCopy() *GatewayStatus {
	if in == nil {
		return nil
	}
	out := new(GatewayStatus)
	in.DeepCopyInto(out)
	return out
}

type GatewayClassSpec struct {
	// ControllerName is the name of the controller that is managing Gateways of
	// this class. The value of this field MUST be a domain prefixed path.
	// Example: "example.net/gateway-controller".
	// This field is not mutable and cannot be empty.
	// Support: Core
	ControllerName string `json:"controllerName"`
	// ParametersRef is a reference to a resource that contains the configuration
	// parameters corresponding to the GatewayClass. This is optional if the
	// controller does not require any additional configuration.
	// ParametersRef can reference a standard Kubernetes resource, i.e. ConfigMap,
	// or an implementation-specific custom resource. The resource can be
	// cluster-scoped or namespace-scoped.
	// If the referent cannot be found, refers to an unsupported kind, or when
	// the data within that resource is malformed, the GatewayClass SHOULD be
	// rejected with the "Accepted" status condition set to "False" and an
	// "InvalidParameters" reason.
	// A Gateway for this GatewayClass may provide its own `parametersRef`. When both are specified,
	// the merging behavior is implementation specific.
	// It is generally recommended that GatewayClass provides defaults that can be overridden by a Gateway.
	// Support: Implementation-specific
	ParametersRef *ParametersReference `json:"parametersRef,omitempty"`
	// Description helps describe a GatewayClass with more details.
	Description string `json:"description,omitempty"`
}

func (in *GatewayClassSpec) DeepCopyInto(out *GatewayClassSpec) {
	*out = *in
	if in.ParametersRef != nil {
		in, out := &in.ParametersRef, &out.ParametersRef
		*out = new(ParametersReference)
		(*in).DeepCopyInto(*out)
	}
}

func (in *GatewayClassSpec) DeepCopy() *GatewayClassSpec {
	if in == nil {
		return nil
	}
	out := new(GatewayClassSpec)
	in.DeepCopyInto(out)
	return out
}

type GatewayClassStatus struct {
	// Conditions is the current status from the controller for
	// this GatewayClass.
	// Controllers should prefer to publish conditions using values
	// of GatewayClassConditionType for the type of each Condition.
	// <gateway:util:excludeFromCRD>
	// Notes for implementors:
	// Conditions are a listType `map`, which means that they function like a
	// map with a key of the `type` field _in the k8s apiserver_.
	// This means that implementations must obey some rules when updating this
	// section.
	// * Implementations MUST perform a read-modify-write cycle on this field
	// before modifying it. That is, when modifying this field, implementations
	// must be confident they have fetched the most recent version of this field,
	// and ensure that changes they make are on that recent version.
	// * Implementations MUST NOT remove or reorder Conditions that they are not
	// directly responsible for. For example, if an implementation sees a Condition
	// with type `special.io/SomeField`, it MUST NOT remove, change or update that
	// Condition.
	// * Implementations MUST always _merge_ changes into Conditions of the same Type,
	// rather than creating more than one Condition of the same Type.
	// * Implementations MUST always update the `observedGeneration` field of the
	// Condition to the `metadata.generation` of the Gateway at the time of update creation.
	// * If the `observedGeneration` of a Condition is _greater than_ the value the
	// implementation knows about, then it MUST NOT perform the update on that Condition,
	// but must wait for a future reconciliation and status update. (The assumption is that
	// the implementation's copy of the object is stale and an update will be re-triggered
	// if relevant.)
	// </gateway:util:excludeFromCRD>
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// SupportedFeatures is the set of features the GatewayClass support.
	// It MUST be sorted in ascending alphabetical order by the Name key.
	SupportedFeatures []SupportedFeature `json:"supportedFeatures,omitempty"`
}

func (in *GatewayClassStatus) DeepCopyInto(out *GatewayClassStatus) {
	*out = *in
	if in.Conditions != nil {
		l := make([]metav1.Condition, len(in.Conditions))
		for i := range in.Conditions {
			in.Conditions[i].DeepCopyInto(&l[i])
		}
		out.Conditions = l
	}
	if in.SupportedFeatures != nil {
		l := make([]SupportedFeature, len(in.SupportedFeatures))
		for i := range in.SupportedFeatures {
			in.SupportedFeatures[i].DeepCopyInto(&l[i])
		}
		out.SupportedFeatures = l
	}
}

func (in *GatewayClassStatus) DeepCopy() *GatewayClassStatus {
	if in == nil {
		return nil
	}
	out := new(GatewayClassStatus)
	in.DeepCopyInto(out)
	return out
}

type HTTPRouteSpec struct {
	CommonRouteSpec `json:",inline"`
	// Hostnames defines a set of hostnames that should match against the HTTP Host
	// header to select a HTTPRoute used to process the request. Implementations
	// MUST ignore any port value specified in the HTTP Host header while
	// performing a match and (absent of any applicable header modification
	// configuration) MUST forward this header unmodified to the backend.
	// Valid values for Hostnames are determined by RFC 1123 definition of a
	// hostname with 2 notable exceptions:
	// 1. IPs are not allowed.
	// 2. A hostname may be prefixed with a wildcard label (`*.`). The wildcard
	// label must appear by itself as the first label.
	// If a hostname is specified by both the Listener and HTTPRoute, there
	// must be at least one intersecting hostname for the HTTPRoute to be
	// attached to the Listener. For example:
	// * A Listener with `test.example.com` as the hostname matches HTTPRoutes
	// that have either not specified any hostnames, or have specified at
	// least one of `test.example.com` or `*.example.com`.
	// * A Listener with `*.example.com` as the hostname matches HTTPRoutes
	// that have either not specified any hostnames or have specified at least
	// one hostname that matches the Listener hostname. For example,
	// `*.example.com`, `test.example.com`, and `foo.test.example.com` would
	// all match. On the other hand, `example.com` and `test.example.net` would
	// not match.
	// Hostnames that are prefixed with a wildcard label (`*.`) are interpreted
	// as a suffix match. That means that a match for `*.example.com` would match
	// both `test.example.com`, and `foo.test.example.com`, but not `example.com`.
	// If both the Listener and HTTPRoute have specified hostnames, any
	// HTTPRoute hostnames that do not match the Listener hostname MUST be
	// ignored. For example, if a Listener specified `*.example.com`, and the
	// HTTPRoute specified `test.example.com` and `test.example.net`,
	// `test.example.net` must not be considered for a match.
	// If both the Listener and HTTPRoute have specified hostnames, and none
	// match with the criteria above, then the HTTPRoute is not accepted. The
	// implementation must raise an 'Accepted' Condition with a status of
	// `False` in the corresponding RouteParentStatus.
	// In the event that multiple HTTPRoutes specify intersecting hostnames (e.g.
	// overlapping wildcard matching and exact matching hostnames), precedence must
	// be given to rules from the HTTPRoute with the largest number of:
	// * Characters in a matching non-wildcard hostname.
	// * Characters in a matching hostname.
	// If ties exist across multiple Routes, the matching precedence rules for
	// HTTPRouteMatches takes over.
	// Support: Core
	Hostnames []string `json:"hostnames,omitempty"`
	// Rules are a list of HTTP matchers, filters and actions.
	// <gateway:experimental:validation:XValidation:message="Rule name must be unique within the route",rule="self.all(l1, !has(l1.name) || self.exists_one(l2, has(l2.name) && l1.name == l2.name))">
	Rules []HTTPRouteRule `json:"rules,omitempty"`
}

func (in *HTTPRouteSpec) DeepCopyInto(out *HTTPRouteSpec) {
	*out = *in
	out.CommonRouteSpec = in.CommonRouteSpec
	if in.Hostnames != nil {
		t := make([]string, len(in.Hostnames))
		copy(t, in.Hostnames)
		out.Hostnames = t
	}
	if in.Rules != nil {
		l := make([]HTTPRouteRule, len(in.Rules))
		for i := range in.Rules {
			in.Rules[i].DeepCopyInto(&l[i])
		}
		out.Rules = l
	}
}

func (in *HTTPRouteSpec) DeepCopy() *HTTPRouteSpec {
	if in == nil {
		return nil
	}
	out := new(HTTPRouteSpec)
	in.DeepCopyInto(out)
	return out
}

type HTTPRouteStatus struct {
	RouteStatus `json:",inline"`
}

func (in *HTTPRouteStatus) DeepCopyInto(out *HTTPRouteStatus) {
	*out = *in
	out.RouteStatus = in.RouteStatus
}

func (in *HTTPRouteStatus) DeepCopy() *HTTPRouteStatus {
	if in == nil {
		return nil
	}
	out := new(HTTPRouteStatus)
	in.DeepCopyInto(out)
	return out
}

type ListenerSetSpec struct {
	// ParentRef references the Gateway that the listeners are attached to.
	ParentRef *ParentGatewayReference `json:"parentRef,omitempty"`
	// Listeners associated with this ListenerSet. Listeners define
	// logical endpoints that are bound on this referenced parent Gateway's addresses.
	// Listeners in a `Gateway` and their attached `ListenerSets` are concatenated
	// as a list when programming the underlying infrastructure. Each listener
	// name does not need to be unique across the Gateway and ListenerSets.
	// See ListenerEntry.Name for more details.
	// Implementations MUST treat the parent Gateway as having the merged
	// list of all listeners from itself and attached ListenerSets using
	// the following precedence:
	// 1. "parent" Gateway
	// 2. ListenerSet ordered by creation time (oldest first)
	// 3. ListenerSet ordered alphabetically by "{namespace}/{name}".
	// An implementation MAY reject listeners by setting the ListenerEntryStatus
	// `Accepted` condition to False with the Reason `TooManyListeners`
	// If a listener has a conflict, this will be reported in the
	// Status.ListenerEntryStatus setting the `Conflicted` condition to True.
	// Implementations SHOULD be cautious about what information from the
	// parent or siblings are reported to avoid accidentally leaking
	// sensitive information that the child would not otherwise have access
	// to. This can include contents of secrets etc.
	Listeners []ListenerEntry `json:"listeners,omitempty"`
}

func (in *ListenerSetSpec) DeepCopyInto(out *ListenerSetSpec) {
	*out = *in
	if in.ParentRef != nil {
		in, out := &in.ParentRef, &out.ParentRef
		*out = new(ParentGatewayReference)
		(*in).DeepCopyInto(*out)
	}
	if in.Listeners != nil {
		l := make([]ListenerEntry, len(in.Listeners))
		for i := range in.Listeners {
			in.Listeners[i].DeepCopyInto(&l[i])
		}
		out.Listeners = l
	}
}

func (in *ListenerSetSpec) DeepCopy() *ListenerSetSpec {
	if in == nil {
		return nil
	}
	out := new(ListenerSetSpec)
	in.DeepCopyInto(out)
	return out
}

type ListenerSetStatus struct {
	// Conditions describe the current conditions of the ListenerSet.
	// Implementations MUST express ListenerSet conditions using the
	// `ListenerSetConditionType` and `ListenerSetConditionReason`
	// constants so that operators and tools can converge on a common
	// vocabulary to describe ListenerSet state.
	// Known condition types are:
	// * "Accepted"
	// * "Programmed"
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// Listeners provide status for each unique listener port defined in the Spec.
	Listeners []ListenerEntryStatus `json:"listeners,omitempty"`
}

func (in *ListenerSetStatus) DeepCopyInto(out *ListenerSetStatus) {
	*out = *in
	if in.Conditions != nil {
		l := make([]metav1.Condition, len(in.Conditions))
		for i := range in.Conditions {
			in.Conditions[i].DeepCopyInto(&l[i])
		}
		out.Conditions = l
	}
	if in.Listeners != nil {
		l := make([]ListenerEntryStatus, len(in.Listeners))
		for i := range in.Listeners {
			in.Listeners[i].DeepCopyInto(&l[i])
		}
		out.Listeners = l
	}
}

func (in *ListenerSetStatus) DeepCopy() *ListenerSetStatus {
	if in == nil {
		return nil
	}
	out := new(ListenerSetStatus)
	in.DeepCopyInto(out)
	return out
}

type ReferenceGrantSpec struct {
	// From describes the trusted namespaces and kinds that can reference the
	// resources described in "To". Each entry in this list MUST be considered
	// to be an additional place that references can be valid from, or to put
	// this another way, entries MUST be combined using OR.
	// Support: Core
	From []ReferenceGrantFrom `json:"from,omitempty"`
	// To describes the resources that may be referenced by the resources
	// described in "From". Each entry in this list MUST be considered to be an
	// additional place that references can be valid to, or to put this another
	// way, entries MUST be combined using OR.
	// Support: Core
	To []ReferenceGrantTo `json:"to,omitempty"`
}

func (in *ReferenceGrantSpec) DeepCopyInto(out *ReferenceGrantSpec) {
	*out = *in
	if in.From != nil {
		l := make([]ReferenceGrantFrom, len(in.From))
		for i := range in.From {
			in.From[i].DeepCopyInto(&l[i])
		}
		out.From = l
	}
	if in.To != nil {
		l := make([]ReferenceGrantTo, len(in.To))
		for i := range in.To {
			in.To[i].DeepCopyInto(&l[i])
		}
		out.To = l
	}
}

func (in *ReferenceGrantSpec) DeepCopy() *ReferenceGrantSpec {
	if in == nil {
		return nil
	}
	out := new(ReferenceGrantSpec)
	in.DeepCopyInto(out)
	return out
}

type TLSRouteSpec struct {
	CommonRouteSpec `json:",inline"`
	// Hostnames defines a set of SNI hostnames that should match against the
	// SNI attribute of TLS ClientHello message in TLS handshake. This matches
	// the RFC 1123 definition of a hostname with 2 notable exceptions:
	// 1. IPs are not allowed in SNI hostnames per RFC 6066.
	// 2. A hostname may be prefixed with a wildcard label (`*.`). The wildcard
	// label must appear by itself as the first label.
	// <gateway:util:excludeFromCRD>
	// If a hostname is specified by both the Listener and TLSRoute, there
	// must be at least one intersecting hostname for the TLSRoute to be
	// attached to the Listener. For example:
	// * A Listener with `test.example.com` as the hostname matches TLSRoutes
	// that have specified at least one of `test.example.com` or
	// `*.example.com`.
	// * A Listener with `*.example.com` as the hostname matches TLSRoutes
	// that have specified at least one hostname that matches the Listener
	// hostname. For example, `test.example.com` and `*.example.com` would both
	// match. On the other hand, `example.com` and `test.example.net` would not
	// match.
	// * A listener with `something.example.com` as the hostname matches a
	// TLSRoute with hostname `*.example.com`.
	// If both the Listener and TLSRoute have specified hostnames, any
	// TLSRoute hostnames that do not match any Listener hostname MUST be
	// ignored. For example, if a Listener specified `*.example.com`, and the
	// TLSRoute specified `test.example.com` and `test.example.net`,
	// `test.example.net` must not be considered for a match.
	// If both the Listener and TLSRoute have specified hostnames, and none
	// match with the criteria above, then the TLSRoute is not accepted for that
	// Listener. If the TLSRoute does not match any Listener on its parent, the
	// implementation must raise an 'Accepted' Condition with a status of
	// `False` in the corresponding RouteParentStatus.
	// A Listener MUST be have protocol set to TLS when a TLSRoute attaches to it. The
	// implementation MUST raise an 'Accepted' Condition with a status of
	// `False` in the corresponding RouteParentStatus with the reason
	// of "UnsupportedValue" in case a Listener of the wrong type is used.
	// Core: Listener with `protocol` `TLS` and `tls.mode` `Passthrough`.
	// Extended: Listener with `protocol` `TLS` and `tls.mode` `Terminate`. The feature name for this Extended feature is `TLSRouteTermination`.
	// </gateway:util:excludeFromCRD>
	Hostnames []string `json:"hostnames,omitempty"`
	// Rules are a list of actions.
	Rules []TLSRouteRule `json:"rules,omitempty"`
}

func (in *TLSRouteSpec) DeepCopyInto(out *TLSRouteSpec) {
	*out = *in
	out.CommonRouteSpec = in.CommonRouteSpec
	if in.Hostnames != nil {
		t := make([]string, len(in.Hostnames))
		copy(t, in.Hostnames)
		out.Hostnames = t
	}
	if in.Rules != nil {
		l := make([]TLSRouteRule, len(in.Rules))
		for i := range in.Rules {
			in.Rules[i].DeepCopyInto(&l[i])
		}
		out.Rules = l
	}
}

func (in *TLSRouteSpec) DeepCopy() *TLSRouteSpec {
	if in == nil {
		return nil
	}
	out := new(TLSRouteSpec)
	in.DeepCopyInto(out)
	return out
}

type TLSRouteStatus struct {
	RouteStatus `json:",inline"`
}

func (in *TLSRouteStatus) DeepCopyInto(out *TLSRouteStatus) {
	*out = *in
	out.RouteStatus = in.RouteStatus
}

func (in *TLSRouteStatus) DeepCopy() *TLSRouteStatus {
	if in == nil {
		return nil
	}
	out := new(TLSRouteStatus)
	in.DeepCopyInto(out)
	return out
}

type LocalPolicyTargetReferenceWithSectionName struct {
	LocalPolicyTargetReference `json:",inline"`
	// SectionName is the name of a section within the target resource. When
	// unspecified, this targetRef targets the entire resource. In the following
	// resources, SectionName is interpreted as the following:
	// * Gateway: Listener name
	// * HTTPRoute: HTTPRouteRule name
	// * Service: Port name
	// If a SectionName is specified, but does not exist on the targeted object,
	// the Policy must fail to attach, and the policy implementation should record
	// a `ResolvedRefs` or similar Condition in the Policy's status.
	SectionName string `json:"sectionName"`
}

func (in *LocalPolicyTargetReferenceWithSectionName) DeepCopyInto(out *LocalPolicyTargetReferenceWithSectionName) {
	*out = *in
	out.LocalPolicyTargetReference = in.LocalPolicyTargetReference
}

func (in *LocalPolicyTargetReferenceWithSectionName) DeepCopy() *LocalPolicyTargetReferenceWithSectionName {
	if in == nil {
		return nil
	}
	out := new(LocalPolicyTargetReferenceWithSectionName)
	in.DeepCopyInto(out)
	return out
}

type BackendTLSPolicyValidation struct {
	// CACertificateRefs contains one or more references to Kubernetes objects that
	// contain a PEM-encoded TLS CA certificate bundle, which is used to
	// validate a TLS handshake between the Gateway and backend Pod.
	// If CACertificateRefs is empty or unspecified, then WellKnownCACertificates must be
	// specified. Only one of CACertificateRefs or WellKnownCACertificates may be specified,
	// not both. If CACertificateRefs is empty or unspecified, the configuration for
	// WellKnownCACertificates MUST be honored instead if supported by the implementation.
	// A CACertificateRef is invalid if:
	// * It refers to a resource that cannot be resolved (e.g., the referenced resource
	// does not exist) or is misconfigured (e.g., a ConfigMap does not contain a key
	// named `ca.crt`). In this case, the Reason must be set to `InvalidCACertificateRef`
	// and the Message of the Condition must indicate which reference is invalid and why.
	// * It refers to an unknown or unsupported kind of resource. In this case, the Reason
	// must be set to `InvalidKind` and the Message of the Condition must explain which
	// kind of resource is unknown or unsupported.
	// * It refers to a resource in another namespace. This may change in future
	// spec updates.
	// Implementations MAY choose to perform further validation of the certificate
	// content (e.g., checking expiry or enforcing specific formats). In such cases,
	// an implementation-specific Reason and Message must be set for the invalid reference.
	// In all cases, the implementation MUST ensure the `ResolvedRefs` Condition on
	// the BackendTLSPolicy is set to `status: False`, with a Reason and Message
	// that indicate the cause of the error. Connections using an invalid
	// CACertificateRef MUST fail, and the client MUST receive an HTTP 5xx error
	// response. If ALL CACertificateRefs are invalid, the implementation MUST also
	// ensure the `Accepted` Condition on the BackendTLSPolicy is set to
	// `status: False`, with a Reason `NoValidCACertificate`.
	// A single CACertificateRef to a Kubernetes ConfigMap kind has "Core" support.
	// Implementations MAY choose to support attaching multiple certificates to
	// a backend, but this behavior is implementation-specific.
	// Support: Core - An optional single reference to a Kubernetes ConfigMap,
	// with the CA certificate in a key named `ca.crt`.
	// Support: Implementation-specific - More than one reference, other kinds
	// of resources, or a single reference that includes multiple certificates.
	CACertificateRefs []LocalObjectReference `json:"caCertificateRefs,omitempty"`
	// WellKnownCACertificates specifies whether a well-known set of CA certificates
	// may be used in the TLS handshake between the gateway and backend pod.
	// If WellKnownCACertificates is unspecified or empty (""), then CACertificateRefs
	// must be specified with at least one entry for a valid configuration. Only one of
	// CACertificateRefs or WellKnownCACertificates may be specified, not both.
	// If an implementation does not support the WellKnownCACertificates field, or
	// the supplied value is not recognized, the implementation MUST ensure the
	// `Accepted` Condition on the BackendTLSPolicy is set to `status: False`, with
	// a Reason `Invalid`.
	// Valid values include:
	// * "System" - indicates that well-known system CA certificates should be used.
	// Implementations MAY define their own sets of CA certificates. Such definitions
	// MUST use an implementation-specific, prefixed name, such as
	// `mycompany.com/my-custom-ca-certificates`.
	// Support: Implementation-specific
	WellKnownCACertificates WellKnownCACertificatesType `json:"wellKnownCACertificates,omitempty"`
	// Hostname is used for two purposes in the connection between Gateways and
	// backends:
	// 1. Hostname MUST be used as the SNI to connect to the backend (RFC 6066).
	// 2. Hostname MUST be used for authentication and MUST match the certificate
	// served by the matching backend, unless SubjectAltNames is specified.
	// 3. If SubjectAltNames are specified, Hostname can be used for certificate selection
	// but MUST NOT be used for authentication. If you want to use the value
	// of the Hostname field for authentication, you MUST add it to the SubjectAltNames list.
	// Support: Core
	Hostname string `json:"hostname"`
	// SubjectAltNames contains one or more Subject Alternative Names.
	// When specified the certificate served from the backend MUST
	// have at least one Subject Alternate Name matching one of the specified SubjectAltNames.
	// Support: Extended
	SubjectAltNames []SubjectAltName `json:"subjectAltNames,omitempty"`
}

func (in *BackendTLSPolicyValidation) DeepCopyInto(out *BackendTLSPolicyValidation) {
	*out = *in
	if in.CACertificateRefs != nil {
		l := make([]LocalObjectReference, len(in.CACertificateRefs))
		for i := range in.CACertificateRefs {
			in.CACertificateRefs[i].DeepCopyInto(&l[i])
		}
		out.CACertificateRefs = l
	}
	if in.SubjectAltNames != nil {
		l := make([]SubjectAltName, len(in.SubjectAltNames))
		for i := range in.SubjectAltNames {
			in.SubjectAltNames[i].DeepCopyInto(&l[i])
		}
		out.SubjectAltNames = l
	}
}

func (in *BackendTLSPolicyValidation) DeepCopy() *BackendTLSPolicyValidation {
	if in == nil {
		return nil
	}
	out := new(BackendTLSPolicyValidation)
	in.DeepCopyInto(out)
	return out
}

type PolicyAncestorStatus struct {
	// AncestorRef corresponds with a ParentRef in the spec that this
	// PolicyAncestorStatus struct describes the status of.
	AncestorRef ParentReference `json:"ancestorRef"`
	// ControllerName is a domain/path string that indicates the name of the
	// controller that wrote this status. This corresponds with the
	// controllerName field on GatewayClass.
	// Example: "example.net/gateway-controller".
	// The format of this field is DOMAIN "/" PATH, where DOMAIN and PATH are
	// valid Kubernetes names
	// (https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names).
	// Controllers MUST populate this field when writing status. Controllers should ensure that
	// entries to status populated with their ControllerName are cleaned up when they are no
	// longer necessary.
	ControllerName string `json:"controllerName"`
	// Conditions describes the status of the Policy with respect to the given Ancestor.
	// <gateway:util:excludeFromCRD>
	// Notes for implementors:
	// Conditions are a listType `map`, which means that they function like a
	// map with a key of the `type` field _in the k8s apiserver_.
	// This means that implementations must obey some rules when updating this
	// section.
	// * Implementations MUST perform a read-modify-write cycle on this field
	// before modifying it. That is, when modifying this field, implementations
	// must be confident they have fetched the most recent version of this field,
	// and ensure that changes they make are on that recent version.
	// * Implementations MUST NOT remove or reorder Conditions that they are not
	// directly responsible for. For example, if an implementation sees a Condition
	// with type `special.io/SomeField`, it MUST NOT remove, change or update that
	// Condition.
	// * Implementations MUST always _merge_ changes into Conditions of the same Type,
	// rather than creating more than one Condition of the same Type.
	// * Implementations MUST always update the `observedGeneration` field of the
	// Condition to the `metadata.generation` of the Gateway at the time of update creation.
	// * If the `observedGeneration` of a Condition is _greater than_ the value the
	// implementation knows about, then it MUST NOT perform the update on that Condition,
	// but must wait for a future reconciliation and status update. (The assumption is that
	// the implementation's copy of the object is stale and an update will be re-triggered
	// if relevant.)
	// </gateway:util:excludeFromCRD>
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func (in *PolicyAncestorStatus) DeepCopyInto(out *PolicyAncestorStatus) {
	*out = *in
	in.AncestorRef.DeepCopyInto(&out.AncestorRef)
	if in.Conditions != nil {
		l := make([]metav1.Condition, len(in.Conditions))
		for i := range in.Conditions {
			in.Conditions[i].DeepCopyInto(&l[i])
		}
		out.Conditions = l
	}
}

func (in *PolicyAncestorStatus) DeepCopy() *PolicyAncestorStatus {
	if in == nil {
		return nil
	}
	out := new(PolicyAncestorStatus)
	in.DeepCopyInto(out)
	return out
}

type CommonRouteSpec struct {
	// ParentRefs references the resources (usually Gateways) that a Route wants
	// to be attached to. Note that the referenced parent resource needs to
	// allow this for the attachment to be complete. For Gateways, that means
	// the Gateway needs to allow attachment from Routes of this kind and
	// namespace. For Services, that means the Service must either be in the same
	// namespace for a "producer" route, or the mesh implementation must support
	// and allow "consumer" routes for the referenced Service. ReferenceGrant is
	// not applicable for governing ParentRefs to Services - it is not possible to
	// create a "producer" route for a Service in a different namespace from the
	// Route.
	// There are two kinds of parent resources with "Core" support:
	// * Gateway (Gateway conformance profile)
	// * Service (Mesh conformance profile, ClusterIP Services only)
	// This API may be extended in the future to support additional kinds of parent
	// resources.
	// ParentRefs must be _distinct_. This means either that:
	// * They select different objects.  If this is the case, then parentRef
	// entries are distinct. In terms of fields, this means that the
	// multi-part key defined by `group`, `kind`, `namespace`, and `name` must
	// be unique across all parentRef entries in the Route.
	// * They do not select different objects, but for each optional field used,
	// each ParentRef that selects the same object must set the same set of
	// optional fields to different values. If one ParentRef sets a
	// combination of optional fields, all must set the same combination.
	// Some examples:
	// * If one ParentRef sets `sectionName`, all ParentRefs referencing the
	// same object must also set `sectionName`.
	// * If one ParentRef sets `port`, all ParentRefs referencing the same
	// object must also set `port`.
	// * If one ParentRef sets `sectionName` and `port`, all ParentRefs
	// referencing the same object must also set `sectionName` and `port`.
	// It is possible to separately reference multiple distinct objects that may
	// be collapsed by an implementation. For example, some implementations may
	// choose to merge compatible Gateway Listeners together. If that is the
	// case, the list of routes attached to those resources should also be
	// merged.
	// Note that for ParentRefs that cross namespace boundaries, there are specific
	// rules. Cross-namespace references are only valid if they are explicitly
	// allowed by something in the namespace they are referring to. For example,
	// Gateway has the AllowedRoutes field, and ReferenceGrant provides a
	// generic way to enable other kinds of cross-namespace reference.
	// <gateway:experimental:description>
	// ParentRefs from a Route to a Service in the same namespace are "producer"
	// routes, which apply default routing rules to inbound connections from
	// any namespace to the Service.
	// ParentRefs from a Route to a Service in a different namespace are
	// "consumer" routes, and these routing rules are only applied to outbound
	// connections originating from the same namespace as the Route, for which
	// the intended destination of the connections are a Service targeted as a
	// ParentRef of the Route.
	// </gateway:experimental:description>
	// <gateway:standard:validation:XValidation:message="sectionName must be specified when parentRefs includes 2 or more references to the same parent",rule="self.all(p1, self.all(p2, p1.group == p2.group && p1.kind == p2.kind && p1.name == p2.name && (((!has(p1.__namespace__) || p1.__namespace__ == ”) && (!has(p2.__namespace__) || p2.__namespace__ == ”)) || (has(p1.__namespace__) && has(p2.__namespace__) && p1.__namespace__ == p2.__namespace__ )) ? ((!has(p1.sectionName) || p1.sectionName == ”) == (!has(p2.sectionName) || p2.sectionName == ”)) : true))">
	// <gateway:standard:validation:XValidation:message="sectionName must be unique when parentRefs includes 2 or more references to the same parent",rule="self.all(p1, self.exists_one(p2, p1.group == p2.group && p1.kind == p2.kind && p1.name == p2.name && (((!has(p1.__namespace__) || p1.__namespace__ == ”) && (!has(p2.__namespace__) || p2.__namespace__ == ”)) || (has(p1.__namespace__) && has(p2.__namespace__) && p1.__namespace__ == p2.__namespace__ )) && (((!has(p1.sectionName) || p1.sectionName == ”) && (!has(p2.sectionName) || p2.sectionName == ”)) || (has(p1.sectionName) && has(p2.sectionName) && p1.sectionName == p2.sectionName))))">
	// <gateway:experimental:validation:XValidation:message="sectionName or port must be specified when parentRefs includes 2 or more references to the same parent",rule="self.all(p1, self.all(p2, p1.group == p2.group && p1.kind == p2.kind && p1.name == p2.name && (((!has(p1.__namespace__) || p1.__namespace__ == ”) && (!has(p2.__namespace__) || p2.__namespace__ == ”)) || (has(p1.__namespace__) && has(p2.__namespace__) && p1.__namespace__ == p2.__namespace__)) ? ((!has(p1.sectionName) || p1.sectionName == ”) == (!has(p2.sectionName) || p2.sectionName == ”) && (!has(p1.port) || p1.port == 0) == (!has(p2.port) || p2.port == 0)): true))">
	// <gateway:experimental:validation:XValidation:message="sectionName or port must be unique when parentRefs includes 2 or more references to the same parent",rule="self.all(p1, self.exists_one(p2, p1.group == p2.group && p1.kind == p2.kind && p1.name == p2.name && (((!has(p1.__namespace__) || p1.__namespace__ == ”) && (!has(p2.__namespace__) || p2.__namespace__ == ”)) || (has(p1.__namespace__) && has(p2.__namespace__) && p1.__namespace__ == p2.__namespace__ )) && (((!has(p1.sectionName) || p1.sectionName == ”) && (!has(p2.sectionName) || p2.sectionName == ”)) || ( has(p1.sectionName) && has(p2.sectionName) && p1.sectionName == p2.sectionName)) && (((!has(p1.port) || p1.port == 0) && (!has(p2.port) || p2.port == 0)) || (has(p1.port) && has(p2.port) && p1.port == p2.port))))">
	ParentRefs []ParentReference `json:"parentRefs,omitempty"`
	// UseDefaultGateways indicates the default Gateway scope to use for this
	// Route. If unset (the default) or set to None, the Route will not be
	// attached to any default Gateway; if set, it will be attached to any
	// default Gateway supporting the named scope, subject to the usual rules
	// about which Routes a Gateway is allowed to claim.
	// Think carefully before using this functionality! The set of default
	// Gateways supporting the requested scope can change over time without
	// any notice to the Route author, and in many situations it will not be
	// appropriate to request a default Gateway for a given Route -- for
	// example, a Route with specific security requirements should almost
	// certainly not use a default Gateway.
	// <gateway:experimental>
	UseDefaultGateways GatewayDefaultScope `json:"useDefaultGateways,omitempty"`
}

func (in *CommonRouteSpec) DeepCopyInto(out *CommonRouteSpec) {
	*out = *in
	if in.ParentRefs != nil {
		l := make([]ParentReference, len(in.ParentRefs))
		for i := range in.ParentRefs {
			in.ParentRefs[i].DeepCopyInto(&l[i])
		}
		out.ParentRefs = l
	}
}

func (in *CommonRouteSpec) DeepCopy() *CommonRouteSpec {
	if in == nil {
		return nil
	}
	out := new(CommonRouteSpec)
	in.DeepCopyInto(out)
	return out
}

type GRPCRouteRule struct {
	// Name is the name of the route rule. This name MUST be unique within a Route if it is set.
	// Support: Extended
	Name string `json:"name"`
	// Matches define conditions used for matching the rule against incoming
	// gRPC requests. Each match is independent, i.e. this rule will be matched
	// if **any** one of the matches is satisfied.
	// For example, take the following matches configuration:
	// ```
	// matches:
	// - method:
	// service: foo.bar
	// headers:
	// values:
	// version: 2
	// - method:
	// service: foo.bar.v2
	// ```
	// For a request to match against this rule, it MUST satisfy
	// EITHER of the two conditions:
	// - service of foo.bar AND contains the header `version: 2`
	// - service of foo.bar.v2
	// See the documentation for GRPCRouteMatch on how to specify multiple
	// match conditions to be ANDed together.
	// If no matches are specified, the implementation MUST match every gRPC request.
	// Proxy or Load Balancer routing configuration generated from GRPCRoutes
	// MUST prioritize rules based on the following criteria, continuing on
	// ties. Merging MUST not be done between GRPCRoutes and HTTPRoutes.
	// Precedence MUST be given to the rule with the largest number of:
	// * Characters in a matching non-wildcard hostname.
	// * Characters in a matching hostname.
	// * Characters in a matching service.
	// * Characters in a matching method.
	// * Header matches.
	// If ties still exist across multiple Routes, matching precedence MUST be
	// determined in order of the following criteria, continuing on ties:
	// * The oldest Route based on creation timestamp.
	// * The Route appearing first in alphabetical order by
	// "{namespace}/{name}".
	// If ties still exist within the Route that has been given precedence,
	// matching precedence MUST be granted to the first matching rule meeting
	// the above criteria.
	Matches []GRPCRouteMatch `json:"matches,omitempty"`
	// Filters define the filters that are applied to requests that match
	// this rule.
	// The effects of ordering of multiple behaviors are currently unspecified.
	// This can change in the future based on feedback during the alpha stage.
	// Conformance-levels at this level are defined based on the type of filter:
	// - ALL core filters MUST be supported by all implementations that support
	// GRPCRoute.
	// - Implementers are encouraged to support extended filters.
	// - Implementation-specific custom filters have no API guarantees across
	// implementations.
	// Specifying the same filter multiple times is not supported unless explicitly
	// indicated in the filter.
	// If an implementation cannot support a combination of filters, it must clearly
	// document that limitation. In cases where incompatible or unsupported
	// filters are specified and cause the `Accepted` condition to be set to status
	// `False`, implementations may use the `IncompatibleFilters` reason to specify
	// this configuration error.
	// Support: Core
	Filters []GRPCRouteFilter `json:"filters,omitempty"`
	// BackendRefs defines the backend(s) where matching requests should be
	// sent.
	// Failure behavior here depends on how many BackendRefs are specified and
	// how many are invalid.
	// If *all* entries in BackendRefs are invalid, and there are also no filters
	// specified in this route rule, *all* traffic which matches this rule MUST
	// receive an `UNAVAILABLE` status.
	// See the GRPCBackendRef definition for the rules about what makes a single
	// GRPCBackendRef invalid.
	// When a GRPCBackendRef is invalid, `UNAVAILABLE` statuses MUST be returned for
	// requests that would have otherwise been routed to an invalid backend. If
	// multiple backends are specified, and some are invalid, the proportion of
	// requests that would otherwise have been routed to an invalid backend
	// MUST receive an `UNAVAILABLE` status.
	// For example, if two backends are specified with equal weights, and one is
	// invalid, 50 percent of traffic MUST receive an `UNAVAILABLE` status.
	// Implementations may choose how that 50 percent is determined.
	// Support: Core for Kubernetes Service
	// Support: Implementation-specific for any other resource
	// Support for weight: Core
	BackendRefs []GRPCBackendRef `json:"backendRefs,omitempty"`
	// SessionPersistence defines and configures session persistence
	// for the route rule.
	// Support: Extended
	// <gateway:experimental>
	SessionPersistence *SessionPersistence `json:"sessionPersistence,omitempty"`
}

func (in *GRPCRouteRule) DeepCopyInto(out *GRPCRouteRule) {
	*out = *in
	if in.Matches != nil {
		l := make([]GRPCRouteMatch, len(in.Matches))
		for i := range in.Matches {
			in.Matches[i].DeepCopyInto(&l[i])
		}
		out.Matches = l
	}
	if in.Filters != nil {
		l := make([]GRPCRouteFilter, len(in.Filters))
		for i := range in.Filters {
			in.Filters[i].DeepCopyInto(&l[i])
		}
		out.Filters = l
	}
	if in.BackendRefs != nil {
		l := make([]GRPCBackendRef, len(in.BackendRefs))
		for i := range in.BackendRefs {
			in.BackendRefs[i].DeepCopyInto(&l[i])
		}
		out.BackendRefs = l
	}
	if in.SessionPersistence != nil {
		in, out := &in.SessionPersistence, &out.SessionPersistence
		*out = new(SessionPersistence)
		(*in).DeepCopyInto(*out)
	}
}

func (in *GRPCRouteRule) DeepCopy() *GRPCRouteRule {
	if in == nil {
		return nil
	}
	out := new(GRPCRouteRule)
	in.DeepCopyInto(out)
	return out
}

type RouteStatus struct {
	// Parents is a list of parent resources (usually Gateways) that are
	// associated with the route, and the status of the route with respect to
	// each parent. When this route attaches to a parent, the controller that
	// manages the parent must add an entry to this list when the controller
	// first sees the route and should update the entry as appropriate when the
	// route or gateway is modified.
	// Note that parent references that cannot be resolved by an implementation
	// of this API will not be added to this list. Implementations of this API
	// can only populate Route status for the Gateways/parent resources they are
	// responsible for.
	// A maximum of 32 Gateways will be represented in this list. An empty list
	// means the route has not been attached to any Gateway.
	// <gateway:util:excludeFromCRD>
	// Notes for implementors:
	// While parents is not a listType `map`, this is due to the fact that the
	// list key is not scalar, and Kubernetes is unable to represent this.
	// Parent status MUST be considered to be namespaced by the combination of
	// the parentRef and controllerName fields, and implementations should keep
	// the following rules in mind when updating this status:
	// * Implementations MUST update only entries that have a matching value of
	// `controllerName` for that implementation.
	// * Implementations MUST NOT update entries with non-matching `controllerName`
	// fields.
	// * Implementations MUST treat each `parentRef“ in the Route separately and
	// update its status based on the relationship with that parent.
	// * Implementations MUST perform a read-modify-write cycle on this field
	// before modifying it. That is, when modifying this field, implementations
	// must be confident they have fetched the most recent version of this field,
	// and ensure that changes they make are on that recent version.
	// </gateway:util:excludeFromCRD>
	Parents []RouteParentStatus `json:"parents,omitempty"`
}

func (in *RouteStatus) DeepCopyInto(out *RouteStatus) {
	*out = *in
	if in.Parents != nil {
		l := make([]RouteParentStatus, len(in.Parents))
		for i := range in.Parents {
			in.Parents[i].DeepCopyInto(&l[i])
		}
		out.Parents = l
	}
}

func (in *RouteStatus) DeepCopy() *RouteStatus {
	if in == nil {
		return nil
	}
	out := new(RouteStatus)
	in.DeepCopyInto(out)
	return out
}

type Listener struct {
	// Name is the name of the Listener. This name MUST be unique within a
	// Gateway.
	// Support: Core
	Name string `json:"name"`
	// Hostname specifies the virtual hostname to match for protocol types that
	// define this concept. When unspecified, all hostnames are matched. This
	// field is ignored for protocols that don't require hostname based
	// matching.
	// Implementations MUST apply Hostname matching appropriately for each of
	// the following protocols:
	// * TLS: The Listener Hostname MUST match the SNI.
	// * HTTP: The Listener Hostname MUST match the Host header of the request.
	// * HTTPS: The Listener Hostname SHOULD match both the SNI and Host header.
	// Note that this does not require the SNI and Host header to be the same.
	// The semantics of this are described in more detail below.
	// To ensure security, Section 11.1 of RFC-6066 emphasizes that server
	// implementations that rely on SNI hostname matching MUST also verify
	// hostnames within the application protocol.
	// Section 9.1.2 of RFC-7540 provides a mechanism for servers to reject the
	// reuse of a connection by responding with the HTTP 421 Misdirected Request
	// status code. This indicates that the origin server has rejected the
	// request because it appears to have been misdirected.
	// To detect misdirected requests, Gateways SHOULD match the authority of
	// the requests with all the SNI hostname(s) configured across all the
	// Gateway Listeners on the same port and protocol:
	// * If another Listener has an exact match or more specific wildcard entry,
	// the Gateway SHOULD return a 421.
	// * If the current Listener (selected by SNI matching during ClientHello)
	// does not match the Host:
	// * If another Listener does match the Host, the Gateway SHOULD return a
	// 421.
	// * If no other Listener matches the Host, the Gateway MUST return a
	// 404.
	// For HTTPRoute and TLSRoute resources, there is an interaction with the
	// `spec.hostnames` array. When both listener and route specify hostnames,
	// there MUST be an intersection between the values for a Route to be
	// accepted. For more information, refer to the Route specific Hostnames
	// documentation.
	// Hostnames that are prefixed with a wildcard label (`*.`) are interpreted
	// as a suffix match. That means that a match for `*.example.com` would match
	// both `test.example.com`, and `foo.test.example.com`, but not `example.com`.
	// Support: Core
	Hostname string `json:"hostname"`
	// Port is the network port. Multiple listeners may use the
	// same port, subject to the Listener compatibility rules.
	// Support: Core
	Port int `json:"port"`
	// Protocol specifies the network protocol this listener expects to receive.
	// Support: Core
	Protocol ProtocolType `json:"protocol"`
	// TLS is the TLS configuration for the Listener. This field is required if
	// the Protocol field is "HTTPS" or "TLS". It is invalid to set this field
	// if the Protocol field is "HTTP", "TCP", or "UDP".
	// The association of SNIs to Certificate defined in ListenerTLSConfig is
	// defined based on the Hostname field for this listener.
	// The GatewayClass MUST use the longest matching SNI out of all
	// available certificates for any TLS handshake.
	// Support: Core
	TLS *ListenerTLSConfig `json:"tls,omitempty"`
	// AllowedRoutes defines the types of routes that MAY be attached to a
	// Listener and the trusted namespaces where those Route resources MAY be
	// present.
	// Although a client request may match multiple route rules, only one rule
	// may ultimately receive the request. Matching precedence MUST be
	// determined in order of the following criteria:
	// * The most specific match as defined by the Route type.
	// * The oldest Route based on creation timestamp. For example, a Route with
	// a creation timestamp of "2020-09-08 01:02:03" is given precedence over
	// a Route with a creation timestamp of "2020-09-08 01:02:04".
	// * If everything else is equivalent, the Route appearing first in
	// alphabetical order (namespace/name) should be given precedence. For
	// example, foo/bar is given precedence over foo/baz.
	// All valid rules within a Route attached to this Listener should be
	// implemented. Invalid Route rules can be ignored (sometimes that will mean
	// the full Route). If a Route rule transitions from valid to invalid,
	// support for that Route rule should be dropped to ensure consistency. For
	// example, even if a filter specified by a Route rule is invalid, the rest
	// of the rules within that Route should still be supported.
	// Support: Core
	AllowedRoutes *AllowedRoutes `json:"allowedRoutes,omitempty"`
}

func (in *Listener) DeepCopyInto(out *Listener) {
	*out = *in
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(ListenerTLSConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.AllowedRoutes != nil {
		in, out := &in.AllowedRoutes, &out.AllowedRoutes
		*out = new(AllowedRoutes)
		(*in).DeepCopyInto(*out)
	}
}

func (in *Listener) DeepCopy() *Listener {
	if in == nil {
		return nil
	}
	out := new(Listener)
	in.DeepCopyInto(out)
	return out
}

type GatewaySpecAddress struct {
	// Type of the address.
	Type AddressType `json:"type,omitempty"`
	// When a value is unspecified, an implementation SHOULD automatically
	// assign an address matching the requested type if possible.
	// If an implementation does not support an empty value, they MUST set the
	// "Programmed" condition in status to False with a reason of "AddressNotAssigned".
	// Examples: `1.2.3.4`, `128::1`, `my-ip-address`.
	Value string `json:"value,omitempty"`
}

func (in *GatewaySpecAddress) DeepCopyInto(out *GatewaySpecAddress) {
	*out = *in
}

func (in *GatewaySpecAddress) DeepCopy() *GatewaySpecAddress {
	if in == nil {
		return nil
	}
	out := new(GatewaySpecAddress)
	in.DeepCopyInto(out)
	return out
}

type GatewayInfrastructure struct {
	// Labels that SHOULD be applied to any resources created in response to this Gateway.
	// For implementations creating other Kubernetes objects, this should be the `metadata.labels` field on resources.
	// For other implementations, this refers to any relevant (implementation specific) "labels" concepts.
	// An implementation may chose to add additional implementation-specific labels as they see fit.
	// If an implementation maps these labels to Pods, or any other resource that would need to be recreated when labels
	// change, it SHOULD clearly warn about this behavior in documentation.
	// Support: Extended
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations that SHOULD be applied to any resources created in response to this Gateway.
	// For implementations creating other Kubernetes objects, this should be the `metadata.annotations` field on resources.
	// For other implementations, this refers to any relevant (implementation specific) "annotations" concepts.
	// An implementation may chose to add additional implementation-specific annotations as they see fit.
	// Support: Extended
	Annotations map[string]string `json:"annotations,omitempty"`
	// ParametersRef is a reference to a resource that contains the configuration
	// parameters corresponding to the Gateway. This is optional if the
	// controller does not require any additional configuration.
	// This follows the same semantics as GatewayClass's `parametersRef`, but on a per-Gateway basis
	// The Gateway's GatewayClass may provide its own `parametersRef`. When both are specified,
	// the merging behavior is implementation specific.
	// It is generally recommended that GatewayClass provides defaults that can be overridden by a Gateway.
	// If the referent cannot be found, refers to an unsupported kind, or when
	// the data within that resource is malformed, the Gateway SHOULD be
	// rejected with the "Accepted" status condition set to "False" and an
	// "InvalidParameters" reason.
	// Support: Implementation-specific
	ParametersRef *LocalParametersReference `json:"parametersRef,omitempty"`
}

func (in *GatewayInfrastructure) DeepCopyInto(out *GatewayInfrastructure) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		maps.Copy(*out, *in)
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		maps.Copy(*out, *in)
	}
	if in.ParametersRef != nil {
		in, out := &in.ParametersRef, &out.ParametersRef
		*out = new(LocalParametersReference)
		(*in).DeepCopyInto(*out)
	}
}

func (in *GatewayInfrastructure) DeepCopy() *GatewayInfrastructure {
	if in == nil {
		return nil
	}
	out := new(GatewayInfrastructure)
	in.DeepCopyInto(out)
	return out
}

type AllowedListeners struct {
	// Namespaces defines which namespaces ListenerSets can be attached to this Gateway.
	// The default value is to allow no ListenerSets.
	Namespaces *ListenerNamespaces `json:"namespaces,omitempty"`
}

func (in *AllowedListeners) DeepCopyInto(out *AllowedListeners) {
	*out = *in
	if in.Namespaces != nil {
		in, out := &in.Namespaces, &out.Namespaces
		*out = new(ListenerNamespaces)
		(*in).DeepCopyInto(*out)
	}
}

func (in *AllowedListeners) DeepCopy() *AllowedListeners {
	if in == nil {
		return nil
	}
	out := new(AllowedListeners)
	in.DeepCopyInto(out)
	return out
}

type GatewayTLSConfig struct {
	// Backend describes TLS configuration for gateway when connecting
	// to backends.
	// Note that this contains only details for the Gateway as a TLS client,
	// and does _not_ imply behavior about how to choose which backend should
	// get a TLS connection. That is determined by the presence of a BackendTLSPolicy.
	// Support: Core
	Backend *GatewayBackendTLS `json:"backend,omitempty"`
	// Frontend describes TLS config when client connects to Gateway.
	// Support: Core
	Frontend *FrontendTLSConfig `json:"frontend,omitempty"`
}

func (in *GatewayTLSConfig) DeepCopyInto(out *GatewayTLSConfig) {
	*out = *in
	if in.Backend != nil {
		in, out := &in.Backend, &out.Backend
		*out = new(GatewayBackendTLS)
		(*in).DeepCopyInto(*out)
	}
	if in.Frontend != nil {
		in, out := &in.Frontend, &out.Frontend
		*out = new(FrontendTLSConfig)
		(*in).DeepCopyInto(*out)
	}
}

func (in *GatewayTLSConfig) DeepCopy() *GatewayTLSConfig {
	if in == nil {
		return nil
	}
	out := new(GatewayTLSConfig)
	in.DeepCopyInto(out)
	return out
}

type GatewayStatusAddress struct {
	// Type of the address.
	Type AddressType `json:"type,omitempty"`
	// Value of the address. The validity of the values will depend
	// on the type and support by the controller.
	// Examples: `1.2.3.4`, `128::1`, `my-ip-address`.
	Value string `json:"value"`
}

func (in *GatewayStatusAddress) DeepCopyInto(out *GatewayStatusAddress) {
	*out = *in
}

func (in *GatewayStatusAddress) DeepCopy() *GatewayStatusAddress {
	if in == nil {
		return nil
	}
	out := new(GatewayStatusAddress)
	in.DeepCopyInto(out)
	return out
}

type ListenerStatus struct {
	// Name is the name of the Listener that this status corresponds to.
	Name string `json:"name"`
	// SupportedKinds is the list indicating the Kinds supported by this
	// listener. This MUST represent the kinds supported by an implementation for
	// that Listener configuration.
	// If kinds are specified in Spec that are not supported, they MUST NOT
	// appear in this list and an implementation MUST set the "ResolvedRefs"
	// condition to "False" with the "InvalidRouteKinds" reason. If both valid
	// and invalid Route kinds are specified, the implementation MUST
	// reference the valid Route kinds that have been specified.
	SupportedKinds []RouteGroupKind `json:"supportedKinds,omitempty"`
	// AttachedRoutes represents the total number of Routes that have been
	// successfully attached to this Listener.
	// Successful attachment of a Route to a Listener is based solely on the
	// combination of the AllowedRoutes field on the corresponding Listener
	// and the Route's ParentRefs field. A Route is successfully attached to
	// a Listener when it is selected by the Listener's AllowedRoutes field
	// AND the Route has a valid ParentRef selecting the whole Gateway
	// resource or a specific Listener as a parent resource (more detail on
	// attachment semantics can be found in the documentation on the various
	// Route kinds ParentRefs fields). Listener or Route status does not impact
	// successful attachment, i.e. the AttachedRoutes field count MUST be set
	// for Listeners, even if the Accepted condition of an individual Listener is set
	// to "False". The AttachedRoutes number represents the number of Routes with
	// the Accepted condition set to "True" that have been attached to this Listener.
	// Routes with any other value for the Accepted condition MUST NOT be included
	// in this count.
	// Uses for this field include troubleshooting Route attachment and
	// measuring blast radius/impact of changes to a Listener.
	AttachedRoutes int `json:"attachedRoutes"`
	// Conditions describe the current condition of this listener.
	// <gateway:util:excludeFromCRD>
	// Notes for implementors:
	// Conditions are a listType `map`, which means that they function like a
	// map with a key of the `type` field _in the k8s apiserver_.
	// This means that implementations must obey some rules when updating this
	// section.
	// * Implementations MUST perform a read-modify-write cycle on this field
	// before modifying it. That is, when modifying this field, implementations
	// must be confident they have fetched the most recent version of this field,
	// and ensure that changes they make are on that recent version.
	// * Implementations MUST NOT remove or reorder Conditions that they are not
	// directly responsible for. For example, if an implementation sees a Condition
	// with type `special.io/SomeField`, it MUST NOT remove, change or update that
	// Condition.
	// * Implementations MUST always _merge_ changes into Conditions of the same Type,
	// rather than creating more than one Condition of the same Type.
	// * Implementations MUST always update the `observedGeneration` field of the
	// Condition to the `metadata.generation` of the Gateway at the time of update creation.
	// * If the `observedGeneration` of a Condition is _greater than_ the value the
	// implementation knows about, then it MUST NOT perform the update on that Condition,
	// but must wait for a future reconciliation and status update. (The assumption is that
	// the implementation's copy of the object is stale and an update will be re-triggered
	// if relevant.)
	// </gateway:util:excludeFromCRD>
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func (in *ListenerStatus) DeepCopyInto(out *ListenerStatus) {
	*out = *in
	if in.SupportedKinds != nil {
		l := make([]RouteGroupKind, len(in.SupportedKinds))
		for i := range in.SupportedKinds {
			in.SupportedKinds[i].DeepCopyInto(&l[i])
		}
		out.SupportedKinds = l
	}
	if in.Conditions != nil {
		l := make([]metav1.Condition, len(in.Conditions))
		for i := range in.Conditions {
			in.Conditions[i].DeepCopyInto(&l[i])
		}
		out.Conditions = l
	}
}

func (in *ListenerStatus) DeepCopy() *ListenerStatus {
	if in == nil {
		return nil
	}
	out := new(ListenerStatus)
	in.DeepCopyInto(out)
	return out
}

type ParametersReference struct {
	// Group is the group of the referent.
	Group string `json:"group"`
	// Kind is kind of the referent.
	Kind string `json:"kind"`
	// Name is the name of the referent.
	Name string `json:"name"`
	// Namespace is the namespace of the referent.
	// This field is required when referring to a Namespace-scoped resource and
	// MUST be unset when referring to a Cluster-scoped resource.
	Namespace string `json:"namespace"`
}

func (in *ParametersReference) DeepCopyInto(out *ParametersReference) {
	*out = *in
}

func (in *ParametersReference) DeepCopy() *ParametersReference {
	if in == nil {
		return nil
	}
	out := new(ParametersReference)
	in.DeepCopyInto(out)
	return out
}

type SupportedFeature struct {
	Name string `json:"name"`
}

func (in *SupportedFeature) DeepCopyInto(out *SupportedFeature) {
	*out = *in
}

func (in *SupportedFeature) DeepCopy() *SupportedFeature {
	if in == nil {
		return nil
	}
	out := new(SupportedFeature)
	in.DeepCopyInto(out)
	return out
}

type HTTPRouteRule struct {
	// Name is the name of the route rule. This name MUST be unique within a Route if it is set.
	// Support: Extended
	Name string `json:"name"`
	// Matches define conditions used for matching the rule against incoming
	// HTTP requests. Each match is independent, i.e. this rule will be matched
	// if **any** one of the matches is satisfied.
	// For example, take the following matches configuration:
	// ```
	// matches:
	// - path:
	// value: "/foo"
	// headers:
	// - name: "version"
	// value: "v2"
	// - path:
	// value: "/v2/foo"
	// ```
	// For a request to match against this rule, a request must satisfy
	// EITHER of the two conditions:
	// - path prefixed with `/foo` AND contains the header `version: v2`
	// - path prefix of `/v2/foo`
	// See the documentation for HTTPRouteMatch on how to specify multiple
	// match conditions that should be ANDed together.
	// If no matches are specified, the default is a prefix
	// path match on "/", which has the effect of matching every
	// HTTP request.
	// Proxy or Load Balancer routing configuration generated from HTTPRoutes
	// MUST prioritize matches based on the following criteria, continuing on
	// ties. Across all rules specified on applicable Routes, precedence must be
	// given to the match having:
	// * "Exact" path match.
	// * "Prefix" path match with largest number of characters.
	// * Method match.
	// * Largest number of header matches.
	// * Largest number of query param matches.
	// Note: The precedence of RegularExpression path matches are implementation-specific.
	// If ties still exist across multiple Routes, matching precedence MUST be
	// determined in order of the following criteria, continuing on ties:
	// * The oldest Route based on creation timestamp.
	// * The Route appearing first in alphabetical order by
	// "{namespace}/{name}".
	// If ties still exist within an HTTPRoute, matching precedence MUST be granted
	// to the FIRST matching rule (in list order) with a match meeting the above
	// criteria.
	// When no rules matching a request have been successfully attached to the
	// parent a request is coming from, a HTTP 404 status code MUST be returned.
	Matches []HTTPRouteMatch `json:"matches,omitempty"`
	// Filters define the filters that are applied to requests that match
	// this rule.
	// Wherever possible, implementations SHOULD implement filters in the order
	// they are specified.
	// Implementations MAY choose to implement this ordering strictly, rejecting
	// any combination or order of filters that cannot be supported. If implementations
	// choose a strict interpretation of filter ordering, they MUST clearly document
	// that behavior.
	// To reject an invalid combination or order of filters, implementations SHOULD
	// consider the Route Rules with this configuration invalid. If all Route Rules
	// in a Route are invalid, the entire Route would be considered invalid. If only
	// a portion of Route Rules are invalid, implementations MUST set the
	// "PartiallyInvalid" condition for the Route.
	// Conformance-levels at this level are defined based on the type of filter:
	// - ALL core filters MUST be supported by all implementations.
	// - Implementers are encouraged to support extended filters.
	// - Implementation-specific custom filters have no API guarantees across
	// implementations.
	// Specifying the same filter multiple times is not supported unless explicitly
	// indicated in the filter.
	// All filters are expected to be compatible with each other except for the
	// URLRewrite and RequestRedirect filters, which may not be combined. If an
	// implementation cannot support other combinations of filters, they must clearly
	// document that limitation. In cases where incompatible or unsupported
	// filters are specified and cause the `Accepted` condition to be set to status
	// `False`, implementations may use the `IncompatibleFilters` reason to specify
	// this configuration error.
	// Support: Core
	Filters []HTTPRouteFilter `json:"filters,omitempty"`
	// BackendRefs defines the backend(s) where matching requests should be
	// sent.
	// Failure behavior here depends on how many BackendRefs are specified and
	// how many are invalid.
	// If *all* entries in BackendRefs are invalid, and there are also no filters
	// specified in this route rule, *all* traffic which matches this rule MUST
	// receive a 500 status code.
	// See the HTTPBackendRef definition for the rules about what makes a single
	// HTTPBackendRef invalid.
	// When a HTTPBackendRef is invalid, 500 status codes MUST be returned for
	// requests that would have otherwise been routed to an invalid backend. If
	// multiple backends are specified, and some are invalid, the proportion of
	// requests that would otherwise have been routed to an invalid backend
	// MUST receive a 500 status code.
	// For example, if two backends are specified with equal weights, and one is
	// invalid, 50 percent of traffic must receive a 500. Implementations may
	// choose how that 50 percent is determined.
	// When a HTTPBackendRef refers to a Service that has no ready endpoints,
	// implementations SHOULD return a 503 for requests to that backend instead.
	// If an implementation chooses to do this, all of the above rules for 500 responses
	// MUST also apply for responses that return a 503.
	// Support: Core for Kubernetes Service
	// Support: Extended for Kubernetes ServiceImport
	// Support: Implementation-specific for any other resource
	// Support for weight: Core
	BackendRefs []HTTPBackendRef `json:"backendRefs,omitempty"`
	// Timeouts defines the timeouts that can be configured for an HTTP request.
	// Support: Extended
	Timeouts *HTTPRouteTimeouts `json:"timeouts,omitempty"`
	// Retry defines the configuration for when to retry an HTTP request.
	// Support: Extended
	// <gateway:experimental>
	Retry *HTTPRouteRetry `json:"retry,omitempty"`
	// SessionPersistence defines and configures session persistence
	// for the route rule.
	// Support: Extended
	// <gateway:experimental>
	SessionPersistence *SessionPersistence `json:"sessionPersistence,omitempty"`
}

func (in *HTTPRouteRule) DeepCopyInto(out *HTTPRouteRule) {
	*out = *in
	if in.Matches != nil {
		l := make([]HTTPRouteMatch, len(in.Matches))
		for i := range in.Matches {
			in.Matches[i].DeepCopyInto(&l[i])
		}
		out.Matches = l
	}
	if in.Filters != nil {
		l := make([]HTTPRouteFilter, len(in.Filters))
		for i := range in.Filters {
			in.Filters[i].DeepCopyInto(&l[i])
		}
		out.Filters = l
	}
	if in.BackendRefs != nil {
		l := make([]HTTPBackendRef, len(in.BackendRefs))
		for i := range in.BackendRefs {
			in.BackendRefs[i].DeepCopyInto(&l[i])
		}
		out.BackendRefs = l
	}
	if in.Timeouts != nil {
		in, out := &in.Timeouts, &out.Timeouts
		*out = new(HTTPRouteTimeouts)
		(*in).DeepCopyInto(*out)
	}
	if in.Retry != nil {
		in, out := &in.Retry, &out.Retry
		*out = new(HTTPRouteRetry)
		(*in).DeepCopyInto(*out)
	}
	if in.SessionPersistence != nil {
		in, out := &in.SessionPersistence, &out.SessionPersistence
		*out = new(SessionPersistence)
		(*in).DeepCopyInto(*out)
	}
}

func (in *HTTPRouteRule) DeepCopy() *HTTPRouteRule {
	if in == nil {
		return nil
	}
	out := new(HTTPRouteRule)
	in.DeepCopyInto(out)
	return out
}

type ParentGatewayReference struct {
	// Group is the group of the referent.
	Group string `json:"group"`
	// Kind is kind of the referent. For example "Gateway".
	Kind string `json:"kind"`
	// Name is the name of the referent.
	Name string `json:"name"`
	// Namespace is the namespace of the referent.  If not present,
	// the namespace of the referent is assumed to be the same as
	// the namespace of the referring object.
	Namespace string `json:"namespace"`
}

func (in *ParentGatewayReference) DeepCopyInto(out *ParentGatewayReference) {
	*out = *in
}

func (in *ParentGatewayReference) DeepCopy() *ParentGatewayReference {
	if in == nil {
		return nil
	}
	out := new(ParentGatewayReference)
	in.DeepCopyInto(out)
	return out
}

type ListenerEntry struct {
	// Name is the name of the Listener. This name MUST be unique within a
	// ListenerSet.
	// Name is not required to be unique across a Gateway and ListenerSets.
	// Routes can attach to a Listener by having a ListenerSet as a parentRef
	// and setting the SectionName
	Name string `json:"name"`
	// Hostname specifies the virtual hostname to match for protocol types that
	// define this concept. When unspecified, all hostnames are matched. This
	// field is ignored for protocols that don't require hostname based
	// matching.
	// Implementations MUST apply Hostname matching appropriately for each of
	// the following protocols:
	// * TLS: The Listener Hostname MUST match the SNI.
	// * HTTP: The Listener Hostname MUST match the Host header of the request.
	// * HTTPS: The Listener Hostname SHOULD match at both the TLS and HTTP
	// protocol layers as described above. If an implementation does not
	// ensure that both the SNI and Host header match the Listener hostname,
	// it MUST clearly document that.
	// For HTTPRoute and TLSRoute resources, there is an interaction with the
	// `spec.hostnames` array. When both listener and route specify hostnames,
	// there MUST be an intersection between the values for a Route to be
	// accepted. For more information, refer to the Route specific Hostnames
	// documentation.
	// Hostnames that are prefixed with a wildcard label (`*.`) are interpreted
	// as a suffix match. That means that a match for `*.example.com` would match
	// both `test.example.com`, and `foo.test.example.com`, but not `example.com`.
	Hostname string `json:"hostname"`
	// Port is the network port. Multiple listeners may use the
	// same port, subject to the Listener compatibility rules.
	Port int `json:"port"`
	// Protocol specifies the network protocol this listener expects to receive.
	Protocol ProtocolType `json:"protocol,omitempty"`
	// TLS is the TLS configuration for the Listener. This field is required if
	// the Protocol field is "HTTPS" or "TLS". It is invalid to set this field
	// if the Protocol field is "HTTP", "TCP", or "UDP".
	// The association of SNIs to Certificate defined in ListenerTLSConfig is
	// defined based on the Hostname field for this listener.
	// The GatewayClass MUST use the longest matching SNI out of all
	// available certificates for any TLS handshake.
	TLS *ListenerTLSConfig `json:"tls,omitempty"`
	// AllowedRoutes defines the types of routes that MAY be attached to a
	// Listener and the trusted namespaces where those Route resources MAY be
	// present.
	// Although a client request may match multiple route rules, only one rule
	// may ultimately receive the request. Matching precedence MUST be
	// determined in order of the following criteria:
	// * The most specific match as defined by the Route type.
	// * The oldest Route based on creation timestamp. For example, a Route with
	// a creation timestamp of "2020-09-08 01:02:03" is given precedence over
	// a Route with a creation timestamp of "2020-09-08 01:02:04".
	// * If everything else is equivalent, the Route appearing first in
	// alphabetical order (namespace/name) should be given precedence. For
	// example, foo/bar is given precedence over foo/baz.
	// All valid rules within a Route attached to this Listener should be
	// implemented. Invalid Route rules can be ignored (sometimes that will mean
	// the full Route). If a Route rule transitions from valid to invalid,
	// support for that Route rule should be dropped to ensure consistency. For
	// example, even if a filter specified by a Route rule is invalid, the rest
	// of the rules within that Route should still be supported.
	AllowedRoutes *AllowedRoutes `json:"allowedRoutes,omitempty"`
}

func (in *ListenerEntry) DeepCopyInto(out *ListenerEntry) {
	*out = *in
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(ListenerTLSConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.AllowedRoutes != nil {
		in, out := &in.AllowedRoutes, &out.AllowedRoutes
		*out = new(AllowedRoutes)
		(*in).DeepCopyInto(*out)
	}
}

func (in *ListenerEntry) DeepCopy() *ListenerEntry {
	if in == nil {
		return nil
	}
	out := new(ListenerEntry)
	in.DeepCopyInto(out)
	return out
}

type ListenerEntryStatus struct {
	// Name is the name of the Listener that this status corresponds to.
	Name string `json:"name"`
	// SupportedKinds is the list indicating the Kinds supported by this
	// listener. This MUST represent the kinds supported by an implementation for
	// that Listener configuration.
	// If kinds are specified in Spec that are not supported, they MUST NOT
	// appear in this list and an implementation MUST set the "ResolvedRefs"
	// condition to "False" with the "InvalidRouteKinds" reason. If both valid
	// and invalid Route kinds are specified, the implementation MUST
	// reference the valid Route kinds that have been specified.
	SupportedKinds []RouteGroupKind `json:"supportedKinds,omitempty"`
	// AttachedRoutes represents the total number of Routes that have been
	// successfully attached to this Listener.
	// Successful attachment of a Route to a Listener is based solely on the
	// combination of the AllowedRoutes field on the corresponding Listener
	// and the Route's ParentRefs field. A Route is successfully attached to
	// a Listener when it is selected by the Listener's AllowedRoutes field
	// AND the Route has a valid ParentRef selecting the whole Gateway
	// resource or a specific Listener as a parent resource (more detail on
	// attachment semantics can be found in the documentation on the various
	// Route kinds ParentRefs fields). Listener status does not impact
	// successful attachment, i.e. the AttachedRoutes field count MUST be set
	// for Listeners, even if the Accepted condition of an individual Listener is set
	// to "False". The AttachedRoutes number represents the number of Routes with
	// the Accepted condition set to "True" that have been attached to this Listener.
	// Routes with any other value for the Accepted condition MUST NOT be included
	// in this count.
	// Uses for this field include troubleshooting Route attachment and
	// measuring blast radius/impact of changes to a Listener.
	AttachedRoutes int `json:"attachedRoutes"`
	// Conditions describe the current condition of this listener.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func (in *ListenerEntryStatus) DeepCopyInto(out *ListenerEntryStatus) {
	*out = *in
	if in.SupportedKinds != nil {
		l := make([]RouteGroupKind, len(in.SupportedKinds))
		for i := range in.SupportedKinds {
			in.SupportedKinds[i].DeepCopyInto(&l[i])
		}
		out.SupportedKinds = l
	}
	if in.Conditions != nil {
		l := make([]metav1.Condition, len(in.Conditions))
		for i := range in.Conditions {
			in.Conditions[i].DeepCopyInto(&l[i])
		}
		out.Conditions = l
	}
}

func (in *ListenerEntryStatus) DeepCopy() *ListenerEntryStatus {
	if in == nil {
		return nil
	}
	out := new(ListenerEntryStatus)
	in.DeepCopyInto(out)
	return out
}

type ReferenceGrantFrom struct {
	// Group is the group of the referent.
	// When empty, the Kubernetes core API group is inferred.
	// Support: Core
	Group string `json:"group"`
	// Kind is the kind of the referent. Although implementations may support
	// additional resources, the following types are part of the "Core"
	// support level for this field.
	// When used to permit a SecretObjectReference:
	// * Gateway
	// When used to permit a BackendObjectReference:
	// * GRPCRoute
	// * HTTPRoute
	// * TCPRoute
	// * TLSRoute
	// * UDPRoute
	Kind string `json:"kind"`
	// Namespace is the namespace of the referent.
	// Support: Core
	Namespace string `json:"namespace"`
}

func (in *ReferenceGrantFrom) DeepCopyInto(out *ReferenceGrantFrom) {
	*out = *in
}

func (in *ReferenceGrantFrom) DeepCopy() *ReferenceGrantFrom {
	if in == nil {
		return nil
	}
	out := new(ReferenceGrantFrom)
	in.DeepCopyInto(out)
	return out
}

type ReferenceGrantTo struct {
	// Group is the group of the referent.
	// When empty, the Kubernetes core API group is inferred.
	// Support: Core
	Group string `json:"group"`
	// Kind is the kind of the referent. Although implementations may support
	// additional resources, the following types are part of the "Core"
	// support level for this field:
	// * Secret when used to permit a SecretObjectReference
	// * Service when used to permit a BackendObjectReference
	Kind string `json:"kind"`
	// Name is the name of the referent. When unspecified, this policy
	// refers to all resources of the specified Group and Kind in the local
	// namespace.
	Name string `json:"name"`
}

func (in *ReferenceGrantTo) DeepCopyInto(out *ReferenceGrantTo) {
	*out = *in
}

func (in *ReferenceGrantTo) DeepCopy() *ReferenceGrantTo {
	if in == nil {
		return nil
	}
	out := new(ReferenceGrantTo)
	in.DeepCopyInto(out)
	return out
}

type TLSRouteRule struct {
	// Name is the name of the route rule. This name MUST be unique within a Route if it is set.
	Name string `json:"name"`
	// BackendRefs defines the backend(s) where matching requests should be
	// sent. If unspecified or invalid (refers to a nonexistent resource or
	// a Service with no endpoints), the rule performs no forwarding; if no
	// filters are specified that would result in a response being sent, the
	// underlying implementation must actively reject request attempts to this
	// backend, by rejecting the connection. Request rejections must respect
	// weight; if an invalid backend is requested to have 80% of requests, then
	// 80% of requests must be rejected instead.
	// Support: Core for Kubernetes Service
	// Support: Extended for Kubernetes ServiceImport
	// Support: Implementation-specific for any other resource
	// Support for weight: Extended
	BackendRefs []BackendRef `json:"backendRefs,omitempty"`
}

func (in *TLSRouteRule) DeepCopyInto(out *TLSRouteRule) {
	*out = *in
	if in.BackendRefs != nil {
		l := make([]BackendRef, len(in.BackendRefs))
		for i := range in.BackendRefs {
			in.BackendRefs[i].DeepCopyInto(&l[i])
		}
		out.BackendRefs = l
	}
}

func (in *TLSRouteRule) DeepCopy() *TLSRouteRule {
	if in == nil {
		return nil
	}
	out := new(TLSRouteRule)
	in.DeepCopyInto(out)
	return out
}

type LocalPolicyTargetReference struct {
	// Group is the group of the target resource.
	Group string `json:"group"`
	// Kind is kind of the target resource.
	Kind string `json:"kind"`
	// Name is the name of the target resource.
	Name string `json:"name"`
}

func (in *LocalPolicyTargetReference) DeepCopyInto(out *LocalPolicyTargetReference) {
	*out = *in
}

func (in *LocalPolicyTargetReference) DeepCopy() *LocalPolicyTargetReference {
	if in == nil {
		return nil
	}
	out := new(LocalPolicyTargetReference)
	in.DeepCopyInto(out)
	return out
}

type LocalObjectReference struct {
	// Group is the group of the referent. For example, "gateway.networking.k8s.io".
	// When unspecified or empty string, core API group is inferred.
	Group string `json:"group"`
	// Kind is kind of the referent. For example "HTTPRoute" or "Service".
	Kind string `json:"kind"`
	// Name is the name of the referent.
	Name string `json:"name"`
}

func (in *LocalObjectReference) DeepCopyInto(out *LocalObjectReference) {
	*out = *in
}

func (in *LocalObjectReference) DeepCopy() *LocalObjectReference {
	if in == nil {
		return nil
	}
	out := new(LocalObjectReference)
	in.DeepCopyInto(out)
	return out
}

type SubjectAltName struct {
	// Type determines the format of the Subject Alternative Name. Always required.
	// Support: Core
	Type SubjectAltNameType `json:"type"`
	// Hostname contains Subject Alternative Name specified in DNS name format.
	// Required when Type is set to Hostname, ignored otherwise.
	// Support: Core
	Hostname string `json:"hostname"`
	// URI contains Subject Alternative Name specified in a full URI format.
	// It MUST include both a scheme (e.g., "http" or "ftp") and a scheme-specific-part.
	// Common values include SPIFFE IDs like "spiffe://mycluster.example.com/ns/myns/sa/svc1sa".
	// Required when Type is set to URI, ignored otherwise.
	// Support: Core
	URI string `json:"uri"`
}

func (in *SubjectAltName) DeepCopyInto(out *SubjectAltName) {
	*out = *in
}

func (in *SubjectAltName) DeepCopy() *SubjectAltName {
	if in == nil {
		return nil
	}
	out := new(SubjectAltName)
	in.DeepCopyInto(out)
	return out
}

type ParentReference struct {
	// Group is the group of the referent.
	// When unspecified, "gateway.networking.k8s.io" is inferred.
	// To set the core API group (such as for a "Service" kind referent),
	// Group must be explicitly set to "" (empty string).
	// Support: Core
	Group string `json:"group"`
	// Kind is kind of the referent.
	// There are two kinds of parent resources with "Core" support:
	// * Gateway (Gateway conformance profile)
	// * Service (Mesh conformance profile, ClusterIP Services only)
	// Support for other resources is Implementation-Specific.
	Kind string `json:"kind"`
	// Namespace is the namespace of the referent. When unspecified, this refers
	// to the local namespace of the Route.
	// Note that there are specific rules for ParentRefs which cross namespace
	// boundaries. Cross-namespace references are only valid if they are explicitly
	// allowed by something in the namespace they are referring to. For example:
	// Gateway has the AllowedRoutes field, and ReferenceGrant provides a
	// generic way to enable any other kind of cross-namespace reference.
	// <gateway:experimental:description>
	// ParentRefs from a Route to a Service in the same namespace are "producer"
	// routes, which apply default routing rules to inbound connections from
	// any namespace to the Service.
	// ParentRefs from a Route to a Service in a different namespace are
	// "consumer" routes, and these routing rules are only applied to outbound
	// connections originating from the same namespace as the Route, for which
	// the intended destination of the connections are a Service targeted as a
	// ParentRef of the Route.
	// </gateway:experimental:description>
	// Support: Core
	Namespace string `json:"namespace"`
	// Name is the name of the referent.
	// Support: Core
	Name string `json:"name"`
	// SectionName is the name of a section within the target resource. In the
	// following resources, SectionName is interpreted as the following:
	// * Gateway: Listener name. When both Port (experimental) and SectionName
	// are specified, the name and port of the selected listener must match
	// both specified values.
	// * Service: Port name. When both Port (experimental) and SectionName
	// are specified, the name and port of the selected listener must match
	// both specified values.
	// Implementations MAY choose to support attaching Routes to other resources.
	// If that is the case, they MUST clearly document how SectionName is
	// interpreted.
	// When unspecified (empty string), this will reference the entire resource.
	// For the purpose of status, an attachment is considered successful if at
	// least one section in the parent resource accepts it. For example, Gateway
	// listeners can restrict which Routes can attach to them by Route kind,
	// namespace, or hostname. If 1 of 2 Gateway listeners accept attachment from
	// the referencing Route, the Route MUST be considered successfully
	// attached. If no Gateway listeners accept attachment from this Route, the
	// Route MUST be considered detached from the Gateway.
	// Support: Core
	SectionName string `json:"sectionName"`
	// Port is the network port this Route targets. It can be interpreted
	// differently based on the type of parent resource.
	// When the parent resource is a Gateway, this targets all listeners
	// listening on the specified port that also support this kind of Route(and
	// select this Route). It's not recommended to set `Port` unless the
	// networking behaviors specified in a Route must apply to a specific port
	// as opposed to a listener(s) whose port(s) may be changed. When both Port
	// and SectionName are specified, the name and port of the selected listener
	// must match both specified values.
	// <gateway:experimental:description>
	// When the parent resource is a Service, this targets a specific port in the
	// Service spec. When both Port (experimental) and SectionName are specified,
	// the name and port of the selected port must match both specified values.
	// </gateway:experimental:description>
	// Implementations MAY choose to support other parent resources.
	// Implementations supporting other types of parent resources MUST clearly
	// document how/if Port is interpreted.
	// For the purpose of status, an attachment is considered successful as
	// long as the parent resource accepts it partially. For example, Gateway
	// listeners can restrict which Routes can attach to them by Route kind,
	// namespace, or hostname. If 1 of 2 Gateway listeners accept attachment
	// from the referencing Route, the Route MUST be considered successfully
	// attached. If no Gateway listeners accept attachment from this Route,
	// the Route MUST be considered detached from the Gateway.
	// Support: Extended
	Port int `json:"port"`
}

func (in *ParentReference) DeepCopyInto(out *ParentReference) {
	*out = *in
}

func (in *ParentReference) DeepCopy() *ParentReference {
	if in == nil {
		return nil
	}
	out := new(ParentReference)
	in.DeepCopyInto(out)
	return out
}

type GRPCRouteMatch struct {
	// Method specifies a gRPC request service/method matcher. If this field is
	// not specified, all services and methods will match.
	Method *GRPCMethodMatch `json:"method,omitempty"`
	// Headers specifies gRPC request header matchers. Multiple match values are
	// ANDed together, meaning, a request MUST match all the specified headers
	// to select the route.
	Headers []GRPCHeaderMatch `json:"headers,omitempty"`
}

func (in *GRPCRouteMatch) DeepCopyInto(out *GRPCRouteMatch) {
	*out = *in
	if in.Method != nil {
		in, out := &in.Method, &out.Method
		*out = new(GRPCMethodMatch)
		(*in).DeepCopyInto(*out)
	}
	if in.Headers != nil {
		l := make([]GRPCHeaderMatch, len(in.Headers))
		for i := range in.Headers {
			in.Headers[i].DeepCopyInto(&l[i])
		}
		out.Headers = l
	}
}

func (in *GRPCRouteMatch) DeepCopy() *GRPCRouteMatch {
	if in == nil {
		return nil
	}
	out := new(GRPCRouteMatch)
	in.DeepCopyInto(out)
	return out
}

type GRPCRouteFilter struct {
	// Type identifies the type of filter to apply. As with other API fields,
	// types are classified into three conformance levels:
	// - Core: Filter types and their corresponding configuration defined by
	// "Support: Core" in this package, e.g. "RequestHeaderModifier". All
	// implementations supporting GRPCRoute MUST support core filters.
	// - Extended: Filter types and their corresponding configuration defined by
	// "Support: Extended" in this package, e.g. "RequestMirror". Implementers
	// are encouraged to support extended filters.
	// - Implementation-specific: Filters that are defined and supported by specific vendors.
	// In the future, filters showing convergence in behavior across multiple
	// implementations will be considered for inclusion in extended or core
	// conformance levels. Filter-specific configuration for such filters
	// is specified using the ExtensionRef field. `Type` MUST be set to
	// "ExtensionRef" for custom filters.
	// Implementers are encouraged to define custom implementation types to
	// extend the core API with implementation-specific behavior.
	// If a reference to a custom filter type cannot be resolved, the filter
	// MUST NOT be skipped. Instead, requests that would have been processed by
	// that filter MUST receive a HTTP error response.
	// <gateway:experimental:validation:Enum=ResponseHeaderModifier;RequestHeaderModifier;RequestMirror;ExtensionRef>
	Type GRPCRouteFilterType `json:"type"`
	// RequestHeaderModifier defines a schema for a filter that modifies request
	// headers.
	// Support: Core
	RequestHeaderModifier *HTTPHeaderFilter `json:"requestHeaderModifier,omitempty"`
	// ResponseHeaderModifier defines a schema for a filter that modifies response
	// headers.
	// Support: Extended
	ResponseHeaderModifier *HTTPHeaderFilter `json:"responseHeaderModifier,omitempty"`
	// RequestMirror defines a schema for a filter that mirrors requests.
	// Requests are sent to the specified destination, but responses from
	// that destination are ignored.
	// This filter can be used multiple times within the same rule. Note that
	// not all implementations will be able to support mirroring to multiple
	// backends.
	// Support: Extended
	RequestMirror *HTTPRequestMirrorFilter `json:"requestMirror,omitempty"`
	// ExtensionRef is an optional, implementation-specific extension to the
	// "filter" behavior.  For example, resource "myroutefilter" in group
	// "networking.example.net"). ExtensionRef MUST NOT be used for core and
	// extended filters.
	// Support: Implementation-specific
	// This filter can be used multiple times within the same rule.
	ExtensionRef *LocalObjectReference `json:"extensionRef,omitempty"`
}

func (in *GRPCRouteFilter) DeepCopyInto(out *GRPCRouteFilter) {
	*out = *in
	if in.RequestHeaderModifier != nil {
		in, out := &in.RequestHeaderModifier, &out.RequestHeaderModifier
		*out = new(HTTPHeaderFilter)
		(*in).DeepCopyInto(*out)
	}
	if in.ResponseHeaderModifier != nil {
		in, out := &in.ResponseHeaderModifier, &out.ResponseHeaderModifier
		*out = new(HTTPHeaderFilter)
		(*in).DeepCopyInto(*out)
	}
	if in.RequestMirror != nil {
		in, out := &in.RequestMirror, &out.RequestMirror
		*out = new(HTTPRequestMirrorFilter)
		(*in).DeepCopyInto(*out)
	}
	if in.ExtensionRef != nil {
		in, out := &in.ExtensionRef, &out.ExtensionRef
		*out = new(LocalObjectReference)
		(*in).DeepCopyInto(*out)
	}
}

func (in *GRPCRouteFilter) DeepCopy() *GRPCRouteFilter {
	if in == nil {
		return nil
	}
	out := new(GRPCRouteFilter)
	in.DeepCopyInto(out)
	return out
}

type GRPCBackendRef struct {
	// BackendRef is a reference to a backend to forward matched requests to.
	// A BackendRef can be invalid for the following reasons. In all cases, the
	// implementation MUST ensure the `ResolvedRefs` Condition on the Route
	// is set to `status: False`, with a Reason and Message that indicate
	// what is the cause of the error.
	// A BackendRef is invalid if:
	// * It refers to an unknown or unsupported kind of resource. In this
	// case, the Reason MUST be set to `InvalidKind` and Message of the
	// Condition MUST explain which kind of resource is unknown or unsupported.
	// * It refers to a resource that does not exist. In this case, the Reason MUST
	// be set to `BackendNotFound` and the Message of the Condition MUST explain
	// which resource does not exist.
	// * It refers a resource in another namespace when the reference has not been
	// explicitly allowed by a ReferenceGrant (or equivalent concept). In this
	// case, the Reason MUST be set to `RefNotPermitted` and the Message of the
	// Condition MUST explain which cross-namespace reference is not allowed.
	// Support: Core for Kubernetes Service
	// Support: Extended for Kubernetes ServiceImport
	// Support: Implementation-specific for any other resource
	// Support for weight: Core
	BackendRef `json:",inline"`
	// Filters defined at this level MUST be executed if and only if the
	// request is being forwarded to the backend defined here.
	// Support: Implementation-specific (For broader support of filters, use the
	// Filters field in GRPCRouteRule.)
	Filters []GRPCRouteFilter `json:"filters,omitempty"`
}

func (in *GRPCBackendRef) DeepCopyInto(out *GRPCBackendRef) {
	*out = *in
	out.BackendRef = in.BackendRef
	if in.Filters != nil {
		l := make([]GRPCRouteFilter, len(in.Filters))
		for i := range in.Filters {
			in.Filters[i].DeepCopyInto(&l[i])
		}
		out.Filters = l
	}
}

func (in *GRPCBackendRef) DeepCopy() *GRPCBackendRef {
	if in == nil {
		return nil
	}
	out := new(GRPCBackendRef)
	in.DeepCopyInto(out)
	return out
}

type SessionPersistence struct {
	// SessionName defines the name of the persistent session token
	// which may be reflected in the cookie or the header. Users
	// should avoid reusing session names to prevent unintended
	// consequences, such as rejection or unpredictable behavior.
	// Support: Implementation-specific
	SessionName string `json:"sessionName,omitempty"`
	// AbsoluteTimeout defines the absolute timeout of the persistent
	// session. Once the AbsoluteTimeout duration has elapsed, the
	// session becomes invalid.
	// Support: Extended
	AbsoluteTimeout string `json:"absoluteTimeout"`
	// IdleTimeout defines the idle timeout of the persistent session.
	// Once the session has been idle for more than the specified
	// IdleTimeout duration, the session becomes invalid.
	// Support: Extended
	IdleTimeout string `json:"idleTimeout"`
	// Type defines the type of session persistence such as through
	// the use of a header or cookie. Defaults to cookie based session
	// persistence.
	// Support: Core for "Cookie" type
	// Support: Extended for "Header" type
	Type SessionPersistenceType `json:"type,omitempty"`
	// CookieConfig provides configuration settings that are specific
	// to cookie-based session persistence.
	// Support: Core
	CookieConfig *CookieConfig `json:"cookieConfig,omitempty"`
}

func (in *SessionPersistence) DeepCopyInto(out *SessionPersistence) {
	*out = *in
	if in.CookieConfig != nil {
		in, out := &in.CookieConfig, &out.CookieConfig
		*out = new(CookieConfig)
		(*in).DeepCopyInto(*out)
	}
}

func (in *SessionPersistence) DeepCopy() *SessionPersistence {
	if in == nil {
		return nil
	}
	out := new(SessionPersistence)
	in.DeepCopyInto(out)
	return out
}

type RouteParentStatus struct {
	// ParentRef corresponds with a ParentRef in the spec that this
	// RouteParentStatus struct describes the status of.
	ParentRef ParentReference `json:"parentRef"`
	// ControllerName is a domain/path string that indicates the name of the
	// controller that wrote this status. This corresponds with the
	// controllerName field on GatewayClass.
	// Example: "example.net/gateway-controller".
	// The format of this field is DOMAIN "/" PATH, where DOMAIN and PATH are
	// valid Kubernetes names
	// (https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names).
	// Controllers MUST populate this field when writing status. Controllers should ensure that
	// entries to status populated with their ControllerName are cleaned up when they are no
	// longer necessary.
	ControllerName string `json:"controllerName"`
	// Conditions describes the status of the route with respect to the Gateway.
	// Note that the route's availability is also subject to the Gateway's own
	// status conditions and listener status.
	// If the Route's ParentRef specifies an existing Gateway that supports
	// Routes of this kind AND that Gateway's controller has sufficient access,
	// then that Gateway's controller MUST set the "Accepted" condition on the
	// Route, to indicate whether the route has been accepted or rejected by the
	// Gateway, and why.
	// A Route MUST be considered "Accepted" if at least one of the Route's
	// rules is implemented by the Gateway.
	// There are a number of cases where the "Accepted" condition may not be set
	// due to lack of controller visibility, that includes when:
	// * The Route refers to a nonexistent parent.
	// * The Route is of a type that the controller does not support.
	// * The Route is in a namespace to which the controller does not have access.
	// <gateway:util:excludeFromCRD>
	// Notes for implementors:
	// Conditions are a listType `map`, which means that they function like a
	// map with a key of the `type` field _in the k8s apiserver_.
	// This means that implementations must obey some rules when updating this
	// section.
	// * Implementations MUST perform a read-modify-write cycle on this field
	// before modifying it. That is, when modifying this field, implementations
	// must be confident they have fetched the most recent version of this field,
	// and ensure that changes they make are on that recent version.
	// * Implementations MUST NOT remove or reorder Conditions that they are not
	// directly responsible for. For example, if an implementation sees a Condition
	// with type `special.io/SomeField`, it MUST NOT remove, change or update that
	// Condition.
	// * Implementations MUST always _merge_ changes into Conditions of the same Type,
	// rather than creating more than one Condition of the same Type.
	// * Implementations MUST always update the `observedGeneration` field of the
	// Condition to the `metadata.generation` of the Gateway at the time of update creation.
	// * If the `observedGeneration` of a Condition is _greater than_ the value the
	// implementation knows about, then it MUST NOT perform the update on that Condition,
	// but must wait for a future reconciliation and status update. (The assumption is that
	// the implementation's copy of the object is stale and an update will be re-triggered
	// if relevant.)
	// </gateway:util:excludeFromCRD>
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func (in *RouteParentStatus) DeepCopyInto(out *RouteParentStatus) {
	*out = *in
	in.ParentRef.DeepCopyInto(&out.ParentRef)
	if in.Conditions != nil {
		l := make([]metav1.Condition, len(in.Conditions))
		for i := range in.Conditions {
			in.Conditions[i].DeepCopyInto(&l[i])
		}
		out.Conditions = l
	}
}

func (in *RouteParentStatus) DeepCopy() *RouteParentStatus {
	if in == nil {
		return nil
	}
	out := new(RouteParentStatus)
	in.DeepCopyInto(out)
	return out
}

type ListenerTLSConfig struct {
	// Mode defines the TLS behavior for the TLS session initiated by the client.
	// There are two possible modes:
	// - Terminate: The TLS session between the downstream client and the
	// Gateway is terminated at the Gateway. This mode requires certificates
	// to be specified in some way, such as populating the certificateRefs
	// field.
	// - Passthrough: The TLS session is NOT terminated by the Gateway. This
	// implies that the Gateway can't decipher the TLS stream except for
	// the ClientHello message of the TLS protocol. The certificateRefs field
	// is ignored in this mode.
	// Support: Core
	Mode TLSModeType `json:"mode,omitempty"`
	// CertificateRefs contains a series of references to Kubernetes objects that
	// contains TLS certificates and private keys. These certificates are used to
	// establish a TLS handshake for requests that match the hostname of the
	// associated listener.
	// A single CertificateRef to a Kubernetes Secret has "Core" support.
	// Implementations MAY choose to support attaching multiple certificates to
	// a Listener, but this behavior is implementation-specific.
	// References to a resource in different namespace are invalid UNLESS there
	// is a ReferenceGrant in the target namespace that allows the certificate
	// to be attached. If a ReferenceGrant does not allow this reference, the
	// "ResolvedRefs" condition MUST be set to False for this listener with the
	// "RefNotPermitted" reason.
	// This field is required to have at least one element when the mode is set
	// to "Terminate" (default) and is optional otherwise.
	// CertificateRefs can reference to standard Kubernetes resources, i.e.
	// Secret, or implementation-specific custom resources.
	// Support: Core - A single reference to a Kubernetes Secret of type kubernetes.io/tls
	// Support: Implementation-specific (More than one reference or other resource types)
	CertificateRefs []SecretObjectReference `json:"certificateRefs,omitempty"`
	// Options are a list of key/value pairs to enable extended TLS
	// configuration for each implementation. For example, configuring the
	// minimum TLS version or supported cipher suites.
	// A set of common keys MAY be defined by the API in the future. To avoid
	// any ambiguity, implementation-specific definitions MUST use
	// domain-prefixed names, such as `example.com/my-custom-option`.
	// Un-prefixed names are reserved for key names defined by Gateway API.
	// Support: Implementation-specific
	Options map[string]string `json:"options,omitempty"`
}

func (in *ListenerTLSConfig) DeepCopyInto(out *ListenerTLSConfig) {
	*out = *in
	if in.CertificateRefs != nil {
		l := make([]SecretObjectReference, len(in.CertificateRefs))
		for i := range in.CertificateRefs {
			in.CertificateRefs[i].DeepCopyInto(&l[i])
		}
		out.CertificateRefs = l
	}
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = make(map[string]string, len(*in))
		maps.Copy(*out, *in)
	}
}

func (in *ListenerTLSConfig) DeepCopy() *ListenerTLSConfig {
	if in == nil {
		return nil
	}
	out := new(ListenerTLSConfig)
	in.DeepCopyInto(out)
	return out
}

type AllowedRoutes struct {
	// Namespaces indicates namespaces from which Routes may be attached to this
	// Listener. This is restricted to the namespace of this Gateway by default.
	// Support: Core
	Namespaces *RouteNamespaces `json:"namespaces,omitempty"`
	// Kinds specifies the groups and kinds of Routes that are allowed to bind
	// to this Gateway Listener. When unspecified or empty, the kinds of Routes
	// selected are determined using the Listener protocol.
	// A RouteGroupKind MUST correspond to kinds of Routes that are compatible
	// with the application protocol specified in the Listener's Protocol field.
	// If an implementation does not support or recognize this resource type, it
	// MUST set the "ResolvedRefs" condition to False for this Listener with the
	// "InvalidRouteKinds" reason.
	// Support: Core
	Kinds []RouteGroupKind `json:"kinds,omitempty"`
}

func (in *AllowedRoutes) DeepCopyInto(out *AllowedRoutes) {
	*out = *in
	if in.Namespaces != nil {
		in, out := &in.Namespaces, &out.Namespaces
		*out = new(RouteNamespaces)
		(*in).DeepCopyInto(*out)
	}
	if in.Kinds != nil {
		l := make([]RouteGroupKind, len(in.Kinds))
		for i := range in.Kinds {
			in.Kinds[i].DeepCopyInto(&l[i])
		}
		out.Kinds = l
	}
}

func (in *AllowedRoutes) DeepCopy() *AllowedRoutes {
	if in == nil {
		return nil
	}
	out := new(AllowedRoutes)
	in.DeepCopyInto(out)
	return out
}

type LocalParametersReference struct {
	// Group is the group of the referent.
	Group string `json:"group"`
	// Kind is kind of the referent.
	Kind string `json:"kind"`
	// Name is the name of the referent.
	Name string `json:"name"`
}

func (in *LocalParametersReference) DeepCopyInto(out *LocalParametersReference) {
	*out = *in
}

func (in *LocalParametersReference) DeepCopy() *LocalParametersReference {
	if in == nil {
		return nil
	}
	out := new(LocalParametersReference)
	in.DeepCopyInto(out)
	return out
}

type ListenerNamespaces struct {
	// From indicates where ListenerSets can attach to this Gateway. Possible
	// values are:
	// * Same: Only ListenerSets in the same namespace may be attached to this Gateway.
	// * Selector: ListenerSets in namespaces selected by the selector may be attached to this Gateway.
	// * All: ListenerSets in all namespaces may be attached to this Gateway.
	// * None: Only listeners defined in the Gateway's spec are allowed
	// The default value None
	From FromNamespaces `json:"from,omitempty"`
	// Selector must be specified when From is set to "Selector". In that case,
	// only ListenerSets in Namespaces matching this Selector will be selected by this
	// Gateway. This field is ignored for other values of "From".
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

func (in *ListenerNamespaces) DeepCopyInto(out *ListenerNamespaces) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

func (in *ListenerNamespaces) DeepCopy() *ListenerNamespaces {
	if in == nil {
		return nil
	}
	out := new(ListenerNamespaces)
	in.DeepCopyInto(out)
	return out
}

type GatewayBackendTLS struct {
	// ClientCertificateRef references an object that contains a client certificate
	// and its associated private key. It can reference standard Kubernetes resources,
	// i.e., Secret, or implementation-specific custom resources.
	// A ClientCertificateRef is considered invalid if:
	// * It refers to a resource that cannot be resolved (e.g., the referenced resource
	// does not exist) or is misconfigured (e.g., a Secret does not contain the keys
	// named `tls.crt` and `tls.key`). In this case, the `ResolvedRefs` condition
	// on the Gateway MUST be set to False with the Reason `InvalidClientCertificateRef`
	// and the Message of the Condition MUST indicate why the reference is invalid.
	// * It refers to a resource in another namespace UNLESS there is a ReferenceGrant
	// in the target namespace that allows the certificate to be attached.
	// If a ReferenceGrant does not allow this reference, the `ResolvedRefs` condition
	// on the Gateway MUST be set to False with the Reason `RefNotPermitted`.
	// Implementations MAY choose to perform further validation of the certificate
	// content (e.g., checking expiry or enforcing specific formats). In such cases,
	// an implementation-specific Reason and Message MUST be set.
	// Support: Core - Reference to a Kubernetes TLS Secret (with the type `kubernetes.io/tls`).
	// Support: Implementation-specific - Other resource kinds or Secrets with a
	// different type (e.g., `Opaque`).
	ClientCertificateRef *SecretObjectReference `json:"clientCertificateRef,omitempty"`
}

func (in *GatewayBackendTLS) DeepCopyInto(out *GatewayBackendTLS) {
	*out = *in
	if in.ClientCertificateRef != nil {
		in, out := &in.ClientCertificateRef, &out.ClientCertificateRef
		*out = new(SecretObjectReference)
		(*in).DeepCopyInto(*out)
	}
}

func (in *GatewayBackendTLS) DeepCopy() *GatewayBackendTLS {
	if in == nil {
		return nil
	}
	out := new(GatewayBackendTLS)
	in.DeepCopyInto(out)
	return out
}

type FrontendTLSConfig struct {
	// Default specifies the default client certificate validation configuration
	// for all Listeners handling HTTPS traffic, unless a per-port configuration
	// is defined.
	// support: Core
	Default TLSConfig `json:"default"`
	// PerPort specifies tls configuration assigned per port.
	// Per port configuration is optional. Once set this configuration overrides
	// the default configuration for all Listeners handling HTTPS traffic
	// that match this port.
	// Each override port requires a unique TLS configuration.
	// support: Core
	PerPort []TLSPortConfig `json:"perPort,omitempty"`
}

func (in *FrontendTLSConfig) DeepCopyInto(out *FrontendTLSConfig) {
	*out = *in
	in.Default.DeepCopyInto(&out.Default)
	if in.PerPort != nil {
		l := make([]TLSPortConfig, len(in.PerPort))
		for i := range in.PerPort {
			in.PerPort[i].DeepCopyInto(&l[i])
		}
		out.PerPort = l
	}
}

func (in *FrontendTLSConfig) DeepCopy() *FrontendTLSConfig {
	if in == nil {
		return nil
	}
	out := new(FrontendTLSConfig)
	in.DeepCopyInto(out)
	return out
}

type RouteGroupKind struct {
	// Group is the group of the Route.
	Group string `json:"group"`
	// Kind is the kind of the Route.
	Kind string `json:"kind"`
}

func (in *RouteGroupKind) DeepCopyInto(out *RouteGroupKind) {
	*out = *in
}

func (in *RouteGroupKind) DeepCopy() *RouteGroupKind {
	if in == nil {
		return nil
	}
	out := new(RouteGroupKind)
	in.DeepCopyInto(out)
	return out
}

type HTTPRouteMatch struct {
	// Path specifies a HTTP request path matcher. If this field is not
	// specified, a default prefix match on the "/" path is provided.
	Path *HTTPPathMatch `json:"path,omitempty"`
	// Headers specifies HTTP request header matchers. Multiple match values are
	// ANDed together, meaning, a request must match all the specified headers
	// to select the route.
	Headers []HTTPHeaderMatch `json:"headers,omitempty"`
	// QueryParams specifies HTTP query parameter matchers. Multiple match
	// values are ANDed together, meaning, a request must match all the
	// specified query parameters to select the route.
	// Support: Extended
	QueryParams []HTTPQueryParamMatch `json:"queryParams,omitempty"`
	// Method specifies HTTP method matcher.
	// When specified, this route will be matched only if the request has the
	// specified method.
	// Support: Extended
	Method HTTPMethod `json:"method,omitempty"`
}

func (in *HTTPRouteMatch) DeepCopyInto(out *HTTPRouteMatch) {
	*out = *in
	if in.Path != nil {
		in, out := &in.Path, &out.Path
		*out = new(HTTPPathMatch)
		(*in).DeepCopyInto(*out)
	}
	if in.Headers != nil {
		l := make([]HTTPHeaderMatch, len(in.Headers))
		for i := range in.Headers {
			in.Headers[i].DeepCopyInto(&l[i])
		}
		out.Headers = l
	}
	if in.QueryParams != nil {
		l := make([]HTTPQueryParamMatch, len(in.QueryParams))
		for i := range in.QueryParams {
			in.QueryParams[i].DeepCopyInto(&l[i])
		}
		out.QueryParams = l
	}
}

func (in *HTTPRouteMatch) DeepCopy() *HTTPRouteMatch {
	if in == nil {
		return nil
	}
	out := new(HTTPRouteMatch)
	in.DeepCopyInto(out)
	return out
}

type HTTPRouteFilter struct {
	// Type identifies the type of filter to apply. As with other API fields,
	// types are classified into three conformance levels:
	// - Core: Filter types and their corresponding configuration defined by
	// "Support: Core" in this package, e.g. "RequestHeaderModifier". All
	// implementations must support core filters.
	// - Extended: Filter types and their corresponding configuration defined by
	// "Support: Extended" in this package, e.g. "RequestMirror". Implementers
	// are encouraged to support extended filters.
	// - Implementation-specific: Filters that are defined and supported by
	// specific vendors.
	// In the future, filters showing convergence in behavior across multiple
	// implementations will be considered for inclusion in extended or core
	// conformance levels. Filter-specific configuration for such filters
	// is specified using the ExtensionRef field. `Type` should be set to
	// "ExtensionRef" for custom filters.
	// Implementers are encouraged to define custom implementation types to
	// extend the core API with implementation-specific behavior.
	// If a reference to a custom filter type cannot be resolved, the filter
	// MUST NOT be skipped. Instead, requests that would have been processed by
	// that filter MUST receive a HTTP error response.
	// Note that values may be added to this enum, implementations
	// must ensure that unknown values will not cause a crash.
	// Unknown values here must result in the implementation setting the
	// Accepted Condition for the Route to `status: False`, with a
	// Reason of `UnsupportedValue`.
	// <gateway:experimental:validation:Enum=RequestHeaderModifier;ResponseHeaderModifier;RequestMirror;RequestRedirect;URLRewrite;ExtensionRef;CORS;ExternalAuth>
	Type HTTPRouteFilterType `json:"type"`
	// RequestHeaderModifier defines a schema for a filter that modifies request
	// headers.
	// Support: Core
	RequestHeaderModifier *HTTPHeaderFilter `json:"requestHeaderModifier,omitempty"`
	// ResponseHeaderModifier defines a schema for a filter that modifies response
	// headers.
	// Support: Extended
	ResponseHeaderModifier *HTTPHeaderFilter `json:"responseHeaderModifier,omitempty"`
	// RequestMirror defines a schema for a filter that mirrors requests.
	// Requests are sent to the specified destination, but responses from
	// that destination are ignored.
	// This filter can be used multiple times within the same rule. Note that
	// not all implementations will be able to support mirroring to multiple
	// backends.
	// Support: Extended
	RequestMirror *HTTPRequestMirrorFilter `json:"requestMirror,omitempty"`
	// RequestRedirect defines a schema for a filter that responds to the
	// request with an HTTP redirection.
	// Support: Core
	RequestRedirect *HTTPRequestRedirectFilter `json:"requestRedirect,omitempty"`
	// URLRewrite defines a schema for a filter that modifies a request during forwarding.
	// Support: Extended
	URLRewrite *HTTPURLRewriteFilter `json:"urlRewrite,omitempty"`
	// Cors defines a schema for a filter that responds to the
	// cross-origin request based on HTTP response header.
	// Support: Extended
	CORS *HTTPCORSFilter `json:"cors,omitempty"`
	// ExternalAuth configures settings related to sending request details
	// to an external auth service. The external service MUST authenticate
	// the request, and MAY authorize the request as well.
	// If there is any problem communicating with the external service,
	// this filter MUST fail closed.
	// Support: Extended
	// <gateway:experimental>
	ExternalAuth *HTTPExternalAuthFilter `json:"externalAuth,omitempty"`
	// ExtensionRef is an optional, implementation-specific extension to the
	// "filter" behavior.  For example, resource "myroutefilter" in group
	// "networking.example.net"). ExtensionRef MUST NOT be used for core and
	// extended filters.
	// This filter can be used multiple times within the same rule.
	// Support: Implementation-specific
	ExtensionRef *LocalObjectReference `json:"extensionRef,omitempty"`
}

func (in *HTTPRouteFilter) DeepCopyInto(out *HTTPRouteFilter) {
	*out = *in
	if in.RequestHeaderModifier != nil {
		in, out := &in.RequestHeaderModifier, &out.RequestHeaderModifier
		*out = new(HTTPHeaderFilter)
		(*in).DeepCopyInto(*out)
	}
	if in.ResponseHeaderModifier != nil {
		in, out := &in.ResponseHeaderModifier, &out.ResponseHeaderModifier
		*out = new(HTTPHeaderFilter)
		(*in).DeepCopyInto(*out)
	}
	if in.RequestMirror != nil {
		in, out := &in.RequestMirror, &out.RequestMirror
		*out = new(HTTPRequestMirrorFilter)
		(*in).DeepCopyInto(*out)
	}
	if in.RequestRedirect != nil {
		in, out := &in.RequestRedirect, &out.RequestRedirect
		*out = new(HTTPRequestRedirectFilter)
		(*in).DeepCopyInto(*out)
	}
	if in.URLRewrite != nil {
		in, out := &in.URLRewrite, &out.URLRewrite
		*out = new(HTTPURLRewriteFilter)
		(*in).DeepCopyInto(*out)
	}
	if in.CORS != nil {
		in, out := &in.CORS, &out.CORS
		*out = new(HTTPCORSFilter)
		(*in).DeepCopyInto(*out)
	}
	if in.ExternalAuth != nil {
		in, out := &in.ExternalAuth, &out.ExternalAuth
		*out = new(HTTPExternalAuthFilter)
		(*in).DeepCopyInto(*out)
	}
	if in.ExtensionRef != nil {
		in, out := &in.ExtensionRef, &out.ExtensionRef
		*out = new(LocalObjectReference)
		(*in).DeepCopyInto(*out)
	}
}

func (in *HTTPRouteFilter) DeepCopy() *HTTPRouteFilter {
	if in == nil {
		return nil
	}
	out := new(HTTPRouteFilter)
	in.DeepCopyInto(out)
	return out
}

type HTTPBackendRef struct {
	// BackendRef is a reference to a backend to forward matched requests to.
	// A BackendRef can be invalid for the following reasons. In all cases, the
	// implementation MUST ensure the `ResolvedRefs` Condition on the Route
	// is set to `status: False`, with a Reason and Message that indicate
	// what is the cause of the error.
	// A BackendRef is invalid if:
	// * It refers to an unknown or unsupported kind of resource. In this
	// case, the Reason must be set to `InvalidKind` and Message of the
	// Condition must explain which kind of resource is unknown or unsupported.
	// * It refers to a resource that does not exist. In this case, the Reason must
	// be set to `BackendNotFound` and the Message of the Condition must explain
	// which resource does not exist.
	// * It refers a resource in another namespace when the reference has not been
	// explicitly allowed by a ReferenceGrant (or equivalent concept). In this
	// case, the Reason must be set to `RefNotPermitted` and the Message of the
	// Condition must explain which cross-namespace reference is not allowed.
	// * It refers to a Kubernetes Service that has an incompatible appProtocol
	// for the given Route type
	// * The BackendTLSPolicy object is installed in the cluster, a BackendTLSPolicy
	// is present that refers to the Service, and the implementation is unable
	// to meet the requirement.
	// Support: Core for Kubernetes Service
	// Support: Implementation-specific for any other resource
	// Support for weight: Core
	// Support for Kubernetes Service appProtocol: Extended
	// Support for BackendTLSPolicy: Extended
	BackendRef `json:",inline"`
	// Filters defined at this level should be executed if and only if the
	// request is being forwarded to the backend defined here.
	// Support: Implementation-specific (For broader support of filters, use the
	// Filters field in HTTPRouteRule.)
	Filters []HTTPRouteFilter `json:"filters,omitempty"`
}

func (in *HTTPBackendRef) DeepCopyInto(out *HTTPBackendRef) {
	*out = *in
	out.BackendRef = in.BackendRef
	if in.Filters != nil {
		l := make([]HTTPRouteFilter, len(in.Filters))
		for i := range in.Filters {
			in.Filters[i].DeepCopyInto(&l[i])
		}
		out.Filters = l
	}
}

func (in *HTTPBackendRef) DeepCopy() *HTTPBackendRef {
	if in == nil {
		return nil
	}
	out := new(HTTPBackendRef)
	in.DeepCopyInto(out)
	return out
}

type HTTPRouteTimeouts struct {
	// Request specifies the maximum duration for a gateway to respond to an HTTP request.
	// If the gateway has not been able to respond before this deadline is met, the gateway
	// MUST return a timeout error.
	// For example, setting the `rules.timeouts.request` field to the value `10s` in an
	// `HTTPRoute` will cause a timeout if a client request is taking longer than 10 seconds
	// to complete.
	// Setting a timeout to the zero duration (e.g. "0s") SHOULD disable the timeout
	// completely. Implementations that cannot completely disable the timeout MUST
	// instead interpret the zero duration as the longest possible value to which
	// the timeout can be set.
	// This timeout is intended to cover as close to the whole request-response transaction
	// as possible although an implementation MAY choose to start the timeout after the entire
	// request stream has been received instead of immediately after the transaction is
	// initiated by the client.
	// The value of Request is a Gateway API Duration string as defined by GEP-2257. When this
	// field is unspecified, request timeout behavior is implementation-specific.
	// Support: Extended
	Request string `json:"request"`
	// BackendRequest specifies a timeout for an individual request from the gateway
	// to a backend. This covers the time from when the request first starts being
	// sent from the gateway to when the full response has been received from the backend.
	// Setting a timeout to the zero duration (e.g. "0s") SHOULD disable the timeout
	// completely. Implementations that cannot completely disable the timeout MUST
	// instead interpret the zero duration as the longest possible value to which
	// the timeout can be set.
	// An entire client HTTP transaction with a gateway, covered by the Request timeout,
	// may result in more than one call from the gateway to the destination backend,
	// for example, if automatic retries are supported.
	// The value of BackendRequest must be a Gateway API Duration string as defined by
	// GEP-2257.  When this field is unspecified, its behavior is implementation-specific;
	// when specified, the value of BackendRequest must be no more than the value of the
	// Request timeout (since the Request timeout encompasses the BackendRequest timeout).
	// Support: Extended
	BackendRequest string `json:"backendRequest"`
}

func (in *HTTPRouteTimeouts) DeepCopyInto(out *HTTPRouteTimeouts) {
	*out = *in
}

func (in *HTTPRouteTimeouts) DeepCopy() *HTTPRouteTimeouts {
	if in == nil {
		return nil
	}
	out := new(HTTPRouteTimeouts)
	in.DeepCopyInto(out)
	return out
}

type HTTPRouteRetry struct {
	// Codes defines the HTTP response status codes for which a backend request
	// should be retried.
	// Support: Extended
	Codes []int `json:"codes,omitempty"`
	// Attempts specifies the maximum number of times an individual request
	// from the gateway to a backend should be retried.
	// If the maximum number of retries has been attempted without a successful
	// response from the backend, the Gateway MUST return an error.
	// When this field is unspecified, the number of times to attempt to retry
	// a backend request is implementation-specific.
	// Support: Extended
	Attempts int `json:"attempts,omitempty"`
	// Backoff specifies the minimum duration a Gateway should wait between
	// retry attempts and is represented in Gateway API Duration formatting.
	// For example, setting the `rules[].retry.backoff` field to the value
	// `100ms` will cause a backend request to first be retried approximately
	// 100 milliseconds after timing out or receiving a response code configured
	// to be retriable.
	// An implementation MAY use an exponential or alternative backoff strategy
	// for subsequent retry attempts, MAY cap the maximum backoff duration to
	// some amount greater than the specified minimum, and MAY add arbitrary
	// jitter to stagger requests, as long as unsuccessful backend requests are
	// not retried before the configured minimum duration.
	// If a Request timeout (`rules[].timeouts.request`) is configured on the
	// route, the entire duration of the initial request and any retry attempts
	// MUST not exceed the Request timeout duration. If any retry attempts are
	// still in progress when the Request timeout duration has been reached,
	// these SHOULD be canceled if possible and the Gateway MUST immediately
	// return a timeout error.
	// If a BackendRequest timeout (`rules[].timeouts.backendRequest`) is
	// configured on the route, any retry attempts which reach the configured
	// BackendRequest timeout duration without a response SHOULD be canceled if
	// possible and the Gateway should wait for at least the specified backoff
	// duration before attempting to retry the backend request again.
	// If a BackendRequest timeout is _not_ configured on the route, retry
	// attempts MAY time out after an implementation default duration, or MAY
	// remain pending until a configured Request timeout or implementation
	// default duration for total request time is reached.
	// When this field is unspecified, the time to wait between retry attempts
	// is implementation-specific.
	// Support: Extended
	Backoff string `json:"backoff"`
}

func (in *HTTPRouteRetry) DeepCopyInto(out *HTTPRouteRetry) {
	*out = *in
	if in.Codes != nil {
		t := make([]int, len(in.Codes))
		copy(t, in.Codes)
		out.Codes = t
	}
}

func (in *HTTPRouteRetry) DeepCopy() *HTTPRouteRetry {
	if in == nil {
		return nil
	}
	out := new(HTTPRouteRetry)
	in.DeepCopyInto(out)
	return out
}

type BackendRef struct {
	// BackendObjectReference references a Kubernetes object.
	BackendObjectReference `json:",inline"`
	// Weight specifies the proportion of requests forwarded to the referenced
	// backend. This is computed as weight/(sum of all weights in this
	// BackendRefs list). For non-zero values, there may be some epsilon from
	// the exact proportion defined here depending on the precision an
	// implementation supports. Weight is not a percentage and the sum of
	// weights does not need to equal 100.
	// If only one backend is specified and it has a weight greater than 0, 100%
	// of the traffic is forwarded to that backend. If weight is set to 0, no
	// traffic should be forwarded for this entry. If unspecified, weight
	// defaults to 1.
	// Support for this field varies based on the context where used.
	Weight int `json:"weight,omitempty"`
}

func (in *BackendRef) DeepCopyInto(out *BackendRef) {
	*out = *in
	out.BackendObjectReference = in.BackendObjectReference
}

func (in *BackendRef) DeepCopy() *BackendRef {
	if in == nil {
		return nil
	}
	out := new(BackendRef)
	in.DeepCopyInto(out)
	return out
}

type GRPCMethodMatch struct {
	// Type specifies how to match against the service and/or method.
	// Support: Core (Exact with service and method specified)
	// Support: Implementation-specific (Exact with method specified but no service specified)
	// Support: Implementation-specific (RegularExpression)
	Type GRPCMethodMatchType `json:"type,omitempty"`
	// Value of the service to match against. If left empty or omitted, will
	// match any service.
	// At least one of Service and Method MUST be a non-empty string.
	Service string `json:"service,omitempty"`
	// Value of the method to match against. If left empty or omitted, will
	// match all services.
	// At least one of Service and Method MUST be a non-empty string.
	Method string `json:"method,omitempty"`
}

func (in *GRPCMethodMatch) DeepCopyInto(out *GRPCMethodMatch) {
	*out = *in
}

func (in *GRPCMethodMatch) DeepCopy() *GRPCMethodMatch {
	if in == nil {
		return nil
	}
	out := new(GRPCMethodMatch)
	in.DeepCopyInto(out)
	return out
}

type GRPCHeaderMatch struct {
	// Type specifies how to match against the value of the header.
	Type GRPCHeaderMatchType `json:"type,omitempty"`
	// Name is the name of the gRPC Header to be matched.
	// If multiple entries specify equivalent header names, only the first
	// entry with an equivalent name MUST be considered for a match. Subsequent
	// entries with an equivalent header name MUST be ignored. Due to the
	// case-insensitivity of header names, "foo" and "Foo" are considered
	// equivalent.
	Name string `json:"name"`
	// Value is the value of the gRPC Header to be matched.
	Value string `json:"value"`
}

func (in *GRPCHeaderMatch) DeepCopyInto(out *GRPCHeaderMatch) {
	*out = *in
}

func (in *GRPCHeaderMatch) DeepCopy() *GRPCHeaderMatch {
	if in == nil {
		return nil
	}
	out := new(GRPCHeaderMatch)
	in.DeepCopyInto(out)
	return out
}

type HTTPHeaderFilter struct {
	// Set overwrites the request with the given header (name, value)
	// before the action.
	// Input:
	// GET /foo HTTP/1.1
	// my-header: foo
	// Config:
	// set:
	// - name: "my-header"
	// value: "bar"
	// Output:
	// GET /foo HTTP/1.1
	// my-header: bar
	Set []HTTPHeader `json:"set,omitempty"`
	// Add adds the given header(s) (name, value) to the request
	// before the action. It appends to any existing values associated
	// with the header name.
	// Input:
	// GET /foo HTTP/1.1
	// my-header: foo
	// Config:
	// add:
	// - name: "my-header"
	// value: "bar,baz"
	// Output:
	// GET /foo HTTP/1.1
	// my-header: foo,bar,baz
	Add []HTTPHeader `json:"add,omitempty"`
	// Remove the given header(s) from the HTTP request before the action. The
	// value of Remove is a list of HTTP header names. Note that the header
	// names are case-insensitive (see
	// https://datatracker.ietf.org/doc/html/rfc2616#section-4.2).
	// Input:
	// GET /foo HTTP/1.1
	// my-header1: foo
	// my-header2: bar
	// my-header3: baz
	// Config:
	// remove: ["my-header1", "my-header3"]
	// Output:
	// GET /foo HTTP/1.1
	// my-header2: bar
	Remove []string `json:"remove,omitempty"`
}

func (in *HTTPHeaderFilter) DeepCopyInto(out *HTTPHeaderFilter) {
	*out = *in
	if in.Set != nil {
		l := make([]HTTPHeader, len(in.Set))
		for i := range in.Set {
			in.Set[i].DeepCopyInto(&l[i])
		}
		out.Set = l
	}
	if in.Add != nil {
		l := make([]HTTPHeader, len(in.Add))
		for i := range in.Add {
			in.Add[i].DeepCopyInto(&l[i])
		}
		out.Add = l
	}
	if in.Remove != nil {
		t := make([]string, len(in.Remove))
		copy(t, in.Remove)
		out.Remove = t
	}
}

func (in *HTTPHeaderFilter) DeepCopy() *HTTPHeaderFilter {
	if in == nil {
		return nil
	}
	out := new(HTTPHeaderFilter)
	in.DeepCopyInto(out)
	return out
}

type HTTPRequestMirrorFilter struct {
	// BackendRef references a resource where mirrored requests are sent.
	// Mirrored requests must be sent only to a single destination endpoint
	// within this BackendRef, irrespective of how many endpoints are present
	// within this BackendRef.
	// If the referent cannot be found, this BackendRef is invalid and must be
	// dropped from the Gateway. The controller must ensure the "ResolvedRefs"
	// condition on the Route status is set to `status: False` and not configure
	// this backend in the underlying implementation.
	// If there is a cross-namespace reference to an *existing* object
	// that is not allowed by a ReferenceGrant, the controller must ensure the
	// "ResolvedRefs"  condition on the Route is set to `status: False`,
	// with the "RefNotPermitted" reason and not configure this backend in the
	// underlying implementation.
	// In either error case, the Message of the `ResolvedRefs` Condition
	// should be used to provide more detail about the problem.
	// Support: Extended for Kubernetes Service
	// Support: Implementation-specific for any other resource
	BackendRef BackendObjectReference `json:"backendRef"`
	// Percent represents the percentage of requests that should be
	// mirrored to BackendRef. Its minimum value is 0 (indicating 0% of
	// requests) and its maximum value is 100 (indicating 100% of requests).
	// Only one of Fraction or Percent may be specified. If neither field
	// is specified, 100% of requests will be mirrored.
	Percent int `json:"percent,omitempty"`
	// Fraction represents the fraction of requests that should be
	// mirrored to BackendRef.
	// Only one of Fraction or Percent may be specified. If neither field
	// is specified, 100% of requests will be mirrored.
	Fraction *Fraction `json:"fraction,omitempty"`
}

func (in *HTTPRequestMirrorFilter) DeepCopyInto(out *HTTPRequestMirrorFilter) {
	*out = *in
	in.BackendRef.DeepCopyInto(&out.BackendRef)
	if in.Fraction != nil {
		in, out := &in.Fraction, &out.Fraction
		*out = new(Fraction)
		(*in).DeepCopyInto(*out)
	}
}

func (in *HTTPRequestMirrorFilter) DeepCopy() *HTTPRequestMirrorFilter {
	if in == nil {
		return nil
	}
	out := new(HTTPRequestMirrorFilter)
	in.DeepCopyInto(out)
	return out
}

type CookieConfig struct {
	// LifetimeType specifies whether the cookie has a permanent or
	// session-based lifetime. A permanent cookie persists until its
	// specified expiry time, defined by the Expires or Max-Age cookie
	// attributes, while a session cookie is deleted when the current
	// session ends.
	// When set to "Permanent", AbsoluteTimeout indicates the
	// cookie's lifetime via the Expires or Max-Age cookie attributes
	// and is required.
	// When set to "Session", AbsoluteTimeout indicates the
	// absolute lifetime of the cookie tracked by the gateway and
	// is optional.
	// Defaults to "Session".
	// Support: Core for "Session" type
	// Support: Extended for "Permanent" type
	LifetimeType CookieLifetimeType `json:"lifetimeType,omitempty"`
}

func (in *CookieConfig) DeepCopyInto(out *CookieConfig) {
	*out = *in
}

func (in *CookieConfig) DeepCopy() *CookieConfig {
	if in == nil {
		return nil
	}
	out := new(CookieConfig)
	in.DeepCopyInto(out)
	return out
}

type SecretObjectReference struct {
	// Group is the group of the referent. For example, "gateway.networking.k8s.io".
	// When unspecified or empty string, core API group is inferred.
	Group string `json:"group"`
	// Kind is kind of the referent. For example "Secret".
	Kind string `json:"kind"`
	// Name is the name of the referent.
	Name string `json:"name"`
	// Namespace is the namespace of the referenced object. When unspecified, the local
	// namespace is inferred.
	// Note that when a namespace different than the local namespace is specified,
	// a ReferenceGrant object is required in the referent namespace to allow that
	// namespace's owner to accept the reference. See the ReferenceGrant
	// documentation for details.
	// Support: Core
	Namespace string `json:"namespace"`
}

func (in *SecretObjectReference) DeepCopyInto(out *SecretObjectReference) {
	*out = *in
}

func (in *SecretObjectReference) DeepCopy() *SecretObjectReference {
	if in == nil {
		return nil
	}
	out := new(SecretObjectReference)
	in.DeepCopyInto(out)
	return out
}

type RouteNamespaces struct {
	// From indicates where Routes will be selected for this Gateway. Possible
	// values are:
	// * All: Routes in all namespaces may be used by this Gateway.
	// * Selector: Routes in namespaces selected by the selector may be used by
	// this Gateway.
	// * Same: Only Routes in the same namespace may be used by this Gateway.
	// Support: Core
	From FromNamespaces `json:"from,omitempty"`
	// Selector must be specified when From is set to "Selector". In that case,
	// only Routes in Namespaces matching this Selector will be selected by this
	// Gateway. This field is ignored for other values of "From".
	// Support: Core
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

func (in *RouteNamespaces) DeepCopyInto(out *RouteNamespaces) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

func (in *RouteNamespaces) DeepCopy() *RouteNamespaces {
	if in == nil {
		return nil
	}
	out := new(RouteNamespaces)
	in.DeepCopyInto(out)
	return out
}

type TLSConfig struct {
	// Validation holds configuration information for validating the frontend (client).
	// Setting this field will result in mutual authentication when connecting to the gateway.
	// In browsers this may result in a dialog appearing
	// that requests a user to specify the client certificate.
	// The maximum depth of a certificate chain accepted in verification is Implementation specific.
	// Support: Core
	Validation *FrontendTLSValidation `json:"validation,omitempty"`
}

func (in *TLSConfig) DeepCopyInto(out *TLSConfig) {
	*out = *in
	if in.Validation != nil {
		in, out := &in.Validation, &out.Validation
		*out = new(FrontendTLSValidation)
		(*in).DeepCopyInto(*out)
	}
}

func (in *TLSConfig) DeepCopy() *TLSConfig {
	if in == nil {
		return nil
	}
	out := new(TLSConfig)
	in.DeepCopyInto(out)
	return out
}

type TLSPortConfig struct {
	// The Port indicates the Port Number to which the TLS configuration will be
	// applied. This configuration will be applied to all Listeners handling HTTPS
	// traffic that match this port.
	// Support: Core
	Port int `json:"port"`
	// TLS store the configuration that will be applied to all Listeners handling
	// HTTPS traffic and matching given port.
	// Support: Core
	TLS TLSConfig `json:"tls"`
}

func (in *TLSPortConfig) DeepCopyInto(out *TLSPortConfig) {
	*out = *in
	in.TLS.DeepCopyInto(&out.TLS)
}

func (in *TLSPortConfig) DeepCopy() *TLSPortConfig {
	if in == nil {
		return nil
	}
	out := new(TLSPortConfig)
	in.DeepCopyInto(out)
	return out
}

type HTTPPathMatch struct {
	// Type specifies how to match against the path Value.
	// Support: Core (Exact, PathPrefix)
	// Support: Implementation-specific (RegularExpression)
	Type PathMatchType `json:"type,omitempty"`
	// Value of the HTTP path to match against.
	Value string `json:"value,omitempty"`
}

func (in *HTTPPathMatch) DeepCopyInto(out *HTTPPathMatch) {
	*out = *in
}

func (in *HTTPPathMatch) DeepCopy() *HTTPPathMatch {
	if in == nil {
		return nil
	}
	out := new(HTTPPathMatch)
	in.DeepCopyInto(out)
	return out
}

type HTTPHeaderMatch struct {
	// Type specifies how to match against the value of the header.
	// Support: Core (Exact)
	// Support: Implementation-specific (RegularExpression)
	// Since RegularExpression HeaderMatchType has implementation-specific
	// conformance, implementations can support POSIX, PCRE or any other dialects
	// of regular expressions. Please read the implementation's documentation to
	// determine the supported dialect.
	Type HeaderMatchType `json:"type,omitempty"`
	// Name is the name of the HTTP Header to be matched. Name matching MUST be
	// case-insensitive. (See https://tools.ietf.org/html/rfc7230#section-3.2).
	// If multiple entries specify equivalent header names, only the first
	// entry with an equivalent name MUST be considered for a match. Subsequent
	// entries with an equivalent header name MUST be ignored. Due to the
	// case-insensitivity of header names, "foo" and "Foo" are considered
	// equivalent.
	// When a header is repeated in an HTTP request, it is
	// implementation-specific behavior as to how this is represented.
	// Generally, proxies should follow the guidance from the RFC:
	// https://www.rfc-editor.org/rfc/rfc7230.html#section-3.2.2 regarding
	// processing a repeated header, with special handling for "Set-Cookie".
	Name string `json:"name"`
	// Value is the value of HTTP Header to be matched.
	// <gateway:experimental:description>
	// Must consist of printable US-ASCII characters, optionally separated
	// by single tabs or spaces. See: https://tools.ietf.org/html/rfc7230#section-3.2
	// </gateway:experimental:description>
	// <gateway:experimental:validation:Pattern=`^[!-~]+([\t ]?[!-~]+)*$`>
	Value string `json:"value"`
}

func (in *HTTPHeaderMatch) DeepCopyInto(out *HTTPHeaderMatch) {
	*out = *in
}

func (in *HTTPHeaderMatch) DeepCopy() *HTTPHeaderMatch {
	if in == nil {
		return nil
	}
	out := new(HTTPHeaderMatch)
	in.DeepCopyInto(out)
	return out
}

type HTTPQueryParamMatch struct {
	// Type specifies how to match against the value of the query parameter.
	// Support: Extended (Exact)
	// Support: Implementation-specific (RegularExpression)
	// Since RegularExpression QueryParamMatchType has Implementation-specific
	// conformance, implementations can support POSIX, PCRE or any other
	// dialects of regular expressions. Please read the implementation's
	// documentation to determine the supported dialect.
	Type QueryParamMatchType `json:"type,omitempty"`
	// Name is the name of the HTTP query param to be matched. This must be an
	// exact string match. (See
	// https://tools.ietf.org/html/rfc7230#section-2.7.3).
	// If multiple entries specify equivalent query param names, only the first
	// entry with an equivalent name MUST be considered for a match. Subsequent
	// entries with an equivalent query param name MUST be ignored.
	// If a query param is repeated in an HTTP request, the behavior is
	// purposely left undefined, since different data planes have different
	// capabilities. However, it is *recommended* that implementations should
	// match against the first value of the param if the data plane supports it,
	// as this behavior is expected in other load balancing contexts outside of
	// the Gateway API.
	// Users SHOULD NOT route traffic based on repeated query params to guard
	// themselves against potential differences in the implementations.
	Name string `json:"name"`
	// Value is the value of HTTP query param to be matched.
	Value string `json:"value"`
}

func (in *HTTPQueryParamMatch) DeepCopyInto(out *HTTPQueryParamMatch) {
	*out = *in
}

func (in *HTTPQueryParamMatch) DeepCopy() *HTTPQueryParamMatch {
	if in == nil {
		return nil
	}
	out := new(HTTPQueryParamMatch)
	in.DeepCopyInto(out)
	return out
}

type HTTPRequestRedirectFilter struct {
	// Scheme is the scheme to be used in the value of the `Location` header in
	// the response. When empty, the scheme of the request is used.
	// Scheme redirects can affect the port of the redirect, for more information,
	// refer to the documentation for the port field of this filter.
	// Note that values may be added to this enum, implementations
	// must ensure that unknown values will not cause a crash.
	// Unknown values here must result in the implementation setting the
	// Accepted Condition for the Route to `status: False`, with a
	// Reason of `UnsupportedValue`.
	// Support: Extended
	Scheme string `json:"scheme,omitempty"`
	// Hostname is the hostname to be used in the value of the `Location`
	// header in the response.
	// When empty, the hostname in the `Host` header of the request is used.
	// Support: Core
	Hostname string `json:"hostname"`
	// Path defines parameters used to modify the path of the incoming request.
	// The modified path is then used to construct the `Location` header. When
	// empty, the request path is used as-is.
	// Support: Extended
	Path *HTTPPathModifier `json:"path,omitempty"`
	// Port is the port to be used in the value of the `Location`
	// header in the response.
	// If no port is specified, the redirect port MUST be derived using the
	// following rules:
	// * If redirect scheme is not-empty, the redirect port MUST be the well-known
	// port associated with the redirect scheme. Specifically "http" to port 80
	// and "https" to port 443. If the redirect scheme does not have a
	// well-known port, the listener port of the Gateway SHOULD be used.
	// * If redirect scheme is empty, the redirect port MUST be the Gateway
	// Listener port.
	// Implementations SHOULD NOT add the port number in the 'Location'
	// header in the following cases:
	// * A Location header that will use HTTP (whether that is determined via
	// the Listener protocol or the Scheme field) _and_ use port 80.
	// * A Location header that will use HTTPS (whether that is determined via
	// the Listener protocol or the Scheme field) _and_ use port 443.
	// Support: Extended
	Port int `json:"port"`
	// StatusCode is the HTTP status code to be used in response.
	// Note that values may be added to this enum, implementations
	// must ensure that unknown values will not cause a crash.
	// Unknown values here must result in the implementation setting the
	// Accepted Condition for the Route to `status: False`, with a
	// Reason of `UnsupportedValue`.
	// Support: Core
	StatusCode int `json:"statusCode,omitempty"`
}

func (in *HTTPRequestRedirectFilter) DeepCopyInto(out *HTTPRequestRedirectFilter) {
	*out = *in
	if in.Path != nil {
		in, out := &in.Path, &out.Path
		*out = new(HTTPPathModifier)
		(*in).DeepCopyInto(*out)
	}
}

func (in *HTTPRequestRedirectFilter) DeepCopy() *HTTPRequestRedirectFilter {
	if in == nil {
		return nil
	}
	out := new(HTTPRequestRedirectFilter)
	in.DeepCopyInto(out)
	return out
}

type HTTPURLRewriteFilter struct {
	// Hostname is the value to be used to replace the Host header value during
	// forwarding.
	// Support: Extended
	Hostname string `json:"hostname"`
	// Path defines a path rewrite.
	// Support: Extended
	Path *HTTPPathModifier `json:"path,omitempty"`
}

func (in *HTTPURLRewriteFilter) DeepCopyInto(out *HTTPURLRewriteFilter) {
	*out = *in
	if in.Path != nil {
		in, out := &in.Path, &out.Path
		*out = new(HTTPPathModifier)
		(*in).DeepCopyInto(*out)
	}
}

func (in *HTTPURLRewriteFilter) DeepCopy() *HTTPURLRewriteFilter {
	if in == nil {
		return nil
	}
	out := new(HTTPURLRewriteFilter)
	in.DeepCopyInto(out)
	return out
}

type HTTPCORSFilter struct {
	// AllowOrigins indicates whether the response can be shared with requested
	// resource from the given `Origin`.
	// The `Origin` consists of a scheme and a host, with an optional port, and
	// takes the form `<scheme>://<host>(:<port>)`.
	// Valid values for scheme are: `http` and `https`.
	// Valid values for port are any integer between 1 and 65535 (the list of
	// available TCP/UDP ports). Note that, if not included, port `80` is
	// assumed for `http` scheme origins, and port `443` is assumed for `https`
	// origins. This may affect origin matching.
	// The host part of the origin may contain the wildcard character `*`. These
	// wildcard characters behave as follows:
	// * `*` is a greedy match to the _left_, including any number of
	// DNS labels to the left of its position. This also means that
	// `*` will include any number of period `.` characters to the
	// left of its position.
	// * A wildcard by itself matches all hosts.
	// An origin value that includes _only_ the `*` character indicates requests
	// from all `Origin`s are allowed.
	// When the `AllowOrigins` field is configured with multiple origins, it
	// means the server supports clients from multiple origins. If the request
	// `Origin` matches the configured allowed origins, the gateway must return
	// the given `Origin` and sets value of the header
	// `Access-Control-Allow-Origin` same as the `Origin` header provided by the
	// client.
	// The status code of a successful response to a "preflight" request is
	// always an OK status (i.e., 204 or 200).
	// If the request `Origin` does not match the configured allowed origins,
	// the gateway returns 204/200 response but doesn't set the relevant
	// cross-origin response headers. Alternatively, the gateway responds with
	// 403 status to the "preflight" request is denied, coupled with omitting
	// the CORS headers. The cross-origin request fails on the client side.
	// Therefore, the client doesn't attempt the actual cross-origin request.
	// Conversely, if the request `Origin` matches one of the configured
	// allowed origins, the gateway sets the response header
	// `Access-Control-Allow-Origin` to the same value as the `Origin`
	// header provided by the client.
	// When config has the wildcard ("*") in allowOrigins, and the request
	// is not credentialed (e.g., it is a preflight request), the
	// `Access-Control-Allow-Origin` response header either contains the
	// wildcard as well or the Origin from the request.
	// When the request is credentialed, the gateway must not specify the `*`
	// wildcard in the `Access-Control-Allow-Origin` response header. When
	// also the `AllowCredentials` field is true and `AllowOrigins` field
	// specified with the `*` wildcard, the gateway must return a single origin
	// in the value of the `Access-Control-Allow-Origin` response header,
	// instead of specifying the `*` wildcard. The value of the header
	// `Access-Control-Allow-Origin` is same as the `Origin` header provided by
	// the client.
	// Support: Extended
	AllowOrigins []string `json:"allowOrigins,omitempty"`
	// AllowCredentials indicates whether the actual cross-origin request allows
	// to include credentials.
	// When set to true, the gateway will include the `Access-Control-Allow-Credentials`
	// response header with value true (case-sensitive).
	// When set to false or omitted the gateway will omit the header
	// `Access-Control-Allow-Credentials` entirely (this is the standard CORS
	// behavior).
	// Support: Extended
	AllowCredentials bool `json:"allowCredentials,omitempty"`
	// AllowMethods indicates which HTTP methods are supported for accessing the
	// requested resource.
	// Valid values are any method defined by RFC9110, along with the special
	// value `*`, which represents all HTTP methods are allowed.
	// Method names are case-sensitive, so these values are also case-sensitive.
	// (See https://www.rfc-editor.org/rfc/rfc2616#section-5.1.1)
	// Multiple method names in the value of the `Access-Control-Allow-Methods`
	// response header are separated by a comma (",").
	// A CORS-safelisted method is a method that is `GET`, `HEAD`, or `POST`.
	// (See https://fetch.spec.whatwg.org/#cors-safelisted-method) The
	// CORS-safelisted methods are always allowed, regardless of whether they
	// are specified in the `AllowMethods` field.
	// When the `AllowMethods` field is configured with one or more methods, the
	// gateway must return the `Access-Control-Allow-Methods` response header
	// which value is present in the `AllowMethods` field.
	// If the HTTP method of the `Access-Control-Request-Method` request header
	// is not included in the list of methods specified by the response header
	// `Access-Control-Allow-Methods`, it will present an error on the client
	// side.
	// If config contains the wildcard "*" in allowMethods and the request is
	// not credentialed, the `Access-Control-Allow-Methods` response header
	// can either use the `*` wildcard or the value of
	// Access-Control-Request-Method from the request.
	// When the request is credentialed, the gateway must not specify the `*`
	// wildcard in the `Access-Control-Allow-Methods` response header. When
	// also the `AllowCredentials` field is true and `AllowMethods` field
	// specified with the `*` wildcard, the gateway must specify one HTTP method
	// in the value of the Access-Control-Allow-Methods response header. The
	// value of the header `Access-Control-Allow-Methods` is same as the
	// `Access-Control-Request-Method` header provided by the client. If the
	// header `Access-Control-Request-Method` is not included in the request,
	// the gateway will omit the `Access-Control-Allow-Methods` response header,
	// instead of specifying the `*` wildcard.
	// Support: Extended
	AllowMethods []string `json:"allowMethods,omitempty"`
	// AllowHeaders indicates which HTTP request headers are supported for
	// accessing the requested resource.
	// Header names are not case-sensitive.
	// Multiple header names in the value of the `Access-Control-Allow-Headers`
	// response header are separated by a comma (",").
	// When the `AllowHeaders` field is configured with one or more headers, the
	// gateway must return the `Access-Control-Allow-Headers` response header
	// which value is present in the `AllowHeaders` field.
	// If any header name in the `Access-Control-Request-Headers` request header
	// is not included in the list of header names specified by the response
	// header `Access-Control-Allow-Headers`, it will present an error on the
	// client side.
	// If any header name in the `Access-Control-Allow-Headers` response header
	// does not recognize by the client, it will also occur an error on the
	// client side.
	// A wildcard indicates that the requests with all HTTP headers are allowed.
	// If config contains the wildcard "*" in allowHeaders and the request is
	// not credentialed, the `Access-Control-Allow-Headers` response header
	// can either use the `*` wildcard or the value of
	// Access-Control-Request-Headers from the request.
	// When the request is credentialed, the gateway must not specify the `*`
	// wildcard in the `Access-Control-Allow-Headers` response header. When
	// also the `AllowCredentials` field is true and `AllowHeaders` field
	// is specified with the `*` wildcard, the gateway must specify one or more
	// HTTP headers in the value of the `Access-Control-Allow-Headers` response
	// header. The value of the header `Access-Control-Allow-Headers` is same as
	// the `Access-Control-Request-Headers` header provided by the client. If
	// the header `Access-Control-Request-Headers` is not included in the
	// request, the gateway will omit the `Access-Control-Allow-Headers`
	// response header, instead of specifying the `*` wildcard.
	// Support: Extended
	AllowHeaders []string `json:"allowHeaders,omitempty"`
	// ExposeHeaders indicates which HTTP response headers can be exposed
	// to client-side scripts in response to a cross-origin request.
	// A CORS-safelisted response header is an HTTP header in a CORS response
	// that it is considered safe to expose to the client scripts.
	// The CORS-safelisted response headers include the following headers:
	// `Cache-Control`
	// `Content-Language`
	// `Content-Length`
	// `Content-Type`
	// `Expires`
	// `Last-Modified`
	// `Pragma`
	// (See https://fetch.spec.whatwg.org/#cors-safelisted-response-header-name)
	// The CORS-safelisted response headers are exposed to client by default.
	// When an HTTP header name is specified using the `ExposeHeaders` field,
	// this additional header will be exposed as part of the response to the
	// client.
	// Header names are not case-sensitive.
	// Multiple header names in the value of the `Access-Control-Expose-Headers`
	// response header are separated by a comma (",").
	// A wildcard indicates that the responses with all HTTP headers are exposed
	// to clients. The `Access-Control-Expose-Headers` response header can only
	// use `*` wildcard as value when the request is not credentialed.
	// When the `exposeHeaders` config field contains the "*" wildcard and
	// the request is credentialed, the gateway cannot use the `*` wildcard in
	// the `Access-Control-Expose-Headers` response header.
	// Support: Extended
	ExposeHeaders []string `json:"exposeHeaders,omitempty"`
	// MaxAge indicates the duration (in seconds) for the client to cache the
	// results of a "preflight" request.
	// The information provided by the `Access-Control-Allow-Methods` and
	// `Access-Control-Allow-Headers` response headers can be cached by the
	// client until the time specified by `Access-Control-Max-Age` elapses.
	// The default value of `Access-Control-Max-Age` response header is 5
	// (seconds).
	// When the `MaxAge` field is unspecified, the gateway sets the response
	// header "Access-Control-Max-Age: 5" by default.
	MaxAge int `json:"maxAge,omitempty"`
}

func (in *HTTPCORSFilter) DeepCopyInto(out *HTTPCORSFilter) {
	*out = *in
	if in.AllowOrigins != nil {
		t := make([]string, len(in.AllowOrigins))
		copy(t, in.AllowOrigins)
		out.AllowOrigins = t
	}
	if in.AllowMethods != nil {
		t := make([]string, len(in.AllowMethods))
		copy(t, in.AllowMethods)
		out.AllowMethods = t
	}
	if in.AllowHeaders != nil {
		t := make([]string, len(in.AllowHeaders))
		copy(t, in.AllowHeaders)
		out.AllowHeaders = t
	}
	if in.ExposeHeaders != nil {
		t := make([]string, len(in.ExposeHeaders))
		copy(t, in.ExposeHeaders)
		out.ExposeHeaders = t
	}
}

func (in *HTTPCORSFilter) DeepCopy() *HTTPCORSFilter {
	if in == nil {
		return nil
	}
	out := new(HTTPCORSFilter)
	in.DeepCopyInto(out)
	return out
}

type HTTPExternalAuthFilter struct {
	// ExternalAuthProtocol describes which protocol to use when communicating with an
	// ext_authz authorization server.
	// When this is set to GRPC, each backend must use the Envoy ext_authz protocol
	// on the port specified in `backendRefs`. Requests and responses are defined
	// in the protobufs explained at:
	// https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto
	// When this is set to HTTP, each backend must respond with a `200` status
	// code in on a successful authorization. Any other code is considered
	// an authorization failure.
	// Feature Names:
	// GRPC Support - HTTPRouteExternalAuthGRPC
	// HTTP Support - HTTPRouteExternalAuthHTTP
	ExternalAuthProtocol HTTPRouteExternalAuthProtocol `json:"protocol,omitempty"`
	// BackendRef is a reference to a backend to send authorization
	// requests to.
	// The backend must speak the selected protocol (GRPC or HTTP) on the
	// referenced port.
	// If the backend service requires TLS, use BackendTLSPolicy to tell the
	// implementation to supply the TLS details to be used to connect to that
	// backend.
	BackendRef *BackendObjectReference `json:"backendRef,omitempty"`
	// GRPCAuthConfig contains configuration for communication with ext_authz
	// protocol-speaking backends.
	// If unset, implementations must assume the default behavior for each
	// included field is intended.
	GRPCAuthConfig *GRPCAuthConfig `json:"grpc,omitempty"`
	// HTTPAuthConfig contains configuration for communication with HTTP-speaking
	// backends.
	// If unset, implementations must assume the default behavior for each
	// included field is intended.
	HTTPAuthConfig *HTTPAuthConfig `json:"http,omitempty"`
	// ForwardBody controls if requests to the authorization server should include
	// the body of the client request; and if so, how big that body is allowed
	// to be.
	// It is expected that implementations will buffer the request body up to
	// `forwardBody.maxSize` bytes. Bodies over that size must be rejected with a
	// 4xx series error (413 or 403 are common examples), and fail processing
	// of the filter.
	// If unset, or `forwardBody.maxSize` is set to `0`, then the body will not
	// be forwarded.
	// Feature Name: HTTPRouteExternalAuthForwardBody
	ForwardBody *ForwardBodyConfig `json:"forwardBody,omitempty"`
}

func (in *HTTPExternalAuthFilter) DeepCopyInto(out *HTTPExternalAuthFilter) {
	*out = *in
	if in.BackendRef != nil {
		in, out := &in.BackendRef, &out.BackendRef
		*out = new(BackendObjectReference)
		(*in).DeepCopyInto(*out)
	}
	if in.GRPCAuthConfig != nil {
		in, out := &in.GRPCAuthConfig, &out.GRPCAuthConfig
		*out = new(GRPCAuthConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTPAuthConfig != nil {
		in, out := &in.HTTPAuthConfig, &out.HTTPAuthConfig
		*out = new(HTTPAuthConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.ForwardBody != nil {
		in, out := &in.ForwardBody, &out.ForwardBody
		*out = new(ForwardBodyConfig)
		(*in).DeepCopyInto(*out)
	}
}

func (in *HTTPExternalAuthFilter) DeepCopy() *HTTPExternalAuthFilter {
	if in == nil {
		return nil
	}
	out := new(HTTPExternalAuthFilter)
	in.DeepCopyInto(out)
	return out
}

type BackendObjectReference struct {
	// Group is the group of the referent. For example, "gateway.networking.k8s.io".
	// When unspecified or empty string, core API group is inferred.
	Group string `json:"group"`
	// Kind is the Kubernetes resource kind of the referent. For example
	// "Service".
	// Defaults to "Service" when not specified.
	// ExternalName services can refer to CNAME DNS records that may live
	// outside of the cluster and as such are difficult to reason about in
	// terms of conformance. They also may not be safe to forward to (see
	// CVE-2021-25740 for more information). Implementations SHOULD NOT
	// support ExternalName Services.
	// Support: Core (Services with a type other than ExternalName)
	// Support: Implementation-specific (Services with type ExternalName)
	Kind string `json:"kind"`
	// Name is the name of the referent.
	Name string `json:"name"`
	// Namespace is the namespace of the backend. When unspecified, the local
	// namespace is inferred.
	// Note that when a namespace different than the local namespace is specified,
	// a ReferenceGrant object is required in the referent namespace to allow that
	// namespace's owner to accept the reference. See the ReferenceGrant
	// documentation for details.
	// Support: Core
	Namespace string `json:"namespace"`
	// Port specifies the destination port number to use for this resource.
	// Port is required when the referent is a Kubernetes Service. In this
	// case, the port number is the service port number, not the target port.
	// For other resources, destination port might be derived from the referent
	// resource or this field.
	Port int `json:"port"`
}

func (in *BackendObjectReference) DeepCopyInto(out *BackendObjectReference) {
	*out = *in
}

func (in *BackendObjectReference) DeepCopy() *BackendObjectReference {
	if in == nil {
		return nil
	}
	out := new(BackendObjectReference)
	in.DeepCopyInto(out)
	return out
}

type HTTPHeader struct {
	// Name is the name of the HTTP Header to be matched. Name matching MUST be
	// case-insensitive. (See https://tools.ietf.org/html/rfc7230#section-3.2).
	// If multiple entries specify equivalent header names, the first entry with
	// an equivalent name MUST be considered for a match. Subsequent entries
	// with an equivalent header name MUST be ignored. Due to the
	// case-insensitivity of header names, "foo" and "Foo" are considered
	// equivalent.
	Name string `json:"name"`
	// Value is the value of HTTP Header to be matched.
	// <gateway:experimental:description>
	// Must consist of printable US-ASCII characters, optionally separated
	// by single tabs or spaces. See: https://tools.ietf.org/html/rfc7230#section-3.2
	// </gateway:experimental:description>
	// <gateway:experimental:validation:Pattern=`^[!-~]+([\t ]?[!-~]+)*$`>
	Value string `json:"value"`
}

func (in *HTTPHeader) DeepCopyInto(out *HTTPHeader) {
	*out = *in
}

func (in *HTTPHeader) DeepCopy() *HTTPHeader {
	if in == nil {
		return nil
	}
	out := new(HTTPHeader)
	in.DeepCopyInto(out)
	return out
}

type Fraction struct {
	Numerator   int `json:"numerator"`
	Denominator int `json:"denominator,omitempty"`
}

func (in *Fraction) DeepCopyInto(out *Fraction) {
	*out = *in
}

func (in *Fraction) DeepCopy() *Fraction {
	if in == nil {
		return nil
	}
	out := new(Fraction)
	in.DeepCopyInto(out)
	return out
}

type FrontendTLSValidation struct {
	// CACertificateRefs contains one or more references to Kubernetes
	// objects that contain a PEM-encoded TLS CA certificate bundle, which
	// is used as a trust anchor to validate the certificates presented by
	// the client.
	// A CACertificateRef is invalid if:
	// * It refers to a resource that cannot be resolved (e.g., the
	// referenced resource does not exist) or is misconfigured (e.g., a
	// ConfigMap does not contain a key named `ca.crt`). In this case, the
	// Reason on all matching HTTPS listeners must be set to `InvalidCACertificateRef`
	// and the Message of the Condition must indicate which reference is invalid and why.
	// * It refers to an unknown or unsupported kind of resource. In this
	// case, the Reason on all matching HTTPS listeners must be set to
	// `InvalidCACertificateKind` and the Message of the Condition must explain
	// which kind of resource is unknown or unsupported.
	// * It refers to a resource in another namespace UNLESS there is a
	// ReferenceGrant in the target namespace that allows the CA
	// certificate to be attached. If a ReferenceGrant does not allow this
	// reference, the `ResolvedRefs` on all matching HTTPS listeners condition
	// MUST be set with the Reason `RefNotPermitted`.
	// Implementations MAY choose to perform further validation of the
	// certificate content (e.g., checking expiry or enforcing specific formats).
	// In such cases, an implementation-specific Reason and Message MUST be set.
	// In all cases, the implementation MUST ensure that the `ResolvedRefs`
	// condition is set to `status: False` on all targeted listeners (i.e.,
	// listeners serving HTTPS on a matching port). The condition MUST
	// include a Reason and Message that indicate the cause of the error. If
	// ALL CACertificateRefs are invalid, the implementation MUST also ensure
	// the `Accepted` condition on the listener is set to `status: False`, with
	// the Reason `NoValidCACertificate`.
	// Implementations MAY choose to support attaching multiple CA certificates
	// to a listener, but this behavior is implementation-specific.
	// Support: Core - A single reference to a Kubernetes ConfigMap, with the
	// CA certificate in a key named `ca.crt`.
	// Support: Implementation-specific - More than one reference, other kinds
	// of resources, or a single reference that includes multiple certificates.
	CACertificateRefs []ObjectReference `json:"caCertificateRefs,omitempty"`
	// FrontendValidationMode defines the mode for validating the client certificate.
	// There are two possible modes:
	// - AllowValidOnly: In this mode, the gateway will accept connections only if
	// the client presents a valid certificate. This certificate must successfully
	// pass validation against the CA certificates specified in `CACertificateRefs`.
	// - AllowInsecureFallback: In this mode, the gateway will accept connections
	// even if the client certificate is not presented or fails verification.
	// This approach delegates client authorization to the backend and introduce
	// a significant security risk. It should be used in testing environments or
	// on a temporary basis in non-testing environments.
	// Defaults to AllowValidOnly.
	// Support: Core
	Mode FrontendValidationModeType `json:"mode,omitempty"`
}

func (in *FrontendTLSValidation) DeepCopyInto(out *FrontendTLSValidation) {
	*out = *in
	if in.CACertificateRefs != nil {
		l := make([]ObjectReference, len(in.CACertificateRefs))
		for i := range in.CACertificateRefs {
			in.CACertificateRefs[i].DeepCopyInto(&l[i])
		}
		out.CACertificateRefs = l
	}
}

func (in *FrontendTLSValidation) DeepCopy() *FrontendTLSValidation {
	if in == nil {
		return nil
	}
	out := new(FrontendTLSValidation)
	in.DeepCopyInto(out)
	return out
}

type HTTPPathModifier struct {
	// Type defines the type of path modifier. Additional types may be
	// added in a future release of the API.
	// Note that values may be added to this enum, implementations
	// must ensure that unknown values will not cause a crash.
	// Unknown values here must result in the implementation setting the
	// Accepted Condition for the Route to `status: False`, with a
	// Reason of `UnsupportedValue`.
	Type HTTPPathModifierType `json:"type"`
	// ReplaceFullPath specifies the value with which to replace the full path
	// of a request during a rewrite or redirect.
	ReplaceFullPath string `json:"replaceFullPath,omitempty"`
	// ReplacePrefixMatch specifies the value with which to replace the prefix
	// match of a request during a rewrite or redirect. For example, a request
	// to "/foo/bar" with a prefix match of "/foo" and a ReplacePrefixMatch
	// of "/xyz" would be modified to "/xyz/bar".
	// Note that this matches the behavior of the PathPrefix match type. This
	// matches full path elements. A path element refers to the list of labels
	// in the path split by the `/` separator. When specified, a trailing `/` is
	// ignored. For example, the paths `/abc`, `/abc/`, and `/abc/def` would all
	// match the prefix `/abc`, but the path `/abcd` would not.
	// ReplacePrefixMatch is only compatible with a `PathPrefix` HTTPRouteMatch.
	// Using any other HTTPRouteMatch type on the same HTTPRouteRule will result in
	// the implementation setting the Accepted Condition for the Route to `status: False`.
	// Request Path | Prefix Match | Replace Prefix | Modified Path
	// -------------|--------------|----------------|----------
	// /foo/bar     | /foo         | /xyz           | /xyz/bar
	// /foo/bar     | /foo         | /xyz/          | /xyz/bar
	// /foo/bar     | /foo/        | /xyz           | /xyz/bar
	// /foo/bar     | /foo/        | /xyz/          | /xyz/bar
	// /foo         | /foo         | /xyz           | /xyz
	// /foo/        | /foo         | /xyz           | /xyz/
	// /foo/bar     | /foo         | <empty string> | /bar
	// /foo/        | /foo         | <empty string> | /
	// /foo         | /foo         | <empty string> | /
	// /foo/        | /foo         | /              | /
	// /foo         | /foo         | /              | /
	ReplacePrefixMatch string `json:"replacePrefixMatch,omitempty"`
}

func (in *HTTPPathModifier) DeepCopyInto(out *HTTPPathModifier) {
	*out = *in
}

func (in *HTTPPathModifier) DeepCopy() *HTTPPathModifier {
	if in == nil {
		return nil
	}
	out := new(HTTPPathModifier)
	in.DeepCopyInto(out)
	return out
}

type GRPCAuthConfig struct {
	// AllowedRequestHeaders specifies what headers from the client request
	// will be sent to the authorization server.
	// If this list is empty, then all headers must be sent.
	// If the list has entries, only those entries must be sent.
	AllowedRequestHeaders []string `json:"allowedHeaders,omitempty"`
}

func (in *GRPCAuthConfig) DeepCopyInto(out *GRPCAuthConfig) {
	*out = *in
	if in.AllowedRequestHeaders != nil {
		t := make([]string, len(in.AllowedRequestHeaders))
		copy(t, in.AllowedRequestHeaders)
		out.AllowedRequestHeaders = t
	}
}

func (in *GRPCAuthConfig) DeepCopy() *GRPCAuthConfig {
	if in == nil {
		return nil
	}
	out := new(GRPCAuthConfig)
	in.DeepCopyInto(out)
	return out
}

type HTTPAuthConfig struct {
	// Path sets the prefix that paths from the client request will have added
	// when forwarded to the authorization server.
	// When empty or unspecified, no prefix is added.
	// Valid values are the same as the "value" regex for path values in the `match`
	// stanza, and the validation regex will screen out invalid paths in the same way.
	// Even with the validation, implementations MUST sanitize this input before using it
	// directly.
	Path string `json:"path,omitempty"`
	// AllowedRequestHeaders specifies what additional headers from the client request
	// will be sent to the authorization server.
	// The following headers must always be sent to the authorization server,
	// regardless of this setting:
	// * `Host`
	// * `Method`
	// * `Path`
	// * `Content-Length`
	// * `Authorization`
	// If this list is empty, then only those headers must be sent.
	// Note that `Content-Length` has a special behavior, in that the length
	// sent must be correct for the actual request to the external authorization
	// server - that is, it must reflect the actual number of bytes sent in the
	// body of the request to the authorization server.
	// So if the `forwardBody` stanza is unset, or `forwardBody.maxSize` is set
	// to `0`, then `Content-Length` must be `0`. If `forwardBody.maxSize` is set
	// to anything other than `0`, then the `Content-Length` of the authorization
	// request must be set to the actual number of bytes forwarded.
	AllowedRequestHeaders []string `json:"allowedHeaders,omitempty"`
	// AllowedResponseHeaders specifies what headers from the authorization response
	// will be copied into the request to the backend.
	// If this list is empty, then all headers from the authorization server
	// except Authority or Host must be copied.
	AllowedResponseHeaders []string `json:"allowedResponseHeaders,omitempty"`
}

func (in *HTTPAuthConfig) DeepCopyInto(out *HTTPAuthConfig) {
	*out = *in
	if in.AllowedRequestHeaders != nil {
		t := make([]string, len(in.AllowedRequestHeaders))
		copy(t, in.AllowedRequestHeaders)
		out.AllowedRequestHeaders = t
	}
	if in.AllowedResponseHeaders != nil {
		t := make([]string, len(in.AllowedResponseHeaders))
		copy(t, in.AllowedResponseHeaders)
		out.AllowedResponseHeaders = t
	}
}

func (in *HTTPAuthConfig) DeepCopy() *HTTPAuthConfig {
	if in == nil {
		return nil
	}
	out := new(HTTPAuthConfig)
	in.DeepCopyInto(out)
	return out
}

type ForwardBodyConfig struct {
	// MaxSize specifies how large in bytes the largest body that will be buffered
	// and sent to the authorization server. If the body size is larger than
	// `maxSize`, then the body sent to the authorization server must be
	// truncated to `maxSize` bytes.
	// Experimental note: This behavior needs to be checked against
	// various dataplanes; it may need to be changed.
	// See https://github.com/kubernetes-sigs/gateway-api/pull/4001#discussion_r2291405746
	// for more.
	// If 0, the body will not be sent to the authorization server.
	MaxSize uint32 `json:"maxSize,omitempty"`
}

func (in *ForwardBodyConfig) DeepCopyInto(out *ForwardBodyConfig) {
	*out = *in
}

func (in *ForwardBodyConfig) DeepCopy() *ForwardBodyConfig {
	if in == nil {
		return nil
	}
	out := new(ForwardBodyConfig)
	in.DeepCopyInto(out)
	return out
}

type ObjectReference struct {
	// Group is the group of the referent. For example, "gateway.networking.k8s.io".
	// When set to the empty string, core API group is inferred.
	Group string `json:"group"`
	// Kind is kind of the referent. For example "ConfigMap" or "Service".
	Kind string `json:"kind"`
	// Name is the name of the referent.
	Name string `json:"name"`
	// Namespace is the namespace of the referenced object. When unspecified, the local
	// namespace is inferred.
	// Note that when a namespace different than the local namespace is specified,
	// a ReferenceGrant object is required in the referent namespace to allow that
	// namespace's owner to accept the reference. See the ReferenceGrant
	// documentation for details.
	// Support: Core
	Namespace string `json:"namespace"`
}

func (in *ObjectReference) DeepCopyInto(out *ObjectReference) {
	*out = *in
}

func (in *ObjectReference) DeepCopy() *ObjectReference {
	if in == nil {
		return nil
	}
	out := new(ObjectReference)
	in.DeepCopyInto(out)
	return out
}
