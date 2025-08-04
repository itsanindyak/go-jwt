package logger

import (
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

var base = log.New(os.Stdout, "", 0)

var (
	infoStyle    = color.New(color.FgCyan).SprintFunc()
	successStyle = color.New(color.FgGreen).SprintFunc()
	warnStyle    = color.New(color.FgYellow).SprintFunc()
	errorStyle   = color.New(color.FgRed).SprintFunc()
)

func timestamp() string {
	return "[" + time.Now().Format("2006-01-02 15:04:05") + "]"
}

func Info(msg string) {
	base.Println(timestamp(), infoStyle("ℹ️ INFO"), infoStyle(msg))
}

func Success(msg string) {
	base.Println(timestamp(), successStyle("✅ SUCCESS"), successStyle(msg))
}

func Warn(msg string) {
	base.Println(timestamp(), warnStyle("⚠️ WARN"), warnStyle(msg))
}

func Error(msg string) {
	base.Println(timestamp(), errorStyle("❌ ERROR"), errorStyle(msg))
}

func Fatal(msg string) {
	base.Println(timestamp(), errorStyle("💀 FATAL"), errorStyle(msg))
	os.Exit(1)
}
