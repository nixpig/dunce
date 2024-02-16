package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/bloggor/internal/pkg/models"
)

func AdminSiteGetHandler(c *fiber.Ctx) error {
	siteData, err := models.Query.Site.Get()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	users, err := models.Query.User.GetAll()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Render("pages/admin/site", &fiber.Map{
		"Name":        siteData.Name,
		"Description": siteData.Description,
		"Url":         siteData.Url,
		"Owner":       siteData.Owner,
		"Users":       users,
	}, "layouts/admin")
}

func AdminSitePostHandler(c *fiber.Ctx) error {
	ownerId, err := strconv.Atoi(c.FormValue("owner"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	siteData := models.SiteData{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Url:         c.FormValue("url"),
		Owner:       ownerId,
	}

	updatedSiteData, err := models.Query.Site.Update(siteData)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	users, err := models.Query.User.GetAll()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Render("fragments/admin/site/site_form", &fiber.Map{
		"Name":        updatedSiteData.Name,
		"Description": updatedSiteData.Description,
		"Url":         updatedSiteData.Url,
		"Owner":       updatedSiteData.Owner,
		"Users":       users,
	})
}
