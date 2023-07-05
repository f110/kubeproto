package coordinationv1

import (
	"go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "coordination.k8s.io"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&Lease{},
		&LeaseList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type Lease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// spec contains the specification of the Lease.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec *LeaseSpec `json:"spec,omitempty"`
}

func (in *Lease) DeepCopyInto(out *Lease) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(LeaseSpec)
		(*in).DeepCopyInto(*out)
	}
}

func (in *Lease) DeepCopy() *Lease {
	if in == nil {
		return nil
	}
	out := new(Lease)
	in.DeepCopyInto(out)
	return out
}

func (in *Lease) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type LeaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Lease `json:"items"`
}

func (in *LeaseList) DeepCopyInto(out *LeaseList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]Lease, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *LeaseList) DeepCopy() *LeaseList {
	if in == nil {
		return nil
	}
	out := new(LeaseList)
	in.DeepCopyInto(out)
	return out
}

func (in *LeaseList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type LeaseSpec struct {
	// holderIdentity contains the identity of the holder of a current lease.
	HolderIdentity string `json:"holderIdentity,omitempty"`
	// leaseDurationSeconds is a duration that candidates for a lease need
	// to wait to force acquire it. This is measure against time of last
	// observed renewTime.
	LeaseDurationSeconds int `json:"leaseDurationSeconds,omitempty"`
	// acquireTime is a time when the current lease was acquired.
	AcquireTime *metav1.MicroTime `json:"acquireTime,omitempty"`
	// renewTime is a time when the current holder of a lease has last
	// updated the lease.
	RenewTime *metav1.MicroTime `json:"renewTime,omitempty"`
	// leaseTransitions is the number of transitions of a lease between
	// holders.
	LeaseTransitions int `json:"leaseTransitions,omitempty"`
}

func (in *LeaseSpec) DeepCopyInto(out *LeaseSpec) {
	*out = *in
	if in.AcquireTime != nil {
		in, out := &in.AcquireTime, &out.AcquireTime
		*out = new(metav1.MicroTime)
		(*in).DeepCopyInto(*out)
	}
	if in.RenewTime != nil {
		in, out := &in.RenewTime, &out.RenewTime
		*out = new(metav1.MicroTime)
		(*in).DeepCopyInto(*out)
	}
}

func (in *LeaseSpec) DeepCopy() *LeaseSpec {
	if in == nil {
		return nil
	}
	out := new(LeaseSpec)
	in.DeepCopyInto(out)
	return out
}
