package property_test

import (
	"github.com/stretchr/testify/assert"
	"std-library/app/property"
	"testing"
)

func TestEnvVarName(t *testing.T) {
	helper := property.OverrideHelper{App: "cloud-app"}
	assert.Equal(t, "CLOUD_APP_VAR1_DOWNLOAD", helper.EnvVarName("var1-download"))
	assert.Equal(t, "CLOUD_APP_VAR2_DOWNLOAD", helper.EnvVarName("var2_download"))

}
