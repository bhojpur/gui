// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// GuiServiceClient is the client API for GuiService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GuiServiceClient interface {
	// StartLocalEngine starts a GUI Engine on the Bhojpur.NET Platform directly.
	// The incoming requests are expected in the following order:
	//   1. metadata
	//   2. all bytes constituting the gui/config.yaml
	//   3. all bytes constituting the Engine YAML that will be executed (that the config.yaml points to)
	//   4. all bytes constituting the gzipped Bhojpur.NET Platform application tar stream
	//   5. the Bhojpur.NET Platform application tar stream done marker
	StartLocalEngine(ctx context.Context, opts ...grpc.CallOption) (GuiService_StartLocalEngineClient, error)
	// StartFromPreviousEngine starts a new Engine based on a previous one.
	// If the previous Engine does not have the can-replay condition set this call will result in an error.
	StartFromPreviousEngine(ctx context.Context, in *StartFromPreviousEngineRequest, opts ...grpc.CallOption) (*StartEngineResponse, error)
	// StartEngineRequest starts a new Engine based on its specification.
	StartEngine(ctx context.Context, in *StartEngineRequest, opts ...grpc.CallOption) (*StartEngineResponse, error)
	// Searches for Engine(s) known to this Engine
	ListEngines(ctx context.Context, in *ListEnginesRequest, opts ...grpc.CallOption) (*ListEnginesResponse, error)
	// Subscribe listens to new Engine(s) updates
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (GuiService_SubscribeClient, error)
	// GetEngine retrieves details of a single Engine
	GetEngine(ctx context.Context, in *GetEngineRequest, opts ...grpc.CallOption) (*GetEngineResponse, error)
	// Listen listens to Engine updates and log output of a running Engine
	Listen(ctx context.Context, in *ListenRequest, opts ...grpc.CallOption) (GuiService_ListenClient, error)
	// StopEngine stops a currently running Engine
	StopEngine(ctx context.Context, in *StopEngineRequest, opts ...grpc.CallOption) (*StopEngineResponse, error)
}

type guiServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGuiServiceClient(cc grpc.ClientConnInterface) GuiServiceClient {
	return &guiServiceClient{cc}
}

func (c *guiServiceClient) StartLocalEngine(ctx context.Context, opts ...grpc.CallOption) (GuiService_StartLocalEngineClient, error) {
	stream, err := c.cc.NewStream(ctx, &GuiService_ServiceDesc.Streams[0], "/v1.GuiService/StartLocalEngine", opts...)
	if err != nil {
		return nil, err
	}
	x := &guiServiceStartLocalEngineClient{stream}
	return x, nil
}

type GuiService_StartLocalEngineClient interface {
	Send(*StartLocalEngineRequest) error
	CloseAndRecv() (*StartEngineResponse, error)
	grpc.ClientStream
}

type guiServiceStartLocalEngineClient struct {
	grpc.ClientStream
}

