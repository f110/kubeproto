package k8stestingclient

import (
	"testing"
	"time"

	"go.f110.dev/kubeproto/go/apis/corev1"
	"go.f110.dev/kubeproto/go/apis/metav1"
	"go.f110.dev/kubeproto/go/internal/assertion"
	"go.f110.dev/kubeproto/go/k8sclient"
	"k8s.io/apimachinery/pkg/labels"
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
