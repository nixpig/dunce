package article

import (
	"github.com/go-playground/validator/v10"
	"github.com/nixpig/dunce/pkg"
)

type ArticleService struct {
	repo     pkg.Repository[Article]
	validate *validator.Validate
	log      pkg.Logger
}

func NewArticleService(
	data pkg.Repository[Article],
	validator *validator.Validate,
	log pkg.Logger,
) ArticleService {
	return ArticleService{
		repo:     data,
		validate: validator,
		log:      log,
	}
}

func (a ArticleService) DeleteById(id int) error {
	return a.repo.DeleteById(id)
}

func (a ArticleService) Create(article *Article) (*Article, error) {
	createdArticle, err := a.repo.Create(article)
	if err != nil {
		a.log.Error(err.Error())
		return nil, err
	}

	return createdArticle, nil
}

func (a ArticleService) GetAll() (*[]Article, error) {
	articles, err := a.repo.GetAll()
	if err != nil {
		a.log.Error(err.Error())
		return nil, err
	}

	return articles, nil
}

func (a ArticleService) GetBySlug(slug string) (*Article, error) {
	article, err := a.repo.GetBySlug(slug)
	if err != nil {
		a.log.Error(err.Error())
		return nil, err
	}

	return article, nil
}

func (a ArticleService) Update(article *Article) (*Article, error) {
	updatedArticle, err := a.repo.Update(article)
	if err != nil {
		a.log.Error(err.Error())
		return nil, err
	}

	return updatedArticle, nil
}
