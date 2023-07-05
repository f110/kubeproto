package policyv1

import (
	"go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilintstr "k8s.io/apimachinery/pkg/util/intstr"
)

const GroupName = "policy"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&Eviction{},
		&EvictionList{},
		&PodDisruptionBudget{},
		&PodDisruptionBudgetList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type UnhealthyPodEvictionPolicyType string

const (
	UnhealthyPodEvictionPolicyTypeIfHealthyBudget UnhealthyPodEvictionPolicyType = "IfHealthyBudget"
	UnhealthyPodEvictionPolicyTypeAlwaysAllow     UnhealthyPodEvictionPolicyType = "AlwaysAllow"
)

type Eviction struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// DeleteOptions may be provided
	DeleteOptions *metav1.DeleteOptions `json:"deleteOptions,omitempty"`
}

func (in *Eviction) DeepCopyInto(out *Eviction) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.DeleteOptions != nil {
		in, out := &in.DeleteOptions, &out.DeleteOptions
		*out = new(metav1.DeleteOptions)
		(*in).DeepCopyInto(*out)
	}
}

func (in *Eviction) DeepCopy() *Eviction {
	if in == nil {
		return nil
	}
	out := new(Eviction)
	in.DeepCopyInto(out)
	return out
}

func (in *Eviction) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type EvictionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Eviction `json:"items"`
}

func (in *EvictionList) DeepCopyInto(out *EvictionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]Eviction, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *EvictionList) DeepCopy() *EvictionList {
	if in == nil {
		return nil
	}
	out := new(EvictionList)
	in.DeepCopyInto(out)
	return out
}

func (in *EvictionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type PodDisruptionBudget struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Specification of the desired behavior of the PodDisruptionBudget.
	Spec *PodDisruptionBudgetSpec `json:"spec,omitempty"`
	// Most recently observed status of the PodDisruptionBudget.
	Status *PodDisruptionBudgetStatus `json:"status,omitempty"`
}

func (in *PodDisruptionBudget) DeepCopyInto(out *PodDisruptionBudget) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(PodDisruptionBudgetSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(PodDisruptionBudgetStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *PodDisruptionBudget) DeepCopy() *PodDisruptionBudget {
	if in == nil {
		return nil
	}
	out := new(PodDisruptionBudget)
	in.DeepCopyInto(out)
	return out
}

func (in *PodDisruptionBudget) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type PodDisruptionBudgetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []PodDisruptionBudget `json:"items"`
}

func (in *PodDisruptionBudgetList) DeepCopyInto(out *PodDisruptionBudgetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]PodDisruptionBudget, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *PodDisruptionBudgetList) DeepCopy() *PodDisruptionBudgetList {
	if in == nil {
		return nil
	}
	out := new(PodDisruptionBudgetList)
	in.DeepCopyInto(out)
	return out
}

