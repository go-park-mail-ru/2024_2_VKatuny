// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package compressmicroservice

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CompressServiceClient is the client API for CompressService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CompressServiceClient interface {
	// rpc CompressAndSaveFile (CompressAndSaveFileInput) returns (Nothing) {}
	// rpc DeleteFile (DeleteFileInput) returns (Nothing) {}
	StartScanCompressDemon(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error)
}

type compressServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCompressServiceClient(cc grpc.ClientConnInterface) CompressServiceClient {
	return &compressServiceClient{cc}
}

func (c *compressServiceClient) StartScanCompressDemon(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := c.cc.Invoke(ctx, "/compressmicroservice.CompressService/StartScanCompressDemon", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CompressServiceServer is the server API for CompressService service.
// All implementations must embed UnimplementedCompressServiceServer
// for forward compatibility
type CompressServiceServer interface {
	// rpc CompressAndSaveFile (CompressAndSaveFileInput) returns (Nothing) {}
	// rpc DeleteFile (DeleteFileInput) returns (Nothing) {}
	StartScanCompressDemon(context.Context, *Nothing) (*Nothing, error)
	mustEmbedUnimplementedCompressServiceServer()
}

// UnimplementedCompressServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCompressServiceServer struct {
}

func (UnimplementedCompressServiceServer) StartScanCompressDemon(context.Context, *Nothing) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartScanCompressDemon not implemented")
}
func (UnimplementedCompressServiceServer) mustEmbedUnimplementedCompressServiceServer() {}

// UnsafeCompressServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CompressServiceServer will
// result in compilation errors.
type UnsafeCompressServiceServer interface {
	mustEmbedUnimplementedCompressServiceServer()
}

func RegisterCompressServiceServer(s grpc.ServiceRegistrar, srv CompressServiceServer) {
	s.RegisterService(&CompressService_ServiceDesc, srv)
}

func _CompressService_StartScanCompressDemon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nothing)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompressServiceServer).StartScanCompressDemon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/compressmicroservice.CompressService/StartScanCompressDemon",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompressServiceServer).StartScanCompressDemon(ctx, req.(*Nothing))
	}
	return interceptor(ctx, in, info, handler)
}

// CompressService_ServiceDesc is the grpc.ServiceDesc for CompressService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CompressService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "compressmicroservice.CompressService",
	HandlerType: (*CompressServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartScanCompressDemon",
			Handler:    _CompressService_StartScanCompressDemon_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "compress.proto",
}
