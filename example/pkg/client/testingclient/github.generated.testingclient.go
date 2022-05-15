package testingclient

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/watch"
	k8stesting "k8s.io/client-go/testing"

	"go.f110.dev/kubeproto/example/pkg/apis/githubv1alpha1"
	"go.f110.dev/kubeproto/example/pkg/apis/githubv1alpha2"
	"go.f110.dev/kubeproto/example/pkg/apis/miniov1alpha1"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

func init() {
	for _, v := range []func(*runtime.Scheme) error{
		githubv1alpha1.AddToScheme,
		githubv1alpha2.AddToScheme,
		miniov1alpha1.AddToScheme,
	} {
		if err := v(scheme); err != nil {
			panic(err)
		}
	}
}

type Set struct {
	k8stesting.Fake

	GrafanaV1alpha1 *TestingGrafanaV1alpha1
	GrafanaV1alpha2 *TestingGrafanaV1alpha2
	MinioV1alpha1   *TestingMinioV1alpha1

	tracker k8stesting.ObjectTracker
}

func NewSet() *Set {
	s := &Set{}
	o := k8stesting.NewObjectTracker(scheme, codecs.UniversalDecoder())
	s.AddReactor("*", "*", k8stesting.ObjectReaction(o))
	s.AddWatchReactor("*", func(action k8stesting.Action) (handled bool, ret watch.Interface, err error) {
		w, err := o.Watch(action.GetResource(), action.GetNamespace())
		if err != nil {
			return false, nil, err
		}
		return true, w, nil
	})

	return &Set{
		GrafanaV1alpha1: NewTestingGrafanaV1alpha1Client(&s.Fake),
		GrafanaV1alpha2: NewTestingGrafanaV1alpha2Client(&s.Fake),
		MinioV1alpha1:   NewTestingMinioV1alpha1Client(&s.Fake),
		tracker:         o,
	}
}

func (s *Set) Tracker() k8stesting.ObjectTracker {
	return s.tracker
}

type TestingGrafanaV1alpha1 struct {
	*k8stesting.Fake
}

func NewTestingGrafanaV1alpha1Client(fake *k8stesting.Fake) *TestingGrafanaV1alpha1 {
	return &TestingGrafanaV1alpha1{Fake: fake}
}

func (c *TestingGrafanaV1alpha1) GetGrafana(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*githubv1alpha1.Grafana, error) {
	obj, err := c.Fake.Invokes(k8stesting.NewGetAction(githubv1alpha1.SchemaGroupVersion.WithResource("grafanas"), namespace, name), &githubv1alpha1.Grafana{})
	if obj == nil {
		return nil, err
	}
	return obj.(*githubv1alpha1.Grafana), err
}

func (c *TestingGrafanaV1alpha1) CreateGrafana(ctx context.Context, v *githubv1alpha1.Grafana, opts metav1.CreateOptions) (*githubv1alpha1.Grafana, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewCreateAction(githubv1alpha1.SchemaGroupVersion.WithResource("grafanas"), v.Namespace, v), &githubv1alpha1.Grafana{})

	if obj == nil {
		return nil, err
	}
	return obj.(*githubv1alpha1.Grafana), err
}

func (c *TestingGrafanaV1alpha1) UpdateGrafana(ctx context.Context, v *githubv1alpha1.Grafana, opts metav1.UpdateOptions) (*githubv1alpha1.Grafana, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewUpdateAction(githubv1alpha1.SchemaGroupVersion.WithResource("grafanas"), v.Namespace, v), &githubv1alpha1.Grafana{})

	if obj == nil {
		return nil, err
	}
	return obj.(*githubv1alpha1.Grafana), err
}

func (c *TestingGrafanaV1alpha1) DeleteGrafana(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(k8stesting.NewDeleteAction(githubv1alpha1.SchemaGroupVersion.WithResource("grafanas"), namespace, name), &githubv1alpha1.Grafana{})

	return err
}

