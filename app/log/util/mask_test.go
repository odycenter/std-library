package util_test

import (
	"github.com/stretchr/testify/assert"
	"std-library/app/log/util"
	"testing"
)

func TestFilterWithoutMask(t *testing.T) {
	value := "{\"field1\": \"value1\"}"
	value = util.Filter(value)
	assert.Equal(t, value, "{\"field1\": \"value1\"}")
}

func TestFilterOneMaskField(t *testing.T) {
	value := "{\"field1\": \"value1\",\n  \"password\" : \"pass123\",\n  \"field2\": \"value2\"\n}"
	actual := util.Filter(value, "password")
	assert.Equal(t, "{\"field1\": \"value1\",\n  \"password\" : \"******\",\n  \"field2\": \"value2\"\n}", actual)

	value = "{\"field1\": \"value1\", \"password\": null, \"field2\": null, \"field3\": null}"
	assert.Equal(t, value, util.Filter(value, "password", "passwordConfirm"))

	value = "{\"field1\": \"value1\", \"password\": {\"field2\": null}}"
	assert.Equal(t, value, util.Filter(value, "password", "passwordConfirm"))

	value = "{\"field1\": [\"secret1\", \"secret2\", \"secret3\"], \"field2\": \"value\"}"
	assert.Equal(t, value, util.Filter(value, "secret1", "secret2", "secret3"))
}

func TestFilterJSONWithNull(t *testing.T) {
	value := "{\"field1\": \"value1\", \"password\": null, \"field2\": \"value2\"}"
	actual := util.Filter(value, "password", "passwordConfirm")
	assert.Equal(t, "{\"field1\": \"value1\", \"password\": null, \"field2\": \"value2\"}", actual)
}

func TestFilterMultipleMaskFields(t *testing.T) {
	value := "{\"field1\": \"value1\", \"password\": \"pass123\", \"passwordConfirm\": \"pass123\", \"field2\": \"value2\", \"nested\": {\"password\": \"pass\\\"123\", \"passwordConfirm\": \"pass123\"}}"
	actual := util.Filter(value, "password", "passwordConfirm")
	assert.Equal(t, "{\"field1\": \"value1\", \"password\": \"******\", \"passwordConfirm\": \"******\", \"field2\": \"value2\", \"nested\": {\"password\": \"******\", \"passwordConfirm\": \"******\"}}", actual)
}

func TestFilterBrokenJSON(t *testing.T) {
	value := "{\"field1\": \"value1\",\n  \"password\": \"pass123\",\n  \"passwordConfirm\": \"pass12"
	value = util.Filter(value, "password", "passwordConfirm")
	assert.NotContains(t, value, "pass123")

	value = "{\"field1\": \"value1\",\n  \"password\": \"pass123\",\n  \"passwordConfirm\""
	value = util.Filter(value, "password", "passwordConfirm")
	assert.NotContains(t, value, "pass123")

	value = "{\"field1\": \"value1\", \"password\": \"pass123"
	value = util.Filter(value, "password", "passwordConfirm")
	assert.Equal(t, value, "{\"field1\": \"value1\", \"password\": \"******")
}
