// Package redis redis操作封装
package redis

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"net"
	"sync"
	"time"
)

var pool sync.Map //map[AliasName]*Cli
var latestCli *Cli
var defaultCtx = context.TODO()
var Nil = redis.Nil
var ErrCmdExec = errors.New("redis exec hsetex failed")

// Opt redis配置信息
type Opt struct {
	AliasName    string
	IsCluster    bool     //default:false. Addrs[0] will only be used
	Addrs        []string //["localhost:6379",...]string{}
	Network      string   //default:tcp
	Username     string   //>redis 6.0
	Password     string   //
	DB           int      //when standalone effected,cluster only use DB 0
	PoolSize     int      //Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	MinIdleConns int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	TLSConfig    *tls.Config                                                              `json:"-"`
	OnConnect    func(ctx context.Context, cn *redis.Conn) error                          `json:"-"`
	Dialer       func(ctx context.Context, network string, addr string) (net.Conn, error) `json:"-"`
}

func (opt *Opt) getAddrs() []string {
	if len(opt.Addrs) == 0 {
		log.Panicln("Redis Addr is Empty")
	}
	if !opt.IsCluster {
		return []string{opt.Addrs[0]}
	}
	return opt.Addrs
}

func (opt *Opt) getMinIdleConns() int {
	if opt.PoolSize == 0 {
		return 10
	}
	return opt.PoolSize
}

func (opt *Opt) getAliasName() string {
	if opt.AliasName == "" {
		return "default"
	}
	return opt.AliasName
}

func (opt *Opt) getOnConnect() func(ctx context.Context, cn *redis.Conn) error {
	if opt.OnConnect == nil {
		return func(ctx context.Context, cn *redis.Conn) error {
			//log.Println("redis Onconnect at:", cn.String())
			return nil
		}
	}
	return opt.OnConnect
}

// WithOnConnect 添加OnConnect回调方法
func (opt *Opt) WithOnConnect(fn func(ctx context.Context, cn *redis.Conn) error) *Opt {
	opt.OnConnect = fn
	return opt
}

// WithDialer 添加Dialer回调方法
func (opt *Opt) WithDialer(fn func(ctx context.Context, network string, addr string) (net.Conn, error)) *Opt {
	opt.Dialer = fn
	return opt
}

// WithTLSConfig 添加自定义TLSConfig
func (opt *Opt) WithTLSConfig(cfg *tls.Config) *Opt {
	opt.TLSConfig = cfg
	return opt
}

// Init 使用传入的Opt初始化redis，支持链接多个Opt
func Init(opts ...*Opt) {
	for _, opt := range opts {
		if !opt.IsCluster {
			pool.Store(opt.getAliasName(), &Cli{opt.getAliasName(), opt, newClient(opt), nil})
			continue
		}
		pool.Store(opt.getAliasName(), &Cli{opt.getAliasName(), opt, newClusterClient(opt), nil})
	}
}

func InitMigration(client redis.UniversalClient) {
	pool.Store("default", &Cli{"default", nil, client, nil})
}

func newClient(opt *Opt) redis.UniversalClient {
	return redis.NewClient(&redis.Options{
		Dialer:       opt.Dialer,
		Network:      opt.Network,
		Addr:         opt.getAddrs()[0],
		OnConnect:    opt.getOnConnect(),
		Username:     opt.Username,
		Password:     opt.Password,
		DB:           opt.DB,
		ReadTimeout:  opt.ReadTimeout,
		WriteTimeout: opt.WriteTimeout,
		PoolSize:     opt.PoolSize,
		MinIdleConns: opt.getMinIdleConns(),
		TLSConfig:    opt.TLSConfig,
	})
}

func newClusterClient(opt *Opt) redis.UniversalClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Dialer: opt.Dialer,
		Addrs:  opt.getAddrs(),
		//ReadOnly:           false,
		//RouteByLatency:     false,
		//RouteRandomly:      false,
		OnConnect:    opt.getOnConnect(),
		Username:     opt.Username,
		Password:     opt.Password,
		ReadTimeout:  opt.ReadTimeout,
		WriteTimeout: opt.WriteTimeout,
		PoolSize:     opt.PoolSize,
		MinIdleConns: opt.getMinIdleConns(),
		TLSConfig:    opt.TLSConfig,
	})
}

