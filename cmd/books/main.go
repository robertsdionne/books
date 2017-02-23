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
	port          = flag.Int("port", 9090, "The port")
	useFilesystem = flag.Bool("use-filesystem", false, "whether to use the filesystem")
)

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalln(err)
	}

	var fs afero.Fs
	switch {
	case *useFilesystem:
		fs = afero.NewOsFs()

	case !*useFilesystem:
		fs = afero.NewMemMapFs()
	}

	grpcServer := grpc.NewServer()
	books.RegisterLibraryServiceServer(grpcServer, books.NewServer(fs))

	grpcServer.Serve(listener)
}
