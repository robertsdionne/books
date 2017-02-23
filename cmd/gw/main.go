package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/robertsdionne/books"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	booksEndpoint = flag.String("books-endpoint", "localhost:9090", "the endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	options := []grpc.DialOption{grpc.WithInsecure()}
	err := books.RegisterLibraryServiceHandlerFromEndpoint(ctx, mux, *booksEndpoint, options)
	if err != nil {
		return err
	}

	http.ListenAndServe(":8080", nil)
	return nil
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
