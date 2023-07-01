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

var protoreflectKindMap = map[protoreflect.Kind]string{
	protoreflect.StringKind: "string",
	protoreflect.Int64Kind:  "int64",
	protoreflect.Int32Kind:  "int",
	protoreflect.BoolKind:   "bool",
	protoreflect.BytesKind:  "[]byte",
}

var ProtoreflectKindToJSONSchemaType = map[protoreflect.Kind]string{
	protoreflect.StringKind: "string",
	protoreflect.Int64Kind:  "integer",
	protoreflect.Int32Kind:  "integer",
	protoreflect.BoolKind:   "boolean",
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
		Package: ImportPackage{Path: "go.f110.dev/kubeproto/go/apis/metav1"},
	}
	MessageObjectMeta = &Message{
		Dep:       true,
		Virtual:   true,
		Name:      ".k8s.io.apimachinery.pkg.apis.meta.v1.ObjectMeta",
		ShortName: "ObjectMeta",
		Package:   ImportPackage{Path: "go.f110.dev/kubeproto/go/apis/metav1"},
	}
	MessageListMeta = &Message{
		Dep:       true,
		Virtual:   true,
		Name:      ".k8s.io.apimachinery.pkg.apis.meta.v1.ListMeta",
		ShortName: "ListMeta",
		Package:   ImportPackage{Path: "go.f110.dev/kubeproto/go/apis/metav1"},
	}
)

type Messages []*Message

func (m Messages) FilterKind() Messages {
	kindMap := make(map[string]*Message)
	var kinds []*Message
	for _, v := range m {
		if v.Dep {
			continue
		}
		if v.messageDescriptor != nil && !isKind(v.messageDescriptor) {
			continue
		}
		kinds = append(kinds, v)
		kindMap[v.ShortName] = v
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

type ScopeType string

const (
	ScopeTypeNamespaced ScopeType = "namespaced"
	ScopeTypeCluster    ScopeType = "cluster"
)

type Message struct {
	// Dep indicates that this message is dependent
	Dep bool
	// Name is a fully qualified message name that includes a package name. (e,g, .k8s.io.apimachinery.pkg.apis.meta.v1.TypeMeta)
	Name string
	// ShortName is a name of message (e,g, TypeMeta)
	ShortName string
	// Kind indicates this message is runtime.Object.
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
	// Scope is a type of this message.
	Scope ScopeType
	// HasTypeMeta indicates this message contains TypeMeta
	HasTypeMeta bool

	fileDescriptor    protoreflect.FileDescriptor
	messageDescriptor protoreflect.MessageDescriptor
}

func NewMessageFromMessageDescriptor(m protoreflect.MessageDescriptor, f protoreflect.FileDescriptor, nsm *PackageNamespaceManager) (*Message, error) {
	var fields []*Field
	for i := 0; i < m.Fields().Len(); i++ {
		v := m.Fields().Get(i)

		var name, fieldName string
		var subResource, inline bool
		e := proto.GetExtension(v.Options(), kubeproto.E_Field)
		ext := e.(*kubeproto.Field)
		if ext != nil {
			name = ext.GetGoName()
			subResource = ext.SubResource
			if ext.ApiFieldName != "" {
				fieldName = ext.ApiFieldName
			}
			inline = ext.Inline
		}
		if name == "" {
			name = stringsutil.ToUpperCamelCase(string(v.Name()))
		}
		if fieldName == "" {
			fieldName = stringsutil.ToLowerCamelCase(string(v.Name()))
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
			FieldName:   fieldName,
			Kind:        v.Kind(),
			Repeated:    repeated,
			MessageName: messageName,
			Description: description,
			Inline:      inline,
			Embed:       inline,
			Optional:    v.HasOptionalKeyword() || v.IsMap(),
			SubResource: subResource,
			descriptor:  v,
		})
	}

	var printerColumns []*kubeproto.PrinterColumn
	messageScope := ScopeTypeNamespaced
	e := proto.GetExtension(m.Options(), kubeproto.E_Kind)
	ext := e.(*kubeproto.Kind)
	if ext != nil {
		printerColumns = ext.AdditionalPrinterColumns
		if ext.Scope == kubeproto.Scope_SCOPE_CLUSTER {
			messageScope = ScopeTypeCluster
		}
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
	if v := proto.GetExtension(fileOpt, kubeproto.E_KubeprotoGoPackage); v != nil {
		if v, ok := v.(string); ok && v != "" {
			goPackage = v
		}
	}
	if strings.HasPrefix(goPackage, "k8s.io/apimachinery") || strings.HasPrefix(goPackage, "k8s.io/api") {
		s := strings.Split(goPackage, "/")
		goPackageAlias = fmt.Sprintf("%s%s", s[len(s)-2], s[len(s)-1])
	} else if !strings.HasPrefix(string(f.Package()), "google.protobuf") && !strings.HasPrefix(string(f.Package()), "k8s.io.api") {
		if i := strings.LastIndex(string(f.Package()), "."); i > 0 {
			goPackageAlias = string(f.Package()[i+1:])
		}
	}
	goPackageAlias = nsm.Add(goPackage, goPackageAlias)
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
		msg.Scope = messageScope
	}
	for _, v := range msg.Fields {
		if strings.HasSuffix(v.MessageName, MessageTypeMeta.Name[1:]) {
			msg.HasTypeMeta = true
			break
		}
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

func (m *Message) IsDefinedSubResource() bool {
	for _, f := range m.Fields {
		if f.SubResource {
			return true
		}
	}

	return false
}

func (m *Message) IsList() bool {
	if len(m.Fields) == 1 && m.Fields[0].Name == "Items" && m.Fields[0].Repeated {
		return true
	}
	return false
}

func (m *Message) ClientName(fqdn bool) string {
	if fqdn && m.Group != "" && m.Group != "." {
		return fmt.Sprintf("%s%s", stringsutil.ToUpperCamelCase(m.Group), stringsutil.ToUpperCamelCase(m.Version))
	}
	if m.SubGroup != "" {
		return fmt.Sprintf("%s%s", stringsutil.ToUpperCamelCase(m.SubGroup), stringsutil.ToUpperCamelCase(m.Version))
	}

	// If Group is empty, we assume that the message belongs to core group (e.g. Pod)
	return fmt.Sprintf("Core%s", stringsutil.ToUpperCamelCase(m.Version))
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
				Kind:        protoreflect.MessageKind,
				Inline:      true,
				Embed:       true,
			},
			{
				Name:        "ObjectMeta",
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
	Kind      protoreflect.Kind
	// Repeated indicates that this field is an array.
	Repeated bool
	// MessageName is a name of Message if Kind is MessageKind
	MessageName string
	// Description is a string of an account of this field
	Description string
	// Inline indicates the embed field
	Inline bool
	// Optional indicates that this field is an optional field.
	Optional bool
	// Embed indicates that this field is embed
	Embed bool
	// SubResource indicates that this field is the sub resource of Kind
	SubResource bool

	importPath   string
	packageAlias string
	typeName     string
	descriptor   protoreflect.FieldDescriptor
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

func (f *Field) IsMap() bool {
	if f.descriptor == nil {
		return false
	}
	return f.descriptor.IsMap()
}

func (f *Field) MapKeyValue() (protoreflect.FieldDescriptor, protoreflect.FieldDescriptor) {
	return f.descriptor.MapKey(), f.descriptor.MapValue()
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
