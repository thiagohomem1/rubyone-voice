package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/controllers"
	"github.com/your-module/backend/middleware"
)

func SetupSubscriptionRoutes(app *fiber.App, controller *controllers.SubscriptionController) {
	api := app.Group("/api/v1")

	// Admin plan management routes
	adminPlans := api.Group("/admin/plans",
		middleware.AuthMiddleware(),
	)

	adminPlans.Post("/",
		middleware.RequirePermission("admin.plan.create"),
		controller.CreatePlan)

	adminPlans.Get("/",
		middleware.RequirePermission("admin.plan.read"),
		controller.GetAllPlans)

	adminPlans.Get("/:id",
		middleware.RequirePermission("admin.plan.read"),
		controller.GetPlanByID)

	adminPlans.Delete("/:id",
		middleware.RequirePermission("admin.plan.delete"),
		controller.DeletePlan)

	// Admin tenant subscription management routes
	adminTenants := api.Group("/admin/tenants",
		middleware.AuthMiddleware(),
	)

	adminTenants.Post("/:tenant_id/subscribe",
		middleware.RequirePermission("admin.tenant.subscribe"),
		controller.SubscribeTenant)

	// Tenant subscription self-service routes
	subscription := api.Group("/subscription",
		middleware.AuthMiddleware(),
		middleware.TenantMiddleware(),
	)

	subscription.Get("/",
		middleware.RequirePermission("subscription.read"),
		controller.GetTenantSubscription)
} 