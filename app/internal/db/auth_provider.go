package internal_db

type AuthProvider interface {
	AccessToken() string
}
