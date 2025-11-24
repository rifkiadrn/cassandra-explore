package repository

import (
	"context"

	context_db "github.com/rifkiadrn/cassandra-explore/internal/context/db"
	"github.com/rifkiadrn/cassandra-explore/internal/entity"
	model_db "github.com/rifkiadrn/cassandra-explore/internal/model/db"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BlogRepository struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewBlogRepository(db *gorm.DB, log *logrus.Logger) BlogRepository {
	return BlogRepository{
		db:  db,
		log: log,
	}
}

func (r *BlogRepository) getDB(ctx context.Context) *gorm.DB {
	if tx := context_db.GetTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

// entityToDBBlog converts domain entity to DB model
func (r BlogRepository) entityToDBBlog(e entity.Blog) model_db.Blog {
	return model_db.Blog{
		ID:       e.ID,
		AuthorID: e.AuthorID,
		Username: e.Username,
		Content:  e.Content,
	}
}

// dbToEntityBlog converts DB model to domain entity pointer
func (r BlogRepository) dbToEntityBlog(db model_db.Blog) *entity.Blog {
	return &entity.Blog{
		ID:       db.ID,
		AuthorID: db.AuthorID,
		Username: db.Username,
		Content:  db.Content,
	}
}

// Create creates a new blog
func (r BlogRepository) Create(ctx context.Context, blog entity.Blog) (*entity.Blog, error) {
	dbBlog := r.entityToDBBlog(blog)

	if err := r.getDB(ctx).Create(&dbBlog).Error; err != nil {
		return nil, err
	}

	return r.dbToEntityBlog(dbBlog), nil
}

// FindAll finds all blogs for a user
func (r BlogRepository) FindAll(ctx context.Context, userID string) ([]*entity.Blog, error) {
	var dbBlogs []model_db.Blog
	if err := r.getDB(ctx).Where("user_id = ?", userID).Find(&dbBlogs).Error; err != nil {
		return nil, err
	}

	// Convert to entities
	blogs := make([]*entity.Blog, len(dbBlogs))
	for i, dbBlog := range dbBlogs {
		blogs[i] = r.dbToEntityBlog(dbBlog)
	}

	return blogs, nil
}
