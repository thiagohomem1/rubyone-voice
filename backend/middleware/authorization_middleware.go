// backend/middleware/authorization_middleware.go
package middleware

import (
	"strings"
	"github.com/gofiber/fiber/v2"
	"rubyone-voice/services"
	"rubyone-voice/database"
)

func RequirePermission(permissionCode string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extrair informações do usuário do contexto (definido pelo AuthMiddleware)
		tenantID, ok := c.Locals("tenant_id").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized: tenant context required",
			})
		}

		roleID, ok := c.Locals("role_id").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized: role context required",
			})
		}

		// Inicializar o RoleService
		roleService := services.NewRoleService(database.GetDB())

		// Buscar o role com suas permissões
		role, err := roleService.GetRoleByID(tenantID, roleID)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "forbidden: role not found or access denied",
			})
		}

		// Verificar se o role tem a permissão necessária (case-insensitive)
		hasPermission := false
		for _, permission := range role.Permissions {
			if strings.EqualFold(permission.Code, permissionCode) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "forbidden: missing required permission",
			})
		}

		return c.Next()
	}
} 