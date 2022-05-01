package k8s

import (
	"io"
	"path"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"

	"go.f110.dev/kubeproto/internal/codegeneration"
	"go.f110.dev/kubeproto/internal/definition"
)

type DeepCopyGenerator struct {
	file   *descriptorpb.FileDescriptorProto
	lister *definition.Lister
}

func NewDeepCopyGenerator(file *descriptorpb.FileDescriptorProto, allProtos []*descriptorpb.FileDescriptorProto) *DeepCopyGenerator {
	return &DeepCopyGenerator{file: file, lister: definition.NewLister([]*descriptorpb.FileDescriptorProto{file}, allProtos)}
}

func (g *DeepCopyGenerator) Generate(out io.Writer) error {
	w := codegeneration.NewWriter()
	messages := g.lister.GetMessages()

	packageName := g.file.GetOptions().GetGoPackage()
	w.F("package %s", path.Base(packageName))

	importPackages := map[string]string{
		"k8s.io/apimachinery/pkg/runtime": "",
	}
	defW := codegeneration.NewWriter()

	enums := g.lister.GetEnums()
	for _, enum := range enums.Own() {
		// Enum definition
		defW.F("type %s string", enum.ShortName)
		defW.F("const (")
		for _, v := range enum.Values {
			defW.F("%s%s %s = %q", enum.ShortName, v, enum.ShortName, v)
		}
		defW.F(")")
	}

	mark := make(map[string]struct{})
	objs := messages.FilterKind()
	for len(objs) > 0 {
		obj := objs[0]
		mark[obj.Name] = struct{}{}

		// Struct definition
		defW.F("type %s struct {", obj.ShortName)
		for _, f := range obj.Fields {
			switch f.Type {
			case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
				m := messages.Find(f.MessageName)
				if m == nil {
					continue
				}
				if m.Package.Path != "" && m.Package.Path != packageName {
					importPackages[m.Package.Path] = m.Package.Alias
				}

				if _, ok := mark[m.Name]; !ok && !m.Dep {
					objs = append(objs, m)
				}
			}
			var name string
			if !f.Embed {
				name = f.Name.CamelCase()
			}
			typ := g.lister.ResolveGoType(packageName, f)
			tag := f.Tag()
			if f.Description != "" {
				d := strings.Replace(f.Description, string(f.Name), f.Name.CamelCase(), 1)
				defW.F("// %s", d)
			}
			defW.F("%s %s %s", name, typ, tag)
		}
		defW.F("}")
		defW.F("")

		// DeepCopy functions (DeepCopyInto / DeepCopy / DeepCopyObject)
		defW.F("func (in *%s) DeepCopyInto(out *%s) {", obj.ShortName, obj.ShortName)
		defW.F("*out = *in")
		for _, f := range obj.Fields {
			switch f.Type {
			case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
				if f.Repeated {
					defW.F("if in.%s != nil {", f.Name.CamelCase())
					defW.F("l := make(%s, len(in.%s))", g.lister.ResolveGoType(packageName, f), f.Name.CamelCase())
					defW.F("for i := range in.%s {", f.Name.CamelCase())
					defW.F("in.%s[i].DeepCopyInto(&l[i])", f.Name.CamelCase())
					defW.F("}")
					defW.F("out.%s = l", f.Name.CamelCase())
					defW.F("}")
					continue
				}
				if f.Optional {
					defW.F("if in.%s != nil {", f.Name.CamelCase())
					defW.F("in, out := &in.%s, &out.%s", f.Name.CamelCase(), f.Name.CamelCase())
					defW.F("*out = new(%s)", f.Name.CamelCase())
					defW.F("}")
				}
				if f.Inline {
					defW.F("out.%s = in.%s", f.Name.CamelCase(), f.Name.CamelCase())
				} else {
					defW.F("in.%s.DeepCopyInto(&out.%s)", f.Name.CamelCase(), f.Name.CamelCase())
				}
			default:
				if f.Repeated {
					defW.F("if in.%s != nil {", f.Name.CamelCase())
					defW.F("t := make(%s, len(in.%s))", g.lister.ResolveGoType(packageName, f), f.Name.CamelCase())
					defW.F("copy(t, in.%s)", f.Name.CamelCase())
					defW.F("out.%s = t", f.Name.CamelCase())
					defW.F("}")
				}
			}
		}
		defW.F("}")
		defW.F("")
		defW.F("func (in *%s) DeepCopy() *%s {", obj.ShortName, obj.ShortName)
		defW.F("if in == nil {\nreturn nil\n}")
		defW.F("out := new(%s)", obj.ShortName)
		defW.F("in.DeepCopyInto(out)")
		defW.F("return out")
		defW.F("}")
		defW.F("")
		// DeepCopyObject should generate only top level object.
		if obj.Kind {
			defW.F("func (in *%s) DeepCopyObject() runtime.Object {", obj.ShortName)
			defW.F("if c := in.DeepCopy(); c != nil {")
			defW.F("return c")
			defW.F("}")
			defW.F("return nil")
			defW.F("}")
		}

		objs = objs[1:]
	}

	w.F("import (")
	for p, a := range importPackages {
		if a != "" {
			w.F("%s %q", a, p)
		} else {
			w.F("%q", p)
		}
	}
	w.F(")")
	w.F("")
	defW.WriteTo(w)

	if err := w.Format(); err != nil {
		return err
	}
	if _, err := w.WriteTo(out); err != nil {
		return err
	}
	return nil
}
