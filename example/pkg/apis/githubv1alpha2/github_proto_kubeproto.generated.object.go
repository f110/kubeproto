package githubv1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "grafana.f110.dev"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1alpha2"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: "grafana.f110.dev", Version: "v1alpha2"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&Grafana{},
		&GrafanaUser{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type Grafana struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              GrafanaSpec   `json:"spec"`
	Status            GrafanaStatus `json:"status"`
}

func (in *Grafana) DeepCopyInto(out *Grafana) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

func (in *Grafana) DeepCopy() *Grafana {
	if in == nil {
		return nil
	}
	out := new(Grafana)
	in.DeepCopyInto(out)
	return out
}

func (in *Grafana) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type GrafanaUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              GrafanaUserSpec   `json:"spec"`
	Status            GrafanaUserStatus `json:"status"`
}

func (in *GrafanaUser) DeepCopyInto(out *GrafanaUser) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

func (in *GrafanaUser) DeepCopy() *GrafanaUser {
	if in == nil {
		return nil
	}
	out := new(GrafanaUser)
	in.DeepCopyInto(out)
	return out
}

func (in *GrafanaUser) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type GrafanaSpec struct {
	AdminUser string `json:"adminUser,omitempty"`
}

func (in *GrafanaSpec) DeepCopyInto(out *GrafanaSpec) {
	*out = *in
}

func (in *GrafanaSpec) DeepCopy() *GrafanaSpec {
	if in == nil {
		return nil
	}
	out := new(GrafanaSpec)
	in.DeepCopyInto(out)
	return out
}

type GrafanaStatus struct {
	ObservedGeneration int64 `json:"observedGeneration"`
	// history is a list of History
	History []History `json:"history"`
}

func (in *GrafanaStatus) DeepCopyInto(out *GrafanaStatus) {
	*out = *in
	if in.History != nil {
		l := make([]History, len(in.History))
		for i := range in.History {
			in.History[i].DeepCopyInto(&l[i])
		}
		out.History = l
	}
}

func (in *GrafanaStatus) DeepCopy() *GrafanaStatus {
	if in == nil {
		return nil
	}
	out := new(GrafanaStatus)
	in.DeepCopyInto(out)
	return out
}

type GrafanaUserSpec struct {
	Email string `json:"email"`
	Admin bool   `json:"admin"`
}

func (in *GrafanaUserSpec) DeepCopyInto(out *GrafanaUserSpec) {
	*out = *in
}

func (in *GrafanaUserSpec) DeepCopy() *GrafanaUserSpec {
	if in == nil {
		return nil
	}
	out := new(GrafanaUserSpec)
	in.DeepCopyInto(out)
	return out
}

type GrafanaUserStatus struct {
	Ready bool `json:"ready"`
}

func (in *GrafanaUserStatus) DeepCopyInto(out *GrafanaUserStatus) {
	*out = *in
}

func (in *GrafanaUserStatus) DeepCopy() *GrafanaUserStatus {
	if in == nil {
		return nil
	}
	out := new(GrafanaUserStatus)
	in.DeepCopyInto(out)
	return out
}

type History struct {
	// message contains the result from the server
	Message string `json:"message"`
}

func (in *History) DeepCopyInto(out *History) {
	*out = *in
}

func (in *History) DeepCopy() *History {
	if in == nil {
		return nil
	}
	out := new(History)
	in.DeepCopyInto(out)
	return out
}
