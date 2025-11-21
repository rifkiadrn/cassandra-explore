package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type GenericHandler struct {
	Log *logrus.Logger
}

func NewGenericHandler(logger *logrus.Logger) *GenericHandler {
	return &GenericHandler{
		Log: logger,
	}
}

// Implement the generated ServerInterface
func (h *GenericHandler) GetPing(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "pong",
	})
}
