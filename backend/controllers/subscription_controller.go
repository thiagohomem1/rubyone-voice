package controllers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/services"
)

type SubscriptionController struct {
	SubscriptionService *services.SubscriptionService
}

func NewSubscriptionController(service *services.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{SubscriptionService: service}
}

func (sc *SubscriptionController) CreatePlan(c *fiber.Ctx) error {
	var req struct {
		Name     string  `json:"name"`
		MaxUsers uint    `json:"max_users"`
		MaxCalls uint    `json:"max_calls"`
		Price    float64 `json:"price"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "plan name is required",
		})
	}

	if req.MaxUsers == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "max users must be greater than 0",
		})
	}

	if req.MaxCalls == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "max calls must be greater than 0",
		})
	}

	if req.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "price must be greater than or equal to 0",
		})
	}

	plan, err := sc.SubscriptionService.CreatePlan(req.Name, req.MaxUsers, req.MaxCalls, req.Price)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "plan created successfully",
		"data":    plan,
	})
}

func (sc *SubscriptionController) GetAllPlans(c *fiber.Ctx) error {
	plans, err := sc.SubscriptionService.GetAllPlans()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "plans retrieved successfully",
		"data":    plans,
	})
}

func (sc *SubscriptionController) GetPlanByID(c *fiber.Ctx) error {
	planIDStr := c.Params("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid plan ID",
		})
	}

	plan, err := sc.SubscriptionService.GetPlanByID(uint(planID))
	if err != nil {
		if err.Error() == "plan not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "plan not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "plan retrieved successfully",
		"data":    plan,
	})
}

func (sc *SubscriptionController) DeletePlan(c *fiber.Ctx) error {
	planIDStr := c.Params("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid plan ID",
		})
	}

	err = sc.SubscriptionService.DeletePlan(uint(planID))
	if err != nil {
		if err.Error() == "plan not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "plan not found",
			})
		}
		if err.Error() == "cannot delete plan with active subscriptions" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "cannot delete plan with active subscriptions",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "plan deleted successfully",
	})
}

func (sc *SubscriptionController) SubscribeTenant(c *fiber.Ctx) error {
	tenantIDStr := c.Params("tenant_id")
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid tenant ID",
		})
	}

	var req struct {
		PlanID uint `json:"plan_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.PlanID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "plan ID is required",
		})
	}

	err = sc.SubscriptionService.SubscribeTenant(uint(tenantID), req.PlanID)
	if err != nil {
		if err.Error() == "tenant not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "tenant not found",
			})
		}
		if err.Error() == "plan not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "plan not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "tenant subscribed successfully",
	})
}

func (sc *SubscriptionController) GetTenantSubscription(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)

	subscription, err := sc.SubscriptionService.GetTenantSubscription(tenantID)
	if err != nil {
		if err.Error() == "subscription not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "subscription not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "subscription retrieved successfully",
		"data":    subscription,
	})
} 