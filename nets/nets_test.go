package nets_test

import (
	"fmt"
	"github.com/odycenter/std-library/nets"
	"testing"
)

func TestNets(t *testing.T) {
	ips := nets.IpInt("127.0.0.1").Int64()
	fmt.Println(ips)
	fmt.Println(nets.IpStr(uint(ips)))
}
