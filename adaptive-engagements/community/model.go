package community

type AdaptiveCommunity string

const (
	Admin      AdaptiveCommunity = "admin"
	HR         AdaptiveCommunity = "hr"
	Coaching   AdaptiveCommunity = "coaching"
	User       AdaptiveCommunity = "user"
	Strategy   AdaptiveCommunity = "strategy"
	Capability AdaptiveCommunity = "capability"
	Initiative AdaptiveCommunity = "initiative"
	Competency AdaptiveCommunity = "competency"
)

var (
	NonStrategyCommunityList = []string{string(Admin), string(HR), string(Coaching), string(User), string(Competency)}
)
