// backend/routes/permission_routes.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/controllers"
	"github.com/your-module/backend/middleware"
)

func SetupPermissionRoutes(app *fiber.App, controller *controllers.PermissionController) {
	api := app.Group("/api/v1")

	permissions := api.Group("/permissions",
		middleware.AuthMiddleware(),
		middleware.TenantMiddleware(),
	)

	permissions.Post("/", controller.CreatePermission)
	permissions.Get("/", controller.GetAllPermissions)
	permissions.Get("/:id", controller.GetPermissionByID)
	permissions.Delete("/:id", controller.DeletePermission)
	permissions.Post("/assign", controller.AssignPermissionToRole)
}
