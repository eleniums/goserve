package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/eleniums/goserve/examples/hello"

	pb "github.com/eleniums/gohost/examples/hello/proto"
	gshttp "github.com/eleniums/goserve/http"
)

func main() {
	// command-line flags
	httpAddr := flag.String("http-addr", "127.0.0.1:9090", "host and port to host the HTTP endpoint")
	// certFile := flag.String("cert-file", "", "cert file for enabling a TLS connection")
	// keyFile := flag.String("key-file", "", "key file for enabling a TLS connection")
	// insecureSkipVerify := flag.Bool("insecure-skip-verify", false, "true to skip verifying the certificate chain and host name")
	flag.Parse()

	// create the service
	service := hello.NewService()

	// create the server
	server := gshttp.New().
		HandleFunc("/v1/hello", func(w http.ResponseWriter, req *http.Request) {
			switch req.Method {
			case http.MethodGet:
				resp, err := service.Hello(context.Background(), &pb.HelloRequest{
					Name: req.URL.Query().Get("name"),
				})
				if err != nil {
					// TODO
				}
				// TODO
			default:
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			}
		}).
		Build()

	// hoster := gohost.NewHoster()
	// hoster.GRPCAddr = *grpcAddr
	// hoster.HTTPAddr = *httpAddr
	// hoster.DebugAddr = *debugAddr
	// hoster.EnableDebug = *enableDebug
	// hoster.CertFile = *certFile
	// hoster.KeyFile = *keyFile
	// hoster.InsecureSkipVerify = *insecureSkipVerify
	// hoster.MaxSendMsgSize = *maxSendMsgSize
	// hoster.MaxRecvMsgSize = *maxRecvMsgSize

	// hoster.RegisterGRPCServer(func(s *grpc.Server) {
	// 	pb.RegisterHelloServiceServer(s, service)
	// })
	// log.Printf("Registered gRPC endpoint at: %v", *grpcAddr)

	// hoster.RegisterHTTPGateway(pb.RegisterHelloServiceHandlerFromEndpoint)
	// log.Printf("Registered HTTP endpoint at: %v", *httpAddr)

	// start listening
	lis, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Listening for HTTP requests at: %v", lis.Addr().String())

	// start the server
	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("Error occurred while serving: %v", err)
	}
}
