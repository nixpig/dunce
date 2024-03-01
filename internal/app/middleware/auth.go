package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/dunce/internal/pkg/models"
)

type Protected struct{}

type ProtectedConfig struct {
	Filter func(c *fiber.Ctx) bool
}

func (p *Protected) New(config ProtectedConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("dunce_jwt")

		claims, err := models.ValidateToken(tokenString)
		if err != nil {
			return c.Redirect("/login")
		}

		claimedRole, err := models.ParseRoleName(claims.UserRole)
		if err != nil {
			return c.Redirect("/login")
		}

		if claimedRole != models.AdminRole {
			return c.Redirect("/login")
		}

		return c.Next()
	}
}
