// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package mgrpc

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

// DomainSerciceClient is the client API for DomainSercice service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DomainSerciceClient interface {
	ListDomains(ctx context.Context, in *EmptyParam, opts ...grpc.CallOption) (*DomainListResp, error)
	DeployDomain(ctx context.Context, in *DomainCreationParam, opts ...grpc.CallOption) (*DomainCreationResp, error)
}

type domainSerciceClient struct {
	cc grpc.ClientConnInterface
}

func NewDomainSerciceClient(cc grpc.ClientConnInterface) DomainSerciceClient {
	return &domainSerciceClient{cc}
}

func (c *domainSerciceClient) ListDomains(ctx context.Context, in *EmptyParam, opts ...grpc.CallOption) (*DomainListResp, error) {
	out := new(DomainListResp)
	err := c.cc.Invoke(ctx, "/DomainSercice/ListDomains", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *domainSerciceClient) DeployDomain(ctx context.Context, in *DomainCreationParam, opts ...grpc.CallOption) (*DomainCreationResp, error) {
	out := new(DomainCreationResp)
	err := c.cc.Invoke(ctx, "/DomainSercice/DeployDomain", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DomainSerciceServer is the server API for DomainSercice service.
// All implementations must embed UnimplementedDomainSerciceServer
// for forward compatibility
type DomainSerciceServer interface {
	ListDomains(context.Context, *EmptyParam) (*DomainListResp, error)
	DeployDomain(context.Context, *DomainCreationParam) (*DomainCreationResp, error)
	mustEmbedUnimplementedDomainSerciceServer()
}

// UnimplementedDomainSerciceServer must be embedded to have forward compatible implementations.
type UnimplementedDomainSerciceServer struct {
}

func (UnimplementedDomainSerciceServer) ListDomains(context.Context, *EmptyParam) (*DomainListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDomains not implemented")
}
func (UnimplementedDomainSerciceServer) DeployDomain(context.Context, *DomainCreationParam) (*DomainCreationResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeployDomain not implemented")
}
func (UnimplementedDomainSerciceServer) mustEmbedUnimplementedDomainSerciceServer() {}

// UnsafeDomainSerciceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DomainSerciceServer will
// result in compilation errors.
type UnsafeDomainSerciceServer interface {
	mustEmbedUnimplementedDomainSerciceServer()
}

func RegisterDomainSerciceServer(s grpc.ServiceRegistrar, srv DomainSerciceServer) {
	s.RegisterService(&DomainSercice_ServiceDesc, srv)
}

func _DomainSercice_ListDomains_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DomainSerciceServer).ListDomains(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/DomainSercice/ListDomains",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DomainSerciceServer).ListDomains(ctx, req.(*EmptyParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _DomainSercice_DeployDomain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DomainCreationParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DomainSerciceServer).DeployDomain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/DomainSercice/DeployDomain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DomainSerciceServer).DeployDomain(ctx, req.(*DomainCreationParam))
	}
	return interceptor(ctx, in, info, handler)
}

// DomainSercice_ServiceDesc is the grpc.ServiceDesc for DomainSercice service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DomainSercice_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "DomainSercice",
	HandlerType: (*DomainSerciceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListDomains",
			Handler:    _DomainSercice_ListDomains_Handler,
		},
		{
			MethodName: "DeployDomain",
			Handler:    _DomainSercice_DeployDomain_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "domain.proto",
}
