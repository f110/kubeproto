package k8s

import (
	"fmt"
	"io"
	"path"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"

	"go.f110.dev/kubeproto/internal/codegeneration"
	"go.f110.dev/kubeproto/internal/definition"
	"go.f110.dev/kubeproto/internal/stringsutil"
)

type ClientGenerator struct {
	files  []*descriptorpb.FileDescriptorProto
	lister *definition.Lister
}

func NewClientGenerator(files []*descriptorpb.FileDescriptorProto, allProtos []*descriptorpb.FileDescriptorProto) *ClientGenerator {
	return &ClientGenerator{files: files, lister: definition.NewLister(files, allProtos)}
}

func (g *ClientGenerator) Generate(out io.Writer, packageName, importPath string) error {
	w := codegeneration.NewWriter()
	w.F("package %s", path.Base(packageName))

	// The key is a package path. The value is an alias.
	importPackages := map[string]string{
		"context":                         "",
		"time":                            "",
		"k8s.io/apimachinery/pkg/runtime": "",
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
	writer.F(")")
	writer.F("")

	writer.F("func init() {")
	writer.F("for _, v := range []func(*runtime.Scheme) error{")
	for _, v := range groupVersions {
		m := v[0]
		writer.F("%s.AddToScheme,", path.Base(m.Package.Path))
	}
	writer.F("} {")
	writer.F("if err := v(Scheme); err != nil {\npanic(err)\n}")
	writer.F("}") // end of for
	writer.F("}") // end of init()
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
			importPackages[m.Package.Path] = m.Package.Alias
		}
	}

	return importPackages
}

