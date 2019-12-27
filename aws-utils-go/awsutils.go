package aws_utils_go

import (
	"fmt"
	"reflect"
	"time"
)

// debug print if enables.
func print(obj fmt.Stringer, debug bool) {
	if !debug {
		return
	}
	var name string
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		name = "*" + t.Elem().Name()
	} else {
		name = t.Name()
	}
	fmt.Printf("[DEBUG] %s\n", name)
	fmt.Println(obj)
}

// generateStatementId() generated unique statement string.
func generateStatementId(sType string) string {
	return fmt.Sprintf("adaptive-statement-%s-%d", sType, time.Now().UnixNano())
}
