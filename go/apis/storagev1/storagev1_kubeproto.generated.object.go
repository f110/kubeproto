package storagev1

import (
	"go.f110.dev/kubeproto/go/apis/corev1"
	"go.f110.dev/kubeproto/go/apis/metav1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "storage.k8s.io"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&CSIDriver{},
		&CSIDriverList{},
		&CSINode{},
		&CSINodeList{},
		&CSIStorageCapacity{},
		&CSIStorageCapacityList{},
		&StorageClass{},
		&StorageClassList{},
		&VolumeAttachment{},
		&VolumeAttachmentList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type FSGroupPolicy string

const (
	FSGroupPolicyReadWriteOnceWithFSType FSGroupPolicy = "ReadWriteOnceWithFSType"
	FSGroupPolicyFile                    FSGroupPolicy = "File"
	FSGroupPolicyNone                    FSGroupPolicy = "None"
)

type VolumeBindingMode string

const (
	VolumeBindingModeImmediate            VolumeBindingMode = "Immediate"
	VolumeBindingModeWaitForFirstConsumer VolumeBindingMode = "WaitForFirstConsumer"
)

type VolumeLifecycleMode string

const (
	VolumeLifecycleModePersistent VolumeLifecycleMode = "Persistent"
	VolumeLifecycleModeEphemeral  VolumeLifecycleMode = "Ephemeral"
)

type CSIDriver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// spec represents the specification of the CSI Driver.
	Spec CSIDriverSpec `json:"spec"`
}

func (in *CSIDriver) DeepCopyInto(out *CSIDriver) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

func (in *CSIDriver) DeepCopy() *CSIDriver {
	if in == nil {
		return nil
	}
	out := new(CSIDriver)
	in.DeepCopyInto(out)
	return out
}

func (in *CSIDriver) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type CSIDriverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []CSIDriver `json:"items"`
}

func (in *CSIDriverList) DeepCopyInto(out *CSIDriverList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]CSIDriver, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *CSIDriverList) DeepCopy() *CSIDriverList {
	if in == nil {
		return nil
	}
	out := new(CSIDriverList)
	in.DeepCopyInto(out)
	return out
}

func (in *CSIDriverList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type CSINode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// spec is the specification of CSINode
	Spec CSINodeSpec `json:"spec"`
}

func (in *CSINode) DeepCopyInto(out *CSINode) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

func (in *CSINode) DeepCopy() *CSINode {
	if in == nil {
		return nil
	}
	out := new(CSINode)
	in.DeepCopyInto(out)
	return out
}

func (in *CSINode) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type CSINodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []CSINode `json:"items"`
}

func (in *CSINodeList) DeepCopyInto(out *CSINodeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]CSINode, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *CSINodeList) DeepCopy() *CSINodeList {
	if in == nil {
		return nil
	}
	out := new(CSINodeList)
	in.DeepCopyInto(out)
	return out
}

func (in *CSINodeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type CSIStorageCapacity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// nodeTopology defines which nodes have access to the storage
	// for which capacity was reported. If not set, the storage is
	// not accessible from any node in the cluster. If empty, the
	// storage is accessible from all nodes. This field is
	// immutable.
	NodeTopology *metav1.LabelSelector `json:"nodeTopology,omitempty"`
	// storageClassName represents the name of the StorageClass that the reported capacity applies to.
	// It must meet the same requirements as the name of a StorageClass
	// object (non-empty, DNS subdomain). If that object no longer exists,
	// the CSIStorageCapacity object is obsolete and should be removed by its
	// creator.
	// This field is immutable.
	StorageClassName string `json:"storageClassName"`
	// capacity is the value reported by the CSI driver in its GetCapacityResponse
	// for a GetCapacityRequest with topology and parameters that match the
	// previous fields.
	// The semantic is currently (CSI spec 1.2) defined as:
	// The available capacity, in bytes, of the storage that can be used
	// to provision volumes. If not set, that information is currently
	// unavailable.
	Capacity *apiresource.Quantity `json:"capacity,omitempty"`
	// maximumVolumeSize is the value reported by the CSI driver in its GetCapacityResponse
	// for a GetCapacityRequest with topology and parameters that match the
	// previous fields.
	// This is defined since CSI spec 1.4.0 as the largest size
	// that may be used in a
	// CreateVolumeRequest.capacity_range.required_bytes field to
	// create a volume with the same parameters as those in
	// GetCapacityRequest. The corresponding value in the Kubernetes
	// API is ResourceRequirements.Requests in a volume claim.
	MaximumVolumeSize *apiresource.Quantity `json:"maximumVolumeSize,omitempty"`
}

