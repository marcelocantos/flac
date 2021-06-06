// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.2
// source: internal/proto/refdata/refdata.proto

package refdata

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RefData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WordList *WordList `protobuf:"bytes,1,opt,name=wordList,proto3" json:"wordList,omitempty"`
	Dict     *CEDict   `protobuf:"bytes,2,opt,name=dict,proto3" json:"dict,omitempty"`
}

func (x *RefData) Reset() {
	*x = RefData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_refdata_refdata_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefData) ProtoMessage() {}

func (x *RefData) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_refdata_refdata_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefData.ProtoReflect.Descriptor instead.
func (*RefData) Descriptor() ([]byte, []int) {
	return file_internal_proto_refdata_refdata_proto_rawDescGZIP(), []int{0}
}

func (x *RefData) GetWordList() *WordList {
	if x != nil {
		return x.WordList
	}
	return nil
}

func (x *RefData) GetDict() *CEDict {
	if x != nil {
		return x.Dict
	}
	return nil
}

type WordList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Words       []string         `protobuf:"bytes,1,rep,name=words,proto3" json:"words,omitempty"`
	Frequencies map[string]int64 `protobuf:"bytes,2,rep,name=frequencies,proto3" json:"frequencies,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Positions   map[string]int64 `protobuf:"bytes,3,rep,name=positions,proto3" json:"positions,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *WordList) Reset() {
	*x = WordList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_refdata_refdata_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WordList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WordList) ProtoMessage() {}

func (x *WordList) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_refdata_refdata_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WordList.ProtoReflect.Descriptor instead.
func (*WordList) Descriptor() ([]byte, []int) {
	return file_internal_proto_refdata_refdata_proto_rawDescGZIP(), []int{1}
}

func (x *WordList) GetWords() []string {
	if x != nil {
		return x.Words
	}
	return nil
}

func (x *WordList) GetFrequencies() map[string]int64 {
	if x != nil {
		return x.Frequencies
	}
	return nil
}

func (x *WordList) GetPositions() map[string]int64 {
	if x != nil {
		return x.Positions
	}
	return nil
}

type CEDict struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entries                 map[string]*CEDict_Entries `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Syllables               map[string]bool            `protobuf:"bytes,2,rep,name=syllables,proto3" json:"syllables,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	TraditionalToSimplified map[string]string          `protobuf:"bytes,3,rep,name=traditionalToSimplified,proto3" json:"traditionalToSimplified,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *CEDict) Reset() {
	*x = CEDict{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_refdata_refdata_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CEDict) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CEDict) ProtoMessage() {}

func (x *CEDict) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_refdata_refdata_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CEDict.ProtoReflect.Descriptor instead.
func (*CEDict) Descriptor() ([]byte, []int) {
	return file_internal_proto_refdata_refdata_proto_rawDescGZIP(), []int{2}
}

func (x *CEDict) GetEntries() map[string]*CEDict_Entries {
	if x != nil {
		return x.Entries
	}
	return nil
}

func (x *CEDict) GetSyllables() map[string]bool {
	if x != nil {
		return x.Syllables
	}
	return nil
}

func (x *CEDict) GetTraditionalToSimplified() map[string]string {
	if x != nil {
		return x.TraditionalToSimplified
	}
	return nil
}

