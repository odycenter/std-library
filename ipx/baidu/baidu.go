package baidu

import (
	"encoding/json"
	"errors"
	"std-library/ipx/ipdto"
	"std-library/stringx"
)

const DriverKey = "baidu"

type Driver struct {
	Api IBaiduApi
}

func (driver *Driver) Info(ip string) (*ipdto.Info, error) {
	if driver.Api == nil {
		driver.Api = &api{}
	}
	data, err := driver.Api.Data(ip)
	if err != nil {
		return nil, err
	}

	bytes := data.([]byte)
	resp := new(ApiResponse)
	err = json.Unmarshal(bytes, resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, errors.New("query fail")
	}

	info := &ipdto.Info{}
	location := resp.Data[0].Location
	arr := stringx.Split(location, " ").Strings()

	size := len(arr)
	if size > 1 {
		info.Isp = arr[1]
	}

	arr1 := stringx.Split(arr[0], "省").Strings()
	if len(arr1) > 1 {
		info.Province = arr1[0] + "省"
		info.City = arr1[1]
	} else {
		info.Province = arr[0]
		info.City = arr[0]
	}
	info.Country = info.Province
	return info, nil
}
