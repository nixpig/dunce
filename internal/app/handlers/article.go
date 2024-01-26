package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/bloggor/internal/pkg/models"
	"github.com/nixpig/bloggor/pkg/api"
)

func AdminArticleGetHandler(c *fiber.Ctx) error {
	a := api.WithContext(c)

	return c.Render("pages/admin/article", &fiber.Map{
		"Api":     a,
		"Context": c,
	}, "layouts/admin")

}

func AdminArticlePostHandler(c *fiber.Ctx) error {
	a := api.WithContext(c)

	title := c.FormValue("title")
	subtitle := c.FormValue("subtitle")
	slug := c.FormValue("slug")
	body := c.FormValue("body")
	typeId := c.FormValue("type_id")
	tagIds := c.FormValue("tag_ids")

	updatedAt := time.Now()
	createdAt := time.Now()

	userId := 3

	newArticle := models.NewArticleData{
		Title:     title,
		Subtitle:  subtitle,
		Slug:      slug,
		Body:      body,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		TypeId:    typeId,
		UserId:    userId,
		TagIds:    tagIds,
	}

	createdArticle, err := models.Query.Article.Create(newArticle)
	if err != nil {
		return err
	}

	return c.Render("pages/admin/article", &fiber.Map{
		"Api":            a,
		"Context":        c,
		"CreatedArticle": createdArticle,
	}, "layouts/admin")
}
