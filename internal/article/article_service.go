package article

import (
	"github.com/go-playground/validator/v10"
	"github.com/nixpig/dunce/pkg"
)

type ArticleService struct {
	data     pkg.Repository[Article]
	validate *validator.Validate
	log      pkg.Logger
}

func NewArticleService(
	data pkg.Repository[Article],
	validator *validator.Validate,
	log pkg.Logger,
) ArticleService {
	return ArticleService{
		data:     data,
		validate: validator,
		log:      log,
	}
}

func (a ArticleService) DeleteById(id int) error {
	return nil
}

func (a ArticleService) Create(article *Article) (*Article, error) {
	createdArticle, err := a.data.Create(article)
	if err != nil {
		a.log.Error(err.Error())
		return nil, err
	}

	return createdArticle, nil
}

func (a ArticleService) GetAll() (*[]Article, error) {
	articles, err := a.data.GetAll()
	if err != nil {
		a.log.Error(err.Error())
		return nil, err
	}

	return articles, nil
}

func (a ArticleService) GetBySlug(slug string) (*Article, error) {
	article, err := a.data.GetBySlug(slug)
	if err != nil {
		a.log.Error(err.Error())
		return nil, err
	}

	return article, nil
}

func (a ArticleService) Update(article *Article) (*Article, error) {
	updatedArticle, err := a.data.Update(article)
	if err != nil {
		a.log.Error(err.Error())
		return nil, err
	}

	return updatedArticle, nil
}
