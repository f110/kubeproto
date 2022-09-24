package goparser

type ProtobufFile struct {
	Package   string
	GoPackage string
	Domain    string
	SubGroup  string
	Version   string

	Messages []*ProtobufMessage
}

type ProtobufMessage struct {
	Name        string
	Fields      []*ProtobufField
	UseFieldsOf string
	Option      *ProtobufMessageOption
}

type ProtobufMessageOption struct{}

func (m *ProtobufMessage) IsRuntimeObject() bool {
	var foundTypeMeta, foundObjectMeta bool
	for _, f := range m.Fields {
		if f.Name == "type_meta" {
			foundTypeMeta = true
			continue
		}
		if f.Name == "object_meta" {
			foundObjectMeta = true
			continue
		}
	}

	return foundTypeMeta && foundObjectMeta
}

type ProtobufField struct {
	Name            string
	GoName          string
	APIFieldName    string
	Kind            string
	IsMap           bool
	InvalidProtobuf bool
	MapKeyKind      string
	MapValueKind    string
	Index           int
	Repeated        bool
	Optional        bool
	Inline          bool
	ExternalPackage string
	Doc             string
}
