package logger

import (
	"fmt"
	"os"
	"time"
)

const (
	LevelInfo  = "INFO"
	LevelError = "ERROR"
	LevelWarn  = "WARN"
	LevelDebug = "DEBUG"
)

func Log(level, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s]  %s - %s\n", level, timestamp, message)
}

func Info(message string) {
	Log(LevelInfo, message)
}

func Error(message string) {
	Log(LevelError, message)
}

func Warn(message string) {
	Log(LevelWarn, message)
}

func Debug(message string) {
	Log(LevelDebug, message)
}

func Fatal(message string) {
	Log(LevelError, message)
	os.Exit(1)
}
