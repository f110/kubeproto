package apis

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "grafana.f110.dev"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1alpha1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	schemaGroupVersion = schema.GroupVersion{Group: "grafana.f110.dev", Version: "v1alpha1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(schemaGroupVersion,
		&Grafana{},
		&GrafanaUser{},
	)
	metav1.AddToGroupVersion(scheme, schemaGroupVersion)
	return nil
}
