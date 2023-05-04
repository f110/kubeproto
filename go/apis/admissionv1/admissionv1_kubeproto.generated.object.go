package admissionv1

import (
	"go.f110.dev/kubeproto/go/apis/metav1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "admission.k8s.io"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}

type Operation string

const (
	OperationCREATE  Operation = "CREATE"
	OperationUPDATE  Operation = "UPDATE"
	OperationDELETE  Operation = "DELETE"
	OperationCONNECT Operation = "CONNECT"
)

type PatchType string

const (
	PatchTypeJSONPatch PatchType = "JSONPatch"
)
