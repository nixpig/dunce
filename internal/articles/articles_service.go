package articles

type ArticleService struct {
	data ArticleDataInterface
}

type ArticleServiceInterface interface {
	create(article *Article) (*Article, error)
}

func NewArticleService(data ArticleDataInterface) ArticleService {
	return ArticleService{data}
}

func (as ArticleService) create(article *Article) (*Article, error) {
	// TODO :validation

	createdArticle, err := as.data.create(article)
	if err != nil {
		return nil, err
	}

	return createdArticle, nil
}
