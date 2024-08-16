package module

import (
	"github.com/beego/beego/v2/server/web"
	"log/slog"
	"os"
	app "std-library/app/conf"
	internalsys "std-library/app/internal/web/sys"
	"std-library/app/log"
)

type LogConfig struct {
	name          string
	moduleContext *Context
	handler       *log.Handler
}

func (c *LogConfig) Initialize(moduleContext *Context, name string) {
	c.name = name
	c.moduleContext = moduleContext

	c.handler = log.NewHandler(os.Stdout)
	logger := slog.New(c.handler)
	logger = logger.With(slog.String("app", app.Name))
	slog.SetDefault(logger)

	controller := internalsys.NewLogController(c.handler)
	web.Handler("/_sys/log", controller)
	web.Handler("/_sys/log/*", controller)
}

func (c *LogConfig) Validate() {

}

func (c *LogConfig) DefaultLevel(level string) {
	c.handler.SetDefaultLevel(level)
}

func (c *LogConfig) AppendToKafka(uri string) {
}

func (c *LogConfig) Appender() {
}

func (c *LogConfig) MaskedFields(fields ...string) {
}
