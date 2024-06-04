// Package pyroscope Pyroscope性能分析工具封装
package pyroscope

import (
	"fmt"
	"github.com/grafana/pyroscope-go"
	"runtime"
	"std-library/logs"
)

var p *pyroscope.Profiler

// Config 配置
type Config struct {
	ApplicationName string //当前应用的名称
	ServerAddress   string //服务地址
	AuthToken       string //授权token
	LogLevel        int    //日志等级
	OpenGMB         bool   //default false,开启会造成过大的性能消耗
}

// Start 使用传入配置启动pyroscope
func Start(cfg *Config) {
	if cfg == nil {
		fmt.Println("pyroscope not configuration")
		return
	}
	logs.Info("[pyroscope] start with address: %s", cfg.ServerAddress)
	if cfg.LogLevel == 0 {
		cfg.LogLevel = LevelDebug
	}
	var profileTypes = []pyroscope.ProfileType{
		pyroscope.ProfileCPU,
		pyroscope.ProfileAllocObjects,
		pyroscope.ProfileAllocSpace,
		pyroscope.ProfileInuseObjects,
		pyroscope.ProfileInuseSpace,
	}
	if cfg.OpenGMB {
		runtime.SetMutexProfileFraction(100) //sampling probability = 1%(1/rate)
		runtime.SetBlockProfileRate(100)     //sampling probability = 1%(1/rate)
		profileTypes = append(profileTypes,
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		)
	}

	var e error
	p, e = pyroscope.Start(pyroscope.Config{
		ApplicationName: cfg.ApplicationName, //beego.AppConfig.String("AppName")
		Tags:            nil,
		ServerAddress:   cfg.ServerAddress,
		AuthToken:       cfg.AuthToken,
		SampleRate:      5,
		Logger:          newLogger(cfg.LogLevel), //pyroscope.StandardLogger
		ProfileTypes:    profileTypes,
		DisableGCRuns:   false,
	})
	if e != nil {
		fmt.Println("pyroscope start failed:", e)
	}
}

// Stop 停止pyroscope
func Stop() {
	if p == nil {
		return
	}
	e := p.Stop()
	if e != nil {
		fmt.Println("pyroscope stop failed:", e)
	}
}
