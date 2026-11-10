package k8s

import (
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
	"gopkg.in/yaml.v2"
	"sigs.k8s.io/structured-merge-diff/v6/schema"

	"go.f110.dev/kubeproto"
	"go.f110.dev/kubeproto/internal/definition"
)

const (
	untypedAtomicTypeName  = "__untyped_atomic_"
	untypedDeducedTypeName = "__untyped_deduced_"
)

// messagesAsScalar is the messages that are encoded as a scalar value in JSON.
var messagesAsScalar = map[string]schema.Scalar{
	"k8s.io.apimachinery.pkg.apis.meta.v1.Time":      schema.String,
	"k8s.io.apimachinery.pkg.apis.meta.v1.MicroTime": schema.String,
	"k8s.io.apimachinery.pkg.apis.meta.v1.Duration":  schema.String,
	"google.protobuf.Timestamp":                      schema.String,
	"k8s.io.apimachinery.pkg.api.resource.Quantity":  schema.Scalar("untyped"),
	"k8s.io.apimachinery.pkg.util.intstr.IntOrString": schema.Scalar("untyped"),
}

// messagesAsUntyped is the messages that hold arbitrary JSON.
var messagesAsUntyped = map[string]string{
	"k8s.io.apimachinery.pkg.runtime.RawExtension":                 untypedAtomicTypeName,
	"k8s.io.apiextensions_apiserver.pkg.apis.apiextensions.v1.JSON": untypedAtomicTypeName,
	"k8s.io.apimachinery.pkg.apis.meta.v1.FieldsV1":                untypedDeducedTypeName,
}

// typeSchemaBuilder builds the schema of structured-merge-diff from the messages.
// The schema is used by the field manager to merge the objects in the same way as Server-Side Apply.
type typeSchemaBuilder struct {
	messages definition.Messages
	types    map[string]schema.TypeDef
	order    []string
}

func newTypeSchemaBuilder(messages definition.Messages) *typeSchemaBuilder {
	return &typeSchemaBuilder{messages: messages, types: make(map[string]schema.TypeDef)}
}

