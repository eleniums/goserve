package main

import (
	"flag"
	"log"
	"net"

	"github.com/eleniums/goserve/examples/hello"

	pb "github.com/eleniums/gohost/examples/hello/proto"
	gsgrpc "github.com/eleniums/goserve/grpc"
)

func main() {
	// command-line flags
	grpcAddr := flag.String("grpc-addr", "127.0.0.1:50051", "host and port to host the gRPC endpoint")
	// httpAddr := flag.String("http-addr", "127.0.0.1:9090", "host and port to host the HTTP endpoint")
	// certFile := flag.String("cert-file", "", "cert file for enabling a TLS connection")
	// keyFile := flag.String("key-file", "", "key file for enabling a TLS connection")
	// insecureSkipVerify := flag.Bool("insecure-skip-verify", false, "true to skip verifying the certificate chain and host name")
	maxSendMsgSize := flag.Int("max-send-msg-size", gsgrpc.DefaultMaxSendMsgSize, "max message size the service is allowed to send")
	maxRecvMsgSize := flag.Int("max-recv-msg-size", gsgrpc.DefaultMaxRecvMsgSize, "max message size the service is allowed to receive")
	flag.Parse()

	// create the service
	service := hello.NewService()

	// create the grpc service
	server := gsgrpc.New().
		Register(pb.RegisterHelloServiceServer, service).
		WithMaxSendMsgSize(*maxSendMsgSize).
		WithMaxRecvMsgSize(*maxRecvMsgSize).
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
	lis, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Listening for gRPC requests at: %v", lis.Addr().String())

	// start the server
	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("Error occurred while serving: %v", err)
	}
}
