// Package nacos nacos连接及操作封装
package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
)

var (
	namingClient naming_client.INamingClient
	configClient config_client.IConfigClient
)

// ServerConfig Nacos连接配置
type ServerConfig struct {
	IpAddr string // Nacos的服务地址
	Port   uint64 // Nacos的服务端口
	Scheme string // Nacos的服务地址前缀，默认http，在2.0中不需要设置
}

func newSrvCfg(srvCfg []*ServerConfig) []constant.ServerConfig {
	var settings []constant.ServerConfig
	for _, cfg := range srvCfg {
		settings = append(settings, constant.ServerConfig{
			Scheme: cfg.Scheme,
			IpAddr: cfg.IpAddr,
			Port:   cfg.Port,
		})
	}
	return settings
}

// ClientConfig nacos客户端配置
type ClientConfig struct {
	TimeoutMs   uint64 // 请求Nacos服务端的超时时间，默认是10000ms
	NamespaceId string // ACM的命名空间Id
	Endpoint    string // 当使用ACM时，需要该配置. https://help.aliyun.com/document_detail/130146.html
	RegionId    string // ACM&KMS的regionId，用于配置中心的鉴权
	AccessKey   string // ACM&KMS的AccessKey，用于配置中心的鉴权
	SecretKey   string // ACM&KMS的SecretKey，用于配置中心的鉴权
	OpenKMS     bool   // 是否开启kms，默认不开启，kms可以参考文档 https://help.aliyun.com/product/28933.html
	// 同时DataId必须以"cipher-"作为前缀才会启动加解密逻辑
	CacheDir             string // 缓存service信息的目录，默认是当前运行目录
	UpdateThreadNum      int    // 监听service变化的并发数，默认20
	NotLoadCacheAtStart  bool   // 在启动的时候不读取缓存在CacheDir的service信息
	UpdateCacheWhenEmpty bool   // 当service返回的实例列表为空时，不更新缓存，用于推空保护
	Username             string // Nacos服务端的API鉴权Username
	Password             string // Nacos服务端的API鉴权Password
	LogDir               string // 日志存储路径
	RotateTime           string // 日志轮转周期，比如：30m, 1h, 24h, 默认是24h
	MaxAge               int64  // 日志最大文件数，默认3
	LogLevel             string // 日志默认级别，值必须是：debug,info,warn,error，默认值是info
}

func newCliCfg(cliCfg *ClientConfig) *constant.ClientConfig {
	return &constant.ClientConfig{
		Endpoint:    cliCfg.Endpoint,
		NamespaceId: cliCfg.NamespaceId,
		RegionId:    cliCfg.RegionId,
		OpenKMS:     true,
		TimeoutMs:   5000,
		LogLevel:    "warn",
		LogDir:      ".",
	}
}

// Init 使用配置初始化nacos
func Init(cliCfg *ClientConfig, srvCfg ...*ServerConfig) {
	namingClient = newNamingClient(newCliCfg(cliCfg), newSrvCfg(srvCfg)...)
	configClient = newConfigClient(newCliCfg(cliCfg), newSrvCfg(srvCfg)...)
}

// 创建服务发现客户端
func newNamingClient(clientConfig *constant.ClientConfig, serverConfigs ...constant.ServerConfig) (cli naming_client.INamingClient) {
	var err error
	cli, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		log.Panicln("[nacos][newNamingClient]new naming client failed:\n", err)
	}
	return cli
}

// 创建动态配置客户端
func newConfigClient(clientConfig *constant.ClientConfig, serverConfigs ...constant.ServerConfig) (cli config_client.IConfigClient) {
	var err error
	cli, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		log.Panic("[nacos][newConfigClient]new config client failed:\n", err)
	}
	return cli
}

/*--------------------------------------------------------Config------------------------------------------------------*/

// GetSrvConfig 获取配置信息
func GetSrvConfig(dataId, group string) (string, error) {
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		log.Print("[nacos][GetSrvConfig]Get Config Failed:\n", err)
	}
	return content, err
}

// PublishConfig 发布配置信息
func PublishConfig(dataId, group, content string) (bool, error) {
	ok, err := configClient.PublishConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: content,
	})
	if err != nil {
		log.Print("[nacos][PublishConfig]Publish Config Failed:\n", err)
	}
	return ok, err
}

// ListenConfig 监听配置变化
func ListenConfig(dataId, group string, onChange func(namespace, group, dataId, data string)) error {
	err := configClient.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: onChange,
	})
	if err != nil {
		log.Print("[nacos][ListenConfig]Listen Config Failed:\n", err)
	}
	return err
}

// CancelListenConfig 取消监听
func CancelListenConfig(dataId, group string) error {
	err := configClient.CancelListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		log.Print("[nacos][CancelListenConfig]Cancel Listen Config Failed:\n", err)
	}
	return err
}

