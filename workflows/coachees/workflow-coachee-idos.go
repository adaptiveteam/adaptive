package coachees

import (
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
	"github.com/adaptiveteam/adaptive/daos/userObjectiveProgress"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"log"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/workflows/exchange"
)

// ViewCoacheesWorkflow - 
var ViewCoacheesWorkflow = exchange.WorkflowInfo{
	Prefix: exchange.CoachingPath,
	Name: Namespace, Init: InitState}

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

// CreateCoacheesWorkflow - constructor.
func CreateCoacheesWorkflow(
	conn DynamoDBConnection,
	logger alog.AdaptiveLogger,
) Workflow {
	logger.Infoln("CoacheesWorkflow")

	return workflowImpl{
		DynamoDBConnection: conn,
		AdaptiveLogger:     logger,
	}
}
const Namespace = exchange.CoacheesNamespace
const InitState wf.State = "init"
// // ViewCoacheeIDOs is a workflow to show list of Coachee idos.
// var ViewCoacheeIDOs = wf.NamedTemplate{
// 	Name: "view-coachee-idos", Template: ViewCoacheeIDOs_Workflow(),
// }

const ShowUpdatesEvent wf.Event = "show-updates"

const ItemShownState wf.State = "item-shown"

const CoacheeListPretext ui.RichText = "The list of coachees:"
const AdvocateListPretext ui.RichText = "The list of advocates:"

const NameLabel = "Name"
const DescriptionLabel = "Description"

const StatusColorLabel = "Status"
const CommentsLabel = "Comments"

const itemIDKey = "itemID"
const progressDateKey = "progressDate"

func (w workflowImpl) GetNamedTemplate() (nt wf.NamedTemplate) {
	nt = wf.NamedTemplate{
		Name: Namespace,
		Template: wf.Template{
			Init: "init",
			FSA: map[struct {
				wf.State
				wf.Event
			}]wf.Handler{
				{State: InitState, Event: ""}:                    wf.SimpleHandler(w.ViewCoacheeIDOs_OnInit(), ItemShownState),
				{State: ItemShownState, Event: ShowUpdatesEvent}: wf.SimpleHandler(w.ViewCoacheeIDOs_OnShowUpdates, "updates-shown"),
			},
			Parser: wf.Parser,
		},
	}
	return
}

// ObjectiveGroup is used to group objectives
type ObjectiveGroup struct {
	GroupID  string
	Elements []userObjective.UserObjective
}

func userIDKey(obj userObjective.UserObjective) string { return obj.UserID }

// ObjectivesGroupByUserID groups objectives by user id and renders each objective as RichText
func ObjectivesGroupByUserID(objectives []userObjective.UserObjective, key func(userObjective.UserObjective) string) (groupedObjectives []ObjectiveGroup) {
	groups := make(map[string]ObjectiveGroup, 0)
	for _, objective := range objectives {
		k := key(objective)
		group, ok := groups[k]
		if !ok {
			group = ObjectiveGroup{GroupID: k}
		}
		group.Elements = append(group.Elements, objective)
		groups[k] = group
	}

	groupedObjectives = make([]ObjectiveGroup, 0, len(groups))

	for _, value := range groups {
		groupedObjectives = append(groupedObjectives, value)
	}
	return
}

func (w workflowImpl) ViewCoacheeIDOs_OnInit() func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		log.Println("ViewCoacheeIDOs_OnInit")
		coachID := ctx.Request.User.ID
		userObjectiveDAO := userObjective.NewDAO(w.DynamoDBConnection.Dynamo, "ViewCoacheeIDOs_OnInit", w.ClientID)
		objectives := userObjectiveDAO.ReadByAccountabilityPartnerUnsafe(coachID)
		coacheeObjectives := filterObjectivesByObjectiveTypeNotCompletedAccepted(objectives, userObjective.IndividualDevelopmentObjective)
		coacheeGroups := ObjectivesGroupByUserID(coacheeObjectives, userIDKey)

		out.Interaction = wf.Interaction{
			Messages: wf.InteractiveMessages(wf.InteractiveMessage{
				PassiveMessage: wf.PassiveMessage{Text: CoacheeListPretext},
				Thread:         renderGroupsAsInteractiveMessages(coacheeGroups),
			}),
		}
		out.NextState = ItemShownState
		return
	}
}

func objectiveFields(item userObjective.UserObjective) (fields []ebm.AttachmentField) {
	fields = []ebm.AttachmentField{
		{Title: NameLabel, Value: item.Name},
		{Title: DescriptionLabel, Value: item.Description},
	}
	return
}