type Cli struct {
	AliasName string
	opt       *Opt
	rdb       redis.UniversalClient
	ctx       context.Context
}

// RDB 返回当前别名的Redis client
// aliasName为空则使用 default
// RDB 返回当前别名的Redis client
// aliasName为空则使用 default
func RDB(aliasName ...string) *Cli {
	name := "default"
	if len(aliasName) != 0 {
		name = aliasName[0]
	}
	if latestCli != nil && latestCli.AliasName == name {
		return latestCli
	}
	v, ok := pool.Load(name)
	if !ok {
		panic(fmt.Errorf("no %s cli in RDB pool", name))
	}
	cli := v.(*Cli)
	return &Cli{
		AliasName: cli.AliasName,
		opt:       cli.opt,
		rdb:       cli.rdb,
		ctx:       nil,
	}
}

// 该方法旨在兼容处理未传入标准时间间隔的情况，例如传入数字
// 局限：
//
//	仅能处理秒级时间参数
//	参数需要小于 time.Second 的大小(< 1000000000)
//
// Warning:
//
//	这只是一个兼容处理，并不是一个好的解决方案，正常情况下请按参数类型传值
func getDuration(t time.Duration) time.Duration {
	if t < time.Second {
		return t * time.Second
	}
	return t
}

// WithCtx 使用自定义的ctx以控制执行流程
func (c *Cli) WithCtx(ctx context.Context) *Cli {
	c.ctx = ctx
	return c
}

func (c *Cli) getCtx() context.Context {
	if c.ctx == nil {
		return defaultCtx
	}
	return c.ctx
}

// Cli 获取一个redis client
func (c *Cli) Cli() redis.UniversalClient {
	return c.rdb
}

// Keys 按照给定规则返回符合条件的key
// Deprecated:	鉴于keys指令会导致查询时间过长而影响redis性能，keys不再被建议使用
// 建议使用 Scan 替换
func (c *Cli) Keys(pattern string) ([]string, error) {
	return c.rdb.Keys(c.getCtx(), pattern).Result()
}

// Scan 按照给定规则返回符合条件的key，支持翻页和指定单页数量
func (c *Cli) Scan(cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.rdb.Scan(c.getCtx(), cursor, match, count).Result()
}

// TTL 返回给定 key 的剩余生存时间(TTL, time to live)
// p nil|false 以秒为单位
// p true 以毫秒为单位
// 当 key 不存在时，返回 -2
// 当 key 存在但没有设置剩余生存时间时，返回 -1
func (c *Cli) TTL(key string, p ...bool) (time.Duration, error) {
	if p != nil && p[0] {
		return c.rdb.PTTL(c.getCtx(), key).Result()
	}
	return c.rdb.TTL(c.getCtx(), key).Result()
}

// Expire 为给定 key 设置生存时间，当 key 过期时(生存时间为 0 )，它会被自动删除
// p nil|false 以秒为单位
// p true 以毫秒为单位
func (c *Cli) Expire(key string, expiration time.Duration, p ...bool) (bool, error) {
	if p != nil && p[0] {
		return c.rdb.PExpire(c.getCtx(), key, getDuration(expiration)).Result()
	}
	return c.rdb.Expire(c.getCtx(), key, getDuration(expiration)).Result()
}

// ExpireAt 为给定 key 设置生存时间，当 key 过期时(生存时间为 0 )，它会被自动删除
// expireTime 为key到期的 UNIX 时间
func (c *Cli) ExpireAt(key string, expireTime time.Time) (bool, error) {
	return c.rdb.ExpireAt(c.getCtx(), key, expireTime).Result()
}

// Del 删除给定的一个或多个 key
// 返回被删除 key 的数量
func (c *Cli) Del(key ...string) (int64, error) {
	return c.rdb.Del(c.getCtx(), key...).Result()
}

