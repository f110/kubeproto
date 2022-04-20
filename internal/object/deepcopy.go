package object

import (
	"bufio"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	"go.f110.dev/kubeproto"
)

type DeepCopyGenerator struct {
	file     *descriptorpb.FileDescriptorProto
	allFiles []*descriptorpb.FileDescriptorProto
}

func NewDeepCopyGenerator(file *descriptorpb.FileDescriptorProto, allProtos []*descriptorpb.FileDescriptorProto) *DeepCopyGenerator {
	return &DeepCopyGenerator{file: file, allFiles: allProtos}
}

func (g *DeepCopyGenerator) Generate(out io.Writer) error {
	w := newWriter()
	messages := g.getMessages()

	packageName := g.file.GetOptions().GetGoPackage()
	w.F("package %s", path.Base(packageName))

	importPackages := make(map[string]string)
	mark := make(map[string]struct{})
	defW := newWriter()
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
				if m.Package.Path != "" {
					importPackages[m.Package.Path] = m.Package.Alias
				}

				if _, ok := mark[m.Name]; !ok && !m.Dep {
					objs = append(objs, m)
				}
			}
			typ := f.TypeName(messages)
			tag := f.Tag()
			defW.F("%s %s %s", f.Name.CamelCase(), typ, tag)
		}
		defW.F("}")
		defW.F("")

		// DeepCopy functions (DeepCopyInto / DeepCopy)
		defW.F("func (in *%s) DeepCopyInto(out *%s) {", obj.ShortName, obj.ShortName)
		defW.F("*out = *in")
		for _, f := range obj.Fields {
			switch f.Type {
			case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
				if f.Repeated {
					defW.F("if in.%s != nil {", f.Name.CamelCase())
					defW.F("l := make(%s, len(in.%s))", f.TypeName(messages), f.Name.CamelCase())
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
					defW.F("t := make(%s, len(in.%s))", f.TypeName(messages), f.Name.CamelCase())
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

	formatted, err := format.Source(w.Bytes())
	if err != nil {
		scanner := bufio.NewScanner(strings.NewReader(w.String()))
		i := 1
		for scanner.Scan() {
			fmt.Fprintf(os.Stderr, "%d: %s\n", i, scanner.Text())
			i++
		}
		return err
	}
	if _, err := out.Write(formatted); err != nil {
		return err
	}
	log.Print(string(formatted))
	return nil
}

func (g *DeepCopyGenerator) getMessages() Messages {
	var msgs Messages
	for _, v := range g.file.GetMessageType() {
		msgs = append(msgs, NewMessage(g.file, v))
	}
	for _, v := range g.allFiles {
		for _, mt := range v.GetMessageType() {
			m := NewMessage(v, mt)
			m.Dep = true
			if exists := msgs.Find(m.Name); exists == nil {
				msgs = append(msgs, m)
			}
		}
	}

	msgs = append(msgs, MessageTypeMeta, MessageObjectMeta)

	return msgs
}

func (g *DeepCopyGenerator) getObjectDescriptors() []*descriptorpb.DescriptorProto {
	var objects []*descriptorpb.DescriptorProto
	for _, v := range g.file.GetMessageType() {
		e := proto.GetExtension(v.GetOptions(), kubeproto.E_Kind)
		if e == nil {
			continue
		}
		ext := e.(*kubeproto.Kind)
		if ext == nil {
			continue
		}
		objects = append(objects, v)
	}

	return objects
}
