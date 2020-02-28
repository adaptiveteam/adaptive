package models

import (
	"encoding/json"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/daos/userAttribute"
	"github.com/adaptiveteam/adaptive/daos/userEngagement"
	
)

// UserEngagement encapsulates an engagement we want to provide to a user
type UserEngagement  = userEngagement.UserEngagement
// type UserEngagement struct {
// 	// PlatformID is the identifier of the platform.
// 	// It's used to get platform token required to send message to Slack/Teams.
// 	PlatformID PlatformID `json:"platform_id"`
// 	// UserId is the Id of the user to send an engagement to
// 	// This usually corresponds to the platform user id
// 	UserID string `json:"user_id"`
// 	// TargetId is the Id of the user for whom this is related to
// 	TargetID string `json:"target_id"`
// 	// Namespace for the engagement
// 	Namespace string `json:"namespace"`
// 	// A unique id to identify the engagement
// 	ID string `json:"id"`
// 	UserEngagementCheckWithValue
// 	// Script that should be sent to a user to start engaging.
// 	// It's a serialized ebm.Message
// 	// Deprecated: Use `Message` directly.
// 	Script string `json:"script"`
// 	// Message is the message we want to send to user.
// 	Message ebm.Message `json:"message"`
// 	// Priority of the engagement
// 	// Urgent priority engagements are immediately sent to a user
// 	// Other priority engagements are queued up in the order of priority to be sent to user in next window
// 	Priority PriorityValue `json:"priority"`
// 	// RFC3339 timestamp at which the engagement was posted to the user
// 	PostedAt string `json:"posted_at,omitempty"`
// 	// A flag indicating is a user has responded to the engagement - 1 for answered, 0 for un-answered
// 	// this is required because, we need to keep the engagement even after a user has answered it
// 	// If the user wants to edit later, we will refer to the same engagement to post to user, like getting survey information
// 	// So, we need a way to differentiate between answered and unanswered engagements
// 	Answered int `json:"answered"`
// 	// Flag indicating if an engagement is ignored, 1 for yes, 0 for no
// 	Ignored            int    `json:"ignored"`
// 	EffectiveStartDate string `json:"effective_start_date,omitempty"`
// 	EffectiveEndDate   string `json:"effective_end_date,omitempty"`
// 	// When a same engagement is written to dynamo, dynamo doesn't update it and it not treated as a new event
// 	// This timestamp will help to identify newer same engagement
// 	CreatedAt string `json:"created_at"`
// 	// Re-scheduled timestamp for the engagement, if any
// 	RescheduledFrom string `json:"rescheduled_from"`
// }

type UserEngagementCheckWithValue struct {
	// Check identifier for the engagement
	CheckIdentifier string `json:"check_identifier,omitempty"`
	CheckValue      bool   `json:"check_value,omitempty"`
}

type UserEngageWithCheckValues struct {
	UserEngage
	CheckValues map[string]bool
}

// UserAttribute encapsulates key-value setting for a user
type UserAttribute = userAttribute.UserAttribute
// type UserAttribute struct {
// 	// Id of the user
// 	UserId string `json:"user_id"`
// 	// Key of the setting
// 	AttrKey string `json:"attr_key"`
// 	// Value of the settings
// 	AttrValue string `json:"attr_value"`
// 	// A flag that tells whether setting is default or is explicitly set
// 	// Every user will have default settings
// 	Default bool `json:"default"`
// }

// UserEngage encapsulates the struct that will be used to trigger engaging with the user
type UserEngage struct {
	// TeamID is the
	TeamID TeamID `json:"platform_id"`
	// UserID is the Id of the user
	UserID string `json:"user_id"`
	// TargetID is the Id of the user who will be affected by this
	TargetID string `json:"target_id,omitempty"`
	// IsNew is a flag to indicate if this is newly added user
	IsNew bool `json:"is_new"`
	// Flag indicating if this engaging is for updating
	Update bool `json:"update"`
	// Channel, if any, to engage user in
	Channel string `json:"channel,omitempty"`
	// ThreadTimestamp if any to engage user in
	ThreadTs string `json:"thread_ts"`
	// Date to emulate for the user
	Date string `json:"date"`
	// Flag to indicate if tne invocation is on-demand
	OnDemand bool `json:"on_demand"`
}

// UserEngagementGetMessage parses `Script` as `Message`
func UserEngagementGetMessage() (eng UserEngagement, postMsg ebm.Message, err error) {
	err = json.Unmarshal([]byte(eng.Script), &postMsg)
	return
}
