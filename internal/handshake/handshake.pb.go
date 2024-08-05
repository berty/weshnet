// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: handshake/handshake.proto

package handshake

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

type BoxEnvelope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Box []byte `protobuf:"bytes,1,opt,name=box,proto3" json:"box,omitempty"`
}

func (x *BoxEnvelope) Reset() {
	*x = BoxEnvelope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_handshake_handshake_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoxEnvelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoxEnvelope) ProtoMessage() {}

func (x *BoxEnvelope) ProtoReflect() protoreflect.Message {
	mi := &file_handshake_handshake_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoxEnvelope.ProtoReflect.Descriptor instead.
func (*BoxEnvelope) Descriptor() ([]byte, []int) {
	return file_handshake_handshake_proto_rawDescGZIP(), []int{0}
}

func (x *BoxEnvelope) GetBox() []byte {
	if x != nil {
		return x.Box
	}
	return nil
}

type HelloPayload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EphemeralPubKey []byte `protobuf:"bytes,1,opt,name=ephemeral_pub_key,json=ephemeralPubKey,proto3" json:"ephemeral_pub_key,omitempty"`
}

func (x *HelloPayload) Reset() {
	*x = HelloPayload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_handshake_handshake_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloPayload) ProtoMessage() {}

func (x *HelloPayload) ProtoReflect() protoreflect.Message {
	mi := &file_handshake_handshake_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloPayload.ProtoReflect.Descriptor instead.
func (*HelloPayload) Descriptor() ([]byte, []int) {
	return file_handshake_handshake_proto_rawDescGZIP(), []int{1}
}

func (x *HelloPayload) GetEphemeralPubKey() []byte {
	if x != nil {
		return x.EphemeralPubKey
	}
	return nil
}

type RequesterAuthenticatePayload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequesterAccountId  []byte `protobuf:"bytes,1,opt,name=requester_account_id,json=requesterAccountId,proto3" json:"requester_account_id,omitempty"`
	RequesterAccountSig []byte `protobuf:"bytes,2,opt,name=requester_account_sig,json=requesterAccountSig,proto3" json:"requester_account_sig,omitempty"`
}

func (x *RequesterAuthenticatePayload) Reset() {
	*x = RequesterAuthenticatePayload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_handshake_handshake_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequesterAuthenticatePayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequesterAuthenticatePayload) ProtoMessage() {}

func (x *RequesterAuthenticatePayload) ProtoReflect() protoreflect.Message {
	mi := &file_handshake_handshake_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequesterAuthenticatePayload.ProtoReflect.Descriptor instead.
func (*RequesterAuthenticatePayload) Descriptor() ([]byte, []int) {
	return file_handshake_handshake_proto_rawDescGZIP(), []int{2}
}

func (x *RequesterAuthenticatePayload) GetRequesterAccountId() []byte {
	if x != nil {
		return x.RequesterAccountId
	}
	return nil
}

func (x *RequesterAuthenticatePayload) GetRequesterAccountSig() []byte {
	if x != nil {
		return x.RequesterAccountSig
	}
	return nil
}

type ResponderAcceptPayload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResponderAccountSig []byte `protobuf:"bytes,1,opt,name=responder_account_sig,json=responderAccountSig,proto3" json:"responder_account_sig,omitempty"`
}

func (x *ResponderAcceptPayload) Reset() {
	*x = ResponderAcceptPayload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_handshake_handshake_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponderAcceptPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponderAcceptPayload) ProtoMessage() {}

func (x *ResponderAcceptPayload) ProtoReflect() protoreflect.Message {
	mi := &file_handshake_handshake_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponderAcceptPayload.ProtoReflect.Descriptor instead.
func (*ResponderAcceptPayload) Descriptor() ([]byte, []int) {
	return file_handshake_handshake_proto_rawDescGZIP(), []int{3}
}

func (x *ResponderAcceptPayload) GetResponderAccountSig() []byte {
	if x != nil {
		return x.ResponderAccountSig
	}
	return nil
}

type RequesterAcknowledgePayload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *RequesterAcknowledgePayload) Reset() {
	*x = RequesterAcknowledgePayload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_handshake_handshake_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequesterAcknowledgePayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequesterAcknowledgePayload) ProtoMessage() {}

