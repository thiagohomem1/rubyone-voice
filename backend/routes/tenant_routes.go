package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/controllers"
	"github.com/your-module/backend/middleware"
)

func SetupTenantRoutes(app *fiber.App, controller *controllers.TenantController) {
	api := app.Group("/api/v1")

	tenants := api.Group("/admin/tenants", 
		middleware.AuthMiddleware(),
	)

	tenants.Post("/", 
		middleware.RequirePermission("admin.tenant.create"),
		controller.CreateTenant)
	
	tenants.Get("/", 
		middleware.RequirePermission("admin.tenant.read"),
		controller.GetAllTenants)
	
	tenants.Get("/:id", 
		middleware.RequirePermission("admin.tenant.read"),
		controller.GetTenantByID)
	
	tenants.Delete("/:id", 
		middleware.RequirePermission("admin.tenant.delete"),
		controller.DeleteTenant)
} 