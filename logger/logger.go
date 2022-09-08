package logger

import (
	"fmt"
	"os"
	"time"
)

type Logger struct{}

var errorFormat = "[%s] ERROR: %s %s\n"
var infoFormat = "[%s] INFO: %s\n"

func (l *Logger) LogError(msg string, err error) {
	fmt.Fprintf(os.Stderr, errorFormat, time.Now().Format("2006-01-02 15:04:05"), msg, err.Error())
}

func (l *Logger) LogInfo(msg string) {
	fmt.Fprintf(os.Stderr, infoFormat, time.Now().Format("2006-01-02 15:04:05"), msg)
}
