package request_coach

import (
	engCommon "github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	engIssues "github.com/adaptiveteam/adaptive/adaptive-engagements/issues"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"

	wfCommon "github.com/adaptiveteam/adaptive/workflows/common"
	// common "github.com/adaptiveteam/adaptive/daos/common"
	issues "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/workflows/exchange"
)

// RequestCoachWorkflow -
var RequestCoachWorkflow = exchange.WorkflowInfo{
	Prefix: exchange.CommunityPath,
	Name: exchange.RequestCoachNamespace, Init: InitState}

// Namespace -
const Namespace = exchange.RequestCoachNamespace

// InitState - the initial state of this workflow.
const InitState wf.State = "init"
// FormShownState -
const FormShownState wf.State = "FormShownState"
const (
	// ConfirmedEvent -
	ConfirmedEvent wf.Event = "Confirmed"
	// RejectedEvent -
	RejectedEvent wf.Event = "Rejected"
)
// UpdateShownState -
const UpdateShownState wf.State = "UpdateShownState"
// DialogShownState -
const DialogShownState wf.State = "DialogShownState"
// Workflow is a public interface of workflow template.
type Workflow interface {
	GetNamedTemplate() wf.NamedTemplate
}

// this can only be created using constructor function. Thus we can guarantee that
// all fields will have values.
type workflowImpl wfCommon.WorkflowContext

// CreateRequestCoachWorkflow - constructor.
func CreateRequestCoachWorkflow(
	conn DynamoDBConnection,
	logger alog.AdaptiveLogger,
) Workflow {
	logger.Infoln("RequestCoachWorkflow")

	return workflowImpl{
		DynamoDBConnection: conn,
		AdaptiveLogger:     logger,
	}
}

func (w workflowImpl) GetNamedTemplate() wf.NamedTemplate {
	nt := wf.NamedTemplate{
		Name: Namespace,
		Template: wf.Template{
			Init: InitState, // initial state is "init". This is used when the user first triggers the workflow
			FSA: map[struct {wf.State; wf.Event}]wf.Handler{
				{State: InitState, Event: ""}:                                 w.OnCoachRequested(),
				{State: InitState, Event: exchange.IssueUpdatedEvent}:         wf.SimpleHandler(w.OnIssueUpdated(), UpdateShownState),
				{State: FormShownState, Event: ConfirmedEvent}:                wf.SimpleHandler(w.OnConfirmed(), wf.DoneState),
				{State: FormShownState, Event: RejectedEvent}:                 wf.SimpleHandler(w.OnRejected(), wf.DoneState),
				{State: UpdateShownState, Event: ConfirmedEvent}:              wf.SimpleHandler(w.OnProvideFeedback(), DialogShownState),
				{State: UpdateShownState, Event: RejectedEvent}:               wf.SimpleHandler(w.OnDismiss(), wf.DoneState),
				{State: UpdateShownState, Event: engIssues.DetailsEvent}:      wf.SimpleHandler(w.OnDetails, UpdateShownState),
				{State: UpdateShownState, Event: engIssues.ProgressShowEvent}: wf.SimpleHandler(w.OnProgressShow, UpdateShownState),
				{State: DialogShownState, Event: wf.DialogSubmittedEvent}:     wf.SimpleHandler(w.OnCommentsSubmitted(), UpdateShownState),
				
			},
			Parser: wf.Parser,
		}}
	return nt
}

