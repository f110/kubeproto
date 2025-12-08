package corev1

import (
	"errors"
	"fmt"

	"go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	EventTypeNormal  string = "Normal"
	EventTypeWarning string = "Warning"
)

func GetObjectReference(scheme *runtime.Scheme, obj runtime.Object) (*ObjectReference, error) {
	if obj == nil {
		return nil, errors.New("can't reference a nil object")
	}

	var listMeta *metav1.ListMeta
	var objectMeta *metav1.ObjectMeta
	if v, ok := obj.(interface{ GetListMeta() *metav1.ListMeta }); ok {
		listMeta = v.GetListMeta()
	} else if v, ok := obj.(interface{ GetObjectMeta() *metav1.ObjectMeta }); ok {
		objectMeta = v.GetObjectMeta()
	}

	gvk := obj.GetObjectKind().GroupVersionKind()
	if gvk.Empty() {
		gvks, _, err := scheme.ObjectKinds(obj)
		if err != nil {
			return nil, err
		}
		if len(gvks) == 0 || gvks[0].Empty() {
			return nil, fmt.Errorf("unexpected gvks registered for object %T: %v", obj, gvks)
		}
		gvk = gvks[0]
	}
	kind := gvk.Kind
	version := gvk.GroupVersion().String()

	// only has list metadata
	if objectMeta == nil {
		return &ObjectReference{
			Kind:            kind,
			APIVersion:      version,
			ResourceVersion: listMeta.ResourceVersion,
		}, nil
	}

	return &ObjectReference{
		Kind:            kind,
		APIVersion:      version,
		Name:            objectMeta.Name,
		Namespace:       objectMeta.Namespace,
		UID:             objectMeta.UID,
		ResourceVersion: objectMeta.ResourceVersion,
	}, nil
}

func (in *Pod) GetObjectKind() schema.ObjectKind {
	return &in.TypeMeta
}
