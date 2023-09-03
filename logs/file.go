package logs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// fileLogWriter implements LoggerInterface.
// 按行数限制、文件大小限制或时间频率写入消息。
type fileLogWriter struct {
	sync.RWMutex // 按顺序和原子增量写入日志 maxLinesCurLines 和 maxSizeCurSize
	Rotate       bool
	Daily        bool
	Hourly       bool
	// The opened file
	Filename   string
	fileWriter *os.File

	// Rotate at line
	MaxLines         int
	maxLinesCurLines int

	MaxFiles         int
	MaxFilesCurFiles int

	// Rotate at size
	MaxSize        int
	maxSizeCurSize int

	// Rotate daily
	MaxDays       int64
	dailyOpenDate int
	dailyOpenTime time.Time

	// Rotate hourly
	MaxHours       int64
	hourlyOpenDate int
	hourlyOpenTime time.Time

	Level LogLevel
	// Permissions for log file
	Perm string
	// Permissions for directory if it is specified in FileName
	DirPerm string

	RotatePerm string

	fileNameOnly, suffix string // like "project.log", project is fileNameOnly and .log is suffix

	logFormatter LogFormatter
	Formatter    string
}

func (w *fileLogWriter) getRotatePerm() string {
	if w.RotatePerm == "" {
		return "0660"
	}
	return w.RotatePerm
}

func (w *fileLogWriter) getPerm() string {
	if w.Perm == "" {
		return "0660"
	}
	return w.Perm
}

func (w *fileLogWriter) getDirPerm() string {
	if w.DirPerm == "" {
		return "0660"
	}
	return w.DirPerm
}

// newFileWriter 创建一个作为 LoggerInterface 返回的 FileLogWriter
func newFileWriter() Logger {
	w := &fileLogWriter{
		Daily:      true,
		MaxDays:    7,
		Hourly:     false,
		MaxHours:   168,
		Rotate:     true,
		RotatePerm: "0440",
		Level:      LevelDebug,
		Perm:       "0660",
		DirPerm:    "0770",
		MaxLines:   10000000,
		MaxFiles:   999,
		MaxSize:    1 << 28,
	}
	w.logFormatter = w
	return w
}

func (*fileLogWriter) Format(lm *Msg) string {
	return lm.Format()
}

func (w *fileLogWriter) SetFormatter(f LogFormatter) {
	w.logFormatter = f
}

// Init 初始化fileLog
func (w *fileLogWriter) Init(opt *Option) error {
	if opt == nil {
		return nil
	}
	w.Rotate = opt.Rotate
	w.Daily = opt.Daily
	w.Hourly = opt.Hourly
	w.Filename = opt.Filename
	w.MaxLines = opt.MaxLines
	w.MaxFiles = opt.MaxFiles
	w.MaxSize = opt.MaxSize
	w.MaxDays = opt.MaxDays
	w.MaxHours = opt.MaxHours
	w.Perm = opt.Perm
	w.Level = opt.LogLevel
	w.Formatter = opt.Formatter
	if w.Filename == "" {
		return errors.New("must have filename")
	}
	w.suffix = filepath.Ext(w.Filename)
	w.fileNameOnly = strings.TrimSuffix(w.Filename, w.suffix)
	if w.suffix == "" {
		w.suffix = ".log"
	}

	if len(w.Formatter) > 0 {
		fmtr, ok := GetFormatter(w.Formatter)
		if !ok {
			return fmt.Errorf("the formatter with name: %s not found", w.Formatter)
		}
		w.logFormatter = fmtr
	}
	err := w.startLogger()
	return err
}

// 启动文件记录器。创建日志文件并设置为 locker-inside 文件编写器。
func (w *fileLogWriter) startLogger() error {
	file, err := w.createLogFile()
	if err != nil {
		return err
	}
	if w.fileWriter != nil {
		_ = w.fileWriter.Close()
	}
	w.fileWriter = file
	return w.initFd()
}

func (w *fileLogWriter) needRotateDaily(day int) bool {
	return (w.MaxLines > 0 && w.maxLinesCurLines >= w.MaxLines) ||
		(w.MaxSize > 0 && w.maxSizeCurSize >= w.MaxSize) ||
		(w.Daily && day != w.dailyOpenDate)
}

