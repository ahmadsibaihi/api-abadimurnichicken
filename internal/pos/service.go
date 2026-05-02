package pos

import (
	"time"
)

// INTERFACE
type PosService interface {
	// Order
	PlaceOrder(order Order) (Order, error)
	GetKdsOrders() ([]Order, error)
	AcceptOrder(id uint) error
	FinishOrder(id uint) error
	GetAverageCookingTime() (float64, error)
	ChangeStatus(id uint, status string) error

	// Product
	AddProduct(p Product) (Product, error)
	GetAllProducts() ([]Product, error)
	GetProductByID(id string) (Product, error)
	UpdateProduct(id string, p Product) (Product, error)
	DeleteProduct(id string) error

	// Category
	CreateCategory(c Category) (Category, error)
	GetAllCategories() ([]Category, error)
	GetCategoryByID(id string) (Category, error)
	UpdateCategory(id string, c Category) (Category, error)
	DeleteCategory(id string) error

	// Variant
	CreateVariant(v ProductVariant) (ProductVariant, error)
	GetVariantsByProductID(productID string) ([]ProductVariant, error)
	UpdateVariant(id string, v ProductVariant) (ProductVariant, error)
	DeleteVariant(id string) error

	// Addon
	CreateAddon(a Addon) (Addon, error)
	GetAllAddons() ([]Addon, error)
	UpdateAddon(id string, a Addon) (Addon, error)
	DeleteAddon(id string) error
	AssignAddonToProduct(productID, addonID uint) error
	RemoveAddonFromProduct(productID, addonID uint) error

	// Combo Package
	CreateCombo(c ComboPackage) (ComboPackage, error)
	GetAllCombos() ([]ComboPackage, error)
	GetComboByID(id string) (ComboPackage, error)
	UpdateCombo(id string, c ComboPackage) (ComboPackage, error)
	DeleteCombo(id string) error

	// Combo Slot
	CreateComboSlot(s ComboSlot) (ComboSlot, error)
	UpdateComboSlot(id string, s ComboSlot) (ComboSlot, error)
	DeleteComboSlot(id string) error

	// Combo Slot Option
	CreateComboSlotOption(o ComboSlotOption) (ComboSlotOption, error)
	DeleteComboSlotOption(id string) error

	// Spicy Level
	GetAllSpicyLevels() ([]SpicyLevel, error)
	GetSpicyLevelByID(id string) (SpicyLevel, error)
	CreateSpicyLevel(level SpicyLevel) (SpicyLevel, error)
	UpdateSpicyLevel(id string, level SpicyLevel) (SpicyLevel, error)
	DeleteSpicyLevel(id string) error
	AssignSpicyLevelToProduct(productID, levelID uint) error
	RemoveSpicyLevelFromProduct(productID, levelID uint) error
	GetProductSpicyLevels(productID uint) ([]SpicyLevel, error)

	// TimeMenu
	GetAllTimeMenus() ([]TimeMenu, error)
	GetTimeMenuByID(id string) (TimeMenu, error)
	CreateTimeMenu(tm TimeMenu) (TimeMenu, error)
	UpdateTimeMenu(id string, tm TimeMenu) (TimeMenu, error)
	DeleteTimeMenu(id string) error
	GetTimeMenusByProduct(productID uint) ([]TimeMenu, error)
	AssignTimeSlotToProduct(productID uint, timeSlot string) error
	RemoveTimeSlotFromProduct(id string) error
}

// SERVICE STRUCT & CONSTRUCTOR

type posService struct {
	repo Repository
}

func NewService(r Repository) PosService {
	return &posService{repo: r}
}

// SPICY LEVEL IMPLEMENTATIONS

func (s *posService) GetAllSpicyLevels() ([]SpicyLevel, error) {
	return s.repo.FindAllSpicyLevels()
}

func (s *posService) GetSpicyLevelByID(id string) (SpicyLevel, error) {
	return s.repo.FindSpicyLevelByID(id)
}

func (s *posService) CreateSpicyLevel(level SpicyLevel) (SpicyLevel, error) {
	return s.repo.CreateSpicyLevel(level)
}

func (s *posService) UpdateSpicyLevel(id string, level SpicyLevel) (SpicyLevel, error) {
	return s.repo.UpdateSpicyLevel(id, level)
}

func (s *posService) DeleteSpicyLevel(id string) error {
	return s.repo.DeleteSpicyLevel(id)
}

func (s *posService) AssignSpicyLevelToProduct(productID, levelID uint) error {
	return s.repo.AssignSpicyLevel(productID, levelID)
}

func (s *posService) RemoveSpicyLevelFromProduct(productID, levelID uint) error {
	return s.repo.RemoveSpicyLevel(productID, levelID)
}

func (s *posService) GetProductSpicyLevels(productID uint) ([]SpicyLevel, error) {
	return s.repo.FindSpicyLevelsByProduct(productID)
}

// ORDER IMPLEMENTATIONS

func (s *posService) PlaceOrder(order Order) (Order, error) {
	order.Status = "pending"
	return s.repo.Create(order)
}

func (s *posService) GetKdsOrders() ([]Order, error) {
	return s.repo.FindActiveOrders()
}

func (s *posService) AcceptOrder(id uint) error {
	now := time.Now()
	return s.repo.UpdateStatus(id, "cooking", &now, nil)
}

func (s *posService) FinishOrder(id uint) error {
	now := time.Now()
	return s.repo.UpdateStatus(id, "ready", nil, &now)
}

func (s *posService) GetAverageCookingTime() (float64, error) {
	return s.repo.CalculateAvgCookingTime()
}

