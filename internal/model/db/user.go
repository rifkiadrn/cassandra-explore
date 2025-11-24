package model_db

import (
	"github.com/google/uuid"
)

// User represents the database model for users
type User struct {
	ID        uuid.UUID `gorm:"column:id;primaryKey;default:gen_random_uuid()"` // Auto-generate UUID
	Name      string    `gorm:"column:name;not null"`
	Username  string    `gorm:"column:username;not null;uniqueIndex"`
	Password  string    `gorm:"column:password;not null"`                        // Never empty
	Token     string    `gorm:"column:token"`                                    // Can be empty
	CreatedAt int64     `gorm:"column:created_at;autoCreateTime"`                // Auto-generated
	UpdatedAt int64     `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"` // Auto-generated
}

func (u *User) TableName() string {
	return "cassandra_users.users"
}
