package config

import (
	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	context_db "github.com/rifkiadrn/cassandra-explore/internal/context/db"
	"github.com/rifkiadrn/cassandra-explore/internal/handler/rest"
	"github.com/rifkiadrn/cassandra-explore/internal/handler/rest/middleware"
	"github.com/rifkiadrn/cassandra-explore/internal/handler/rest/router"
	"github.com/rifkiadrn/cassandra-explore/internal/repository"
	"github.com/rifkiadrn/cassandra-explore/internal/usecase"
	"github.com/rifkiadrn/cassandra-explore/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	NoSQLDB  *gocql.Session
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {

	// setup repositories
	userRepository := repository.NewUserRepository(config.DB, config.Log)
	userRepositoryNoSQL := repository.NewUserRepositoryNoSQL(config.NoSQLDB)
	blogRepository := repository.NewBlogRepository(config.DB, config.Log)

	// setup JWT manager
	jwtManager := utils.NewJWTManager(config.Config.GetString("SECRET_KEY")) // TODO: move to config

	// setup use cases
	// dbTrx/unitOfWork
	unitOfWork := context_db.NewGormUnitOfWork(config.DB)

	userUseCase := usecase.NewUserUseCase(unitOfWork, config.Log, config.Validate, userRepository, userRepositoryNoSQL, jwtManager)

	userHandler := rest.NewUserHandler(userUseCase, config.Log)

	blogUsecase := usecase.NewBlogUseCase(unitOfWork, config.Log, config.Validate, blogRepository)

	blogHandler := rest.NewBlogHandler(blogUsecase, config.Log)

	genericHandler := rest.NewGenericHandler(config.Log)

	// setup handler
	apiHandler := rest.NewAPIHandler(genericHandler, userHandler, blogHandler)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase, config.Log)

	routerConfig := router.RouterConfig{
		App:            config.App,
		Log:            config.Log,
		APIHandler:     *apiHandler,
		AuthMiddleware: authMiddleware,
	}
	routerConfig.Setup()
}
