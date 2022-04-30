package definition

import (
	"fmt"
	"path"
	"path/filepath"
	"reflect"
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

var (
	MessageTypeMeta = &Message{
		Dep:       true,
		Virtual:   true,
		Name:      ".k8s.io.apimachinery.pkg.apis.meta.v1.TypeMeta",
		ShortName: "TypeMeta",
		Fields: []*Field{
			{
				Name:        "Kind",
				FieldName:   "kind",
				Optional:    true,
				Type:        descriptorpb.FieldDescriptorProto_TYPE_STRING,
				Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated.",
			},
			{
				Name:        "APIVersion",
				FieldName:   "apiVersion",
				Optional:    true,
				Type:        descriptorpb.FieldDescriptorProto_TYPE_STRING,
				Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values.",
			},
		},
		Package: ImportPackage{Path: "k8s.io/apimachinery/pkg/apis/meta/v1", Alias: "metav1"},
	}
	MessageObjectMeta = &Message{
		Dep:       true,
		Virtual:   true,
		Name:      ".k8s.io.apimachinery.pkg.apis.meta.v1.ObjectMeta",
		ShortName: "ObjectMeta",
		Package:   ImportPackage{Path: "k8s.io/apimachinery/pkg/apis/meta/v1", Alias: "metav1"},
	}
	MessageListMeta = &Message{
		Dep:       true,
		Virtual:   true,
		Name:      ".k8s.io.apimachinery.pkg.apis.meta.v1.ListMeta",
		ShortName: "ListMeta",
		Package:   ImportPackage{Path: "k8s.io/apimachinery/pkg/apis/meta/v1", Alias: "metav1"},
	}
)

type Messages []*Message

func (m *Messages) FilterKind() Messages {
	kindMap := make(map[string]*Message)
	var kinds []*Message
	for _, v := range *m {
		if v.Dep {
			continue
		}
		if !isKind(v.descriptor) {
			continue
		}
		kinds = append(kinds, v)
		kindMap[v.ShortName] = v
	}
	for name, v := range kindMap {
		if _, ok := kindMap[name+"List"]; !ok {
			listMessage := &Message{
				Name:      fmt.Sprintf("%sList", v.Name),
				ShortName: fmt.Sprintf("%sList", v.ShortName),
				Fields: []*Field{
					{
						Name:        "TypeMeta",
						MessageName: ".k8s.io.apimachinery.pkg.apis.meta.v1.TypeMeta",
						Type:        descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
						Inline:      true,
						Embed:       true,
					},
					{
						Name:        "ListMeta",
						FieldName:   "metadata",
						MessageName: ".k8s.io.apimachinery.pkg.apis.meta.v1.ListMeta",
						Type:        descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
						Embed:       true,
					},
					{
						Name:        "Items",
						FieldName:   "items",
						Type:        descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
						Repeated:    true,
						MessageName: v.Name,
					},
				},
				Kind:    true,
				Virtual: true,
			}
			*m = append(*m, listMessage)
			kinds = append(kinds, listMessage)
		}
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
	// Dep indicates that this message is dependent
	Dep bool
	// Name is a fully qualified message name that includes a package name. (e,g, .k8s.io.apimachinery.pkg.apis.meta.v1.TypeMeta)
	Name string
	// ShortName is a name of message (e,g, TypeMeta)
	ShortName string
	// Kind indicates has this message is runtime.Object.
	Kind bool
	// Fields has all fields of the message.
	Fields Fields
	// Virtual indicates that the message is not defined protobuf.
	Virtual                  bool
	AdditionalPrinterColumns []*kubeproto.PrinterColumn
	Package                  ImportPackage
	// Group is the api group (e,g, authorization.k8s.io)
	Group    string
	SubGroup string
	// Version is the api version (e,g, v1alpha1)
	Version string

	descriptor     *descriptorpb.DescriptorProto
	fileDescriptor *descriptorpb.FileDescriptorProto
}

func NewMessage(f *descriptorpb.FileDescriptorProto, desc *descriptorpb.DescriptorProto) *Message {
	var messageIndex int
	for i, v := range f.MessageType {
		if v == desc {
			messageIndex = i
			break
		}
	}

	var fields []*Field
	for i, v := range desc.Field {
		var name string
		var subResource bool
		e := proto.GetExtension(v.GetOptions(), kubeproto.E_Field)
		ext := e.(*kubeproto.Field)
		if ext != nil {
			name = ext.GetGoName()
			subResource = ext.SubResource
		}
		if name == "" {
			name = v.GetName()
		}

		var repeated bool
		if v.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
			repeated = true
		}

		var description string
		fieldPath := []int32{4, int32(messageIndex), 2, int32(i)}
		for _, s := range f.SourceCodeInfo.GetLocation() {
			if reflect.DeepEqual(s.Path, fieldPath) {
				description = s.GetLeadingComments()
				break
			}
		}

		fields = append(fields, &Field{
			Name:        Name(name),
			FieldName:   stringsutil.ToLowerCamelCase(v.GetName()),
			Type:        v.GetType(),
			Repeated:    repeated,
			MessageName: v.GetTypeName(),
			Description: description,
			Optional:    v.GetProto3Optional(),
			SubResource: subResource,
		})
	}

	var printerColumns []*kubeproto.PrinterColumn
	e := proto.GetExtension(desc.GetOptions(), kubeproto.E_Kind)
	ext := e.(*kubeproto.Kind)
	if ext != nil {
		printerColumns = ext.AdditionalPrinterColumns
	}

	var group, subGroup, version string
	e = proto.GetExtension(f.GetOptions(), kubeproto.E_K8S)
	k8sExt := e.(*kubeproto.Kubernetes)
	if k8sExt != nil {
		group = fmt.Sprintf("%s.%s", k8sExt.SubGroup, k8sExt.Domain)
		subGroup = k8sExt.SubGroup
		version = k8sExt.Version
	}

	m := &Message{
		Name:                     fmt.Sprintf(".%s.%s", f.GetPackage(), desc.GetName()),
		ShortName:                desc.GetName(),
		Fields:                   fields,
		AdditionalPrinterColumns: printerColumns,
		Group:                    group,
		SubGroup:                 subGroup,
		Version:                  version,
		Package: ImportPackage{
			Name: path.Base(f.GetOptions().GetGoPackage()),
			Path: f.GetOptions().GetGoPackage(),
		},
		descriptor:     desc,
		fileDescriptor: f,
	}

	if isKind(desc) {
		extendAsKind(m)
	}

	return m
}

func (m *Message) Kubernetes() (*kubeproto.Kubernetes, error) {
	e := proto.GetExtension(m.fileDescriptor.GetOptions(), kubeproto.E_K8S)
	ext := e.(*kubeproto.Kubernetes)
	if ext == nil {
		return nil, fmt.Errorf("%s is not extended by kubeproto.Kubernetes", m.ShortName)
	}

	return ext, nil
}

func extendAsKind(m *Message) {
	m.Kind = true
	found := false
	for _, v := range m.Fields {
		if v.Name == "TypeMeta" {
			found = true
			break
		}
	}
	if !found {
		m.Fields = append([]*Field{
			{
				Name:        "TypeMeta",
				MessageName: ".k8s.io.apimachinery.pkg.apis.meta.v1.TypeMeta",
				Type:        descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
				Inline:      true,
				Embed:       true,
			},
			{
				Name:        "ObjectMeta",
				FieldName:   "metadata",
				MessageName: ".k8s.io.apimachinery.pkg.apis.meta.v1.ObjectMeta",
				Type:        descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
				Embed:       true,
			},
		}, m.Fields...)
	}
}

type Field struct {
	// Name is a struct field name
	Name Name
	// FieldName is a json tag name
	FieldName string
	Type      descriptorpb.FieldDescriptorProto_Type
	// Repeated indicates that this field is an array.
	Repeated bool
	// MessageName is a name of Message if Type is FieldDescriptorProto_TYPE_MESSAGE
	MessageName string
	// Description is a string of an account of this field
	Description string
	// Inline indicates the embed field
	Inline bool
	// Optional indicates that this field is an optional field.
	Optional    bool
	Embed       bool
	SubResource bool

	typeName string
}

func (f *Field) TypeName(packageName string, messages Messages) string {
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
		if m.Package.Path != "" && m.Package.Path != packageName {
			alias := m.Package.Alias
			if alias == "" {
				alias = filepath.Base(m.Package.Path)
			}
			typ = fmt.Sprintf("%s.%s", alias, m.ShortName)
		} else {
			typ = m.ShortName
		}
		if f.Optional {
			typ = "*" + typ
		}
	default:
		typ = descriptorTypeMap[f.Type]
	}

	if f.Repeated {
		typ = "[]" + typ
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
	Name  string
	Path  string
	Alias string
}

type Name string

func (n Name) CamelCase() string {
	return stringsutil.ToUpperCamelCase(string(n))
}
