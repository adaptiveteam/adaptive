package common


import (
	"log"
	// alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	issuesUtils "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
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

type GlobalDialogContext = string
type GlobalDialogContextByType = map[IssueType]GlobalDialogContext

var contexts = map[DialogSituationIDWithoutIssueType]GlobalDialogContextByType{
	DescriptionContext: {
		IDO:        "dialog/ido/language-coaching/description",
		SObjective: "",
		Initiative: "",
	},
	CloseoutAgreementContext: {
		IDO:        "dialog/ido/language-coaching/close-out-agreement",
		SObjective: "dialog/strategy/language-coaching/objective/close-out-agreement",
		Initiative: "dialog/strategy/language-coaching/initiative/close-out-agreement",
	},
	CloseoutDisagreementContext: {
		IDO:        "dialog/ido/language-coaching/close-out-disagreement",
		SObjective: "dialog/strategy/language-coaching/objective/close-out-disagreement",
		Initiative: "dialog/strategy/language-coaching/initiative/close-out-disagreement",
	},
	UpdateContext: {
		IDO:        "dialog/ido/language-coaching/update",
		SObjective: "dialog/strategy/language-coaching/objective/update",
		Initiative: "dialog/strategy/language-coaching/initiative/update",
	},
	UpdateResponseContext: {
		IDO:        "dialog/ido/language-coaching/update-response",
		SObjective: "dialog/strategy/language-coaching/objective/update-response",
		Initiative: "dialog/strategy/language-coaching/initiative/update-response",
	},
	CoachingRequestRejectionContext: {
		IDO:        "dialog/ido/language-coaching/coaching-request-rejection",
		SObjective: "",
		Initiative: "",
	},
	ProgressUpdateContext: { // TODO: provide progress update contexts
		IDO:        "dialog/ido/language-coaching/update",
		SObjective: "dialog/strategy/language-coaching/objective/update",
		Initiative: "dialog/strategy/language-coaching/initiative/update",
	},
}

// GetDialogContext returns dialog context.
// In case of errors logs, and returns empty context.
func GetDialogContext(dialogSituationID DialogSituationIDWithoutIssueType, itype IssueType) (context string) {
	contextsForSituation, ok := contexts[dialogSituationID]
	if ok {
		context, ok = contextsForSituation[itype]
		if !ok {
			log.Printf("Missing dialog context (dialogSituationID=%s, itype=%s)", dialogSituationID, itype)
		}
	} else {
		log.Printf("Missing dialog situation (dialogSituationID=%s)", dialogSituationID)
	}
	return
}
