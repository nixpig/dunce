package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/dunce/internal/pkg/models"
)

func AdminLoginGetHandler(c *fiber.Ctx) error {
	return c.Render("pages/admin/login", &fiber.Map{}, "layouts/admin")
}

func AdminLoginPostHandler(c *fiber.Ctx) error {
	user := models.LoginDetails{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}

	token, err := models.Query.Login.WithUsernamePassword(&user)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	fmt.Println("TOKEN: ", token)

	return c.SendStatus(fiber.StatusOK)
}