// DeleteConfig 删除配置
func DeleteConfig(dataId, group string) (bool, error) {
	ok, err := configClient.DeleteConfig(vo.ConfigParam{DataId: dataId,
		Group: group,
	})
	if err != nil {
		log.Print("[nacos][DeleteConfig]Delete Config Failed:\n", err)
	}
	return ok, err
}

// SearchConfig 搜索配置
func SearchConfig(search, dataId, group string, page, pageSize int) (*model.ConfigPage, error) {
	cfg, err := configClient.SearchConfig(vo.SearchConfigParam{
		Search:   search,
		DataId:   dataId,
		Group:    group,
		PageNo:   page,
		PageSize: pageSize,
	})
	if err != nil {
		log.Print("[nacos][SearchConfig]Search Config Failed:\n", err)
	}
	return cfg, err
}

/*--------------------------------------------------------Naming------------------------------------------------------*/

// RegisterInstance 注册实例
func RegisterInstance(serviceName, group, ip string, port uint64) (bool, error) {
	// 调用权重
	weight := 10
	ok, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        port,
		Weight:      float64(weight),
		Enable:      true,
		Healthy:     true,
		ServiceName: serviceName,
		GroupName:   group,
		Ephemeral:   true,
		Metadata: map[string]string{ // 自定义Nacos异常实例需要3s就剔除
			"preserved.heart.beat.interval": "2000",
			"preserved.heart.beat.timeout":  "6000",
			"preserved.ip.delete.timeout":   "6000",
		},
	})
	if err != nil {
		log.Panicln("[nacos][RegisterNaming]Register Naming Failed:", err)
	}
	return ok, err
}

// DeregisterInstance 注销实例
func DeregisterInstance(serviceName, group, ip string, port uint64) (bool, error) {
	ok, err := namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serviceName,
		GroupName:   group,
		Ephemeral:   true,
	})
	if err != nil {
		log.Panic("[nacos][DeregisterInstance]Deregister Naming Failed:", err)
	}
	return ok, err
}

// GetService 获取服务信息
func GetService(serviceName, group string) (model.Service, error) {
	service, err := namingClient.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
		GroupName:   group,
	})
	if err != nil {
		log.Println("[nacos][GetService]GetService Failed:", err)
	}
	return service, err
}

// SelectAllInstances 获取所有的实例列表
func SelectAllInstances(serviceName, group string) ([]model.Instance, error) {
	instance, err := namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: serviceName,
		GroupName:   group,
	})
	if err != nil {
		log.Println("[nacos][SelectAllInstances]Select All Instances Failed:", err)
	}
	return instance, err
}

// SelectInstances 获取实例列表
func SelectInstances(serviceName, group string, healthOnly bool) ([]model.Instance, error) {
	instance, err := namingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   group,
		HealthyOnly: healthOnly,
	})
	if err != nil {
		log.Println("[nacos][SelectInstances]Select Instances Failed:", err)
	}
	return instance, err
}

// SelectOneHealthyInstance 获取一个健康的实例（加权随机轮询）
func SelectOneHealthyInstance(serviceName, group string) (*model.Instance, error) {
	instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		GroupName:   group,
	})
	if err != nil {
		log.Println("[nacos][SelectOneHealthyInstance]Select One Healthy Instance Failed:", err)
	}
	return instance, err
}

// GetAllServicesInfo 获取服务名列表
func GetAllServicesInfo(nameSpace, group string, page, pageSize uint32) (model.ServiceList, error) {
	instance, err := namingClient.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		NameSpace: nameSpace,
		GroupName: group,
		PageNo:    page,
		PageSize:  pageSize,
	})
	if err != nil {
		log.Println("[nacos][GetAllServicesInfo]Get All Services Info Failed:", err)
	}
	return instance, err
}

// Subscribe 监听服务变化
func Subscribe(serviceName, group string, fn func(services []model.Instance, err error)) error {
	err := namingClient.Subscribe(&vo.SubscribeParam{
		ServiceName:       serviceName,
		GroupName:         group,
		SubscribeCallback: fn,
	})
	if err != nil {
		log.Println("[nacos][Subscribe]Subscribe Failed:", err)
	}
	return err
}

// Unsubscribe 取消监听服务变化
func Unsubscribe(serviceName, group string, fn func(services []model.Instance, err error)) error {
	err := namingClient.Unsubscribe(&vo.SubscribeParam{
		ServiceName:       serviceName,
		GroupName:         group,
		SubscribeCallback: fn,
	})
	if err != nil {
		log.Println("[nacos][Unsubscribe]Unsubscribe Failed:", err)
	}
	return err
}
