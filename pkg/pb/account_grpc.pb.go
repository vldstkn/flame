// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: account.proto

package pb

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
	Account_Register_FullMethodName         = "/Account/Register"
	Account_Login_FullMethodName            = "/Account/Login"
	Account_GetTokens_FullMethodName        = "/Account/GetTokens"
	Account_UpdateProfile_FullMethodName    = "/Account/UpdateProfile"
	Account_GetProfile_FullMethodName       = "/Account/GetProfile"
	Account_UploadPhoto_FullMethodName      = "/Account/UploadPhoto"
	Account_DeletePhoto_FullMethodName      = "/Account/DeletePhoto"
	Account_GetMatchingUsers_FullMethodName = "/Account/GetMatchingUsers"
)

// AccountClient is the client API for Account service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountClient interface {
	Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterRes, error)
	Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginRes, error)
	GetTokens(ctx context.Context, in *GetTokensReq, opts ...grpc.CallOption) (*GetTokensRes, error)
	UpdateProfile(ctx context.Context, in *UpdateProfileReq, opts ...grpc.CallOption) (*UpdateProfileRes, error)
	GetProfile(ctx context.Context, in *GetProfileReq, opts ...grpc.CallOption) (*GetProfileRes, error)
	UploadPhoto(ctx context.Context, in *UploadPhotoReq, opts ...grpc.CallOption) (*UploadPhotoRes, error)
	DeletePhoto(ctx context.Context, in *DeletePhotoReq, opts ...grpc.CallOption) (*DeletePhotoRes, error)
	GetMatchingUsers(ctx context.Context, in *GetMatchingUsersReq, opts ...grpc.CallOption) (*GetMatchingUsersRes, error)
}

type accountClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountClient(cc grpc.ClientConnInterface) AccountClient {
	return &accountClient{cc}
}

func (c *accountClient) Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterRes)
	err := c.cc.Invoke(ctx, Account_Register_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginRes)
	err := c.cc.Invoke(ctx, Account_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) GetTokens(ctx context.Context, in *GetTokensReq, opts ...grpc.CallOption) (*GetTokensRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTokensRes)
	err := c.cc.Invoke(ctx, Account_GetTokens_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) UpdateProfile(ctx context.Context, in *UpdateProfileReq, opts ...grpc.CallOption) (*UpdateProfileRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateProfileRes)
	err := c.cc.Invoke(ctx, Account_UpdateProfile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) GetProfile(ctx context.Context, in *GetProfileReq, opts ...grpc.CallOption) (*GetProfileRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetProfileRes)
	err := c.cc.Invoke(ctx, Account_GetProfile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) UploadPhoto(ctx context.Context, in *UploadPhotoReq, opts ...grpc.CallOption) (*UploadPhotoRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UploadPhotoRes)
	err := c.cc.Invoke(ctx, Account_UploadPhoto_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) DeletePhoto(ctx context.Context, in *DeletePhotoReq, opts ...grpc.CallOption) (*DeletePhotoRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeletePhotoRes)
	err := c.cc.Invoke(ctx, Account_DeletePhoto_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountClient) GetMatchingUsers(ctx context.Context, in *GetMatchingUsersReq, opts ...grpc.CallOption) (*GetMatchingUsersRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetMatchingUsersRes)
	err := c.cc.Invoke(ctx, Account_GetMatchingUsers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountServer is the server API for Account service.
// All implementations must embed UnimplementedAccountServer
// for forward compatibility.
type AccountServer interface {
	Register(context.Context, *RegisterReq) (*RegisterRes, error)
	Login(context.Context, *LoginReq) (*LoginRes, error)
	GetTokens(context.Context, *GetTokensReq) (*GetTokensRes, error)
	UpdateProfile(context.Context, *UpdateProfileReq) (*UpdateProfileRes, error)
	GetProfile(context.Context, *GetProfileReq) (*GetProfileRes, error)
	UploadPhoto(context.Context, *UploadPhotoReq) (*UploadPhotoRes, error)
	DeletePhoto(context.Context, *DeletePhotoReq) (*DeletePhotoRes, error)
	GetMatchingUsers(context.Context, *GetMatchingUsersReq) (*GetMatchingUsersRes, error)
	mustEmbedUnimplementedAccountServer()
}

// UnimplementedAccountServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAccountServer struct{}

func (UnimplementedAccountServer) Register(context.Context, *RegisterReq) (*RegisterRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedAccountServer) Login(context.Context, *LoginReq) (*LoginRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAccountServer) GetTokens(context.Context, *GetTokensReq) (*GetTokensRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTokens not implemented")
}
func (UnimplementedAccountServer) UpdateProfile(context.Context, *UpdateProfileReq) (*UpdateProfileRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProfile not implemented")
}
func (UnimplementedAccountServer) GetProfile(context.Context, *GetProfileReq) (*GetProfileRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfile not implemented")
}
func (UnimplementedAccountServer) UploadPhoto(context.Context, *UploadPhotoReq) (*UploadPhotoRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadPhoto not implemented")
}
func (UnimplementedAccountServer) DeletePhoto(context.Context, *DeletePhotoReq) (*DeletePhotoRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePhoto not implemented")
}
func (UnimplementedAccountServer) GetMatchingUsers(context.Context, *GetMatchingUsersReq) (*GetMatchingUsersRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMatchingUsers not implemented")
}
func (UnimplementedAccountServer) mustEmbedUnimplementedAccountServer() {}
func (UnimplementedAccountServer) testEmbeddedByValue()                 {}

// UnsafeAccountServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountServer will
// result in compilation errors.
type UnsafeAccountServer interface {
	mustEmbedUnimplementedAccountServer()
}

func RegisterAccountServer(s grpc.ServiceRegistrar, srv AccountServer) {
	// If the following call pancis, it indicates UnimplementedAccountServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Account_ServiceDesc, srv)
}

func _Account_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Account_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).Register(ctx, req.(*RegisterReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Account_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).Login(ctx, req.(*LoginReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_GetTokens_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTokensReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).GetTokens(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Account_GetTokens_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).GetTokens(ctx, req.(*GetTokensReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_UpdateProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProfileReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).UpdateProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Account_UpdateProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).UpdateProfile(ctx, req.(*UpdateProfileReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_GetProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProfileReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).GetProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Account_GetProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).GetProfile(ctx, req.(*GetProfileReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_UploadPhoto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadPhotoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).UploadPhoto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Account_UploadPhoto_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).UploadPhoto(ctx, req.(*UploadPhotoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_DeletePhoto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePhotoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).DeletePhoto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Account_DeletePhoto_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).DeletePhoto(ctx, req.(*DeletePhotoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Account_GetMatchingUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMatchingUsersReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).GetMatchingUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Account_GetMatchingUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).GetMatchingUsers(ctx, req.(*GetMatchingUsersReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Account_ServiceDesc is the grpc.ServiceDesc for Account service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Account_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Account",
	HandlerType: (*AccountServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Account_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Account_Login_Handler,
		},
		{
			MethodName: "GetTokens",
			Handler:    _Account_GetTokens_Handler,
		},
		{
			MethodName: "UpdateProfile",
			Handler:    _Account_UpdateProfile_Handler,
		},
		{
			MethodName: "GetProfile",
			Handler:    _Account_GetProfile_Handler,
		},
		{
			MethodName: "UploadPhoto",
			Handler:    _Account_UploadPhoto_Handler,
		},
		{
			MethodName: "DeletePhoto",
			Handler:    _Account_DeletePhoto_Handler,
		},
		{
			MethodName: "GetMatchingUsers",
			Handler:    _Account_GetMatchingUsers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "account.proto",
}
