// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: outofstoremessagetypes/outofstoremessage.proto

package outofstoremessagetypes

import (
	protocoltypes "berty.tech/weshnet/v2/pkg/protocoltypes"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_outofstoremessagetypes_outofstoremessage_proto protoreflect.FileDescriptor

var file_outofstoremessagetypes_outofstoremessage_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x6f, 0x75, 0x74, 0x6f, 0x66, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x6f, 0x75, 0x74, 0x6f, 0x66, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x1c, 0x77, 0x65, 0x73, 0x68, 0x6e, 0x65, 0x74, 0x2e, 0x6f, 0x75, 0x74, 0x6f, 0x66, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x13,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x32, 0x8d, 0x01, 0x0a, 0x18, 0x4f, 0x75, 0x74, 0x4f, 0x66, 0x53, 0x74, 0x6f,
	0x72, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x71, 0x0a, 0x11, 0x4f, 0x75, 0x74, 0x4f, 0x66, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x12, 0x2e, 0x2e, 0x77, 0x65, 0x73, 0x68, 0x6e, 0x65, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x75, 0x74, 0x4f,
	0x66, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x77, 0x65, 0x73, 0x68, 0x6e, 0x65, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x75, 0x74, 0x4f,
	0x66, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x2e, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x42, 0x32, 0x5a, 0x30, 0x62, 0x65, 0x72, 0x74, 0x79, 0x2e, 0x74, 0x65, 0x63,
	0x68, 0x2f, 0x77, 0x65, 0x73, 0x68, 0x6e, 0x65, 0x74, 0x2f, 0x76, 0x32, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x6f, 0x75, 0x74, 0x6f, 0x66, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x74, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_outofstoremessagetypes_outofstoremessage_proto_goTypes = []any{
	(*protocoltypes.OutOfStoreReceive_Request)(nil), // 0: weshnet.protocol.v1.OutOfStoreReceive.Request
	(*protocoltypes.OutOfStoreReceive_Reply)(nil),   // 1: weshnet.protocol.v1.OutOfStoreReceive.Reply
}
var file_outofstoremessagetypes_outofstoremessage_proto_depIdxs = []int32{
	0, // 0: weshnet.outofstoremessage.v1.OutOfStoreMessageService.OutOfStoreReceive:input_type -> weshnet.protocol.v1.OutOfStoreReceive.Request
	1, // 1: weshnet.outofstoremessage.v1.OutOfStoreMessageService.OutOfStoreReceive:output_type -> weshnet.protocol.v1.OutOfStoreReceive.Reply
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_outofstoremessagetypes_outofstoremessage_proto_init() }
func file_outofstoremessagetypes_outofstoremessage_proto_init() {
	if File_outofstoremessagetypes_outofstoremessage_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_outofstoremessagetypes_outofstoremessage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_outofstoremessagetypes_outofstoremessage_proto_goTypes,
		DependencyIndexes: file_outofstoremessagetypes_outofstoremessage_proto_depIdxs,
	}.Build()
	File_outofstoremessagetypes_outofstoremessage_proto = out.File
	file_outofstoremessagetypes_outofstoremessage_proto_rawDesc = nil
	file_outofstoremessagetypes_outofstoremessage_proto_goTypes = nil
	file_outofstoremessagetypes_outofstoremessage_proto_depIdxs = nil
}
