// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proto/crypto/merkle/types.proto

package merkle

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

type Proof struct {
	Total    int64    `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Index    int64    `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	LeafHash []byte   `protobuf:"bytes,3,opt,name=leaf_hash,json=leafHash,proto3" json:"leaf_hash,omitempty"`
	Aunts    [][]byte `protobuf:"bytes,4,rep,name=aunts,proto3" json:"aunts,omitempty"`
}

func (m *Proof) Reset()         { *m = Proof{} }
func (m *Proof) String() string { return proto.CompactTextString(m) }
func (*Proof) ProtoMessage()    {}
func (*Proof) Descriptor() ([]byte, []int) {
	return fileDescriptor_57e39eefdaf7ae96, []int{0}
}
func (m *Proof) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Proof) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Proof.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Proof) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Proof.Merge(m, src)
}
func (m *Proof) XXX_Size() int {
	return m.Size()
}
func (m *Proof) XXX_DiscardUnknown() {
	xxx_messageInfo_Proof.DiscardUnknown(m)
}

var xxx_messageInfo_Proof proto.InternalMessageInfo

func (m *Proof) GetTotal() int64 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *Proof) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *Proof) GetLeafHash() []byte {
	if m != nil {
		return m.LeafHash
	}
	return nil
}

func (m *Proof) GetAunts() [][]byte {
	if m != nil {
		return m.Aunts
	}
	return nil
}

type ValueOp struct {
	// Encoded in ProofOp.Key.
	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// To encode in ProofOp.Data
	Proof *Proof `protobuf:"bytes,2,opt,name=proof,proto3" json:"proof,omitempty"`
}

func (m *ValueOp) Reset()         { *m = ValueOp{} }
func (m *ValueOp) String() string { return proto.CompactTextString(m) }
func (*ValueOp) ProtoMessage()    {}
func (*ValueOp) Descriptor() ([]byte, []int) {
	return fileDescriptor_57e39eefdaf7ae96, []int{1}
}
func (m *ValueOp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ValueOp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ValueOp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ValueOp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValueOp.Merge(m, src)
}
func (m *ValueOp) XXX_Size() int {
	return m.Size()
}
func (m *ValueOp) XXX_DiscardUnknown() {
	xxx_messageInfo_ValueOp.DiscardUnknown(m)
}

var xxx_messageInfo_ValueOp proto.InternalMessageInfo

func (m *ValueOp) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *ValueOp) GetProof() *Proof {
	if m != nil {
		return m.Proof
	}
	return nil
}

type DominoOp struct {
	Key    string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Input  string `protobuf:"bytes,2,opt,name=input,proto3" json:"input,omitempty"`
	Output string `protobuf:"bytes,3,opt,name=output,proto3" json:"output,omitempty"`
}

func (m *DominoOp) Reset()         { *m = DominoOp{} }
func (m *DominoOp) String() string { return proto.CompactTextString(m) }
func (*DominoOp) ProtoMessage()    {}
func (*DominoOp) Descriptor() ([]byte, []int) {
	return fileDescriptor_57e39eefdaf7ae96, []int{2}
}
func (m *DominoOp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DominoOp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DominoOp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DominoOp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DominoOp.Merge(m, src)
}
func (m *DominoOp) XXX_Size() int {
	return m.Size()
}
func (m *DominoOp) XXX_DiscardUnknown() {
	xxx_messageInfo_DominoOp.DiscardUnknown(m)
}

var xxx_messageInfo_DominoOp proto.InternalMessageInfo

func (m *DominoOp) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *DominoOp) GetInput() string {
	if m != nil {
		return m.Input
	}
	return ""
}

func (m *DominoOp) GetOutput() string {
	if m != nil {
		return m.Output
	}
	return ""
}

// ProofOp defines an operation used for calculating Merkle root
// The data could be arbitrary format, providing nessecary data
// for example neighbouring node hash
type ProofOp struct {
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Key  []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Data []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *ProofOp) Reset()         { *m = ProofOp{} }
func (m *ProofOp) String() string { return proto.CompactTextString(m) }
func (*ProofOp) ProtoMessage()    {}
func (*ProofOp) Descriptor() ([]byte, []int) {
	return fileDescriptor_57e39eefdaf7ae96, []int{3}
}
func (m *ProofOp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProofOp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProofOp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProofOp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProofOp.Merge(m, src)
}
func (m *ProofOp) XXX_Size() int {
	return m.Size()
}
func (m *ProofOp) XXX_DiscardUnknown() {
	xxx_messageInfo_ProofOp.DiscardUnknown(m)
}

var xxx_messageInfo_ProofOp proto.InternalMessageInfo

func (m *ProofOp) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ProofOp) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *ProofOp) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// ProofOps is Merkle proof defined by the list of ProofOps
type ProofOps struct {
	Ops []ProofOp `protobuf:"bytes,1,rep,name=ops,proto3" json:"ops"`
}

