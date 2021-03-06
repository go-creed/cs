// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/upload/upload.proto

package go_micro_cs_service_upload

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
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

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Upload service

func NewUploadEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Upload service

type UploadService interface {
	WriteBytes(ctx context.Context, opts ...client.CallOption) (Upload_WriteBytesService, error)
	FileDetail(ctx context.Context, in *FileMate, opts ...client.CallOption) (*FileMate, error)
	FileChunk(ctx context.Context, in *ChunkRequest, opts ...client.CallOption) (*ChunkResponse, error)
	FileMerge(ctx context.Context, in *MergeRequest, opts ...client.CallOption) (*FileMate, error)
	FileChunkVerify(ctx context.Context, in *ChunkRequest, opts ...client.CallOption) (*ChunkResponse, error)
}

type uploadService struct {
	c    client.Client
	name string
}

func NewUploadService(name string, c client.Client) UploadService {
	return &uploadService{
		c:    c,
		name: name,
	}
}

func (c *uploadService) WriteBytes(ctx context.Context, opts ...client.CallOption) (Upload_WriteBytesService, error) {
	req := c.c.NewRequest(c.name, "Upload.WriteBytes", &Bytes{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &uploadServiceWriteBytes{stream}, nil
}

type Upload_WriteBytesService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Bytes) error
	Recv() (*StreamingResponse, error)
}

type uploadServiceWriteBytes struct {
	stream client.Stream
}

func (x *uploadServiceWriteBytes) Close() error {
	return x.stream.Close()
}

func (x *uploadServiceWriteBytes) Context() context.Context {
	return x.stream.Context()
}

func (x *uploadServiceWriteBytes) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *uploadServiceWriteBytes) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *uploadServiceWriteBytes) Send(m *Bytes) error {
	return x.stream.Send(m)
}

func (x *uploadServiceWriteBytes) Recv() (*StreamingResponse, error) {
	m := new(StreamingResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *uploadService) FileDetail(ctx context.Context, in *FileMate, opts ...client.CallOption) (*FileMate, error) {
	req := c.c.NewRequest(c.name, "Upload.FileDetail", in)
	out := new(FileMate)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uploadService) FileChunk(ctx context.Context, in *ChunkRequest, opts ...client.CallOption) (*ChunkResponse, error) {
	req := c.c.NewRequest(c.name, "Upload.FileChunk", in)
	out := new(ChunkResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uploadService) FileMerge(ctx context.Context, in *MergeRequest, opts ...client.CallOption) (*FileMate, error) {
	req := c.c.NewRequest(c.name, "Upload.FileMerge", in)
	out := new(FileMate)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uploadService) FileChunkVerify(ctx context.Context, in *ChunkRequest, opts ...client.CallOption) (*ChunkResponse, error) {
	req := c.c.NewRequest(c.name, "Upload.FileChunkVerify", in)
	out := new(ChunkResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Upload service

type UploadHandler interface {
	WriteBytes(context.Context, Upload_WriteBytesStream) error
	FileDetail(context.Context, *FileMate, *FileMate) error
	FileChunk(context.Context, *ChunkRequest, *ChunkResponse) error
	FileMerge(context.Context, *MergeRequest, *FileMate) error
	FileChunkVerify(context.Context, *ChunkRequest, *ChunkResponse) error
}

func RegisterUploadHandler(s server.Server, hdlr UploadHandler, opts ...server.HandlerOption) error {
	type upload interface {
		WriteBytes(ctx context.Context, stream server.Stream) error
		FileDetail(ctx context.Context, in *FileMate, out *FileMate) error
		FileChunk(ctx context.Context, in *ChunkRequest, out *ChunkResponse) error
		FileMerge(ctx context.Context, in *MergeRequest, out *FileMate) error
		FileChunkVerify(ctx context.Context, in *ChunkRequest, out *ChunkResponse) error
	}
	type Upload struct {
		upload
	}
	h := &uploadHandler{hdlr}
	return s.Handle(s.NewHandler(&Upload{h}, opts...))
}

type uploadHandler struct {
	UploadHandler
}

func (h *uploadHandler) WriteBytes(ctx context.Context, stream server.Stream) error {
	return h.UploadHandler.WriteBytes(ctx, &uploadWriteBytesStream{stream})
}

type Upload_WriteBytesStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamingResponse) error
	Recv() (*Bytes, error)
}

type uploadWriteBytesStream struct {
	stream server.Stream
}

func (x *uploadWriteBytesStream) Close() error {
	return x.stream.Close()
}

func (x *uploadWriteBytesStream) Context() context.Context {
	return x.stream.Context()
}

func (x *uploadWriteBytesStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *uploadWriteBytesStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *uploadWriteBytesStream) Send(m *StreamingResponse) error {
	return x.stream.Send(m)
}

func (x *uploadWriteBytesStream) Recv() (*Bytes, error) {
	m := new(Bytes)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (h *uploadHandler) FileDetail(ctx context.Context, in *FileMate, out *FileMate) error {
	return h.UploadHandler.FileDetail(ctx, in, out)
}

func (h *uploadHandler) FileChunk(ctx context.Context, in *ChunkRequest, out *ChunkResponse) error {
	return h.UploadHandler.FileChunk(ctx, in, out)
}

func (h *uploadHandler) FileMerge(ctx context.Context, in *MergeRequest, out *FileMate) error {
	return h.UploadHandler.FileMerge(ctx, in, out)
}

func (h *uploadHandler) FileChunkVerify(ctx context.Context, in *ChunkRequest, out *ChunkResponse) error {
	return h.UploadHandler.FileChunkVerify(ctx, in, out)
}
