package api

import (
	"fmt"
	"strconv"

	"github.com/nixpig/dunce/internal/pkg/models"
)

func GetArticles() map[string]models.Article {
	articles, err := models.Query.Article.GetAll()
	if err != nil {
		fmt.Println(fmt.Errorf("unable to get articles: %v", err))
		return nil
	}

	articlemap := make(map[string]models.Article)

	for index, item := range *articles {
		articlemap[strconv.Itoa(index)] = item
	}

	return articlemap
}

func GetArticlesByTypeName(typeName string) *map[string]models.Article {
	articles, err := models.Query.Article.GetByTypeName(typeName)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to get articles by type name: %v", err))
		return nil
	}

	articlemap := make(map[string]models.Article)

	for index, article := range *articles {
		articlemap[strconv.Itoa(index)] = article

	}

	return &articlemap
}

func GetArticleBySlug(slug string) *models.Article {
	article, err := models.Query.Article.GetBySlug(slug)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to get article by slug: %v", err))
		return nil
	}

	return article
}