func (x *RequesterAcknowledgePayload) ProtoReflect() protoreflect.Message {
	mi := &file_handshake_handshake_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequesterAcknowledgePayload.ProtoReflect.Descriptor instead.
func (*RequesterAcknowledgePayload) Descriptor() ([]byte, []int) {
	return file_handshake_handshake_proto_rawDescGZIP(), []int{4}
}

func (x *RequesterAcknowledgePayload) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_handshake_handshake_proto protoreflect.FileDescriptor

var file_handshake_handshake_proto_rawDesc = []byte{
	0x0a, 0x19, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x2f, 0x68, 0x61, 0x6e, 0x64,
	0x73, 0x68, 0x61, 0x6b, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x68, 0x61, 0x6e,
	0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x22, 0x1f, 0x0a, 0x0b, 0x42, 0x6f, 0x78, 0x45, 0x6e, 0x76,
	0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x62, 0x6f, 0x78, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x03, 0x62, 0x6f, 0x78, 0x22, 0x3a, 0x0a, 0x0c, 0x48, 0x65, 0x6c, 0x6c, 0x6f,
	0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x65, 0x70, 0x68, 0x65, 0x6d,
	0x65, 0x72, 0x61, 0x6c, 0x5f, 0x70, 0x75, 0x62, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x0f, 0x65, 0x70, 0x68, 0x65, 0x6d, 0x65, 0x72, 0x61, 0x6c, 0x50, 0x75, 0x62,
	0x4b, 0x65, 0x79, 0x22, 0x84, 0x01, 0x0a, 0x1c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65,
	0x72, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x50, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x12, 0x30, 0x0a, 0x14, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65,
	0x72, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x12, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x72, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x32, 0x0a, 0x15, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x65, 0x72, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x69, 0x67, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x13, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x72,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x69, 0x67, 0x22, 0x4c, 0x0a, 0x16, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x64, 0x65, 0x72, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x50, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x12, 0x32, 0x0a, 0x15, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x65,
	0x72, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x69, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x13, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x65, 0x72, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x69, 0x67, 0x22, 0x37, 0x0a, 0x1b, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x65, 0x72, 0x41, 0x63, 0x6b, 0x6e, 0x6f, 0x77, 0x6c, 0x65, 0x64, 0x67, 0x65,
	0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x42, 0x27, 0x5a, 0x25, 0x62, 0x65, 0x72, 0x74, 0x79, 0x2e, 0x74, 0x65, 0x63, 0x68, 0x2f,
	0x77, 0x65, 0x73, 0x68, 0x6e, 0x65, 0x74, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_handshake_handshake_proto_rawDescOnce sync.Once
	file_handshake_handshake_proto_rawDescData = file_handshake_handshake_proto_rawDesc
)

func file_handshake_handshake_proto_rawDescGZIP() []byte {
	file_handshake_handshake_proto_rawDescOnce.Do(func() {
		file_handshake_handshake_proto_rawDescData = protoimpl.X.CompressGZIP(file_handshake_handshake_proto_rawDescData)
	})
	return file_handshake_handshake_proto_rawDescData
}

var file_handshake_handshake_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_handshake_handshake_proto_goTypes = []interface{}{
	(*BoxEnvelope)(nil),                  // 0: handshake.BoxEnvelope
	(*HelloPayload)(nil),                 // 1: handshake.HelloPayload
	(*RequesterAuthenticatePayload)(nil), // 2: handshake.RequesterAuthenticatePayload
	(*ResponderAcceptPayload)(nil),       // 3: handshake.ResponderAcceptPayload
	(*RequesterAcknowledgePayload)(nil),  // 4: handshake.RequesterAcknowledgePayload
}
var file_handshake_handshake_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_handshake_handshake_proto_init() }
func file_handshake_handshake_proto_init() {
	if File_handshake_handshake_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_handshake_handshake_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoxEnvelope); i {
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
		file_handshake_handshake_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HelloPayload); i {
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
		file_handshake_handshake_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequesterAuthenticatePayload); i {
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
		file_handshake_handshake_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponderAcceptPayload); i {
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
		file_handshake_handshake_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequesterAcknowledgePayload); i {
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
			RawDescriptor: file_handshake_handshake_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_handshake_handshake_proto_goTypes,
		DependencyIndexes: file_handshake_handshake_proto_depIdxs,
		MessageInfos:      file_handshake_handshake_proto_msgTypes,
	}.Build()
	File_handshake_handshake_proto = out.File
	file_handshake_handshake_proto_rawDesc = nil
	file_handshake_handshake_proto_goTypes = nil
	file_handshake_handshake_proto_depIdxs = nil
}