func (g *restClientGenerator) WriteTo(writer *codegeneration.Writer) error {
	for _, v := range g.groupVersions {
		m := v[0]
		clientName := fmt.Sprintf("%s%s", stringsutil.ToUpperCamelCase(m.SubGroup), stringsutil.ToUpperCamelCase(m.Version))
		writer.F("type %s struct {", clientName)
		writer.F("client *rest.RESTClient")
		writer.F("}")
		writer.F("")

		writer.F("func New%sClient(c *rest.Config) (*%s, error) {", clientName, clientName)
		writer.F("client, err := rest.RESTClientFor(c)")
		writer.F("if err != nil {")
		writer.F("return nil, err")
		writer.F("}")
		writer.F("return &%s{", clientName)
		writer.F("client: client,")
		writer.F("}, nil")
		writer.F("}")
		writer.F("")

		for _, m := range v {
			structNameWithPkg := fmt.Sprintf("%s.%s", m.Package.Name, m.ShortName)
			// GetXXX
			writer.F("func(c *%s) Get%s(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*%s, error) {", clientName, m.ShortName, structNameWithPkg)
			writer.F("result := &%s{}", structNameWithPkg)
			writer.F("err := c.client.Get().")
			writer.F("Namespace(namespace).")
			writer.F("Resource(%q).", strings.ToLower(stringsutil.Plural(m.ShortName)))
			writer.F("Name(name).")
			writer.F("VersionedParams(&opts, ParameterCodec).")
			writer.F("Do(ctx).")
			writer.F("Into(result)")
			writer.F("if err != nil {\nreturn nil, err\n}")
			writer.F("return result, nil")
			writer.F("}") // end of GetXXX
			writer.F("")

			// CreateXXX
			writer.F("func (c *%s) Create%s(ctx context.Context, v *%s, opts metav1.CreateOptions) (*%s, error) {", clientName, m.ShortName, structNameWithPkg, structNameWithPkg)
			writer.F("result := &%s{}", structNameWithPkg)
			writer.F("err := c.client.Post().")
			writer.F("Namespace(v.Namespace).")
			writer.F("Resource(%q).", strings.ToLower(stringsutil.Plural(m.ShortName)))
			writer.F("VersionedParams(&opts, ParameterCodec).")
			writer.F("Body(v).")
			writer.F("Do(ctx).")
			writer.F("Into(result)")
			writer.F("if err != nil {\nreturn nil, err\n}")
			writer.F("return result, nil")
			writer.F("}") // end of CreateXXX
			writer.F("")

			// UpdateXXX
			writer.F("func (c *%s) Update%s(ctx context.Context, v *%s, opts metav1.UpdateOptions) (*%s, error) {", clientName, m.ShortName, structNameWithPkg, structNameWithPkg)
			writer.F("result := &%s{}", structNameWithPkg)
			writer.F("err := c.client.Put().")
			writer.F("Namespace(v.Namespace).")
			writer.F("Resource(%q).", strings.ToLower(stringsutil.Plural(m.ShortName)))
			writer.F("Name(v.Name).")
			writer.F("VersionedParams(&opts, ParameterCodec).")
			writer.F("Body(v).")
			writer.F("Do(ctx).")
			writer.F("Into(result)")
			writer.F("if err != nil {\nreturn nil, err\n}")
			writer.F("return result, nil")
			writer.F("}") // end of UpdateXXX
			writer.F("")

			// DeleteXXX
			writer.F("func (c *%s) Delete%s(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {", clientName, m.ShortName)
			writer.F("return c.client.Delete().")
			writer.F("Namespace(namespace).")
			writer.F("Resource(%q).", strings.ToLower(stringsutil.Plural(m.ShortName)))
			writer.F("Name(name).")
			writer.F("Body(&opts).")
			writer.F("Do(ctx).")
			writer.F("Error()")
			writer.F("}") // end of DeleteXXX
			writer.F("")

			// ListXXX
			writer.F("func (c *%s) List%s(ctx context.Context, namespace string, opts metav1.ListOptions) (*%s.%sList, error) {", clientName, m.ShortName, m.Package.Name, m.ShortName)
			writer.F("var timeout time.Duration")
			writer.F("if opts.TimeoutSeconds != nil {")
			writer.F("timeout = time.Duration(*opts.TimeoutSeconds) * time.Second")
			writer.F("}")
			writer.F("result := &%s.%sList{}", m.Package.Name, m.ShortName)
			writer.F("err := c.client.Get().")
			writer.F("Namespace(namespace).")
			writer.F("Resource(%q).", strings.ToLower(stringsutil.Plural(m.ShortName)))
			writer.F("VersionedParams(&opts, ParameterCodec).")
			writer.F("Timeout(timeout).")
			writer.F("Do(ctx).")
			writer.F("Into(result)")
			writer.F("if err != nil {\nreturn nil, err\n}")
			writer.F("return result, nil")
			writer.F("}") // end of ListXXX
			writer.F("")

			// WatchXXX
			writer.F("func (c *%s) Watch%s(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {", clientName, m.ShortName)
			writer.F("var timeout time.Duration")
			writer.F("if opts.TimeoutSeconds != nil {")
			writer.F("timeout = time.Duration(*opts.TimeoutSeconds) * time.Second")
			writer.F("}")
			writer.F("opts.Watch = true")
			writer.F("return c.client.Get().")
			writer.F("Namespace(namespace).")
			writer.F("Resource(%q).", strings.ToLower(stringsutil.Plural(m.ShortName)))
			writer.F("VersionedParams(&opts, ParameterCodec).")
			writer.F("Timeout(timeout).")
			writer.F("Watch(ctx)")
			writer.F("}") // end of WatchXXX
			writer.F("")
		}
	}

	return nil
}

type informerGenerator struct {
	groupVersions map[string][]*definition.Message
}

func (g *informerGenerator) Import() map[string]string {
	importPackages := map[string]string{
		"reflect":                              "",
		"sync":                                 "",
		"context":                              "",
		"time":                                 "",
		"k8s.io/client-go/rest":                "",
		"k8s.io/apimachinery/pkg/apis/meta/v1": "metav1",
		"k8s.io/apimachinery/pkg/watch":        "",
		"k8s.io/apimachinery/pkg/runtime":      "",
		"k8s.io/client-go/tools/cache":         "",
	}
	for _, v := range g.groupVersions {
		for _, m := range v {
			importPackages[m.Package.Path] = m.Package.Alias
		}
	}

	return importPackages
}

