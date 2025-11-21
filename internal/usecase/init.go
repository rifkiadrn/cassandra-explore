package usecase

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// DBContextKey is the context key for database transactions
const DBContextKey contextKey = "db"
