// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: entity.proto

package bertyprotocol

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	crypto "berty.tech/go/internal/crypto"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
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

type Contact_ContactStatus int32

const (
	Contact_Unknown    Contact_ContactStatus = 0
	Contact_Myself     Contact_ContactStatus = 1
	Contact_Contact    Contact_ContactStatus = 2
	Contact_Trusted    Contact_ContactStatus = 3
	Contact_Blocked    Contact_ContactStatus = 4
	Contact_RequestIn  Contact_ContactStatus = 5
	Contact_RequestOut Contact_ContactStatus = 6
)

var Contact_ContactStatus_name = map[int32]string{
	0: "Unknown",
	1: "Myself",
	2: "Contact",
	3: "Trusted",
	4: "Blocked",
	5: "RequestIn",
	6: "RequestOut",
}

var Contact_ContactStatus_value = map[string]int32{
	"Unknown":    0,
	"Myself":     1,
	"Contact":    2,
	"Trusted":    3,
	"Blocked":    4,
	"RequestIn":  5,
	"RequestOut": 6,
}

func (x Contact_ContactStatus) String() string {
	return proto.EnumName(Contact_ContactStatus_name, int32(x))
}

func (Contact_ContactStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{1, 0}
}

type Device struct {
	DevicePubKey  []byte `protobuf:"bytes,1,opt,name=device_pub_key,json=devicePubKey,proto3" json:"device_pub_key,omitempty"`
	AccountPubKey []byte `protobuf:"bytes,2,opt,name=account_pub_key,json=accountPubKey,proto3" json:"account_pub_key,omitempty"`
}

func (m *Device) Reset()         { *m = Device{} }
func (m *Device) String() string { return proto.CompactTextString(m) }
func (*Device) ProtoMessage()    {}
func (*Device) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{0}
}
func (m *Device) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Device) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Device.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Device) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Device.Merge(m, src)
}
func (m *Device) XXX_Size() int {
	return m.Size()
}
func (m *Device) XXX_DiscardUnknown() {
	xxx_messageInfo_Device.DiscardUnknown(m)
}

var xxx_messageInfo_Device proto.InternalMessageInfo

func (m *Device) GetDevicePubKey() []byte {
	if m != nil {
		return m.DevicePubKey
	}
	return nil
}

func (m *Device) GetAccountPubKey() []byte {
	if m != nil {
		return m.AccountPubKey
	}
	return nil
}

