package k8sclient

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
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"go.f110.dev/kubeproto/go/apis/admissionregistrationv1"
	"go.f110.dev/kubeproto/go/apis/appsv1"
	"go.f110.dev/kubeproto/go/apis/authenticationv1"
	"go.f110.dev/kubeproto/go/apis/authorizationv1"
	"go.f110.dev/kubeproto/go/apis/autoscalingv1"
	"go.f110.dev/kubeproto/go/apis/batchv1"
	"go.f110.dev/kubeproto/go/apis/certificatesv1"
	"go.f110.dev/kubeproto/go/apis/corev1"
	"go.f110.dev/kubeproto/go/apis/discoveryv1"
	"go.f110.dev/kubeproto/go/apis/networkingv1"
	"go.f110.dev/kubeproto/go/apis/policyv1"
	"go.f110.dev/kubeproto/go/apis/rbacv1"
)

var (
	Scheme         = runtime.NewScheme()
	ParameterCodec = runtime.NewParameterCodec(Scheme)
	Codecs         = serializer.NewCodecFactory(Scheme)
	AddToScheme    = localSchemeBuilder.AddToScheme
)

var localSchemeBuilder = runtime.SchemeBuilder{
	corev1.AddToScheme,
	admissionregistrationv1.AddToScheme,
	appsv1.AddToScheme,
	authenticationv1.AddToScheme,
	authorizationv1.AddToScheme,
	autoscalingv1.AddToScheme,
	certificatesv1.AddToScheme,
	discoveryv1.AddToScheme,
	networkingv1.AddToScheme,
	policyv1.AddToScheme,
	rbacv1.AddToScheme,
}

func init() {
	for _, v := range []func(*runtime.Scheme) error{
		corev1.AddToScheme,
		admissionregistrationv1.AddToScheme,
		appsv1.AddToScheme,
		authenticationv1.AddToScheme,
		authorizationv1.AddToScheme,
		autoscalingv1.AddToScheme,
		certificatesv1.AddToScheme,
		discoveryv1.AddToScheme,
		networkingv1.AddToScheme,
		policyv1.AddToScheme,
		rbacv1.AddToScheme,
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
	CoreV1                       *CoreV1
	AdmissionregistrationK8sIoV1 *AdmissionregistrationK8sIoV1
	AppsV1                       *AppsV1
	AuthenticationK8sIoV1        *AuthenticationK8sIoV1
	AuthorizationK8sIoV1         *AuthorizationK8sIoV1
	AutoscalingV1                *AutoscalingV1
	CertificatesK8sIoV1          *CertificatesK8sIoV1
	DiscoveryK8sIoV1             *DiscoveryK8sIoV1
	NetworkingK8sIoV1            *NetworkingK8sIoV1
	PolicyV1                     *PolicyV1
	RbacAuthorizationK8sIoV1     *RbacAuthorizationK8sIoV1
}

func NewSet(cfg *rest.Config) (*Set, error) {
	s := &Set{}
	{
		conf := *cfg
		conf.GroupVersion = &corev1.SchemaGroupVersion
		conf.APIPath = "/apis"
		conf.NegotiatedSerializer = Codecs.WithoutConversion()
		c, err := rest.RESTClientFor(&conf)
		if err != nil {
			return nil, err
		}
		s.CoreV1 = NewCoreV1Client(&restBackend{client: c})
	}
	{
		conf := *cfg
		conf.GroupVersion = &admissionregistrationv1.SchemaGroupVersion
		conf.APIPath = "/apis"
		conf.NegotiatedSerializer = Codecs.WithoutConversion()
		c, err := rest.RESTClientFor(&conf)
		if err != nil {
			return nil, err
		}
		s.AdmissionregistrationK8sIoV1 = NewAdmissionregistrationK8sIoV1Client(&restBackend{client: c})
	}
	{
		conf := *cfg
		conf.GroupVersion = &appsv1.SchemaGroupVersion
		conf.APIPath = "/apis"
		conf.NegotiatedSerializer = Codecs.WithoutConversion()
		c, err := rest.RESTClientFor(&conf)
		if err != nil {
			return nil, err
		}
		s.AppsV1 = NewAppsV1Client(&restBackend{client: c})
	}
	{
		conf := *cfg
		conf.GroupVersion = &authenticationv1.SchemaGroupVersion
		conf.APIPath = "/apis"
		conf.NegotiatedSerializer = Codecs.WithoutConversion()
		c, err := rest.RESTClientFor(&conf)
		if err != nil {
			return nil, err
		}
		s.AuthenticationK8sIoV1 = NewAuthenticationK8sIoV1Client(&restBackend{client: c})
	}
	{
		conf := *cfg
		conf.GroupVersion = &authorizationv1.SchemaGroupVersion
		conf.APIPath = "/apis"
		conf.NegotiatedSerializer = Codecs.WithoutConversion()
		c, err := rest.RESTClientFor(&conf)
		if err != nil {
			return nil, err
		}
		s.AuthorizationK8sIoV1 = NewAuthorizationK8sIoV1Client(&restBackend{client: c})
	}
	{
		conf := *cfg
		conf.GroupVersion = &autoscalingv1.SchemaGroupVersion
		conf.APIPath = "/apis"
		conf.NegotiatedSerializer = Codecs.WithoutConversion()
		c, err := rest.RESTClientFor(&conf)
		if err != nil {
			return nil, err
		}
		s.AutoscalingV1 = NewAutoscalingV1Client(&restBackend{client: c})
	}
	{
		conf := *cfg
		conf.GroupVersion = &certificatesv1.SchemaGroupVersion
		conf.APIPath = "/apis"
		conf.NegotiatedSerializer = Codecs.WithoutConversion()
		c, err := rest.RESTClientFor(&conf)
		if err != nil {
			return nil, err
		}
		s.CertificatesK8sIoV1 = NewCertificatesK8sIoV1Client(&restBackend{client: c})
	}
	{
		conf := *cfg
		conf.GroupVersion = &discoveryv1.SchemaGroupVersion
		conf.APIPath = "/apis"
		conf.NegotiatedSerializer = Codecs.WithoutConversion()
		c, err := rest.RESTClientFor(&conf)
		if err != nil {
			return nil, err
		}
		s.DiscoveryK8sIoV1 = NewDiscoveryK8sIoV1Client(&restBackend{client: c})
	}
	{
		conf := *cfg
		conf.GroupVersion = &networkingv1.SchemaGroupVersion
		conf.APIPath = "/apis"
		conf.NegotiatedSerializer = Codecs.WithoutConversion()
		c, err := rest.RESTClientFor(&conf)
		if err != nil {
			return nil, err
		}
		s.NetworkingK8sIoV1 = NewNetworkingK8sIoV1Client(&restBackend{client: c})
	}
	{
		conf := *cfg
		conf.GroupVersion = &policyv1.SchemaGroupVersion
		conf.APIPath = "/apis"
		conf.NegotiatedSerializer = Codecs.WithoutConversion()
		c, err := rest.RESTClientFor(&conf)
		if err != nil {
			return nil, err
		}
		s.PolicyV1 = NewPolicyV1Client(&restBackend{client: c})
	}
	{
		conf := *cfg
		conf.GroupVersion = &rbacv1.SchemaGroupVersion
		conf.APIPath = "/apis"
		conf.NegotiatedSerializer = Codecs.WithoutConversion()
		c, err := rest.RESTClientFor(&conf)
		if err != nil {
			return nil, err
		}
		s.RbacAuthorizationK8sIoV1 = NewRbacAuthorizationK8sIoV1Client(&restBackend{client: c})
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
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
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
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return r.client.Get().
		Resource(gvr.Resource).
		VersionedParams(&opts, ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

type CoreV1 struct {
	backend Backend
}

func NewCoreV1Client(b Backend) *CoreV1 {
	return &CoreV1{backend: b}
}

func (c *CoreV1) GetBinding(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Binding, error) {
	result, err := c.backend.Get(ctx, "bindings", "Binding", namespace, name, opts, &corev1.Binding{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Binding), nil
}

func (c *CoreV1) CreateBinding(ctx context.Context, v *corev1.Binding, opts metav1.CreateOptions) (*corev1.Binding, error) {
	result, err := c.backend.Create(ctx, "bindings", "Binding", v, opts, &corev1.Binding{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Binding), nil
}

func (c *CoreV1) UpdateBinding(ctx context.Context, v *corev1.Binding, opts metav1.UpdateOptions) (*corev1.Binding, error) {
	result, err := c.backend.Update(ctx, "bindings", "Binding", v, opts, &corev1.Binding{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Binding), nil
}

func (c *CoreV1) DeleteBinding(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "bindings"}, namespace, name, opts)
}

func (c *CoreV1) ListBinding(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.BindingList, error) {
	result, err := c.backend.List(ctx, "bindings", "Binding", namespace, opts, &corev1.BindingList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.BindingList), nil
}

func (c *CoreV1) WatchBinding(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "bindings"}, namespace, opts)
}

func (c *CoreV1) GetComponentStatus(ctx context.Context, name string, opts metav1.GetOptions) (*corev1.ComponentStatus, error) {
	result, err := c.backend.GetClusterScoped(ctx, "componentstatuses", "ComponentStatus", name, opts, &corev1.ComponentStatus{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ComponentStatus), nil
}

func (c *CoreV1) CreateComponentStatus(ctx context.Context, v *corev1.ComponentStatus, opts metav1.CreateOptions) (*corev1.ComponentStatus, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "componentstatuses", "ComponentStatus", v, opts, &corev1.ComponentStatus{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ComponentStatus), nil
}

func (c *CoreV1) UpdateComponentStatus(ctx context.Context, v *corev1.ComponentStatus, opts metav1.UpdateOptions) (*corev1.ComponentStatus, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "componentstatuses", "ComponentStatus", v, opts, &corev1.ComponentStatus{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ComponentStatus), nil
}

func (c *CoreV1) DeleteComponentStatus(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "componentstatuses"}, name, opts)
}

func (c *CoreV1) ListComponentStatus(ctx context.Context, opts metav1.ListOptions) (*corev1.ComponentStatusList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "componentstatuses", "ComponentStatus", opts, &corev1.ComponentStatusList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ComponentStatusList), nil
}

func (c *CoreV1) WatchComponentStatus(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "componentstatuses"}, opts)
}

func (c *CoreV1) GetConfigMap(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.ConfigMap, error) {
	result, err := c.backend.Get(ctx, "configmaps", "ConfigMap", namespace, name, opts, &corev1.ConfigMap{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ConfigMap), nil
}

func (c *CoreV1) CreateConfigMap(ctx context.Context, v *corev1.ConfigMap, opts metav1.CreateOptions) (*corev1.ConfigMap, error) {
	result, err := c.backend.Create(ctx, "configmaps", "ConfigMap", v, opts, &corev1.ConfigMap{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ConfigMap), nil
}

func (c *CoreV1) UpdateConfigMap(ctx context.Context, v *corev1.ConfigMap, opts metav1.UpdateOptions) (*corev1.ConfigMap, error) {
	result, err := c.backend.Update(ctx, "configmaps", "ConfigMap", v, opts, &corev1.ConfigMap{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ConfigMap), nil
}

func (c *CoreV1) DeleteConfigMap(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "configmaps"}, namespace, name, opts)
}

func (c *CoreV1) ListConfigMap(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.ConfigMapList, error) {
	result, err := c.backend.List(ctx, "configmaps", "ConfigMap", namespace, opts, &corev1.ConfigMapList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ConfigMapList), nil
}

func (c *CoreV1) WatchConfigMap(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "configmaps"}, namespace, opts)
}

func (c *CoreV1) GetEndpoints(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Endpoints, error) {
	result, err := c.backend.Get(ctx, "endpoints", "Endpoints", namespace, name, opts, &corev1.Endpoints{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Endpoints), nil
}

func (c *CoreV1) CreateEndpoints(ctx context.Context, v *corev1.Endpoints, opts metav1.CreateOptions) (*corev1.Endpoints, error) {
	result, err := c.backend.Create(ctx, "endpoints", "Endpoints", v, opts, &corev1.Endpoints{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Endpoints), nil
}

func (c *CoreV1) UpdateEndpoints(ctx context.Context, v *corev1.Endpoints, opts metav1.UpdateOptions) (*corev1.Endpoints, error) {
	result, err := c.backend.Update(ctx, "endpoints", "Endpoints", v, opts, &corev1.Endpoints{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Endpoints), nil
}

func (c *CoreV1) DeleteEndpoints(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "endpoints"}, namespace, name, opts)
}

func (c *CoreV1) ListEndpoints(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.EndpointsList, error) {
	result, err := c.backend.List(ctx, "endpoints", "Endpoints", namespace, opts, &corev1.EndpointsList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.EndpointsList), nil
}

func (c *CoreV1) WatchEndpoints(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "endpoints"}, namespace, opts)
}

func (c *CoreV1) GetEvent(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Event, error) {
	result, err := c.backend.Get(ctx, "events", "Event", namespace, name, opts, &corev1.Event{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Event), nil
}

func (c *CoreV1) CreateEvent(ctx context.Context, v *corev1.Event, opts metav1.CreateOptions) (*corev1.Event, error) {
	result, err := c.backend.Create(ctx, "events", "Event", v, opts, &corev1.Event{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Event), nil
}

func (c *CoreV1) UpdateEvent(ctx context.Context, v *corev1.Event, opts metav1.UpdateOptions) (*corev1.Event, error) {
	result, err := c.backend.Update(ctx, "events", "Event", v, opts, &corev1.Event{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Event), nil
}

func (c *CoreV1) DeleteEvent(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "events"}, namespace, name, opts)
}

func (c *CoreV1) ListEvent(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.EventList, error) {
	result, err := c.backend.List(ctx, "events", "Event", namespace, opts, &corev1.EventList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.EventList), nil
}

func (c *CoreV1) WatchEvent(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "events"}, namespace, opts)
}

func (c *CoreV1) GetLimitRange(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.LimitRange, error) {
	result, err := c.backend.Get(ctx, "limitranges", "LimitRange", namespace, name, opts, &corev1.LimitRange{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.LimitRange), nil
}

func (c *CoreV1) CreateLimitRange(ctx context.Context, v *corev1.LimitRange, opts metav1.CreateOptions) (*corev1.LimitRange, error) {
	result, err := c.backend.Create(ctx, "limitranges", "LimitRange", v, opts, &corev1.LimitRange{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.LimitRange), nil
}

func (c *CoreV1) UpdateLimitRange(ctx context.Context, v *corev1.LimitRange, opts metav1.UpdateOptions) (*corev1.LimitRange, error) {
	result, err := c.backend.Update(ctx, "limitranges", "LimitRange", v, opts, &corev1.LimitRange{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.LimitRange), nil
}

func (c *CoreV1) DeleteLimitRange(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "limitranges"}, namespace, name, opts)
}

func (c *CoreV1) ListLimitRange(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.LimitRangeList, error) {
	result, err := c.backend.List(ctx, "limitranges", "LimitRange", namespace, opts, &corev1.LimitRangeList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.LimitRangeList), nil
}

func (c *CoreV1) WatchLimitRange(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "limitranges"}, namespace, opts)
}

func (c *CoreV1) GetNamespace(ctx context.Context, name string, opts metav1.GetOptions) (*corev1.Namespace, error) {
	result, err := c.backend.GetClusterScoped(ctx, "namespaces", "Namespace", name, opts, &corev1.Namespace{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Namespace), nil
}

func (c *CoreV1) CreateNamespace(ctx context.Context, v *corev1.Namespace, opts metav1.CreateOptions) (*corev1.Namespace, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "namespaces", "Namespace", v, opts, &corev1.Namespace{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Namespace), nil
}

func (c *CoreV1) UpdateNamespace(ctx context.Context, v *corev1.Namespace, opts metav1.UpdateOptions) (*corev1.Namespace, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "namespaces", "Namespace", v, opts, &corev1.Namespace{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Namespace), nil
}

func (c *CoreV1) DeleteNamespace(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "namespaces"}, name, opts)
}

func (c *CoreV1) ListNamespace(ctx context.Context, opts metav1.ListOptions) (*corev1.NamespaceList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "namespaces", "Namespace", opts, &corev1.NamespaceList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.NamespaceList), nil
}

func (c *CoreV1) WatchNamespace(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "namespaces"}, opts)
}

func (c *CoreV1) GetNode(ctx context.Context, name string, opts metav1.GetOptions) (*corev1.Node, error) {
	result, err := c.backend.GetClusterScoped(ctx, "nodes", "Node", name, opts, &corev1.Node{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Node), nil
}

func (c *CoreV1) CreateNode(ctx context.Context, v *corev1.Node, opts metav1.CreateOptions) (*corev1.Node, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "nodes", "Node", v, opts, &corev1.Node{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Node), nil
}

func (c *CoreV1) UpdateNode(ctx context.Context, v *corev1.Node, opts metav1.UpdateOptions) (*corev1.Node, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "nodes", "Node", v, opts, &corev1.Node{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Node), nil
}

func (c *CoreV1) DeleteNode(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "nodes"}, name, opts)
}

func (c *CoreV1) ListNode(ctx context.Context, opts metav1.ListOptions) (*corev1.NodeList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "nodes", "Node", opts, &corev1.NodeList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.NodeList), nil
}

func (c *CoreV1) WatchNode(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "nodes"}, opts)
}

func (c *CoreV1) GetPersistentVolume(ctx context.Context, name string, opts metav1.GetOptions) (*corev1.PersistentVolume, error) {
	result, err := c.backend.GetClusterScoped(ctx, "persistentvolumes", "PersistentVolume", name, opts, &corev1.PersistentVolume{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PersistentVolume), nil
}

func (c *CoreV1) CreatePersistentVolume(ctx context.Context, v *corev1.PersistentVolume, opts metav1.CreateOptions) (*corev1.PersistentVolume, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "persistentvolumes", "PersistentVolume", v, opts, &corev1.PersistentVolume{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PersistentVolume), nil
}

func (c *CoreV1) UpdatePersistentVolume(ctx context.Context, v *corev1.PersistentVolume, opts metav1.UpdateOptions) (*corev1.PersistentVolume, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "persistentvolumes", "PersistentVolume", v, opts, &corev1.PersistentVolume{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PersistentVolume), nil
}

func (c *CoreV1) DeletePersistentVolume(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "persistentvolumes"}, name, opts)
}

func (c *CoreV1) ListPersistentVolume(ctx context.Context, opts metav1.ListOptions) (*corev1.PersistentVolumeList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "persistentvolumes", "PersistentVolume", opts, &corev1.PersistentVolumeList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PersistentVolumeList), nil
}

func (c *CoreV1) WatchPersistentVolume(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "persistentvolumes"}, opts)
}

func (c *CoreV1) GetPersistentVolumeClaim(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.PersistentVolumeClaim, error) {
	result, err := c.backend.Get(ctx, "persistentvolumeclaims", "PersistentVolumeClaim", namespace, name, opts, &corev1.PersistentVolumeClaim{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PersistentVolumeClaim), nil
}

func (c *CoreV1) CreatePersistentVolumeClaim(ctx context.Context, v *corev1.PersistentVolumeClaim, opts metav1.CreateOptions) (*corev1.PersistentVolumeClaim, error) {
	result, err := c.backend.Create(ctx, "persistentvolumeclaims", "PersistentVolumeClaim", v, opts, &corev1.PersistentVolumeClaim{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PersistentVolumeClaim), nil
}

func (c *CoreV1) UpdatePersistentVolumeClaim(ctx context.Context, v *corev1.PersistentVolumeClaim, opts metav1.UpdateOptions) (*corev1.PersistentVolumeClaim, error) {
	result, err := c.backend.Update(ctx, "persistentvolumeclaims", "PersistentVolumeClaim", v, opts, &corev1.PersistentVolumeClaim{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PersistentVolumeClaim), nil
}

func (c *CoreV1) DeletePersistentVolumeClaim(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "persistentvolumeclaims"}, namespace, name, opts)
}

func (c *CoreV1) ListPersistentVolumeClaim(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.PersistentVolumeClaimList, error) {
	result, err := c.backend.List(ctx, "persistentvolumeclaims", "PersistentVolumeClaim", namespace, opts, &corev1.PersistentVolumeClaimList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PersistentVolumeClaimList), nil
}

func (c *CoreV1) WatchPersistentVolumeClaim(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "persistentvolumeclaims"}, namespace, opts)
}

func (c *CoreV1) GetPod(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Pod, error) {
	result, err := c.backend.Get(ctx, "pods", "Pod", namespace, name, opts, &corev1.Pod{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Pod), nil
}

func (c *CoreV1) CreatePod(ctx context.Context, v *corev1.Pod, opts metav1.CreateOptions) (*corev1.Pod, error) {
	result, err := c.backend.Create(ctx, "pods", "Pod", v, opts, &corev1.Pod{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Pod), nil
}

func (c *CoreV1) UpdatePod(ctx context.Context, v *corev1.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error) {
	result, err := c.backend.Update(ctx, "pods", "Pod", v, opts, &corev1.Pod{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Pod), nil
}

func (c *CoreV1) DeletePod(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "pods"}, namespace, name, opts)
}

func (c *CoreV1) ListPod(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.PodList, error) {
	result, err := c.backend.List(ctx, "pods", "Pod", namespace, opts, &corev1.PodList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PodList), nil
}

func (c *CoreV1) WatchPod(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "pods"}, namespace, opts)
}

