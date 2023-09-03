package nets_test

import (
	"fmt"
	"testing"

	"github.com/odycenter/std-library/nets"
)

func TestNets(t *testing.T) {
	ips := nets.IpInt("127.0.0.1").Int64()
	fmt.Println(ips)
	fmt.Println(nets.IpStr(uint(ips)))
}
