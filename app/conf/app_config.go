package app

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"os"
	"std-library/app/property"
	"std-library/logs"
	"sync"
)

var once sync.Once
var envOnce sync.Once
var Name = "" // TODO chris, refactor later
var env = ""

func SetName(name string) {
	if name == "" {
		return
	}
	once.Do(func() {
		Name = name // TODO refactor later
		logs.AppName = name
	})
}

type Config struct {
	propertyManager property.Manager
	validator       property.Validator
}

func NewConfig() *Config {
	return &Config{
		propertyManager: *property.NewManager(),
		validator:       *property.NewValidator(),
	}
}

func Local() bool {
	return Env() == ""
}

func Env() string {
	envOnce.Do(func() {
		var envStr = "env"
		var envVar = flag.String(envStr, "", "env: dev, cqa, prod")
		flag.Parse()
		env = *envVar

		envVarName := property.EnvVarName(envStr)
		envVarValue := os.Getenv(envVarName)
		if envVarValue != "" {
			logs.Warn("found overridden property by os.env var [%s], key=%s, value=%s", envVarName, envStr, envVarValue)
			env = envVarValue
		}
	})
	return env
}

func LoadProperties(envFS map[string]embed.FS, propertiesFileName string) *Config {
	env = Env()
	logs.Info("loadProperties by env: %s, propertiesFileName: %s", env, propertiesFileName)
	var conf = NewConfig()
	files, ok := envFS[env]
	if !ok {
		log.Fatal("loadProperties error! Invalid environment: " + env)
	}

	defaultFS, ok := envFS[""]
	if !ok {
		log.Fatal("loadProperties error! Default environment not found!")
	}
	conf.LoadPropertiesByFS(files, propertiesFileName, defaultFS)
	return conf
}

func (c *Config) LoadPropertiesByFS(properties embed.FS, propertyFile string, defaultFS embed.FS) {
	f, err := properties.Open(propertyFile)
	if err != nil {
		if os.IsNotExist(err) {
			logs.Warn("propertyFile not found!  load default properties, env: %v, fileName: %s ", Env(), propertyFile)
			f, err = defaultFS.Open(propertyFile)
		}

		if err != nil {
			log.Fatal("propertyFile not found!", propertyFile, err)
		}
	}
	defer f.Close()
	c.propertyManager.LoadProperties(f)
}

func (c *Config) LoadProperties(file fs.File) {
	c.propertyManager.LoadProperties(file)
}

func (c *Config) LoadPropertiesByPath(path string) {
	c.propertyManager.LoadPropertiesByPath(path)
}

func (c *Config) Property(key string) string {
	c.validator.Add(key)
	return c.propertyManager.Get(key)
}

func (c *Config) RequiredProperty(key string) string {
	c.validator.Add(key)
	return c.propertyManager.Get(key, true)
}

func (c *Config) Validate() {
	keys := c.propertyManager.Keys()
	c.validator.Validate(keys)
}
