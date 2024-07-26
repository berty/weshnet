// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: errcode/errcode.proto

package errcode

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

type ErrCode int32

const (
	ErrCode_Undefined                                            ErrCode = 0   // default value, should never be set manually
	ErrCode_TODO                                                 ErrCode = 666 // indicates that you plan to create an error later
	ErrCode_ErrNotImplemented                                    ErrCode = 777 // indicates that a method is not implemented yet
	ErrCode_ErrInternal                                          ErrCode = 888 // indicates an unknown error (without Code), i.e. in gRPC
	ErrCode_ErrInvalidInput                                      ErrCode = 100
	ErrCode_ErrInvalidRange                                      ErrCode = 101
	ErrCode_ErrMissingInput                                      ErrCode = 102
	ErrCode_ErrSerialization                                     ErrCode = 103
	ErrCode_ErrDeserialization                                   ErrCode = 104
	ErrCode_ErrStreamRead                                        ErrCode = 105
	ErrCode_ErrStreamWrite                                       ErrCode = 106
	ErrCode_ErrStreamTransform                                   ErrCode = 110
	ErrCode_ErrStreamSendAndClose                                ErrCode = 111
	ErrCode_ErrStreamHeaderWrite                                 ErrCode = 112
	ErrCode_ErrStreamHeaderRead                                  ErrCode = 115
	ErrCode_ErrStreamSink                                        ErrCode = 113
	ErrCode_ErrStreamCloseAndRecv                                ErrCode = 114
	ErrCode_ErrMissingMapKey                                     ErrCode = 107
	ErrCode_ErrDBWrite                                           ErrCode = 108
	ErrCode_ErrDBRead                                            ErrCode = 109
	ErrCode_ErrDBDestroy                                         ErrCode = 120
	ErrCode_ErrDBMigrate                                         ErrCode = 121
	ErrCode_ErrDBReplay                                          ErrCode = 122
	ErrCode_ErrDBRestore                                         ErrCode = 123
	ErrCode_ErrDBOpen                                            ErrCode = 124
	ErrCode_ErrDBClose                                           ErrCode = 125
	ErrCode_ErrCryptoRandomGeneration                            ErrCode = 200
	ErrCode_ErrCryptoKeyGeneration                               ErrCode = 201
	ErrCode_ErrCryptoNonceGeneration                             ErrCode = 202
	ErrCode_ErrCryptoSignature                                   ErrCode = 203
	ErrCode_ErrCryptoSignatureVerification                       ErrCode = 204
	ErrCode_ErrCryptoDecrypt                                     ErrCode = 205
	ErrCode_ErrCryptoDecryptPayload                              ErrCode = 206
	ErrCode_ErrCryptoEncrypt                                     ErrCode = 207
	ErrCode_ErrCryptoKeyConversion                               ErrCode = 208
	ErrCode_ErrCryptoCipherInit                                  ErrCode = 209
	ErrCode_ErrCryptoKeyDerivation                               ErrCode = 210
	ErrCode_ErrMap                                               ErrCode = 300
	ErrCode_ErrForEach                                           ErrCode = 301
	ErrCode_ErrKeystoreGet                                       ErrCode = 400
	ErrCode_ErrKeystorePut                                       ErrCode = 401
	ErrCode_ErrNotFound                                          ErrCode = 404 // generic
	ErrCode_ErrOrbitDBInit                                       ErrCode = 1000
	ErrCode_ErrOrbitDBOpen                                       ErrCode = 1001
	ErrCode_ErrOrbitDBAppend                                     ErrCode = 1002
	ErrCode_ErrOrbitDBDeserialization                            ErrCode = 1003
	ErrCode_ErrOrbitDBStoreCast                                  ErrCode = 1004
	ErrCode_ErrHandshakeOwnEphemeralKeyGenSend                   ErrCode = 1100
	ErrCode_ErrHandshakePeerEphemeralKeyRecv                     ErrCode = 1101
	ErrCode_ErrHandshakeRequesterAuthenticateBoxKeyGen           ErrCode = 1102
	ErrCode_ErrHandshakeResponderAcceptBoxKeyGen                 ErrCode = 1103
	ErrCode_ErrHandshakeRequesterHello                           ErrCode = 1104
	ErrCode_ErrHandshakeResponderHello                           ErrCode = 1105
	ErrCode_ErrHandshakeRequesterAuthenticate                    ErrCode = 1106
	ErrCode_ErrHandshakeResponderAccept                          ErrCode = 1107
	ErrCode_ErrHandshakeRequesterAcknowledge                     ErrCode = 1108
	ErrCode_ErrContactRequestSameAccount                         ErrCode = 1200
	ErrCode_ErrContactRequestContactAlreadyAdded                 ErrCode = 1201
	ErrCode_ErrContactRequestContactBlocked                      ErrCode = 1202
	ErrCode_ErrContactRequestContactUndefined                    ErrCode = 1203
	ErrCode_ErrContactRequestIncomingAlreadyReceived             ErrCode = 1204
	ErrCode_ErrGroupMemberLogEventOpen                           ErrCode = 1300
	ErrCode_ErrGroupMemberLogEventSignature                      ErrCode = 1301
	ErrCode_ErrGroupMemberUnknownGroupID                         ErrCode = 1302
	ErrCode_ErrGroupSecretOtherDestMember                        ErrCode = 1303
	ErrCode_ErrGroupSecretAlreadySentToMember                    ErrCode = 1304
	ErrCode_ErrGroupInvalidType                                  ErrCode = 1305
	ErrCode_ErrGroupMissing                                      ErrCode = 1306
	ErrCode_ErrGroupActivate                                     ErrCode = 1307
	ErrCode_ErrGroupDeactivate                                   ErrCode = 1308
	ErrCode_ErrGroupInfo                                         ErrCode = 1309
	ErrCode_ErrGroupUnknown                                      ErrCode = 1310
	ErrCode_ErrGroupOpen                                         ErrCode = 1311
	ErrCode_ErrMessageKeyPersistencePut                          ErrCode = 1500
	ErrCode_ErrMessageKeyPersistenceGet                          ErrCode = 1501
	ErrCode_ErrServiceReplication                                ErrCode = 4100
	ErrCode_ErrServiceReplicationServer                          ErrCode = 4101
	ErrCode_ErrServiceReplicationMissingEndpoint                 ErrCode = 4102
	ErrCode_ErrServicesDirectory                                 ErrCode = 4200
	ErrCode_ErrServicesDirectoryInvalidVerifiedCredentialSubject ErrCode = 4201
	ErrCode_ErrServicesDirectoryExistingRecordNotFound           ErrCode = 4202
	ErrCode_ErrServicesDirectoryRecordLockedAndCantBeReplaced    ErrCode = 4203
	ErrCode_ErrServicesDirectoryExplicitReplaceFlagRequired      ErrCode = 4204
	ErrCode_ErrServicesDirectoryInvalidVerifiedCredential        ErrCode = 4205
	ErrCode_ErrServicesDirectoryExpiredVerifiedCredential        ErrCode = 4206
	ErrCode_ErrServicesDirectoryInvalidVerifiedCredentialID      ErrCode = 4207
)

