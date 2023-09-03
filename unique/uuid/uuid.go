// Package uuid UUID唯一值算法
package uuid

import "github.com/nacos-group/nacos-sdk-go/v2/inner/uuid"

// Gen 生成UUID
func Gen() uuid.UUID {
	u, _ := uuid.NewV1()
	return u
}
