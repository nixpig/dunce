package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/bloggor/internal/pkg/models"
)

func WithContext(c *fiber.Ctx) map[string]interface{} {
	return map[string]interface{}{
		// Users
		"GetUser":  func() *models.UserData { return GetUser(c) },
		"GetUsers": GetUsers,
		// TODO: GetUserById
		// TODO:  GetUserByUsername
		// TODO: GetUserByEmail
		// TODO: GetLoggedInUser

		// TODO: Login
		// TODO: Logout

		// Tags
		"GetTags": GetTags,

		// Types
		"GetTypes": GetTypes,

		// Articles
		"GetArticles": GetArticles,
		// TODO: GetArticleById
		"GetArticleBySlug": func() *models.ArticleData { return GetArticleBySlug(c.Params("article_slug")) },
		// TODO: GetAllArticlesByAuthor
		// TODO: GetAllArticlesByTag
		"GetArticlesByTypeName": func() *map[string]models.ArticleData { return GetArticlesByTypeName(c.Params("article_type")) },

		// Site
		// TODO: SiteName
		// TODO: SiteDescription
		// TODO: SiteUrl
		// TODO: SiteOwner
	}
}
