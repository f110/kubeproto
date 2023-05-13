package k8s

import (
	"fmt"
	"io"
	"path"
	"sort"
	"strings"

	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	"go.f110.dev/kubeproto/internal/codegeneration"
	"go.f110.dev/kubeproto/internal/definition"
	"go.f110.dev/kubeproto/internal/stringsutil"
)

type ClientGenerator struct {
	files                   []*descriptorpb.FileDescriptorProto
	lister                  *definition.Lister
	packageNamespaceManager *definition.PackageNamespaceManager
}

func NewClientGenerator(fileToGenerate []string, files *protoregistry.Files) *ClientGenerator {
	nsm := definition.NewPackageNamespaceManager()
	return &ClientGenerator{
		files:                   nil,
		lister:                  definition.NewLister(fileToGenerate, files, nsm),
		packageNamespaceManager: nsm,
	}
}

func (g *ClientGenerator) Generate(out io.Writer, packageName, importPath string) error {
	w := codegeneration.NewWriter()
	w.F("package %s", path.Base(packageName))

	// The key is a package path. The value is an alias.
	importPackages := map[string]string{
		"errors":                                 "",
		"context":                                "",
		"time":                                   "",
		"k8s.io/apimachinery/pkg/runtime":        "",
		"k8s.io/apimachinery/pkg/runtime/schema": "",
		"k8s.io/apimachinery/pkg/runtime/serializer": "",
	}
	for k, v := range importPackages {
		g.packageNamespaceManager.Add(k, v)
	}

	messages := g.lister.GetMessages()
	kinds := make(map[string][]*definition.Message)
	groupVersions := make(map[string][]*definition.Message)
	for _, m := range messages.FilterKind() {
		if m.Virtual {
			continue
		}
		kinds[m.ShortName] = append(kinds[m.ShortName], m)
		gv := fmt.Sprintf("%s/%s", m.Group, m.Version)
		groupVersions[gv] = append(groupVersions[gv], m)
	}

	writer := codegeneration.NewWriter()
	writer.F("var (")
	writer.F("Scheme = runtime.NewScheme()")
	writer.F("ParameterCodec = runtime.NewParameterCodec(Scheme)")
	writer.F("Codecs = serializer.NewCodecFactory(Scheme)")
	writer.F("AddToScheme = localSchemeBuilder.AddToScheme")
	writer.F(")")
	writer.F("")

	writer.F("var localSchemeBuilder = runtime.SchemeBuilder{")
	for _, key := range keys(groupVersions) {
		v := groupVersions[key]
		m := v[0]
		writer.F("%s.AddToScheme,", path.Base(m.Package.Alias))
	}
	writer.F("}")

	writer.F("func init() {")
	writer.F("for _, v := range []func(*runtime.Scheme) error{")
	for _, key := range keys(groupVersions) {
		v := groupVersions[key]
		m := v[0]
		writer.F("%s.AddToScheme,", path.Base(m.Package.Alias))
	}
	writer.F("} {")
	writer.F("if err := v(Scheme); err != nil {\npanic(err)\n}")
	writer.F("}") // end of for
	writer.F("}") // end of init()
	writer.F("")

	writer.F(`
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
}`)

	writer.F("type Set struct {")
	for _, key := range keys(groupVersions) {
		m := groupVersions[key][0]
		clientName := m.ClientName()
		writer.F("%s *%s", clientName, clientName)
	}
	writer.F("}")
	writer.F("")
	writer.F("func NewSet(cfg *rest.Config) (*Set,error) {")
	writer.F("s := &Set{}")
	for _, key := range keys(groupVersions) {
		m := groupVersions[key][0]
		clientName := m.ClientName()
		writer.F("{")
		writer.F("conf := *cfg")
		writer.F("conf.GroupVersion = &%s.SchemaGroupVersion", m.Package.Alias)
		writer.F("conf.APIPath = \"/apis\"")
		writer.F("conf.NegotiatedSerializer = Codecs.WithoutConversion()")
		writer.F("c, err := rest.RESTClientFor(&conf)")
		writer.F("if err != nil {")
		writer.F("return nil, err")
		writer.F("}")
		writer.F("s.%s = New%sClient(&restBackend{client: c})", clientName, clientName)
		writer.F("}")
	}
	writer.F("")
	writer.F("return s, nil")
	writer.F("}") // end of NewSet
	writer.F("")

	restClient := newRestClientGenerator(groupVersions)
	if err := restClient.WriteTo(writer); err != nil {
		return err
	}
	for p, a := range restClient.Import() {
		importPackages[p] = a
	}
	informer := newInformerGenerator(groupVersions)
	if err := informer.WriteTo(writer); err != nil {
		return err
	}
	for p, a := range informer.Import() {
		importPackages[p] = a
	}
	lister := newListerGenerator(groupVersions)
	if err := lister.WriteTo(writer); err != nil {
		return err
	}
	for p, a := range lister.Import() {
		importPackages[p] = a
	}

	w.F("import (")
	core, libs, proj := sortImports(importPackages, importPath)
	for _, v := range []map[string]string{core, libs, proj} {
		for p, a := range v {
			if a != "" {
				w.F("%s %q", a, p)
			} else {
				w.F("%q", p)
			}
		}
		w.F("")
	}
	w.F(")")
	writer.WriteTo(w)

	if err := w.Format(); err != nil {
		return err
	}
	if _, err := w.WriteTo(out); err != nil {
		return err
	}
	return nil
}

