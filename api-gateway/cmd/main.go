package main

import (
	"api-gateway/api/v1/routes"
	"api-gateway/pkg/auth/pb"
	"api-gateway/pkg/config"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed at config ", err)
	}

	conn, err := grpc.NewClient(c.AuthServiceClient, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	authService := pb.NewAuthServiceClient(conn)

	app := fiber.New()

	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!1")
	})

	app.Get("message", func(c *fiber.Ctx) error {
		req := &pb.TestRequest{
			Name: "World1",
		}
		if res, err := authService.Test(context.Background(), req); err == nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": res.GetMessage(),
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")
	auth := v1.Group("/auth")
	routes.AuthRouter(auth, authService)

	log.Fatal(app.Listen(c.Port))
}
