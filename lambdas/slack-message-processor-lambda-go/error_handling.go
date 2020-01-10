package lambda

import (
	"github.com/pkg/errors"
)

// recoverToErrorVar recovers and places the recovered error into the given variable
func recoverToErrorVar(name string, err *error) {
	if err2 := recover(); err2 != nil {
		switch err2.(type) {
		case error:
			err3 := err2.(error)
			err4 := errors.Wrap(err3, "Recover " + name + " from panic in issues workflow")
			err = &err4
		case string:
			err3 := err2.(string)
			err4 := errors.New("Recover " + name + " from string-panic in issues workflow: "+err3)
			err = &err4
		}
	}
}
