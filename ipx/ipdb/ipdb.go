package ipdb

import (
	"embed"
	"github.com/ipipdotnet/ipdb-go"
	"std-library/ipx/ipdto"
	"std-library/logs"
	"strings"
)

const DriverKey = "ipdb"

type Driver struct{}

//go:embed ipv4_china.ipdb
var ipdbIpv4 embed.FS
var dbIpv4 *ipdb.City

//go:embed ipv6_china_cn.ipdb
var ipdbIpv6 embed.FS
var dbIpv6 *ipdb.City

func init() {
	fileIpv4, err := ipdbIpv4.ReadFile("ipv4_china.ipdb")
	if err != nil {
		panic(err)
	}
	dbIpv4, err = ipdb.NewCityFromBytes(fileIpv4)
	if err != nil {
		panic(err)
	}
	fileIpv6, err := ipdbIpv6.ReadFile("ipv6_china_cn.ipdb")
	if err != nil {
		panic(err)
	}
	dbIpv6, err = ipdb.NewCityFromBytes(fileIpv6)
	if err != nil {
		panic(err)
	}
}

func (driver *Driver) Info(ip string) (info *ipdto.Info, err error) {
	var db *ipdb.City
	if strings.ContainsAny(ip, ".") {
		db = dbIpv4
	} else if strings.ContainsAny(ip, ":") {
		db = dbIpv6
	}

	record, err := db.FindInfo(ip, "CN")
	if err != nil {
		logs.Error("[Info] get ip %s info err", ip, err)
		return nil, err
	}

	return &ipdto.Info{
		Isp:      record.IspDomain,
		Country:  record.CountryName,
		Province: record.RegionName,
		City:     record.CityName,
	}, nil
}
