// Code generated by protoc-gen-go. DO NOT EDIT.
// source: messages/messages.proto

package messages

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type DuplicateTest struct {
	// Unique random ID for the goroutine that published this message.
	GoroutineId string `protobuf:"bytes,1,opt,name=goroutine_id,json=goroutineId,proto3" json:"goroutine_id,omitempty"`
	// An increasing sequence number for the goroutine that published this message.
	Sequence int64 `protobuf:"varint,2,opt,name=sequence,proto3" json:"sequence,omitempty"`
	// The time this message was created, immediately before it is published.
	Created              *timestamp.Timestamp `protobuf:"bytes,3,opt,name=created,proto3" json:"created,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *DuplicateTest) Reset()         { *m = DuplicateTest{} }
func (m *DuplicateTest) String() string { return proto.CompactTextString(m) }
func (*DuplicateTest) ProtoMessage()    {}
func (*DuplicateTest) Descriptor() ([]byte, []int) {
	return fileDescriptor_83994550f81e9f35, []int{0}
}

func (m *DuplicateTest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DuplicateTest.Unmarshal(m, b)
}
func (m *DuplicateTest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DuplicateTest.Marshal(b, m, deterministic)
}
func (m *DuplicateTest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DuplicateTest.Merge(m, src)
}
func (m *DuplicateTest) XXX_Size() int {
	return xxx_messageInfo_DuplicateTest.Size(m)
}
func (m *DuplicateTest) XXX_DiscardUnknown() {
	xxx_messageInfo_DuplicateTest.DiscardUnknown(m)
}

var xxx_messageInfo_DuplicateTest proto.InternalMessageInfo

func (m *DuplicateTest) GetGoroutineId() string {
	if m != nil {
		return m.GoroutineId
	}
	return ""
}

func (m *DuplicateTest) GetSequence() int64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *DuplicateTest) GetCreated() *timestamp.Timestamp {
	if m != nil {
		return m.Created
	}
	return nil
}

func init() {
	proto.RegisterType((*DuplicateTest)(nil), "messages.DuplicateTest")
}

func init() { proto.RegisterFile("messages/messages.proto", fileDescriptor_83994550f81e9f35) }

var fileDescriptor_83994550f81e9f35 = []byte{
	// 204 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x3c, 0x8f, 0xb1, 0x4a, 0xc6, 0x30,
	0x14, 0x85, 0x89, 0x3f, 0x68, 0x4d, 0x75, 0xc9, 0x62, 0xe9, 0x62, 0x15, 0xc4, 0x4e, 0x09, 0xa8,
	0x4f, 0x20, 0x2e, 0xae, 0xa5, 0x93, 0x8b, 0x24, 0xe9, 0x35, 0x8d, 0x34, 0x49, 0x4d, 0x6e, 0x7c,
	0x02, 0x1f, 0x5c, 0x68, 0x49, 0xb7, 0x7b, 0x0e, 0x1f, 0x97, 0xef, 0xd0, 0x1b, 0x07, 0x29, 0x49,
	0x03, 0x49, 0x94, 0x83, 0xaf, 0x31, 0x60, 0x60, 0x55, 0xc9, 0xed, 0xad, 0x09, 0xc1, 0x2c, 0x20,
	0xb6, 0x5e, 0xe5, 0x2f, 0x81, 0xd6, 0x41, 0x42, 0xe9, 0xd6, 0x1d, 0xbd, 0xff, 0x23, 0xf4, 0xfa,
	0x2d, 0xaf, 0x8b, 0xd5, 0x12, 0x61, 0x84, 0x84, 0xec, 0x8e, 0x5e, 0x99, 0x10, 0x43, 0x46, 0xeb,
	0xe1, 0xd3, 0x4e, 0x0d, 0xe9, 0x48, 0x7f, 0x39, 0xd4, 0x47, 0xf7, 0x3e, 0xb1, 0x96, 0x56, 0x09,
	0x7e, 0x32, 0x78, 0x0d, 0xcd, 0x59, 0x47, 0xfa, 0xd3, 0x70, 0x64, 0xf6, 0x42, 0x2f, 0x74, 0x04,
	0x89, 0x30, 0x35, 0xa7, 0x8e, 0xf4, 0xf5, 0x53, 0xcb, 0x77, 0x07, 0x5e, 0x1c, 0xf8, 0x58, 0x1c,
	0x86, 0x82, 0xbe, 0x3e, 0x7e, 0x3c, 0x18, 0x8b, 0x73, 0x56, 0x5c, 0x07, 0x27, 0xe0, 0x57, 0xfa,
	0x6f, 0x91, 0x30, 0x82, 0x74, 0x0a, 0xbc, 0x9e, 0x8f, 0x81, 0xea, 0x7c, 0xfb, 0xf2, 0xfc, 0x1f,
	0x00, 0x00, 0xff, 0xff, 0x90, 0x05, 0xfb, 0x6c, 0xfc, 0x00, 0x00, 0x00,
}