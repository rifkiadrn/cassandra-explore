package config

import (
	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

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
	userRepository := repository.NewUserRepository(config.Log)

	// setup JWT manager
	jwtManager := utils.NewJWTManager(config.Config.GetString("SECRET_KEY")) // TODO: move to config

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.NoSQLDB, config.Log, config.Validate, userRepository, jwtManager)

	userHandler := rest.NewUserHandler(userUseCase, config.Log)

	blogRepository := repository.NewBlogRepository(config.Log)

	blogUsecase := usecase.NewBlogUseCase(config.DB, config.NoSQLDB, config.Log, config.Validate, blogRepository)

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
