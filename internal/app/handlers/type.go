package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/bloggor/internal/pkg/models"
	"github.com/nixpig/bloggor/pkg/api"
)

func AdminTypeGetHandler(c *fiber.Ctx) error {
	a := api.WithContext(c)

	types, err := models.Query.Type.GetAll()
	if err != nil {
		return err
	}

	return c.Render("pages/admin/type", &fiber.Map{
		"Api":     a,
		"Context": c,
		"Types":   types,
	}, "layouts/admin")
}

func AdminTypePostHander(c *fiber.Ctx) error {
	a := api.WithContext(c)

	name := c.FormValue("name")
	template := c.FormValue("template")
	slug := c.FormValue("slug")

	newType := models.NewTypeData{
		Name:     name,
		Template: template,
		Slug:     slug,
	}

	createdType, err := models.Query.Type.Create(newType)
	if err != nil {
		return err
	}

	return c.Render("pages/admin/type", &fiber.Map{
		"Api":         a,
		"Context":     c,
		"CreatedType": createdType,
	}, "layouts/admin")
}
