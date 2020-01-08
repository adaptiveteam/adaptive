package models

type AdaptiveCommunity struct {
	Id          string `json:"id"`
	PlatformId  string `json:"platform_id"`
	Channel     string `json:"channel"`
	Active      bool   `json:"active"`
	RequestedBy string `json:"requested_by"`
	CreatedAt   string `json:"created_at"`
}

// AdaptiveCommunityTableSchema is schema of Dynamo table with org-community info.
type AdaptiveCommunityTableSchema struct {
	Name string
	ChannelIndex string
	PlatformIndex string
}

// AdaptiveCommunityTableSchemaForClientID creates table schema given client id
func AdaptiveCommunityTableSchemaForClientID(clientID string) AdaptiveCommunityTableSchema {
	return AdaptiveCommunityTableSchema{
		Name: clientID + "_communities",
		ChannelIndex: "UserCommunityChannelIndex",
		PlatformIndex: "UserCommunityPlatformIndex",
	}
}

type AdaptiveCommunityUser2 struct {
	ChannelId   string `json:"channel_id"`
	UserId      string `json:"user_id"`
	PlatformId  string `json:"platform_id"`
	CommunityId string `json:"community_id"`
}

type AdaptiveCommunityUser3 struct {
	ChannelID   string     `json:"channel_id"`
	UserID      string     `json:"user_id"`
	PlatformID  PlatformID `json:"platform_id"`
	CommunityID string     `json:"community_id"`
}

// CommunityUsersTableSchema is schema of Dynamo table with community user info.
type CommunityUsersTableSchema struct {
	Name string
	ChannelIndex string
	UserCommunityIndex string
	UserIndex string
	CommunityIndex string
}

// CommunityUsersTableSchemaForClientID creates table schema given client id
func CommunityUsersTableSchemaForClientID(clientID string) CommunityUsersTableSchema {
	return CommunityUsersTableSchema{
		Name: clientID + "_community_users",
		ChannelIndex: "CommunityUsersChannelIndex",
		UserCommunityIndex: "CommunityUsersUserCommunityIndex",
		UserIndex: "CommunityUsersUserIndex",
		CommunityIndex: "CommunityUsersCommunityIndex",
	}
}