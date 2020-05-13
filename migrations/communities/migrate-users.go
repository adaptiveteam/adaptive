package communities


import (
	"github.com/adaptiveteam/adaptive/daos/migration"
	"github.com/adaptiveteam/adaptive/daos/channelMember"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	"github.com/adaptiveteam/adaptive/daos/common"
	"fmt"
)

// MigrateUsers - copies community users to channel users.
func MigrateUsers(conn common.DynamoDBConnection, m *migration.Migration) (err error) {
	var acus []adaptiveCommunityUser.AdaptiveCommunityUser
	err = conn.Dynamo.ScanTable(adaptiveCommunityUser.TableName(conn.ClientID), &acus)
	total := 0
	if err == nil {
		for _, acu := range acus {
			if acu.PlatformID == conn.PlatformID || conn.PlatformID == "-all" {
				total ++
				fmt.Printf("channel %s user %s: ", acu.ChannelID, acu.UserID)
				cm := channelMember.ChannelMember{
					PlatformID: acu.PlatformID,
					ChannelID: acu.ChannelID,
					UserID: acu.UserID,
				}
				err = channelMember.CreateOrUpdate(cm)(conn)
				if err == nil {
					m.SuccessCount ++
					fmt.Printf("ok\n")
				} else {
					fmt.Printf("FAILED: %+v\n", err)
					m.FailuresCount ++
				}
			}
			if err != nil {
				break
			}
		}
	}
	fmt.Printf("ChannelMember: success %d + failures %d of total %d.\n", m.SuccessCount, m.FailuresCount, total)
	return
}