func (c *CoreV1) GetPodStatusResult(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.PodStatusResult, error) {
	result, err := c.backend.Get(ctx, "podstatusresults", "PodStatusResult", namespace, name, opts, &corev1.PodStatusResult{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PodStatusResult), nil
}

func (c *CoreV1) CreatePodStatusResult(ctx context.Context, v *corev1.PodStatusResult, opts metav1.CreateOptions) (*corev1.PodStatusResult, error) {
	result, err := c.backend.Create(ctx, "podstatusresults", "PodStatusResult", v, opts, &corev1.PodStatusResult{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PodStatusResult), nil
}

func (c *CoreV1) UpdatePodStatusResult(ctx context.Context, v *corev1.PodStatusResult, opts metav1.UpdateOptions) (*corev1.PodStatusResult, error) {
	result, err := c.backend.Update(ctx, "podstatusresults", "PodStatusResult", v, opts, &corev1.PodStatusResult{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PodStatusResult), nil
}

func (c *CoreV1) DeletePodStatusResult(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "podstatusresults"}, namespace, name, opts)
}

func (c *CoreV1) ListPodStatusResult(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.PodStatusResultList, error) {
	result, err := c.backend.List(ctx, "podstatusresults", "PodStatusResult", namespace, opts, &corev1.PodStatusResultList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PodStatusResultList), nil
}

func (c *CoreV1) WatchPodStatusResult(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "podstatusresults"}, namespace, opts)
}

func (c *CoreV1) GetPodTemplate(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.PodTemplate, error) {
	result, err := c.backend.Get(ctx, "podtemplates", "PodTemplate", namespace, name, opts, &corev1.PodTemplate{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PodTemplate), nil
}

func (c *CoreV1) CreatePodTemplate(ctx context.Context, v *corev1.PodTemplate, opts metav1.CreateOptions) (*corev1.PodTemplate, error) {
	result, err := c.backend.Create(ctx, "podtemplates", "PodTemplate", v, opts, &corev1.PodTemplate{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PodTemplate), nil
}

func (c *CoreV1) UpdatePodTemplate(ctx context.Context, v *corev1.PodTemplate, opts metav1.UpdateOptions) (*corev1.PodTemplate, error) {
	result, err := c.backend.Update(ctx, "podtemplates", "PodTemplate", v, opts, &corev1.PodTemplate{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PodTemplate), nil
}

func (c *CoreV1) DeletePodTemplate(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "podtemplates"}, namespace, name, opts)
}

func (c *CoreV1) ListPodTemplate(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.PodTemplateList, error) {
	result, err := c.backend.List(ctx, "podtemplates", "PodTemplate", namespace, opts, &corev1.PodTemplateList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.PodTemplateList), nil
}

func (c *CoreV1) WatchPodTemplate(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "podtemplates"}, namespace, opts)
}

func (c *CoreV1) GetRangeAllocation(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.RangeAllocation, error) {
	result, err := c.backend.Get(ctx, "rangeallocations", "RangeAllocation", namespace, name, opts, &corev1.RangeAllocation{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.RangeAllocation), nil
}

func (c *CoreV1) CreateRangeAllocation(ctx context.Context, v *corev1.RangeAllocation, opts metav1.CreateOptions) (*corev1.RangeAllocation, error) {
	result, err := c.backend.Create(ctx, "rangeallocations", "RangeAllocation", v, opts, &corev1.RangeAllocation{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.RangeAllocation), nil
}

func (c *CoreV1) UpdateRangeAllocation(ctx context.Context, v *corev1.RangeAllocation, opts metav1.UpdateOptions) (*corev1.RangeAllocation, error) {
	result, err := c.backend.Update(ctx, "rangeallocations", "RangeAllocation", v, opts, &corev1.RangeAllocation{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.RangeAllocation), nil
}

func (c *CoreV1) DeleteRangeAllocation(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "rangeallocations"}, namespace, name, opts)
}

func (c *CoreV1) ListRangeAllocation(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.RangeAllocationList, error) {
	result, err := c.backend.List(ctx, "rangeallocations", "RangeAllocation", namespace, opts, &corev1.RangeAllocationList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.RangeAllocationList), nil
}

func (c *CoreV1) WatchRangeAllocation(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "rangeallocations"}, namespace, opts)
}

func (c *CoreV1) GetReplicationController(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.ReplicationController, error) {
	result, err := c.backend.Get(ctx, "replicationcontrollers", "ReplicationController", namespace, name, opts, &corev1.ReplicationController{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ReplicationController), nil
}

func (c *CoreV1) CreateReplicationController(ctx context.Context, v *corev1.ReplicationController, opts metav1.CreateOptions) (*corev1.ReplicationController, error) {
	result, err := c.backend.Create(ctx, "replicationcontrollers", "ReplicationController", v, opts, &corev1.ReplicationController{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ReplicationController), nil
}

func (c *CoreV1) UpdateReplicationController(ctx context.Context, v *corev1.ReplicationController, opts metav1.UpdateOptions) (*corev1.ReplicationController, error) {
	result, err := c.backend.Update(ctx, "replicationcontrollers", "ReplicationController", v, opts, &corev1.ReplicationController{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ReplicationController), nil
}

func (c *CoreV1) DeleteReplicationController(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "replicationcontrollers"}, namespace, name, opts)
}

func (c *CoreV1) ListReplicationController(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.ReplicationControllerList, error) {
	result, err := c.backend.List(ctx, "replicationcontrollers", "ReplicationController", namespace, opts, &corev1.ReplicationControllerList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ReplicationControllerList), nil
}

func (c *CoreV1) WatchReplicationController(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "replicationcontrollers"}, namespace, opts)
}

func (c *CoreV1) GetResourceQuota(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.ResourceQuota, error) {
	result, err := c.backend.Get(ctx, "resourcequotas", "ResourceQuota", namespace, name, opts, &corev1.ResourceQuota{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ResourceQuota), nil
}

func (c *CoreV1) CreateResourceQuota(ctx context.Context, v *corev1.ResourceQuota, opts metav1.CreateOptions) (*corev1.ResourceQuota, error) {
	result, err := c.backend.Create(ctx, "resourcequotas", "ResourceQuota", v, opts, &corev1.ResourceQuota{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ResourceQuota), nil
}

func (c *CoreV1) UpdateResourceQuota(ctx context.Context, v *corev1.ResourceQuota, opts metav1.UpdateOptions) (*corev1.ResourceQuota, error) {
	result, err := c.backend.Update(ctx, "resourcequotas", "ResourceQuota", v, opts, &corev1.ResourceQuota{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ResourceQuota), nil
}

func (c *CoreV1) DeleteResourceQuota(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "resourcequotas"}, namespace, name, opts)
}

func (c *CoreV1) ListResourceQuota(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.ResourceQuotaList, error) {
	result, err := c.backend.List(ctx, "resourcequotas", "ResourceQuota", namespace, opts, &corev1.ResourceQuotaList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ResourceQuotaList), nil
}

func (c *CoreV1) WatchResourceQuota(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "resourcequotas"}, namespace, opts)
}

func (c *CoreV1) GetSecret(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Secret, error) {
	result, err := c.backend.Get(ctx, "secrets", "Secret", namespace, name, opts, &corev1.Secret{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Secret), nil
}

func (c *CoreV1) CreateSecret(ctx context.Context, v *corev1.Secret, opts metav1.CreateOptions) (*corev1.Secret, error) {
	result, err := c.backend.Create(ctx, "secrets", "Secret", v, opts, &corev1.Secret{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Secret), nil
}

func (c *CoreV1) UpdateSecret(ctx context.Context, v *corev1.Secret, opts metav1.UpdateOptions) (*corev1.Secret, error) {
	result, err := c.backend.Update(ctx, "secrets", "Secret", v, opts, &corev1.Secret{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Secret), nil
}

func (c *CoreV1) DeleteSecret(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "secrets"}, namespace, name, opts)
}

func (c *CoreV1) ListSecret(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.SecretList, error) {
	result, err := c.backend.List(ctx, "secrets", "Secret", namespace, opts, &corev1.SecretList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.SecretList), nil
}

func (c *CoreV1) WatchSecret(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "secrets"}, namespace, opts)
}

func (c *CoreV1) GetService(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Service, error) {
	result, err := c.backend.Get(ctx, "services", "Service", namespace, name, opts, &corev1.Service{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Service), nil
}

func (c *CoreV1) CreateService(ctx context.Context, v *corev1.Service, opts metav1.CreateOptions) (*corev1.Service, error) {
	result, err := c.backend.Create(ctx, "services", "Service", v, opts, &corev1.Service{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Service), nil
}

func (c *CoreV1) UpdateService(ctx context.Context, v *corev1.Service, opts metav1.UpdateOptions) (*corev1.Service, error) {
	result, err := c.backend.Update(ctx, "services", "Service", v, opts, &corev1.Service{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.Service), nil
}

func (c *CoreV1) DeleteService(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "services"}, namespace, name, opts)
}

func (c *CoreV1) ListService(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.ServiceList, error) {
	result, err := c.backend.List(ctx, "services", "Service", namespace, opts, &corev1.ServiceList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ServiceList), nil
}

func (c *CoreV1) WatchService(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "services"}, namespace, opts)
}

func (c *CoreV1) GetServiceAccount(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.ServiceAccount, error) {
	result, err := c.backend.Get(ctx, "serviceaccounts", "ServiceAccount", namespace, name, opts, &corev1.ServiceAccount{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ServiceAccount), nil
}

func (c *CoreV1) CreateServiceAccount(ctx context.Context, v *corev1.ServiceAccount, opts metav1.CreateOptions) (*corev1.ServiceAccount, error) {
	result, err := c.backend.Create(ctx, "serviceaccounts", "ServiceAccount", v, opts, &corev1.ServiceAccount{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ServiceAccount), nil
}

func (c *CoreV1) UpdateServiceAccount(ctx context.Context, v *corev1.ServiceAccount, opts metav1.UpdateOptions) (*corev1.ServiceAccount, error) {
	result, err := c.backend.Update(ctx, "serviceaccounts", "ServiceAccount", v, opts, &corev1.ServiceAccount{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ServiceAccount), nil
}

func (c *CoreV1) DeleteServiceAccount(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "serviceaccounts"}, namespace, name, opts)
}

func (c *CoreV1) ListServiceAccount(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.ServiceAccountList, error) {
	result, err := c.backend.List(ctx, "serviceaccounts", "ServiceAccount", namespace, opts, &corev1.ServiceAccountList{})
	if err != nil {
		return nil, err
	}
	return result.(*corev1.ServiceAccountList), nil
}

func (c *CoreV1) WatchServiceAccount(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".", Version: "v1", Resource: "serviceaccounts"}, namespace, opts)
}

type AdmissionregistrationK8sIoV1 struct {
	backend Backend
}

func NewAdmissionregistrationK8sIoV1Client(b Backend) *AdmissionregistrationK8sIoV1 {
	return &AdmissionregistrationK8sIoV1{backend: b}
}

func (c *AdmissionregistrationK8sIoV1) GetMutatingWebhookConfiguration(ctx context.Context, name string, opts metav1.GetOptions) (*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	result, err := c.backend.GetClusterScoped(ctx, "mutatingwebhookconfigurations", "MutatingWebhookConfiguration", name, opts, &admissionregistrationv1.MutatingWebhookConfiguration{})
	if err != nil {
		return nil, err
	}
	return result.(*admissionregistrationv1.MutatingWebhookConfiguration), nil
}

func (c *AdmissionregistrationK8sIoV1) CreateMutatingWebhookConfiguration(ctx context.Context, v *admissionregistrationv1.MutatingWebhookConfiguration, opts metav1.CreateOptions) (*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "mutatingwebhookconfigurations", "MutatingWebhookConfiguration", v, opts, &admissionregistrationv1.MutatingWebhookConfiguration{})
	if err != nil {
		return nil, err
	}
	return result.(*admissionregistrationv1.MutatingWebhookConfiguration), nil
}

func (c *AdmissionregistrationK8sIoV1) UpdateMutatingWebhookConfiguration(ctx context.Context, v *admissionregistrationv1.MutatingWebhookConfiguration, opts metav1.UpdateOptions) (*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "mutatingwebhookconfigurations", "MutatingWebhookConfiguration", v, opts, &admissionregistrationv1.MutatingWebhookConfiguration{})
	if err != nil {
		return nil, err
	}
	return result.(*admissionregistrationv1.MutatingWebhookConfiguration), nil
}

func (c *AdmissionregistrationK8sIoV1) DeleteMutatingWebhookConfiguration(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: ".admissionregistration.k8s.io", Version: "v1", Resource: "mutatingwebhookconfigurations"}, name, opts)
}

func (c *AdmissionregistrationK8sIoV1) ListMutatingWebhookConfiguration(ctx context.Context, opts metav1.ListOptions) (*admissionregistrationv1.MutatingWebhookConfigurationList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "mutatingwebhookconfigurations", "MutatingWebhookConfiguration", opts, &admissionregistrationv1.MutatingWebhookConfigurationList{})
	if err != nil {
		return nil, err
	}
	return result.(*admissionregistrationv1.MutatingWebhookConfigurationList), nil
}

func (c *AdmissionregistrationK8sIoV1) WatchMutatingWebhookConfiguration(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: ".admissionregistration.k8s.io", Version: "v1", Resource: "mutatingwebhookconfigurations"}, opts)
}

func (c *AdmissionregistrationK8sIoV1) GetValidatingWebhookConfiguration(ctx context.Context, name string, opts metav1.GetOptions) (*admissionregistrationv1.ValidatingWebhookConfiguration, error) {
	result, err := c.backend.GetClusterScoped(ctx, "validatingwebhookconfigurations", "ValidatingWebhookConfiguration", name, opts, &admissionregistrationv1.ValidatingWebhookConfiguration{})
	if err != nil {
		return nil, err
	}
	return result.(*admissionregistrationv1.ValidatingWebhookConfiguration), nil
}

func (c *AdmissionregistrationK8sIoV1) CreateValidatingWebhookConfiguration(ctx context.Context, v *admissionregistrationv1.ValidatingWebhookConfiguration, opts metav1.CreateOptions) (*admissionregistrationv1.ValidatingWebhookConfiguration, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "validatingwebhookconfigurations", "ValidatingWebhookConfiguration", v, opts, &admissionregistrationv1.ValidatingWebhookConfiguration{})
	if err != nil {
		return nil, err
	}
	return result.(*admissionregistrationv1.ValidatingWebhookConfiguration), nil
}

func (c *AdmissionregistrationK8sIoV1) UpdateValidatingWebhookConfiguration(ctx context.Context, v *admissionregistrationv1.ValidatingWebhookConfiguration, opts metav1.UpdateOptions) (*admissionregistrationv1.ValidatingWebhookConfiguration, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "validatingwebhookconfigurations", "ValidatingWebhookConfiguration", v, opts, &admissionregistrationv1.ValidatingWebhookConfiguration{})
	if err != nil {
		return nil, err
	}
	return result.(*admissionregistrationv1.ValidatingWebhookConfiguration), nil
}

func (c *AdmissionregistrationK8sIoV1) DeleteValidatingWebhookConfiguration(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: ".admissionregistration.k8s.io", Version: "v1", Resource: "validatingwebhookconfigurations"}, name, opts)
}

func (c *AdmissionregistrationK8sIoV1) ListValidatingWebhookConfiguration(ctx context.Context, opts metav1.ListOptions) (*admissionregistrationv1.ValidatingWebhookConfigurationList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "validatingwebhookconfigurations", "ValidatingWebhookConfiguration", opts, &admissionregistrationv1.ValidatingWebhookConfigurationList{})
	if err != nil {
		return nil, err
	}
	return result.(*admissionregistrationv1.ValidatingWebhookConfigurationList), nil
}

func (c *AdmissionregistrationK8sIoV1) WatchValidatingWebhookConfiguration(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: ".admissionregistration.k8s.io", Version: "v1", Resource: "validatingwebhookconfigurations"}, opts)
}

type AppsV1 struct {
	backend Backend
}

func NewAppsV1Client(b Backend) *AppsV1 {
	return &AppsV1{backend: b}
}

func (c *AppsV1) GetControllerRevision(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*appsv1.ControllerRevision, error) {
	result, err := c.backend.Get(ctx, "controllerrevisions", "ControllerRevision", namespace, name, opts, &appsv1.ControllerRevision{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.ControllerRevision), nil
}

func (c *AppsV1) CreateControllerRevision(ctx context.Context, v *appsv1.ControllerRevision, opts metav1.CreateOptions) (*appsv1.ControllerRevision, error) {
	result, err := c.backend.Create(ctx, "controllerrevisions", "ControllerRevision", v, opts, &appsv1.ControllerRevision{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.ControllerRevision), nil
}

func (c *AppsV1) UpdateControllerRevision(ctx context.Context, v *appsv1.ControllerRevision, opts metav1.UpdateOptions) (*appsv1.ControllerRevision, error) {
	result, err := c.backend.Update(ctx, "controllerrevisions", "ControllerRevision", v, opts, &appsv1.ControllerRevision{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.ControllerRevision), nil
}

func (c *AppsV1) DeleteControllerRevision(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".apps", Version: "v1", Resource: "controllerrevisions"}, namespace, name, opts)
}

func (c *AppsV1) ListControllerRevision(ctx context.Context, namespace string, opts metav1.ListOptions) (*appsv1.ControllerRevisionList, error) {
	result, err := c.backend.List(ctx, "controllerrevisions", "ControllerRevision", namespace, opts, &appsv1.ControllerRevisionList{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.ControllerRevisionList), nil
}

func (c *AppsV1) WatchControllerRevision(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".apps", Version: "v1", Resource: "controllerrevisions"}, namespace, opts)
}

func (c *AppsV1) GetCronJob(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*batchv1.CronJob, error) {
	result, err := c.backend.Get(ctx, "cronjobs", "CronJob", namespace, name, opts, &batchv1.CronJob{})
	if err != nil {
		return nil, err
	}
	return result.(*batchv1.CronJob), nil
}

func (c *AppsV1) CreateCronJob(ctx context.Context, v *batchv1.CronJob, opts metav1.CreateOptions) (*batchv1.CronJob, error) {
	result, err := c.backend.Create(ctx, "cronjobs", "CronJob", v, opts, &batchv1.CronJob{})
	if err != nil {
		return nil, err
	}
	return result.(*batchv1.CronJob), nil
}

func (c *AppsV1) UpdateCronJob(ctx context.Context, v *batchv1.CronJob, opts metav1.UpdateOptions) (*batchv1.CronJob, error) {
	result, err := c.backend.Update(ctx, "cronjobs", "CronJob", v, opts, &batchv1.CronJob{})
	if err != nil {
		return nil, err
	}
	return result.(*batchv1.CronJob), nil
}

func (c *AppsV1) DeleteCronJob(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".apps", Version: "v1", Resource: "cronjobs"}, namespace, name, opts)
}

func (c *AppsV1) ListCronJob(ctx context.Context, namespace string, opts metav1.ListOptions) (*batchv1.CronJobList, error) {
	result, err := c.backend.List(ctx, "cronjobs", "CronJob", namespace, opts, &batchv1.CronJobList{})
	if err != nil {
		return nil, err
	}
	return result.(*batchv1.CronJobList), nil
}

func (c *AppsV1) WatchCronJob(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".apps", Version: "v1", Resource: "cronjobs"}, namespace, opts)
}

func (c *AppsV1) GetDaemonSet(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*appsv1.DaemonSet, error) {
	result, err := c.backend.Get(ctx, "daemonsets", "DaemonSet", namespace, name, opts, &appsv1.DaemonSet{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.DaemonSet), nil
}

func (c *AppsV1) CreateDaemonSet(ctx context.Context, v *appsv1.DaemonSet, opts metav1.CreateOptions) (*appsv1.DaemonSet, error) {
	result, err := c.backend.Create(ctx, "daemonsets", "DaemonSet", v, opts, &appsv1.DaemonSet{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.DaemonSet), nil
}

func (c *AppsV1) UpdateDaemonSet(ctx context.Context, v *appsv1.DaemonSet, opts metav1.UpdateOptions) (*appsv1.DaemonSet, error) {
	result, err := c.backend.Update(ctx, "daemonsets", "DaemonSet", v, opts, &appsv1.DaemonSet{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.DaemonSet), nil
}

func (c *AppsV1) DeleteDaemonSet(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".apps", Version: "v1", Resource: "daemonsets"}, namespace, name, opts)
}

func (c *AppsV1) ListDaemonSet(ctx context.Context, namespace string, opts metav1.ListOptions) (*appsv1.DaemonSetList, error) {
	result, err := c.backend.List(ctx, "daemonsets", "DaemonSet", namespace, opts, &appsv1.DaemonSetList{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.DaemonSetList), nil
}

func (c *AppsV1) WatchDaemonSet(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".apps", Version: "v1", Resource: "daemonsets"}, namespace, opts)
}

