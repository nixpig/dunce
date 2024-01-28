package handlers

import (
	"fmt"
	"strconv"
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
	tagIds := c.FormValue("tag_ids")

	updatedAt := time.Now()
	createdAt := time.Now()

	typeId, err := strconv.Atoi(c.FormValue("type_id"))
	if err != nil {
		return err
	}

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

func ArticleHandler(c *fiber.Ctx) error {
	a := api.WithContext(c)
	articleTypeName := c.Params("article_type")

	return c.Render(fmt.Sprintf("pages/public/%s", articleTypeName), &fiber.Map{
		"Api":     a,
		"Context": c,
	}, "layouts/public")

}
