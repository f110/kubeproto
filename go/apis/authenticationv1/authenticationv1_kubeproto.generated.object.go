package authenticationv1

import (
	"go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "authentication.k8s.io"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&TokenRequest{},
		&TokenRequestList{},
		&TokenReview{},
		&TokenReviewList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type TokenRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec holds information about the request being evaluated
	Spec TokenRequestSpec `json:"spec"`
	// Status is filled in by the server and indicates whether the token can be authenticated.
	Status *TokenRequestStatus `json:"status,omitempty"`
}

func (in *TokenRequest) DeepCopyInto(out *TokenRequest) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(TokenRequestStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *TokenRequest) DeepCopy() *TokenRequest {
	if in == nil {
		return nil
	}
	out := new(TokenRequest)
	in.DeepCopyInto(out)
	return out
}

func (in *TokenRequest) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type TokenRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []TokenRequest `json:"items"`
}

func (in *TokenRequestList) DeepCopyInto(out *TokenRequestList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]TokenRequest, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *TokenRequestList) DeepCopy() *TokenRequestList {
	if in == nil {
		return nil
	}
	out := new(TokenRequestList)
	in.DeepCopyInto(out)
	return out
}

func (in *TokenRequestList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type TokenReview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec holds information about the request being evaluated
	Spec TokenReviewSpec `json:"spec"`
	// Status is filled in by the server and indicates whether the request can be authenticated.
	Status *TokenReviewStatus `json:"status,omitempty"`
}

func (in *TokenReview) DeepCopyInto(out *TokenReview) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(TokenReviewStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *TokenReview) DeepCopy() *TokenReview {
	if in == nil {
		return nil
	}
	out := new(TokenReview)
	in.DeepCopyInto(out)
	return out
}

func (in *TokenReview) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type TokenReviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []TokenReview `json:"items"`
}

func (in *TokenReviewList) DeepCopyInto(out *TokenReviewList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]TokenReview, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *TokenReviewList) DeepCopy() *TokenReviewList {
	if in == nil {
		return nil
	}
	out := new(TokenReviewList)
	in.DeepCopyInto(out)
	return out
}

