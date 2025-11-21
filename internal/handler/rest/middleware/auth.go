package middleware

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	model_api "github.com/rifkiadrn/cassandra-explore/internal/model/api"
	"github.com/sirupsen/logrus"
)

type IUserUseCase interface {
	Verify(ctx context.Context, request model_api.VerifyUserRequest) (model_api.Auth, error)
}

func NewAuth(userUseCase IUserUseCase, logger *logrus.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		apiKey := ctx.Get("X-API-KEY")
		bearer := ctx.Get("Authorization")

		// bearer := c.Get(rest.BearerAuthScopes)
		fmt.Println("bearer", bearer)
		var token string
		if bearer != "" {
			// Try to get Bearer token from Authorization header
			authHeader := ctx.Get("Authorization", "")
			if authHeader == "" {
				logger.Warn("Missing Authorization header")
				return fiber.ErrUnauthorized
			}

			// Typically "Bearer <token>"
			fmt.Sscanf(authHeader, "Bearer %s", &token)
			if token == "" {
				logger.Warn("Empty bearer token")
				return fiber.ErrUnauthorized
			}
		}

		// apiKey := c.Get(rest.ApiKeyAuthScopes)
		fmt.Println("apiKey", apiKey)
		if apiKey != "" {
			// Try to get API key from X-API-KEY header
			token = ctx.Get("X-API-KEY", "")
			if token == "" {
				logger.Warn("Missing X-API-KEY header")
				return fiber.ErrUnauthorized
			}
		}

		// Verify token (could be JWT or API key, depending on your logic)
		fmt.Println("token", token)
		auth, err := userUseCase.Verify(ctx.Context(), model_api.VerifyUserRequest{Token: token})
		if err != nil {
			logger.Warnf("Invalid token: %+v", err)
			return fiber.ErrUnauthorized
		}

		fmt.Println("auth", auth)

		ctx.Locals("auth", auth)
		ctx.Locals("user_id", auth.ID)
		// ALSO set in User Context (for business logic)
		newCtx := context.WithValue(ctx.UserContext(), "auth", auth)
		ctx.SetUserContext(newCtx)
		return ctx.Next()
	}
}
