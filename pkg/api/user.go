package api

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/bloggor/internal/pkg/models"
)

func GetUser(c *fiber.Ctx) *models.UserData {
	user_id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		fmt.Println("Invalid ID: ", err)
		return nil
	}

	user, err := models.Query.User.GetById(user_id)
	if err != nil {
		fmt.Println("Unable to get user: ", user_id, err)
		return nil
	}

	return user
}

func GetUsers() map[string]models.UserData {
	users, err := models.Query.User.GetAll()
	if err != nil {
		fmt.Println("Error getting users: ", err)
		return nil
	}

	usermap := make(map[string]models.UserData)

	for index, item := range *users {
		usermap[strconv.Itoa(index)] = item
	}

	return usermap
}
