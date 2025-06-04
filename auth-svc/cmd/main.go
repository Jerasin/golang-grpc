package main

import (
	"auth-svc/pkg/config"
	"auth-svc/pkg/db"
	"auth-svc/pkg/pb"
	"auth-svc/pkg/repositories"
	"auth-svc/pkg/services"
	"auth-svc/pkg/utils"

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
	defer cancel()

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalln("failed at listening : ", err)
	}

	grpcServer := grpc.NewServer()
	userCollection := con.Database(c.DbName).Collection(c.UserCollection)
	userRepo := repositories.NewUserRepository(userCollection)
	jwt := utils.JWTWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "go-grpc-auth-svc",
		ExpirationHours: 24 * 365,
	}

	s := services.Server{
		UserRepo: userRepo,
		Jwt:      &jwt,
	}
	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}

}
