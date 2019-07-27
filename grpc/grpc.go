package grpc

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
)

// Builder is used to construct a gRPC server.
type Builder struct {
	// Servers is used to register server handlers.
	Servers map[interface{}]interface{}

	// TLSConfig stores the TLS configuration if a secure endpoint is desired.
	TLSConfig *tls.Config

	// Options is an array of server options for customizing the server further.
	Options []grpc.ServerOption

	// UnaryInterceptors is an array of unary interceptors. They will be executed in order, from first to last.
	UnaryInterceptors []grpc.UnaryServerInterceptor

	// StreamInterceptors is an array of stream interceptors. They will be executed in order, from first to last.
	StreamInterceptors []grpc.StreamServerInterceptor
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
	// TODO: register grpc servers
	// for k, v := range b.Servers {
	// 	reflectFunc := reflect.TypeOf(k)
	// 	server := reflect.TypeOf(v)
	// 	// TODO: how to call a reflected method?
	// }

	return s
}

func (b *Builder) Register(registerFunc interface{}, server interface{}) *Builder {
	b.Servers[registerFunc] = server
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