func (c *TestingGrafanaV1alpha1) ListGrafana(ctx context.Context, namespace string, opts metav1.ListOptions) (*githubv1alpha1.GrafanaList, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewListAction(githubv1alpha1.SchemaGroupVersion.WithResource("grafanas"), githubv1alpha1.SchemaGroupVersion.WithKind("Grafana"), namespace, opts), &githubv1alpha1.GrafanaList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := k8stesting.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &githubv1alpha1.GrafanaList{ListMeta: obj.(*githubv1alpha1.GrafanaList).ListMeta}
	for _, item := range obj.(*githubv1alpha1.GrafanaList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

func (c *TestingGrafanaV1alpha1) WatchGrafana(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(k8stesting.NewWatchAction(githubv1alpha1.SchemaGroupVersion.WithResource("grafanas"), namespace, opts))
}

func (c *TestingGrafanaV1alpha1) GetGrafanaUser(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*githubv1alpha1.GrafanaUser, error) {
	obj, err := c.Fake.Invokes(k8stesting.NewGetAction(githubv1alpha1.SchemaGroupVersion.WithResource("grafanausers"), namespace, name), &githubv1alpha1.GrafanaUser{})
	if obj == nil {
		return nil, err
	}
	return obj.(*githubv1alpha1.GrafanaUser), err
}

func (c *TestingGrafanaV1alpha1) CreateGrafanaUser(ctx context.Context, v *githubv1alpha1.GrafanaUser, opts metav1.CreateOptions) (*githubv1alpha1.GrafanaUser, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewCreateAction(githubv1alpha1.SchemaGroupVersion.WithResource("grafanausers"), v.Namespace, v), &githubv1alpha1.GrafanaUser{})

	if obj == nil {
		return nil, err
	}
	return obj.(*githubv1alpha1.GrafanaUser), err
}

func (c *TestingGrafanaV1alpha1) UpdateGrafanaUser(ctx context.Context, v *githubv1alpha1.GrafanaUser, opts metav1.UpdateOptions) (*githubv1alpha1.GrafanaUser, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewUpdateAction(githubv1alpha1.SchemaGroupVersion.WithResource("grafanausers"), v.Namespace, v), &githubv1alpha1.GrafanaUser{})

	if obj == nil {
		return nil, err
	}
	return obj.(*githubv1alpha1.GrafanaUser), err
}

func (c *TestingGrafanaV1alpha1) DeleteGrafanaUser(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(k8stesting.NewDeleteAction(githubv1alpha1.SchemaGroupVersion.WithResource("grafanausers"), namespace, name), &githubv1alpha1.GrafanaUser{})

	return err
}

