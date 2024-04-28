package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func NewValidator() (*validator.Validate, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.RegisterValidation("slug", validateSlug); err != nil {
		return nil, err
	}

	return validate, nil
}

func validateSlug(slug validator.FieldLevel) bool {
	slugRegexString := "^[a-zA-Z0-9\\-]+$"
	slugRegex := regexp.MustCompile(slugRegexString)

	return slugRegex.MatchString(slug.Field().String())
}
