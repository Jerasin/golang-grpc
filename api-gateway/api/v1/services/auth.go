package services

import (
	"api-gateway/pkg/auth/pb"
	"api-gateway/pkg/constant"
	"api-gateway/pkg/models"
	"api-gateway/pkg/utils"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/goforj/godump"
)

func RegisterUser(service pb.AuthServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer utils.PanicHandler(c)

		var req pb.RegisterRequest

		if err := c.BodyParser(&req); err != nil {
			errMsg := "Failed to parse request body"
			utils.PanicException(constant.BadRequest, &errMsg)
		}

		user := models.User{
			Email:    req.Email,
			Password: req.Password,
		}

		utils.Validate(&user)

		_, err := service.RegisterUser(context.Background(), &req)
		if err != nil {
			godump.Dump(err.Error())
			utils.GrpcPanicException(err, nil)
		}

		return c.JSON(fiber.Map{
			"message": "User registered successfully",
		})
	}

}

func LoginUser(service pb.AuthServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer utils.PanicHandler(c)

		var req pb.LoginRequest

		if err := c.BodyParser(&req); err != nil {
			errMsg := "Failed to parse request body"
			utils.PanicException(constant.BadRequest, &errMsg)
		}

		user := models.User{
			Email:    req.Email,
			Password: req.Password,
		}

		utils.Validate(&user)

		res, err := service.LoginUser(context.Background(), &req)
		if err != nil {
			godump.Dump(err.Error())
			utils.GrpcPanicException(err, nil)
		}

		return c.JSON(fiber.Map{
			"message": "User logged in successfully",
			"data":    res,
		})
	}

}
