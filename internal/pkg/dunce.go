package dunce

import "github.com/gofiber/fiber/v2"

type Controller interface {
	HandleCreate(c *fiber.Ctx) error
	HandleGetAll(c *fiber.Ctx) error
	HandleGetOne(c *fiber.Ctx) error
}

type Service[Request any, Response any] interface {
	Create(record Request) (*Response, error)
	GetAll() (*[]Response, error)
}

type Data[Request any, Response any] interface {
	Create(record Request) (*Response, error)
	GetAll() (*[]Response, error)
}
