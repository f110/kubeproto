package metav1

import (
	"fmt"
	"time"

	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

var Unversioned = schema.GroupVersion{Group: "", Version: "v1"}

const (
	NamespaceDefault = "default"
	NamespaceAll     = ""
	NamespaceNone    = ""
	NamespaceSystem  = "kube-system"
	NamespacePublic  = "kube-public"
)

func (in *TypeMeta) GetObjectKind() schema.ObjectKind { return in }

func (in *TypeMeta) SetGroupVersionKind(gvk schema.GroupVersionKind) {
	in.APIVersion, in.Kind = gvk.ToAPIVersionAndKind()
}

func (in *TypeMeta) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(in.APIVersion, in.Kind)
}

type InternalEvent watch.Event

func (e *InternalEvent) GetObjectKind() schema.ObjectKind { return schema.EmptyObjectKind }

func (e *InternalEvent) DeepCopyObject() runtime.Object {
	if c := e.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

func (e *InternalEvent) DeepCopy() *InternalEvent {
	if e == nil {
		return nil
	}
	out := new(InternalEvent)
	e.DeepCopyInto(out)
	return out
}

func (e *InternalEvent) DeepCopyInto(out *InternalEvent) {
	*out = *e
	if e.Object != nil {
		out.Object = e.Object.DeepCopyObject()
	}
	return
}

func (in *WatchEvent) GetObjectKind() schema.ObjectKind { return schema.EmptyObjectKind }

func (in *WatchEvent) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func AddToGroupVersion(scheme *runtime.Scheme, groupVersion schema.GroupVersion) {
	scheme.AddKnownTypeWithName(groupVersion.WithKind("WatchEventKind"), &WatchEvent{})
	scheme.AddKnownTypeWithName(
		schema.GroupVersion{Group: groupVersion.Group, Version: runtime.APIVersionInternal}.WithKind("WatchEventKind"),
		&InternalEvent{},
	)
	// Supports legacy code paths, most callers should use metav1.ParameterCodec for now
	scheme.AddKnownTypes(groupVersion, &ListOptions{},
		&GetOptions{},
		&DeleteOptions{},
		&CreateOptions{},
		&UpdateOptions{},
		&PatchOptions{},
	)
	// Register Unversioned types under their own special group
	scheme.AddUnversionedTypes(Unversioned,
		&Status{},
		&APIVersions{},
		&APIGroupList{},
		&APIGroup{},
		&APIResourceList{},
	)
}

func Now() Time {
	return NewTime(time.Now())
}

func NewTime(t time.Time) Time {
	return Time{Time: t}
}

func (in *Time) IsZero() bool {
	if in == nil {
		return true
	}
	return in.Time.IsZero()
}

func (in *Time) Equal(u *Time) bool {
	if in == nil && u == nil {
		return true
	}
	if in != nil && u != nil {
		return in.Time.Equal(u.Time)
	}
	return false
}

func (in *Duration) TimeDuration() time.Duration {
	return time.Duration(in.Duration)
}

type Object interface {
	GetObjectMeta() *ObjectMeta
}

type ListObject interface {
	GetListMeta() *ListMeta
}

func NewControllerRef(owner ObjectMeta, gvk schema.GroupVersionKind) OwnerReference {
	return OwnerReference{
		APIVersion:         gvk.GroupVersion().String(),
		Kind:               gvk.Kind,
		Name:               owner.Name,
		UID:                owner.UID,
		BlockOwnerDeletion: true,
		Controller:         true,
	}
}

func IsControlledBy(obj Object, owner Object) bool {
	objMeta := obj.GetObjectMeta()
	ownerMeta := owner.GetObjectMeta()
	for _, v := range objMeta.OwnerReferences {
		if v.Controller && v.UID == ownerMeta.UID {
			return true
		}
	}
	return false
}

func HasAnnotation(obj ObjectMeta, key string) bool {
	_, ok := obj.Annotations[key]
	return ok
}

func SetMetadataAnnotation(obj *ObjectMeta, key string, value string) {
	if obj.Annotations == nil {
		obj.Annotations = make(map[string]string)
	}
	obj.Annotations[key] = value
}

func (in *ObjectMeta) GetObjectMeta() *ObjectMeta {
	return in
}

func (in *ObjectMeta) OwnedBy(ref OwnerReference) {
	for _, v := range in.OwnerReferences {
		if v.UID == ref.UID {
			return
		}
	}
	in.OwnerReferences = append(in.OwnerReferences, ref)
}

func (in *ObjectMeta) GetNamespace() string {
	return in.Namespace
}
func (in *ObjectMeta) SetNamespace(namespace string) {
	in.Namespace = namespace
}
func (in *ObjectMeta) GetName() string {
	return in.Name
}

func (in *ObjectMeta) SetName(name string) {
	in.Name = name
}

func (in *ObjectMeta) GetGenerateName() string {
	return in.GenerateName
}

func (in *ObjectMeta) SetGenerateName(name string) {
	in.GenerateName = name
}

func (in *ObjectMeta) GetUID() types.UID {
	return types.UID(in.UID)
}

func (in *ObjectMeta) SetUID(uid types.UID) {
	in.UID = string(uid)
}

func (in *ObjectMeta) GetResourceVersion() string {
	return in.ResourceVersion
}

func (in *ObjectMeta) SetResourceVersion(version string) {
	in.ResourceVersion = version
}

func (in *ObjectMeta) GetGeneration() int64 {
	return in.Generation
}

func (in *ObjectMeta) SetGeneration(generation int64) {
	in.Generation = generation
}

func (in *ObjectMeta) GetSelfLink() string {
	return in.SelfLink
}

func (in *ObjectMeta) SetSelfLink(selfLink string) {
	in.SelfLink = selfLink
}

func (in *ObjectMeta) GetCreationTimestamp() k8smetav1.Time {
	if in.CreationTimestamp == nil {
		return k8smetav1.Time{}
	}
	return k8smetav1.NewTime(in.CreationTimestamp.Time)
}

func (in *ObjectMeta) SetCreationTimestamp(timestamp k8smetav1.Time) {
	ts := NewTime(timestamp.Time)
	in.CreationTimestamp = &ts
}

func (in *ObjectMeta) GetDeletionTimestamp() *k8smetav1.Time {
	if in.DeletionTimestamp == nil {
		return nil
	}
	ts := k8smetav1.NewTime(in.DeletionTimestamp.Time)
	return &ts
}

func (in *ObjectMeta) SetDeletionTimestamp(timestamp *k8smetav1.Time) {
	ts := NewTime(timestamp.Time)
	in.DeletionTimestamp = &ts
}

func (in *ObjectMeta) GetDeletionGracePeriodSeconds() *int64 {
	return &in.DeletionGracePeriodSeconds
}

func (in *ObjectMeta) SetDeletionGracePeriodSeconds(v *int64) {
	in.DeletionGracePeriodSeconds = *v
}

func (in *ObjectMeta) GetLabels() map[string]string {
	return in.Labels
}

func (in *ObjectMeta) SetLabels(labels map[string]string) {
	in.Labels = labels
}

func (in *ObjectMeta) GetAnnotations() map[string]string {
	return in.Annotations
}

func (in *ObjectMeta) SetAnnotations(annotations map[string]string) {
	in.Annotations = annotations
}

func (in *ObjectMeta) GetFinalizers() []string {
	return in.Finalizers
}

func (in *ObjectMeta) SetFinalizers(finalizers []string) {
	in.Finalizers = finalizers
}

func (in *ObjectMeta) GetOwnerReferences() []k8smetav1.OwnerReference {
	ref := make([]k8smetav1.OwnerReference, len(in.OwnerReferences))
	for i := 0; i < len(in.OwnerReferences); i++ {
		ref[i] = k8smetav1.OwnerReference{
			APIVersion:         in.OwnerReferences[i].APIVersion,
			Kind:               in.OwnerReferences[i].Kind,
			Name:               in.OwnerReferences[i].Name,
			UID:                types.UID(in.OwnerReferences[i].UID),
			Controller:         &in.OwnerReferences[i].Controller,
			BlockOwnerDeletion: &in.OwnerReferences[i].BlockOwnerDeletion,
		}
	}
	return ref
}

func (in *ObjectMeta) SetOwnerReferences(ownerReferences []k8smetav1.OwnerReference) {
	ref := make([]OwnerReference, len(ownerReferences))
	for i := 0; i < len(ownerReferences); i++ {
		controller := false
		if ownerReferences[i].Controller != nil {
			controller = *ownerReferences[i].Controller
		}
		blockOwnerDeletion := false
		if ownerReferences[i].BlockOwnerDeletion != nil {
			blockOwnerDeletion = *ownerReferences[i].BlockOwnerDeletion
		}
		ref[i] = OwnerReference{
			APIVersion:         ownerReferences[i].APIVersion,
			Kind:               ownerReferences[i].Kind,
			Name:               ownerReferences[i].Name,
			UID:                string(ownerReferences[i].UID),
			Controller:         controller,
			BlockOwnerDeletion: blockOwnerDeletion,
		}
	}
	in.OwnerReferences = ref
}

func (in *ObjectMeta) GetManagedFields() []k8smetav1.ManagedFieldsEntry {
	entries := make([]k8smetav1.ManagedFieldsEntry, len(in.ManagedFields))
	for i := 0; i < len(in.ManagedFields); i++ {
		var t *k8smetav1.Time
		if in.ManagedFields[i].Time != nil {
			ts := k8smetav1.NewTime(in.ManagedFields[i].Time.Time)
			t = &ts
		}
		var fields *k8smetav1.FieldsV1
		if in.ManagedFields[i].FieldsV1 != nil {
			fields = &k8smetav1.FieldsV1{Raw: in.ManagedFields[i].FieldsV1.Raw}
		}
		entries[i] = k8smetav1.ManagedFieldsEntry{
			Manager:     in.ManagedFields[i].Manager,
			Operation:   k8smetav1.ManagedFieldsOperationType(in.ManagedFields[i].Operation),
			APIVersion:  in.ManagedFields[i].APIVersion,
			Time:        t,
			FieldsType:  in.ManagedFields[i].FieldsType,
			FieldsV1:    fields,
			Subresource: in.ManagedFields[i].Subresource,
		}
	}
	return entries
}

func (in *ObjectMeta) SetManagedFields(managedFields []k8smetav1.ManagedFieldsEntry) {
	entries := make([]ManagedFieldsEntry, len(managedFields))
	for i := 0; i < len(managedFields); i++ {
		var t *Time
		if managedFields[i].Time != nil {
			ts := NewTime(managedFields[i].Time.Time)
			t = &ts
		}
		var fields *FieldsV1
		if managedFields[i].FieldsV1 != nil {
			fields = &FieldsV1{Raw: managedFields[i].FieldsV1.Raw}
		}
		entries[i] = ManagedFieldsEntry{
			Manager:     managedFields[i].Manager,
			Operation:   ManagedFieldsOperationType(managedFields[i].Operation),
			APIVersion:  managedFields[i].APIVersion,
			Time:        t,
			FieldsType:  managedFields[i].FieldsType,
			FieldsV1:    fields,
			Subresource: managedFields[i].Subresource,
		}
	}
	in.ManagedFields = entries
}

func (in *ListMeta) GetListMeta() *ListMeta {
	return in
}

func LabelSelectorAsSelector(ps *LabelSelector) (labels.Selector, error) {
	if ps == nil {
		return labels.Nothing(), nil
	}
	if len(ps.MatchLabels)+len(ps.MatchExpressions) == 0 {
		return labels.Everything(), nil
	}
	requirements := make([]labels.Requirement, 0, len(ps.MatchLabels)+len(ps.MatchExpressions))
	for k, v := range ps.MatchLabels {
		r, err := labels.NewRequirement(k, selection.Equals, []string{v})
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, *r)
	}
	for _, expr := range ps.MatchExpressions {
		var op selection.Operator
		switch expr.Operator {
		case LabelSelectorOperatorIn:
			op = selection.In
		case LabelSelectorOperatorNotIn:
			op = selection.NotIn
		case LabelSelectorOperatorExists:
			op = selection.Exists
		case LabelSelectorOperatorDoesNotExist:
			op = selection.DoesNotExist
		default:
			return nil, fmt.Errorf("%q is not a valid label selector operator", expr.Operator)
		}
		r, err := labels.NewRequirement(expr.Key, op, append([]string(nil), expr.Values...))
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, *r)
	}
	selector := labels.NewSelector()
	selector = selector.Add(requirements...)
	return selector, nil
}
