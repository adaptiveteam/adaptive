package exchange

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
)

const IssueIDKey = "iid"
const IssueTypeKey = "itype"
const DialogSituationKey = "dsit"

// const ItemIDKey = "itemID"
const CapCommIDKey = "cid"//"capCommID"
const IsShowingDetailsKey = "isd"
const IsShowingProgressKey = "isp"

func GetIssueID(ctx wf.EventHandlingContext) string {
	return ctx.Data[IssueIDKey]
}

func GetIssueType(ctx wf.EventHandlingContext) IssueType {
	return IssueType(ctx.Data[IssueTypeKey])
}

func GetDialogSituation(ctx wf.EventHandlingContext) DialogSituationIDWithoutIssueType {
	return DialogSituationIDWithoutIssueType(ctx.Data[DialogSituationKey])
}

func IsShowingDetails(ctx wf.EventHandlingContext) (res bool) {
	_, res = ctx.Data[IsShowingDetailsKey]
	return
}

func IsShowingProgress(ctx wf.EventHandlingContext) (res bool) {
	_, res = ctx.Data[IsShowingProgressKey]
	return
}

const CommunityNamespace = "community"

var CommunityPath models.Path = models.ParsePath("/" + CommunityNamespace)

const FeedbackNamespace = "feedback"

var CoachingPath models.Path = models.ParsePath("/"+FeedbackNamespace)

// RequestCoachNamespace -
const RequestCoachNamespace = "request_coach"

