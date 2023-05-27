package autoscalingv1

import (
	metav1 "go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "autoscaling"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&HorizontalPodAutoscaler{},
		&HorizontalPodAutoscalerList{},
		&Scale{},
		&ScaleList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type HorizontalPodAutoscalerConditionType string

const (
	HorizontalPodAutoscalerConditionTypeScalingActive  HorizontalPodAutoscalerConditionType = "ScalingActive"
	HorizontalPodAutoscalerConditionTypeAbleToScale    HorizontalPodAutoscalerConditionType = "AbleToScale"
	HorizontalPodAutoscalerConditionTypeScalingLimited HorizontalPodAutoscalerConditionType = "ScalingLimited"
)

type MetricSourceType string

const (
	MetricSourceTypeObject            MetricSourceType = "Object"
	MetricSourceTypePods              MetricSourceType = "Pods"
	MetricSourceTypeResource          MetricSourceType = "Resource"
	MetricSourceTypeContainerResource MetricSourceType = "ContainerResource"
	MetricSourceTypeExternal          MetricSourceType = "External"
)

type HorizontalPodAutoscaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// spec defines the behaviour of autoscaler. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status.
	Spec *HorizontalPodAutoscalerSpec `json:"spec,omitempty"`
	// status is the current information about the autoscaler.
	Status *HorizontalPodAutoscalerStatus `json:"status,omitempty"`
}

func (in *HorizontalPodAutoscaler) DeepCopyInto(out *HorizontalPodAutoscaler) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(HorizontalPodAutoscalerSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(HorizontalPodAutoscalerStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *HorizontalPodAutoscaler) DeepCopy() *HorizontalPodAutoscaler {
	if in == nil {
		return nil
	}
	out := new(HorizontalPodAutoscaler)
	in.DeepCopyInto(out)
	return out
}

func (in *HorizontalPodAutoscaler) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type HorizontalPodAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []HorizontalPodAutoscaler `json:"items"`
}

func (in *HorizontalPodAutoscalerList) DeepCopyInto(out *HorizontalPodAutoscalerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]HorizontalPodAutoscaler, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *HorizontalPodAutoscalerList) DeepCopy() *HorizontalPodAutoscalerList {
	if in == nil {
		return nil
	}
	out := new(HorizontalPodAutoscalerList)
	in.DeepCopyInto(out)
	return out
}

func (in *HorizontalPodAutoscalerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type Scale struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// spec defines the behavior of the scale. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status.
	Spec *ScaleSpec `json:"spec,omitempty"`
	// status is the current status of the scale. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status. Read-only.
	Status *ScaleStatus `json:"status,omitempty"`
}

func (in *Scale) DeepCopyInto(out *Scale) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(ScaleSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(ScaleStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *Scale) DeepCopy() *Scale {
	if in == nil {
		return nil
	}
	out := new(Scale)
	in.DeepCopyInto(out)
	return out
}

func (in *Scale) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type ScaleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Scale `json:"items"`
}

func (in *ScaleList) DeepCopyInto(out *ScaleList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]Scale, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *ScaleList) DeepCopy() *ScaleList {
	if in == nil {
		return nil
	}
	out := new(ScaleList)
	in.DeepCopyInto(out)
	return out
}

