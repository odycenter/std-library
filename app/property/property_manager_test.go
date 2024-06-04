package property_test

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"std-library/app/property"
	"testing"
)

func TestProperty(t *testing.T) {
	manager := property.NewManager()
	file, err := os.Open("test.properties")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	manager.LoadProperties(file)
	assert.Equal(t, "123", manager.Get("abcKey1"))
	assert.Equal(t, "123=456", manager.Get("key3"))
	assert.Equal(t, "", manager.Get("AcsPassword1"))

	err = os.Setenv("ABCKEY1", "new123")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "new123", manager.Get("abcKey1"))

	assert.Equal(t, []string{"abcKey1", "app.key2", "key3"}, manager.Keys())
}
