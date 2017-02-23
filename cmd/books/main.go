package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/robertsdionne/books"
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

	/*var fs afero.Fs
	switch {
	case *useFilesystem:
		fs = afero.NewOsFs()

	case !*useFilesystem:
		fs = afero.NewMemMapFs()
	}*/

	db, err := sql.Open("postgres", "postgres://root@gauss:30257/books?sslmode=disable")
	if err != nil {
		return
	}

	grpcServer := grpc.NewServer()
	books.RegisterLibraryServiceServer(grpcServer, books.NewServer(db))

	grpcServer.Serve(listener)
}
