// Package logs logs操作统一封装
package logs

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"std-library/app/log/consts/logKey"
	"strings"
	"sync"
	"time"
)

// DefaultLog 默认创建的Logger
type DefaultLog struct {
	lock                sync.Mutex
	wg                  sync.WaitGroup
	level               LogLevel
	init                bool
	loggerFuncCallDepth int //如果值 <0 则不再打印，=0为
	asynchronous        bool
	msgChan             chan *Msg
	signalChan          chan string
	msgChanLen          int
	outputs             []*nameLogger
}

func (dl *DefaultLog) setLogger(adapterName Adapter, opts ...*Option) error {
	for _, output := range dl.outputs {
		if output.name == adapterName {
			return fmt.Errorf("logs: duplicate adaptername %q (you have set this Logger before)", adapterName)
		}
	}

	logAdapter, ok := adapters.Load(adapterName)
	if !ok {
		return fmt.Errorf("logs: unknown Adapter %q (forgotten Register?)", adapterName)
	}
	lg := logAdapter.(newLoggerFunc)()
	opts = append(opts, &Option{LogLevel: LevelDebug})
	err := lg.Init(opts[0])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "logs.defaultLogger.SetLogger: "+err.Error())
		return err
	}
	dl.outputs = append(dl.outputs, &nameLogger{
		name:   adapterName,
		Logger: lg,
	})
	return nil
}

// NewLogger 创建日志操作类
func NewLogger(asyncMsgLen ...int) *DefaultLog {
	dl := new(DefaultLog)
	dl.level = LevelDebug
	dl.loggerFuncCallDepth = 2
	dl.msgChanLen = append(asyncMsgLen, 0)[0]
	if dl.msgChanLen <= 0 {
		dl.msgChanLen = defaultAsyncMsgLen
	}
	dl.signalChan = make(chan string, 1)
	_ = dl.setLogger(AdapterConsole)
	return dl
}

// SetLogger 添加日志实现
func (dl *DefaultLog) SetLogger(adapterName Adapter, opts ...*Option) error {
	dl.lock.Lock()
	defer dl.lock.Unlock()
	if !dl.init {
		dl.outputs = []*nameLogger{}
		dl.init = true
	}
	return dl.setLogger(adapterName, opts...)
}

// Async 异步日志打印
func (dl *DefaultLog) Async(msgLen ...int) *DefaultLog {
	dl.lock.Lock()
	defer dl.lock.Unlock()
	if dl.asynchronous {
		return dl
	}
	dl.asynchronous = true
	if len(msgLen) > 0 && msgLen[0] > 0 {
		dl.msgChanLen = msgLen[0]
	}
	dl.msgChan = make(chan *Msg, dl.msgChanLen)
	msgPool = &sync.Pool{
		New: func() any {
			return &Msg{}
		},
	}
	dl.wg.Add(1)
	go dl.startLogger()
	return dl
}

// DelLogger 从 DefaultLogger 中删除一个adapter
func (dl *DefaultLog) DelLogger(adapterName Adapter) error {
	dl.lock.Lock()
	defer dl.lock.Unlock()
	var outputs []*nameLogger
	for _, lg := range dl.outputs {
		if lg.name == adapterName {
			lg.Destroy()
		} else {
			outputs = append(outputs, lg)
		}
	}
	if len(outputs) == len(dl.outputs) {
		return fmt.Errorf("logs: unknown adaptername %q (forgotten Register?)", adapterName)
	}
	dl.outputs = outputs
	return nil
}

func (dl *DefaultLog) writeToLoggers(lm *Msg) {
	for _, l := range dl.outputs {
		err := l.WriteMsg(lm)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "unable to WriteMsg to Adapter:%v,error:%v\n", l.name, err)
		}
	}
}

func (dl *DefaultLog) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	//补充 '\n'
	if p[len(p)-1] == '\n' {
		p = p[0 : len(p)-1]
	}
	lm := &Msg{
		Msg:   string(p),
		Level: LevelError,
		When:  time.Now(),
	}
	//设置 levelLoggerImpl 以确保所有日志消息都将被写出
	err = dl.writeMsg(lm)
	if err == nil {
		return len(p), nil
	}
	return 0, err
}

