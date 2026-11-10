package corev1

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime"

	"go.f110.dev/kubeproto/go/internal/assertion"
)

func TestOpenAPIModelName(t *testing.T) {
	scheme := runtime.NewScheme()
	assertion.MustNoError(t, AddToScheme(scheme))

	name, err := scheme.ToOpenAPIDefinitionName(SchemaGroupVersion.WithKind("Pod"))
	assertion.MustNoError(t, err)
	assertion.Equal(t, "k8s.io.api.core.v1.Pod", name)
}
