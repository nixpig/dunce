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
