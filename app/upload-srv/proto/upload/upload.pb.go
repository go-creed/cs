// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/upload/upload.proto

package go_micro_cs_service_upload

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
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

type Bytes struct {
	Content              []byte   `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	Size                 int64    `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Bytes) Reset()         { *m = Bytes{} }
func (m *Bytes) String() string { return proto.CompactTextString(m) }
func (*Bytes) ProtoMessage()    {}
func (*Bytes) Descriptor() ([]byte, []int) {
	return fileDescriptor_5508b024e4d0ced2, []int{0}
}

func (m *Bytes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Bytes.Unmarshal(m, b)
}
func (m *Bytes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Bytes.Marshal(b, m, deterministic)
}
func (m *Bytes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bytes.Merge(m, src)
}
func (m *Bytes) XXX_Size() int {
	return xxx_messageInfo_Bytes.Size(m)
}
func (m *Bytes) XXX_DiscardUnknown() {
	xxx_messageInfo_Bytes.DiscardUnknown(m)
}

var xxx_messageInfo_Bytes proto.InternalMessageInfo

func (m *Bytes) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func (m *Bytes) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

type FileInfo struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	FileName             string   `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	Size                 int64    `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	Filesha256           string   `protobuf:"bytes,4,opt,name=filesha256,proto3" json:"filesha256,omitempty"`
	Location             string   `protobuf:"bytes,5,opt,name=location,proto3" json:"location,omitempty"`
	CreatedAt            int64    `protobuf:"varint,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt            int64    `protobuf:"varint,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt            int64    `protobuf:"varint,8,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FileInfo) Reset()         { *m = FileInfo{} }
func (m *FileInfo) String() string { return proto.CompactTextString(m) }
func (*FileInfo) ProtoMessage()    {}
func (*FileInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_5508b024e4d0ced2, []int{1}
}

func (m *FileInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileInfo.Unmarshal(m, b)
}
func (m *FileInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileInfo.Marshal(b, m, deterministic)
}
func (m *FileInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileInfo.Merge(m, src)
}
func (m *FileInfo) XXX_Size() int {
	return xxx_messageInfo_FileInfo.Size(m)
}
func (m *FileInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_FileInfo.DiscardUnknown(m)
}

var xxx_messageInfo_FileInfo proto.InternalMessageInfo

func (m *FileInfo) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *FileInfo) GetFileName() string {
	if m != nil {
		return m.FileName
	}
	return ""
}

func (m *FileInfo) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *FileInfo) GetFilesha256() string {
	if m != nil {
		return m.Filesha256
	}
	return ""
}

func (m *FileInfo) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

func (m *FileInfo) GetCreatedAt() int64 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

func (m *FileInfo) GetUpdatedAt() int64 {
	if m != nil {
		return m.UpdatedAt
	}
	return 0
}

func (m *FileInfo) GetDeletedAt() int64 {
	if m != nil {
		return m.DeletedAt
	}
	return 0
}

