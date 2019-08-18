# goserve
Collection of packages for easily hosting gRPC, HTTP, or other endpoints.

## Installation
```
go get -u github.com/eleniums/goserve
```

# HTTP
```
import (
	gs "github.com/eleniums/goserve"
	gshttp "github.com/eleniums/goserve/http"
)

func serveHTTP(tlsConfig *tls.Config) error {
	// create the server
	srv := httpserver.NewServer()

	// build the server
	s := gshttp.New().
		HandleFunc("/ping", srv.Ping).
		HandleFunc("/items", srv.GetItems).
		WithTLS(tlsConfig).
		WithMiddleware(httpserver.Telemetry).
		Build()

	// start the server
	lis, err := tls.Listen("tcp", httpAddr, tlsConfig)
	if err != nil {
		return err
	}

	return s.Serve(lis)
}
```

# gRPC
