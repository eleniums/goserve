package grpc

import (
	"crypto/tls"
	"fmt"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
)

const (
	// DefaultMaxSendMsgSize is the default max send message size, per gRPC
	DefaultMaxSendMsgSize = 1024 * 1024 * 4

	// DefaultMaxRecvMsgSize is the default max receive message size, per gRPC
	DefaultMaxRecvMsgSize = 1024 * 1024 * 4
)

// Builder is used to construct a gRPC server.
type Builder struct {
	servers            []server
	options            []grpc.ServerOption
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
}

type server struct {
	RegisterFunc interface{}
	Server       interface{}
}

// New will create a new gRPC server builder.
func New() *Builder {
	return &Builder{}
}

// Build a gRPC server.
func (b *Builder) Build() *grpc.Server {
	if len(b.unaryInterceptors) > 0 {
		b.options = append(b.options, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(b.unaryInterceptors...)))
	}

	if len(b.streamInterceptors) > 0 {
		b.options = append(b.options, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(b.streamInterceptors...)))
	}

	s := grpc.NewServer(b.options...)

	for _, v := range b.servers {
		reflect.ValueOf(v.RegisterFunc).Call([]reflect.Value{
			reflect.ValueOf(s),
			reflect.ValueOf(v.Server),
		})
	}

	return s
}

// Register adds a gRPC registration function and an associated server implementation. The registerFunc should be something akin to RegisterSampleServer(s *grpc.Server, srv SampleServer) and srv should be an implementation that satisfies the srv interface in registerFunc.
func (b *Builder) Register(registerFunc interface{}, srv interface{}) *Builder {
	var s *grpc.Server
	registerFuncType := reflect.TypeOf(registerFunc)
	if registerFuncType.Kind() != reflect.Func || registerFuncType.NumIn() != 2 || registerFuncType.In(0) != reflect.TypeOf(s) || registerFuncType.In(1).Kind() != reflect.Interface {
		panic(fmt.Errorf("registerFunc is not a grpc registration function: %v, ex: RegisterSampleServer(s *grpc.Server, srv SampleServer)", registerFuncType))
	}

	serverType := reflect.TypeOf(srv)
	if !serverType.Implements(registerFuncType.In(1)) {
		panic(fmt.Errorf("Incorrect type for server: %v does not implement %v", serverType, registerFuncType.In(1)))
	}

	b.servers = append(b.servers, server{
		RegisterFunc: registerFunc,
		Server:       srv,
	})

	return b
}

// WithTLS adds configuration to provide secure communications via Transport Layer Security.
func (b *Builder) WithTLS(config *tls.Config) *Builder {
	creds := credentials.NewTLS(config)
	b.options = append(b.options, grpc.Creds(creds))
	return b
}

// WithOptions adds additional server options for customizing the server further.
func (b *Builder) WithOptions(options ...grpc.ServerOption) *Builder {
	b.options = append(b.options, options...)
	return b
}

// WithUnaryInterceptors adds unary interceptors to be used by the service. They will be executed in order, from first to last.
func (b *Builder) WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) *Builder {
	b.unaryInterceptors = append(b.unaryInterceptors, interceptors...)
	return b
}

// WithStreamInterceptors adds stream interceptors to be used by the service. They will be executed in order, from first to last.
func (b *Builder) WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) *Builder {
	b.streamInterceptors = append(b.streamInterceptors, interceptors...)
	return b
}

// WithMaxRecvMsgSize will change the size of messages that can be received by the service.
func (b *Builder) WithMaxRecvMsgSize(size int) *Builder {
	b.options = append(b.options, grpc.MaxRecvMsgSize(size))
	return b
}

// WithMaxSendMsgSize will change the size of messages that can be sent from the service.
func (b *Builder) WithMaxSendMsgSize(size int) *Builder {
	b.options = append(b.options, grpc.MaxSendMsgSize(size))
	return b
}
