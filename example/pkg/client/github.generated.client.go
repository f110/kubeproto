package client

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"go.f110.dev/kubeproto/example/pkg/apis/githubv1alpha1"
	"go.f110.dev/kubeproto/example/pkg/apis/githubv1alpha2"
	"go.f110.dev/kubeproto/example/pkg/apis/miniov1alpha1"
)

var (
	Scheme         = runtime.NewScheme()
	ParameterCodec = runtime.NewParameterCodec(Scheme)
)

func init() {
	for _, v := range []func(*runtime.Scheme) error{
		githubv1alpha1.AddToScheme,
		githubv1alpha2.AddToScheme,
		miniov1alpha1.AddToScheme,
	} {
		if err := v(Scheme); err != nil {
			panic(err)
		}
	}
}

type Backend interface {
	Get(ctx context.Context, resourceName, kindName, namespace, name string, opts metav1.GetOptions, result runtime.Object) (runtime.Object, error)
	List(ctx context.Context, resourceName, kindName, namespace string, opts metav1.ListOptions, result runtime.Object) (runtime.Object, error)
	Create(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.CreateOptions, result runtime.Object) (runtime.Object, error)
	Update(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error)
	UpdateStatus(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error)
	Delete(ctx context.Context, gvr schema.GroupVersionResource, namespace, name string, opts metav1.DeleteOptions) error
	Watch(ctx context.Context, gvr schema.GroupVersionResource, namespace string, opts metav1.ListOptions) (watch.Interface, error)
}
type Set struct {
	GrafanaV1alpha1 *GrafanaV1alpha1
	GrafanaV1alpha2 *GrafanaV1alpha2
	MinioV1alpha1   *MinioV1alpha1
}

func NewSet(cfg *rest.Config) (*Set, error) {
	c, err := rest.RESTClientFor(cfg)
	if err != nil {
		return nil, err
	}
	b := &restBackend{client: c}
	s := &Set{
		GrafanaV1alpha1: NewGrafanaV1alpha1Client(b),
		GrafanaV1alpha2: NewGrafanaV1alpha2Client(b),
		MinioV1alpha1:   NewMinioV1alpha1Client(b),
	}

	return s, nil
}

type restBackend struct {
	client *rest.RESTClient
}

func (r *restBackend) Get(ctx context.Context, resourceName, kindName, namespace, name string, opts metav1.GetOptions, result runtime.Object) (runtime.Object, error) {
	return result, r.client.Get().
		Namespace(namespace).
		Resource(resourceName).
		Name(name).
		VersionedParams(&opts, ParameterCodec).
		Do(ctx).
		Into(result)
}

func (r *restBackend) List(ctx context.Context, resourceName, kindName, namespace string, opts metav1.ListOptions, result runtime.Object) (runtime.Object, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	return result, r.client.Get().
		Namespace(namespace).
		Resource(resourceName).
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
}

func (r *restBackend) Create(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.CreateOptions, result runtime.Object) (runtime.Object, error) {
	m := obj.(metav1.Object)
	if m == nil {
		return nil, errors.New("obj is not implement metav1.Object")
	}
	return result, r.client.Post().
		Namespace(m.GetNamespace()).
		Resource(resourceName).
		VersionedParams(&opts, ParameterCodec).
		Body(obj).
		Do(ctx).
		Into(result)
}

func (r *restBackend) Update(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	m := obj.(metav1.Object)
	if m == nil {
		return nil, errors.New("obj is not implement metav1.Object")
	}
	return result, r.client.Put().
		Namespace(m.GetNamespace()).
		Resource(resourceName).
		Name(m.GetName()).
		VersionedParams(&opts, ParameterCodec).
		Body(obj).
		Do(ctx).
		Into(result)
}

func (r *restBackend) UpdateStatus(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	m := obj.(metav1.Object)
	if m == nil {
		return nil, errors.New("obj is not implement metav1.Object")
	}
	return result, r.client.Put().
		Namespace(m.GetNamespace()).
		Resource(resourceName).
		Name(m.GetName()).
		SubResource("status").
		VersionedParams(&opts, ParameterCodec).
		Body(obj).
		Do(ctx).
		Into(result)
}

