// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: verifiablecredstypes/bertyverifiablecreds.proto

package verifiablecredstypes

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

type FlowType int32

const (
	FlowType_FlowTypeUndefined FlowType = 0
	// FlowTypeCode asks users a code sent on a side channel
	FlowType_FlowTypeCode FlowType = 1
	// FlowTypeAuth currently unimplemented
	FlowType_FlowTypeAuth FlowType = 2
	// FlowTypeProof currently unimplemented
	FlowType_FlowTypeProof FlowType = 3
)

// Enum value maps for FlowType.
var (
	FlowType_name = map[int32]string{
		0: "FlowTypeUndefined",
		1: "FlowTypeCode",
		2: "FlowTypeAuth",
		3: "FlowTypeProof",
	}
	FlowType_value = map[string]int32{
		"FlowTypeUndefined": 0,
		"FlowTypeCode":      1,
		"FlowTypeAuth":      2,
		"FlowTypeProof":     3,
	}
)

func (x FlowType) Enum() *FlowType {
	p := new(FlowType)
	*p = x
	return p
}

func (x FlowType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FlowType) Descriptor() protoreflect.EnumDescriptor {
	return file_verifiablecredstypes_bertyverifiablecreds_proto_enumTypes[0].Descriptor()
}

func (FlowType) Type() protoreflect.EnumType {
	return &file_verifiablecredstypes_bertyverifiablecreds_proto_enumTypes[0]
}

func (x FlowType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FlowType.Descriptor instead.
func (FlowType) EnumDescriptor() ([]byte, []int) {
	return file_verifiablecredstypes_bertyverifiablecreds_proto_rawDescGZIP(), []int{0}
}

type CodeStrategy int32

const (
	CodeStrategy_CodeStrategyUndefined CodeStrategy = 0
	// CodeStrategy6Digits currently unimplemented
	CodeStrategy_CodeStrategy6Digits CodeStrategy = 1
	// CodeStrategy10Chars currently unimplemented
	CodeStrategy_CodeStrategy10Chars CodeStrategy = 2
	// CodeStrategyMocked6Zeroes must only be used in testing
	CodeStrategy_CodeStrategyMocked6Zeroes CodeStrategy = 999
)

// Enum value maps for CodeStrategy.
var (
	CodeStrategy_name = map[int32]string{
		0:   "CodeStrategyUndefined",
		1:   "CodeStrategy6Digits",
		2:   "CodeStrategy10Chars",
		999: "CodeStrategyMocked6Zeroes",
	}
	CodeStrategy_value = map[string]int32{
		"CodeStrategyUndefined":     0,
		"CodeStrategy6Digits":       1,
		"CodeStrategy10Chars":       2,
		"CodeStrategyMocked6Zeroes": 999,
	}
)

func (x CodeStrategy) Enum() *CodeStrategy {
	p := new(CodeStrategy)
	*p = x
	return p
}

func (x CodeStrategy) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CodeStrategy) Descriptor() protoreflect.EnumDescriptor {
	return file_verifiablecredstypes_bertyverifiablecreds_proto_enumTypes[1].Descriptor()
}

func (CodeStrategy) Type() protoreflect.EnumType {
	return &file_verifiablecredstypes_bertyverifiablecreds_proto_enumTypes[1]
}

func (x CodeStrategy) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CodeStrategy.Descriptor instead.
func (CodeStrategy) EnumDescriptor() ([]byte, []int) {
	return file_verifiablecredstypes_bertyverifiablecreds_proto_rawDescGZIP(), []int{1}
}

