# goserve
Collection of packages for easily hosting gRPC, HTTP, or other endpoints.

## Installation
```
go get -u github.com/eleniums/goserve
```

# HTTP
Import the relevant packages:
```
import (
	gs "github.com/eleniums/goserve"
	gshttp "github.com/eleniums/goserve/http"
)
```

Build an HTTP server and listen:
```
// create TLS configuration (optional)
tlsConfig, err = gs.NewTLS(certFile, keyFile)
if err != nil {
	return err
}

// build the server
s := gshttp.New().
	HandleFunc("/ping", example.Ping).
	HandleFunc("/items", example.GetItems).
	WithTLS(tlsConfig).
	WithMiddleware(example.Telemetry).
	Build()

// create a listener
lis, err := tls.Listen("tcp", "127.0.0.1:9090", tlsConfig)
if err != nil {
	return err
}

// serve the endpoint
err = s.Serve(lis)
if err != nil {
    return err
}
```

# gRPC