func sortImports(v map[string]string, projectPackageName string) (core map[string]string, libs map[string]string, proj map[string]string) {
	core = make(map[string]string)
	libs = make(map[string]string)
	proj = make(map[string]string)

	s := strings.Split(projectPackageName, "/")
	projPkg := fmt.Sprintf("%s/%s", s[0], s[1])

	for p, a := range v {
		if isCorePackage(p) {
			core[p] = a
			continue
		}
		if isProjectPackage(p, projPkg) {
			proj[p] = a
			continue
		}
		libs[p] = a
	}

	return
}

func isCorePackage(v string) bool {
	if strings.Contains(v, ".") {
		return false
	}
	return true
}

func isProjectPackage(v, projPkg string) bool {
	if strings.HasPrefix(v, projPkg) {
		return true
	}
	return false
}

type restClientGenerator struct {
	groupVersions map[string][]*definition.Message
}

func newRestClientGenerator(groupVersions map[string][]*definition.Message) *restClientGenerator {
	return &restClientGenerator{groupVersions: groupVersions}
}

func (g *restClientGenerator) Import() map[string]string {
	importPackages := map[string]string{
		"k8s.io/client-go/rest":                "",
		"k8s.io/apimachinery/pkg/apis/meta/v1": "metav1",
		"k8s.io/apimachinery/pkg/watch":        "",
	}
	for _, v := range g.groupVersions {
		for _, m := range v {
			_, p := path.Split(m.Package.Path)
			alias := m.Package.Alias
			if p == m.Package.Alias {
				alias = ""
			}
			importPackages[m.Package.Path] = alias
		}
	}

	return importPackages
}