func (in *CSIStorageCapacity) DeepCopyInto(out *CSIStorageCapacity) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.NodeTopology != nil {
		in, out := &in.NodeTopology, &out.NodeTopology
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Capacity != nil {
		in, out := &in.Capacity, &out.Capacity
		*out = new(apiresource.Quantity)
		(*in).DeepCopyInto(*out)
	}
	if in.MaximumVolumeSize != nil {
		in, out := &in.MaximumVolumeSize, &out.MaximumVolumeSize
		*out = new(apiresource.Quantity)
		(*in).DeepCopyInto(*out)
	}
}

func (in *CSIStorageCapacity) DeepCopy() *CSIStorageCapacity {
	if in == nil {
		return nil
	}
	out := new(CSIStorageCapacity)
	in.DeepCopyInto(out)
	return out
}

func (in *CSIStorageCapacity) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type CSIStorageCapacityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []CSIStorageCapacity `json:"items"`
}

func (in *CSIStorageCapacityList) DeepCopyInto(out *CSIStorageCapacityList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]CSIStorageCapacity, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *CSIStorageCapacityList) DeepCopy() *CSIStorageCapacityList {
	if in == nil {
		return nil
	}
	out := new(CSIStorageCapacityList)
	in.DeepCopyInto(out)
	return out
}

func (in *CSIStorageCapacityList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type StorageClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// provisioner indicates the type of the provisioner.
	Provisioner string `json:"provisioner"`
	// parameters holds the parameters for the provisioner that should
	// create volumes of this storage class.
	Parameters map[string]string `json:"parameters,omitempty"`
	// reclaimPolicy controls the reclaimPolicy for dynamically provisioned PersistentVolumes of this storage class.
	// Defaults to Delete.
	ReclaimPolicy corev1.PersistentVolumeReclaimPolicy `json:"reclaimPolicy,omitempty"`
	// mountOptions controls the mountOptions for dynamically provisioned PersistentVolumes of this storage class.
	// e.g. ["ro", "soft"]. Not validated -
	// mount of the PVs will simply fail if one is invalid.
	MountOptions []string `json:"mountOptions"`
	// allowVolumeExpansion shows whether the storage class allow volume expand.
	AllowVolumeExpansion bool `json:"allowVolumeExpansion,omitempty"`
	// volumeBindingMode indicates how PersistentVolumeClaims should be
	// provisioned and bound.  When unset, VolumeBindingImmediate is used.
	// This field is only honored by servers that enable the VolumeScheduling feature.
	VolumeBindingMode VolumeBindingMode `json:"volumeBindingMode,omitempty"`
	// allowedTopologies restrict the node topologies where volumes can be dynamically provisioned.
	// Each volume plugin defines its own supported topology specifications.
	// An empty TopologySelectorTerm list means there is no topology restriction.
	// This field is only honored by servers that enable the VolumeScheduling feature.
	AllowedTopologies []corev1.TopologySelectorTerm `json:"allowedTopologies"`
}

func (in *StorageClass) DeepCopyInto(out *StorageClass) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Parameters != nil {
		in, out := &in.Parameters, &out.Parameters
		*out = make(map[string]string, len(*in))
		for k, v := range *in {
			(*out)[k] = v
		}
	}
	if in.MountOptions != nil {
		t := make([]string, len(in.MountOptions))
		copy(t, in.MountOptions)
		out.MountOptions = t
	}
	if in.AllowedTopologies != nil {
		l := make([]corev1.TopologySelectorTerm, len(in.AllowedTopologies))
		for i := range in.AllowedTopologies {
			in.AllowedTopologies[i].DeepCopyInto(&l[i])
		}
		out.AllowedTopologies = l
	}
}

func (in *StorageClass) DeepCopy() *StorageClass {
	if in == nil {
		return nil
	}
	out := new(StorageClass)
	in.DeepCopyInto(out)
	return out
}

func (in *StorageClass) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type StorageClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []StorageClass `json:"items"`
}

func (in *StorageClassList) DeepCopyInto(out *StorageClassList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]StorageClass, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *StorageClassList) DeepCopy() *StorageClassList {
	if in == nil {
		return nil
	}
	out := new(StorageClassList)
	in.DeepCopyInto(out)
	return out
}

