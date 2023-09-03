// Package snowflake 雪花唯一值算法
package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"log"
)

var sfNode *snowflake.Node

// New 创建一个雪花算法节点
// 如用于分布式系统唯一ID则node参数需全局唯一
// warning:请勿回拨时钟，否则可能重复
func New(node int64) (err error) {
	sfNode, err = snowflake.NewNode(node)
	if err != nil {
		log.Println("snowflake node init failed")
		return err
	}
	return nil
}

// Gen 生成雪花算法ID
func Gen() snowflake.ID {
	if sfNode == nil {
		return 0
	}
	return sfNode.Generate()
}
