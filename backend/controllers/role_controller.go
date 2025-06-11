package controllers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/services"
)

type RoleController struct {
	RoleService *services.RoleService
}

func NewRoleController(service *services.RoleService) *RoleController {
	return &RoleController{RoleService: service}
}

func (rc *RoleController) CreateRole(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)
	
	var req struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "role name is required",
		})
	}

	role, err := rc.RoleService.CreateRole(tenantID, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "role created successfully",
		"data":    role,
	})
}

func (rc *RoleController) GetRoles(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)

	roles, err := rc.RoleService.GetRolesByTenant(tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "roles retrieved successfully",
		"data":    roles,
	})
}

func (rc *RoleController) GetRole(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)
	
	roleIDStr := c.Params("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid role ID",
		})
	}

	role, err := rc.RoleService.GetRoleByID(tenantID, uint(roleID))
	if err != nil {
		if err.Error() == "role not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "role not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "role retrieved successfully",
		"data":    role,
	})
}

func (rc *RoleController) DeleteRole(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)
	
	roleIDStr := c.Params("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid role ID",
		})
	}

	err = rc.RoleService.DeleteRole(tenantID, uint(roleID))
	if err != nil {
		if err.Error() == "role not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "role not found",
			})
		}
		if err.Error() == "cannot delete role that is assigned to users" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "cannot delete role that is assigned to users",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "role deleted successfully",
	})
}

func (rc *RoleController) AssignPermissions(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)
	
	roleIDStr := c.Params("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid role ID",
		})
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	err = rc.RoleService.AssignPermissionsToRole(tenantID, uint(roleID), req.PermissionIDs)
	if err != nil {
		if err.Error() == "role not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "role not found",
			})
		}
		if err.Error() == "one or more permissions not found" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "one or more permissions not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "permissions assigned to role successfully",
	})
} 