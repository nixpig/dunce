package api

import (
	"fmt"
	"strconv"

	"github.com/nixpig/bloggor/internal/pkg/models"
)

func GetArticles() map[string]models.ArticleData {
	articles, err := models.Query.Article.GetAll()
	if err != nil {
		fmt.Println(fmt.Errorf("unable to get articles: %v", err))
		return nil
	}

	articlemap := make(map[string]models.ArticleData)

	for index, item := range *articles {
		articlemap[strconv.Itoa(index)] = item
	}

	return articlemap
}

func GetArticlesByTypeName(typeName string) *map[string]models.ArticleData {
	articles, err := models.Query.Article.GetByTypeName(typeName)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to get articles by type name: %v", err))
		return nil
	}

	articlemap := make(map[string]models.ArticleData)

	for index, article := range *articles {
		articlemap[strconv.Itoa(index)] = article

	}

	return &articlemap
}

func GetArticleBySlug(slug string) *models.ArticleData {
	article, err := models.Query.Article.GetBySlug(slug)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to get article by slug: %v", err))
		return nil
	}

	return article
}
