package routes

import (
	"saas-backend/controllers"
	"saas-backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupUsageRoutes(app *fiber.App, usageController *controllers.UsageController) {
	api := app.Group("/api/v1/usage")
	
	api.Use(middleware.AuthMiddleware())
	api.Use(middleware.TenantMiddleware())
	api.Use(middleware.RequirePermission("usage.read"))

	api.Get("/", usageController.GetUsage)
} 