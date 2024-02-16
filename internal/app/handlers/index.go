package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/dunce/pkg/api"
)

func IndexHandler(c *fiber.Ctx) error {
	a := api.WithContext(c)

	return c.Render("index", fiber.Map{
		"Api":   a,
		"Title": "Hello, pig!",
	})
}
