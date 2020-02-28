package models

import (
	"strings"

	"github.com/adaptiveteam/adaptive/daos/common"
)

// TeamID is an identifier of the team to which we communicate.
// This structure is needed during the transition period from AppID to TeamID.
//
type TeamID struct {
	// AppID is Slack's App identifier
	// Deprecated:
	AppID  common.PlatformID
	TeamID common.PlatformID
}

// ToPlatformID returns only one of the platform ids
func (p TeamID) ToPlatformID() (res common.PlatformID) {
	if p.TeamID == "" {
		res = p.AppID
	} else {
		res = p.TeamID
	}
	return
}

// ToString converts PlatformID to string representation
func (p TeamID) ToString() string {
	return string(p.ToPlatformID())
}

// ParseTeamID converts the value obtained from DB to TeamID
func ParseTeamID(platformID common.PlatformID) (res TeamID) {
	s := string(platformID)
	if strings.HasPrefix(s, "T") {
		res = TeamID{TeamID: platformID}
	} else {
		res = TeamID{AppID: platformID}
	}
	return
}

// IsEmpty checks
func (p TeamID) IsEmpty() bool {
	return p.TeamID == "" && p.AppID == ""
}
