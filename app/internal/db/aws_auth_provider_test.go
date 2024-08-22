package internal_db

import (
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAWSAuthProvider_DataSourceNameByToken(t *testing.T) {
	c, e := mysql.ParseDSN("cloud_ap_use:@tcp(privatelink-18130664.xcfjzr1bqus1.clusters.tidb-cloud.com:4000)/cloud")
	c.Passwd = "token"
	c.AllowCleartextPasswords = true
	assert.Nil(t, e)
	assert.Equal(t, "tcp", c.Net)
	assert.Equal(t, "cloud_ap_use", c.User)
	assert.Equal(t, "privatelink-18130664.xcfjzr1bqus1.clusters.tidb-cloud.com:4000", c.Addr)
	assert.Equal(t, "cloud", c.DBName)
	assert.Equal(t, "token", c.Passwd)
	println(c.FormatDSN())

	//cloud_iam:token@tcp(localhost:3306)/cloud?allowCleartextPasswords=true&tls=rds&charset=utf8mb4
	//cloud_iam:token@tcp(localhost:3306)/cloud?allowCleartextPasswords=true&tls=rds&charset=utf8mb4
}
