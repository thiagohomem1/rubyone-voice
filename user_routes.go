package routes

import (
	"saas-backend/controllers"
	"saas-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, userController *controllers.UserController, authMiddleware, tenantMiddleware middleware.Handler) {
	userGroup := app.Group("/api/v1/users")

	userGroup.Use(authMiddleware)
	userGroup.Use(tenantMiddleware)

	userGroup.Post("/", middleware.RequirePermission("user.create"), userController.CreateUser)
	userGroup.Get("/", middleware.RequirePermission("user.read"), userController.GetAllUsers)
	userGroup.Get("/:id", middleware.RequirePermission("user.read"), userController.GetUserByID)
	userGroup.Delete("/:id", middleware.RequirePermission("user.delete"), userController.DeleteUser)
} 