func (dl *DefaultLog) writeMsg(lm *Msg) error {
	if !dl.init {
		dl.lock.Lock()
		_ = dl.setLogger(AdapterConsole)
		dl.lock.Unlock()
	}
	var (
		file string
		line int
		ok   bool
	)
	_, file, line, ok = runtime.Caller(dl.loggerFuncCallDepth)
	if !ok {
		file = "???"
		line = 0
	}
	lm.FilePath = file
	lm.LineNumber = line

	lm.enableFuncCallDepth = dl.loggerFuncCallDepth > 0
	if dl.asynchronous {
		logM := msgPool.Get().(*Msg)
		logM.Level = lm.Level
		logM.Msg = lm.Msg
		logM.When = lm.When
		logM.Args = lm.Args
		logM.FilePath = lm.FilePath
		logM.LineNumber = lm.LineNumber
		if dl.outputs != nil {
			dl.msgChan <- lm
		} else {
			msgPool.Put(lm)
		}
	} else {
		dl.writeToLoggers(lm)
	}
	return nil
}

// SetLevel 设置日志消息级别
// 如果消息级别（如 LevelDebug）高于记录器级别（如 LevelWarning），日志提供者将不会发送消息
func (dl *DefaultLog) SetLevel(l LogLevel) {
	dl.level = l
}

// GetLevel 获取当前log的等级(int)
func (dl *DefaultLog) GetLevel() int {
	return int(dl.level)
}

// SetLogFuncCallDepth 设置日志 funcCallDepth
func (dl *DefaultLog) SetLogFuncCallDepth(d int) {
	dl.loggerFuncCallDepth = d
}

// GetLogFuncCallDepth 返回日志调用处 wrapper 的 funcCallDepth
func (dl *DefaultLog) GetLogFuncCallDepth() int {
	return dl.loggerFuncCallDepth
}

// 开始log消耗
func (dl *DefaultLog) startLogger() {
	gameOver := false
	for {
		select {
		case bm := <-dl.msgChan:
			dl.writeToLoggers(bm)
			msgPool.Put(bm)
		case sg := <-dl.signalChan:
			// Now should only send "flush" or "close" to dl.signalChan
			dl.flush()
			if sg == "close" {
				for _, l := range dl.outputs {
					l.Destroy()
				}
				dl.outputs = nil
				gameOver = true
			}
			dl.wg.Done()
		}
		if gameOver {
			break
		}
	}
}

// Error Log ERROR level 日志。
// Deprecated: Use Log instead
func (dl *DefaultLog) Error(format string, v ...any) {
	if LevelError > dl.level {
		return
	}
	lm := &Msg{
		Level: LevelError,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}

	_ = dl.writeMsg(lm)
}

// Warn Log WARN level 日志。
// Deprecated: Use Log instead
func (dl *DefaultLog) Warn(format string, v ...any) {
	if LevelWarning > dl.level {
		return
	}
	lm := &Msg{
		Level: LevelWarning,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}

	_ = dl.writeMsg(lm)
}

// Notice Log NOTICE level 日志。
// Deprecated: Use Log instead
func (dl *DefaultLog) Notice(format string, v ...any) {
	if LevelNotice > dl.level {
		return
	}
	lm := &Msg{
		Level: LevelNotice,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}

	_ = dl.writeMsg(lm)
}

// Info Log INFO level 日志。
// Deprecated: Use Log instead
func (dl *DefaultLog) Info(format string, v ...any) {
	if LevelInformation > dl.level {
		return
	}
	lm := &Msg{
		Level: LevelInformation,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}

	_ = dl.writeMsg(lm)
}

// Debug Log DEBUG level 日志。
// Deprecated: Use Log instead
func (dl *DefaultLog) Debug(format string, v ...any) {
	if LevelDebug > dl.level {
		return
	}
	lm := &Msg{
		Level: LevelDebug,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}

	_ = dl.writeMsg(lm)
}

func (dl *DefaultLog) Log(ctx context.Context, lv LogLevel, f any, v ...any) {
	if lv > dl.level {
		return
	}
	message := formatPattern(f, v...)

	var id string
	if ctx != nil {
		val := ctx.Value(logKey.Id)
		ctxId, ok := val.(string)
		if ok {
			id = ctxId
		}
	}

	lm := &Msg{
		Level: lv,
		Msg:   message,
		When:  time.Now(),
		Args:  v,
		ID:    id,
	}
	_ = dl.writeMsg(lm)
}