// Exists 检查给定 key 是否存在
func (c *Cli) Exists(key ...string) (bool, error) {
	r, err := c.rdb.Exists(c.getCtx(), key...).Result()
	return r == int64(len(key)), err
}

// Move 将当前数据库的 key 移动到给定的数据库 db 当中
func (c *Cli) Move(key string, db int) (bool, error) {
	return c.rdb.Move(c.getCtx(), key, db).Result()
}

// Get 返回与键 key 相关联的字符串值
func (c *Cli) Get(key string) (string, error) {
	return c.rdb.Get(c.getCtx(), key).Result()
}

// MGet 返回与键 key 相关联的字符串值数组
func (c *Cli) MGet(key ...string) ([]any, error) {
	return c.rdb.MGet(c.getCtx(), key...).Result()
}

const KeepTTL = redis.KeepTTL

// Set 将字符串值 value 关联到 key
// expiration >0则设置过期时间，0不设置过期时间，需要set时保留之前的TTL则使用 KeepTTL
func (c *Cli) Set(key string, val any, expiration ...time.Duration) error {
	expiration = append(expiration, 0)
	return c.rdb.Set(c.getCtx(), key, val, getDuration(expiration[0])).Err()
}

// SetNx 只在键 key 不存在的情况下，将键 key 的值设置为 value
func (c *Cli) SetNx(key string, val any, expiration time.Duration) (bool, error) {
	return c.rdb.SetNX(c.getCtx(), key, val, getDuration(expiration)).Result()
}

// SetEx 将键 key 的值设置为 value ，并将键 key 的生存时间设置为 expiration 秒钟
func (c *Cli) SetEx(key string, val any, expiration time.Duration) (string, error) {
	return c.rdb.SetEx(c.getCtx(), key, val, getDuration(expiration)).Result()
}

// Incr 为键 key 储存的数字值加上一
// 值限制在 64 位(bit)有符号数字表示之内
// 如果键 key 不存在，那么它的值会先被初始化为 0 ，然后再执行 INCR 命令
// 如果键 key 储存的值不能被解释为数字，那么 INCR 命令将返回一个错误
// 返回键 key 在执行加一操作之后的值
func (c *Cli) Incr(key string) (int64, error) {
	return c.rdb.Incr(c.getCtx(), key).Result()
}

// IncrBy 为键 key 储存的数字值加上增量 increment
// 值限制在 64 位(bit)有符号数字表示之内
// 如果键 key 不存在，那么它的值会先被初始化为 0 ，然后再执行 INCR 命令
// 如果键 key 储存的值不能被解释为数字，那么 INCR 命令将返回一个错误
// 返回键 key 在执行加一操作之后的值
func (c *Cli) IncrBy(key string, increment int64) (int64, error) {
	return c.rdb.IncrBy(c.getCtx(), key, increment).Result()
}

// IncrByFloat 为键 key 储存的数字值加上增量 increment
// 无论加法计算所得的浮点数的实际精度有多长，计算结果最多只保留小数点的后十七位
// 如果键 key 不存在，那么它的值会先被初始化为 0 ，然后再执行 INCR 命令
// 如果键 key 储存的值不能被解释为数字，那么 INCR 命令将返回一个错误
// 返回键 key 在执行加一操作之后的值
func (c *Cli) IncrByFloat(key string, increment float64) (float64, error) {
	return c.rdb.IncrByFloat(c.getCtx(), key, increment).Result()
}

// Decr 为键 key 储存的数字值减去一
// 值限制在 64 位(bit)有符号数字表示之内
// 如果键 key 不存在，那么键 key 的值会先被初始化为 0 ，然后再执行 DECR 操作
// 如果键 key 储存的值不能被解释为数字，那么 DECR 命令将返回一个错误
func (c *Cli) Decr(key string) (int64, error) {
	return c.rdb.Decr(c.getCtx(), key).Result()
}

// DecrBy 将键 key 储存的整数值减去减量 decrement
// 如果键 key 不存在，那么键 key 的值会先被初始化为 0 ，然后再执行 DECR 操作
// 如果键 key 储存的值不能被解释为数字，那么 DECR 命令将返回一个错误
func (c *Cli) DecrBy(key string, decrement int64) (int64, error) {
	return c.rdb.DecrBy(c.getCtx(), key, decrement).Result()
}