func (in *ScaleList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type HorizontalPodAutoscalerSpec struct {
	// reference to scaled resource; horizontal pod autoscaler will learn the current resource consumption
	// and will set the desired number of pods by using its Scale subresource.
	ScaleTargetRef CrossVersionObjectReference `json:"scaleTargetRef"`
	// minReplicas is the lower limit for the number of replicas to which the autoscaler
	// can scale down.  It defaults to 1 pod.  minReplicas is allowed to be 0 if the
	// alpha feature gate HPAScaleToZero is enabled and at least one Object or External
	// metric is configured.  Scaling is active as long as at least one metric value is
	// available.
	MinReplicas int `json:"minReplicas,omitempty"`
	// maxReplicas is the upper limit for the number of pods that can be set by the autoscaler; cannot be smaller than MinReplicas.
	MaxReplicas int `json:"maxReplicas"`
	// targetCPUUtilizationPercentage is the target average CPU utilization (represented as a percentage of requested CPU) over all the pods;
	// if not specified the default autoscaling policy will be used.
	TargetCPUUtilizationPercentage int `json:"targetCPUUtilizationPercentage,omitempty"`
}

func (in *HorizontalPodAutoscalerSpec) DeepCopyInto(out *HorizontalPodAutoscalerSpec) {
	*out = *in
	in.ScaleTargetRef.DeepCopyInto(&out.ScaleTargetRef)
}

func (in *HorizontalPodAutoscalerSpec) DeepCopy() *HorizontalPodAutoscalerSpec {
	if in == nil {
		return nil
	}
	out := new(HorizontalPodAutoscalerSpec)
	in.DeepCopyInto(out)
	return out
}

type HorizontalPodAutoscalerStatus struct {
	// observedGeneration is the most recent generation observed by this autoscaler.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// lastScaleTime is the last time the HorizontalPodAutoscaler scaled the number of pods;
	// used by the autoscaler to control how often the number of pods is changed.
	LastScaleTime *metav1.Time `json:"lastScaleTime,omitempty"`
	// currentReplicas is the current number of replicas of pods managed by this autoscaler.
	CurrentReplicas int `json:"currentReplicas"`
	// desiredReplicas is the  desired number of replicas of pods managed by this autoscaler.
	DesiredReplicas int `json:"desiredReplicas"`
	// currentCPUUtilizationPercentage is the current average CPU utilization over all pods, represented as a percentage of requested CPU,
	// e.g. 70 means that an average pod is using now 70% of its requested CPU.
	CurrentCPUUtilizationPercentage int `json:"currentCPUUtilizationPercentage,omitempty"`
}

func (in *HorizontalPodAutoscalerStatus) DeepCopyInto(out *HorizontalPodAutoscalerStatus) {
	*out = *in
	if in.LastScaleTime != nil {
		in, out := &in.LastScaleTime, &out.LastScaleTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

func (in *HorizontalPodAutoscalerStatus) DeepCopy() *HorizontalPodAutoscalerStatus {
	if in == nil {
		return nil
	}
	out := new(HorizontalPodAutoscalerStatus)
	in.DeepCopyInto(out)
	return out
}

type ScaleSpec struct {
	// replicas is the desired number of instances for the scaled object.
	Replicas int `json:"replicas,omitempty"`
}

func (in *ScaleSpec) DeepCopyInto(out *ScaleSpec) {
	*out = *in
}

func (in *ScaleSpec) DeepCopy() *ScaleSpec {
	if in == nil {
		return nil
	}
	out := new(ScaleSpec)
	in.DeepCopyInto(out)
	return out
}

type ScaleStatus struct {
	// replicas is the actual number of observed instances of the scaled object.
	Replicas int `json:"replicas"`
	// selector is the label query over pods that should match the replicas count. This is same
	// as the label selector but in the string format to avoid introspection
	// by clients. The string will be in the same format as the query-param syntax.
	// More info about label selectors: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/
	Selector string `json:"selector,omitempty"`
}

func (in *ScaleStatus) DeepCopyInto(out *ScaleStatus) {
	*out = *in
}

func (in *ScaleStatus) DeepCopy() *ScaleStatus {
	if in == nil {
		return nil
	}
	out := new(ScaleStatus)
	in.DeepCopyInto(out)
	return out
}

type CrossVersionObjectReference struct {
	// kind is the kind of the referent; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind string `json:"kind"`
	// name is the name of the referent; More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name"`
	// apiVersion is the API version of the referent
	APIVersion string `json:"apiVersion,omitempty"`
}

func (in *CrossVersionObjectReference) DeepCopyInto(out *CrossVersionObjectReference) {
	*out = *in
}

func (in *CrossVersionObjectReference) DeepCopy() *CrossVersionObjectReference {
	if in == nil {
		return nil
	}
	out := new(CrossVersionObjectReference)
	in.DeepCopyInto(out)
	return out
}
