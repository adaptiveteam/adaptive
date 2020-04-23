package workflow

import (
	"github.com/pkg/errors"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunity"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/daos/common"
)

// ResolveCommunityImpl is a reference implementation of community resolution function.
func ResolveCommunityImpl(conn common.DynamoDBConnection) ResolveCommunity {
	return func(communityID string) (conversationID platform.ConversationID, err error) {
		var comms []adaptiveCommunity.AdaptiveCommunity
		comms, err = adaptiveCommunity.ReadByPlatformID(conn.PlatformID)(conn)
		if err == nil {
			for _, comm := range comms {
				if comm.ID == communityID {
					conversationID = platform.ConversationID(comm.ChannelID)
					return
				}
			}
			err = errors.Errorf("Couldn't find community %s in team %v", communityID, conn.PlatformID)
		}
		return
	}
}
