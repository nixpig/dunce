package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/bloggor/internal/pkg/models"
	"github.com/nixpig/bloggor/pkg/api"
)

func AdminTagGetHandler(c *fiber.Ctx) error {
	a := api.WithContext(c)

	return c.Render("pages/admin/tags", &fiber.Map{
		"Context": c,
		"Api":     a,
	}, "layouts/admin")
}

func AdminTagDeleteHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	if err := models.Query.Tag.Delete(id); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).Send([]byte{})
}

func AdminTagPostHandler(c *fiber.Ctx) error {
	a := api.WithContext(c)

	name := c.FormValue("name")
	slug := c.FormValue("slug")

	newTag := models.NewTagData{
		Name: name,
		Slug: slug,
	}

	createdTag, err := models.Query.Tag.Create(newTag)
	if err != nil {
		return err
	}

	return c.Render("pages/admin/tags", &fiber.Map{
		"Api":        a,
		"Context":    c,
		"CreatedTag": createdTag,
	})
}
