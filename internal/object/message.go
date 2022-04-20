package object

import (
	"fmt"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	"go.f110.dev/kubeproto"
	"go.f110.dev/kubeproto/internal/stringsutil"
)

var descriptorTypeMap = map[descriptorpb.FieldDescriptorProto_Type]string{
	descriptorpb.FieldDescriptorProto_TYPE_STRING: "string",
	descriptorpb.FieldDescriptorProto_TYPE_INT64:  "int64",
	descriptorpb.FieldDescriptorProto_TYPE_INT32:  "int",
	descriptorpb.FieldDescriptorProto_TYPE_BOOL:   "bool",
}

type Messages []*Message

func (m Messages) FilterKind() Messages {
	var kinds []*Message
	for _, v := range m {
		if v.Dep {
			continue
		}
		if !isKind(v.descriptor) {
			continue
		}
		kinds = append(kinds, v)
	}

	return kinds
}

func (m Messages) Find(name string) *Message {
	for _, v := range m {
		if v.Name == name {
			return v
		}
	}

	return nil
}

func isKind(desc *descriptorpb.DescriptorProto) bool {
	e := proto.GetExtension(desc.GetOptions(), kubeproto.E_Kind)
	if e == nil {
		return false
	}
	ext := e.(*kubeproto.Kind)
	if ext == nil {
		return false
	}

	return true
}

type Message struct {
	Dep       bool
	Name      string
	ShortName string
	Fields    Fields
	Virtual   bool
	Package   ImportPackage

	descriptor *descriptorpb.DescriptorProto
}

var (
	MessageTypeMeta = &Message{
		Dep:       true,
		Virtual:   true,
		Name:      ".k8s.io.apimachinery.pkg.apis.meta.v1.TypeMeta",
		ShortName: "TypeMeta",
		Package:   ImportPackage{Path: "k8s.io/apimachinery/pkg/apis/meta/v1", Alias: "metav1"},
	}
	MessageObjectMeta = &Message{
		Dep:       true,
		Virtual:   true,
		Name:      ".k8s.io.apimachinery.pkg.apis.meta.v1.ObjectMeta",
		ShortName: "ObjectMeta",
		Package:   ImportPackage{Path: "k8s.io/apimachinery/pkg/apis/meta/v1", Alias: "metav1"},
	}
)

func NewMessage(f *descriptorpb.FileDescriptorProto, desc *descriptorpb.DescriptorProto) *Message {
	var fields []*Field
	for _, v := range desc.Field {
		fields = append(fields, &Field{
			Name:        Name(v.GetName()),
			FieldName:   stringsutil.ToLowerCamelCase(v.GetName()),
			Type:        v.GetType(),
			MessageName: v.GetTypeName(),
			Optional:    v.GetProto3Optional(),
		})
	}

	if isKind(desc) {
		fields = append([]*Field{
			{
				Name:        "TypeMeta",
				MessageName: ".k8s.io.apimachinery.pkg.apis.meta.v1.TypeMeta",
				Type:        descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
				Inline:      true,
			},
			{
				Name:        "ObjectMeta",
				FieldName:   "metadata",
				MessageName: ".k8s.io.apimachinery.pkg.apis.meta.v1.ObjectMeta",
				Type:        descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
			},
		}, fields...)
	}

	return &Message{
		Name:       fmt.Sprintf(".%s.%s", f.GetPackage(), desc.GetName()),
		ShortName:  desc.GetName(),
		Fields:     fields,
		descriptor: desc,
	}
}

type Field struct {
	// Name is a struct field name
	Name Name
	// FieldName is a json tag name
	FieldName string
	Type      descriptorpb.FieldDescriptorProto_Type
	// MessageName is a name of Message if Type is FieldDescriptorProto_TYPE_MESSAGE
	MessageName string
	// Inline indicates the embed field
	Inline   bool
	Optional bool

	typeName string
}

func (f *Field) TypeName(messages Messages) string {
	if f.typeName != "" {
		return f.typeName
	}

	var typ string
	switch f.Type {
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		m := messages.Find(f.MessageName)
		if m == nil {
			return ""
		}
		if m.Package.Path != "" {
			alias := m.Package.Alias
			if alias == "" {
				alias = filepath.Base(m.Package.Path)
			}
			typ = fmt.Sprintf("%s.%s", m.Package.Alias, m.ShortName)
		} else {
			typ = m.ShortName
		}
		if f.Optional {
			typ = "*" + typ
		}
	default:
		typ = descriptorTypeMap[f.Type]
	}

	f.typeName = typ
	return typ
}

func (f *Field) Tag() string {
	if f.Inline {
		return "`json:\",inline\"`"
	}

	s := strings.Builder{}
	s.WriteString("`json:\"")
	if f.FieldName != "" {
		s.WriteString(f.FieldName)
	}
	if f.Optional {
		s.WriteString(",omitempty")
	}
	s.WriteString("\"`")

	// `json:""` is effectively an empty tag
	if s.String() == "`json:\"\"`" {
		return ""
	}
	return s.String()
}

type Fields []*Field

type ImportPackage struct {
	Path  string
	Alias string
}

type Name string

func (n Name) CamelCase() string {
	return stringsutil.ToUpperCamelCase(string(n))
}
