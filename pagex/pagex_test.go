package pagex_test

import (
	"fmt"
	"github.com/odycenter/std-library/pagex"
	"testing"
)

func TestNew(t *testing.T) {
	page := pagex.New(2, 100, 20000)
	fmt.Println(page.Offset(), page.CurrPage)
}
