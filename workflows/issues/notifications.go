package issues

import (
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)
// notifications about changes in issue
// Triggers Analysis!
func (w workflowImpl) onNewOrUpdatedItemAvailable(ctx wf.EventHandlingContext,
	tc IssueTypeClass,
	newAndOldIssues NewAndOldIssues,
	dialogSituationID DialogSituationIDWithoutIssueType, //ContextByType GlobalDialogContextByType,
	eventDescription ui.RichText,
	isShowingProgress bool,
) (out wf.EventOutput, err error) {
	// keeping the old item so that we'll be able to show it again after analysis.
	ctx.RuntimeData = runtimeData(newAndOldIssues)
	out, err = w.standardView(ctx)
	out.ImmediateEvent = MessageIDAvailableEventInContext(dialogSituationID) // this is needed to post analysis
	if err == nil {
		if newAndOldIssues.NewIssue.GetIssueType() != IDO {

			notification, err2 := w.notifyStrategy(ctx, tc, newAndOldIssues, eventDescription, isShowingProgress)
			if err2 == nil {
				out.Responses = append(out.Responses, notification)
			} else {
				w.AdaptiveLogger.WithError(err2).Error("onNewOrUpdatedItemAvailable/notifyStrategy")
			}
		}
	}
	// We also set output runtime date
	out.RuntimeData = runtimeData(newAndOldIssues)
	return
}

func (w workflowImpl) notifyStrategy(ctx wf.EventHandlingContext, tc IssueTypeClass,
	newAndOldIssues NewAndOldIssues,
	eventDescription ui.RichText,
	isShowingProgress bool,
) (notification platform.Response, err error) {
	var strategyCommunityConversation platform.ConversationID
	strategyCommunityConversation, err = findStrategyCommunityConversation(w, ctx)
	if err != nil {
		return
	}

	msgToStrategyCommunity := viewObjectiveReadonly(w, tc, newAndOldIssues, isShowingProgress)
	// userID := ctx.Request.User.ID
	msgToStrategyCommunity.Message = eventDescription//ui.Sprintf("Below objective has been %s by <@%s>", eventDescription, userID)
	notification = platform.Post(strategyCommunityConversation, msgToStrategyCommunity)
	return
}
