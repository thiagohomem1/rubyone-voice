package controllers

import (
	"saas-backend/services"
	"github.com/gofiber/fiber/v2"
)

type UsageController struct {
	usageService *services.UsageService
}

func NewUsageController(usageService *services.UsageService) *UsageController {
	return &UsageController{
		usageService: usageService,
	}
}

func (c *UsageController) GetUsage(ctx *fiber.Ctx) error {
	tenantID, ok := ctx.Locals("tenantID").(uint)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant context",
		})
	}

	usage, err := c.usageService.GetTenantUsage(tenantID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve usage data",
		})
	}

	return ctx.JSON(fiber.Map{
		"usage": usage,
	})
} 