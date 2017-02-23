package books

import (
	"database/sql"
	"log"

	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/spf13/afero"
	"golang.org/x/net/context"
)

//go:generate protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. books.proto
//go:generate protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. books.proto
//go:generate protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. books.proto

type bookServer struct {
	Fs afero.Fs `inject:""`
	Db *sql.DB
}

func NewServer(db *sql.DB) *bookServer {
	return &bookServer{Db: db}
}

func (s *bookServer) ListShelves(ctx context.Context, request *ListShelvesRequest) (response *ListShelvesResponse, err error) {
	log.Println("ListShelves", request)

	rows, err := s.Db.Query(`
	SELECT name FROM shelves;
	`)
	if err != nil {
		return
	}
	defer rows.Close()

	response = &ListShelvesResponse{}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return
		}

		response.Shelves = append(response.Shelves, &Shelf{
			Name: name,
		})
	}

	return
}

func (s *bookServer) ListBooks(ctx context.Context, request *ListBooksRequest) (response *ListBooksResponse, err error) {
	log.Println("ListBooks", request)

	rows, err := s.Db.Query(`
	SELECT name, text FROM books WHERE parent=$1;
	`, request.Parent)
	if err != nil {
		return
	}
	defer rows.Close()

	response = &ListBooksResponse{}
	for rows.Next() {
		var name, text string
		err = rows.Scan(&name, &text)
		if err != nil {
			return
		}

		response.Books = append(response.Books, &Book{
			Name: name,
			Text: text,
		})
	}

	return
}

func (s *bookServer) GetBook(ctx context.Context, request *GetBookRequest) (response *Book, err error) {
	log.Println("GetBook", request)

	var name, text string
	err = s.Db.QueryRow(`
	SELECT name, text FROM books WHERE name=$1;
	`, request.Name).Scan(&name, &text)
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

	_, err = s.Db.Exec(`
	INSERT INTO books VALUES ($1, $2, $3);
	`, request.Parent, request.Book.Name, request.Book.Text)
	if err != nil {
		return
	}

	response = request.Book
	return
}

func (s *bookServer) CreateShelf(ctx context.Context, request *CreateShelfRequest) (response *Shelf, err error) {
	log.Println("CreateShelf", request)

	_, err = s.Db.Exec(`
	INSERT INTO shelves VALUES ($1);
	`, request.Shelf.Name)
	if err != nil {
		return
	}

	response = request.Shelf
	return
}

func (s *bookServer) UpdateBook(ctx context.Context, request *UpdateBookRequest) (response *Book, err error) {
	log.Println("UpdateBook", request)

	_, err = s.Db.Exec(`
	UPDATE books SET text=$2 WHERE name=$1;
	`, request.Book.Name, request.Book.Text)
	if err != nil {
		return
	}

	response = request.Book
	return
}

func (s *bookServer) DeleteShelf(ctx context.Context, request *DeleteShelfRequest) (response *google_protobuf.Empty, err error) {
	log.Println("DeleteShelf", request)

	_, err = s.Db.Exec(`
	DELETE FROM shelves WHERE name=$1;
	`, request.Name)
	if err != nil {
		return
	}

	response = &google_protobuf.Empty{}
	return
}

func (s *bookServer) DeleteBook(ctx context.Context, request *DeleteBookRequest) (response *google_protobuf.Empty, err error) {
	log.Println("DeleteBook", request)

	_, err = s.Db.Exec(`
	DELETE FROM books WHERE name=$1;
	`, request.Name)
	if err != nil {
		return
	}

	response = &google_protobuf.Empty{}
	return
}
