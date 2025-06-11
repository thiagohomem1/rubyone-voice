package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/models"
	"github.com/your-module/backend/services"
	"gorm.io/gorm"
)

func CheckUserQuota(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Locals("tenant_id")
		if tenantID == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "tenant context required",
			})
		}

		tenantIDUint, ok := tenantID.(uint)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid tenant ID",
			})
		}

		// Get active subscription for tenant
		subscriptionService := services.NewSubscriptionService(db)
		subscription, err := subscriptionService.GetTenantSubscription(tenantIDUint)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "subscription not found" {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "no active subscription found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to check subscription",
			})
		}

		// Count current active users for the tenant
		var userCount int64
		if err := db.Model(&models.User{}).
			Where("tenant_id = ? AND deleted_at IS NULL", tenantIDUint).
			Count(&userCount).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to count users",
			})
		}

		// Check if adding a new user would exceed the quota
		if uint(userCount) >= subscription.Plan.MaxUsers {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "quota exceeded: max users reached",
			})
		}

		return c.Next()
	}
}

func CheckCallQuota(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Locals("tenant_id")
		if tenantID == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "tenant context required",
			})
		}

		tenantIDUint, ok := tenantID.(uint)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid tenant ID",
			})
		}

		// Get active subscription for tenant
		subscriptionService := services.NewSubscriptionService(db)
		subscription, err := subscriptionService.GetTenantSubscription(tenantIDUint)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "subscription not found" {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "no active subscription found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to check subscription",
			})
		}

		// Count total calls for the tenant
		var callCount int64
		if err := db.Model(&models.Call{}).
			Where("tenant_id = ? AND deleted_at IS NULL", tenantIDUint).
			Count(&callCount).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to count calls",
			})
		}

		// Check if adding a new call would exceed the quota
		if uint(callCount) >= subscription.Plan.MaxCalls {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "quota exceeded: max calls reached",
			})
		}

		return c.Next()
	}
} 