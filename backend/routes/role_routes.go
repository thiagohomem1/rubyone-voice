package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/controllers"
	"github.com/your-module/backend/middleware"
)

func SetupRoleRoutes(app *fiber.App, controller *controllers.RoleController) {
	api := app.Group("/api/v1")

	roles := api.Group("/roles", 
		middleware.AuthMiddleware(),
		middleware.TenantMiddleware(),
	)

	roles.Post("/", controller.CreateRole)
	roles.Get("/", controller.GetRoles)
	roles.Get("/:id", controller.GetRole)
	roles.Delete("/:id", controller.DeleteRole)
	roles.Post("/:id/permissions", controller.AssignPermissions)
} 