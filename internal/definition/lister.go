package definition

import (
	"google.golang.org/protobuf/types/descriptorpb"
)

type Lister struct {
	file     *descriptorpb.FileDescriptorProto
	allFiles []*descriptorpb.FileDescriptorProto
}

func NewLister(file *descriptorpb.FileDescriptorProto, allProtos []*descriptorpb.FileDescriptorProto) *Lister {
	return &Lister{file: file, allFiles: allProtos}
}

func (l *Lister) GetMessages() Messages {
	var msgs Messages
	for _, v := range l.file.GetMessageType() {
		msgs = append(msgs, NewMessage(l.file, v))
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
	return msgs
}
