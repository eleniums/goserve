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
	// Servers is used to register server handlers.
	Servers []Server

	// TLSConfig stores the TLS configuration if a secure endpoint is desired.
	TLSConfig *tls.Config

	// Options is an array of server options for customizing the server further.
	Options []grpc.ServerOption

	// UnaryInterceptors is an array of unary interceptors. They will be executed in order, from first to last.
	UnaryInterceptors []grpc.UnaryServerInterceptor

	// StreamInterceptors is an array of stream interceptors. They will be executed in order, from first to last.
	StreamInterceptors []grpc.StreamServerInterceptor
}

type Server struct {
	RegisterFunc interface{}
	Server       interface{}
}

// New will create a GRPC instance with default values.
func New() *Builder {
	return &Builder{}
}

func (b *Builder) Build() *grpc.Server {
	if len(b.UnaryInterceptors) > 0 {
		b.Options = append(b.Options, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(b.UnaryInterceptors...)))
	}

	if len(b.StreamInterceptors) > 0 {
		b.Options = append(b.Options, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(b.StreamInterceptors...)))
	}

	s := grpc.NewServer(b.Options...)

	for _, v := range b.Servers {
		err := registerServer(s, v.RegisterFunc, v.Server)
		if err != nil {
			panic(err)
		}
	}

	return s
}

func (b *Builder) Register(registerFunc interface{}, server interface{}) *Builder {
	b.Servers = append(b.Servers, Server{
		RegisterFunc: registerFunc,
		Server:       server,
	})
	return b
}

func (b *Builder) WithTLS(config *tls.Config) *Builder {
	b.TLSConfig = config
	creds := credentials.NewTLS(config)
	b.Options = append(b.Options, grpc.Creds(creds))
	return b
}

func (b *Builder) WithUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) *Builder {
	b.UnaryInterceptors = append(b.UnaryInterceptors, interceptor)
	return b
}

func (b *Builder) WithStreamInterceptor(interceptor grpc.StreamServerInterceptor) *Builder {
	b.StreamInterceptors = append(b.StreamInterceptors, interceptor)
	return b
}

func (b *Builder) WithMaxRecvMsgSize(size int) *Builder {
	b.Options = append(b.Options, grpc.MaxRecvMsgSize(size))
	return b
}

func (b *Builder) WithMaxSendMsgSize(size int) *Builder {
	b.Options = append(b.Options, grpc.MaxSendMsgSize(size))
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