type StreamingResponse struct {
	Size                 int64    `protobuf:"varint,1,opt,name=size,proto3" json:"size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StreamingResponse) Reset()         { *m = StreamingResponse{} }
func (m *StreamingResponse) String() string { return proto.CompactTextString(m) }
func (*StreamingResponse) ProtoMessage()    {}
func (*StreamingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5508b024e4d0ced2, []int{2}
}

func (m *StreamingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamingResponse.Unmarshal(m, b)
}
func (m *StreamingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamingResponse.Marshal(b, m, deterministic)
}
func (m *StreamingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamingResponse.Merge(m, src)
}
func (m *StreamingResponse) XXX_Size() int {
	return xxx_messageInfo_StreamingResponse.Size(m)
}
func (m *StreamingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StreamingResponse proto.InternalMessageInfo

func (m *StreamingResponse) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func init() {
	proto.RegisterType((*Bytes)(nil), "go.micro.cs.service.upload.Bytes")
	proto.RegisterType((*FileInfo)(nil), "go.micro.cs.service.upload.FileInfo")
	proto.RegisterType((*StreamingResponse)(nil), "go.micro.cs.service.upload.StreamingResponse")
}

func init() {
	proto.RegisterFile("proto/upload/upload.proto", fileDescriptor_5508b024e4d0ced2)
}

var fileDescriptor_5508b024e4d0ced2 = []byte{
	// 324 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0x3f, 0x4f, 0xc3, 0x30,
	0x10, 0xc5, 0x71, 0xff, 0xa6, 0x27, 0x84, 0x84, 0xa7, 0x10, 0x04, 0x94, 0x2c, 0x74, 0xc1, 0x45,
	0x45, 0x65, 0x2f, 0x03, 0x52, 0x17, 0x86, 0x20, 0xc4, 0x58, 0xb9, 0xf1, 0x35, 0x58, 0x8a, 0xed,
	0x28, 0x76, 0x91, 0xe0, 0xf3, 0xf2, 0x41, 0x50, 0xed, 0x94, 0x20, 0x21, 0x98, 0x92, 0x7b, 0xbf,
	0x7b, 0xcf, 0xba, 0x3b, 0x38, 0xa9, 0x6a, 0xe3, 0xcc, 0x74, 0x5b, 0x95, 0x86, 0x8b, 0xe6, 0xc3,
	0xbc, 0x46, 0x93, 0xc2, 0x30, 0x25, 0xf3, 0xda, 0xb0, 0xdc, 0x32, 0x8b, 0xf5, 0x9b, 0xcc, 0x91,
	0x85, 0x8e, 0xe4, 0xa2, 0x30, 0xa6, 0x28, 0x71, 0xea, 0x3b, 0xd7, 0xdb, 0xcd, 0xd4, 0x49, 0x85,
	0xd6, 0x71, 0x55, 0x05, 0x73, 0x3a, 0x87, 0xfe, 0xfd, 0xbb, 0x43, 0x4b, 0x63, 0x18, 0xe6, 0x46,
	0x3b, 0xd4, 0x2e, 0x26, 0x63, 0x32, 0x39, 0xcc, 0xf6, 0x25, 0xa5, 0xd0, 0xb3, 0xf2, 0x03, 0xe3,
	0xce, 0x98, 0x4c, 0xba, 0x99, 0xff, 0x4f, 0x3f, 0x09, 0x44, 0x0f, 0xb2, 0xc4, 0xa5, 0xde, 0x18,
	0x7a, 0x04, 0x1d, 0x29, 0xbc, 0xab, 0x9b, 0x75, 0xa4, 0xa0, 0xa7, 0x30, 0xda, 0xc8, 0x12, 0x57,
	0x9a, 0xab, 0xe0, 0x1a, 0x65, 0xd1, 0x4e, 0x78, 0xe4, 0x0a, 0xbf, 0xd3, 0xba, 0x6d, 0x1a, 0x3d,
	0x07, 0xd8, 0x71, 0xfb, 0xca, 0x67, 0xf3, 0xbb, 0xb8, 0xe7, 0x1d, 0x3f, 0x14, 0x9a, 0x40, 0x54,
	0x9a, 0x9c, 0x3b, 0x69, 0x74, 0xdc, 0x0f, 0x79, 0xfb, 0x9a, 0x9e, 0x01, 0xe4, 0x35, 0x72, 0x87,
	0x62, 0xc5, 0x5d, 0x3c, 0xf0, 0xa9, 0xa3, 0x46, 0x59, 0xb8, 0x1d, 0xde, 0x56, 0x62, 0x8f, 0x87,
	0x01, 0x37, 0x4a, 0xc0, 0x02, 0x4b, 0x6c, 0x70, 0x14, 0x70, 0xa3, 0x2c, 0x5c, 0x7a, 0x05, 0xc7,
	0x4f, 0xae, 0x46, 0xae, 0xa4, 0x2e, 0x32, 0xb4, 0x95, 0xd1, 0xb6, 0x9d, 0x80, 0xb4, 0x13, 0xcc,
	0x34, 0x0c, 0x9e, 0xfd, 0xc6, 0xa9, 0x00, 0x78, 0xa9, 0xa5, 0xc3, 0xa5, 0xe2, 0x05, 0xd2, 0x4b,
	0xf6, 0xf7, 0x71, 0x98, 0x5f, 0x7c, 0x72, 0xfd, 0x5f, 0xcb, 0xaf, 0xd7, 0xd3, 0x83, 0x09, 0xb9,
	0x21, 0xeb, 0x81, 0xbf, 0xde, 0xed, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc0, 0x81, 0x3f, 0xee,
	0x17, 0x02, 0x00, 0x00,
}