func (r *restBackend) Delete(ctx context.Context, gvr schema.GroupVersionResource, namespace, name string, opts metav1.DeleteOptions) error {
	return r.client.Delete().
		Namespace(namespace).
		Resource(gvr.Resource).
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (r *restBackend) Watch(ctx context.Context, gvr schema.GroupVersionResource, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return r.client.Get().
		Namespace(namespace).
		Resource(gvr.Resource).
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

type GrafanaV1alpha1 struct {
	backend Backend
}

func NewGrafanaV1alpha1Client(b Backend) *GrafanaV1alpha1 {
	return &GrafanaV1alpha1{backend: b}
}

func (c *GrafanaV1alpha1) GetGrafana(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*githubv1alpha1.Grafana, error) {
	result, err := c.backend.Get(ctx, "grafanas", "Grafana", namespace, name, opts, &githubv1alpha1.Grafana{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha1.Grafana), nil
}

func (c *GrafanaV1alpha1) CreateGrafana(ctx context.Context, v *githubv1alpha1.Grafana, opts metav1.CreateOptions) (*githubv1alpha1.Grafana, error) {
	result, err := c.backend.Create(ctx, "grafanas", "Grafana", v, opts, &githubv1alpha1.Grafana{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha1.Grafana), nil
}

func (c *GrafanaV1alpha1) UpdateGrafana(ctx context.Context, v *githubv1alpha1.Grafana, opts metav1.UpdateOptions) (*githubv1alpha1.Grafana, error) {
	result, err := c.backend.Update(ctx, "grafanas", "Grafana", v, opts, &githubv1alpha1.Grafana{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha1.Grafana), nil
}

func (c *GrafanaV1alpha1) DeleteGrafana(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: "grafana.f110.dev", Version: "v1alpha1", Resource: "grafanas"}, namespace, name, opts)
}

func (c *GrafanaV1alpha1) ListGrafana(ctx context.Context, namespace string, opts metav1.ListOptions) (*githubv1alpha1.GrafanaList, error) {
	result, err := c.backend.List(ctx, "grafanas", "Grafana", namespace, opts, &githubv1alpha1.GrafanaList{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha1.GrafanaList), nil
}

func (c *GrafanaV1alpha1) WatchGrafana(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: "grafana.f110.dev", Version: "v1alpha1", Resource: "grafanas"}, namespace, opts)
}

func (c *GrafanaV1alpha1) GetGrafanaUser(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*githubv1alpha1.GrafanaUser, error) {
	result, err := c.backend.Get(ctx, "grafanausers", "GrafanaUser", namespace, name, opts, &githubv1alpha1.GrafanaUser{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha1.GrafanaUser), nil
}

func (c *GrafanaV1alpha1) CreateGrafanaUser(ctx context.Context, v *githubv1alpha1.GrafanaUser, opts metav1.CreateOptions) (*githubv1alpha1.GrafanaUser, error) {
	result, err := c.backend.Create(ctx, "grafanausers", "GrafanaUser", v, opts, &githubv1alpha1.GrafanaUser{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha1.GrafanaUser), nil
}

func (c *GrafanaV1alpha1) UpdateGrafanaUser(ctx context.Context, v *githubv1alpha1.GrafanaUser, opts metav1.UpdateOptions) (*githubv1alpha1.GrafanaUser, error) {
	result, err := c.backend.Update(ctx, "grafanausers", "GrafanaUser", v, opts, &githubv1alpha1.GrafanaUser{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha1.GrafanaUser), nil
}

func (c *GrafanaV1alpha1) DeleteGrafanaUser(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: "grafana.f110.dev", Version: "v1alpha1", Resource: "grafanausers"}, namespace, name, opts)
}

func (c *GrafanaV1alpha1) ListGrafanaUser(ctx context.Context, namespace string, opts metav1.ListOptions) (*githubv1alpha1.GrafanaUserList, error) {
	result, err := c.backend.List(ctx, "grafanausers", "GrafanaUser", namespace, opts, &githubv1alpha1.GrafanaUserList{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha1.GrafanaUserList), nil
}

func (c *GrafanaV1alpha1) WatchGrafanaUser(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: "grafana.f110.dev", Version: "v1alpha1", Resource: "grafanausers"}, namespace, opts)
}

type GrafanaV1alpha2 struct {
	backend Backend
}

func NewGrafanaV1alpha2Client(b Backend) *GrafanaV1alpha2 {
	return &GrafanaV1alpha2{backend: b}
}

func (c *GrafanaV1alpha2) GetGrafana(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*githubv1alpha2.Grafana, error) {
	result, err := c.backend.Get(ctx, "grafanas", "Grafana", namespace, name, opts, &githubv1alpha2.Grafana{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha2.Grafana), nil
}

func (c *GrafanaV1alpha2) CreateGrafana(ctx context.Context, v *githubv1alpha2.Grafana, opts metav1.CreateOptions) (*githubv1alpha2.Grafana, error) {
	result, err := c.backend.Create(ctx, "grafanas", "Grafana", v, opts, &githubv1alpha2.Grafana{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha2.Grafana), nil
}

func (c *GrafanaV1alpha2) UpdateGrafana(ctx context.Context, v *githubv1alpha2.Grafana, opts metav1.UpdateOptions) (*githubv1alpha2.Grafana, error) {
	result, err := c.backend.Update(ctx, "grafanas", "Grafana", v, opts, &githubv1alpha2.Grafana{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha2.Grafana), nil
}

func (c *GrafanaV1alpha2) DeleteGrafana(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: "grafana.f110.dev", Version: "v1alpha2", Resource: "grafanas"}, namespace, name, opts)
}

func (c *GrafanaV1alpha2) ListGrafana(ctx context.Context, namespace string, opts metav1.ListOptions) (*githubv1alpha2.GrafanaList, error) {
	result, err := c.backend.List(ctx, "grafanas", "Grafana", namespace, opts, &githubv1alpha2.GrafanaList{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha2.GrafanaList), nil
}

func (c *GrafanaV1alpha2) WatchGrafana(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: "grafana.f110.dev", Version: "v1alpha2", Resource: "grafanas"}, namespace, opts)
}

func (c *GrafanaV1alpha2) GetGrafanaUser(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*githubv1alpha2.GrafanaUser, error) {
	result, err := c.backend.Get(ctx, "grafanausers", "GrafanaUser", namespace, name, opts, &githubv1alpha2.GrafanaUser{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha2.GrafanaUser), nil
}

func (c *GrafanaV1alpha2) CreateGrafanaUser(ctx context.Context, v *githubv1alpha2.GrafanaUser, opts metav1.CreateOptions) (*githubv1alpha2.GrafanaUser, error) {
	result, err := c.backend.Create(ctx, "grafanausers", "GrafanaUser", v, opts, &githubv1alpha2.GrafanaUser{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha2.GrafanaUser), nil
}

func (c *GrafanaV1alpha2) UpdateGrafanaUser(ctx context.Context, v *githubv1alpha2.GrafanaUser, opts metav1.UpdateOptions) (*githubv1alpha2.GrafanaUser, error) {
	result, err := c.backend.Update(ctx, "grafanausers", "GrafanaUser", v, opts, &githubv1alpha2.GrafanaUser{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha2.GrafanaUser), nil
}

func (c *GrafanaV1alpha2) DeleteGrafanaUser(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: "grafana.f110.dev", Version: "v1alpha2", Resource: "grafanausers"}, namespace, name, opts)
}

func (c *GrafanaV1alpha2) ListGrafanaUser(ctx context.Context, namespace string, opts metav1.ListOptions) (*githubv1alpha2.GrafanaUserList, error) {
	result, err := c.backend.List(ctx, "grafanausers", "GrafanaUser", namespace, opts, &githubv1alpha2.GrafanaUserList{})
	if err != nil {
		return nil, err
	}
	return result.(*githubv1alpha2.GrafanaUserList), nil
}

func (c *GrafanaV1alpha2) WatchGrafanaUser(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: "grafana.f110.dev", Version: "v1alpha2", Resource: "grafanausers"}, namespace, opts)
}

type MinioV1alpha1 struct {
	backend Backend
}

func NewMinioV1alpha1Client(b Backend) *MinioV1alpha1 {
	return &MinioV1alpha1{backend: b}
}

func (c *MinioV1alpha1) GetMinIOBucket(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*miniov1alpha1.MinIOBucket, error) {
	result, err := c.backend.Get(ctx, "miniobuckets", "MinIOBucket", namespace, name, opts, &miniov1alpha1.MinIOBucket{})
	if err != nil {
		return nil, err
	}
	return result.(*miniov1alpha1.MinIOBucket), nil
}

func (c *MinioV1alpha1) CreateMinIOBucket(ctx context.Context, v *miniov1alpha1.MinIOBucket, opts metav1.CreateOptions) (*miniov1alpha1.MinIOBucket, error) {
	result, err := c.backend.Create(ctx, "miniobuckets", "MinIOBucket", v, opts, &miniov1alpha1.MinIOBucket{})
	if err != nil {
		return nil, err
	}
	return result.(*miniov1alpha1.MinIOBucket), nil
}

func (c *MinioV1alpha1) UpdateMinIOBucket(ctx context.Context, v *miniov1alpha1.MinIOBucket, opts metav1.UpdateOptions) (*miniov1alpha1.MinIOBucket, error) {
	result, err := c.backend.Update(ctx, "miniobuckets", "MinIOBucket", v, opts, &miniov1alpha1.MinIOBucket{})
	if err != nil {
		return nil, err
	}
	return result.(*miniov1alpha1.MinIOBucket), nil
}

func (c *MinioV1alpha1) UpdateStatusMinIOBucket(ctx context.Context, v *miniov1alpha1.MinIOBucket, opts metav1.UpdateOptions) (*miniov1alpha1.MinIOBucket, error) {
	result, err := c.backend.UpdateStatus(ctx, "miniobuckets", "MinIOBucket", v, opts, &miniov1alpha1.MinIOBucket{})
	if err != nil {
		return nil, err
	}
	return result.(*miniov1alpha1.MinIOBucket), nil
}

func (c *MinioV1alpha1) DeleteMinIOBucket(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: "minio.f110.dev", Version: "v1alpha1", Resource: "miniobuckets"}, namespace, name, opts)
}

