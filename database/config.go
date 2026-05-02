package database

import (
    "fmt"
    "go-serve-pos/internal/expenditure" // <-- tambahkan import ini
    "go-serve-pos/internal/user"   
    "go-serve-pos/internal/pos"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
    dsn := "root:root@tcp(127.0.0.1:8889)/pos-kds?charset=utf8mb4&parseTime=True&loc=Local"
    
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
    fmt.Println("Berhasil konek ke MAMP MySQL! 🚀")
}