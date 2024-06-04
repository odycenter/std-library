package cloud

import (
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAWSAuthProvider_DataSourceNameByToken(t *testing.T) {
	a := &AWSAuthProvider{
		DBEndpoint: "localhost:3306",
		DBName:     "cloud",
		Region:     "ap-northeast-1",
	}
	actual := a.dataSourceName("token")
	excepted := "cloud_iam:token@tcp(localhost:3306)/cloud?allowCleartextPasswords=true&tls=rds&charset=utf8mb4"
	assert.Equal(t, excepted, actual)

	a.Register()
	cfg, err := mysql.ParseDSN(actual)
	assert.Nil(t, err)
	actual = cfg.FormatDSN()
	assert.Equal(t, excepted, actual)
}