func (c *AppsV1) GetDeployment(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*appsv1.Deployment, error) {
	result, err := c.backend.Get(ctx, "deployments", "Deployment", namespace, name, opts, &appsv1.Deployment{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.Deployment), nil
}

func (c *AppsV1) CreateDeployment(ctx context.Context, v *appsv1.Deployment, opts metav1.CreateOptions) (*appsv1.Deployment, error) {
	result, err := c.backend.Create(ctx, "deployments", "Deployment", v, opts, &appsv1.Deployment{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.Deployment), nil
}

func (c *AppsV1) UpdateDeployment(ctx context.Context, v *appsv1.Deployment, opts metav1.UpdateOptions) (*appsv1.Deployment, error) {
	result, err := c.backend.Update(ctx, "deployments", "Deployment", v, opts, &appsv1.Deployment{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.Deployment), nil
}

func (c *AppsV1) DeleteDeployment(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".apps", Version: "v1", Resource: "deployments"}, namespace, name, opts)
}

func (c *AppsV1) ListDeployment(ctx context.Context, namespace string, opts metav1.ListOptions) (*appsv1.DeploymentList, error) {
	result, err := c.backend.List(ctx, "deployments", "Deployment", namespace, opts, &appsv1.DeploymentList{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.DeploymentList), nil
}

func (c *AppsV1) WatchDeployment(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".apps", Version: "v1", Resource: "deployments"}, namespace, opts)
}

func (c *AppsV1) GetJob(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*batchv1.Job, error) {
	result, err := c.backend.Get(ctx, "jobs", "Job", namespace, name, opts, &batchv1.Job{})
	if err != nil {
		return nil, err
	}
	return result.(*batchv1.Job), nil
}

func (c *AppsV1) CreateJob(ctx context.Context, v *batchv1.Job, opts metav1.CreateOptions) (*batchv1.Job, error) {
	result, err := c.backend.Create(ctx, "jobs", "Job", v, opts, &batchv1.Job{})
	if err != nil {
		return nil, err
	}
	return result.(*batchv1.Job), nil
}

func (c *AppsV1) UpdateJob(ctx context.Context, v *batchv1.Job, opts metav1.UpdateOptions) (*batchv1.Job, error) {
	result, err := c.backend.Update(ctx, "jobs", "Job", v, opts, &batchv1.Job{})
	if err != nil {
		return nil, err
	}
	return result.(*batchv1.Job), nil
}

func (c *AppsV1) DeleteJob(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".apps", Version: "v1", Resource: "jobs"}, namespace, name, opts)
}

func (c *AppsV1) ListJob(ctx context.Context, namespace string, opts metav1.ListOptions) (*batchv1.JobList, error) {
	result, err := c.backend.List(ctx, "jobs", "Job", namespace, opts, &batchv1.JobList{})
	if err != nil {
		return nil, err
	}
	return result.(*batchv1.JobList), nil
}

func (c *AppsV1) WatchJob(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".apps", Version: "v1", Resource: "jobs"}, namespace, opts)
}

func (c *AppsV1) GetReplicaSet(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*appsv1.ReplicaSet, error) {
	result, err := c.backend.Get(ctx, "replicasets", "ReplicaSet", namespace, name, opts, &appsv1.ReplicaSet{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.ReplicaSet), nil
}

func (c *AppsV1) CreateReplicaSet(ctx context.Context, v *appsv1.ReplicaSet, opts metav1.CreateOptions) (*appsv1.ReplicaSet, error) {
	result, err := c.backend.Create(ctx, "replicasets", "ReplicaSet", v, opts, &appsv1.ReplicaSet{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.ReplicaSet), nil
}

func (c *AppsV1) UpdateReplicaSet(ctx context.Context, v *appsv1.ReplicaSet, opts metav1.UpdateOptions) (*appsv1.ReplicaSet, error) {
	result, err := c.backend.Update(ctx, "replicasets", "ReplicaSet", v, opts, &appsv1.ReplicaSet{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.ReplicaSet), nil
}

func (c *AppsV1) DeleteReplicaSet(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".apps", Version: "v1", Resource: "replicasets"}, namespace, name, opts)
}

func (c *AppsV1) ListReplicaSet(ctx context.Context, namespace string, opts metav1.ListOptions) (*appsv1.ReplicaSetList, error) {
	result, err := c.backend.List(ctx, "replicasets", "ReplicaSet", namespace, opts, &appsv1.ReplicaSetList{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.ReplicaSetList), nil
}

func (c *AppsV1) WatchReplicaSet(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".apps", Version: "v1", Resource: "replicasets"}, namespace, opts)
}

func (c *AppsV1) GetStatefulSet(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*appsv1.StatefulSet, error) {
	result, err := c.backend.Get(ctx, "statefulsets", "StatefulSet", namespace, name, opts, &appsv1.StatefulSet{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.StatefulSet), nil
}

func (c *AppsV1) CreateStatefulSet(ctx context.Context, v *appsv1.StatefulSet, opts metav1.CreateOptions) (*appsv1.StatefulSet, error) {
	result, err := c.backend.Create(ctx, "statefulsets", "StatefulSet", v, opts, &appsv1.StatefulSet{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.StatefulSet), nil
}

func (c *AppsV1) UpdateStatefulSet(ctx context.Context, v *appsv1.StatefulSet, opts metav1.UpdateOptions) (*appsv1.StatefulSet, error) {
	result, err := c.backend.Update(ctx, "statefulsets", "StatefulSet", v, opts, &appsv1.StatefulSet{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.StatefulSet), nil
}

func (c *AppsV1) DeleteStatefulSet(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".apps", Version: "v1", Resource: "statefulsets"}, namespace, name, opts)
}

func (c *AppsV1) ListStatefulSet(ctx context.Context, namespace string, opts metav1.ListOptions) (*appsv1.StatefulSetList, error) {
	result, err := c.backend.List(ctx, "statefulsets", "StatefulSet", namespace, opts, &appsv1.StatefulSetList{})
	if err != nil {
		return nil, err
	}
	return result.(*appsv1.StatefulSetList), nil
}

func (c *AppsV1) WatchStatefulSet(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".apps", Version: "v1", Resource: "statefulsets"}, namespace, opts)
}

type AuthenticationK8sIoV1 struct {
	backend Backend
}

func NewAuthenticationK8sIoV1Client(b Backend) *AuthenticationK8sIoV1 {
	return &AuthenticationK8sIoV1{backend: b}
}

func (c *AuthenticationK8sIoV1) GetTokenRequest(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*authenticationv1.TokenRequest, error) {
	result, err := c.backend.Get(ctx, "tokenrequests", "TokenRequest", namespace, name, opts, &authenticationv1.TokenRequest{})
	if err != nil {
		return nil, err
	}
	return result.(*authenticationv1.TokenRequest), nil
}

func (c *AuthenticationK8sIoV1) CreateTokenRequest(ctx context.Context, v *authenticationv1.TokenRequest, opts metav1.CreateOptions) (*authenticationv1.TokenRequest, error) {
	result, err := c.backend.Create(ctx, "tokenrequests", "TokenRequest", v, opts, &authenticationv1.TokenRequest{})
	if err != nil {
		return nil, err
	}
	return result.(*authenticationv1.TokenRequest), nil
}

func (c *AuthenticationK8sIoV1) UpdateTokenRequest(ctx context.Context, v *authenticationv1.TokenRequest, opts metav1.UpdateOptions) (*authenticationv1.TokenRequest, error) {
	result, err := c.backend.Update(ctx, "tokenrequests", "TokenRequest", v, opts, &authenticationv1.TokenRequest{})
	if err != nil {
		return nil, err
	}
	return result.(*authenticationv1.TokenRequest), nil
}

func (c *AuthenticationK8sIoV1) DeleteTokenRequest(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".authentication.k8s.io", Version: "v1", Resource: "tokenrequests"}, namespace, name, opts)
}

func (c *AuthenticationK8sIoV1) ListTokenRequest(ctx context.Context, namespace string, opts metav1.ListOptions) (*authenticationv1.TokenRequestList, error) {
	result, err := c.backend.List(ctx, "tokenrequests", "TokenRequest", namespace, opts, &authenticationv1.TokenRequestList{})
	if err != nil {
		return nil, err
	}
	return result.(*authenticationv1.TokenRequestList), nil
}

func (c *AuthenticationK8sIoV1) WatchTokenRequest(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".authentication.k8s.io", Version: "v1", Resource: "tokenrequests"}, namespace, opts)
}

func (c *AuthenticationK8sIoV1) GetTokenReview(ctx context.Context, name string, opts metav1.GetOptions) (*authenticationv1.TokenReview, error) {
	result, err := c.backend.GetClusterScoped(ctx, "tokenreviews", "TokenReview", name, opts, &authenticationv1.TokenReview{})
	if err != nil {
		return nil, err
	}
	return result.(*authenticationv1.TokenReview), nil
}

func (c *AuthenticationK8sIoV1) CreateTokenReview(ctx context.Context, v *authenticationv1.TokenReview, opts metav1.CreateOptions) (*authenticationv1.TokenReview, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "tokenreviews", "TokenReview", v, opts, &authenticationv1.TokenReview{})
	if err != nil {
		return nil, err
	}
	return result.(*authenticationv1.TokenReview), nil
}

func (c *AuthenticationK8sIoV1) UpdateTokenReview(ctx context.Context, v *authenticationv1.TokenReview, opts metav1.UpdateOptions) (*authenticationv1.TokenReview, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "tokenreviews", "TokenReview", v, opts, &authenticationv1.TokenReview{})
	if err != nil {
		return nil, err
	}
	return result.(*authenticationv1.TokenReview), nil
}

func (c *AuthenticationK8sIoV1) DeleteTokenReview(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: ".authentication.k8s.io", Version: "v1", Resource: "tokenreviews"}, name, opts)
}

func (c *AuthenticationK8sIoV1) ListTokenReview(ctx context.Context, opts metav1.ListOptions) (*authenticationv1.TokenReviewList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "tokenreviews", "TokenReview", opts, &authenticationv1.TokenReviewList{})
	if err != nil {
		return nil, err
	}
	return result.(*authenticationv1.TokenReviewList), nil
}

func (c *AuthenticationK8sIoV1) WatchTokenReview(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: ".authentication.k8s.io", Version: "v1", Resource: "tokenreviews"}, opts)
}

type AuthorizationK8sIoV1 struct {
	backend Backend
}

func NewAuthorizationK8sIoV1Client(b Backend) *AuthorizationK8sIoV1 {
	return &AuthorizationK8sIoV1{backend: b}
}

func (c *AuthorizationK8sIoV1) GetLocalSubjectAccessReview(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*authorizationv1.LocalSubjectAccessReview, error) {
	result, err := c.backend.Get(ctx, "localsubjectaccessreviews", "LocalSubjectAccessReview", namespace, name, opts, &authorizationv1.LocalSubjectAccessReview{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.LocalSubjectAccessReview), nil
}

func (c *AuthorizationK8sIoV1) CreateLocalSubjectAccessReview(ctx context.Context, v *authorizationv1.LocalSubjectAccessReview, opts metav1.CreateOptions) (*authorizationv1.LocalSubjectAccessReview, error) {
	result, err := c.backend.Create(ctx, "localsubjectaccessreviews", "LocalSubjectAccessReview", v, opts, &authorizationv1.LocalSubjectAccessReview{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.LocalSubjectAccessReview), nil
}

func (c *AuthorizationK8sIoV1) UpdateLocalSubjectAccessReview(ctx context.Context, v *authorizationv1.LocalSubjectAccessReview, opts metav1.UpdateOptions) (*authorizationv1.LocalSubjectAccessReview, error) {
	result, err := c.backend.Update(ctx, "localsubjectaccessreviews", "LocalSubjectAccessReview", v, opts, &authorizationv1.LocalSubjectAccessReview{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.LocalSubjectAccessReview), nil
}

func (c *AuthorizationK8sIoV1) DeleteLocalSubjectAccessReview(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".authorization.k8s.io", Version: "v1", Resource: "localsubjectaccessreviews"}, namespace, name, opts)
}

func (c *AuthorizationK8sIoV1) ListLocalSubjectAccessReview(ctx context.Context, namespace string, opts metav1.ListOptions) (*authorizationv1.LocalSubjectAccessReviewList, error) {
	result, err := c.backend.List(ctx, "localsubjectaccessreviews", "LocalSubjectAccessReview", namespace, opts, &authorizationv1.LocalSubjectAccessReviewList{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.LocalSubjectAccessReviewList), nil
}

func (c *AuthorizationK8sIoV1) WatchLocalSubjectAccessReview(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".authorization.k8s.io", Version: "v1", Resource: "localsubjectaccessreviews"}, namespace, opts)
}

func (c *AuthorizationK8sIoV1) GetSelfSubjectAccessReview(ctx context.Context, name string, opts metav1.GetOptions) (*authorizationv1.SelfSubjectAccessReview, error) {
	result, err := c.backend.GetClusterScoped(ctx, "selfsubjectaccessreviews", "SelfSubjectAccessReview", name, opts, &authorizationv1.SelfSubjectAccessReview{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.SelfSubjectAccessReview), nil
}

func (c *AuthorizationK8sIoV1) CreateSelfSubjectAccessReview(ctx context.Context, v *authorizationv1.SelfSubjectAccessReview, opts metav1.CreateOptions) (*authorizationv1.SelfSubjectAccessReview, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "selfsubjectaccessreviews", "SelfSubjectAccessReview", v, opts, &authorizationv1.SelfSubjectAccessReview{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.SelfSubjectAccessReview), nil
}

func (c *AuthorizationK8sIoV1) UpdateSelfSubjectAccessReview(ctx context.Context, v *authorizationv1.SelfSubjectAccessReview, opts metav1.UpdateOptions) (*authorizationv1.SelfSubjectAccessReview, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "selfsubjectaccessreviews", "SelfSubjectAccessReview", v, opts, &authorizationv1.SelfSubjectAccessReview{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.SelfSubjectAccessReview), nil
}

func (c *AuthorizationK8sIoV1) DeleteSelfSubjectAccessReview(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: ".authorization.k8s.io", Version: "v1", Resource: "selfsubjectaccessreviews"}, name, opts)
}

func (c *AuthorizationK8sIoV1) ListSelfSubjectAccessReview(ctx context.Context, opts metav1.ListOptions) (*authorizationv1.SelfSubjectAccessReviewList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "selfsubjectaccessreviews", "SelfSubjectAccessReview", opts, &authorizationv1.SelfSubjectAccessReviewList{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.SelfSubjectAccessReviewList), nil
}

func (c *AuthorizationK8sIoV1) WatchSelfSubjectAccessReview(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: ".authorization.k8s.io", Version: "v1", Resource: "selfsubjectaccessreviews"}, opts)
}

func (c *AuthorizationK8sIoV1) GetSelfSubjectRulesReview(ctx context.Context, name string, opts metav1.GetOptions) (*authorizationv1.SelfSubjectRulesReview, error) {
	result, err := c.backend.GetClusterScoped(ctx, "selfsubjectrulesreviews", "SelfSubjectRulesReview", name, opts, &authorizationv1.SelfSubjectRulesReview{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.SelfSubjectRulesReview), nil
}

func (c *AuthorizationK8sIoV1) CreateSelfSubjectRulesReview(ctx context.Context, v *authorizationv1.SelfSubjectRulesReview, opts metav1.CreateOptions) (*authorizationv1.SelfSubjectRulesReview, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "selfsubjectrulesreviews", "SelfSubjectRulesReview", v, opts, &authorizationv1.SelfSubjectRulesReview{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.SelfSubjectRulesReview), nil
}

func (c *AuthorizationK8sIoV1) UpdateSelfSubjectRulesReview(ctx context.Context, v *authorizationv1.SelfSubjectRulesReview, opts metav1.UpdateOptions) (*authorizationv1.SelfSubjectRulesReview, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "selfsubjectrulesreviews", "SelfSubjectRulesReview", v, opts, &authorizationv1.SelfSubjectRulesReview{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.SelfSubjectRulesReview), nil
}

func (c *AuthorizationK8sIoV1) DeleteSelfSubjectRulesReview(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: ".authorization.k8s.io", Version: "v1", Resource: "selfsubjectrulesreviews"}, name, opts)
}

func (c *AuthorizationK8sIoV1) ListSelfSubjectRulesReview(ctx context.Context, opts metav1.ListOptions) (*authorizationv1.SelfSubjectRulesReviewList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "selfsubjectrulesreviews", "SelfSubjectRulesReview", opts, &authorizationv1.SelfSubjectRulesReviewList{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.SelfSubjectRulesReviewList), nil
}

func (c *AuthorizationK8sIoV1) WatchSelfSubjectRulesReview(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: ".authorization.k8s.io", Version: "v1", Resource: "selfsubjectrulesreviews"}, opts)
}

func (c *AuthorizationK8sIoV1) GetSubjectAccessReview(ctx context.Context, name string, opts metav1.GetOptions) (*authorizationv1.SubjectAccessReview, error) {
	result, err := c.backend.GetClusterScoped(ctx, "subjectaccessreviews", "SubjectAccessReview", name, opts, &authorizationv1.SubjectAccessReview{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.SubjectAccessReview), nil
}

func (c *AuthorizationK8sIoV1) CreateSubjectAccessReview(ctx context.Context, v *authorizationv1.SubjectAccessReview, opts metav1.CreateOptions) (*authorizationv1.SubjectAccessReview, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "subjectaccessreviews", "SubjectAccessReview", v, opts, &authorizationv1.SubjectAccessReview{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.SubjectAccessReview), nil
}

func (c *AuthorizationK8sIoV1) UpdateSubjectAccessReview(ctx context.Context, v *authorizationv1.SubjectAccessReview, opts metav1.UpdateOptions) (*authorizationv1.SubjectAccessReview, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "subjectaccessreviews", "SubjectAccessReview", v, opts, &authorizationv1.SubjectAccessReview{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.SubjectAccessReview), nil
}

func (c *AuthorizationK8sIoV1) DeleteSubjectAccessReview(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: ".authorization.k8s.io", Version: "v1", Resource: "subjectaccessreviews"}, name, opts)
}

func (c *AuthorizationK8sIoV1) ListSubjectAccessReview(ctx context.Context, opts metav1.ListOptions) (*authorizationv1.SubjectAccessReviewList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "subjectaccessreviews", "SubjectAccessReview", opts, &authorizationv1.SubjectAccessReviewList{})
	if err != nil {
		return nil, err
	}
	return result.(*authorizationv1.SubjectAccessReviewList), nil
}

func (c *AuthorizationK8sIoV1) WatchSubjectAccessReview(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: ".authorization.k8s.io", Version: "v1", Resource: "subjectaccessreviews"}, opts)
}

type AutoscalingV1 struct {
	backend Backend
}

func NewAutoscalingV1Client(b Backend) *AutoscalingV1 {
	return &AutoscalingV1{backend: b}
}

func (c *AutoscalingV1) GetHorizontalPodAutoscaler(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*autoscalingv1.HorizontalPodAutoscaler, error) {
	result, err := c.backend.Get(ctx, "horizontalpodautoscalers", "HorizontalPodAutoscaler", namespace, name, opts, &autoscalingv1.HorizontalPodAutoscaler{})
	if err != nil {
		return nil, err
	}
	return result.(*autoscalingv1.HorizontalPodAutoscaler), nil
}

func (c *AutoscalingV1) CreateHorizontalPodAutoscaler(ctx context.Context, v *autoscalingv1.HorizontalPodAutoscaler, opts metav1.CreateOptions) (*autoscalingv1.HorizontalPodAutoscaler, error) {
	result, err := c.backend.Create(ctx, "horizontalpodautoscalers", "HorizontalPodAutoscaler", v, opts, &autoscalingv1.HorizontalPodAutoscaler{})
	if err != nil {
		return nil, err
	}
	return result.(*autoscalingv1.HorizontalPodAutoscaler), nil
}

func (c *AutoscalingV1) UpdateHorizontalPodAutoscaler(ctx context.Context, v *autoscalingv1.HorizontalPodAutoscaler, opts metav1.UpdateOptions) (*autoscalingv1.HorizontalPodAutoscaler, error) {
	result, err := c.backend.Update(ctx, "horizontalpodautoscalers", "HorizontalPodAutoscaler", v, opts, &autoscalingv1.HorizontalPodAutoscaler{})
	if err != nil {
		return nil, err
	}
	return result.(*autoscalingv1.HorizontalPodAutoscaler), nil
}

