package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/dunce/internal/pkg/models"
)

func AdminTagGetHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")

	if len(idParam) == 0 {
		tags, err := models.Query.Tag.GetAll()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		pathParts := strings.Split(c.Path(), "/")

		page := pathParts[len(pathParts)-1]

		return c.Render("pages/admin/tags", &fiber.Map{
			"Page": page,
			"Tags": tags,
		}, "layouts/admin")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	tagData, err := models.Query.Tag.GetById(id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	editable := c.Query("edit")
	if editable == "true" {
		return c.Render("fragments/admin/tags/tag_table_row_edit", &fiber.Map{
			"Id":   tagData.Id,
			"Name": tagData.Name,
			"Slug": tagData.Slug,
		})
	}

	return c.Render("fragments/admin/tags/tag_table_row_view", &fiber.Map{
		"Id":   tagData.Id,
		"Name": tagData.Name,
		"Slug": tagData.Slug,
	})
}

func AdminTagUpdateHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	fmt.Println("AdminTagUpdateHandler: ", id)

	tag := models.UpdateTagData{
		Name: c.FormValue("name"),
		Slug: c.FormValue("slug"),
	}

	updatedTag, err := models.Query.Tag.UpdateById(id, tag)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("fragments/admin/tags/tag_update_errors", &fiber.Map{
			"Errors": []string{err.Error()},
		})
	}

	return c.Render("fragments/admin/tags/tag_table_row_view", &fiber.Map{
		"Id":   updatedTag.Id,
		"Name": updatedTag.Name,
		"Slug": updatedTag.Slug,
	})
}

func AdminTagDeleteHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := models.Query.Tag.Delete(id); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).Send([]byte{})
}

func AdminTagPostHandler(c *fiber.Ctx) error {
	name := c.FormValue("name")
	slug := c.FormValue("slug")

	newTag := models.NewTagData{
		Name: name,
		Slug: slug,
	}

	createdTag, err := models.Query.Tag.Create(newTag)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Render("fragments/admin/shared/error_list", &fiber.Map{
			"Errors": []string{err.Error()},
		})
	}

	return c.Render("fragments/admin/tags/tag_table_row_view", &fiber.Map{
		"Id":   createdTag.Id,
		"Name": createdTag.Name,
		"Slug": createdTag.Slug,
	})
}