func (g *restClientGenerator) WriteTo(writer *codegeneration.Writer) error {
	writer.F(`type restBackend struct {
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
`)

	for _, k := range keys(g.groupVersions) {
		v := g.groupVersions[k]
		m := v[0]
		clientName := m.ClientName()
		writer.F(`type %s struct {
	backend Backend
}`, clientName)
		writer.F("")

		writer.F("func New%sClient(b Backend) *%s {", clientName, clientName)
		writer.F("return &%s{backend: b}", clientName)
		writer.F("}")
		writer.F("")

		for _, m := range v {
			structNameWithPkg := fmt.Sprintf("%s.%s", m.Package.Alias, m.ShortName)
			// GetXXX
			if m.Scope == definition.ScopeTypeCluster {
				writer.F("func(c *%s) Get%s(ctx context.Context, name string, opts metav1.GetOptions) (*%s, error) {", clientName, m.ShortName, structNameWithPkg)
				writer.F("result, err := c.backend.GetClusterScoped(ctx, %q, %q, name, opts, &%s{})", strings.ToLower(stringsutil.Plural(m.ShortName)), m.ShortName, structNameWithPkg)
				writer.F("if err != nil {")
				writer.F("return nil, err")
				writer.F("}")
				writer.F("return result.(*%s), nil", structNameWithPkg)
				writer.F("}")
				writer.F("")
			} else {
				writer.F("func(c *%s) Get%s(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*%s, error) {", clientName, m.ShortName, structNameWithPkg)
				writer.F("result, err := c.backend.Get(ctx, %q, %q, namespace, name, opts, &%s{})", strings.ToLower(stringsutil.Plural(m.ShortName)), m.ShortName, structNameWithPkg)
				writer.F("if err != nil {")
				writer.F("return nil, err")
				writer.F("}")
				writer.F("return result.(*%s), nil", structNameWithPkg)
				writer.F("}")
				writer.F("")
			}

			// CreateXXX
			if m.Scope == definition.ScopeTypeCluster {
				writer.F("func (c *%s) Create%s(ctx context.Context, v *%s, opts metav1.CreateOptions) (*%s, error) {", clientName, m.ShortName, structNameWithPkg, structNameWithPkg)
				writer.F("result, err := c.backend.CreateClusterScoped(ctx, %q, %q, v, opts, &%s{})", strings.ToLower(stringsutil.Plural(m.ShortName)), m.ShortName, structNameWithPkg)
				writer.F("if err != nil {")
				writer.F("return nil, err")
				writer.F("}")
				writer.F("return result.(*%s), nil", structNameWithPkg)
				writer.F("}")
				writer.F("")
			} else {
				writer.F("func (c *%s) Create%s(ctx context.Context, v *%s, opts metav1.CreateOptions) (*%s, error) {", clientName, m.ShortName, structNameWithPkg, structNameWithPkg)
				writer.F("result, err := c.backend.Create(ctx, %q, %q, v, opts, &%s{})", strings.ToLower(stringsutil.Plural(m.ShortName)), m.ShortName, structNameWithPkg)
				writer.F("if err != nil {")
				writer.F("return nil, err")
				writer.F("}")
				writer.F("return result.(*%s), nil", structNameWithPkg)
				writer.F("}")
				writer.F("")
			}

			// UpdateXXX
			if m.Scope == definition.ScopeTypeCluster {
				writer.F("func (c *%s) Update%s(ctx context.Context, v *%s, opts metav1.UpdateOptions) (*%s, error) {", clientName, m.ShortName, structNameWithPkg, structNameWithPkg)
				writer.F("result, err := c.backend.UpdateClusterScoped(ctx, %q, %q, v, opts, &%s{})", strings.ToLower(stringsutil.Plural(m.ShortName)), m.ShortName, structNameWithPkg)
				writer.F("if err != nil {")
				writer.F("return nil, err")
				writer.F("}")
				writer.F("return result.(*%s), nil", structNameWithPkg)
				writer.F("}")
				writer.F("")
			} else {
				writer.F("func (c *%s) Update%s(ctx context.Context, v *%s, opts metav1.UpdateOptions) (*%s, error) {", clientName, m.ShortName, structNameWithPkg, structNameWithPkg)
				writer.F("result, err := c.backend.Update(ctx, %q, %q, v, opts, &%s{})", strings.ToLower(stringsutil.Plural(m.ShortName)), m.ShortName, structNameWithPkg)
				writer.F("if err != nil {")
				writer.F("return nil, err")
				writer.F("}")
				writer.F("return result.(*%s), nil", structNameWithPkg)
				writer.F("}")
				writer.F("")
			}

			// UpdateStatusXXX
			if m.IsDefinedSubResource() {
				if m.Scope == definition.ScopeTypeCluster {
					writer.F("func (c *%s) UpdateStatus%s(ctx context.Context, v *%s, opts metav1.UpdateOptions) (*%s, error) {", clientName, m.ShortName, structNameWithPkg, structNameWithPkg)
					writer.F("result, err := c.backend.UpdateStatusClusterScoped(ctx, %q, %q, v, opts, &%s{})", strings.ToLower(stringsutil.Plural(m.ShortName)), m.ShortName, structNameWithPkg)
					writer.F("if err != nil {")
					writer.F("return nil, err")
					writer.F("}")
					writer.F("return result.(*%s), nil", structNameWithPkg)
					writer.F("}")
					writer.F("")
				} else {
					writer.F("func (c *%s) UpdateStatus%s(ctx context.Context, v *%s, opts metav1.UpdateOptions) (*%s, error) {", clientName, m.ShortName, structNameWithPkg, structNameWithPkg)
					writer.F("result, err := c.backend.UpdateStatus(ctx, %q, %q, v, opts, &%s{})", strings.ToLower(stringsutil.Plural(m.ShortName)), m.ShortName, structNameWithPkg)
					writer.F("if err != nil {")
					writer.F("return nil, err")
					writer.F("}")
					writer.F("return result.(*%s), nil", structNameWithPkg)
					writer.F("}")
					writer.F("")
				}
			}

			// DeleteXXX
			if m.Scope == definition.ScopeTypeCluster {
				writer.F("func (c *%s) Delete%s(ctx context.Context, name string, opts metav1.DeleteOptions) error {", clientName, m.ShortName)
				writer.F("return c.backend.DeleteClusterScoped(ctx, schema.GroupVersionResource{Group:%q, Version:%q, Resource:%q}, name, opts)", m.Group, m.Version, strings.ToLower(stringsutil.Plural(m.ShortName)))
				writer.F("}")
				writer.F("")
			} else {
				writer.F("func (c *%s) Delete%s(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {", clientName, m.ShortName)
				writer.F("return c.backend.Delete(ctx, schema.GroupVersionResource{Group:%q, Version:%q, Resource:%q}, namespace, name, opts)", m.Group, m.Version, strings.ToLower(stringsutil.Plural(m.ShortName)))
				writer.F("}")
				writer.F("")
			}

			// ListXXX
			if m.Scope == definition.ScopeTypeCluster {
				writer.F("func (c *%s) List%s(ctx context.Context, opts metav1.ListOptions) (*%s.%sList, error) {", clientName, m.ShortName, m.Package.Alias, m.ShortName)
				writer.F("result, err := c.backend.ListClusterScoped(ctx, %q, %q, opts, &%s.%sList{})", strings.ToLower(stringsutil.Plural(m.ShortName)), m.ShortName, m.Package.Alias, m.ShortName)
				writer.F("if err != nil {")
				writer.F("return nil, err")
				writer.F("}")
				writer.F("return result.(*%s.%sList), nil", m.Package.Alias, m.ShortName)
				writer.F("}")
				writer.F("")
			} else {
				writer.F("func (c *%s) List%s(ctx context.Context, namespace string, opts metav1.ListOptions) (*%s.%sList, error) {", clientName, m.ShortName, m.Package.Alias, m.ShortName)
				writer.F("result, err := c.backend.List(ctx, %q, %q, namespace, opts, &%s.%sList{})", strings.ToLower(stringsutil.Plural(m.ShortName)), m.ShortName, m.Package.Alias, m.ShortName)
				writer.F("if err != nil {")
				writer.F("return nil, err")
				writer.F("}")
				writer.F("return result.(*%s.%sList), nil", m.Package.Alias, m.ShortName)
				writer.F("}")
				writer.F("")
			}

			// WatchXXX
			if m.Scope == definition.ScopeTypeCluster {
				writer.F("func (c *%s) Watch%s(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {", clientName, m.ShortName)
				writer.F("return c.backend.WatchClusterScoped(ctx, schema.GroupVersionResource{Group:%q, Version:%q, Resource:%q}, opts)", m.Group, m.Version, strings.ToLower(stringsutil.Plural(m.ShortName)))
				writer.F("}")
				writer.F("")
			} else {
				writer.F("func (c *%s) Watch%s(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {", clientName, m.ShortName)
				writer.F("return c.backend.Watch(ctx, schema.GroupVersionResource{Group:%q, Version:%q, Resource:%q}, namespace, opts)", m.Group, m.Version, strings.ToLower(stringsutil.Plural(m.ShortName)))
				writer.F("}")
				writer.F("")
			}
		}
	}

	return nil
}

