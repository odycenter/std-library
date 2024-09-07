package jwt_test

import (
	"github.com/odycenter/std-library/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	tk := jwt.Gen(map[string]any{
		"Browser":          "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
		"CreateTime":       1685501436,
		"HId":              0,
		"Id":               1113,
		"Ip":               "54.255.114.140",
		"UseNewPermission": true,
		"exp":              1717037436,
	}, "TokenGenSalt", time.Minute*1440*365*10)
	result := jwt.Parse(tk, "TokenGenSalt")

	assert.Equal(t, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36", result["Browser"])
	assert.Equal(t, float64(1685501436), result["CreateTime"])
	assert.Equal(t, float64(0), result["HId"])
	assert.Equal(t, float64(1113), result["Id"])
	assert.Equal(t, "54.255.114.140", result["Ip"])
	assert.Equal(t, true, result["UseNewPermission"])

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJCcm93c2VyIjoiTW96aWxsYS81LjAgKFdpbmRvd3MgTlQgMTAuMDsgV2luNjQ7IHg2NCkgQXBwbGVXZWJLaXQvNTM3LjM2IChLSFRNTCwgbGlrZSBHZWNrbykgQ2hyb21lLzExMi4wLjAuMCBTYWZhcmkvNTM3LjM2IiwiQ3JlYXRlVGltZSI6MTY4NTUwMTQzNiwiSElkIjowLCJJZCI6MTExMywiSXAiOiI1NC4yNTUuMTE0LjE0MCIsIlVzZU5ld1Blcm1pc3Npb24iOnRydWUsImV4cCI6MjAzNTE4NDQ0MH0.VTO-Zj7TI48vGa_aHThJOr8ERA3CKcb0ilpAwdA7gU0"
	result2 := jwt.Parse(token, "TokenGenSalt")
	assert.Equal(t, result2["Browser"], result["Browser"])
	assert.Equal(t, result2["CreateTime"], result["CreateTime"])
	assert.Equal(t, result2["HId"], result["HId"])
	assert.Equal(t, result2["Id"], result["Id"])
	assert.Equal(t, result2["Ip"], result["Ip"])
	assert.Equal(t, result2["UseNewPermission"], result["UseNewPermission"])
}

func TestParseExpiredToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJCcm93c2VyIjoiTW96aWxsYS81LjAgKFdpbmRvd3MgTlQgMTAuMDsgV2luNjQ7IHg2NCkgQXBwbGVXZWJLaXQvNTM3LjM2IChLSFRNTCwgbGlrZSBHZWNrbykgQ2hyb21lLzExMi4wLjAuMCBTYWZhcmkvNTM3LjM2IiwiQ3JlYXRlVGltZSI6MTY4NTUwMTQzNiwiSElkIjowLCJJZCI6MTExMywiSXAiOiI1NC4yNTUuMTE0LjE0MCIsIlVzZU5ld1Blcm1pc3Npb24iOnRydWUsImV4cCI6MTcxOTgyNDAzM30.IatFWSPcIMKD5PJZFNynzZShEhtsb__myFb1_LtJihE"
	var claims map[string]any
	result := jwt.Parse(token, "TokenGenSalt")
	assert.Equal(t, claims, result)
}
