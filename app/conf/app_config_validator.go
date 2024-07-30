package app

import (
	"fmt"
	"github.com/stretchr/testify/require"
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

func (era *EnvResourceAssert) OverridesDefaultResources(t *testing.T) {
	require.DirExists(t, era.mainResources)
	require.DirExists(t, era.confPath)

	resourceDirs, err := era.resourceDirs()
	require.NoError(t, err)

	era.validateMainResourcesProperties(t)

	for _, resourceDir := range resourceDirs {
		assert.DirExists(t, resourceDir)
		era.AssertOverridesDefault(t, resourceDir)
		era.AssertPropertyOverridesDefault(t, resourceDir)
	}
}

func (era *EnvResourceAssert) resourceDirs() ([]string, error) {
	files, err := os.ReadDir(era.confPath)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, filepath.Join(era.confPath, f.Name(), "resources"))
		}
	}

	return dirs, nil
}

func (era *EnvResourceAssert) AssertOverridesDefault(t *testing.T, resourceDir string) {
	err := filepath.Walk(resourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() || path == resourceDir {
			return nil
		}

		assert.FileExists(t, filepath.Join(era.mainResources, filepath.Base(path)))
		return nil
	})

	assert.NoError(t, err)
}

func (era *EnvResourceAssert) AssertPropertyOverridesDefault(t *testing.T, resourceDir string) {
	err := filepath.Walk(resourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".properties" {
			return nil
		}

		return era.validatePropertyFile(t, path)
	})

	assert.NoError(t, err)
}

func (era *EnvResourceAssert) validatePropertyFile(t *testing.T, path string) error {
	defaultPropertyFile := filepath.Join(era.mainResources, filepath.Base(path))
	if _, err := os.Stat(defaultPropertyFile); os.IsNotExist(err) {
		return fmt.Errorf("default property file does not exist: %s", defaultPropertyFile)
	}

	defaultProperties := loadProperties(defaultPropertyFile)
	envProperties := loadProperties(path)

	keys := defaultProperties.Keys()
	envKeys := envProperties.Keys()

	era.validateKeyOrder(t, envKeys, path)

	assert.Equal(t, keys, envKeys, "%v must override %v", path, defaultPropertyFile)
	return nil
}

func (era *EnvResourceAssert) validateMainResourcesProperties(t *testing.T) {
	err := filepath.Walk(era.mainResources, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".properties" {
			properties := loadProperties(path)
			era.validateKeyOrder(t, properties.Keys(), path)
		}
		return nil
	})
	assert.NoError(t, err)
}

func (era *EnvResourceAssert) validateKeyOrder(t *testing.T, keys []string, propertyFile string) {
	for i := 1; i < len(keys); i++ {
		assert.True(t, strings.Compare(keys[i-1], keys[i]) < 0,
			"Property key '%s' is not in the right order with previous property '%s' in %v",
			keys[i], keys[i-1], propertyFile)
	}
}

func loadProperties(path string) *property.Manager {
	manager := property.NewManager()
	manager.LoadPropertiesByPath(path)
	return manager
}
