package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/nixpig/dunce/internal/pkg/models"
)

func AdminTypeGetHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")

	if len(idParam) == 0 {
		types, err := models.Query.Type.GetAll()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		pathParts := strings.Split(c.Path(), "/")

		page := pathParts[len(pathParts)-1]

		return c.Render("pages/admin/types", &fiber.Map{
			"Page":  page,
			"Types": types,
		}, "layouts/admin")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	typeData, err := models.Query.Type.GetById(id)
	if err != nil {
		fmt.Println("handling: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	editable := c.Query("edit")

	if editable == "true" {
		return c.Render("fragments/admin/types/type_table_row_edit", &fiber.Map{
			"Id":       typeData.Id,
			"Name":     typeData.Name,
			"Template": typeData.Template,
			"Slug":     typeData.Slug,
		})
	}

	return c.Render("fragments/admin/types/type_table_row_view", &fiber.Map{
		"Id":       typeData.Id,
		"Name":     typeData.Name,
		"Template": typeData.Template,
		"Slug":     typeData.Slug,
	})
}

func AdminTypePostHander(c *fiber.Ctx) error {
	name := c.FormValue("name")
	template := c.FormValue("template")
	slug := c.FormValue("slug")

	newType := models.TypeData{
		Name:     name,
		Template: template,
		Slug:     slug,
	}

	createdType, err := models.Query.Type.Create(newType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("fragments/admin/shared/error_list", &fiber.Map{
			"Errors": []string{err.Error()},
		})
	}

	return c.Render("fragments/admin/types/type_table_row_view", &fiber.Map{
		"Id":       createdType.Id,
		"Name":     createdType.Name,
		"Template": createdType.Template,
		"Slug":     createdType.Slug,
	})
}

func AdminTypeDeleteHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := models.Query.Type.DeleteById(id); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).Send([]byte{})
}

func AdminTypePutHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	typeData := models.TypeData{
		Name:     c.FormValue("name"),
		Template: c.FormValue("template"),
		Slug:     c.FormValue("slug"),
	}

	updatedType, err := models.Query.Type.UpdateById(id, typeData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("fragments/admin/shared/admin_table_errors", &fiber.Map{
			"Errors": []string{err.Error()},
		})
	}

	return c.Status(fiber.StatusOK).Render("fragments/admin/types/type_table_row_view", &fiber.Map{
		"Id":       updatedType.Id,
		"Name":     updatedType.Name,
		"Template": updatedType.Template,
		"Slug":     updatedType.Slug,
	})
}
