package article

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ArticleService interface {
	DeleteById(id int) error
	Create(article *ArticleNewRequestDto) (*ArticleResponseDto, error)
	GetAll() (*[]ArticleResponseDto, error)
	GetManyByAttribute(attr, value string) (*[]ArticleResponseDto, error)
	GetByAttribute(attr, value string) (*ArticleResponseDto, error)
	Update(article *ArticleUpdateRequestDto) (*ArticleResponseDto, error)
}

type ArticleServiceImpl struct {
	repo     ArticleRepository
	validate *validator.Validate
}

func NewArticleService(
	data ArticleRepository,
	validator *validator.Validate,
) ArticleServiceImpl {
	return ArticleServiceImpl{
		repo:     data,
		validate: validator,
	}
}

func (a ArticleServiceImpl) DeleteById(id int) error {
	return a.repo.DeleteById(id)
}

func (a ArticleServiceImpl) Create(article *ArticleNewRequestDto) (*ArticleResponseDto, error) {
	articleToCreate := ArticleNew{
		Title:     article.Title,
		Subtitle:  article.Subtitle,
		Slug:      article.Slug,
		Body:      article.Body,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
		TagIds:    article.TagIds,
	}

	if err := a.validate.Struct(articleToCreate); err != nil {
		return nil, err
	}

	if len(article.TagIds) == 0 {
		minTagsError := errors.New("article must have at least one tag")
		return nil, minTagsError
	}

	createdArticle, err := a.repo.Create(&articleToCreate)
	if err != nil {
		return nil, err
	}

	return &ArticleResponseDto{
		Id:        createdArticle.Id,
		Title:     createdArticle.Title,
		Subtitle:  createdArticle.Subtitle,
		Slug:      createdArticle.Slug,
		Body:      createdArticle.Body,
		CreatedAt: createdArticle.CreatedAt,
		UpdatedAt: createdArticle.UpdatedAt,
		Tags:      createdArticle.Tags,
	}, nil
}

func (a ArticleServiceImpl) GetAll() (*[]ArticleResponseDto, error) {
	articles, err := a.repo.GetAll()
	if err != nil {
		return nil, err
	}

	allArticles := make([]ArticleResponseDto, len(*articles))

	for index, article := range *articles {
		allArticles[index] = ArticleResponseDto{
			Id:        article.Id,
			Title:     article.Title,
			Subtitle:  article.Subtitle,
			Slug:      article.Slug,
			Body:      article.Body,
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
			Tags:      article.Tags,
		}
	}

	return &allArticles, nil
}

func (a ArticleServiceImpl) GetManyByAttribute(attr, value string) (*[]ArticleResponseDto, error) {
	articles, err := a.repo.GetManyByAttribute(attr, value)
	if err != nil {
		return nil, err
	}

	allArticles := make([]ArticleResponseDto, len(*articles))

	for index, article := range *articles {
		allArticles[index] = ArticleResponseDto{
			Id:        article.Id,
			Title:     article.Title,
			Subtitle:  article.Subtitle,
			Slug:      article.Slug,
			Body:      article.Body,
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
			Tags:      article.Tags,
		}
	}

	return &allArticles, nil
}

func (a ArticleServiceImpl) GetByAttribute(attr, value string) (*ArticleResponseDto, error) {
	article, err := a.repo.GetByAttribute(attr, value)
	if err != nil {
		return nil, err
	}

	return &ArticleResponseDto{
		Id:        article.Id,
		Title:     article.Title,
		Subtitle:  article.Subtitle,
		Slug:      article.Slug,
		Body:      article.Body,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
		Tags:      article.Tags,
	}, nil
}

func (a ArticleServiceImpl) Update(article *ArticleUpdateRequestDto) (*ArticleResponseDto, error) {
	articleToUpdate := UpdateArticle{
		Id:        article.Id,
		Title:     article.Title,
		Subtitle:  article.Subtitle,
		Slug:      article.Slug,
		Body:      article.Body,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
		TagIds:    article.TagIds,
	}

	if err := a.validate.Struct(articleToUpdate); err != nil {
		return nil, err
	}

	updatedArticle, err := a.repo.Update(&articleToUpdate)
	if err != nil {
		return nil, err
	}

	return &ArticleResponseDto{
		Id:        updatedArticle.Id,
		Title:     updatedArticle.Title,
		Subtitle:  updatedArticle.Subtitle,
		Slug:      updatedArticle.Slug,
		Body:      updatedArticle.Body,
		CreatedAt: updatedArticle.CreatedAt,
		UpdatedAt: updatedArticle.UpdatedAt,
		Tags:      updatedArticle.Tags,
	}, nil
}