func (in *PodDisruptionBudgetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type PodDisruptionBudgetSpec struct {
	// An eviction is allowed if at least "minAvailable" pods selected by
	// "selector" will still be available after the eviction, i.e. even in the
	// absence of the evicted pod.  So for example you can prevent all voluntary
	// evictions by specifying "100%".
	MinAvailable *utilintstr.IntOrString `json:"minAvailable,omitempty"`
	// Label query over pods whose evictions are managed by the disruption
	// budget.
	// A null selector will match no pods, while an empty ({}) selector will select
	// all pods within the namespace.
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
	// An eviction is allowed if at most "maxUnavailable" pods selected by
	// "selector" are unavailable after the eviction, i.e. even in absence of
	// the evicted pod. For example, one can prevent all voluntary evictions
	// by specifying 0. This is a mutually exclusive setting with "minAvailable".
	MaxUnavailable *utilintstr.IntOrString `json:"maxUnavailable,omitempty"`
	// UnhealthyPodEvictionPolicy defines the criteria for when unhealthy pods
	// should be considered for eviction. Current implementation considers healthy pods,
	// as pods that have status.conditions item with type="Ready",status="True".
	// Valid policies are IfHealthyBudget and AlwaysAllow.
	// If no policy is specified, the default behavior will be used,
	// which corresponds to the IfHealthyBudget policy.
	// IfHealthyBudget policy means that running pods (status.phase="Running"),
	// but not yet healthy can be evicted only if the guarded application is not
	// disrupted (status.currentHealthy is at least equal to status.desiredHealthy).
	// Healthy pods will be subject to the PDB for eviction.
	// AlwaysAllow policy means that all running pods (status.phase="Running"),
	// but not yet healthy are considered disrupted and can be evicted regardless
	// of whether the criteria in a PDB is met. This means perspective running
	// pods of a disrupted application might not get a chance to become healthy.
	// Healthy pods will be subject to the PDB for eviction.
	// Additional policies may be added in the future.
	// Clients making eviction decisions should disallow eviction of unhealthy pods
	// if they encounter an unrecognized policy in this field.
	// This field is beta-level. The eviction API uses this field when
	// the feature gate PDBUnhealthyPodEvictionPolicy is enabled (enabled by default).
	UnhealthyPodEvictionPolicy UnhealthyPodEvictionPolicyType `json:"unhealthyPodEvictionPolicy,omitempty"`
}

func (in *PodDisruptionBudgetSpec) DeepCopyInto(out *PodDisruptionBudgetSpec) {
	*out = *in
	if in.MinAvailable != nil {
		in, out := &in.MinAvailable, &out.MinAvailable
		*out = new(utilintstr.IntOrString)
		*out = *in
	}
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.MaxUnavailable != nil {
		in, out := &in.MaxUnavailable, &out.MaxUnavailable
		*out = new(utilintstr.IntOrString)
		*out = *in
	}
}

func (in *PodDisruptionBudgetSpec) DeepCopy() *PodDisruptionBudgetSpec {
	if in == nil {
		return nil
	}
	out := new(PodDisruptionBudgetSpec)
	in.DeepCopyInto(out)
	return out
}

type PodDisruptionBudgetStatus struct {
	// Most recent generation observed when updating this PDB status. DisruptionsAllowed and other
	// status information is valid only if observedGeneration equals to PDB's object generation.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// DisruptedPods contains information about pods whose eviction was
	// processed by the API server eviction subresource handler but has not
	// yet been observed by the PodDisruptionBudget controller.
	// A pod will be in this map from the time when the API server processed the
	// eviction request to the time when the pod is seen by PDB controller
	// as having been marked for deletion (or after a timeout). The key in the map is the name of the pod
	// and the value is the time when the API server processed the eviction request. If
	// the deletion didn't occur and a pod is still there it will be removed from
	// the list automatically by PodDisruptionBudget controller after some time.
	// If everything goes smooth this map should be empty for the most of the time.
	// Large number of entries in the map may indicate problems with pod deletions.
	DisruptedPods map[string]metav1.Time `json:"disruptedPods,omitempty"`
	// Number of pod disruptions that are currently allowed.
	DisruptionsAllowed int `json:"disruptionsAllowed"`
	// current number of healthy pods
	CurrentHealthy int `json:"currentHealthy"`
	// minimum desired number of healthy pods
	DesiredHealthy int `json:"desiredHealthy"`
	// total number of pods counted by this disruption budget
	ExpectedPods int `json:"expectedPods"`
	// Conditions contain conditions for PDB. The disruption controller sets the
	// DisruptionAllowed condition. The following are known values for the reason field
	// (additional reasons could be added in the future):
	// - SyncFailed: The controller encountered an error and wasn't able to compute
	// the number of allowed disruptions. Therefore no disruptions are
	// allowed and the status of the condition will be False.
	// - InsufficientPods: The number of pods are either at or below the number
	// required by the PodDisruptionBudget. No disruptions are
	// allowed and the status of the condition will be False.
	// - SufficientPods: There are more pods than required by the PodDisruptionBudget.
	// The condition will be True, and the number of allowed
	// disruptions are provided by the disruptionsAllowed property.
	Conditions []metav1.Condition `json:"conditions"`
}

func (in *PodDisruptionBudgetStatus) DeepCopyInto(out *PodDisruptionBudgetStatus) {
	*out = *in
	if in.DisruptedPods != nil {
		in, out := &in.DisruptedPods, &out.DisruptedPods
		*out = make(map[string]metav1.Time, len(*in))
		for k, v := range *in {
			(*out)[k] = v
		}
	}
	if in.Conditions != nil {
		l := make([]metav1.Condition, len(in.Conditions))
		for i := range in.Conditions {
			in.Conditions[i].DeepCopyInto(&l[i])
		}
		out.Conditions = l
	}
}

func (in *PodDisruptionBudgetStatus) DeepCopy() *PodDisruptionBudgetStatus {
	if in == nil {
		return nil
	}
	out := new(PodDisruptionBudgetStatus)
	in.DeepCopyInto(out)
	return out
}
