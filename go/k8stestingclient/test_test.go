package k8stestingclient

import (
	"encoding/json"
	"testing"
	"time"

	"go.f110.dev/kubeproto/go/apis/appsv1"
	"go.f110.dev/kubeproto/go/apis/corev1"
	"go.f110.dev/kubeproto/go/apis/metav1"
	"go.f110.dev/kubeproto/go/internal/assertion"
	"go.f110.dev/kubeproto/go/k8sclient"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	k8stesting "k8s.io/client-go/testing"
)

func TestTestingClient(t *testing.T) {
	s := NewSet()
	err := s.Tracker().Add(&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test-1", Namespace: metav1.NamespaceDefault}})
	assertion.MustNoError(t, err)
	err = s.Tracker().Add(&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test-2", Namespace: metav1.NamespaceDefault}})
	assertion.MustNoError(t, err)

	pod, err := s.CoreV1.GetPod(t.Context(), metav1.NamespaceDefault, "test-1", metav1.GetOptions{})
	assertion.MustNoError(t, err)
	assertion.Equal(t, "test-1", pod.Name)
	pods, err := s.CoreV1.ListPod(t.Context(), metav1.NamespaceDefault, metav1.ListOptions{})
	assertion.MustNoError(t, err)
	assertion.Len(t, pods.Items, 2)
	err = s.CoreV1.DeletePod(t.Context(), metav1.NamespaceDefault, "test-1", metav1.DeleteOptions{})
	assertion.MustNoError(t, err)
	_, err = s.CoreV1.WatchPod(t.Context(), metav1.NamespaceDefault, metav1.ListOptions{})
	assertion.MustNoError(t, err)

	sharedInformerFactory := k8sclient.NewInformerFactory(&s.Set, k8sclient.NewInformerCache(), metav1.NamespaceAll, 30*time.Second)
	err = sharedInformerFactory.InformerFor(&corev1.Pod{}).GetIndexer().Add(&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test-1", Namespace: metav1.NamespaceDefault}})
	assertion.MustNoError(t, err)
	err = sharedInformerFactory.InformerFor(&corev1.Pod{}).GetIndexer().Add(&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test-2", Namespace: metav1.NamespaceDefault}})
	assertion.MustNoError(t, err)
	corev1Informer := k8sclient.NewCoreV1Informer(sharedInformerFactory.Cache(), s.CoreV1, metav1.NamespaceAll, 30*time.Second)
	podsFromLister, err := corev1Informer.PodLister().List(metav1.NamespaceDefault, labels.Everything())
	assertion.MustNoError(t, err)
	assertion.Len(t, podsFromLister, 2)
}

func TestApply(t *testing.T) {
	s := NewSet()
	err := s.Tracker().Add(&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test-1", Namespace: metav1.NamespaceDefault}})
	assertion.MustNoError(t, err)

	pod, err := s.CoreV1.ApplyPod(
		t.Context(),
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test-1", Namespace: metav1.NamespaceDefault, Labels: map[string]string{"app": "test"}}},
		metav1.ApplyOptions{FieldManager: "test", Force: true},
	)
	assertion.MustNoError(t, err)
	assertion.Equal(t, "test", pod.Labels["app"])

	actions := s.Actions()
	assertion.Len(t, actions, 1)
	action, ok := actions[0].(k8stesting.PatchAction)
	assertion.Equal(t, true, ok)
	assertion.Equal(t, types.ApplyPatchType, action.GetPatchType())
	assertion.Equal(t, "test-1", action.GetName())
	assertion.Equal(t, "pods", action.GetResource().Resource)
	patchOpts := action.(k8stesting.PatchActionImpl).PatchOptions
	assertion.Equal(t, "test", patchOpts.FieldManager)
	assertion.Equal(t, true, *patchOpts.Force)

	var applied map[string]any
	assertion.MustNoError(t, json.Unmarshal(action.GetPatch(), &applied))
	assertion.Equal(t, "v1", applied["apiVersion"])
	assertion.Equal(t, "Pod", applied["kind"])
}

func TestApplyWithFieldManager(t *testing.T) {
	s := NewSet()
	deployment := func(containers ...corev1.Container) *appsv1.Deployment {
		return &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: metav1.NamespaceDefault},
			Spec:       &appsv1.DeploymentSpec{Template: corev1.PodTemplateSpec{Spec: &corev1.PodSpec{Containers: containers}}},
		}
	}

	// Apply creates the object when it doesn't exist.
	d, err := s.AppsV1.ApplyDeployment(t.Context(), deployment(corev1.Container{Name: "app", Image: "app:v1"}), metav1.ApplyOptions{FieldManager: "operator"})
	assertion.MustNoError(t, err)
	assertion.Len(t, d.ManagedFields, 1)
	assertion.Equal(t, "operator", d.ManagedFields[0].Manager)

	// kubectl rollout restart adds the annotation to the template.
	restart := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: metav1.NamespaceDefault},
		Spec: &appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: &metav1.ObjectMeta{Annotations: map[string]string{"kubectl.kubernetes.io/restartedAt": "2026-11-10T00:00:00Z"}},
			},
		},
	}
	_, err = s.AppsV1.ApplyDeployment(t.Context(), restart, metav1.ApplyOptions{FieldManager: "kubectl"})
	assertion.MustNoError(t, err)
	// Another controller injects the sidecar container.
	_, err = s.AppsV1.ApplyDeployment(t.Context(), deployment(corev1.Container{Name: "sidecar", Image: "sidecar:v1"}), metav1.ApplyOptions{FieldManager: "injector"})
	assertion.MustNoError(t, err)

	// The fields of the other managers are kept when the operator applies again.
	d, err = s.AppsV1.ApplyDeployment(t.Context(), deployment(corev1.Container{Name: "app", Image: "app:v2"}), metav1.ApplyOptions{FieldManager: "operator"})
	assertion.MustNoError(t, err)
	assertion.Equal(t, "2026-11-10T00:00:00Z", d.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"])
	assertion.Len(t, d.Spec.Template.Spec.Containers, 2)
	for _, c := range d.Spec.Template.Spec.Containers {
		switch c.Name {
		case "app":
			assertion.Equal(t, "app:v2", c.Image)
		case "sidecar":
			assertion.Equal(t, "sidecar:v1", c.Image)
		default:
			t.Errorf("unexpected container: %s", c.Name)
		}
	}

	// Changing the field owned by the other manager is a conflict.
	_, err = s.AppsV1.ApplyDeployment(t.Context(), deployment(corev1.Container{Name: "app", Image: "app:v3"}), metav1.ApplyOptions{FieldManager: "injector"})
	assertion.Equal(t, true, apierrors.IsConflict(err))
}

func TestApplyIgnoresStatusSubresource(t *testing.T) {
	s := NewSet()
	err := s.Tracker().Add(&appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: metav1.NamespaceDefault},
		Status:     &appsv1.DeploymentStatus{Replicas: 3},
	})
	assertion.MustNoError(t, err)

	d, err := s.AppsV1.ApplyDeployment(
		t.Context(),
		&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: metav1.NamespaceDefault, Labels: map[string]string{"app": "test"}},
			Status:     &appsv1.DeploymentStatus{Replicas: 1},
		},
		metav1.ApplyOptions{FieldManager: "test"},
	)
	assertion.MustNoError(t, err)
	assertion.Equal(t, "test", d.Labels["app"])
	assertion.Equal(t, 3, d.Status.Replicas)

	d, err = s.AppsV1.GetDeployment(t.Context(), metav1.NamespaceDefault, "test", metav1.GetOptions{})
	assertion.MustNoError(t, err)
	assertion.Equal(t, 3, d.Status.Replicas)
}
