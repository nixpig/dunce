package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/bloggor/internal/pkg/models"
	"github.com/nixpig/bloggor/pkg/api"
)

func UserGetHandler(c *fiber.Ctx) error {
	a := api.WithContext(c)

	return c.Render("user", fiber.Map{
		"Context": c,
		"Api":     a,
	})
}

func AdminUserGetHandler(c *fiber.Ctx) error {
	a := api.WithContext(c)

	return c.Render("admin_user", fiber.Map{
		"Context": c,
		"Api":     a,
		"IsEditable": func(userId int, editId string) bool {
			editIdConv, err := strconv.Atoi(editId)
			if err != nil {
				return false
			}

			return editIdConv == userId
		},
	}, "layouts/admin")
}

func AdminUserPostHandler(c *fiber.Ctx) error {
	a := api.WithContext(c)

	username := c.FormValue("username")
	email := c.FormValue("email")
	link := c.FormValue("link")
	password := c.FormValue("password")

	role, err := models.ParseRoleName(c.FormValue("role"))
	if err != nil {
		return err
	}

	createdUser, err := models.Query.User.Create(&models.NewUserData{
		Username: username,
		Email:    email,
		Link:     link,
		Password: password,
		Role:     role,
	})
	if err != nil {
		fmt.Println("ERROR")
		return c.Render("admin_user", fiber.Map{
			"Api":     a,
			"Context": c,
			"IsEditable": func(userId int, editId string) bool {
				editIdConv, err := strconv.Atoi(editId)
				if err != nil {
					return false
				}

				return editIdConv == userId
			},
			"Errors": err,
		}, "layouts/admin")
	}

	return c.Render("admin_user", fiber.Map{
		"Api":     a,
		"Context": c,
		"IsEditable": func(userId int, editId string) bool {
			editIdConv, err := strconv.Atoi(editId)
			if err != nil {
				return false
			}

			return editIdConv == userId
		},
		"CreatedUser": createdUser,
	}, "layouts/admin")
}
