package model_api

import "github.com/google/uuid"

// VerifyUserRequest represents the internal API request for user verification
// This is an internal model since it's not defined in the OpenAPI spec
type VerifyUserRequest struct {
	Token string `json:"token" validate:"required"`
}

// Auth represents authentication information for internal API responses
// This is an internal model since it's not defined in the OpenAPI spec
type Auth struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Token    string    `json:"token"`
}
