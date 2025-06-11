// backend/middleware/auth_middleware.go
package middleware

import (
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/utils"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "authorization header required",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		token := tokenParts[1]
		claims, err := utils.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("tenant_id", claims.TenantID)
		c.Locals("role_id", claims.RoleID)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}

func TenantMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Locals("tenant_id")
		if tenantID == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "tenant context required",
			})
		}
		return c.Next()
	}
} 