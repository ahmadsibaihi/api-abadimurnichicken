package pos

import (
	"time"
)

// CATEGORY
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`

	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

// PRODUCT
type Product struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	CategoryID   uint    `json:"category_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Stock        int     `json:"stock"`
	CookingTime  int     `json:"cooking_time"` // dalam menit
	Image        string  `json:"image"`
	IsActive     bool    `gorm:"default:true" json:"is_active"`
	IsBestSeller bool    `gorm:"default:false" json:"is_best_seller"`
	IsNew        bool    `gorm:"default:false" json:"is_new"`
	CreatedAt    time.Time `json:"created_at"`

	Category  Category         `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Variants  []ProductVariant `json:"variants,omitempty" gorm:"foreignKey:ProductID"`
	Addons    []Addon          `json:"addons,omitempty" gorm:"many2many:product_addons"`
	HasSpicy  bool             `json:"-" gorm:"-"` // diisi manual dari product_spicy
}

// PRODUCT VARIANT (Sambal Ijo, Matah, dll)
type ProductVariant struct {
	ID              uint    `gorm:"primaryKey" json:"id"`
	ProductID       uint    `json:"product_id"`
	Name            string  `json:"name"`            // "Sambal Ijo", "Tepung Pedes"
	AdditionalPrice float64 `json:"additional_price"` // harga tambahan, 0 jika sama
	IsActive        bool    `gorm:"default:true" json:"is_active"`
}

// ADDON / TOPPING
type Addon struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Name     string  `json:"name"`  // "Extra Keju", "Extra Sambal"
	Price    float64 `json:"price"`
	IsActive bool    `gorm:"default:true" json:"is_active"`
}

// Tabel pivot product_addons (many2many)
type ProductAddon struct {
	ProductID uint `json:"product_id"`
	AddonID   uint `json:"addon_id"`
}

// SPICY LEVEL
type SpicyLevel struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Level int    `json:"level"` 
	Label string `json:"label"` 
	Emoji string `json:"emoji"` 
}
type ProductSpicy struct {
    ProductID    uint `gorm:"primaryKey"`
    SpicyLevelID uint `gorm:"primaryKey"`
}
type TimeMenu struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ProductID uint   `json:"product_id"`
	TimeSlot  string `json:"time_slot"` // breakfast/lunch/dinner
}

// COMBO PACKAGE
type ComboPackage struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BasePrice   float64 `json:"base_price"`
	Image       string  `json:"image"`
	IsActive    bool    `gorm:"default:true" json:"is_active"`

	Slots []ComboSlot `json:"slots,omitempty" gorm:"foreignKey:ComboID"`
}

type ComboSlot struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	ComboID    uint   `json:"combo_id"`
	SlotName   string `json:"slot_name"`   // "Pilih Ayam", "Pilih Minuman"
	IsRequired bool   `gorm:"default:true" json:"is_required"`

	Options []ComboSlotOption `json:"options,omitempty" gorm:"foreignKey:SlotID"`
}

type ComboSlotOption struct {
	ID              uint    `gorm:"primaryKey" json:"id"`
	SlotID          uint    `json:"slot_id"`
	ProductID       uint    `json:"product_id"`
	AdditionalPrice float64 `json:"additional_price"`

	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

// ORDER
type Order struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	TableNumber string      `json:"table_number"`
	OrderType   string      `json:"order_type"`   // "Dine In" | "Take Away"
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"` // "pending" | "cooking" | "ready" | "served"
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`

	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at"`  // klik "Terima" di KDS
	FinishedAt *time.Time `json:"finished_at"` // klik "Selesai" di KDS
}

// ORDER ITEM
type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	VariantID *uint   `json:"variant_id"` // nullable, jika ada varian
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`    // snapshot harga saat order
	Subtotal  float64 `json:"subtotal"` // quantity * price + addons
	SpicyLevel *int   `json:"spicy_level"` // nullable, 1-10
	Notes     string  `json:"notes"`

	Product  Product        `json:"product" gorm:"foreignKey:ProductID"`
	Variant  *ProductVariant `json:"variant,omitempty" gorm:"foreignKey:VariantID"`
	Addons   []OrderItemAddon `json:"addons,omitempty" gorm:"foreignKey:OrderItemID"`
}

// Add-on yang dipilih per item order
type OrderItemAddon struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	OrderItemID uint    `json:"order_item_id"`
	AddonID     uint    `json:"addon_id"`
	Price       float64 `json:"price"`    // snapshot harga addon saat order
	Quantity    int     `gorm:"default:1" json:"quantity"`

	Addon Addon `json:"addon,omitempty" gorm:"foreignKey:AddonID"`
}