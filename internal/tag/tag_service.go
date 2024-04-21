package tag

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type TagService struct {
	data TagDataInterface
}

func NewTagService(data TagDataInterface) TagService {
	return TagService{data}
}

func (ts *TagService) Create(tag *Tag) (*Tag, error) {
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

	fmt.Println(" >>> does tag exist?", exists)

	if exists {
		return nil, errors.New("tag name and/or slug already exists")
	}

	createdTag, err := ts.data.create(tag)
	if err != nil {
		return nil, err
	}

	return createdTag, nil
}

func (ts *TagService) DeleteById(id int) error {
	if err := ts.data.deleteById(id); err != nil {
		return err
	}

	return nil
}

func ValidateSlug(slug validator.FieldLevel) bool {
	slugRegexString := "^[a-zA-Z0-9\\-]+$"
	slugRegex := regexp.MustCompile(slugRegexString)

	return slugRegex.MatchString(slug.Field().String())
}
