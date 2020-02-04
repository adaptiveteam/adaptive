package tests

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

func getFieldString(i interface{}, field string) string {
	// Create a value for the slice.
	v := reflect.ValueOf(i)
	// Get the field of the slice element that we want to set.
	f := v.FieldByName(field)
	// Get value
	return f.String()
}

func TestAsKv(t *testing.T) {
	str := "2019-06-30"
	test, err := time.Parse(string(core.ISODateLayout), str)
	if err != nil {
		t.Error(err)
	}
	log.Printf("test=%s", test)
}

func TestTimeFormatChange(t *testing.T) {
	oldTimeStr := "2019-05-30T14:12:51Z"
	res := timeFormatChange(oldTimeStr, core.TimestampLayout, core.ISODateLayout)
	log.Printf("res=%s", res)
}

func timeFormatChange(str string, oldFormat, newFormat core.AdaptiveDateLayout) string {
	t, err := oldFormat.ChangeLayout(str, newFormat)
	core.ErrorHandler(err, "", fmt.Sprintf("Could not parse time %s using format %s", str, oldFormat))
	return t
}
