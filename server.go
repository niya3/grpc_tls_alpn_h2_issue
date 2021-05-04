package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"

	"example.com/server/pkg/credit"
	"google.golang.org/grpc"
)

var withH2Flag = flag.Bool("with-h2", true, "specify h2 as next proto to TLS listener")

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

	if *withH2Flag {
		log.Println("Run server with h2")

		lis = tls.NewListener(lis, &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS13,
			NextProtos: []string{
				"h2", // enable ALPN needed for python grpcio client
			},
		})
	} else {
		log.Println("Run server WITHOUT h2")
		lis = tls.NewListener(lis, &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS13,
		})
	}

	srv := grpc.NewServer()
	credit.RegisterCreditServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(lis))
}

func (s *server) Credit(ctx context.Context, request *credit.CreditRequest) (*credit.CreditResponse, error) {
	log.Printf("Request: %g", request.GetAmount())

	return &credit.CreditResponse{Confirmation: fmt.Sprintf("Credited %g", request.GetAmount())}, nil
}
