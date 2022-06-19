package githubv1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "grafana.f110.dev"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1alpha1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: "grafana.f110.dev", Version: "v1alpha1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&Grafana{},
		&GrafanaUser{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type GrafanaPhase string

const (
	GrafanaPhaseCreated  GrafanaPhase = "Created"
	GrafanaPhaseCreating GrafanaPhase = "Creating"
)

type GrafanaUserPhase string

const (
	GrafanaUserPhaseCreated GrafanaUserPhase = "Created"
)

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
	AdminUser    string               `json:"adminUser,omitempty"`
	APIVersion   string               `json:"apiVersion"`
	FeatureGates []string             `json:"featureGates"`
	Volumes      []Volume             `json:"volumes"`
	UserSelector metav1.LabelSelector `json:"userSelector"`
}

func (in *GrafanaSpec) DeepCopyInto(out *GrafanaSpec) {
	*out = *in
	if in.FeatureGates != nil {
		t := make([]string, len(in.FeatureGates))
		copy(t, in.FeatureGates)
		out.FeatureGates = t
	}
	if in.Volumes != nil {
		l := make([]Volume, len(in.Volumes))
		for i := range in.Volumes {
			in.Volumes[i].DeepCopyInto(&l[i])
		}
		out.Volumes = l
	}
	in.UserSelector.DeepCopyInto(&out.UserSelector)
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
	ObservedGeneration int64        `json:"observedGeneration"`
	Phase              GrafanaPhase `json:"phase"`
	ObservedTime       *metav1.Time `json:"observedTime,omitempty"`
}

func (in *GrafanaStatus) DeepCopyInto(out *GrafanaStatus) {
	*out = *in
	if in.ObservedTime != nil {
		in, out := &in.ObservedTime, &out.ObservedTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
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
	// admin indicates that the user has the privilege.
	//  If admin is true, the user is an administrator.
	Admin bool `json:"admin"`
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

type Volume struct {
	Name string `json:"name"`
}

func (in *Volume) DeepCopyInto(out *Volume) {
	*out = *in
}

func (in *Volume) DeepCopy() *Volume {
	if in == nil {
		return nil
	}
	out := new(Volume)
	in.DeepCopyInto(out)
	return out
}
