package controllers

import (
	"saas-backend/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	RoleID    uint   `json:"role_id" validate:"required"`
}

func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	tenantID := ctx.Locals("tenant_id").(uint)

	var req CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" || req.RoleID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields are required",
		})
	}

	user, err := c.userService.CreateUser(tenantID, req.RoleID, req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    user,
	})
}

func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	tenantID := ctx.Locals("tenant_id").(uint)

	users, err := c.userService.GetAllUsers(tenantID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Users retrieved successfully",
		"data":    users,
	})
}

func (c *UserController) GetUserByID(ctx *fiber.Ctx) error {
	tenantID := ctx.Locals("tenant_id").(uint)

	userIDStr := ctx.Params("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := c.userService.GetUserByID(tenantID, uint(userID))
	if err != nil {
		if err.Error() == "user not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "User retrieved successfully",
		"data":    user,
	})
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	tenantID := ctx.Locals("tenant_id").(uint)

	userIDStr := ctx.Params("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	err = c.userService.DeleteUser(tenantID, uint(userID))
	if err != nil {
		if err.Error() == "user not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
} 