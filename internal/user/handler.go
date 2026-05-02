package user

import (
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
	"strconv"
    "golang.org/x/crypto/bcrypt"
)

// GetKaryawan mengembalikan fungsi handler yang membutuhkan koneksi DB
func GetKaryawan(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var users []User
        err := db.Model(&User{}).
            Where("role = ?", "karyawan").
            Select("id, name, email, role").
            Find(&users).Error
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data karyawan"})
        }
        return c.JSON(users)
    }
}

func CreateUser(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var req struct {
            Name     string `json:"name"`
            Email    string `json:"email"`
            Password string `json:"password"`
            Role     string `json:"role"` // "karyawan", "admin", "superadmin"
        }
        if err := c.BodyParser(&req); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
        }
        // Validasi role yang diperbolehkan (admin/superadmin hanya bisa buat karyawan dan admin, bukan superadmin)
        userRole := c.Locals("user_role").(string)
        if userRole != "superadmin" && req.Role == "superadmin" {
            return c.Status(403).JSON(fiber.Map{"error": "Only superadmin can create superadmin"})
        }
        hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
        user := User{
            Name:     req.Name,
            Email:    req.Email,
            Password: string(hashed),
            Role:     req.Role,
            IsActive: true,
        }
        if err := db.Create(&user).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Email already exists"})
        }
        return c.Status(201).JSON(fiber.Map{"message": "User created", "user": user})
    }
}

// GetAllUsers - daftar semua user (khusus admin/superadmin)
func GetAllUsers(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var users []User
        db.Select("id, name, email, role, is_active, created_at").Find(&users)
        return c.JSON(users)
    }
}

// UpdateUser - edit nama, email, role, status
func UpdateUser(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
        var req struct {
            Name     string `json:"name"`
            Email    string `json:"email"`
            Role     string `json:"role"`
            IsActive *bool  `json:"is_active"`
        }
        if err := c.BodyParser(&req); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
        }
        updates := map[string]interface{}{}
        if req.Name != "" { updates["name"] = req.Name }
        if req.Email != "" { updates["email"] = req.Email }
        if req.Role != "" { updates["role"] = req.Role }
        if req.IsActive != nil { updates["is_active"] = *req.IsActive }
        if err := db.Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Update failed"})
        }
        return c.JSON(fiber.Map{"message": "User updated"})
    }
}

// DeleteUser - soft delete? Atau hard delete. Kita hard delete saja.
func DeleteUser(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
        if err := db.Delete(&User{}, id).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Delete failed"})
        }
        return c.JSON(fiber.Map{"message": "User deleted"})
    }
}