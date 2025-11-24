package router

import (
	"context"
	"errors"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gofiber/fiber/v2"
	fiberMiddleware "github.com/oapi-codegen/fiber-middleware"
	"github.com/rifkiadrn/cassandra-explore/internal/handler/rest"
	"github.com/sirupsen/logrus"
)

type RouterConfig struct {
	App            *fiber.App
	APIHandler     rest.APIHandler
	AuthMiddleware fiber.Handler
	Log            *logrus.Logger
}

func (r *RouterConfig) Setup() {
	r.App.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})

	// Internal routes (not in OpenAPI spec)
	internal := r.App.Group("/internal")

	internal.Post("/users", r.APIHandler.RegisterUser)

	// API exposes: /internal/healthz
	internal.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// API exposes: /interal/metrics
	internal.Get("/metrics", func(c *fiber.Ctx) error {
		// Expose Prometheus or other metrics
		return c.SendString("prometheus metrics here")
	})

	swagger, err := rest.GetSwagger()
	if err != nil {
		r.Log.Fatalf("failed to get swagger: %v", err)
	}
	swagger.Servers = nil // ignore server URL check
	// Add base path to validator
	swagger.AddServer(&openapi3.Server{URL: "/api/v1"})

	api := r.App.Group("/api/v1")

	api.Use(func(c *fiber.Ctx) error {
		// Run OAPI validator manually with injected AuthenticationFunc
		validator := fiberMiddleware.OapiRequestValidatorWithOptions(swagger, &fiberMiddleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: func(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
					// Save auth type directly in Fiber context
					c.Locals("authType", ai.SecuritySchemeName)

					switch ai.SecuritySchemeName {
					case "ApiKeyAuth":
						apiKey := ai.RequestValidationInput.Request.Header.Get("X-API-KEY")
						if apiKey == "" {
							return errors.New("missing api key")
						}
						// validate API key here
						return nil
					case "BearerAuth":
						authHeader := ai.RequestValidationInput.Request.Header.Get("Authorization")
						if authHeader == "" {
							return errors.New("missing bearer token")
						}
						// validate JWT here
						return nil
					default:
						return nil
					}
				},
			},
		})

		// Continue to the validator middleware
		return validator(c)
	})

	authSkipper := func(c *fiber.Ctx) error {
		fmt.Println("authSkipper", c.Path())
		authType, _ := c.Locals("authType").(string)
		fmt.Println("authType", authType)

		if authType == "" {
			return c.Next()
		}

		return r.AuthMiddleware(c)
	}

	// API exposes: openapi /api/v1
	rest.RegisterHandlersWithOptions(api, &r.APIHandler, rest.FiberServerOptions{
		Middlewares: []rest.MiddlewareFunc{
			authSkipper,
		},
	})
}
