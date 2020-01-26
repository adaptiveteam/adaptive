package closeout

import (
	"time"
	engCommon "github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	engIssues "github.com/adaptiveteam/adaptive/adaptive-engagements/issues"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"

	// common "github.com/adaptiveteam/adaptive/daos/common"
	issues "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/workflows/exchange"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)
// RequestCloseoutWorkflow - 
var RequestCloseoutWorkflow = exchange.WorkflowInfo{
	Prefix: exchange.CommunityPath,
	Name: exchange.RequestCloseoutNamespace, Init: InitState}

// Namespace -
const Namespace = exchange.RequestCloseoutNamespace

// InitState - the initial state of this workflow.
const InitState wf.State = "init"
// FormShownState
const FormShownState wf.State = "FormShownState"
const (
	ConfirmedEvent wf.Event = "Confirmed"
	RejectedEvent wf.Event = "Rejected"
)
const CommentsDialogShown wf.State = "CommentsDialogShown"
const (
	CommentsSubmittedEvent wf.Event = wf.DialogSubmittedEvent
)
// Workflow is a public interface of workflow template.
type Workflow interface {
	GetNamedTemplate() wf.NamedTemplate
}

// this can only be created using constructor function. Thus we can guarantee that
// all fields will have values.
type workflowImpl struct {
	DynamoDBConnection
	alog.AdaptiveLogger
}

// CreateRequestCloseoutWorkflow - constructor.
func CreateRequestCloseoutWorkflow(
	conn DynamoDBConnection,
	logger alog.AdaptiveLogger,
) Workflow {
	logger.Infoln("RequestCloseoutWorkflow")

	return workflowImpl{
		DynamoDBConnection: conn,
		AdaptiveLogger:     logger,
	}
}

func (w workflowImpl) GetNamedTemplate() (nt wf.NamedTemplate) {
	nt = wf.NamedTemplate{
		Name: Namespace,
		Template: wf.Template{
			Init: InitState, // initial state is "init". This is used when the user first triggers the workflow
			FSA: map[struct {
				wf.State
				wf.Event
			}]wf.Handler{
				{State: InitState, Event: ""}:                  w.OnCloseoutRequested(),
				{State: FormShownState, Event: ConfirmedEvent}: wf.SimpleHandler(w.OnCloseoutConfirmed(), wf.DoneState),
				{State: FormShownState, Event: RejectedEvent}:  wf.SimpleHandler(w.OnCloseoutRejected(), CommentsDialogShown),
				{State: CommentsDialogShown, Event: wf.SurveySubmitted}: 
					wf.SimpleHandler(w.OnCommentsSubmitted(), wf.DoneState),
				{State: CommentsDialogShown, Event: wf.SurveyCancelled}: 
					wf.SimpleHandler(w.OnCommentsCancelled(), wf.DoneState),
				
			},
			Parser: wf.Parser,
		}}
	return
}

// OnCloseoutRequested is triggered by some other workflow.
// Data should contain IssueIDKey and IssueTypeKey
func (w workflowImpl) OnCloseoutRequested() wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		issueID := ctx.Data[IssueIDKey]
		issueType := issues.IssueType(ctx.Data[IssueTypeKey])
		log := w.AdaptiveLogger.WithField("issueID", issueID).WithField("issueType", issueType).WithField("Handler", "OnCloseoutRequested")
		log.Info("Start")
		var issue issues.Issue
		issue, err = issues.Read(issueType, issueID)(w.DynamoDBConnection)
		if err != nil {
			return
		}
		out.Interaction = wf.Interaction{
			Messages: []wf.InteractiveMessage{
				{
					PassiveMessage: areYouOkToCloseoutView(issue),
					InteractiveElements: []wf.InteractiveElement{
						wf.Button(ConfirmedEvent, "I agree"),
						wf.Button(RejectedEvent, "I tend to disagree"),
					},
				},
			},
		}
		out.NextState = FormShownState
		return
	}
}

func areYouOkToCloseoutView(issue Issue) wf.PassiveMessage {
	view := engIssues.GetView(issue.GetIssueType())
	newAndOldIssues := NewAndOldIssues{NewIssue: issue, OldIssue: issue}
	fields := view.GetMainFields(newAndOldIssues)
	return wf.PassiveMessage{
		AttachmentText: ui.Sprintf("%s wants to close the following %s. Are you ok with that?",
			engCommon.TaggedUser(issue.UserObjective.UserID),
			issue.GetIssueType().Template()),
		Fields: ebm.OmitEmpty(fields),
	}
}

