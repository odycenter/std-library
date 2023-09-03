package mongo

import "go.mongodb.org/mongo-driver/mongo"

// IsDuplicateKeyError 是否为主键重复Error
func IsDuplicateKeyError(err error) bool {
	return mongo.IsDuplicateKeyError(err)
}

// IsTimeout 是否为超时Error
func IsTimeout(err error) bool {
	return mongo.IsTimeout(err)
}

// IsNetworkError 是否为网络Error
func IsNetworkError(err error) bool {
	return mongo.IsNetworkError(err)
}