func (x *guiServiceStartLocalEngineClient) Send(m *StartLocalEngineRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *guiServiceStartLocalEngineClient) CloseAndRecv() (*StartEngineResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(StartEngineResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *guiServiceClient) StartFromPreviousEngine(ctx context.Context, in *StartFromPreviousEngineRequest, opts ...grpc.CallOption) (*StartEngineResponse, error) {
	out := new(StartEngineResponse)
	err := c.cc.Invoke(ctx, "/v1.GuiService/StartFromPreviousEngine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *guiServiceClient) StartEngine(ctx context.Context, in *StartEngineRequest, opts ...grpc.CallOption) (*StartEngineResponse, error) {
	out := new(StartEngineResponse)
	err := c.cc.Invoke(ctx, "/v1.GuiService/StartEngine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *guiServiceClient) ListEngines(ctx context.Context, in *ListEnginesRequest, opts ...grpc.CallOption) (*ListEnginesResponse, error) {
	out := new(ListEnginesResponse)
	err := c.cc.Invoke(ctx, "/v1.GuiService/ListEngines", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *guiServiceClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (GuiService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &GuiService_ServiceDesc.Streams[1], "/v1.GuiService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &guiServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GuiService_SubscribeClient interface {
	Recv() (*SubscribeResponse, error)
	grpc.ClientStream
}

type guiServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *guiServiceSubscribeClient) Recv() (*SubscribeResponse, error) {
	m := new(SubscribeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *guiServiceClient) GetEngine(ctx context.Context, in *GetEngineRequest, opts ...grpc.CallOption) (*GetEngineResponse, error) {
	out := new(GetEngineResponse)
	err := c.cc.Invoke(ctx, "/v1.GuiService/GetEngine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *guiServiceClient) Listen(ctx context.Context, in *ListenRequest, opts ...grpc.CallOption) (GuiService_ListenClient, error) {
	stream, err := c.cc.NewStream(ctx, &GuiService_ServiceDesc.Streams[2], "/v1.GuiService/Listen", opts...)
	if err != nil {
		return nil, err
	}
	x := &guiServiceListenClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GuiService_ListenClient interface {
	Recv() (*ListenResponse, error)
	grpc.ClientStream
}

type guiServiceListenClient struct {
	grpc.ClientStream
}

func (x *guiServiceListenClient) Recv() (*ListenResponse, error) {
	m := new(ListenResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *guiServiceClient) StopEngine(ctx context.Context, in *StopEngineRequest, opts ...grpc.CallOption) (*StopEngineResponse, error) {
	out := new(StopEngineResponse)
	err := c.cc.Invoke(ctx, "/v1.GuiService/StopEngine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GuiServiceServer is the server API for GuiService service.
// All implementations must embed UnimplementedGuiServiceServer
// for forward compatibility
type GuiServiceServer interface {
	// StartLocalEngine starts a GUI Engine on the Bhojpur.NET Platform directly.
	// The incoming requests are expected in the following order:
	//   1. metadata
	//   2. all bytes constituting the gui/config.yaml
	//   3. all bytes constituting the Engine YAML that will be executed (that the config.yaml points to)
	//   4. all bytes constituting the gzipped Bhojpur.NET Platform application tar stream
	//   5. the Bhojpur.NET Platform application tar stream done marker
	StartLocalEngine(GuiService_StartLocalEngineServer) error
	// StartFromPreviousEngine starts a new Engine based on a previous one.
	// If the previous Engine does not have the can-replay condition set this call will result in an error.
	StartFromPreviousEngine(context.Context, *StartFromPreviousEngineRequest) (*StartEngineResponse, error)
	// StartEngineRequest starts a new Engine based on its specification.
	StartEngine(context.Context, *StartEngineRequest) (*StartEngineResponse, error)
	// Searches for Engine(s) known to this Engine
	ListEngines(context.Context, *ListEnginesRequest) (*ListEnginesResponse, error)
	// Subscribe listens to new Engine(s) updates
	Subscribe(*SubscribeRequest, GuiService_SubscribeServer) error
	// GetEngine retrieves details of a single Engine
	GetEngine(context.Context, *GetEngineRequest) (*GetEngineResponse, error)
	// Listen listens to Engine updates and log output of a running Engine
	Listen(*ListenRequest, GuiService_ListenServer) error
	// StopEngine stops a currently running Engine
	StopEngine(context.Context, *StopEngineRequest) (*StopEngineResponse, error)
	mustEmbedUnimplementedGuiServiceServer()
}

// UnimplementedGuiServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGuiServiceServer struct {
}

func (UnimplementedGuiServiceServer) StartLocalEngine(GuiService_StartLocalEngineServer) error {
	return status.Errorf(codes.Unimplemented, "method StartLocalEngine not implemented")
}
func (UnimplementedGuiServiceServer) StartFromPreviousEngine(context.Context, *StartFromPreviousEngineRequest) (*StartEngineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartFromPreviousEngine not implemented")
}
func (UnimplementedGuiServiceServer) StartEngine(context.Context, *StartEngineRequest) (*StartEngineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartEngine not implemented")
}
func (UnimplementedGuiServiceServer) ListEngines(context.Context, *ListEnginesRequest) (*ListEnginesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEngines not implemented")
}
func (UnimplementedGuiServiceServer) Subscribe(*SubscribeRequest, GuiService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedGuiServiceServer) GetEngine(context.Context, *GetEngineRequest) (*GetEngineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEngine not implemented")
}
func (UnimplementedGuiServiceServer) Listen(*ListenRequest, GuiService_ListenServer) error {
	return status.Errorf(codes.Unimplemented, "method Listen not implemented")
}
func (UnimplementedGuiServiceServer) StopEngine(context.Context, *StopEngineRequest) (*StopEngineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopEngine not implemented")
}
func (UnimplementedGuiServiceServer) mustEmbedUnimplementedGuiServiceServer() {}

// UnsafeGuiServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GuiServiceServer will
// result in compilation errors.
type UnsafeGuiServiceServer interface {
	mustEmbedUnimplementedGuiServiceServer()
}

func RegisterGuiServiceServer(s grpc.ServiceRegistrar, srv GuiServiceServer) {
	s.RegisterService(&GuiService_ServiceDesc, srv)
}

func _GuiService_StartLocalEngine_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GuiServiceServer).StartLocalEngine(&guiServiceStartLocalEngineServer{stream})
}

type GuiService_StartLocalEngineServer interface {
	SendAndClose(*StartEngineResponse) error
	Recv() (*StartLocalEngineRequest, error)
	grpc.ServerStream
}

type guiServiceStartLocalEngineServer struct {
	grpc.ServerStream
}

func (x *guiServiceStartLocalEngineServer) SendAndClose(m *StartEngineResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *guiServiceStartLocalEngineServer) Recv() (*StartLocalEngineRequest, error) {
	m := new(StartLocalEngineRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GuiService_StartFromPreviousEngine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartFromPreviousEngineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GuiServiceServer).StartFromPreviousEngine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.GuiService/StartFromPreviousEngine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GuiServiceServer).StartFromPreviousEngine(ctx, req.(*StartFromPreviousEngineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GuiService_StartEngine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartEngineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GuiServiceServer).StartEngine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.GuiService/StartEngine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GuiServiceServer).StartEngine(ctx, req.(*StartEngineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GuiService_ListEngines_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEnginesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GuiServiceServer).ListEngines(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.GuiService/ListEngines",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GuiServiceServer).ListEngines(ctx, req.(*ListEnginesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GuiService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GuiServiceServer).Subscribe(m, &guiServiceSubscribeServer{stream})
}

type GuiService_SubscribeServer interface {
	Send(*SubscribeResponse) error
	grpc.ServerStream
}

type guiServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *guiServiceSubscribeServer) Send(m *SubscribeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _GuiService_GetEngine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEngineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GuiServiceServer).GetEngine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.GuiService/GetEngine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GuiServiceServer).GetEngine(ctx, req.(*GetEngineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GuiService_Listen_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListenRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GuiServiceServer).Listen(m, &guiServiceListenServer{stream})
}

type GuiService_ListenServer interface {
	Send(*ListenResponse) error
	grpc.ServerStream
}

type guiServiceListenServer struct {
	grpc.ServerStream
}

func (x *guiServiceListenServer) Send(m *ListenResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _GuiService_StopEngine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopEngineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GuiServiceServer).StopEngine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.GuiService/StopEngine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GuiServiceServer).StopEngine(ctx, req.(*StopEngineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GuiService_ServiceDesc is the grpc.ServiceDesc for GuiService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GuiService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.GuiService",
	HandlerType: (*GuiServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartFromPreviousEngine",
			Handler:    _GuiService_StartFromPreviousEngine_Handler,
		},
		{
			MethodName: "StartEngine",
			Handler:    _GuiService_StartEngine_Handler,
		},
		{
			MethodName: "ListEngines",
			Handler:    _GuiService_ListEngines_Handler,
		},
		{
			MethodName: "GetEngine",
			Handler:    _GuiService_GetEngine_Handler,
		},
		{
			MethodName: "StopEngine",
			Handler:    _GuiService_StopEngine_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StartLocalEngine",
			Handler:       _GuiService_StartLocalEngine_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _GuiService_Subscribe_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Listen",
			Handler:       _GuiService_Listen_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "gui.proto",
}