func newInformerGenerator(groupVersions map[string][]*definition.Message) *informerGenerator {
	return &informerGenerator{groupVersions: groupVersions}
}

func (g *informerGenerator) WriteTo(writer *codegeneration.Writer) error {
	writer.F("var Factory = NewInformerFactory()")
	writer.F("")
	writer.F("type InformerFactory struct {")
	writer.F("mu        sync.Mutex")
	writer.F("informers map[reflect.Type]cache.SharedIndexInformer")
	writer.F("once sync.Once")
	writer.F("ctx  context.Context")
	writer.F("}")
	writer.F("")
	writer.F("func NewInformerFactory() *InformerFactory {")
	writer.F("return &InformerFactory{informers: make(map[reflect.Type]cache.SharedIndexInformer)}")
	writer.F("}")
	writer.F("")
	writer.F("func (f *InformerFactory) InformerFor(obj runtime.Object, newFunc func() cache.SharedIndexInformer) cache.SharedIndexInformer {")
	writer.F("f.mu.Lock()")
	writer.F("defer f.mu.Unlock()")
	writer.F("")
	writer.F("typ := reflect.TypeOf(obj)")
	writer.F("if v, ok := f.informers[typ]; ok {")
	writer.F("return v")
	writer.F("}")
	writer.F("informer := newFunc()")
	writer.F("f.informers[typ] = informer")
	writer.F("if f.ctx != nil {")
	writer.F("go informer.Run(f.ctx.Done())")
	writer.F("}")
	writer.F("return informer")
	writer.F("}")
	writer.F("func (f *InformerFactory) Run(ctx context.Context) {")
	writer.F("f.mu.Lock()")
	writer.F("f.once.Do(func() {")
	writer.F("for _, v := range f.informers {")
	writer.F("go v.Run(ctx.Done())")
	writer.F("}")
	writer.F("f.ctx = ctx")
	writer.F("})")
	writer.F("f.mu.Unlock()")
	writer.F("}")

	for _, v := range g.groupVersions {
		m := v[0]
		clientName := fmt.Sprintf("%s%s", stringsutil.ToUpperCamelCase(m.SubGroup), stringsutil.ToUpperCamelCase(m.Version))

		writer.F("type %sInformer struct {", clientName)
		writer.F("factory *InformerFactory")
		writer.F("client  *%s", clientName)
		writer.F("}")
		writer.F("")
		writer.F("func New%sInformer(f *InformerFactory, client *%s) *%sInformer {", clientName, clientName, clientName)
		writer.F("return &%sInformer{factory: f, client: client}", clientName)
		writer.F("}")
		writer.F("")

		for _, m := range v {
			writer.F(
				"func(f *%sInformer) %sInformer(namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer{",
				clientName,
				m.ShortName,
			)
			writer.F("return f.factory.InformerFor(&%s.%s{}, func () cache.SharedIndexInformer{", m.Package.Name, m.ShortName)
			writer.F("return cache.NewSharedIndexInformer(")
			writer.F("&cache.ListWatch{")
			writer.F("ListFunc: func (options metav1.ListOptions) (runtime.Object, error){")
			writer.F("return f.client.List%s(context.TODO(), namespace, metav1.ListOptions{})", m.ShortName)
			writer.F("},") // end of ListFunc
			writer.F("WatchFunc: func (options metav1.ListOptions) (watch.Interface, error){")
			writer.F("return f.client.Watch%s(context.TODO(), namespace, metav1.ListOptions{})", m.ShortName)
			writer.F("},") // end of WatchFunc
			writer.F("},")
			writer.F("&%s.%s{},", m.Package.Name, m.ShortName)
			writer.F("resyncPeriod,")
			writer.F("indexers,")
			writer.F(")")
			writer.F("},")
			writer.F(")")
			writer.F("}") // end of XXXInformer
		}
	}

	return nil
}
