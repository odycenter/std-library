package baidu

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"net/http"
)

type IBaiduApi interface {
	Data(string) (interface{}, error)
}

var ApiUrl = "https://sp0.baidu.com/8aQDcjqpAAV3otqbppnN2DJv/api.php"

type api struct{}

func (a *api) Data(ip string) (interface{}, error) {
	url := fmt.Sprintf("%s?query=%v&resource_id=6006", ApiUrl, ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	decodeBytes, err := simplifiedchinese.GB18030.NewDecoder().Bytes(bytes)
	if err != nil {
		return nil, err
	}
	return decodeBytes, nil
}

type ApiResponse struct {
	Status string     `json:"status"`
	Data   []*ApiData `json:"data"`
}

type ApiData struct {
	Location string `json:"location"`
}
