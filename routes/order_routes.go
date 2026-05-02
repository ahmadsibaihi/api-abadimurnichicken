package routes

import (
	"github.com/gofiber/fiber/v2"
	"go-serve-pos/internal/pos"
)

func OrderRoutes(api fiber.Router, h *pos.PosHandler) {
	// Endpoint untuk Kasir membuat pesanan baru
	api.Post("/orders", h.CreateOrder)

	// Endpoint untuk KDS melihat antrean masak
	api.Get("/kds/active", h.GetActiveOrders)

	// Endpoint Spesifik untuk KDS (Gantikan UpdateStatus)
	api.Patch("/orders/:id/accept", h.AcceptOrder) // Koki klik "Terima/Masak"
	api.Patch("/orders/:id/finish", h.FinishOrder) // Koki klik "Selesai/Siap"	
}