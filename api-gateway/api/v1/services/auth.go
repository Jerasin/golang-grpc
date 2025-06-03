package services

import (
	"api-gateway/pkg/auth/pb"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/goforj/godump"
)

func RegisterUser(service pb.AuthServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Logic for registering a user
		var req pb.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		_, err := service.Register(context.Background(), &req)
		if err != nil {
			godump.Dump(err.Error())
			if err.Error() == "rpc error: code = AlreadyExists desc = Email already exists" {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "Email already exists",
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to register user",
			})
		}

		return c.JSON(fiber.Map{
			"message": "User registered successfully",
		})
	}

}
