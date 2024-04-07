package user

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	dunce "github.com/nixpig/dunce/internal/pkg"
	"github.com/nixpig/dunce/internal/pkg/models"
)

type UserController struct {
	service dunce.Service[UserRequest, UserResponse]
}

func NewUserController(service dunce.Service[UserRequest, UserResponse]) UserController {
	return UserController{service}
}

func (u *UserController) HandleCreate(c *fiber.Ctx) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	link := c.FormValue("link")
	password := c.FormValue("password")
	role := c.FormValue("role")

	createdUser, err := u.service.Create(UserRequest{
		Username: username,
		Email:    email,
		Link:     link,
		Password: password,
		Role:     role,
	})
	if err != nil {
		// TODO: check for type of error and return correct status
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

func (u *UserController) HandleGetAll(c *fiber.Ctx) error {
	users, err := u.service.GetAll()
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
