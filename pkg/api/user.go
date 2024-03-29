package api

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/dunce/internal/pkg/models"
)

func GetUser(c *fiber.Ctx) *models.User {
	user_id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		fmt.Println(fmt.Errorf("unable to convert ID: %v", err))
		return nil
	}

	user, err := models.Query.User.GetById(user_id)
	if err != nil {
		fmt.Println(fmt.Errorf("Unable to get user: %v, %v", user_id, err))
		return nil
	}

	return user
}

func GetUsers() map[string]models.User {
	users, err := models.Query.User.GetAll()
	if err != nil {
		fmt.Println(fmt.Errorf("Error getting users: %v", err))
		return nil
	}

	usermap := make(map[string]models.User)

	for index, item := range *users {
		usermap[strconv.Itoa(index)] = item
	}

	return usermap
}
