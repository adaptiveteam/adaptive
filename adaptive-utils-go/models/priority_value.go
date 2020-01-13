package models

import (
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	common "github.com/adaptiveteam/adaptive/daos/common"
)

type PriorityValue = common.PriorityValue
 
const (
	UrgentPriority = common.UrgentPriority
	HighPriority   = common.HighPriority
	MediumPriority = common.MediumPriority
	LowPriority    = common.LowPriority
)

// Urgency converts boolean value to Priority
func Urgency(urgent bool) PriorityValue {
	return core.IfThenElse(urgent, UrgentPriority, HighPriority).(PriorityValue)
}
