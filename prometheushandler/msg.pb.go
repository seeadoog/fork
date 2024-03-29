// Code generated by protoc-gen-go. DO NOT EDIT.
// source: msg.proto

/*
Package prometheushandler is a generated protocol buffer package.

It is generated from these files:
	msg.proto

It has these top-level messages:
	Request
	Response
*/
package prometheushandler

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Request struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Age  int32  `protobuf:"varint,2,opt,name=age" json:"age,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Request) GetAge() int32 {
	if m != nil {
		return m.Age
	}
	return 0
}

type Response struct {
	Code int32  `protobuf:"varint,1,opt,name=code" json:"code,omitempty"`
	Sid  string `protobuf:"bytes,2,opt,name=sid" json:"sid,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Response) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Response) GetSid() string {
	if m != nil {
		return m.Sid
	}
	return ""
}

func init() {
	proto.RegisterType((*Request)(nil), "prometheushandler.request")
	proto.RegisterType((*Response)(nil), "prometheushandler.response")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for RpcServer service

type RpcServerClient interface {
	CreateUser(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type rpcServerClient struct {
	cc *grpc.ClientConn
}

func NewRpcServerClient(cc *grpc.ClientConn) RpcServerClient {
	return &rpcServerClient{cc}
}

func (c *rpcServerClient) CreateUser(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := grpc.Invoke(ctx, "/prometheushandler.rpc_server/CreateUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RpcServer service

type RpcServerServer interface {
	CreateUser(context.Context, *Request) (*Response, error)
}

func RegisterRpcServerServer(s *grpc.Server, srv RpcServerServer) {
	s.RegisterService(&_RpcServer_serviceDesc, srv)
}

func _RpcServer_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcServerServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/prometheushandler.rpc_server/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcServerServer).CreateUser(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _RpcServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "prometheushandler.rpc_server",
	HandlerType: (*RpcServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _RpcServer_CreateUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "msg.proto",
}

func init() { proto.RegisterFile("msg.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 174 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x8f, 0x31, 0x0f, 0x82, 0x30,
	0x10, 0x46, 0x45, 0x45, 0xe5, 0x26, 0xed, 0x44, 0x70, 0x21, 0x4c, 0x4c, 0x68, 0xf4, 0x27, 0x38,
	0xb8, 0x37, 0x61, 0x36, 0x15, 0x2e, 0x60, 0x22, 0x6d, 0xbd, 0x2b, 0xfe, 0x7e, 0xd3, 0xe2, 0xa6,
	0xdb, 0x4b, 0x2e, 0x2f, 0xf7, 0x3e, 0x48, 0x06, 0xee, 0x2a, 0x4b, 0xc6, 0x19, 0xb1, 0xb3, 0x64,
	0x06, 0x74, 0x3d, 0x8e, 0xdc, 0x2b, 0xdd, 0x3e, 0x91, 0x8a, 0x03, 0xac, 0x09, 0x5f, 0x23, 0xb2,
	0x13, 0x02, 0x96, 0x5a, 0x0d, 0x98, 0x46, 0x79, 0x54, 0x26, 0x32, 0xb0, 0xd8, 0xc2, 0x42, 0x75,
	0x98, 0xce, 0xf3, 0xa8, 0x8c, 0xa5, 0xc7, 0xe2, 0x08, 0x1b, 0x42, 0xb6, 0x46, 0x33, 0x7a, 0xa3,
	0x31, 0xed, 0x64, 0xc4, 0x32, 0xb0, 0x37, 0xf8, 0xd1, 0x06, 0x23, 0x91, 0x1e, 0x4f, 0x35, 0x00,
	0xd9, 0xe6, 0xc6, 0x48, 0x6f, 0x24, 0x71, 0x05, 0xb8, 0x10, 0x2a, 0x87, 0x35, 0x23, 0x89, 0xac,
	0xfa, 0x49, 0xaa, 0xbe, 0x3d, 0xd9, 0xfe, 0xef, 0x6d, 0x7a, 0x5d, 0xcc, 0xee, 0xab, 0xb0, 0xe9,
	0xfc, 0x09, 0x00, 0x00, 0xff, 0xff, 0x91, 0x9b, 0x83, 0x72, 0xe0, 0x00, 0x00, 0x00,
}