func (m *ProofOps) Reset()         { *m = ProofOps{} }
func (m *ProofOps) String() string { return proto.CompactTextString(m) }
func (*ProofOps) ProtoMessage()    {}
func (*ProofOps) Descriptor() ([]byte, []int) {
	return fileDescriptor_57e39eefdaf7ae96, []int{4}
}
func (m *ProofOps) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProofOps) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProofOps.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProofOps) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProofOps.Merge(m, src)
}
func (m *ProofOps) XXX_Size() int {
	return m.Size()
}
func (m *ProofOps) XXX_DiscardUnknown() {
	xxx_messageInfo_ProofOps.DiscardUnknown(m)
}

var xxx_messageInfo_ProofOps proto.InternalMessageInfo

func (m *ProofOps) GetOps() []ProofOp {
	if m != nil {
		return m.Ops
	}
	return nil
}

func init() {
	proto.RegisterType((*Proof)(nil), "tendermint.crypto.merkle.Proof")
	proto.RegisterType((*ValueOp)(nil), "tendermint.crypto.merkle.ValueOp")
	proto.RegisterType((*DominoOp)(nil), "tendermint.crypto.merkle.DominoOp")
	proto.RegisterType((*ProofOp)(nil), "tendermint.crypto.merkle.ProofOp")
	proto.RegisterType((*ProofOps)(nil), "tendermint.crypto.merkle.ProofOps")
}

func init() { proto.RegisterFile("proto/crypto/merkle/types.proto", fileDescriptor_57e39eefdaf7ae96) }

var fileDescriptor_57e39eefdaf7ae96 = []byte{
	// 371 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xbf, 0x4e, 0xeb, 0x30,
	0x18, 0xc5, 0x93, 0x26, 0xfd, 0xe7, 0x76, 0xb8, 0xb2, 0xae, 0xae, 0xa2, 0x8b, 0x94, 0x86, 0x0c,
	0x28, 0x53, 0x22, 0x15, 0x18, 0x58, 0x0b, 0x48, 0x88, 0x81, 0x22, 0x0f, 0x0c, 0x2c, 0x95, 0xdb,
	0xb8, 0x4d, 0xd4, 0x24, 0xb6, 0x1c, 0x47, 0x22, 0x6f, 0xc1, 0x63, 0x75, 0xec, 0xc8, 0x84, 0x50,
	0xfb, 0x22, 0xc8, 0x76, 0x50, 0x41, 0x02, 0xb1, 0x9d, 0x73, 0xf2, 0xf9, 0xe7, 0xf3, 0x45, 0x06,
	0x23, 0xc6, 0xa9, 0xa0, 0xd1, 0x82, 0xd7, 0x4c, 0xd0, 0x28, 0x27, 0x7c, 0x9d, 0x91, 0x48, 0xd4,
	0x8c, 0x94, 0xa1, 0xfa, 0x02, 0x1d, 0x41, 0x8a, 0x98, 0xf0, 0x3c, 0x2d, 0x44, 0xa8, 0xa7, 0x42,
	0x3d, 0xf5, 0xff, 0x44, 0x24, 0x29, 0x8f, 0x67, 0x0c, 0x73, 0x51, 0x47, 0x1a, 0xb3, 0xa2, 0x2b,
	0x7a, 0x50, 0x9a, 0xe0, 0x2f, 0x41, 0xfb, 0x9e, 0x53, 0xba, 0x84, 0x7f, 0x41, 0x5b, 0x50, 0x81,
	0x33, 0xc7, 0xf4, 0xcc, 0xc0, 0x42, 0xda, 0xc8, 0x34, 0x2d, 0x62, 0xf2, 0xe4, 0xb4, 0x74, 0xaa,
	0x0c, 0x3c, 0x02, 0xfd, 0x8c, 0xe0, 0xe5, 0x2c, 0xc1, 0x65, 0xe2, 0x58, 0x9e, 0x19, 0x0c, 0x51,
	0x4f, 0x06, 0x37, 0xb8, 0x4c, 0xe4, 0x11, 0x5c, 0x15, 0xa2, 0x74, 0x6c, 0xcf, 0x0a, 0x86, 0x48,
	0x1b, 0x1f, 0x81, 0xee, 0x03, 0xce, 0x2a, 0x32, 0x65, 0xf0, 0x0f, 0xb0, 0xd6, 0xa4, 0x56, 0xf7,
	0x0c, 0x91, 0x94, 0xf0, 0x1c, 0xb4, 0x99, 0x2c, 0xa1, 0x6e, 0x19, 0x8c, 0x47, 0xe1, 0x4f, 0x6b,
	0x85, 0xaa, 0x2b, 0xd2, 0xd3, 0xfe, 0x2d, 0xe8, 0x5d, 0xd1, 0x3c, 0x2d, 0xe8, 0x57, 0x68, 0x5f,
	0x43, 0x55, 0x75, 0x56, 0x09, 0x05, 0xed, 0x23, 0x6d, 0xe0, 0x3f, 0xd0, 0xa1, 0x95, 0x90, 0xb1,
	0xa5, 0xe2, 0xc6, 0xf9, 0x97, 0xa0, 0xab, 0xd8, 0x53, 0x06, 0x21, 0xb0, 0xe5, 0x3f, 0x6e, 0x58,
	0x4a, 0x7f, 0xe0, 0x5b, 0x87, 0xce, 0x10, 0xd8, 0x31, 0x16, 0xb8, 0x59, 0x5f, 0x69, 0xff, 0x1a,
	0xf4, 0x1a, 0x48, 0x09, 0x2f, 0x80, 0x45, 0x59, 0xe9, 0x98, 0x9e, 0x15, 0x0c, 0xc6, 0xc7, 0xbf,
	0x6c, 0x34, 0x65, 0x13, 0x7b, 0xf3, 0x3a, 0x32, 0x90, 0x3c, 0x33, 0xb9, 0xdb, 0xec, 0x5c, 0x73,
	0xbb, 0x73, 0xcd, 0xb7, 0x9d, 0x6b, 0x3e, 0xef, 0x5d, 0x63, 0xbb, 0x77, 0x8d, 0x97, 0xbd, 0x6b,
	0x3c, 0x9e, 0xad, 0x52, 0x91, 0x54, 0xf3, 0x70, 0x41, 0xf3, 0xe8, 0x40, 0xfc, 0x2c, 0xbf, 0x79,
	0x31, 0xf3, 0x8e, 0x0a, 0x4f, 0xdf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xe3, 0x94, 0xcb, 0xb9, 0x4f,
	0x02, 0x00, 0x00,
}

