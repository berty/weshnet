// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: errcode/errcode.proto

package errcode

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
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
	Undefined                                            ErrCode = 0
	TODO                                                 ErrCode = 666
	ErrNotImplemented                                    ErrCode = 777
	ErrInternal                                          ErrCode = 888
	ErrInvalidInput                                      ErrCode = 100
	ErrInvalidRange                                      ErrCode = 101
	ErrMissingInput                                      ErrCode = 102
	ErrSerialization                                     ErrCode = 103
	ErrDeserialization                                   ErrCode = 104
	ErrStreamRead                                        ErrCode = 105
	ErrStreamWrite                                       ErrCode = 106
	ErrStreamTransform                                   ErrCode = 110
	ErrStreamSendAndClose                                ErrCode = 111
	ErrStreamHeaderWrite                                 ErrCode = 112
	ErrStreamHeaderRead                                  ErrCode = 115
	ErrStreamSink                                        ErrCode = 113
	ErrStreamCloseAndRecv                                ErrCode = 114
	ErrMissingMapKey                                     ErrCode = 107
	ErrDBWrite                                           ErrCode = 108
	ErrDBRead                                            ErrCode = 109
	ErrDBDestroy                                         ErrCode = 120
	ErrDBMigrate                                         ErrCode = 121
	ErrDBReplay                                          ErrCode = 122
	ErrDBRestore                                         ErrCode = 123
	ErrDBOpen                                            ErrCode = 124
	ErrDBClose                                           ErrCode = 125
	ErrCryptoRandomGeneration                            ErrCode = 200
	ErrCryptoKeyGeneration                               ErrCode = 201
	ErrCryptoNonceGeneration                             ErrCode = 202
	ErrCryptoSignature                                   ErrCode = 203
	ErrCryptoSignatureVerification                       ErrCode = 204
	ErrCryptoDecrypt                                     ErrCode = 205
	ErrCryptoDecryptPayload                              ErrCode = 206
	ErrCryptoEncrypt                                     ErrCode = 207
	ErrCryptoKeyConversion                               ErrCode = 208
	ErrCryptoCipherInit                                  ErrCode = 209
	ErrCryptoKeyDerivation                               ErrCode = 210
	ErrMap                                               ErrCode = 300
	ErrForEach                                           ErrCode = 301
	ErrKeystoreGet                                       ErrCode = 400
	ErrKeystorePut                                       ErrCode = 401
	ErrNotFound                                          ErrCode = 404
	ErrOrbitDBInit                                       ErrCode = 1000
	ErrOrbitDBOpen                                       ErrCode = 1001
	ErrOrbitDBAppend                                     ErrCode = 1002
	ErrOrbitDBDeserialization                            ErrCode = 1003
	ErrOrbitDBStoreCast                                  ErrCode = 1004
	ErrHandshakeOwnEphemeralKeyGenSend                   ErrCode = 1100
	ErrHandshakePeerEphemeralKeyRecv                     ErrCode = 1101
	ErrHandshakeRequesterAuthenticateBoxKeyGen           ErrCode = 1102
	ErrHandshakeResponderAcceptBoxKeyGen                 ErrCode = 1103
	ErrHandshakeRequesterHello                           ErrCode = 1104
	ErrHandshakeResponderHello                           ErrCode = 1105
	ErrHandshakeRequesterAuthenticate                    ErrCode = 1106
	ErrHandshakeResponderAccept                          ErrCode = 1107
	ErrHandshakeRequesterAcknowledge                     ErrCode = 1108
	ErrContactRequestSameAccount                         ErrCode = 1200
	ErrContactRequestContactAlreadyAdded                 ErrCode = 1201
	ErrContactRequestContactBlocked                      ErrCode = 1202
	ErrContactRequestContactUndefined                    ErrCode = 1203
	ErrContactRequestIncomingAlreadyReceived             ErrCode = 1204
	ErrGroupMemberLogEventOpen                           ErrCode = 1300
	ErrGroupMemberLogEventSignature                      ErrCode = 1301
	ErrGroupMemberUnknownGroupID                         ErrCode = 1302
	ErrGroupSecretOtherDestMember                        ErrCode = 1303
	ErrGroupSecretAlreadySentToMember                    ErrCode = 1304
	ErrGroupInvalidType                                  ErrCode = 1305
	ErrGroupMissing                                      ErrCode = 1306
	ErrGroupActivate                                     ErrCode = 1307
	ErrGroupDeactivate                                   ErrCode = 1308
	ErrGroupInfo                                         ErrCode = 1309
	ErrGroupUnknown                                      ErrCode = 1310
	ErrGroupOpen                                         ErrCode = 1311
	ErrMessageKeyPersistencePut                          ErrCode = 1500
	ErrMessageKeyPersistenceGet                          ErrCode = 1501
	ErrServicesAuth                                      ErrCode = 4000
	ErrServicesAuthNotInitialized                        ErrCode = 4001
	ErrServicesAuthWrongState                            ErrCode = 4002
	ErrServicesAuthInvalidResponse                       ErrCode = 4003
	ErrServicesAuthServer                                ErrCode = 4004
	ErrServicesAuthCodeChallenge                         ErrCode = 4005
	ErrServicesAuthServiceInvalidToken                   ErrCode = 4006
	ErrServicesAuthServiceNotSupported                   ErrCode = 4007
	ErrServicesAuthUnknownToken                          ErrCode = 4008
	ErrServicesAuthInvalidURL                            ErrCode = 4009
	ErrServiceReplication                                ErrCode = 4100
	ErrServiceReplicationServer                          ErrCode = 4101
	ErrServiceReplicationMissingEndpoint                 ErrCode = 4102
	ErrServicesDirectory                                 ErrCode = 4200
	ErrServicesDirectoryInvalidVerifiedCredentialSubject ErrCode = 4201
	ErrServicesDirectoryExistingRecordNotFound           ErrCode = 4202
	ErrServicesDirectoryRecordLockedAndCantBeReplaced    ErrCode = 4203
	ErrServicesDirectoryExplicitReplaceFlagRequired      ErrCode = 4204
	ErrServicesDirectoryInvalidVerifiedCredential        ErrCode = 4205
	ErrServicesDirectoryExpiredVerifiedCredential        ErrCode = 4206
	ErrServicesDirectoryInvalidVerifiedCredentialID      ErrCode = 4207
	ErrPush                                              ErrCode = 6000
	ErrPushWrongAccount                                  ErrCode = 6001
	ErrPushUnableToDecrypt                               ErrCode = 6002
	ErrPushInvalidPayload                                ErrCode = 6003
	ErrPushInvalidServerConfig                           ErrCode = 6004
	ErrPushMissingBundleID                               ErrCode = 6005
	ErrPushUnknownDestination                            ErrCode = 6006
	ErrPushProvider                                      ErrCode = 6007
	ErrPushUnknownProvider                               ErrCode = 6008
	ErrNoProvidersConfigured                             ErrCode = 6009
	ErrInvalidPrivateKey                                 ErrCode = 6010
)

