package pkg

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestValidators(t *testing.T) {
	scenarios := map[string]func(t *testing.T, v *validator.Validate){
		"test validate slug": testValidatorSlug,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			v, err := NewValidator()
			if err != nil {
				t.Fatal("failed to create validator")
			}

			fn(t, v)
		})
	}
}

func testValidatorSlug(t *testing.T, v *validator.Validate) {
	slugWithSpaces := "test slug spaces"
	slugWithSpecials := "test+$%slug()spec!@l$"
	validSlug := "valid-test-slug"

	var err error

	err = v.Var(slugWithSpaces, "slug")
	require.Error(t, err, "should not be valid")

	err = v.Var(slugWithSpecials, "slug")
	require.Error(t, err, "should not be valid")

	err = v.Var(validSlug, "slug")
	require.Nil(t, err, "should not return error")
}
