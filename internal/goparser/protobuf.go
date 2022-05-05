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
}

type ProtobufField struct {
	Name            string
	GoName          string
	APIFieldName    string
	Kind            string
	Index           int
	Repeated        bool
	Optional        bool
	Inline          bool
	ExternalPackage string
}
