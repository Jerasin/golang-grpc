package main

import (
	"auth-svc/pkg/pb"
	"auth-svc/pkg/services"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	// c, err := config.LoadConfig()

	// fmt.Println("Auth svc config loaded", c.Port)

	// if err != nil {
	// 	log.Fatalln("failed at config ", err)
	// }

	fmt.Println("Starting Auth Service...")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln("failed at listening : ", err)
	}

	grpcServer := grpc.NewServer()
	s := services.Server{}
	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
