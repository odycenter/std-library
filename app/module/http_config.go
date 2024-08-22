package module

import (
	"bytes"
	"compress/gzip"
	"context"
	"embed"
	"github.com/beego/beego/v2/server/web"
	"io"
	"log"
	"log/slog"
	internal "std-library/app/internal/module"
	"std-library/app/internal/web"
	internalHttp "std-library/app/internal/web/http"
	internalSys "std-library/app/internal/web/sys"
	"std-library/app/web/beego"
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
	c.moduleContext.AddListenPort(c.server.HttpHost.Port)
	c.moduleContext.httpConfigAdded = true
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
	slog.Info("allow /_sys/api access", "cidrs", cidrs)
	c.apiController.AccessControl.Allow = internalHttp.NewIPv4Ranges(cidrs)
}

func loadAndCompressApiJson(env string, embedFS embed.FS) error {
	file, err := embedFS.Open("api.json")
	if err != nil {
		slog.Warn("api.json not found", "env", env)
		return err
	}
	defer file.Close()

	internalSys.ApiJsonContent, err = io.ReadAll(file)
	if err != nil {
		slog.Error("Error reading api.json", "env", env, "error", err)
		return err
	}

	gzippedContent, err := compressContent(internalSys.ApiJsonContent)
	if err != nil {
		slog.Error("Error compressing api.json", "env", env, "error", err)
		return err
	}

	internalSys.ApiJsonGzipped = gzippedContent

	slog.Info("api.json loaded and compressed successfully", "env", env)
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
