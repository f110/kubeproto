package client

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"go.f110.dev/kubeproto/go/apis/metav1"

	"go.f110.dev/kubeproto/example/pkg/apis/blogv1alpha1"
	"go.f110.dev/kubeproto/example/pkg/apis/blogv1alpha2"
)

var (
	Scheme         = runtime.NewScheme()
	ParameterCodec = runtime.NewParameterCodec(Scheme)
	Codecs         = serializer.NewCodecFactory(Scheme)
	AddToScheme    = localSchemeBuilder.AddToScheme
)

var localSchemeBuilder = runtime.SchemeBuilder{
	blogv1alpha1.AddToScheme,
	blogv1alpha2.AddToScheme,
}

func init() {
	for _, v := range []func(*runtime.Scheme) error{
		blogv1alpha1.AddToScheme,
		blogv1alpha2.AddToScheme,
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
	GetClusterScoped(ctx context.Context, resourceName, kindName, name string, opts metav1.GetOptions, result runtime.Object) (runtime.Object, error)
	ListClusterScoped(ctx context.Context, resourceName, kindName string, opts metav1.ListOptions, result runtime.Object) (runtime.Object, error)
	CreateClusterScoped(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.CreateOptions, result runtime.Object) (runtime.Object, error)
	UpdateClusterScoped(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error)
	UpdateStatusClusterScoped(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error)
	DeleteClusterScoped(ctx context.Context, gvr schema.GroupVersionResource, name string, opts metav1.DeleteOptions) error
	WatchClusterScoped(ctx context.Context, gvr schema.GroupVersionResource, opts metav1.ListOptions) (watch.Interface, error)
}

type Set struct {
	BlogV1alpha1 *BlogV1alpha1
	BlogV1alpha2 *BlogV1alpha2
}

func NewSet(cfg *rest.Config) (*Set, error) {
	s := &Set{}
	{
		conf := *cfg
		conf.GroupVersion = &blogv1alpha1.SchemaGroupVersion
		conf.APIPath = "/apis"
		conf.NegotiatedSerializer = Codecs.WithoutConversion()
		c, err := rest.RESTClientFor(&conf)
		if err != nil {
			return nil, err
		}
		s.BlogV1alpha1 = NewBlogV1alpha1Client(&restBackend{client: c})
	}
	{
		conf := *cfg
		conf.GroupVersion = &blogv1alpha2.SchemaGroupVersion
		conf.APIPath = "/apis"
		conf.NegotiatedSerializer = Codecs.WithoutConversion()
		c, err := rest.RESTClientFor(&conf)
		if err != nil {
			return nil, err
		}
		s.BlogV1alpha2 = NewBlogV1alpha2Client(&restBackend{client: c})
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
	if opts.TimeoutSeconds > 0 {
		timeout = time.Duration(opts.TimeoutSeconds) * time.Second
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
	if opts.TimeoutSeconds > 0 {
		timeout = time.Duration(opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return r.client.Get().
		Namespace(namespace).
		Resource(gvr.Resource).
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

func (r *restBackend) GetClusterScoped(ctx context.Context, resourceName, kindName, name string, opts metav1.GetOptions, result runtime.Object) (runtime.Object, error) {
	return result, r.client.Get().
		Resource(resourceName).
		Name(name).
		VersionedParams(&opts, ParameterCodec).
		Do(ctx).
		Into(result)
}

func (r *restBackend) ListClusterScoped(ctx context.Context, resourceName, kindName string, opts metav1.ListOptions, result runtime.Object) (runtime.Object, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds > 0 {
		timeout = time.Duration(opts.TimeoutSeconds) * time.Second
	}
	return result, r.client.Get().
		Resource(resourceName).
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
}

func (r *restBackend) CreateClusterScoped(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.CreateOptions, result runtime.Object) (runtime.Object, error) {
	return result, r.client.Post().
		Resource(resourceName).
		VersionedParams(&opts, ParameterCodec).
		Body(obj).
		Do(ctx).
		Into(result)
}

func (r *restBackend) UpdateClusterScoped(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	m := obj.(metav1.Object)
	if m == nil {
		return nil, errors.New("obj is not implement metav1.Object")
	}
	return result, r.client.Put().
		Resource(resourceName).
		Name(m.GetName()).
		VersionedParams(&opts, ParameterCodec).
		Body(obj).
		Do(ctx).
		Into(result)
}

func (r *restBackend) UpdateStatusClusterScoped(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	m := obj.(metav1.Object)
	if m == nil {
		return nil, errors.New("obj is not implement metav1.Object")
	}
	return result, r.client.Put().
		Resource(resourceName).
		Name(m.GetName()).
		SubResource("status").
		VersionedParams(&opts, ParameterCodec).
		Body(obj).
		Do(ctx).
		Into(result)
}

func (r *restBackend) DeleteClusterScoped(ctx context.Context, gvr schema.GroupVersionResource, name string, opts metav1.DeleteOptions) error {
	return r.client.Delete().
		Resource(gvr.Resource).
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (r *restBackend) WatchClusterScoped(ctx context.Context, gvr schema.GroupVersionResource, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds > 0 {
		timeout = time.Duration(opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return r.client.Get().
		Resource(gvr.Resource).
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

type BlogV1alpha1 struct {
	backend Backend
}

func NewBlogV1alpha1Client(b Backend) *BlogV1alpha1 {
	return &BlogV1alpha1{backend: b}
}

func (c *BlogV1alpha1) GetBlog(ctx context.Context, name string, opts metav1.GetOptions) (*blogv1alpha1.Blog, error) {
	result, err := c.backend.GetClusterScoped(ctx, "blogs", "Blog", name, opts, &blogv1alpha1.Blog{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha1.Blog), nil
}

func (c *BlogV1alpha1) CreateBlog(ctx context.Context, v *blogv1alpha1.Blog, opts metav1.CreateOptions) (*blogv1alpha1.Blog, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "blogs", "Blog", v, opts, &blogv1alpha1.Blog{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha1.Blog), nil
}

func (c *BlogV1alpha1) UpdateBlog(ctx context.Context, v *blogv1alpha1.Blog, opts metav1.UpdateOptions) (*blogv1alpha1.Blog, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "blogs", "Blog", v, opts, &blogv1alpha1.Blog{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha1.Blog), nil
}

func (c *BlogV1alpha1) DeleteBlog(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: "blog.f110.dev", Version: "v1alpha1", Resource: "blogs"}, name, opts)
}

func (c *BlogV1alpha1) ListBlog(ctx context.Context, opts metav1.ListOptions) (*blogv1alpha1.BlogList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "blogs", "Blog", opts, &blogv1alpha1.BlogList{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha1.BlogList), nil
}

func (c *BlogV1alpha1) WatchBlog(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: "blog.f110.dev", Version: "v1alpha1", Resource: "blogs"}, opts)
}

func (c *BlogV1alpha1) GetPost(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*blogv1alpha1.Post, error) {
	result, err := c.backend.Get(ctx, "posts", "Post", namespace, name, opts, &blogv1alpha1.Post{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha1.Post), nil
}

func (c *BlogV1alpha1) CreatePost(ctx context.Context, v *blogv1alpha1.Post, opts metav1.CreateOptions) (*blogv1alpha1.Post, error) {
	result, err := c.backend.Create(ctx, "posts", "Post", v, opts, &blogv1alpha1.Post{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha1.Post), nil
}

func (c *BlogV1alpha1) UpdatePost(ctx context.Context, v *blogv1alpha1.Post, opts metav1.UpdateOptions) (*blogv1alpha1.Post, error) {
	result, err := c.backend.Update(ctx, "posts", "Post", v, opts, &blogv1alpha1.Post{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha1.Post), nil
}

func (c *BlogV1alpha1) DeletePost(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: "blog.f110.dev", Version: "v1alpha1", Resource: "posts"}, namespace, name, opts)
}

func (c *BlogV1alpha1) ListPost(ctx context.Context, namespace string, opts metav1.ListOptions) (*blogv1alpha1.PostList, error) {
	result, err := c.backend.List(ctx, "posts", "Post", namespace, opts, &blogv1alpha1.PostList{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha1.PostList), nil
}

func (c *BlogV1alpha1) WatchPost(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: "blog.f110.dev", Version: "v1alpha1", Resource: "posts"}, namespace, opts)
}

type BlogV1alpha2 struct {
	backend Backend
}

func NewBlogV1alpha2Client(b Backend) *BlogV1alpha2 {
	return &BlogV1alpha2{backend: b}
}

func (c *BlogV1alpha2) GetAuthor(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*blogv1alpha2.Author, error) {
	result, err := c.backend.Get(ctx, "authors", "Author", namespace, name, opts, &blogv1alpha2.Author{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha2.Author), nil
}

func (c *BlogV1alpha2) CreateAuthor(ctx context.Context, v *blogv1alpha2.Author, opts metav1.CreateOptions) (*blogv1alpha2.Author, error) {
	result, err := c.backend.Create(ctx, "authors", "Author", v, opts, &blogv1alpha2.Author{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha2.Author), nil
}

func (c *BlogV1alpha2) UpdateAuthor(ctx context.Context, v *blogv1alpha2.Author, opts metav1.UpdateOptions) (*blogv1alpha2.Author, error) {
	result, err := c.backend.Update(ctx, "authors", "Author", v, opts, &blogv1alpha2.Author{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha2.Author), nil
}

func (c *BlogV1alpha2) DeleteAuthor(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: "blog.f110.dev", Version: "v1alpha2", Resource: "authors"}, namespace, name, opts)
}

func (c *BlogV1alpha2) ListAuthor(ctx context.Context, namespace string, opts metav1.ListOptions) (*blogv1alpha2.AuthorList, error) {
	result, err := c.backend.List(ctx, "authors", "Author", namespace, opts, &blogv1alpha2.AuthorList{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha2.AuthorList), nil
}

func (c *BlogV1alpha2) WatchAuthor(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: "blog.f110.dev", Version: "v1alpha2", Resource: "authors"}, namespace, opts)
}

func (c *BlogV1alpha2) GetBlog(ctx context.Context, name string, opts metav1.GetOptions) (*blogv1alpha2.Blog, error) {
	result, err := c.backend.GetClusterScoped(ctx, "blogs", "Blog", name, opts, &blogv1alpha2.Blog{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha2.Blog), nil
}

func (c *BlogV1alpha2) CreateBlog(ctx context.Context, v *blogv1alpha2.Blog, opts metav1.CreateOptions) (*blogv1alpha2.Blog, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "blogs", "Blog", v, opts, &blogv1alpha2.Blog{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha2.Blog), nil
}

func (c *BlogV1alpha2) UpdateBlog(ctx context.Context, v *blogv1alpha2.Blog, opts metav1.UpdateOptions) (*blogv1alpha2.Blog, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "blogs", "Blog", v, opts, &blogv1alpha2.Blog{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha2.Blog), nil
}

func (c *BlogV1alpha2) DeleteBlog(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: "blog.f110.dev", Version: "v1alpha2", Resource: "blogs"}, name, opts)
}

func (c *BlogV1alpha2) ListBlog(ctx context.Context, opts metav1.ListOptions) (*blogv1alpha2.BlogList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "blogs", "Blog", opts, &blogv1alpha2.BlogList{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha2.BlogList), nil
}

func (c *BlogV1alpha2) WatchBlog(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: "blog.f110.dev", Version: "v1alpha2", Resource: "blogs"}, opts)
}

func (c *BlogV1alpha2) GetPost(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*blogv1alpha2.Post, error) {
	result, err := c.backend.Get(ctx, "posts", "Post", namespace, name, opts, &blogv1alpha2.Post{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha2.Post), nil
}

func (c *BlogV1alpha2) CreatePost(ctx context.Context, v *blogv1alpha2.Post, opts metav1.CreateOptions) (*blogv1alpha2.Post, error) {
	result, err := c.backend.Create(ctx, "posts", "Post", v, opts, &blogv1alpha2.Post{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha2.Post), nil
}

func (c *BlogV1alpha2) UpdatePost(ctx context.Context, v *blogv1alpha2.Post, opts metav1.UpdateOptions) (*blogv1alpha2.Post, error) {
	result, err := c.backend.Update(ctx, "posts", "Post", v, opts, &blogv1alpha2.Post{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha2.Post), nil
}

func (c *BlogV1alpha2) DeletePost(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: "blog.f110.dev", Version: "v1alpha2", Resource: "posts"}, namespace, name, opts)
}

func (c *BlogV1alpha2) ListPost(ctx context.Context, namespace string, opts metav1.ListOptions) (*blogv1alpha2.PostList, error) {
	result, err := c.backend.List(ctx, "posts", "Post", namespace, opts, &blogv1alpha2.PostList{})
	if err != nil {
		return nil, err
	}
	return result.(*blogv1alpha2.PostList), nil
}

func (c *BlogV1alpha2) WatchPost(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: "blog.f110.dev", Version: "v1alpha2", Resource: "posts"}, namespace, opts)
}

type InformerCache struct {
	mu        sync.Mutex
	informers map[reflect.Type]cache.SharedIndexInformer
}

func NewInformerCache() *InformerCache {
	return &InformerCache{informers: make(map[reflect.Type]cache.SharedIndexInformer)}
}

func (c *InformerCache) Write(obj runtime.Object, newFunc func() cache.SharedIndexInformer) cache.SharedIndexInformer {
	c.mu.Lock()
	defer c.mu.Unlock()

	typ := reflect.TypeOf(obj)
	if v, ok := c.informers[typ]; ok {
		return v
	}
	informer := newFunc()
	c.informers[typ] = informer

	return informer
}

func (c *InformerCache) Informers() []cache.SharedIndexInformer {
	c.mu.Lock()
	defer c.mu.Unlock()

	a := make([]cache.SharedIndexInformer, 0, len(c.informers))
	for _, v := range c.informers {
		a = append(a, v)
	}

	return a
}

type InformerFactory struct {
	set   *Set
	cache *InformerCache

	namespace    string
	resyncPeriod time.Duration
}

func NewInformerFactory(s *Set, c *InformerCache, namespace string, resyncPeriod time.Duration) *InformerFactory {
	return &InformerFactory{set: s, cache: c, namespace: namespace, resyncPeriod: resyncPeriod}
}

func (f *InformerFactory) Cache() *InformerCache {
	return f.cache
}

func (f *InformerFactory) InformerFor(obj runtime.Object) cache.SharedIndexInformer {
	switch obj.(type) {
	case *blogv1alpha1.Blog:
		return NewBlogV1alpha1Informer(f.cache, f.set.BlogV1alpha1, f.namespace, f.resyncPeriod).BlogInformer()
	case *blogv1alpha1.Post:
		return NewBlogV1alpha1Informer(f.cache, f.set.BlogV1alpha1, f.namespace, f.resyncPeriod).PostInformer()
	case *blogv1alpha2.Author:
		return NewBlogV1alpha2Informer(f.cache, f.set.BlogV1alpha2, f.namespace, f.resyncPeriod).AuthorInformer()
	case *blogv1alpha2.Blog:
		return NewBlogV1alpha2Informer(f.cache, f.set.BlogV1alpha2, f.namespace, f.resyncPeriod).BlogInformer()
	case *blogv1alpha2.Post:
		return NewBlogV1alpha2Informer(f.cache, f.set.BlogV1alpha2, f.namespace, f.resyncPeriod).PostInformer()
	default:
		return nil
	}
}

func (f *InformerFactory) InformerForResource(gvr schema.GroupVersionResource) cache.SharedIndexInformer {
	switch gvr {
	case blogv1alpha1.SchemaGroupVersion.WithResource("blogs"):
		return NewBlogV1alpha1Informer(f.cache, f.set.BlogV1alpha1, f.namespace, f.resyncPeriod).BlogInformer()
	case blogv1alpha1.SchemaGroupVersion.WithResource("posts"):
		return NewBlogV1alpha1Informer(f.cache, f.set.BlogV1alpha1, f.namespace, f.resyncPeriod).PostInformer()
	case blogv1alpha2.SchemaGroupVersion.WithResource("authors"):
		return NewBlogV1alpha2Informer(f.cache, f.set.BlogV1alpha2, f.namespace, f.resyncPeriod).AuthorInformer()
	case blogv1alpha2.SchemaGroupVersion.WithResource("blogs"):
		return NewBlogV1alpha2Informer(f.cache, f.set.BlogV1alpha2, f.namespace, f.resyncPeriod).BlogInformer()
	case blogv1alpha2.SchemaGroupVersion.WithResource("posts"):
		return NewBlogV1alpha2Informer(f.cache, f.set.BlogV1alpha2, f.namespace, f.resyncPeriod).PostInformer()
	default:
		return nil
	}
}

func (f *InformerFactory) Run(ctx context.Context) {
	for _, v := range f.cache.Informers() {
		go v.Run(ctx.Done())
	}
}

type BlogV1alpha1Informer struct {
	cache        *InformerCache
	client       *BlogV1alpha1
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewBlogV1alpha1Informer(c *InformerCache, client *BlogV1alpha1, namespace string, resyncPeriod time.Duration) *BlogV1alpha1Informer {
	return &BlogV1alpha1Informer{
		cache:        c,
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *BlogV1alpha1Informer) BlogInformer() cache.SharedIndexInformer {
	return f.cache.Write(&blogv1alpha1.Blog{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options k8smetav1.ListOptions) (runtime.Object, error) {
					return f.client.ListBlog(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options k8smetav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchBlog(context.TODO(), metav1.ListOptions{})
				},
			},
			&blogv1alpha1.Blog{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *BlogV1alpha1Informer) BlogLister() *BlogV1alpha1BlogLister {
	return NewBlogV1alpha1BlogLister(f.BlogInformer().GetIndexer())
}

func (f *BlogV1alpha1Informer) PostInformer() cache.SharedIndexInformer {
	return f.cache.Write(&blogv1alpha1.Post{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options k8smetav1.ListOptions) (runtime.Object, error) {
					return f.client.ListPost(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options k8smetav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchPost(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&blogv1alpha1.Post{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *BlogV1alpha1Informer) PostLister() *BlogV1alpha1PostLister {
	return NewBlogV1alpha1PostLister(f.PostInformer().GetIndexer())
}

type BlogV1alpha2Informer struct {
	cache        *InformerCache
	client       *BlogV1alpha2
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewBlogV1alpha2Informer(c *InformerCache, client *BlogV1alpha2, namespace string, resyncPeriod time.Duration) *BlogV1alpha2Informer {
	return &BlogV1alpha2Informer{
		cache:        c,
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *BlogV1alpha2Informer) AuthorInformer() cache.SharedIndexInformer {
	return f.cache.Write(&blogv1alpha2.Author{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options k8smetav1.ListOptions) (runtime.Object, error) {
					return f.client.ListAuthor(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options k8smetav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchAuthor(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&blogv1alpha2.Author{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *BlogV1alpha2Informer) AuthorLister() *BlogV1alpha2AuthorLister {
	return NewBlogV1alpha2AuthorLister(f.AuthorInformer().GetIndexer())
}

func (f *BlogV1alpha2Informer) BlogInformer() cache.SharedIndexInformer {
	return f.cache.Write(&blogv1alpha2.Blog{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options k8smetav1.ListOptions) (runtime.Object, error) {
					return f.client.ListBlog(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options k8smetav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchBlog(context.TODO(), metav1.ListOptions{})
				},
			},
			&blogv1alpha2.Blog{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *BlogV1alpha2Informer) BlogLister() *BlogV1alpha2BlogLister {
	return NewBlogV1alpha2BlogLister(f.BlogInformer().GetIndexer())
}

func (f *BlogV1alpha2Informer) PostInformer() cache.SharedIndexInformer {
	return f.cache.Write(&blogv1alpha2.Post{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options k8smetav1.ListOptions) (runtime.Object, error) {
					return f.client.ListPost(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options k8smetav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchPost(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&blogv1alpha2.Post{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *BlogV1alpha2Informer) PostLister() *BlogV1alpha2PostLister {
	return NewBlogV1alpha2PostLister(f.PostInformer().GetIndexer())
}

type BlogV1alpha1BlogLister struct {
	indexer cache.Indexer
}

func NewBlogV1alpha1BlogLister(indexer cache.Indexer) *BlogV1alpha1BlogLister {
	return &BlogV1alpha1BlogLister{indexer: indexer}
}

func (x *BlogV1alpha1BlogLister) List(selector labels.Selector) ([]*blogv1alpha1.Blog, error) {
	var ret []*blogv1alpha1.Blog
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*blogv1alpha1.Blog).DeepCopy())
	})
	return ret, err
}

func (x *BlogV1alpha1BlogLister) Get(name string) (*blogv1alpha1.Blog, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(blogv1alpha1.SchemaGroupVersion.WithResource("blog").GroupResource(), name)
	}
	return obj.(*blogv1alpha1.Blog).DeepCopy(), nil
}

type BlogV1alpha1PostLister struct {
	indexer cache.Indexer
}

func NewBlogV1alpha1PostLister(indexer cache.Indexer) *BlogV1alpha1PostLister {
	return &BlogV1alpha1PostLister{indexer: indexer}
}

func (x *BlogV1alpha1PostLister) List(namespace string, selector labels.Selector) ([]*blogv1alpha1.Post, error) {
	var ret []*blogv1alpha1.Post
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*blogv1alpha1.Post).DeepCopy())
	})
	return ret, err
}

func (x *BlogV1alpha1PostLister) Get(namespace, name string) (*blogv1alpha1.Post, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(blogv1alpha1.SchemaGroupVersion.WithResource("post").GroupResource(), name)
	}
	return obj.(*blogv1alpha1.Post).DeepCopy(), nil
}

type BlogV1alpha2AuthorLister struct {
	indexer cache.Indexer
}

func NewBlogV1alpha2AuthorLister(indexer cache.Indexer) *BlogV1alpha2AuthorLister {
	return &BlogV1alpha2AuthorLister{indexer: indexer}
}

func (x *BlogV1alpha2AuthorLister) List(namespace string, selector labels.Selector) ([]*blogv1alpha2.Author, error) {
	var ret []*blogv1alpha2.Author
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*blogv1alpha2.Author).DeepCopy())
	})
	return ret, err
}

func (x *BlogV1alpha2AuthorLister) Get(namespace, name string) (*blogv1alpha2.Author, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(blogv1alpha2.SchemaGroupVersion.WithResource("author").GroupResource(), name)
	}
	return obj.(*blogv1alpha2.Author).DeepCopy(), nil
}

type BlogV1alpha2BlogLister struct {
	indexer cache.Indexer
}

func NewBlogV1alpha2BlogLister(indexer cache.Indexer) *BlogV1alpha2BlogLister {
	return &BlogV1alpha2BlogLister{indexer: indexer}
}

func (x *BlogV1alpha2BlogLister) List(selector labels.Selector) ([]*blogv1alpha2.Blog, error) {
	var ret []*blogv1alpha2.Blog
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*blogv1alpha2.Blog).DeepCopy())
	})
	return ret, err
}

func (x *BlogV1alpha2BlogLister) Get(name string) (*blogv1alpha2.Blog, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(blogv1alpha2.SchemaGroupVersion.WithResource("blog").GroupResource(), name)
	}
	return obj.(*blogv1alpha2.Blog).DeepCopy(), nil
}

type BlogV1alpha2PostLister struct {
	indexer cache.Indexer
}

func NewBlogV1alpha2PostLister(indexer cache.Indexer) *BlogV1alpha2PostLister {
	return &BlogV1alpha2PostLister{indexer: indexer}
}

func (x *BlogV1alpha2PostLister) List(namespace string, selector labels.Selector) ([]*blogv1alpha2.Post, error) {
	var ret []*blogv1alpha2.Post
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*blogv1alpha2.Post).DeepCopy())
	})
	return ret, err
}

func (x *BlogV1alpha2PostLister) Get(namespace, name string) (*blogv1alpha2.Post, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(blogv1alpha2.SchemaGroupVersion.WithResource("post").GroupResource(), name)
	}
	return obj.(*blogv1alpha2.Post).DeepCopy(), nil
}