// StateChallenge serialized and signed state used when requesting a challenge
type StateChallenge struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp   []byte `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Nonce       []byte `protobuf:"bytes,2,opt,name=nonce,proto3" json:"nonce,omitempty"`
	BertyLink   string `protobuf:"bytes,3,opt,name=berty_link,json=bertyLink,proto3" json:"berty_link,omitempty"`
	RedirectUri string `protobuf:"bytes,4,opt,name=redirect_uri,json=redirectUri,proto3" json:"redirect_uri,omitempty"`
	State       string `protobuf:"bytes,5,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *StateChallenge) Reset() {
	*x = StateChallenge{}
	if protoimpl.UnsafeEnabled {
		mi := &file_verifiablecredstypes_bertyverifiablecreds_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StateChallenge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateChallenge) ProtoMessage() {}

func (x *StateChallenge) ProtoReflect() protoreflect.Message {
	mi := &file_verifiablecredstypes_bertyverifiablecreds_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateChallenge.ProtoReflect.Descriptor instead.
func (*StateChallenge) Descriptor() ([]byte, []int) {
	return file_verifiablecredstypes_bertyverifiablecreds_proto_rawDescGZIP(), []int{0}
}

func (x *StateChallenge) GetTimestamp() []byte {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *StateChallenge) GetNonce() []byte {
	if x != nil {
		return x.Nonce
	}
	return nil
}

func (x *StateChallenge) GetBertyLink() string {
	if x != nil {
		return x.BertyLink
	}
	return ""
}

func (x *StateChallenge) GetRedirectUri() string {
	if x != nil {
		return x.RedirectUri
	}
	return ""
}

func (x *StateChallenge) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

// StateCode serialized and signed state used when requesting a code
type StateCode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp    []byte       `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	BertyLink    string       `protobuf:"bytes,2,opt,name=berty_link,json=bertyLink,proto3" json:"berty_link,omitempty"`
	CodeStrategy CodeStrategy `protobuf:"varint,3,opt,name=code_strategy,json=codeStrategy,proto3,enum=weshnet.account.v1.CodeStrategy" json:"code_strategy,omitempty"`
	Identifier   string       `protobuf:"bytes,4,opt,name=identifier,proto3" json:"identifier,omitempty"`
	Code         string       `protobuf:"bytes,5,opt,name=code,proto3" json:"code,omitempty"`
	RedirectUri  string       `protobuf:"bytes,6,opt,name=redirect_uri,json=redirectUri,proto3" json:"redirect_uri,omitempty"`
	State        string       `protobuf:"bytes,7,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *StateCode) Reset() {
	*x = StateCode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_verifiablecredstypes_bertyverifiablecreds_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StateCode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateCode) ProtoMessage() {}

func (x *StateCode) ProtoReflect() protoreflect.Message {
	mi := &file_verifiablecredstypes_bertyverifiablecreds_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateCode.ProtoReflect.Descriptor instead.
func (*StateCode) Descriptor() ([]byte, []int) {
	return file_verifiablecredstypes_bertyverifiablecreds_proto_rawDescGZIP(), []int{1}
}

func (x *StateCode) GetTimestamp() []byte {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *StateCode) GetBertyLink() string {
	if x != nil {
		return x.BertyLink
	}
	return ""
}

func (x *StateCode) GetCodeStrategy() CodeStrategy {
	if x != nil {
		return x.CodeStrategy
	}
	return CodeStrategy_CodeStrategyUndefined
}

func (x *StateCode) GetIdentifier() string {
	if x != nil {
		return x.Identifier
	}
	return ""
}

func (x *StateCode) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *StateCode) GetRedirectUri() string {
	if x != nil {
		return x.RedirectUri
	}
	return ""
}

func (x *StateCode) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

type AccountCryptoChallenge struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Challenge string `protobuf:"bytes,1,opt,name=challenge,proto3" json:"challenge,omitempty"`
}

