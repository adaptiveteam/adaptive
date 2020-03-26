package core_utils_go

import (
	"fmt"
	"log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var zlogger = func()*zap.Logger{
	res, err2 := zap.NewProduction()
	ErrorHandler(err2, "zap", "logger initialization")
	return res
}()
// ErrorHandler is a universal handler that logs error and then panics unless err is nil
func ErrorHandler(err error, namespace, msg string) {
	if err != nil {
		panic(errors.Wrapf(err, "ERROR in '%s': %s\n", namespace, msg))
	}
}

// RecoverAsFalse recovers from panic and returns false
func RecoverAsFalse(name string, res *bool) {
	if err2 := recover(); err2 != nil {
		switch err2.(type){
		case error:
			zlogger.Error(name, zap.Error(err2.(error)))
		default:
		}
		log.Printf("IGNORING ERROR in %s (returning false): %+v", name, err2)
		*res = false
	}
	return
} 

// RecoverAsLogError is a universal error recovery that logs error and
func RecoverAsLogError(label string) {
	if err2 := recover(); err2 != nil {
		switch err2.(type){
		case error:
			zlogger.Error(label, zap.Error(err2.(error)))
		default:
		}
		log.Printf("IGNORING ERROR in %s: %+v\n", label, err2)
	}
}

// RecoverAsLogErrorf is a universal error recovery that logs error and
func RecoverAsLogErrorf(format string, args ... interface{}) {
	RecoverAsLogError(fmt.Sprintf(format, args...))
}

// RecoverToErrorVar recovers and places the recovered error into the given variable
func RecoverToErrorVar(name string, err *error) {
	err2 := recover()
	if err2 != nil {
		log.Printf("RecoverToErrorVar2 (%s) (err=%+v), (err2: %+v\n", name, *err, err2)
		switch err2.(type) {
		case error:
			err3 := err2.(error)
			err4 := errors.Wrapf(err3, "%s: Recover from panic", name)
			*err = err4
		case string:
			err3 := err2.(string)
			err4 := errors.New(name + ": Recover from string-panic: " + err3)
			*err = err4
		default:
			err4 := errors.New(fmt.Sprintf("%s: Recover from unknown-panic: %+v", name, err2))
			*err = err4
		}
	}
}
