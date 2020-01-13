package models

import (
	"github.com/adaptiveteam/adaptive/daos/capabilityCommunity"
)

type StrategyObjectiveType string

const (
	CustomerStrategyObjective   StrategyObjectiveType = "Customer"
	FinancialStrategyObjective  StrategyObjectiveType = "Financial"
	CapabilityStrategyObjective StrategyObjectiveType = "Capability"
)

type StrategyObjective struct {
	ID                     string                `json:"id"`          // hash
	PlatformID             string                `json:"platform_id"` // range key
	Name                   string                `json:"name"`
	Description            string                `json:"description"`
	AsMeasuredBy           string                `json:"as_measured_by"`
	Targets                string                `json:"targets"`
	Type                   StrategyObjectiveType `json:"type"`
	Advocate               string                `json:"advocate"`
	CapabilityCommunityIDs []string              `json:"capability_community_ids,omitempty"` // community id not require d for customer/financial objectives
	ExpectedEndDate        string                `json:"expected_end_date"`
	CreatedBy              string                `json:"created_by"`
	CreatedAt              string                `json:"created_at"`
}

type StrategyInitiative struct {
	ID                    string `json:"id"`
	PlatformID            string `json:"platform_id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	DefinitionOfVictory   string `json:"definition_of_victory"`
	Advocate              string `json:"advocate"`
	InitiativeCommunityID string `json:"initiative_community_id"`
	Budget                string `json:"budget"`
	ExpectedEndDate       string `json:"expected_end_date"`
	CapabilityObjective   string `json:"capability_objective"`
	CreatedAt             string `json:"created_at"`
	CreatedBy             string `json:"created_by"`
	ModifiedAt            string `json:"modified_at"`
	ModifiedBy            string `json:"modified_by"`
}

type VisionMission struct {
	ID         string `json:"id"`          // hash
	PlatformID string `json:"platform_id"` // range key
	Mission    string `json:"mission"`
	Vision     string `json:"vision"`
	Advocate   string `json:"advocate"`
	CreatedBy  string `json:"created_by"`
	CreatedAt  string `json:"created_at"`
}

type StrategyCommunity struct {
	ID                       string            `json:"id"` // hash key
	PlatformID               string            `json:"platform_id"`
	Advocate                 string            `json:"advocate"`
	Community                AdaptiveCommunity `json:"community"`
	ChannelID                string            `json:"channel_id"`
	ChannelCreated           int               `json:"channel_created"` // 0 for false, 1 for true
	AccountabilityPartner    string            `json:"accountability_partner"`
	ParentCommunity          AdaptiveCommunity `json:"parent_community"`
	ParentCommunityChannelID string            `json:"parent_community_channel_id"`
	CreatedAt                string            `json:"created_at"`
}

type CapabilityCommunity = capabilityCommunity.CapabilityCommunity
//  struct {
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
