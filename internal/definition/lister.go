package definition

import (
	"google.golang.org/protobuf/types/descriptorpb"
)

type Lister struct {
	files    []*descriptorpb.FileDescriptorProto
	allFiles []*descriptorpb.FileDescriptorProto
}

func NewLister(files []*descriptorpb.FileDescriptorProto, allProtos []*descriptorpb.FileDescriptorProto) *Lister {
	return &Lister{files: files, allFiles: allProtos}
}

func (l *Lister) GetMessages() Messages {
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
	return msgs
}
