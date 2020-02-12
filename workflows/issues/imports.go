package issues

import (
	issuesUtils "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	// wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	// "github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

type Issue = issuesUtils.Issue
type IssueType = issuesUtils.IssueType
type IssueProgressID = issuesUtils.IssueProgressID
type NewAndOldIssues = issuesUtils.NewAndOldIssues
type IssuePredicate = issuesUtils.IssuePredicate
type DialogSituationIDWithoutIssueType = issuesUtils.DialogSituationIDWithoutIssueType

const IDO = issuesUtils.IDO
const SObjective = issuesUtils.SObjective
const Initiative = issuesUtils.Initiative

const (
	DescriptionContext              DialogSituationIDWithoutIssueType = issuesUtils.DescriptionContext
	CloseoutDisagreementContext     DialogSituationIDWithoutIssueType = issuesUtils.CloseoutDisagreementContext
	CloseoutAgreementContext        DialogSituationIDWithoutIssueType = issuesUtils.CloseoutAgreementContext
	UpdateContext                   DialogSituationIDWithoutIssueType = issuesUtils.UpdateContext
	UpdateResponseContext           DialogSituationIDWithoutIssueType = issuesUtils.UpdateResponseContext
	CoachingRequestRejectionContext DialogSituationIDWithoutIssueType = issuesUtils.CoachingRequestRejectionContext
	ProgressUpdateContext           DialogSituationIDWithoutIssueType = issuesUtils.ProgressUpdateContext
)