func (in *StorageClassList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type VolumeAttachment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	// spec represents specification of the desired attach/detach volume behavior.
	// Populated by the Kubernetes system.
	Spec VolumeAttachmentSpec `json:"spec"`
	// status represents status of the VolumeAttachment request.
	// Populated by the entity completing the attach or detach
	// operation, i.e. the external-attacher.
	Status *VolumeAttachmentStatus `json:"status,omitempty"`
}

func (in *VolumeAttachment) DeepCopyInto(out *VolumeAttachment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(VolumeAttachmentStatus)
		(*in).DeepCopyInto(*out)
	}
}

func (in *VolumeAttachment) DeepCopy() *VolumeAttachment {
	if in == nil {
		return nil
	}
	out := new(VolumeAttachment)
	in.DeepCopyInto(out)
	return out
}

func (in *VolumeAttachment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type VolumeAttachmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []VolumeAttachment `json:"items"`
}

func (in *VolumeAttachmentList) DeepCopyInto(out *VolumeAttachmentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		l := make([]VolumeAttachment, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&l[i])
		}
		out.Items = l
	}
}

func (in *VolumeAttachmentList) DeepCopy() *VolumeAttachmentList {
	if in == nil {
		return nil
	}
	out := new(VolumeAttachmentList)
	in.DeepCopyInto(out)
	return out
}

