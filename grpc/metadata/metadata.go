package metadata

import "google.golang.org/grpc/metadata"

// Pairs
// 组织一个由键、值映射形成的MD ...如果 len(kv) 为奇数，则 Pairs 会出现panic。
// 键中只允许使用以下 ASCII 字符：
// 数字：0-9
// 大写字母：AZ（标准化为小写）
// 小写字母：az
// 特殊字符： -_.
// 大写字母会自动转换为小写字母。
// 以“grpc-”开头的密钥仅供 grpc 内部使用，如果在元数据中设置，可能会导致错误。
func Pairs(kv ...string) metadata.MD {
	return metadata.Pairs(kv...)
}

// Map
// 根据给定的键值映射创建MD 。
// 键中只允许使用以下 ASCII 字符：
// 数字：0-9
// 大写字母：AZ（标准化为小写）
// 小写字母：az
// 特殊字符： -_.
// 大写字母会自动转换为小写字母。
// 以“grpc-”开头的密钥仅供 grpc 内部使用，如果在元数据中设置，可能会导致错误。
func Map(m map[string]string) metadata.MD {
	return metadata.New(m)
}