func (c *AutoscalingV1) DeleteHorizontalPodAutoscaler(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".autoscaling", Version: "v1", Resource: "horizontalpodautoscalers"}, namespace, name, opts)
}

func (c *AutoscalingV1) ListHorizontalPodAutoscaler(ctx context.Context, namespace string, opts metav1.ListOptions) (*autoscalingv1.HorizontalPodAutoscalerList, error) {
	result, err := c.backend.List(ctx, "horizontalpodautoscalers", "HorizontalPodAutoscaler", namespace, opts, &autoscalingv1.HorizontalPodAutoscalerList{})
	if err != nil {
		return nil, err
	}
	return result.(*autoscalingv1.HorizontalPodAutoscalerList), nil
}

func (c *AutoscalingV1) WatchHorizontalPodAutoscaler(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".autoscaling", Version: "v1", Resource: "horizontalpodautoscalers"}, namespace, opts)
}

func (c *AutoscalingV1) GetScale(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*autoscalingv1.Scale, error) {
	result, err := c.backend.Get(ctx, "scales", "Scale", namespace, name, opts, &autoscalingv1.Scale{})
	if err != nil {
		return nil, err
	}
	return result.(*autoscalingv1.Scale), nil
}

func (c *AutoscalingV1) CreateScale(ctx context.Context, v *autoscalingv1.Scale, opts metav1.CreateOptions) (*autoscalingv1.Scale, error) {
	result, err := c.backend.Create(ctx, "scales", "Scale", v, opts, &autoscalingv1.Scale{})
	if err != nil {
		return nil, err
	}
	return result.(*autoscalingv1.Scale), nil
}

func (c *AutoscalingV1) UpdateScale(ctx context.Context, v *autoscalingv1.Scale, opts metav1.UpdateOptions) (*autoscalingv1.Scale, error) {
	result, err := c.backend.Update(ctx, "scales", "Scale", v, opts, &autoscalingv1.Scale{})
	if err != nil {
		return nil, err
	}
	return result.(*autoscalingv1.Scale), nil
}

func (c *AutoscalingV1) DeleteScale(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".autoscaling", Version: "v1", Resource: "scales"}, namespace, name, opts)
}

func (c *AutoscalingV1) ListScale(ctx context.Context, namespace string, opts metav1.ListOptions) (*autoscalingv1.ScaleList, error) {
	result, err := c.backend.List(ctx, "scales", "Scale", namespace, opts, &autoscalingv1.ScaleList{})
	if err != nil {
		return nil, err
	}
	return result.(*autoscalingv1.ScaleList), nil
}

func (c *AutoscalingV1) WatchScale(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".autoscaling", Version: "v1", Resource: "scales"}, namespace, opts)
}

type CertificatesK8sIoV1 struct {
	backend Backend
}

func NewCertificatesK8sIoV1Client(b Backend) *CertificatesK8sIoV1 {
	return &CertificatesK8sIoV1{backend: b}
}

func (c *CertificatesK8sIoV1) GetCertificateSigningRequest(ctx context.Context, name string, opts metav1.GetOptions) (*certificatesv1.CertificateSigningRequest, error) {
	result, err := c.backend.GetClusterScoped(ctx, "certificatesigningrequests", "CertificateSigningRequest", name, opts, &certificatesv1.CertificateSigningRequest{})
	if err != nil {
		return nil, err
	}
	return result.(*certificatesv1.CertificateSigningRequest), nil
}

func (c *CertificatesK8sIoV1) CreateCertificateSigningRequest(ctx context.Context, v *certificatesv1.CertificateSigningRequest, opts metav1.CreateOptions) (*certificatesv1.CertificateSigningRequest, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "certificatesigningrequests", "CertificateSigningRequest", v, opts, &certificatesv1.CertificateSigningRequest{})
	if err != nil {
		return nil, err
	}
	return result.(*certificatesv1.CertificateSigningRequest), nil
}

func (c *CertificatesK8sIoV1) UpdateCertificateSigningRequest(ctx context.Context, v *certificatesv1.CertificateSigningRequest, opts metav1.UpdateOptions) (*certificatesv1.CertificateSigningRequest, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "certificatesigningrequests", "CertificateSigningRequest", v, opts, &certificatesv1.CertificateSigningRequest{})
	if err != nil {
		return nil, err
	}
	return result.(*certificatesv1.CertificateSigningRequest), nil
}

func (c *CertificatesK8sIoV1) DeleteCertificateSigningRequest(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: ".certificates.k8s.io", Version: "v1", Resource: "certificatesigningrequests"}, name, opts)
}

func (c *CertificatesK8sIoV1) ListCertificateSigningRequest(ctx context.Context, opts metav1.ListOptions) (*certificatesv1.CertificateSigningRequestList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "certificatesigningrequests", "CertificateSigningRequest", opts, &certificatesv1.CertificateSigningRequestList{})
	if err != nil {
		return nil, err
	}
	return result.(*certificatesv1.CertificateSigningRequestList), nil
}

func (c *CertificatesK8sIoV1) WatchCertificateSigningRequest(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: ".certificates.k8s.io", Version: "v1", Resource: "certificatesigningrequests"}, opts)
}

type DiscoveryK8sIoV1 struct {
	backend Backend
}

func NewDiscoveryK8sIoV1Client(b Backend) *DiscoveryK8sIoV1 {
	return &DiscoveryK8sIoV1{backend: b}
}

func (c *DiscoveryK8sIoV1) GetEndpointSlice(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*discoveryv1.EndpointSlice, error) {
	result, err := c.backend.Get(ctx, "endpointslices", "EndpointSlice", namespace, name, opts, &discoveryv1.EndpointSlice{})
	if err != nil {
		return nil, err
	}
	return result.(*discoveryv1.EndpointSlice), nil
}

func (c *DiscoveryK8sIoV1) CreateEndpointSlice(ctx context.Context, v *discoveryv1.EndpointSlice, opts metav1.CreateOptions) (*discoveryv1.EndpointSlice, error) {
	result, err := c.backend.Create(ctx, "endpointslices", "EndpointSlice", v, opts, &discoveryv1.EndpointSlice{})
	if err != nil {
		return nil, err
	}
	return result.(*discoveryv1.EndpointSlice), nil
}

func (c *DiscoveryK8sIoV1) UpdateEndpointSlice(ctx context.Context, v *discoveryv1.EndpointSlice, opts metav1.UpdateOptions) (*discoveryv1.EndpointSlice, error) {
	result, err := c.backend.Update(ctx, "endpointslices", "EndpointSlice", v, opts, &discoveryv1.EndpointSlice{})
	if err != nil {
		return nil, err
	}
	return result.(*discoveryv1.EndpointSlice), nil
}

func (c *DiscoveryK8sIoV1) DeleteEndpointSlice(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".discovery.k8s.io", Version: "v1", Resource: "endpointslices"}, namespace, name, opts)
}

func (c *DiscoveryK8sIoV1) ListEndpointSlice(ctx context.Context, namespace string, opts metav1.ListOptions) (*discoveryv1.EndpointSliceList, error) {
	result, err := c.backend.List(ctx, "endpointslices", "EndpointSlice", namespace, opts, &discoveryv1.EndpointSliceList{})
	if err != nil {
		return nil, err
	}
	return result.(*discoveryv1.EndpointSliceList), nil
}

func (c *DiscoveryK8sIoV1) WatchEndpointSlice(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".discovery.k8s.io", Version: "v1", Resource: "endpointslices"}, namespace, opts)
}

type NetworkingK8sIoV1 struct {
	backend Backend
}

func NewNetworkingK8sIoV1Client(b Backend) *NetworkingK8sIoV1 {
	return &NetworkingK8sIoV1{backend: b}
}

func (c *NetworkingK8sIoV1) GetIngress(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*networkingv1.Ingress, error) {
	result, err := c.backend.Get(ctx, "ingresses", "Ingress", namespace, name, opts, &networkingv1.Ingress{})
	if err != nil {
		return nil, err
	}
	return result.(*networkingv1.Ingress), nil
}

func (c *NetworkingK8sIoV1) CreateIngress(ctx context.Context, v *networkingv1.Ingress, opts metav1.CreateOptions) (*networkingv1.Ingress, error) {
	result, err := c.backend.Create(ctx, "ingresses", "Ingress", v, opts, &networkingv1.Ingress{})
	if err != nil {
		return nil, err
	}
	return result.(*networkingv1.Ingress), nil
}

func (c *NetworkingK8sIoV1) UpdateIngress(ctx context.Context, v *networkingv1.Ingress, opts metav1.UpdateOptions) (*networkingv1.Ingress, error) {
	result, err := c.backend.Update(ctx, "ingresses", "Ingress", v, opts, &networkingv1.Ingress{})
	if err != nil {
		return nil, err
	}
	return result.(*networkingv1.Ingress), nil
}

func (c *NetworkingK8sIoV1) DeleteIngress(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".networking.k8s.io", Version: "v1", Resource: "ingresses"}, namespace, name, opts)
}

func (c *NetworkingK8sIoV1) ListIngress(ctx context.Context, namespace string, opts metav1.ListOptions) (*networkingv1.IngressList, error) {
	result, err := c.backend.List(ctx, "ingresses", "Ingress", namespace, opts, &networkingv1.IngressList{})
	if err != nil {
		return nil, err
	}
	return result.(*networkingv1.IngressList), nil
}

func (c *NetworkingK8sIoV1) WatchIngress(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".networking.k8s.io", Version: "v1", Resource: "ingresses"}, namespace, opts)
}

func (c *NetworkingK8sIoV1) GetIngressClass(ctx context.Context, name string, opts metav1.GetOptions) (*networkingv1.IngressClass, error) {
	result, err := c.backend.GetClusterScoped(ctx, "ingressclasses", "IngressClass", name, opts, &networkingv1.IngressClass{})
	if err != nil {
		return nil, err
	}
	return result.(*networkingv1.IngressClass), nil
}

func (c *NetworkingK8sIoV1) CreateIngressClass(ctx context.Context, v *networkingv1.IngressClass, opts metav1.CreateOptions) (*networkingv1.IngressClass, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "ingressclasses", "IngressClass", v, opts, &networkingv1.IngressClass{})
	if err != nil {
		return nil, err
	}
	return result.(*networkingv1.IngressClass), nil
}

func (c *NetworkingK8sIoV1) UpdateIngressClass(ctx context.Context, v *networkingv1.IngressClass, opts metav1.UpdateOptions) (*networkingv1.IngressClass, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "ingressclasses", "IngressClass", v, opts, &networkingv1.IngressClass{})
	if err != nil {
		return nil, err
	}
	return result.(*networkingv1.IngressClass), nil
}

func (c *NetworkingK8sIoV1) DeleteIngressClass(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: ".networking.k8s.io", Version: "v1", Resource: "ingressclasses"}, name, opts)
}

func (c *NetworkingK8sIoV1) ListIngressClass(ctx context.Context, opts metav1.ListOptions) (*networkingv1.IngressClassList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "ingressclasses", "IngressClass", opts, &networkingv1.IngressClassList{})
	if err != nil {
		return nil, err
	}
	return result.(*networkingv1.IngressClassList), nil
}

func (c *NetworkingK8sIoV1) WatchIngressClass(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: ".networking.k8s.io", Version: "v1", Resource: "ingressclasses"}, opts)
}

func (c *NetworkingK8sIoV1) GetNetworkPolicy(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*networkingv1.NetworkPolicy, error) {
	result, err := c.backend.Get(ctx, "networkpolicies", "NetworkPolicy", namespace, name, opts, &networkingv1.NetworkPolicy{})
	if err != nil {
		return nil, err
	}
	return result.(*networkingv1.NetworkPolicy), nil
}

func (c *NetworkingK8sIoV1) CreateNetworkPolicy(ctx context.Context, v *networkingv1.NetworkPolicy, opts metav1.CreateOptions) (*networkingv1.NetworkPolicy, error) {
	result, err := c.backend.Create(ctx, "networkpolicies", "NetworkPolicy", v, opts, &networkingv1.NetworkPolicy{})
	if err != nil {
		return nil, err
	}
	return result.(*networkingv1.NetworkPolicy), nil
}

func (c *NetworkingK8sIoV1) UpdateNetworkPolicy(ctx context.Context, v *networkingv1.NetworkPolicy, opts metav1.UpdateOptions) (*networkingv1.NetworkPolicy, error) {
	result, err := c.backend.Update(ctx, "networkpolicies", "NetworkPolicy", v, opts, &networkingv1.NetworkPolicy{})
	if err != nil {
		return nil, err
	}
	return result.(*networkingv1.NetworkPolicy), nil
}

func (c *NetworkingK8sIoV1) DeleteNetworkPolicy(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".networking.k8s.io", Version: "v1", Resource: "networkpolicies"}, namespace, name, opts)
}

func (c *NetworkingK8sIoV1) ListNetworkPolicy(ctx context.Context, namespace string, opts metav1.ListOptions) (*networkingv1.NetworkPolicyList, error) {
	result, err := c.backend.List(ctx, "networkpolicies", "NetworkPolicy", namespace, opts, &networkingv1.NetworkPolicyList{})
	if err != nil {
		return nil, err
	}
	return result.(*networkingv1.NetworkPolicyList), nil
}

func (c *NetworkingK8sIoV1) WatchNetworkPolicy(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".networking.k8s.io", Version: "v1", Resource: "networkpolicies"}, namespace, opts)
}

type PolicyV1 struct {
	backend Backend
}

func NewPolicyV1Client(b Backend) *PolicyV1 {
	return &PolicyV1{backend: b}
}

func (c *PolicyV1) GetEviction(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*policyv1.Eviction, error) {
	result, err := c.backend.Get(ctx, "evictions", "Eviction", namespace, name, opts, &policyv1.Eviction{})
	if err != nil {
		return nil, err
	}
	return result.(*policyv1.Eviction), nil
}

func (c *PolicyV1) CreateEviction(ctx context.Context, v *policyv1.Eviction, opts metav1.CreateOptions) (*policyv1.Eviction, error) {
	result, err := c.backend.Create(ctx, "evictions", "Eviction", v, opts, &policyv1.Eviction{})
	if err != nil {
		return nil, err
	}
	return result.(*policyv1.Eviction), nil
}

func (c *PolicyV1) UpdateEviction(ctx context.Context, v *policyv1.Eviction, opts metav1.UpdateOptions) (*policyv1.Eviction, error) {
	result, err := c.backend.Update(ctx, "evictions", "Eviction", v, opts, &policyv1.Eviction{})
	if err != nil {
		return nil, err
	}
	return result.(*policyv1.Eviction), nil
}

func (c *PolicyV1) DeleteEviction(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".policy", Version: "v1", Resource: "evictions"}, namespace, name, opts)
}

func (c *PolicyV1) ListEviction(ctx context.Context, namespace string, opts metav1.ListOptions) (*policyv1.EvictionList, error) {
	result, err := c.backend.List(ctx, "evictions", "Eviction", namespace, opts, &policyv1.EvictionList{})
	if err != nil {
		return nil, err
	}
	return result.(*policyv1.EvictionList), nil
}

func (c *PolicyV1) WatchEviction(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".policy", Version: "v1", Resource: "evictions"}, namespace, opts)
}

func (c *PolicyV1) GetPodDisruptionBudget(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*policyv1.PodDisruptionBudget, error) {
	result, err := c.backend.Get(ctx, "poddisruptionbudgets", "PodDisruptionBudget", namespace, name, opts, &policyv1.PodDisruptionBudget{})
	if err != nil {
		return nil, err
	}
	return result.(*policyv1.PodDisruptionBudget), nil
}

func (c *PolicyV1) CreatePodDisruptionBudget(ctx context.Context, v *policyv1.PodDisruptionBudget, opts metav1.CreateOptions) (*policyv1.PodDisruptionBudget, error) {
	result, err := c.backend.Create(ctx, "poddisruptionbudgets", "PodDisruptionBudget", v, opts, &policyv1.PodDisruptionBudget{})
	if err != nil {
		return nil, err
	}
	return result.(*policyv1.PodDisruptionBudget), nil
}

func (c *PolicyV1) UpdatePodDisruptionBudget(ctx context.Context, v *policyv1.PodDisruptionBudget, opts metav1.UpdateOptions) (*policyv1.PodDisruptionBudget, error) {
	result, err := c.backend.Update(ctx, "poddisruptionbudgets", "PodDisruptionBudget", v, opts, &policyv1.PodDisruptionBudget{})
	if err != nil {
		return nil, err
	}
	return result.(*policyv1.PodDisruptionBudget), nil
}

func (c *PolicyV1) DeletePodDisruptionBudget(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".policy", Version: "v1", Resource: "poddisruptionbudgets"}, namespace, name, opts)
}

func (c *PolicyV1) ListPodDisruptionBudget(ctx context.Context, namespace string, opts metav1.ListOptions) (*policyv1.PodDisruptionBudgetList, error) {
	result, err := c.backend.List(ctx, "poddisruptionbudgets", "PodDisruptionBudget", namespace, opts, &policyv1.PodDisruptionBudgetList{})
	if err != nil {
		return nil, err
	}
	return result.(*policyv1.PodDisruptionBudgetList), nil
}

func (c *PolicyV1) WatchPodDisruptionBudget(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".policy", Version: "v1", Resource: "poddisruptionbudgets"}, namespace, opts)
}

type RbacAuthorizationK8sIoV1 struct {
	backend Backend
}

func NewRbacAuthorizationK8sIoV1Client(b Backend) *RbacAuthorizationK8sIoV1 {
	return &RbacAuthorizationK8sIoV1{backend: b}
}

func (c *RbacAuthorizationK8sIoV1) GetClusterRole(ctx context.Context, name string, opts metav1.GetOptions) (*rbacv1.ClusterRole, error) {
	result, err := c.backend.GetClusterScoped(ctx, "clusterroles", "ClusterRole", name, opts, &rbacv1.ClusterRole{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.ClusterRole), nil
}

func (c *RbacAuthorizationK8sIoV1) CreateClusterRole(ctx context.Context, v *rbacv1.ClusterRole, opts metav1.CreateOptions) (*rbacv1.ClusterRole, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "clusterroles", "ClusterRole", v, opts, &rbacv1.ClusterRole{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.ClusterRole), nil
}

func (c *RbacAuthorizationK8sIoV1) UpdateClusterRole(ctx context.Context, v *rbacv1.ClusterRole, opts metav1.UpdateOptions) (*rbacv1.ClusterRole, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "clusterroles", "ClusterRole", v, opts, &rbacv1.ClusterRole{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.ClusterRole), nil
}

func (c *RbacAuthorizationK8sIoV1) DeleteClusterRole(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: ".rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"}, name, opts)
}

func (c *RbacAuthorizationK8sIoV1) ListClusterRole(ctx context.Context, opts metav1.ListOptions) (*rbacv1.ClusterRoleList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "clusterroles", "ClusterRole", opts, &rbacv1.ClusterRoleList{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.ClusterRoleList), nil
}

func (c *RbacAuthorizationK8sIoV1) WatchClusterRole(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: ".rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"}, opts)
}

func (c *RbacAuthorizationK8sIoV1) GetClusterRoleBinding(ctx context.Context, name string, opts metav1.GetOptions) (*rbacv1.ClusterRoleBinding, error) {
	result, err := c.backend.GetClusterScoped(ctx, "clusterrolebindings", "ClusterRoleBinding", name, opts, &rbacv1.ClusterRoleBinding{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.ClusterRoleBinding), nil
}

func (c *RbacAuthorizationK8sIoV1) CreateClusterRoleBinding(ctx context.Context, v *rbacv1.ClusterRoleBinding, opts metav1.CreateOptions) (*rbacv1.ClusterRoleBinding, error) {
	result, err := c.backend.CreateClusterScoped(ctx, "clusterrolebindings", "ClusterRoleBinding", v, opts, &rbacv1.ClusterRoleBinding{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.ClusterRoleBinding), nil
}

func (c *RbacAuthorizationK8sIoV1) UpdateClusterRoleBinding(ctx context.Context, v *rbacv1.ClusterRoleBinding, opts metav1.UpdateOptions) (*rbacv1.ClusterRoleBinding, error) {
	result, err := c.backend.UpdateClusterScoped(ctx, "clusterrolebindings", "ClusterRoleBinding", v, opts, &rbacv1.ClusterRoleBinding{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.ClusterRoleBinding), nil
}

func (c *RbacAuthorizationK8sIoV1) DeleteClusterRoleBinding(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group: ".rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings"}, name, opts)
}

func (c *RbacAuthorizationK8sIoV1) ListClusterRoleBinding(ctx context.Context, opts metav1.ListOptions) (*rbacv1.ClusterRoleBindingList, error) {
	result, err := c.backend.ListClusterScoped(ctx, "clusterrolebindings", "ClusterRoleBinding", opts, &rbacv1.ClusterRoleBindingList{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.ClusterRoleBindingList), nil
}

func (c *RbacAuthorizationK8sIoV1) WatchClusterRoleBinding(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group: ".rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings"}, opts)
}

