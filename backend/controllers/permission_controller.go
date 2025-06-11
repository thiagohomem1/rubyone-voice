// backend/controllers/permission_controller.go
package controllers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/services"
)

type PermissionController struct {
	PermissionService *services.PermissionService
}

func NewPermissionController(service *services.PermissionService) *PermissionController {
	return &PermissionController{PermissionService: service}
}

func (pc *PermissionController) CreatePermission(c *fiber.Ctx) error {
	var req struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Code == "" || req.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "code and description are required",
		})
	}

	permission, err := pc.PermissionService.CreatePermission(req.Code, req.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "permission created successfully",
		"data":    permission,
	})
}

func (pc *PermissionController) GetAllPermissions(c *fiber.Ctx) error {
	permissions, err := pc.PermissionService.GetAllPermissions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "permissions retrieved successfully",
		"data":    permissions,
	})
}

func (pc *PermissionController) GetPermissionByID(c *fiber.Ctx) error {
	permissionIDStr := c.Params("id")
	permissionID, err := strconv.ParseUint(permissionIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid permission ID",
		})
	}

	permission, err := pc.PermissionService.GetPermissionByID(uint(permissionID))
	if err != nil {
		if err.Error() == "permission not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "permission not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "permission retrieved successfully",
		"data":    permission,
	})
}

func (pc *PermissionController) DeletePermission(c *fiber.Ctx) error {
	permissionIDStr := c.Params("id")
	permissionID, err := strconv.ParseUint(permissionIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid permission ID",
		})
	}

	err = pc.PermissionService.DeletePermission(uint(permissionID))
	if err != nil {
		if err.Error() == "permission not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "permission not found",
			})
		}
		if err.Error() == "cannot delete permission that is assigned to roles" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "cannot delete permission that is assigned to roles",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "permission deleted successfully",
	})
}

func (pc *PermissionController) AssignPermissionToRole(c *fiber.Ctx) error {
	var req struct {
		RoleID       uint `json:"role_id"`
		PermissionID uint `json:"permission_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	err := pc.PermissionService.AssignPermissionToRole(req.RoleID, req.PermissionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "permission assigned to role successfully",
	})
}