func (w *fileLogWriter) needRotateHourly(hour int) bool {
	return (w.MaxLines > 0 && w.maxLinesCurLines >= w.MaxLines) ||
		(w.MaxSize > 0 && w.maxSizeCurSize >= w.MaxSize) ||
		(w.Hourly && hour != w.hourlyOpenDate)
}

// WriteMsg 将log写入文件
func (w *fileLogWriter) WriteMsg(lm *Msg) error {
	if lm.Level > w.Level {
		return nil
	}

	_, d, h := formatTimeHeader(lm.When)

	msg := w.logFormatter.Format(lm) + "\n"
	if w.Rotate {
		w.RLock()
		if w.needRotateHourly(h) {
			w.RUnlock()
			w.Lock()
			if w.needRotateHourly(h) {
				if err := w.doRotate(lm.When); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.Filename, err)
				}
			}
			w.Unlock()
		} else if w.needRotateDaily(d) {
			w.RUnlock()
			w.Lock()
			if w.needRotateDaily(d) {
				if err := w.doRotate(lm.When); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.Filename, err)
				}
			}
			w.Unlock()
		} else {
			w.RUnlock()
		}
	}

	w.Lock()
	_, err := w.fileWriter.Write([]byte(msg))
	if err == nil {
		w.maxLinesCurLines++
		w.maxSizeCurSize += len(msg)
	}
	w.Unlock()
	return err
}

func (w *fileLogWriter) createLogFile() (*os.File, error) {
	// Open the log file
	perm, err := strconv.ParseInt(w.getPerm(), 8, 64)
	if err != nil {
		return nil, err
	}

	dirPerm, err := strconv.ParseInt(w.getDirPerm(), 8, 64)
	if err != nil {
		return nil, err
	}

	fPath := path.Dir(w.Filename)
	_ = os.MkdirAll(fPath, os.FileMode(dirPerm))

	fd, err := os.OpenFile(w.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if err == nil {
		// Make sure file perm is user set perm cause of `os.OpenFile` will obey umask
		_ = os.Chmod(w.Filename, os.FileMode(perm))
	}
	return fd, err
}

func (w *fileLogWriter) initFd() error {
	fd := w.fileWriter
	fInfo, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s", err)
	}
	w.maxSizeCurSize = int(fInfo.Size())
	w.dailyOpenTime = time.Now()
	w.dailyOpenDate = w.dailyOpenTime.Day()
	w.hourlyOpenTime = time.Now()
	w.hourlyOpenDate = w.hourlyOpenTime.Hour()
	w.maxLinesCurLines = 0
	if w.Hourly {
		go w.hourlyRotate(w.hourlyOpenTime)
	} else if w.Daily {
		go w.dailyRotate(w.dailyOpenTime)
	}
	if fInfo.Size() > 0 && w.MaxLines > 0 {
		count, err := w.lines()
		if err != nil {
			return err
		}
		w.maxLinesCurLines = count
	}
	return nil
}

func (w *fileLogWriter) dailyRotate(openTime time.Time) {
	y, m, d := openTime.Add(24 * time.Hour).Date()
	nextDay := time.Date(y, m, d, 0, 0, 0, 0, openTime.Location())
	tm := time.NewTimer(time.Duration(nextDay.UnixNano() - openTime.UnixNano() + 100))
	<-tm.C
	w.Lock()
	if w.needRotateDaily(time.Now().Day()) {
		if err := w.doRotate(time.Now()); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.Filename, err)
		}
	}
	w.Unlock()
}

func (w *fileLogWriter) hourlyRotate(openTime time.Time) {
	y, m, d := openTime.Add(1 * time.Hour).Date()
	h, _, _ := openTime.Add(1 * time.Hour).Clock()
	nextHour := time.Date(y, m, d, h, 0, 0, 0, openTime.Location())
	tm := time.NewTimer(time.Duration(nextHour.UnixNano() - openTime.UnixNano() + 100))
	<-tm.C
	w.Lock()
	if w.needRotateHourly(time.Now().Hour()) {
		if err := w.doRotate(time.Now()); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.Filename, err)
		}
	}
	w.Unlock()
}

