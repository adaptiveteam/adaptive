package workflow


import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
)

// PostponeEvent saves an event to a database for further processing.
// The database will be eventually evaluated for a particular user and
// the event will be triggered.
type PostponeEvent = func(teamID models.TeamID, postponedEvent PostponeEventForAnotherUser) error

// ResolveCommunity an environment function that resolves a community. If not available, error is returned
type ResolveCommunity = func (communityID string) (conversationID platform.ConversationID, err error)

// Environment contains mechanisms to deal with external world
type Environment struct {
	// this is provided from outside as a context. When we want to
	// have a callback routed to our instance, we should prepend this prefix.
	Prefix         models.Path
	GetPlatformAPI PlatformAPIForTeamID
	LogInfof       LogInfof
	PostponeEvent
	ResolveCommunity 
}
