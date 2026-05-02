package routes

import (
    "go-serve-pos/internal/pos"
    "github.com/gofiber/fiber/v2"
)

func CategoryRoutes(api fiber.Router, h *pos.PosHandler) {
    api.Get("/categories", h.GetCategories)
    api.Get("/categories/:id", h.GetCategoryByID)
    api.Post("/categories", h.CreateCategory)
    api.Put("/categories/:id", h.UpdateCategory)
    api.Delete("/categories/:id", h.DeleteCategory)
}

func VariantRoutes(api fiber.Router, h *pos.PosHandler) {
    api.Get("/products/:product_id/variants", h.GetVariants)
    api.Post("/products/:product_id/variants", h.CreateVariant)
    api.Put("/variants/:id", h.UpdateVariant)
    api.Delete("/variants/:id", h.DeleteVariant)
}

func AddonRoutes(api fiber.Router, h *pos.PosHandler) {
    api.Get("/addons", h.GetAddons)
    api.Post("/addons", h.CreateAddon)
    api.Put("/addons/:id", h.UpdateAddon)
    api.Delete("/addons/:id", h.DeleteAddon)
    api.Post("/products/:product_id/addons/:addon_id", h.AssignAddonToProduct)
    api.Delete("/products/:product_id/addons/:addon_id", h.RemoveAddonFromProduct)
}

func ComboRoutes(api fiber.Router, h *pos.PosHandler) {
    api.Get("/combos", h.GetCombos)
    api.Get("/combos/:id", h.GetComboByID)
    api.Post("/combos", h.CreateCombo)
    api.Put("/combos/:id", h.UpdateCombo)
    api.Delete("/combos/:id", h.DeleteCombo)
    api.Post("/combos/:combo_id/slots", h.CreateComboSlot)
    api.Put("/slots/:id", h.UpdateComboSlot)
    api.Delete("/slots/:id", h.DeleteComboSlot)
    api.Post("/slots/:slot_id/options", h.CreateComboSlotOption)
    api.Delete("/options/:id", h.DeleteComboSlotOption)
}

func SpicyLevelRoutes(api fiber.Router, h *pos.PosHandler) {
	// CRUD SpicyLevel
	api.Get("/spicy-levels", h.GetSpicyLevels)
	api.Get("/spicy-levels/:id", h.GetSpicyLevelByID)
	api.Post("/spicy-levels", h.CreateSpicyLevel)
	api.Put("/spicy-levels/:id", h.UpdateSpicyLevel)
	api.Delete("/spicy-levels/:id", h.DeleteSpicyLevel)

	// Assign/Remove SpicyLevel to Product
	api.Get("/products/:product_id/spicy-levels", h.GetProductSpicyLevels)
	api.Post("/products/:product_id/spicy-levels/:level_id", h.AssignSpicyLevelToProduct)
	api.Delete("/products/:product_id/spicy-levels/:level_id", h.RemoveSpicyLevelFromProduct)
}

func TimeMenuRoutes(api fiber.Router, h *pos.PosHandler) {
    api.Get("/time-menus", h.GetTimeMenus)
    api.Get("/time-menus/:id", h.GetTimeMenuByID)
    api.Post("/time-menus", h.CreateTimeMenu)
    api.Put("/time-menus/:id", h.UpdateTimeMenu)
    api.Delete("/time-menus/:id", h.DeleteTimeMenu)
    api.Get("/products/:product_id/time-menus", h.GetProductTimeMenus)
    api.Post("/products/:product_id/time-slots", h.AssignTimeSlotToProduct)
    api.Delete("/time-menus/:id", h.RemoveTimeSlotFromProduct)
}