package main

import (
	"auth-svc/pkg/config"
	"auth-svc/pkg/pb"
	"auth-svc/pkg/services"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("failed at config ", err)
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("failed at listening : ", err)
	}
	fmt.Println("Auth svc on ", c.Port)

	grpcServer := grpc.NewServer()
	s := services.Server{}

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
