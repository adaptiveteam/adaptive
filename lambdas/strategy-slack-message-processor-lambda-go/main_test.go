package lambda

import (
	"reflect"
	"testing"
	"time"
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
	str := "06-30-2019"
	test, err := time.Parse(string(DateFormat), str)
	if err != nil {
		t.Fail()
	}
	logger.Info(test)
}

func TestTimeFormatChange(t *testing.T) {
	oldTimeStr := "2019-05-30T14:12:51Z"
	oldTimeFormat := TimestampFormat
	newTimeFormat := DateFormat
	res := timeFormatChange(oldTimeStr, oldTimeFormat, newTimeFormat)
	logger.Info(res)
}
