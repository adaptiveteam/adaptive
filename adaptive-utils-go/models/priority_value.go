package models

import (
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

type PriorityValue string
 
const (
	UrgentPriority PriorityValue = "Urgent"
	HighPriority   PriorityValue = "High"
	MediumPriority PriorityValue = "Medium"
	LowPriority    PriorityValue = "Low"
)

// Urgency converts boolean value to Priority
func Urgency(urgent bool) PriorityValue {
	return core.IfThenElse(urgent, UrgentPriority, HighPriority).(PriorityValue)
}
