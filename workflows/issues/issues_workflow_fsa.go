package issues

import (
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	dialogFetcher "github.com/adaptiveteam/adaptive/dialog-fetcher"

	// "encoding/json"
	// "fmt"
	// "github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	// "github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	// "github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	// "github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	// "github.com/thoas/go-funk"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	// "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	// "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	// core "github.com/adaptiveteam/adaptive/core-utils-go"
	// ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	// "github.com/adaptiveteam/adaptive/engagement-builder/ui"
	// mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	// "github.com/adaptiveteam/adaptive/daos/strategyObjective"
	ex "github.com/adaptiveteam/adaptive/workflows/exchange"
)

const issueIDKey = ex.IssueIDKey
const capCommIDKey = "cid"
const initCommIDKey = "icid"
const issueTypeKey = ex.IssueTypeKey
const isShowingDetailsKey = "isd"
const isShowingProgressKey = "isp"
const dialogSituationIDKey = "sid"

// IssuesWorkflow is a description of an issues workflow
var IssuesWorkflow = ex.WorkflowInfo{Name: IssuesNamespace, Init: InitState}

const IssuesNamespace = "issues"

const InitState wf.State = "init"
const (
	DefaultEvent            wf.Event = ""
	// ViewListOfIssuesEvent   wf.Event = "ViewListOfIssuesEvent"
	// ViewMyListOfIssuesEvent wf.Event = "ViewMyListOfIssuesEvent"
)

func CreateIssueByTypeEvent(itype IssueType) wf.Event {
	return wf.Event("Cr" + string(itype) + "Evt")
}

func eventByType(name string, itype IssueType) wf.Event {
	return wf.Event(name + "(" + string(itype) + ")")
}
func ViewListOfIssuesByTypeEvent(itype IssueType) wf.Event {
	return eventByType("VLOfIssuesByType", itype)
}
func ViewMyListOfIssuesByTypeEvent(itype IssueType) wf.Event {
	return eventByType("VMyLOfIssuesByType", itype)
}
func ViewListOfStaleIssuesByTypeEvent(itype IssueType) wf.Event {
	return eventByType("VLOfStaleIssuesByType", itype)
}
func ViewListOfAdvocacyIssuesByTypeEvent(itype IssueType) wf.Event {
	return eventByType("ViewListOfAdvocacyIssuesByTypeEvent", itype)
}

const MessagePostedState wf.State = "MessagePostedState"
const (
	EditEvent                 wf.Event = "EditEvent"
	AddAnotherEvent           wf.Event = "AddAnotherEvent"
	DetailsEvent              wf.Event = "DetailsEvent"
	CancelEvent               wf.Event = "CancelEvent"
	ProgressShowEvent         wf.Event = "ProgressShowEvent"
	ProgressIntermediateEvent wf.Event = "ProgressIntermediateEvent"
	ProgressCloseoutEvent     wf.Event = "ProgressCloseoutEvent"
)

func MessageIDAvailableEventInContext(context string) wf.Event { 
	return wf.Event("MessageIDAvailableEventInContext(" + context+")")
}

const FormShownState wf.State = "FormShownState"

const CommunitySelectingState wf.State = "CommunitySelectingState"
const (
	CommunitySelectedEvent wf.Event = "CommunitySelectedEvent"
)
const ObjectiveShownState wf.State = "ObjectiveShownState"
const ProgressFormShownState wf.State = "ProgressFormShownState"
const DoneState wf.State = "DoneState"

// Workflow is a public interface of workflow template.
type Workflow interface {
	GetNamedTemplate() wf.NamedTemplate
}

// this can only be created using constructor function. Thus we can guarantee that
// all fields will have values.
type workflowImpl struct {
	DynamoDBConnection
	// IssueDAO
	// IssueProgressDAO
	// AdaptiveCommunityDAO
	// UserDAO
	// CompetencyDAO
	// StrategyObjectiveDAO
	// StrategyInitiativeDAO
	// StrategyCommunityDAO
	// StrategyInitiativeCommunityDAO

	// CapabilityCommunityDAO

	DialogFetcherDAO dialogFetcher.DAO

	// Queries
	alog.AdaptiveLogger
	// UserHasWriteAccessToIssues
	// /*
	// 	itype := getIssueTypeFromContext(ctx)
	// 	typLabel := string(itype)
	// 	mc := models.MessageCallback{ // TODO: generate the correct MessageCallback for closeoutEng
	// 		Module: "objectives",
	// 		Target: uo.ID,
	// 		Source: uo.UserID,
	// 		Action: "ask_closeout", // will be replaced with `closeout`
	// 		Topic:  "init",
	// 	}
	// 	// send a notification to the partner
	// 	objectives.ObjectiveCloseoutEng(engagementTable, mc,
	// 		uo.AccountabilityPartner,
	// 		fmt.Sprintf("<@%s> wants to close the following %s. Are you ok with that?", ctx.Request.User.ID, typLabel),
	// 		fmt.Sprintf("*%s*: %s \n *%s*: %s",
	// 			NameLabel, uo.Name,
	// 			DescriptionLabel, uo.Description),
	// 		string(closeoutLabel(uo.ID)), objectiveCloseoutPath, false, dns, common2.EngagementEmptyCheck, ctx.PlatformID)
	// */
	// OnCloseout func(issue Issue) (err error)
}

