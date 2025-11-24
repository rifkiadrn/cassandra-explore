package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	viper.AutomaticEnv()
	username := viper.GetString("database.username")
	if viper.GetString("DB_USERNAME") != "" {
		username = viper.GetString("DB_USERNAME")
	}

	password := viper.GetString("database.password")
	if viper.GetString("DB_PASSWORD") != "" {
		password = viper.GetString("DB_PASSWORD")
	}

	host := viper.GetString("database.host")
	if viper.GetString("DB_HOST") != "" {
		host = viper.GetString("DB_HOST")
	}

	port := viper.GetInt("database.port")
	if viper.GetInt("DB_PORT") != 0 {
		port = viper.GetInt("DB_PORT")
	}

	database := viper.GetString("database.name")
	if viper.GetString("DB_NAME") != "" {
		database = viper.GetString("DB_NAME")
	}

	idleConnection := viper.GetInt("database.pool.idle")
	maxConnection := viper.GetInt("database.pool.max")
	maxLifeTimeConnection := viper.GetInt("database.pool.lifetime")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta", host, username, password, database, port)
	log.Printf(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	return db
}