type informerGenerator struct {
	groupVersions map[string][]*definition.Message
}

func (g *informerGenerator) Import() map[string]string {
	importPackages := map[string]string{
		"reflect":                                "",
		"sync":                                   "",
		"context":                                "",
		"time":                                   "",
		"k8s.io/client-go/rest":                  "",
		"k8s.io/apimachinery/pkg/apis/meta/v1":   "metav1",
		"k8s.io/apimachinery/pkg/watch":          "",
		"k8s.io/apimachinery/pkg/runtime":        "",
		"k8s.io/apimachinery/pkg/runtime/schema": "",
		"k8s.io/client-go/tools/cache":           "",
	}
	for _, v := range g.groupVersions {
		for _, m := range v {
			_, p := path.Split(m.Package.Path)
			alias := m.Package.Alias
			if p == m.Package.Alias {
				alias = ""
			}
			importPackages[m.Package.Path] = alias
		}
	}

	return importPackages
}

func newInformerGenerator(groupVersions map[string][]*definition.Message) *informerGenerator {
	return &informerGenerator{groupVersions: groupVersions}
}

func (g *informerGenerator) WriteTo(writer *codegeneration.Writer) error {
	writer.F("type InformerCache struct {")
	writer.F("mu sync.Mutex")
	writer.F("informers map[reflect.Type]cache.SharedIndexInformer")
	writer.F("}")
	writer.F("func NewInformerCache() *InformerCache {")
	writer.F("return &InformerCache{informers: make(map[reflect.Type]cache.SharedIndexInformer)}")
	writer.F("}")
	writer.F("")
	writer.F("func (c *InformerCache) Write(obj runtime.Object, newFunc func() cache.SharedIndexInformer) cache.SharedIndexInformer {")
	writer.F("c.mu.Lock()")
	writer.F("defer c.mu.Unlock()")
	writer.F("")
	writer.F("typ := reflect.TypeOf(obj)")
	writer.F("if v, ok := c.informers[typ]; ok {")
	writer.F("return v")
	writer.F("}")
	writer.F("informer := newFunc()")
	writer.F("c.informers[typ] = informer")
	writer.F("")
	writer.F("return informer")
	writer.F("}") // end of Write
	writer.F("")
	writer.F("func (c *InformerCache) Informers() []cache.SharedIndexInformer {")
	writer.F("c.mu.Lock()")
	writer.F("defer c.mu.Unlock()")
	writer.F("")
	writer.F("a := make([]cache.SharedIndexInformer, 0, len(c.informers))")
	writer.F("for _, v := range c.informers {")
	writer.F("a = append(a, v)")
	writer.F("}")
	writer.F("")
	writer.F("return a")
	writer.F("}")
	writer.F("")

	writer.F("type InformerFactory struct {")
	writer.F("set *Set")
	writer.F("cache *InformerCache")
	writer.F("")
	writer.F("namespace string")
	writer.F("resyncPeriod time.Duration")
	writer.F("}")
	writer.F("")
	writer.F("func NewInformerFactory(s *Set, c *InformerCache, namespace string, resyncPeriod time.Duration) *InformerFactory {")
	writer.F("return &InformerFactory{set: s, cache: c, namespace: namespace, resyncPeriod: resyncPeriod}")
	writer.F("}") // end of NewInformerFactory
	writer.F("")
	writer.F("func (f *InformerFactory) Cache() *InformerCache {")
	writer.F("return f.cache")
	writer.F("}") // end of Cache
	writer.F("")

	writer.F("func (f *InformerFactory) InformerFor(obj runtime.Object) cache.SharedIndexInformer {")
	writer.F("switch obj.(type) {")
	for _, k := range keys(g.groupVersions) {
		v := g.groupVersions[k]
		for _, m := range v {
			clientName := m.ClientName()
			writer.F("case *%s.%s:", m.Package.Alias, m.ShortName)
			writer.F("return New%sInformer(f.cache, f.set.%s, f.namespace, f.resyncPeriod).%sInformer()", clientName, clientName, m.ShortName)
		}
	}
	writer.F("default:")
	writer.F("return nil")
	writer.F("}")
	writer.F("}") // end of InformerFor
	writer.F("")

	writer.F("func (f *InformerFactory) InformerForResource(gvr schema.GroupVersionResource) cache.SharedIndexInformer {")
	writer.F("switch gvr {")
	for _, k := range keys(g.groupVersions) {
		v := g.groupVersions[k]
		for _, m := range v {
			clientName := m.ClientName()
			writer.F("case %s.SchemaGroupVersion.WithResource(%q):", m.Package.Alias, strings.ToLower(stringsutil.Plural(m.ShortName)))
			writer.F("return New%sInformer(f.cache, f.set.%s, f.namespace, f.resyncPeriod).%sInformer()", clientName, clientName, m.ShortName)
		}
	}
	writer.F("default:")
	writer.F("return nil")
	writer.F("}")
	writer.F("}") // end of InformerForResource
	writer.F("")

	writer.F("func (f *InformerFactory) Run(ctx context.Context) {")
	writer.F("for _, v := range f.cache.Informers() {")
	writer.F("go v.Run(ctx.Done())")
	writer.F("}")
	writer.F("}") // end of Run
	writer.F("")

	for _, k := range keys(g.groupVersions) {
		v := g.groupVersions[k]
		m := v[0]
		clientName := m.ClientName()

		writer.F("type %sInformer struct {", clientName)
		writer.F("cache *InformerCache")
		writer.F("client  *%s", clientName)
		writer.F("namespace string")
		writer.F("resyncPeriod time.Duration")
		writer.F("indexers cache.Indexers")
		writer.F("}")
		writer.F("")
		writer.F("func New%sInformer(c *InformerCache, client *%s, namespace string, resyncPeriod time.Duration) *%sInformer {", clientName, clientName, clientName)
		writer.F("return &%sInformer{", clientName)
		writer.F("cache: c,")
		writer.F("client: client,")
		writer.F("namespace: namespace,")
		writer.F("resyncPeriod: resyncPeriod,")
		writer.F("indexers: cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},")
		writer.F("}")
		writer.F("}") // end of NewXXXInformer
		writer.F("")

		for _, m := range v {
			writer.F(
				"func (f *%sInformer) %sInformer() cache.SharedIndexInformer{",
				clientName,
				m.ShortName,
			)
			writer.F("return f.cache.Write(&%s.%s{}, func () cache.SharedIndexInformer{", m.Package.Alias, m.ShortName)
			if m.Scope == definition.ScopeTypeCluster {
				writer.F("return cache.NewSharedIndexInformer(")
				writer.F("&cache.ListWatch{")
				writer.F("ListFunc: func (options metav1.ListOptions) (runtime.Object, error){")
				writer.F("return f.client.List%s(context.TODO(), metav1.ListOptions{})", m.ShortName)
				writer.F("},") // end of ListFunc
				writer.F("WatchFunc: func (options metav1.ListOptions) (watch.Interface, error){")
				writer.F("return f.client.Watch%s(context.TODO(), metav1.ListOptions{})", m.ShortName)
				writer.F("},") // end of WatchFunc
				writer.F("},")
				writer.F("&%s.%s{},", m.Package.Alias, m.ShortName)
				writer.F("f.resyncPeriod,")
				writer.F("f.indexers,")
				writer.F(")")
			} else {
				writer.F("return cache.NewSharedIndexInformer(")
				writer.F("&cache.ListWatch{")
				writer.F("ListFunc: func (options metav1.ListOptions) (runtime.Object, error){")
				writer.F("return f.client.List%s(context.TODO(), f.namespace, metav1.ListOptions{})", m.ShortName)
				writer.F("},") // end of ListFunc
				writer.F("WatchFunc: func (options metav1.ListOptions) (watch.Interface, error){")
				writer.F("return f.client.Watch%s(context.TODO(), f.namespace, metav1.ListOptions{})", m.ShortName)
				writer.F("},") // end of WatchFunc
				writer.F("},")
				writer.F("&%s.%s{},", m.Package.Alias, m.ShortName)
				writer.F("f.resyncPeriod,")
				writer.F("f.indexers,")
				writer.F(")")
			}
			writer.F("})")
			writer.F("}") // end of NewXXXInformer
			writer.F("")

			writer.F("func (f *%sInformer) %sLister() *%s%sLister {", clientName, m.ShortName, clientName, m.ShortName)
			writer.F("return New%s%sLister(f.%sInformer().GetIndexer())", clientName, m.ShortName, m.ShortName)
			writer.F("}")
			writer.F("")
		}
	}

	return nil
}

