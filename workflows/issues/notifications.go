package issues

import (
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/adaptiveteam/adaptive/workflows/exchange"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	engIssues"github.com/adaptiveteam/adaptive/adaptive-engagements/issues"
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
	if err == nil {
		out.ImmediateEvent = MessageIDAvailableEventInContext(dialogSituationID) // this is needed to post analysis
		out = out.
			WithPostponedEvent(
				exchange.NotifyAboutUpdatesForIssue(newAndOldIssues, dialogSituationID),
			).
			WithRuntimeData(newAndOldIssues). // We also set output runtime date
			WithResponse(
				w.notifyStrategyIfNeeded(ctx, tc, newAndOldIssues, eventDescription, isShowingProgress)...
			)
	}
	return
}

func (w workflowImpl) notifyStrategyIfNeeded(ctx wf.EventHandlingContext, tc IssueTypeClass,
	newAndOldIssues NewAndOldIssues,
	eventDescription ui.RichText,
	isShowingProgress bool,
) (notifications []platform.Response) {
	if newAndOldIssues.NewIssue.GetIssueType() != IDO {
		strategyCommunityConversation, err2 := findStrategyCommunityConversation(w, ctx)
		if err2 == nil {
			notifications = append(notifications, 
				platform.Post(strategyCommunityConversation, platform.MessageContent{
					Message: eventDescription,
					Attachments: []ebm.Attachment{{
						Fields: ebm.OmitEmpty(
							engIssues.OrdinaryViewFields(newAndOldIssues, 
								engIssues.ViewState{IsShowingProgress: isShowingProgress})),
					}},
				}))
		} else {
			w.AdaptiveLogger.WithError(err2).Error("onNewOrUpdatedItemAvailable/notifyStrategy")
		}
	}
	return
}
