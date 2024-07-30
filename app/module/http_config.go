package module

import (
	"bytes"
	"compress/gzip"
	"context"
	"embed"
	"github.com/beego/beego/v2/server/web"
	"io"
	"log"
	internal "std-library/app/internal/module"
	"std-library/app/internal/web"
	internalHttp "std-library/app/internal/web/http"
	internalSys "std-library/app/internal/web/sys"
	"std-library/app/web/beego"
	"std-library/logs"
)

type HTTPConfig struct {
	moduleContext *Context
	server        *beego.HTTPServer
	apiController *internalSys.APIController
}

func (c *HTTPConfig) Initialize(moduleContext *Context, name string) {
	c.moduleContext = moduleContext
	c.server = beego.NewHTTPServer()

	c.moduleContext.StartupHook.StartStage2 = append(c.moduleContext.StartupHook.StartStage2, c.server)

	c.moduleContext.ShutdownHook.Add(internal.STAGE_0, func(ctx context.Context, timeoutInMs int64) {
		c.server.Shutdown(ctx)
	})
	c.moduleContext.ShutdownHook.Add(internal.STAGE_1, func(ctx context.Context, timeoutInMs int64) {
		c.server.AwaitRequestCompletion(ctx, timeoutInMs)
	})
	c.moduleContext.ShutdownHook.Add(internal.STAGE_8, func(ctx context.Context, timeoutInMs int64) {
		c.server.AwaitTermination(ctx)
	})

	c.apiController = internalSys.NewAPIController()
	web.Handler("/_sys/api/*", c.apiController)
}

func (c *HTTPConfig) Validate() {
	if c.server.HttpHost == nil {
		log.Fatal("http listen is not configured, please configure first")
	}
}

func (c *HTTPConfig) GZip() {
	c.server.GZip()
}

func (c *HTTPConfig) Listen(listen string) {
	c.server.HttpHost = internal_web.Parse(listen)
}

func (c *HTTPConfig) ErrorWithOkStatus(val bool) {
	c.server.ErrorWithOkStatus(val)
}

func (c *HTTPConfig) CustomErrorResponseMessage(f func(code int, message string) map[string]interface{}) {
	c.server.CustomErrorResponseMessage(f)
}

func (c *HTTPConfig) APIContent(envFS *map[string]embed.FS) {
	for env, embedFS := range *envFS {
		if err := loadAndCompressApiJson(env, embedFS); err == nil {
			return
		}
	}
}

func (c *HTTPConfig) AllowAPI(cidrs []string) {
	logs.Info("allow /_sys/api access, cidrs=%v", cidrs)
	c.apiController.AccessControl.Allow = internalHttp.NewIPv4Ranges(cidrs)
}

func loadAndCompressApiJson(env string, embedFS embed.FS) error {
	file, err := embedFS.Open("api.json")
	if err != nil {
		logs.Warn("api.json not found in env: %s", env)
		return err
	}
	defer file.Close()

	internalSys.ApiJsonContent, err = io.ReadAll(file)
	if err != nil {
		logs.Error("Error reading api.json from env: %s: %v", env, err)
		return err
	}

	gzippedContent, err := compressContent(internalSys.ApiJsonContent)
	if err != nil {
		logs.Error("Error compressing api.json from env: %s: %v", env, err)
		return err
	}

	internalSys.ApiJsonGzipped = gzippedContent

	logs.Info("api.json loaded and compressed successfully from env: %s", env)
	return nil
}

func compressContent(content []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzWriter, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	defer gzWriter.Close()

	if _, err := gzWriter.Write(content); err != nil {
		return nil, err
	}

	if err := gzWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
