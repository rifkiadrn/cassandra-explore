package usecase

import (
	"context"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	context_db "github.com/rifkiadrn/cassandra-explore/internal/context/db"
	"github.com/rifkiadrn/cassandra-explore/internal/entity"
	"github.com/rifkiadrn/cassandra-explore/internal/model"
	model_api "github.com/rifkiadrn/cassandra-explore/internal/model/api"
	"github.com/rifkiadrn/cassandra-explore/internal/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserRepo interface {
	Create(ctx context.Context, user entity.User) (*entity.User, error)
	FindById(ctx context.Context, userID string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, existingUser entity.User, updatedUser entity.User) (*entity.User, error)
}

type UserUseCase struct {
	DB             *gorm.DB
	NoSQLDB        *gocql.Session
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository IUserRepo
	JWTManager     *utils.JWTManager
}

func NewUserUseCase(db *gorm.DB, noSQLDB *gocql.Session, logger *logrus.Logger, validate *validator.Validate,
	userRepository IUserRepo, jwtManager *utils.JWTManager) UserUseCase {
	return UserUseCase{
		DB:             db,
		NoSQLDB:        noSQLDB,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
		JWTManager:     jwtManager,
	}
}

func (userUC UserUseCase) Verify(ctx context.Context, request model_api.VerifyUserRequest) (model_api.Auth, error) {
	// Validate request
	if err := userUC.Validate.Struct(request); err != nil {
		userUC.Log.Warnf("Invalid request body : %+v", err)
		return model_api.Auth{}, fiber.ErrBadRequest
	}

	// Verify JWT token
	claims, err := userUC.JWTManager.ValidateToken(request.Token)
	if err != nil {
		userUC.Log.Warnf("Invalid JWT token: %+v", err)
		return model_api.Auth{}, fiber.ErrUnauthorized
	}

	// Start transaction
	tx, txCtx := context_db.BeginTxWithContext(ctx, userUC.DB)
	defer tx.Rollback()

	// Find user by ID
	userEntity, err := userUC.UserRepository.FindById(txCtx, claims.UserID.String())
	if err != nil {
		userUC.Log.Warnf("User not found for token subject: %+v", err)
		return model_api.Auth{}, fiber.ErrUnauthorized
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		userUC.Log.Warnf("Failed commit transaction : %+v", err)
		return model_api.Auth{}, fiber.ErrInternalServerError
	}

	return model_api.Auth{
		ID:       (*userEntity).ID,
		Username: (*userEntity).Username,
		Token:    request.Token,
	}, nil
}

func (userUC UserUseCase) Register(ctx context.Context, request model.RegisterUser) (model.User, error) {
	// Validate request
	if err := userUC.Validate.Struct(request); err != nil {
		userUC.Log.Warnf("Invalid request body : %+v", err)
		return model.User{}, fiber.ErrBadRequest
	}

	// Start transaction
	tx, txCtx := context_db.BeginTxWithContext(ctx, userUC.DB)
	defer tx.Rollback()

	// Check if username already exists
	_, err := userUC.UserRepository.FindByUsername(txCtx, request.Username)
	if err == nil {
		return model.User{}, fiber.ErrConflict
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		userUC.Log.Warnf("Failed to hash password : %+v", err)
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
	createdUser, err := userUC.UserRepository.Create(txCtx, userEntity)
	if err != nil {
		userUC.Log.Warnf("Failed create user : %+v", err)
		return model.User{}, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		userUC.Log.Warnf("Failed commit transaction : %+v", err)
		return model.User{}, fiber.ErrInternalServerError
	}

	id := gocql.TimeUUID()

	// Create user in Cassandra

	if err := userUC.NoSQLDB.Query(`INSERT INTO users (id, name, username) VALUES (?, ?, ?)`, id, request.Name, request.Username).Exec(); err != nil {
		return model.User{}, err
	}

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
	if err := userUC.Validate.Struct(request); err != nil {
		userUC.Log.Warnf("Invalid request body : %+v", err)
		return model.LoginResponse{}, fiber.ErrBadRequest
	}

	// Start transaction
	tx, txCtx := context_db.BeginTxWithContext(ctx, userUC.DB)
	defer tx.Rollback()

	// Find user by username
	userEntity, err := userUC.UserRepository.FindByUsername(txCtx, request.Username)
	if err != nil {
		userUC.Log.Warnf("Failed find user by username : %+v", err)
		return model.LoginResponse{}, fiber.ErrUnauthorized
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte((*userEntity).Password), []byte(request.Password)); err != nil {
		userUC.Log.Warnf("Invalid password : %+v", err)
		return model.LoginResponse{}, fiber.ErrUnauthorized
	}

	// Generate JWT token
	token, err := userUC.JWTManager.GenerateToken((*userEntity).ID, (*userEntity).Username)
	if err != nil {
		userUC.Log.Warnf("Failed to generate token : %+v", err)
		return model.LoginResponse{}, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		userUC.Log.Warnf("Failed commit transaction : %+v", err)
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
	tx, txCtx := context_db.BeginTxWithContext(ctx, userUC.DB)
	defer tx.Rollback()

	user, err := userUC.UserRepository.FindById(txCtx, userID)
	if err != nil {
		return entity.User{}, err
	}
	return *user, nil
}

func (userUC UserUseCase) UpdateUser(ctx context.Context, userID string, user entity.User) (entity.User, error) {
	tx, txCtx := context_db.BeginTxWithContext(ctx, userUC.DB)
	defer tx.Rollback()

	existingUser, err := userUC.UserRepository.FindById(txCtx, userID)
	if err != nil {
		return entity.User{}, err
	}
	updatedUser, err := userUC.UserRepository.Update(txCtx, *existingUser, user)
	if err != nil {
		return entity.User{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return entity.User{}, err
	}

	return *updatedUser, nil
}
