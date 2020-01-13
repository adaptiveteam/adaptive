package models

import (
	"encoding/json"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

// AnnualHolidayJSONUnmarshalUnsafe - unmarshalls AnnualHoliday
func AnnualHolidayJSONUnmarshalUnsafe(jsMessage string, namespace string) AnnualHoliday {
	res, err := AnnualHolidayJSONUnmarshal(jsMessage)
	core.ErrorHandler(err, namespace, "AnnualHoliday not unmarshaled from '" + jsMessage + "'")
	return res
}

// AnnualHolidayJSONUnmarshal - unmarshalls AnnualHoliday
func AnnualHolidayJSONUnmarshal(jsMessage string) (res AnnualHoliday, err error) {
	err = json.Unmarshal([]byte(jsMessage), &res)
	return
}

// ToJSON returns json string
func (l *AnnualHoliday) ToJSON() (string, error) {
	b, err := json.Marshal(&l)
	return string(b), err
}

// ToJSONUnsafe returns json string and panics in case of any errors
func (l *AnnualHoliday) ToJSONUnsafe(namespace string) string {
	str, err := l.ToJSON()
	core.ErrorHandler(err, namespace, "AnnualHoliday failed to marshal")
	return str
}



// AdHocHolidayJSONUnmarshalUnsafe - unmarshalls AdHocHoliday
func AdHocHolidayJSONUnmarshalUnsafe(jsMessage string, namespace string) AdHocHoliday {
	res, err := AdHocHolidayJSONUnmarshal(jsMessage)
	core.ErrorHandler(err, namespace, "AdHocHoliday not unmarshaled from '" + jsMessage + "'")
	return res
}

// AdHocHolidayJSONUnmarshal - unmarshalls AdHocHoliday
func AdHocHolidayJSONUnmarshal(jsMessage string) (res AdHocHoliday, err error) {
	err = json.Unmarshal([]byte(jsMessage), &res)
	return res, err
}

// // ToJSON returns json string
// func (l *AdHocHoliday) ToJSON() (string, error) {
// 	var b, err = json.Marshal(&l)
// 	return string(b), err
// }

// // ToJSONUnsafe returns json string and panics in case of any errors
// func (l *AdHocHoliday) ToJSONUnsafe(namespace string) string {
// 	var str, err = l.ToJSON()
// 	core.ErrorHandler(err, namespace, "AdHocHoliday failed to marshal")
// 	return str
// }