func (c *TestingGrafanaV1alpha1) ListGrafanaUser(ctx context.Context, namespace string, opts metav1.ListOptions) (*githubv1alpha1.GrafanaUserList, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewListAction(githubv1alpha1.SchemaGroupVersion.WithResource("grafanausers"), githubv1alpha1.SchemaGroupVersion.WithKind("GrafanaUser"), namespace, opts), &githubv1alpha1.GrafanaUserList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := k8stesting.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &githubv1alpha1.GrafanaUserList{ListMeta: obj.(*githubv1alpha1.GrafanaUserList).ListMeta}
	for _, item := range obj.(*githubv1alpha1.GrafanaUserList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

func (c *TestingGrafanaV1alpha1) WatchGrafanaUser(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(k8stesting.NewWatchAction(githubv1alpha1.SchemaGroupVersion.WithResource("grafanausers"), namespace, opts))
}

type TestingGrafanaV1alpha2 struct {
	*k8stesting.Fake
}

func NewTestingGrafanaV1alpha2Client(fake *k8stesting.Fake) *TestingGrafanaV1alpha2 {
	return &TestingGrafanaV1alpha2{Fake: fake}
}

func (c *TestingGrafanaV1alpha2) GetGrafana(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*githubv1alpha2.Grafana, error) {
	obj, err := c.Fake.Invokes(k8stesting.NewGetAction(githubv1alpha2.SchemaGroupVersion.WithResource("grafanas"), namespace, name), &githubv1alpha2.Grafana{})
	if obj == nil {
		return nil, err
	}
	return obj.(*githubv1alpha2.Grafana), err
}

func (c *TestingGrafanaV1alpha2) CreateGrafana(ctx context.Context, v *githubv1alpha2.Grafana, opts metav1.CreateOptions) (*githubv1alpha2.Grafana, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewCreateAction(githubv1alpha2.SchemaGroupVersion.WithResource("grafanas"), v.Namespace, v), &githubv1alpha2.Grafana{})

	if obj == nil {
		return nil, err
	}
	return obj.(*githubv1alpha2.Grafana), err
}

func (c *TestingGrafanaV1alpha2) UpdateGrafana(ctx context.Context, v *githubv1alpha2.Grafana, opts metav1.UpdateOptions) (*githubv1alpha2.Grafana, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewUpdateAction(githubv1alpha2.SchemaGroupVersion.WithResource("grafanas"), v.Namespace, v), &githubv1alpha2.Grafana{})

	if obj == nil {
		return nil, err
	}
	return obj.(*githubv1alpha2.Grafana), err
}

func (c *TestingGrafanaV1alpha2) DeleteGrafana(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(k8stesting.NewDeleteAction(githubv1alpha2.SchemaGroupVersion.WithResource("grafanas"), namespace, name), &githubv1alpha2.Grafana{})

	return err
}

func (c *TestingGrafanaV1alpha2) ListGrafana(ctx context.Context, namespace string, opts metav1.ListOptions) (*githubv1alpha2.GrafanaList, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewListAction(githubv1alpha2.SchemaGroupVersion.WithResource("grafanas"), githubv1alpha2.SchemaGroupVersion.WithKind("Grafana"), namespace, opts), &githubv1alpha2.GrafanaList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := k8stesting.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &githubv1alpha2.GrafanaList{ListMeta: obj.(*githubv1alpha2.GrafanaList).ListMeta}
	for _, item := range obj.(*githubv1alpha2.GrafanaList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

func (c *TestingGrafanaV1alpha2) WatchGrafana(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(k8stesting.NewWatchAction(githubv1alpha2.SchemaGroupVersion.WithResource("grafanas"), namespace, opts))
}

func (c *TestingGrafanaV1alpha2) GetGrafanaUser(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*githubv1alpha2.GrafanaUser, error) {
	obj, err := c.Fake.Invokes(k8stesting.NewGetAction(githubv1alpha2.SchemaGroupVersion.WithResource("grafanausers"), namespace, name), &githubv1alpha2.GrafanaUser{})
	if obj == nil {
		return nil, err
	}
	return obj.(*githubv1alpha2.GrafanaUser), err
}

func (c *TestingGrafanaV1alpha2) CreateGrafanaUser(ctx context.Context, v *githubv1alpha2.GrafanaUser, opts metav1.CreateOptions) (*githubv1alpha2.GrafanaUser, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewCreateAction(githubv1alpha2.SchemaGroupVersion.WithResource("grafanausers"), v.Namespace, v), &githubv1alpha2.GrafanaUser{})

	if obj == nil {
		return nil, err
	}
	return obj.(*githubv1alpha2.GrafanaUser), err
}

func (c *TestingGrafanaV1alpha2) UpdateGrafanaUser(ctx context.Context, v *githubv1alpha2.GrafanaUser, opts metav1.UpdateOptions) (*githubv1alpha2.GrafanaUser, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewUpdateAction(githubv1alpha2.SchemaGroupVersion.WithResource("grafanausers"), v.Namespace, v), &githubv1alpha2.GrafanaUser{})

	if obj == nil {
		return nil, err
	}
	return obj.(*githubv1alpha2.GrafanaUser), err
}

