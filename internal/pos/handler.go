package pos

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PosHandler struct {
	service PosService
}

func NewHandler(s PosService) *PosHandler {
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, 0755)
		if err != nil {
			fmt.Println("Gagal membuat folder uploads:", err)
		} else {
			fmt.Println("Folder 'uploads' berhasil dibuat otomatis.")
		}
	}
	return &PosHandler{service: s}
}

// ─────────────────────────────────────────
// PRODUCT
// ─────────────────────────────────────────

func (h *PosHandler) CreateProduct(c *fiber.Ctx) error {
	name := c.FormValue("name")
	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)
	stock, _ := strconv.Atoi(c.FormValue("stock"))
	cookingTime, _ := strconv.Atoi(c.FormValue("cooking_time"))
	categoryID, _ := strconv.ParseUint(c.FormValue("category_id"), 10, 64)
	description := c.FormValue("description")
	isBestSeller := c.FormValue("is_best_seller") == "true"
	isNew := c.FormValue("is_new") == "true"

	file, err := c.FormFile("image")
	var fileName string
	if err == nil {
		fileName = fmt.Sprintf("%s-%s", time.Now().Format("20060102150405"), file.Filename)
		c.SaveFile(file, filepath.Join("./uploads", fileName))
	}

	newProduct := Product{
		Name:         name,
		Description:  description,
		Price:        price,
		Stock:        stock,
		CookingTime:  cookingTime,
		CategoryID:   uint(categoryID),
		Image:        fileName,
		IsActive:     true,
		IsBestSeller: isBestSeller,
		IsNew:        isNew,
	}

	result, err := h.service.AddProduct(newProduct)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(result)
}

func (h *PosHandler) GetProducts(c *fiber.Ctx) error {
	products, err := h.service.GetAllProducts()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

func (h *PosHandler) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	product, err := h.service.GetProductByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}
	return c.JSON(product)
}

func (h *PosHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	existingProduct, err := h.service.GetProductByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	name := c.FormValue("name")
	price, _ := strconv.ParseFloat(c.FormValue("price"), 64)
	stock, _ := strconv.Atoi(c.FormValue("stock"))
	cookingTime, _ := strconv.Atoi(c.FormValue("cooking_time"))
	categoryID, _ := strconv.ParseUint(c.FormValue("category_id"), 10, 64)
	description := c.FormValue("description")
	isBestSeller := c.FormValue("is_best_seller") == "true"
	isNew := c.FormValue("is_new") == "true"
	isActive := c.FormValue("is_active") != "false"

	file, err := c.FormFile("image")
	var fileName string
	if err == nil {
		if existingProduct.Image != "" {
			oldImagePath := filepath.Join("./uploads", existingProduct.Image)
			if _, err := os.Stat(oldImagePath); err == nil {
				os.Remove(oldImagePath)
			}
		}
		fileName = fmt.Sprintf("%s-%s", time.Now().Format("20060102150405"), file.Filename)
		c.SaveFile(file, filepath.Join("./uploads", fileName))
	} else {
		fileName = existingProduct.Image
	}

	updatedProduct := Product{
		Name:         name,
		Description:  description,
		Price:        price,
		Stock:        stock,
		CookingTime:  cookingTime,
		CategoryID:   uint(categoryID),
		Image:        fileName,
		IsActive:     isActive,
		IsBestSeller: isBestSeller,
		IsNew:        isNew,
	}

	result, err := h.service.UpdateProduct(id, updatedProduct)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *PosHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	existingProduct, err := h.service.GetProductByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	err = h.service.DeleteProduct(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if existingProduct.Image != "" {
		imagePath := filepath.Join("./uploads", existingProduct.Image)
		if _, err := os.Stat(imagePath); err == nil {
			os.Remove(imagePath)
		}
	}

	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}

// ─────────────────────────────────────────
// CATEGORY
// ─────────────────────────────────────────

func (h *PosHandler) CreateCategory(c *fiber.Ctx) error {
	var category Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	result, err := h.service.CreateCategory(category)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(result)
}

func (h *PosHandler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(categories)
}

func (h *PosHandler) GetCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}
	return c.JSON(category)
}

func (h *PosHandler) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := h.service.GetCategoryByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	var category Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	result, err := h.service.UpdateCategory(id, category)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *PosHandler) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := h.service.GetCategoryByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	err = h.service.DeleteCategory(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Category deleted successfully"})
}

// ─────────────────────────────────────────
// VARIANT
// ─────────────────────────────────────────

func (h *PosHandler) CreateVariant(c *fiber.Ctx) error {
	productID, _ := strconv.ParseUint(c.Params("product_id"), 10, 64)

	var variant ProductVariant
	if err := c.BodyParser(&variant); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	variant.ProductID = uint(productID)

	result, err := h.service.CreateVariant(variant)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(result)
}

