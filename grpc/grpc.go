package grpc

import (
	"crypto/tls"
	"math"
	"net"

	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
)

// GRPC contains properties needed to host a gRPC server.
type GRPC struct {
	// Server is used to register server handlers.
	Server *grpc.Server

	// TLSConfig stores the TLS configuration if a secure endpoint is desired.
	TLSConfig *tls.Config

	// Options is an array of server options for customizing the server further.
	Options []grpc.ServerOption

	// UnaryInterceptors is an array of unary interceptors. They will be executed in order, from first to last.
	UnaryInterceptors []grpc.UnaryServerInterceptor

	// StreamInterceptors is an array of stream interceptors. They will be executed in order, from first to last.
	StreamInterceptors []grpc.StreamServerInterceptor

	// MaxSendMsgSize will change the size of the message that can be sent from the service.
	MaxSendMsgSize int

	// MaxRecvMsgSize will change the size of the message that can be received by the service.
	MaxRecvMsgSize int
}

// New will create a GRPC instance with default values.
func New() *GRPC {
	return &GRPC{
		MaxSendMsgSize: math.MaxInt32,
		MaxRecvMsgSize: math.MaxInt32,
	}
}

// Serve will accept incoming connections on the given listener.
func (g *GRPC) Serve(lis net.Listener) error {
	return g.Server.Serve(lis)
}

// Initialize will create a gRPC server based on the properties already set.
func (g *GRPC) Initialize() {
	opt := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(g.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(g.MaxSendMsgSize),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(g.UnaryInterceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(g.StreamInterceptors...)),
	}

	opt = append(opt, g.Options...)

	g.Server = grpc.NewServer(opt...)
}