func (c *TestingGrafanaV1alpha2) DeleteGrafanaUser(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(k8stesting.NewDeleteAction(githubv1alpha2.SchemaGroupVersion.WithResource("grafanausers"), namespace, name), &githubv1alpha2.GrafanaUser{})

	return err
}

func (c *TestingGrafanaV1alpha2) ListGrafanaUser(ctx context.Context, namespace string, opts metav1.ListOptions) (*githubv1alpha2.GrafanaUserList, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewListAction(githubv1alpha2.SchemaGroupVersion.WithResource("grafanausers"), githubv1alpha2.SchemaGroupVersion.WithKind("GrafanaUser"), namespace, opts), &githubv1alpha2.GrafanaUserList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := k8stesting.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &githubv1alpha2.GrafanaUserList{ListMeta: obj.(*githubv1alpha2.GrafanaUserList).ListMeta}
	for _, item := range obj.(*githubv1alpha2.GrafanaUserList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

func (c *TestingGrafanaV1alpha2) WatchGrafanaUser(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(k8stesting.NewWatchAction(githubv1alpha2.SchemaGroupVersion.WithResource("grafanausers"), namespace, opts))
}

type TestingMinioV1alpha1 struct {
	*k8stesting.Fake
}

func NewTestingMinioV1alpha1Client(fake *k8stesting.Fake) *TestingMinioV1alpha1 {
	return &TestingMinioV1alpha1{Fake: fake}
}

func (c *TestingMinioV1alpha1) GetMinIOBucket(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*miniov1alpha1.MinIOBucket, error) {
	obj, err := c.Fake.Invokes(k8stesting.NewGetAction(miniov1alpha1.SchemaGroupVersion.WithResource("miniobuckets"), namespace, name), &miniov1alpha1.MinIOBucket{})
	if obj == nil {
		return nil, err
	}
	return obj.(*miniov1alpha1.MinIOBucket), err
}

func (c *TestingMinioV1alpha1) CreateMinIOBucket(ctx context.Context, v *miniov1alpha1.MinIOBucket, opts metav1.CreateOptions) (*miniov1alpha1.MinIOBucket, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewCreateAction(miniov1alpha1.SchemaGroupVersion.WithResource("miniobuckets"), v.Namespace, v), &miniov1alpha1.MinIOBucket{})

	if obj == nil {
		return nil, err
	}
	return obj.(*miniov1alpha1.MinIOBucket), err
}

func (c *TestingMinioV1alpha1) UpdateMinIOBucket(ctx context.Context, v *miniov1alpha1.MinIOBucket, opts metav1.UpdateOptions) (*miniov1alpha1.MinIOBucket, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewUpdateAction(miniov1alpha1.SchemaGroupVersion.WithResource("miniobuckets"), v.Namespace, v), &miniov1alpha1.MinIOBucket{})

	if obj == nil {
		return nil, err
	}
	return obj.(*miniov1alpha1.MinIOBucket), err
}

func (c *TestingMinioV1alpha1) UpdateStatusMinIOBucket(ctx context.Context, v *miniov1alpha1.MinIOBucket, opts metav1.UpdateOptions) (*miniov1alpha1.MinIOBucket, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewUpdateSubresourceAction(miniov1alpha1.SchemaGroupVersion.WithResource("miniobuckets"), "status", v.Namespace, v), &miniov1alpha1.MinIOBucket{})

	if obj == nil {
		return nil, err
	}
	return obj.(*miniov1alpha1.MinIOBucket), err
}

func (c *TestingMinioV1alpha1) DeleteMinIOBucket(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(k8stesting.NewDeleteAction(miniov1alpha1.SchemaGroupVersion.WithResource("miniobuckets"), namespace, name), &miniov1alpha1.MinIOBucket{})

	return err
}

