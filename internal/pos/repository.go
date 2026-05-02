package pos

import (
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	// Order
	Create(order Order) (Order, error)
	FindActiveOrders() ([]Order, error)
	UpdateStatus(id uint, status string, start *time.Time, finish *time.Time) error
	CalculateAvgCookingTime() (float64, error)

	// Product
	CreateProduct(p Product) (Product, error)
	FindAllProducts() ([]Product, error)
	FindProductByID(id string) (Product, error)
	UpdateProduct(id string, p Product) (Product, error)
	DeleteProduct(id string) error

	// Category
	CreateCategory(c Category) (Category, error)
	FindAllCategories() ([]Category, error)
	FindCategoryByID(id string) (Category, error)
	UpdateCategory(id string, c Category) (Category, error)
	DeleteCategory(id string) error

	// Variant
	CreateVariant(v ProductVariant) (ProductVariant, error)
	FindVariantsByProductID(productID string) ([]ProductVariant, error)
	FindVariantByID(id string) (ProductVariant, error)
	UpdateVariant(id string, v ProductVariant) (ProductVariant, error)
	DeleteVariant(id string) error

	// Addon
	CreateAddon(a Addon) (Addon, error)
	FindAllAddons() ([]Addon, error)
	FindAddonByID(id string) (Addon, error)
	UpdateAddon(id string, a Addon) (Addon, error)
	DeleteAddon(id string) error
	AssignAddonToProduct(productID, addonID uint) error
	RemoveAddonFromProduct(productID, addonID uint) error

	// Combo
	CreateCombo(c ComboPackage) (ComboPackage, error)
	FindAllCombos() ([]ComboPackage, error)
	FindComboByID(id string) (ComboPackage, error)
	UpdateCombo(id string, c ComboPackage) (ComboPackage, error)
	DeleteCombo(id string) error

	// Combo Slot
	CreateComboSlot(s ComboSlot) (ComboSlot, error)
	UpdateComboSlot(id string, s ComboSlot) (ComboSlot, error)
	DeleteComboSlot(id string) error

	// Combo Slot Option
	CreateComboSlotOption(o ComboSlotOption) (ComboSlotOption, error)
	DeleteComboSlotOption(id string) error

	FindAllSpicyLevels() ([]SpicyLevel, error)
    FindSpicyLevelByID(id string) (SpicyLevel, error)
    CreateSpicyLevel(level SpicyLevel) (SpicyLevel, error)
    UpdateSpicyLevel(id string, level SpicyLevel) (SpicyLevel, error)
    DeleteSpicyLevel(id string) error
    AssignSpicyLevel(productID, levelID uint) error
    RemoveSpicyLevel(productID, levelID uint) error
    FindSpicyLevelsByProduct(productID uint) ([]SpicyLevel, error)

	FindAllTimeMenus() ([]TimeMenu, error)
	FindTimeMenuByID(id string) (TimeMenu, error)
	CreateTimeMenu(tm TimeMenu) (TimeMenu, error)
	UpdateTimeMenu(id string, tm TimeMenu) (TimeMenu, error)
	DeleteTimeMenu(id string) error
	FindTimeMenusByProduct(productID uint) ([]TimeMenu, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

// ─────────────────────────────────────────
// ORDER
// ─────────────────────────────────────────

func (r *repository) Create(order Order) (Order, error) {
	err := r.db.Create(&order).Error
	return order, err
}

func (r *repository) FindActiveOrders() ([]Order, error) {
	var orders []Order
	err := r.db.Preload("Items.Product").
		Where("status IN ?", []string{"pending", "cooking"}).
		Find(&orders).Error
	return orders, err
}

func (r *repository) UpdateStatus(id uint, status string, start *time.Time, finish *time.Time) error {
	updates := map[string]interface{}{"status": status}
	if start != nil {
		updates["started_at"] = start
	}
	if finish != nil {
		updates["finished_at"] = finish
	}
	return r.db.Model(&Order{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repository) CalculateAvgCookingTime() (float64, error) {
	var avgMinutes float64
	err := r.db.Raw(`
		SELECT AVG(TIMESTAMPDIFF(MINUTE, started_at, finished_at)) 
		FROM orders WHERE finished_at IS NOT NULL
	`).Scan(&avgMinutes).Error
	return avgMinutes, err
}

// ─────────────────────────────────────────
// PRODUCT
// ─────────────────────────────────────────

func (r *repository) CreateProduct(p Product) (Product, error) {
	err := r.db.Create(&p).Error
	return p, err
}

func (r *repository) FindAllProducts() ([]Product, error) {
	var products []Product
	err := r.db.Preload("Category").Preload("Variants").Find(&products).Error
	return products, err
}

func (r *repository) FindProductByID(id string) (Product, error) {
	var product Product
	err := r.db.Preload("Category").Preload("Variants").Preload("Addons").
		Where("id = ?", id).First(&product).Error
	return product, err
}

func (r *repository) UpdateProduct(id string, p Product) (Product, error) {
	updates := map[string]interface{}{}
	if p.Name != "" {
		updates["name"] = p.Name
	}
	if p.Price != 0 {
		updates["price"] = p.Price
	}
	if p.Stock != 0 {
		updates["stock"] = p.Stock
	}
	if p.CookingTime != 0 {
		updates["cooking_time"] = p.CookingTime
	}
	if p.Image != "" {
		updates["image"] = p.Image
	}
	if p.Description != "" {
		updates["description"] = p.Description
	}
	if p.CategoryID != 0 {
		updates["category_id"] = p.CategoryID
	}
	updates["is_active"] = p.IsActive
	updates["is_best_seller"] = p.IsBestSeller
	updates["is_new"] = p.IsNew

	err := r.db.Model(&Product{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return Product{}, err
	}
	return r.FindProductByID(id)
}

func (r *repository) DeleteProduct(id string) error {
	return r.db.Delete(&Product{}, "id = ?", id).Error
}

// ─────────────────────────────────────────
// CATEGORY
// ─────────────────────────────────────────

func (r *repository) CreateCategory(c Category) (Category, error) {
	err := r.db.Create(&c).Error
	return c, err
}

func (r *repository) FindAllCategories() ([]Category, error) {
	var categories []Category
	err := r.db.Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

func (r *repository) FindCategoryByID(id string) (Category, error) {
	var category Category
	err := r.db.Where("id = ?", id).First(&category).Error
	return category, err
}

func (r *repository) UpdateCategory(id string, c Category) (Category, error) {
	updates := map[string]interface{}{}
	if c.Name != "" {
		updates["name"] = c.Name
	}
	if c.Icon != "" {
		updates["icon"] = c.Icon
	}
	updates["sort_order"] = c.SortOrder
	updates["is_active"] = c.IsActive

	err := r.db.Model(&Category{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return Category{}, err
	}
	return r.FindCategoryByID(id)
}

func (r *repository) DeleteCategory(id string) error {
	return r.db.Delete(&Category{}, "id = ?", id).Error
}

// ─────────────────────────────────────────
// VARIANT
// ─────────────────────────────────────────

func (r *repository) CreateVariant(v ProductVariant) (ProductVariant, error) {
	err := r.db.Create(&v).Error
	return v, err
}

func (r *repository) FindVariantsByProductID(productID string) ([]ProductVariant, error) {
	var variants []ProductVariant
	err := r.db.Where("product_id = ?", productID).Find(&variants).Error
	return variants, err
}

func (r *repository) FindVariantByID(id string) (ProductVariant, error) {
	var variant ProductVariant
	err := r.db.Where("id = ?", id).First(&variant).Error
	return variant, err
}

func (r *repository) UpdateVariant(id string, v ProductVariant) (ProductVariant, error) {
	updates := map[string]interface{}{}
	if v.Name != "" {
		updates["name"] = v.Name
	}
	updates["additional_price"] = v.AdditionalPrice
	updates["is_active"] = v.IsActive

	err := r.db.Model(&ProductVariant{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return ProductVariant{}, err
	}
	return r.FindVariantByID(id)
}

func (r *repository) DeleteVariant(id string) error {
	return r.db.Delete(&ProductVariant{}, "id = ?", id).Error
}

// ─────────────────────────────────────────
// ADDON
// ─────────────────────────────────────────

func (r *repository) CreateAddon(a Addon) (Addon, error) {
	err := r.db.Create(&a).Error
	return a, err
}

func (r *repository) FindAllAddons() ([]Addon, error) {
	var addons []Addon
	err := r.db.Find(&addons).Error
	return addons, err
}

func (r *repository) FindAddonByID(id string) (Addon, error) {
	var addon Addon
	err := r.db.Where("id = ?", id).First(&addon).Error
	return addon, err
}

func (r *repository) UpdateAddon(id string, a Addon) (Addon, error) {
	updates := map[string]interface{}{}
	if a.Name != "" {
		updates["name"] = a.Name
	}
	updates["price"] = a.Price
	updates["is_active"] = a.IsActive

	err := r.db.Model(&Addon{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return Addon{}, err
	}
	return r.FindAddonByID(id)
}

func (r *repository) DeleteAddon(id string) error {
	return r.db.Delete(&Addon{}, "id = ?", id).Error
}

func (r *repository) AssignAddonToProduct(productID, addonID uint) error {
	pa := ProductAddon{ProductID: productID, AddonID: addonID}
	return r.db.Where(pa).FirstOrCreate(&pa).Error
}

func (r *repository) RemoveAddonFromProduct(productID, addonID uint) error {
	return r.db.Where("product_id = ? AND addon_id = ?", productID, addonID).
		Delete(&ProductAddon{}).Error
}

// ─────────────────────────────────────────
// COMBO
// ─────────────────────────────────────────

func (r *repository) CreateCombo(c ComboPackage) (ComboPackage, error) {
	err := r.db.Create(&c).Error
	return c, err
}

func (r *repository) FindAllCombos() ([]ComboPackage, error) {
	var combos []ComboPackage
	err := r.db.Preload("Slots.Options.Product").Find(&combos).Error
	return combos, err
}

func (r *repository) FindComboByID(id string) (ComboPackage, error) {
	var combo ComboPackage
	err := r.db.Preload("Slots.Options.Product").
		Where("id = ?", id).First(&combo).Error
	return combo, err
}

func (r *repository) UpdateCombo(id string, c ComboPackage) (ComboPackage, error) {
	updates := map[string]interface{}{}
	if c.Name != "" {
		updates["name"] = c.Name
	}
	if c.Description != "" {
		updates["description"] = c.Description
	}
	if c.BasePrice != 0 {
		updates["base_price"] = c.BasePrice
	}
	if c.Image != "" {
		updates["image"] = c.Image
	}
	updates["is_active"] = c.IsActive

	err := r.db.Model(&ComboPackage{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return ComboPackage{}, err
	}
	return r.FindComboByID(id)
}

func (r *repository) DeleteCombo(id string) error {
	return r.db.Delete(&ComboPackage{}, "id = ?", id).Error
}

// ─────────────────────────────────────────
// COMBO SLOT
// ─────────────────────────────────────────

func (r *repository) CreateComboSlot(s ComboSlot) (ComboSlot, error) {
	err := r.db.Create(&s).Error
	return s, err
}

func (r *repository) UpdateComboSlot(id string, s ComboSlot) (ComboSlot, error) {
	updates := map[string]interface{}{}
	if s.SlotName != "" {
		updates["slot_name"] = s.SlotName
	}
	updates["is_required"] = s.IsRequired

	err := r.db.Model(&ComboSlot{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return ComboSlot{}, err
	}

	var slot ComboSlot
	r.db.Where("id = ?", id).First(&slot)
	return slot, nil
}

func (r *repository) DeleteComboSlot(id string) error {
	return r.db.Delete(&ComboSlot{}, "id = ?", id).Error
}

// ─────────────────────────────────────────
// COMBO SLOT OPTION
// ─────────────────────────────────────────

func (r *repository) CreateComboSlotOption(o ComboSlotOption) (ComboSlotOption, error) {
	err := r.db.Create(&o).Error
	return o, err
}

func (r *repository) DeleteComboSlotOption(id string) error {
	return r.db.Delete(&ComboSlotOption{}, "id = ?", id).Error
}

// ──────────────────────────────────────────────────────────
// SPICY LEVEL REPOSITORY METHODS
// ──────────────────────────────────────────────────────────

func (r *repository) FindAllSpicyLevels() ([]SpicyLevel, error) {
	var levels []SpicyLevel
	err := r.db.Order("level asc").Find(&levels).Error
	return levels, err
}

func (r *repository) FindSpicyLevelByID(id string) (SpicyLevel, error) {
	var level SpicyLevel
	err := r.db.First(&level, id).Error
	return level, err
}

func (r *repository) CreateSpicyLevel(level SpicyLevel) (SpicyLevel, error) {
	err := r.db.Create(&level).Error
	return level, err
}

func (r *repository) UpdateSpicyLevel(id string, level SpicyLevel) (SpicyLevel, error) {
	err := r.db.Model(&SpicyLevel{}).Where("id = ?", id).Updates(level).Error
	return level, err
}

func (r *repository) DeleteSpicyLevel(id string) error {
	return r.db.Delete(&SpicyLevel{}, id).Error
}

// AssignSpicyLevelToProduct (pivot table)
func (r *repository) AssignSpicyLevel(productID, levelID uint) error {
	ps := ProductSpicy{ProductID: productID, SpicyLevelID: levelID}
	return r.db.FirstOrCreate(&ps, ps).Error
}

func (r *repository) RemoveSpicyLevel(productID, levelID uint) error {
	return r.db.Where("product_id = ? AND spicy_level_id = ?", productID, levelID).
		Delete(&ProductSpicy{}).Error
}

func (r *repository) FindSpicyLevelsByProduct(productID uint) ([]SpicyLevel, error) {
	var levels []SpicyLevel
	err := r.db.Table("spicy_levels").
		Joins("JOIN product_spicy ON product_spicy.spicy_level_id = spicy_levels.id").
		Where("product_spicy.product_id = ?", productID).
		Order("spicy_levels.level asc").
		Find(&levels).Error
	return levels, err
}

func (r *repository) FindAllTimeMenus() ([]TimeMenu, error) {
    var menus []TimeMenu
    err := r.db.Find(&menus).Error
    return menus, err
}
func (r *repository) FindTimeMenuByID(id string) (TimeMenu, error) {
    var tm TimeMenu
    err := r.db.First(&tm, id).Error
    return tm, err
}
func (r *repository) CreateTimeMenu(tm TimeMenu) (TimeMenu, error) {
    err := r.db.Create(&tm).Error
    return tm, err
}
func (r *repository) UpdateTimeMenu(id string, tm TimeMenu) (TimeMenu, error) {
    err := r.db.Model(&TimeMenu{}).Where("id = ?", id).Updates(tm).Error
    return tm, err
}
func (r *repository) DeleteTimeMenu(id string) error {
    return r.db.Delete(&TimeMenu{}, id).Error
}
func (r *repository) FindTimeMenusByProduct(productID uint) ([]TimeMenu, error) {
    var menus []TimeMenu
    err := r.db.Where("product_id = ?", productID).Find(&menus).Error
    return menus, err
}