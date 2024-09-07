package snowflake_test

import (
	"fmt"
	"github.com/odycenter/std-library/unique/snowflake"
	"testing"
)

func TestSnowFlake(t *testing.T) {
	if err := snowflake.New(1); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(snowflake.Gen().String())
}
