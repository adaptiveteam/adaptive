package core_utils_go

import (
	"log"
	"time"
)

// IfThenElse evaluates a condition, if true returns the first parameter otherwise the second
func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

// LogTimeConsumed logs the time consumed for a block of code
// It's typically used as 'defer LogTimeConsumed(time.Now(), "identifier")'
func LogTimeConsumed(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("TimeConsumed: %s took %s", name, elapsed)
}
