package common
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.

type PlatformID string

type PriorityValue string
const (
	UrgentPriority PriorityValue = "Urgent"
	HighPriority PriorityValue = "High"
	MediumPriority PriorityValue = "Medium"
	LowPriority PriorityValue = "Low"
)

type ObjectiveStatusColor string
const (
	ObjectiveStatusRedKey ObjectiveStatusColor = "Red"
	ObjectiveStatusYellowKey ObjectiveStatusColor = "Yellow"
	ObjectiveStatusGreenKey ObjectiveStatusColor = "Green"
)

type PlatformName string
const (
	SlackPlatform PlatformName = "slack"
	MsTeamsPlatform PlatformName = "ms-teams"
)

type AdaptiveCommunityID string
const (
	Admin AdaptiveCommunityID = "admin"
	HR AdaptiveCommunityID = "hr"
	Coaching AdaptiveCommunityID = "coaching"
	User AdaptiveCommunityID = "user"
	Strategy AdaptiveCommunityID = "strategy"
	Capability AdaptiveCommunityID = "capability"
	Initiative AdaptiveCommunityID = "initiative"
	Competency AdaptiveCommunityID = "competency"
)

type CommunityKind string
const (
	AdminCommunity CommunityKind = "admin"
	HRCommunity CommunityKind = "hr"
	CoachingCommunity CommunityKind = "coaching"
	UserCommunity CommunityKind = "user"
	StrategyCommunity CommunityKind = "strategy"
	CompetencyCommunity CommunityKind = "competency"
	ObjectiveCommunity CommunityKind = "objective"
	InitiativeCommunity CommunityKind = "initiative"
)
