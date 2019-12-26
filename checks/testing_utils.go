package checks

import (
business_time "github.com/adaptiveteam/business-time"
)

var SimpleTestProfile = CheckFunctionMap{
	"ReturnsTrue":  ReturnsTrue,
	"ReturnsFalse": ReturnsFalse,
}

func ReturnsTrue(_ string, _ business_time.Date) (rv bool) {
	rv = true
	return rv
}

func ReturnsFalse(_ string, _ business_time.Date) (rv bool) {
	rv = false
	return rv
}

