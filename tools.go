//go:build tools
// +build tools

package tools

import (
	_ "github.com/go-playground/validator/v10"
	_ "github.com/gofiber/fiber/v2"
	_ "github.com/google/uuid"
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/sirupsen/logrus"
	_ "github.com/spf13/viper"
	_ "github.com/stretchr/testify"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
)
