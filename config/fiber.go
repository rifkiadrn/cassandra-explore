package config

import (
	"errors"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      config.GetString("app.name"),
		ErrorHandler: NewErrorHandler(),
	})

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		var e *fiber.Error

		if se, ok := err.(*openapi3filter.SecurityRequirementsError); ok {
			if len(se.Errors) > 0 {
				return c.Status(fiber.StatusUnauthorized).SendString(se.Errors[0].Error())
			}
		}

		if errors.As(err, &e) {
			code = e.Code
		}

		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(code).SendString(err.Error())
	}
}
