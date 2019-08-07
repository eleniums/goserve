package http

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/justinas/alice"
)

// Builder is used to construct a gRPC server.
type Builder struct {
	middleware []alice.Constructor
	tlsConfig  *tls.Config
}

// New will create a new gRPC server builder.
func New() *Builder {
	return &Builder{}
}

// Build a gRPC server.
func (b *Builder) Build() *http.Server {
	mux := http.NewServeMux()
	chain := alice.New(b.middleware...).Then(mux)

	s := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		TLSConfig:    b.tlsConfig,
		Handler:      chain,
	}

	return s
}

// WithTLS adds configuration to provide secure communications via TLS (Transport Layer Security). Use server.Serve with a TLS listener or server.ServeTLS with a regular listener.
func (b *Builder) WithTLS(config *tls.Config) *Builder {
	b.tlsConfig = config
	return b
}

// WithMiddleware adds middleware to be used by the service. They will be executed in order, from first to last.
func (b *Builder) WithMiddleware(middleware ...func(http.Handler) http.Handler) *Builder {
	for _, v := range middleware {
		b.middleware = append(b.middleware, v)
	}
	return b
}
