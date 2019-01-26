package goserve

import (
	"net"
)

// Server is an interface for something that accepts incoming connections.
type Server interface {
	// Serve will accept incoming connections on the given listener.
	Serve(lis net.Listener) error
}
