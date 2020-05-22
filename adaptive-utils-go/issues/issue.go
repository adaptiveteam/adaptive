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
	return DetectIssueType(issue.UserObjective)
}

// DetectIssueType is the reference mechanism to detect issue type
func DetectIssueType(uo userObjective.UserObjective) (itype IssueType) {
	itype = IDO
	switch uo.ObjectiveType {
	case userObjective.IndividualDevelopmentObjective:
		itype = IDO
	case userObjective.StrategyDevelopmentObjectiveIssue:
		itype = SObjective
	case userObjective.StrategyDevelopmentInitiative:
		itype = Initiative
	case userObjective.StrategyDevelopmentObjective:
		log.Printf("WARN using old-style issue type detection")
		itype = SObjective
		switch uo.StrategyAlignmentEntityType {
		case userObjective.ObjectiveStrategyObjectiveAlignment:
			itype = SObjective
		case userObjective.ObjectiveStrategyInitiativeAlignment:
			itype = Initiative
		default:
			log.Printf("WARN (1) Couldn't determine issue type for %s. ObjectiveType=%s, StrategyAlignmentEntityType=%s\n", uo.ID, uo.ObjectiveType, uo.StrategyAlignmentEntityType)
		}
	default:
		log.Printf("WARN (2) Couldn't determine issue type for %s. ObjectiveType=%s, StrategyAlignmentEntityType=%s\n", uo.ID, uo.ObjectiveType, uo.StrategyAlignmentEntityType)
	}
	return
}

// GetObjectiveType -
func (itype IssueType) GetObjectiveType() (otype userObjective.DevelopmentObjectiveType) {
	switch itype {
	case IDO:
		otype = userObjective.IndividualDevelopmentObjective
	case SObjective:
		otype = userObjective.StrategyDevelopmentObjectiveIssue
	case Initiative:
		otype = userObjective.StrategyDevelopmentInitiative
	}
	return
}

// SetIssueType -
func SetIssueType(uo *userObjective.UserObjective, itype IssueType) {
	uo.ObjectiveType = itype.GetObjectiveType()
}

// GetIssueID returns issue.UserObjective.ID
func (issue Issue) GetIssueID() string {
	return issue.UserObjective.ID
}

func (itype IssueType)Template() (text ui.PlainText) {
	return ui.PlainText(itype.FoldString("Individual Development Objective", "Strategy Objective", "Strategy Initiative"))
}

func (itype IssueType)FoldString(ido, sobjective, initiative string) (text string) {
	switch itype {
	case IDO: text = ido
	case SObjective: text = sobjective
	case Initiative: text = initiative
	}
	return
}

// ObjectiveTypeLabel returns a user-friendly objective type label
func ObjectiveTypeLabel(userObj models.UserObjective) ui.PlainText {
	return DetectIssueType(userObj).Template()
}

type IssueProgressID struct {
	IssueID string
	Date string
}


type DialogSituationIDWithoutIssueType = string

const (
	DescriptionContext              DialogSituationIDWithoutIssueType = "description"
	CloseoutDisagreementContext     DialogSituationIDWithoutIssueType = "close-out-disagreement"
	CloseoutAgreementContext        DialogSituationIDWithoutIssueType = "close-out-agreement"
	UpdateContext                   DialogSituationIDWithoutIssueType = "update"
	UpdateResponseContext           DialogSituationIDWithoutIssueType = "update-response"
	CoachingRequestRejectionContext DialogSituationIDWithoutIssueType = "coaching-request-rejection"
	ProgressUpdateContext           DialogSituationIDWithoutIssueType = UpdateContext
	ProgressUpdateResponseContext   DialogSituationIDWithoutIssueType = UpdateResponseContext
)