func (c *RbacAuthorizationK8sIoV1) GetRole(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*rbacv1.Role, error) {
	result, err := c.backend.Get(ctx, "roles", "Role", namespace, name, opts, &rbacv1.Role{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.Role), nil
}

func (c *RbacAuthorizationK8sIoV1) CreateRole(ctx context.Context, v *rbacv1.Role, opts metav1.CreateOptions) (*rbacv1.Role, error) {
	result, err := c.backend.Create(ctx, "roles", "Role", v, opts, &rbacv1.Role{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.Role), nil
}

func (c *RbacAuthorizationK8sIoV1) UpdateRole(ctx context.Context, v *rbacv1.Role, opts metav1.UpdateOptions) (*rbacv1.Role, error) {
	result, err := c.backend.Update(ctx, "roles", "Role", v, opts, &rbacv1.Role{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.Role), nil
}

func (c *RbacAuthorizationK8sIoV1) DeleteRole(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".rbac.authorization.k8s.io", Version: "v1", Resource: "roles"}, namespace, name, opts)
}

func (c *RbacAuthorizationK8sIoV1) ListRole(ctx context.Context, namespace string, opts metav1.ListOptions) (*rbacv1.RoleList, error) {
	result, err := c.backend.List(ctx, "roles", "Role", namespace, opts, &rbacv1.RoleList{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.RoleList), nil
}

func (c *RbacAuthorizationK8sIoV1) WatchRole(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".rbac.authorization.k8s.io", Version: "v1", Resource: "roles"}, namespace, opts)
}

func (c *RbacAuthorizationK8sIoV1) GetRoleBinding(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*rbacv1.RoleBinding, error) {
	result, err := c.backend.Get(ctx, "rolebindings", "RoleBinding", namespace, name, opts, &rbacv1.RoleBinding{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.RoleBinding), nil
}

func (c *RbacAuthorizationK8sIoV1) CreateRoleBinding(ctx context.Context, v *rbacv1.RoleBinding, opts metav1.CreateOptions) (*rbacv1.RoleBinding, error) {
	result, err := c.backend.Create(ctx, "rolebindings", "RoleBinding", v, opts, &rbacv1.RoleBinding{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.RoleBinding), nil
}

func (c *RbacAuthorizationK8sIoV1) UpdateRoleBinding(ctx context.Context, v *rbacv1.RoleBinding, opts metav1.UpdateOptions) (*rbacv1.RoleBinding, error) {
	result, err := c.backend.Update(ctx, "rolebindings", "RoleBinding", v, opts, &rbacv1.RoleBinding{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.RoleBinding), nil
}

func (c *RbacAuthorizationK8sIoV1) DeleteRoleBinding(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.backend.Delete(ctx, schema.GroupVersionResource{Group: ".rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"}, namespace, name, opts)
}

func (c *RbacAuthorizationK8sIoV1) ListRoleBinding(ctx context.Context, namespace string, opts metav1.ListOptions) (*rbacv1.RoleBindingList, error) {
	result, err := c.backend.List(ctx, "rolebindings", "RoleBinding", namespace, opts, &rbacv1.RoleBindingList{})
	if err != nil {
		return nil, err
	}
	return result.(*rbacv1.RoleBindingList), nil
}

func (c *RbacAuthorizationK8sIoV1) WatchRoleBinding(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.backend.Watch(ctx, schema.GroupVersionResource{Group: ".rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"}, namespace, opts)
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
	case *corev1.Binding:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).BindingInformer()
	case *corev1.ComponentStatus:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).ComponentStatusInformer()
	case *corev1.ConfigMap:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).ConfigMapInformer()
	case *corev1.Endpoints:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).EndpointsInformer()
	case *corev1.Event:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).EventInformer()
	case *corev1.LimitRange:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).LimitRangeInformer()
	case *corev1.Namespace:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).NamespaceInformer()
	case *corev1.Node:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).NodeInformer()
	case *corev1.PersistentVolume:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).PersistentVolumeInformer()
	case *corev1.PersistentVolumeClaim:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).PersistentVolumeClaimInformer()
	case *corev1.Pod:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).PodInformer()
	case *corev1.PodStatusResult:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).PodStatusResultInformer()
	case *corev1.PodTemplate:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).PodTemplateInformer()
	case *corev1.RangeAllocation:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).RangeAllocationInformer()
	case *corev1.ReplicationController:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).ReplicationControllerInformer()
	case *corev1.ResourceQuota:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).ResourceQuotaInformer()
	case *corev1.Secret:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).SecretInformer()
	case *corev1.Service:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).ServiceInformer()
	case *corev1.ServiceAccount:
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).ServiceAccountInformer()
	case *admissionregistrationv1.MutatingWebhookConfiguration:
		return NewAdmissionregistrationK8sIoV1Informer(f.cache, f.set.AdmissionregistrationK8sIoV1, f.namespace, f.resyncPeriod).MutatingWebhookConfigurationInformer()
	case *admissionregistrationv1.ValidatingWebhookConfiguration:
		return NewAdmissionregistrationK8sIoV1Informer(f.cache, f.set.AdmissionregistrationK8sIoV1, f.namespace, f.resyncPeriod).ValidatingWebhookConfigurationInformer()
	case *appsv1.ControllerRevision:
		return NewAppsV1Informer(f.cache, f.set.AppsV1, f.namespace, f.resyncPeriod).ControllerRevisionInformer()
	case *batchv1.CronJob:
		return NewAppsV1Informer(f.cache, f.set.AppsV1, f.namespace, f.resyncPeriod).CronJobInformer()
	case *appsv1.DaemonSet:
		return NewAppsV1Informer(f.cache, f.set.AppsV1, f.namespace, f.resyncPeriod).DaemonSetInformer()
	case *appsv1.Deployment:
		return NewAppsV1Informer(f.cache, f.set.AppsV1, f.namespace, f.resyncPeriod).DeploymentInformer()
	case *batchv1.Job:
		return NewAppsV1Informer(f.cache, f.set.AppsV1, f.namespace, f.resyncPeriod).JobInformer()
	case *appsv1.ReplicaSet:
		return NewAppsV1Informer(f.cache, f.set.AppsV1, f.namespace, f.resyncPeriod).ReplicaSetInformer()
	case *appsv1.StatefulSet:
		return NewAppsV1Informer(f.cache, f.set.AppsV1, f.namespace, f.resyncPeriod).StatefulSetInformer()
	case *authenticationv1.TokenRequest:
		return NewAuthenticationK8sIoV1Informer(f.cache, f.set.AuthenticationK8sIoV1, f.namespace, f.resyncPeriod).TokenRequestInformer()
	case *authenticationv1.TokenReview:
		return NewAuthenticationK8sIoV1Informer(f.cache, f.set.AuthenticationK8sIoV1, f.namespace, f.resyncPeriod).TokenReviewInformer()
	case *authorizationv1.LocalSubjectAccessReview:
		return NewAuthorizationK8sIoV1Informer(f.cache, f.set.AuthorizationK8sIoV1, f.namespace, f.resyncPeriod).LocalSubjectAccessReviewInformer()
	case *authorizationv1.SelfSubjectAccessReview:
		return NewAuthorizationK8sIoV1Informer(f.cache, f.set.AuthorizationK8sIoV1, f.namespace, f.resyncPeriod).SelfSubjectAccessReviewInformer()
	case *authorizationv1.SelfSubjectRulesReview:
		return NewAuthorizationK8sIoV1Informer(f.cache, f.set.AuthorizationK8sIoV1, f.namespace, f.resyncPeriod).SelfSubjectRulesReviewInformer()
	case *authorizationv1.SubjectAccessReview:
		return NewAuthorizationK8sIoV1Informer(f.cache, f.set.AuthorizationK8sIoV1, f.namespace, f.resyncPeriod).SubjectAccessReviewInformer()
	case *autoscalingv1.HorizontalPodAutoscaler:
		return NewAutoscalingV1Informer(f.cache, f.set.AutoscalingV1, f.namespace, f.resyncPeriod).HorizontalPodAutoscalerInformer()
	case *autoscalingv1.Scale:
		return NewAutoscalingV1Informer(f.cache, f.set.AutoscalingV1, f.namespace, f.resyncPeriod).ScaleInformer()
	case *certificatesv1.CertificateSigningRequest:
		return NewCertificatesK8sIoV1Informer(f.cache, f.set.CertificatesK8sIoV1, f.namespace, f.resyncPeriod).CertificateSigningRequestInformer()
	case *discoveryv1.EndpointSlice:
		return NewDiscoveryK8sIoV1Informer(f.cache, f.set.DiscoveryK8sIoV1, f.namespace, f.resyncPeriod).EndpointSliceInformer()
	case *networkingv1.Ingress:
		return NewNetworkingK8sIoV1Informer(f.cache, f.set.NetworkingK8sIoV1, f.namespace, f.resyncPeriod).IngressInformer()
	case *networkingv1.IngressClass:
		return NewNetworkingK8sIoV1Informer(f.cache, f.set.NetworkingK8sIoV1, f.namespace, f.resyncPeriod).IngressClassInformer()
	case *networkingv1.NetworkPolicy:
		return NewNetworkingK8sIoV1Informer(f.cache, f.set.NetworkingK8sIoV1, f.namespace, f.resyncPeriod).NetworkPolicyInformer()
	case *policyv1.Eviction:
		return NewPolicyV1Informer(f.cache, f.set.PolicyV1, f.namespace, f.resyncPeriod).EvictionInformer()
	case *policyv1.PodDisruptionBudget:
		return NewPolicyV1Informer(f.cache, f.set.PolicyV1, f.namespace, f.resyncPeriod).PodDisruptionBudgetInformer()
	case *rbacv1.ClusterRole:
		return NewRbacAuthorizationK8sIoV1Informer(f.cache, f.set.RbacAuthorizationK8sIoV1, f.namespace, f.resyncPeriod).ClusterRoleInformer()
	case *rbacv1.ClusterRoleBinding:
		return NewRbacAuthorizationK8sIoV1Informer(f.cache, f.set.RbacAuthorizationK8sIoV1, f.namespace, f.resyncPeriod).ClusterRoleBindingInformer()
	case *rbacv1.Role:
		return NewRbacAuthorizationK8sIoV1Informer(f.cache, f.set.RbacAuthorizationK8sIoV1, f.namespace, f.resyncPeriod).RoleInformer()
	case *rbacv1.RoleBinding:
		return NewRbacAuthorizationK8sIoV1Informer(f.cache, f.set.RbacAuthorizationK8sIoV1, f.namespace, f.resyncPeriod).RoleBindingInformer()
	default:
		return nil
	}
}

func (f *InformerFactory) InformerForResource(gvr schema.GroupVersionResource) cache.SharedIndexInformer {
	switch gvr {
	case corev1.SchemaGroupVersion.WithResource("bindings"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).BindingInformer()
	case corev1.SchemaGroupVersion.WithResource("componentstatuses"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).ComponentStatusInformer()
	case corev1.SchemaGroupVersion.WithResource("configmaps"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).ConfigMapInformer()
	case corev1.SchemaGroupVersion.WithResource("endpoints"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).EndpointsInformer()
	case corev1.SchemaGroupVersion.WithResource("events"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).EventInformer()
	case corev1.SchemaGroupVersion.WithResource("limitranges"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).LimitRangeInformer()
	case corev1.SchemaGroupVersion.WithResource("namespaces"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).NamespaceInformer()
	case corev1.SchemaGroupVersion.WithResource("nodes"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).NodeInformer()
	case corev1.SchemaGroupVersion.WithResource("persistentvolumes"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).PersistentVolumeInformer()
	case corev1.SchemaGroupVersion.WithResource("persistentvolumeclaims"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).PersistentVolumeClaimInformer()
	case corev1.SchemaGroupVersion.WithResource("pods"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).PodInformer()
	case corev1.SchemaGroupVersion.WithResource("podstatusresults"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).PodStatusResultInformer()
	case corev1.SchemaGroupVersion.WithResource("podtemplates"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).PodTemplateInformer()
	case corev1.SchemaGroupVersion.WithResource("rangeallocations"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).RangeAllocationInformer()
	case corev1.SchemaGroupVersion.WithResource("replicationcontrollers"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).ReplicationControllerInformer()
	case corev1.SchemaGroupVersion.WithResource("resourcequotas"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).ResourceQuotaInformer()
	case corev1.SchemaGroupVersion.WithResource("secrets"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).SecretInformer()
	case corev1.SchemaGroupVersion.WithResource("services"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).ServiceInformer()
	case corev1.SchemaGroupVersion.WithResource("serviceaccounts"):
		return NewCoreV1Informer(f.cache, f.set.CoreV1, f.namespace, f.resyncPeriod).ServiceAccountInformer()
	case admissionregistrationv1.SchemaGroupVersion.WithResource("mutatingwebhookconfigurations"):
		return NewAdmissionregistrationK8sIoV1Informer(f.cache, f.set.AdmissionregistrationK8sIoV1, f.namespace, f.resyncPeriod).MutatingWebhookConfigurationInformer()
	case admissionregistrationv1.SchemaGroupVersion.WithResource("validatingwebhookconfigurations"):
		return NewAdmissionregistrationK8sIoV1Informer(f.cache, f.set.AdmissionregistrationK8sIoV1, f.namespace, f.resyncPeriod).ValidatingWebhookConfigurationInformer()
	case appsv1.SchemaGroupVersion.WithResource("controllerrevisions"):
		return NewAppsV1Informer(f.cache, f.set.AppsV1, f.namespace, f.resyncPeriod).ControllerRevisionInformer()
	case batchv1.SchemaGroupVersion.WithResource("cronjobs"):
		return NewAppsV1Informer(f.cache, f.set.AppsV1, f.namespace, f.resyncPeriod).CronJobInformer()
	case appsv1.SchemaGroupVersion.WithResource("daemonsets"):
		return NewAppsV1Informer(f.cache, f.set.AppsV1, f.namespace, f.resyncPeriod).DaemonSetInformer()
	case appsv1.SchemaGroupVersion.WithResource("deployments"):
		return NewAppsV1Informer(f.cache, f.set.AppsV1, f.namespace, f.resyncPeriod).DeploymentInformer()
	case batchv1.SchemaGroupVersion.WithResource("jobs"):
		return NewAppsV1Informer(f.cache, f.set.AppsV1, f.namespace, f.resyncPeriod).JobInformer()
	case appsv1.SchemaGroupVersion.WithResource("replicasets"):
		return NewAppsV1Informer(f.cache, f.set.AppsV1, f.namespace, f.resyncPeriod).ReplicaSetInformer()
	case appsv1.SchemaGroupVersion.WithResource("statefulsets"):
		return NewAppsV1Informer(f.cache, f.set.AppsV1, f.namespace, f.resyncPeriod).StatefulSetInformer()
	case authenticationv1.SchemaGroupVersion.WithResource("tokenrequests"):
		return NewAuthenticationK8sIoV1Informer(f.cache, f.set.AuthenticationK8sIoV1, f.namespace, f.resyncPeriod).TokenRequestInformer()
	case authenticationv1.SchemaGroupVersion.WithResource("tokenreviews"):
		return NewAuthenticationK8sIoV1Informer(f.cache, f.set.AuthenticationK8sIoV1, f.namespace, f.resyncPeriod).TokenReviewInformer()
	case authorizationv1.SchemaGroupVersion.WithResource("localsubjectaccessreviews"):
		return NewAuthorizationK8sIoV1Informer(f.cache, f.set.AuthorizationK8sIoV1, f.namespace, f.resyncPeriod).LocalSubjectAccessReviewInformer()
	case authorizationv1.SchemaGroupVersion.WithResource("selfsubjectaccessreviews"):
		return NewAuthorizationK8sIoV1Informer(f.cache, f.set.AuthorizationK8sIoV1, f.namespace, f.resyncPeriod).SelfSubjectAccessReviewInformer()
	case authorizationv1.SchemaGroupVersion.WithResource("selfsubjectrulesreviews"):
		return NewAuthorizationK8sIoV1Informer(f.cache, f.set.AuthorizationK8sIoV1, f.namespace, f.resyncPeriod).SelfSubjectRulesReviewInformer()
	case authorizationv1.SchemaGroupVersion.WithResource("subjectaccessreviews"):
		return NewAuthorizationK8sIoV1Informer(f.cache, f.set.AuthorizationK8sIoV1, f.namespace, f.resyncPeriod).SubjectAccessReviewInformer()
	case autoscalingv1.SchemaGroupVersion.WithResource("horizontalpodautoscalers"):
		return NewAutoscalingV1Informer(f.cache, f.set.AutoscalingV1, f.namespace, f.resyncPeriod).HorizontalPodAutoscalerInformer()
	case autoscalingv1.SchemaGroupVersion.WithResource("scales"):
		return NewAutoscalingV1Informer(f.cache, f.set.AutoscalingV1, f.namespace, f.resyncPeriod).ScaleInformer()
	case certificatesv1.SchemaGroupVersion.WithResource("certificatesigningrequests"):
		return NewCertificatesK8sIoV1Informer(f.cache, f.set.CertificatesK8sIoV1, f.namespace, f.resyncPeriod).CertificateSigningRequestInformer()
	case discoveryv1.SchemaGroupVersion.WithResource("endpointslices"):
		return NewDiscoveryK8sIoV1Informer(f.cache, f.set.DiscoveryK8sIoV1, f.namespace, f.resyncPeriod).EndpointSliceInformer()
	case networkingv1.SchemaGroupVersion.WithResource("ingresses"):
		return NewNetworkingK8sIoV1Informer(f.cache, f.set.NetworkingK8sIoV1, f.namespace, f.resyncPeriod).IngressInformer()
	case networkingv1.SchemaGroupVersion.WithResource("ingressclasses"):
		return NewNetworkingK8sIoV1Informer(f.cache, f.set.NetworkingK8sIoV1, f.namespace, f.resyncPeriod).IngressClassInformer()
	case networkingv1.SchemaGroupVersion.WithResource("networkpolicies"):
		return NewNetworkingK8sIoV1Informer(f.cache, f.set.NetworkingK8sIoV1, f.namespace, f.resyncPeriod).NetworkPolicyInformer()
	case policyv1.SchemaGroupVersion.WithResource("evictions"):
		return NewPolicyV1Informer(f.cache, f.set.PolicyV1, f.namespace, f.resyncPeriod).EvictionInformer()
	case policyv1.SchemaGroupVersion.WithResource("poddisruptionbudgets"):
		return NewPolicyV1Informer(f.cache, f.set.PolicyV1, f.namespace, f.resyncPeriod).PodDisruptionBudgetInformer()
	case rbacv1.SchemaGroupVersion.WithResource("clusterroles"):
		return NewRbacAuthorizationK8sIoV1Informer(f.cache, f.set.RbacAuthorizationK8sIoV1, f.namespace, f.resyncPeriod).ClusterRoleInformer()
	case rbacv1.SchemaGroupVersion.WithResource("clusterrolebindings"):
		return NewRbacAuthorizationK8sIoV1Informer(f.cache, f.set.RbacAuthorizationK8sIoV1, f.namespace, f.resyncPeriod).ClusterRoleBindingInformer()
	case rbacv1.SchemaGroupVersion.WithResource("roles"):
		return NewRbacAuthorizationK8sIoV1Informer(f.cache, f.set.RbacAuthorizationK8sIoV1, f.namespace, f.resyncPeriod).RoleInformer()
	case rbacv1.SchemaGroupVersion.WithResource("rolebindings"):
		return NewRbacAuthorizationK8sIoV1Informer(f.cache, f.set.RbacAuthorizationK8sIoV1, f.namespace, f.resyncPeriod).RoleBindingInformer()
	default:
		return nil
	}
}

func (f *InformerFactory) Run(ctx context.Context) {
	for _, v := range f.cache.Informers() {
		go v.Run(ctx.Done())
	}
}

