package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"

	"github.com/eleniums/goserve/examples/hello"

	pb "github.com/eleniums/gohost/examples/hello/proto"
	gs "github.com/eleniums/goserve"
	gsgrpc "github.com/eleniums/goserve/grpc"
)

func main() {
	// command-line flags
	grpcAddr := flag.String("grpc-addr", "127.0.0.1:50051", "host and port to host the gRPC endpoint")
	certFile := flag.String("cert-file", "", "cert file for enabling a TLS connection")
	keyFile := flag.String("key-file", "", "key file for enabling a TLS connection")
	insecureSkipVerify := flag.Bool("insecure-skip-verify", false, "true to skip verifying the certificate chain and host name")
	maxSendMsgSize := flag.Int("max-send-msg-size", gsgrpc.DefaultMaxSendMsgSize, "max message size the service is allowed to send")
	maxRecvMsgSize := flag.Int("max-recv-msg-size", gsgrpc.DefaultMaxRecvMsgSize, "max message size the service is allowed to receive")
	flag.Parse()

	// create the service
	service := hello.NewService()

	// add TLS config if requested
	var tlsConfig *tls.Config
	if *certFile != "" && *keyFile != "" {
		var err error
		tlsConfig, err = gs.NewTLS(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Error creating TLS config: %v", err)
		}
		tlsConfig.InsecureSkipVerify = *insecureSkipVerify
	}

	// create the grpc server
	server := gsgrpc.New().
		Register(pb.RegisterHelloServiceServer, service).
		WithTLS(tlsConfig).
		WithMaxSendMsgSize(*maxSendMsgSize).
		WithMaxRecvMsgSize(*maxRecvMsgSize).
		Build()

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
