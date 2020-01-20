package issues

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

// recoverToErrorVar recovers and places the recovered error into the given variable
func recoverToErrorVar(name string, err *error) {
	if err2 := recover(); err2 != nil {
		switch err2.(type) {
		case error:
			err3 := err2.(error)
			err4 := errors.Wrap(err3, "Recover "+name+" from panic in issues workflow")
			err = &err4
		case string:
			err3 := err2.(string)
			err4 := errors.New("Recover " + name + " from string-panic in issues workflow: " + err3)
			err = &err4
		}
	}
}

func (w workflowImpl) recoverToErrorVar(name string, err *error) {
	if err2 := recover(); err2 != nil {
		if err != nil {
			w.AdaptiveLogger.WithError(*err).Errorln("Before recoverToErrorVar " + name)
		}
		switch err2.(type) {
		case error:
			err3 := err2.(error)
			err4 := fmt.Errorf("Recover %s from panic in workflow HandleRequest: %+v", name, err3)
			//err4 := errors.Wrap(err3, "Recover from panic in workflow HandleRequest")
			err = &err4
			w.AdaptiveLogger.WithError(err3).Errorln("recoverToErrorVar")
		case string:
			err3 := err2.(string)
			err4 := errors.New("Recover " + name + " from string-panic in workflow HandleRequest: " + err3)
			err = &err4
			w.AdaptiveLogger.WithError(err4).Errorln("recoverStringToErrorVar " + name)
		}
	}
}

// RecoverToLog recovers and places the recovered error into the given variable
func RecoverToLog(name string) {
	if err2 := recover(); err2 != nil {
		log.Printf(name+": recovered from %+v\n", err2)
	}
}