func (in *VolumeAttachmentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

type CSIDriverSpec struct {
	// attachRequired indicates this CSI volume driver requires an attach
	// operation (because it implements the CSI ControllerPublishVolume()
	// method), and that the Kubernetes attach detach controller should call
	// the attach volume interface which checks the volumeattachment status
	// and waits until the volume is attached before proceeding to mounting.
	// The CSI external-attacher coordinates with CSI volume driver and updates
	// the volumeattachment status when the attach operation is complete.
	// If the CSIDriverRegistry feature gate is enabled and the value is
	// specified to false, the attach operation will be skipped.
	// Otherwise the attach operation will be called.
	// This field is immutable.
	AttachRequired bool `json:"attachRequired,omitempty"`
	// podInfoOnMount indicates this CSI volume driver requires additional pod information (like podName, podUID, etc.)
	// during mount operations, if set to true.
	// If set to false, pod information will not be passed on mount.
	// Default is false.
	// The CSI driver specifies podInfoOnMount as part of driver deployment.
	// If true, Kubelet will pass pod information as VolumeContext in the CSI NodePublishVolume() calls.
	// The CSI driver is responsible for parsing and validating the information passed in as VolumeContext.
	// The following VolumeConext will be passed if podInfoOnMount is set to true.
	// This list might grow, but the prefix will be used.
	// "csi.storage.k8s.io/pod.name": pod.Name
	// "csi.storage.k8s.io/pod.namespace": pod.Namespace
	// "csi.storage.k8s.io/pod.uid": string(pod.UID)
	// "csi.storage.k8s.io/ephemeral": "true" if the volume is an ephemeral inline volume
	// defined by a CSIVolumeSource, otherwise "false"
	// "csi.storage.k8s.io/ephemeral" is a new feature in Kubernetes 1.16. It is only
	// required for drivers which support both the "Persistent" and "Ephemeral" VolumeLifecycleMode.
	// Other drivers can leave pod info disabled and/or ignore this field.
	// As Kubernetes 1.15 doesn't support this field, drivers can only support one mode when
	// deployed on such a cluster and the deployment determines which mode that is, for example
	// via a command line parameter of the driver.
	// This field is immutable.
	PodInfoOnMount bool `json:"podInfoOnMount,omitempty"`
	// volumeLifecycleModes defines what kind of volumes this CSI volume driver supports.
	// The default if the list is empty is "Persistent", which is the usage defined by the
	// CSI specification and implemented in Kubernetes via the usual PV/PVC mechanism.
	// The other mode is "Ephemeral". In this mode, volumes are defined inline inside the pod spec
	// with CSIVolumeSource and their lifecycle is tied to the lifecycle of that pod.
	// A driver has to be aware of this because it is only going to get a NodePublishVolume call for such a volume.
	// For more information about implementing this mode, see
	// https://kubernetes-csi.github.io/docs/ephemeral-local-volumes.html
	// A driver can support one or more of these modes and more modes may be added in the future.
	// This field is beta.
	// This field is immutable.
	VolumeLifecycleModes []VolumeLifecycleMode `json:"volumeLifecycleModes"`
	// storageCapacity indicates that the CSI volume driver wants pod scheduling to consider the storage
	// capacity that the driver deployment will report by creating
	// CSIStorageCapacity objects with capacity information, if set to true.
	// The check can be enabled immediately when deploying a driver.
	// In that case, provisioning new volumes with late binding
	// will pause until the driver deployment has published
	// some suitable CSIStorageCapacity object.
	// Alternatively, the driver can be deployed with the field
	// unset or false and it can be flipped later when storage
	// capacity information has been published.
	// This field was immutable in Kubernetes <= 1.22 and now is mutable.
	StorageCapacity bool `json:"storageCapacity,omitempty"`
	// fsGroupPolicy defines if the underlying volume supports changing ownership and
	// permission of the volume before being mounted.
	// Refer to the specific FSGroupPolicy values for additional details.
	// This field is immutable.
	// Defaults to ReadWriteOnceWithFSType, which will examine each volume
	// to determine if Kubernetes should modify ownership and permissions of the volume.
	// With the default policy the defined fsGroup will only be applied
	// if a fstype is defined and the volume's access mode contains ReadWriteOnce.
	FSGroupPolicy FSGroupPolicy `json:"fsGroupPolicy,omitempty"`
	// tokenRequests indicates the CSI driver needs pods' service account
	// tokens it is mounting volume for to do necessary authentication. Kubelet
	// will pass the tokens in VolumeContext in the CSI NodePublishVolume calls.
	// The CSI driver should parse and validate the following VolumeContext:
	// "csi.storage.k8s.io/serviceAccount.tokens": {
	// "<audience>": {
	// "token": <token>,
	// "expirationTimestamp": <expiration timestamp in RFC3339>,
	// },
	// ...
	// }
	// Note: Audience in each TokenRequest should be different and at
	// most one token is empty string. To receive a new token after expiry,
	// RequiresRepublish can be used to trigger NodePublishVolume periodically.
	TokenRequests []TokenRequest `json:"tokenRequests"`
	// requiresRepublish indicates the CSI driver wants `NodePublishVolume`
	// being periodically called to reflect any possible change in the mounted
	// volume. This field defaults to false.
	// Note: After a successful initial NodePublishVolume call, subsequent calls
	// to NodePublishVolume should only update the contents of the volume. New
	// mount points will not be seen by a running container.
	RequiresRepublish bool `json:"requiresRepublish,omitempty"`
	// seLinuxMount specifies if the CSI driver supports "-o context"
	// mount option.
	// When "true", the CSI driver must ensure that all volumes provided by this CSI
	// driver can be mounted separately with different `-o context` options. This is
	// typical for storage backends that provide volumes as filesystems on block
	// devices or as independent shared volumes.
	// Kubernetes will call NodeStage / NodePublish with "-o context=xyz" mount
	// option when mounting a ReadWriteOncePod volume used in Pod that has
	// explicitly set SELinux context. In the future, it may be expanded to other
	// volume AccessModes. In any case, Kubernetes will ensure that the volume is
	// mounted only with a single SELinux context.
	// When "false", Kubernetes won't pass any special SELinux mount options to the driver.
	// This is typical for volumes that represent subdirectories of a bigger shared filesystem.
	// Default is "false".
	SELinuxMount bool `json:"seLinuxMount,omitempty"`
}

func (in *CSIDriverSpec) DeepCopyInto(out *CSIDriverSpec) {
	*out = *in
	if in.VolumeLifecycleModes != nil {
		t := make([]VolumeLifecycleMode, len(in.VolumeLifecycleModes))
		copy(t, in.VolumeLifecycleModes)
		out.VolumeLifecycleModes = t
	}
	if in.TokenRequests != nil {
		l := make([]TokenRequest, len(in.TokenRequests))
		for i := range in.TokenRequests {
			in.TokenRequests[i].DeepCopyInto(&l[i])
		}
		out.TokenRequests = l
	}
}

func (in *CSIDriverSpec) DeepCopy() *CSIDriverSpec {
	if in == nil {
		return nil
	}
	out := new(CSIDriverSpec)
	in.DeepCopyInto(out)
	return out
}

type CSINodeSpec struct {
	// drivers is a list of information of all CSI Drivers existing on a node.
	// If all drivers in the list are uninstalled, this can become empty.
	Drivers []CSINodeDriver `json:"drivers"`
}

func (in *CSINodeSpec) DeepCopyInto(out *CSINodeSpec) {
	*out = *in
	if in.Drivers != nil {
		l := make([]CSINodeDriver, len(in.Drivers))
		for i := range in.Drivers {
			in.Drivers[i].DeepCopyInto(&l[i])
		}
		out.Drivers = l
	}
}

func (in *CSINodeSpec) DeepCopy() *CSINodeSpec {
	if in == nil {
		return nil
	}
	out := new(CSINodeSpec)
	in.DeepCopyInto(out)
	return out
}

type VolumeAttachmentSpec struct {
	// attacher indicates the name of the volume driver that MUST handle this
	// request. This is the name returned by GetPluginName().
	Attacher string `json:"attacher"`
	// source represents the volume that should be attached.
	Source VolumeAttachmentSource `json:"source"`
	// nodeName represents the node that the volume should be attached to.
	NodeName string `json:"nodeName"`
}

func (in *VolumeAttachmentSpec) DeepCopyInto(out *VolumeAttachmentSpec) {
	*out = *in
	in.Source.DeepCopyInto(&out.Source)
}

func (in *VolumeAttachmentSpec) DeepCopy() *VolumeAttachmentSpec {
	if in == nil {
		return nil
	}
	out := new(VolumeAttachmentSpec)
	in.DeepCopyInto(out)
	return out
}

type VolumeAttachmentStatus struct {
	// attached indicates the volume is successfully attached.
	// This field must only be set by the entity completing the attach
	// operation, i.e. the external-attacher.
	Attached bool `json:"attached"`
	// attachmentMetadata is populated with any
	// information returned by the attach operation, upon successful attach, that must be passed
	// into subsequent WaitForAttach or Mount calls.
	// This field must only be set by the entity completing the attach
	// operation, i.e. the external-attacher.
	AttachmentMetadata map[string]string `json:"attachmentMetadata,omitempty"`
	// attachError represents the last error encountered during attach operation, if any.
	// This field must only be set by the entity completing the attach
	// operation, i.e. the external-attacher.
	AttachError *VolumeError `json:"attachError,omitempty"`
	// detachError represents the last error encountered during detach operation, if any.
	// This field must only be set by the entity completing the detach
	// operation, i.e. the external-attacher.
	DetachError *VolumeError `json:"detachError,omitempty"`
}

func (in *VolumeAttachmentStatus) DeepCopyInto(out *VolumeAttachmentStatus) {
	*out = *in
	if in.AttachmentMetadata != nil {
		in, out := &in.AttachmentMetadata, &out.AttachmentMetadata
		*out = make(map[string]string, len(*in))
		for k, v := range *in {
			(*out)[k] = v
		}
	}
	if in.AttachError != nil {
		in, out := &in.AttachError, &out.AttachError
		*out = new(VolumeError)
		(*in).DeepCopyInto(*out)
	}
	if in.DetachError != nil {
		in, out := &in.DetachError, &out.DetachError
		*out = new(VolumeError)
		(*in).DeepCopyInto(*out)
	}
}

func (in *VolumeAttachmentStatus) DeepCopy() *VolumeAttachmentStatus {
	if in == nil {
		return nil
	}
	out := new(VolumeAttachmentStatus)
	in.DeepCopyInto(out)
	return out
}

type TokenRequest struct {
	// audience is the intended audience of the token in "TokenRequestSpec".
	// It will default to the audiences of kube apiserver.
	Audience string `json:"audience"`
	// expirationSeconds is the duration of validity of the token in "TokenRequestSpec".
	// It has the same default value of "ExpirationSeconds" in "TokenRequestSpec".
	ExpirationSeconds int64 `json:"expirationSeconds,omitempty"`
}

func (in *TokenRequest) DeepCopyInto(out *TokenRequest) {
	*out = *in
}

func (in *TokenRequest) DeepCopy() *TokenRequest {
	if in == nil {
		return nil
	}
	out := new(TokenRequest)
	in.DeepCopyInto(out)
	return out
}

type CSINodeDriver struct {
	// name represents the name of the CSI driver that this object refers to.
	// This MUST be the same name returned by the CSI GetPluginName() call for
	// that driver.
	Name string `json:"name"`
	// nodeID of the node from the driver point of view.
	// This field enables Kubernetes to communicate with storage systems that do
	// not share the same nomenclature for nodes. For example, Kubernetes may
	// refer to a given node as "node1", but the storage system may refer to
	// the same node as "nodeA". When Kubernetes issues a command to the storage
	// system to attach a volume to a specific node, it can use this field to
	// refer to the node name using the ID that the storage system will
	// understand, e.g. "nodeA" instead of "node1". This field is required.
	NodeID string `json:"nodeID"`
	// topologyKeys is the list of keys supported by the driver.
	// When a driver is initialized on a cluster, it provides a set of topology
	// keys that it understands (e.g. "company.com/zone", "company.com/region").
	// When a driver is initialized on a node, it provides the same topology keys
	// along with values. Kubelet will expose these topology keys as labels
	// on its own node object.
	// When Kubernetes does topology aware provisioning, it can use this list to
	// determine which labels it should retrieve from the node object and pass
	// back to the driver.
	// It is possible for different nodes to use different topology keys.
	// This can be empty if driver does not support topology.
	TopologyKeys []string `json:"topologyKeys"`
	// allocatable represents the volume resources of a node that are available for scheduling.
	// This field is beta.
	Allocatable *VolumeNodeResources `json:"allocatable,omitempty"`
}

func (in *CSINodeDriver) DeepCopyInto(out *CSINodeDriver) {
	*out = *in
	if in.TopologyKeys != nil {
		t := make([]string, len(in.TopologyKeys))
		copy(t, in.TopologyKeys)
		out.TopologyKeys = t
	}
	if in.Allocatable != nil {
		in, out := &in.Allocatable, &out.Allocatable
		*out = new(VolumeNodeResources)
		(*in).DeepCopyInto(*out)
	}
}

func (in *CSINodeDriver) DeepCopy() *CSINodeDriver {
	if in == nil {
		return nil
	}
	out := new(CSINodeDriver)
	in.DeepCopyInto(out)
	return out
}

type VolumeAttachmentSource struct {
	// persistentVolumeName represents the name of the persistent volume to attach.
	PersistentVolumeName string `json:"persistentVolumeName,omitempty"`
	// inlineVolumeSpec contains all the information necessary to attach
	// a persistent volume defined by a pod's inline VolumeSource. This field
	// is populated only for the CSIMigration feature. It contains
	// translated fields from a pod's inline VolumeSource to a
	// PersistentVolumeSpec. This field is beta-level and is only
	// honored by servers that enabled the CSIMigration feature.
	InlineVolumeSpec *corev1.PersistentVolumeSpec `json:"inlineVolumeSpec,omitempty"`
}

func (in *VolumeAttachmentSource) DeepCopyInto(out *VolumeAttachmentSource) {
	*out = *in
	if in.InlineVolumeSpec != nil {
		in, out := &in.InlineVolumeSpec, &out.InlineVolumeSpec
		*out = new(corev1.PersistentVolumeSpec)
		(*in).DeepCopyInto(*out)
	}
}

func (in *VolumeAttachmentSource) DeepCopy() *VolumeAttachmentSource {
	if in == nil {
		return nil
	}
	out := new(VolumeAttachmentSource)
	in.DeepCopyInto(out)
	return out
}

type VolumeError struct {
	// time represents the time the error was encountered.
	Time *metav1.Time `json:"time,omitempty"`
	// message represents the error encountered during Attach or Detach operation.
	// This string may be logged, so it should not contain sensitive
	// information.
	Message string `json:"message,omitempty"`
}

func (in *VolumeError) DeepCopyInto(out *VolumeError) {
	*out = *in
	if in.Time != nil {
		in, out := &in.Time, &out.Time
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

func (in *VolumeError) DeepCopy() *VolumeError {
	if in == nil {
		return nil
	}
	out := new(VolumeError)
	in.DeepCopyInto(out)
	return out
}

type VolumeNodeResources struct {
	// count indicates the maximum number of unique volumes managed by the CSI driver that can be used on a node.
	// A volume that is both attached and mounted on a node is considered to be used once, not twice.
	// The same rule applies for a unique volume that is shared among multiple pods on the same node.
	// If this field is not specified, then the supported number of volumes on this node is unbounded.
	Count int `json:"count,omitempty"`
}

func (in *VolumeNodeResources) DeepCopyInto(out *VolumeNodeResources) {
	*out = *in
}

func (in *VolumeNodeResources) DeepCopy() *VolumeNodeResources {
	if in == nil {
		return nil
	}
	out := new(VolumeNodeResources)
	in.DeepCopyInto(out)
	return out
}
