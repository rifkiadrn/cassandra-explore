package repository

import (
	"context"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/rifkiadrn/cassandra-explore/internal/entity"
)

type UserRepositoryNoSQL struct {
	db *gocql.Session
}

func NewUserRepositoryNoSQL(db *gocql.Session) UserRepositoryNoSQL {
	return UserRepositoryNoSQL{
		db: db,
	}
}

// Create creates a new blog
func (r UserRepositoryNoSQL) Create(ctx context.Context, userEntity entity.User) (*entity.User, error) {
	// Create blog in Cassandra
	userId, _ := gocql.ParseUUID(userEntity.ID.String())

	if err := r.db.Query(`INSERT INTO users (id, name, username) VALUES (?, ?, ?)`, userId, userEntity.Name, userEntity.Username).Exec(); err != nil {
		return nil, err
	}

	return &userEntity, nil
}
