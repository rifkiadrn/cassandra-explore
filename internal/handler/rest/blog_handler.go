package rest

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/rifkiadrn/cassandra-explore/internal/entity"
	model "github.com/rifkiadrn/cassandra-explore/internal/model"
	"github.com/sirupsen/logrus"
)

type IBlogUseCase interface {
	CreateBlog(ctx context.Context, request entity.Blog) (entity.Blog, error)
	GetBlogs(ctx context.Context) ([]entity.Blog, error)
}

type BlogHandler struct {
	Log     *logrus.Logger
	UseCase IBlogUseCase
}

func NewBlogHandler(useCase IBlogUseCase, logger *logrus.Logger) *BlogHandler {
	return &BlogHandler{
		Log:     logger,
		UseCase: useCase,
	}
}

func (h *BlogHandler) CreateBlog(c *fiber.Ctx) error {
	request := model.CreateBlogRequest{}
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	fmt.Println("request", request)

	var blogInput entity.Blog
	err := copier.Copy(&blogInput, &request)
	if err != nil {
		return err
	}

	fmt.Println("blogInput", blogInput)

	blog, err := h.UseCase.CreateBlog(c.Context(), blogInput)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": convertToBlogResponse(blog),
	})
}

func (h *BlogHandler) Blogs(c *fiber.Ctx) error {
	blogs, err := h.UseCase.GetBlogs(c.Context())
	if err != nil {
		return err
	}

	var blogsResponse []model.Blog
	for _, blog := range blogs {
		blogsResponse = append(blogsResponse, convertToBlogResponse(blog))
	}

	return c.JSON(fiber.Map{
		"data": blogsResponse,
	})
}

func convertToBlogResponse(blog entity.Blog) model.Blog {
	authorId := blog.AuthorID.String()

	return model.Blog{
		Id:       blog.ID.String(),
		Content:  blog.Content,
		AuthorId: &authorId,
		Username: blog.Username,
		Ts:       blog.Ts.Unix(),
	}
}
