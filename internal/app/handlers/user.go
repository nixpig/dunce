package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/dunce/internal/pkg/models"
	"github.com/nixpig/dunce/pkg/api"
)

func UserGetHandler(c *fiber.Ctx) error {
	a := api.WithContext(c)

	return c.Render("user", fiber.Map{
		"Context": c,
		"Api":     a,
	})
}

func AdminUserGetHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")

	if len(idParam) == 0 {
		users, err := models.Query.User.GetAll()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		pathParts := strings.Split(c.Path(), "/")

		page := pathParts[len(pathParts)-1]

		return c.Render("pages/admin/users", fiber.Map{
			"Page":  page,
			"Users": users,
			"Roles": models.RoleNames,
		}, "layouts/admin")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := models.Query.User.GetById(id)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	editable := c.Query("edit")
	if editable == "true" {
		return c.Render("fragments/admin/users/user_table_row_edit", &fiber.Map{
			"Id":       user.Id,
			"Username": user.Username,
			"Role":     user.Role,
			"Email":    user.Email,
			"Link":     user.Link,
			"Roles":    models.RoleNames,
		})
	}

	fmt.Println(">>> SIX <<<")
	return c.Render("fragments/admin/users/user_table_row_view", &fiber.Map{
		"Id":       user.Id,
		"Username": user.Username,
		"Role":     user.Role,
		"Email":    user.Email,
		"Link":     user.Link,
		"Roles":    models.RoleNames,
	})

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

	newUser := models.UserData{
		Username: username,
		Email:    email,
		Link:     link,
		Role:     role,
	}

	createdUser, err := models.Query.User.Create(&newUser, password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("fragments/admin/shared/error_list", fiber.Map{
			"Errors": err,
		})
	}

	return c.Render("fragments/admin/users/user_table_row_view", fiber.Map{
		"Id":       createdUser.Id,
		"Role":     createdUser.Role,
		"Username": createdUser.Username,
		"Email":    createdUser.Email,
		"Link":     createdUser.Link,
	})
}

func AdminUserPutHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	username := c.FormValue("username")
	email := c.FormValue("email")
	link := c.FormValue("link")

	role, err := models.ParseRoleName(c.FormValue("role"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user := models.UserData{
		Username: username,
		Email:    email,
		Link:     link,
		Role:     role,
	}

	updatedUser, err := models.Query.User.Update(id, &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("fragments/admin/shared/admin_table_errors", &fiber.Map{
			"Errors": []string{err.Error()},
		})
	}

	return c.Render("fragments/admin/users/user_table_row_view", &fiber.Map{
		"Id":       updatedUser.Id,
		"Role":     updatedUser.Role,
		"Username": updatedUser.Username,
		"Email":    updatedUser.Email,
	})
}

func AdminUserDeleteHander(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	models.Query.User.Delete(id)

	return c.Status(200).Send([]byte{})
}