// Watch 监视一个(或多个) key ，如果在事务执行之前这个(或这些) key 被其他命令所改动，那么事务将被打断。
func (c *Cli) Watch(fn func(tx *Tx) error, keys ...string) error {
	err := c.rdb.Watch(c.getCtx(), func(tx *redis.Tx) error {
		return fn(&Tx{tx})
	}, keys...)
	return err
}

// HGet 返回哈希表中给定域的值
func (c *Cli) HGet(key, field string) (string, error) {
	return c.rdb.HGet(c.getCtx(), key, field).Result()
}

// HMGet 返回哈希表 key 中，一个或多个给定域的值
// 返回值 一个包含多个给定域的关联值的表，表值的排列顺序和给定域参数的请求顺序一样
// 不存在的域返回nil值
func (c *Cli) HMGet(key string, field ...string) ([]any, error) {
	return c.rdb.HMGet(c.getCtx(), key, field...).Result()
}

// HMSet 同时将多个 field-value (域-值)对设置到哈希表 key 中
// HSet("hash", "key1", "value1", "key2", "value2")
// HSet("hash", []string{"key1", "value1", "key2", "value2"})
// HSet("hash", map[string]any{"key1": "value1", "key2": "value2"})
func (c *Cli) HMSet(key string, values ...any) (int64, error) {
	return c.rdb.HSet(c.getCtx(), key, values...).Result()
}

// HGetAll 返回哈希表 key 中，所有的域和值
func (c *Cli) HGetAll(key string) (map[string]string, error) {
	return c.rdb.HGetAll(c.getCtx(), key).Result()
}

// HSet 向指定 key 的 hash 表插入值
func (c *Cli) HSet(key, field string, values any) (int64, error) {
	return c.rdb.HSet(c.getCtx(), key, field, values).Result()
}

// HSetEx 并将键 key 的生存时间设置为 expiration 秒钟
// 并非严谨实现方式，有可能造成 Expire 设置失败
func (c *Cli) HSetEx(key, field string, values any, expiration time.Duration) (int64, error) {
	script := `
redis.call('HSET', KEYS[1], ARGV[1], ARGV[2])
redis.call('EXPIRE', KEYS[1], tonumber(ARGV[3]))
return 1
`
	return c.rdb.Eval(c.getCtx(), script, []string{key}, field, values, expiration/time.Second).Int64()
}

// HExists 检查给定域 field 是否存在于哈希表 hash 当中
func (c *Cli) HExists(key, field string) (bool, error) {
	return c.rdb.HExists(c.getCtx(), key, field).Result()
}

// HDel 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略
func (c *Cli) HDel(key, field string) (int64, error) {
	return c.rdb.HDel(c.getCtx(), key, field).Result()
}

// HIncrBy 为哈希表中的字段值加上指定增量值。增量也可以为负数，相当于对指定字段进行减法操作
func (c *Cli) HIncrBy(key, field string, incr int64) (int64, error) {
	return c.rdb.HIncrBy(c.getCtx(), key, field, incr).Result()
}

// HIncrByFloat 为哈希表中的字段值加上指定浮点类型增量值。增量也可以为负数，相当于对指定字段进行减法操作
func (c *Cli) HIncrByFloat(key, field string, incr float64) (float64, error) {
	return c.rdb.HIncrByFloat(c.getCtx(), key, field, incr).Result()
}

// LPush 将一个或多个值 value 插入到列表 key 的表头
func (c *Cli) LPush(key string, value ...any) (int64, error) {
	return c.rdb.LPush(c.getCtx(), key, value...).Result()
}

// LSet 将一个值 value 插入到列表 key 的index位置
func (c *Cli) LSet(key string, index int64, value any) (string, error) {
	return c.rdb.LSet(c.getCtx(), key, index, value).Result()
}

// RPush 将一个或多个值 value 插入到列表 key 的表尾(最右边)
func (c *Cli) RPush(key string, value any) (int64, error) {
	return c.rdb.RPush(c.getCtx(), key, value).Result()
}

