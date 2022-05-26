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

	var msgs Messages
	addMessage := func(m protoreflect.MessageDescriptor, desc protoreflect.FileDescriptor) bool {
		msg, err := NewMessageFromMessageDescriptor(m, desc)
		if err != nil {
			return false
		}
		if _, ok := l.files[desc.Path()]; !ok {
			msg.Dep = true
		}
		if exists := msgs.Find(msg.Name); exists == nil {
			msgs = append(msgs, msg)
		}

		return true
	}
	l.allFiles.RangeFiles(func(desc protoreflect.FileDescriptor) bool {
		for i := 0; i < desc.Messages().Len(); i++ {
			m := desc.Messages().Get(i)
			if !addMessage(m, desc) {
				return false
			}
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
