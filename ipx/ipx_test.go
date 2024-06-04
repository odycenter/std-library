package ipx_test

import (
	"github.com/stretchr/testify/assert"
	"std-library/ipx"
	"testing"
)

func TestInfo(t *testing.T) {
	result, err := ipx.Info("52.139.172.204")
	assert.Nil(t, err)
	assert.Equal(t, "中国", result.Country)
	assert.Equal(t, "香港", result.Province)

	result, err = ipx.Info("127.0.0.1")
	assert.Nil(t, err)
	assert.Equal(t, "本机IP", result.Country)
	assert.Equal(t, "本机IP", result.Province)

	result, err = ipx.Info("2400:8902::f03c:91ff:febd:be52")
	assert.Nil(t, err)
	assert.Equal(t, "日本", result.Country)
	assert.Equal(t, "日本", result.Province)
}