// LPop 移除并返回列表 key 的头元素
func (c *Cli) LPop(key string) (string, error) {
	return c.rdb.LPop(c.getCtx(), key).Result()
}

// RPop 移除并返回列表 key 的尾元素
func (c *Cli) RPop(key string) (string, error) {
	return c.rdb.RPop(c.getCtx(), key).Result()
}

// LLen 返回列表 key 的长度
func (c *Cli) LLen(key string) (int64, error) {
	return c.rdb.LLen(c.getCtx(), key).Result()
}

// LRange 返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定
func (c *Cli) LRange(key string, start, stop int64) ([]string, error) {
	return c.rdb.LRange(c.getCtx(), key, start, stop).Result()
}

// Z SortSet结构
type Z struct {
	Score  float64
	Member any
}

func getZ(z []Z) (s []redis.Z) {
	for _, e := range z {
		s = append(s, redis.Z{
			Score:  e.Score,
			Member: e.Member,
		})
	}
	return
}

// ZAdd 将一个或多个 member 元素及其 score 值加入到有序集 key 当中
// 返回被成功添加的新成员的数量，不包括那些被更新的、已经存在的成员
func (c *Cli) ZAdd(key string, members ...Z) (int64, error) {
	return c.rdb.ZAdd(c.getCtx(), key, getZ(members)...).Result()
}

// ZRange 返回有序集 key 中，指定区间内的成员
// 其中成员的位置按 score 值递增(从小到大)来排序。具有相同 score 值的成员按字典序(lexicographical order )来排列。
// 返回指定区间内的有序集成员的列表
func (c *Cli) ZRange(key string, start, stop int64) ([]string, error) {
	return c.rdb.ZRange(c.getCtx(), key, start, stop).Result()
}

// ZRangeWithScore 返回有序集 key 中，指定区间内的成员
// 其中成员的位置按 score 值递增(从小到大)来排序。具有相同 score 值的成员按字典序(lexicographical order )来排列。
// 返回指定区间内，带有 score 值的有序集成员的列表
func (c *Cli) ZRangeWithScore(key string, start, stop int64) ([]redis.Z, error) {
	return c.rdb.ZRangeWithScores(c.getCtx(), key, start, stop).Result()
}

// ZRevRange 返回有序集 key 中，指定区间内的成员
// 其中成员的位置按 score 值递减(从大到小)来排列。具有相同 score 值的成员按字典序的逆序
// 返回指定区间内的有序集成员的列表
func (c *Cli) ZRevRange(key string, start, stop int64) ([]string, error) {
	return c.rdb.ZRevRange(c.getCtx(), key, start, stop).Result()
}

// ZRevRangeWithScore 返回有序集 key 中，指定区间内的成员
// 其中成员的位置按 score 值递减(从大到小)来排列。具有相同 score 值的成员按字典序的逆序
// 返回指定区间内，带有 score 值的有序集成员的列表
func (c *Cli) ZRevRangeWithScore(key string, start, stop int64) ([]redis.Z, error) {
	return c.rdb.ZRevRangeWithScores(c.getCtx(), key, start, stop).Result()
}

// ZRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)来排序
func (c *Cli) ZRank(key string, member string) (int64, error) {
	return c.rdb.ZRank(c.getCtx(), key, member).Result()
}

// ZRevRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序
func (c *Cli) ZRevRank(key string, member string) (int64, error) {
	return c.rdb.ZRevRank(c.getCtx(), key, member).Result()
}

// ZRem 移除有序集 key 中的一个或多个成员，不存在的成员将被忽略
func (c *Cli) ZRem(key string, member ...any) (int64, error) {
	return c.rdb.ZRem(c.getCtx(), key, member...).Result()
}

// ZRemRangeByScore 移除有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max)的成员
func (c *Cli) ZRemRangeByScore(key string, min, max string) (int64, error) {
	return c.rdb.ZRemRangeByScore(c.getCtx(), key, min, max).Result()
}

