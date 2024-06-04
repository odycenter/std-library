package ipx

import (
	"net"
	"std-library/ipx/baidu"
	"std-library/ipx/ipdb"
	"std-library/ipx/ipdto"
)

var defaultDriver = ipdb.DriverKey
var privateIPBlocks = []*net.IPNet{
	{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
	{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)},
	{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
}
var drivers = map[string]IpGeolocation{
	baidu.DriverKey: &baidu.Driver{},
	ipdb.DriverKey:  &ipdb.Driver{},
}

type IpGeolocation interface {
	Info(string) (*ipdto.Info, error)
}

func Driver(key string) IpGeolocation {
	return drivers[key]
}

func Info(ipAddr string) (*ipdto.Info, error) {
	ip := net.ParseIP(ipAddr)
	var info = privateIPInfo(ip)
	if info != "" {
		return &ipdto.Info{
			Country:  info,
			Province: info,
			City:     info,
		}, nil
	}

	v4 := ip.To4()
	if v4 != nil {
		return Driver(defaultDriver).Info(v4.String())
	}

	return Driver(defaultDriver).Info(ipAddr)
}

func privateIPInfo(ip net.IP) string {
	if ip == nil {
		return "其他"
	}
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return "本机IP"
	}
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return "内网IP"
		}
	}
	return ""
}
