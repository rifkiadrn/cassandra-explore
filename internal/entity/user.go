package entity

import (
	"time"

	"github.com/google/uuid"
)

// User is a struct that represents a user domain entity
type User struct {
	ID        uuid.UUID `json:"id,omitempty"` // Omit if zero UUID
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Token     string    `json:"token,omitempty"`      // Omit if empty
	Password  string    `json:"-"`                    // Never include in JSON
	IsOnline  bool      `json:"is_online,omitempty"`  // Omit if false
	LastSeen  time.Time `json:"last_seen,omitempty"`  // Omit if zero time
	CreatedAt time.Time `json:"created_at,omitempty"` // Omit if zero time
	UpdatedAt time.Time `json:"updated_at,omitempty"` // Omit if zero time
}

type Auth struct {
	ID       string
	Username string
	Token    string
}