// ZRemRangeByRank 移除有序集 key 中，指定排名(rank)区间内的所有成员
func (c *Cli) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	return c.rdb.ZRemRangeByRank(c.getCtx(), key, start, stop).Result()
}

// PFAdd 将任意数量的元素添加到指定的 HyperLogLog 里面
// 如果 HyperLogLog 的内部储存被修改了，那么返回 1，否则返回 0
func (c *Cli) PFAdd(key string, els ...any) (int64, error) {
	return c.rdb.PFAdd(c.getCtx(), key, els...).Result()
}

// PFCount
// 当传入单个键时，返回储存在给定键的 HyperLogLog 的近似基数，如果键不存在，那么返回 0
// 当传入多个键时，返回所有给定 HyperLogLog 的并集的近似基数，这个近似基数是通过将所有给定 HyperLogLog 合并至一个临时 HyperLogLog 来计算得出的
// 整数回复：给定 HyperLogLog 包含的唯一元素的近似数量
func (c *Cli) PFCount(key ...string) (int64, error) {
	return c.rdb.PFCount(c.getCtx(), key...).Result()
}

// PFMerge 将多个 HyperLogLog 合并（merge）为一个 HyperLogLog ，
// 合并后的 HyperLogLog 的基数接近于所有输入 HyperLogLog 的可见集合（observed set）的并集。
func (c *Cli) PFMerge(dest string, key ...string) (string, error) {
	return c.rdb.PFMerge(c.getCtx(), dest, key...).Result()
}

// SAdd 向集合添加一个或多个成员
func (c *Cli) SAdd(key string, members ...any) (int64, error) {
	return c.rdb.SAdd(c.getCtx(), key, members...).Result()
}

// SCard 获取集合的成员数
func (c *Cli) SCard(key string) (int64, error) {
	return c.rdb.SCard(c.getCtx(), key).Result()
}

// SDiff 返回第一个集合与其他集合之间的差异。
func (c *Cli) SDiff(keys ...string) ([]string, error) {
	return c.rdb.SDiff(c.getCtx(), keys...).Result()
}

// SDiffStore 返回给定所有集合的差集并存储在 destination 中
func (c *Cli) SDiffStore(destination string, keys ...string) (int64, error) {
	return c.rdb.SDiffStore(c.getCtx(), destination, keys...).Result()
}

// SInter 返回给定所有集合的差集
func (c *Cli) SInter(keys ...string) ([]string, error) {
	return c.rdb.SInter(c.getCtx(), keys...).Result()
}

// SInterStore 返回给定所有集合的差集并存储在 destination 中
func (c *Cli) SInterStore(destination string, keys ...string) (int64, error) {
	return c.rdb.SInterStore(c.getCtx(), destination, keys...).Result()
}

// SInterCard 获取集合差集的成员数
func (c *Cli) SInterCard(limit int64, keys ...string) (int64, error) {
	return c.rdb.SInterCard(c.getCtx(), limit, keys...).Result()
}

// SIsMember 判断 member 元素是否是集合 key 的成员
func (c *Cli) SIsMember(key string, member string) (bool, error) {
	return c.rdb.SIsMember(c.getCtx(), key, member).Result()
}

// SMIsMember 判断 members 元素是否都是集合 key 的成员
func (c *Cli) SMIsMember(key string, members ...any) ([]bool, error) {
	return c.rdb.SMIsMember(c.getCtx(), key, members...).Result()
}

// SMembers 返回集合中的所有成员
func (c *Cli) SMembers(key string) ([]string, error) {
	return c.rdb.SMembers(c.getCtx(), key).Result()
}

// SMembersMap 以map的形式返回集合中的所有成员
func (c *Cli) SMembersMap(key string) (map[string]struct{}, error) {
	return c.rdb.SMembersMap(c.getCtx(), key).Result()
}

// SPop 移除并返回集合中的一个随机元素
func (c *Cli) SPop(key string) (string, error) {
	return c.rdb.SPop(c.getCtx(), key).Result()
}

// SPopN 移除并返回集合中的N个随机元素
func (c *Cli) SPopN(key string, count int64) ([]string, error) {
	return c.rdb.SPopN(c.getCtx(), key, count).Result()
}

