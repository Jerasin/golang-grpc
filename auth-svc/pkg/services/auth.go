package services

import (
	"auth-svc/pkg/pb"
	"context"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
}

func (s *Server) Test(_ context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	return &pb.TestResponse{
		Message: "Hello " + req.Name,
	}, nil
}
