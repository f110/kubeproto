package resourcev1

import (
	"go.f110.dev/kubeproto/go/apis/corev1"
	"go.f110.dev/kubeproto/go/apis/metav1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "resource"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&DeviceClass{},
		&DeviceClassList{},
		&ResourceClaim{},
		&ResourceClaimList{},
		&ResourceClaimTemplate{},
		&ResourceClaimTemplateList{},
		&ResourceSlice{},
		&ResourceSliceList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type AllocationConfigSource string

const (
	AllocationConfigSourceFromClass AllocationConfigSource = "FromClass"
	AllocationConfigSourceFromClaim AllocationConfigSource = "FromClaim"
)

type DeviceAllocationMode string

const (
	DeviceAllocationModeExactCount DeviceAllocationMode = "ExactCount"
	DeviceAllocationModeAll        DeviceAllocationMode = "All"
)

type DeviceTaintEffect string

const (
	DeviceTaintEffectNone       DeviceTaintEffect = "None"
	DeviceTaintEffectNoSchedule DeviceTaintEffect = "NoSchedule"
	DeviceTaintEffectNoExecute  DeviceTaintEffect = "NoExecute"
)

type DeviceTolerationOperator string

const (
	DeviceTolerationOperatorExists DeviceTolerationOperator = "Exists"
	DeviceTolerationOperatorEqual  DeviceTolerationOperator = "Equal"
)

type DeviceClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec defines what can be allocated and how to configure it.
	// This is mutable. Consumers have to be prepared for classes changing
	// at any time, either because they get updated or replaced. Claim
	// allocations are done once based on whatever was set in classes at
	// the time of allocation.
	// Changing the spec automatically increments the metadata.generation number.
	Spec DeviceClassSpec `json:"spec"`
}

func (in *DeviceClass) DeepCopyInto(out *DeviceClass) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

func (in *DeviceClass) DeepCopy() *DeviceClass {
	if in == nil {
		return nil
	}
	out := new(DeviceClass)
	in.DeepCopyInto(out)
	return out
}

func (in *DeviceClass) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type DeviceClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []DeviceClass `json:"items"`
}

func (in *DeviceClassList) DeepCopyInto(out *DeviceClassList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]DeviceClass, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *DeviceClassList) DeepCopy() *DeviceClassList {
	if in == nil {
		return nil
	}
	out := new(DeviceClassList)
	in.DeepCopyInto(out)
	return out
}

func (in *DeviceClassList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type ResourceClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Spec describes what is being requested and how to configure it.
	// The spec is immutable.
	Spec ResourceClaimSpec `json:"spec"`
	// Status describes whether the claim is ready to use and what has been allocated.
	Status *ResourceClaimStatus `json:"status,omitempty"`
}

func (in *ResourceClaim) DeepCopyInto(out *ResourceClaim) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(ResourceClaimStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *ResourceClaim) DeepCopy() *ResourceClaim {
	if in == nil {
		return nil
	}
	out := new(ResourceClaim)
	in.DeepCopyInto(out)
	return out
}

func (in *ResourceClaim) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type ResourceClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ResourceClaim `json:"items"`
}

func (in *ResourceClaimList) DeepCopyInto(out *ResourceClaimList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]ResourceClaim, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *ResourceClaimList) DeepCopy() *ResourceClaimList {
	if in == nil {
		return nil
	}
	out := new(ResourceClaimList)
	in.DeepCopyInto(out)
	return out
}

func (in *ResourceClaimList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type ResourceClaimTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Describes the ResourceClaim that is to be generated.
	// This field is immutable. A ResourceClaim will get created by the
	// control plane for a Pod when needed and then not get updated
	// anymore.
	Spec ResourceClaimTemplateSpec `json:"spec"`
}

func (in *ResourceClaimTemplate) DeepCopyInto(out *ResourceClaimTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

func (in *ResourceClaimTemplate) DeepCopy() *ResourceClaimTemplate {
	if in == nil {
		return nil
	}
	out := new(ResourceClaimTemplate)
	in.DeepCopyInto(out)
	return out
}

func (in *ResourceClaimTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type ResourceClaimTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ResourceClaimTemplate `json:"items"`
}

func (in *ResourceClaimTemplateList) DeepCopyInto(out *ResourceClaimTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]ResourceClaimTemplate, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *ResourceClaimTemplateList) DeepCopy() *ResourceClaimTemplateList {
	if in == nil {
		return nil
	}
	out := new(ResourceClaimTemplateList)
	in.DeepCopyInto(out)
	return out
}

func (in *ResourceClaimTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type ResourceSlice struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// Contains the information published by the driver.
	// Changing the spec automatically increments the metadata.generation number.
	Spec ResourceSliceSpec `json:"spec"`
}

func (in *ResourceSlice) DeepCopyInto(out *ResourceSlice) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

func (in *ResourceSlice) DeepCopy() *ResourceSlice {
	if in == nil {
		return nil
	}
	out := new(ResourceSlice)
	in.DeepCopyInto(out)
	return out
}

func (in *ResourceSlice) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type ResourceSliceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ResourceSlice `json:"items"`
}

func (in *ResourceSliceList) DeepCopyInto(out *ResourceSliceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]ResourceSlice, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *ResourceSliceList) DeepCopy() *ResourceSliceList {
	if in == nil {
		return nil
	}
	out := new(ResourceSliceList)
	in.DeepCopyInto(out)
	return out
}

func (in *ResourceSliceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type DeviceClassSpec struct {
	// Each selector must be satisfied by a device which is claimed via this class.
	Selectors []DeviceSelector `json:"selectors"`
	// Config defines configuration parameters that apply to each device that is claimed via this class.
	// Some classses may potentially be satisfied by multiple drivers, so each instance of a vendor
	// configuration applies to exactly one driver.
	// They are passed to the driver, but are not considered while allocating the claim.
	Config []DeviceClassConfiguration `json:"config"`
	// ExtendedResourceName is the extended resource name for the devices of this class.
	// The devices of this class can be used to satisfy a pod's extended resource requests.
	// It has the same format as the name of a pod's extended resource.
	// It should be unique among all the device classes in a cluster.
	// If two device classes have the same name, then the class created later
	// is picked to satisfy a pod's extended resource requests.
	// If two classes are created at the same time, then the name of the class
	// lexicographically sorted first is picked.
	// This is a beta field.
	ExtendedResourceName string `json:"extendedResourceName,omitempty"`
}

func (in *DeviceClassSpec) DeepCopyInto(out *DeviceClassSpec) {
	*out = *in
	if in.Selectors != nil {
		l := make([]DeviceSelector, len(in.Selectors))
		for i := range in.Selectors {
			in.Selectors[i].DeepCopyInto(&l[i])
		}
		out.Selectors = l
	}
	if in.Config != nil {
		l := make([]DeviceClassConfiguration, len(in.Config))
		for i := range in.Config {
			in.Config[i].DeepCopyInto(&l[i])
		}
		out.Config = l
	}
}

func (in *DeviceClassSpec) DeepCopy() *DeviceClassSpec {
	if in == nil {
		return nil
	}
	out := new(DeviceClassSpec)
	in.DeepCopyInto(out)
	return out
}

type ResourceClaimSpec struct {
	// Devices defines how to request devices.
	Devices DeviceClaim `json:"devices"`
}

func (in *ResourceClaimSpec) DeepCopyInto(out *ResourceClaimSpec) {
	*out = *in
	in.Devices.DeepCopyInto(&out.Devices)
}

func (in *ResourceClaimSpec) DeepCopy() *ResourceClaimSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceClaimSpec)
	in.DeepCopyInto(out)
	return out
}

type ResourceClaimStatus struct {
	// Allocation is set once the claim has been allocated successfully.
	Allocation *AllocationResult `json:"allocation,omitempty"`
	// ReservedFor indicates which entities are currently allowed to use
	// the claim. A Pod which references a ResourceClaim which is not
	// reserved for that Pod will not be started. A claim that is in
	// use or might be in use because it has been reserved must not get
	// deallocated.
	// In a cluster with multiple scheduler instances, two pods might get
	// scheduled concurrently by different schedulers. When they reference
	// the same ResourceClaim which already has reached its maximum number
	// of consumers, only one pod can be scheduled.
	// Both schedulers try to add their pod to the claim.status.reservedFor
	// field, but only the update that reaches the API server first gets
	// stored. The other one fails with an error and the scheduler
	// which issued it knows that it must put the pod back into the queue,
	// waiting for the ResourceClaim to become usable again.
	// There can be at most 256 such reservations. This may get increased in
	// the future, but not reduced.
	ReservedFor []ResourceClaimConsumerReference `json:"reservedFor"`
	// Devices contains the status of each device allocated for this
	// claim, as reported by the driver. This can include driver-specific
	// information. Entries are owned by their respective drivers.
	Devices []AllocatedDeviceStatus `json:"devices"`
}