// SRandMember 返回集合中一个随机数
func (c *Cli) SRandMember(key string) (string, error) {
	return c.rdb.SRandMember(c.getCtx(), key).Result()
}

// SRandMemberN 返回集合中多个随机数
func (c *Cli) SRandMemberN(key string, count int64) ([]string, error) {
	return c.rdb.SRandMemberN(c.getCtx(), key, count).Result()
}

// SRem 移除集合中一个或多个成员
func (c *Cli) SRem(key string, members ...any) (int64, error) {
	return c.rdb.SRem(c.getCtx(), key, members...).Result()
}

// SUnion 返回所有给定集合的并集
func (c *Cli) SUnion(keys ...string) ([]string, error) {
	return c.rdb.SUnion(c.getCtx(), keys...).Result()
}

// SUnionStore 所有给定集合的并集存储在 destination 集合中
func (c *Cli) SUnionStore(destination string, keys ...string) (int64, error) {
	return c.rdb.SUnionStore(c.getCtx(), destination, keys...).Result()
}

// SScan 所有给定集合的并集存储在 destination 集合中
func (c *Cli) SScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.rdb.SScan(c.getCtx(), key, cursor, match, count).Result()
}

// GeoLocation GOE位置对象
type GeoLocation struct {
	Name                      string
	Longitude, Latitude, Dist float64
	GeoHash                   int64
}

func getGeoLocation(geoLocation []*GeoLocation) (s []*redis.GeoLocation) {
	for _, e := range geoLocation {
		s = append(s, &redis.GeoLocation{
			Name:      e.Name,
			Longitude: e.Longitude,
			Latitude:  e.Latitude,
			Dist:      e.Dist,
			GeoHash:   e.GeoHash,
		})
	}
	return
}

// GeoRadiusQuery GOE半径对象
type GeoRadiusQuery struct {
	Radius float64
	// Can be m, km, ft, or mi. Default is km.
	Unit        string
	WithCoord   bool
	WithDist    bool
	WithGeoHash bool
	Count       int
	// Can be ASC or DESC. Default is no sort order.
	Sort      string
	Store     string
	StoreDist string
}

// Redis GEO GeoRadiusQuery Unit
const (
	M  = "m"
	KM = "km"
	FT = "ft"
	MI = "mi"
)

// Redis GEO GeoRadiusQuery Sort
const (
	ASC  = "ASC"
	DESC = "DESC"
)

func getGeoRadiusQuery(e *GeoRadiusQuery) *redis.GeoRadiusQuery {
	return &redis.GeoRadiusQuery{
		Radius:      e.Radius,
		Unit:        e.Unit,
		WithCoord:   e.WithCoord,
		WithDist:    e.WithDist,
		WithGeoHash: e.WithGeoHash,
		Count:       e.Count,
		Sort:        e.Sort,
		Store:       e.Store,
		StoreDist:   e.StoreDist,
	}
}

// GeoAdd 将给定的空间元素（纬度、经度、名字）添加到指定的键里面。
// 这些数据会以有序集合的形式被储存在键里面，从而使得像 GEORADIUS 和 GEORADIUSBYMEMBER 这样的命令可以在之后通过位置查询取得这些元素
// 返回新添加到键里面的空间元素数量，不包括那些已经存在但是被更新的元素
func (c *Cli) GeoAdd(key string, geoLocation ...*GeoLocation) (int64, error) {
	return c.rdb.GeoAdd(c.getCtx(), key, getGeoLocation(geoLocation)...).Result()
}

// GeoPos 从键里面返回所有给定位置元素的位置（经度和纬度）
// 返回一个数组， 数组中的每个项都由两个元素组成： 第一个元素为给定位置元素的经度， 而第二个元素则为给定位置元素的纬度
func (c *Cli) GeoPos(key string, member ...string) ([]*redis.GeoPos, error) {
	return c.rdb.GeoPos(c.getCtx(), key, member...).Result()
}