func (c *MinioV1alpha1) ListMinIOBucket(ctx context.Context, namespace string, opts metav1.ListOptions) (*miniov1alpha1.MinIOBucketList, error) {
	result, err := c.backend.List(ctx, "miniobuckets", "MinIOBucket", namespace, opts, &miniov1alpha1.MinIOBucketList{})
	if err != nil {
		return nil, err
	}
	return result.(*miniov1alpha1.MinIOBucketList), nil
}

func (c *MinioV1alpha1) WatchMinIOBucket(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: "minio.f110.dev", Version: "v1alpha1", Resource: "miniobuckets"}, namespace, opts)
}

func (c *MinioV1alpha1) GetMinIOUser(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*miniov1alpha1.MinIOUser, error) {
	result, err := c.backend.Get(ctx, "miniousers", "MinIOUser", namespace, name, opts, &miniov1alpha1.MinIOUser{})
	if err != nil {
		return nil, err
	}
	return result.(*miniov1alpha1.MinIOUser), nil
}

func (c *MinioV1alpha1) CreateMinIOUser(ctx context.Context, v *miniov1alpha1.MinIOUser, opts metav1.CreateOptions) (*miniov1alpha1.MinIOUser, error) {
	result, err := c.backend.Create(ctx, "miniousers", "MinIOUser", v, opts, &miniov1alpha1.MinIOUser{})
	if err != nil {
		return nil, err
	}
	return result.(*miniov1alpha1.MinIOUser), nil
}

