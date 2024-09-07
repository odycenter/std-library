package app

import (
	"flag"
	"fmt"
	"github.com/odycenter/std-library/app/property"
	"github.com/odycenter/std-library/logs"
	"log/slog"
	"os"
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

func Local() bool {
	return Env() == ""
}

func LocalHostName() string {
	hostname, err := os.Hostname()
	if err == nil {
		return hostname
	}
	return "local"
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
			slog.Warn(fmt.Sprintf("found overridden property by os.env var [%s], key=%s, value=%s", envVarName, envStr, envVarValue))
			env = envVarValue
		}
	})
	return env
}
