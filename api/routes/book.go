package routes

import (
	"golang-grpc/api/handlers"
	"golang-grpc/pkg/book"

	"github.com/gofiber/fiber/v2"
)

func BookRouter(app fiber.Router, service book.Service) {
	app.Get("/books", handlers.GetBooks(service))
	app.Post("/books", handlers.AddBook(service))
}
