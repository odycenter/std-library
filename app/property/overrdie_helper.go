package property

import (
	"os"
	"std-library/logs"
	"strings"
)

type OverrideHelper struct {
	App string
}

func (o *OverrideHelper) Get(key string) string {
	if o.App == "" {
		return ""
	}

	envVarName := o.EnvVarName(key)
	envVarValue := os.Getenv(envVarName)
	if envVarValue != "" {
		logs.Warn("found local overridden property by env var %s, key=%s, value=%s", envVarName, key, MaskValue(key, envVarValue))
		return envVarValue
	}
	return ""
}

func (o *OverrideHelper) EnvVarName(key string) string {
	envVarName := EnvVarName(o.App) + "_" + EnvVarName(key)
	return strings.ReplaceAll(envVarName, "-", "_")
}
