// backend/controllers/auth_controller.go
package controllers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

type RegisterTenantRequest struct {
	TenantName string `json:"tenant_name" validate:"required"`
	Domain     string `json:"domain" validate:"required"`
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required,min=6"`
}

type RegisterUserRequest struct {
	TenantID uint   `json:"tenant_id" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	RoleID   uint   `json:"role_id" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (ac *AuthController) RegisterTenant(c *fiber.Ctx) error {
	var req RegisterTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user, token, err := ac.authService.RegisterTenant(req.TenantName, req.Domain, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "tenant and admin user created successfully",
		"user":    user,
		"token":   token,
	})
}

func (ac *AuthController) RegisterUser(c *fiber.Ctx) error {
	var req RegisterUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user, token, err := ac.authService.RegisterUser(req.TenantID, req.Username, req.Password, req.RoleID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user created successfully",
		"user":    user,
		"token":   token,
	})
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user, token, err := ac.authService.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "login successful",
		"user":    user,
		"token":   token,
	})
}

func (ac *AuthController) Logout(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logout successful",
	})
}

func (ac *AuthController) Profile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	tenantID := c.Locals("tenant_id").(uint)
	roleID := c.Locals("role_id").(uint)
	username := c.Locals("username").(string)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":   userID,
		"tenant_id": tenantID,
		"role_id":   roleID,
		"username":  username,
	})
} 