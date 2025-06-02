package book

import (
	"golang-grpc/api/presenter"
	"golang-grpc/pkg/entities"
)

// Service is an interface from which our api module can access our repository of all our models
type Service interface {
	InsertBook(book *entities.Book) (*entities.Book, error)
	GetListBook() (*[]presenter.Book, error)
	// UpdateBook(book *entities.Book) (*entities.Book, error)
	// RemoveBook(ID string) error
}

type service struct {
	repository Repository
}

// NewService is used to create a single instance of the service
func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

// InsertBook is a service layer that helps insert book in BookShop
func (s *service) InsertBook(book *entities.Book) (*entities.Book, error) {
	return s.repository.CreateBook(book)
}

// GetListBook is a service layer that helps fetch all books in BookShop
func (s *service) GetListBook() (*[]presenter.Book, error) {
	return s.repository.GetListBook()
}
