package adaptiveCommunityUser
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.

type FieldName string
const (
	ChannelID FieldName = "channel_id"
	UserID FieldName = "user_id"
	PlatformID FieldName = "platform_id"
	CommunityID FieldName = "community_id"
)

type IndexName string
const (
	ChannelIDIndex IndexName = "ChannelIDIndex"
	UserIDCommunityIDIndex IndexName = "UserIDCommunityIDIndex"
	UserIDIndex IndexName = "UserIDIndex"
	PlatformIDCommunityIDIndex IndexName = "PlatformIDCommunityIDIndex"
)