type CoreV1Informer struct {
	cache        *InformerCache
	client       *CoreV1
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewCoreV1Informer(c *InformerCache, client *CoreV1, namespace string, resyncPeriod time.Duration) *CoreV1Informer {
	return &CoreV1Informer{
		cache:        c,
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *CoreV1Informer) BindingInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.Binding{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListBinding(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchBinding(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&corev1.Binding{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) BindingLister() *CoreV1BindingLister {
	return NewCoreV1BindingLister(f.BindingInformer().GetIndexer())
}

func (f *CoreV1Informer) ComponentStatusInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.ComponentStatus{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListComponentStatus(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchComponentStatus(context.TODO(), metav1.ListOptions{})
				},
			},
			&corev1.ComponentStatus{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) ComponentStatusLister() *CoreV1ComponentStatusLister {
	return NewCoreV1ComponentStatusLister(f.ComponentStatusInformer().GetIndexer())
}

func (f *CoreV1Informer) ConfigMapInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.ConfigMap{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListConfigMap(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchConfigMap(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&corev1.ConfigMap{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) ConfigMapLister() *CoreV1ConfigMapLister {
	return NewCoreV1ConfigMapLister(f.ConfigMapInformer().GetIndexer())
}

func (f *CoreV1Informer) EndpointsInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.Endpoints{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListEndpoints(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchEndpoints(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&corev1.Endpoints{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) EndpointsLister() *CoreV1EndpointsLister {
	return NewCoreV1EndpointsLister(f.EndpointsInformer().GetIndexer())
}

func (f *CoreV1Informer) EventInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.Event{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListEvent(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchEvent(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&corev1.Event{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) EventLister() *CoreV1EventLister {
	return NewCoreV1EventLister(f.EventInformer().GetIndexer())
}

func (f *CoreV1Informer) LimitRangeInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.LimitRange{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListLimitRange(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchLimitRange(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&corev1.LimitRange{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) LimitRangeLister() *CoreV1LimitRangeLister {
	return NewCoreV1LimitRangeLister(f.LimitRangeInformer().GetIndexer())
}

func (f *CoreV1Informer) NamespaceInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.Namespace{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListNamespace(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchNamespace(context.TODO(), metav1.ListOptions{})
				},
			},
			&corev1.Namespace{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) NamespaceLister() *CoreV1NamespaceLister {
	return NewCoreV1NamespaceLister(f.NamespaceInformer().GetIndexer())
}

func (f *CoreV1Informer) NodeInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.Node{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListNode(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchNode(context.TODO(), metav1.ListOptions{})
				},
			},
			&corev1.Node{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) NodeLister() *CoreV1NodeLister {
	return NewCoreV1NodeLister(f.NodeInformer().GetIndexer())
}

func (f *CoreV1Informer) PersistentVolumeInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.PersistentVolume{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListPersistentVolume(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchPersistentVolume(context.TODO(), metav1.ListOptions{})
				},
			},
			&corev1.PersistentVolume{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) PersistentVolumeLister() *CoreV1PersistentVolumeLister {
	return NewCoreV1PersistentVolumeLister(f.PersistentVolumeInformer().GetIndexer())
}

func (f *CoreV1Informer) PersistentVolumeClaimInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.PersistentVolumeClaim{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListPersistentVolumeClaim(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchPersistentVolumeClaim(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&corev1.PersistentVolumeClaim{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) PersistentVolumeClaimLister() *CoreV1PersistentVolumeClaimLister {
	return NewCoreV1PersistentVolumeClaimLister(f.PersistentVolumeClaimInformer().GetIndexer())
}

func (f *CoreV1Informer) PodInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.Pod{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListPod(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchPod(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&corev1.Pod{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) PodLister() *CoreV1PodLister {
	return NewCoreV1PodLister(f.PodInformer().GetIndexer())
}

func (f *CoreV1Informer) PodStatusResultInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.PodStatusResult{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListPodStatusResult(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchPodStatusResult(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&corev1.PodStatusResult{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) PodStatusResultLister() *CoreV1PodStatusResultLister {
	return NewCoreV1PodStatusResultLister(f.PodStatusResultInformer().GetIndexer())
}

func (f *CoreV1Informer) PodTemplateInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.PodTemplate{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListPodTemplate(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchPodTemplate(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&corev1.PodTemplate{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) PodTemplateLister() *CoreV1PodTemplateLister {
	return NewCoreV1PodTemplateLister(f.PodTemplateInformer().GetIndexer())
}

func (f *CoreV1Informer) RangeAllocationInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.RangeAllocation{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListRangeAllocation(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchRangeAllocation(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&corev1.RangeAllocation{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) RangeAllocationLister() *CoreV1RangeAllocationLister {
	return NewCoreV1RangeAllocationLister(f.RangeAllocationInformer().GetIndexer())
}

func (f *CoreV1Informer) ReplicationControllerInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.ReplicationController{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListReplicationController(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchReplicationController(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&corev1.ReplicationController{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) ReplicationControllerLister() *CoreV1ReplicationControllerLister {
	return NewCoreV1ReplicationControllerLister(f.ReplicationControllerInformer().GetIndexer())
}

func (f *CoreV1Informer) ResourceQuotaInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.ResourceQuota{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListResourceQuota(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchResourceQuota(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&corev1.ResourceQuota{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) ResourceQuotaLister() *CoreV1ResourceQuotaLister {
	return NewCoreV1ResourceQuotaLister(f.ResourceQuotaInformer().GetIndexer())
}

func (f *CoreV1Informer) SecretInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.Secret{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListSecret(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchSecret(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&corev1.Secret{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) SecretLister() *CoreV1SecretLister {
	return NewCoreV1SecretLister(f.SecretInformer().GetIndexer())
}

func (f *CoreV1Informer) ServiceInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.Service{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListService(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchService(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&corev1.Service{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) ServiceLister() *CoreV1ServiceLister {
	return NewCoreV1ServiceLister(f.ServiceInformer().GetIndexer())
}

func (f *CoreV1Informer) ServiceAccountInformer() cache.SharedIndexInformer {
	return f.cache.Write(&corev1.ServiceAccount{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListServiceAccount(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchServiceAccount(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&corev1.ServiceAccount{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CoreV1Informer) ServiceAccountLister() *CoreV1ServiceAccountLister {
	return NewCoreV1ServiceAccountLister(f.ServiceAccountInformer().GetIndexer())
}

type AdmissionregistrationK8sIoV1Informer struct {
	cache        *InformerCache
	client       *AdmissionregistrationK8sIoV1
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewAdmissionregistrationK8sIoV1Informer(c *InformerCache, client *AdmissionregistrationK8sIoV1, namespace string, resyncPeriod time.Duration) *AdmissionregistrationK8sIoV1Informer {
	return &AdmissionregistrationK8sIoV1Informer{
		cache:        c,
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *AdmissionregistrationK8sIoV1Informer) MutatingWebhookConfigurationInformer() cache.SharedIndexInformer {
	return f.cache.Write(&admissionregistrationv1.MutatingWebhookConfiguration{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListMutatingWebhookConfiguration(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchMutatingWebhookConfiguration(context.TODO(), metav1.ListOptions{})
				},
			},
			&admissionregistrationv1.MutatingWebhookConfiguration{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AdmissionregistrationK8sIoV1Informer) MutatingWebhookConfigurationLister() *AdmissionregistrationK8sIoV1MutatingWebhookConfigurationLister {
	return NewAdmissionregistrationK8sIoV1MutatingWebhookConfigurationLister(f.MutatingWebhookConfigurationInformer().GetIndexer())
}

func (f *AdmissionregistrationK8sIoV1Informer) ValidatingWebhookConfigurationInformer() cache.SharedIndexInformer {
	return f.cache.Write(&admissionregistrationv1.ValidatingWebhookConfiguration{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListValidatingWebhookConfiguration(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchValidatingWebhookConfiguration(context.TODO(), metav1.ListOptions{})
				},
			},
			&admissionregistrationv1.ValidatingWebhookConfiguration{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AdmissionregistrationK8sIoV1Informer) ValidatingWebhookConfigurationLister() *AdmissionregistrationK8sIoV1ValidatingWebhookConfigurationLister {
	return NewAdmissionregistrationK8sIoV1ValidatingWebhookConfigurationLister(f.ValidatingWebhookConfigurationInformer().GetIndexer())
}

type AppsV1Informer struct {
	cache        *InformerCache
	client       *AppsV1
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewAppsV1Informer(c *InformerCache, client *AppsV1, namespace string, resyncPeriod time.Duration) *AppsV1Informer {
	return &AppsV1Informer{
		cache:        c,
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *AppsV1Informer) ControllerRevisionInformer() cache.SharedIndexInformer {
	return f.cache.Write(&appsv1.ControllerRevision{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListControllerRevision(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchControllerRevision(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&appsv1.ControllerRevision{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AppsV1Informer) ControllerRevisionLister() *AppsV1ControllerRevisionLister {
	return NewAppsV1ControllerRevisionLister(f.ControllerRevisionInformer().GetIndexer())
}

func (f *AppsV1Informer) CronJobInformer() cache.SharedIndexInformer {
	return f.cache.Write(&batchv1.CronJob{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListCronJob(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchCronJob(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&batchv1.CronJob{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AppsV1Informer) CronJobLister() *AppsV1CronJobLister {
	return NewAppsV1CronJobLister(f.CronJobInformer().GetIndexer())
}

func (f *AppsV1Informer) DaemonSetInformer() cache.SharedIndexInformer {
	return f.cache.Write(&appsv1.DaemonSet{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListDaemonSet(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchDaemonSet(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&appsv1.DaemonSet{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AppsV1Informer) DaemonSetLister() *AppsV1DaemonSetLister {
	return NewAppsV1DaemonSetLister(f.DaemonSetInformer().GetIndexer())
}

func (f *AppsV1Informer) DeploymentInformer() cache.SharedIndexInformer {
	return f.cache.Write(&appsv1.Deployment{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListDeployment(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchDeployment(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&appsv1.Deployment{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AppsV1Informer) DeploymentLister() *AppsV1DeploymentLister {
	return NewAppsV1DeploymentLister(f.DeploymentInformer().GetIndexer())
}

func (f *AppsV1Informer) JobInformer() cache.SharedIndexInformer {
	return f.cache.Write(&batchv1.Job{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListJob(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchJob(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&batchv1.Job{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AppsV1Informer) JobLister() *AppsV1JobLister {
	return NewAppsV1JobLister(f.JobInformer().GetIndexer())
}

func (f *AppsV1Informer) ReplicaSetInformer() cache.SharedIndexInformer {
	return f.cache.Write(&appsv1.ReplicaSet{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListReplicaSet(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchReplicaSet(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&appsv1.ReplicaSet{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AppsV1Informer) ReplicaSetLister() *AppsV1ReplicaSetLister {
	return NewAppsV1ReplicaSetLister(f.ReplicaSetInformer().GetIndexer())
}

func (f *AppsV1Informer) StatefulSetInformer() cache.SharedIndexInformer {
	return f.cache.Write(&appsv1.StatefulSet{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListStatefulSet(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchStatefulSet(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&appsv1.StatefulSet{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AppsV1Informer) StatefulSetLister() *AppsV1StatefulSetLister {
	return NewAppsV1StatefulSetLister(f.StatefulSetInformer().GetIndexer())
}

type AuthenticationK8sIoV1Informer struct {
	cache        *InformerCache
	client       *AuthenticationK8sIoV1
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewAuthenticationK8sIoV1Informer(c *InformerCache, client *AuthenticationK8sIoV1, namespace string, resyncPeriod time.Duration) *AuthenticationK8sIoV1Informer {
	return &AuthenticationK8sIoV1Informer{
		cache:        c,
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *AuthenticationK8sIoV1Informer) TokenRequestInformer() cache.SharedIndexInformer {
	return f.cache.Write(&authenticationv1.TokenRequest{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListTokenRequest(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchTokenRequest(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&authenticationv1.TokenRequest{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AuthenticationK8sIoV1Informer) TokenRequestLister() *AuthenticationK8sIoV1TokenRequestLister {
	return NewAuthenticationK8sIoV1TokenRequestLister(f.TokenRequestInformer().GetIndexer())
}

func (f *AuthenticationK8sIoV1Informer) TokenReviewInformer() cache.SharedIndexInformer {
	return f.cache.Write(&authenticationv1.TokenReview{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListTokenReview(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchTokenReview(context.TODO(), metav1.ListOptions{})
				},
			},
			&authenticationv1.TokenReview{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AuthenticationK8sIoV1Informer) TokenReviewLister() *AuthenticationK8sIoV1TokenReviewLister {
	return NewAuthenticationK8sIoV1TokenReviewLister(f.TokenReviewInformer().GetIndexer())
}

type AuthorizationK8sIoV1Informer struct {
	cache        *InformerCache
	client       *AuthorizationK8sIoV1
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewAuthorizationK8sIoV1Informer(c *InformerCache, client *AuthorizationK8sIoV1, namespace string, resyncPeriod time.Duration) *AuthorizationK8sIoV1Informer {
	return &AuthorizationK8sIoV1Informer{
		cache:        c,
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *AuthorizationK8sIoV1Informer) LocalSubjectAccessReviewInformer() cache.SharedIndexInformer {
	return f.cache.Write(&authorizationv1.LocalSubjectAccessReview{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListLocalSubjectAccessReview(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchLocalSubjectAccessReview(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&authorizationv1.LocalSubjectAccessReview{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AuthorizationK8sIoV1Informer) LocalSubjectAccessReviewLister() *AuthorizationK8sIoV1LocalSubjectAccessReviewLister {
	return NewAuthorizationK8sIoV1LocalSubjectAccessReviewLister(f.LocalSubjectAccessReviewInformer().GetIndexer())
}

func (f *AuthorizationK8sIoV1Informer) SelfSubjectAccessReviewInformer() cache.SharedIndexInformer {
	return f.cache.Write(&authorizationv1.SelfSubjectAccessReview{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListSelfSubjectAccessReview(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchSelfSubjectAccessReview(context.TODO(), metav1.ListOptions{})
				},
			},
			&authorizationv1.SelfSubjectAccessReview{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AuthorizationK8sIoV1Informer) SelfSubjectAccessReviewLister() *AuthorizationK8sIoV1SelfSubjectAccessReviewLister {
	return NewAuthorizationK8sIoV1SelfSubjectAccessReviewLister(f.SelfSubjectAccessReviewInformer().GetIndexer())
}

func (f *AuthorizationK8sIoV1Informer) SelfSubjectRulesReviewInformer() cache.SharedIndexInformer {
	return f.cache.Write(&authorizationv1.SelfSubjectRulesReview{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListSelfSubjectRulesReview(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchSelfSubjectRulesReview(context.TODO(), metav1.ListOptions{})
				},
			},
			&authorizationv1.SelfSubjectRulesReview{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AuthorizationK8sIoV1Informer) SelfSubjectRulesReviewLister() *AuthorizationK8sIoV1SelfSubjectRulesReviewLister {
	return NewAuthorizationK8sIoV1SelfSubjectRulesReviewLister(f.SelfSubjectRulesReviewInformer().GetIndexer())
}

func (f *AuthorizationK8sIoV1Informer) SubjectAccessReviewInformer() cache.SharedIndexInformer {
	return f.cache.Write(&authorizationv1.SubjectAccessReview{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListSubjectAccessReview(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchSubjectAccessReview(context.TODO(), metav1.ListOptions{})
				},
			},
			&authorizationv1.SubjectAccessReview{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AuthorizationK8sIoV1Informer) SubjectAccessReviewLister() *AuthorizationK8sIoV1SubjectAccessReviewLister {
	return NewAuthorizationK8sIoV1SubjectAccessReviewLister(f.SubjectAccessReviewInformer().GetIndexer())
}

type AutoscalingV1Informer struct {
	cache        *InformerCache
	client       *AutoscalingV1
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewAutoscalingV1Informer(c *InformerCache, client *AutoscalingV1, namespace string, resyncPeriod time.Duration) *AutoscalingV1Informer {
	return &AutoscalingV1Informer{
		cache:        c,
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *AutoscalingV1Informer) HorizontalPodAutoscalerInformer() cache.SharedIndexInformer {
	return f.cache.Write(&autoscalingv1.HorizontalPodAutoscaler{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListHorizontalPodAutoscaler(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchHorizontalPodAutoscaler(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&autoscalingv1.HorizontalPodAutoscaler{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AutoscalingV1Informer) HorizontalPodAutoscalerLister() *AutoscalingV1HorizontalPodAutoscalerLister {
	return NewAutoscalingV1HorizontalPodAutoscalerLister(f.HorizontalPodAutoscalerInformer().GetIndexer())
}

func (f *AutoscalingV1Informer) ScaleInformer() cache.SharedIndexInformer {
	return f.cache.Write(&autoscalingv1.Scale{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListScale(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchScale(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&autoscalingv1.Scale{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *AutoscalingV1Informer) ScaleLister() *AutoscalingV1ScaleLister {
	return NewAutoscalingV1ScaleLister(f.ScaleInformer().GetIndexer())
}

type CertificatesK8sIoV1Informer struct {
	cache        *InformerCache
	client       *CertificatesK8sIoV1
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewCertificatesK8sIoV1Informer(c *InformerCache, client *CertificatesK8sIoV1, namespace string, resyncPeriod time.Duration) *CertificatesK8sIoV1Informer {
	return &CertificatesK8sIoV1Informer{
		cache:        c,
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *CertificatesK8sIoV1Informer) CertificateSigningRequestInformer() cache.SharedIndexInformer {
	return f.cache.Write(&certificatesv1.CertificateSigningRequest{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListCertificateSigningRequest(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchCertificateSigningRequest(context.TODO(), metav1.ListOptions{})
				},
			},
			&certificatesv1.CertificateSigningRequest{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *CertificatesK8sIoV1Informer) CertificateSigningRequestLister() *CertificatesK8sIoV1CertificateSigningRequestLister {
	return NewCertificatesK8sIoV1CertificateSigningRequestLister(f.CertificateSigningRequestInformer().GetIndexer())
}

type DiscoveryK8sIoV1Informer struct {
	cache        *InformerCache
	client       *DiscoveryK8sIoV1
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewDiscoveryK8sIoV1Informer(c *InformerCache, client *DiscoveryK8sIoV1, namespace string, resyncPeriod time.Duration) *DiscoveryK8sIoV1Informer {
	return &DiscoveryK8sIoV1Informer{
		cache:        c,
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *DiscoveryK8sIoV1Informer) EndpointSliceInformer() cache.SharedIndexInformer {
	return f.cache.Write(&discoveryv1.EndpointSlice{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListEndpointSlice(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchEndpointSlice(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&discoveryv1.EndpointSlice{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *DiscoveryK8sIoV1Informer) EndpointSliceLister() *DiscoveryK8sIoV1EndpointSliceLister {
	return NewDiscoveryK8sIoV1EndpointSliceLister(f.EndpointSliceInformer().GetIndexer())
}

type NetworkingK8sIoV1Informer struct {
	cache        *InformerCache
	client       *NetworkingK8sIoV1
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewNetworkingK8sIoV1Informer(c *InformerCache, client *NetworkingK8sIoV1, namespace string, resyncPeriod time.Duration) *NetworkingK8sIoV1Informer {
	return &NetworkingK8sIoV1Informer{
		cache:        c,
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *NetworkingK8sIoV1Informer) IngressInformer() cache.SharedIndexInformer {
	return f.cache.Write(&networkingv1.Ingress{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListIngress(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchIngress(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&networkingv1.Ingress{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *NetworkingK8sIoV1Informer) IngressLister() *NetworkingK8sIoV1IngressLister {
	return NewNetworkingK8sIoV1IngressLister(f.IngressInformer().GetIndexer())
}

func (f *NetworkingK8sIoV1Informer) IngressClassInformer() cache.SharedIndexInformer {
	return f.cache.Write(&networkingv1.IngressClass{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListIngressClass(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchIngressClass(context.TODO(), metav1.ListOptions{})
				},
			},
			&networkingv1.IngressClass{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *NetworkingK8sIoV1Informer) IngressClassLister() *NetworkingK8sIoV1IngressClassLister {
	return NewNetworkingK8sIoV1IngressClassLister(f.IngressClassInformer().GetIndexer())
}

func (f *NetworkingK8sIoV1Informer) NetworkPolicyInformer() cache.SharedIndexInformer {
	return f.cache.Write(&networkingv1.NetworkPolicy{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListNetworkPolicy(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchNetworkPolicy(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&networkingv1.NetworkPolicy{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *NetworkingK8sIoV1Informer) NetworkPolicyLister() *NetworkingK8sIoV1NetworkPolicyLister {
	return NewNetworkingK8sIoV1NetworkPolicyLister(f.NetworkPolicyInformer().GetIndexer())
}

type PolicyV1Informer struct {
	cache        *InformerCache
	client       *PolicyV1
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewPolicyV1Informer(c *InformerCache, client *PolicyV1, namespace string, resyncPeriod time.Duration) *PolicyV1Informer {
	return &PolicyV1Informer{
		cache:        c,
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *PolicyV1Informer) EvictionInformer() cache.SharedIndexInformer {
	return f.cache.Write(&policyv1.Eviction{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListEviction(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchEviction(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&policyv1.Eviction{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *PolicyV1Informer) EvictionLister() *PolicyV1EvictionLister {
	return NewPolicyV1EvictionLister(f.EvictionInformer().GetIndexer())
}

func (f *PolicyV1Informer) PodDisruptionBudgetInformer() cache.SharedIndexInformer {
	return f.cache.Write(&policyv1.PodDisruptionBudget{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListPodDisruptionBudget(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchPodDisruptionBudget(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&policyv1.PodDisruptionBudget{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *PolicyV1Informer) PodDisruptionBudgetLister() *PolicyV1PodDisruptionBudgetLister {
	return NewPolicyV1PodDisruptionBudgetLister(f.PodDisruptionBudgetInformer().GetIndexer())
}

type RbacAuthorizationK8sIoV1Informer struct {
	cache        *InformerCache
	client       *RbacAuthorizationK8sIoV1
	namespace    string
	resyncPeriod time.Duration
	indexers     cache.Indexers
}

func NewRbacAuthorizationK8sIoV1Informer(c *InformerCache, client *RbacAuthorizationK8sIoV1, namespace string, resyncPeriod time.Duration) *RbacAuthorizationK8sIoV1Informer {
	return &RbacAuthorizationK8sIoV1Informer{
		cache:        c,
		client:       client,
		namespace:    namespace,
		resyncPeriod: resyncPeriod,
		indexers:     cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	}
}

func (f *RbacAuthorizationK8sIoV1Informer) ClusterRoleInformer() cache.SharedIndexInformer {
	return f.cache.Write(&rbacv1.ClusterRole{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListClusterRole(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchClusterRole(context.TODO(), metav1.ListOptions{})
				},
			},
			&rbacv1.ClusterRole{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *RbacAuthorizationK8sIoV1Informer) ClusterRoleLister() *RbacAuthorizationK8sIoV1ClusterRoleLister {
	return NewRbacAuthorizationK8sIoV1ClusterRoleLister(f.ClusterRoleInformer().GetIndexer())
}

func (f *RbacAuthorizationK8sIoV1Informer) ClusterRoleBindingInformer() cache.SharedIndexInformer {
	return f.cache.Write(&rbacv1.ClusterRoleBinding{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListClusterRoleBinding(context.TODO(), metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchClusterRoleBinding(context.TODO(), metav1.ListOptions{})
				},
			},
			&rbacv1.ClusterRoleBinding{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *RbacAuthorizationK8sIoV1Informer) ClusterRoleBindingLister() *RbacAuthorizationK8sIoV1ClusterRoleBindingLister {
	return NewRbacAuthorizationK8sIoV1ClusterRoleBindingLister(f.ClusterRoleBindingInformer().GetIndexer())
}

func (f *RbacAuthorizationK8sIoV1Informer) RoleInformer() cache.SharedIndexInformer {
	return f.cache.Write(&rbacv1.Role{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListRole(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchRole(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&rbacv1.Role{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *RbacAuthorizationK8sIoV1Informer) RoleLister() *RbacAuthorizationK8sIoV1RoleLister {
	return NewRbacAuthorizationK8sIoV1RoleLister(f.RoleInformer().GetIndexer())
}

func (f *RbacAuthorizationK8sIoV1Informer) RoleBindingInformer() cache.SharedIndexInformer {
	return f.cache.Write(&rbacv1.RoleBinding{}, func() cache.SharedIndexInformer {
		return cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return f.client.ListRoleBinding(context.TODO(), f.namespace, metav1.ListOptions{})
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return f.client.WatchRoleBinding(context.TODO(), f.namespace, metav1.ListOptions{})
				},
			},
			&rbacv1.RoleBinding{},
			f.resyncPeriod,
			f.indexers,
		)
	})
}

func (f *RbacAuthorizationK8sIoV1Informer) RoleBindingLister() *RbacAuthorizationK8sIoV1RoleBindingLister {
	return NewRbacAuthorizationK8sIoV1RoleBindingLister(f.RoleBindingInformer().GetIndexer())
}

type CoreV1BindingLister struct {
	indexer cache.Indexer
}

func NewCoreV1BindingLister(indexer cache.Indexer) *CoreV1BindingLister {
	return &CoreV1BindingLister{indexer: indexer}
}

func (x *CoreV1BindingLister) List(namespace string, selector labels.Selector) ([]*corev1.Binding, error) {
	var ret []*corev1.Binding
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.Binding).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1BindingLister) Get(namespace, name string) (*corev1.Binding, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("binding").GroupResource(), name)
	}
	return obj.(*corev1.Binding).DeepCopy(), nil
}

type CoreV1ComponentStatusLister struct {
	indexer cache.Indexer
}

func NewCoreV1ComponentStatusLister(indexer cache.Indexer) *CoreV1ComponentStatusLister {
	return &CoreV1ComponentStatusLister{indexer: indexer}
}

func (x *CoreV1ComponentStatusLister) List(selector labels.Selector) ([]*corev1.ComponentStatus, error) {
	var ret []*corev1.ComponentStatus
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.ComponentStatus).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1ComponentStatusLister) Get(name string) (*corev1.ComponentStatus, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("componentstatus").GroupResource(), name)
	}
	return obj.(*corev1.ComponentStatus).DeepCopy(), nil
}

type CoreV1ConfigMapLister struct {
	indexer cache.Indexer
}

func NewCoreV1ConfigMapLister(indexer cache.Indexer) *CoreV1ConfigMapLister {
	return &CoreV1ConfigMapLister{indexer: indexer}
}

func (x *CoreV1ConfigMapLister) List(namespace string, selector labels.Selector) ([]*corev1.ConfigMap, error) {
	var ret []*corev1.ConfigMap
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.ConfigMap).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1ConfigMapLister) Get(namespace, name string) (*corev1.ConfigMap, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("configmap").GroupResource(), name)
	}
	return obj.(*corev1.ConfigMap).DeepCopy(), nil
}

type CoreV1EndpointsLister struct {
	indexer cache.Indexer
}

func NewCoreV1EndpointsLister(indexer cache.Indexer) *CoreV1EndpointsLister {
	return &CoreV1EndpointsLister{indexer: indexer}
}

func (x *CoreV1EndpointsLister) List(namespace string, selector labels.Selector) ([]*corev1.Endpoints, error) {
	var ret []*corev1.Endpoints
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.Endpoints).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1EndpointsLister) Get(namespace, name string) (*corev1.Endpoints, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("endpoints").GroupResource(), name)
	}
	return obj.(*corev1.Endpoints).DeepCopy(), nil
}

type CoreV1EventLister struct {
	indexer cache.Indexer
}

func NewCoreV1EventLister(indexer cache.Indexer) *CoreV1EventLister {
	return &CoreV1EventLister{indexer: indexer}
}

func (x *CoreV1EventLister) List(namespace string, selector labels.Selector) ([]*corev1.Event, error) {
	var ret []*corev1.Event
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.Event).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1EventLister) Get(namespace, name string) (*corev1.Event, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("event").GroupResource(), name)
	}
	return obj.(*corev1.Event).DeepCopy(), nil
}

type CoreV1LimitRangeLister struct {
	indexer cache.Indexer
}

func NewCoreV1LimitRangeLister(indexer cache.Indexer) *CoreV1LimitRangeLister {
	return &CoreV1LimitRangeLister{indexer: indexer}
}

func (x *CoreV1LimitRangeLister) List(namespace string, selector labels.Selector) ([]*corev1.LimitRange, error) {
	var ret []*corev1.LimitRange
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.LimitRange).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1LimitRangeLister) Get(namespace, name string) (*corev1.LimitRange, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("limitrange").GroupResource(), name)
	}
	return obj.(*corev1.LimitRange).DeepCopy(), nil
}

type CoreV1NamespaceLister struct {
	indexer cache.Indexer
}

func NewCoreV1NamespaceLister(indexer cache.Indexer) *CoreV1NamespaceLister {
	return &CoreV1NamespaceLister{indexer: indexer}
}

func (x *CoreV1NamespaceLister) List(selector labels.Selector) ([]*corev1.Namespace, error) {
	var ret []*corev1.Namespace
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.Namespace).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1NamespaceLister) Get(name string) (*corev1.Namespace, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("namespace").GroupResource(), name)
	}
	return obj.(*corev1.Namespace).DeepCopy(), nil
}

type CoreV1NodeLister struct {
	indexer cache.Indexer
}

func NewCoreV1NodeLister(indexer cache.Indexer) *CoreV1NodeLister {
	return &CoreV1NodeLister{indexer: indexer}
}

func (x *CoreV1NodeLister) List(selector labels.Selector) ([]*corev1.Node, error) {
	var ret []*corev1.Node
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.Node).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1NodeLister) Get(name string) (*corev1.Node, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("node").GroupResource(), name)
	}
	return obj.(*corev1.Node).DeepCopy(), nil
}

type CoreV1PersistentVolumeLister struct {
	indexer cache.Indexer
}

func NewCoreV1PersistentVolumeLister(indexer cache.Indexer) *CoreV1PersistentVolumeLister {
	return &CoreV1PersistentVolumeLister{indexer: indexer}
}

func (x *CoreV1PersistentVolumeLister) List(selector labels.Selector) ([]*corev1.PersistentVolume, error) {
	var ret []*corev1.PersistentVolume
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.PersistentVolume).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1PersistentVolumeLister) Get(name string) (*corev1.PersistentVolume, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("persistentvolume").GroupResource(), name)
	}
	return obj.(*corev1.PersistentVolume).DeepCopy(), nil
}

type CoreV1PersistentVolumeClaimLister struct {
	indexer cache.Indexer
}

func NewCoreV1PersistentVolumeClaimLister(indexer cache.Indexer) *CoreV1PersistentVolumeClaimLister {
	return &CoreV1PersistentVolumeClaimLister{indexer: indexer}
}

func (x *CoreV1PersistentVolumeClaimLister) List(namespace string, selector labels.Selector) ([]*corev1.PersistentVolumeClaim, error) {
	var ret []*corev1.PersistentVolumeClaim
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.PersistentVolumeClaim).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1PersistentVolumeClaimLister) Get(namespace, name string) (*corev1.PersistentVolumeClaim, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("persistentvolumeclaim").GroupResource(), name)
	}
	return obj.(*corev1.PersistentVolumeClaim).DeepCopy(), nil
}

type CoreV1PodLister struct {
	indexer cache.Indexer
}

func NewCoreV1PodLister(indexer cache.Indexer) *CoreV1PodLister {
	return &CoreV1PodLister{indexer: indexer}
}

func (x *CoreV1PodLister) List(namespace string, selector labels.Selector) ([]*corev1.Pod, error) {
	var ret []*corev1.Pod
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.Pod).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1PodLister) Get(namespace, name string) (*corev1.Pod, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("pod").GroupResource(), name)
	}
	return obj.(*corev1.Pod).DeepCopy(), nil
}

type CoreV1PodStatusResultLister struct {
	indexer cache.Indexer
}

func NewCoreV1PodStatusResultLister(indexer cache.Indexer) *CoreV1PodStatusResultLister {
	return &CoreV1PodStatusResultLister{indexer: indexer}
}

func (x *CoreV1PodStatusResultLister) List(namespace string, selector labels.Selector) ([]*corev1.PodStatusResult, error) {
	var ret []*corev1.PodStatusResult
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.PodStatusResult).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1PodStatusResultLister) Get(namespace, name string) (*corev1.PodStatusResult, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("podstatusresult").GroupResource(), name)
	}
	return obj.(*corev1.PodStatusResult).DeepCopy(), nil
}

type CoreV1PodTemplateLister struct {
	indexer cache.Indexer
}

func NewCoreV1PodTemplateLister(indexer cache.Indexer) *CoreV1PodTemplateLister {
	return &CoreV1PodTemplateLister{indexer: indexer}
}

func (x *CoreV1PodTemplateLister) List(namespace string, selector labels.Selector) ([]*corev1.PodTemplate, error) {
	var ret []*corev1.PodTemplate
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.PodTemplate).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1PodTemplateLister) Get(namespace, name string) (*corev1.PodTemplate, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("podtemplate").GroupResource(), name)
	}
	return obj.(*corev1.PodTemplate).DeepCopy(), nil
}

type CoreV1RangeAllocationLister struct {
	indexer cache.Indexer
}

func NewCoreV1RangeAllocationLister(indexer cache.Indexer) *CoreV1RangeAllocationLister {
	return &CoreV1RangeAllocationLister{indexer: indexer}
}

func (x *CoreV1RangeAllocationLister) List(namespace string, selector labels.Selector) ([]*corev1.RangeAllocation, error) {
	var ret []*corev1.RangeAllocation
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.RangeAllocation).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1RangeAllocationLister) Get(namespace, name string) (*corev1.RangeAllocation, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("rangeallocation").GroupResource(), name)
	}
	return obj.(*corev1.RangeAllocation).DeepCopy(), nil
}

type CoreV1ReplicationControllerLister struct {
	indexer cache.Indexer
}

func NewCoreV1ReplicationControllerLister(indexer cache.Indexer) *CoreV1ReplicationControllerLister {
	return &CoreV1ReplicationControllerLister{indexer: indexer}
}

func (x *CoreV1ReplicationControllerLister) List(namespace string, selector labels.Selector) ([]*corev1.ReplicationController, error) {
	var ret []*corev1.ReplicationController
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.ReplicationController).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1ReplicationControllerLister) Get(namespace, name string) (*corev1.ReplicationController, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("replicationcontroller").GroupResource(), name)
	}
	return obj.(*corev1.ReplicationController).DeepCopy(), nil
}

type CoreV1ResourceQuotaLister struct {
	indexer cache.Indexer
}

func NewCoreV1ResourceQuotaLister(indexer cache.Indexer) *CoreV1ResourceQuotaLister {
	return &CoreV1ResourceQuotaLister{indexer: indexer}
}

func (x *CoreV1ResourceQuotaLister) List(namespace string, selector labels.Selector) ([]*corev1.ResourceQuota, error) {
	var ret []*corev1.ResourceQuota
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.ResourceQuota).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1ResourceQuotaLister) Get(namespace, name string) (*corev1.ResourceQuota, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("resourcequota").GroupResource(), name)
	}
	return obj.(*corev1.ResourceQuota).DeepCopy(), nil
}

type CoreV1SecretLister struct {
	indexer cache.Indexer
}

func NewCoreV1SecretLister(indexer cache.Indexer) *CoreV1SecretLister {
	return &CoreV1SecretLister{indexer: indexer}
}

func (x *CoreV1SecretLister) List(namespace string, selector labels.Selector) ([]*corev1.Secret, error) {
	var ret []*corev1.Secret
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.Secret).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1SecretLister) Get(namespace, name string) (*corev1.Secret, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("secret").GroupResource(), name)
	}
	return obj.(*corev1.Secret).DeepCopy(), nil
}

type CoreV1ServiceLister struct {
	indexer cache.Indexer
}

func NewCoreV1ServiceLister(indexer cache.Indexer) *CoreV1ServiceLister {
	return &CoreV1ServiceLister{indexer: indexer}
}

func (x *CoreV1ServiceLister) List(namespace string, selector labels.Selector) ([]*corev1.Service, error) {
	var ret []*corev1.Service
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.Service).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1ServiceLister) Get(namespace, name string) (*corev1.Service, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("service").GroupResource(), name)
	}
	return obj.(*corev1.Service).DeepCopy(), nil
}

type CoreV1ServiceAccountLister struct {
	indexer cache.Indexer
}

func NewCoreV1ServiceAccountLister(indexer cache.Indexer) *CoreV1ServiceAccountLister {
	return &CoreV1ServiceAccountLister{indexer: indexer}
}

func (x *CoreV1ServiceAccountLister) List(namespace string, selector labels.Selector) ([]*corev1.ServiceAccount, error) {
	var ret []*corev1.ServiceAccount
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.ServiceAccount).DeepCopy())
	})
	return ret, err
}

func (x *CoreV1ServiceAccountLister) Get(namespace, name string) (*corev1.ServiceAccount, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(corev1.SchemaGroupVersion.WithResource("serviceaccount").GroupResource(), name)
	}
	return obj.(*corev1.ServiceAccount).DeepCopy(), nil
}

type AdmissionregistrationK8sIoV1MutatingWebhookConfigurationLister struct {
	indexer cache.Indexer
}

func NewAdmissionregistrationK8sIoV1MutatingWebhookConfigurationLister(indexer cache.Indexer) *AdmissionregistrationK8sIoV1MutatingWebhookConfigurationLister {
	return &AdmissionregistrationK8sIoV1MutatingWebhookConfigurationLister{indexer: indexer}
}

func (x *AdmissionregistrationK8sIoV1MutatingWebhookConfigurationLister) List(selector labels.Selector) ([]*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	var ret []*admissionregistrationv1.MutatingWebhookConfiguration
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*admissionregistrationv1.MutatingWebhookConfiguration).DeepCopy())
	})
	return ret, err
}

func (x *AdmissionregistrationK8sIoV1MutatingWebhookConfigurationLister) Get(name string) (*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(admissionregistrationv1.SchemaGroupVersion.WithResource("mutatingwebhookconfiguration").GroupResource(), name)
	}
	return obj.(*admissionregistrationv1.MutatingWebhookConfiguration).DeepCopy(), nil
}

type AdmissionregistrationK8sIoV1ValidatingWebhookConfigurationLister struct {
	indexer cache.Indexer
}

func NewAdmissionregistrationK8sIoV1ValidatingWebhookConfigurationLister(indexer cache.Indexer) *AdmissionregistrationK8sIoV1ValidatingWebhookConfigurationLister {
	return &AdmissionregistrationK8sIoV1ValidatingWebhookConfigurationLister{indexer: indexer}
}

func (x *AdmissionregistrationK8sIoV1ValidatingWebhookConfigurationLister) List(selector labels.Selector) ([]*admissionregistrationv1.ValidatingWebhookConfiguration, error) {
	var ret []*admissionregistrationv1.ValidatingWebhookConfiguration
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*admissionregistrationv1.ValidatingWebhookConfiguration).DeepCopy())
	})
	return ret, err
}

func (x *AdmissionregistrationK8sIoV1ValidatingWebhookConfigurationLister) Get(name string) (*admissionregistrationv1.ValidatingWebhookConfiguration, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(admissionregistrationv1.SchemaGroupVersion.WithResource("validatingwebhookconfiguration").GroupResource(), name)
	}
	return obj.(*admissionregistrationv1.ValidatingWebhookConfiguration).DeepCopy(), nil
}

type AppsV1ControllerRevisionLister struct {
	indexer cache.Indexer
}

func NewAppsV1ControllerRevisionLister(indexer cache.Indexer) *AppsV1ControllerRevisionLister {
	return &AppsV1ControllerRevisionLister{indexer: indexer}
}

func (x *AppsV1ControllerRevisionLister) List(namespace string, selector labels.Selector) ([]*appsv1.ControllerRevision, error) {
	var ret []*appsv1.ControllerRevision
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*appsv1.ControllerRevision).DeepCopy())
	})
	return ret, err
}

func (x *AppsV1ControllerRevisionLister) Get(namespace, name string) (*appsv1.ControllerRevision, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(appsv1.SchemaGroupVersion.WithResource("controllerrevision").GroupResource(), name)
	}
	return obj.(*appsv1.ControllerRevision).DeepCopy(), nil
}

type AppsV1CronJobLister struct {
	indexer cache.Indexer
}

func NewAppsV1CronJobLister(indexer cache.Indexer) *AppsV1CronJobLister {
	return &AppsV1CronJobLister{indexer: indexer}
}

func (x *AppsV1CronJobLister) List(namespace string, selector labels.Selector) ([]*batchv1.CronJob, error) {
	var ret []*batchv1.CronJob
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*batchv1.CronJob).DeepCopy())
	})
	return ret, err
}

func (x *AppsV1CronJobLister) Get(namespace, name string) (*batchv1.CronJob, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(batchv1.SchemaGroupVersion.WithResource("cronjob").GroupResource(), name)
	}
	return obj.(*batchv1.CronJob).DeepCopy(), nil
}

type AppsV1DaemonSetLister struct {
	indexer cache.Indexer
}

func NewAppsV1DaemonSetLister(indexer cache.Indexer) *AppsV1DaemonSetLister {
	return &AppsV1DaemonSetLister{indexer: indexer}
}

func (x *AppsV1DaemonSetLister) List(namespace string, selector labels.Selector) ([]*appsv1.DaemonSet, error) {
	var ret []*appsv1.DaemonSet
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*appsv1.DaemonSet).DeepCopy())
	})
	return ret, err
}

func (x *AppsV1DaemonSetLister) Get(namespace, name string) (*appsv1.DaemonSet, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(appsv1.SchemaGroupVersion.WithResource("daemonset").GroupResource(), name)
	}
	return obj.(*appsv1.DaemonSet).DeepCopy(), nil
}

type AppsV1DeploymentLister struct {
	indexer cache.Indexer
}

func NewAppsV1DeploymentLister(indexer cache.Indexer) *AppsV1DeploymentLister {
	return &AppsV1DeploymentLister{indexer: indexer}
}

func (x *AppsV1DeploymentLister) List(namespace string, selector labels.Selector) ([]*appsv1.Deployment, error) {
	var ret []*appsv1.Deployment
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*appsv1.Deployment).DeepCopy())
	})
	return ret, err
}

func (x *AppsV1DeploymentLister) Get(namespace, name string) (*appsv1.Deployment, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(appsv1.SchemaGroupVersion.WithResource("deployment").GroupResource(), name)
	}
	return obj.(*appsv1.Deployment).DeepCopy(), nil
}

type AppsV1JobLister struct {
	indexer cache.Indexer
}

func NewAppsV1JobLister(indexer cache.Indexer) *AppsV1JobLister {
	return &AppsV1JobLister{indexer: indexer}
}

func (x *AppsV1JobLister) List(namespace string, selector labels.Selector) ([]*batchv1.Job, error) {
	var ret []*batchv1.Job
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*batchv1.Job).DeepCopy())
	})
	return ret, err
}

func (x *AppsV1JobLister) Get(namespace, name string) (*batchv1.Job, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(batchv1.SchemaGroupVersion.WithResource("job").GroupResource(), name)
	}
	return obj.(*batchv1.Job).DeepCopy(), nil
}

type AppsV1ReplicaSetLister struct {
	indexer cache.Indexer
}

func NewAppsV1ReplicaSetLister(indexer cache.Indexer) *AppsV1ReplicaSetLister {
	return &AppsV1ReplicaSetLister{indexer: indexer}
}

func (x *AppsV1ReplicaSetLister) List(namespace string, selector labels.Selector) ([]*appsv1.ReplicaSet, error) {
	var ret []*appsv1.ReplicaSet
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*appsv1.ReplicaSet).DeepCopy())
	})
	return ret, err
}

func (x *AppsV1ReplicaSetLister) Get(namespace, name string) (*appsv1.ReplicaSet, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(appsv1.SchemaGroupVersion.WithResource("replicaset").GroupResource(), name)
	}
	return obj.(*appsv1.ReplicaSet).DeepCopy(), nil
}

type AppsV1StatefulSetLister struct {
	indexer cache.Indexer
}

func NewAppsV1StatefulSetLister(indexer cache.Indexer) *AppsV1StatefulSetLister {
	return &AppsV1StatefulSetLister{indexer: indexer}
}

func (x *AppsV1StatefulSetLister) List(namespace string, selector labels.Selector) ([]*appsv1.StatefulSet, error) {
	var ret []*appsv1.StatefulSet
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*appsv1.StatefulSet).DeepCopy())
	})
	return ret, err
}

func (x *AppsV1StatefulSetLister) Get(namespace, name string) (*appsv1.StatefulSet, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(appsv1.SchemaGroupVersion.WithResource("statefulset").GroupResource(), name)
	}
	return obj.(*appsv1.StatefulSet).DeepCopy(), nil
}

type AuthenticationK8sIoV1TokenRequestLister struct {
	indexer cache.Indexer
}

func NewAuthenticationK8sIoV1TokenRequestLister(indexer cache.Indexer) *AuthenticationK8sIoV1TokenRequestLister {
	return &AuthenticationK8sIoV1TokenRequestLister{indexer: indexer}
}

func (x *AuthenticationK8sIoV1TokenRequestLister) List(namespace string, selector labels.Selector) ([]*authenticationv1.TokenRequest, error) {
	var ret []*authenticationv1.TokenRequest
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*authenticationv1.TokenRequest).DeepCopy())
	})
	return ret, err
}

func (x *AuthenticationK8sIoV1TokenRequestLister) Get(namespace, name string) (*authenticationv1.TokenRequest, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(authenticationv1.SchemaGroupVersion.WithResource("tokenrequest").GroupResource(), name)
	}
	return obj.(*authenticationv1.TokenRequest).DeepCopy(), nil
}

type AuthenticationK8sIoV1TokenReviewLister struct {
	indexer cache.Indexer
}

func NewAuthenticationK8sIoV1TokenReviewLister(indexer cache.Indexer) *AuthenticationK8sIoV1TokenReviewLister {
	return &AuthenticationK8sIoV1TokenReviewLister{indexer: indexer}
}

func (x *AuthenticationK8sIoV1TokenReviewLister) List(selector labels.Selector) ([]*authenticationv1.TokenReview, error) {
	var ret []*authenticationv1.TokenReview
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*authenticationv1.TokenReview).DeepCopy())
	})
	return ret, err
}

func (x *AuthenticationK8sIoV1TokenReviewLister) Get(name string) (*authenticationv1.TokenReview, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(authenticationv1.SchemaGroupVersion.WithResource("tokenreview").GroupResource(), name)
	}
	return obj.(*authenticationv1.TokenReview).DeepCopy(), nil
}

type AuthorizationK8sIoV1LocalSubjectAccessReviewLister struct {
	indexer cache.Indexer
}

func NewAuthorizationK8sIoV1LocalSubjectAccessReviewLister(indexer cache.Indexer) *AuthorizationK8sIoV1LocalSubjectAccessReviewLister {
	return &AuthorizationK8sIoV1LocalSubjectAccessReviewLister{indexer: indexer}
}

func (x *AuthorizationK8sIoV1LocalSubjectAccessReviewLister) List(namespace string, selector labels.Selector) ([]*authorizationv1.LocalSubjectAccessReview, error) {
	var ret []*authorizationv1.LocalSubjectAccessReview
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*authorizationv1.LocalSubjectAccessReview).DeepCopy())
	})
	return ret, err
}

func (x *AuthorizationK8sIoV1LocalSubjectAccessReviewLister) Get(namespace, name string) (*authorizationv1.LocalSubjectAccessReview, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(authorizationv1.SchemaGroupVersion.WithResource("localsubjectaccessreview").GroupResource(), name)
	}
	return obj.(*authorizationv1.LocalSubjectAccessReview).DeepCopy(), nil
}

type AuthorizationK8sIoV1SelfSubjectAccessReviewLister struct {
	indexer cache.Indexer
}

func NewAuthorizationK8sIoV1SelfSubjectAccessReviewLister(indexer cache.Indexer) *AuthorizationK8sIoV1SelfSubjectAccessReviewLister {
	return &AuthorizationK8sIoV1SelfSubjectAccessReviewLister{indexer: indexer}
}

func (x *AuthorizationK8sIoV1SelfSubjectAccessReviewLister) List(selector labels.Selector) ([]*authorizationv1.SelfSubjectAccessReview, error) {
	var ret []*authorizationv1.SelfSubjectAccessReview
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*authorizationv1.SelfSubjectAccessReview).DeepCopy())
	})
	return ret, err
}

func (x *AuthorizationK8sIoV1SelfSubjectAccessReviewLister) Get(name string) (*authorizationv1.SelfSubjectAccessReview, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(authorizationv1.SchemaGroupVersion.WithResource("selfsubjectaccessreview").GroupResource(), name)
	}
	return obj.(*authorizationv1.SelfSubjectAccessReview).DeepCopy(), nil
}

type AuthorizationK8sIoV1SelfSubjectRulesReviewLister struct {
	indexer cache.Indexer
}

func NewAuthorizationK8sIoV1SelfSubjectRulesReviewLister(indexer cache.Indexer) *AuthorizationK8sIoV1SelfSubjectRulesReviewLister {
	return &AuthorizationK8sIoV1SelfSubjectRulesReviewLister{indexer: indexer}
}

func (x *AuthorizationK8sIoV1SelfSubjectRulesReviewLister) List(selector labels.Selector) ([]*authorizationv1.SelfSubjectRulesReview, error) {
	var ret []*authorizationv1.SelfSubjectRulesReview
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*authorizationv1.SelfSubjectRulesReview).DeepCopy())
	})
	return ret, err
}

func (x *AuthorizationK8sIoV1SelfSubjectRulesReviewLister) Get(name string) (*authorizationv1.SelfSubjectRulesReview, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(authorizationv1.SchemaGroupVersion.WithResource("selfsubjectrulesreview").GroupResource(), name)
	}
	return obj.(*authorizationv1.SelfSubjectRulesReview).DeepCopy(), nil
}

type AuthorizationK8sIoV1SubjectAccessReviewLister struct {
	indexer cache.Indexer
}

func NewAuthorizationK8sIoV1SubjectAccessReviewLister(indexer cache.Indexer) *AuthorizationK8sIoV1SubjectAccessReviewLister {
	return &AuthorizationK8sIoV1SubjectAccessReviewLister{indexer: indexer}
}

func (x *AuthorizationK8sIoV1SubjectAccessReviewLister) List(selector labels.Selector) ([]*authorizationv1.SubjectAccessReview, error) {
	var ret []*authorizationv1.SubjectAccessReview
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*authorizationv1.SubjectAccessReview).DeepCopy())
	})
	return ret, err
}

func (x *AuthorizationK8sIoV1SubjectAccessReviewLister) Get(name string) (*authorizationv1.SubjectAccessReview, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(authorizationv1.SchemaGroupVersion.WithResource("subjectaccessreview").GroupResource(), name)
	}
	return obj.(*authorizationv1.SubjectAccessReview).DeepCopy(), nil
}

type AutoscalingV1HorizontalPodAutoscalerLister struct {
	indexer cache.Indexer
}

func NewAutoscalingV1HorizontalPodAutoscalerLister(indexer cache.Indexer) *AutoscalingV1HorizontalPodAutoscalerLister {
	return &AutoscalingV1HorizontalPodAutoscalerLister{indexer: indexer}
}

func (x *AutoscalingV1HorizontalPodAutoscalerLister) List(namespace string, selector labels.Selector) ([]*autoscalingv1.HorizontalPodAutoscaler, error) {
	var ret []*autoscalingv1.HorizontalPodAutoscaler
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*autoscalingv1.HorizontalPodAutoscaler).DeepCopy())
	})
	return ret, err
}

func (x *AutoscalingV1HorizontalPodAutoscalerLister) Get(namespace, name string) (*autoscalingv1.HorizontalPodAutoscaler, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(autoscalingv1.SchemaGroupVersion.WithResource("horizontalpodautoscaler").GroupResource(), name)
	}
	return obj.(*autoscalingv1.HorizontalPodAutoscaler).DeepCopy(), nil
}

type AutoscalingV1ScaleLister struct {
	indexer cache.Indexer
}

func NewAutoscalingV1ScaleLister(indexer cache.Indexer) *AutoscalingV1ScaleLister {
	return &AutoscalingV1ScaleLister{indexer: indexer}
}

func (x *AutoscalingV1ScaleLister) List(namespace string, selector labels.Selector) ([]*autoscalingv1.Scale, error) {
	var ret []*autoscalingv1.Scale
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*autoscalingv1.Scale).DeepCopy())
	})
	return ret, err
}

func (x *AutoscalingV1ScaleLister) Get(namespace, name string) (*autoscalingv1.Scale, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(autoscalingv1.SchemaGroupVersion.WithResource("scale").GroupResource(), name)
	}
	return obj.(*autoscalingv1.Scale).DeepCopy(), nil
}

type CertificatesK8sIoV1CertificateSigningRequestLister struct {
	indexer cache.Indexer
}

func NewCertificatesK8sIoV1CertificateSigningRequestLister(indexer cache.Indexer) *CertificatesK8sIoV1CertificateSigningRequestLister {
	return &CertificatesK8sIoV1CertificateSigningRequestLister{indexer: indexer}
}

func (x *CertificatesK8sIoV1CertificateSigningRequestLister) List(selector labels.Selector) ([]*certificatesv1.CertificateSigningRequest, error) {
	var ret []*certificatesv1.CertificateSigningRequest
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*certificatesv1.CertificateSigningRequest).DeepCopy())
	})
	return ret, err
}

func (x *CertificatesK8sIoV1CertificateSigningRequestLister) Get(name string) (*certificatesv1.CertificateSigningRequest, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(certificatesv1.SchemaGroupVersion.WithResource("certificatesigningrequest").GroupResource(), name)
	}
	return obj.(*certificatesv1.CertificateSigningRequest).DeepCopy(), nil
}

type DiscoveryK8sIoV1EndpointSliceLister struct {
	indexer cache.Indexer
}

func NewDiscoveryK8sIoV1EndpointSliceLister(indexer cache.Indexer) *DiscoveryK8sIoV1EndpointSliceLister {
	return &DiscoveryK8sIoV1EndpointSliceLister{indexer: indexer}
}

func (x *DiscoveryK8sIoV1EndpointSliceLister) List(namespace string, selector labels.Selector) ([]*discoveryv1.EndpointSlice, error) {
	var ret []*discoveryv1.EndpointSlice
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*discoveryv1.EndpointSlice).DeepCopy())
	})
	return ret, err
}

func (x *DiscoveryK8sIoV1EndpointSliceLister) Get(namespace, name string) (*discoveryv1.EndpointSlice, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(discoveryv1.SchemaGroupVersion.WithResource("endpointslice").GroupResource(), name)
	}
	return obj.(*discoveryv1.EndpointSlice).DeepCopy(), nil
}

type NetworkingK8sIoV1IngressLister struct {
	indexer cache.Indexer
}

func NewNetworkingK8sIoV1IngressLister(indexer cache.Indexer) *NetworkingK8sIoV1IngressLister {
	return &NetworkingK8sIoV1IngressLister{indexer: indexer}
}

func (x *NetworkingK8sIoV1IngressLister) List(namespace string, selector labels.Selector) ([]*networkingv1.Ingress, error) {
	var ret []*networkingv1.Ingress
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*networkingv1.Ingress).DeepCopy())
	})
	return ret, err
}

func (x *NetworkingK8sIoV1IngressLister) Get(namespace, name string) (*networkingv1.Ingress, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(networkingv1.SchemaGroupVersion.WithResource("ingress").GroupResource(), name)
	}
	return obj.(*networkingv1.Ingress).DeepCopy(), nil
}

type NetworkingK8sIoV1IngressClassLister struct {
	indexer cache.Indexer
}

func NewNetworkingK8sIoV1IngressClassLister(indexer cache.Indexer) *NetworkingK8sIoV1IngressClassLister {
	return &NetworkingK8sIoV1IngressClassLister{indexer: indexer}
}

func (x *NetworkingK8sIoV1IngressClassLister) List(selector labels.Selector) ([]*networkingv1.IngressClass, error) {
	var ret []*networkingv1.IngressClass
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*networkingv1.IngressClass).DeepCopy())
	})
	return ret, err
}

func (x *NetworkingK8sIoV1IngressClassLister) Get(name string) (*networkingv1.IngressClass, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(networkingv1.SchemaGroupVersion.WithResource("ingressclass").GroupResource(), name)
	}
	return obj.(*networkingv1.IngressClass).DeepCopy(), nil
}

type NetworkingK8sIoV1NetworkPolicyLister struct {
	indexer cache.Indexer
}

func NewNetworkingK8sIoV1NetworkPolicyLister(indexer cache.Indexer) *NetworkingK8sIoV1NetworkPolicyLister {
	return &NetworkingK8sIoV1NetworkPolicyLister{indexer: indexer}
}

func (x *NetworkingK8sIoV1NetworkPolicyLister) List(namespace string, selector labels.Selector) ([]*networkingv1.NetworkPolicy, error) {
	var ret []*networkingv1.NetworkPolicy
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*networkingv1.NetworkPolicy).DeepCopy())
	})
	return ret, err
}

func (x *NetworkingK8sIoV1NetworkPolicyLister) Get(namespace, name string) (*networkingv1.NetworkPolicy, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(networkingv1.SchemaGroupVersion.WithResource("networkpolicy").GroupResource(), name)
	}
	return obj.(*networkingv1.NetworkPolicy).DeepCopy(), nil
}