// // WorkflowEnvironment - All external dependencies
// type WorkflowEnvironment struct {
// 	UserHasWriteAccessToIssues
// 	SelectFromCapabilityCommunityJoinStrategyCommunityWhereChannelCreated
// 	CommunityById // func(community.AdaptiveCommunity, models.PlatformID) community.AdaptiveCommunity
// }

func CreateIssueWorkflow(
	conn DynamoDBConnection,
	logger alog.AdaptiveLogger,
) Workflow {
	logger.Infoln("IssueWorkflow")
	return CreateWorkflowImpl(logger)(conn)
}

func (w workflowImpl) GetNamedTemplate() wf.NamedTemplate {
	nt := wf.NamedTemplate{
		Name: IssuesNamespace,
		Template: wf.Template{
			Init: InitState, // initial state is "init". This is used when the user first triggers the workflow
			FSA: map[struct {
				wf.State
				wf.Event
			}]wf.Handler{
				{State: InitState, Event: CreateIssueByTypeEvent(IDO)}:          w.OnCreateItem(true, IDO),
				{State: InitState, Event: ""}:                                   w.OnCreateItem(true, IDO),// TODO: remove after integration period
				{State: InitState, Event: CreateIssueByTypeEvent(SObjective)}:   w.OnCreateItem(true, SObjective),
				{State: InitState, Event: CreateIssueByTypeEvent(Initiative)}:   w.OnCreateItem(true, Initiative),
				{State: CommunitySelectingState, Event: CommunitySelectedEvent}: wf.SimpleHandler(w.OnCommunitySelected, FormShownState),
				{State: FormShownState, Event: wf.SurveySubmitted}:              wf.SimpleHandler(w.OnDialogSubmitted, MessagePostedState),
				
				{State: MessagePostedState, Event: MessageIDAvailableEventInContext(DescriptionContext)}:              wf.SimpleHandler(w.OnFieldsShown(ExtractDescription, DescriptionContext), MessagePostedState), // returning to the same state for other events to trigger
				{State: MessagePostedState, Event: MessageIDAvailableEventInContext(CloseoutAgreementContext)}:        wf.SimpleHandler(w.OnFieldsShown(ExtractDescription, CloseoutAgreementContext), MessagePostedState), // returning to the same state for other events to trigger
				{State: MessagePostedState, Event: MessageIDAvailableEventInContext(CloseoutDisagreementContext)}:     wf.SimpleHandler(w.OnFieldsShown(ExtractDescription, CloseoutDisagreementContext), MessagePostedState), // returning to the same state for other events to trigger
				{State: MessagePostedState, Event: MessageIDAvailableEventInContext(UpdateContext)}:                   wf.SimpleHandler(w.OnFieldsShown(ExtractDescription, UpdateContext), MessagePostedState), // returning to the same state for other events to trigger
				{State: MessagePostedState, Event: MessageIDAvailableEventInContext(UpdateResponseContext)}:           wf.SimpleHandler(w.OnFieldsShown(ExtractDescription, UpdateResponseContext), MessagePostedState), // returning to the same state for other events to trigger
				{State: MessagePostedState, Event: MessageIDAvailableEventInContext(CoachingRequestRejectionContext)}: wf.SimpleHandler(w.OnFieldsShown(ExtractDescription, CoachingRequestRejectionContext), MessagePostedState), // returning to the same state for other events to trigger
				{State: MessagePostedState, Event: MessageIDAvailableEventInContext(ProgressUpdateContext)}:           wf.SimpleHandler(w.OnFieldsShown(ExtractDescription, ProgressUpdateContext), MessagePostedState), // returning to the same state for other events to trigger

				// the following events are for buttons. Will be invoked not immediately
				{State: MessagePostedState, Event: EditEvent}:                 wf.SimpleHandler(w.OnEdit, FormShownState),
				{State: MessagePostedState, Event: DetailsEvent}:              wf.SimpleHandler(w.OnDetails, MessagePostedState),
				{State: MessagePostedState, Event: CancelEvent}:               wf.SimpleHandler(w.OnProgressCancel, DoneState),
				{State: MessagePostedState, Event: ProgressShowEvent}:         wf.SimpleHandler(w.OnProgressShow, MessagePostedState),
				{State: MessagePostedState, Event: ProgressIntermediateEvent}: wf.SimpleHandler(w.OnProgressIntermediate, ProgressFormShownState),
				{State: MessagePostedState, Event: ProgressCloseoutEvent}:     wf.SimpleHandler(w.OnProgressCloseout, MessagePostedState),
				// {State: MessagePostedState, Event: "delete"}: wf.SimpleHandler(OnDelete, DoneState),
				// {State: MessagePostedState, Event: AddAnotherEvent}:             w.OnCreateItem(false),
				{State: ProgressFormShownState, Event: wf.SurveySubmitted}:     wf.SimpleHandler(w.OnProgressFormSubmitted, MessagePostedState),
				{State: ProgressFormShownState, Event: wf.SurveyCancelled}:     wf.SimpleHandler(w.OnDialogCancelled, DoneState), // NB! we handle on cancel using the same method
				{State: FormShownState, Event: wf.SurveyCancelled}:             wf.SimpleHandler(w.OnDialogCancelled, DoneState),
				{State: InitState, Event: ViewListOfIssuesByTypeEvent(IDO)}:    wf.SimpleHandler(w.OnViewListOfIssues(IDO, unfiltered), ObjectiveShownState),
				// {State: InitState, Event: "view-idos"                     }:    wf.SimpleHandler(w.OnViewListOfIssues(IDO, unfiltered), ObjectiveShownState),// TODO: remove after integration period
				{State: InitState, Event: ViewMyListOfIssuesByTypeEvent(IDO)}:  wf.SimpleHandler(w.OnViewListOfIssues(IDO, filterUserAdvocate), ObjectiveShownState),
				{State: InitState, Event: ViewListOfIssuesByTypeEvent(SObjective)}:    wf.SimpleHandler(w.OnViewListOfIssues(SObjective, unfiltered), ObjectiveShownState),
				{State: InitState, Event: ViewMyListOfIssuesByTypeEvent(SObjective)}:  wf.SimpleHandler(w.OnViewListOfIssues(SObjective, filterUserAdvocate), ObjectiveShownState),
				{State: InitState, Event: ViewListOfIssuesByTypeEvent(Initiative)}:    wf.SimpleHandler(w.OnViewListOfIssues(Initiative, unfiltered), ObjectiveShownState),
				{State: InitState, Event: ViewMyListOfIssuesByTypeEvent(Initiative)}:  wf.SimpleHandler(w.OnViewListOfIssues(Initiative, filterUserAdvocate), ObjectiveShownState),
				{State: InitState, Event: ViewListOfStaleIssuesByTypeEvent(IDO)}:             wf.SimpleHandler(w.OnViewListOfQueryIssues(IDO,        StaleObjectivesQuery, "Stale Individual Development Objectives"), ObjectiveShownState),
				{State: InitState, Event: ViewListOfStaleIssuesByTypeEvent(SObjective)}:      wf.SimpleHandler(w.OnViewListOfQueryIssues(SObjective, StaleObjectivesQuery, "Stale Objectives"), ObjectiveShownState),
				{State: InitState, Event: ViewListOfStaleIssuesByTypeEvent(Initiative)}:      wf.SimpleHandler(w.OnViewListOfQueryIssues(Initiative, StaleObjectivesQuery, "Stale Initiatives"), ObjectiveShownState),

				{State: InitState, Event: ViewListOfAdvocacyIssuesByTypeEvent(IDO)}:          wf.SimpleHandler(w.OnViewListOfQueryIssues(IDO,        AdvocacyIssuesQuery, "Individual Development Objectives you are coaching"), ObjectiveShownState),
				{State: InitState, Event: ViewListOfAdvocacyIssuesByTypeEvent(SObjective)}:   wf.SimpleHandler(w.OnViewListOfQueryIssues(SObjective, AdvocacyIssuesQuery, "Objectives you are an advocate for"),  ObjectiveShownState),
				{State: InitState, Event: ViewListOfAdvocacyIssuesByTypeEvent(Initiative)}:   wf.SimpleHandler(w.OnViewListOfQueryIssues(Initiative, AdvocacyIssuesQuery, "Initiatives you are an advocate for"), ObjectiveShownState),

				{State: ObjectiveShownState, Event: EditEvent}:                 wf.SimpleHandler(w.OnEdit, FormShownState),
				{State: ObjectiveShownState, Event: DetailsEvent}:              wf.SimpleHandler(w.OnDetails, MessagePostedState),
				{State: ObjectiveShownState, Event: CancelEvent}:               wf.SimpleHandler(w.OnProgressCancel, DoneState),
				{State: ObjectiveShownState, Event: ProgressShowEvent}:         wf.SimpleHandler(w.OnProgressShow, MessagePostedState),
				{State: ObjectiveShownState, Event: ProgressIntermediateEvent}: wf.SimpleHandler(w.OnProgressIntermediate, ProgressFormShownState),
				{State: ObjectiveShownState, Event: ProgressCloseoutEvent}:     wf.SimpleHandler(w.OnProgressCloseout, MessagePostedState),
				// {State: ObjectiveShownState, Event: AddAnotherEvent}:            w.OnCreateItem(false),
			},
			Parser: wf.Parser,
		}}
	return nt
}