type listerGenerator struct {
	groupVersions map[string][]*definition.Message
}

func newListerGenerator(groupVersions map[string][]*definition.Message) *listerGenerator {
	return &listerGenerator{groupVersions: groupVersions}
}

func (g *listerGenerator) Import() map[string]string {
	importPackages := map[string]string{
		"k8s.io/client-go/tools/cache":       "",
		"k8s.io/apimachinery/pkg/labels":     "",
		"k8s.io/apimachinery/pkg/api/errors": "k8serrors",
	}
	for _, v := range g.groupVersions {
		for _, m := range v {
			_, p := path.Split(m.Package.Path)
			alias := m.Package.Alias
			if p == m.Package.Alias {
				alias = ""
			}
			importPackages[m.Package.Path] = alias
		}
	}

	return importPackages
}

func (g *listerGenerator) WriteTo(writer *codegeneration.Writer) error {
	for _, k := range keys(g.groupVersions) {
		v := g.groupVersions[k]
		m := v[0]
		clientName := m.ClientName()

		for _, m := range v {
			writer.F("type %s%sLister struct {", clientName, m.ShortName)
			writer.F("indexer cache.Indexer")
			writer.F("}")
			writer.F("")
			writer.F("func New%s%sLister(indexer cache.Indexer) *%s%sLister {", clientName, m.ShortName, clientName, m.ShortName)
			writer.F("return &%s%sLister{indexer: indexer}", clientName, m.ShortName)
			writer.F("}")
			writer.F("")

			// ListXXX
			if m.Scope == definition.ScopeTypeCluster {
				writer.F("func (x *%s%sLister) List(selector labels.Selector) ([]*%s.%s, error) {", clientName, m.ShortName, m.Package.Alias, m.ShortName)
				writer.F("var ret []*%s.%s", m.Package.Alias, m.ShortName)
				writer.F("err := cache.ListAll(x.indexer, selector, func(m interface{}) {")
				writer.F("ret = append(ret, m.(*%s.%s).DeepCopy())", m.Package.Alias, m.ShortName)
				writer.F("})")
				writer.F("return ret, err")
				writer.F("}")
				writer.F("")
			} else {
				writer.F("func (x *%s%sLister) List(namespace string, selector labels.Selector) ([]*%s.%s, error) {", clientName, m.ShortName, m.Package.Alias, m.ShortName)
				writer.F("var ret []*%s.%s", m.Package.Alias, m.ShortName)
				writer.F("err := cache.ListAllByNamespace(x.indexer, namespace, selector, func(m interface{}) {")
				writer.F("ret = append(ret, m.(*%s.%s).DeepCopy())", m.Package.Alias, m.ShortName)
				writer.F("})")
				writer.F("return ret, err")
				writer.F("}")
				writer.F("")
			}

			// GetXXX
			if m.Scope == definition.ScopeTypeCluster {
				writer.F("func (x *%s%sLister) Get(name string) (*%s.%s, error) {", clientName, m.ShortName, m.Package.Alias, m.ShortName)
				writer.F("obj, exists, err := x.indexer.GetByKey(\"/\" + name)")
				writer.F("if err != nil {")
				writer.F("return nil, err")
				writer.F("}")
				writer.F("if !exists {")
				writer.F("return nil, k8serrors.NewNotFound(%s.SchemaGroupVersion.WithResource(%q).GroupResource(), name)", m.Package.Alias, strings.ToLower(m.ShortName))
				writer.F("}")
				writer.F("return obj.(*%s.%s).DeepCopy(), nil", m.Package.Alias, m.ShortName)
				writer.F("}")
				writer.F("")
			} else {
				writer.F("func (x *%s%sLister) Get(namespace, name string) (*%s.%s, error) {", clientName, m.ShortName, m.Package.Alias, m.ShortName)
				writer.F("obj, exists, err := x.indexer.GetByKey(namespace + \"/\" + name)")
				writer.F("if err != nil {")
				writer.F("return nil, err")
				writer.F("}")
				writer.F("if !exists {")
				writer.F("return nil, k8serrors.NewNotFound(%s.SchemaGroupVersion.WithResource(%q).GroupResource(), name)", m.Package.Alias, strings.ToLower(m.ShortName))
				writer.F("}")
				writer.F("return obj.(*%s.%s).DeepCopy(), nil", m.Package.Alias, m.ShortName)
				writer.F("}")
				writer.F("")
			}
		}
	}

	return nil
}

func keys[V any](in map[string]V) []string {
	keys := make([]string, 0, len(in))
	for k := range in {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}
