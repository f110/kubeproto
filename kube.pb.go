// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.20.0--rc2
// source: kube.proto

package kubeproto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Kind struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AdditionalPrinterColumns []*PrinterColumn `protobuf:"bytes,1,rep,name=additional_printer_columns,json=additionalPrinterColumns,proto3" json:"additional_printer_columns,omitempty"`
}

func (x *Kind) Reset() {
	*x = Kind{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kube_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Kind) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Kind) ProtoMessage() {}

func (x *Kind) ProtoReflect() protoreflect.Message {
	mi := &file_kube_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Kind.ProtoReflect.Descriptor instead.
func (*Kind) Descriptor() ([]byte, []int) {
	return file_kube_proto_rawDescGZIP(), []int{0}
}

func (x *Kind) GetAdditionalPrinterColumns() []*PrinterColumn {
	if x != nil {
		return x.AdditionalPrinterColumns
	}
	return nil
}

type Field struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GoName       string `protobuf:"bytes,1,opt,name=go_name,json=goName,proto3" json:"go_name,omitempty"`
	Inline       bool   `protobuf:"varint,2,opt,name=inline,proto3" json:"inline,omitempty"`
	SubResource  bool   `protobuf:"varint,3,opt,name=sub_resource,json=subResource,proto3" json:"sub_resource,omitempty"`
	ApiFieldName string `protobuf:"bytes,4,opt,name=api_field_name,json=apiFieldName,proto3" json:"api_field_name,omitempty"`
}

func (x *Field) Reset() {
	*x = Field{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kube_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Field) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Field) ProtoMessage() {}

func (x *Field) ProtoReflect() protoreflect.Message {
	mi := &file_kube_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Field.ProtoReflect.Descriptor instead.
func (*Field) Descriptor() ([]byte, []int) {
	return file_kube_proto_rawDescGZIP(), []int{1}
}

func (x *Field) GetGoName() string {
	if x != nil {
		return x.GoName
	}
	return ""
}

func (x *Field) GetInline() bool {
	if x != nil {
		return x.Inline
	}
	return false
}

func (x *Field) GetSubResource() bool {
	if x != nil {
		return x.SubResource
	}
	return false
}

func (x *Field) GetApiFieldName() string {
	if x != nil {
		return x.ApiFieldName
	}
	return ""
}

type Kubernetes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Domain   string `protobuf:"bytes,1,opt,name=domain,proto3" json:"domain,omitempty"`
	SubGroup string `protobuf:"bytes,2,opt,name=sub_group,json=subGroup,proto3" json:"sub_group,omitempty"`
	Version  string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	Served   bool   `protobuf:"varint,4,opt,name=served,proto3" json:"served,omitempty"`
	Storage  bool   `protobuf:"varint,5,opt,name=storage,proto3" json:"storage,omitempty"`
}

func (x *Kubernetes) Reset() {
	*x = Kubernetes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kube_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Kubernetes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Kubernetes) ProtoMessage() {}

func (x *Kubernetes) ProtoReflect() protoreflect.Message {
	mi := &file_kube_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Kubernetes.ProtoReflect.Descriptor instead.
func (*Kubernetes) Descriptor() ([]byte, []int) {
	return file_kube_proto_rawDescGZIP(), []int{2}
}

func (x *Kubernetes) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *Kubernetes) GetSubGroup() string {
	if x != nil {
		return x.SubGroup
	}
	return ""
}

func (x *Kubernetes) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Kubernetes) GetServed() bool {
	if x != nil {
		return x.Served
	}
	return false
}

func (x *Kubernetes) GetStorage() bool {
	if x != nil {
		return x.Storage
	}
	return false
}

