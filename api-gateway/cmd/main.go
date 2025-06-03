package main

import (
	"api-gateway/pkg/auth/pb"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	auth := pb.NewAuthServiceClient(conn)

	app := fiber.New()

	app.Use(logger.New())

	app.Get("message", func(c *fiber.Ctx) error {
		req := &pb.TestRequest{
			Name: "World",
		}
		if res, err := auth.Test(context.Background(), req); err == nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": res.GetMessage(),
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	})

	log.Fatal(app.Listen(":3000"))
}
