// Package nets 网络类实现
package nets

import (
	"fmt"
	"github.com/mssola/useragent"
	"io"
	"net"
	"net/http"
	"std-library/json"
	"std-library/stringx"
	"strconv"
	"strings"
)

// T 返回值类型转换结构
type T struct {
	v string
}

func (t *T) String() string {
	return t.v
}

func (t *T) Int64() int64 {
	bits := strings.Split(t.v, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

// IpInt IP地址转数字
func IpInt(ip string) *T {
	return &T{ip}
}

// IpStr 数字IP转v4Ip
func IpStr(n uint) string {
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(n >> 24)
	ip[1] = byte(n >> 16)
	ip[2] = byte(n >> 8)
	ip[3] = byte(n)
	return ip.String()
}

// IP 获取IP地址
func IP(req *http.Request) *T {
	if len(req.Header.Get("x-forwarded-for")) > 0 {
		// 代理转发ip 格式：180.158.93.171,128.18.31.52
		ipList := strings.Split(req.Header.Get("x-forwarded-for"), ",")
		if len(ipList) > 0 {
			return &T{strings.TrimSpace(ipList[0])}
		}
		// 防止被突破
	} else if len(req.Header.Get("X-App-Real-IP")) > 0 {
		// 转发过程强制设置的一个变量
		return &T{req.Header.Get("X-App-Real-IP")}
	} else if len(req.Header.Get("X-Real-Ip")) > 0 {
		// 转发过程强制设置的一个变量
		return &T{req.Header.Get("X-Real-Ip")}
	}
	// "IP:port" "192.168.1.150:8889"
	return &T{strings.Split(req.RemoteAddr, ":")[0]}
}

// LocalIP 获取本机IP
func LocalIP() *T {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return &T{""}
	}
	for _, address := range addresses {
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return &T{ip.IP.String()}
			}
		}
	}
	return &T{""}
}

var IpUrl = "https://sp0.baidu.com/8aQDcjqpAAV3otqbppnN2DJv/api.php?query=%v&resource_id=6006"

// IPAddress 获取IP定位
func IPAddress(ip string) string {
	if ip == "127.0.0.1" || ip == "0.0.0.0" {
		return "本机IP"
	} else if strings.HasPrefix(ip, "192.168") || strings.HasPrefix(ip, "10.") {
		return "内网IP"
	} else if len(ip) == 0 {
		return "其他"
	}

	resp, err := http.Get(fmt.Sprintf(IpUrl, ip))
	if err != nil {
		return "其他"
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	address := json.Get(stringx.Encode(bytes, stringx.GB18030), "data.0.location").String()

	return address
}

// ParseUserAgent UserAgent
// 解析给定的 User-Agent 字符串并获取结果UserAgent对象。
// 返回一个在解析给定 User-Agent 字符串后已初始化的UserAgent对象。
func ParseUserAgent(userAgent string) (string, string, string, string, string) {
	ua := useragent.New(userAgent)
	info := ua.OSInfo()
	browserName, browserVersion := ua.Browser()
	return ua.Platform(), info.Name, info.Version, browserName, browserVersion
}
