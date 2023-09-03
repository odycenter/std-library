// Package elastic ElasticSearch连接和操作方法封装
package elastic

import (
	"context"
	"crypto/tls"
	"github.com/olivere/elastic/v7"
	"log"
	"net/http"
	"sync"
	"time"
)

var clients clientMap

func init() {
	clients = clientMap{
		m: make(map[string]*Client),
	}
}

type clientMap struct {
	m map[string]*Client
	sync.RWMutex
}

func (c *clientMap) load(k string) (*Client, bool) {
	c.RLock()
	v, ok := c.m[k]
	c.RUnlock()
	return v, ok
}

func (c *clientMap) store(k string, v *Client) {
	c.Lock()
	c.m[k] = v
	c.Unlock()
}

// Client ES连接结构体
type Client struct {
	opt *Option
	cli *elastic.Client
}

// Option Elasticsearch 配置
type Option struct {
	AliasName            string        //别名
	URLs                 []string      //url 127.0.0.1:9200
	Scheme               string        //http or https
	Username             string        //用户名
	Password             string        //密码
	Shards               int           //分片
	Replicas             int           //副本
	Sniff                bool          //是否开启嗅探功能
	SniffInterval        time.Duration //嗅探间隔
	SniffIntervalTimeOut time.Duration //嗅探超时时间
	Healthcheck          bool          //是否开启健康监测
	HealthcheckInterval  time.Duration //健康检测间隔时间
	HealthcheckTimeout   time.Duration //健康检测超时时间
	BackoffRetrierTicks  []int         //回退重试间隔,e.i.:[]int{200,200,200}三次重试间隔200ms
	Gzip                 bool          //是否开启gzip压缩
	SkipTLS              bool          //跳过TLS验证
}

func (opt *Option) getAliasName() string {
	if opt.AliasName == "" {
		return "default"
	}
	return opt.AliasName
}

func (opt *Option) getScheme() string {
	if opt.Scheme == "" {
		return "https"
	}
	return opt.Scheme
}

// NewClient 创建ES客户端
func NewClient(opts ...*Option) {
	for _, opt := range opts {
		newClient(opt)
	}
}

func newClient(opt *Option) {
	var cof []elastic.ClientOptionFunc
	cof = append(cof,
		elastic.SetURL(opt.URLs...),
		elastic.SetScheme(opt.getScheme()),
		elastic.SetBasicAuth(opt.Username, opt.Password),
		elastic.SetHealthcheck(opt.Healthcheck),
		elastic.SetHealthcheckInterval(opt.HealthcheckInterval),
		elastic.SetHealthcheckTimeout(opt.HealthcheckTimeout),
		elastic.SetSniff(opt.Sniff),
		elastic.SetSnifferInterval(opt.SniffInterval),
		elastic.SetSnifferTimeout(opt.SniffIntervalTimeOut),
		elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewSimpleBackoff(opt.BackoffRetrierTicks...))),
		elastic.SetGzip(opt.Gzip),
	)
	if opt.SkipTLS {
		cof = append(cof, elastic.SetHttpClient(&http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}))
	}
	cli, err := elastic.NewClient(cof...)
	if err != nil {
		log.Panicf("ES new client <%s> failed \n", err.Error())
	}
	clients.store(opt.getAliasName(), &Client{opt, cli})
}

// Cli 根据别名获取ES客户端
func Cli(aliasName ...string) *Client {
	name := "default"
	if aliasName != nil {
		name = aliasName[0]
	}
	v, ok := clients.load(name)
	if !ok {
		log.Panicf("elastic client <%s> not exists(need create)\n", name)
	}
	return v
}

// Version 获取ES版本
func (es *Client) Version(url string) (string, error) {
	return es.cli.ElasticsearchVersion(url)
}

// IndexExists 索引是否已存在
func (es *Client) IndexExists(ctx context.Context, indexName ...string) (bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return es.cli.IndexExists(indexName...).Do(ctx)
}

// CreateIndex 创建索引
func (es *Client) CreateIndex(ctx context.Context, indexName, bodyJsonMap string) (*elastic.IndicesCreateResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return es.cli.CreateIndex(indexName).BodyString(bodyJsonMap).Do(ctx)
}