var ErrCode_name = map[int32]string{
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
	4000: "ErrServicesAuth",
	4001: "ErrServicesAuthNotInitialized",
	4002: "ErrServicesAuthWrongState",
	4003: "ErrServicesAuthInvalidResponse",
	4004: "ErrServicesAuthServer",
	4005: "ErrServicesAuthCodeChallenge",
	4006: "ErrServicesAuthServiceInvalidToken",
	4007: "ErrServicesAuthServiceNotSupported",
	4008: "ErrServicesAuthUnknownToken",
	4009: "ErrServicesAuthInvalidURL",
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
	6000: "ErrPush",
	6001: "ErrPushWrongAccount",
	6002: "ErrPushUnableToDecrypt",
	6003: "ErrPushInvalidPayload",
	6004: "ErrPushInvalidServerConfig",
	6005: "ErrPushMissingBundleID",
	6006: "ErrPushUnknownDestination",
	6007: "ErrPushProvider",
	6008: "ErrPushUnknownProvider",
	6009: "ErrNoProvidersConfigured",
	6010: "ErrInvalidPrivateKey",
}

var ErrCode_value = map[string]int32{
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
	"ErrServicesAuth":                                      4000,
	"ErrServicesAuthNotInitialized":                        4001,
	"ErrServicesAuthWrongState":                            4002,
	"ErrServicesAuthInvalidResponse":                       4003,
	"ErrServicesAuthServer":                                4004,
	"ErrServicesAuthCodeChallenge":                         4005,
	"ErrServicesAuthServiceInvalidToken":                   4006,
	"ErrServicesAuthServiceNotSupported":                   4007,
	"ErrServicesAuthUnknownToken":                          4008,
	"ErrServicesAuthInvalidURL":                            4009,
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
	"ErrPush":                    6000,
	"ErrPushWrongAccount":        6001,
	"ErrPushUnableToDecrypt":     6002,
	"ErrPushInvalidPayload":      6003,
	"ErrPushInvalidServerConfig": 6004,
	"ErrPushMissingBundleID":     6005,
	"ErrPushUnknownDestination":  6006,
	"ErrPushProvider":            6007,
	"ErrPushUnknownProvider":     6008,
	"ErrNoProvidersConfigured":   6009,
	"ErrInvalidPrivateKey":       6010,
}

