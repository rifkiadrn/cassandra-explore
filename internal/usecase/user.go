package usecase

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rifkiadrn/cassandra-explore/internal/entity"
	"github.com/rifkiadrn/cassandra-explore/internal/model"
	model_api "github.com/rifkiadrn/cassandra-explore/internal/model/api"
	"github.com/rifkiadrn/cassandra-explore/internal/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type IUserRepo interface {
	Create(ctx context.Context, user entity.User) (*entity.User, error)
	FindById(ctx context.Context, userID string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, existingUser entity.User, updatedUser entity.User) (*entity.User, error)
}

type IUserRepoNoSQL interface {
	Create(ctx context.Context, user entity.User) (*entity.User, error)
}

type UserUseCase struct {
	uow                 UnitOfWork
	log                 *logrus.Logger
	validate            *validator.Validate
	userRepository      IUserRepo
	userRepositoryNoSQL IUserRepoNoSQL
	jwtManager          *utils.JWTManager
}

func NewUserUseCase(uow UnitOfWork, logger *logrus.Logger, validate *validator.Validate,
	userRepository IUserRepo, userRepositoryNoSQL IUserRepoNoSQL, jwtManager *utils.JWTManager) UserUseCase {
	return UserUseCase{
		uow:                 uow,
		log:                 logger,
		validate:            validate,
		userRepository:      userRepository,
		userRepositoryNoSQL: userRepositoryNoSQL,
		jwtManager:          jwtManager,
	}
}

func (userUC UserUseCase) Verify(ctx context.Context, request model_api.VerifyUserRequest) (model_api.Auth, error) {
	// Validate request
	if err := userUC.validate.Struct(request); err != nil {
		userUC.log.Warnf("Invalid request body : %+v", err)
		return model_api.Auth{}, fiber.ErrBadRequest
	}

	// Verify JWT token
	claims, err := userUC.jwtManager.ValidateToken(request.Token)
	if err != nil {
		userUC.log.Warnf("Invalid JWT token: %+v", err)
		return model_api.Auth{}, fiber.ErrUnauthorized
	}

	// Find user by ID
	userEntity, err := userUC.userRepository.FindById(ctx, claims.UserID.String())
	if err != nil {
		userUC.log.Warnf("User not found for token subject: %+v", err)
		return model_api.Auth{}, fiber.ErrUnauthorized
	}

	return model_api.Auth{
		ID:       (*userEntity).ID,
		Username: (*userEntity).Username,
		Token:    request.Token,
	}, nil
}

func (userUC UserUseCase) Register(ctx context.Context, request model.RegisterUser) (model.User, error) {
	// Validate request
	if err := userUC.validate.Struct(request); err != nil {
		userUC.log.Warnf("Invalid request body : %+v", err)
		return model.User{}, fiber.ErrBadRequest
	}

	// Start transaction
	tx, txCtx, err := userUC.uow.Begin(ctx)
	if err != nil {
		return model.User{}, err
	}
	defer tx.Rollback()

	// Check if username already exists
	_, err = userUC.userRepository.FindByUsername(txCtx, request.Username)
	if err == nil {
		return model.User{}, fiber.ErrConflict
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		userUC.log.Warnf("Failed to hash password : %+v", err)
		return model.User{}, fiber.ErrInternalServerError
	}

	// Create domain entity
	userEntity := entity.User{
		ID:        uuid.New(),
		Name:      request.Name,
		Username:  request.Username,
		Password:  string(hashed),
		Token:     uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create user via repository
	createdUser, err := userUC.userRepository.Create(txCtx, userEntity)
	if err != nil {
		userUC.log.Warnf("Failed create user : %+v", err)
		return model.User{}, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		userUC.log.Warnf("Failed commit transaction : %+v", err)
		return model.User{}, fiber.ErrInternalServerError
	}

	// double create to cassandra
	_, _ = userUC.userRepositoryNoSQL.Create(ctx, *createdUser)

	// Create user in Cassandra

	// Convert domain entity to response DTO
	return model.User{
		Id:       (*createdUser).ID.String(),
		Name:     (*createdUser).Name,
		Username: (*createdUser).Username,
		Token:    (*createdUser).Token,
	}, nil
}

func (userUC UserUseCase) Login(ctx context.Context, request model.LoginUser) (model.LoginResponse, error) {
	// Validate request
	if err := userUC.validate.Struct(request); err != nil {
		userUC.log.Warnf("Invalid request body : %+v", err)
		return model.LoginResponse{}, fiber.ErrBadRequest
	}

	// Find user by username
	userEntity, err := userUC.userRepository.FindByUsername(ctx, request.Username)
	if err != nil {
		userUC.log.Warnf("Failed find user by username : %+v", err)
		return model.LoginResponse{}, fiber.ErrUnauthorized
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte((*userEntity).Password), []byte(request.Password)); err != nil {
		userUC.log.Warnf("Invalid password : %+v", err)
		return model.LoginResponse{}, fiber.ErrUnauthorized
	}

	// Generate JWT token
	token, err := userUC.jwtManager.GenerateToken((*userEntity).ID, (*userEntity).Username)
	if err != nil {
		userUC.log.Warnf("Failed to generate token : %+v", err)
		return model.LoginResponse{}, fiber.ErrInternalServerError
	}

	// Convert domain entity to response DTO
	return model.LoginResponse{
		Token: token,
		User: model.User{
			Id:       (*userEntity).ID.String(),
			Name:     (*userEntity).Name,
			Username: (*userEntity).Username,
			Token:    (*userEntity).Token,
		},
	}, nil
}

func (userUC UserUseCase) FindById(ctx context.Context, userID string) (entity.User, error) {
	user, err := userUC.userRepository.FindById(ctx, userID)
	if err != nil {
		return entity.User{}, err
	}
	return *user, nil
}

func (userUC UserUseCase) UpdateUser(ctx context.Context, userID string, user entity.User) (entity.User, error) {
	// Start transaction
	tx, txCtx, err := userUC.uow.Begin(ctx)
	if err != nil {
		return entity.User{}, err
	}
	defer tx.Rollback()

	existingUser, err := userUC.userRepository.FindById(txCtx, userID)
	if err != nil {
		return entity.User{}, err
	}
	updatedUser, err := userUC.userRepository.Update(txCtx, *existingUser, user)
	if err != nil {
		return entity.User{}, err
	}

	if err := tx.Commit(); err != nil {
		return entity.User{}, err
	}

	return *updatedUser, nil
}
