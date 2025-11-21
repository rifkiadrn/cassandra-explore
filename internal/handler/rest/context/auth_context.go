package context

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	model_api "github.com/rifkiadrn/cassandra-explore/internal/model/api"
)

func GetUserFromContext(ctx context.Context) (model_api.Auth, error) {
	authInterface := ctx.Value("auth")
	if authInterface == nil {
		return model_api.Auth{}, fiber.ErrUnauthorized
	}

	auth, ok := authInterface.(model_api.Auth)
	if !ok {
		return model_api.Auth{}, fiber.ErrUnauthorized
	}

	fmt.Println("GetUserFromContext: ", auth)

	return auth, nil
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	auth, err := GetUserFromContext(ctx)
	if err != nil {
		return uuid.UUID{}, err
	}
	return auth.ID, nil
}