// OnCoachRequested is triggered by some other workflow.
// Data should contain IssueIDKey
// as well as Issue type
func (w workflowImpl) OnCoachRequested() wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		issueID := ctx.Data[IssueIDKey]
		issueType := issues.IssueType(ctx.Data[IssueTypeKey])
		log := w.AdaptiveLogger.WithField("issueID", issueID).WithField("issueType", issueType).WithField("Handler", "OnCoachRequested")
		log.Info("Start")
		var issue issues.Issue
		issue, err = issues.Read(issueType, issueID)(w.DynamoDBConnection)
		if err != nil {
			return
		}
		ap := issue.UserObjective.AccountabilityPartner
		if ap == "" || ap == "none" || issue.UserObjective.Accepted == 0 {
			out = out.WithInteractiveMessage(wf.InteractiveMessage{
				PassiveMessage: wf.PassiveMessage{
					AttachmentText: ui.Sprintf("%s is requesting your coaching for the below %s. "+
						"Are you available to partner with and guide your colleague with this effort?",
						engCommon.TaggedUser(issue.UserObjective.UserID),
						issue.GetIssueType().Template()),
					Fields: shortViewFields(issue),
				},
				InteractiveElements: []wf.InteractiveElement{
					wf.Button(ConfirmedEvent, "I agree"),
					wf.Button(RejectedEvent, "I tend to disagree"),
				},
			})
			out.NextState = FormShownState
		} else {
			log.WithField("AccountabilityPartner", ap).Info("AccountabilityPartner already assigned")
			out.NextState = wf.DoneState
		}
		return
	}
}

func shortViewFields(issue Issue) []ebm.AttachmentField {
	view := engIssues.GetView(issue.GetIssueType())
	newAndOldIssues := NewAndOldIssues{NewIssue: issue, OldIssue: issue}
	fields := view.GetMainFields(newAndOldIssues)
	return ebm.OmitEmpty(fields)
}

func (w workflowImpl) OnConfirmed() wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		issueID := ctx.Data[IssueIDKey]
		issueType := issues.IssueType(ctx.Data[IssueTypeKey])
		log := w.AdaptiveLogger.WithField("issueID", issueID).WithField("issueType", issueType).WithField("Handler", "OnConfirmed")
		log.Info("Start")
		var issue issues.Issue
		issue, err = issues.Read(issueType, issueID)(w.DynamoDBConnection)
		issueOwner := issue.UserObjective.UserID
		view := engIssues.GetView(issue.GetIssueType())
		objectiveView := view.GetTextView(issue)
		if err == nil {
			if issue.UserObjective.Accepted == 1 {
				var who ui.RichText
				switch issue.UserObjective.AccountabilityPartner {
				case ctx.Request.User.ID:
					who = "You"
				default:
					who = "Someone else"
				}
				out = ctx.Reply(ui.Sprintf("%s have already accepted the coaching request from <@%s> about the following %s:\n%s", 
					who,
					issueOwner, 
					issue.GetIssueType().Template(),
					objectiveView))
			} else {
				issue.UserObjective.AccountabilityPartner = ctx.Request.User.ID
				issue.UserObjective.Accepted = 1
				err = issues.Save(issue)(w.DynamoDBConnection)
				if err == nil {				
					out = ctx.Reply(ui.Sprintf("You have accepted the request from <@%s> about the following %s:\n%s", 
						issueOwner, 
						issue.GetIssueType().Template(),
						objectiveView))
					out.Responses = append(out.Responses, 
						platform.Post(platform.ConversationID(issue.UserID), 
							platform.Message(ui.Sprintf("<@%s> has accepted your request about the following %s:\n%s",
								ctx.Request.User.ID, 
								issue.GetIssueType().Template(),
								objectiveView)),
						))
				}
			}
		}
		return
	}
}

func (w workflowImpl) OnRejected() wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		issueID := ctx.Data[IssueIDKey]
		issueType := issues.IssueType(ctx.Data[IssueTypeKey])
		log := w.AdaptiveLogger.WithField("issueID", issueID).WithField("issueType", issueType).WithField("Handler", "OnCoachRequested")
		log.Info("Start")
		var issue issues.Issue
		issue, err = issues.Read(issueType, issueID)(w.DynamoDBConnection)
		if err == nil {
			issueOwner := issue.UserObjective.UserID
			view := engIssues.GetView(issue.GetIssueType())
			objectiveView := view.GetTextView(issue)
		
			out = ctx.Reply(ui.Sprintf("I'll notify <@%s> that you have rejected the coaching request", issueOwner))
			out.Responses = append(out.Responses,
				platform.Post(platform.ConversationID(issueOwner),
					platform.MessageContent{Message: ui.Sprintf(
						"<@%s> has just rejected your coaching request of the below %s:\n%s",
						ctx.Request.User.ID,
						issue.GetIssueType().Template(),
						objectiveView, 
					)}),
			)
		}
		return
	}
}
