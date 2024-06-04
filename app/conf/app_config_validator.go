package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"std-library/app/property"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type EnvResourceAssert struct {
	confPath      string
	mainResources string
}

func NewEnvResourceAssert() *EnvResourceAssert {
	return &EnvResourceAssert{
		confPath:      filepath.Join("conf"),
		mainResources: filepath.Join("resources"),
	}
}

func (e *EnvResourceAssert) OverridesDefaultResources(t *testing.T) {
	assert := assert.New(t)

	assert.DirExists(e.confPath)

	resourceDirs, err := e.resourceDirs()
	if err != nil {
		t.Fatal(err)
	}

	for _, resourceDir := range resourceDirs {
		assert.DirExists(resourceDir)
		e.AssertOverridesDefault(resourceDir, assert)
		e.AssertPropertyOverridesDefault(resourceDir, assert)
	}

}

func (e *EnvResourceAssert) resourceDirs() ([]string, error) {
	files, err := os.ReadDir(e.confPath)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, filepath.Join(e.confPath, f.Name(), "resources"))
		}
	}

	return dirs, nil
}

func (e *EnvResourceAssert) AssertOverridesDefault(resourceDir string, assert *assert.Assertions) {
	err := filepath.Walk(resourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		base := filepath.Base(path)
		if base == filepath.Base(resourceDir) {
			return nil
		}
		defaultFile := filepath.Join(e.mainResources, base)
		assert.DirExists(defaultFile)

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (e *EnvResourceAssert) AssertPropertyOverridesDefault(resourceDir string, assert *assert.Assertions) {
	err := filepath.Walk(resourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".properties" {
			return nil
		}

		defaultPropertyFile := filepath.Join(e.mainResources, filepath.Base(path))
		if _, err := os.Stat(defaultPropertyFile); os.IsNotExist(err) {
			return err
		}

		defaultProperties := loadProperties(defaultPropertyFile)
		envProperties := loadProperties(path)
		keys := defaultProperties.Keys()
		envKeys := envProperties.Keys()
		validKeyOrder(keys, defaultPropertyFile)
		validKeyOrder(envKeys, path)

		assert.Equal(keys, envProperties.Keys(), "%v must override %v", path, defaultPropertyFile)

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func validKeyOrder(keys []string, propertyFile string) {
	for i, key := range keys {
		if i == 0 {
			continue
		}
		if strings.Compare(key, keys[i-1]) < 0 {
			panic(errors.New(fmt.Sprintf("property key '%s' is not in the right order with previous property '%s', %v", key, keys[i-1], propertyFile)))
		}
	}
}

func loadProperties(path string) *property.Manager {
	manager := property.NewManager()
	manager.LoadPropertiesByPath(path)
	return manager
}
