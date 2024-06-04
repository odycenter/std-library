package baidu_test

import (
	"github.com/stretchr/testify/assert"
	"std-library/ipx/baidu"
	"testing"
)

func TestInfo(t *testing.T) {
	stub := &apiTestStub{}
	driver := &baidu.Driver{Api: stub}
	result, err := driver.Info("any")
	assert.Nil(t, err)
	assert.Equal(t, "台灣省", result.Country)
	assert.Equal(t, "台北市", result.City)
}

type apiTestStub struct{}

func (s *apiTestStub) Data(ip string) (interface{}, error) {
	return []byte("{\"status\":\"0\",\"data\":[{\"location\":\"台灣省台北市\"}]}"), nil
}
