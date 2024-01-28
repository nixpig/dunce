package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/bloggor/internal/pkg/models"
)

func WithContext(c *fiber.Ctx) map[string]interface{} {
	return map[string]interface{}{
		"GetUser":               func() *models.UserData { return GetUser(c) },
		"GetUsers":              GetUsers,
		"GetTags":               GetTags,
		"GetTypes":              GetTypes,
		"GetArticles":           GetArticles,
		"GetArticleBySlug":      func() *models.ArticleData { return GetArticleBySlug(c.Params("article_slug")) },
		"GetArticlesByTypeName": func() map[string]models.ArticleData { return GetArticlesByTypeName(c.Params("article_type")) },
	}
}
