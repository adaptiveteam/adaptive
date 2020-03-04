package core_utils_go

import (
	"log"
	"github.com/pkg/errors"
)

// ErrorHandler is a universal handler that logs error and then panics unless err is nil
func ErrorHandler(err error, namespace, msg string) {
	if err != nil {
		panic(errors.Wrapf(err, "ERROR in %s : %s\n", namespace, msg))
	}
}

// RecoverAsLogError is a universal error recovery that logs error and
func RecoverAsLogError(label string) {
	err2 := recover()
	if err2 != nil {
		log.Printf("IGNORING ERROR in %s: %+v", label, err2)
	}
}

// RecoverToErrorVar recovers and places the recovered error into the given variable
func RecoverToErrorVar(name string, err *error) {
	if err2 := recover(); err2 != nil {
		switch err2.(type) {
		case error:
			err3 := err2.(error)
			err4 := errors.Wrap(err3, "Recover "+name+" from panic")
			err = &err4
		case string:
			err3 := err2.(string)
			err4 := errors.New("Recover " + name + " from string-panic: " + err3)
			err = &err4
		}
	}
}