// UpdateIndex 更新索引
func (es *Client) UpdateIndex(ctx context.Context, indexName, bodyJsonMap string) (*elastic.PutMappingResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return es.cli.PutMapping().Index(indexName).BodyString(bodyJsonMap).Do(ctx)
}

// Insert 插入数据
func (es *Client) Insert(ctx context.Context, indexName, indexId string, body any) (*elastic.IndexResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return es.cli.Index().Index(indexName).Id(indexId).BodyJson(body).Do(ctx)
}

// BulkInsert 批量插入数据
func (es *Client) BulkInsert(ctx context.Context, indexName string, bodyMap map[string]any) (*elastic.BulkResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	bulk := es.cli.Bulk()
	for k, v := range bodyMap {
		bulk.Add(elastic.NewBulkIndexRequest().Index(indexName).Id(k).Doc(v))
	}
	return bulk.Do(ctx)
}

// 查询方式
const (
	Term = queryKey(iota)
	Range
)

// SearchAdvance 拼装高级搜索
func (es *Client) SearchAdvance(ctx context.Context, query *SearchQuery) (*elastic.SearchResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if query == nil {
		return nil, nil
	}
	if query.Index == nil {
		log.Panicf("elestic search need index\n")
	}
	searcher := es.cli.Search().Index(query.Index...)
	if query.From != nil {
		searcher = searcher.From(*query.From)
	}
	if query.Size != nil {
		searcher = searcher.Size(*query.Size)
	}
	if query.Pretty != nil {
		searcher = searcher.Pretty(*query.Pretty)
	}
	if query.Sort != nil {
		for _, sort := range query.Sort {
			searcher = searcher.Sort(sort.Field, sort.Ascending)
		}
	}
	if query.QueryEntity != nil {
		for _, queryEntity := range query.QueryEntity {
			switch queryEntity.Key {
			case Term:
				for _, queryBody := range queryEntity.Values {
					searcher = searcher.Query(elastic.NewTermQuery(queryBody.Key, queryBody.Value))

				}
			case Range:
				for _, queryBody := range queryEntity.Values {
					rangeQuery := elastic.NewRangeQuery(queryBody.Key)
					if queryBody.Gte != nil {
						rangeQuery.Gte(queryBody.Gte)
					}
					if queryBody.Lte != nil {
						rangeQuery.Lte(queryBody.Lte)
					}
					if queryBody.Gt != nil {
						rangeQuery.Gt(queryBody.Gt)
					}
					if queryBody.Gt != nil {
						rangeQuery.Lt(queryBody.Lt)
					}
					searcher = searcher.Query(rangeQuery)
				}
			}
		}
	}
	return searcher.Do(ctx)
}

// SearchSimple 简单查询ES数据,精确查询某条数据
// indexName 索引名 - 对应MySQL库名
// indexId 索引ID - 默认为logID
func (es *Client) SearchSimple(ctx context.Context, indexName, indexId string) (*elastic.GetResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return es.cli.Get().Index(indexName).Id(indexId).Do(ctx)
}

// SearchSimpleWithBody 简单查询ES数据,按条件查询
// query SimpleSearchQuery简单查询条件
func (es *Client) SearchSimpleWithBody(ctx context.Context, query *SimpleSearchQuery) (*elastic.SearchResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	searcher := es.cli.Search()
	searcher = searcher.Index(query.Index...)
	if query.IDs != nil {
		searcher = searcher.Query(elastic.NewIdsQuery().Ids(query.IDs...))
	}
	if query.From != nil {
		searcher = searcher.From(*query.From)
	}
	searcher = searcher.Size(*query.getSize())
	if query.Sort != nil {
		for _, sort := range query.Sort {
			searcher = searcher.Sort(sort.Field, sort.Ascending)
		}
	}
	if query.Pretty != nil {
		searcher = searcher.Pretty(*query.Pretty)
	}
	return searcher.Do(ctx)
}

// DelIndex 删除指定的索引
func (es *Client) DelIndex(ctx context.Context, indexName ...string) (*elastic.IndicesDeleteResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return es.cli.DeleteIndex(indexName...).Do(ctx)
}
