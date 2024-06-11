package handler

import (
	"github.com/gofiber/fiber/v2"
)

type SuccessResponse struct {
	Status  int         `json:"status"`
	Success interface{} `json:"success"`
}

type FailureResponse struct {
	Status int         `json:"status"`
	Fail   FailDetails `json:"fail"`
}

type FailDetails struct {
	Message    string `json:"msg"`
	Reason     string `json:"reason"`
	CustomCode int    `json:"code,omitempty"`
}

func handleError(c *fiber.Ctx, statusCode int, message, reason string) error {
	return c.Status(statusCode).JSON(FailureResponse{
		Status: statusCode,
		Fail: FailDetails{
			Message: message,
			Reason:  reason,
		},
	})
}

func handleSuccess(c *fiber.Ctx, statusCode int, data interface{}) error {
	a := SuccessResponse{
		Status:  statusCode,
		Success: data,
	}
	return c.Status(statusCode).JSON(a)
}