func (in *ResourceClaimStatus) DeepCopyInto(out *ResourceClaimStatus) {
	*out = *in
	if in.Allocation != nil {
		in, out := &in.Allocation, &out.Allocation
		*out = new(AllocationResult)
		(*in).DeepCopyInto(*out)
	}
	if in.ReservedFor != nil {
		l := make([]ResourceClaimConsumerReference, len(in.ReservedFor))
		for i := range in.ReservedFor {
			in.ReservedFor[i].DeepCopyInto(&l[i])
		}
		out.ReservedFor = l
	}
	if in.Devices != nil {
		l := make([]AllocatedDeviceStatus, len(in.Devices))
		for i := range in.Devices {
			in.Devices[i].DeepCopyInto(&l[i])
		}
		out.Devices = l
	}
}

func (in *ResourceClaimStatus) DeepCopy() *ResourceClaimStatus {
	if in == nil {
		return nil
	}
	out := new(ResourceClaimStatus)
	in.DeepCopyInto(out)
	return out
}

type ResourceClaimTemplateSpec struct {
	// ObjectMeta may contain labels and annotations that will be copied into the ResourceClaim
	// when creating it. No other fields are allowed and will be rejected during
	// validation.
	ObjectMeta *metav1.ObjectMeta `json:"metadata,omitempty"`
	// Spec for the ResourceClaim. The entire content is copied unchanged
	// into the ResourceClaim that gets created from this template. The
	// same fields as in a ResourceClaim are also valid here.
	Spec ResourceClaimSpec `json:"spec"`
}

func (in *ResourceClaimTemplateSpec) DeepCopyInto(out *ResourceClaimTemplateSpec) {
	*out = *in
	if in.ObjectMeta != nil {
		in, out := &in.ObjectMeta, &out.ObjectMeta
		*out = new(metav1.ObjectMeta)
		(*in).DeepCopyInto(*out)
	}
	in.Spec.DeepCopyInto(&out.Spec)
}

func (in *ResourceClaimTemplateSpec) DeepCopy() *ResourceClaimTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceClaimTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

type ResourceSliceSpec struct {
	// Driver identifies the DRA driver providing the capacity information.
	// A field selector can be used to list only ResourceSlice
	// objects with a certain driver name.
	// Must be a DNS subdomain and should end with a DNS domain owned by the
	// vendor of the driver. It should use only lower case characters.
	// This field is immutable.
	Driver string `json:"driver"`
	// Pool describes the pool that this ResourceSlice belongs to.
	Pool ResourcePool `json:"pool"`
	// NodeName identifies the node which provides the resources in this pool.
	// A field selector can be used to list only ResourceSlice
	// objects belonging to a certain node.
	// This field can be used to limit access from nodes to ResourceSlices with
	// the same node name. It also indicates to autoscalers that adding
	// new nodes of the same type as some old node might also make new
	// resources available.
	// Exactly one of NodeName, NodeSelector, AllNodes, and PerDeviceNodeSelection must be set.
	// This field is immutable.
	NodeName string `json:"nodeName,omitempty"`
	// NodeSelector defines which nodes have access to the resources in the pool,
	// when that pool is not limited to a single node.
	// Must use exactly one term.
	// Exactly one of NodeName, NodeSelector, AllNodes, and PerDeviceNodeSelection must be set.
	NodeSelector *corev1.NodeSelector `json:"nodeSelector,omitempty"`
	// AllNodes indicates that all nodes have access to the resources in the pool.
	// Exactly one of NodeName, NodeSelector, AllNodes, and PerDeviceNodeSelection must be set.
	AllNodes bool `json:"allNodes,omitempty"`
	// Devices lists some or all of the devices in this pool.
	// Must not have more than 128 entries. If any device uses taints or consumes counters the limit is 64.
	// Only one of Devices and SharedCounters can be set in a ResourceSlice.
	Devices []Device `json:"devices"`
	// PerDeviceNodeSelection defines whether the access from nodes to
	// resources in the pool is set on the ResourceSlice level or on each
	// device. If it is set to true, every device defined the ResourceSlice
	// must specify this individually.
	// Exactly one of NodeName, NodeSelector, AllNodes, and PerDeviceNodeSelection must be set.
	PerDeviceNodeSelection bool `json:"perDeviceNodeSelection,omitempty"`
	// SharedCounters defines a list of counter sets, each of which
	// has a name and a list of counters available.
	// The names of the counter sets must be unique in the ResourcePool.
	// Only one of Devices and SharedCounters can be set in a ResourceSlice.
	// The maximum number of counter sets is 8.
	SharedCounters []CounterSet `json:"sharedCounters"`
}

func (in *ResourceSliceSpec) DeepCopyInto(out *ResourceSliceSpec) {
	*out = *in
	in.Pool.DeepCopyInto(&out.Pool)
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = new(corev1.NodeSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Devices != nil {
		l := make([]Device, len(in.Devices))
		for i := range in.Devices {
			in.Devices[i].DeepCopyInto(&l[i])
		}
		out.Devices = l
	}
	if in.SharedCounters != nil {
		l := make([]CounterSet, len(in.SharedCounters))
		for i := range in.SharedCounters {
			in.SharedCounters[i].DeepCopyInto(&l[i])
		}
		out.SharedCounters = l
	}
}

func (in *ResourceSliceSpec) DeepCopy() *ResourceSliceSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceSliceSpec)
	in.DeepCopyInto(out)
	return out
}

type DeviceSelector struct {
	// CEL contains a CEL expression for selecting a device.
	CEL *CELDeviceSelector `json:"cel,omitempty"`
}

func (in *DeviceSelector) DeepCopyInto(out *DeviceSelector) {
	*out = *in
	if in.CEL != nil {
		in, out := &in.CEL, &out.CEL
		*out = new(CELDeviceSelector)
		(*in).DeepCopyInto(*out)
	}
}

func (in *DeviceSelector) DeepCopy() *DeviceSelector {
	if in == nil {
		return nil
	}
	out := new(DeviceSelector)
	in.DeepCopyInto(out)
	return out
}

type DeviceClassConfiguration struct {
	DeviceConfiguration `json:",inline"`
}

func (in *DeviceClassConfiguration) DeepCopyInto(out *DeviceClassConfiguration) {
	*out = *in
	out.DeviceConfiguration = in.DeviceConfiguration
}

func (in *DeviceClassConfiguration) DeepCopy() *DeviceClassConfiguration {
	if in == nil {
		return nil
	}
	out := new(DeviceClassConfiguration)
	in.DeepCopyInto(out)
	return out
}

type DeviceClaim struct {
	// Requests represent individual requests for distinct devices which
	// must all be satisfied. If empty, nothing needs to be allocated.
	Requests []DeviceRequest `json:"requests"`
	// These constraints must be satisfied by the set of devices that get
	// allocated for the claim.
	Constraints []DeviceConstraint `json:"constraints"`
	// This field holds configuration for multiple potential drivers which
	// could satisfy requests in this claim. It is ignored while allocating
	// the claim.
	Config []DeviceClaimConfiguration `json:"config"`
}

func (in *DeviceClaim) DeepCopyInto(out *DeviceClaim) {
	*out = *in
	if in.Requests != nil {
		l := make([]DeviceRequest, len(in.Requests))
		for i := range in.Requests {
			in.Requests[i].DeepCopyInto(&l[i])
		}
		out.Requests = l
	}
	if in.Constraints != nil {
		l := make([]DeviceConstraint, len(in.Constraints))
		for i := range in.Constraints {
			in.Constraints[i].DeepCopyInto(&l[i])
		}
		out.Constraints = l
	}
	if in.Config != nil {
		l := make([]DeviceClaimConfiguration, len(in.Config))
		for i := range in.Config {
			in.Config[i].DeepCopyInto(&l[i])
		}
		out.Config = l
	}
}

func (in *DeviceClaim) DeepCopy() *DeviceClaim {
	if in == nil {
		return nil
	}
	out := new(DeviceClaim)
	in.DeepCopyInto(out)
	return out
}

type AllocationResult struct {
	// Devices is the result of allocating devices.
	Devices *DeviceAllocationResult `json:"devices,omitempty"`
	// NodeSelector defines where the allocated resources are available. If
	// unset, they are available everywhere.
	NodeSelector *corev1.NodeSelector `json:"nodeSelector,omitempty"`
	// AllocationTimestamp stores the time when the resources were allocated.
	// This field is not guaranteed to be set, in which case that time is unknown.
	// This is a beta field and requires enabling the DRADeviceBindingConditions and DRAResourceClaimDeviceStatus
	// feature gate.
	AllocationTimestamp *metav1.Time `json:"allocationTimestamp,omitempty"`
}

