package property_test

import (
	"github.com/odycenter/std-library/app/property"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOverrideEnvVarName(t *testing.T) {
	assert.Equal(t, "ABCKEY1", property.EnvVarName("abcKey1"))
	assert.Equal(t, "ABC_KEY1", property.EnvVarName("abc.Key1"))
	assert.Equal(t, "ABC_DEF_KEY1", property.EnvVarName("abc.def.Key1"))
	assert.Equal(t, "ABC_DEF_KEY1", property.EnvVarName("abc.def.key1"))
	assert.Equal(t, "ENV", property.EnvVarName("env"))
	assert.Equal(t, "MICRO_CLOUD_WITHDRAW", property.EnvVarName("micro_cloud_withdraw"))
	assert.Equal(t, "CLOUD-MICRO-DOWNLOAD", property.EnvVarName("cloud-micro-download"))

}
