package services

import (
	"auth-svc/pkg/models"
	"auth-svc/pkg/pb"
	"auth-svc/pkg/repositories"
	"auth-svc/pkg/utils"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	UserRepo *repositories.UserRepository
}

func (s *Server) Test(_ context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	return &pb.TestResponse{
		Message: "Hello " + req.Name,
	}, nil
}

func (s *Server) Register(_ context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	if _, err := s.UserRepo.IsExist(map[string]any{"email": req.Email}); err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "Email already exists")
	}

	s.UserRepo.Register(&models.User{
		Email:    req.Email,
		Password: utils.HashPassword(req.Password),
	})

	return &pb.RegisterResponse{
		Status: 200,
	}, nil
}