// Enum value maps for ErrCode.
var (
	ErrCode_name = map[int32]string{
		0:    "Undefined",
		666:  "TODO",
		777:  "ErrNotImplemented",
		888:  "ErrInternal",
		100:  "ErrInvalidInput",
		101:  "ErrInvalidRange",
		102:  "ErrMissingInput",
		103:  "ErrSerialization",
		104:  "ErrDeserialization",
		105:  "ErrStreamRead",
		106:  "ErrStreamWrite",
		110:  "ErrStreamTransform",
		111:  "ErrStreamSendAndClose",
		112:  "ErrStreamHeaderWrite",
		115:  "ErrStreamHeaderRead",
		113:  "ErrStreamSink",
		114:  "ErrStreamCloseAndRecv",
		107:  "ErrMissingMapKey",
		108:  "ErrDBWrite",
		109:  "ErrDBRead",
		120:  "ErrDBDestroy",
		121:  "ErrDBMigrate",
		122:  "ErrDBReplay",
		123:  "ErrDBRestore",
		124:  "ErrDBOpen",
		125:  "ErrDBClose",
		200:  "ErrCryptoRandomGeneration",
		201:  "ErrCryptoKeyGeneration",
		202:  "ErrCryptoNonceGeneration",
		203:  "ErrCryptoSignature",
		204:  "ErrCryptoSignatureVerification",
		205:  "ErrCryptoDecrypt",
		206:  "ErrCryptoDecryptPayload",
		207:  "ErrCryptoEncrypt",
		208:  "ErrCryptoKeyConversion",
		209:  "ErrCryptoCipherInit",
		210:  "ErrCryptoKeyDerivation",
		300:  "ErrMap",
		301:  "ErrForEach",
		400:  "ErrKeystoreGet",
		401:  "ErrKeystorePut",
		404:  "ErrNotFound",
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
		1108: "ErrHandshakeRequesterAcknowledge",
		1200: "ErrContactRequestSameAccount",
		1201: "ErrContactRequestContactAlreadyAdded",
		1202: "ErrContactRequestContactBlocked",
		1203: "ErrContactRequestContactUndefined",
		1204: "ErrContactRequestIncomingAlreadyReceived",
		1300: "ErrGroupMemberLogEventOpen",
		1301: "ErrGroupMemberLogEventSignature",
		1302: "ErrGroupMemberUnknownGroupID",
		1303: "ErrGroupSecretOtherDestMember",
		1304: "ErrGroupSecretAlreadySentToMember",
		1305: "ErrGroupInvalidType",
		1306: "ErrGroupMissing",
		1307: "ErrGroupActivate",
		1308: "ErrGroupDeactivate",
		1309: "ErrGroupInfo",
		1310: "ErrGroupUnknown",
		1311: "ErrGroupOpen",
		1500: "ErrMessageKeyPersistencePut",
		1501: "ErrMessageKeyPersistenceGet",
		4100: "ErrServiceReplication",
		4101: "ErrServiceReplicationServer",
		4102: "ErrServiceReplicationMissingEndpoint",
		4200: "ErrServicesDirectory",
		4201: "ErrServicesDirectoryInvalidVerifiedCredentialSubject",
		4202: "ErrServicesDirectoryExistingRecordNotFound",
		4203: "ErrServicesDirectoryRecordLockedAndCantBeReplaced",
		4204: "ErrServicesDirectoryExplicitReplaceFlagRequired",
		4205: "ErrServicesDirectoryInvalidVerifiedCredential",
		4206: "ErrServicesDirectoryExpiredVerifiedCredential",
		4207: "ErrServicesDirectoryInvalidVerifiedCredentialID",
	}
	ErrCode_value = map[string]int32{
		"Undefined":                          0,
		"TODO":                               666,
		"ErrNotImplemented":                  777,
		"ErrInternal":                        888,
		"ErrInvalidInput":                    100,
		"ErrInvalidRange":                    101,
		"ErrMissingInput":                    102,
		"ErrSerialization":                   103,
		"ErrDeserialization":                 104,
		"ErrStreamRead":                      105,
		"ErrStreamWrite":                     106,
		"ErrStreamTransform":                 110,
		"ErrStreamSendAndClose":              111,
		"ErrStreamHeaderWrite":               112,
		"ErrStreamHeaderRead":                115,
		"ErrStreamSink":                      113,
		"ErrStreamCloseAndRecv":              114,
		"ErrMissingMapKey":                   107,
		"ErrDBWrite":                         108,
		"ErrDBRead":                          109,
		"ErrDBDestroy":                       120,
		"ErrDBMigrate":                       121,
		"ErrDBReplay":                        122,
		"ErrDBRestore":                       123,
		"ErrDBOpen":                          124,
		"ErrDBClose":                         125,
		"ErrCryptoRandomGeneration":          200,
		"ErrCryptoKeyGeneration":             201,
		"ErrCryptoNonceGeneration":           202,
		"ErrCryptoSignature":                 203,
		"ErrCryptoSignatureVerification":     204,
		"ErrCryptoDecrypt":                   205,
		"ErrCryptoDecryptPayload":            206,
		"ErrCryptoEncrypt":                   207,
		"ErrCryptoKeyConversion":             208,
		"ErrCryptoCipherInit":                209,
		"ErrCryptoKeyDerivation":             210,
		"ErrMap":                             300,
		"ErrForEach":                         301,
		"ErrKeystoreGet":                     400,
		"ErrKeystorePut":                     401,
		"ErrNotFound":                        404,
		"ErrOrbitDBInit":                     1000,
		"ErrOrbitDBOpen":                     1001,
		"ErrOrbitDBAppend":                   1002,
		"ErrOrbitDBDeserialization":          1003,
		"ErrOrbitDBStoreCast":                1004,
		"ErrHandshakeOwnEphemeralKeyGenSend": 1100,
		"ErrHandshakePeerEphemeralKeyRecv":   1101,
		"ErrHandshakeRequesterAuthenticateBoxKeyGen":           1102,
		"ErrHandshakeResponderAcceptBoxKeyGen":                 1103,
		"ErrHandshakeRequesterHello":                           1104,
		"ErrHandshakeResponderHello":                           1105,
		"ErrHandshakeRequesterAuthenticate":                    1106,
		"ErrHandshakeResponderAccept":                          1107,
		"ErrHandshakeRequesterAcknowledge":                     1108,
		"ErrContactRequestSameAccount":                         1200,
		"ErrContactRequestContactAlreadyAdded":                 1201,
		"ErrContactRequestContactBlocked":                      1202,
		"ErrContactRequestContactUndefined":                    1203,
		"ErrContactRequestIncomingAlreadyReceived":             1204,
		"ErrGroupMemberLogEventOpen":                           1300,
		"ErrGroupMemberLogEventSignature":                      1301,
		"ErrGroupMemberUnknownGroupID":                         1302,
		"ErrGroupSecretOtherDestMember":                        1303,
		"ErrGroupSecretAlreadySentToMember":                    1304,
		"ErrGroupInvalidType":                                  1305,
		"ErrGroupMissing":                                      1306,
		"ErrGroupActivate":                                     1307,
		"ErrGroupDeactivate":                                   1308,
		"ErrGroupInfo":                                         1309,
		"ErrGroupUnknown":                                      1310,
		"ErrGroupOpen":                                         1311,
		"ErrMessageKeyPersistencePut":                          1500,
		"ErrMessageKeyPersistenceGet":                          1501,
		"ErrServiceReplication":                                4100,
		"ErrServiceReplicationServer":                          4101,
		"ErrServiceReplicationMissingEndpoint":                 4102,
		"ErrServicesDirectory":                                 4200,
		"ErrServicesDirectoryInvalidVerifiedCredentialSubject": 4201,
		"ErrServicesDirectoryExistingRecordNotFound":           4202,
		"ErrServicesDirectoryRecordLockedAndCantBeReplaced":    4203,
		"ErrServicesDirectoryExplicitReplaceFlagRequired":      4204,
		"ErrServicesDirectoryInvalidVerifiedCredential":        4205,
		"ErrServicesDirectoryExpiredVerifiedCredential":        4206,
		"ErrServicesDirectoryInvalidVerifiedCredentialID":      4207,
	}
)