func (x ErrCode) String() string {
	return proto.EnumName(ErrCode_name, int32(x))
}

func (ErrCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_fb5abb189af31c1a, []int{0}
}

type ErrDetails struct {
	Codes                []ErrCode `protobuf:"varint,1,rep,packed,name=codes,proto3,enum=weshnet.errcode.ErrCode" json:"codes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ErrDetails) Reset()         { *m = ErrDetails{} }
func (m *ErrDetails) String() string { return proto.CompactTextString(m) }
func (*ErrDetails) ProtoMessage()    {}
func (*ErrDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb5abb189af31c1a, []int{0}
}
func (m *ErrDetails) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ErrDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ErrDetails.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ErrDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrDetails.Merge(m, src)
}
func (m *ErrDetails) XXX_Size() int {
	return m.Size()
}
func (m *ErrDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrDetails.DiscardUnknown(m)
}

var xxx_messageInfo_ErrDetails proto.InternalMessageInfo

func (m *ErrDetails) GetCodes() []ErrCode {
	if m != nil {
		return m.Codes
	}
	return nil
}

func init() {
	proto.RegisterEnum("weshnet.errcode.ErrCode", ErrCode_name, ErrCode_value)
	proto.RegisterType((*ErrDetails)(nil), "weshnet.errcode.ErrDetails")
}

func init() { proto.RegisterFile("errcode/errcode.proto", fileDescriptor_fb5abb189af31c1a) }

var fileDescriptor_fb5abb189af31c1a = []byte{
	// 1552 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x57, 0xdb, 0x6f, 0x9c, 0x47,
	0x15, 0xef, 0x36, 0xd4, 0x97, 0x71, 0x1b, 0x9f, 0x4c, 0x9c, 0xc4, 0x49, 0x5b, 0xef, 0x36, 0x14,
	0x28, 0x95, 0xba, 0x86, 0x52, 0x21, 0x21, 0xf1, 0xb2, 0xde, 0xdd, 0x24, 0x2b, 0xc7, 0x17, 0xed,
	0xda, 0x54, 0xe2, 0x6d, 0xfc, 0xcd, 0xf1, 0xb7, 0x53, 0x7f, 0x3b, 0xf3, 0x75, 0xbe, 0xd9, 0x8d,
	0xb7, 0xc0, 0x1b, 0x20, 0xf1, 0x06, 0x52, 0xb9, 0x95, 0xfb, 0x1d, 0x24, 0x90, 0xb8, 0xfd, 0x11,
	0x05, 0xda, 0xc4, 0x09, 0xbc, 0x20, 0x81, 0x04, 0x79, 0xca, 0x8d, 0xfb, 0x2d, 0xf0, 0x84, 0xe6,
	0xb2, 0xeb, 0x6f, 0xed, 0x75, 0x10, 0x4f, 0xf6, 0x77, 0xce, 0xef, 0xdc, 0xe6, 0x9c, 0xf9, 0xcd,
	0x59, 0x72, 0x0a, 0xb5, 0x8e, 0x14, 0xc7, 0xc5, 0xf0, 0xb7, 0x9c, 0x6a, 0x65, 0x14, 0x9d, 0xbd,
	0x82, 0x59, 0x5b, 0xa2, 0x29, 0x07, 0xf1, 0xb9, 0xb9, 0x58, 0xc5, 0xca, 0xe9, 0x16, 0xed, 0x7f,
	0x1e, 0x76, 0xfe, 0xfd, 0x84, 0xd4, 0xb5, 0xae, 0xa1, 0x61, 0x22, 0xc9, 0x68, 0x99, 0x3c, 0x62,
	0xb1, 0xd9, 0x7c, 0xa1, 0x74, 0xec, 0x99, 0xe3, 0xcf, 0xcf, 0x97, 0x0f, 0x38, 0x29, 0xd7, 0xb5,
	0xae, 0x2a, 0x8e, 0x4d, 0x0f, 0x7b, 0xf6, 0xd7, 0xf3, 0x64, 0x32, 0x88, 0xe8, 0x63, 0x64, 0x7a,
	0x53, 0x72, 0xdc, 0x16, 0x12, 0x39, 0x3c, 0x44, 0xa7, 0xc9, 0x5b, 0x36, 0xd6, 0x6a, 0x6b, 0xf0,
	0xda, 0x23, 0xf4, 0x34, 0x39, 0x51, 0xd7, 0x7a, 0x55, 0x99, 0x46, 0x27, 0x4d, 0xb0, 0x83, 0xd2,
	0x20, 0x87, 0x4f, 0x4c, 0x50, 0x20, 0x33, 0x75, 0xad, 0x1b, 0xd2, 0xa0, 0x96, 0x2c, 0x81, 0xfb,
	0x13, 0xf4, 0x24, 0x99, 0x75, 0x92, 0x1e, 0x4b, 0x04, 0x6f, 0xc8, 0xb4, 0x6b, 0x80, 0x8f, 0x0a,
	0x9b, 0x4c, 0xc6, 0x08, 0x18, 0x84, 0x2b, 0x22, 0xcb, 0x84, 0x8c, 0x3d, 0x72, 0x9b, 0xce, 0x11,
	0xa8, 0x6b, 0xdd, 0x42, 0x2d, 0x58, 0x22, 0x5e, 0x61, 0x46, 0x28, 0x09, 0x31, 0x3d, 0x4d, 0xa8,
	0x2b, 0x31, 0x1b, 0x91, 0xb7, 0xe9, 0x09, 0xf2, 0x98, 0x45, 0x1b, 0x8d, 0xac, 0xd3, 0x44, 0xc6,
	0x41, 0x50, 0x4a, 0x8e, 0x0f, 0x45, 0x2f, 0x6a, 0x61, 0x10, 0x5e, 0x0a, 0xe6, 0x5e, 0xb6, 0xa1,
	0x99, 0xcc, 0xb6, 0x95, 0xee, 0x80, 0xa4, 0x67, 0xc9, 0xa9, 0xa1, 0xbc, 0x85, 0x92, 0x57, 0x24,
	0xaf, 0x26, 0x2a, 0x43, 0x50, 0x74, 0x9e, 0xcc, 0x0d, 0x55, 0x97, 0x90, 0x71, 0xd4, 0xde, 0x59,
	0x4a, 0xcf, 0x90, 0x93, 0x07, 0x34, 0x2e, 0x72, 0x36, 0x92, 0x4c, 0x4b, 0xc8, 0x1d, 0x78, 0x79,
	0x24, 0x80, 0xf3, 0x5c, 0x91, 0xbc, 0x89, 0x51, 0x0f, 0x74, 0x28, 0x34, 0x54, 0xbf, 0xc2, 0xd2,
	0x65, 0xec, 0xc3, 0x0e, 0x3d, 0xee, 0x7b, 0xb9, 0xe4, 0x83, 0x25, 0xb6, 0x23, 0xee, 0xdb, 0x85,
	0xe8, 0x50, 0x20, 0x8f, 0xba, 0xcf, 0x1a, 0x66, 0x46, 0xab, 0x3e, 0xec, 0x0e, 0x25, 0x2b, 0x22,
	0xd6, 0xcc, 0x20, 0xf4, 0xe9, 0xac, 0x6b, 0x89, 0x35, 0x49, 0x13, 0xd6, 0x87, 0x57, 0x86, 0x90,
	0x26, 0x66, 0x46, 0x69, 0x84, 0x0f, 0x0d, 0xbd, 0xae, 0xa5, 0x28, 0xe1, 0xc3, 0xc3, 0xa0, 0xbe,
	0xf6, 0x8f, 0xd0, 0x05, 0x72, 0xd6, 0x4e, 0x84, 0xee, 0xa7, 0x46, 0x35, 0x99, 0xe4, 0xaa, 0x73,
	0x11, 0x25, 0x6a, 0x7f, 0xe8, 0xaf, 0x17, 0xe8, 0xe3, 0xe4, 0xf4, 0x50, 0xbf, 0x8c, 0xfd, 0x9c,
	0xf2, 0x67, 0x05, 0xfa, 0x24, 0x99, 0x1f, 0x2a, 0x57, 0x95, 0x8c, 0x30, 0xa7, 0xfe, 0x79, 0x81,
	0x9e, 0x71, 0xad, 0xf0, 0xea, 0x96, 0x88, 0x25, 0x33, 0x5d, 0x8d, 0xf0, 0x8b, 0x02, 0x7d, 0x2b,
	0x59, 0x38, 0xac, 0xf8, 0x00, 0x6a, 0xb1, 0x2d, 0x22, 0x6f, 0xfd, 0x46, 0x81, 0x9e, 0x72, 0x87,
	0xe6, 0x41, 0x35, 0x8c, 0xec, 0x5f, 0x78, 0xb3, 0x40, 0x9f, 0x20, 0x67, 0x0e, 0x8a, 0xd7, 0x59,
	0x3f, 0x51, 0x8c, 0xc3, 0xd5, 0x51, 0xa3, 0xba, 0xf4, 0x46, 0xd7, 0x0e, 0x55, 0x51, 0x55, 0xb2,
	0x87, 0x3a, 0xb3, 0x81, 0xf6, 0x0a, 0x74, 0xde, 0x35, 0xd9, 0x2b, 0xab, 0x22, 0x6d, 0xa3, 0x6e,
	0x48, 0x61, 0xe0, 0xfa, 0x21, 0xb3, 0x1a, 0x6a, 0xd1, 0xf3, 0xf9, 0xdd, 0x28, 0xd0, 0x19, 0x32,
	0x61, 0x9b, 0xca, 0x52, 0xf8, 0xfe, 0xc3, 0x74, 0xd6, 0x1d, 0xeb, 0x05, 0xa5, 0xeb, 0x2c, 0x6a,
	0xc3, 0x0f, 0x1e, 0xa6, 0x27, 0xdd, 0x68, 0x2e, 0x63, 0xdf, 0xf5, 0xe1, 0x22, 0x1a, 0xf8, 0xe4,
	0xb1, 0x03, 0xc2, 0xf5, 0xae, 0x81, 0x4f, 0x1d, 0x0b, 0xd7, 0x6a, 0x55, 0x99, 0x0b, 0xaa, 0x2b,
	0x39, 0xbc, 0x3a, 0x80, 0xad, 0xe9, 0x2d, 0x61, 0x6a, 0x4b, 0x2e, 0x97, 0x5b, 0x93, 0xa3, 0x42,
	0xd7, 0xcc, 0xdb, 0x93, 0xa1, 0xdc, 0x20, 0xac, 0xa4, 0x29, 0x4a, 0x0e, 0x77, 0x26, 0x43, 0x53,
	0x83, 0xf8, 0xe0, 0x4d, 0xba, 0x3b, 0x19, 0x2a, 0x0e, 0xfa, 0x96, 0xcd, 0xa5, 0xca, 0x32, 0x03,
	0xf7, 0x26, 0xe9, 0x3b, 0xc8, 0xf9, 0xba, 0xd6, 0x97, 0x98, 0xe4, 0x59, 0x9b, 0xed, 0xe0, 0xda,
	0x15, 0x59, 0x4f, 0xdb, 0xd8, 0x41, 0xcd, 0x12, 0xdf, 0x7d, 0x7b, 0x75, 0xe0, 0x8d, 0x29, 0xfa,
	0x36, 0x52, 0xca, 0x03, 0xd7, 0x11, 0x75, 0x1e, 0xe9, 0x06, 0xff, 0xcd, 0x29, 0xba, 0x48, 0x9e,
	0xcd, 0xc3, 0x9a, 0xf8, 0x72, 0x17, 0x33, 0x83, 0xba, 0xd2, 0x35, 0x6d, 0x94, 0xc6, 0xb6, 0x1b,
	0x97, 0xd4, 0xae, 0xf7, 0x0d, 0x57, 0xa7, 0xe8, 0x3b, 0xc9, 0xd3, 0xa3, 0x06, 0x59, 0xaa, 0x24,
	0x47, 0x5d, 0x89, 0x22, 0x4c, 0xcd, 0x3e, 0xf4, 0xda, 0x14, 0x2d, 0x92, 0x73, 0x63, 0x7d, 0x5f,
	0xc2, 0x24, 0x51, 0xb0, 0x37, 0x06, 0x10, 0x7c, 0x79, 0xc0, 0xf5, 0x29, 0xfa, 0x76, 0xf2, 0xd4,
	0xff, 0xcc, 0x0e, 0x6e, 0x4c, 0xd1, 0x12, 0x79, 0xfc, 0x01, 0x49, 0xc1, 0x2f, 0x0f, 0x1d, 0xc7,
	0xbe, 0xa7, 0x68, 0x47, 0xaa, 0x2b, 0x09, 0xf2, 0x18, 0xe1, 0x57, 0x53, 0xf4, 0x29, 0xf2, 0x84,
	0xe3, 0x5f, 0x69, 0x58, 0x64, 0x02, 0xa8, 0xc5, 0x3a, 0x58, 0x89, 0x22, 0xd5, 0x95, 0x06, 0x7e,
	0x38, 0x1d, 0x0e, 0x60, 0x14, 0x12, 0xbe, 0x2a, 0x89, 0x46, 0xc6, 0xfb, 0x15, 0xce, 0x91, 0xc3,
	0x8f, 0xa6, 0xe9, 0xd3, 0xa4, 0x78, 0x14, 0x74, 0x29, 0x51, 0xd1, 0x0e, 0x72, 0xf8, 0xf1, 0x74,
	0x28, 0x72, 0x2c, 0x6a, 0xff, 0x01, 0xf8, 0xc9, 0x34, 0x7d, 0x8e, 0x3c, 0x73, 0x08, 0xd7, 0x90,
	0x91, 0xea, 0x08, 0x19, 0x87, 0xc8, 0x4d, 0x8c, 0x50, 0xf4, 0x90, 0xc3, 0x4f, 0xa7, 0xc3, 0xe1,
	0x5e, 0xd4, 0xaa, 0x9b, 0xae, 0x60, 0x67, 0x0b, 0xf5, 0x65, 0x15, 0xd7, 0x7b, 0x28, 0x8d, 0x9b,
	0xcd, 0x57, 0x49, 0xc8, 0x6e, 0x0c, 0x60, 0x9f, 0x0a, 0x3e, 0x4d, 0xc2, 0x89, 0xe4, 0x50, 0x9b,
	0xd2, 0x9e, 0x98, 0x74, 0x92, 0x46, 0x0d, 0x3e, 0x43, 0xe8, 0x79, 0xf2, 0xe4, 0x00, 0xd2, 0xc2,
	0x48, 0xa3, 0x59, 0x33, 0x6d, 0xb4, 0x0f, 0x84, 0xf1, 0x16, 0xf0, 0x59, 0x12, 0x8a, 0xcc, 0x61,
	0x42, 0xc6, 0x2d, 0x94, 0x66, 0x43, 0x05, 0xdc, 0xe7, 0x48, 0x98, 0x7c, 0xef, 0xdc, 0xbf, 0x50,
	0x1b, 0xfd, 0x14, 0xe1, 0xf3, 0x84, 0xce, 0xb9, 0x17, 0xca, 0x27, 0xe2, 0x89, 0x1a, 0x5e, 0x23,
	0xe1, 0x82, 0x39, 0x69, 0x25, 0x32, 0xf6, 0xf6, 0x23, 0x7c, 0x81, 0x04, 0x66, 0x73, 0xe2, 0x1a,
	0xb2, 0x81, 0xe2, 0x8b, 0x84, 0x9e, 0x70, 0xfc, 0x1b, 0xfc, 0x6f, 0x2b, 0xf8, 0xd2, 0x88, 0xe3,
	0x50, 0x1b, 0x7c, 0x79, 0x04, 0xe8, 0x0e, 0xec, 0x2b, 0x24, 0x4c, 0xd9, 0x0a, 0x66, 0x19, 0x8b,
	0x71, 0x19, 0xfb, 0xeb, 0x96, 0xa2, 0x32, 0x83, 0x32, 0x72, 0x54, 0xf1, 0x9b, 0x99, 0x07, 0x21,
	0x2c, 0xc3, 0xfc, 0x76, 0x26, 0x04, 0x6b, 0xa1, 0xee, 0x89, 0x08, 0x33, 0x3b, 0xc7, 0xf0, 0xd5,
	0x62, 0x38, 0xc1, 0xbc, 0xd4, 0xbe, 0xee, 0x52, 0x18, 0x47, 0x0b, 0xc8, 0xe1, 0x6b, 0xc5, 0xc0,
	0x19, 0x79, 0xcc, 0x8b, 0x5a, 0xc9, 0xb8, 0x65, 0x6c, 0x65, 0x5f, 0x2f, 0x06, 0xce, 0xce, 0xeb,
	0x07, 0x4f, 0xbc, 0xbb, 0x0d, 0x19, 0xc2, 0x37, 0x8a, 0xf4, 0x9c, 0x7f, 0x03, 0x73, 0x20, 0xfb,
	0x3f, 0x6a, 0xf8, 0x66, 0x31, 0x74, 0x3a, 0xaf, 0xb3, 0x7b, 0x48, 0xb5, 0xcd, 0x92, 0x04, 0xed,
	0x92, 0xf0, 0xad, 0x62, 0x60, 0x9f, 0x83, 0xe6, 0x22, 0xc2, 0x41, 0xaf, 0xd4, 0x0e, 0x4a, 0xf8,
	0xf6, 0x03, 0x80, 0xab, 0xca, 0xb4, 0xba, 0x69, 0xaa, 0xb4, 0xdd, 0x59, 0xbe, 0x53, 0x0c, 0x27,
	0x96, 0x07, 0x86, 0x1e, 0x78, 0x57, 0xdf, 0x1d, 0x57, 0x77, 0x08, 0xb6, 0xd9, 0xbc, 0x0c, 0xdf,
	0x3b, 0x50, 0x92, 0x7d, 0x67, 0x07, 0x4f, 0xd4, 0x47, 0x4b, 0xa3, 0xde, 0x73, 0xba, 0x50, 0xf4,
	0xc7, 0x4a, 0xe1, 0x36, 0x1f, 0x46, 0x84, 0x11, 0xab, 0x4b, 0x9e, 0x2a, 0x21, 0x0d, 0x7c, 0xbc,
	0x44, 0xcf, 0xfa, 0x2d, 0x24, 0x24, 0x52, 0x13, 0x1a, 0x23, 0xa3, 0x74, 0x1f, 0x6e, 0x95, 0xe8,
	0xfb, 0xc8, 0x0b, 0xe3, 0x54, 0x21, 0x51, 0xff, 0x70, 0x22, 0xaf, 0x6a, 0xe4, 0x96, 0xb5, 0x58,
	0xd2, 0xea, 0x6e, 0xbd, 0x84, 0x91, 0x81, 0xdb, 0xa5, 0x40, 0xc0, 0x87, 0x4c, 0xeb, 0xbb, 0x22,
	0x33, 0x42, 0xc6, 0x4d, 0x8c, 0x94, 0xe6, 0xc3, 0xc7, 0xe7, 0x4e, 0x89, 0xbe, 0x97, 0xbc, 0x7b,
	0x9c, 0x81, 0x07, 0x5e, 0x76, 0xac, 0x62, 0x57, 0x27, 0x26, 0xcd, 0x92, 0x2b, 0x88, 0x45, 0xc8,
	0xe1, 0x6e, 0x89, 0xbe, 0x40, 0x16, 0xc7, 0x07, 0xb2, 0x35, 0x0b, 0x13, 0xa0, 0x17, 0x12, 0x16,
	0x5b, 0x6e, 0x11, 0x1a, 0x39, 0xdc, 0x2b, 0xd1, 0xe7, 0xc9, 0x73, 0xff, 0x57, 0x65, 0xf0, 0x87,
	0x23, 0x6d, 0xea, 0xbb, 0xa9, 0xf5, 0x3a, 0xc6, 0xe6, 0x8f, 0x47, 0x66, 0x77, 0x64, 0x9c, 0x46,
	0x0d, 0xfe, 0x54, 0xa2, 0x8f, 0xba, 0x75, 0x79, 0xbd, 0x9b, 0xb5, 0xe1, 0xcf, 0xe5, 0xc0, 0x1d,
	0xf6, 0xcb, 0xdd, 0x8c, 0x01, 0x67, 0xff, 0xa5, 0x1c, 0xf6, 0x04, 0xab, 0xd9, 0x94, 0x6c, 0x2b,
	0xc1, 0x8d, 0xe1, 0xc2, 0xf2, 0xd7, 0x72, 0x18, 0x20, 0xab, 0x0c, 0xd1, 0x06, 0xeb, 0xca, 0xdf,
	0xca, 0x81, 0x44, 0x73, 0x3a, 0x3f, 0x39, 0x55, 0x25, 0xb7, 0x45, 0x0c, 0x7f, 0xcf, 0x7b, 0x0e,
	0x13, 0xb3, 0xd4, 0x95, 0x3c, 0xc1, 0x46, 0x0d, 0xfe, 0x51, 0x0e, 0xa3, 0xeb, 0xc3, 0xba, 0xa1,
	0xb6, 0x9c, 0x28, 0xa4, 0x1f, 0xcf, 0x7f, 0x96, 0x03, 0x19, 0x58, 0xfd, 0xba, 0x56, 0x3d, 0xc1,
	0x51, 0xc3, 0xbf, 0x46, 0x93, 0x75, 0x56, 0x43, 0xe5, 0xfd, 0x72, 0xd8, 0xe8, 0x56, 0xd5, 0x40,
	0x96, 0xf9, 0x54, 0xba, 0xb6, 0x5d, 0xff, 0x2e, 0x87, 0x19, 0x1d, 0xd4, 0xe1, 0xb6, 0x21, 0xcb,
	0x43, 0xf0, 0x9f, 0xf2, 0xd2, 0xbb, 0xf6, 0x7e, 0xbf, 0xf0, 0xd0, 0xeb, 0x37, 0x17, 0x0a, 0x7b,
	0x37, 0x17, 0x0a, 0xbf, 0xbb, 0xb9, 0x50, 0xf8, 0xe0, 0xc2, 0x16, 0x6a, 0xd3, 0x2f, 0x1b, 0x8c,
	0xda, 0x8b, 0xe1, 0x87, 0xc9, 0x62, 0xba, 0x13, 0x0f, 0x7e, 0xf8, 0x6c, 0x4d, 0xb8, 0x9f, 0x34,
	0xef, 0xf9, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x32, 0x17, 0x7f, 0x89, 0x12, 0x0d, 0x00, 0x00,
}

func (m *ErrDetails) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ErrDetails) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ErrDetails) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Codes) > 0 {
		dAtA2 := make([]byte, len(m.Codes)*10)
		var j1 int
		for _, num := range m.Codes {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		i -= j1
		copy(dAtA[i:], dAtA2[:j1])
		i = encodeVarintErrcode(dAtA, i, uint64(j1))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintErrcode(dAtA []byte, offset int, v uint64) int {
	offset -= sovErrcode(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ErrDetails) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Codes) > 0 {
		l = 0
		for _, e := range m.Codes {
			l += sovErrcode(uint64(e))
		}
		n += 1 + sovErrcode(uint64(l)) + l
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovErrcode(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozErrcode(x uint64) (n int) {
	return sovErrcode(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ErrDetails) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowErrcode
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ErrDetails: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ErrDetails: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v ErrCode
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowErrcode
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= ErrCode(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Codes = append(m.Codes, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowErrcode
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthErrcode
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthErrcode
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				if elementCount != 0 && len(m.Codes) == 0 {
					m.Codes = make([]ErrCode, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v ErrCode
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowErrcode
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= ErrCode(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Codes = append(m.Codes, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Codes", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipErrcode(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthErrcode
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipErrcode(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowErrcode
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowErrcode
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowErrcode
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthErrcode
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupErrcode
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthErrcode
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthErrcode        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowErrcode          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupErrcode = fmt.Errorf("proto: unexpected end of group")
)
