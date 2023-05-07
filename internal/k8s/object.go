package k8s

import (
	"bufio"
	"io"
	"os"
	"path"
	"sort"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	"go.f110.dev/kubeproto"
	"go.f110.dev/kubeproto/internal/codegeneration"
	"go.f110.dev/kubeproto/internal/definition"
	"go.f110.dev/kubeproto/internal/stringsutil"
)

type ObjectGenerator struct {
	file                    protoreflect.FileDescriptor
	lister                  *definition.Lister
	packageNamespaceManager *definition.PackageNamespaceManager
}

func NewObjectGenerator(fileToGenerate []string, files *protoregistry.Files) (*ObjectGenerator, error) {
	desc, err := files.FindFileByPath(fileToGenerate[0])
	if err != nil {
		return nil, err
	}
	nsm := definition.NewPackageNamespaceManager()

	return &ObjectGenerator{
		file:                    desc.(protoreflect.FileDescriptor),
		lister:                  definition.NewLister(fileToGenerate, files, nsm),
		packageNamespaceManager: nsm,
	}, nil
}

func (g *ObjectGenerator) Generate(out io.Writer) error {
	importPackages := map[string]string{
		"k8s.io/apimachinery/pkg/runtime": "",
	}
	for k, v := range importPackages {
		g.packageNamespaceManager.Add(k, v)
	}

	w := codegeneration.NewWriter()
	messages := g.lister.GetMessages()

	defW := codegeneration.NewWriter()
	fileOpt := g.file.Options().(*descriptorpb.FileOptions)
	e := proto.GetExtension(g.file.Options(), kubeproto.E_K8S)
	ext := e.(*kubeproto.Kubernetes)
	hasRuntimeObject := ext != nil

	packageName := fileOpt.GetGoPackage()
	kubeprotoGoPackage := proto.GetExtension(g.file.Options(), kubeproto.E_KubeprotoGoPackage)
	if kubeprotoGoPackage != nil {
		packageName = kubeprotoGoPackage.(string)
	}
	w.F("package %s", path.Base(packageName))

	if hasRuntimeObject {
		importPackages["k8s.io/apimachinery/pkg/runtime/schema"] = ""
		importPackages["go.f110.dev/kubeproto/go/apis/metav1"] = ""

		if ext.SubGroup != "" {
			defW.F("const GroupName = \"%s.%s\"", ext.SubGroup, ext.Domain)
		} else {
			defW.F("const GroupName = \"%s\"", ext.Domain)
		}
		defW.F("")
		defW.F("var (")
		defW.F("GroupVersion = metav1.GroupVersion{Group: GroupName, Version: %q}", ext.Version)
		defW.F("SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)")
		defW.F("AddToScheme = SchemeBuilder.AddToScheme")
		defW.F("SchemaGroupVersion = schema.GroupVersion{Group: GroupName, Version: %q}", ext.Version)
		defW.F(")")

		defW.F("func addKnownTypes(scheme *runtime.Scheme) error {")
		defW.F("scheme.AddKnownTypes(SchemaGroupVersion,")
		for _, m := range messages.FilterKind() {
			defW.F("&%s{},", m.ShortName)
		}
		defW.F(")")
		defW.F("metav1.AddToGroupVersion(scheme, SchemaGroupVersion)")
		defW.F("return nil")
		defW.F("}")
	}

	enums := g.lister.GetEnums()
	for _, enum := range enums.Own() {
		if len(enum.Values) == 0 {
			continue
		}

		// Enum definition
		defW.F("type %s string", enum.ShortName)
		defW.F("const (")
		for _, v := range enum.Values {
			defW.F("%s%s %s = %q", enum.ShortName, stringsutil.Letterize(stringsutil.ToUpperCamelCase(v)), enum.ShortName, v)
		}
		defW.F(")")
	}

	mark := make(map[string]struct{})
	var objs definition.Messages
	if os.Getenv("KUBEPROTO_OPTS") == "all" {
		objMap := make(map[string]*definition.Message)
		for _, v := range messages {
			if !v.Dep {
				objMap[v.ShortName] = v
			}
		}
		for _, v := range objMap {
			objs = append(objs, v)
		}
		sort.Slice(objs, func(i, j int) bool {
			return objs[i].ShortName < objs[j].ShortName
		})
	} else {
		objs = messages.FilterKind()
	}
	for len(objs) > 0 {
		obj := objs[0]
		if _, ok := mark[obj.Name]; ok {
			objs = objs[1:]
			continue
		}
		mark[obj.Name] = struct{}{}

		if obj.IsList() {
			importPath, packageAlias, typ := g.lister.ResolveGoType(packageName, obj.Fields[0])
			if importPath != "" {
				importPackages[importPath] = packageAlias
			}
			defW.F("type %s %s", obj.ShortName, typ)
		} else {
			// Struct definition
			defW.F("type %s struct {", obj.ShortName)
			for _, f := range obj.Fields {
				switch f.Kind {
				case protoreflect.MessageKind:
					m := messages.Find(f.MessageName)
					if m != nil {
						if _, ok := mark[m.Name]; !ok && !m.Dep {
							objs = append(objs, m)
						}
					}

					if f.IsMap() {
						key, value := f.MapKeyValue()
						if key.Kind() == protoreflect.MessageKind {
							m := messages.Find(string(key.Message().FullName()))
							if m != nil {
								if _, ok := mark[m.Name]; !ok && !m.Dep {
									objs = append(objs, m)
								}
							}
						}
						if value.Kind() == protoreflect.MessageKind {
							m := messages.Find(string(value.Message().FullName()))
							if m != nil {
								if _, ok := mark[m.Name]; !ok && !m.Dep {
									objs = append(objs, m)
								}
							}
						}
					}
				}

				var name string
				if !f.Embed {
					name = string(f.Name)
				}
				importPath, packageName, typ := g.lister.ResolveGoType(packageName, f)
				if importPath != "" {
					importPackages[importPath] = packageName
				}
				tag := f.Tag()
				if f.Description != "" {
					d := strings.Replace(f.Description, string(f.Name), f.Name.CamelCase(), 1)
					scanner := bufio.NewScanner(strings.NewReader(d))
					for scanner.Scan() {
						line := scanner.Text()
						line = strings.TrimSpace(line)
						defW.F("// %s", line)
					}
				}
				defW.F("%s %s %s", name, typ, tag)
			}
			defW.F("}")
			defW.F("")

			// DeepCopy functions (DeepCopyInto / DeepCopy / DeepCopyObject)
			defW.F("func (in *%s) DeepCopyInto(out *%s) {", obj.ShortName, obj.ShortName)
			defW.F("*out = *in")
			for _, f := range obj.Fields {
				switch f.Kind {
				case protoreflect.MessageKind:
					if f.Repeated {
						defW.F("if in.%s != nil {", f.Name)
						_, _, typ := g.lister.ResolveGoType(packageName, f)
						defW.F("l := make(%s, len(in.%s))", typ, f.Name)
						defW.F("for i := range in.%s {", f.Name)
						defW.F("in.%s[i].DeepCopyInto(&l[i])", f.Name)
						defW.F("}")
						defW.F("out.%s = l", f.Name)
						defW.F("}")
						continue
					}
					if f.Optional {
						defW.F("if in.%s != nil {", f.Name)
						if f.IsMap() {
							defW.F("in, out := &in.%s, &out.%s", f.Name, f.Name)
							_, _, typ := g.lister.ResolveGoType(packageName, f)
							defW.F("*out = make(%s, len(*in))", typ)
							defW.F("for k, v := range *in {")
							defW.F("(*out)[k] = v")
							defW.F("}")
						} else {
							defW.F("in, out := &in.%s, &out.%s", f.Name, f.Name)
							importPath, _, typ := g.lister.ResolveGoType(packageName, f)
							defW.F("*out = new(%s)", typ[1:])
							switch importPath {
							case "k8s.io/apimachinery/pkg/util/intstr":
								defW.F("*out = *in")
							default:
								defW.F("(*in).DeepCopyInto(*out)")
							}
						}
						defW.F("}")
						continue
					}
					if f.Inline {
						_, alias, typ := g.lister.ResolveGoType(packageName, f)
						if alias != "" {
							typ = strings.TrimPrefix(typ, alias+".")
						}
						defW.F("out.%s = in.%s", typ, typ)
					} else {
						importPath, _, _ := g.lister.ResolveGoType(packageName, f)
						switch importPath {
						case "k8s.io/apimachinery/pkg/util/intstr":
							defW.F("in = out")
						default:
							defW.F("in.%s.DeepCopyInto(&out.%s)", f.Name, f.Name)
						}
					}
				default:
					if f.Repeated {
						defW.F("if in.%s != nil {", f.Name)
						_, _, typ := g.lister.ResolveGoType(packageName, f)
						defW.F("t := make(%s, len(in.%s))", typ, f.Name)
						defW.F("copy(t, in.%s)", f.Name)
						defW.F("out.%s = t", f.Name)
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
			// TypeMeta partially implements runtime.Object.
			if obj.HasTypeMeta {
				defW.F("func (in *%s) DeepCopyObject() runtime.Object {", obj.ShortName)
				defW.F("if c := in.DeepCopy(); c != nil {")
				defW.F("return c")
				defW.F("}")
				defW.F("return nil")
				defW.F("}")
			}
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
