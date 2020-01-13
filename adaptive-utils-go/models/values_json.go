package models

import (
	"encoding/json"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

// AdaptiveValueJSONUnmarshalUnsafe - unmarshalls AdaptiveValue
func AdaptiveValueJSONUnmarshalUnsafe(jsMessage string, namespace string) AdaptiveValue {
	res, err := AdaptiveValueJSONUnmarshal(jsMessage)
	core.ErrorHandler(err, namespace, "AdaptiveValue not unmarshaled from '" + jsMessage + "'")
	return res
}

// AdaptiveValueJSONUnmarshal - unmarshalls AdaptiveValue
func AdaptiveValueJSONUnmarshal(jsMessage string) (res AdaptiveValue, err error) {
	err = json.Unmarshal([]byte(jsMessage), &res)
	return res, err
}

// // ToJSON returns json string
// func (l *AdaptiveValue) ToJSON() (string, error) {
// 	b, err := json.Marshal(&l)
// 	return string(b), err
// }

// // ToJSONUnsafe returns json string and panics in case of any errors
// func (l *AdaptiveValue) ToJSONUnsafe(namespace string) string {
// 	str, err := l.ToJSON()
// 	core.ErrorHandler(err, namespace, "AdaptiveValue failed to marshal")
// 	return str
// }

