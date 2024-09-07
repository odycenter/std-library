package util_test

import (
	"context"
	"github.com/odycenter/std-library/app/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScheduleParse(t *testing.T) {
	emptyFunc := func(ctx context.Context) {}
	s := util.NewScheduler()
	s.PanicOnAnyAddError(true)
	_, err := s.AddFunc("0 0 3,9,15,21 * * ?", emptyFunc)
	assert.Nil(t, err)

	_, err = s.AddFunc("*/10 * * * * ?", emptyFunc)
	assert.Nil(t, err)
	s.Start()
	s.JobsInfo()
}
