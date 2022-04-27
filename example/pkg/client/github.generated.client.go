package client

import (
	"context"
	"reflect"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"go.f110.dev/kubeproto/example/pkg/apis/githubv1alpha1"
	"go.f110.dev/kubeproto/example/pkg/apis/githubv1alpha2"
)

var (
	Scheme         = runtime.NewScheme()
	ParameterCodec = runtime.NewParameterCodec(Scheme)
)

func init() {
	for _, v := range []func(*runtime.Scheme) error{
		githubv1alpha1.AddToScheme,
		githubv1alpha2.AddToScheme,
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
	factory *InformerFactory
	client  *GrafanaV1alpha1
}

func NewGrafanaV1alpha1Informer(f *InformerFactory, client *GrafanaV1alpha1) *GrafanaV1alpha1Informer {
	return &GrafanaV1alpha1Informer{factory: f, client: client}
}

func (f *GrafanaV1alpha1Informer) GrafanaInformer(namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return f.factory.InformerFor(&githubv1alpha1.Grafana{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListGrafana(context.TODO(), namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchGrafana(context.TODO(), namespace, metav1.ListOptions{})
				},
			},
			&githubv1alpha1.Grafana{},
			resyncPeriod,
			indexers,
		)
	},
	)
}
func (f *GrafanaV1alpha1Informer) GrafanaUserInformer(namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return f.factory.InformerFor(&githubv1alpha1.GrafanaUser{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListGrafanaUser(context.TODO(), namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchGrafanaUser(context.TODO(), namespace, metav1.ListOptions{})
				},
			},
			&githubv1alpha1.GrafanaUser{},
			resyncPeriod,
			indexers,
		)
	},
	)
}

type GrafanaV1alpha2Informer struct {
	factory *InformerFactory
	client  *GrafanaV1alpha2
}

func NewGrafanaV1alpha2Informer(f *InformerFactory, client *GrafanaV1alpha2) *GrafanaV1alpha2Informer {
	return &GrafanaV1alpha2Informer{factory: f, client: client}
}

func (f *GrafanaV1alpha2Informer) GrafanaInformer(namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return f.factory.InformerFor(&githubv1alpha2.Grafana{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListGrafana(context.TODO(), namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchGrafana(context.TODO(), namespace, metav1.ListOptions{})
				},
			},
			&githubv1alpha2.Grafana{},
			resyncPeriod,
			indexers,
		)
	},
	)
}
func (f *GrafanaV1alpha2Informer) GrafanaUserInformer(namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return f.factory.InformerFor(&githubv1alpha2.GrafanaUser{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListGrafanaUser(context.TODO(), namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchGrafanaUser(context.TODO(), namespace, metav1.ListOptions{})
				},
			},
			&githubv1alpha2.GrafanaUser{},
			resyncPeriod,
			indexers,
		)
	},
	)
}
