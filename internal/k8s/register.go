package k8s

import (
	"fmt"
	"io"
	"path"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	"go.f110.dev/kubeproto"
	"go.f110.dev/kubeproto/internal/codegeneration"
	"go.f110.dev/kubeproto/internal/definition"
)

type RegisterGenerator struct {
	file   *descriptorpb.FileDescriptorProto
	lister *definition.Lister
}

func NewRegisterGenerator(file *descriptorpb.FileDescriptorProto, allProtos []*descriptorpb.FileDescriptorProto) *RegisterGenerator {
	return &RegisterGenerator{file: file, lister: definition.NewLister([]*descriptorpb.FileDescriptorProto{file}, allProtos)}
}

func (g *RegisterGenerator) Generate(out io.Writer) error {
	w := codegeneration.NewWriter()

	e := proto.GetExtension(g.file.GetOptions(), kubeproto.E_K8S)
	ext := e.(*kubeproto.Kubernetes)
	if ext == nil {
		return fmt.Errorf("%s is not extended by kubeproto.Kubernetes", g.file.GetName())
	}

	packageName := g.file.GetOptions().GetGoPackage()
	w.F("package %s", path.Base(packageName))
	w.F("import (")
	w.F("metav1 \"k8s.io/apimachinery/pkg/apis/meta/v1\"")
	w.F("\"k8s.io/apimachinery/pkg/runtime\"")
	w.F("\"k8s.io/apimachinery/pkg/runtime/schema\"")
	w.F(")")
	w.F("const GroupName = \"%s.%s\"", ext.SubGroup, ext.Domain)
	w.F("")
	w.F("var (")
	w.F("GroupVersion = metav1.GroupVersion{Group: GroupName, Version: %q}", ext.Version)
	w.F("SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)")
	w.F("AddToScheme = SchemeBuilder.AddToScheme")
	w.F("SchemaGroupVersion = schema.GroupVersion{Group: \"%s.%s\", Version: %q}", ext.SubGroup, ext.Domain, ext.Version)
	w.F(")")

	w.F("func addKnownTypes(scheme *runtime.Scheme) error {")
	w.F("scheme.AddKnownTypes(SchemaGroupVersion,")
	messages := g.lister.GetMessages()
	for _, m := range messages.FilterKind() {
		w.F("&%s{},", m.ShortName)
	}
	w.F(")")
	w.F("metav1.AddToGroupVersion(scheme, SchemaGroupVersion)")
	w.F("return nil")
	w.F("}")

	if err := w.Format(); err != nil {
		return err
	}
	if _, err := w.WriteTo(out); err != nil {
		return err
	}
	return nil
}
