package schedulingv1

import (
	"go.f110.dev/kubeproto/go/apis/corev1"
	"go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "scheduling.k8s.io"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&PriorityClass{},
		&PriorityClassList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type PriorityClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// value represents the integer value of this priority class. This is the actual priority that pods
	// receive when they have the name of this class in their pod spec.
	Value int `json:"value"`
	// globalDefault specifies whether this PriorityClass should be considered as
	// the default priority for pods that do not have any priority class.
	// Only one PriorityClass can be marked as `globalDefault`. However, if more than
	// one PriorityClasses exists with their `globalDefault` field set to true,
	// the smallest value of such global default PriorityClasses will be used as the default priority.
	GlobalDefault bool `json:"globalDefault,omitempty"`
	// description is an arbitrary string that usually provides guidelines on
	// when this priority class should be used.
	Description string `json:"description,omitempty"`
	// preemptionPolicy is the Policy for preempting pods with lower priority.
	// One of Never, PreemptLowerPriority.
	// Defaults to PreemptLowerPriority if unset.
	PreemptionPolicy corev1.PreemptionPolicy `json:"preemptionPolicy,omitempty"`
}

func (in *PriorityClass) DeepCopyInto(out *PriorityClass) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
}

func (in *PriorityClass) DeepCopy() *PriorityClass {
	if in == nil {
		return nil
	}
	out := new(PriorityClass)
	in.DeepCopyInto(out)
	return out
}

func (in *PriorityClass) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type PriorityClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []PriorityClass `json:"items"`
}

func (in *PriorityClassList) DeepCopyInto(out *PriorityClassList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]PriorityClass, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *PriorityClassList) DeepCopy() *PriorityClassList {
	if in == nil {
		return nil
	}
	out := new(PriorityClassList)
	in.DeepCopyInto(out)
	return out
}

func (in *PriorityClassList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
