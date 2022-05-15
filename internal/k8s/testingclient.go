package k8s

import (
	"fmt"
	"io"
	"path"
	"strings"

	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	"go.f110.dev/kubeproto/internal/codegeneration"
	"go.f110.dev/kubeproto/internal/definition"
	"go.f110.dev/kubeproto/internal/stringsutil"
)

type FakeClientGenerator struct {
	files  []*descriptorpb.FileDescriptorProto
	lister *definition.Lister
}

func NewFakeClientGenerator(fileToGenerate []string, files *protoregistry.Files) *FakeClientGenerator {
	return &FakeClientGenerator{
		files:  nil,
		lister: definition.NewLister(fileToGenerate, files),
	}
}

func (g *FakeClientGenerator) Generate(out io.Writer, packageName, importPath string) error {
	w := codegeneration.NewWriter()
	w.F("package %s", path.Base(packageName))

	// The key is a package path. The value is an alias.
	importPackages := map[string]string{
		"context": "",
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

	restClient := newRestFakeClientGenerator(groupVersions)
	if err := restClient.WriteTo(writer); err != nil {
		return err
	}
	for p, a := range restClient.Import() {
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

type restFakeClientGenerator struct {
	groupVersions map[string][]*definition.Message
}

func newRestFakeClientGenerator(groupVersions map[string][]*definition.Message) *restFakeClientGenerator {
	return &restFakeClientGenerator{groupVersions: groupVersions}
}

func (g *restFakeClientGenerator) Import() map[string]string {
	importPackages := map[string]string{
		"k8s.io/apimachinery/pkg/apis/meta/v1":       "metav1",
		"k8s.io/apimachinery/pkg/watch":              "",
		"k8s.io/apimachinery/pkg/labels":             "",
		"k8s.io/apimachinery/pkg/runtime":            "",
		"k8s.io/apimachinery/pkg/runtime/serializer": "",
		"k8s.io/client-go/testing":                   "k8stesting",
	}
	for _, v := range g.groupVersions {
		for _, m := range v {
			importPackages[m.Package.Path] = m.Package.Alias
		}
	}

	return importPackages
}

func (g *restFakeClientGenerator) WriteTo(writer *codegeneration.Writer) error {
	writer.F("var (")
	writer.F("scheme = runtime.NewScheme()")
	writer.F("codecs = serializer.NewCodecFactory(scheme)")
	writer.F(")")
	writer.F("")

	writer.F("func init() {")
	writer.F("for _, v := range []func(*runtime.Scheme) error{")
	for _, key := range keys(g.groupVersions) {
		v := g.groupVersions[key]
		m := v[0]
		writer.F("%s.AddToScheme,", path.Base(m.Package.Path))
	}
	writer.F("} {")
	writer.F("if err := v(scheme); err != nil {\npanic(err)\n}")
	writer.F("}") // end of for
	writer.F("}") // end of init()

	writer.F("type Set struct {")
	writer.F("k8stesting.Fake")
	writer.F("")
	for _, k := range keys(g.groupVersions) {
		m := g.groupVersions[k][0]
		clientName := fmt.Sprintf("%s%s", stringsutil.ToUpperCamelCase(m.SubGroup), stringsutil.ToUpperCamelCase(m.Version))
		writer.F("%s *Testing%s", clientName, clientName)
	}
	writer.F("")
	writer.F("tracker k8stesting.ObjectTracker")
	writer.F("}")
	writer.F("")
	writer.F("func NewSet() *Set {")
	writer.F("s := &Set{}")
	writer.F("o := k8stesting.NewObjectTracker(scheme, codecs.UniversalDecoder())")
	writer.F("s.AddReactor(\"*\", \"*\", k8stesting.ObjectReaction(o))")
	writer.F("s.AddWatchReactor(\"*\", func(action k8stesting.Action) (handled bool, ret watch.Interface, err error) {")
	writer.F("w, err := o.Watch(action.GetResource(), action.GetNamespace())")
	writer.F("if err != nil {")
	writer.F("return false, nil, err")
	writer.F("}")
	writer.F("return true, w, nil")
	writer.F("})")
	writer.F("")
	writer.F("return &Set{")
	for _, k := range keys(g.groupVersions) {
		m := g.groupVersions[k][0]
		clientName := fmt.Sprintf("%s%s", stringsutil.ToUpperCamelCase(m.SubGroup), stringsutil.ToUpperCamelCase(m.Version))
		writer.F("%s: NewTesting%sClient(&s.Fake),", clientName, clientName)
	}
	writer.F("tracker: o,")
	writer.F("}")
	writer.F("}")
	writer.F("")
	writer.F("func (s *Set) Tracker() k8stesting.ObjectTracker {")
	writer.F("return s.tracker")
	writer.F("}") // end of Tracker
	writer.F("")

	for _, k := range keys(g.groupVersions) {
		v := g.groupVersions[k]
		m := v[0]
		clientName := fmt.Sprintf("%s%s", stringsutil.ToUpperCamelCase(m.SubGroup), stringsutil.ToUpperCamelCase(m.Version))
		writer.F("type Testing%s struct {", clientName)
		writer.F("*k8stesting.Fake")
		writer.F("}")
		writer.F("")

		writer.F("func NewTesting%sClient(fake *k8stesting.Fake) *Testing%s {", clientName, clientName)
		writer.F("return &Testing%s{Fake: fake}", clientName)
		writer.F("}") // end of NewFakeXXXClient
		writer.F("")

		for _, m := range v {
			structNameWithPkg := fmt.Sprintf("%s.%s", m.Package.Name, m.ShortName)
			resourceName := strings.ToLower(stringsutil.Plural(m.ShortName))
			// GetXXX
			writer.F("func(c *Testing%s) Get%s(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*%s, error) {", clientName, m.ShortName, structNameWithPkg)
			writer.F(
				"obj, err := c.Fake.Invokes(k8stesting.NewGetAction(%s.SchemaGroupVersion.WithResource(%q), namespace, name), &%s{})",
				m.Package.Name,
				resourceName,
				structNameWithPkg,
			)
			writer.F("if obj == nil {")
			writer.F("return nil, err")
			writer.F("}")
			writer.F("return obj.(*%s), err", structNameWithPkg)
			writer.F("}") // end of GetXXX
			writer.F("")

			// CreateXXX
			writer.F("func (c *Testing%s) Create%s(ctx context.Context, v *%s, opts metav1.CreateOptions) (*%s, error) {", clientName, m.ShortName, structNameWithPkg, structNameWithPkg)
			writer.F("obj, err := c.Fake.")
			writer.F("Invokes(k8stesting.NewCreateAction(%s.SchemaGroupVersion.WithResource(%q), v.Namespace, v), &%s{})", m.Package.Name, resourceName, structNameWithPkg)
			writer.F("")
			writer.F("if obj == nil {")
			writer.F("return nil, err")
			writer.F("}")
			writer.F("return obj.(*%s), err", structNameWithPkg)
			writer.F("}") // end of CreateXXX
			writer.F("")

			// UpdateXXX
			writer.F("func (c *Testing%s) Update%s(ctx context.Context, v *%s, opts metav1.UpdateOptions) (*%s, error) {", clientName, m.ShortName, structNameWithPkg, structNameWithPkg)
			writer.F("obj, err := c.Fake.")
			writer.F("Invokes(k8stesting.NewUpdateAction(%s.SchemaGroupVersion.WithResource(%q), v.Namespace, v), &%s{})", m.Package.Name, resourceName, structNameWithPkg)
			writer.F("")
			writer.F("if obj == nil {")
			writer.F("return nil, err")
			writer.F("}")
			writer.F("return obj.(*%s), err", structNameWithPkg)
			writer.F("}") // end of UpdateXXX
			writer.F("")

			// UpdateStatusXXX
			if m.IsDefinedSubResource() {
				writer.F("func (c *Testing%s) UpdateStatus%s(ctx context.Context, v *%s, opts metav1.UpdateOptions) (*%s, error) {", clientName, m.ShortName, structNameWithPkg, structNameWithPkg)
				writer.F("obj, err := c.Fake.")
				writer.F("Invokes(k8stesting.NewUpdateSubresourceAction(%s.SchemaGroupVersion.WithResource(%q), \"status\", v.Namespace, v), &%s{})", m.Package.Name, resourceName, structNameWithPkg)
				writer.F("")
				writer.F("if obj == nil {")
				writer.F("return nil, err")
				writer.F("}")
				writer.F("return obj.(*%s), err", structNameWithPkg)
				writer.F("}") // end of UpdateStatusXXX
				writer.F("")
			}

			// DeleteXXX
			writer.F("func (c *Testing%s) Delete%s(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {", clientName, m.ShortName)
			writer.F("_, err := c.Fake.")
			writer.F("Invokes(k8stesting.NewDeleteAction(%s.SchemaGroupVersion.WithResource(%q), namespace, name), &%s{})", m.Package.Name, resourceName, structNameWithPkg)
			writer.F("")
			writer.F("return err")
			writer.F("}") // end of DeleteXXX
			writer.F("")

			// ListXXX
			writer.F("func (c *Testing%s) List%s(ctx context.Context, namespace string, opts metav1.ListOptions) (*%s.%sList, error) {", clientName, m.ShortName, m.Package.Name, m.ShortName)
			writer.F("obj, err := c.Fake.")
			writer.F("Invokes(k8stesting.NewListAction(%s.SchemaGroupVersion.WithResource(%q), %s.SchemaGroupVersion.WithKind(%q), namespace, opts), &%sList{})", m.Package.Name, resourceName, m.Package.Name, m.ShortName, structNameWithPkg)
			writer.F("")
			writer.F("if obj == nil {")
			writer.F("return nil, err")
			writer.F("}")
			writer.F("")
			writer.F("label, _, _ := k8stesting.ExtractFromListOptions(opts)")
			writer.F("if label == nil {")
			writer.F("label = labels.Everything()")
			writer.F("}")
			writer.F("list := &%sList{ListMeta: obj.(*%sList).ListMeta}", structNameWithPkg, structNameWithPkg)
			writer.F("for _, item := range obj.(*%sList).Items {", structNameWithPkg)
			writer.F("if label.Matches(labels.Set(item.Labels)) {")
			writer.F("list.Items = append(list.Items, item)")
			writer.F("}")
			writer.F("}")
			writer.F("return list, err")
			writer.F("}") // end of ListXXX
			writer.F("")

			// WatchXXX
			writer.F("func (c *Testing%s) Watch%s(ctx context.Context, namespace string, opts metav1.ListOptions) (watch.Interface, error) {", clientName, m.ShortName)
			writer.F("return c.Fake.InvokesWatch(k8stesting.NewWatchAction(%s.SchemaGroupVersion.WithResource(%q), namespace, opts))", m.Package.Name, resourceName)
			writer.F("}") // end of WatchXXX
			writer.F("")
		}
	}

	return nil
}
