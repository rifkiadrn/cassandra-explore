package entity

import (
	"time"

	"github.com/google/uuid"
)

type Blog struct {
	ID       uuid.UUID `json:"id,omitempty"` // Omit if zero UUID
	Content  string    `json:"content"`
	AuthorID uuid.UUID `json:"author_id,omitempty"` // Omit if zero UUID
	Username string    `json:"username"`
	Ts       time.Time `json:"ts,omitempty"` // Omit if nil
}
