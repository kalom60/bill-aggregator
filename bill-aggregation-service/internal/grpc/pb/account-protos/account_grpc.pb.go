// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.0
// source: account.proto

package account_protos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	AccountService_GetLinkedAccounts_FullMethodName = "/account.AccountService/GetLinkedAccounts"
)

// AccountServiceClient is the client API for AccountService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountServiceClient interface {
	GetLinkedAccounts(ctx context.Context, in *AccountRequest, opts ...grpc.CallOption) (*AccountResponse, error)
}

type accountServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountServiceClient(cc grpc.ClientConnInterface) AccountServiceClient {
	return &accountServiceClient{cc}
}

func (c *accountServiceClient) GetLinkedAccounts(ctx context.Context, in *AccountRequest, opts ...grpc.CallOption) (*AccountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AccountResponse)
	err := c.cc.Invoke(ctx, AccountService_GetLinkedAccounts_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountServiceServer is the server API for AccountService service.
// All implementations must embed UnimplementedAccountServiceServer
// for forward compatibility.
type AccountServiceServer interface {
	GetLinkedAccounts(context.Context, *AccountRequest) (*AccountResponse, error)
	mustEmbedUnimplementedAccountServiceServer()
}

// UnimplementedAccountServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAccountServiceServer struct{}

func (UnimplementedAccountServiceServer) GetLinkedAccounts(context.Context, *AccountRequest) (*AccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLinkedAccounts not implemented")
}
func (UnimplementedAccountServiceServer) mustEmbedUnimplementedAccountServiceServer() {}
func (UnimplementedAccountServiceServer) testEmbeddedByValue()                        {}

// UnsafeAccountServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountServiceServer will
// result in compilation errors.
type UnsafeAccountServiceServer interface {
	mustEmbedUnimplementedAccountServiceServer()
}

func RegisterAccountServiceServer(s grpc.ServiceRegistrar, srv AccountServiceServer) {
	// If the following call pancis, it indicates UnimplementedAccountServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AccountService_ServiceDesc, srv)
}

func _AccountService_GetLinkedAccounts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetLinkedAccounts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_GetLinkedAccounts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetLinkedAccounts(ctx, req.(*AccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccountService_ServiceDesc is the grpc.ServiceDesc for AccountService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccountService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "account.AccountService",
	HandlerType: (*AccountServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLinkedAccounts",
			Handler:    _AccountService_GetLinkedAccounts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "account.proto",
}
