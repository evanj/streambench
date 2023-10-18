// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.2
// source: messages/messages.proto

package messages

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DuplicateTest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique random ID for the goroutine that published this message.
	GoroutineId string `protobuf:"bytes,1,opt,name=goroutine_id,json=goroutineId,proto3" json:"goroutine_id,omitempty"`
	// An increasing sequence number for the goroutine that published this message.
	Sequence int64 `protobuf:"varint,2,opt,name=sequence,proto3" json:"sequence,omitempty"`
	// The time this message was created, immediately before it is published.
	Created *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created,proto3" json:"created,omitempty"`
}

func (x *DuplicateTest) Reset() {
	*x = DuplicateTest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DuplicateTest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DuplicateTest) ProtoMessage() {}

func (x *DuplicateTest) ProtoReflect() protoreflect.Message {
	mi := &file_messages_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DuplicateTest.ProtoReflect.Descriptor instead.
func (*DuplicateTest) Descriptor() ([]byte, []int) {
	return file_messages_messages_proto_rawDescGZIP(), []int{0}
}

func (x *DuplicateTest) GetGoroutineId() string {
	if x != nil {
		return x.GoroutineId
	}
	return ""
}

func (x *DuplicateTest) GetSequence() int64 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *DuplicateTest) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

var File_messages_messages_proto protoreflect.FileDescriptor

var file_messages_messages_proto_rawDesc = []byte{
	0x0a, 0x17, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x84, 0x01, 0x0a, 0x0d, 0x44, 0x75, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x67, 0x6f, 0x72, 0x6f, 0x75, 0x74,
	0x69, 0x6e, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x67, 0x6f,
	0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x65, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x71,
	0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x71,
	0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x34, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x42, 0x27, 0x5a, 0x25, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x76, 0x61, 0x6e, 0x6a, 0x2f,
	0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x2f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_messages_messages_proto_rawDescOnce sync.Once
	file_messages_messages_proto_rawDescData = file_messages_messages_proto_rawDesc
)

func file_messages_messages_proto_rawDescGZIP() []byte {
	file_messages_messages_proto_rawDescOnce.Do(func() {
		file_messages_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_messages_messages_proto_rawDescData)
	})
	return file_messages_messages_proto_rawDescData
}

var file_messages_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_messages_messages_proto_goTypes = []interface{}{
	(*DuplicateTest)(nil),         // 0: messages.DuplicateTest
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_messages_messages_proto_depIdxs = []int32{
	1, // 0: messages.DuplicateTest.created:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_messages_messages_proto_init() }
func file_messages_messages_proto_init() {
	if File_messages_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_messages_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DuplicateTest); i {
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
			RawDescriptor: file_messages_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_messages_messages_proto_goTypes,
		DependencyIndexes: file_messages_messages_proto_depIdxs,
		MessageInfos:      file_messages_messages_proto_msgTypes,
	}.Build()
	File_messages_messages_proto = out.File
	file_messages_messages_proto_rawDesc = nil
	file_messages_messages_proto_goTypes = nil
	file_messages_messages_proto_depIdxs = nil
}
