package k8s

import (
	"fmt"
	"io"
	"log"
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
	log.Print(files[0].GetOptions().GetGoPackage())
	return &ClientGenerator{files: files, lister: definition.NewLister(files, allProtos)}
}

func (g *ClientGenerator) Generate(out io.Writer, packageName string) error {
	w := codegeneration.NewWriter()
	w.F("package %s", path.Base(packageName))

	// The key is a package path. The value is an alias.
	importPackages := map[string]string{
		"k8s.io/apimachinery/pkg/runtime": "",
	}
	importProjectPackages := make(map[string]string)

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
		importProjectPackages[p] = a
	}

	w.F("import (")
	w.F("\"context\"")
	w.F("\"time\"")
	w.F("")
	for p, a := range importPackages {
		if a != "" {
			w.F("%s %q", a, p)
		} else {
			w.F("%q", p)
		}
	}
	w.F("")
	for p, a := range importProjectPackages {
		if a != "" {
			w.F("%s %q", a, p)
		} else {
			w.F("%q", p)
		}
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
			pkgName := path.Base(m.Package.Path)
			structNameWithPkg := fmt.Sprintf("%s.%s", pkgName, m.ShortName)
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
			writer.F("func (c *%s) List%s(ctx context.Context, namespace string, opts metav1.ListOptions) (*%s.%sList, error) {", clientName, m.ShortName, pkgName, m.ShortName)
			writer.F("var timeout time.Duration")
			writer.F("if opts.TimeoutSeconds != nil {")
			writer.F("timeout = time.Duration(*opts.TimeoutSeconds) * time.Second")
			writer.F("}")
			writer.F("result := &%s.%sList{}", pkgName, m.ShortName)
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
