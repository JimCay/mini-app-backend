package log

import (
	"errors"
	"fmt"
	"github.com/gggwvg/logrotate"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
	"strings"
	"tg-backend/config"
	"time"
)

type Level uint32

var AllLevels = []Level{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

func Setup(conf config.LogConfig) error {

	logPath := conf.Path
	if logPath == "" {
		logPath = "logs/"
	}
	err := os.MkdirAll(logPath, os.ModePerm)
	if err != nil {
		return errors.New("mkdir failed")
	}

	now := time.Now()
	fileDate := now.Format("2006-01-02")
	filename := fmt.Sprintf("%s%s.log", logPath, fileDate)
	opts := []logrotate.Option{
		logrotate.File(filename),
	}

	if conf.Size == "" {
		conf.Size = "30m"
	}
	opts = append(opts, logrotate.RotateSize(conf.Size))
	logger, err := logrotate.NewLogger(opts...)
	if err != nil {
		return errors.New("new logger is error")
	}
	level := conf.Level
	logrus.SetLevel(logrus.Level(level))
	logrus.SetOutput(os.Stdout)
	logrus.SetOutput(logger)
	writers := []io.Writer{
		logger,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	logrus.SetOutput(fileAndStdoutWriter)
	logrus.SetFormatter(new(LogFormatter))
	logger.Close()
	return nil
}

type LogFormatter struct{}

// Format 格式化日志输出
func (s *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05.000")
	// 日志输出格式
	msg := fmt.Sprintf("[%s][%s][%s][RoutineID=%s]%s\n", strings.ToUpper(entry.Level.String()), timestamp, getFileLine(), getGoroutinedID(), entry.Message)
	return []byte(msg), nil
}

func getFileLine() string {
	fl := ""
	if _, file, line, ok := runtime.Caller(10); ok {
		i := strings.LastIndex(file, "/")
		fl = fmt.Sprintf("%s:%d", file[i:], line)
	}
	return fl
}

func getGoroutinedID() string {
	buf := [64]byte{}
	n := runtime.Stack(buf[:], false)
	stk := strings.TrimPrefix(string(buf[:n]), "goroutine")
	firstField := strings.Fields(stk)[0]
	return firstField
}

func Info(format string, a ...interface{}) {
	logrus.Infof(format, a...)
}

func Warn(format string, a ...interface{}) {
	logrus.Warnf(format, a...)
}

func Debug(format string, a ...interface{}) {
	logrus.Debugf(format, a...)
}

func Error(format string, a ...interface{}) {
	logrus.Errorf(format, a...)
}

func Fatal(format string, a ...interface{}) {
	logrus.Fatalf(format, a...)
}

func Panic(format string, a ...interface{}) {
	logrus.Panicf(format, a...)
}
