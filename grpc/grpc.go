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
	// servers is used to register server handlers.
	servers []server

	// options is an array of server options for customizing the server further.
	options []grpc.ServerOption

	// unaryInterceptors is an array of unary interceptors. They will be executed in order, from first to last.
	unaryInterceptors []grpc.UnaryServerInterceptor

	// streamInterceptors is an array of stream interceptors. They will be executed in order, from first to last.
	streamInterceptors []grpc.StreamServerInterceptor
}

type server struct {
	RegisterFunc interface{}
	Server       interface{}
}

// New will create a GRPC instance with default values.
func New() *Builder {
	return &Builder{}
}

func (b *Builder) Build() *grpc.Server {
	if len(b.unaryInterceptors) > 0 {
		b.options = append(b.options, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(b.unaryInterceptors...)))
	}

	if len(b.streamInterceptors) > 0 {
		b.options = append(b.options, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(b.streamInterceptors...)))
	}

	s := grpc.NewServer(b.options...)

	for _, v := range b.servers {
		err := registerServer(s, v.RegisterFunc, v.Server)
		if err != nil {
			panic(err)
		}
	}

	return s
}

func (b *Builder) Register(registerFunc interface{}, srv interface{}) *Builder {
	b.servers = append(b.servers, server{
		RegisterFunc: registerFunc,
		Server:       srv,
	})
	return b
}

func (b *Builder) WithTLS(config *tls.Config) *Builder {
	creds := credentials.NewTLS(config)
	b.options = append(b.options, grpc.Creds(creds))
	return b
}

func (b *Builder) WithUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) *Builder {
	b.unaryInterceptors = append(b.unaryInterceptors, interceptor)
	return b
}

func (b *Builder) WithStreamInterceptor(interceptor grpc.StreamServerInterceptor) *Builder {
	b.streamInterceptors = append(b.streamInterceptors, interceptor)
	return b
}

func (b *Builder) WithMaxRecvMsgSize(size int) *Builder {
	b.options = append(b.options, grpc.MaxRecvMsgSize(size))
	return b
}

func (b *Builder) WithMaxSendMsgSize(size int) *Builder {
	b.options = append(b.options, grpc.MaxSendMsgSize(size))
	return b
}

func registerServer(s *grpc.Server, registerFunc interface{}, server interface{}) error {
	registerFuncValue := reflect.ValueOf(registerFunc)
	if registerFuncValue.Kind() != reflect.Func ||
		registerFuncValue.Type().NumIn() != 2 ||
		registerFuncValue.Type().In(0) != reflect.TypeOf(s) ||
		registerFuncValue.Type().In(1).Kind() != reflect.Interface {
		return fmt.Errorf("registerFunc is not a grpc registration function: %v, ex: RegisterSampleServer(s *grpc.Server, srv SampleServer)", registerFuncValue.Type())
	}

	serverValue := reflect.ValueOf(server)
	if !serverValue.Type().Implements(registerFuncValue.Type().In(1)) {
		return fmt.Errorf("Incorrect type for server: %v does not implement %v", serverValue.Type(), registerFuncValue.Type().In(1))
	}

	registerFuncValue.Call([]reflect.Value{
		reflect.ValueOf(s),
		serverValue,
	})

	return nil
}