func (w *fileLogWriter) lines() (int, error) {
	fd, err := os.Open(w.Filename)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	buf := make([]byte, 32768) // 32k
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := fd.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

// DoRotate 意味着它需要将日志写入一个新文件。新文件名，
// 如 xx.2013-01-01.log（每日）或 xx.001.log（按行或大小）
func (w *fileLogWriter) doRotate(logTime time.Time) error {
	// file exists
	// Find the next available number
	num := w.MaxFilesCurFiles + 1
	fName := ""
	format := ""
	var openTime time.Time
	rotatePerm, err := strconv.ParseInt(w.RotatePerm, 8, 64)
	if err != nil {
		return err
	}

	_, err = os.Lstat(w.Filename)
	if err != nil {
		// even if the file is not exist or other ,we should RESTART the logger
		goto RestartLogger
	}

	if w.Hourly {
		format = "2006010215"
		openTime = w.hourlyOpenTime
	} else if w.Daily {
		format = "2006-01-02"
		openTime = w.dailyOpenTime
	}

	// only when one of them be setted, then the file would be splited
	if w.MaxLines > 0 || w.MaxSize > 0 {
		for ; err == nil && num <= w.MaxFiles; num++ {
			fName = w.fileNameOnly + fmt.Sprintf(".%s.%03d%s", logTime.Format(format), num, w.suffix)
			_, err = os.Lstat(fName)
		}
	} else {
		fName = w.fileNameOnly + fmt.Sprintf(".%s.%03d%s", openTime.Format(format), num, w.suffix)
		_, err = os.Lstat(fName)
		w.MaxFilesCurFiles = num
	}

	// return error if the last file checked still existed
	if err == nil {
		return fmt.Errorf("rotate: Cannot find free log number to rename %s", w.Filename)
	}

	// close fileWriter before rename
	_ = w.fileWriter.Close()

	// Rename the file to its new found name
	// even if occurs error,we MUST guarantee to  restart new logger
	err = os.Rename(w.Filename, fName)
	if err != nil {
		goto RestartLogger
	}

	err = os.Chmod(fName, os.FileMode(rotatePerm))

RestartLogger:

	startLoggerErr := w.startLogger()
	go w.deleteOldLog()

	if startLoggerErr != nil {
		return fmt.Errorf("rotate StartLogger: %s", startLoggerErr)
	}
	if err != nil {
		return fmt.Errorf("rotate: %s", err)
	}
	return nil
}

func (w *fileLogWriter) deleteOldLog() {
	dir := filepath.Dir(w.Filename)
	absolutePath, err := filepath.EvalSymlinks(w.Filename)
	if err == nil {
		dir = filepath.Dir(absolutePath)
	}
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Unable to delete old log '%s', error: %v\n", path, r)
			}
		}()

		if info == nil {
			return
		}
		if w.Hourly {
			if !info.IsDir() && info.ModTime().Add(1*time.Hour*time.Duration(w.MaxHours)).Before(time.Now()) {
				if strings.HasPrefix(filepath.Base(path), filepath.Base(w.fileNameOnly)) &&
					strings.HasSuffix(filepath.Base(path), w.suffix) {
					_ = os.Remove(path)
				}
			}
		} else if w.Daily {
			if !info.IsDir() && info.ModTime().Add(24*time.Hour*time.Duration(w.MaxDays)).Before(time.Now()) {
				if strings.HasPrefix(filepath.Base(path), filepath.Base(w.fileNameOnly)) &&
					strings.HasSuffix(filepath.Base(path), w.suffix) {
					_ = os.Remove(path)
				}
			}
		}
		return
	})
}

// Destroy 关闭文件
func (w *fileLogWriter) Destroy() {
	_ = w.fileWriter.Close()
}

// Flush 刷新文件记录器。
// 内存中的文件记录器中没有缓冲消息。
// 刷新文件意味着从磁盘同步文件。
func (w *fileLogWriter) Flush() {
	_ = w.fileWriter.Sync()
}

func init() {
	Register(AdapterFile, newFileWriter)
}
