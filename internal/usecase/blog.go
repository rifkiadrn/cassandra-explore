package usecase

import (
	"context"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	context_db "github.com/rifkiadrn/cassandra-explore/internal/context/db"
	"github.com/rifkiadrn/cassandra-explore/internal/entity"
	authContext "github.com/rifkiadrn/cassandra-explore/internal/handler/rest/context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IBlog interface {
	Create(ctx context.Context, blog entity.Blog) (*entity.Blog, error)
	FindAll(ctx context.Context, userID string) ([]*entity.Blog, error)
}

type BlogUseCase struct {
	DB             *gorm.DB
	NoSQLDB        *gocql.Session
	Log            *logrus.Logger
	Validate       *validator.Validate
	BlogRepository IBlog
}

func NewBlogUseCase(db *gorm.DB, noSQLDB *gocql.Session, logger *logrus.Logger, validate *validator.Validate,
	blogRepository IBlog) BlogUseCase {
	return BlogUseCase{
		DB:             db,
		NoSQLDB:        noSQLDB,
		Log:            logger,
		Validate:       validate,
		BlogRepository: blogRepository,
	}
}

func (b BlogUseCase) CreateBlog(ctx context.Context, request entity.Blog) (entity.Blog, error) {
	// Get authenticated user
	user, err := authContext.GetUserFromContext(ctx)
	if err != nil {
		return entity.Blog{}, err
	}

	time := gocql.TimeUUID()

	// Create domain entity
	blogEntity := entity.Blog{
		ID:       uuid.New(),
		AuthorID: user.ID,
		Username: user.Username,
		Content:  request.Content,
		Ts:       time.Time(),
	}

	// Validate request
	if err := b.Validate.Struct(request); err != nil {
		b.Log.Warnf("Invalid request body : %+v", err)
		return entity.Blog{}, fiber.ErrBadRequest
	}

	// Create blog in Cassandra

	authorIdStr := blogEntity.AuthorID.String()

	authorId, _ := gocql.ParseUUID(authorIdStr)

	blogId, _ := gocql.ParseUUID(blogEntity.ID.String())

	if err := b.NoSQLDB.Query(`INSERT INTO blogs.blogs_by_author(author_id, username, id, content, ts) VALUES (?, ?, ?, ?, ?)`, authorId, blogEntity.Username, blogId, blogEntity.Content, gocql.UUIDFromTime(blogEntity.Ts)).Exec(); err != nil {
		return entity.Blog{}, err
	}

	return blogEntity, nil

	// // Start transaction
	// tx, txCtx := context_db.BeginTxWithContext(ctx, b.DB)
	// defer tx.Rollback()

	// // Create blog via repository
	// createdBlog, err := b.BlogRepository.Create(txCtx, blogEntity)
	// if err != nil {
	// 	b.Log.Warnf("Failed create blog : %+v", err)
	// 	return entity.Blog{}, fiber.ErrInternalServerError
	// }

	// // Commit transaction
	// if err := tx.Commit().Error; err != nil {
	// 	b.Log.Warnf("Failed commit transaction : %+v", err)
	// 	return entity.Blog{}, fiber.ErrInternalServerError
	// }

	// return *createdBlog, nil
}

func (b BlogUseCase) GetBlogs(ctx context.Context) ([]entity.Blog, error) {
	// Get authenticated user
	user, err := authContext.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Start transaction
	tx, txCtx := context_db.BeginTxWithContext(ctx, b.DB)
	defer tx.Rollback()

	// Get blogs via repository
	blogs, err := b.BlogRepository.FindAll(txCtx, user.ID.String())
	if err != nil {
		return nil, err
	}

	// Dereference pointers to return values
	result := make([]entity.Blog, len(blogs))
	for i, blog := range blogs {
		result[i] = *blog
	}

	return result, nil
}
