package main

import (
	"context"
	"golang-grpc/proto-gen/go/service"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	service.UnimplementedAddServiceServer
}

func main() {
	lis, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer()
	service.RegisterAddServiceServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(lis); e != nil {
		panic(e)
	}
}

func (s *server) GetMessages(_ context.Context, request *service.Message) (*service.MessageResponse, error) {
	return &service.MessageResponse{
		Messages: request.Text,
	}, nil
}

func (s *server) Add(_ context.Context, request *service.Request) (*service.Response, error) {
	a, b := request.GetA(), request.GetB()

	result := a + b

	return &service.Response{Result: result}, nil
}

func (s *server) Multiply(_ context.Context, request *service.Request) (*service.Response, error) {
	a, b := request.GetA(), request.GetB()

	result := a * b

	return &service.Response{Result: result}, nil
}
