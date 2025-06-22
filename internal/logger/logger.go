package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func parseLogLevel(levelStr string) LogLevel {
	switch strings.ToUpper(strings.TrimSpace(levelStr)) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN", "WARNING":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return WARN // 默认返回 WARN 级别
	}
}

type Logger struct {
	level  LogLevel
	logger *log.Logger
}

func NewLogger() *Logger {
	HOME := os.Getenv("HOME")
	logDir := filepath.Join(HOME, ".config")
	logFile := filepath.Join(logDir, "zsh_yakumo.log")

	os.MkdirAll(logDir, 0755)

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		file = os.Stderr
	}

	// 从环境变量读取日志级别，默认为 WARN
	logLevelStr := os.Getenv("LOG_LEVEL")
	logLevel := parseLogLevel(logLevelStr)
	if logLevelStr == "" {
		logLevel = WARN // 明确设置默认级别为 WARN
	}

	logger := log.New(file, "", 0)
	return &Logger{
		level:  logLevel,
		logger: logger,
	}
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level >= l.level {
		timestamp := time.Now().Format("2006-01-02 15:04:05.000")
		message := fmt.Sprintf(format, args...)
		l.logger.Printf("[%s] %s: %s", timestamp, level.String(), message)
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}