type CEDict_Entries struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Definitions map[string]*CEDict_Definitions `protobuf:"bytes,1,rep,name=definitions,proto3" json:"definitions,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Traditional string                         `protobuf:"bytes,2,opt,name=traditional,proto3" json:"traditional,omitempty"`
}

func (x *CEDict_Entries) Reset() {
	*x = CEDict_Entries{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_refdata_refdata_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CEDict_Entries) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CEDict_Entries) ProtoMessage() {}

func (x *CEDict_Entries) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_refdata_refdata_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CEDict_Entries.ProtoReflect.Descriptor instead.
func (*CEDict_Entries) Descriptor() ([]byte, []int) {
	return file_internal_proto_refdata_refdata_proto_rawDescGZIP(), []int{2, 3}
}

func (x *CEDict_Entries) GetDefinitions() map[string]*CEDict_Definitions {
	if x != nil {
		return x.Definitions
	}
	return nil
}

func (x *CEDict_Entries) GetTraditional() string {
	if x != nil {
		return x.Traditional
	}
	return ""
}

type CEDict_Definitions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Definitions []string `protobuf:"bytes,1,rep,name=definitions,proto3" json:"definitions,omitempty"`
}

func (x *CEDict_Definitions) Reset() {
	*x = CEDict_Definitions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_refdata_refdata_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CEDict_Definitions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CEDict_Definitions) ProtoMessage() {}

func (x *CEDict_Definitions) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_refdata_refdata_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CEDict_Definitions.ProtoReflect.Descriptor instead.
func (*CEDict_Definitions) Descriptor() ([]byte, []int) {
	return file_internal_proto_refdata_refdata_proto_rawDescGZIP(), []int{2, 4}
}

func (x *CEDict_Definitions) GetDefinitions() []string {
	if x != nil {
		return x.Definitions
	}
	return nil
}

var File_internal_proto_refdata_refdata_proto protoreflect.FileDescriptor

var file_internal_proto_refdata_refdata_proto_rawDesc = []byte{
	0x0a, 0x24, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x72, 0x65, 0x66, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x72, 0x65, 0x66, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4d, 0x0a, 0x07, 0x52, 0x65, 0x66, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x25, 0x0a, 0x08, 0x77, 0x6f, 0x72, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x57, 0x6f, 0x72, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x08,
	0x77, 0x6f, 0x72, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x04, 0x64, 0x69, 0x63, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x43, 0x45, 0x44, 0x69, 0x63, 0x74, 0x52,
	0x04, 0x64, 0x69, 0x63, 0x74, 0x22, 0x94, 0x02, 0x0a, 0x08, 0x57, 0x6f, 0x72, 0x64, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x05, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x3c, 0x0a, 0x0b, 0x66, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x57, 0x6f, 0x72, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x2e, 0x46, 0x72, 0x65, 0x71, 0x75, 0x65, 0x6e,
	0x63, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x66, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x12, 0x36, 0x0a, 0x09, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x57, 0x6f, 0x72, 0x64,
	0x4c, 0x69, 0x73, 0x74, 0x2e, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x09, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x3e,
	0x0a, 0x10, 0x46, 0x72, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3c,
	0x0a, 0x0e, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x9d, 0x05, 0x0a,
	0x06, 0x43, 0x45, 0x44, 0x69, 0x63, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x43, 0x45, 0x44, 0x69, 0x63,
	0x74, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07,
	0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x34, 0x0a, 0x09, 0x73, 0x79, 0x6c, 0x6c, 0x61,
	0x62, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x43, 0x45, 0x44,
	0x69, 0x63, 0x74, 0x2e, 0x53, 0x79, 0x6c, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x09, 0x73, 0x79, 0x6c, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x12, 0x5e, 0x0a,
	0x17, 0x74, 0x72, 0x61, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x54, 0x6f, 0x53, 0x69,
	0x6d, 0x70, 0x6c, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24,
	0x2e, 0x43, 0x45, 0x44, 0x69, 0x63, 0x74, 0x2e, 0x54, 0x72, 0x61, 0x64, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x61, 0x6c, 0x54, 0x6f, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x69, 0x66, 0x69, 0x65, 0x64, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x17, 0x74, 0x72, 0x61, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61,
	0x6c, 0x54, 0x6f, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x69, 0x66, 0x69, 0x65, 0x64, 0x1a, 0x4b, 0x0a,
	0x0c, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x25, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x43, 0x45, 0x44, 0x69, 0x63, 0x74, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3c, 0x0a, 0x0e, 0x53, 0x79,
	0x6c, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x4a, 0x0a, 0x1c, 0x54, 0x72, 0x61, 0x64,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x54, 0x6f, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x69, 0x66,
	0x69, 0x65, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x1a, 0xc4, 0x01, 0x0a, 0x07, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x12, 0x42, 0x0a, 0x0b, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x43, 0x45, 0x44, 0x69, 0x63, 0x74, 0x2e, 0x45,
	0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x74, 0x72, 0x61, 0x64, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x74, 0x72, 0x61, 0x64, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x1a, 0x53, 0x0a, 0x10, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x29, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x43, 0x45,
	0x44, 0x69, 0x63, 0x74, 0x2e, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x2f, 0x0a, 0x0b, 0x44,
	0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65,
	0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x0b, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x18, 0x5a, 0x16,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72,
	0x65, 0x66, 0x64, 0x61, 0x74, 0x61, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_proto_refdata_refdata_proto_rawDescOnce sync.Once
	file_internal_proto_refdata_refdata_proto_rawDescData = file_internal_proto_refdata_refdata_proto_rawDesc
)

func file_internal_proto_refdata_refdata_proto_rawDescGZIP() []byte {
	file_internal_proto_refdata_refdata_proto_rawDescOnce.Do(func() {
		file_internal_proto_refdata_refdata_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_proto_refdata_refdata_proto_rawDescData)
	})
	return file_internal_proto_refdata_refdata_proto_rawDescData
}

var file_internal_proto_refdata_refdata_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_internal_proto_refdata_refdata_proto_goTypes = []interface{}{
	(*RefData)(nil),            // 0: RefData
	(*WordList)(nil),           // 1: WordList
	(*CEDict)(nil),             // 2: CEDict
	nil,                        // 3: WordList.FrequenciesEntry
	nil,                        // 4: WordList.PositionsEntry
	nil,                        // 5: CEDict.EntriesEntry
	nil,                        // 6: CEDict.SyllablesEntry
	nil,                        // 7: CEDict.TraditionalToSimplifiedEntry
	(*CEDict_Entries)(nil),     // 8: CEDict.Entries
	(*CEDict_Definitions)(nil), // 9: CEDict.Definitions
	nil,                        // 10: CEDict.Entries.DefinitionsEntry
}
var file_internal_proto_refdata_refdata_proto_depIdxs = []int32{
	1,  // 0: RefData.wordList:type_name -> WordList
	2,  // 1: RefData.dict:type_name -> CEDict
	3,  // 2: WordList.frequencies:type_name -> WordList.FrequenciesEntry
	4,  // 3: WordList.positions:type_name -> WordList.PositionsEntry
	5,  // 4: CEDict.entries:type_name -> CEDict.EntriesEntry
	6,  // 5: CEDict.syllables:type_name -> CEDict.SyllablesEntry
	7,  // 6: CEDict.traditionalToSimplified:type_name -> CEDict.TraditionalToSimplifiedEntry
	8,  // 7: CEDict.EntriesEntry.value:type_name -> CEDict.Entries
	10, // 8: CEDict.Entries.definitions:type_name -> CEDict.Entries.DefinitionsEntry
	9,  // 9: CEDict.Entries.DefinitionsEntry.value:type_name -> CEDict.Definitions
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_internal_proto_refdata_refdata_proto_init() }
func file_internal_proto_refdata_refdata_proto_init() {
	if File_internal_proto_refdata_refdata_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_proto_refdata_refdata_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RefData); i {
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
		file_internal_proto_refdata_refdata_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WordList); i {
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
		file_internal_proto_refdata_refdata_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CEDict); i {
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
		file_internal_proto_refdata_refdata_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CEDict_Entries); i {
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
		file_internal_proto_refdata_refdata_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CEDict_Definitions); i {
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
			RawDescriptor: file_internal_proto_refdata_refdata_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_proto_refdata_refdata_proto_goTypes,
		DependencyIndexes: file_internal_proto_refdata_refdata_proto_depIdxs,
		MessageInfos:      file_internal_proto_refdata_refdata_proto_msgTypes,
	}.Build()
	File_internal_proto_refdata_refdata_proto = out.File
	file_internal_proto_refdata_refdata_proto_rawDesc = nil
	file_internal_proto_refdata_refdata_proto_goTypes = nil
	file_internal_proto_refdata_refdata_proto_depIdxs = nil
}
