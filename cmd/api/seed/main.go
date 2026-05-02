package main

import (
    "fmt"
    "go-serve-pos/database"
    "go-serve-pos/internal/user"
)

func main() {
    database.ConnectDB()
    var count int64
    database.DB.Model(&user.User{}).Count(&count)
    if count == 0 {
        superAdmin := user.User{
            Name:     "Super Admin",
            Email:    "super@kasir.com",
            Password: "admin123",
            Role:     "superadmin",
            IsActive: true,
        }
        superAdmin.HashPassword()
        database.DB.Create(&superAdmin)
        fmt.Println("Superadmin created: super@kasir.com / admin123")
    } else {
        fmt.Println("User already exists, skip seeding.")
    }
}