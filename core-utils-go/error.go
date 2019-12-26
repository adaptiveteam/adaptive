package core_utils_go

import (
	"log"
)

// ErrorHandler is a universal handler that logs error and then panics unless err is nil
func ErrorHandler(err error, namespace, msg string) {
	if err != nil {
		log.Panicf("ERROR in %s : %s [%s]\n", namespace, msg, err.Error())
	}
}

// RecoverAsLogError is a universal error recovery that logs error and 
func RecoverAsLogError(label string) {
	err := recover()
	if err != nil {
		log.Printf("IGNORING ERROR in %s: %v", label, err)
	}
}
