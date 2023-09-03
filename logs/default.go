// Package logs logs操作统一封装
package logs

import (
	"fmt"
	"log"
	"os"
	"runtime"
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
	enableFullFilePath  bool
	asynchronous        bool
	msgChan             chan *Msg
	signalChan          chan string
	msgChanLen          int
	outputs             []*nameLogger
	globalFormatter     string
	prefix              string
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
	// 全局的格式化程序覆盖默认设置格式化程序
	if len(dl.globalFormatter) > 0 {
		ft, ok := GetFormatter(dl.globalFormatter)
		if !ok {
			return fmt.Errorf("the formatter with name: %s not found", dl.globalFormatter)
		}
		lg.SetFormatter(ft)
	}
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
		Level: LevelEmergency,
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

	lm.enableFullFilePath = dl.enableFullFilePath
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

// SetPrefix 设置前缀
func (dl *DefaultLog) SetPrefix(s string) {
	dl.prefix = s
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

func (dl *DefaultLog) setGlobalFormatter(formatter string) error {
	dl.globalFormatter = formatter
	return nil
}

// SetGlobalFormatter 为所有日志适配器设置全局格式化程序 不要忘记通过调用 RegisterFormatter 来注册格式化程序
func SetGlobalFormatter(formatter string) error {
	return defaultLogger.setGlobalFormatter(formatter)
}

// Emergency Log EMERGENCY level 日志。
func (dl *DefaultLog) Emergency(format string, v ...any) {
	if LevelEmergency > dl.level {
		return
	}

	lm := &Msg{
		Level: LevelEmergency,
		Msg:   format,
		When:  time.Now(),
	}
	if len(v) > 0 {
		lm.Msg = fmt.Sprintf(lm.Msg, v...)
	}

	_ = dl.writeMsg(lm)
}

// Alert Log ALERT level 日志。
func (dl *DefaultLog) Alert(format string, v ...any) {
	if LevelAlert > dl.level {
		return
	}

	lm := &Msg{
		Level: LevelAlert,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}
	_ = dl.writeMsg(lm)
}

// Critical Log CRITICAL level 日志。
func (dl *DefaultLog) Critical(format string, v ...any) {
	if LevelCritical > dl.level {
		return
	}
	lm := &Msg{
		Level: LevelCritical,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}

	_ = dl.writeMsg(lm)
}

// Error Log ERROR level 日志。
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

// ID Log参数填充ID
func (dl *DefaultLog) ID(lv LogLevel, ID string, format string, v ...any) {
	if lv > dl.level {
		return
	}
	lm := &Msg{
		Level: lv,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
		ID:    ID,
	}
	_ = dl.writeMsg(lm)
}

// Ex Log扩展参数填充，
// [0]Title string
// [1]ExecDur int64
// [2]ID string
func (dl *DefaultLog) Ex(lv LogLevel, ex map[string]any, format string, v ...any) {
	if lv > dl.level {
		return
	}
	lm := &Msg{
		Level: lv,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}
	if title, ok := ex["Title"]; ok {
		lm.Title = title.(string)
	}
	if execDur, ok := ex["ExecDur"]; ok {
		lm.ExecDur = int64(execDur.(int))
	}
	if id, ok := ex["ID"]; ok {
		lm.ID = id.(string)
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
	for _, l := range dl.outputs {
		l.Flush()
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

// EnableFullFilePath 启用完整文件路径日志记录。默认禁用
// e.g "/home/Documents/GitHub/beego/mainapp/" instead of "mainapp"
func EnableFullFilePath(b bool) {
	defaultLogger.enableFullFilePath = b
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

// SetPrefix 设置前缀
func SetPrefix(s string) {
	defaultLogger.SetPrefix(s)
}

// SetLogFuncCallDepth 设置日志 funcCallDepth
func SetLogFuncCallDepth(d int) {
	defaultLogger.loggerFuncCallDepth = d
}

// SetLogger 设置一个新的 logger.
func SetLogger(adapter Adapter, opts ...*Option) error {
	return defaultLogger.SetLogger(adapter, opts...)
}

// Emergency Log EMERGENCY level 日志。
func Emergency(f any, v ...any) {
	defaultLogger.Emergency(formatPattern(f, v...), v...)
}

// Alert Log ALERT level 日志。
func Alert(f any, v ...any) {
	defaultLogger.Alert(formatPattern(f, v...), v...)
}

// Critical Log CRITICAL level 日志。
func Critical(f any, v ...any) {
	defaultLogger.Critical(formatPattern(f, v...), v...)
}

// Error Log ERROR level 日志。
func Error(f any, v ...any) {
	defaultLogger.Error(formatPattern(f, v...), v...)
}

// Warn Log WARN level 日志。
func Warn(f any, v ...any) {
	defaultLogger.Warn(formatPattern(f, v...), v...)
}

// Notice Log NOTICE level 日志。
func Notice(f any, v ...any) {
	defaultLogger.Notice(formatPattern(f, v...), v...)
}

// Info Log INFO level 日志。
func Info(f any, v ...any) {
	defaultLogger.Info(formatPattern(f, v...), v...)
}

// Debug Log DEBUG level 日志。
func Debug(f any, v ...any) {
	defaultLogger.Debug(formatPattern(f, v...), v...)
}

// Ex Log Ex 额外参数填充的 level 日志。
// 目前EX可包含以下字段来作为扩展信息（注意字段类型）：
// Title	string
// ExecDur	int64
// ID		string
func Ex(lv LogLevel, ex map[string]any, f any, v ...any) {
	defaultLogger.Ex(lv, ex, formatPattern(f, v...), v...)
}

// ID Log ID 为log 填充ID
func ID(lv LogLevel, ID string, f any, v ...any) {
	defaultLogger.ID(lv, ID, formatPattern(f, v...), v...)
}
