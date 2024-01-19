package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var message string
	var info string

	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	switch code {
	case fiber.StatusInternalServerError:
		info = "Internal Server Error"
	case fiber.StatusNotFound:
		info = "Not Found"
	}

	return c.Status(code).Render("error", fiber.Map{
		"Code":    code,
		"Message": message,
		"Info":    info,
	})
}
