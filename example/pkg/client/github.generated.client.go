package client

import (
	"context"
	"reflect"
	"sync"
	"time"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
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

type GrafanaV1alpha1 struct {
	client *rest.RESTClient
}

func NewGrafanaV1alpha1Client(c *rest.Config) (*GrafanaV1alpha1, error) {
	client, err := rest.RESTClientFor(c)
	if err != nil {
		return nil, err
	}
	return &GrafanaV1alpha1{
		client: client,
	}, nil
}

func (c *GrafanaV1alpha1) GetGrafana(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*githubv1alpha1.Grafana, error) {
	result := &githubv1alpha1.Grafana{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("grafanas").
		Name(name).
		VersionedParams(&opts, ParameterCodec).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha1) CreateGrafana(ctx context.Context, v *githubv1alpha1.Grafana, opts metav1.CreateOptions) (*githubv1alpha1.Grafana, error) {
	result := &githubv1alpha1.Grafana{}
	err := c.client.Post().
		Namespace(v.Namespace).
		Resource("grafanas").
		VersionedParams(&opts, ParameterCodec).
		Body(v).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha1) UpdateGrafana(ctx context.Context, v *githubv1alpha1.Grafana, opts metav1.UpdateOptions) (*githubv1alpha1.Grafana, error) {
	result := &githubv1alpha1.Grafana{}
	err := c.client.Put().
		Namespace(v.Namespace).
		Resource("grafanas").
		Name(v.Name).
		VersionedParams(&opts, ParameterCodec).
		Body(v).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha1) DeleteGrafana(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(namespace).
		Resource("grafanas").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *GrafanaV1alpha1) ListGrafana(ctx context.Context, namespace string, opts metav1.ListOptions) (*githubv1alpha1.GrafanaList, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result := &githubv1alpha1.GrafanaList{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("grafanas").
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha1) WatchGrafana(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(namespace).
		Resource("grafanas").
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

func (c *GrafanaV1alpha1) GetGrafanaUser(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*githubv1alpha1.GrafanaUser, error) {
	result := &githubv1alpha1.GrafanaUser{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("grafanausers").
		Name(name).
		VersionedParams(&opts, ParameterCodec).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha1) CreateGrafanaUser(ctx context.Context, v *githubv1alpha1.GrafanaUser, opts metav1.CreateOptions) (*githubv1alpha1.GrafanaUser, error) {
	result := &githubv1alpha1.GrafanaUser{}
	err := c.client.Post().
		Namespace(v.Namespace).
		Resource("grafanausers").
		VersionedParams(&opts, ParameterCodec).
		Body(v).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha1) UpdateGrafanaUser(ctx context.Context, v *githubv1alpha1.GrafanaUser, opts metav1.UpdateOptions) (*githubv1alpha1.GrafanaUser, error) {
	result := &githubv1alpha1.GrafanaUser{}
	err := c.client.Put().
		Namespace(v.Namespace).
		Resource("grafanausers").
		Name(v.Name).
		VersionedParams(&opts, ParameterCodec).
		Body(v).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha1) DeleteGrafanaUser(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(namespace).
		Resource("grafanausers").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *GrafanaV1alpha1) ListGrafanaUser(ctx context.Context, namespace string, opts metav1.ListOptions) (*githubv1alpha1.GrafanaUserList, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result := &githubv1alpha1.GrafanaUserList{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("grafanausers").
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha1) WatchGrafanaUser(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(namespace).
		Resource("grafanausers").
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

type GrafanaV1alpha2 struct {
	client *rest.RESTClient
}

func NewGrafanaV1alpha2Client(c *rest.Config) (*GrafanaV1alpha2, error) {
	client, err := rest.RESTClientFor(c)
	if err != nil {
		return nil, err
	}
	return &GrafanaV1alpha2{
		client: client,
	}, nil
}

func (c *GrafanaV1alpha2) GetGrafana(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*githubv1alpha2.Grafana, error) {
	result := &githubv1alpha2.Grafana{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("grafanas").
		Name(name).
		VersionedParams(&opts, ParameterCodec).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha2) CreateGrafana(ctx context.Context, v *githubv1alpha2.Grafana, opts metav1.CreateOptions) (*githubv1alpha2.Grafana, error) {
	result := &githubv1alpha2.Grafana{}
	err := c.client.Post().
		Namespace(v.Namespace).
		Resource("grafanas").
		VersionedParams(&opts, ParameterCodec).
		Body(v).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha2) UpdateGrafana(ctx context.Context, v *githubv1alpha2.Grafana, opts metav1.UpdateOptions) (*githubv1alpha2.Grafana, error) {
	result := &githubv1alpha2.Grafana{}
	err := c.client.Put().
		Namespace(v.Namespace).
		Resource("grafanas").
		Name(v.Name).
		VersionedParams(&opts, ParameterCodec).
		Body(v).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha2) DeleteGrafana(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(namespace).
		Resource("grafanas").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *GrafanaV1alpha2) ListGrafana(ctx context.Context, namespace string, opts metav1.ListOptions) (*githubv1alpha2.GrafanaList, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result := &githubv1alpha2.GrafanaList{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("grafanas").
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha2) WatchGrafana(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(namespace).
		Resource("grafanas").
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

func (c *GrafanaV1alpha2) GetGrafanaUser(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*githubv1alpha2.GrafanaUser, error) {
	result := &githubv1alpha2.GrafanaUser{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("grafanausers").
		Name(name).
		VersionedParams(&opts, ParameterCodec).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha2) CreateGrafanaUser(ctx context.Context, v *githubv1alpha2.GrafanaUser, opts metav1.CreateOptions) (*githubv1alpha2.GrafanaUser, error) {
	result := &githubv1alpha2.GrafanaUser{}
	err := c.client.Post().
		Namespace(v.Namespace).
		Resource("grafanausers").
		VersionedParams(&opts, ParameterCodec).
		Body(v).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha2) UpdateGrafanaUser(ctx context.Context, v *githubv1alpha2.GrafanaUser, opts metav1.UpdateOptions) (*githubv1alpha2.GrafanaUser, error) {
	result := &githubv1alpha2.GrafanaUser{}
	err := c.client.Put().
		Namespace(v.Namespace).
		Resource("grafanausers").
		Name(v.Name).
		VersionedParams(&opts, ParameterCodec).
		Body(v).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha2) DeleteGrafanaUser(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(namespace).
		Resource("grafanausers").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *GrafanaV1alpha2) ListGrafanaUser(ctx context.Context, namespace string, opts metav1.ListOptions) (*githubv1alpha2.GrafanaUserList, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result := &githubv1alpha2.GrafanaUserList{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("grafanausers").
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *GrafanaV1alpha2) WatchGrafanaUser(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(namespace).
		Resource("grafanausers").
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

type MinioV1alpha1 struct {
	client *rest.RESTClient
}

func NewMinioV1alpha1Client(c *rest.Config) (*MinioV1alpha1, error) {
	client, err := rest.RESTClientFor(c)
	if err != nil {
		return nil, err
	}
	return &MinioV1alpha1{
		client: client,
	}, nil
}

func (c *MinioV1alpha1) GetMinIOBucket(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*miniov1alpha1.MinIOBucket, error) {
	result := &miniov1alpha1.MinIOBucket{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("miniobuckets").
		Name(name).
		VersionedParams(&opts, ParameterCodec).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MinioV1alpha1) CreateMinIOBucket(ctx context.Context, v *miniov1alpha1.MinIOBucket, opts metav1.CreateOptions) (*miniov1alpha1.MinIOBucket, error) {
	result := &miniov1alpha1.MinIOBucket{}
	err := c.client.Post().
		Namespace(v.Namespace).
		Resource("miniobuckets").
		VersionedParams(&opts, ParameterCodec).
		Body(v).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MinioV1alpha1) UpdateMinIOBucket(ctx context.Context, v *miniov1alpha1.MinIOBucket, opts metav1.UpdateOptions) (*miniov1alpha1.MinIOBucket, error) {
	result := &miniov1alpha1.MinIOBucket{}
	err := c.client.Put().
		Namespace(v.Namespace).
		Resource("miniobuckets").
		Name(v.Name).
		VersionedParams(&opts, ParameterCodec).
		Body(v).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MinioV1alpha1) DeleteMinIOBucket(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(namespace).
		Resource("miniobuckets").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *MinioV1alpha1) ListMinIOBucket(ctx context.Context, namespace string, opts metav1.ListOptions) (*miniov1alpha1.MinIOBucketList, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result := &miniov1alpha1.MinIOBucketList{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("miniobuckets").
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MinioV1alpha1) WatchMinIOBucket(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(namespace).
		Resource("miniobuckets").
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

func (c *MinioV1alpha1) GetMinIOUser(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*miniov1alpha1.MinIOUser, error) {
	result := &miniov1alpha1.MinIOUser{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("miniousers").
		Name(name).
		VersionedParams(&opts, ParameterCodec).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MinioV1alpha1) CreateMinIOUser(ctx context.Context, v *miniov1alpha1.MinIOUser, opts metav1.CreateOptions) (*miniov1alpha1.MinIOUser, error) {
	result := &miniov1alpha1.MinIOUser{}
	err := c.client.Post().
		Namespace(v.Namespace).
		Resource("miniousers").
		VersionedParams(&opts, ParameterCodec).
		Body(v).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MinioV1alpha1) UpdateMinIOUser(ctx context.Context, v *miniov1alpha1.MinIOUser, opts metav1.UpdateOptions) (*miniov1alpha1.MinIOUser, error) {
	result := &miniov1alpha1.MinIOUser{}
	err := c.client.Put().
		Namespace(v.Namespace).
		Resource("miniousers").
		Name(v.Name).
		VersionedParams(&opts, ParameterCodec).
		Body(v).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MinioV1alpha1) DeleteMinIOUser(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(namespace).
		Resource("miniousers").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *MinioV1alpha1) ListMinIOUser(ctx context.Context, namespace string, opts metav1.ListOptions) (*miniov1alpha1.MinIOUserList, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result := &miniov1alpha1.MinIOUserList{}
	err := c.client.Get().
		Namespace(namespace).
		Resource("miniousers").
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MinioV1alpha1) WatchMinIOUser(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(namespace).
		Resource("miniousers").
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
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

func NewGrafanaV1alpha1Informer(f *InformerFactory, client *GrafanaV1alpha1, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) *GrafanaV1alpha1Informer {
	return &GrafanaV1alpha1Informer{factory: f, client: client, namespace: namespace, resyncPeriod: resyncPeriod, indexers: indexers}
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

func NewGrafanaV1alpha2Informer(f *InformerFactory, client *GrafanaV1alpha2, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) *GrafanaV1alpha2Informer {
	return &GrafanaV1alpha2Informer{factory: f, client: client, namespace: namespace, resyncPeriod: resyncPeriod, indexers: indexers}
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

func NewMinioV1alpha1Informer(f *InformerFactory, client *MinioV1alpha1, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) *MinioV1alpha1Informer {
	return &MinioV1alpha1Informer{factory: f, client: client, namespace: namespace, resyncPeriod: resyncPeriod, indexers: indexers}
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
		ret = append(ret, m.(*githubv1alpha1.Grafana))
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
	return obj.(*githubv1alpha1.Grafana), nil
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
		ret = append(ret, m.(*githubv1alpha1.GrafanaUser))
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
	return obj.(*githubv1alpha1.GrafanaUser), nil
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
		ret = append(ret, m.(*githubv1alpha2.Grafana))
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
	return obj.(*githubv1alpha2.Grafana), nil
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
		ret = append(ret, m.(*githubv1alpha2.GrafanaUser))
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
	return obj.(*githubv1alpha2.GrafanaUser), nil
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
		ret = append(ret, m.(*miniov1alpha1.MinIOBucket))
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
	return obj.(*miniov1alpha1.MinIOBucket), nil
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
		ret = append(ret, m.(*miniov1alpha1.MinIOUser))
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
	return obj.(*miniov1alpha1.MinIOUser), nil
}
