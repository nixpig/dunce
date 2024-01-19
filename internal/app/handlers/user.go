package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/bloggor/internal/pkg/models"
)

func UserHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	if len(id) == 0 {
		users, err := models.Query.User.GetAll()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Render("500", fiber.Map{
				"Status": "Something went wrong",
			})
		}

		return c.Render("user", fiber.Map{
			"Users": users,
		})
	}

	user_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("500", fiber.Map{
			"Status": "Invalid user id provided",
		})
	}

	user, err := models.Query.User.GetById(user_id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("500", fiber.Map{
			"Status": "Couldn't get user",
		})
	}

	return c.Render("user", fiber.Map{
		"User": user,
	})
}
