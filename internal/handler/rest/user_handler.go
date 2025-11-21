package rest

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/rifkiadrn/cassandra-explore/internal/model"
	"github.com/sirupsen/logrus"
)

type IUserUseCase interface {
	Register(ctx context.Context, request model.RegisterUser) (model.User, error)
	Login(ctx context.Context, request model.LoginUser) (model.LoginResponse, error)
}

type UserHandler struct {
	Log     *logrus.Logger
	UseCase IUserUseCase
}

func NewUserHandler(useCase IUserUseCase, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		Log:     logger,
		UseCase: useCase,
	}
}

func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	var req model.RegisterUser
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	user, err := h.UseCase.Register(c.Context(), req)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
	var req model.LoginUser
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	response, err := h.UseCase.Login(c.Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(response)
}
