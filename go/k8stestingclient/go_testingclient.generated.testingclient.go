package k8stestingclient

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/watch"
	k8stesting "k8s.io/client-go/testing"

	"go.f110.dev/kubeproto/go/apis/metav1"
	"go.f110.dev/kubeproto/go/k8sclient"
)

var (
	codecs = serializer.NewCodecFactory(k8sclient.Scheme)
)

type Set struct {
	k8sclient.Set

	fake    k8stesting.Fake
	tracker k8stesting.ObjectTracker
}

func NewSet() *Set {
	s := &Set{}
	s.tracker = k8stesting.NewObjectTracker(k8sclient.Scheme, codecs.UniversalDecoder())
	s.fake.AddReactor("*", "*", k8stesting.ObjectReaction(s.tracker))
	s.fake.AddWatchReactor("*", func(action k8stesting.Action) (handled bool, ret watch.Interface, err error) {
		w, err := s.tracker.Watch(action.GetResource(), action.GetNamespace())
		if err != nil {
			return false, nil, err
		}
		return true, w, nil
	})

	s.CoreV1 = k8sclient.NewCoreV1Client(&fakerBackend{fake: &s.fake})
	s.AdmissionregistrationK8sIoV1 = k8sclient.NewAdmissionregistrationK8sIoV1Client(&fakerBackend{fake: &s.fake})
	s.AppsV1 = k8sclient.NewAppsV1Client(&fakerBackend{fake: &s.fake})
	s.AuthenticationK8sIoV1 = k8sclient.NewAuthenticationK8sIoV1Client(&fakerBackend{fake: &s.fake})
	s.AuthorizationK8sIoV1 = k8sclient.NewAuthorizationK8sIoV1Client(&fakerBackend{fake: &s.fake})
	s.AutoscalingV1 = k8sclient.NewAutoscalingV1Client(&fakerBackend{fake: &s.fake})
	s.AutoscalingV2 = k8sclient.NewAutoscalingV2Client(&fakerBackend{fake: &s.fake})
	s.BatchV1 = k8sclient.NewBatchV1Client(&fakerBackend{fake: &s.fake})
	s.CertificatesK8sIoV1 = k8sclient.NewCertificatesK8sIoV1Client(&fakerBackend{fake: &s.fake})
	s.CoordinationK8sIoV1 = k8sclient.NewCoordinationK8sIoV1Client(&fakerBackend{fake: &s.fake})
	s.DiscoveryK8sIoV1 = k8sclient.NewDiscoveryK8sIoV1Client(&fakerBackend{fake: &s.fake})
	s.EventsK8sIoV1 = k8sclient.NewEventsK8sIoV1Client(&fakerBackend{fake: &s.fake})
	s.NetworkingK8sIoV1 = k8sclient.NewNetworkingK8sIoV1Client(&fakerBackend{fake: &s.fake})
	s.PolicyV1 = k8sclient.NewPolicyV1Client(&fakerBackend{fake: &s.fake})
	s.RbacAuthorizationK8sIoV1 = k8sclient.NewRbacAuthorizationK8sIoV1Client(&fakerBackend{fake: &s.fake})
	s.SchedulingK8sIoV1 = k8sclient.NewSchedulingK8sIoV1Client(&fakerBackend{fake: &s.fake})
	s.StorageK8sIoV1 = k8sclient.NewStorageK8sIoV1Client(&fakerBackend{fake: &s.fake})
	return s
}

func (s *Set) Tracker() k8stesting.ObjectTracker {
	return s.tracker
}

func (s *Set) Actions() []k8stesting.Action {
	return s.fake.Actions()
}

type fakerBackend struct {
	fake *k8stesting.Fake
}

func (f *fakerBackend) Get(ctx context.Context, resourceName, namespace, name string, opts metav1.GetOptions, result runtime.Object) (runtime.Object, error) {
	gvks, _, err := k8sclient.Scheme.ObjectKinds(result)
	if err != nil {
		return nil, err
	}
	gvk := gvks[0]
	obj, err := f.fake.Invokes(k8stesting.NewGetAction(gvk.GroupVersion().WithResource(resourceName), namespace, name), result)
	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), nil
}