func viewObjective(item userObjective.UserObjective) (message wf.InteractiveMessage) {
	fields := objectiveFields(item)
	showUpdates := wf.Button(ShowUpdatesEvent, "Show updates")
	return wf.InteractiveMessage{
		PassiveMessage:      wf.PassiveMessage{Fields: fields},
		InteractiveElements: wf.InteractiveElements(showUpdates),
		DataOverride:        wf.Data{itemIDKey: item.ID},
	}
}

func renderGroupsAsInteractiveMessages(coacheeGroups []ObjectiveGroup) (messages []wf.InteractiveMessage) {
	for _, group := range coacheeGroups {
		// TODO: group.Sort() 
		threadMessages := wf.InteractiveMessages()
		for _, item := range group.Elements {
			threadMessages = append(threadMessages, viewObjective(item))
		}
		messages = append(messages, wf.InteractiveMessage{
			PassiveMessage: wf.PassiveMessage{Text: ui.Sprintf("<@%s>", group.GroupID)},
		})
		messages = append(messages, threadMessages...)
	}
	return
}

// func onViewCoacheeIDOsOld(coachID string) []platform.Response {
// 	objectives := userObjectiveDAO.ReadByAccountabilityPartnerUnsafe(coachID)
// 	coacheeObjectives := filterObjectivesByObjectiveTypeNotCompletedAccepted(objectives, userObjective.IndividualDevelopmentObjective)
// 	coacheeInfos := GroupByUserID(coacheeObjectives, FormatObjectiveItem)
// 	names := renderGroupsAsIDAndElementsAsSubList(coacheeInfos)
// 	text := CoacheeListPretext + "\n" + ui.ListItems(names...)
// 	return []platform.Response{
// 		platform.Post(platform.ConversationID(coachID), platform.MessageContent{Message: text}),
// 	}
// }
func showProgress(progress userObjectiveProgress.UserObjectiveProgress) wf.InteractiveMessage {
	return wf.InteractiveMessage{
		PassiveMessage: wf.PassiveMessage{Fields: []ebm.AttachmentField{
			{Title: StatusColorLabel, Value: string(models.ObjectiveStatusColorLabels[progress.StatusColor])},
			{Title: CommentsLabel, Value: progress.Comments},
		}},
		DataOverride: wf.Data{itemIDKey: progress.ID, progressDateKey: progress.CreatedOn},
	}
}

var userObjectiveProgressTableName             = func(clientID string) string { return clientID + "_user_objectives_progress" }

// ViewCoacheeIDOs_OnShowUpdates -
func (w workflowImpl) ViewCoacheeIDOs_OnShowUpdates(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	itemID := ctx.Data[itemIDKey]
	w.AdaptiveLogger.Infof("ViewCoacheeIDOs_OnShowUpdates itemID=%s", itemID)
	userObjectiveDAO := userObjective.NewDAO(w.DynamoDBConnection.Dynamo, "ViewCoacheeIDOs_OnShowUpdates", w.ClientID)
	userObjectiveProgressDAO := userObjectiveProgress.NewDAOByTableName(w.DynamoDBConnection.Dynamo, "ViewCoacheeIDOs_OnShowUpdates", userObjectiveProgressTableName(w.ClientID))
	var items []userObjective.UserObjective
	items, err = userObjectiveDAO.ReadOrEmpty(itemID)
	if err == nil {
		if len(items) > 0 {
			item := items[0]
			var progress []userObjectiveProgress.UserObjectiveProgress
			progress, err = userObjectiveProgressDAO.ReadByID(item.ID)

			threadMessages := wf.InteractiveMessages()
			if err == nil {
				for _, p := range progress {
					threadMessages = append(threadMessages, showProgress(p))
				}
			}
			var title string
			if len(progress) == 0 {
				title = "There are no updates"
			} else {
				title = "You can find the status updates of this Individual Development Objective in the thread. :point_down:"
			}
			out.Interaction.Messages = []wf.InteractiveMessage{
				{
					PassiveMessage: wf.PassiveMessage{
						Footer: ebm.AttachmentFooter{Text: title},
						Fields: objectiveFields(item),
					},
					Thread: threadMessages,
				},
			}
		} else {
			w.AdaptiveLogger.Warnf("Couldn't find objective %s", itemID)
		}
	}
	out.KeepOriginal = true // this event comes from item view.
	return
}

func filterObjectivesByObjectiveTypeNotCompletedAccepted(objectives []userObjective.UserObjective, objectiveType userObjective.DevelopmentObjectiveType) (res []userObjective.UserObjective) {
	for _, objective := range objectives {
		if objective.ObjectiveType == objectiveType &&
		   objective.Completed == 0 &&
		   objective.Accepted == 1 {
			res = append(res, objective)
		}
	}
	return
}