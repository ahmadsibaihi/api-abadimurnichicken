package routes

import (
    "os"
    "go-serve-pos/database"
    "go-serve-pos/internal/expenditure"
    "go-serve-pos/internal/pos"
    "go-serve-pos/internal/auth"
	"go-serve-pos/internal/user"  
    "go-serve-pos/middleware"
    "github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
    // Buat folder uploads
    if err := os.MkdirAll("./uploads/expenditures", 0755); err != nil {
        panic("Gagal membuat folder uploads: " + err.Error())
    }
    app.Static("/uploads", "./uploads")

    // Auth routes (tanpa middleware)
    authGroup := app.Group("/api/auth")
    authGroup.Post("/login", auth.Login)
    authGroup.Post("/register", auth.Register)

    // API v1 dengan auth middleware
    v1 := app.Group("/api/v1")
    v1.Use(middleware.AuthMiddleware())

    // POS existing
    posRepo := pos.NewRepository(database.DB)
    posService := pos.NewService(posRepo)
    posHandler := pos.NewHandler(posService)

    CategoryRoutes(v1, posHandler)
    VariantRoutes(v1, posHandler)
    AddonRoutes(v1, posHandler)
    ComboRoutes(v1, posHandler)
    SpicyLevelRoutes(v1, posHandler)
    TimeMenuRoutes(v1, posHandler)

    // Expenditure - hanya admin & superadmin
    expRepo := expenditure.NewRepository(database.DB)
    expService := expenditure.NewService(expRepo)
    expHandler := expenditure.NewHandler(expService)

    expGroup := v1.Group("/expenditures")
    expGroup.Use(middleware.RoleMiddleware("admin", "superadmin"))
    expenditure.ExpenditureRoutes(expGroup, expHandler)

    // Get list karyawan - hanya admin & superadmin
    userGroup := v1.Group("/users")
	userGroup.Use(middleware.RoleMiddleware("admin", "superadmin"))
	userGroup.Get("/", user.GetAllUsers(database.DB))
	userGroup.Post("/", user.CreateUser(database.DB))
	userGroup.Put("/:id", user.UpdateUser(database.DB))
	userGroup.Delete("/:id", user.DeleteUser(database.DB))
	userGroup.Get("/karyawan", user.GetKaryawan(database.DB)) 
}