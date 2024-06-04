package demo_test

import (
	"github.com/beego/beego/v2/server/web"
	"google.golang.org/grpc"
	"net/http"
	_ "net/http/pprof"
	"os"
	"std-library/app/module"
	"std-library/demo"
	"std-library/logs"
	"std-library/mongo"
	"std-library/redis"
	"testing"
	"time"
)

import (
	"embed"
)

//go:embed sys.properties
var Conf embed.FS

func TestAppStart(t *testing.T) {
	os.Setenv("SHUTDOWN_DELAY_IN_SEC", "10")
	os.Setenv("SHUTDOWN_TIMEOUT_IN_SEC", "45")

	app := CoreApp{}
	app.Configure()
	app.Initialize()
	app.Start()
}

type CoreApp struct {
	module.App
}

func (a *CoreApp) Initialize() {
	envProperties := map[string]embed.FS{
		"": Conf,
	}
	a.Load(&module.SystemModule{EnvProperties: envProperties})
	a.Mongo().ForceEarlyStart()
	result, err := mongo.DB().InsertOne("test", "ChrisTest", map[string]interface{}{"name": "Chris"})
	if err != nil {
		panic(err)
	}
	logs.Info(result.InsertedID)

	a.Redis().ForceEarlyStart() // use ForceEarlyStart to aquire redis instance before startup stage
	redis.RDB().Del("app_test.cachetest:rwww")
	// a.Pyroscope().ForceLocalStart()

	a.Load(&demo.HookModule{})
	//a.Load(&demo.KafkaModule{})
	//a.Load(&demo.ScheduleModule{})
	a.Load(&demo.CacheModule{})

	a.Grpc().AddOpt(grpc.MaxRecvMsgSize(4 << 30)).Server()
	a.Http()
	// _ = a.Cache().Add(demo.CacheTest{}, time.Second*300)  // enable to get "found duplicate cache name" panic on startup

	web.Handler("/sleep", &sleep50SHandler{})
}

type sleep50SHandler struct {
}

func (h *sleep50SHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logs.Info("r: " + r.Method + ":" + r.URL.Path)
	time.Sleep(50 * time.Second)
}
