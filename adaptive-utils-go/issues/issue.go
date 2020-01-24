package issues

import (
	"log"
	"github.com/adaptiveteam/adaptive/daos/adaptiveValue"
	userObjective "github.com/adaptiveteam/adaptive/daos/userObjective"
	userObjectiveProgress "github.com/adaptiveteam/adaptive/daos/userObjectiveProgress"
	// wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

// IssueType is one of IDO, objective, initiative
type IssueType string
const (
	IDO IssueType = "IDO"
	SObjective IssueType = "SObjective"
	Initiative IssueType = "Initiative"
)

// PrefetchedData contains information from other connected tables that is used
// to represent the issue.
type PrefetchedData struct {
	Progress []userObjectiveProgress.UserObjectiveProgress
	/*
	func readUserDisplayName(userID string) (displayName ui.PlainText) {
		accountabilityPartner, err2 := utils.UserToken(userID, userProfileLambda, region, namespace)

		if err == nil {
			displayName = ui.PlainText(accountabilityPartner.DisplayName)
		} else {
			displayName = "Unknown"
			logger.Infof("Couldn't find AccountabilityPartner @" + userID)
		}
		return
	}

	*/
	AccountabilityPartner models.User
	AlignedCapabilityObjective models.StrategyObjective
	AlignedCapabilityInitiative models.StrategyInitiative
	AlignedCompetency adaptiveValue.AdaptiveValue
	AlignedCapabilityCommunity models.CapabilityCommunity
	AlignedInitiativeCommunity models.StrategyInitiativeCommunity
}
// This file contains a generic mechanism for handling the creation of IDOs, strategy objectives, initiatives.
// Issue is the type of entity this workflow mostly deals with.
type Issue struct {
	userObjective.UserObjective
	models.StrategyObjective
	models.StrategyInitiative
	PrefetchedData
}
// NewAndOldIssues contains two versions of the issue
type NewAndOldIssues struct {
	NewIssue Issue
	OldIssue Issue
	Updated bool
}
// IssuePredicate is a predicate on the issue
type IssuePredicate = func (issue Issue) bool

// GetIssueType detects the issue type of the existing issue
func (issue Issue) GetIssueType() (itype IssueType) {
	switch issue.UserObjective.ObjectiveType {
	case userObjective.IndividualDevelopmentObjective:
		itype = IDO
	case userObjective.StrategyDevelopmentObjective:
		itype = SObjective
		switch issue.StrategyAlignmentEntityType {
		case userObjective.ObjectiveStrategyObjectiveAlignment:
			itype = SObjective
		case userObjective.ObjectiveStrategyInitiativeAlignment:
			itype = Initiative
		}
	default:
		log.Printf("Couldn't determine issue type for %s. ObjectiveType=%s, StrategyAlignmentEntityType=%s\n", issue.UserObjective.ID, issue.UserObjective.ObjectiveType, issue.StrategyAlignmentEntityType)
	}
	return
}

// GetIssueID returns issue.UserObjective.ID
func (issue Issue) GetIssueID() string {
	return issue.UserObjective.ID
}

func (itype IssueType)Template() (text ui.PlainText) {
	switch itype {
	case IDO: text = "Individual Development Objective"
	case SObjective: text = "Strategy Objective"
	case Initiative: text = "Strategy Initiative"
	}
	return
}

type IssueProgressID struct {
	IssueID string
	Date string
}