func (in *TokenReviewList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type TokenRequestSpec struct {
	// Audiences are the intendend audiences of the token. A recipient of a
	// token must identify themself with an identifier in the list of
	// audiences of the token, and otherwise should reject the token. A
	// token issued for multiple audiences may be used to authenticate
	// against any of the audiences listed but implies a high degree of
	// trust between the target audiences.
	Audiences []string `json:"audiences"`
	// ExpirationSeconds is the requested duration of validity of the request. The
	// token issuer may return a token with a different validity duration so a
	// client needs to check the 'expiration' field in a response.
	ExpirationSeconds int64 `json:"expirationSeconds,omitempty"`
	// BoundObjectRef is a reference to an object that the token will be bound to.
	// The token will only be valid for as long as the bound object exists.
	// NOTE: The API server's TokenReview endpoint will validate the
	// BoundObjectRef, but other audiences may not. Keep ExpirationSeconds
	// small if you want prompt revocation.
	BoundObjectRef *BoundObjectReference `json:"boundObjectRef,omitempty"`
}

func (in *TokenRequestSpec) DeepCopyInto(out *TokenRequestSpec) {
	*out = *in
	if in.Audiences != nil {
		t := make([]string, len(in.Audiences))
		copy(t, in.Audiences)
		out.Audiences = t
	}
	if in.BoundObjectRef != nil {
		in, out := &in.BoundObjectRef, &out.BoundObjectRef
		*out = new(BoundObjectReference)
		(*in).DeepCopyInto(*out)
	}
}

func (in *TokenRequestSpec) DeepCopy() *TokenRequestSpec {
	if in == nil {
		return nil
	}
	out := new(TokenRequestSpec)
	in.DeepCopyInto(out)
	return out
}

type TokenRequestStatus struct {
	// Token is the opaque bearer token.
	Token string `json:"token"`
	// ExpirationTimestamp is the time of expiration of the returned token.
	ExpirationTimestamp metav1.Time `json:"expirationTimestamp"`
}

func (in *TokenRequestStatus) DeepCopyInto(out *TokenRequestStatus) {
	*out = *in
	in.ExpirationTimestamp.DeepCopyInto(&out.ExpirationTimestamp)
}

func (in *TokenRequestStatus) DeepCopy() *TokenRequestStatus {
	if in == nil {
		return nil
	}
	out := new(TokenRequestStatus)
	in.DeepCopyInto(out)
	return out
}

type TokenReviewSpec struct {
	// Token is the opaque bearer token.
	Token string `json:"token,omitempty"`
	// Audiences is a list of the identifiers that the resource server presented
	// with the token identifies as. Audience-aware token authenticators will
	// verify that the token was intended for at least one of the audiences in
	// this list. If no audiences are provided, the audience will default to the
	// audience of the Kubernetes apiserver.
	Audiences []string `json:"audiences"`
}

func (in *TokenReviewSpec) DeepCopyInto(out *TokenReviewSpec) {
	*out = *in
	if in.Audiences != nil {
		t := make([]string, len(in.Audiences))
		copy(t, in.Audiences)
		out.Audiences = t
	}
}

func (in *TokenReviewSpec) DeepCopy() *TokenReviewSpec {
	if in == nil {
		return nil
	}
	out := new(TokenReviewSpec)
	in.DeepCopyInto(out)
	return out
}

type TokenReviewStatus struct {
	// Authenticated indicates that the token was associated with a known user.
	Authenticated bool `json:"authenticated,omitempty"`
	// User is the UserInfo associated with the provided token.
	User *UserInfo `json:"user,omitempty"`
	// Audiences are audience identifiers chosen by the authenticator that are
	// compatible with both the TokenReview and token. An identifier is any
	// identifier in the intersection of the TokenReviewSpec audiences and the
	// token's audiences. A client of the TokenReview API that sets the
	// spec.audiences field should validate that a compatible audience identifier
	// is returned in the status.audiences field to ensure that the TokenReview
	// server is audience aware. If a TokenReview returns an empty
	// status.audience field where status.authenticated is "true", the token is
	// valid against the audience of the Kubernetes API server.
	Audiences []string `json:"audiences"`
	// Error indicates that the token couldn't be checked
	Error string `json:"error,omitempty"`
}

func (in *TokenReviewStatus) DeepCopyInto(out *TokenReviewStatus) {
	*out = *in
	if in.User != nil {
		in, out := &in.User, &out.User
		*out = new(UserInfo)
		(*in).DeepCopyInto(*out)
	}
	if in.Audiences != nil {
		t := make([]string, len(in.Audiences))
		copy(t, in.Audiences)
		out.Audiences = t
	}
}

func (in *TokenReviewStatus) DeepCopy() *TokenReviewStatus {
	if in == nil {
		return nil
	}
	out := new(TokenReviewStatus)
	in.DeepCopyInto(out)
	return out
}

type BoundObjectReference struct {
	// Kind of the referent. Valid kinds are 'Pod' and 'Secret'.
	Kind string `json:"kind,omitempty"`
	// API version of the referent.
	APIVersion string `json:"apiVersion,omitempty"`
	// Name of the referent.
	Name string `json:"name,omitempty"`
	// UID of the referent.
	UID string `json:"uid,omitempty"`
}

func (in *BoundObjectReference) DeepCopyInto(out *BoundObjectReference) {
	*out = *in
}

func (in *BoundObjectReference) DeepCopy() *BoundObjectReference {
	if in == nil {
		return nil
	}
	out := new(BoundObjectReference)
	in.DeepCopyInto(out)
	return out
}

type UserInfo struct {
	// The name that uniquely identifies this user among all active users.
	Username string `json:"username,omitempty"`
	// A unique value that identifies this user across time. If this user is
	// deleted and another user by the same name is added, they will have
	// different UIDs.
	UID string `json:"uid,omitempty"`
	// The names of groups this user is a part of.
	Groups []string `json:"groups"`
	// Any additional information provided by the authenticator.
	Extra map[string]ExtraValue `json:"extra,omitempty"`
}

func (in *UserInfo) DeepCopyInto(out *UserInfo) {
	*out = *in
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

func (in *UserInfo) DeepCopy() *UserInfo {
	if in == nil {
		return nil
	}
	out := new(UserInfo)
	in.DeepCopyInto(out)
	return out
}

type ExtraValue []string
