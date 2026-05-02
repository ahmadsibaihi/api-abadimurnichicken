package auth

import (
    "go-serve-pos/database"
    "go-serve-pos/internal/user"
    "github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type RegisterRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Role     string `json:"role"` // "admin" or "superadmin"
}

func Login(c *fiber.Ctx) error {
    var req LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    var user user.User
    if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
    }

    if !user.CheckPassword(req.Password) {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
    }

    if !user.IsActive {
        return c.Status(403).JSON(fiber.Map{"error": "Account disabled"})
    }

    token, err := GenerateToken(user.ID, user.Email, user.Role)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
    }

    return c.JSON(fiber.Map{
        "token": token,
        "user": fiber.Map{
            "id":    user.ID,
            "name":  user.Name,
            "email": user.Email,
            "role":  user.Role,
        },
    })
}

func Register(c *fiber.Ctx) error {
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    // Hash password
    newUser := user.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
        Role:     req.Role,
        IsActive: true,
    }
    if err := newUser.HashPassword(); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
    }

    if err := database.DB.Create(&newUser).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Email already exists or invalid data"})
    }

    return c.Status(201).JSON(fiber.Map{"message": "User created successfully"})
}