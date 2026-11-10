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
	// StructType is the value of +structType marker (e.g. atomic)
	StructType string
}

type ProtobufMessageOption struct {
	ClusterScope bool
}

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
	SubResource     bool
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
	// ListType is the value of +listType marker (e.g. atomic, map, set)
	ListType    string
	ListMapKeys []string
	// MapType is the value of +mapType marker (e.g. atomic, granular)
	MapType string
}