type PolicyV1EvictionLister struct {
	indexer cache.Indexer
}

func NewPolicyV1EvictionLister(indexer cache.Indexer) *PolicyV1EvictionLister {
	return &PolicyV1EvictionLister{indexer: indexer}
}

func (x *PolicyV1EvictionLister) List(namespace string, selector labels.Selector) ([]*policyv1.Eviction, error) {
	var ret []*policyv1.Eviction
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*policyv1.Eviction).DeepCopy())
	})
	return ret, err
}

func (x *PolicyV1EvictionLister) Get(namespace, name string) (*policyv1.Eviction, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(policyv1.SchemaGroupVersion.WithResource("eviction").GroupResource(), name)
	}
	return obj.(*policyv1.Eviction).DeepCopy(), nil
}

type PolicyV1PodDisruptionBudgetLister struct {
	indexer cache.Indexer
}

func NewPolicyV1PodDisruptionBudgetLister(indexer cache.Indexer) *PolicyV1PodDisruptionBudgetLister {
	return &PolicyV1PodDisruptionBudgetLister{indexer: indexer}
}

func (x *PolicyV1PodDisruptionBudgetLister) List(namespace string, selector labels.Selector) ([]*policyv1.PodDisruptionBudget, error) {
	var ret []*policyv1.PodDisruptionBudget
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*policyv1.PodDisruptionBudget).DeepCopy())
	})
	return ret, err
}

