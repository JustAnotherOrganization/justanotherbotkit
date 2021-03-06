// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: channel.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	channel.proto
	event.proto

It has these top-level messages:
	BaseChannel
	BaseEvent
*/
package pb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// BaseChannel is a network channel.
type BaseChannel struct {
	ID        string                 `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name      string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	MemberIDs []string               `protobuf:"bytes,10,rep,name=MemberIDs" json:"MemberIDs,omitempty"`
	Archived  bool                   `protobuf:"varint,11,opt,name=archived,proto3" json:"archived,omitempty"`
	Topic     *BaseChannel_Attribute `protobuf:"bytes,20,opt,name=topic" json:"topic,omitempty"`
	Purpose   *BaseChannel_Attribute `protobuf:"bytes,21,opt,name=purpose" json:"purpose,omitempty"`
}

func (m *BaseChannel) Reset()                    { *m = BaseChannel{} }
func (m *BaseChannel) String() string            { return proto.CompactTextString(m) }
func (*BaseChannel) ProtoMessage()               {}
func (*BaseChannel) Descriptor() ([]byte, []int) { return fileDescriptorChannel, []int{0} }

func (m *BaseChannel) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *BaseChannel) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *BaseChannel) GetMemberIDs() []string {
	if m != nil {
		return m.MemberIDs
	}
	return nil
}

func (m *BaseChannel) GetArchived() bool {
	if m != nil {
		return m.Archived
	}
	return false
}

func (m *BaseChannel) GetTopic() *BaseChannel_Attribute {
	if m != nil {
		return m.Topic
	}
	return nil
}

func (m *BaseChannel) GetPurpose() *BaseChannel_Attribute {
	if m != nil {
		return m.Purpose
	}
	return nil
}

type BaseChannel_Attribute struct {
	Value     string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	CreatorID string `protobuf:"bytes,2,opt,name=CreatorID,proto3" json:"CreatorID,omitempty"`
}

func (m *BaseChannel_Attribute) Reset()                    { *m = BaseChannel_Attribute{} }
func (m *BaseChannel_Attribute) String() string            { return proto.CompactTextString(m) }
func (*BaseChannel_Attribute) ProtoMessage()               {}
func (*BaseChannel_Attribute) Descriptor() ([]byte, []int) { return fileDescriptorChannel, []int{0, 0} }

func (m *BaseChannel_Attribute) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *BaseChannel_Attribute) GetCreatorID() string {
	if m != nil {
		return m.CreatorID
	}
	return ""
}

func init() {
	proto.RegisterType((*BaseChannel)(nil), "transport_proto.BaseChannel")
	proto.RegisterType((*BaseChannel_Attribute)(nil), "transport_proto.BaseChannel.Attribute")
}

func init() { proto.RegisterFile("channel.proto", fileDescriptorChannel) }

var fileDescriptorChannel = []byte{
	// 235 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x8f, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0xc9, 0xda, 0x6a, 0x33, 0x41, 0x85, 0xa1, 0xc2, 0x52, 0x3c, 0x04, 0x0f, 0x92, 0x53,
	0x0e, 0x7a, 0x15, 0xd4, 0x36, 0x97, 0x1c, 0xbc, 0xe4, 0xe8, 0x45, 0x36, 0x71, 0xa0, 0x81, 0x76,
	0x77, 0x99, 0x9d, 0xf4, 0x57, 0xf9, 0x23, 0xc5, 0x4d, 0x6d, 0xc5, 0x5b, 0x6f, 0xfb, 0x3e, 0xf6,
	0x31, 0xef, 0x83, 0xcb, 0x6e, 0x6d, 0xac, 0xa5, 0x4d, 0xe9, 0xd9, 0x89, 0xc3, 0x6b, 0x61, 0x63,
	0x83, 0x77, 0x2c, 0x1f, 0x11, 0xdc, 0x7d, 0x29, 0xc8, 0x96, 0x26, 0xd0, 0x6a, 0xfc, 0x86, 0x57,
	0xa0, 0xea, 0x4a, 0x27, 0x79, 0x52, 0xa4, 0x8d, 0xaa, 0x2b, 0x44, 0x98, 0x58, 0xb3, 0x25, 0xad,
	0x22, 0x89, 0x6f, 0xbc, 0x85, 0xf4, 0x8d, 0xb6, 0x2d, 0x71, 0x5d, 0x05, 0x0d, 0xf9, 0x59, 0x91,
	0x36, 0x47, 0x80, 0x0b, 0x98, 0x19, 0xee, 0xd6, 0xfd, 0x8e, 0x3e, 0x75, 0x96, 0x27, 0xc5, 0xac,
	0x39, 0x64, 0x7c, 0x82, 0xa9, 0x38, 0xdf, 0x77, 0x7a, 0x9e, 0x27, 0x45, 0xf6, 0x70, 0x5f, 0xfe,
	0x9b, 0x53, 0xfe, 0x99, 0x52, 0xbe, 0x8a, 0x70, 0xdf, 0x0e, 0x42, 0xcd, 0x58, 0xc2, 0x17, 0xb8,
	0xf0, 0x03, 0x7b, 0x17, 0x48, 0xdf, 0x9c, 0xd4, 0xff, 0xad, 0x2d, 0x9e, 0x21, 0x3d, 0x50, 0x9c,
	0xc3, 0x74, 0x67, 0x36, 0x03, 0xed, 0x6d, 0xc7, 0xf0, 0x23, 0xb7, 0x62, 0x32, 0xe2, 0xb8, 0xae,
	0xf6, 0xd6, 0x47, 0xb0, 0x9c, 0xbc, 0x2b, 0xdf, 0xb6, 0xe7, 0xf1, 0xd8, 0xe3, 0x77, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xee, 0x7b, 0x4e, 0xe4, 0x5d, 0x01, 0x00, 0x00,
}
