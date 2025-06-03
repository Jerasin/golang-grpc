package main

import (
	"auth-svc/pkg/config"
	"auth-svc/pkg/db"
	"auth-svc/pkg/pb"
	"auth-svc/pkg/services"

	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed at config ", err)
	}

	con, cancel, err := db.MongoConnection()
	if err != nil {
		log.Fatal("Database Connection Error: ", err)
	}

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalln("failed at listening : ", err)
	}

	grpcServer := grpc.NewServer()
	s := services.Server{
		MongoCon: con,
	}
	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}

	defer cancel()
}
