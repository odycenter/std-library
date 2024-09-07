package kafka_test

import (
	"github.com/odycenter/std-library/app/kafka"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUriParse(t *testing.T) {
	uri := kafka.Uri{}
	uri.Uri("localhost:9092,localhost2,localhost3:9099")
	uris := uri.Parse()
	assert.Equal(t, []string{"localhost:9092", "localhost2:9092", "localhost3:9099"}, uris)
}
