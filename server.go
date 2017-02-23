package books

import (
	"fmt"
	"log"
	"strings"

	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/spf13/afero"
	"golang.org/x/net/context"
)

//go:generate protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. books.proto
//go:generate protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. books.proto
//go:generate protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. books.proto

type bookServer struct {
	Fs afero.Fs `inject:""`
}

func NewServer(fs afero.Fs) *bookServer {
	return &bookServer{Fs: fs}
}

func (s *bookServer) ListShelves(ctx context.Context, request *ListShelvesRequest) (response *ListShelvesResponse, err error) {
	log.Println("ListShelves", request)

	response = &ListShelvesResponse{}
	shelves, err := afero.ReadDir(s.Fs, "shelves")
	if err != nil {
		return
	}

	response = &ListShelvesResponse{}
	for i := range shelves {
		response.Shelves = append(response.Shelves, &Shelf{
			Name: fmt.Sprintf("shelves/%s", shelves[i].Name()),
		})
	}

	return
}

func (s *bookServer) ListBooks(ctx context.Context, request *ListBooksRequest) (response *ListBooksResponse, err error) {
	log.Println("ListBooks", request)

	response = &ListBooksResponse{}
	books, err := afero.ReadDir(s.Fs, fmt.Sprintf("%s/books", request.Parent))
	if err != nil {
		return
	}

	response = &ListBooksResponse{}
	for i := range books {
		name := fmt.Sprintf("%s/books/%s", request.Parent, books[i].Name())

		var text []byte
		text, err = afero.ReadFile(s.Fs, name)
		if err != nil {
			return
		}

		response.Books = append(response.Books, &Book{
			Name: name,
			Text: string(text),
		})
	}

	return
}

func (s *bookServer) GetBook(ctx context.Context, request *GetBookRequest) (response *Book, err error) {
	log.Println("GetBook", request)

	text, err := afero.ReadFile(s.Fs, request.Name)
	if err != nil {
		return
	}

	response = &Book{
		Name: request.Name,
		Text: string(text),
	}
	return
}

func (s *bookServer) CreateBook(ctx context.Context, request *CreateBookRequest) (response *Book, err error) {
	log.Println("CreateBook", request)

	books := fmt.Sprintf("%s/books", request.Parent)
	exists, err := afero.DirExists(s.Fs, books)
	if err != nil {
		return
	}

	if !exists {
		err = s.Fs.Mkdir(books, 0755)
		if err != nil {
			return
		}
	}

	if !strings.HasPrefix(request.Book.Name, books) {
		err = fmt.Errorf("Expected book name to start with %s", books)
		return
	}

	exists, err = afero.Exists(s.Fs, request.Book.Name)
	if err != nil {
		return
	}

	if exists {
		err = fmt.Errorf("Expected book %s not to exist", request.Book.Name)
		return
	}

	err = afero.WriteFile(s.Fs, request.Book.Name, []byte(request.Book.Text), 0655)
	if err != nil {
		return
	}

	response = request.Book
	return
}

func (s *bookServer) CreateShelf(ctx context.Context, request *CreateShelfRequest) (response *Shelf, err error) {
	log.Println("CreateShelf", request)
	err = s.Fs.MkdirAll(request.Shelf.Name, 0755)
	if err != nil {
		return
	}

	response = request.Shelf
	return
}

func (s *bookServer) UpdateBook(ctx context.Context, request *UpdateBookRequest) (response *Book, err error) {
	log.Println("UpdateBook", request)

	exists, err := afero.Exists(s.Fs, request.Book.Name)
	if err != nil {
		return
	}

	if !exists {
		err = fmt.Errorf("Book %s not found", request.Book.Name)
		return
	}

	err = afero.WriteFile(s.Fs, request.Book.Name, []byte(request.Book.Text), 0655)
	if err != nil {
		return
	}

	response = request.Book
	return
}

func (s *bookServer) DeleteShelf(ctx context.Context, request *DeleteShelfRequest) (response *google_protobuf.Empty, err error) {
	log.Println("DeleteShelf", request)

	err = s.Fs.RemoveAll(request.Name)
	if err != nil {
		return
	}

	response = &google_protobuf.Empty{}
	return
}

func (s *bookServer) DeleteBook(ctx context.Context, request *DeleteBookRequest) (response *google_protobuf.Empty, err error) {
	log.Println("DeleteBook", request)

	err = s.Fs.Remove(request.Name)
	if err != nil {
		return
	}

	response = &google_protobuf.Empty{}
	return
}
