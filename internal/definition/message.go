package definition

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"go.f110.dev/kubeproto"
	"go.f110.dev/kubeproto/internal/stringsutil"
)

var descriptorTypeMap = map[protoreflect.Kind]string{
	protoreflect.StringKind: "string",
	protoreflect.Int64Kind:  "int64",
	protoreflect.Int32Kind:  "int",
	protoreflect.BoolKind:   "bool",
}

var (
	MessageTypeMeta = &Message{
		Dep:       true,
		Virtual:   true,
		Name:      ".k8s.io.apimachinery.pkg.apis.meta.v1.TypeMeta",
		ShortName: "TypeMeta",
		Fields: []*Field{
			{
				Name:        "kind",
				FieldName:   "kind",
				Optional:    true,
				Kind:        protoreflect.StringKind,
				Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated.",
			},
			{
				Name:        "api_version",
				FieldName:   "apiVersion",
				Optional:    true,
				Kind:        protoreflect.StringKind,
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
		if !isKind(v.messageDescriptor) {
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
						Name:        "type_meta",
						MessageName: ".k8s.io.apimachinery.pkg.apis.meta.v1.TypeMeta",
						Kind:        protoreflect.MessageKind,
						Inline:      true,
						Embed:       true,
					},
					{
						Name:        "list_meta",
						FieldName:   "metadata",
						MessageName: ".k8s.io.apimachinery.pkg.apis.meta.v1.ListMeta",
						Kind:        protoreflect.MessageKind,
						Embed:       true,
					},
					{
						Name:        "items",
						FieldName:   "items",
						Kind:        protoreflect.MessageKind,
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

	sort.Slice(kinds, func(i, j int) bool {
		return kinds[i].ShortName < kinds[j].ShortName
	})
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

func isKind(desc protoreflect.MessageDescriptor) bool {
	e := proto.GetExtension(desc.Options(), kubeproto.E_Kind)
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

	fileDescriptor    protoreflect.FileDescriptor
	messageDescriptor protoreflect.MessageDescriptor
}

func NewMessageFromMessageDescriptor(m protoreflect.MessageDescriptor, f protoreflect.FileDescriptor) (*Message, error) {
	var fields []*Field
	for i := 0; i < m.Fields().Len(); i++ {
		v := m.Fields().Get(i)

		var name string
		var subResource bool
		e := proto.GetExtension(v.Options(), kubeproto.E_Field)
		ext := e.(*kubeproto.Field)
		if ext != nil {
			name = ext.GetGoName()
			subResource = ext.SubResource
		}
		if name == "" {
			name = string(v.Name())
		}

		repeated := v.IsList()
		var description string
		if location := f.SourceLocations().ByDescriptor(v); location.LeadingComments != "" {
			description = strings.TrimSuffix(strings.TrimPrefix(location.LeadingComments, " "), "\n")
		}

		var messageName string
		switch v.Kind() {
		case protoreflect.MessageKind:
			messageName = string(v.Message().FullName())
		case protoreflect.EnumKind:
			messageName = string(v.Enum().FullName())
		}
		fields = append(fields, &Field{
			Name:        Name(name),
			FieldName:   stringsutil.ToLowerCamelCase(string(v.Name())),
			Kind:        v.Kind(),
			Repeated:    repeated,
			MessageName: messageName,
			Description: description,
			Optional:    v.HasOptionalKeyword(),
			SubResource: subResource,
		})
	}

	var printerColumns []*kubeproto.PrinterColumn
	e := proto.GetExtension(m.Options(), kubeproto.E_Kind)
	ext := e.(*kubeproto.Kind)
	if ext != nil {
		printerColumns = ext.AdditionalPrinterColumns
	}

	var group, subGroup, version string
	e = proto.GetExtension(f.Options(), kubeproto.E_K8S)
	k8sExt := e.(*kubeproto.Kubernetes)
	if k8sExt != nil {
		group = fmt.Sprintf("%s.%s", k8sExt.SubGroup, k8sExt.Domain)
		subGroup = k8sExt.SubGroup
		version = k8sExt.Version
	}

	fileOpt := f.Options()
	var goPackage, goPackageAlias string
	if v, ok := fileOpt.(*descriptorpb.FileOptions); ok {
		goPackage = v.GetGoPackage()
	}
	if strings.HasPrefix(goPackage, "k8s.io/apimachinery") || strings.HasPrefix(goPackage, "k8s.io/api") {
		s := strings.Split(goPackage, "/")
		goPackageAlias = fmt.Sprintf("%s%s", s[len(s)-2], s[len(s)-1])
	}
	msg := &Message{
		Name:                     string(m.FullName()),
		ShortName:                string(m.Name()),
		Fields:                   fields,
		AdditionalPrinterColumns: printerColumns,
		Group:                    group,
		SubGroup:                 subGroup,
		Version:                  version,
		Package: ImportPackage{
			Name:  path.Base(goPackage),
			Path:  goPackage,
			Alias: goPackageAlias,
		},
		fileDescriptor:    f,
		messageDescriptor: m,
	}

	if isKind(m) {
		extendAsKind(msg)
	}
	return msg, nil
}

func (m *Message) Kubernetes() (*kubeproto.Kubernetes, error) {
	e := proto.GetExtension(m.fileDescriptor.Options(), kubeproto.E_K8S)
	ext := e.(*kubeproto.Kubernetes)
	if ext == nil {
		return nil, fmt.Errorf("%s is not extended by kubeproto.Kubernetes", m.ShortName)
	}

	return ext, nil
}

func isEqualProtoPath(a, b []int32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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
				Name:        "type_meta",
				MessageName: ".k8s.io.apimachinery.pkg.apis.meta.v1.TypeMeta",
				Kind:        protoreflect.MessageKind,
				Inline:      true,
				Embed:       true,
			},
			{
				Name:        "object_meta",
				FieldName:   "metadata",
				MessageName: ".k8s.io.apimachinery.pkg.apis.meta.v1.ObjectMeta",
				Kind:        protoreflect.MessageKind,
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
	Kind      protoreflect.Kind
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
