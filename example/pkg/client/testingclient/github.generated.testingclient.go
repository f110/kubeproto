package testingclient

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/watch"
	k8stesting "k8s.io/client-go/testing"

	"go.f110.dev/kubeproto/example/pkg/apis/githubv1alpha1"
	"go.f110.dev/kubeproto/example/pkg/apis/miniov1alpha1"
	"go.f110.dev/kubeproto/example/pkg/client"
)

var (
	codecs = serializer.NewCodecFactory(client.Scheme)
)

type Set struct {
	client.Set

	fake    k8stesting.Fake
	tracker k8stesting.ObjectTracker
}

func NewSet() *Set {
	s := &Set{}
	o := k8stesting.NewObjectTracker(client.Scheme, codecs.UniversalDecoder())
	s.fake.AddReactor("*", "*", k8stesting.ObjectReaction(o))
	s.fake.AddWatchReactor("*", func(action k8stesting.Action) (handled bool, ret watch.Interface, err error) {
		w, err := o.Watch(action.GetResource(), action.GetNamespace())
		if err != nil {
			return false, nil, err
		}
		return true, w, nil
	})

	s.GrafanaV1alpha1 = client.NewGrafanaV1alpha1Client(&fakerBackend{fake: &s.fake})
	s.GrafanaV1alpha2 = client.NewGrafanaV1alpha2Client(&fakerBackend{fake: &s.fake})
	s.MinioV1alpha1 = client.NewMinioV1alpha1Client(&fakerBackend{fake: &s.fake})
	return s
}

func (s *Set) Tracker() k8stesting.ObjectTracker {
	return s.tracker
}

type fakerBackend struct {
	fake *k8stesting.Fake
}

func (f *fakerBackend) Get(ctx context.Context, resourceName, kindName, namespace, name string, opts metav1.GetOptions, result runtime.Object) (runtime.Object, error) {
	obj, err := f.fake.Invokes(k8stesting.NewGetAction(githubv1alpha1.SchemaGroupVersion.WithResource(resourceName), namespace, name), result)
	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), nil
}

func (f *fakerBackend) List(ctx context.Context, resourceName, kindName, namespace string, opts metav1.ListOptions, result runtime.Object) (runtime.Object, error) {
	obj, err := f.fake.
		Invokes(k8stesting.NewListAction(githubv1alpha1.SchemaGroupVersion.WithResource(resourceName), githubv1alpha1.SchemaGroupVersion.WithKind(kindName), namespace, opts), result)

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
		if label.Matches(labels.Set(m.GetLabels())) {
			filtered = append(filtered, item)
		}
	}
	if err := meta.SetList(obj, filtered); err != nil {
		return nil, err
	}
	return obj.DeepCopyObject(), err
}

func (f *fakerBackend) Create(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.CreateOptions, result runtime.Object) (runtime.Object, error) {
	m := obj.(metav1.Object)
	obj, err := f.fake.
		Invokes(k8stesting.NewCreateAction(githubv1alpha1.SchemaGroupVersion.WithResource(resourceName), m.GetNamespace(), obj), result)

	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), err
}

func (f *fakerBackend) Update(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	m := obj.(metav1.Object)
	obj, err := f.fake.
		Invokes(k8stesting.NewUpdateAction(githubv1alpha1.SchemaGroupVersion.WithResource(resourceName), m.GetNamespace(), obj), result)

	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), err

}

func (f *fakerBackend) UpdateStatus(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	m := obj.(metav1.Object)
	obj, err := f.fake.
		Invokes(k8stesting.NewUpdateSubresourceAction(miniov1alpha1.SchemaGroupVersion.WithResource(resourceName), "status", m.GetNamespace(), obj), result)

	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), err
}

func (f *fakerBackend) Delete(ctx context.Context, resourceName, kindName, namespace, name string, opts metav1.DeleteOptions) error {
	_, err := f.fake.
		Invokes(k8stesting.NewDeleteAction(miniov1alpha1.SchemaGroupVersion.WithResource(resourceName), namespace, name), nil)

	return err
}

func (f *fakerBackend) Watch(ctx context.Context, resourceName, kindName, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return f.fake.InvokesWatch(k8stesting.NewWatchAction(miniov1alpha1.SchemaGroupVersion.WithResource(resourceName), namespace, opts))
}
