// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: go-internal/records.proto

package tinder

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

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

type Records struct {
	Records              []*Record `protobuf:"bytes,1,rep,name=records,proto3" json:"records,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Records) Reset()         { *m = Records{} }
func (m *Records) String() string { return proto.CompactTextString(m) }
func (*Records) ProtoMessage()    {}
func (*Records) Descriptor() ([]byte, []int) {
	return fileDescriptor_4ec813e1227eb54c, []int{0}
}
func (m *Records) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Records) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Records.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Records) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Records.Merge(m, src)
}
func (m *Records) XXX_Size() int {
	return m.Size()
}
func (m *Records) XXX_DiscardUnknown() {
	xxx_messageInfo_Records.DiscardUnknown(m)
}

var xxx_messageInfo_Records proto.InternalMessageInfo

func (m *Records) GetRecords() []*Record {
	if m != nil {
		return m.Records
	}
	return nil
}

type Record struct {
	Cid                  string   `protobuf:"bytes,1,opt,name=cid,proto3" json:"cid,omitempty"`
	Expire               int64    `protobuf:"varint,2,opt,name=expire,proto3" json:"expire,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Record) Reset()         { *m = Record{} }
func (m *Record) String() string { return proto.CompactTextString(m) }
func (*Record) ProtoMessage()    {}
func (*Record) Descriptor() ([]byte, []int) {
	return fileDescriptor_4ec813e1227eb54c, []int{1}
}
func (m *Record) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Record) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Record.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Record) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Record.Merge(m, src)
}
func (m *Record) XXX_Size() int {
	return m.Size()
}
func (m *Record) XXX_DiscardUnknown() {
	xxx_messageInfo_Record.DiscardUnknown(m)
}

var xxx_messageInfo_Record proto.InternalMessageInfo

func (m *Record) GetCid() string {
	if m != nil {
		return m.Cid
	}
	return ""
}

func (m *Record) GetExpire() int64 {
	if m != nil {
		return m.Expire
	}
	return 0
}

func init() {
	proto.RegisterType((*Records)(nil), "tinder.Records")
	proto.RegisterType((*Record)(nil), "tinder.Record")
}

func init() { proto.RegisterFile("go-internal/records.proto", fileDescriptor_4ec813e1227eb54c) }

var fileDescriptor_4ec813e1227eb54c = []byte{
	// 189 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4c, 0xcf, 0xd7, 0xcd,
	0xcc, 0x2b, 0x49, 0x2d, 0xca, 0x4b, 0xcc, 0xd1, 0x2f, 0x4a, 0x4d, 0xce, 0x2f, 0x4a, 0x29, 0xd6,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2b, 0xc9, 0xcc, 0x4b, 0x49, 0x2d, 0x92, 0x12, 0x49,
	0xcf, 0x4f, 0xcf, 0x07, 0x0b, 0xe9, 0x83, 0x58, 0x10, 0x59, 0x25, 0x63, 0x2e, 0xf6, 0x20, 0x88,
	0x72, 0x21, 0x0d, 0x2e, 0x76, 0xa8, 0x4e, 0x09, 0x46, 0x05, 0x66, 0x0d, 0x6e, 0x23, 0x3e, 0x3d,
	0x88, 0x56, 0x3d, 0x88, 0x8a, 0x20, 0x98, 0xb4, 0x92, 0x11, 0x17, 0x1b, 0x44, 0x48, 0x48, 0x80,
	0x8b, 0x39, 0x39, 0x33, 0x45, 0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x08, 0xc4, 0x14, 0x12, 0xe3,
	0x62, 0x4b, 0xad, 0x28, 0xc8, 0x2c, 0x4a, 0x95, 0x60, 0x52, 0x60, 0xd4, 0x60, 0x0e, 0x82, 0xf2,
	0x9c, 0x4c, 0x2f, 0x3c, 0x94, 0x63, 0x38, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07,
	0x8f, 0xe4, 0x18, 0xa3, 0x94, 0x93, 0x52, 0x8b, 0x4a, 0x2a, 0xf5, 0x4a, 0x52, 0x93, 0x33, 0xf4,
	0xc1, 0x4c, 0xfd, 0xf4, 0x7c, 0x7d, 0xb8, 0x27, 0x20, 0x56, 0x27, 0xb1, 0x81, 0x9d, 0x69, 0x0c,
	0x08, 0x00, 0x00, 0xff, 0xff, 0xa8, 0x94, 0xa8, 0x61, 0xe1, 0x00, 0x00, 0x00,
}

func (m *Records) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Records) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Records) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Records) > 0 {
		for iNdEx := len(m.Records) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Records[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintRecords(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Record) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Record) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Record) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Expire != 0 {
		i = encodeVarintRecords(dAtA, i, uint64(m.Expire))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Cid) > 0 {
		i -= len(m.Cid)
		copy(dAtA[i:], m.Cid)
		i = encodeVarintRecords(dAtA, i, uint64(len(m.Cid)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintRecords(dAtA []byte, offset int, v uint64) int {
	offset -= sovRecords(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Records) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Records) > 0 {
		for _, e := range m.Records {
			l = e.Size()
			n += 1 + l + sovRecords(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Record) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Cid)
	if l > 0 {
		n += 1 + l + sovRecords(uint64(l))
	}
	if m.Expire != 0 {
		n += 1 + sovRecords(uint64(m.Expire))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovRecords(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRecords(x uint64) (n int) {
	return sovRecords(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Records) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRecords
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
			return fmt.Errorf("proto: Records: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Records: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Records", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRecords
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
				return ErrInvalidLengthRecords
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRecords
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Records = append(m.Records, &Record{})
			if err := m.Records[len(m.Records)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRecords(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRecords
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRecords
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
func (m *Record) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRecords
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
			return fmt.Errorf("proto: Record: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Record: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRecords
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRecords
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRecords
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Cid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Expire", wireType)
			}
			m.Expire = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRecords
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Expire |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRecords(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRecords
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRecords
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
func skipRecords(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRecords
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
					return 0, ErrIntOverflowRecords
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
					return 0, ErrIntOverflowRecords
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
				return 0, ErrInvalidLengthRecords
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRecords
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRecords
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRecords        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRecords          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRecords = fmt.Errorf("proto: unexpected end of group")
)
