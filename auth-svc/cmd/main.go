package main

import (
	"auth-svc/pkg/config"
	"auth-svc/pkg/db"
	"auth-svc/pkg/middleware"
	"auth-svc/pkg/pb"
	"auth-svc/pkg/repositories"
	"auth-svc/pkg/services"
	"auth-svc/pkg/utils"
	"fmt"

	"context"
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

	jwt := utils.JWTWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "go-grpc-auth-svc",
		ExpirationHours: 24 * 365,
	}

	db := con.Database(c.DbName)
	fmt.Println("AuditLogCollection", c.AuditLogCollection)
	fmt.Println("UserCollection", c.UserCollection)

	auditLogCollection := db.Collection(c.AuditLogCollection)
	auditLogRepo := repositories.NewAuditLogRepository(auditLogCollection)

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
			return middleware.LoggerInterceptor(ctx, req, info, handler, auditLogRepo, &[]middleware.CustomParameterMiddleware{})
		},
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
			return middleware.AuthorizationInterceptor(ctx, req, info, handler, jwt, &[]middleware.CustomParameterMiddleware{{SkipPath: "/auth.AuthService/LoginUser"}})
		},
	))

	userCollection := db.Collection(c.UserCollection)
	userRepo := repositories.NewUserRepository(userCollection)

	s := services.Server{
		UserRepo: userRepo,
		Jwt:      &jwt,
	}
	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}

}
