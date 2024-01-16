package logging

import (
	"fmt"
	"time"
)

type LogLevel int

const (
    DEBUG LogLevel = iota
    INFO 
	WARNING 
	ERROR 
)

func (level LogLevel) String() string {
	names := []string{"DEBUG", "INFO", "WARNING", "ERROR"}
	return names[level]
}

// ANSI color codes for terminal output
const (
	colorReset  = "\033[0m"
    colorCyan = "\033[96m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen   = "\033[92m"
)

type Logger struct {
	Prefix string
    MinLogLevel LogLevel
}

func NewLogger(prefix string, minLogLevel LogLevel) *Logger {
	return &Logger{Prefix: prefix, MinLogLevel: minLogLevel}
}

func (logger *Logger) log(level LogLevel, msg string) {
    if level >= logger.MinLogLevel  {
        var color string
        switch level {
        case DEBUG:
            color = colorCyan
        case INFO:
            color = colorGreen
        case WARNING:
            color = colorYellow
        case ERROR:
            color = colorRed
        }
        currentTime := time.Now().Format("2006-01-02 15:04:05")
        fmt.Printf("%s[%s]\t[%s]\t[%s]: %s%s\n", color, currentTime, logger.Prefix, level, msg, colorReset)
    }
}

func (logger *Logger) Debug(msg string) {
    logger.log(DEBUG, msg)
}

func (logger *Logger) Info(msg string) {
	logger.log(INFO, msg)
}

func (logger *Logger) Warn(msg string) {
	logger.log(WARNING, msg)
}

func (logger *Logger) Error(msg string) {
	logger.log(ERROR, msg)
}
