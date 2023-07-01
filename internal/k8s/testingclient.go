package k8s

import (
	"fmt"
	"io"
	"path"

	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	"go.f110.dev/kubeproto/internal/codegeneration"
	"go.f110.dev/kubeproto/internal/definition"
)

type FakeClientGenerator struct {
	files                   []*descriptorpb.FileDescriptorProto
	lister                  *definition.Lister
	packageNamespaceManager *definition.PackageNamespaceManager
}

func NewFakeClientGenerator(fileToGenerate []string, files *protoregistry.Files) *FakeClientGenerator {
	nsm := definition.NewPackageNamespaceManager()
	return &FakeClientGenerator{
		files:                   nil,
		lister:                  definition.NewLister(fileToGenerate, files, nsm),
		packageNamespaceManager: nsm,
	}
}

func (g *FakeClientGenerator) Generate(out io.Writer, packageName, importPath, clientPath string, fqdnSetName bool) error {
	w := codegeneration.NewWriter()
	w.F("package %s", path.Base(packageName))

	// The key is a package path. The value is an alias.
	importPackages := map[string]string{
		"context": "",
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

	restClient := newRestFakeClientGenerator(groupVersions, clientPath)
	if err := restClient.WriteTo(writer, fqdnSetName); err != nil {
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
		"k8s.io/apimachinery/pkg/runtime/schema":     "",
		"k8s.io/apimachinery/pkg/runtime/serializer": "",
		"k8s.io/client-go/testing":                   "k8stesting",
		g.clientPath:                                 "",
	}

	return importPackages
}

func (g *restFakeClientGenerator) WriteTo(writer *codegeneration.Writer, fqdn bool) error {
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
	writer.F("s.tracker = k8stesting.NewObjectTracker(%s.Scheme, codecs.UniversalDecoder())", clientPackageName)
	writer.F("s.fake.AddReactor(\"*\", \"*\", k8stesting.ObjectReaction(s.tracker))")
	writer.F("s.fake.AddWatchReactor(\"*\", func(action k8stesting.Action) (handled bool, ret watch.Interface, err error) {")
	writer.F("w, err := s.tracker.Watch(action.GetResource(), action.GetNamespace())")
	writer.F("if err != nil {")
	writer.F("return false, nil, err")
	writer.F("}")
	writer.F("return true, w, nil")
	writer.F("})")
	writer.F("")
	for _, k := range keys(g.groupVersions) {
		m := g.groupVersions[k][0]
		clientName := m.ClientName(fqdn)
		writer.F("s.%s = %s.New%sClient(&fakerBackend{fake: &s.fake})", clientName, clientPackageName, clientName)
	}
	writer.F("return s")
	writer.F("}") // end of NewSet
	writer.F("")
	writer.F("func (s *Set) Tracker() k8stesting.ObjectTracker {")
	writer.F("return s.tracker")
	writer.F("}") // end of Tracker
	writer.F("")
	writer.F("func (s *Set) Actions() []k8stesting.Action {")
	writer.F("return s.fake.Actions()")
	writer.F("}")
	writer.F("")

	writer.F(`
type fakerBackend struct {
	fake *k8stesting.Fake
}
`)

	writer.F(`func (f *fakerBackend) Get(ctx context.Context, resourceName, kindName, namespace, name string, opts metav1.GetOptions, result runtime.Object) (runtime.Object, error) {
	gvks, _, err := %s.Scheme.ObjectKinds(result)
	if err != nil {
		return nil, err
	}
	gvk := gvks[0]
	obj, err := f.fake.Invokes(k8stesting.NewGetAction(gvk.GroupVersion().WithResource(resourceName), namespace, name), result)
	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), nil
}`, clientPackageName)

	writer.F(`func (f *fakerBackend) List(ctx context.Context, resourceName, kindName, namespace string, opts metav1.ListOptions, result runtime.Object) (runtime.Object, error) {
	gvks, _, err := %s.Scheme.ObjectKinds(result)
	if err != nil {
		return nil, err
	}
	gvk := gvks[0]
	obj, err := f.fake.Invokes(k8stesting.NewListAction(gvk.GroupVersion().WithResource(resourceName), gvk, namespace, opts), result)

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
}`, clientPackageName)

	writer.F(`func (f *fakerBackend) Create(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.CreateOptions, result runtime.Object) (runtime.Object, error) {
	gvks, _, err := %s.Scheme.ObjectKinds(result)
	if err != nil {
		return nil, err
	}
	gvk := gvks[0]
	m := obj.(metav1.Object)
	obj, err = f.fake.Invokes(k8stesting.NewCreateAction(gvk.GroupVersion().WithResource(resourceName), m.GetNamespace(), obj), result)

	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), err
}`, clientPackageName)

	writer.F(`func (f *fakerBackend) Update(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	gvks, _, err := %s.Scheme.ObjectKinds(result)
	if err != nil {
		return nil, err
	}
	gvk := gvks[0]
	m := obj.(metav1.Object)
	obj, err = f.fake.Invokes(k8stesting.NewUpdateAction(gvk.GroupVersion().WithResource(resourceName), m.GetNamespace(), obj), result)

	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), err
}`, clientPackageName)

	writer.F(`func (f *fakerBackend) UpdateStatus(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	gvks, _, err := %s.Scheme.ObjectKinds(result)
	if err != nil {
		return nil, err
	}
	gvk := gvks[0]
	m := obj.(metav1.Object)
	obj, err = f.fake.Invokes(k8stesting.NewUpdateSubresourceAction(gvk.GroupVersion().WithResource(resourceName), "status", m.GetNamespace(), obj), result)

	if obj == nil {
		return nil, err
	}
	return obj.DeepCopyObject(), err
}`, clientPackageName)

	writer.F(`func (f *fakerBackend) Delete(ctx context.Context, gvr schema.GroupVersionResource, namespace, name string, opts metav1.DeleteOptions) error {
	_, err := f.fake.Invokes(k8stesting.NewDeleteAction(gvr, namespace, name), nil)

	return err
}`)

	writer.F(`func (f *fakerBackend) Watch(ctx context.Context, gvr schema.GroupVersionResource, namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return f.fake.InvokesWatch(k8stesting.NewWatchAction(gvr, namespace, opts))
}
`)

	// For non-namespaced resource
	writer.F(`func (f *fakerBackend) GetClusterScoped(ctx context.Context, resourceName, kindName, name string, opts metav1.GetOptions, result runtime.Object) (runtime.Object, error) {
	return f.Get(ctx, resourceName, kindName, "", name, opts, result)
}

func (f *fakerBackend) ListClusterScoped(ctx context.Context, resourceName, kindName string, opts metav1.ListOptions, result runtime.Object) (runtime.Object, error) {
	return f.List(ctx, resourceName, kindName, "", opts, result)
}

func (f *fakerBackend) CreateClusterScoped(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.CreateOptions, result runtime.Object) (runtime.Object, error) {
	return f.Create(ctx, resourceName, kindName, obj, opts, result)
}

func (f *fakerBackend) UpdateClusterScoped(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	return f.Update(ctx, resourceName, kindName, obj, opts, result)
}

func (f *fakerBackend) UpdateStatusClusterScoped(ctx context.Context, resourceName, kindName string, obj runtime.Object, opts metav1.UpdateOptions, result runtime.Object) (runtime.Object, error) {
	return f.UpdateStatus(ctx, resourceName, kindName, obj,  opts, result)
}

func (f *fakerBackend) DeleteClusterScoped(ctx context.Context, gvr schema.GroupVersionResource, name string, opts metav1.DeleteOptions) error {
	return f.Delete(ctx, gvr, "", name, opts)
}

func (f *fakerBackend) WatchClusterScoped(ctx context.Context, gvr schema.GroupVersionResource, opts metav1.ListOptions) (watch.Interface, error) {
	return f.Watch(ctx, gvr, "", opts)
}
`)

	return nil
}
