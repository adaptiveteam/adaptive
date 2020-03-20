package request_coach

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"time"

	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/adaptiveteam/adaptive/workflows/common"
	"github.com/adaptiveteam/adaptive/workflows/exchange"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

const AnalysisChannelKey = "AnalysisChannelKey"
const NewAndOldIssuesKey = "NewAndOldIssuesKey"

// starts analysis and request an immediate callback from workflow engine
func (w workflowImpl) onNewOrUpdatedCoachCommentAvailable(ctx wf.EventHandlingContext, newAndOldIssues NewAndOldIssues) (out wf.EventOutput, err error) {
	// p := newAndOldIssues.NewIssue.Progress[0]
	analysisInput := w.textAnalysisInput(ctx, newAndOldIssues, ExtractProgressCoachComment, "")
	analysisChan := utils.AnalyzeTextAsync(analysisInput, w.DynamoDBConnection)
	out, err = w.standardView(ctx)
	if err == nil {
		out.KeepOriginal = true
		out = out.
			WithRuntimeData(AnalysisChannelKey, analysisChan).
			WithRuntimeData(NewAndOldIssuesKey, newAndOldIssues)
		thankYou := wf.InteractiveMessage{
			PassiveMessage: wf.PassiveMessage{
				Text:             ui.Sprintf("Thank you for providing the feedback. I'll send your comments to <@%s>", newAndOldIssues.NewIssue.UserID),
				OverrideOriginal: false, // we don't want to override the same message again. `view` will override the original message.
			},
		}
		out = out.
			WithPrependInteractiveMessage(thankYou).
			WithPostponedEvent(
				exchange.NotifyOwnerAboutFeedbackOnUpdatesForIssue(newAndOldIssues.NewIssue),
			)
		out.ImmediateEvent = MessageIDAvailableEvent // this is needed to post analysis
	}
	return
}

func (w workflowImpl) OnNewOrUpdatedCoachCommentAvailableOnMessageIDAvailableEvent(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	analysisChanI, ok := ctx.TryGetRuntimeData(AnalysisChannelKey)
	if ok {
		analysisChan := analysisChanI.(utils.ChanTextAnalysisResultsAsync)
		var analysisResults utils.TextAnalysisResults
		analysisResults, err = analysisChan.Read(2 * time.Second) // Slack gives 3 seconds to process dialog submissions. We've spent some time to post the message

		if err == nil {
			msgID := toMapperMessageID(ctx.TargetMessageID)
			convCtx := conversationContext(ctx.Request, msgID)
			out, err = w.standardView(ctx)
			originalMessage := out.Messages[0]
			out = wf.EventOutput{}
			resp := wf.PresentTextAnalysisResults(
				convCtx, analysisResults, originalMessage,
			)
			resp.OverrideOriginal = true

			out.Messages = wf.InteractiveMessages(resp)
		} else {
			w.AdaptiveLogger.WithError(err).Errorf("OnNewOrUpdatedCoachCommentAvailableOnMessageIDAvailableEvent, Analysis results are not ready")
		}
	} else {
		err = errors.New("Couldn't find analysis channel")
	}
	out.KeepOriginal = true
	return
}

type TextExtractor = func(newAndOldIssues NewAndOldIssues) (text string)

func ExtractProgressCoachComment(newAndOldIssues NewAndOldIssues) (text string) {
	p := newAndOldIssues.NewIssue.Progress
	if len(p) > 0 {
		text = p[0].PartnerComments
	}
	return
}

func (w workflowImpl) textAnalysisInput(ctx wf.EventHandlingContext,
	newAndOldIssues NewAndOldIssues, textExtractor TextExtractor,
	dialogSituationID DialogSituationIDWithoutIssueType,
) (analysisInput utils.TextAnalysisInput) {
	itype := newAndOldIssues.NewIssue.GetIssueType()

	context := common.GetDialogContext(dialogSituationID, itype)

	text := textExtractor(newAndOldIssues)
	analysisInput = utils.TextAnalysisInput{
		Text:                       text,
		OriginalMessageAttachments: []ebm.Attachment{},
		Namespace:                  "onNewOrUpdatedCoachCommentAvailable",
		Context:                    context,
	}
	return
}

func conversationContext(request slack.InteractionCallback, msgID mapper.MessageID) utils.ConversationContext {
	ctx := utils.ConversationContext{
		UserID:            request.User.ID,
		ConversationID:    string(wf.GetConversationID(request)),
		OriginalMessageTs: msgID.Ts,
		ThreadTs:          msgID.Ts,
	}
	return ctx
}

func toMapperMessageID(id platform.TargetMessageID) mapper.MessageID {
	return mapper.MessageID{
		ConversationID: id.ConversationID,
		Ts:             id.Ts,
	}
}
