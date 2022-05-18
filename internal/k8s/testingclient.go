package k8s

import (
	"fmt"
	"io"
	"path"

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

func (g *FakeClientGenerator) Generate(out io.Writer, packageName, importPath, clientPath string) error {
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

	restClient := newRestFakeClientGenerator(groupVersions, clientPath)
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
	clientPath string

	groupVersions map[string][]*definition.Message
}

func newRestFakeClientGenerator(groupVersions map[string][]*definition.Message, clientPath string) *restFakeClientGenerator {
	return &restFakeClientGenerator{clientPath: clientPath, groupVersions: groupVersions}
}

func (g *restFakeClientGenerator) Import() map[string]string {
	importPackages := map[string]string{
		"k8s.io/apimachinery/pkg/api/meta":           "",
		"k8s.io/apimachinery/pkg/apis/meta/v1":       "metav1",
		"k8s.io/apimachinery/pkg/watch":              "",
		"k8s.io/apimachinery/pkg/labels":             "",
		"k8s.io/apimachinery/pkg/runtime":            "",
		"k8s.io/apimachinery/pkg/runtime/serializer": "",
		"k8s.io/client-go/testing":                   "k8stesting",
		g.clientPath:                                 "",
	}
	for _, v := range g.groupVersions {
		for _, m := range v {
			importPackages[m.Package.Path] = m.Package.Alias
		}
	}

	return importPackages
}

func (g *restFakeClientGenerator) WriteTo(writer *codegeneration.Writer) error {
	clientPackageName := path.Base(g.clientPath)

	writer.F("var (")
	writer.F("codecs = serializer.NewCodecFactory(%s.Scheme)", clientPackageName)
	writer.F(")")
	writer.F("")

	writer.F("type Set struct {")
	writer.F("%s.Set", clientPackageName)
	writer.F("")
	writer.F("fake k8stesting.Fake")
	writer.F("tracker k8stesting.ObjectTracker")
	writer.F("}")
	writer.F("")
	writer.F("func NewSet() *Set {")
	writer.F("s := &Set{}")
	writer.F("o := k8stesting.NewObjectTracker(%s.Scheme, codecs.UniversalDecoder())", clientPackageName)
	writer.F("s.fake.AddReactor(\"*\", \"*\", k8stesting.ObjectReaction(o))")
	writer.F("s.fake.AddWatchReactor(\"*\", func(action k8stesting.Action) (handled bool, ret watch.Interface, err error) {")
	writer.F("w, err := o.Watch(action.GetResource(), action.GetNamespace())")
	writer.F("if err != nil {")
	writer.F("return false, nil, err")
	writer.F("}")
	writer.F("return true, w, nil")
	writer.F("})")
	writer.F("")
	for _, k := range keys(g.groupVersions) {
		m := g.groupVersions[k][0]
		clientName := fmt.Sprintf("%s%s", stringsutil.ToUpperCamelCase(m.SubGroup), stringsutil.ToUpperCamelCase(m.Version))
		writer.F("s.%s = %s.New%sClient(&fakerBackend{fake: &s.fake})", clientName, clientPackageName, clientName)
	}
	writer.F("return s")
	writer.F("}") // end of NewSet
	writer.F("")
	writer.F("func (s *Set) Tracker() k8stesting.ObjectTracker {")
	writer.F("return s.tracker")
	writer.F("}") // end of Tracker
	writer.F("")

	writer.F(`
type fakerBackend struct {
	fake *k8stesting.Fake
}

func (f *fakerBackend) Get(ctx context.Context, resourceName, kindName, namespace, name string, opts metav1.GetOptions, result runtime.Object) (runtime.Object, error) {
	obj, err := f.fake.Invokes(k8stesting.NewGetAction(githubv1alpha1.SchemaGroupVersion.WithResource(resourceName), namespace, name), result)
	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), nil
}

func (f *fakerBackend) List(ctx context.Context, resourceName, kindName, namespace string, opts metav1.ListOptions, result runtime.Object) (runtime.Object, error) {
	obj, err := f.fake.
		Invokes(k8stesting.NewListAction(githubv1alpha1.SchemaGroupVersion.WithResource(resourceName), githubv1alpha1.SchemaGroupVersion.WithKind(kindName), namespace, opts), result)

	if obj == nil {
		return nil, err
	}

	label, _, _ := k8stesting.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	objs, err := meta.ExtractList(obj)
	if err != nil {
		return nil, err
	}
	filtered := make([]runtime.Object, 0)
	for _, item := range objs {
		m := item.(metav1.Object)
		if label.Matches(labels.Set(m.GetLabels())) {
			filtered = append(filtered, item)
		}
	}
	if err := meta.SetList(obj, filtered); err != nil {
		return nil, err
	}
	return obj.DeepCopyObject(), err
}

func (f *fakerBackend) Create(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.CreateOptions, result runtime.Object) (runtime.Object, error) {
	m := obj.(metav1.Object)
	obj, err := f.fake.
		Invokes(k8stesting.NewCreateAction(githubv1alpha1.SchemaGroupVersion.WithResource(resourceName), m.GetNamespace(), obj), result)

	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), err
}

func (f *fakerBackend) Update(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	m := obj.(metav1.Object)
	obj, err := f.fake.
		Invokes(k8stesting.NewUpdateAction(githubv1alpha1.SchemaGroupVersion.WithResource(resourceName), m.GetNamespace(), obj), result)

	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), err

}

func (f *fakerBackend) UpdateStatus(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	m := obj.(metav1.Object)
	obj, err := f.fake.
		Invokes(k8stesting.NewUpdateSubresourceAction(miniov1alpha1.SchemaGroupVersion.WithResource(resourceName), "status", m.GetNamespace(), obj), result)

	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), err
}

func (f *fakerBackend) Delete(ctx context.Context, resourceName, kindName, namespace, name string, opts metav1.DeleteOptions) error {
	_, err := f.fake.
		Invokes(k8stesting.NewDeleteAction(miniov1alpha1.SchemaGroupVersion.WithResource(resourceName), namespace, name), nil)

	return err
}

func (f *fakerBackend) Watch(ctx context.Context, resourceName, kindName, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return f.fake.InvokesWatch(k8stesting.NewWatchAction(miniov1alpha1.SchemaGroupVersion.WithResource(resourceName), namespace, opts))
}
`)

	return nil
}
