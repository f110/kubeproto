package definition

import (
	"fmt"
	"path/filepath"

	"google.golang.org/protobuf/types/descriptorpb"
)

type Lister struct {
	files    []*descriptorpb.FileDescriptorProto
	allFiles []*descriptorpb.FileDescriptorProto

	messages Messages
	enums    Enums
}

func NewLister(files []*descriptorpb.FileDescriptorProto, allProtos []*descriptorpb.FileDescriptorProto) *Lister {
	return &Lister{files: files, allFiles: allProtos}
}

func (l *Lister) GetMessages() Messages {
	if l.messages != nil {
		return l.messages
	}

	var msgs Messages
	for _, f := range l.files {
		for _, v := range f.GetMessageType() {
			msgs = append(msgs, NewMessage(f, v))
		}
	}
	for _, v := range l.allFiles {
		for _, mt := range v.GetMessageType() {
			m := NewMessage(v, mt)
			m.Dep = true
			if exists := msgs.Find(m.Name); exists == nil {
				msgs = append(msgs, m)
			}
		}
	}

	msgs = append(msgs, MessageTypeMeta, MessageObjectMeta, MessageListMeta)
	l.messages = msgs
	return msgs
}

func (l *Lister) GetEnums() Enums {
	if l.enums != nil {
		return l.enums
	}

	var enums []*Enum
	for _, f := range l.files {
		for _, v := range f.GetEnumType() {
			enums = append(enums, NewEnum(f, v))
		}
	}
	for _, f := range l.allFiles {
		for _, v := range f.GetEnumType() {
			enums = append(enums, NewEnum(f, v))
		}
	}

	l.enums = enums
	return enums
}

func (l *Lister) ResolveGoType(packageName string, f *Field) string {
	if f.typeName != "" {
		return f.typeName
	}

	var typ string
	switch f.Type {
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
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
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		e := l.GetEnums().Find(f.MessageName)
		typ = e.ShortName
	default:
		typ = descriptorTypeMap[f.Type]
	}

	if f.Repeated {
		typ = "[]" + typ
	}

	f.typeName = typ
	return typ
}