func (x *AccountCryptoChallenge) Reset() {
	*x = AccountCryptoChallenge{}
	if protoimpl.UnsafeEnabled {
		mi := &file_verifiablecredstypes_bertyverifiablecreds_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccountCryptoChallenge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccountCryptoChallenge) ProtoMessage() {}

func (x *AccountCryptoChallenge) ProtoReflect() protoreflect.Message {
	mi := &file_verifiablecredstypes_bertyverifiablecreds_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccountCryptoChallenge.ProtoReflect.Descriptor instead.
func (*AccountCryptoChallenge) Descriptor() ([]byte, []int) {
	return file_verifiablecredstypes_bertyverifiablecreds_proto_rawDescGZIP(), []int{2}
}

func (x *AccountCryptoChallenge) GetChallenge() string {
	if x != nil {
		return x.Challenge
	}
	return ""
}

var File_verifiablecredstypes_bertyverifiablecreds_proto protoreflect.FileDescriptor

var file_verifiablecredstypes_bertyverifiablecreds_proto_rawDesc = []byte{
	0x0a, 0x2f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x63, 0x72, 0x65, 0x64,
	0x73, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x62, 0x65, 0x72, 0x74, 0x79, 0x76, 0x65, 0x72, 0x69,
	0x66, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x63, 0x72, 0x65, 0x64, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x12, 0x77, 0x65, 0x73, 0x68, 0x6e, 0x65, 0x74, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x2e, 0x76, 0x31, 0x22, 0x9c, 0x01, 0x0a, 0x0e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x43,
	0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x62, 0x65, 0x72, 0x74, 0x79, 0x5f, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x62, 0x65, 0x72, 0x74, 0x79, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x21, 0x0a, 0x0c, 0x72,
	0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x55, 0x72, 0x69, 0x12, 0x14,
	0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x22, 0xfc, 0x01, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x74, 0x65, 0x43, 0x6f,
	0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x65, 0x72, 0x74, 0x79, 0x5f, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x62, 0x65, 0x72, 0x74, 0x79, 0x4c, 0x69, 0x6e, 0x6b, 0x12,
	0x45, 0x0a, 0x0d, 0x63, 0x6f, 0x64, 0x65, 0x5f, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x77, 0x65, 0x73, 0x68, 0x6e, 0x65, 0x74,
	0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x64, 0x65,
	0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x52, 0x0c, 0x63, 0x6f, 0x64, 0x65, 0x53, 0x74,
	0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65,
	0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x55, 0x72, 0x69, 0x12, 0x14, 0x0a,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x22, 0x36, 0x0a, 0x16, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x43, 0x72,
	0x79, 0x70, 0x74, 0x6f, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x2a, 0x58, 0x0a, 0x08, 0x46,
	0x6c, 0x6f, 0x77, 0x54, 0x79, 0x70, 0x65, 0x12, 0x15, 0x0a, 0x11, 0x46, 0x6c, 0x6f, 0x77, 0x54,
	0x79, 0x70, 0x65, 0x55, 0x6e, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x10, 0x00, 0x12, 0x10,
	0x0a, 0x0c, 0x46, 0x6c, 0x6f, 0x77, 0x54, 0x79, 0x70, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x10, 0x01,
	0x12, 0x10, 0x0a, 0x0c, 0x46, 0x6c, 0x6f, 0x77, 0x54, 0x79, 0x70, 0x65, 0x41, 0x75, 0x74, 0x68,
	0x10, 0x02, 0x12, 0x11, 0x0a, 0x0d, 0x46, 0x6c, 0x6f, 0x77, 0x54, 0x79, 0x70, 0x65, 0x50, 0x72,
	0x6f, 0x6f, 0x66, 0x10, 0x03, 0x2a, 0x7b, 0x0a, 0x0c, 0x43, 0x6f, 0x64, 0x65, 0x53, 0x74, 0x72,
	0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x19, 0x0a, 0x15, 0x43, 0x6f, 0x64, 0x65, 0x53, 0x74, 0x72,
	0x61, 0x74, 0x65, 0x67, 0x79, 0x55, 0x6e, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x10, 0x00,
	0x12, 0x17, 0x0a, 0x13, 0x43, 0x6f, 0x64, 0x65, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79,
	0x36, 0x44, 0x69, 0x67, 0x69, 0x74, 0x73, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x6f, 0x64,
	0x65, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x31, 0x30, 0x43, 0x68, 0x61, 0x72, 0x73,
	0x10, 0x02, 0x12, 0x1e, 0x0a, 0x19, 0x43, 0x6f, 0x64, 0x65, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65,
	0x67, 0x79, 0x4d, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x36, 0x5a, 0x65, 0x72, 0x6f, 0x65, 0x73, 0x10,
	0xe7, 0x07, 0x42, 0x2d, 0x5a, 0x2b, 0x62, 0x65, 0x72, 0x74, 0x79, 0x2e, 0x74, 0x65, 0x63, 0x68,
	0x2f, 0x77, 0x65, 0x73, 0x68, 0x6e, 0x65, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x76, 0x65, 0x72,
	0x69, 0x66, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x63, 0x72, 0x65, 0x64, 0x73, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_verifiablecredstypes_bertyverifiablecreds_proto_rawDescOnce sync.Once
	file_verifiablecredstypes_bertyverifiablecreds_proto_rawDescData = file_verifiablecredstypes_bertyverifiablecreds_proto_rawDesc
)

func file_verifiablecredstypes_bertyverifiablecreds_proto_rawDescGZIP() []byte {
	file_verifiablecredstypes_bertyverifiablecreds_proto_rawDescOnce.Do(func() {
		file_verifiablecredstypes_bertyverifiablecreds_proto_rawDescData = protoimpl.X.CompressGZIP(file_verifiablecredstypes_bertyverifiablecreds_proto_rawDescData)
	})
	return file_verifiablecredstypes_bertyverifiablecreds_proto_rawDescData
}

var file_verifiablecredstypes_bertyverifiablecreds_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_verifiablecredstypes_bertyverifiablecreds_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_verifiablecredstypes_bertyverifiablecreds_proto_goTypes = []interface{}{
	(FlowType)(0),                  // 0: weshnet.account.v1.FlowType
	(CodeStrategy)(0),              // 1: weshnet.account.v1.CodeStrategy
	(*StateChallenge)(nil),         // 2: weshnet.account.v1.StateChallenge
	(*StateCode)(nil),              // 3: weshnet.account.v1.StateCode
	(*AccountCryptoChallenge)(nil), // 4: weshnet.account.v1.AccountCryptoChallenge
}
var file_verifiablecredstypes_bertyverifiablecreds_proto_depIdxs = []int32{
	1, // 0: weshnet.account.v1.StateCode.code_strategy:type_name -> weshnet.account.v1.CodeStrategy
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_verifiablecredstypes_bertyverifiablecreds_proto_init() }
func file_verifiablecredstypes_bertyverifiablecreds_proto_init() {
	if File_verifiablecredstypes_bertyverifiablecreds_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_verifiablecredstypes_bertyverifiablecreds_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StateChallenge); i {
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
		file_verifiablecredstypes_bertyverifiablecreds_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StateCode); i {
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
		file_verifiablecredstypes_bertyverifiablecreds_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccountCryptoChallenge); i {
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
			RawDescriptor: file_verifiablecredstypes_bertyverifiablecreds_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_verifiablecredstypes_bertyverifiablecreds_proto_goTypes,
		DependencyIndexes: file_verifiablecredstypes_bertyverifiablecreds_proto_depIdxs,
		EnumInfos:         file_verifiablecredstypes_bertyverifiablecreds_proto_enumTypes,
		MessageInfos:      file_verifiablecredstypes_bertyverifiablecreds_proto_msgTypes,
	}.Build()
	File_verifiablecredstypes_bertyverifiablecreds_proto = out.File
	file_verifiablecredstypes_bertyverifiablecreds_proto_rawDesc = nil
	file_verifiablecredstypes_bertyverifiablecreds_proto_goTypes = nil
	file_verifiablecredstypes_bertyverifiablecreds_proto_depIdxs = nil
}