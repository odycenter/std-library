package demo_test

import (
	"context"
	"github.com/beego/beego/v2/server/web"
	"github.com/odycenter/std-library/app/async"
	"github.com/odycenter/std-library/app/module"
	"github.com/odycenter/std-library/dbase"
	"github.com/odycenter/std-library/demo"
	"github.com/odycenter/std-library/demo/test"
	"github.com/odycenter/std-library/mongo"
	"github.com/odycenter/std-library/redis"
	"google.golang.org/grpc"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"testing"
	"time"
)

import (
	"embed"
)

//go:embed sys.properties api.json
var Conf embed.FS

func TestAppStart(t *testing.T) {
	os.Setenv("SHUTDOWN_DELAY_IN_SEC", "10")
	os.Setenv("SHUTDOWN_TIMEOUT_IN_SEC", "45")

	module.Start(&CoreApp{})
}

type CoreApp struct {
	module.App
}

func (a *CoreApp) Initialize() {
	envProperties := map[string]embed.FS{
		"": Conf,
	}
	a.Load(&module.SystemModule{EnvProperties: envProperties})
	a.Log().MaskedFields("Password")
	a.DB().ForceEarlyStart()
	var data time.Time
	e := dbase.Orm().Raw("select now() from dual;").QueryRow(&data)
	slog.Info("orm", "result", data, "err", e)

	a.DB("default@tidb").
		Url("@tcp(10.15.39.64:4000)/cloud").
		User("dev-tidb-rd-use").
		Password("wpI7ENNysS").
		ForceEarlyStart()
	e = dbase.Orm("default@tidb").Raw("select now() from dual;").QueryRow(&data)
	slog.Info("orm", "result", data, "err", e)

	a.Mongo().SlowOperationThreshold(40 * time.Millisecond)
	a.Mongo().ForceEarlyStart()
	result, err := mongo.DB().InsertOne("test", "InsertTest", map[string]interface{}{"name": "name"})
	if err != nil {
		panic(err)
	}
	slog.Info("after insert", "id", result.InsertedID)

	readOnly := a.Redis("read-only")
	readOnly.Host(a.RequiredProperty("sys.redis.readonly.host"))
	readOnly.ForceEarlyStart()
	redis.RDB("redis:read-only").Get("app_test.cachetest:rwww")
	a.Redis().ForceEarlyStart() // use ForceEarlyStart to aquire redis instance before startup stage
	redis.RDB().Del("app_test.cachetest:rwww")
	// a.Pyroscope().ForceLocalStart()
	aa := async.New("test", 10)
	aa.Submit(nil, "test", func(ctx context.Context) {
		slog.InfoContext(ctx, "test:%s", "player", "asa")
	})

	a.Load(&demo.HookModule{})
	// a.Load(&demo.KafkaModule{})
	a.Load(&demo.ScheduleModule{})
	a.Load(&demo.CacheModule{})
	a.Metric().Listen("8001")
	a.Grpc().MaxConnections(1)
	server := a.Grpc().AddOpt(grpc.MaxRecvMsgSize(4 << 30)).Server()
	test.RegisterTestServiceServer(server, &demo.HelloController{})
	a.Http()
	// _ = a.Cache().Add(demo.CacheTest{}, time.Second*300)  // enable to get "found duplicate cache name" panic on startup

	web.Handler("/sleep", &sleep50SHandler{})
}

type sleep50SHandler struct {
}

func (h *sleep50SHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("r: " + r.Method + ":" + r.URL.Path)
	time.Sleep(10 * time.Second)
}
