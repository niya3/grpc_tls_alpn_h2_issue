package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc/credentials"

	"example.com/server/pkg/credit"
	"google.golang.org/grpc"
)

var withH2Flag = flag.Bool("with-h2", false, "specify h2 as next proto to TLS listener")
var viaGrpcCreds = flag.Bool("via-grpc-creds", false, "pass TLS config via grpc.Creds when `true` or via tls.NewListener when `false`")

type server struct {
	credit.UnimplementedCreditServiceServer
}

func main() {
	flag.Parse()

	cert, err := tls.LoadX509KeyPair("./root.crt", "./root.key")
	if err != nil {
		log.Fatalln(err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	}

	if *withH2Flag {
		log.Println("TLS config includes h2")
		tlsConfig.NextProtos = []string{"h2"} // enable ALPN needed for python grpcio client
	} else {
		log.Println("TLS config doesn't include h2")
	}

	var opts []grpc.ServerOption
	if *viaGrpcCreds {
		log.Println("Pass TLS config via grpc.Creds")
		tlsCredentials := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.Creds(tlsCredentials))

	} else {
		log.Println("Pass TLS config via tls.NewListener")
		lis = tls.NewListener(lis, tlsConfig)
	}

	srv := grpc.NewServer(opts...)
	credit.RegisterCreditServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(lis))
}

func (s *server) Credit(ctx context.Context, request *credit.CreditRequest) (*credit.CreditResponse, error) {
	log.Printf("Request: %g", request.GetAmount())

	return &credit.CreditResponse{Confirmation: fmt.Sprintf("Credited %g", request.GetAmount())}, nil
}
