// backend/routes/auth_routes.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/controllers"
	"github.com/your-module/backend/middleware"
)

func SetupAuthRoutes(app *fiber.App, authController *controllers.AuthController) {
	api := app.Group("/api/v1")
	
	auth := api.Group("/auth")
	auth.Post("/register-tenant", authController.RegisterTenant)
	auth.Post("/register-user", authController.RegisterUser)
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.Logout)
	
	protected := auth.Group("/", middleware.AuthMiddleware())
	protected.Get("/profile", authController.Profile)
} 