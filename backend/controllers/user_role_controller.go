package controllers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/services"
)

type UserRoleController struct {
	UserRoleService *services.UserRoleService
}

func NewUserRoleController(service *services.UserRoleService) *UserRoleController {
	return &UserRoleController{UserRoleService: service}
}

func (urc *UserRoleController) AssignRolesToUser(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)
	
	userIDStr := c.Params("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}

	var req struct {
		RoleIDs []uint `json:"role_ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if len(req.RoleIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "at least one role ID is required",
		})
	}

	err = urc.UserRoleService.AssignRolesToUser(tenantID, uint(userID), req.RoleIDs)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		if err.Error() == "one or more roles not found or don't belong to tenant" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "one or more roles not found or don't belong to tenant",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "roles assigned to user successfully",
	})
}

func (urc *UserRoleController) GetUserRoles(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)
	
	userIDStr := c.Params("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}

	roles, err := urc.UserRoleService.GetUserRoles(tenantID, uint(userID))
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "user roles retrieved successfully",
		"data":    roles,
	})
}

func (urc *UserRoleController) RemoveRoleFromUser(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)
	
	userIDStr := c.Params("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}

	roleIDStr := c.Params("role_id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid role ID",
		})
	}

	err = urc.UserRoleService.RemoveRoleFromUser(tenantID, uint(userID), uint(roleID))
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		if err.Error() == "role not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "role not found",
			})
		}
		if err.Error() == "user role assignment not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user role assignment not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "role removed from user successfully",
	})
} 