// backend/middleware/permission_middleware.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/services"
)

func PermissionMiddleware(permissionCode string, permissionService *services.PermissionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uint)

		hasPermission, err := permissionService.HasPermission(userID, permissionCode)
		if err != nil || !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "forbidden - missing permission",
			})
		}

		return c.Next()
	}
}
