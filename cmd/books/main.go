package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/robertsdionne/books"
	"github.com/spf13/afero"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 9090, "The port")
)

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	books.RegisterLibraryServiceServer(grpcServer, books.NewServer(afero.NewMemMapFs()))

	grpcServer.Serve(listener)
}
