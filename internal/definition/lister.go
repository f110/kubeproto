package definition

import (
	"fmt"
	"path/filepath"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type Lister struct {
	files    map[string]struct{}
	allFiles *protoregistry.Files

	messages Messages
	enums    Enums
}

func NewLister(files []string, all *protoregistry.Files) *Lister {
	m := make(map[string]struct{})
	for _, v := range files {
		m[v] = struct{}{}
	}
	return &Lister{files: m, allFiles: all}
}

func (l *Lister) GetMessages() Messages {
	if l.messages != nil {
		return l.messages
	}

	kinds := make(map[string]struct{})
	l.allFiles.RangeFiles(func(desc protoreflect.FileDescriptor) bool {
		for i := 0; i < desc.Messages().Len(); i++ {
			m := desc.Messages().Get(i)
			if isKind(m) {
				kinds[string(m.Name())] = struct{}{}
			}
		}
		return true
	})

	var msgs Messages
	addMessage := func(m protoreflect.MessageDescriptor, desc protoreflect.FileDescriptor) bool {
		msg, err := NewMessageFromMessageDescriptor(m, desc)
		if err != nil {
			return false
		}
		if _, ok := l.files[desc.Path()]; !ok {
			msg.Dep = true
		}
		if exists := msgs.Find(msg.Name); exists != nil {
			return true
		}
		msgs = append(msgs, msg)

		if _, ok := kinds[msg.ShortName+"List"]; !ok && msg.Kind {
			listMessage := &Message{
				Name:      fmt.Sprintf("%sList", msg.Name),
				ShortName: fmt.Sprintf("%sList", msg.ShortName),
				Fields: []*Field{
					{
						Name:        "TypeMeta",
						MessageName: ".k8s.io.apimachinery.pkg.apis.meta.v1.TypeMeta",
						Kind:        protoreflect.MessageKind,
						Inline:      true,
						Embed:       true,
					},
					{
						Name:        "ListMeta",
						FieldName:   "metadata",
						MessageName: ".k8s.io.apimachinery.pkg.apis.meta.v1.ListMeta",
						Kind:        protoreflect.MessageKind,
						Embed:       true,
					},
					{
						Name:        "Items",
						FieldName:   "items",
						Kind:        protoreflect.MessageKind,
						Repeated:    true,
						MessageName: msg.Name,
					},
				},
				Kind:    true,
				Virtual: true,
				Dep:     msg.Dep,
			}
			msgs = append(msgs, listMessage)
		}

		return true
	}

	l.allFiles.RangeFiles(func(desc protoreflect.FileDescriptor) bool {
		for i := 0; i < desc.Messages().Len(); i++ {
			m := desc.Messages().Get(i)
			if !addMessage(m, desc) {
				return false
			}
			// For nested message declarations
			for i := 0; i < m.Messages().Len(); i++ {
				if !addMessage(m.Messages().Get(i), desc) {
					return false
				}
			}
		}

		return true
	})

	msgs = append(msgs, MessageTypeMeta, MessageObjectMeta, MessageListMeta)
	l.messages = msgs
	return msgs
}

func (l *Lister) GetEnums() Enums {
	if l.enums != nil {
		return l.enums
	}

	var enums []*Enum
	l.allFiles.RangeFiles(func(desc protoreflect.FileDescriptor) bool {
		_, own := l.files[desc.Path()]

		for i := 0; i < desc.Enums().Len(); i++ {
			e := desc.Enums().Get(i)
			enums = append(enums, NewEnumFromEnumDescriptor(e, desc, !own))
		}

		return true
	})

	l.enums = enums
	return enums
}

func (l *Lister) ResolveGoType(packageName string, f *Field) string {
	if f.typeName != "" {
		return f.typeName
	}

	if f.IsMap() {
		key, value := f.MapKeyValue()
		keyTyp := protoreflectKindMap[key.Kind()]
		valTyp := protoreflectKindMap[value.Kind()]
		return fmt.Sprintf("map[%s]%s", keyTyp, valTyp)
	}

	var typ string
	switch f.Kind {
	case protoreflect.MessageKind:
		m := l.GetMessages().Find(f.MessageName)
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
	case protoreflect.EnumKind:
		e := l.GetEnums().Find(f.MessageName)
		typ = e.ShortName
	default:
		typ = protoreflectKindMap[f.Kind]
	}

	if f.Repeated {
		typ = "[]" + typ
	}

	f.typeName = typ
	return typ
}
