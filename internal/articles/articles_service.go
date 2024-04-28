package articles

import (
	"github.com/go-playground/validator/v10"
	"github.com/nixpig/dunce/pkg/logging"
)

type ArticleService struct {
	data     ArticleDataInterface
	validate *validator.Validate
	log      logging.Logger
}

type ArticleServiceInterface interface {
	Create(article *Article) (*Article, error)
	GetAll() (*[]Article, error)
	GetBySlug(slug string) (*Article, error)
	Update(article *Article) (*Article, error)
}

func NewArticleService(
	data ArticleDataInterface,
	validator *validator.Validate,
	log logging.Logger,
) ArticleService {
	return ArticleService{
		data:     data,
		validate: validator,
		log:      log,
	}
}

func (as ArticleService) Create(article *Article) (*Article, error) {
	createdArticle, err := as.data.Create(article)
	if err != nil {
		as.log.Error(err.Error())
		return nil, err
	}

	return createdArticle, nil
}

func (as ArticleService) GetAll() (*[]Article, error) {
	articles, err := as.data.GetAll()
	if err != nil {
		as.log.Error(err.Error())
		return nil, err
	}

	return articles, nil
}

func (as ArticleService) GetBySlug(slug string) (*Article, error) {
	article, err := as.data.GetBySlug(slug)
	if err != nil {
		as.log.Error(err.Error())
		return nil, err
	}

	return article, nil
}

func (as ArticleService) Update(article *Article) (*Article, error) {
	updatedArticle, err := as.data.Update(article)
	if err != nil {
		as.log.Error(err.Error())
		return nil, err
	}

	return updatedArticle, nil
}
