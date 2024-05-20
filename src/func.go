package src

import (
	"fmt"
	"io"
	"os"
	"time"
)

func CheckPath(path string, isDir bool) bool {
	if file, err := os.Stat(path); err == nil {
		if isDir {
			return file.IsDir()
		} else {
			return true
		}
	}
	return false
}

/**
 * Logger Class that logs messages to a writer
 * Logs Format:
 * time=2021-07-07T15:04:05.999999999Z level=info package="chromium" source="" msg="This is a log message"
 */

var Log = NewLogger()

type Logger struct {
	writer io.Writer
	level  string // debug, info, warn, error
}

func (l *Logger) SetOutput(w io.Writer) {
	l.writer = w
}

func (l *Logger) SetLevel(level string) {
	l.level = level
}

func (l *Logger) Log(level string, pkg string, source string, msg string) {
	log := fmt.Sprintf("time=%s level=%s package=%s source=%s msg=%s\n", time.Now().Format(time.RFC3339), level, pkg, source, msg)
	if level == "debug" && l.level != "debug" {
		return
	} else if level == "info" && !(l.level == "info" || l.level == "debug") {
		return
	} else if level == "warn" && !(l.level == "warn" || l.level == "info" || l.level == "debug") {
		return
	}

	l.writer.Write([]byte(log))
}

func NewLogger() *Logger {
	return &Logger{}
}
