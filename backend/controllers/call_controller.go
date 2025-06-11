package controllers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/backend/services"
)

type CallController struct {
	CallService *services.CallService
}

func NewCallController(service *services.CallService) *CallController {
	return &CallController{CallService: service}
}

func (cc *CallController) CreateCall(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)
	
	var req struct {
		Caller       string `json:"caller"`
		Callee       string `json:"callee"`
		Duration     uint   `json:"duration"`
		RecordingURL string `json:"recording_url"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Caller == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "caller is required",
		})
	}

	if req.Callee == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "callee is required",
		})
	}

	call, err := cc.CallService.CreateCall(tenantID, req.Caller, req.Callee, req.Duration, req.RecordingURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "call created successfully",
		"data":    call,
	})
}

func (cc *CallController) GetAllCalls(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)

	calls, err := cc.CallService.GetAllCalls(tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "calls retrieved successfully",
		"data":    calls,
	})
}

func (cc *CallController) GetCallByID(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)
	
	callIDStr := c.Params("id")
	callID, err := strconv.ParseUint(callIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid call ID",
		})
	}

	call, err := cc.CallService.GetCallByID(tenantID, uint(callID))
	if err != nil {
		if err.Error() == "call not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "call not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "call retrieved successfully",
		"data":    call,
	})
}

func (cc *CallController) DeleteCall(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)
	
	callIDStr := c.Params("id")
	callID, err := strconv.ParseUint(callIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid call ID",
		})
	}

	err = cc.CallService.DeleteCall(tenantID, uint(callID))
	if err != nil {
		if err.Error() == "call not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "call not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "call deleted successfully",
	})
} 