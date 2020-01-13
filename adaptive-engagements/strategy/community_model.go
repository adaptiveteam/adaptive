package strategy

import (
	"github.com/adaptiveteam/adaptive/daos/capabilityCommunity"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
)

type StrategyCommunity struct {
	ID                       string                      `json:"id"` // hash key
	PlatformID               string                      `json:"platform_id"`
	Advocate                 string                      `json:"advocate"`
	Community                community.AdaptiveCommunity `json:"community"`
	ChannelID                string                      `json:"channel_id"`
	ChannelCreated           int                         `json:"channel_created"` // 0 for false, 1 for true
	AccountabilityPartner    string                      `json:"accountability_partner"`
	ParentCommunity          community.AdaptiveCommunity `json:"parent_community"`
	ParentCommunityChannelID string                      `json:"parent_community_channel_id"`
	CreatedAt                string                      `json:"created_at"`
}

type CapabilityCommunity = capabilityCommunity.CapabilityCommunity
// type CapabilityCommunity struct {
// 	ID          string `json:"id"`          // hash key
// 	PlatformID  string `json:"platform_id"` // range key
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	Advocate    string `json:"advocate"`
// 	CreatedBy   string `json:"created_by"`
// 	CreatedAt   string `json:"created_at"`
// }

type StrategyInitiativeCommunity struct {
	ID                    string `json:"id"`          // hash key
	PlatformID            string `json:"platform_id"` // range key
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Advocate              string `json:"advocate"`
	CapabilityCommunityID string `json:"capability_community_id"`
	CreatedBy             string `json:"created_by"`
	CreatedAt             string `json:"created_at"`
}