func (s *posService) ChangeStatus(id uint, status string) error {
	// bisa diimplementasikan nanti sesuai kebutuhan
	return nil
}

// PRODUCT IMPLEMENTATIONS

func (s *posService) AddProduct(p Product) (Product, error) {
	return s.repo.CreateProduct(p)
}

func (s *posService) GetAllProducts() ([]Product, error) {
	return s.repo.FindAllProducts()
}

func (s *posService) GetProductByID(id string) (Product, error) {
	return s.repo.FindProductByID(id)
}

func (s *posService) UpdateProduct(id string, p Product) (Product, error) {
	return s.repo.UpdateProduct(id, p)
}

func (s *posService) DeleteProduct(id string) error {
	return s.repo.DeleteProduct(id)
}

// TIME MENU IMPLEMENTATIONS
func (s *posService) GetAllTimeMenus() ([]TimeMenu, error) {
    return s.repo.FindAllTimeMenus()
}
func (s *posService) GetTimeMenuByID(id string) (TimeMenu, error) {
    return s.repo.FindTimeMenuByID(id)
}
func (s *posService) CreateTimeMenu(tm TimeMenu) (TimeMenu, error) {
    return s.repo.CreateTimeMenu(tm)
}
func (s *posService) UpdateTimeMenu(id string, tm TimeMenu) (TimeMenu, error) {
    return s.repo.UpdateTimeMenu(id, tm)
}
func (s *posService) DeleteTimeMenu(id string) error {
    return s.repo.DeleteTimeMenu(id)
}
func (s *posService) GetTimeMenusByProduct(productID uint) ([]TimeMenu, error) {
    return s.repo.FindTimeMenusByProduct(productID)
}
func (s *posService) AssignTimeSlotToProduct(productID uint, timeSlot string) error {
    tm := TimeMenu{ProductID: productID, TimeSlot: timeSlot}
    _, err := s.repo.CreateTimeMenu(tm)
    return err
}
func (s *posService) RemoveTimeSlotFromProduct(id string) error {
    return s.repo.DeleteTimeMenu(id)
}

// CATEGORY IMPLEMENTATIONS

func (s *posService) CreateCategory(c Category) (Category, error) {
	return s.repo.CreateCategory(c)
}

func (s *posService) GetAllCategories() ([]Category, error) {
	return s.repo.FindAllCategories()
}

func (s *posService) GetCategoryByID(id string) (Category, error) {
	return s.repo.FindCategoryByID(id)
}

func (s *posService) UpdateCategory(id string, c Category) (Category, error) {
	return s.repo.UpdateCategory(id, c)
}

func (s *posService) DeleteCategory(id string) error {
	return s.repo.DeleteCategory(id)
}

// VARIANT IMPLEMENTATIONS

func (s *posService) CreateVariant(v ProductVariant) (ProductVariant, error) {
	return s.repo.CreateVariant(v)
}

func (s *posService) GetVariantsByProductID(productID string) ([]ProductVariant, error) {
	return s.repo.FindVariantsByProductID(productID)
}

func (s *posService) UpdateVariant(id string, v ProductVariant) (ProductVariant, error) {
	return s.repo.UpdateVariant(id, v)
}

func (s *posService) DeleteVariant(id string) error {
	return s.repo.DeleteVariant(id)
}

// ADDON IMPLEMENTATIONS

func (s *posService) CreateAddon(a Addon) (Addon, error) {
	return s.repo.CreateAddon(a)
}

func (s *posService) GetAllAddons() ([]Addon, error) {
	return s.repo.FindAllAddons()
}

func (s *posService) UpdateAddon(id string, a Addon) (Addon, error) {
	return s.repo.UpdateAddon(id, a)
}

func (s *posService) DeleteAddon(id string) error {
	return s.repo.DeleteAddon(id)
}

func (s *posService) AssignAddonToProduct(productID, addonID uint) error {
	return s.repo.AssignAddonToProduct(productID, addonID)
}

func (s *posService) RemoveAddonFromProduct(productID, addonID uint) error {
	return s.repo.RemoveAddonFromProduct(productID, addonID)
}

// COMBO PACKAGE IMPLEMENTATIONS

func (s *posService) CreateCombo(c ComboPackage) (ComboPackage, error) {
	return s.repo.CreateCombo(c)
}

func (s *posService) GetAllCombos() ([]ComboPackage, error) {
	return s.repo.FindAllCombos()
}

func (s *posService) GetComboByID(id string) (ComboPackage, error) {
	return s.repo.FindComboByID(id)
}

func (s *posService) UpdateCombo(id string, c ComboPackage) (ComboPackage, error) {
	return s.repo.UpdateCombo(id, c)
}

func (s *posService) DeleteCombo(id string) error {
	return s.repo.DeleteCombo(id)
}

// COMBO SLOT IMPLEMENTATIONS

func (s *posService) CreateComboSlot(sl ComboSlot) (ComboSlot, error) {
	return s.repo.CreateComboSlot(sl)
}

func (s *posService) UpdateComboSlot(id string, sl ComboSlot) (ComboSlot, error) {
	return s.repo.UpdateComboSlot(id, sl)
}

func (s *posService) DeleteComboSlot(id string) error {
	return s.repo.DeleteComboSlot(id)
}

// COMBO SLOT OPTION IMPLEMENTATIONS

func (s *posService) CreateComboSlotOption(o ComboSlotOption) (ComboSlotOption, error) {
	return s.repo.CreateComboSlotOption(o)
}

func (s *posService) DeleteComboSlotOption(id string) error {
	return s.repo.DeleteComboSlotOption(id)
}