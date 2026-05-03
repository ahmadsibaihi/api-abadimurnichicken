package database

import (
	"fmt"
	"go-serve-pos/internal/expenditure"
	"go-serve-pos/internal/pos"
	"go-serve-pos/internal/user"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func ConnectDB() {
	host := getEnv("DB_HOST", "mysql_container")
	port := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "amc-gobackend")
	dbName := getEnv("DB_NAME", "db_pos-kds")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		password,
		host,
		port,
		dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Gagal konek: %v\n", err)
		panic("Gagal terhubung ke database")
	}

	err = db.AutoMigrate(
		&pos.Category{},
		&pos.Product{},
		&pos.ProductVariant{},
		&pos.Addon{},
		&pos.ProductAddon{},
		&pos.SpicyLevel{},
		&pos.ProductSpicy{},
		&pos.ComboPackage{},
		&pos.ComboSlot{},
		&pos.ComboSlotOption{},
		&pos.Order{},
		&pos.OrderItem{},
		&pos.OrderItemAddon{},
		&pos.TimeMenu{},
		&expenditure.Expenditure{},
		&expenditure.ExpenditureAuditLog{},
		&user.User{},
	)
	if err != nil {
		fmt.Printf("AutoMigrate error: %v\n", err)
		panic("Gagal migrasi database")
	}

	DB = db
	fmt.Println("Berhasil konek ke database! 🚀")
}