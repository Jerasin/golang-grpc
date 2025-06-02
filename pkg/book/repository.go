package book

import (
	"context"
	"golang-grpc/api/presenter"
	"golang-grpc/pkg/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateBook(book *entities.Book) (*entities.Book, error)
	GetListBook() (*[]presenter.Book, error)
	// UpdateBook(book *entities.Book) (*entities.Book, error)
	// DeleteBook(ID string) error
}

type repository struct {
	Collection *mongo.Collection
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}

// CreateBook is a mongo repository that helps to create books
func (r *repository) CreateBook(book *entities.Book) (*entities.Book, error) {
	book.ID = primitive.NewObjectID()
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	_, err := r.Collection.InsertOne(context.Background(), book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

// GetListBook is a mongo repository that helps to fetch books
func (r *repository) GetListBook() (*[]presenter.Book, error) {
	var books []presenter.Book
	cursor, err := r.Collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var book presenter.Book
		_ = cursor.Decode(&book)
		books = append(books, book)
	}
	return &books, nil
}
