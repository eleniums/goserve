package grpc

import (
	"crypto/tls"
	"math"
	"net"

	"google.golang.org/grpc"
)

type GRPC struct {
	// Server is used to register server handlers.
	Server *grpc.Server

	// TLSConfig stores the TLS configuration if a secure endpoint is desired.
	TLSConfig *tls.Config

	// Opt is an array of server options for customizing the server further.
	Opt []grpc.ServerOption

	// UnaryInterceptors is an array of unary interceptors. They will be executed in order, from first to last.
	UnaryInterceptors []grpc.UnaryServerInterceptor

	// StreamInterceptors is an array of stream interceptors. They will be executed in order, from first to last.
	StreamInterceptors []grpc.StreamServerInterceptor

	// MaxSendMsgSize will change the size of the message that can be sent from the service.
	MaxSendMsgSize int

	// MaxRecvMsgSize will change the size of the message that can be received by the service.
	MaxRecvMsgSize int
}

func New() *GRPC {
	return &GRPC{
		MaxSendMsgSize: math.MaxInt32,
		MaxRecvMsgSize: math.MaxInt32,
	}
}

func (g *GRPC) Serve(l net.Listener) error {
	return g.Server.Serve(l)
}

func (g *GRPC) Initialize() {
	opt := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(g.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(g.MaxSendMsgSize),
	}

	g.Server = grpc.NewServer(opt...)
}
