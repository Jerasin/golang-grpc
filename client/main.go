package main

import (
	"context"
	"fmt"
	"golang-grpc/proto-gen/go/service"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:4040", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := service.NewAddServiceClient(conn)

	// g := gin.Default()
	app := fiber.New()

	app.Use(logger.New())

	app.Get("/messages", func(c *fiber.Ctx) error {
		text := c.Query("text")
		if text == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Text query parameter is required",
			})
		}
		req := &service.Message{Text: text}
		if res, err := client.GetMessages(context.Background(), req); err == nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"messages": res.Messages,
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	})

	app.Get("/add/:a/:b", func(c *fiber.Ctx) error {
		a, err := strconv.ParseUint(c.Params("a"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid argument A",
			})
		}
		b, err := strconv.ParseUint(c.Params("b"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid argument B",
			})
		}
		req := &service.Request{A: int64(a), B: int64(b)}
		if res, err := client.Add(context.Background(), req); err == nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"result": fmt.Sprint(res.Result),
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	})

	app.Get("/mult/:a/:b", func(c *fiber.Ctx) error {
		a, err := strconv.ParseUint(c.Params("a"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid argument A",
			})
		}
		b, err := strconv.ParseUint(c.Params("b"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid argument B",
			})
		}
		req := &service.Request{A: int64(a), B: int64(b)}
		if res, err := client.Multiply(context.Background(), req); err == nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"result": fmt.Sprint(res.Result),
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})

		}

	})
	log.Fatal(app.Listen(":3000"))
}
