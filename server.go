package books

import "golang.org/x/net/context"

//go:generate protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. books.proto
//go:generate protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. books.proto
//go:generate protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. books.proto

type bookServer struct{}

func (s *bookServer) ListBooks(context.Context, *ListBooksRequest) (*ListBooksResponse, error) {
}

func (s *bookServer) GetBook(context.Context, *GetBookRequest) (*Book, error) {
}

func (s *bookServer) CreateBook(context.Context, *CreateBookRequest) (*Book, error) {
}

func (s *bookServer) CreateShelf(context.Context, *CreateShelfRequest) (*Shelf, error) {
}

func (s *bookServer) UpdateBook(context.Context, *UpdateBookRequest) (*Book, error) {
}