func (f *fakerBackend) List(ctx context.Context, resourceName, namespace string, opts metav1.ListOptions, result runtime.Object) (runtime.Object, error) {
	gvks, _, err := k8sclient.Scheme.ObjectKinds(result)
	if err != nil {
		return nil, err
	}
	gvk := gvks[0]
	obj, err := f.fake.Invokes(k8stesting.NewListAction(gvk.GroupVersion().WithResource(resourceName), gvk, namespace, opts), result)

	if obj == nil {
		return nil, err
	}

	label, _, _ := k8stesting.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	objs, err := meta.ExtractList(obj)
	if err != nil {
		return nil, err
	}
	filtered := make([]runtime.Object, 0)
	for _, item := range objs {
		m := item.(metav1.Object)
		objMeta := m.GetObjectMeta()
		if label.Matches(labels.Set(objMeta.Labels)) {
			filtered = append(filtered, item)
		}
	}
	if err := meta.SetList(obj, filtered); err != nil {
		return nil, err
	}
	return obj.DeepCopyObject(), err
}

func (f *fakerBackend) Create(ctx context.Context, resourceName string, obj runtime.Object, opts metav1.CreateOptions, result runtime.Object) (runtime.Object, error) {
	gvks, _, err := k8sclient.Scheme.ObjectKinds(result)
	if err != nil {
		return nil, err
	}
	gvk := gvks[0]
	m := obj.(metav1.Object)
	objMeta := m.GetObjectMeta()
	obj, err = f.fake.Invokes(k8stesting.NewCreateAction(gvk.GroupVersion().WithResource(resourceName), objMeta.Namespace, obj), result)

	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), err
}

func (f *fakerBackend) Update(ctx context.Context, resourceName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	gvks, _, err := k8sclient.Scheme.ObjectKinds(result)
	if err != nil {
		return nil, err
	}
	gvk := gvks[0]
	m := obj.(metav1.Object)
	objMeta := m.GetObjectMeta()
	obj, err = f.fake.Invokes(k8stesting.NewUpdateAction(gvk.GroupVersion().WithResource(resourceName), objMeta.Namespace, obj), result)

	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), err
}
func (f *fakerBackend) UpdateStatus(ctx context.Context, resourceName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	gvks, _, err := k8sclient.Scheme.ObjectKinds(result)
	if err != nil {
		return nil, err
	}
	gvk := gvks[0]
	m := obj.(metav1.Object)
	objMeta := m.GetObjectMeta()
	obj, err = f.fake.Invokes(k8stesting.NewUpdateSubresourceAction(gvk.GroupVersion().WithResource(resourceName), "status", objMeta.Namespace, obj), result)

	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), err
}
func (f *fakerBackend) Delete(ctx context.Context, gvr schema.GroupVersionResource, namespace, name string, opts metav1.DeleteOptions) error {
	_, err := f.fake.Invokes(k8stesting.NewDeleteAction(gvr, namespace, name), nil)

	return err
}
func (f *fakerBackend) Watch(ctx context.Context, gvr schema.GroupVersionResource, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return f.fake.InvokesWatch(k8stesting.NewWatchAction(gvr, namespace, opts))
}
func (f *fakerBackend) GetClusterScoped(ctx context.Context, resourceName, name string, opts metav1.GetOptions, result runtime.Object) (runtime.Object, error) {
	return f.Get(ctx, resourceName, "", name, opts, result)
}

func (f *fakerBackend) ListClusterScoped(ctx context.Context, resourceName string, opts metav1.ListOptions, result runtime.Object) (runtime.Object, error) {
	return f.List(ctx, resourceName, "", opts, result)
}

func (f *fakerBackend) CreateClusterScoped(ctx context.Context, resourceName string, obj runtime.Object, opts metav1.CreateOptions, result runtime.Object) (runtime.Object, error) {
	return f.Create(ctx, resourceName, obj, opts, result)
}

func (f *fakerBackend) UpdateClusterScoped(ctx context.Context, resourceName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	return f.Update(ctx, resourceName, obj, opts, result)
}

func (f *fakerBackend) UpdateStatusClusterScoped(ctx context.Context, resourceName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	return f.UpdateStatus(ctx, resourceName, obj, opts, result)
}

func (f *fakerBackend) DeleteClusterScoped(ctx context.Context, gvr schema.GroupVersionResource, name string, opts metav1.DeleteOptions) error {
	return f.Delete(ctx, gvr, "", name, opts)
}

func (f *fakerBackend) WatchClusterScoped(ctx context.Context, gvr schema.GroupVersionResource, opts metav1.ListOptions) (watch.Interface, error) {
	return f.Watch(ctx, gvr, "", opts)
}
