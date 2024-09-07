package nacos_test

import (
	"fmt"
	"github.com/odycenter/std-library/nacos"
	"testing"
	"time"
)

func TestGetConfigClient(t *testing.T) {
	nacos.Init(&nacos.ClientConfig{
		NamespaceId: "",
		Endpoint:    "127.0.0.1:8848",
	}, &nacos.ServerConfig{
		IpAddr: "127.0.0.1",
		Port:   8848,
		Scheme: "http",
	})
	dynamicDataId := fmt.Sprint("publish", time.Now().Unix())
	err := nacos.ListenConfig(dynamicDataId, "test", func(namespace, group, dataId, data string) {
		fmt.Println("changed:", namespace, group, dataId, data)
	})
	if err != nil {
		return
	}
	nacos.GetSrvConfig("game_system_config.go", "cloud")
	fmt.Println(nacos.PublishConfig(dynamicDataId, "cloud", "test1 = 11\ntest2 = AAA"))
	fmt.Println(nacos.PublishConfig(dynamicDataId, "test", "test1 = 22\ntest2 = BBB"))
	fmt.Println(nacos.DeleteConfig(dynamicDataId, "cloud"))
	fmt.Println(nacos.PublishConfig(dynamicDataId, "cloud", "test1 = 33\ntest2 = CCA"))
	fmt.Println(nacos.CancelListenConfig(dynamicDataId, "test"))
	fmt.Println(nacos.PublishConfig(dynamicDataId, "test", "test1 = 44\ntest2 = DDD"))
	for range time.Tick(time.Second) {

	}
}
