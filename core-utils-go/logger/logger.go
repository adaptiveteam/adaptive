package logger

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"strings"
)

// Inspired from here: https://github.com/ysugimoto/ginger/blob/master/logger/logger.go

var stdout = colorable.NewColorableStdout()

// Logger is the struct that outputs colored log with namespace.
type Logger struct {
	ns string
}

func WithNamespace(ns string) *Logger {
	return &Logger{
		ns: ns,
	}
}

// AddNamespace() adds more namespace string on output log.
func (l *Logger) AddNamespace(ns string) {
	l.ns += "." + ns
}
func (l *Logger) RemoveNamespace(ns string) {
	index := strings.Index(l.ns, "."+ns)
	if index != -1 {
		l.ns = l.ns[0:index]
	}
}

// Info() outputs information log with green color.
func (l *Logger) Info(message ...interface{}) {
	fmt.Fprintln(stdout, Green("INFO:["+l.ns+"] "+fmt.Sprint(message...)))
}

// Infof() outputs formatted information log with green color.
func (l *Logger) Infof(format string, args ...interface{}) {
	fmt.Fprintf(stdout, Green("INFO:["+l.ns+"] "+format), args...)
}

// Warn() outputs warning log with yellow color.
func (l *Logger) Warn(message ...interface{}) {
	fmt.Fprintln(stdout, Yellow("WARN:["+l.ns+"] "+fmt.Sprint(message...)))
}

// Warnf() outputs formatted warning log with yellow color.
func (l *Logger) Warnf(format string, args ...interface{}) {
	fmt.Fprintf(stdout, Yellow("WARN:["+l.ns+"] "+format), args...)
}

// Error() outputs error log with red color.
func (l *Logger) Error(message ...interface{}) {
	fmt.Fprintln(stdout, Red("ERROR:["+l.ns+"] "+fmt.Sprint(message...)))
}

// Errorf() outputs formatted error log with red color.
func (l *Logger) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(stdout, Red("ERROR:["+l.ns+"] "+format), args...)
}

// Print() outputs log with default color.
func (l *Logger) Print(message ...interface{}) {
	fmt.Fprintln(stdout, "["+l.ns+"] "+fmt.Sprint(message...))
}

// Printf() outputs formatted log with default color.
func (l *Logger) Printf(format string, args ...interface{}) {
	fmt.Fprintf(stdout, "["+l.ns+"] "+format, args...)
}
