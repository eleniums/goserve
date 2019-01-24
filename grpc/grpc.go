package grpc

import (
	"net"

	"google.golang.org/grpc"
)

type GRPC struct {
	Server *grpc.Server
}

func New() *GRPC {
	return &GRPC{}
}

func (g *GRPC) Serve(l net.Listener) error {
	// TODO: implement serve
	return nil
}