// Build returns the schema in YAML that contains the kinds and all messages referenced from them.
func (b *typeSchemaBuilder) Build(kinds []*definition.Message) (string, error) {
	for _, m := range kinds {
		b.addMessage(m)
	}

	s := schema.Schema{}
	for _, name := range b.order {
		s.Types = append(s.Types, b.types[name])
	}
	s.Types = append(s.Types, untypedTypeDefs()...)
	buf, err := yaml.Marshal(&s)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// typeName returns the name of the type in the schema for the message.
// It is the same as OpenAPIModelName of the generated type so that managedfields.NewSchemeTypeConverter can find the type.
func (b *typeSchemaBuilder) typeName(m *definition.Message) string {
	return strings.TrimPrefix(m.Name, ".")
}

func (b *typeSchemaBuilder) addMessage(m *definition.Message) {
	name := b.typeName(m)
	if _, ok := b.types[name]; ok {
		return
	}
	// Reserve the name before resolving the fields because the message may refer itself.
	b.types[name] = schema.TypeDef{Name: name}
	b.order = append(b.order, name)

	typeMap := &schema.Map{Fields: b.structFields(m)}
	if m.StructType == kubeproto.StructType_STRUCT_TYPE_ATOMIC {
		typeMap.ElementRelationship = schema.Atomic
	}
	b.types[name] = schema.TypeDef{Name: name, Atom: schema.Atom{Map: typeMap}}
}

func (b *typeSchemaBuilder) structFields(m *definition.Message) []schema.StructField {
	var fields []schema.StructField
	for _, f := range m.Fields {
		if f.Inline && f.Kind == protoreflect.MessageKind && !f.IsMap() {
			if child := b.findMessage(f.MessageName); child != nil {
				fields = append(fields, b.structFields(child)...)
			}
			continue
		}
		if f.FieldName == "" || f.FieldName == "-" {
			continue
		}
		fields = append(fields, schema.StructField{Name: f.FieldName, Type: b.fieldTypeRef(f)})
	}
	return fields
}

func (b *typeSchemaBuilder) fieldTypeRef(f *definition.Field) schema.TypeRef {
	if f.IsMap() {
		_, value := f.MapKeyValue()
		typeMap := &schema.Map{ElementType: b.descriptorTypeRef(value)}
		if f.MapType == kubeproto.MapType_MAP_TYPE_ATOMIC {
			typeMap.ElementRelationship = schema.Atomic
		}
		return schema.TypeRef{Inlined: schema.Atom{Map: typeMap}}
	}

	var elem schema.TypeRef
	switch f.Kind {
	case protoreflect.MessageKind:
		elem = b.messageTypeRef(f.MessageName)
	default:
		elem = scalarTypeRef(f.Kind)
	}
	if !f.Repeated {
		return elem
	}

	list := &schema.List{ElementType: elem, ElementRelationship: schema.Atomic}
	switch f.ListType {
	case kubeproto.ListType_LIST_TYPE_SET:
		list.ElementRelationship = schema.Associative
	case kubeproto.ListType_LIST_TYPE_MAP:
		list.ElementRelationship = schema.Associative
		list.Keys = f.ListMapKeys
	}
	return schema.TypeRef{Inlined: schema.Atom{List: list}}
}

func (b *typeSchemaBuilder) descriptorTypeRef(d protoreflect.FieldDescriptor) schema.TypeRef {
	if d.Kind() == protoreflect.MessageKind {
		return b.messageTypeRef(string(d.Message().FullName()))
	}
	return scalarTypeRef(d.Kind())
}

func (b *typeSchemaBuilder) messageTypeRef(messageName string) schema.TypeRef {
	name := strings.TrimPrefix(messageName, ".")
	if s, ok := messagesAsScalar[name]; ok {
		return schema.TypeRef{Inlined: schema.Atom{Scalar: &s}}
	}
	if v, ok := messagesAsUntyped[name]; ok {
		return schema.TypeRef{NamedType: &v}
	}

	child := b.findMessage(name)
	if child == nil {
		n := untypedDeducedTypeName
		return schema.TypeRef{NamedType: &n}
	}
	b.addMessage(child)
	n := b.typeName(child)
	return schema.TypeRef{NamedType: &n}
}

// findMessage returns the message that is defined in proto.
// The virtual messages (e.g. ObjectMeta for Kind) are not returned because they don't have the fields.
func (b *typeSchemaBuilder) findMessage(name string) *definition.Message {
	m := b.messages.Find(strings.TrimPrefix(name, "."))
	if m == nil || m.Virtual {
		return nil
	}
	return m
}

func scalarTypeRef(kind protoreflect.Kind) schema.TypeRef {
	var s schema.Scalar
	switch kind {
	case protoreflect.BoolKind:
		s = schema.Boolean
	case protoreflect.StringKind, protoreflect.BytesKind, protoreflect.EnumKind:
		s = schema.String
	default:
		s = schema.Numeric
	}
	return schema.TypeRef{Inlined: schema.Atom{Scalar: &s}}
}

// untypedTypeDefs returns the types for the values that don't have the schema.
// These are the same as the types that structured-merge-diff and Kubernetes use.
func untypedTypeDefs() []schema.TypeDef {
	untyped := schema.Scalar("untyped")
	atomic, deduced := untypedAtomicTypeName, untypedDeducedTypeName
	return []schema.TypeDef{
		{
			Name: untypedAtomicTypeName,
			Atom: schema.Atom{
				Scalar: &untyped,
				List:   &schema.List{ElementType: schema.TypeRef{NamedType: &atomic}, ElementRelationship: schema.Atomic},
				Map:    &schema.Map{ElementType: schema.TypeRef{NamedType: &atomic}, ElementRelationship: schema.Atomic},
			},
		},
		{
			Name: untypedDeducedTypeName,
			Atom: schema.Atom{
				Scalar: &untyped,
				List:   &schema.List{ElementType: schema.TypeRef{NamedType: &atomic}, ElementRelationship: schema.Atomic},
				Map:    &schema.Map{ElementType: schema.TypeRef{NamedType: &deduced}, ElementRelationship: schema.Separable},
			},
		},
	}
}
