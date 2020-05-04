package models

import (
	"github.com/adaptiveteam/adaptive/daos/user"
	"github.com/slack-go/slack"
	"github.com/pkg/errors"
)

type User = user.User

// type User struct {
// 	UserProfile
// 	// Platform of the user
// 	PlatformId    string `json:"platform_id"`
// 	PlatformOrg   string `json:"platform_org"`
// 	IsAdmin       bool   `json:"is_admin"`
// 	IsAdaptiveBot bool   `json:"is_adaptive_bot,omitempty"`
// 	Deleted       bool   `json:"deleted"`
// 	CreatedAt     string `json:"created_at"`
// 	ModifiedAt    string `json:"modified_at,omitempty"`
// 	// This indicates if the user is shared among a group. This is typically for channels, groups, conversations etc.
// 	IsShared bool `json:"is_shared"`
// }

type UserProfile struct {
	// Id of the user, this is the platform specific id
	Id                         string `json:"id"`
	DisplayName                string `json:"display_name"`
	FirstName                  string `json:"first_name,omitempty"`
	LastName                   string `json:"last_name,omitempty"`
	Timezone                   string `json:"timezone"`
	TimezoneOffset             int    `json:"timezone_offset"`
	AdaptiveScheduledTime      string `json:"adaptive_scheduled_time,omitempty"` // in 24 hr format, localtime
	AdaptiveScheduledTimeInUTC string `json:"adaptive_scheduled_time_in_utc,omitempty"`
}

type UserToken struct {
	UserProfile
	ClientPlatform
	ClientPlatformRequest
}

// TeamIDUnsafe extracts TeamID and ensures that it's nonempty.
func (ut UserToken) TeamIDUnsafe() TeamID {
	teamID := ut.ClientPlatformRequest.TeamID
	if teamID.IsEmpty() {
		panic(errors.New("Team ID is empty"))
	}
	return teamID
}

// SlackAPI returns Slack client for the given platform token
func (ut UserToken) SlackAPI() *slack.Client {
	return slack.New(ut.PlatformToken)
}

// AdaptiveUsersTableSchema is schema of Dynamo table with user info.
type AdaptiveUsersTableSchema struct {
	Name                       string
	PlatformIndex              string
	PlatformTZOffsetIndex      string
	PlatformScheduledTimeIndex string
}

// AdaptiveUsersTableSchemaForClientID creates table schema given client id
func AdaptiveUsersTableSchemaForClientID(clientID string) AdaptiveUsersTableSchema {
	return AdaptiveUsersTableSchema{
		Name:                       clientID + "_adaptive_users",
		PlatformIndex:              "PlatformIDIndex",
		PlatformTZOffsetIndex:      "PlatformIDTimezoneOffsetIndex",
		PlatformScheduledTimeIndex: "PlatformIDAdaptiveScheduledTimeInUTCIndex",
	}
}
