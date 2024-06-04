package ipdb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInfo(t *testing.T) {
	driver := &Driver{}
	result, err := driver.Info("2400:8902::f03c:91ff:febd:be52")
	assert.Nil(t, err)
	assert.Equal(t, "中国", result.Country)
	assert.Equal(t, "山东", result.Province)

	result2, err := driver.Info("160.16.207.213")
	assert.Nil(t, err)
	assert.Equal(t, "日本", result2.Country)
	assert.Equal(t, "日本", result2.Province)

}
