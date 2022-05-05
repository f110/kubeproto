package githubv1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

type GrafanaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Grafana `json:"items"`
}

func (in *GrafanaList) DeepCopyInto(out *GrafanaList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]Grafana, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *GrafanaList) DeepCopy() *GrafanaList {
	if in == nil {
		return nil
	}
	out := new(GrafanaList)
	in.DeepCopyInto(out)
	return out
}

func (in *GrafanaList) DeepCopyObject() runtime.Object {
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

type GrafanaUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []GrafanaUser `json:"items"`
}

func (in *GrafanaUserList) DeepCopyInto(out *GrafanaUserList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]GrafanaUser, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *GrafanaUserList) DeepCopy() *GrafanaUserList {
	if in == nil {
		return nil
	}
	out := new(GrafanaUserList)
	in.DeepCopyInto(out)
	return out
}

func (in *GrafanaUserList) DeepCopyObject() runtime.Object {
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
}

func (in *GrafanaStatus) DeepCopyInto(out *GrafanaStatus) {
	*out = *in
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