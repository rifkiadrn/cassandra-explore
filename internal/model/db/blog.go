package model_db

import (
	"github.com/google/uuid"
)

type Blog struct {
	ID       uuid.UUID `gorm:"column:id;primaryKey;default:gen_random_uuid()"` // Auto-generate UUID
	Content  string    `gorm:"column:content;not null"`
	AuthorID uuid.UUID `gorm:"column:user_id;not null"`
	Username string    `gorm:"column:username;not null"`
	Ts       int64     `gorm:"column:ts;autoCreateTime"`
}
