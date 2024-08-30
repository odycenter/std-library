package web_test

import (
	"github.com/stretchr/testify/assert"
	"std-library/app/web"
	"testing"
)

func TestParse(t *testing.T) {
	assert.Equal(t, "0.0.0.0:8080", web.Parse("8080").String())
	assert.Panics(t, func() {
		web.Parse(":8080")
	})
	assert.Equal(t, "123.123.123.123:8080", web.Parse("123.123.123.123:8080").String())
}
