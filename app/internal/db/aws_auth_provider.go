package internal_db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/go-sql-driver/mysql"
	"io"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type cachedToken struct {
	Token     string
	expiresAt time.Time
}

type AWSAuthProvider struct {
	User        string
	DBEndpoint  string
	Region      string
	cachedToken *cachedToken
}

var certsOnce sync.Once

func RegisterCerts() {
	certsOnce.Do(func() {
		RegisterRDSMysqlCerts()
	})
}

func (a *AWSAuthProvider) AccessToken() string {
	start := time.Now()
	if a.cachedToken != nil && time.Now().Before(a.cachedToken.expiresAt) {
		slog.Info("[AWSAuthProvider] get aws access token from cache", "elapse", time.Since(start))
		return a.cachedToken.Token
	}

	if a.Region == "" {
		log.Fatalf("DB region can not be empty, endpoint=%s, user=%s", a.DBEndpoint, a.User)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(a.Region))
	if err != nil {
		log.Panic("configuration error", err)
	}
	authenticationToken, err := auth.BuildAuthToken(context.TODO(), a.DBEndpoint, a.Region, a.User, cfg.Credentials)
	if err != nil {
		log.Panic("failed to create authentication token: " + err.Error())
	}
	elapse := time.Since(start)
	slog.Info("[AWSAuthProvider] get aws access token", "elapse", elapse)

	a.cachedToken = &cachedToken{
		Token:     authenticationToken,
		expiresAt: time.Now().Add(10 * time.Minute), // aws iam token expires in 15 minutes
	}

	return authenticationToken
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
