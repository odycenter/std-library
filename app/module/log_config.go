package module

import (
	"github.com/beego/beego/v2/server/web"
	app "github.com/odycenter/std-library/app/conf"
	internalLog "github.com/odycenter/std-library/app/internal/log"
	internalsys "github.com/odycenter/std-library/app/internal/web/sys"
	"log/slog"
	"os"
)

type LogConfig struct {
	name          string
	moduleContext *Context
	handler       *internalLog.Handler
}

func (c *LogConfig) Initialize(moduleContext *Context, name string) {
	c.name = name
	c.moduleContext = moduleContext

	c.handler = internalLog.NewHandler(os.Stdout)
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
	internalLog.AddMaskedField(fields...)
}
