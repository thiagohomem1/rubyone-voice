package controllers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/services"
)

type TenantController struct {
	TenantService *services.TenantService
}

func NewTenantController(service *services.TenantService) *TenantController {
	return &TenantController{TenantService: service}
}

func (tc *TenantController) CreateTenant(c *fiber.Ctx) error {
	var req struct {
		Name   string `json:"name"`
		Domain string `json:"domain"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Name == "" || req.Domain == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name and domain are required",
		})
	}

	tenant, err := tc.TenantService.CreateTenant(req.Name, req.Domain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "tenant created successfully",
		"data":    tenant,
	})
}

func (tc *TenantController) GetAllTenants(c *fiber.Ctx) error {
	tenants, err := tc.TenantService.GetAllTenants()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "tenants retrieved successfully",
		"data":    tenants,
	})
}

func (tc *TenantController) GetTenantByID(c *fiber.Ctx) error {
	tenantIDStr := c.Params("id")
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid tenant ID",
		})
	}

	tenant, err := tc.TenantService.GetTenantByID(uint(tenantID))
	if err != nil {
		if err.Error() == "tenant not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "tenant not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "tenant retrieved successfully",
		"data":    tenant,
	})
}

func (tc *TenantController) DeleteTenant(c *fiber.Ctx) error {
	tenantIDStr := c.Params("id")
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid tenant ID",
		})
	}

	err = tc.TenantService.DeleteTenant(uint(tenantID))
	if err != nil {
		if err.Error() == "tenant not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "tenant not found",
			})
		}
		if err.Error() == "cannot delete tenant that has associated users" ||
		   err.Error() == "cannot delete tenant that has associated roles" ||
		   err.Error() == "cannot delete tenant that has associated calls" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "tenant deleted successfully",
	})
} 