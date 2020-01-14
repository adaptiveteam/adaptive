package models

import (
	"github.com/adaptiveteam/adaptive/daos/adaptiveValue"
)

// AdaptiveValue is a value for a client (Reliability, Skill, Contribution, and Productivity)
type AdaptiveValue = adaptiveValue.AdaptiveValue
// struct {
// 	ID            string `json:"id"`
// 	PlatformID    string `json:"platform_id"`
// 	Name          string `json:"value_name"`
// 	ValueType     string `json:"value_type"`
// 	Description   string `json:"description"`
// 	DeactivatedOn string `json:"deactivated_on,omitempty"` // if empty, then it's not deleted/deactivated, otherwise contains timestamp when it was deactivated
// }

var (
	// DefaultAdaptiveValueLevels default set of possible answers to value level
	DefaultAdaptiveValueLevels = map[int]string{
		5: "Exceeds",
		4: "Meets",
		3: "Approaching",
		2: "Below",
		1: "Does not meet",
	}
)

const (
	ValueTypeEnumPerformance = "performance"
	ValueTypeEnumRelationship = "relationship"
)

var (
	ValueTypeEnumValues = []string{ValueTypeEnumPerformance, ValueTypeEnumRelationship}
)

// AdaptiveValuesTableSchema - schema of adaptive values
type AdaptiveValuesTableSchema struct {
	Name string
	PlatformIDIndex string
}

// AdaptiveValuesTableSchemaForClientID creates table schema given client id
func AdaptiveValuesTableSchemaForClientID(clientID string) AdaptiveValuesTableSchema {
	return AdaptiveValuesTableSchema{
		Name: clientID + "_adaptive_value",
		PlatformIDIndex: "PlatformIDIndex",
	}
}
// ValueTypeMapping mapping from value name to value type
type ValueTypeMapping = map[string]string

// ConvertValuesToValueTypesMapping converts slice of adaptive values
// to mapping from value name to value type.
func ConvertValuesToValueTypesMapping(values []AdaptiveValue) ValueTypeMapping {
	topics := make(map[string]string,0)
	for i:= 0; i < len(values); i++ {
		topic := values[i].Name
		topicType := values[i].ValueType
		topics[topic]=topicType
	}
	return topics
}

// AdaptiveValueFilterActive removes deactivated values
func AdaptiveValueFilterActive(values []AdaptiveValue) (res []AdaptiveValue) {
	for _, v := range values {
		if v.DeactivatedOn == "" {
			res = append(res, v)
		}
	}
	return
}