func (m *Proof) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Proof) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Proof) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Aunts) > 0 {
		for iNdEx := len(m.Aunts) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Aunts[iNdEx])
			copy(dAtA[i:], m.Aunts[iNdEx])
			i = encodeVarintTypes(dAtA, i, uint64(len(m.Aunts[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.LeafHash) > 0 {
		i -= len(m.LeafHash)
		copy(dAtA[i:], m.LeafHash)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.LeafHash)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Index != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x10
	}
	if m.Total != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Total))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ValueOp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ValueOp) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ValueOp) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Proof != nil {
		{
			size, err := m.Proof.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DominoOp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DominoOp) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DominoOp) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Output) > 0 {
		i -= len(m.Output)
		copy(dAtA[i:], m.Output)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Output)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Input) > 0 {
		i -= len(m.Input)
		copy(dAtA[i:], m.Input)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Input)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ProofOp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProofOp) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProofOp) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Type) > 0 {
		i -= len(m.Type)
		copy(dAtA[i:], m.Type)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Type)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ProofOps) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProofOps) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProofOps) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Ops) > 0 {
		for iNdEx := len(m.Ops) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Ops[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Proof) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Total != 0 {
		n += 1 + sovTypes(uint64(m.Total))
	}
	if m.Index != 0 {
		n += 1 + sovTypes(uint64(m.Index))
	}
	l = len(m.LeafHash)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if len(m.Aunts) > 0 {
		for _, b := range m.Aunts {
			l = len(b)
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	return n
}

func (m *ValueOp) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.Proof != nil {
		l = m.Proof.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *DominoOp) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Input)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Output)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *ProofOp) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Type)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *ProofOps) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Ops) > 0 {
		for _, e := range m.Ops {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Proof) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: Proof: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Proof: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Total", wireType)
			}
			m.Total = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Total |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LeafHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LeafHash = append(m.LeafHash[:0], dAtA[iNdEx:postIndex]...)
			if m.LeafHash == nil {
				m.LeafHash = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Aunts", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Aunts = append(m.Aunts, make([]byte, postIndex-iNdEx))
			copy(m.Aunts[len(m.Aunts)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *ValueOp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: ValueOp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ValueOp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = append(m.Key[:0], dAtA[iNdEx:postIndex]...)
			if m.Key == nil {
				m.Key = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Proof", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Proof == nil {
				m.Proof = &Proof{}
			}
			if err := m.Proof.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *DominoOp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: DominoOp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DominoOp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Input", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Input = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Output", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Output = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *ProofOp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: ProofOp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProofOp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = append(m.Key[:0], dAtA[iNdEx:postIndex]...)
			if m.Key == nil {
				m.Key = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *ProofOps) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: ProofOps: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProofOps: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ops", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ops = append(m.Ops, ProofOp{})
			if err := m.Ops[len(m.Ops)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
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
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
