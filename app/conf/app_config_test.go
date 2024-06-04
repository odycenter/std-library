package app_test

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"std-library/app/conf"
	"testing"
)

//go:embed test2.properties
var dev embed.FS

//go:embed test.properties
var defaultFS embed.FS

func TestLoadProperties(t *testing.T) {
	envFS := map[string]embed.FS{"dev": dev, "": defaultFS}
	os.Setenv("ENV", "dev")
	conf := app.LoadProperties(envFS, "test.properties")
	assert.Equal(t, "name", conf.RequiredProperty("app"))

	conf = app.LoadProperties(envFS, "test2.properties")
	assert.Equal(t, "override", conf.RequiredProperty("app"))

}

func TestProperty(t *testing.T) {
	conf := app.NewConfig()
	conf.LoadPropertiesByPath("test.properties")
	assert.Equal(t, "name", conf.Property("app"))

	f, err := dev.Open("test2.properties")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	conf.LoadProperties(f)

	assert.Equal(t, "override", conf.Property("app"))

	assert.Panics(t, func() {
		conf.Validate()
	})

	conf.Property("abc1")
	assert.NotPanics(t, func() {
		conf.Validate()
	})
}