func (x ErrCode) Enum() *ErrCode {
	p := new(ErrCode)
	*p = x
	return p
}

func (x ErrCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrCode) Descriptor() protoreflect.EnumDescriptor {
	return file_errcode_errcode_proto_enumTypes[0].Descriptor()
}

func (ErrCode) Type() protoreflect.EnumType {
	return &file_errcode_errcode_proto_enumTypes[0]
}

func (x ErrCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrCode.Descriptor instead.
func (ErrCode) EnumDescriptor() ([]byte, []int) {
	return file_errcode_errcode_proto_rawDescGZIP(), []int{0}
}

type ErrDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Codes []ErrCode `protobuf:"varint,1,rep,packed,name=codes,proto3,enum=weshnet.errcode.ErrCode" json:"codes,omitempty"`
}

func (x *ErrDetails) Reset() {
	*x = ErrDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_errcode_errcode_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrDetails) ProtoMessage() {}

func (x *ErrDetails) ProtoReflect() protoreflect.Message {
	mi := &file_errcode_errcode_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrDetails.ProtoReflect.Descriptor instead.
func (*ErrDetails) Descriptor() ([]byte, []int) {
	return file_errcode_errcode_proto_rawDescGZIP(), []int{0}
}

func (x *ErrDetails) GetCodes() []ErrCode {
	if x != nil {
		return x.Codes
	}
	return nil
}

var File_errcode_errcode_proto protoreflect.FileDescriptor

var file_errcode_errcode_proto_rawDesc = []byte{
	0x0a, 0x15, 0x65, 0x72, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x2f, 0x65, 0x72, 0x72, 0x63, 0x6f, 0x64,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x77, 0x65, 0x73, 0x68, 0x6e, 0x65, 0x74,
	0x2e, 0x65, 0x72, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x3c, 0x0a, 0x0a, 0x45, 0x72, 0x72, 0x44,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x2e, 0x0a, 0x05, 0x63, 0x6f, 0x64, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x77, 0x65, 0x73, 0x68, 0x6e, 0x65, 0x74, 0x2e,
	0x65, 0x72, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52,
	0x05, 0x63, 0x6f, 0x64, 0x65, 0x73, 0x2a, 0xdb, 0x13, 0x0a, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f,
	0x64, 0x65, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x6e, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x10,
	0x00, 0x12, 0x09, 0x0a, 0x04, 0x54, 0x4f, 0x44, 0x4f, 0x10, 0x9a, 0x05, 0x12, 0x16, 0x0a, 0x11,
	0x45, 0x72, 0x72, 0x4e, 0x6f, 0x74, 0x49, 0x6d, 0x70, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x65,
	0x64, 0x10, 0x89, 0x06, 0x12, 0x10, 0x0a, 0x0b, 0x45, 0x72, 0x72, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x10, 0xf8, 0x06, 0x12, 0x13, 0x0a, 0x0f, 0x45, 0x72, 0x72, 0x49, 0x6e, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x10, 0x64, 0x12, 0x13, 0x0a, 0x0f, 0x45,
	0x72, 0x72, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x10, 0x65,
	0x12, 0x13, 0x0a, 0x0f, 0x45, 0x72, 0x72, 0x4d, 0x69, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x49, 0x6e,
	0x70, 0x75, 0x74, 0x10, 0x66, 0x12, 0x14, 0x0a, 0x10, 0x45, 0x72, 0x72, 0x53, 0x65, 0x72, 0x69,
	0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0x67, 0x12, 0x16, 0x0a, 0x12, 0x45,
	0x72, 0x72, 0x44, 0x65, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x10, 0x68, 0x12, 0x11, 0x0a, 0x0d, 0x45, 0x72, 0x72, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x52, 0x65, 0x61, 0x64, 0x10, 0x69, 0x12, 0x12, 0x0a, 0x0e, 0x45, 0x72, 0x72, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x57, 0x72, 0x69, 0x74, 0x65, 0x10, 0x6a, 0x12, 0x16, 0x0a, 0x12, 0x45, 0x72,
	0x72, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d,
	0x10, 0x6e, 0x12, 0x19, 0x0a, 0x15, 0x45, 0x72, 0x72, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53,
	0x65, 0x6e, 0x64, 0x41, 0x6e, 0x64, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x10, 0x6f, 0x12, 0x18, 0x0a,
	0x14, 0x45, 0x72, 0x72, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x57, 0x72, 0x69, 0x74, 0x65, 0x10, 0x70, 0x12, 0x17, 0x0a, 0x13, 0x45, 0x72, 0x72, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x65, 0x61, 0x64, 0x10, 0x73,
	0x12, 0x11, 0x0a, 0x0d, 0x45, 0x72, 0x72, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x69, 0x6e,
	0x6b, 0x10, 0x71, 0x12, 0x19, 0x0a, 0x15, 0x45, 0x72, 0x72, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x43, 0x6c, 0x6f, 0x73, 0x65, 0x41, 0x6e, 0x64, 0x52, 0x65, 0x63, 0x76, 0x10, 0x72, 0x12, 0x14,
	0x0a, 0x10, 0x45, 0x72, 0x72, 0x4d, 0x69, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x4d, 0x61, 0x70, 0x4b,
	0x65, 0x79, 0x10, 0x6b, 0x12, 0x0e, 0x0a, 0x0a, 0x45, 0x72, 0x72, 0x44, 0x42, 0x57, 0x72, 0x69,
	0x74, 0x65, 0x10, 0x6c, 0x12, 0x0d, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x44, 0x42, 0x52, 0x65, 0x61,
	0x64, 0x10, 0x6d, 0x12, 0x10, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x44, 0x42, 0x44, 0x65, 0x73, 0x74,
	0x72, 0x6f, 0x79, 0x10, 0x78, 0x12, 0x10, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x44, 0x42, 0x4d, 0x69,
	0x67, 0x72, 0x61, 0x74, 0x65, 0x10, 0x79, 0x12, 0x0f, 0x0a, 0x0b, 0x45, 0x72, 0x72, 0x44, 0x42,
	0x52, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x10, 0x7a, 0x12, 0x10, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x44,
	0x42, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x10, 0x7b, 0x12, 0x0d, 0x0a, 0x09, 0x45, 0x72,
	0x72, 0x44, 0x42, 0x4f, 0x70, 0x65, 0x6e, 0x10, 0x7c, 0x12, 0x0e, 0x0a, 0x0a, 0x45, 0x72, 0x72,
	0x44, 0x42, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x10, 0x7d, 0x12, 0x1e, 0x0a, 0x19, 0x45, 0x72, 0x72,
	0x43, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x47, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0xc8, 0x01, 0x12, 0x1b, 0x0a, 0x16, 0x45, 0x72, 0x72,
	0x43, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x4b, 0x65, 0x79, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x10, 0xc9, 0x01, 0x12, 0x1d, 0x0a, 0x18, 0x45, 0x72, 0x72, 0x43, 0x72, 0x79,
	0x70, 0x74, 0x6f, 0x4e, 0x6f, 0x6e, 0x63, 0x65, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x10, 0xca, 0x01, 0x12, 0x17, 0x0a, 0x12, 0x45, 0x72, 0x72, 0x43, 0x72, 0x79, 0x70,
	0x74, 0x6f, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x10, 0xcb, 0x01, 0x12, 0x23,
	0x0a, 0x1e, 0x45, 0x72, 0x72, 0x43, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x53, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x10, 0xcc, 0x01, 0x12, 0x15, 0x0a, 0x10, 0x45, 0x72, 0x72, 0x43, 0x72, 0x79, 0x70, 0x74, 0x6f,
	0x44, 0x65, 0x63, 0x72, 0x79, 0x70, 0x74, 0x10, 0xcd, 0x01, 0x12, 0x1c, 0x0a, 0x17, 0x45, 0x72,
	0x72, 0x43, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x44, 0x65, 0x63, 0x72, 0x79, 0x70, 0x74, 0x50, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x10, 0xce, 0x01, 0x12, 0x15, 0x0a, 0x10, 0x45, 0x72, 0x72, 0x43,
	0x72, 0x79, 0x70, 0x74, 0x6f, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x10, 0xcf, 0x01, 0x12,
	0x1b, 0x0a, 0x16, 0x45, 0x72, 0x72, 0x43, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x4b, 0x65, 0x79, 0x43,
	0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x10, 0xd0, 0x01, 0x12, 0x18, 0x0a, 0x13,
	0x45, 0x72, 0x72, 0x43, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x43, 0x69, 0x70, 0x68, 0x65, 0x72, 0x49,
	0x6e, 0x69, 0x74, 0x10, 0xd1, 0x01, 0x12, 0x1b, 0x0a, 0x16, 0x45, 0x72, 0x72, 0x43, 0x72, 0x79,
	0x70, 0x74, 0x6f, 0x4b, 0x65, 0x79, 0x44, 0x65, 0x72, 0x69, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x10, 0xd2, 0x01, 0x12, 0x0b, 0x0a, 0x06, 0x45, 0x72, 0x72, 0x4d, 0x61, 0x70, 0x10, 0xac, 0x02,
	0x12, 0x0f, 0x0a, 0x0a, 0x45, 0x72, 0x72, 0x46, 0x6f, 0x72, 0x45, 0x61, 0x63, 0x68, 0x10, 0xad,
	0x02, 0x12, 0x13, 0x0a, 0x0e, 0x45, 0x72, 0x72, 0x4b, 0x65, 0x79, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x47, 0x65, 0x74, 0x10, 0x90, 0x03, 0x12, 0x13, 0x0a, 0x0e, 0x45, 0x72, 0x72, 0x4b, 0x65, 0x79,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x50, 0x75, 0x74, 0x10, 0x91, 0x03, 0x12, 0x10, 0x0a, 0x0b, 0x45,
	0x72, 0x72, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10, 0x94, 0x03, 0x12, 0x13, 0x0a,
	0x0e, 0x45, 0x72, 0x72, 0x4f, 0x72, 0x62, 0x69, 0x74, 0x44, 0x42, 0x49, 0x6e, 0x69, 0x74, 0x10,
	0xe8, 0x07, 0x12, 0x13, 0x0a, 0x0e, 0x45, 0x72, 0x72, 0x4f, 0x72, 0x62, 0x69, 0x74, 0x44, 0x42,
	0x4f, 0x70, 0x65, 0x6e, 0x10, 0xe9, 0x07, 0x12, 0x15, 0x0a, 0x10, 0x45, 0x72, 0x72, 0x4f, 0x72,
	0x62, 0x69, 0x74, 0x44, 0x42, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x10, 0xea, 0x07, 0x12, 0x1e,
	0x0a, 0x19, 0x45, 0x72, 0x72, 0x4f, 0x72, 0x62, 0x69, 0x74, 0x44, 0x42, 0x44, 0x65, 0x73, 0x65,
	0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0xeb, 0x07, 0x12, 0x18,
	0x0a, 0x13, 0x45, 0x72, 0x72, 0x4f, 0x72, 0x62, 0x69, 0x74, 0x44, 0x42, 0x53, 0x74, 0x6f, 0x72,
	0x65, 0x43, 0x61, 0x73, 0x74, 0x10, 0xec, 0x07, 0x12, 0x27, 0x0a, 0x22, 0x45, 0x72, 0x72, 0x48,
	0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x4f, 0x77, 0x6e, 0x45, 0x70, 0x68, 0x65, 0x6d,
	0x65, 0x72, 0x61, 0x6c, 0x4b, 0x65, 0x79, 0x47, 0x65, 0x6e, 0x53, 0x65, 0x6e, 0x64, 0x10, 0xcc,
	0x08, 0x12, 0x25, 0x0a, 0x20, 0x45, 0x72, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b,
	0x65, 0x50, 0x65, 0x65, 0x72, 0x45, 0x70, 0x68, 0x65, 0x6d, 0x65, 0x72, 0x61, 0x6c, 0x4b, 0x65,
	0x79, 0x52, 0x65, 0x63, 0x76, 0x10, 0xcd, 0x08, 0x12, 0x2f, 0x0a, 0x2a, 0x45, 0x72, 0x72, 0x48,
	0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65,
	0x72, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x42, 0x6f, 0x78,
	0x4b, 0x65, 0x79, 0x47, 0x65, 0x6e, 0x10, 0xce, 0x08, 0x12, 0x29, 0x0a, 0x24, 0x45, 0x72, 0x72,
	0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64,
	0x65, 0x72, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x42, 0x6f, 0x78, 0x4b, 0x65, 0x79, 0x47, 0x65,
	0x6e, 0x10, 0xcf, 0x08, 0x12, 0x1f, 0x0a, 0x1a, 0x45, 0x72, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x73,
	0x68, 0x61, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x72, 0x48, 0x65, 0x6c,
	0x6c, 0x6f, 0x10, 0xd0, 0x08, 0x12, 0x1f, 0x0a, 0x1a, 0x45, 0x72, 0x72, 0x48, 0x61, 0x6e, 0x64,
	0x73, 0x68, 0x61, 0x6b, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x65, 0x72, 0x48, 0x65,
	0x6c, 0x6c, 0x6f, 0x10, 0xd1, 0x08, 0x12, 0x26, 0x0a, 0x21, 0x45, 0x72, 0x72, 0x48, 0x61, 0x6e,
	0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x72, 0x41,
	0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x10, 0xd2, 0x08, 0x12, 0x20,
	0x0a, 0x1b, 0x45, 0x72, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x64, 0x65, 0x72, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x10, 0xd3, 0x08,
	0x12, 0x25, 0x0a, 0x20, 0x45, 0x72, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x72, 0x41, 0x63, 0x6b, 0x6e, 0x6f, 0x77, 0x6c,
	0x65, 0x64, 0x67, 0x65, 0x10, 0xd4, 0x08, 0x12, 0x21, 0x0a, 0x1c, 0x45, 0x72, 0x72, 0x43, 0x6f,
	0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x61, 0x6d, 0x65,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x10, 0xb0, 0x09, 0x12, 0x29, 0x0a, 0x24, 0x45, 0x72,
	0x72, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x43,
	0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x41, 0x6c, 0x72, 0x65, 0x61, 0x64, 0x79, 0x41, 0x64, 0x64,
	0x65, 0x64, 0x10, 0xb1, 0x09, 0x12, 0x24, 0x0a, 0x1f, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x6e, 0x74,
	0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63,
	0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x10, 0xb2, 0x09, 0x12, 0x26, 0x0a, 0x21, 0x45,
	0x72, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x55, 0x6e, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64,
	0x10, 0xb3, 0x09, 0x12, 0x2d, 0x0a, 0x28, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x6e, 0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67,
	0x41, 0x6c, 0x72, 0x65, 0x61, 0x64, 0x79, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x10,
	0xb4, 0x09, 0x12, 0x1f, 0x0a, 0x1a, 0x45, 0x72, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x4c, 0x6f, 0x67, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4f, 0x70, 0x65, 0x6e,
	0x10, 0x94, 0x0a, 0x12, 0x24, 0x0a, 0x1f, 0x45, 0x72, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x4c, 0x6f, 0x67, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x10, 0x95, 0x0a, 0x12, 0x21, 0x0a, 0x1c, 0x45, 0x72, 0x72,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x55, 0x6e, 0x6b, 0x6e, 0x6f,
	0x77, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x44, 0x10, 0x96, 0x0a, 0x12, 0x22, 0x0a, 0x1d,
	0x45, 0x72, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x4f, 0x74,
	0x68, 0x65, 0x72, 0x44, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x10, 0x97, 0x0a,
	0x12, 0x26, 0x0a, 0x21, 0x45, 0x72, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x53, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x41, 0x6c, 0x72, 0x65, 0x61, 0x64, 0x79, 0x53, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x10, 0x98, 0x0a, 0x12, 0x18, 0x0a, 0x13, 0x45, 0x72, 0x72, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x54, 0x79, 0x70, 0x65, 0x10,
	0x99, 0x0a, 0x12, 0x14, 0x0a, 0x0f, 0x45, 0x72, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4d, 0x69,
	0x73, 0x73, 0x69, 0x6e, 0x67, 0x10, 0x9a, 0x0a, 0x12, 0x15, 0x0a, 0x10, 0x45, 0x72, 0x72, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x41, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x10, 0x9b, 0x0a, 0x12,
	0x17, 0x0a, 0x12, 0x45, 0x72, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x44, 0x65, 0x61, 0x63, 0x74,
	0x69, 0x76, 0x61, 0x74, 0x65, 0x10, 0x9c, 0x0a, 0x12, 0x11, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x10, 0x9d, 0x0a, 0x12, 0x14, 0x0a, 0x0f, 0x45,
	0x72, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x9e,
	0x0a, 0x12, 0x11, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4f, 0x70, 0x65,
	0x6e, 0x10, 0x9f, 0x0a, 0x12, 0x20, 0x0a, 0x1b, 0x45, 0x72, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x4b, 0x65, 0x79, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65,
	0x50, 0x75, 0x74, 0x10, 0xdc, 0x0b, 0x12, 0x20, 0x0a, 0x1b, 0x45, 0x72, 0x72, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x4b, 0x65, 0x79, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e,
	0x63, 0x65, 0x47, 0x65, 0x74, 0x10, 0xdd, 0x0b, 0x12, 0x1a, 0x0a, 0x15, 0x45, 0x72, 0x72, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x10, 0x84, 0x20, 0x12, 0x20, 0x0a, 0x1b, 0x45, 0x72, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x10, 0x85, 0x20, 0x12, 0x29, 0x0a, 0x24, 0x45, 0x72, 0x72, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d,
	0x69, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x10, 0x86,
	0x20, 0x12, 0x19, 0x0a, 0x14, 0x45, 0x72, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x10, 0xe8, 0x20, 0x12, 0x39, 0x0a, 0x34,
	0x45, 0x72, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x44, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x79, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x56, 0x65, 0x72, 0x69, 0x66,
	0x69, 0x65, 0x64, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x53, 0x75, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x10, 0xe9, 0x20, 0x12, 0x2f, 0x0a, 0x2a, 0x45, 0x72, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x45,
	0x78, 0x69, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x4e, 0x6f, 0x74,
	0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10, 0xea, 0x20, 0x12, 0x36, 0x0a, 0x31, 0x45, 0x72, 0x72, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x4c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x41, 0x6e, 0x64, 0x43,
	0x61, 0x6e, 0x74, 0x42, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x64, 0x10, 0xeb, 0x20,
	0x12, 0x34, 0x0a, 0x2f, 0x45, 0x72, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x44,
	0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x45, 0x78, 0x70, 0x6c, 0x69, 0x63, 0x69, 0x74,
	0x52, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x46, 0x6c, 0x61, 0x67, 0x52, 0x65, 0x71, 0x75, 0x69,
	0x72, 0x65, 0x64, 0x10, 0xec, 0x20, 0x12, 0x32, 0x0a, 0x2d, 0x45, 0x72, 0x72, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x49, 0x6e,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x43, 0x72, 0x65,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x10, 0xed, 0x20, 0x12, 0x32, 0x0a, 0x2d, 0x45, 0x72,
	0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x79, 0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65,
	0x64, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x10, 0xee, 0x20, 0x12, 0x34,
	0x0a, 0x2f, 0x45, 0x72, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x44, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x56, 0x65, 0x72,
	0x69, 0x66, 0x69, 0x65, 0x64, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x49,
	0x44, 0x10, 0xef, 0x20, 0x42, 0x20, 0x5a, 0x1e, 0x62, 0x65, 0x72, 0x74, 0x79, 0x2e, 0x74, 0x65,
	0x63, 0x68, 0x2f, 0x77, 0x65, 0x73, 0x68, 0x6e, 0x65, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x65,
	0x72, 0x72, 0x63, 0x6f, 0x64, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_errcode_errcode_proto_rawDescOnce sync.Once
	file_errcode_errcode_proto_rawDescData = file_errcode_errcode_proto_rawDesc
)

func file_errcode_errcode_proto_rawDescGZIP() []byte {
	file_errcode_errcode_proto_rawDescOnce.Do(func() {
		file_errcode_errcode_proto_rawDescData = protoimpl.X.CompressGZIP(file_errcode_errcode_proto_rawDescData)
	})
	return file_errcode_errcode_proto_rawDescData
}

var file_errcode_errcode_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_errcode_errcode_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_errcode_errcode_proto_goTypes = []interface{}{
	(ErrCode)(0),       // 0: weshnet.errcode.ErrCode
	(*ErrDetails)(nil), // 1: weshnet.errcode.ErrDetails
}
var file_errcode_errcode_proto_depIdxs = []int32{
	0, // 0: weshnet.errcode.ErrDetails.codes:type_name -> weshnet.errcode.ErrCode
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_errcode_errcode_proto_init() }
func file_errcode_errcode_proto_init() {
	if File_errcode_errcode_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_errcode_errcode_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ErrDetails); i {
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
			RawDescriptor: file_errcode_errcode_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_errcode_errcode_proto_goTypes,
		DependencyIndexes: file_errcode_errcode_proto_depIdxs,
		EnumInfos:         file_errcode_errcode_proto_enumTypes,
		MessageInfos:      file_errcode_errcode_proto_msgTypes,
	}.Build()
	File_errcode_errcode_proto = out.File
	file_errcode_errcode_proto_rawDesc = nil
	file_errcode_errcode_proto_goTypes = nil
	file_errcode_errcode_proto_depIdxs = nil
}
