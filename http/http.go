package http

import (
	"crypto/tls"
	"net/http"
	"time"
)

// Builder is used to construct a gRPC server.
type Builder struct {
	middleware []func(http.Handler) http.Handler
	tlsConfig  *tls.Config
	handle     map[string]http.Handler
	handleFunc map[string]func(http.ResponseWriter, *http.Request)
}

// New will create a new gRPC server builder.
func New() *Builder {
	return &Builder{
		handle:     map[string]http.Handler{},
		handleFunc: map[string]func(http.ResponseWriter, *http.Request){},
	}
}

// Build a gRPC server.
func (b *Builder) Build() *http.Server {
	mux := http.NewServeMux()

	var chain http.Handler = mux
	for i := len(b.middleware) - 1; i >= 0; i-- {
		chain = b.middleware[i](chain)
	}

	for p, h := range b.handle {
		mux.Handle(p, h)
	}

	for p, h := range b.handleFunc {
		mux.HandleFunc(p, h)
	}

	s := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		TLSConfig:    b.tlsConfig,
		Handler:      chain,
	}

	return s
}

// Handle registers the handler for the given pattern.
func (b *Builder) Handle(pattern string, handler http.Handler) *Builder {
	b.handle[pattern] = handler
	return b
}

// HandleFunc registers the handler function for the given pattern.
func (b *Builder) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) *Builder {
	b.handleFunc[pattern] = handler
	return b
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
