package router

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	"github.com/rifkiadrn/cassandra-explore/internal/handler/rest"
	"github.com/sirupsen/logrus"
)

type RouterConfig struct {
	App            *fiber.App
	APIHandler     rest.APIHandler
	AuthMiddleware fiber.Handler
	Log            *logrus.Logger
}

type User struct {
	Id       gocql.UUID `json:"id"`
	Name     string     `json:"name"`
	Username string     `json:"username"`
}

func (r *RouterConfig) Setup() {
	r.App.Get("/ping", func(c *fiber.Ctx) error {

		fmt.Println("validator middleware hit:", c.Path())
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
}
