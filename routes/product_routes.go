package routes

import (
    "github.com/gofiber/fiber/v2"
    "go-serve-pos/internal/pos"
)

func ProductRoutes(api fiber.Router, h *pos.PosHandler) {
    api.Get("/products", h.GetProducts)    
    api.Post("/products", h.CreateProduct)
    api.Put("/products/:id", h.UpdateProduct)
    api.Delete("/products/:id", h.DeleteProduct)
}