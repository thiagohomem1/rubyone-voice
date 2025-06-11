package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/controllers"
	"github.com/your-module/backend/middleware"
)

func SetupCallRoutes(app *fiber.App, controller *controllers.CallController) {
	api := app.Group("/api/v1")

	calls := api.Group("/calls", 
		middleware.AuthMiddleware(),
		middleware.TenantMiddleware(),
	)

	calls.Post("/", 
		middleware.RequirePermission("call.create"),
		controller.CreateCall,
	)
	calls.Get("/", 
		middleware.RequirePermission("call.read"),
		controller.GetAllCalls,
	)
	calls.Get("/:id", 
		middleware.RequirePermission("call.read"),
		controller.GetCallByID,
	)
	calls.Delete("/:id", 
		middleware.RequirePermission("call.delete"),
		controller.DeleteCall,
	)
} 