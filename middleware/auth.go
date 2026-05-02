package middleware

import (
    "strings"
    "go-serve-pos/internal/auth"
    "github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(401).JSON(fiber.Map{"error": "Missing token"})
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.Status(401).JSON(fiber.Map{"error": "Invalid token format"})
        }

        claims, err := auth.ValidateToken(parts[1])
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
        }

        c.Locals("user_id", claims.UserID)
        c.Locals("user_email", claims.Email)
        c.Locals("user_role", claims.Role)
        return c.Next()
    }
}

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        role := c.Locals("user_role").(string)
        for _, r := range allowedRoles {
            if r == role {
                return c.Next()
            }
        }
        return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
    }
}