func (c *MinioV1alpha1) UpdateMinIOUser(ctx context.Context, v *miniov1alpha1.MinIOUser, opts metav1.UpdateOptions) (*miniov1alpha1.MinIOUser, error) {
	result, err := c.backend.Update(ctx, "miniousers", "MinIOUser", v, opts, &miniov1alpha1.MinIOUser{})
	if err != nil {
		return nil, err
	}
	return result.(*miniov1alpha1.MinIOUser), nil
}

func (c *MinioV1alpha1) UpdateStatusMinIOUser(ctx context.Context, v *miniov1alpha1.MinIOUser, opts metav1.UpdateOptions) (*miniov1alpha1.MinIOUser, error) {
	result, err := c.backend.UpdateStatus(ctx, "miniousers", "MinIOUser", v, opts, &miniov1alpha1.MinIOUser{})
	if err != nil {
		return nil, err
	}
	return result.(*miniov1alpha1.MinIOUser), nil
}

func (c *MinioV1alpha1) DeleteMinIOUser(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: "minio.f110.dev", Version: "v1alpha1", Resource: "miniousers"}, namespace, name, opts)
}

func (c *MinioV1alpha1) ListMinIOUser(ctx context.Context, namespace string, opts metav1.ListOptions) (*miniov1alpha1.MinIOUserList, error) {
	result, err := c.backend.List(ctx, "miniousers", "MinIOUser", namespace, opts, &miniov1alpha1.MinIOUserList{})
	if err != nil {
		return nil, err
	}
	return result.(*miniov1alpha1.MinIOUserList), nil
}

