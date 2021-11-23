// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package apiv1

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

// BasicPluginServiceClient is the client API for BasicPluginService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BasicPluginServiceClient interface {
	Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*InitResponse, error)
	Start(ctx context.Context, in *StartRequest, opts ...grpc.CallOption) (*StartResponse, error)
	ProjectInit(ctx context.Context, in *ProjectInitRequest, opts ...grpc.CallOption) (*ProjectInitResponse, error)
}

type basicPluginServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBasicPluginServiceClient(cc grpc.ClientConnInterface) BasicPluginServiceClient {
	return &basicPluginServiceClient{cc}
}

func (c *basicPluginServiceClient) Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*InitResponse, error) {
	out := new(InitResponse)
	err := c.cc.Invoke(ctx, "/api.v1.BasicPluginService/Init", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *basicPluginServiceClient) Start(ctx context.Context, in *StartRequest, opts ...grpc.CallOption) (*StartResponse, error) {
	out := new(StartResponse)
	err := c.cc.Invoke(ctx, "/api.v1.BasicPluginService/Start", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *basicPluginServiceClient) ProjectInit(ctx context.Context, in *ProjectInitRequest, opts ...grpc.CallOption) (*ProjectInitResponse, error) {
	out := new(ProjectInitResponse)
	err := c.cc.Invoke(ctx, "/api.v1.BasicPluginService/ProjectInit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BasicPluginServiceServer is the server API for BasicPluginService service.
// All implementations should embed UnimplementedBasicPluginServiceServer
// for forward compatibility
type BasicPluginServiceServer interface {
	Init(context.Context, *InitRequest) (*InitResponse, error)
	Start(context.Context, *StartRequest) (*StartResponse, error)
	ProjectInit(context.Context, *ProjectInitRequest) (*ProjectInitResponse, error)
}

// UnimplementedBasicPluginServiceServer should be embedded to have forward compatible implementations.
type UnimplementedBasicPluginServiceServer struct {
}

func (UnimplementedBasicPluginServiceServer) Init(context.Context, *InitRequest) (*InitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Init not implemented")
}
func (UnimplementedBasicPluginServiceServer) Start(context.Context, *StartRequest) (*StartResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Start not implemented")
}
func (UnimplementedBasicPluginServiceServer) ProjectInit(context.Context, *ProjectInitRequest) (*ProjectInitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProjectInit not implemented")
}

// UnsafeBasicPluginServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BasicPluginServiceServer will
// result in compilation errors.
type UnsafeBasicPluginServiceServer interface {
	mustEmbedUnimplementedBasicPluginServiceServer()
}

func RegisterBasicPluginServiceServer(s grpc.ServiceRegistrar, srv BasicPluginServiceServer) {
	s.RegisterService(&BasicPluginService_ServiceDesc, srv)
}

func _BasicPluginService_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicPluginServiceServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.v1.BasicPluginService/Init",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicPluginServiceServer).Init(ctx, req.(*InitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BasicPluginService_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicPluginServiceServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.v1.BasicPluginService/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicPluginServiceServer).Start(ctx, req.(*StartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BasicPluginService_ProjectInit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectInitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicPluginServiceServer).ProjectInit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.v1.BasicPluginService/ProjectInit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicPluginServiceServer).ProjectInit(ctx, req.(*ProjectInitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BasicPluginService_ServiceDesc is the grpc.ServiceDesc for BasicPluginService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BasicPluginService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.BasicPluginService",
	HandlerType: (*BasicPluginServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Init",
			Handler:    _BasicPluginService_Init_Handler,
		},
		{
			MethodName: "Start",
			Handler:    _BasicPluginService_Start_Handler,
		},
		{
			MethodName: "ProjectInit",
			Handler:    _BasicPluginService_ProjectInit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/plugin.proto",
}

// StatePluginServiceClient is the client API for StatePluginService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StatePluginServiceClient interface {
	GetState(ctx context.Context, in *GetStateRequest, opts ...grpc.CallOption) (StatePluginService_GetStateClient, error)
	SaveState(ctx context.Context, in *SaveStateRequest, opts ...grpc.CallOption) (*SaveStateResponse, error)
	ReleaseStateLock(ctx context.Context, in *ReleaseStateLockRequest, opts ...grpc.CallOption) (*ReleaseStateLockResponse, error)
}

type statePluginServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStatePluginServiceClient(cc grpc.ClientConnInterface) StatePluginServiceClient {
	return &statePluginServiceClient{cc}
}

func (c *statePluginServiceClient) GetState(ctx context.Context, in *GetStateRequest, opts ...grpc.CallOption) (StatePluginService_GetStateClient, error) {
	stream, err := c.cc.NewStream(ctx, &StatePluginService_ServiceDesc.Streams[0], "/api.v1.StatePluginService/GetState", opts...)
	if err != nil {
		return nil, err
	}
	x := &statePluginServiceGetStateClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StatePluginService_GetStateClient interface {
	Recv() (*GetStateResponse, error)
	grpc.ClientStream
}

type statePluginServiceGetStateClient struct {
	grpc.ClientStream
}

func (x *statePluginServiceGetStateClient) Recv() (*GetStateResponse, error) {
	m := new(GetStateResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *statePluginServiceClient) SaveState(ctx context.Context, in *SaveStateRequest, opts ...grpc.CallOption) (*SaveStateResponse, error) {
	out := new(SaveStateResponse)
	err := c.cc.Invoke(ctx, "/api.v1.StatePluginService/SaveState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statePluginServiceClient) ReleaseStateLock(ctx context.Context, in *ReleaseStateLockRequest, opts ...grpc.CallOption) (*ReleaseStateLockResponse, error) {
	out := new(ReleaseStateLockResponse)
	err := c.cc.Invoke(ctx, "/api.v1.StatePluginService/ReleaseStateLock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatePluginServiceServer is the server API for StatePluginService service.
// All implementations should embed UnimplementedStatePluginServiceServer
// for forward compatibility
type StatePluginServiceServer interface {
	GetState(*GetStateRequest, StatePluginService_GetStateServer) error
	SaveState(context.Context, *SaveStateRequest) (*SaveStateResponse, error)
	ReleaseStateLock(context.Context, *ReleaseStateLockRequest) (*ReleaseStateLockResponse, error)
}

// UnimplementedStatePluginServiceServer should be embedded to have forward compatible implementations.
type UnimplementedStatePluginServiceServer struct {
}

func (UnimplementedStatePluginServiceServer) GetState(*GetStateRequest, StatePluginService_GetStateServer) error {
	return status.Errorf(codes.Unimplemented, "method GetState not implemented")
}
func (UnimplementedStatePluginServiceServer) SaveState(context.Context, *SaveStateRequest) (*SaveStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveState not implemented")
}
func (UnimplementedStatePluginServiceServer) ReleaseStateLock(context.Context, *ReleaseStateLockRequest) (*ReleaseStateLockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseStateLock not implemented")
}

// UnsafeStatePluginServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StatePluginServiceServer will
// result in compilation errors.
type UnsafeStatePluginServiceServer interface {
	mustEmbedUnimplementedStatePluginServiceServer()
}

func RegisterStatePluginServiceServer(s grpc.ServiceRegistrar, srv StatePluginServiceServer) {
	s.RegisterService(&StatePluginService_ServiceDesc, srv)
}

func _StatePluginService_GetState_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetStateRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StatePluginServiceServer).GetState(m, &statePluginServiceGetStateServer{stream})
}

type StatePluginService_GetStateServer interface {
	Send(*GetStateResponse) error
	grpc.ServerStream
}

type statePluginServiceGetStateServer struct {
	grpc.ServerStream
}

func (x *statePluginServiceGetStateServer) Send(m *GetStateResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _StatePluginService_SaveState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatePluginServiceServer).SaveState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.v1.StatePluginService/SaveState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatePluginServiceServer).SaveState(ctx, req.(*SaveStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatePluginService_ReleaseStateLock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReleaseStateLockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatePluginServiceServer).ReleaseStateLock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.v1.StatePluginService/ReleaseStateLock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatePluginServiceServer).ReleaseStateLock(ctx, req.(*ReleaseStateLockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// StatePluginService_ServiceDesc is the grpc.ServiceDesc for StatePluginService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StatePluginService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.StatePluginService",
	HandlerType: (*StatePluginServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SaveState",
			Handler:    _StatePluginService_SaveState_Handler,
		},
		{
			MethodName: "ReleaseStateLock",
			Handler:    _StatePluginService_ReleaseStateLock_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetState",
			Handler:       _StatePluginService_GetState_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/v1/plugin.proto",
}

// LockingPluginServiceClient is the client API for LockingPluginService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LockingPluginServiceClient interface {
	AcquireLocks(ctx context.Context, in *AcquireLocksRequest, opts ...grpc.CallOption) (*AcquireLocksResponse, error)
	ReleaseLocks(ctx context.Context, in *ReleaseLocksRequest, opts ...grpc.CallOption) (*ReleaseLocksResponse, error)
}

type lockingPluginServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLockingPluginServiceClient(cc grpc.ClientConnInterface) LockingPluginServiceClient {
	return &lockingPluginServiceClient{cc}
}

func (c *lockingPluginServiceClient) AcquireLocks(ctx context.Context, in *AcquireLocksRequest, opts ...grpc.CallOption) (*AcquireLocksResponse, error) {
	out := new(AcquireLocksResponse)
	err := c.cc.Invoke(ctx, "/api.v1.LockingPluginService/AcquireLocks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lockingPluginServiceClient) ReleaseLocks(ctx context.Context, in *ReleaseLocksRequest, opts ...grpc.CallOption) (*ReleaseLocksResponse, error) {
	out := new(ReleaseLocksResponse)
	err := c.cc.Invoke(ctx, "/api.v1.LockingPluginService/ReleaseLocks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LockingPluginServiceServer is the server API for LockingPluginService service.
// All implementations should embed UnimplementedLockingPluginServiceServer
// for forward compatibility
type LockingPluginServiceServer interface {
	AcquireLocks(context.Context, *AcquireLocksRequest) (*AcquireLocksResponse, error)
	ReleaseLocks(context.Context, *ReleaseLocksRequest) (*ReleaseLocksResponse, error)
}

// UnimplementedLockingPluginServiceServer should be embedded to have forward compatible implementations.
type UnimplementedLockingPluginServiceServer struct {
}

func (UnimplementedLockingPluginServiceServer) AcquireLocks(context.Context, *AcquireLocksRequest) (*AcquireLocksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcquireLocks not implemented")
}
func (UnimplementedLockingPluginServiceServer) ReleaseLocks(context.Context, *ReleaseLocksRequest) (*ReleaseLocksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseLocks not implemented")
}

// UnsafeLockingPluginServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LockingPluginServiceServer will
// result in compilation errors.
type UnsafeLockingPluginServiceServer interface {
	mustEmbedUnimplementedLockingPluginServiceServer()
}

func RegisterLockingPluginServiceServer(s grpc.ServiceRegistrar, srv LockingPluginServiceServer) {
	s.RegisterService(&LockingPluginService_ServiceDesc, srv)
}

func _LockingPluginService_AcquireLocks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcquireLocksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LockingPluginServiceServer).AcquireLocks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.v1.LockingPluginService/AcquireLocks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LockingPluginServiceServer).AcquireLocks(ctx, req.(*AcquireLocksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LockingPluginService_ReleaseLocks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReleaseLocksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LockingPluginServiceServer).ReleaseLocks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.v1.LockingPluginService/ReleaseLocks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LockingPluginServiceServer).ReleaseLocks(ctx, req.(*ReleaseLocksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LockingPluginService_ServiceDesc is the grpc.ServiceDesc for LockingPluginService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LockingPluginService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.LockingPluginService",
	HandlerType: (*LockingPluginServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AcquireLocks",
			Handler:    _LockingPluginService_AcquireLocks_Handler,
		},
		{
			MethodName: "ReleaseLocks",
			Handler:    _LockingPluginService_ReleaseLocks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/plugin.proto",
}

// DeployPluginServiceClient is the client API for DeployPluginService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeployPluginServiceClient interface {
	Plan(ctx context.Context, in *PlanRequest, opts ...grpc.CallOption) (*PlanResponse, error)
	Apply(ctx context.Context, in *ApplyRequest, opts ...grpc.CallOption) (DeployPluginService_ApplyClient, error)
}

type deployPluginServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDeployPluginServiceClient(cc grpc.ClientConnInterface) DeployPluginServiceClient {
	return &deployPluginServiceClient{cc}
}

func (c *deployPluginServiceClient) Plan(ctx context.Context, in *PlanRequest, opts ...grpc.CallOption) (*PlanResponse, error) {
	out := new(PlanResponse)
	err := c.cc.Invoke(ctx, "/api.v1.DeployPluginService/Plan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deployPluginServiceClient) Apply(ctx context.Context, in *ApplyRequest, opts ...grpc.CallOption) (DeployPluginService_ApplyClient, error) {
	stream, err := c.cc.NewStream(ctx, &DeployPluginService_ServiceDesc.Streams[0], "/api.v1.DeployPluginService/Apply", opts...)
	if err != nil {
		return nil, err
	}
	x := &deployPluginServiceApplyClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DeployPluginService_ApplyClient interface {
	Recv() (*ApplyResponse, error)
	grpc.ClientStream
}

type deployPluginServiceApplyClient struct {
	grpc.ClientStream
}

func (x *deployPluginServiceApplyClient) Recv() (*ApplyResponse, error) {
	m := new(ApplyResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DeployPluginServiceServer is the server API for DeployPluginService service.
// All implementations should embed UnimplementedDeployPluginServiceServer
// for forward compatibility
type DeployPluginServiceServer interface {
	Plan(context.Context, *PlanRequest) (*PlanResponse, error)
	Apply(*ApplyRequest, DeployPluginService_ApplyServer) error
}

// UnimplementedDeployPluginServiceServer should be embedded to have forward compatible implementations.
type UnimplementedDeployPluginServiceServer struct {
}

func (UnimplementedDeployPluginServiceServer) Plan(context.Context, *PlanRequest) (*PlanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Plan not implemented")
}
func (UnimplementedDeployPluginServiceServer) Apply(*ApplyRequest, DeployPluginService_ApplyServer) error {
	return status.Errorf(codes.Unimplemented, "method Apply not implemented")
}

// UnsafeDeployPluginServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeployPluginServiceServer will
// result in compilation errors.
type UnsafeDeployPluginServiceServer interface {
	mustEmbedUnimplementedDeployPluginServiceServer()
}

func RegisterDeployPluginServiceServer(s grpc.ServiceRegistrar, srv DeployPluginServiceServer) {
	s.RegisterService(&DeployPluginService_ServiceDesc, srv)
}

func _DeployPluginService_Plan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeployPluginServiceServer).Plan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.v1.DeployPluginService/Plan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeployPluginServiceServer).Plan(ctx, req.(*PlanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeployPluginService_Apply_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ApplyRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DeployPluginServiceServer).Apply(m, &deployPluginServiceApplyServer{stream})
}

type DeployPluginService_ApplyServer interface {
	Send(*ApplyResponse) error
	grpc.ServerStream
}

type deployPluginServiceApplyServer struct {
	grpc.ServerStream
}

func (x *deployPluginServiceApplyServer) Send(m *ApplyResponse) error {
	return x.ServerStream.SendMsg(m)
}

// DeployPluginService_ServiceDesc is the grpc.ServiceDesc for DeployPluginService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DeployPluginService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.DeployPluginService",
	HandlerType: (*DeployPluginServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Plan",
			Handler:    _DeployPluginService_Plan_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Apply",
			Handler:       _DeployPluginService_Apply_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/v1/plugin.proto",
}

// RunPluginServiceClient is the client API for RunPluginService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RunPluginServiceClient interface {
	Run(ctx context.Context, in *RunRequest, opts ...grpc.CallOption) (RunPluginService_RunClient, error)
}

type runPluginServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRunPluginServiceClient(cc grpc.ClientConnInterface) RunPluginServiceClient {
	return &runPluginServiceClient{cc}
}

func (c *runPluginServiceClient) Run(ctx context.Context, in *RunRequest, opts ...grpc.CallOption) (RunPluginService_RunClient, error) {
	stream, err := c.cc.NewStream(ctx, &RunPluginService_ServiceDesc.Streams[0], "/api.v1.RunPluginService/Run", opts...)
	if err != nil {
		return nil, err
	}
	x := &runPluginServiceRunClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RunPluginService_RunClient interface {
	Recv() (*RunResponse, error)
	grpc.ClientStream
}

type runPluginServiceRunClient struct {
	grpc.ClientStream
}

func (x *runPluginServiceRunClient) Recv() (*RunResponse, error) {
	m := new(RunResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RunPluginServiceServer is the server API for RunPluginService service.
// All implementations should embed UnimplementedRunPluginServiceServer
// for forward compatibility
type RunPluginServiceServer interface {
	Run(*RunRequest, RunPluginService_RunServer) error
}

// UnimplementedRunPluginServiceServer should be embedded to have forward compatible implementations.
type UnimplementedRunPluginServiceServer struct {
}

func (UnimplementedRunPluginServiceServer) Run(*RunRequest, RunPluginService_RunServer) error {
	return status.Errorf(codes.Unimplemented, "method Run not implemented")
}

// UnsafeRunPluginServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RunPluginServiceServer will
// result in compilation errors.
type UnsafeRunPluginServiceServer interface {
	mustEmbedUnimplementedRunPluginServiceServer()
}

func RegisterRunPluginServiceServer(s grpc.ServiceRegistrar, srv RunPluginServiceServer) {
	s.RegisterService(&RunPluginService_ServiceDesc, srv)
}

func _RunPluginService_Run_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RunRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RunPluginServiceServer).Run(m, &runPluginServiceRunServer{stream})
}

type RunPluginService_RunServer interface {
	Send(*RunResponse) error
	grpc.ServerStream
}

type runPluginServiceRunServer struct {
	grpc.ServerStream
}

func (x *runPluginServiceRunServer) Send(m *RunResponse) error {
	return x.ServerStream.SendMsg(m)
}

// RunPluginService_ServiceDesc is the grpc.ServiceDesc for RunPluginService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RunPluginService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.RunPluginService",
	HandlerType: (*RunPluginServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Run",
			Handler:       _RunPluginService_Run_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/v1/plugin.proto",
}

// CommandPluginServiceClient is the client API for CommandPluginService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommandPluginServiceClient interface {
	Command(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (CommandPluginService_CommandClient, error)
}

type commandPluginServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCommandPluginServiceClient(cc grpc.ClientConnInterface) CommandPluginServiceClient {
	return &commandPluginServiceClient{cc}
}

func (c *commandPluginServiceClient) Command(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (CommandPluginService_CommandClient, error) {
	stream, err := c.cc.NewStream(ctx, &CommandPluginService_ServiceDesc.Streams[0], "/api.v1.CommandPluginService/Command", opts...)
	if err != nil {
		return nil, err
	}
	x := &commandPluginServiceCommandClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CommandPluginService_CommandClient interface {
	Recv() (*CommandResponse, error)
	grpc.ClientStream
}

type commandPluginServiceCommandClient struct {
	grpc.ClientStream
}

func (x *commandPluginServiceCommandClient) Recv() (*CommandResponse, error) {
	m := new(CommandResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CommandPluginServiceServer is the server API for CommandPluginService service.
// All implementations should embed UnimplementedCommandPluginServiceServer
// for forward compatibility
type CommandPluginServiceServer interface {
	Command(*CommandRequest, CommandPluginService_CommandServer) error
}

// UnimplementedCommandPluginServiceServer should be embedded to have forward compatible implementations.
type UnimplementedCommandPluginServiceServer struct {
}

func (UnimplementedCommandPluginServiceServer) Command(*CommandRequest, CommandPluginService_CommandServer) error {
	return status.Errorf(codes.Unimplemented, "method Command not implemented")
}

// UnsafeCommandPluginServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommandPluginServiceServer will
// result in compilation errors.
type UnsafeCommandPluginServiceServer interface {
	mustEmbedUnimplementedCommandPluginServiceServer()
}

func RegisterCommandPluginServiceServer(s grpc.ServiceRegistrar, srv CommandPluginServiceServer) {
	s.RegisterService(&CommandPluginService_ServiceDesc, srv)
}

func _CommandPluginService_Command_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CommandRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CommandPluginServiceServer).Command(m, &commandPluginServiceCommandServer{stream})
}

type CommandPluginService_CommandServer interface {
	Send(*CommandResponse) error
	grpc.ServerStream
}

type commandPluginServiceCommandServer struct {
	grpc.ServerStream
}

func (x *commandPluginServiceCommandServer) Send(m *CommandResponse) error {
	return x.ServerStream.SendMsg(m)
}

// CommandPluginService_ServiceDesc is the grpc.ServiceDesc for CommandPluginService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CommandPluginService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.CommandPluginService",
	HandlerType: (*CommandPluginServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Command",
			Handler:       _CommandPluginService_Command_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/v1/plugin.proto",
}