func (h *PosHandler) GetVariants(c *fiber.Ctx) error {
	productID := c.Params("product_id")
	variants, err := h.service.GetVariantsByProductID(productID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(variants)
}

func (h *PosHandler) UpdateVariant(c *fiber.Ctx) error {
	id := c.Params("id")

	var variant ProductVariant
	if err := c.BodyParser(&variant); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	result, err := h.service.UpdateVariant(id, variant)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *PosHandler) DeleteVariant(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.DeleteVariant(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Variant deleted successfully"})
}

// ─────────────────────────────────────────
// ADDON
// ─────────────────────────────────────────

func (h *PosHandler) CreateAddon(c *fiber.Ctx) error {
	var addon Addon
	if err := c.BodyParser(&addon); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	result, err := h.service.CreateAddon(addon)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(result)
}

func (h *PosHandler) GetAddons(c *fiber.Ctx) error {
	addons, err := h.service.GetAllAddons()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(addons)
}

func (h *PosHandler) UpdateAddon(c *fiber.Ctx) error {
	id := c.Params("id")

	var addon Addon
	if err := c.BodyParser(&addon); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	result, err := h.service.UpdateAddon(id, addon)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *PosHandler) DeleteAddon(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.DeleteAddon(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Addon deleted successfully"})
}

func (h *PosHandler) AssignAddonToProduct(c *fiber.Ctx) error {
	productID, _ := strconv.ParseUint(c.Params("product_id"), 10, 64)
	addonID, _ := strconv.ParseUint(c.Params("addon_id"), 10, 64)

	err := h.service.AssignAddonToProduct(uint(productID), uint(addonID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Addon assigned to product"})
}

func (h *PosHandler) RemoveAddonFromProduct(c *fiber.Ctx) error {
	productID, _ := strconv.ParseUint(c.Params("product_id"), 10, 64)
	addonID, _ := strconv.ParseUint(c.Params("addon_id"), 10, 64)

	err := h.service.RemoveAddonFromProduct(uint(productID), uint(addonID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Addon removed from product"})
}

// ─────────────────────────────────────────
// COMBO
// ─────────────────────────────────────────

func (h *PosHandler) CreateCombo(c *fiber.Ctx) error {
	var combo ComboPackage
	if err := c.BodyParser(&combo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	result, err := h.service.CreateCombo(combo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(result)
}

func (h *PosHandler) GetCombos(c *fiber.Ctx) error {
	combos, err := h.service.GetAllCombos()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(combos)
}

func (h *PosHandler) GetComboByID(c *fiber.Ctx) error {
	id := c.Params("id")
	combo, err := h.service.GetComboByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Combo not found"})
	}
	return c.JSON(combo)
}

func (h *PosHandler) UpdateCombo(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := h.service.GetComboByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Combo not found"})
	}

	var combo ComboPackage
	if err := c.BodyParser(&combo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	result, err := h.service.UpdateCombo(id, combo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *PosHandler) DeleteCombo(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := h.service.GetComboByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Combo not found"})
	}

	err = h.service.DeleteCombo(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Combo deleted successfully"})
}

// ─────────────────────────────────────────
// COMBO SLOT
// ─────────────────────────────────────────

func (h *PosHandler) CreateComboSlot(c *fiber.Ctx) error {
	comboID, _ := strconv.ParseUint(c.Params("combo_id"), 10, 64)

	var slot ComboSlot
	if err := c.BodyParser(&slot); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	slot.ComboID = uint(comboID)

	result, err := h.service.CreateComboSlot(slot)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(result)
}

func (h *PosHandler) UpdateComboSlot(c *fiber.Ctx) error {
	id := c.Params("id")

	var slot ComboSlot
	if err := c.BodyParser(&slot); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	result, err := h.service.UpdateComboSlot(id, slot)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *PosHandler) DeleteComboSlot(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.DeleteComboSlot(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Slot deleted successfully"})
}

// ─────────────────────────────────────────
// COMBO SLOT OPTION
// ─────────────────────────────────────────

func (h *PosHandler) CreateComboSlotOption(c *fiber.Ctx) error {
	slotID, _ := strconv.ParseUint(c.Params("slot_id"), 10, 64)

	var option ComboSlotOption
	if err := c.BodyParser(&option); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	option.SlotID = uint(slotID)

	result, err := h.service.CreateComboSlotOption(option)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(result)
}

func (h *PosHandler) DeleteComboSlotOption(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.DeleteComboSlotOption(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Option deleted successfully"})
}

// ─────────────────────────────────────────
// ORDER & KDS
// ─────────────────────────────────────────

func (h *PosHandler) CreateOrder(c *fiber.Ctx) error {
	var order Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	newOrder, err := h.service.PlaceOrder(order)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(newOrder)
}

func (h *PosHandler) GetActiveOrders(c *fiber.Ctx) error {
	orders, err := h.service.GetKdsOrders()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data KDS"})
	}
	return c.JSON(orders)
}

func (h *PosHandler) AcceptOrder(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	err := h.service.AcceptOrder(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menerima pesanan"})
	}
	return c.JSON(fiber.Map{"message": "Pesanan sedang dimasak"})
}

func (h *PosHandler) FinishOrder(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	err := h.service.FinishOrder(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyelesaikan pesanan"})
	}
	return c.JSON(fiber.Map{"message": "Pesanan siap disajikan!"})
}

func (h *PosHandler) GetAverageCookingTime(c *fiber.Ctx) error {
	avgTime, err := h.service.GetAverageCookingTime()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"average_cooking_time_minutes": avgTime})
}

// ─────────────────────────────────────────
// SPICY LEVEL HANDLERS
// ─────────────────────────────────────────

func (h *PosHandler) GetSpicyLevels(c *fiber.Ctx) error {
	levels, err := h.service.GetAllSpicyLevels()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(levels)
}

func (h *PosHandler) GetSpicyLevelByID(c *fiber.Ctx) error {
	id := c.Params("id")
	level, err := h.service.GetSpicyLevelByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Spicy level not found"})
	}
	return c.JSON(level)
}

func (h *PosHandler) CreateSpicyLevel(c *fiber.Ctx) error {
	var level SpicyLevel
	if err := c.BodyParser(&level); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	result, err := h.service.CreateSpicyLevel(level)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(result)
}

func (h *PosHandler) UpdateSpicyLevel(c *fiber.Ctx) error {
	id := c.Params("id")
	var level SpicyLevel
	if err := c.BodyParser(&level); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	result, err := h.service.UpdateSpicyLevel(id, level)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *PosHandler) DeleteSpicyLevel(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.DeleteSpicyLevel(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Spicy level deleted successfully"})
}

func (h *PosHandler) AssignSpicyLevelToProduct(c *fiber.Ctx) error {
	productID, err := c.ParamsInt("product_id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid product_id"})
	}
	levelID, err := c.ParamsInt("level_id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid level_id"})
	}
	err = h.service.AssignSpicyLevelToProduct(uint(productID), uint(levelID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Spicy level assigned to product"})
}

func (h *PosHandler) RemoveSpicyLevelFromProduct(c *fiber.Ctx) error {
	productID, err := c.ParamsInt("product_id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid product_id"})
	}
	levelID, err := c.ParamsInt("level_id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid level_id"})
	}
	err = h.service.RemoveSpicyLevelFromProduct(uint(productID), uint(levelID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Spicy level removed from product"})
}

func (h *PosHandler) GetProductSpicyLevels(c *fiber.Ctx) error {
	productID, err := c.ParamsInt("product_id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid product_id"})
	}
	levels, err := h.service.GetProductSpicyLevels(uint(productID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(levels)
}

// TIME MENU HANDLERS
func (h *PosHandler) GetTimeMenus(c *fiber.Ctx) error {
    menus, err := h.service.GetAllTimeMenus()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(menus)
}
func (h *PosHandler) GetTimeMenuByID(c *fiber.Ctx) error {
    id := c.Params("id")
    menu, err := h.service.GetTimeMenuByID(id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Time menu not found"})
    }
    return c.JSON(menu)
}
func (h *PosHandler) CreateTimeMenu(c *fiber.Ctx) error {
    var tm TimeMenu
    if err := c.BodyParser(&tm); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
    }
    result, err := h.service.CreateTimeMenu(tm)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.Status(201).JSON(result)
}
func (h *PosHandler) UpdateTimeMenu(c *fiber.Ctx) error {
    id := c.Params("id")
    var tm TimeMenu
    if err := c.BodyParser(&tm); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
    }
    result, err := h.service.UpdateTimeMenu(id, tm)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(result)
}
func (h *PosHandler) DeleteTimeMenu(c *fiber.Ctx) error {
    id := c.Params("id")
    err := h.service.DeleteTimeMenu(id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"message": "Time menu deleted"})
}
func (h *PosHandler) GetProductTimeMenus(c *fiber.Ctx) error {
    productID, err := c.ParamsInt("product_id")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid product_id"})
    }
    menus, err := h.service.GetTimeMenusByProduct(uint(productID))
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(menus)
}
func (h *PosHandler) AssignTimeSlotToProduct(c *fiber.Ctx) error {
    productID, err := c.ParamsInt("product_id")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid product_id"})
    }
    var body struct {
        TimeSlot string `json:"time_slot"`
    }
    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
    }
    err = h.service.AssignTimeSlotToProduct(uint(productID), body.TimeSlot)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"message": "Time slot assigned"})
}
func (h *PosHandler) RemoveTimeSlotFromProduct(c *fiber.Ctx) error {
    id := c.Params("id")
    err := h.service.RemoveTimeSlotFromProduct(id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"message": "Time slot removed"})
}