func (c *MinioV1alpha1) WatchMinIOUser(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: "minio.f110.dev", Version: "v1alpha1", Resource: "miniousers"}, namespace, opts)
}

var Factory = NewInformerFactory()

type InformerFactory struct {
	mu        sync.Mutex
	informers map[reflect.Type]cache.SharedIndexInformer
	once      sync.Once
	ctx       context.Context
}

func NewInformerFactory() *InformerFactory {
	return &InformerFactory{informers: make(map[reflect.Type]cache.SharedIndexInformer)}
}

func (f *InformerFactory) InformerFor(obj runtime.Object, newFunc func() cache.SharedIndexInformer) cache.SharedIndexInformer {
	f.mu.Lock()
	defer f.mu.Unlock()

	typ := reflect.TypeOf(obj)
	if v, ok := f.informers[typ]; ok {
		return v
	}
	informer := newFunc()
	f.informers[typ] = informer
	if f.ctx != nil {
		go informer.Run(f.ctx.Done())
	}
	return informer
}

func (f *InformerFactory) Run(ctx context.Context) {
	f.mu.Lock()
	f.once.Do(func() {
		for _, v := range f.informers {
			go v.Run(ctx.Done())
		}
		f.ctx = ctx
	})
	f.mu.Unlock()
}

type GrafanaV1alpha1Informer struct {
	factory      *InformerFactory
	client       *GrafanaV1alpha1
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewGrafanaV1alpha1Informer(f *InformerFactory, client *GrafanaV1alpha1, namespace string, resyncPeriod time.Duration) *GrafanaV1alpha1Informer {
	return &GrafanaV1alpha1Informer{
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *GrafanaV1alpha1Informer) GrafanaInformer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&githubv1alpha1.Grafana{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListGrafana(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchGrafana(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&githubv1alpha1.Grafana{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *GrafanaV1alpha1Informer) GrafanaLister() *GrafanaV1alpha1GrafanaLister {
	return NewGrafanaV1alpha1GrafanaLister(f.GrafanaInformer().GetIndexer())
}

func (f *GrafanaV1alpha1Informer) GrafanaUserInformer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&githubv1alpha1.GrafanaUser{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListGrafanaUser(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchGrafanaUser(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&githubv1alpha1.GrafanaUser{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *GrafanaV1alpha1Informer) GrafanaUserLister() *GrafanaV1alpha1GrafanaUserLister {
	return NewGrafanaV1alpha1GrafanaUserLister(f.GrafanaUserInformer().GetIndexer())
}

type GrafanaV1alpha2Informer struct {
	factory      *InformerFactory
	client       *GrafanaV1alpha2
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewGrafanaV1alpha2Informer(f *InformerFactory, client *GrafanaV1alpha2, namespace string, resyncPeriod time.Duration) *GrafanaV1alpha2Informer {
	return &GrafanaV1alpha2Informer{
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *GrafanaV1alpha2Informer) GrafanaInformer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&githubv1alpha2.Grafana{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListGrafana(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchGrafana(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&githubv1alpha2.Grafana{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *GrafanaV1alpha2Informer) GrafanaLister() *GrafanaV1alpha2GrafanaLister {
	return NewGrafanaV1alpha2GrafanaLister(f.GrafanaInformer().GetIndexer())
}

func (f *GrafanaV1alpha2Informer) GrafanaUserInformer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&githubv1alpha2.GrafanaUser{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListGrafanaUser(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchGrafanaUser(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&githubv1alpha2.GrafanaUser{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *GrafanaV1alpha2Informer) GrafanaUserLister() *GrafanaV1alpha2GrafanaUserLister {
	return NewGrafanaV1alpha2GrafanaUserLister(f.GrafanaUserInformer().GetIndexer())
}

type MinioV1alpha1Informer struct {
	factory      *InformerFactory
	client       *MinioV1alpha1
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewMinioV1alpha1Informer(f *InformerFactory, client *MinioV1alpha1, namespace string, resyncPeriod time.Duration) *MinioV1alpha1Informer {
	return &MinioV1alpha1Informer{
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *MinioV1alpha1Informer) MinIOBucketInformer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&miniov1alpha1.MinIOBucket{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListMinIOBucket(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchMinIOBucket(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&miniov1alpha1.MinIOBucket{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *MinioV1alpha1Informer) MinIOBucketLister() *MinioV1alpha1MinIOBucketLister {
	return NewMinioV1alpha1MinIOBucketLister(f.MinIOBucketInformer().GetIndexer())
}

func (f *MinioV1alpha1Informer) MinIOUserInformer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&miniov1alpha1.MinIOUser{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListMinIOUser(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchMinIOUser(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&miniov1alpha1.MinIOUser{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *MinioV1alpha1Informer) MinIOUserLister() *MinioV1alpha1MinIOUserLister {
	return NewMinioV1alpha1MinIOUserLister(f.MinIOUserInformer().GetIndexer())
}

type GrafanaV1alpha1GrafanaLister struct {
	indexer cache.Indexer
}

func NewGrafanaV1alpha1GrafanaLister(indexer cache.Indexer) *GrafanaV1alpha1GrafanaLister {
	return &GrafanaV1alpha1GrafanaLister{indexer: indexer}
}

func (x *GrafanaV1alpha1GrafanaLister) List(namespace string, selector labels.Selector) ([]*githubv1alpha1.Grafana, error) {
	var ret []*githubv1alpha1.Grafana
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*githubv1alpha1.Grafana).DeepCopy())
	})
	return ret, err
}

func (x *GrafanaV1alpha1GrafanaLister) Get(namespace, name string) (*githubv1alpha1.Grafana, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(githubv1alpha1.SchemaGroupVersion.WithResource("grafana").GroupResource(), name)
	}
	return obj.(*githubv1alpha1.Grafana).DeepCopy(), nil
}

type GrafanaV1alpha1GrafanaUserLister struct {
	indexer cache.Indexer
}

func NewGrafanaV1alpha1GrafanaUserLister(indexer cache.Indexer) *GrafanaV1alpha1GrafanaUserLister {
	return &GrafanaV1alpha1GrafanaUserLister{indexer: indexer}
}

func (x *GrafanaV1alpha1GrafanaUserLister) List(namespace string, selector labels.Selector) ([]*githubv1alpha1.GrafanaUser, error) {
	var ret []*githubv1alpha1.GrafanaUser
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*githubv1alpha1.GrafanaUser).DeepCopy())
	})
	return ret, err
}

func (x *GrafanaV1alpha1GrafanaUserLister) Get(namespace, name string) (*githubv1alpha1.GrafanaUser, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(githubv1alpha1.SchemaGroupVersion.WithResource("grafanauser").GroupResource(), name)
	}
	return obj.(*githubv1alpha1.GrafanaUser).DeepCopy(), nil
}

type GrafanaV1alpha2GrafanaLister struct {
	indexer cache.Indexer
}

func NewGrafanaV1alpha2GrafanaLister(indexer cache.Indexer) *GrafanaV1alpha2GrafanaLister {
	return &GrafanaV1alpha2GrafanaLister{indexer: indexer}
}

func (x *GrafanaV1alpha2GrafanaLister) List(namespace string, selector labels.Selector) ([]*githubv1alpha2.Grafana, error) {
	var ret []*githubv1alpha2.Grafana
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*githubv1alpha2.Grafana).DeepCopy())
	})
	return ret, err
}

func (x *GrafanaV1alpha2GrafanaLister) Get(namespace, name string) (*githubv1alpha2.Grafana, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(githubv1alpha2.SchemaGroupVersion.WithResource("grafana").GroupResource(), name)
	}
	return obj.(*githubv1alpha2.Grafana).DeepCopy(), nil
}

type GrafanaV1alpha2GrafanaUserLister struct {
	indexer cache.Indexer
}

func NewGrafanaV1alpha2GrafanaUserLister(indexer cache.Indexer) *GrafanaV1alpha2GrafanaUserLister {
	return &GrafanaV1alpha2GrafanaUserLister{indexer: indexer}
}

func (x *GrafanaV1alpha2GrafanaUserLister) List(namespace string, selector labels.Selector) ([]*githubv1alpha2.GrafanaUser, error) {
	var ret []*githubv1alpha2.GrafanaUser
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*githubv1alpha2.GrafanaUser).DeepCopy())
	})
	return ret, err
}

func (x *GrafanaV1alpha2GrafanaUserLister) Get(namespace, name string) (*githubv1alpha2.GrafanaUser, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(githubv1alpha2.SchemaGroupVersion.WithResource("grafanauser").GroupResource(), name)
	}
	return obj.(*githubv1alpha2.GrafanaUser).DeepCopy(), nil
}

type MinioV1alpha1MinIOBucketLister struct {
	indexer cache.Indexer
}

func NewMinioV1alpha1MinIOBucketLister(indexer cache.Indexer) *MinioV1alpha1MinIOBucketLister {
	return &MinioV1alpha1MinIOBucketLister{indexer: indexer}
}

func (x *MinioV1alpha1MinIOBucketLister) List(namespace string, selector labels.Selector) ([]*miniov1alpha1.MinIOBucket, error) {
	var ret []*miniov1alpha1.MinIOBucket
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*miniov1alpha1.MinIOBucket).DeepCopy())
	})
	return ret, err
}

func (x *MinioV1alpha1MinIOBucketLister) Get(namespace, name string) (*miniov1alpha1.MinIOBucket, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(miniov1alpha1.SchemaGroupVersion.WithResource("miniobucket").GroupResource(), name)
	}
	return obj.(*miniov1alpha1.MinIOBucket).DeepCopy(), nil
}

type MinioV1alpha1MinIOUserLister struct {
	indexer cache.Indexer
}

func NewMinioV1alpha1MinIOUserLister(indexer cache.Indexer) *MinioV1alpha1MinIOUserLister {
	return &MinioV1alpha1MinIOUserLister{indexer: indexer}
}

func (x *MinioV1alpha1MinIOUserLister) List(namespace string, selector labels.Selector) ([]*miniov1alpha1.MinIOUser, error) {
	var ret []*miniov1alpha1.MinIOUser
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*miniov1alpha1.MinIOUser).DeepCopy())
	})
	return ret, err
}

func (x *MinioV1alpha1MinIOUserLister) Get(namespace, name string) (*miniov1alpha1.MinIOUser, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(miniov1alpha1.SchemaGroupVersion.WithResource("miniouser").GroupResource(), name)
	}
	return obj.(*miniov1alpha1.MinIOUser).DeepCopy(), nil
}
