package routes

import (
	"api-gateway/api/v1/services"
	"api-gateway/pkg/auth/pb"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(app fiber.Router, service pb.AuthServiceClient) {
	app.Post("/register", services.RegisterUser(service))

}
