package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/controllers"
	"github.com/your-module/backend/middleware"
)

func SetupUserRoleRoutes(app *fiber.App, controller *controllers.UserRoleController) {
	api := app.Group("/api/v1")

	userRoles := api.Group("/user-roles", 
		middleware.AuthMiddleware(),
		middleware.TenantMiddleware(),
	)

	userRoles.Post("/:user_id/roles", 
		middleware.RequirePermission("userrole.assign"),
		controller.AssignRolesToUser)
	
	userRoles.Get("/:user_id/roles", 
		middleware.RequirePermission("userrole.read"),
		controller.GetUserRoles)
	
	userRoles.Delete("/:user_id/roles/:role_id", 
		middleware.RequirePermission("userrole.remove"),
		controller.RemoveRoleFromUser)
} 