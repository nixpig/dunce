package tags

import (
	"errors"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type TagService struct {
	data TagDataInterface
}

type TagServiceInterface interface {
	create(tag *Tag) (*Tag, error)
	deleteById(id int) error
	getAll() (*[]Tag, error)
	getBySlug(slug string) (*Tag, error)
	update(tag *Tag) (*Tag, error)
}

func NewTagService(data TagDataInterface) TagService {
	return TagService{data}
}

func (ts TagService) create(tag *Tag) (*Tag, error) {
	// TODO: maybe inject validator at point of struct initialisation?
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.RegisterValidation("slug", ValidateSlug); err != nil {
		return nil, err
	}

	if err := validate.Struct(tag); err != nil {
		return nil, err
	}

	exists, err := ts.data.exists(tag)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("tag name and/or slug already exists")
	}

	createdTag, err := ts.data.create(tag)
	if err != nil {
		return nil, err
	}

	return createdTag, nil
}

func (ts TagService) deleteById(id int) error {
	if err := ts.data.deleteById(id); err != nil {
		return err
	}

	return nil
}

func (ts TagService) getAll() (*[]Tag, error) {
	tags, err := ts.data.getAll()
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (ts TagService) getBySlug(slug string) (*Tag, error) {
	tag, err := ts.data.getBySlug(slug)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (ts TagService) update(tag *Tag) (*Tag, error) {
	// TODO: maybe inject validator at point of struct initialisation?
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.RegisterValidation("slug", ValidateSlug); err != nil {
		return nil, err
	}

	if err := validate.Struct(tag); err != nil {
		return nil, err
	}

	exists, err := ts.data.exists(tag)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("tag name and/or slug already exists")
	}

	updatedTag, err := ts.data.update(tag)
	if err != nil {
		return nil, err
	}

	return updatedTag, nil
}

func ValidateSlug(slug validator.FieldLevel) bool {
	slugRegexString := "^[a-zA-Z0-9\\-]+$"
	slugRegex := regexp.MustCompile(slugRegexString)

	return slugRegex.MatchString(slug.Field().String())
}
