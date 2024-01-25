package handlers

import (
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

	return c.Render("pages/admin/user", fiber.Map{
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
	username := c.FormValue("username")
	email := c.FormValue("email")
	link := c.FormValue("link")
	password := c.FormValue("password")

	role, err := models.ParseRoleName(c.FormValue("role"))
	if err != nil {
		return err
	}

	newUser := models.NewUserData{
		Username: username,
		Email:    email,
		Link:     link,
		Password: password,
		Role:     role,
	}

	createdUser, err := models.Query.User.Create(&newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("fragments/admin/user/add_user_error", fiber.Map{
			"Errors": []error{err},
		})
	}

	return c.Render("fragments/admin/user/add_user_success", fiber.Map{
		"CreatedUser": createdUser,
	})
}

func AdminUserPutHandler(c *fiber.Ctx) error {
	a := api.WithContext(c)

	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return err
	}

	username := c.FormValue("username")
	email := c.FormValue("email")
	link := c.FormValue("link")

	role, err := models.ParseRoleName(c.FormValue("role"))
	if err != nil {
		return err
	}

	user := models.UpdateUserData{
		Id:       id,
		Username: username,
		Email:    email,
		Link:     link,
		Role:     role,
	}

	updatedUser, err := models.Query.User.Update(&user)
	if err != nil {
		return err
	}

	return c.Render("pages/admin/user", &fiber.Map{
		"Api":         a,
		"Context":     c,
		"UpdatedUser": updatedUser,
	}, "layouts/admin")
}

func AdminUserDeleteHander(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.FormValue("delete"))
	if err != nil {
		return err
	}

	models.Query.User.Delete(id)

	return c.Status(200).Send([]byte{})
}
