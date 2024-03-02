package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/dunce/internal/pkg/models"
	"github.com/nixpig/dunce/pkg/api"
)

const longFormat = "2006-01-02 15:04:05.999999999 -0700 MST"

func AdminArticleGetHandler(c *fiber.Ctx) error {
	pathParts := strings.Split(c.Path(), "/")
	page := pathParts[len(pathParts)-1]

	paramId := c.Params("id")
	if len(paramId) > 0 {
		id, err := strconv.Atoi(paramId)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		article, err := models.Query.Article.GetById(id)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var editable bool
		if c.Query("edit") == "true" {
			editable = true
		}

		if editable {
			return c.Render("pages/admin/articles", &fiber.Map{
				"Article":  article,
				"Page":     page,
				"Editable": editable,
			}, "layouts/admin")

		}

		return c.Render("fragments/admin/articles/article_table_row_view", &fiber.Map{
			"Article": article,
		})
	}

	articles, err := models.Query.Article.GetAll()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Render("pages/admin/articles", &fiber.Map{
		"Page":     page,
		"Articles": articles,
	}, "layouts/admin")
}

func AdminArticleDeleteHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := models.Query.Article.DeleteById(id); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).Send([]byte{})
}

func AdminArticlePutHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Println("error converting 'id' param to an int: ", err, idParam)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	fmt.Println("ID: ", id)

	typeId, err := strconv.Atoi(c.FormValue("type_id"))
	if err != nil {
		fmt.Println("error converting 'type_id' form value to an int: ", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	userId, err := strconv.Atoi(c.FormValue("user_id"))
	if err != nil {
		fmt.Println("error converting 'user_id' form value to an int: ", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	fmt.Println("CREATED_AT: ", c.FormValue("created_at"))
	createdAt, err := time.Parse(longFormat, c.FormValue("created_at"))
	if err != nil {
		fmt.Println("error converting 'created_at' form value to a time: ", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	articleUpdate := models.Article{
		Id: id,
		ArticleData: models.ArticleData{
			Title:     c.FormValue("title"),
			Subtitle:  c.FormValue("subtitle"),
			Slug:      c.FormValue("slug"),
			Body:      c.FormValue("body"),
			UpdatedAt: time.Now(),
			CreatedAt: createdAt,
			TypeId:    typeId,
			UserId:    userId,
			TagIds:    c.FormValue("tag_ids"),
		},
	}

	updated, err := models.Query.Article.UpdateById(id, articleUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("fragments/admin/shared/admin_table_errors", &fiber.Map{
			"Errors": []string{err.Error()},
		})
	}

	return c.Render("fragments/admin/articles/article_table_row_view", &fiber.Map{
		"Article": updated,
	})
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

	userId, err := strconv.Atoi(c.FormValue("user_id"))
	if err != nil {
		return err
	}

	newArticle := models.ArticleData{
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

	return c.Render("pages/admin/articles", &fiber.Map{
		"Api":            a,
		"Context":        c,
		"CreatedArticle": createdArticle,
	}, "layouts/admin")
}

func AdminArticleCreateGetHandler(c *fiber.Ctx) error {
	page := "articles"

	return c.Render("pages/admin/articles", &fiber.Map{
		"Page":     page,
		"New":      true,
		"Editable": nil,
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
