package property_test

import (
	"github.com/odycenter/std-library/app/property"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator(t *testing.T) {
	validator := property.NewValidator()
	assert.Panics(t, func() { validator.Validate([]string{"abcKey1"}) })

	validator.Add("abcKey1")
	assert.NotPanics(t, func() { validator.Validate([]string{"abcKey1"}) })
}
