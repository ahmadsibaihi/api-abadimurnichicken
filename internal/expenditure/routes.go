package expenditure

import "github.com/gofiber/fiber/v2"

func ExpenditureRoutes(router fiber.Router, h *Handler) {
    router.Get("/", h.GetExpenditures)
    router.Post("/", h.CreateExpenditure)
    router.Get("/daily-summary", h.DailySummary)
    router.Get("/:id", h.GetExpenditureByID)
    router.Put("/:id", h.UpdateExpenditure)
    router.Delete("/:id", h.SoftDeleteExpenditure)
    router.Post("/:id/approve", h.ApproveExpenditure)
}