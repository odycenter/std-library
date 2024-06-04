package cloud

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type cachedToken struct {
	Token     string
	expiresAt time.Time
}

type AWSAuthProvider struct {
	Region      string
	DBEndpoint  string
	DBName      string
	cachedToken *cachedToken
}

var once sync.Once

func (a *AWSAuthProvider) Register() {
	once.Do(func() {
		RegisterRDSMysqlCerts()
	})
}

func (a *AWSAuthProvider) DataSourceName() string {
	return a.dataSourceName("emptyToken")
}

func (a *AWSAuthProvider) dataSourceName(token string) string {
	log.Println("[AWSAuthProvider] DataSourceName, addr: ", a.DBEndpoint, ", DBName: ", a.DBName)
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?allowCleartextPasswords=true&tls=rds&charset=utf8mb4", a.DBUser(), token, a.DBEndpoint, a.DBName)
}

func (a *AWSAuthProvider) AccessToken() string {
	start := time.Now()
	if a.cachedToken != nil && time.Now().Before(a.cachedToken.expiresAt) {
		log.Println("[AWSAuthProvider] get aws access token from cache, elapse: ", time.Since(start))
		return a.cachedToken.Token
	}

	if a.Region == "" {
		region := os.Getenv("RDS_REGION")
		if region != "" {
			log.Println("[AWSAuthProvider] region is empty, use ENV[RDS_REGION] region: ", region)
			a.Region = region
		} else {
			log.Println("[AWSAuthProvider] region is empty, use default region ap-northeast-1")
			a.Region = "ap-northeast-1"
		}
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(a.Region))
	if err != nil {
		panic("configuration error: " + err.Error())
	}
	authenticationToken, err := auth.BuildAuthToken(context.TODO(), a.DBEndpoint, a.Region, a.DBUser(), cfg.Credentials)
	if err != nil {
		panic("failed to create authentication token: " + err.Error())
	}
	elapse := time.Since(start)
	log.Println("[AWSAuthProvider] get aws access token elapse: ", elapse)

	a.cachedToken = &cachedToken{
		Token:     authenticationToken,
		expiresAt: time.Now().Add(10 * time.Minute), // aws iam token expires in 15 minutes
	}

	return authenticationToken
}

func (a *AWSAuthProvider) DBUser() string {
	return "cloud_iam"
}

func RegisterRDSMysqlCerts() {
	resp, err := http.DefaultClient.Get("https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem")
	if err != nil {
		panic("failed to download RDS certificate: " + err.Error())
	}

	pem, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("failed to read RDS certificate: " + err.Error())
	}

	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		panic("failed to append RDS certificate to root cert pool")
	}
	err = mysql.RegisterTLSConfig("rds", &tls.Config{RootCAs: rootCertPool, InsecureSkipVerify: true})
	if err != nil {
		panic("failed to register RDS certificate: " + err.Error())
	}
}
