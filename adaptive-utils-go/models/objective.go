package models

import (
	common "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
	"github.com/adaptiveteam/adaptive/daos/userObjectiveProgress"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

type DevelopmentObjectiveType string

const (
	IndividualDevelopmentObjective DevelopmentObjectiveType = "individual"
	StrategyDevelopmentObjective   DevelopmentObjectiveType = "strategy"
)

type AlignedStrategyType string

const (
	ObjectiveStrategyObjectiveAlignment  AlignedStrategyType = "strategy_objective"
	ObjectiveStrategyInitiativeAlignment AlignedStrategyType = "strategy_initiative"
	ObjectiveCompetencyAlignment         AlignedStrategyType = "competency"
	ObjectiveNoStrategyAlignment         AlignedStrategyType = "none"
)

type ObjectiveStatusColor = common.ObjectiveStatusColor

const (
	ObjectiveStatusRedKey    = common.ObjectiveStatusRedKey
	ObjectiveStatusYellowKey = common.ObjectiveStatusYellowKey
	ObjectiveStatusGreenKey  = common.ObjectiveStatusGreenKey

	ObjectiveStatusRedLabel    ui.PlainText = "Off Track" // "Red"
	ObjectiveStatusYellowLabel ui.PlainText = "At Risk"   // "Yellow"
	ObjectiveStatusGreenLabel  ui.PlainText = "On Track"  // "Green"
)

var (
	ObjectiveStatusColorKeys = []ObjectiveStatusColor{
		ObjectiveStatusRedKey,
		ObjectiveStatusYellowKey,
		ObjectiveStatusGreenKey,
	}
	// KvPair is converted to drop-down in an incorrect way. See utils.AttachActionElementOptions
	ObjectiveStatusColorKeyValues = []KvPair{
		{Value: string(ObjectiveStatusRedKey), Key: string(ObjectiveStatusRedLabel)},
		{Value: string(ObjectiveStatusYellowKey), Key: string(ObjectiveStatusYellowLabel)},
		{Value: string(ObjectiveStatusGreenKey), Key: string(ObjectiveStatusGreenLabel)},
	}
	ObjectiveStatusColorLabels = map[ObjectiveStatusColor]ui.PlainText{
		"":                       "", // this works anyway.
		ObjectiveStatusRedKey:    ObjectiveStatusRedLabel,
		ObjectiveStatusYellowKey: ObjectiveStatusYellowLabel,
		ObjectiveStatusGreenKey:  ObjectiveStatusGreenLabel,
	}
)
type UserObjective = userObjective.UserObjective
// type UserObjective struct {
// 	UserID                        string                   `json:"user_id"`
// 	Name                          string                   `json:"name"`
// 	ID                            string                   `json:"id"`
// 	Description                   string                   `json:"description"`
// 	AccountabilityPartner         string                   `json:"accountability_partner"`
// 	Accepted                      int                      `json:"accepted"` // 1 for true, 0 for false
// 	Type                          DevelopmentObjectiveType `json:"type"`
// 	StrategyAlignmentEntityID     string                   `json:"strategy_alignment_entity_id"`
// 	StrategyAlignmentEntityType   AlignedStrategyType      `json:"strategy_alignment_entity_type"`
// 	Quarter                       int                      `json:"quarter"`
// 	Year                          int                      `json:"year"`
// 	CreatedDate                   string                   `json:"created_date"`
// 	ExpectedEndDate               string                   `json:"expected_end_date"`
// 	Completed                     int                      `json:"completed"` // 1 for true, 0 for false
// 	PartnerVerifiedCompletion     bool                     `json:"partner_verified_completion"`
// 	CompletedDate                 string                   `json:"completed_date,omitempty"`
// 	PartnerVerifiedCompletionDate string                   `json:"partner_verified_completion_date,omitempty"`
// 	Comments                      string                   `json:"comments"`
// 	Cancelled                     int                      `json:"cancelled"` // 1 for true, 0 for false
// 	PlatformID                    PlatformID               `json:"platform_id"`
// }

type UserObjectiveProgress = userObjectiveProgress.UserObjectiveProgress
// type UserObjectiveProgress struct {
// 	ID                      string               `json:"id"`
// 	CreatedOn               string               `json:"created_on"`
// 	UserID                  string               `json:"user_id"`
// 	Comments                string               `json:"comments"`
// 	Closeout                int                  `json:"closeout"` // 1 for true, 0 for false
// 	PercentTimeLapsed       string               `json:"percent_time_lapsed"`
// 	StatusColor             ObjectiveStatusColor `json:"status_color"`
// 	PartnerID               string               `json:"partner_id"`
// 	ReviewedByPartner       bool                 `json:"reviewed_by_partner"`
// 	PartnerComments         string               `json:"partner_comments"`
// 	PartnerReportedProgress string               `json:"partner_reported_progress"`
// 	PlatformID              PlatformID           `json:"platform_id"`
// }

type UserObjectiveWithProgress struct {
	Objective UserObjective           `json:"objective"`
	Progress  []UserObjectiveProgress `json:"progress"`
}

type UserObjectiveNotAccepted struct {
	ObjectiveID             string `json:"objective_id"`
	AccountabilityPartnerID string `json:"accountability_partner_id"`
	Comments                string `json:"comments"`
	Timestamp               string `json:"timestamp"`
}

type AccountabilityPartnerShipRejection struct {
	ObjectiveID             string `json:"objective_id"` // hash
	CreatedOn               string `json:"created_on"`   // range
	UserID                  string `json:"user_id"`
	AccountabilityPartnerID string `json:"accountability_partner_id"`
	Comments                string `json:"comments"`
}
