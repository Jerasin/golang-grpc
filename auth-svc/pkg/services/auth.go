package services

import (
	"auth-svc/pkg/models"
	"auth-svc/pkg/pb"
	"context"

	"github.com/goforj/godump"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	MongoCon *mongo.Client
}

func (s *Server) Test(_ context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	bookCollection := s.MongoCon.Database("books").Collection("books")
	if bookCollection == nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to database")
	}

	var books []models.Book
	cursor, err := bookCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var book models.Book
		_ = cursor.Decode(&book)
		books = append(books, book)
	}

	godump.Dump(books)
	return &pb.TestResponse{
		Message: "Hello " + req.Name,
	}, nil
}
