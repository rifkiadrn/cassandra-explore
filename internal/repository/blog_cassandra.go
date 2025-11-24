package repository

import (
	"context"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/rifkiadrn/cassandra-explore/internal/entity"
)

type BlogRepositoryNoSQL struct {
	db *gocql.Session
}

func NewBlogRepositoryNoSQL(db *gocql.Session) BlogRepositoryNoSQL {
	return BlogRepositoryNoSQL{
		db: db,
	}
}

// Create creates a new blog
func (r BlogRepositoryNoSQL) Create(ctx context.Context, blogEntity entity.Blog) (*entity.Blog, error) {
	// Create blog in Cassandra

	authorIdStr := blogEntity.AuthorID.String()

	authorId, _ := gocql.ParseUUID(authorIdStr)

	blogId, _ := gocql.ParseUUID(blogEntity.ID.String())

	if err := r.db.Query(`INSERT INTO blogs.blogs_by_author(author_id, username, id, content, ts) VALUES (?, ?, ?, ?, ?)`, authorId, blogEntity.Username, blogId, blogEntity.Content, gocql.UUIDFromTime(blogEntity.Ts)).Exec(); err != nil {
		return nil, err
	}

	return &blogEntity, nil
}