func (x *PolicyV1PodDisruptionBudgetLister) Get(namespace, name string) (*policyv1.PodDisruptionBudget, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(policyv1.SchemaGroupVersion.WithResource("poddisruptionbudget").GroupResource(), name)
	}
	return obj.(*policyv1.PodDisruptionBudget).DeepCopy(), nil
}

type RbacAuthorizationK8sIoV1ClusterRoleLister struct {
	indexer cache.Indexer
}

func NewRbacAuthorizationK8sIoV1ClusterRoleLister(indexer cache.Indexer) *RbacAuthorizationK8sIoV1ClusterRoleLister {
	return &RbacAuthorizationK8sIoV1ClusterRoleLister{indexer: indexer}
}

func (x *RbacAuthorizationK8sIoV1ClusterRoleLister) List(selector labels.Selector) ([]*rbacv1.ClusterRole, error) {
	var ret []*rbacv1.ClusterRole
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*rbacv1.ClusterRole).DeepCopy())
	})
	return ret, err
}

func (x *RbacAuthorizationK8sIoV1ClusterRoleLister) Get(name string) (*rbacv1.ClusterRole, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(rbacv1.SchemaGroupVersion.WithResource("clusterrole").GroupResource(), name)
	}
	return obj.(*rbacv1.ClusterRole).DeepCopy(), nil
}

type RbacAuthorizationK8sIoV1ClusterRoleBindingLister struct {
	indexer cache.Indexer
}

func NewRbacAuthorizationK8sIoV1ClusterRoleBindingLister(indexer cache.Indexer) *RbacAuthorizationK8sIoV1ClusterRoleBindingLister {
	return &RbacAuthorizationK8sIoV1ClusterRoleBindingLister{indexer: indexer}
}

func (x *RbacAuthorizationK8sIoV1ClusterRoleBindingLister) List(selector labels.Selector) ([]*rbacv1.ClusterRoleBinding, error) {
	var ret []*rbacv1.ClusterRoleBinding
	err := cache.ListAll(x.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*rbacv1.ClusterRoleBinding).DeepCopy())
	})
	return ret, err
}

func (x *RbacAuthorizationK8sIoV1ClusterRoleBindingLister) Get(name string) (*rbacv1.ClusterRoleBinding, error) {
	obj, exists, err := x.indexer.GetByKey("/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(rbacv1.SchemaGroupVersion.WithResource("clusterrolebinding").GroupResource(), name)
	}
	return obj.(*rbacv1.ClusterRoleBinding).DeepCopy(), nil
}

type RbacAuthorizationK8sIoV1RoleLister struct {
	indexer cache.Indexer
}

func NewRbacAuthorizationK8sIoV1RoleLister(indexer cache.Indexer) *RbacAuthorizationK8sIoV1RoleLister {
	return &RbacAuthorizationK8sIoV1RoleLister{indexer: indexer}
}

func (x *RbacAuthorizationK8sIoV1RoleLister) List(namespace string, selector labels.Selector) ([]*rbacv1.Role, error) {
	var ret []*rbacv1.Role
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*rbacv1.Role).DeepCopy())
	})
	return ret, err
}

func (x *RbacAuthorizationK8sIoV1RoleLister) Get(namespace, name string) (*rbacv1.Role, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(rbacv1.SchemaGroupVersion.WithResource("role").GroupResource(), name)
	}
	return obj.(*rbacv1.Role).DeepCopy(), nil
}

type RbacAuthorizationK8sIoV1RoleBindingLister struct {
	indexer cache.Indexer
}

func NewRbacAuthorizationK8sIoV1RoleBindingLister(indexer cache.Indexer) *RbacAuthorizationK8sIoV1RoleBindingLister {
	return &RbacAuthorizationK8sIoV1RoleBindingLister{indexer: indexer}
}

func (x *RbacAuthorizationK8sIoV1RoleBindingLister) List(namespace string, selector labels.Selector) ([]*rbacv1.RoleBinding, error) {
	var ret []*rbacv1.RoleBinding
	err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*rbacv1.RoleBinding).DeepCopy())
	})
	return ret, err
}

func (x *RbacAuthorizationK8sIoV1RoleBindingLister) Get(namespace, name string) (*rbacv1.RoleBinding, error) {
	obj, exists, err := x.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, k8serrors.NewNotFound(rbacv1.SchemaGroupVersion.WithResource("rolebinding").GroupResource(), name)
	}
	return obj.(*rbacv1.RoleBinding).DeepCopy(), nil
}