func (dl *DefaultLog) flush() {
	if dl.asynchronous {
		for {
			if len(dl.msgChan) > 0 {
				bm := <-dl.msgChan
				dl.writeToLoggers(bm)
				msgPool.Put(bm)
				continue
			}
			break
		}
	}
}

// Flush flush all chan data.
func (dl *DefaultLog) Flush() {
	if dl.asynchronous {
		dl.signalChan <- "flush"
		dl.wg.Wait()
		dl.wg.Add(1)
		return
	}
	dl.flush()
}

// Close 关闭logger，刷新所有通道数据并销毁 DefaultLogger 中的所有适配器。
func (dl *DefaultLog) Close() {
	if dl.asynchronous {
		dl.signalChan <- "close"
		dl.wg.Wait()
		close(dl.msgChan)
	} else {
		dl.flush()
		for _, l := range dl.outputs {
			l.Destroy()
		}
		dl.outputs = nil
	}
	close(dl.signalChan)
}

// Reset 关闭所有输出，并将 dl.outputs 设置为 nil
func (dl *DefaultLog) Reset() {
	dl.Flush()
	for _, l := range dl.outputs {
		l.Destroy()
	}
	dl.outputs = nil
}

// GetLogger 返回默认 Logger
func GetLogger(prefixes ...string) *log.Logger {
	prefix := append(prefixes, "")[0]
	if prefix != "" {
		prefix = fmt.Sprintf(`[%s] `, strings.ToUpper(prefix))
	}
	defaultLoggerMap.RLock()
	l, ok := defaultLoggerMap.logs[prefix]
	if ok {
		defaultLoggerMap.RUnlock()
		return l
	}
	defaultLoggerMap.RUnlock()
	defaultLoggerMap.Lock()
	defer defaultLoggerMap.Unlock()
	l, ok = defaultLoggerMap.logs[prefix]
	if !ok {
		l = log.New(defaultLogger, prefix, 0)
		defaultLoggerMap.logs[prefix] = l
	}
	return l
}

// GetDefaultLogger 返回默认 Logger
func GetDefaultLogger() *DefaultLog {
	return defaultLogger
}

// Reset 删除所有适配器
func Reset() {
	defaultLogger.Reset()
}

// Async 使用异步模式设置 defaultLogger 并保留 msglen 消息
func Async(msgLen ...int) *DefaultLog {
	return defaultLogger.Async(msgLen...)
}

// SetLevel 设置使用的全局日志级别
func SetLevel(l LogLevel) {
	defaultLogger.SetLevel(l)
}

// SetLogFuncCallDepth 设置日志 funcCallDepth
func SetLogFuncCallDepth(d int) {
	defaultLogger.loggerFuncCallDepth = d
}

// SetLogger 设置一个新的 logger.
func SetLogger(adapter Adapter, opts ...*Option) error {
	return defaultLogger.SetLogger(adapter, opts...)
}

// Deprecated: Use slog.Error instead
func Error(f any, v ...any) {
	defaultLogger.Log(nil, LevelError, f, v...)
}

// Deprecated: Use slog.Warn instead
func Warn(f any, v ...any) {
	defaultLogger.Log(nil, LevelWarning, f, v...)
}

// Deprecated: Use slog.Warn instead
func Notice(f any, v ...any) {
	defaultLogger.Log(nil, LevelNotice, f, v...)
}

// Deprecated: Use slog.Info instead
func Info(f any, v ...any) {
	defaultLogger.Log(nil, LevelInformation, f, v...)
}

// Deprecated: Use slog.Debug instead
func Debug(f any, v ...any) {
	defaultLogger.Log(nil, LevelDebug, f, v...)
}

// Deprecated: Use slog.DebugContext instead
func DebugWithCtx(ctx context.Context, f any, v ...any) {
	defaultLogger.Log(ctx, LevelDebug, f, v...)
}

// Deprecated: Use slog.InfoContext instead
func InfoWithCtx(ctx context.Context, f any, v ...any) {
	defaultLogger.Log(ctx, LevelInformation, f, v...)
}

// Deprecated: Use slog.WarnContext instead
func WarnWithCtx(ctx context.Context, f any, v ...any) {
	defaultLogger.Log(ctx, LevelWarning, f, v...)
}

// Deprecated: Use slog.ErrorContext instead
func ErrorWithCtx(ctx context.Context, f any, v ...any) {
	defaultLogger.Log(ctx, LevelError, f, v...)
}
