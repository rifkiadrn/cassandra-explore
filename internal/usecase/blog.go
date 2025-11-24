package usecase

import (
	"context"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rifkiadrn/cassandra-explore/internal/entity"
	authContext "github.com/rifkiadrn/cassandra-explore/internal/handler/rest/context"
	"github.com/sirupsen/logrus"
)

type IBlog interface {
	Create(ctx context.Context, blog entity.Blog) (*entity.Blog, error)
	FindAll(ctx context.Context, userID string) ([]*entity.Blog, error)
}

type BlogUseCase struct {
	uow            UnitOfWork
	log            *logrus.Logger
	validate       *validator.Validate
	blogRepository IBlog
}

func NewBlogUseCase(uow UnitOfWork, logger *logrus.Logger, validate *validator.Validate,
	blogRepository IBlog) BlogUseCase {
	return BlogUseCase{
		uow:            uow,
		log:            logger,
		validate:       validate,
		blogRepository: blogRepository,
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
	if err := b.validate.Struct(request); err != nil {
		b.log.Warnf("Invalid request body : %+v", err)
		return entity.Blog{}, fiber.ErrBadRequest
	}

	res, err := b.blogRepository.Create(ctx, blogEntity)
	if err != nil {
		return entity.Blog{}, err
	}

	return *res, nil
}

func (b BlogUseCase) GetBlogs(ctx context.Context) ([]entity.Blog, error) {
	// Get authenticated user
	user, err := authContext.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get blogs via repository
	blogs, err := b.blogRepository.FindAll(ctx, user.ID.String())
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