// GeoDist 返回两个给定位置之间的距离
// 计算出的距离会以双精度浮点数的形式被返回。 如果给定的位置元素不存在， 那么命令返回空值
func (c *Cli) GeoDist(key, member1, member2, unit string) (float64, error) {
	return c.rdb.GeoDist(c.getCtx(), key, member1, member2, unit).Result()
}

// GeoRadius 以给定的经纬度为中心， 返回键包含的位置元素当中， 与中心的距离不超过给定最大距离的所有位置元素
// 在给定以下可选项时， 命令会返回额外的信息：
// WithDist ： 在返回位置元素的同时， 将位置元素与中心之间的距离也一并返回。 距离的单位和用户给定的范围单位保持一致。
// WithCoord ： 将位置元素的经度和维度也一并返回。
// WithHash ： 以 52 位有符号整数的形式， 返回位置元素经过原始 geohash 编码的有序集合分值。 这个选项主要用于底层应用或者调试， 实际中的作用并不大。
func (c *Cli) GeoRadius(key string, longitude, latitude float64, query *GeoRadiusQuery) ([]redis.GeoLocation, error) {
	return c.rdb.GeoRadius(c.getCtx(), key, longitude, latitude, getGeoRadiusQuery(query)).Result()
}

// GeoRadiusByMember 和 GeoRadius一样， 都可以找出位于指定范围内的元素
// 但是 GeoRadiusByMember 的中心点是由给定的位置元素决定的
func (c *Cli) GeoRadiusByMember(key, member string, query *GeoRadiusQuery) ([]redis.GeoLocation, error) {
	return c.rdb.GeoRadiusByMember(c.getCtx(), key, member, getGeoRadiusQuery(query)).Result()
}

// ScriptLoad 将脚本 script 添加到脚本缓存中，但并不立即执行这个脚本
func (c *Cli) ScriptLoad(script string) (string, error) {
	return c.rdb.ScriptLoad(c.getCtx(), script).Result()
}

// ScriptExists 给定一个或多个脚本的SHA1校验和
// 返回一个包含 0 和 1 的列表，表示校验和所指定的脚本是否已经被保存在缓存当中
func (c *Cli) ScriptExists(script string) ([]bool, error) {
	return c.rdb.ScriptExists(c.getCtx(), script).Result()
}

// ScriptKill 杀死当前正在运行的 Lua 脚本，当且仅当这个脚本没有执行过任何写操作时，这个命令才生效
// 执行成功返回 OK ，否则返回一个错误
// ERR Sorry the script already executed write commands against the dataset.
// You can either wait the script termination or kill the server in an hard way using the SHUTDOWN NOSAVE command.
func (c *Cli) ScriptKill() (string, error) {
	return c.rdb.ScriptKill(c.getCtx()).Result()
}

// EvalSha 根据给定的 sha1 校验码，对缓存在服务器中的脚本进行求值
func (c *Cli) EvalSha(sha1 string, keys []string, args ...any) (any, error) {
	return c.rdb.EvalSha(c.getCtx(), sha1, keys, args...).Result()
}

// Eval 使用 EVAL 命令对 Lua 脚本进行求值
func (c *Cli) Eval(script string, keys []string, args ...any) (any, error) {
	return c.rdb.Eval(c.getCtx(), script, keys, args...).Result()
}

func (c *Cli) EvalCmd(script string, keys []string, args ...any) *redis.Cmd {
	return c.rdb.Eval(c.getCtx(), script, keys, args...)
}

// Publish 发布
func (c *Cli) Publish(channel string, message any) (int64, error) {
	return c.rdb.Publish(c.getCtx(), channel, message).Result()
}

// Subscribe 订阅
func (c *Cli) Subscribe(channel ...string) *redis.PubSub {
	return c.rdb.Subscribe(c.getCtx(), channel...)
}

// Pipeline 管道
func (c *Cli) Pipeline() Pipeliner {
	return Pipeliner{c.rdb.Pipeline()}
}

// Pipelined 管道回调方式
func (c *Cli) Pipelined(ctx context.Context, fn func(pipe Pipeliner) error) ([]redis.Cmder, error) {
	return c.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		return fn(Pipeliner{pipe})
	})
}
