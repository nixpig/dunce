package dunce

import "github.com/gofiber/fiber/v2"

type Controller interface {
	HandleSave(c *fiber.Ctx) error
	HandleGetAll(c *fiber.Ctx) error
	HandleFindOne(c *fiber.Ctx) error
}

type Service[Request any, Response any] interface {
	Save(record Request) (*Response, error)
	GetAll() (*[]Response, error)
}

type Data[Request any, Response any] interface {
	Save(record Request) (*Response, error)
	GetAll() (*[]Response, error)
}
