package snowflake_test

import (
	"fmt"
	"testing"

	"github.com/odycenter/std-library/unique/snowflake"
)

func TestSnowFlake(t *testing.T) {
	if err := snowflake.New(1); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(snowflake.Gen().String())
}
