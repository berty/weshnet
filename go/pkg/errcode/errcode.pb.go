// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: errcode.proto

package errcode

import (
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ErrCode int32

const (
	Undefined                                  ErrCode = 0
	TODO                                       ErrCode = 666
	ErrNotImplemented                          ErrCode = 777
	ErrInternal                                ErrCode = 888
	ErrInvalidInput                            ErrCode = 100
	ErrMissingInput                            ErrCode = 101
	ErrSerialization                           ErrCode = 102
	ErrDeserialization                         ErrCode = 103
	ErrStreamRead                              ErrCode = 104
	ErrStreamWrite                             ErrCode = 105
	ErrMissingMapKey                           ErrCode = 106
	ErrCryptoRandomGeneration                  ErrCode = 200
	ErrCryptoKeyGeneration                     ErrCode = 201
	ErrCryptoNonceGeneration                   ErrCode = 202
	ErrCryptoSignature                         ErrCode = 203
	ErrCryptoSignatureVerification             ErrCode = 204
	ErrCryptoDecrypt                           ErrCode = 205
	ErrCryptoEncrypt                           ErrCode = 206
	ErrCryptoKeyConversion                     ErrCode = 207
	ErrOrbitDBInit                             ErrCode = 1000
	ErrOrbitDBOpen                             ErrCode = 1001
	ErrOrbitDBAppend                           ErrCode = 1002
	ErrOrbitDBDeserialization                  ErrCode = 1003
	ErrOrbitDBStoreCast                        ErrCode = 1004
	ErrHandshakeOwnEphemeralKeyGenSend         ErrCode = 1100
	ErrHandshakePeerEphemeralKeyRecv           ErrCode = 1101
	ErrHandshakeRequesterAuthenticateBoxKeyGen ErrCode = 1102
	ErrHandshakeResponderAcceptBoxKeyGen       ErrCode = 1103
	ErrHandshakeRequesterHello                 ErrCode = 1104
	ErrHandshakeResponderHello                 ErrCode = 1105
	ErrHandshakeRequesterAuthenticate          ErrCode = 1106
	ErrHandshakeResponderAccept                ErrCode = 1107
	ErrGroupMemberLogEventOpen                 ErrCode = 1200
	ErrGroupMemberLogEventSignature            ErrCode = 1201
	ErrGroupMemberUnknownGroupID               ErrCode = 1202
	ErrGroupSecretOtherDestMember              ErrCode = 1203
	ErrGroupSecretAlreadySentToMember          ErrCode = 1204
	ErrGroupInvalidType                        ErrCode = 1205
	ErrGroupMissing                            ErrCode = 1206
	ErrMessageKeyPersistencePut                ErrCode = 1300
	ErrMessageKeyPersistenceGet                ErrCode = 1301
	ErrBridgeInterrupted                       ErrCode = 1400
	ErrBridgeNotRunning                        ErrCode = 1401
)

var ErrCode_name = map[int32]string{
	0:    "Undefined",
	666:  "TODO",
	777:  "ErrNotImplemented",
	888:  "ErrInternal",
	100:  "ErrInvalidInput",
	101:  "ErrMissingInput",
	102:  "ErrSerialization",
	103:  "ErrDeserialization",
	104:  "ErrStreamRead",
	105:  "ErrStreamWrite",
	106:  "ErrMissingMapKey",
	200:  "ErrCryptoRandomGeneration",
	201:  "ErrCryptoKeyGeneration",
	202:  "ErrCryptoNonceGeneration",
	203:  "ErrCryptoSignature",
	204:  "ErrCryptoSignatureVerification",
	205:  "ErrCryptoDecrypt",
	206:  "ErrCryptoEncrypt",
	207:  "ErrCryptoKeyConversion",
	1000: "ErrOrbitDBInit",
	1001: "ErrOrbitDBOpen",
	1002: "ErrOrbitDBAppend",
	1003: "ErrOrbitDBDeserialization",
	1004: "ErrOrbitDBStoreCast",
	1100: "ErrHandshakeOwnEphemeralKeyGenSend",
	1101: "ErrHandshakePeerEphemeralKeyRecv",
	1102: "ErrHandshakeRequesterAuthenticateBoxKeyGen",
	1103: "ErrHandshakeResponderAcceptBoxKeyGen",
	1104: "ErrHandshakeRequesterHello",
	1105: "ErrHandshakeResponderHello",
	1106: "ErrHandshakeRequesterAuthenticate",
	1107: "ErrHandshakeResponderAccept",
	1200: "ErrGroupMemberLogEventOpen",
	1201: "ErrGroupMemberLogEventSignature",
	1202: "ErrGroupMemberUnknownGroupID",
	1203: "ErrGroupSecretOtherDestMember",
	1204: "ErrGroupSecretAlreadySentToMember",
	1205: "ErrGroupInvalidType",
	1206: "ErrGroupMissing",
	1300: "ErrMessageKeyPersistencePut",
	1301: "ErrMessageKeyPersistenceGet",
	1400: "ErrBridgeInterrupted",
	1401: "ErrBridgeNotRunning",
}

var ErrCode_value = map[string]int32{
	"Undefined":                                  0,
	"TODO":                                       666,
	"ErrNotImplemented":                          777,
	"ErrInternal":                                888,
	"ErrInvalidInput":                            100,
	"ErrMissingInput":                            101,
	"ErrSerialization":                           102,
	"ErrDeserialization":                         103,
	"ErrStreamRead":                              104,
	"ErrStreamWrite":                             105,
	"ErrMissingMapKey":                           106,
	"ErrCryptoRandomGeneration":                  200,
	"ErrCryptoKeyGeneration":                     201,
	"ErrCryptoNonceGeneration":                   202,
	"ErrCryptoSignature":                         203,
	"ErrCryptoSignatureVerification":             204,
	"ErrCryptoDecrypt":                           205,
	"ErrCryptoEncrypt":                           206,
	"ErrCryptoKeyConversion":                     207,
	"ErrOrbitDBInit":                             1000,
	"ErrOrbitDBOpen":                             1001,
	"ErrOrbitDBAppend":                           1002,
	"ErrOrbitDBDeserialization":                  1003,
	"ErrOrbitDBStoreCast":                        1004,
	"ErrHandshakeOwnEphemeralKeyGenSend":         1100,
	"ErrHandshakePeerEphemeralKeyRecv":           1101,
	"ErrHandshakeRequesterAuthenticateBoxKeyGen": 1102,
	"ErrHandshakeResponderAcceptBoxKeyGen":       1103,
	"ErrHandshakeRequesterHello":                 1104,
	"ErrHandshakeResponderHello":                 1105,
	"ErrHandshakeRequesterAuthenticate":          1106,
	"ErrHandshakeResponderAccept":                1107,
	"ErrGroupMemberLogEventOpen":                 1200,
	"ErrGroupMemberLogEventSignature":            1201,
	"ErrGroupMemberUnknownGroupID":               1202,
	"ErrGroupSecretOtherDestMember":              1203,
	"ErrGroupSecretAlreadySentToMember":          1204,
	"ErrGroupInvalidType":                        1205,
	"ErrGroupMissing":                            1206,
	"ErrMessageKeyPersistencePut":                1300,
	"ErrMessageKeyPersistenceGet":                1301,
	"ErrBridgeInterrupted":                       1400,
	"ErrBridgeNotRunning":                        1401,
}

func (x ErrCode) String() string {
	return proto.EnumName(ErrCode_name, int32(x))
}

func (ErrCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4240057316120df7, []int{0}
}

func init() {
	proto.RegisterEnum("berty.errcode.ErrCode", ErrCode_name, ErrCode_value)
}

func init() { proto.RegisterFile("errcode.proto", fileDescriptor_4240057316120df7) }

var fileDescriptor_4240057316120df7 = []byte{
	// 752 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0x4b, 0x53, 0xe3, 0x46,
	0x10, 0xc6, 0x55, 0x09, 0xb6, 0x26, 0xe5, 0x64, 0x68, 0x08, 0xe1, 0x11, 0x9e, 0x81, 0xa4, 0x42,
	0x55, 0xf0, 0x21, 0xbf, 0xc0, 0xc6, 0x2a, 0x70, 0x11, 0x63, 0xca, 0x86, 0xa4, 0x2a, 0x37, 0x59,
	0x6a, 0xcb, 0x13, 0xe4, 0x19, 0xa5, 0x35, 0x32, 0x71, 0xfe, 0x41, 0xee, 0xc9, 0x25, 0xbf, 0x22,
	0xef, 0xdf, 0x90, 0x07, 0xaf, 0xec, 0x9e, 0xf6, 0xb8, 0xb7, 0x7d, 0xfd, 0x00, 0xf6, 0xb6, 0x25,
	0x8d, 0x00, 0x7b, 0xa1, 0xd8, 0x93, 0xa4, 0xef, 0xfb, 0xfa, 0xeb, 0x56, 0xcf, 0x74, 0xb3, 0x22,
	0x12, 0xb9, 0xca, 0xc3, 0xcd, 0x90, 0x94, 0x56, 0x50, 0x6c, 0x23, 0xe9, 0xc1, 0x66, 0x06, 0xce,
	0x7d, 0xe6, 0x0b, 0xdd, 0x8d, 0xdb, 0x9b, 0xae, 0xea, 0x95, 0x7c, 0xe5, 0xab, 0x52, 0xaa, 0x6a,
	0xc7, 0x9d, 0xf4, 0x2b, 0xfd, 0x48, 0xdf, 0x4c, 0xf4, 0xc6, 0x23, 0x8b, 0xe5, 0x6d, 0xa2, 0x2d,
	0xe5, 0x21, 0x14, 0x99, 0x75, 0x28, 0x3d, 0xec, 0x08, 0x89, 0x1e, 0x1f, 0x03, 0x8b, 0xbd, 0x75,
	0xd0, 0xa8, 0x36, 0xf8, 0xcf, 0x6f, 0xc3, 0x34, 0x9b, 0xb0, 0x89, 0xf6, 0x94, 0xae, 0xf5, 0xc2,
	0x00, 0x7b, 0x28, 0x35, 0x7a, 0xfc, 0x87, 0x71, 0xe0, 0xec, 0x1d, 0x9b, 0xa8, 0x26, 0x35, 0x92,
	0x74, 0x02, 0x7e, 0x39, 0x0e, 0x93, 0xec, 0xbd, 0x14, 0xe9, 0x3b, 0x81, 0xf0, 0x6a, 0x32, 0x8c,
	0x35, 0xf7, 0x32, 0xb0, 0x2e, 0xa2, 0x48, 0x48, 0xdf, 0x80, 0x08, 0x53, 0x8c, 0xdb, 0x44, 0x2d,
	0x24, 0xe1, 0x04, 0xe2, 0x7b, 0x47, 0x0b, 0x25, 0x79, 0x07, 0xa6, 0x19, 0xd8, 0x44, 0x55, 0x8c,
	0x46, 0x70, 0x1f, 0x26, 0x58, 0x31, 0x51, 0x6b, 0x42, 0xa7, 0xd7, 0x44, 0xc7, 0xe3, 0x5d, 0x00,
	0xf6, 0xee, 0x35, 0xf4, 0x15, 0x09, 0x8d, 0x5c, 0x64, 0xa6, 0x59, 0xa6, 0xba, 0x13, 0xee, 0xe2,
	0x80, 0x7f, 0x03, 0x8b, 0x6c, 0x36, 0xf9, 0x47, 0x1a, 0x84, 0x5a, 0x35, 0x1d, 0xe9, 0xa9, 0xde,
	0x36, 0x4a, 0x24, 0xe3, 0xfd, 0x77, 0x0e, 0xe6, 0xd9, 0xf4, 0x35, 0xbf, 0x8b, 0x83, 0x21, 0xf2,
	0x9f, 0x1c, 0x2c, 0xb0, 0x99, 0x6b, 0x72, 0x4f, 0x49, 0x17, 0x87, 0xe8, 0x7f, 0x73, 0xf0, 0x41,
	0x5a, 0xb0, 0xa1, 0x5b, 0xc2, 0x97, 0x8e, 0x8e, 0x09, 0xf9, 0x7f, 0x39, 0xf8, 0x88, 0x2d, 0xde,
	0x26, 0xbe, 0x44, 0x12, 0x1d, 0xe1, 0x9a, 0xe8, 0x93, 0x1c, 0xbc, 0x9f, 0xd6, 0x6b, 0x44, 0x55,
	0x74, 0x93, 0x27, 0x3f, 0x1d, 0x85, 0x6d, 0x69, 0xe0, 0xb3, 0x5b, 0x75, 0x6e, 0x29, 0xd9, 0x47,
	0x8a, 0x12, 0xab, 0xf3, 0x1c, 0x4c, 0xa6, 0xed, 0x68, 0x50, 0x5b, 0xe8, 0x6a, 0xa5, 0x26, 0x85,
	0xe6, 0x4f, 0xf2, 0xa3, 0x60, 0x23, 0x44, 0xc9, 0x9f, 0xe6, 0x33, 0xf7, 0x0c, 0x2c, 0x87, 0x21,
	0x4a, 0x8f, 0x3f, 0xcb, 0x67, 0x5d, 0xca, 0xe0, 0xd7, 0x4f, 0xe0, 0x79, 0x1e, 0x66, 0xd8, 0xe4,
	0x0d, 0xdf, 0xd2, 0x8a, 0x70, 0xcb, 0x89, 0x34, 0x7f, 0x91, 0x87, 0x4f, 0xd8, 0xaa, 0x4d, 0xb4,
	0xe3, 0x48, 0x2f, 0xea, 0x3a, 0x47, 0xd8, 0x38, 0x96, 0x76, 0xd8, 0xc5, 0x1e, 0x92, 0x13, 0x98,
	0x76, 0xb6, 0x92, 0x14, 0x27, 0x05, 0x58, 0x67, 0xcb, 0xc3, 0xc2, 0x7d, 0x44, 0x1a, 0x56, 0x36,
	0xd1, 0xed, 0xf3, 0xd3, 0x02, 0x94, 0xd8, 0xc6, 0xb0, 0xac, 0x89, 0xdf, 0xc6, 0x18, 0x69, 0xa4,
	0x72, 0xac, 0xbb, 0x28, 0x75, 0xd2, 0x3f, 0xac, 0xa8, 0xef, 0x8c, 0x37, 0x3f, 0x2b, 0xc0, 0xa7,
	0x6c, 0x6d, 0x34, 0x20, 0x0a, 0x95, 0xf4, 0x90, 0xca, 0xae, 0x8b, 0xa1, 0xbe, 0x91, 0x9e, 0x17,
	0x60, 0x89, 0xcd, 0xdd, 0xe9, 0xbd, 0x83, 0x41, 0xa0, 0xf8, 0xc5, 0x1d, 0x82, 0xcc, 0xcb, 0x08,
	0xfe, 0x2f, 0xc0, 0xc7, 0x6c, 0xe5, 0x8d, 0xd5, 0xf1, 0x07, 0x05, 0x58, 0x66, 0xf3, 0xf7, 0x14,
	0xc5, 0x1f, 0x5e, 0xa5, 0xda, 0x26, 0x15, 0x87, 0x75, 0xec, 0xb5, 0x91, 0xbe, 0x50, 0xbe, 0xdd,
	0x47, 0xa9, 0xd3, 0x93, 0xfa, 0xc5, 0x82, 0x35, 0xb6, 0x74, 0xb7, 0xe0, 0xe6, 0xa6, 0xfd, 0x6a,
	0xc1, 0x0a, 0xfb, 0x70, 0x54, 0x75, 0x28, 0x8f, 0xa4, 0x3a, 0x96, 0x29, 0x52, 0xab, 0xf2, 0xdf,
	0x2c, 0x58, 0x65, 0x0b, 0x57, 0x92, 0x16, 0xba, 0x84, 0xba, 0xa1, 0xbb, 0x98, 0x8c, 0x99, 0x36,
	0x11, 0xfc, 0x77, 0x2b, 0xfb, 0xaf, 0x21, 0x4d, 0x39, 0x20, 0x74, 0xbc, 0x41, 0x0b, 0xa5, 0x3e,
	0x50, 0x99, 0xee, 0x0f, 0x2b, 0xbb, 0x07, 0xc6, 0xdc, 0xcc, 0xf9, 0xc1, 0x20, 0x44, 0xfe, 0xa7,
	0x05, 0x53, 0xe9, 0x9c, 0x9b, 0x42, 0xcc, 0x08, 0xf2, 0xbf, 0xac, 0xac, 0x0f, 0x75, 0x8c, 0x22,
	0xc7, 0xc7, 0x5d, 0x1c, 0xec, 0x27, 0x77, 0x36, 0xd2, 0x28, 0x5d, 0xdc, 0x8f, 0x35, 0xff, 0x91,
	0xdd, 0xa7, 0xd8, 0x46, 0xcd, 0x7f, 0x62, 0x30, 0xcb, 0xa6, 0x6c, 0xa2, 0x0a, 0x09, 0xcf, 0xc7,
	0x74, 0xdd, 0x50, 0x1c, 0x26, 0x3b, 0xe8, 0x92, 0x65, 0xe5, 0x18, 0x6a, 0x4f, 0xe9, 0x66, 0x2c,
	0x65, 0x92, 0xf8, 0x25, 0xab, 0xac, 0x5f, 0x3c, 0x5e, 0x1c, 0xfb, 0x7a, 0xc9, 0xec, 0x47, 0x8d,
	0x6e, 0xb7, 0x94, 0xbe, 0x96, 0x92, 0xa5, 0x78, 0xe4, 0x97, 0xb2, 0x8d, 0xd9, 0x1e, 0x4f, 0x37,
	0xe1, 0xe7, 0xaf, 0x02, 0x00, 0x00, 0xff, 0xff, 0x34, 0x2d, 0xff, 0x99, 0x58, 0x05, 0x00, 0x00,
}