type Contact struct {
	// AccountID = sig_chain.entries[0].pub_key
	SigChain            *crypto.SigChain      `protobuf:"bytes,1,opt,name=sig_chain,json=sigChain,proto3" json:"sig_chain,omitempty"`
	Metadata            []byte                `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
	RendezvousPointSeed []byte                `protobuf:"bytes,3,opt,name=rendezvous_point_seed,json=rendezvousPointSeed,proto3" json:"rendezvous_point_seed,omitempty"`
	SharedSecret        []byte                `protobuf:"bytes,4,opt,name=shared_secret,json=sharedSecret,proto3" json:"shared_secret,omitempty"`
	ContactStatus       Contact_ContactStatus `protobuf:"varint,5,opt,name=contact_status,json=contactStatus,proto3,enum=Contact_ContactStatus" json:"contact_status,omitempty"`
}

func (m *Contact) Reset()         { *m = Contact{} }
func (m *Contact) String() string { return proto.CompactTextString(m) }
func (*Contact) ProtoMessage()    {}
func (*Contact) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{1}
}
func (m *Contact) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Contact) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Contact.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Contact) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Contact.Merge(m, src)
}
func (m *Contact) XXX_Size() int {
	return m.Size()
}
func (m *Contact) XXX_DiscardUnknown() {
	xxx_messageInfo_Contact.DiscardUnknown(m)
}

var xxx_messageInfo_Contact proto.InternalMessageInfo

func (m *Contact) GetSigChain() *crypto.SigChain {
	if m != nil {
		return m.SigChain
	}
	return nil
}

func (m *Contact) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *Contact) GetRendezvousPointSeed() []byte {
	if m != nil {
		return m.RendezvousPointSeed
	}
	return nil
}

func (m *Contact) GetSharedSecret() []byte {
	if m != nil {
		return m.SharedSecret
	}
	return nil
}

func (m *Contact) GetContactStatus() Contact_ContactStatus {
	if m != nil {
		return m.ContactStatus
	}
	return Contact_Unknown
}

type AccountConfig struct {
	RendezvousPointEnabled bool `protobuf:"varint,1,opt,name=rendezvous_point_enabled,json=rendezvousPointEnabled,proto3" json:"rendezvous_point_enabled,omitempty"`
}

func (m *AccountConfig) Reset()         { *m = AccountConfig{} }
func (m *AccountConfig) String() string { return proto.CompactTextString(m) }
func (*AccountConfig) ProtoMessage()    {}
func (*AccountConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{2}
}
func (m *AccountConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccountConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccountConfig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AccountConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountConfig.Merge(m, src)
}
func (m *AccountConfig) XXX_Size() int {
	return m.Size()
}
func (m *AccountConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountConfig.DiscardUnknown(m)
}

var xxx_messageInfo_AccountConfig proto.InternalMessageInfo

func (m *AccountConfig) GetRendezvousPointEnabled() bool {
	if m != nil {
		return m.RendezvousPointEnabled
	}
	return false
}

type DeviceConfig struct {
}

func (m *DeviceConfig) Reset()         { *m = DeviceConfig{} }
func (m *DeviceConfig) String() string { return proto.CompactTextString(m) }
func (*DeviceConfig) ProtoMessage()    {}
func (*DeviceConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf50d946d740d100, []int{3}
}
func (m *DeviceConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DeviceConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DeviceConfig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DeviceConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeviceConfig.Merge(m, src)
}
func (m *DeviceConfig) XXX_Size() int {
	return m.Size()
}
func (m *DeviceConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_DeviceConfig.DiscardUnknown(m)
}

var xxx_messageInfo_DeviceConfig proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("Contact_ContactStatus", Contact_ContactStatus_name, Contact_ContactStatus_value)
	proto.RegisterType((*Device)(nil), "Device")
	proto.RegisterType((*Contact)(nil), "Contact")
	proto.RegisterType((*AccountConfig)(nil), "AccountConfig")
	proto.RegisterType((*DeviceConfig)(nil), "DeviceConfig")
}

func init() { proto.RegisterFile("entity.proto", fileDescriptor_cf50d946d740d100) }

var fileDescriptor_cf50d946d740d100 = []byte{
	// 472 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x92, 0x41, 0x6f, 0xd3, 0x40,
	0x10, 0x85, 0xe3, 0xb4, 0x4d, 0x93, 0x69, 0x1c, 0xac, 0x45, 0x94, 0xa8, 0x07, 0xb7, 0x0a, 0xa8,
	0xea, 0x85, 0x58, 0x0a, 0x17, 0x38, 0x70, 0xa0, 0x81, 0x43, 0x85, 0x10, 0x55, 0x02, 0x1c, 0xb8,
	0x58, 0xeb, 0xf5, 0xc4, 0x59, 0xc5, 0xd9, 0x0d, 0xde, 0x71, 0x51, 0xf8, 0x15, 0xfc, 0x2c, 0x8e,
	0xbd, 0x20, 0x71, 0x44, 0xc9, 0x1f, 0x41, 0xde, 0x35, 0x2d, 0xd0, 0x93, 0xe7, 0xbd, 0xf7, 0x69,
	0xbc, 0x7e, 0x5e, 0xe8, 0xa2, 0x22, 0x49, 0xeb, 0xe1, 0xaa, 0xd0, 0xa4, 0x8f, 0x8e, 0x33, 0xad,
	0xb3, 0x1c, 0x23, 0xab, 0x92, 0x72, 0x16, 0x91, 0x5c, 0xa2, 0x21, 0xbe, 0x5c, 0xd5, 0xc0, 0x93,
	0x4c, 0xd2, 0xbc, 0x4c, 0x86, 0x42, 0x2f, 0xa3, 0x4c, 0x67, 0xfa, 0x96, 0xac, 0x94, 0x15, 0x76,
	0xaa, 0xf1, 0x87, 0x52, 0x11, 0x16, 0x8a, 0xe7, 0x91, 0x91, 0x99, 0x98, 0x73, 0xa9, 0x5c, 0x30,
	0xf8, 0x08, 0xad, 0x57, 0x78, 0x25, 0x05, 0xb2, 0xc7, 0xd0, 0x4b, 0xed, 0x14, 0xaf, 0xca, 0x24,
	0x5e, 0xe0, 0xba, 0xef, 0x9d, 0x78, 0x67, 0xdd, 0x49, 0xd7, 0xb9, 0x97, 0x65, 0xf2, 0x06, 0xd7,
	0xec, 0x14, 0xee, 0x71, 0x21, 0x74, 0xa9, 0xe8, 0x06, 0x6b, 0x5a, 0xcc, 0xaf, 0x6d, 0xc7, 0x0d,
	0x7e, 0x34, 0x61, 0x7f, 0xac, 0x15, 0x71, 0x41, 0xec, 0x14, 0x3a, 0x46, 0x66, 0xb1, 0x7d, 0xad,
	0x5d, 0x7a, 0x30, 0xea, 0x0c, 0xa7, 0x32, 0x1b, 0x57, 0xc6, 0xa4, 0x6d, 0xea, 0x89, 0x1d, 0x41,
	0x7b, 0x89, 0xc4, 0x53, 0x4e, 0xbc, 0x5e, 0x7a, 0xa3, 0xd9, 0x08, 0x1e, 0x14, 0xa8, 0x52, 0xfc,
	0x7a, 0xa5, 0x4b, 0x13, 0xaf, 0xb4, 0x54, 0x14, 0x1b, 0xc4, 0xb4, 0xbf, 0x63, 0xc1, 0xfb, 0xb7,
	0xe1, 0x65, 0x95, 0x4d, 0x11, 0x53, 0xf6, 0x08, 0x7c, 0x33, 0xe7, 0x05, 0xa6, 0xb1, 0x41, 0x51,
	0x20, 0xf5, 0x77, 0xdd, 0x07, 0x39, 0x73, 0x6a, 0x3d, 0xf6, 0x02, 0x7a, 0xc2, 0x9d, 0x33, 0x36,
	0xc4, 0xa9, 0x34, 0xfd, 0xbd, 0x13, 0xef, 0xac, 0x37, 0x3a, 0x1c, 0xd6, 0xc7, 0xff, 0xf3, 0x9c,
	0xda, 0x74, 0xe2, 0x8b, 0xbf, 0xe5, 0x40, 0x81, 0xff, 0x4f, 0xce, 0x0e, 0x60, 0xff, 0x83, 0x5a,
	0x28, 0xfd, 0x45, 0x05, 0x0d, 0x06, 0xd0, 0x7a, 0xbb, 0x36, 0x98, 0xcf, 0x02, 0xaf, 0x0a, 0x6a,
	0x32, 0x68, 0x56, 0xe2, 0x7d, 0x51, 0x1a, 0xc2, 0x34, 0xd8, 0xa9, 0xc4, 0x79, 0xae, 0xc5, 0x02,
	0xd3, 0x60, 0x97, 0xf9, 0xd0, 0x99, 0xe0, 0xe7, 0x12, 0x0d, 0x5d, 0xa8, 0x60, 0x8f, 0xf5, 0x00,
	0x6a, 0xf9, 0xae, 0xa4, 0xa0, 0x35, 0xb8, 0x00, 0xff, 0xa5, 0x2b, 0x7a, 0xac, 0xd5, 0x4c, 0x66,
	0xec, 0x19, 0xf4, 0xef, 0x14, 0x83, 0x8a, 0x27, 0x39, 0xa6, 0xb6, 0xeb, 0xf6, 0xe4, 0xf0, 0xbf,
	0x6e, 0x5e, 0xbb, 0x74, 0xd0, 0x83, 0xae, 0xfb, 0xf5, 0x6e, 0xd3, 0xf9, 0xf3, 0xef, 0x9b, 0xd0,
	0xbb, 0xde, 0x84, 0xde, 0xaf, 0x4d, 0xe8, 0x7d, 0xdb, 0x86, 0x8d, 0xeb, 0x6d, 0xd8, 0xf8, 0xb9,
	0x0d, 0x1b, 0x9f, 0x8e, 0x13, 0x2c, 0x68, 0x3d, 0x24, 0x14, 0xf3, 0xa8, 0xba, 0x6a, 0x8b, 0x2c,
	0xb2, 0x8e, 0xbd, 0x42, 0x42, 0xe7, 0x49, 0xcb, 0x4e, 0x4f, 0x7f, 0x07, 0x00, 0x00, 0xff, 0xff,
	0xbe, 0x6f, 0x5f, 0xb9, 0xc5, 0x02, 0x00, 0x00,
}

func (m *Device) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Device) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Device) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.AccountPubKey) > 0 {
		i -= len(m.AccountPubKey)
		copy(dAtA[i:], m.AccountPubKey)
		i = encodeVarintEntity(dAtA, i, uint64(len(m.AccountPubKey)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.DevicePubKey) > 0 {
		i -= len(m.DevicePubKey)
		copy(dAtA[i:], m.DevicePubKey)
		i = encodeVarintEntity(dAtA, i, uint64(len(m.DevicePubKey)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Contact) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Contact) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Contact) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ContactStatus != 0 {
		i = encodeVarintEntity(dAtA, i, uint64(m.ContactStatus))
		i--
		dAtA[i] = 0x28
	}
	if len(m.SharedSecret) > 0 {
		i -= len(m.SharedSecret)
		copy(dAtA[i:], m.SharedSecret)
		i = encodeVarintEntity(dAtA, i, uint64(len(m.SharedSecret)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.RendezvousPointSeed) > 0 {
		i -= len(m.RendezvousPointSeed)
		copy(dAtA[i:], m.RendezvousPointSeed)
		i = encodeVarintEntity(dAtA, i, uint64(len(m.RendezvousPointSeed)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Metadata) > 0 {
		i -= len(m.Metadata)
		copy(dAtA[i:], m.Metadata)
		i = encodeVarintEntity(dAtA, i, uint64(len(m.Metadata)))
		i--
		dAtA[i] = 0x12
	}
	if m.SigChain != nil {
		{
			size, err := m.SigChain.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintEntity(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AccountConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccountConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AccountConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RendezvousPointEnabled {
		i--
		if m.RendezvousPointEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *DeviceConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DeviceConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DeviceConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintEntity(dAtA []byte, offset int, v uint64) int {
	offset -= sovEntity(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Device) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DevicePubKey)
	if l > 0 {
		n += 1 + l + sovEntity(uint64(l))
	}
	l = len(m.AccountPubKey)
	if l > 0 {
		n += 1 + l + sovEntity(uint64(l))
	}
	return n
}

func (m *Contact) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SigChain != nil {
		l = m.SigChain.Size()
		n += 1 + l + sovEntity(uint64(l))
	}
	l = len(m.Metadata)
	if l > 0 {
		n += 1 + l + sovEntity(uint64(l))
	}
	l = len(m.RendezvousPointSeed)
	if l > 0 {
		n += 1 + l + sovEntity(uint64(l))
	}
	l = len(m.SharedSecret)
	if l > 0 {
		n += 1 + l + sovEntity(uint64(l))
	}
	if m.ContactStatus != 0 {
		n += 1 + sovEntity(uint64(m.ContactStatus))
	}
	return n
}

func (m *AccountConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RendezvousPointEnabled {
		n += 2
	}
	return n
}

func (m *DeviceConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovEntity(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEntity(x uint64) (n int) {
	return sovEntity(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Device) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEntity
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
			return fmt.Errorf("proto: Device: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Device: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DevicePubKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthEntity
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthEntity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DevicePubKey = append(m.DevicePubKey[:0], dAtA[iNdEx:postIndex]...)
			if m.DevicePubKey == nil {
				m.DevicePubKey = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AccountPubKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthEntity
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthEntity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AccountPubKey = append(m.AccountPubKey[:0], dAtA[iNdEx:postIndex]...)
			if m.AccountPubKey == nil {
				m.AccountPubKey = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEntity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEntity
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthEntity
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Contact) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEntity
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
			return fmt.Errorf("proto: Contact: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Contact: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigChain", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthEntity
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEntity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SigChain == nil {
				m.SigChain = &crypto.SigChain{}
			}
			if err := m.SigChain.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Metadata", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthEntity
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthEntity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Metadata = append(m.Metadata[:0], dAtA[iNdEx:postIndex]...)
			if m.Metadata == nil {
				m.Metadata = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RendezvousPointSeed", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthEntity
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthEntity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RendezvousPointSeed = append(m.RendezvousPointSeed[:0], dAtA[iNdEx:postIndex]...)
			if m.RendezvousPointSeed == nil {
				m.RendezvousPointSeed = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SharedSecret", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthEntity
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthEntity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SharedSecret = append(m.SharedSecret[:0], dAtA[iNdEx:postIndex]...)
			if m.SharedSecret == nil {
				m.SharedSecret = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContactStatus", wireType)
			}
			m.ContactStatus = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ContactStatus |= Contact_ContactStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipEntity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEntity
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthEntity
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *AccountConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEntity
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
			return fmt.Errorf("proto: AccountConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AccountConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RendezvousPointEnabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.RendezvousPointEnabled = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipEntity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEntity
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthEntity
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DeviceConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEntity
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
			return fmt.Errorf("proto: DeviceConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DeviceConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipEntity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthEntity
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthEntity
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipEntity(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEntity
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
					return 0, ErrIntOverflowEntity
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEntity
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
				return 0, ErrInvalidLengthEntity
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthEntity
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowEntity
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipEntity(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthEntity
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthEntity = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEntity   = fmt.Errorf("proto: integer overflow")
)