type PrinterColumn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Description string `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	JsonPath    string `protobuf:"bytes,3,opt,name=json_path,json=jsonPath,proto3" json:"json_path,omitempty"`
	Priority    int32  `protobuf:"varint,4,opt,name=priority,proto3" json:"priority,omitempty"`
	Type        string `protobuf:"bytes,5,opt,name=type,proto3" json:"type,omitempty"`
	Format      string `protobuf:"bytes,6,opt,name=format,proto3" json:"format,omitempty"`
}

func (x *PrinterColumn) Reset() {
	*x = PrinterColumn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kube_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrinterColumn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrinterColumn) ProtoMessage() {}

func (x *PrinterColumn) ProtoReflect() protoreflect.Message {
	mi := &file_kube_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrinterColumn.ProtoReflect.Descriptor instead.
func (*PrinterColumn) Descriptor() ([]byte, []int) {
	return file_kube_proto_rawDescGZIP(), []int{3}
}

func (x *PrinterColumn) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *PrinterColumn) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PrinterColumn) GetJsonPath() string {
	if x != nil {
		return x.JsonPath
	}
	return ""
}

func (x *PrinterColumn) GetPriority() int32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

func (x *PrinterColumn) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *PrinterColumn) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

type EnumValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *EnumValue) Reset() {
	*x = EnumValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kube_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnumValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnumValue) ProtoMessage() {}

func (x *EnumValue) ProtoReflect() protoreflect.Message {
	mi := &file_kube_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnumValue.ProtoReflect.Descriptor instead.
func (*EnumValue) Descriptor() ([]byte, []int) {
	return file_kube_proto_rawDescGZIP(), []int{4}
}

func (x *EnumValue) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var file_kube_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*Kind)(nil),
		Field:         60010,
		Name:          "dev.f110.kubeproto.kind",
		Tag:           "bytes,60010,opt,name=kind",
		Filename:      "kube.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*Field)(nil),
		Field:         60010,
		Name:          "dev.f110.kubeproto.field",
		Tag:           "bytes,60010,opt,name=field",
		Filename:      "kube.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*Kubernetes)(nil),
		Field:         60010,
		Name:          "dev.f110.kubeproto.k8s",
		Tag:           "bytes,60010,opt,name=k8s",
		Filename:      "kube.proto",
	},
	{
		ExtendedType:  (*descriptorpb.EnumValueOptions)(nil),
		ExtensionType: (*EnumValue)(nil),
		Field:         60010,
		Name:          "dev.f110.kubeproto.value",
		Tag:           "bytes,60010,opt,name=value",
		Filename:      "kube.proto",
	},
}

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional dev.f110.kubeproto.Kind kind = 60010;
	E_Kind = &file_kube_proto_extTypes[0]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional dev.f110.kubeproto.Field field = 60010;
	E_Field = &file_kube_proto_extTypes[1]
)

// Extension fields to descriptorpb.FileOptions.
var (
	// optional dev.f110.kubeproto.Kubernetes k8s = 60010;
	E_K8S = &file_kube_proto_extTypes[2]
)

// Extension fields to descriptorpb.EnumValueOptions.
var (
	// optional dev.f110.kubeproto.EnumValue value = 60010;
	E_Value = &file_kube_proto_extTypes[3]
)

var File_kube_proto protoreflect.FileDescriptor

var file_kube_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6b, 0x75, 0x62, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x64, 0x65,
	0x76, 0x2e, 0x66, 0x31, 0x31, 0x30, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x67, 0x0a, 0x04, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x5f, 0x0a, 0x1a, 0x61, 0x64,
	0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x5f, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x5f, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x64, 0x65, 0x76, 0x2e, 0x66, 0x31, 0x31, 0x30, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6c, 0x75, 0x6d,
	0x6e, 0x52, 0x18, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x50, 0x72, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x22, 0x81, 0x01, 0x0a, 0x05,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x67, 0x6f, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x67, 0x6f, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x69, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
	0x69, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x75, 0x62, 0x5f, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x73, 0x75,
	0x62, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x24, 0x0a, 0x0e, 0x61, 0x70, 0x69,
	0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x61, 0x70, 0x69, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x22,
	0x8d, 0x01, 0x0a, 0x0a, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x12, 0x16,
	0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x75, 0x62, 0x5f, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x75, 0x62, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x22,
	0xaa, 0x01, 0x0a, 0x0d, 0x50, 0x72, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6c, 0x75, 0x6d,
	0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6a, 0x73, 0x6f, 0x6e, 0x5f,
	0x70, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6a, 0x73, 0x6f, 0x6e,
	0x50, 0x61, 0x74, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x22, 0x21, 0x0a, 0x09,
	0x45, 0x6e, 0x75, 0x6d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x4f, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xea, 0xd4, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x18, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x66, 0x31, 0x31, 0x30, 0x2e, 0x6b, 0x75, 0x62, 0x65,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64,
	0x3a, 0x50, 0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c,
	0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xea, 0xd4, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x66, 0x31, 0x31, 0x30, 0x2e, 0x6b, 0x75, 0x62, 0x65,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x52, 0x05, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x3a, 0x50, 0x0a, 0x03, 0x6b, 0x38, 0x73, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xea, 0xd4, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1e, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x66, 0x31, 0x31, 0x30, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x52,
	0x03, 0x6b, 0x38, 0x73, 0x3a, 0x58, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x21, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6e, 0x75, 0x6d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0xea, 0xd4, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x66,
	0x31, 0x31, 0x30, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6e,
	0x75, 0x6d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kube_proto_rawDescOnce sync.Once
	file_kube_proto_rawDescData = file_kube_proto_rawDesc
)

func file_kube_proto_rawDescGZIP() []byte {
	file_kube_proto_rawDescOnce.Do(func() {
		file_kube_proto_rawDescData = protoimpl.X.CompressGZIP(file_kube_proto_rawDescData)
	})
	return file_kube_proto_rawDescData
}

var file_kube_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_kube_proto_goTypes = []interface{}{
	(*Kind)(nil),                          // 0: dev.f110.kubeproto.Kind
	(*Field)(nil),                         // 1: dev.f110.kubeproto.Field
	(*Kubernetes)(nil),                    // 2: dev.f110.kubeproto.Kubernetes
	(*PrinterColumn)(nil),                 // 3: dev.f110.kubeproto.PrinterColumn
	(*EnumValue)(nil),                     // 4: dev.f110.kubeproto.EnumValue
	(*descriptorpb.MessageOptions)(nil),   // 5: google.protobuf.MessageOptions
	(*descriptorpb.FieldOptions)(nil),     // 6: google.protobuf.FieldOptions
	(*descriptorpb.FileOptions)(nil),      // 7: google.protobuf.FileOptions
	(*descriptorpb.EnumValueOptions)(nil), // 8: google.protobuf.EnumValueOptions
}
var file_kube_proto_depIdxs = []int32{
	3, // 0: dev.f110.kubeproto.Kind.additional_printer_columns:type_name -> dev.f110.kubeproto.PrinterColumn
	5, // 1: dev.f110.kubeproto.kind:extendee -> google.protobuf.MessageOptions
	6, // 2: dev.f110.kubeproto.field:extendee -> google.protobuf.FieldOptions
	7, // 3: dev.f110.kubeproto.k8s:extendee -> google.protobuf.FileOptions
	8, // 4: dev.f110.kubeproto.value:extendee -> google.protobuf.EnumValueOptions
	0, // 5: dev.f110.kubeproto.kind:type_name -> dev.f110.kubeproto.Kind
	1, // 6: dev.f110.kubeproto.field:type_name -> dev.f110.kubeproto.Field
	2, // 7: dev.f110.kubeproto.k8s:type_name -> dev.f110.kubeproto.Kubernetes
	4, // 8: dev.f110.kubeproto.value:type_name -> dev.f110.kubeproto.EnumValue
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	5, // [5:9] is the sub-list for extension type_name
	1, // [1:5] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_kube_proto_init() }
func file_kube_proto_init() {
	if File_kube_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_kube_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Kind); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_kube_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Field); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_kube_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Kubernetes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_kube_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrinterColumn); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_kube_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnumValue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_kube_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 4,
			NumServices:   0,
		},
		GoTypes:           file_kube_proto_goTypes,
		DependencyIndexes: file_kube_proto_depIdxs,
		MessageInfos:      file_kube_proto_msgTypes,
		ExtensionInfos:    file_kube_proto_extTypes,
	}.Build()
	File_kube_proto = out.File
	file_kube_proto_rawDesc = nil
	file_kube_proto_goTypes = nil
	file_kube_proto_depIdxs = nil
}
