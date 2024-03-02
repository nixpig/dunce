package handlers

import (
	"time"

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

	c.Cookie(&fiber.Cookie{
		Name:     "dunce_jwt",
		Value:    token,
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 1),
		SameSite: "strict",
		HTTPOnly: true,
	})

	return c.SendStatus(fiber.StatusOK)
}

func AdminLogoutHandler(c *fiber.Ctx) error {
	tokenString := c.Cookies("dunce_jwt")

	if tokenString != "" {
		claims, err := models.ValidateToken(tokenString)
		if err != nil {
			return err
		}

		if err := models.Query.Login.Logout(claims.UserId); err != nil {
			return err
		}

	}

	// set cookie expiry to now
	c.Cookie(&fiber.Cookie{
		Name:     "dunce_jwt",
		Value:    "",
		Secure:   true,
		Expires:  time.Now(),
		SameSite: "strict",
		HTTPOnly: true,
	})

	return c.Redirect("/login")
}
