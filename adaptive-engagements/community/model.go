package community

import (
	"github.com/adaptiveteam/adaptive/daos/common"
)

type AdaptiveCommunity = common.AdaptiveCommunityID

const (
	Admin      = common.Admin
	HR         = common.HR
	Coaching   = common.Coaching
	User       = common.User
	Strategy   = common.Strategy
	Capability = common.Capability
	Initiative = common.Initiative
	Competency = common.Competency
)

// AdaptiveCommunityShow - 
func AdaptiveCommunityShow(community AdaptiveCommunity) (res string) {
	switch community {
	case Capability: res = "Objective"
	default:
		res = string(community)
	}
	return
}

var (
	NonStrategyCommunityList = []string{string(Admin), string(HR), string(Coaching), string(User), string(Competency)}
)