func (w workflowImpl) OnCloseoutConfirmed() wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		issueID := ctx.Data[IssueIDKey]
		issueType := issues.IssueType(ctx.Data[IssueTypeKey])
		log := w.AdaptiveLogger.WithField("issueID", issueID).WithField("issueType", issueType).WithField("Handler", "OnCloseoutConfirmed")
		log.Info("Start")
		var issue issues.Issue
		issue, err = issues.Read(issueType, issueID)(w.DynamoDBConnection)
		issueOwner := issue.UserObjective.UserID
		view := engIssues.GetView(issue.GetIssueType())
		objectiveView := view.GetTextView(issue)
		if err == nil {
			issue.UserObjective.PartnerVerifiedCompletion = true
			issue.UserObjective.PartnerVerifiedCompletionDate = core.ISODateLayout.Format(time.Now())
			err = issues.Save(issue)(w.DynamoDBConnection)
			if err == nil {				
				out = ctx.Reply(ui.Sprintf("You have approved the closeout request from <@%s> about the following %s:\n%s", 
					issueOwner, 
					issue.GetIssueType().Template(),
					objectiveView))
				out.Responses = append(out.Responses, platform.Post(platform.ConversationID(issueOwner), platform.MessageContent{
					Message: ui.Sprintf("<@%s> has approved the closeout of the following %s:\n%s", 
						issue.UserObjective.AccountabilityPartner, 
						issue.GetIssueType().Template(),
						objectiveView),
				}))
			}
		}
		return
	}
}

func (w workflowImpl) OnCloseoutRejected() wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		issueID := ctx.Data[IssueIDKey]
		issueType := issues.IssueType(ctx.Data[IssueTypeKey])
		log := w.AdaptiveLogger.WithField("issueID", issueID).WithField("issueType", issueType).WithField("Handler", "OnCloseoutRejected")
		log.Info("Start")
		out.Interaction = wf.OpenSurvey(
			CloseoutRequestRejectedCommentsDialog(),
		)
		out.KeepOriginal = true // we want to keep the original message until dialog is submitted or cancelled
		return
	}
}

const CommentsName = "Comments"

func CloseoutRequestRejectedCommentsDialog() (survey ebm.AttachmentActionSurvey) {
	survey = ebm.AttachmentActionSurvey{
		Title: "Closeout rejected",
		SubmitLabel: "Submit",
		Elements: []ebm.AttachmentActionTextElement{
			{
				Name: CommentsName,
				Label: "Why are you disagreeing with closeout?",
				ElemType: string(ebm.ElemTypeTextArea),
			},
		},
	}
	return
}

// OnCommentsSubmitted - sends the comments about rejected closeout
func (w workflowImpl) OnCommentsSubmitted() wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		issueID := ctx.Data[IssueIDKey]
		issueType := issues.IssueType(ctx.Data[IssueTypeKey])
		form := ctx.Request.DialogSubmissionCallback.Submission

		comments := form[CommentsName]
		log := w.AdaptiveLogger.WithField("issueID", issueID).WithField("issueType", issueType).WithField("Handler", "OnCommentsCancelled")
		log.Info("Start")
		var issue issues.Issue
		issue, err = issues.Read(issueType, issueID)(w.DynamoDBConnection)
		if err == nil {
			issueOwner := issue.UserObjective.UserID
			view := engIssues.GetView(issueType)
			objectiveView := view.GetTextView(issue)
			out.Responses = append(out.Responses, platform.Post(platform.ConversationID(issueOwner), platform.MessageContent{
				Message: ui.Sprintf("<@%s> has rejected the closeout of the following %s:\n%s\nwith the following comments:\n%s", 
					issue.UserObjective.AccountabilityPartner, 
					issueType.Template(),
					objectiveView,
					comments),
			}))
		}
		return
	}
}

// OnCommentsCancelled repeats the closeout request so that the coach
// is able to answer it later.
func (w workflowImpl) OnCommentsCancelled() wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		issueID := ctx.Data[IssueIDKey]
		issueType := issues.IssueType(ctx.Data[IssueTypeKey])
		log := w.AdaptiveLogger.WithField("issueID", issueID).WithField("issueType", issueType).WithField("Handler", "OnCommentsCancelled")
		log.Info("Start")
		var issue issues.Issue
		issue, err = issues.Read(issueType, issueID)(w.DynamoDBConnection)
		if err == nil {
			out.PostponedEvents = []wf.PostponeEventForAnotherUser{
				exchange.RequestCloseoutForIssue(issue),
			}
		}
		return
	}
}
