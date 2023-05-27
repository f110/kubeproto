package autoscalingv2

import (
	corev1 "go.f110.dev/kubeproto/go/apis/corev1"
	metav1 "go.f110.dev/kubeproto/go/apis/metav1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "autoscaling"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v2"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v2"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&HorizontalPodAutoscaler{},
		&HorizontalPodAutoscalerList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type HPAScalingPolicyType string

const (
	HPAScalingPolicyTypePods    HPAScalingPolicyType = "Pods"
	HPAScalingPolicyTypePercent HPAScalingPolicyType = "Percent"
)

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

type MetricTargetType string

const (
	MetricTargetTypeUtilization  MetricTargetType = "Utilization"
	MetricTargetTypeValue        MetricTargetType = "Value"
	MetricTargetTypeAverageValue MetricTargetType = "AverageValue"
)

type ScalingPolicySelect string

const (
	ScalingPolicySelectMax      ScalingPolicySelect = "Max"
	ScalingPolicySelectMin      ScalingPolicySelect = "Min"
	ScalingPolicySelectDisabled ScalingPolicySelect = "Disabled"
)

type HorizontalPodAutoscaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// spec is the specification for the behaviour of the autoscaler.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status.
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

type HorizontalPodAutoscalerSpec struct {
	// scaleTargetRef points to the target resource to scale, and is used to the pods for which metrics
	// should be collected, as well as to actually change the replica count.
	ScaleTargetRef CrossVersionObjectReference `json:"scaleTargetRef"`
	// minReplicas is the lower limit for the number of replicas to which the autoscaler
	// can scale down.  It defaults to 1 pod.  minReplicas is allowed to be 0 if the
	// alpha feature gate HPAScaleToZero is enabled and at least one Object or External
	// metric is configured.  Scaling is active as long as at least one metric value is
	// available.
	MinReplicas int `json:"minReplicas,omitempty"`
	// maxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up.
	// It cannot be less that minReplicas.
	MaxReplicas int `json:"maxReplicas"`
	// metrics contains the specifications for which to use to calculate the
	// desired replica count (the maximum replica count across all metrics will
	// be used).  The desired replica count is calculated multiplying the
	// ratio between the target value and the current value by the current
	// number of pods.  Ergo, metrics used must decrease as the pod count is
	// increased, and vice-versa.  See the individual metric source types for
	// more information about how each type of metric must respond.
	// If not set, the default metric will be set to 80% average CPU utilization.
	Metrics []MetricSpec `json:"metrics"`
	// behavior configures the scaling behavior of the target
	// in both Up and Down directions (scaleUp and scaleDown fields respectively).
	// If not set, the default HPAScalingRules for scale up and scale down are used.
	Behavior *HorizontalPodAutoscalerBehavior `json:"behavior,omitempty"`
}

func (in *HorizontalPodAutoscalerSpec) DeepCopyInto(out *HorizontalPodAutoscalerSpec) {
	*out = *in
	in.ScaleTargetRef.DeepCopyInto(&out.ScaleTargetRef)
	if in.Metrics != nil {
		l := make([]MetricSpec, len(in.Metrics))
		for i := range in.Metrics {
			in.Metrics[i].DeepCopyInto(&l[i])
		}
		out.Metrics = l
	}
	if in.Behavior != nil {
		in, out := &in.Behavior, &out.Behavior
		*out = new(HorizontalPodAutoscalerBehavior)
		(*in).DeepCopyInto(*out)
	}
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
	// lastScaleTime is the last time the HorizontalPodAutoscaler scaled the number of pods,
	// used by the autoscaler to control how often the number of pods is changed.
	LastScaleTime *metav1.Time `json:"lastScaleTime,omitempty"`
	// currentReplicas is current number of replicas of pods managed by this autoscaler,
	// as last seen by the autoscaler.
	CurrentReplicas int `json:"currentReplicas,omitempty"`
	// desiredReplicas is the desired number of replicas of pods managed by this autoscaler,
	// as last calculated by the autoscaler.
	DesiredReplicas int `json:"desiredReplicas"`
	// currentMetrics is the last read state of the metrics used by this autoscaler.
	CurrentMetrics []MetricStatus `json:"currentMetrics"`
	// conditions is the set of conditions required for this autoscaler to scale its target,
	// and indicates whether or not those conditions are met.
	Conditions []HorizontalPodAutoscalerCondition `json:"conditions"`
}

func (in *HorizontalPodAutoscalerStatus) DeepCopyInto(out *HorizontalPodAutoscalerStatus) {
	*out = *in
	if in.LastScaleTime != nil {
		in, out := &in.LastScaleTime, &out.LastScaleTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
	if in.CurrentMetrics != nil {
		l := make([]MetricStatus, len(in.CurrentMetrics))
		for i := range in.CurrentMetrics {
			in.CurrentMetrics[i].DeepCopyInto(&l[i])
		}
		out.CurrentMetrics = l
	}
	if in.Conditions != nil {
		l := make([]HorizontalPodAutoscalerCondition, len(in.Conditions))
		for i := range in.Conditions {
			in.Conditions[i].DeepCopyInto(&l[i])
		}
		out.Conditions = l
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

type MetricSpec struct {
	// type is the type of metric source.  It should be one of "ContainerResource", "External",
	// "Object", "Pods" or "Resource", each mapping to a matching field in the object.
	// Note: "ContainerResource" type is available on when the feature-gate
	// HPAContainerMetrics is enabled
	Type MetricSourceType `json:"type"`
	// object refers to a metric describing a single kubernetes object
	// (for example, hits-per-second on an Ingress object).
	Object *ObjectMetricSource `json:"object,omitempty"`
	// pods refers to a metric describing each pod in the current scale target
	// (for example, transactions-processed-per-second).  The values will be
	// averaged together before being compared to the target value.
	Pods *PodsMetricSource `json:"pods,omitempty"`
	// resource refers to a resource metric (such as those specified in
	// requests and limits) known to Kubernetes describing each pod in the
	// current scale target (e.g. CPU or memory). Such metrics are built in to
	// Kubernetes, and have special scaling options on top of those available
	// to normal per-pod metrics using the "pods" source.
	Resource *ResourceMetricSource `json:"resource,omitempty"`
	// containerResource refers to a resource metric (such as those specified in
	// requests and limits) known to Kubernetes describing a single container in
	// each pod of the current scale target (e.g. CPU or memory). Such metrics are
	// built in to Kubernetes, and have special scaling options on top of those
	// available to normal per-pod metrics using the "pods" source.
	// This is an alpha feature and can be enabled by the HPAContainerMetrics feature flag.
	ContainerResource *ContainerResourceMetricSource `json:"containerResource,omitempty"`
	// external refers to a global metric that is not associated
	// with any Kubernetes object. It allows autoscaling based on information
	// coming from components running outside of cluster
	// (for example length of queue in cloud messaging service, or
	// QPS from loadbalancer running outside of cluster).
	External *ExternalMetricSource `json:"external,omitempty"`
}

func (in *MetricSpec) DeepCopyInto(out *MetricSpec) {
	*out = *in
	if in.Object != nil {
		in, out := &in.Object, &out.Object
		*out = new(ObjectMetricSource)
		(*in).DeepCopyInto(*out)
	}
	if in.Pods != nil {
		in, out := &in.Pods, &out.Pods
		*out = new(PodsMetricSource)
		(*in).DeepCopyInto(*out)
	}
	if in.Resource != nil {
		in, out := &in.Resource, &out.Resource
		*out = new(ResourceMetricSource)
		(*in).DeepCopyInto(*out)
	}
	if in.ContainerResource != nil {
		in, out := &in.ContainerResource, &out.ContainerResource
		*out = new(ContainerResourceMetricSource)
		(*in).DeepCopyInto(*out)
	}
	if in.External != nil {
		in, out := &in.External, &out.External
		*out = new(ExternalMetricSource)
		(*in).DeepCopyInto(*out)
	}
}

func (in *MetricSpec) DeepCopy() *MetricSpec {
	if in == nil {
		return nil
	}
	out := new(MetricSpec)
	in.DeepCopyInto(out)
	return out
}

type HorizontalPodAutoscalerBehavior struct {
	// scaleUp is scaling policy for scaling Up.
	// If not set, the default value is the higher of:
	// * increase no more than 4 pods per 60 seconds
	// * double the number of pods per 60 seconds
	// No stabilization is used.
	ScaleUp *HPAScalingRules `json:"scaleUp,omitempty"`
	// scaleDown is scaling policy for scaling Down.
	// If not set, the default value is to allow to scale down to minReplicas pods, with a
	// 300 second stabilization window (i.e., the highest recommendation for
	// the last 300sec is used).
	ScaleDown *HPAScalingRules `json:"scaleDown,omitempty"`
}

func (in *HorizontalPodAutoscalerBehavior) DeepCopyInto(out *HorizontalPodAutoscalerBehavior) {
	*out = *in
	if in.ScaleUp != nil {
		in, out := &in.ScaleUp, &out.ScaleUp
		*out = new(HPAScalingRules)
		(*in).DeepCopyInto(*out)
	}
	if in.ScaleDown != nil {
		in, out := &in.ScaleDown, &out.ScaleDown
		*out = new(HPAScalingRules)
		(*in).DeepCopyInto(*out)
	}
}

func (in *HorizontalPodAutoscalerBehavior) DeepCopy() *HorizontalPodAutoscalerBehavior {
	if in == nil {
		return nil
	}
	out := new(HorizontalPodAutoscalerBehavior)
	in.DeepCopyInto(out)
	return out
}

type MetricStatus struct {
	// type is the type of metric source.  It will be one of "ContainerResource", "External",
	// "Object", "Pods" or "Resource", each corresponds to a matching field in the object.
	// Note: "ContainerResource" type is available on when the feature-gate
	// HPAContainerMetrics is enabled
	Type MetricSourceType `json:"type"`
	// object refers to a metric describing a single kubernetes object
	// (for example, hits-per-second on an Ingress object).
	Object *ObjectMetricStatus `json:"object,omitempty"`
	// pods refers to a metric describing each pod in the current scale target
	// (for example, transactions-processed-per-second).  The values will be
	// averaged together before being compared to the target value.
	Pods *PodsMetricStatus `json:"pods,omitempty"`
	// resource refers to a resource metric (such as those specified in
	// requests and limits) known to Kubernetes describing each pod in the
	// current scale target (e.g. CPU or memory). Such metrics are built in to
	// Kubernetes, and have special scaling options on top of those available
	// to normal per-pod metrics using the "pods" source.
	Resource *ResourceMetricStatus `json:"resource,omitempty"`
	// container resource refers to a resource metric (such as those specified in
	// requests and limits) known to Kubernetes describing a single container in each pod in the
	// current scale target (e.g. CPU or memory). Such metrics are built in to
	// Kubernetes, and have special scaling options on top of those available
	// to normal per-pod metrics using the "pods" source.
	ContainerResource *ContainerResourceMetricStatus `json:"containerResource,omitempty"`
	// external refers to a global metric that is not associated
	// with any Kubernetes object. It allows autoscaling based on information
	// coming from components running outside of cluster
	// (for example length of queue in cloud messaging service, or
	// QPS from loadbalancer running outside of cluster).
	External *ExternalMetricStatus `json:"external,omitempty"`
}

func (in *MetricStatus) DeepCopyInto(out *MetricStatus) {
	*out = *in
	if in.Object != nil {
		in, out := &in.Object, &out.Object
		*out = new(ObjectMetricStatus)
		(*in).DeepCopyInto(*out)
	}
	if in.Pods != nil {
		in, out := &in.Pods, &out.Pods
		*out = new(PodsMetricStatus)
		(*in).DeepCopyInto(*out)
	}
	if in.Resource != nil {
		in, out := &in.Resource, &out.Resource
		*out = new(ResourceMetricStatus)
		(*in).DeepCopyInto(*out)
	}
	if in.ContainerResource != nil {
		in, out := &in.ContainerResource, &out.ContainerResource
		*out = new(ContainerResourceMetricStatus)
		(*in).DeepCopyInto(*out)
	}
	if in.External != nil {
		in, out := &in.External, &out.External
		*out = new(ExternalMetricStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *MetricStatus) DeepCopy() *MetricStatus {
	if in == nil {
		return nil
	}
	out := new(MetricStatus)
	in.DeepCopyInto(out)
	return out
}

type HorizontalPodAutoscalerCondition struct {
	// type describes the current condition
	Type HorizontalPodAutoscalerConditionType `json:"type"`
	// status is the status of the condition (True, False, Unknown)
	Status corev1.ConditionStatus `json:"status"`
	// lastTransitionTime is the last time the condition transitioned from
	// one status to another
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`
	// reason is the reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// message is a human-readable explanation containing details about
	// the transition
	Message string `json:"message,omitempty"`
}

func (in *HorizontalPodAutoscalerCondition) DeepCopyInto(out *HorizontalPodAutoscalerCondition) {
	*out = *in
	if in.LastTransitionTime != nil {
		in, out := &in.LastTransitionTime, &out.LastTransitionTime
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

func (in *HorizontalPodAutoscalerCondition) DeepCopy() *HorizontalPodAutoscalerCondition {
	if in == nil {
		return nil
	}
	out := new(HorizontalPodAutoscalerCondition)
	in.DeepCopyInto(out)
	return out
}

type ObjectMetricSource struct {
	// describedObject specifies the descriptions of a object,such as kind,name apiVersion
	DescribedObject CrossVersionObjectReference `json:"describedObject"`
	// target specifies the target value for the given metric
	Target MetricTarget `json:"target"`
	// metric identifies the target metric by name and selector
	Metric MetricIdentifier `json:"metric"`
}

func (in *ObjectMetricSource) DeepCopyInto(out *ObjectMetricSource) {
	*out = *in
	in.DescribedObject.DeepCopyInto(&out.DescribedObject)
	in.Target.DeepCopyInto(&out.Target)
	in.Metric.DeepCopyInto(&out.Metric)
}

func (in *ObjectMetricSource) DeepCopy() *ObjectMetricSource {
	if in == nil {
		return nil
	}
	out := new(ObjectMetricSource)
	in.DeepCopyInto(out)
	return out
}

type PodsMetricSource struct {
	// metric identifies the target metric by name and selector
	Metric MetricIdentifier `json:"metric"`
	// target specifies the target value for the given metric
	Target MetricTarget `json:"target"`
}

func (in *PodsMetricSource) DeepCopyInto(out *PodsMetricSource) {
	*out = *in
	in.Metric.DeepCopyInto(&out.Metric)
	in.Target.DeepCopyInto(&out.Target)
}

func (in *PodsMetricSource) DeepCopy() *PodsMetricSource {
	if in == nil {
		return nil
	}
	out := new(PodsMetricSource)
	in.DeepCopyInto(out)
	return out
}

type ResourceMetricSource struct {
	// name is the name of the resource in question.
	Name corev1.ResourceName `json:"name"`
	// target specifies the target value for the given metric
	Target MetricTarget `json:"target"`
}

func (in *ResourceMetricSource) DeepCopyInto(out *ResourceMetricSource) {
	*out = *in
	in.Target.DeepCopyInto(&out.Target)
}

func (in *ResourceMetricSource) DeepCopy() *ResourceMetricSource {
	if in == nil {
		return nil
	}
	out := new(ResourceMetricSource)
	in.DeepCopyInto(out)
	return out
}

type ContainerResourceMetricSource struct {
	// name is the name of the resource in question.
	Name corev1.ResourceName `json:"name"`
	// target specifies the target value for the given metric
	Target MetricTarget `json:"target"`
	// container is the name of the container in the pods of the scaling target
	Container string `json:"container"`
}

func (in *ContainerResourceMetricSource) DeepCopyInto(out *ContainerResourceMetricSource) {
	*out = *in
	in.Target.DeepCopyInto(&out.Target)
}

func (in *ContainerResourceMetricSource) DeepCopy() *ContainerResourceMetricSource {
	if in == nil {
		return nil
	}
	out := new(ContainerResourceMetricSource)
	in.DeepCopyInto(out)
	return out
}

type ExternalMetricSource struct {
	// metric identifies the target metric by name and selector
	Metric MetricIdentifier `json:"metric"`
	// target specifies the target value for the given metric
	Target MetricTarget `json:"target"`
}

func (in *ExternalMetricSource) DeepCopyInto(out *ExternalMetricSource) {
	*out = *in
	in.Metric.DeepCopyInto(&out.Metric)
	in.Target.DeepCopyInto(&out.Target)
}

func (in *ExternalMetricSource) DeepCopy() *ExternalMetricSource {
	if in == nil {
		return nil
	}
	out := new(ExternalMetricSource)
	in.DeepCopyInto(out)
	return out
}

type HPAScalingRules struct {
	// stabilizationWindowSeconds is the number of seconds for which past recommendations should be
	// considered while scaling up or scaling down.
	// StabilizationWindowSeconds must be greater than or equal to zero and less than or equal to 3600 (one hour).
	// If not set, use the default values:
	// - For scale up: 0 (i.e. no stabilization is done).
	// - For scale down: 300 (i.e. the stabilization window is 300 seconds long).
	StabilizationWindowSeconds int `json:"stabilizationWindowSeconds,omitempty"`
	// selectPolicy is used to specify which policy should be used.
	// If not set, the default value Max is used.
	SelectPolicy ScalingPolicySelect `json:"selectPolicy,omitempty"`
	// policies is a list of potential scaling polices which can be used during scaling.
	// At least one policy must be specified, otherwise the HPAScalingRules will be discarded as invalid
	Policies []HPAScalingPolicy `json:"policies"`
}

func (in *HPAScalingRules) DeepCopyInto(out *HPAScalingRules) {
	*out = *in
	if in.Policies != nil {
		l := make([]HPAScalingPolicy, len(in.Policies))
		for i := range in.Policies {
			in.Policies[i].DeepCopyInto(&l[i])
		}
		out.Policies = l
	}
}

func (in *HPAScalingRules) DeepCopy() *HPAScalingRules {
	if in == nil {
		return nil
	}
	out := new(HPAScalingRules)
	in.DeepCopyInto(out)
	return out
}

type ObjectMetricStatus struct {
	// metric identifies the target metric by name and selector
	Metric MetricIdentifier `json:"metric"`
	// current contains the current value for the given metric
	Current MetricValueStatus `json:"current"`
	// DescribedObject specifies the descriptions of a object,such as kind,name apiVersion
	DescribedObject CrossVersionObjectReference `json:"describedObject"`
}

func (in *ObjectMetricStatus) DeepCopyInto(out *ObjectMetricStatus) {
	*out = *in
	in.Metric.DeepCopyInto(&out.Metric)
	in.Current.DeepCopyInto(&out.Current)
	in.DescribedObject.DeepCopyInto(&out.DescribedObject)
}

func (in *ObjectMetricStatus) DeepCopy() *ObjectMetricStatus {
	if in == nil {
		return nil
	}
	out := new(ObjectMetricStatus)
	in.DeepCopyInto(out)
	return out
}

type PodsMetricStatus struct {
	// metric identifies the target metric by name and selector
	Metric MetricIdentifier `json:"metric"`
	// current contains the current value for the given metric
	Current MetricValueStatus `json:"current"`
}

func (in *PodsMetricStatus) DeepCopyInto(out *PodsMetricStatus) {
	*out = *in
	in.Metric.DeepCopyInto(&out.Metric)
	in.Current.DeepCopyInto(&out.Current)
}

func (in *PodsMetricStatus) DeepCopy() *PodsMetricStatus {
	if in == nil {
		return nil
	}
	out := new(PodsMetricStatus)
	in.DeepCopyInto(out)
	return out
}

type ResourceMetricStatus struct {
	// name is the name of the resource in question.
	Name corev1.ResourceName `json:"name"`
	// current contains the current value for the given metric
	Current MetricValueStatus `json:"current"`
}

func (in *ResourceMetricStatus) DeepCopyInto(out *ResourceMetricStatus) {
	*out = *in
	in.Current.DeepCopyInto(&out.Current)
}

func (in *ResourceMetricStatus) DeepCopy() *ResourceMetricStatus {
	if in == nil {
		return nil
	}
	out := new(ResourceMetricStatus)
	in.DeepCopyInto(out)
	return out
}

type ContainerResourceMetricStatus struct {
	// name is the name of the resource in question.
	Name corev1.ResourceName `json:"name"`
	// current contains the current value for the given metric
	Current MetricValueStatus `json:"current"`
	// container is the name of the container in the pods of the scaling target
	Container string `json:"container"`
}

func (in *ContainerResourceMetricStatus) DeepCopyInto(out *ContainerResourceMetricStatus) {
	*out = *in
	in.Current.DeepCopyInto(&out.Current)
}

func (in *ContainerResourceMetricStatus) DeepCopy() *ContainerResourceMetricStatus {
	if in == nil {
		return nil
	}
	out := new(ContainerResourceMetricStatus)
	in.DeepCopyInto(out)
	return out
}

type ExternalMetricStatus struct {
	// metric identifies the target metric by name and selector
	Metric MetricIdentifier `json:"metric"`
	// current contains the current value for the given metric
	Current MetricValueStatus `json:"current"`
}

func (in *ExternalMetricStatus) DeepCopyInto(out *ExternalMetricStatus) {
	*out = *in
	in.Metric.DeepCopyInto(&out.Metric)
	in.Current.DeepCopyInto(&out.Current)
}

func (in *ExternalMetricStatus) DeepCopy() *ExternalMetricStatus {
	if in == nil {
		return nil
	}
	out := new(ExternalMetricStatus)
	in.DeepCopyInto(out)
	return out
}

type MetricTarget struct {
	// type represents whether the metric type is Utilization, Value, or AverageValue
	Type MetricTargetType `json:"type"`
	// value is the target value of the metric (as a quantity).
	Value *apiresource.Quantity `json:"value,omitempty"`
	// averageValue is the target value of the average of the
	// metric across all relevant pods (as a quantity)
	AverageValue *apiresource.Quantity `json:"averageValue,omitempty"`
	// averageUtilization is the target value of the average of the
	// resource metric across all relevant pods, represented as a percentage of
	// the requested value of the resource for the pods.
	// Currently only valid for Resource metric source type
	AverageUtilization int `json:"averageUtilization,omitempty"`
}

func (in *MetricTarget) DeepCopyInto(out *MetricTarget) {
	*out = *in
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = new(apiresource.Quantity)
		(*in).DeepCopyInto(*out)
	}
	if in.AverageValue != nil {
		in, out := &in.AverageValue, &out.AverageValue
		*out = new(apiresource.Quantity)
		(*in).DeepCopyInto(*out)
	}
}

func (in *MetricTarget) DeepCopy() *MetricTarget {
	if in == nil {
		return nil
	}
	out := new(MetricTarget)
	in.DeepCopyInto(out)
	return out
}

type MetricIdentifier struct {
	// name is the name of the given metric
	Name string `json:"name"`
	// selector is the string-encoded form of a standard kubernetes label selector for the given metric
	// When set, it is passed as an additional parameter to the metrics server for more specific metrics scoping.
	// When unset, just the metricName will be used to gather metrics.
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

func (in *MetricIdentifier) DeepCopyInto(out *MetricIdentifier) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

func (in *MetricIdentifier) DeepCopy() *MetricIdentifier {
	if in == nil {
		return nil
	}
	out := new(MetricIdentifier)
	in.DeepCopyInto(out)
	return out
}

type HPAScalingPolicy struct {
	// type is used to specify the scaling policy.
	Type HPAScalingPolicyType `json:"type"`
	// value contains the amount of change which is permitted by the policy.
	// It must be greater than zero
	Value int `json:"value"`
	// periodSeconds specifies the window of time for which the policy should hold true.
	// PeriodSeconds must be greater than zero and less than or equal to 1800 (30 min).
	PeriodSeconds int `json:"periodSeconds"`
}

func (in *HPAScalingPolicy) DeepCopyInto(out *HPAScalingPolicy) {
	*out = *in
}

func (in *HPAScalingPolicy) DeepCopy() *HPAScalingPolicy {
	if in == nil {
		return nil
	}
	out := new(HPAScalingPolicy)
	in.DeepCopyInto(out)
	return out
}

type MetricValueStatus struct {
	// value is the current value of the metric (as a quantity).
	Value *apiresource.Quantity `json:"value,omitempty"`
	// averageValue is the current value of the average of the
	// metric across all relevant pods (as a quantity)
	AverageValue *apiresource.Quantity `json:"averageValue,omitempty"`
	// currentAverageUtilization is the current value of the average of the
	// resource metric across all relevant pods, represented as a percentage of
	// the requested value of the resource for the pods.
	AverageUtilization int `json:"averageUtilization,omitempty"`
}

func (in *MetricValueStatus) DeepCopyInto(out *MetricValueStatus) {
	*out = *in
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = new(apiresource.Quantity)
		(*in).DeepCopyInto(*out)
	}
	if in.AverageValue != nil {
		in, out := &in.AverageValue, &out.AverageValue
		*out = new(apiresource.Quantity)
		(*in).DeepCopyInto(*out)
	}
}

func (in *MetricValueStatus) DeepCopy() *MetricValueStatus {
	if in == nil {
		return nil
	}
	out := new(MetricValueStatus)
	in.DeepCopyInto(out)
	return out
}