func (in *AllocationResult) DeepCopyInto(out *AllocationResult) {
	*out = *in
	if in.Devices != nil {
		in, out := &in.Devices, &out.Devices
		*out = new(DeviceAllocationResult)
		(*in).DeepCopyInto(*out)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = new(corev1.NodeSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.AllocationTimestamp != nil {
		in, out := &in.AllocationTimestamp, &out.AllocationTimestamp
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

func (in *AllocationResult) DeepCopy() *AllocationResult {
	if in == nil {
		return nil
	}
	out := new(AllocationResult)
	in.DeepCopyInto(out)
	return out
}

type ResourceClaimConsumerReference struct {
	// APIGroup is the group for the resource being referenced. It is
	// empty for the core API. This matches the group in the APIVersion
	// that is used when creating the resources.
	APIGroup string `json:"apiGroup,omitempty"`
	// Resource is the type of resource being referenced, for example "pods".
	Resource string `json:"resource"`
	// Name is the name of resource being referenced.
	Name string `json:"name"`
	// UID identifies exactly one incarnation of the resource.
	UID string `json:"uid"`
}

func (in *ResourceClaimConsumerReference) DeepCopyInto(out *ResourceClaimConsumerReference) {
	*out = *in
}

func (in *ResourceClaimConsumerReference) DeepCopy() *ResourceClaimConsumerReference {
	if in == nil {
		return nil
	}
	out := new(ResourceClaimConsumerReference)
	in.DeepCopyInto(out)
	return out
}

type AllocatedDeviceStatus struct {
	// Driver specifies the name of the DRA driver whose kubelet
	// plugin should be invoked to process the allocation once the claim is
	// needed on a node.
	// Must be a DNS subdomain and should end with a DNS domain owned by the
	// vendor of the driver. It should use only lower case characters.
	Driver string `json:"driver"`
	// This name together with the driver name and the device name field
	// identify which device was allocated (`<driver name>/<pool name>/<device name>`).
	// Must not be longer than 253 characters and may contain one or more
	// DNS sub-domains separated by slashes.
	Pool string `json:"pool"`
	// Device references one device instance via its name in the driver's
	// resource pool. It must be a DNS label.
	Device string `json:"device"`
	// ShareID uniquely identifies an individual allocation share of the device.
	ShareID string `json:"shareID,omitempty"`
	// Conditions contains the latest observation of the device's state.
	// If the device has been configured according to the class and claim
	// config references, the `Ready` condition should be True.
	// Must not contain more than 8 entries.
	Conditions []metav1.Condition `json:"conditions"`
	// Data contains arbitrary driver-specific data.
	// The length of the raw data must be smaller or equal to 10 Ki.
	Data *runtime.RawExtension `json:"data,omitempty"`
	// NetworkData contains network-related information specific to the device.
	NetworkData *NetworkDeviceData `json:"networkData,omitempty"`
}

func (in *AllocatedDeviceStatus) DeepCopyInto(out *AllocatedDeviceStatus) {
	*out = *in
	if in.Conditions != nil {
		l := make([]metav1.Condition, len(in.Conditions))
		for i := range in.Conditions {
			in.Conditions[i].DeepCopyInto(&l[i])
		}
		out.Conditions = l
	}
	if in.Data != nil {
		in, out := &in.Data, &out.Data
		*out = new(runtime.RawExtension)
		(*in).DeepCopyInto(*out)
	}
	if in.NetworkData != nil {
		in, out := &in.NetworkData, &out.NetworkData
		*out = new(NetworkDeviceData)
		(*in).DeepCopyInto(*out)
	}
}

func (in *AllocatedDeviceStatus) DeepCopy() *AllocatedDeviceStatus {
	if in == nil {
		return nil
	}
	out := new(AllocatedDeviceStatus)
	in.DeepCopyInto(out)
	return out
}

type ResourcePool struct {
	// Name is used to identify the pool. For node-local devices, this
	// is often the node name, but this is not required.
	// It must not be longer than 253 characters and must consist of one or more DNS sub-domains
	// separated by slashes. This field is immutable.
	Name string `json:"name"`
	// Generation tracks the change in a pool over time. Whenever a driver
	// changes something about one or more of the resources in a pool, it
	// must change the generation in all ResourceSlices which are part of
	// that pool. Consumers of ResourceSlices should only consider
	// resources from the pool with the highest generation number. The
	// generation may be reset by drivers, which should be fine for
	// consumers, assuming that all ResourceSlices in a pool are updated to
	// match or deleted.
	// Combined with ResourceSliceCount, this mechanism enables consumers to
	// detect pools which are comprised of multiple ResourceSlices and are
	// in an incomplete state.
	Generation int64 `json:"generation"`
	// ResourceSliceCount is the total number of ResourceSlices in the pool at this
	// generation number. Must be greater than zero.
	// Consumers can use this to check whether they have seen all ResourceSlices
	// belonging to the same pool.
	ResourceSliceCount int64 `json:"resourceSliceCount"`
}

func (in *ResourcePool) DeepCopyInto(out *ResourcePool) {
	*out = *in
}

func (in *ResourcePool) DeepCopy() *ResourcePool {
	if in == nil {
		return nil
	}
	out := new(ResourcePool)
	in.DeepCopyInto(out)
	return out
}

type Device struct {
	// Name is unique identifier among all devices managed by
	// the driver in the pool. It must be a DNS label.
	Name string `json:"name"`
	// Attributes defines the set of attributes for this device.
	// The name of each attribute must be unique in that set.
	// The maximum number of attributes and capacities combined is 32.
	Attributes map[string]DeviceAttribute `json:"attributes,omitempty"`
	// Capacity defines the set of capacities for this device.
	// The name of each capacity must be unique in that set.
	// The maximum number of attributes and capacities combined is 32.
	Capacity map[string]DeviceCapacity `json:"capacity,omitempty"`
	// ConsumesCounters defines a list of references to sharedCounters
	// and the set of counters that the device will
	// consume from those counter sets.
	// There can only be a single entry per counterSet.
	// The maximum number of device counter consumptions per
	// device is 2.
	ConsumesCounters []DeviceCounterConsumption `json:"consumesCounters"`
	// NodeName identifies the node where the device is available.
	// Must only be set if Spec.PerDeviceNodeSelection is set to true.
	// At most one of NodeName, NodeSelector and AllNodes can be set.
	NodeName string `json:"nodeName,omitempty"`
	// NodeSelector defines the nodes where the device is available.
	// Must use exactly one term.
	// Must only be set if Spec.PerDeviceNodeSelection is set to true.
	// At most one of NodeName, NodeSelector and AllNodes can be set.
	NodeSelector *corev1.NodeSelector `json:"nodeSelector,omitempty"`
	// AllNodes indicates that all nodes have access to the device.
	// Must only be set if Spec.PerDeviceNodeSelection is set to true.
	// At most one of NodeName, NodeSelector and AllNodes can be set.
	AllNodes bool `json:"allNodes,omitempty"`
	// If specified, these are the driver-defined taints.
	// The maximum number of taints is 16. If taints are set for
	// any device in a ResourceSlice, then the maximum number of
	// allowed devices per ResourceSlice is 64 instead of 128.
	// This is a beta field and requires enabling the DRADeviceTaints
	// feature gate.
	Taints []DeviceTaint `json:"taints"`
	// BindsToNode indicates if the usage of an allocation involving this device
	// has to be limited to exactly the node that was chosen when allocating the claim.
	// If set to true, the scheduler will set the ResourceClaim.Status.Allocation.NodeSelector
	// to match the node where the allocation was made.
	// This is a beta field and requires enabling the DRADeviceBindingConditions and DRAResourceClaimDeviceStatus
	// feature gates.
	BindsToNode bool `json:"bindsToNode,omitempty"`
	// BindingConditions defines the conditions for proceeding with binding.
	// All of these conditions must be set in the per-device status
	// conditions with a value of True to proceed with binding the pod to the node
	// while scheduling the pod.
	// The maximum number of binding conditions is 4.
	// The conditions must be a valid condition type string.
	// This is a beta field and requires enabling the DRADeviceBindingConditions and DRAResourceClaimDeviceStatus
	// feature gates.
	BindingConditions []string `json:"bindingConditions"`
	// BindingFailureConditions defines the conditions for binding failure.
	// They may be set in the per-device status conditions.
	// If any is set to "True", a binding failure occurred.
	// The maximum number of binding failure conditions is 4.
	// The conditions must be a valid condition type string.
	// This is a beta field and requires enabling the DRADeviceBindingConditions and DRAResourceClaimDeviceStatus
	// feature gates.
	BindingFailureConditions []string `json:"bindingFailureConditions"`
	// AllowMultipleAllocations marks whether the device is allowed to be allocated to multiple DeviceRequests.
	// If AllowMultipleAllocations is set to true, the device can be allocated more than once,
	// and all of its capacity is consumable, regardless of whether the requestPolicy is defined or not.
	AllowMultipleAllocations bool `json:"allowMultipleAllocations,omitempty"`
	// NodeAllocatableResourceMappings defines the mapping of node resources
	// that are managed by the DRA driver exposing this device. This includes resources currently
	// reported in v1.Node `status.allocatable` that are not extended resources
	// (see https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#extended-resources).
	// Examples include "cpu", "memory", "ephemeral-storage", and hugepages.
	// In addition to standard requests made through the Pod `spec`, these resources
	// can also be requested through claims and allocated by the DRA driver.
	// For example, a CPU DRA driver might allocate exclusive CPUs or auxiliary node memory
	// dependencies of an accelerator device.
	// The keys of this map are the node-allocatable resource names (e.g., "cpu", "memory").
	// Extended resource names are not permitted as keys.
	NodeAllocatableResourceMappings map[string]NodeAllocatableResourceMapping `json:"nodeAllocatableResourceMappings,omitempty"`
}

func (in *Device) DeepCopyInto(out *Device) {
	*out = *in
	if in.Attributes != nil {
		in, out := &in.Attributes, &out.Attributes
		*out = make(map[string]DeviceAttribute, len(*in))
		for k, v := range *in {
			(*out)[k] = v
		}
	}
	if in.Capacity != nil {
		in, out := &in.Capacity, &out.Capacity
		*out = make(map[string]DeviceCapacity, len(*in))
		for k, v := range *in {
			(*out)[k] = v
		}
	}
	if in.ConsumesCounters != nil {
		l := make([]DeviceCounterConsumption, len(in.ConsumesCounters))
		for i := range in.ConsumesCounters {
			in.ConsumesCounters[i].DeepCopyInto(&l[i])
		}
		out.ConsumesCounters = l
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = new(corev1.NodeSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Taints != nil {
		l := make([]DeviceTaint, len(in.Taints))
		for i := range in.Taints {
			in.Taints[i].DeepCopyInto(&l[i])
		}
		out.Taints = l
	}
	if in.BindingConditions != nil {
		t := make([]string, len(in.BindingConditions))
		copy(t, in.BindingConditions)
		out.BindingConditions = t
	}
	if in.BindingFailureConditions != nil {
		t := make([]string, len(in.BindingFailureConditions))
		copy(t, in.BindingFailureConditions)
		out.BindingFailureConditions = t
	}
	if in.NodeAllocatableResourceMappings != nil {
		in, out := &in.NodeAllocatableResourceMappings, &out.NodeAllocatableResourceMappings
		*out = make(map[string]NodeAllocatableResourceMapping, len(*in))
		for k, v := range *in {
			(*out)[k] = v
		}
	}
}

func (in *Device) DeepCopy() *Device {
	if in == nil {
		return nil
	}
	out := new(Device)
	in.DeepCopyInto(out)
	return out
}

type CounterSet struct {
	// Name defines the name of the counter set.
	// It must be a DNS label.
	Name string `json:"name"`
	// Counters defines the set of counters for this CounterSet
	// The name of each counter must be unique in that set and must be a DNS label.
	// The maximum number of counters is 32.
	Counters map[string]Counter `json:"counters,omitempty"`
}

func (in *CounterSet) DeepCopyInto(out *CounterSet) {
	*out = *in
	if in.Counters != nil {
		in, out := &in.Counters, &out.Counters
		*out = make(map[string]Counter, len(*in))
		for k, v := range *in {
			(*out)[k] = v
		}
	}
}

func (in *CounterSet) DeepCopy() *CounterSet {
	if in == nil {
		return nil
	}
	out := new(CounterSet)
	in.DeepCopyInto(out)
	return out
}

type CELDeviceSelector struct {
	// Expression is a CEL expression which evaluates a single device. It
	// must evaluate to true when the device under consideration satisfies
	// the desired criteria, and false when it does not. Any other result
	// is an error and causes allocation of devices to abort.
	// The expression's input is an object named "device", which carries
	// the following properties:
	// - driver (string): the name of the driver which defines this device.
	// - attributes (map[string]object): the device's attributes, grouped by prefix
	// (e.g. device.attributes["dra.example.com"] evaluates to an object with all
	// of the attributes which were prefixed by "dra.example.com".
	// - capacity (map[string]object): the device's capacities, grouped by prefix.
	// - allowMultipleAllocations (bool): the allowMultipleAllocations property of the device
	// (v1.34+ with the DRAConsumableCapacity feature enabled).
	// Example: Consider a device with driver="dra.example.com", which exposes
	// two attributes named "model" and "ext.example.com/family" and which
	// exposes one capacity named "modules". This input to this expression
	// would have the following fields:
	// device.driver
	// device.attributes["dra.example.com"].model
	// device.attributes["ext.example.com"].family
	// device.capacity["dra.example.com"].modules
	// The device.driver field can be used to check for a specific driver,
	// either as a high-level precondition (i.e. you only want to consider
	// devices from this driver) or as part of a multi-clause expression
	// that is meant to consider devices from different drivers.
	// The value type of each attribute is defined by the device
	// definition, and users who write these expressions must consult the
	// documentation for their specific drivers. The value type of each
	// capacity is Quantity.
	// If an unknown prefix is used as a lookup in either device.attributes
	// or device.capacity, an empty map will be returned. Any reference to
	// an unknown field will cause an evaluation error and allocation to
	// abort.
	// A robust expression should check for the existence of attributes
	// before referencing them.
	// For ease of use, the cel.bind() function is enabled, and can be used
	// to simplify expressions that access multiple attributes with the
	// same domain. For example:
	// cel.bind(dra, device.attributes["dra.example.com"], dra.someBool && dra.anotherBool)
	// When the DRAListTypeAttributes feature gate is enabled,
	// the includes() helper is available and it can work for both scalar
	// and list-type attributes. It was introduced to support smooth migration
	// from scalar attributes to list-type attributes while keeping
	// CEL expressions simple. For example:
	// device.attributes["dra.example.com"].models.includes("some-model")
	// The length of the expression must be smaller or equal to 10 Ki. The
	// cost of evaluating it is also limited based on the estimated number
	// of logical steps.
	Expression string `json:"expression"`
}

func (in *CELDeviceSelector) DeepCopyInto(out *CELDeviceSelector) {
	*out = *in
}

func (in *CELDeviceSelector) DeepCopy() *CELDeviceSelector {
	if in == nil {
		return nil
	}
	out := new(CELDeviceSelector)
	in.DeepCopyInto(out)
	return out
}

type DeviceConfiguration struct {
	// Opaque provides driver-specific configuration parameters.
	Opaque *OpaqueDeviceConfiguration `json:"opaque,omitempty"`
}

func (in *DeviceConfiguration) DeepCopyInto(out *DeviceConfiguration) {
	*out = *in
	if in.Opaque != nil {
		in, out := &in.Opaque, &out.Opaque
		*out = new(OpaqueDeviceConfiguration)
		(*in).DeepCopyInto(*out)
	}
}

func (in *DeviceConfiguration) DeepCopy() *DeviceConfiguration {
	if in == nil {
		return nil
	}
	out := new(DeviceConfiguration)
	in.DeepCopyInto(out)
	return out
}

type DeviceRequest struct {
	// Name can be used to reference this request in a pod.spec.containers[].resources.claims
	// entry and in a constraint of the claim.
	// References using the name in the DeviceRequest will uniquely
	// identify a request when the Exactly field is set. When the
	// FirstAvailable field is set, a reference to the name of the
	// DeviceRequest will match whatever subrequest is chosen by the
	// scheduler.
	// Must be a DNS label.
	Name string `json:"name"`
	// Exactly specifies the details for a single request that must
	// be met exactly for the request to be satisfied.
	// One of Exactly or FirstAvailable must be set.
	Exactly *ExactDeviceRequest `json:"exactly,omitempty"`
	// FirstAvailable contains subrequests, of which exactly one will be
	// selected by the scheduler. It tries to
	// satisfy them in the order in which they are listed here. So if
	// there are two entries in the list, the scheduler will only check
	// the second one if it determines that the first one can not be used.
	// DRA does not yet implement scoring, so the scheduler will
	// select the first set of devices that satisfies all the
	// requests in the claim. And if the requirements can
	// be satisfied on more than one node, other scheduling features
	// will determine which node is chosen. This means that the set of
	// devices allocated to a claim might not be the optimal set
	// available to the cluster. Scoring will be implemented later.
	FirstAvailable []DeviceSubRequest `json:"firstAvailable"`
}

func (in *DeviceRequest) DeepCopyInto(out *DeviceRequest) {
	*out = *in
	if in.Exactly != nil {
		in, out := &in.Exactly, &out.Exactly
		*out = new(ExactDeviceRequest)
		(*in).DeepCopyInto(*out)
	}
	if in.FirstAvailable != nil {
		l := make([]DeviceSubRequest, len(in.FirstAvailable))
		for i := range in.FirstAvailable {
			in.FirstAvailable[i].DeepCopyInto(&l[i])
		}
		out.FirstAvailable = l
	}
}

func (in *DeviceRequest) DeepCopy() *DeviceRequest {
	if in == nil {
		return nil
	}
	out := new(DeviceRequest)
	in.DeepCopyInto(out)
	return out
}

type DeviceConstraint struct {
	// Requests is a list of the one or more requests in this claim which
	// must co-satisfy this constraint. If a request is fulfilled by
	// multiple devices, then all of the devices must satisfy the
	// constraint. If this is not specified, this constraint applies to all
	// requests in this claim.
	// References to subrequests must include the name of the main request
	// and may include the subrequest using the format <main request>[/<subrequest>]. If just
	// the main request is given, the constraint applies to all subrequests.
	Requests []string `json:"requests"`
	// MatchAttribute requires that all devices in question have this
	// attribute and that its type and value are the same across those
	// devices.
	// For example, if you specified "dra.example.com/numa" (a hypothetical example!),
	// then only devices in the same NUMA node will be chosen. A device which
	// does not have that attribute will not be chosen. All devices should
	// use a value of the same type for this attribute because that is part of
	// its specification, but if one device doesn't, then it also will not be
	// chosen.
	// When the DRAListTypeAttributes feature gate is enabled, comparison uses
	// set semantics(i.e., element order and duplicates are ignored): list-valued attributes
	// match when the intersection across all devices is non-empty.
	// Scalar values are treated as single-element lists for backward compatibility.
	// Must include the domain qualifier.
	MatchAttribute string `json:"matchAttribute"`
	// DistinctAttribute requires that all devices in question have this
	// attribute and that its type and value are unique across those devices.
	// When the DRAListTypeAttributes feature gate is enabled, comparison uses
	// set semantics (i.e., element order and duplicates are ignored):
	// list-valued attributes must be pairwise disjoint across devices.
	// Scalar values are treated as singleton sets for backward compatibility.
	// This acts as the inverse of MatchAttribute.
	// This constraint is used to avoid allocating multiple requests to the same device
	// by ensuring attribute-level differentiation.
	// This is useful for scenarios where resource requests must be fulfilled by separate physical devices.
	// For example, a container requests two network interfaces that must be allocated from two different physical NICs.
	DistinctAttribute string `json:"distinctAttribute"`
}

func (in *DeviceConstraint) DeepCopyInto(out *DeviceConstraint) {
	*out = *in
	if in.Requests != nil {
		t := make([]string, len(in.Requests))
		copy(t, in.Requests)
		out.Requests = t
	}
}

func (in *DeviceConstraint) DeepCopy() *DeviceConstraint {
	if in == nil {
		return nil
	}
	out := new(DeviceConstraint)
	in.DeepCopyInto(out)
	return out
}

type DeviceClaimConfiguration struct {
	// Requests lists the names of requests where the configuration applies.
	// If empty, it applies to all requests.
	// References to subrequests must include the name of the main request
	// and may include the subrequest using the format <main request>[/<subrequest>]. If just
	// the main request is given, the configuration applies to all subrequests.
	Requests            []string `json:"requests"`
	DeviceConfiguration `json:",inline"`
}

func (in *DeviceClaimConfiguration) DeepCopyInto(out *DeviceClaimConfiguration) {
	*out = *in
	if in.Requests != nil {
		t := make([]string, len(in.Requests))
		copy(t, in.Requests)
		out.Requests = t
	}
	out.DeviceConfiguration = in.DeviceConfiguration
}

func (in *DeviceClaimConfiguration) DeepCopy() *DeviceClaimConfiguration {
	if in == nil {
		return nil
	}
	out := new(DeviceClaimConfiguration)
	in.DeepCopyInto(out)
	return out
}

type DeviceAllocationResult struct {
	// Results lists all allocated devices.
	Results []DeviceRequestAllocationResult `json:"results"`
	// This field is a combination of all the claim and class configuration parameters.
	// Drivers can distinguish between those based on a flag.
	// This includes configuration parameters for drivers which have no allocated
	// devices in the result because it is up to the drivers which configuration
	// parameters they support. They can silently ignore unknown configuration
	// parameters.
	Config []DeviceAllocationConfiguration `json:"config"`
}

func (in *DeviceAllocationResult) DeepCopyInto(out *DeviceAllocationResult) {
	*out = *in
	if in.Results != nil {
		l := make([]DeviceRequestAllocationResult, len(in.Results))
		for i := range in.Results {
			in.Results[i].DeepCopyInto(&l[i])
		}
		out.Results = l
	}
	if in.Config != nil {
		l := make([]DeviceAllocationConfiguration, len(in.Config))
		for i := range in.Config {
			in.Config[i].DeepCopyInto(&l[i])
		}
		out.Config = l
	}
}

func (in *DeviceAllocationResult) DeepCopy() *DeviceAllocationResult {
	if in == nil {
		return nil
	}
	out := new(DeviceAllocationResult)
	in.DeepCopyInto(out)
	return out
}

type NetworkDeviceData struct {
	// InterfaceName specifies the name of the network interface associated with
	// the allocated device. This might be the name of a physical or virtual
	// network interface being configured in the pod.
	// Must not be longer than 256 bytes.
	InterfaceName string `json:"interfaceName,omitempty"`
	// IPs lists the network addresses assigned to the device's network interface.
	// This can include both IPv4 and IPv6 addresses.
	// The IPs are in the CIDR notation, which includes both the address and the
	// associated subnet mask.
	// e.g.: "192.0.2.5/24" for IPv4 and "2001:db8::5/64" for IPv6.
	IPs []string `json:"ips"`
	// HardwareAddress represents the hardware address (e.g. MAC Address) of the device's network interface.
	// Must not be longer than 128 bytes.
	HardwareAddress string `json:"hardwareAddress,omitempty"`
}

func (in *NetworkDeviceData) DeepCopyInto(out *NetworkDeviceData) {
	*out = *in
	if in.IPs != nil {
		t := make([]string, len(in.IPs))
		copy(t, in.IPs)
		out.IPs = t
	}
}

func (in *NetworkDeviceData) DeepCopy() *NetworkDeviceData {
	if in == nil {
		return nil
	}
	out := new(NetworkDeviceData)
	in.DeepCopyInto(out)
	return out
}

type DeviceAttribute struct {
	// IntValue is a number.
	IntValue int64 `json:"int,omitempty"`
	// BoolValue is a true/false value.
	BoolValue bool `json:"bool,omitempty"`
	// StringValue is a string. Must not be longer than 64 characters.
	StringValue string `json:"string,omitempty"`
	// VersionValue is a semantic version according to semver.org spec 2.0.0.
	// Must not be longer than 64 characters.
	VersionValue string `json:"version,omitempty"`
	// IntValues is a non-empty list of numbers.
	// This is an alpha field and requires enabling the DRAListTypeAttributes feature gate.
	IntValues []int64 `json:"ints"`
	// BoolValues is a non-empty list of true/false values.
	BoolValues []bool `json:"bools"`
	// StringValues is a non-empty list of strings.
	// Each string must not be longer than 64 characters.
	// This is an alpha field and requires enabling the DRAListTypeAttributes feature gate.
	StringValues []string `json:"strings"`
	// VersionValues is a non-empty list of semantic versions according to semver.org spec 2.0.0.
	// Each version string must not be longer than 64 characters.
	// This is an alpha field and requires enabling the DRAListTypeAttributes feature gate.
	VersionValues []string `json:"versions"`
}

func (in *DeviceAttribute) DeepCopyInto(out *DeviceAttribute) {
	*out = *in
	if in.IntValues != nil {
		t := make([]int64, len(in.IntValues))
		copy(t, in.IntValues)
		out.IntValues = t
	}
	if in.BoolValues != nil {
		t := make([]bool, len(in.BoolValues))
		copy(t, in.BoolValues)
		out.BoolValues = t
	}
	if in.StringValues != nil {
		t := make([]string, len(in.StringValues))
		copy(t, in.StringValues)
		out.StringValues = t
	}
	if in.VersionValues != nil {
		t := make([]string, len(in.VersionValues))
		copy(t, in.VersionValues)
		out.VersionValues = t
	}
}

func (in *DeviceAttribute) DeepCopy() *DeviceAttribute {
	if in == nil {
		return nil
	}
	out := new(DeviceAttribute)
	in.DeepCopyInto(out)
	return out
}

type DeviceCapacity struct {
	// Value defines how much of a certain capacity that device has.
	// This field reflects the fixed total capacity and does not change.
	// The consumed amount is tracked separately by scheduler
	// and does not affect this value.
	Value apiresource.Quantity `json:"value"`
	// RequestPolicy defines how this DeviceCapacity must be consumed
	// when the device is allowed to be shared by multiple allocations.
	// The Device must have allowMultipleAllocations set to true in order to set a requestPolicy.
	// If unset, capacity requests are unconstrained:
	// requests can consume any amount of capacity, as long as the total consumed
	// across all allocations does not exceed the device's defined capacity.
	// If request is also unset, default is the full capacity value.
	RequestPolicy *CapacityRequestPolicy `json:"requestPolicy,omitempty"`
}

func (in *DeviceCapacity) DeepCopyInto(out *DeviceCapacity) {
	*out = *in
	in.Value.DeepCopyInto(&out.Value)
	if in.RequestPolicy != nil {
		in, out := &in.RequestPolicy, &out.RequestPolicy
		*out = new(CapacityRequestPolicy)
		(*in).DeepCopyInto(*out)
	}
}

func (in *DeviceCapacity) DeepCopy() *DeviceCapacity {
	if in == nil {
		return nil
	}
	out := new(DeviceCapacity)
	in.DeepCopyInto(out)
	return out
}

type DeviceCounterConsumption struct {
	// CounterSet is the name of the set from which the
	// counters defined will be consumed.
	CounterSet string `json:"counterSet"`
	// Counters defines the counters that will be consumed by the device.
	// The maximum number of counters is 32.
	Counters map[string]Counter `json:"counters,omitempty"`
}

func (in *DeviceCounterConsumption) DeepCopyInto(out *DeviceCounterConsumption) {
	*out = *in
	if in.Counters != nil {
		in, out := &in.Counters, &out.Counters
		*out = make(map[string]Counter, len(*in))
		for k, v := range *in {
			(*out)[k] = v
		}
	}
}

func (in *DeviceCounterConsumption) DeepCopy() *DeviceCounterConsumption {
	if in == nil {
		return nil
	}
	out := new(DeviceCounterConsumption)
	in.DeepCopyInto(out)
	return out
}

type DeviceTaint struct {
	// The taint key to be applied to a device.
	// Must be a label name.
	Key string `json:"key"`
	// The taint value corresponding to the taint key.
	// Must be a label value.
	Value string `json:"value,omitempty"`
	// The effect of the taint on claims that do not tolerate the taint
	// and through such claims on the pods using them.
	// Valid effects are None, NoSchedule and NoExecute. PreferNoSchedule as used for
	// nodes is not valid here. More effects may get added in the future.
	// Consumers must treat unknown effects like None.
	Effect DeviceTaintEffect `json:"effect"`
	// TimeAdded represents the time at which the taint was added or
	// (only in a DeviceTaintRule) the effect was modified.
	// Added automatically during create or update if not set.
	// In addition, in a DeviceTaintRule a value provided during
	// an update gets replaced with the current time if the provided
	// value is the same as the old one and the new effect is different.
	// Changing the key and/or value while keeping the effect unchanged
	// is possible and does not update the time stamp because the eviction
	// which uses it is either already started (NoExecute) or
	// not started yet (NoEffect, NoSchedule).
	TimeAdded *metav1.Time `json:"timeAdded,omitempty"`
}

func (in *DeviceTaint) DeepCopyInto(out *DeviceTaint) {
	*out = *in
	if in.TimeAdded != nil {
		in, out := &in.TimeAdded, &out.TimeAdded
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

func (in *DeviceTaint) DeepCopy() *DeviceTaint {
	if in == nil {
		return nil
	}
	out := new(DeviceTaint)
	in.DeepCopyInto(out)
	return out
}

type NodeAllocatableResourceMapping struct {
	// CapacityKey references a capacity name defined as a key in the
	// `spec.devices[*].capacity` map. When this field is set, the value associated with
	// this key in the `status.allocation.devices.results[*].consumedCapacity` map
	// (for a specific claim allocation) determines the base quantity for
	// the node allocatable resource. If `allocationMultiplier` is also set, it is
	// multiplied with the base quantity.
	// For example, if `spec.devices[*].capacity` has an entry "dra.example.com/memory": "128Gi",
	// and this field is set to "dra.example.com/memory", then for a claim allocation
	// that consumes { "dra.example.com/memory": "4Gi" } the base quantity for the
	// node allocatable resource mapping will be "4Gi", and `allocationMultiplier` should
	// be omitted or set to "1".
	CapacityKey string `json:"capacityKey"`
	// AllocationMultiplier is used as a multiplier for the allocated device count or the allocated capacity in the claim.
	// It defaults to 1 if not specified. How the field is used also depends on whether `capacityKey` is set.
	// 1.  If `capacityKey` is NOT set: `allocationMultiplier` multiplies the device count allocated to the claim.
	// a. A DRA driver representing each CPU core as a device would have
	// {ResourceName: "cpu", allocationMultiplier: "2"} in its
	// `nodeAllocatableResourceMappings`. If 4 devices are allocated to the claim,
	// 4 * 2 CPUs would be considered as allocated and subtracted from the node's capacity.
	// b. A GPU device that needs additional node memory per GPU allocation would
	// have {ResourceName: "memory", allocationMultiplier: "2Gi"}.  Each allocated
	// GPU device instance of this type will account for 2Gi of memory.
	// 2.  If `capacityKey` IS set: `allocationMultiplier` is multiplied by the amount of that capacity consumed.
	// The final node allocatable resource amount is `consumedCapacity[capacityKey]` * `allocationMultiplier`.
	// For example, if a Device's capacity "dra.example.com/cores" is consumed,
	// and each "core" provides 2 "cpu"s, the mapping would be:
	// {ResourceName: "cpu", capacityKey: "dra.example.com/cores", allocationMultiplier: "2"}.
	// If a claim consumes 8 "dra.example.com/cores", the CPU footprint is 8 * 2 = 16.
	AllocationMultiplier *apiresource.Quantity `json:"allocationMultiplier,omitempty"`
}

func (in *NodeAllocatableResourceMapping) DeepCopyInto(out *NodeAllocatableResourceMapping) {
	*out = *in
	if in.AllocationMultiplier != nil {
		in, out := &in.AllocationMultiplier, &out.AllocationMultiplier
		*out = new(apiresource.Quantity)
		(*in).DeepCopyInto(*out)
	}
}

func (in *NodeAllocatableResourceMapping) DeepCopy() *NodeAllocatableResourceMapping {
	if in == nil {
		return nil
	}
	out := new(NodeAllocatableResourceMapping)
	in.DeepCopyInto(out)
	return out
}

type Counter struct {
	// Value defines how much of a certain device counter is available.
	Value apiresource.Quantity `json:"value"`
}

func (in *Counter) DeepCopyInto(out *Counter) {
	*out = *in
	in.Value.DeepCopyInto(&out.Value)
}

func (in *Counter) DeepCopy() *Counter {
	if in == nil {
		return nil
	}
	out := new(Counter)
	in.DeepCopyInto(out)
	return out
}

type OpaqueDeviceConfiguration struct {
	// Driver is used to determine which kubelet plugin needs
	// to be passed these configuration parameters.
	// An admission policy provided by the driver developer could use this
	// to decide whether it needs to validate them.
	// Must be a DNS subdomain and should end with a DNS domain owned by the
	// vendor of the driver. It should use only lower case characters.
	Driver string `json:"driver"`
	// Parameters can contain arbitrary data. It is the responsibility of
	// the driver developer to handle validation and versioning. Typically this
	// includes self-identification and a version ("kind" + "apiVersion" for
	// Kubernetes types), with conversion between different versions.
	// The length of the raw data must be smaller or equal to 10 Ki.
	Parameters runtime.RawExtension `json:"parameters"`
}

func (in *OpaqueDeviceConfiguration) DeepCopyInto(out *OpaqueDeviceConfiguration) {
	*out = *in
	in.Parameters.DeepCopyInto(&out.Parameters)
}

func (in *OpaqueDeviceConfiguration) DeepCopy() *OpaqueDeviceConfiguration {
	if in == nil {
		return nil
	}
	out := new(OpaqueDeviceConfiguration)
	in.DeepCopyInto(out)
	return out
}

type ExactDeviceRequest struct {
	// DeviceClassName references a specific DeviceClass, which can define
	// additional configuration and selectors to be inherited by this
	// request.
	// A DeviceClassName is required.
	// Administrators may use this to restrict which devices may get
	// requested by only installing classes with selectors for permitted
	// devices. If users are free to request anything without restrictions,
	// then administrators can create an empty DeviceClass for users
	// to reference.
	DeviceClassName string `json:"deviceClassName"`
	// Selectors define criteria which must be satisfied by a specific
	// device in order for that device to be considered for this
	// request. All selectors must be satisfied for a device to be
	// considered.
	Selectors []DeviceSelector `json:"selectors"`
	// AllocationMode and its related fields define how devices are allocated
	// to satisfy this request. Supported values are:
	// - ExactCount: This request is for a specific number of devices.
	// This is the default. The exact number is provided in the
	// count field.
	// - All: This request is for all of the matching devices in a pool.
	// At least one device must exist on the node for the allocation to succeed.
	// Allocation will fail if some devices are already allocated,
	// unless adminAccess is requested.
	// If AllocationMode is not specified, the default mode is ExactCount. If
	// the mode is ExactCount and count is not specified, the default count is
	// one. Any other requests must specify this field.
	// More modes may get added in the future. Clients must refuse to handle
	// requests with unknown modes.
	AllocationMode DeviceAllocationMode `json:"allocationMode,omitempty"`
	// Count is used only when the count mode is "ExactCount". Must be greater than zero.
	// If AllocationMode is ExactCount and this field is not specified, the default is one.
	Count int64 `json:"count,omitempty"`
	// AdminAccess indicates that this is a claim for administrative access
	// to the device(s). Claims with AdminAccess are expected to be used for
	// monitoring or other management services for a device.  They ignore
	// all ordinary claims to the device with respect to access modes and
	// any resource allocations.
	// Admin access is disabled if this field is unset or
	// set to false, otherwise it is enabled.
	AdminAccess bool `json:"adminAccess,omitempty"`
	// If specified, the request's tolerations.
	// Tolerations for NoSchedule are required to allocate a
	// device which has a taint with that effect. The same applies
	// to NoExecute.
	// In addition, should any of the allocated devices get tainted
	// with NoExecute after allocation and that effect is not tolerated,
	// then all pods consuming the ResourceClaim get deleted to evict
	// them. The scheduler will not let new pods reserve the claim while
	// it has these tainted devices. Once all pods are evicted, the
	// claim will get deallocated.
	// The maximum number of tolerations is 16.
	// This is a beta field and requires enabling the DRADeviceTaints
	// feature gate.
	Tolerations []DeviceToleration `json:"tolerations"`
	// Capacity define resource requirements against each capacity.
	// If this field is unset and the device supports multiple allocations,
	// the default value will be applied to each capacity according to requestPolicy.
	// For the capacity that has no requestPolicy, default is the full capacity value.
	// Applies to each device allocation.
	// If Count > 1,
	// the request fails if there aren't enough devices that meet the requirements.
	// If AllocationMode is set to All,
	// the request fails if there are devices that otherwise match the request,
	// and have this capacity, with a value >= the requested amount, but which cannot be allocated to this request.
	Capacity *CapacityRequirements `json:"capacity,omitempty"`
}

func (in *ExactDeviceRequest) DeepCopyInto(out *ExactDeviceRequest) {
	*out = *in
	if in.Selectors != nil {
		l := make([]DeviceSelector, len(in.Selectors))
		for i := range in.Selectors {
			in.Selectors[i].DeepCopyInto(&l[i])
		}
		out.Selectors = l
	}
	if in.Tolerations != nil {
		l := make([]DeviceToleration, len(in.Tolerations))
		for i := range in.Tolerations {
			in.Tolerations[i].DeepCopyInto(&l[i])
		}
		out.Tolerations = l
	}
	if in.Capacity != nil {
		in, out := &in.Capacity, &out.Capacity
		*out = new(CapacityRequirements)
		(*in).DeepCopyInto(*out)
	}
}

func (in *ExactDeviceRequest) DeepCopy() *ExactDeviceRequest {
	if in == nil {
		return nil
	}
	out := new(ExactDeviceRequest)
	in.DeepCopyInto(out)
	return out
}

type DeviceSubRequest struct {
	// Name can be used to reference this subrequest in the list of constraints
	// or the list of configurations for the claim. References must use the
	// format <main request>/<subrequest>.
	// Must be a DNS label.
	Name string `json:"name"`
	// DeviceClassName references a specific DeviceClass, which can define
	// additional configuration and selectors to be inherited by this
	// subrequest.
	// A class is required. Which classes are available depends on the cluster.
	// Administrators may use this to restrict which devices may get
	// requested by only installing classes with selectors for permitted
	// devices. If users are free to request anything without restrictions,
	// then administrators can create an empty DeviceClass for users
	// to reference.
	DeviceClassName string `json:"deviceClassName"`
	// Selectors define criteria which must be satisfied by a specific
	// device in order for that device to be considered for this
	// subrequest. All selectors must be satisfied for a device to be
	// considered.
	Selectors []DeviceSelector `json:"selectors"`
	// AllocationMode and its related fields define how devices are allocated
	// to satisfy this subrequest. Supported values are:
	// - ExactCount: This request is for a specific number of devices.
	// This is the default. The exact number is provided in the
	// count field.
	// - All: This subrequest is for all of the matching devices in a pool.
	// Allocation will fail if some devices are already allocated,
	// unless adminAccess is requested.
	// If AllocationMode is not specified, the default mode is ExactCount. If
	// the mode is ExactCount and count is not specified, the default count is
	// one. Any other subrequests must specify this field.
	// More modes may get added in the future. Clients must refuse to handle
	// requests with unknown modes.
	AllocationMode DeviceAllocationMode `json:"allocationMode,omitempty"`
	// Count is used only when the count mode is "ExactCount". Must be greater than zero.
	// If AllocationMode is ExactCount and this field is not specified, the default is one.
	Count int64 `json:"count,omitempty"`
	// If specified, the request's tolerations.
	// Tolerations for NoSchedule are required to allocate a
	// device which has a taint with that effect. The same applies
	// to NoExecute.
	// In addition, should any of the allocated devices get tainted
	// with NoExecute after allocation and that effect is not tolerated,
	// then all pods consuming the ResourceClaim get deleted to evict
	// them. The scheduler will not let new pods reserve the claim while
	// it has these tainted devices. Once all pods are evicted, the
	// claim will get deallocated.
	// The maximum number of tolerations is 16.
	// This is a beta field and requires enabling the DRADeviceTaints
	// feature gate.
	Tolerations []DeviceToleration `json:"tolerations"`
	// Capacity define resource requirements against each capacity.
	// If this field is unset and the device supports multiple allocations,
	// the default value will be applied to each capacity according to requestPolicy.
	// For the capacity that has no requestPolicy, default is the full capacity value.
	// Applies to each device allocation.
	// If Count > 1,
	// the request fails if there aren't enough devices that meet the requirements.
	// If AllocationMode is set to All,
	// the request fails if there are devices that otherwise match the request,
	// and have this capacity, with a value >= the requested amount, but which cannot be allocated to this request.
	Capacity *CapacityRequirements `json:"capacity,omitempty"`
}

func (in *DeviceSubRequest) DeepCopyInto(out *DeviceSubRequest) {
	*out = *in
	if in.Selectors != nil {
		l := make([]DeviceSelector, len(in.Selectors))
		for i := range in.Selectors {
			in.Selectors[i].DeepCopyInto(&l[i])
		}
		out.Selectors = l
	}
	if in.Tolerations != nil {
		l := make([]DeviceToleration, len(in.Tolerations))
		for i := range in.Tolerations {
			in.Tolerations[i].DeepCopyInto(&l[i])
		}
		out.Tolerations = l
	}
	if in.Capacity != nil {
		in, out := &in.Capacity, &out.Capacity
		*out = new(CapacityRequirements)
		(*in).DeepCopyInto(*out)
	}
}

func (in *DeviceSubRequest) DeepCopy() *DeviceSubRequest {
	if in == nil {
		return nil
	}
	out := new(DeviceSubRequest)
	in.DeepCopyInto(out)
	return out
}

type DeviceRequestAllocationResult struct {
	// Request is the name of the request in the claim which caused this
	// device to be allocated. If it references a subrequest in the
	// firstAvailable list on a DeviceRequest, this field must
	// include both the name of the main request and the subrequest
	// using the format <main request>/<subrequest>.
	// Multiple devices may have been allocated per request.
	Request string `json:"request"`
	// Driver specifies the name of the DRA driver whose kubelet
	// plugin should be invoked to process the allocation once the claim is
	// needed on a node.
	// Must be a DNS subdomain and should end with a DNS domain owned by the
	// vendor of the driver. It should use only lower case characters.
	Driver string `json:"driver"`
	// This name together with the driver name and the device name field
	// identify which device was allocated (`<driver name>/<pool name>/<device name>`).
	// Must not be longer than 253 characters and may contain one or more
	// DNS sub-domains separated by slashes.
	Pool string `json:"pool"`
	// Device references one device instance via its name in the driver's
	// resource pool. It must be a DNS label.
	Device string `json:"device"`
	// AdminAccess indicates that this device was allocated for
	// administrative access. See the corresponding request field
	// for a definition of mode.
	// Admin access is disabled if this field is unset or
	// set to false, otherwise it is enabled.
	AdminAccess bool `json:"adminAccess,omitempty"`
	// A copy of all tolerations specified in the request at the time
	// when the device got allocated.
	// The maximum number of tolerations is 16.
	// This is a beta field and requires enabling the DRADeviceTaints
	// feature gate.
	Tolerations []DeviceToleration `json:"tolerations"`
	// BindingConditions contains a copy of the BindingConditions
	// from the corresponding ResourceSlice at the time of allocation.
	// This is a beta field and requires enabling the DRADeviceBindingConditions and DRAResourceClaimDeviceStatus
	// feature gates.
	BindingConditions []string `json:"bindingConditions"`
	// BindingFailureConditions contains a copy of the BindingFailureConditions
	// from the corresponding ResourceSlice at the time of allocation.
	// This is a beta field and requires enabling the DRADeviceBindingConditions and DRAResourceClaimDeviceStatus
	// feature gates.
	BindingFailureConditions []string `json:"bindingFailureConditions"`
	// ShareID uniquely identifies an individual allocation share of the device,
	// used when the device supports multiple simultaneous allocations.
	// It serves as an additional map key to differentiate concurrent shares
	// of the same device.
	ShareID string `json:"shareID,omitempty"`
	// ConsumedCapacity tracks the amount of capacity consumed per device as part of the claim request.
	// The consumed amount may differ from the requested amount: it is rounded up to the nearest valid
	// value based on the device’s requestPolicy if applicable (i.e., may not be less than the requested amount).
	// The total consumed capacity for each device must not exceed the DeviceCapacity's Value.
	// This field is populated only for devices that allow multiple allocations.
	// All capacity entries are included, even if the consumed amount is zero.
	ConsumedCapacity map[string]apiresource.Quantity `json:"consumedCapacity,omitempty"`
}

func (in *DeviceRequestAllocationResult) DeepCopyInto(out *DeviceRequestAllocationResult) {
	*out = *in
	if in.Tolerations != nil {
		l := make([]DeviceToleration, len(in.Tolerations))
		for i := range in.Tolerations {
			in.Tolerations[i].DeepCopyInto(&l[i])
		}
		out.Tolerations = l
	}
	if in.BindingConditions != nil {
		t := make([]string, len(in.BindingConditions))
		copy(t, in.BindingConditions)
		out.BindingConditions = t
	}
	if in.BindingFailureConditions != nil {
		t := make([]string, len(in.BindingFailureConditions))
		copy(t, in.BindingFailureConditions)
		out.BindingFailureConditions = t
	}
	if in.ConsumedCapacity != nil {
		in, out := &in.ConsumedCapacity, &out.ConsumedCapacity
		*out = make(map[string]apiresource.Quantity, len(*in))
		for k, v := range *in {
			(*out)[k] = v
		}
	}
}

func (in *DeviceRequestAllocationResult) DeepCopy() *DeviceRequestAllocationResult {
	if in == nil {
		return nil
	}
	out := new(DeviceRequestAllocationResult)
	in.DeepCopyInto(out)
	return out
}

type DeviceAllocationConfiguration struct {
	// Source records whether the configuration comes from a class and thus
	// is not something that a normal user would have been able to set
	// or from a claim.
	Source AllocationConfigSource `json:"source"`
	// Requests lists the names of requests where the configuration applies.
	// If empty, its applies to all requests.
	// References to subrequests must include the name of the main request
	// and may include the subrequest using the format <main request>[/<subrequest>]. If just
	// the main request is given, the configuration applies to all subrequests.
	Requests            []string `json:"requests"`
	DeviceConfiguration `json:",inline"`
}

func (in *DeviceAllocationConfiguration) DeepCopyInto(out *DeviceAllocationConfiguration) {
	*out = *in
	if in.Requests != nil {
		t := make([]string, len(in.Requests))
		copy(t, in.Requests)
		out.Requests = t
	}
	out.DeviceConfiguration = in.DeviceConfiguration
}

func (in *DeviceAllocationConfiguration) DeepCopy() *DeviceAllocationConfiguration {
	if in == nil {
		return nil
	}
	out := new(DeviceAllocationConfiguration)
	in.DeepCopyInto(out)
	return out
}

type CapacityRequestPolicy struct {
	// Default specifies how much of this capacity is consumed by a request
	// that does not contain an entry for it in DeviceRequest's Capacity.
	Default *apiresource.Quantity `json:"default,omitempty"`
	// ValidValues defines a set of acceptable quantity values in consuming requests.
	// Must not contain more than 10 entries.
	// Must be sorted in ascending order.
	// If this field is set,
	// Default must be defined and it must be included in ValidValues list.
	// If the requested amount does not match any valid value but smaller than some valid values,
	// the scheduler calculates the smallest valid value that is greater than or equal to the request.
	// That is: min(ceil(requestedValue) ∈ validValues), where requestedValue ≤ max(validValues).
	// If the requested amount exceeds all valid values, the request violates the policy,
	// and this device cannot be allocated.
	ValidValues []apiresource.Quantity `json:"validValues"`
	// ValidRange defines an acceptable quantity value range in consuming requests.
	// If this field is set,
	// Default must be defined and it must fall within the defined ValidRange.
	// If the requested amount does not fall within the defined range, the request violates the policy,
	// and this device cannot be allocated.
	// If the request doesn't contain this capacity entry, Default value is used.
	ValidRange *CapacityRequestPolicyRange `json:"validRange,omitempty"`
}

func (in *CapacityRequestPolicy) DeepCopyInto(out *CapacityRequestPolicy) {
	*out = *in
	if in.Default != nil {
		in, out := &in.Default, &out.Default
		*out = new(apiresource.Quantity)
		(*in).DeepCopyInto(*out)
	}
	if in.ValidValues != nil {
		l := make([]apiresource.Quantity, len(in.ValidValues))
		for i := range in.ValidValues {
			in.ValidValues[i].DeepCopyInto(&l[i])
		}
		out.ValidValues = l
	}
	if in.ValidRange != nil {
		in, out := &in.ValidRange, &out.ValidRange
		*out = new(CapacityRequestPolicyRange)
		(*in).DeepCopyInto(*out)
	}
}

func (in *CapacityRequestPolicy) DeepCopy() *CapacityRequestPolicy {
	if in == nil {
		return nil
	}
	out := new(CapacityRequestPolicy)
	in.DeepCopyInto(out)
	return out
}

type DeviceToleration struct {
	// Key is the taint key that the toleration applies to. Empty means match all taint keys.
	// If the key is empty, operator must be Exists; this combination means to match all values and all keys.
	// Must be a label name.
	Key string `json:"key,omitempty"`
	// Operator represents a key's relationship to the value.
	// Valid operators are Exists and Equal. Defaults to Equal.
	// Exists is equivalent to wildcard for value, so that a ResourceClaim can
	// tolerate all taints of a particular category.
	Operator DeviceTolerationOperator `json:"operator,omitempty"`
	// Value is the taint value the toleration matches to.
	// If the operator is Exists, the value must be empty, otherwise just a regular string.
	// Must be a label value.
	Value string `json:"value,omitempty"`
	// Effect indicates the taint effect to match. Empty means match all taint effects.
	// When specified, allowed values are NoSchedule and NoExecute.
	Effect DeviceTaintEffect `json:"effect,omitempty"`
	// TolerationSeconds represents the period of time the toleration (which must be
	// of effect NoExecute, otherwise this field is ignored) tolerates the taint. By default,
	// it is not set, which means tolerate the taint forever (do not evict). Zero and
	// negative values will be treated as 0 (evict immediately) by the system.
	// If larger than zero, the time when the pod needs to be evicted is calculated as <time when
	// taint was adedd> + <toleration seconds>.
	TolerationSeconds int64 `json:"tolerationSeconds,omitempty"`
}

func (in *DeviceToleration) DeepCopyInto(out *DeviceToleration) {
	*out = *in
}

func (in *DeviceToleration) DeepCopy() *DeviceToleration {
	if in == nil {
		return nil
	}
	out := new(DeviceToleration)
	in.DeepCopyInto(out)
	return out
}

type CapacityRequirements struct {
	// Requests represent individual device resource requests for distinct resources,
	// all of which must be provided by the device.
	// This value is used as an additional filtering condition against the available capacity on the device.
	// This is semantically equivalent to a CEL selector with
	// `device.capacity[<domain>].<name>.compareTo(quantity(<request quantity>)) >= 0`.
	// For example, device.capacity['test-driver.cdi.k8s.io'].counters.compareTo(quantity('2')) >= 0.
	// When a requestPolicy is defined, the requested amount is adjusted upward
	// to the nearest valid value based on the policy.
	// If the requested amount cannot be adjusted to a valid value—because it exceeds what the requestPolicy allows—
	// the device is considered ineligible for allocation.
	// For any capacity that is not explicitly requested:
	// - If no requestPolicy is set, the default consumed capacity is equal to the full device capacity
	// (i.e., the whole device is claimed).
	// - If a requestPolicy is set, the default consumed capacity is determined according to that policy.
	// If the device allows multiple allocation,
	// the aggregated amount across all requests must not exceed the capacity value.
	// The consumed capacity, which may be adjusted based on the requestPolicy if defined,
	// is recorded in the resource claim’s status.devices[*].consumedCapacity field.
	Requests map[string]apiresource.Quantity `json:"requests,omitempty"`
}

func (in *CapacityRequirements) DeepCopyInto(out *CapacityRequirements) {
	*out = *in
	if in.Requests != nil {
		in, out := &in.Requests, &out.Requests
		*out = make(map[string]apiresource.Quantity, len(*in))
		for k, v := range *in {
			(*out)[k] = v
		}
	}
}

func (in *CapacityRequirements) DeepCopy() *CapacityRequirements {
	if in == nil {
		return nil
	}
	out := new(CapacityRequirements)
	in.DeepCopyInto(out)
	return out
}

type CapacityRequestPolicyRange struct {
	// Min specifies the minimum capacity allowed for a consumption request.
	// Min must be greater than or equal to zero,
	// and less than or equal to the capacity value.
	// requestPolicy.default must be more than or equal to the minimum.
	Min *apiresource.Quantity `json:"min,omitempty"`
	// Max defines the upper limit for capacity that can be requested.
	// Max must be less than or equal to the capacity value.
	// Min and requestPolicy.default must be less than or equal to the maximum.
	Max *apiresource.Quantity `json:"max,omitempty"`
	// Step defines the step size between valid capacity amounts within the range.
	// Max (if set) and requestPolicy.default must be a multiple of Step.
	// Min + Step must be less than or equal to the capacity value.
	Step *apiresource.Quantity `json:"step,omitempty"`
}

func (in *CapacityRequestPolicyRange) DeepCopyInto(out *CapacityRequestPolicyRange) {
	*out = *in
	if in.Min != nil {
		in, out := &in.Min, &out.Min
		*out = new(apiresource.Quantity)
		(*in).DeepCopyInto(*out)
	}
	if in.Max != nil {
		in, out := &in.Max, &out.Max
		*out = new(apiresource.Quantity)
		(*in).DeepCopyInto(*out)
	}
	if in.Step != nil {
		in, out := &in.Step, &out.Step
		*out = new(apiresource.Quantity)
		(*in).DeepCopyInto(*out)
	}
}

func (in *CapacityRequestPolicyRange) DeepCopy() *CapacityRequestPolicyRange {
	if in == nil {
		return nil
	}
	out := new(CapacityRequestPolicyRange)
	in.DeepCopyInto(out)
	return out
}
