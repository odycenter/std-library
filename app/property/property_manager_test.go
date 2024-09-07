package property_test

import (
	"github.com/odycenter/std-library/app/property"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
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

	assert.Equal(t, "123=456k4,kk4", manager.Get("key4"))
	assert.Equal(t, "key5,//k5,#kk5", manager.Get("key5"))
	assert.Equal(t, "123", manager.Get("key6"))

	assert.Equal(t, []string{"abcKey1", "app.key2", "key3", "key4", "key5", "key6"}, manager.Keys())
}