func (c *TestingMinioV1alpha1) ListMinIOBucket(ctx context.Context, namespace string, opts metav1.ListOptions) (*miniov1alpha1.MinIOBucketList, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewListAction(miniov1alpha1.SchemaGroupVersion.WithResource("miniobuckets"), miniov1alpha1.SchemaGroupVersion.WithKind("MinIOBucket"), namespace, opts), &miniov1alpha1.MinIOBucketList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := k8stesting.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &miniov1alpha1.MinIOBucketList{ListMeta: obj.(*miniov1alpha1.MinIOBucketList).ListMeta}
	for _, item := range obj.(*miniov1alpha1.MinIOBucketList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

func (c *TestingMinioV1alpha1) WatchMinIOBucket(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(k8stesting.NewWatchAction(miniov1alpha1.SchemaGroupVersion.WithResource("miniobuckets"), namespace, opts))
}

func (c *TestingMinioV1alpha1) GetMinIOUser(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*miniov1alpha1.MinIOUser, error) {
	obj, err := c.Fake.Invokes(k8stesting.NewGetAction(miniov1alpha1.SchemaGroupVersion.WithResource("miniousers"), namespace, name), &miniov1alpha1.MinIOUser{})
	if obj == nil {
		return nil, err
	}
	return obj.(*miniov1alpha1.MinIOUser), err
}

func (c *TestingMinioV1alpha1) CreateMinIOUser(ctx context.Context, v *miniov1alpha1.MinIOUser, opts metav1.CreateOptions) (*miniov1alpha1.MinIOUser, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewCreateAction(miniov1alpha1.SchemaGroupVersion.WithResource("miniousers"), v.Namespace, v), &miniov1alpha1.MinIOUser{})

	if obj == nil {
		return nil, err
	}
	return obj.(*miniov1alpha1.MinIOUser), err
}

func (c *TestingMinioV1alpha1) UpdateMinIOUser(ctx context.Context, v *miniov1alpha1.MinIOUser, opts metav1.UpdateOptions) (*miniov1alpha1.MinIOUser, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewUpdateAction(miniov1alpha1.SchemaGroupVersion.WithResource("miniousers"), v.Namespace, v), &miniov1alpha1.MinIOUser{})

	if obj == nil {
		return nil, err
	}
	return obj.(*miniov1alpha1.MinIOUser), err
}

func (c *TestingMinioV1alpha1) UpdateStatusMinIOUser(ctx context.Context, v *miniov1alpha1.MinIOUser, opts metav1.UpdateOptions) (*miniov1alpha1.MinIOUser, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewUpdateSubresourceAction(miniov1alpha1.SchemaGroupVersion.WithResource("miniousers"), "status", v.Namespace, v), &miniov1alpha1.MinIOUser{})

	if obj == nil {
		return nil, err
	}
	return obj.(*miniov1alpha1.MinIOUser), err
}

func (c *TestingMinioV1alpha1) DeleteMinIOUser(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(k8stesting.NewDeleteAction(miniov1alpha1.SchemaGroupVersion.WithResource("miniousers"), namespace, name), &miniov1alpha1.MinIOUser{})

	return err
}

func (c *TestingMinioV1alpha1) ListMinIOUser(ctx context.Context, namespace string, opts metav1.ListOptions) (*miniov1alpha1.MinIOUserList, error) {
	obj, err := c.Fake.
		Invokes(k8stesting.NewListAction(miniov1alpha1.SchemaGroupVersion.WithResource("miniousers"), miniov1alpha1.SchemaGroupVersion.WithKind("MinIOUser"), namespace, opts), &miniov1alpha1.MinIOUserList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := k8stesting.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &miniov1alpha1.MinIOUserList{ListMeta: obj.(*miniov1alpha1.MinIOUserList).ListMeta}
	for _, item := range obj.(*miniov1alpha1.MinIOUserList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

func (c *TestingMinioV1alpha1) WatchMinIOUser(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(k8stesting.NewWatchAction(miniov1alpha1.SchemaGroupVersion.WithResource("miniousers"), namespace, opts))
}
