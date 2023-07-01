package metav1

import (
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

var Unversioned = schema.GroupVersion{Group: "", Version: "v1"}

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
	return Time{
		Seconds: t.Unix(),
		Nanos:   t.Nanosecond(),
	}
}

func (in *Time) IsZero() bool {
	if in == nil {
		return true
	}
	t := time.Unix(in.Seconds, int64(in.Nanos))
	return t.IsZero()
}

func (in *Time) Before(u *Time) bool {
	if in != nil && u != nil {
		return in.Time().Before(u.Time())
	}
	return false
}

func (in *Time) After(u *Time) bool {
	if in != nil && u != nil {
		return in.Time().After(u.Time())
	}
	return false
}

func (in *Time) Equal(u *Time) bool {
	if in == nil && u == nil {
		return true
	}
	if in != nil && u != nil {
		return in.Time().Equal(u.Time())
	}
	return false
}

func (in *Time) Time() time.Time {
	return time.Unix(in.Seconds, int64(in.Nanos))
}

func (in *Time) Unix() int64 {
	return in.Time().Unix()
}
