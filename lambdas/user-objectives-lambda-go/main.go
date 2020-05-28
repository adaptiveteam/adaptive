package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	"github.com/adaptiveteam/adaptive/daos/strategyInitiative"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	"github.com/adaptiveteam/adaptive/daos/strategyObjective"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	daosUser "github.com/adaptiveteam/adaptive/daos/user"
	"github.com/pkg/errors"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	issues "github.com/adaptiveteam/adaptive/workflows/issues"
	common2 "github.com/adaptiveteam/adaptive/adaptive-engagement-scheduling/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/values"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsIssues "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	wfRoutes "github.com/adaptiveteam/adaptive/workflows"
	
)

const (
	AdaptiveDateFormat         = core.ISODateLayout
	AdaptiveReadableDateFormat = core.USDateLayout
	DateFormat                 = string(AdaptiveDateFormat)
	ReadableDateFormat         = string(AdaptiveReadableDateFormat)
)

const (
	IDODescriptionContext    = "dialog/ido/language-coaching/description"
	IDOProgressUpdateContext = "dialog/ido/language-coaching/update"

	IDOCloseoutDisagreementContext                 = "dialog/ido/language-coaching/close-out-disagreement"
	InitiativeCloseoutDisagreementContext          = "dialog/strategy/language-coaching/initiative/close-out-disagreement"
	CapabilityObjectiveCloseoutDisagreementContext = "dialog/strategy/language-coaching/objective/close-out-disagreement"

	IDOCloseoutAgreementContext                 = "dialog/ido/language-coaching/close-out-agreement"
	InitiativeCloseoutAgreementContext          = "dialog/strategy/language-coaching/initiative/close-out-agreement"
	CapabilityObjectiveCloseoutAgreementContext = "dialog/strategy/language-coaching/objective/close-out-agreement"

	IDOCoachingRejectionContext       = "dialog/ido/language-coaching/coaching-request-rejection"
	IDOResponseObjectiveUpdateContext = "dialog/ido/language-coaching/update-response"

	CapabilityObjectiveProgressUpdateContext = "dialog/strategy/language-coaching/objective/update"
	InitiativeProgressUpdateContext          = "dialog/strategy/language-coaching/initiative/update"
	CapabilityObjectiveUpdateResponseContext = "dialog/strategy/language-coaching/objective/update-response"
	InitiativeUpdateResponseContext          = "dialog/strategy/language-coaching/initiative/update-response"
	BlueDiamondEmoji                         = ":small_blue_diamond:"
)

const (
	ObjectiveStatusColor       = "objective_status_color"
	ObjectiveCloseoutComment   = "objective_closeout_comment"
	ObjectiveNoCloseoutComment = "objective_no_closeout_comment"

	ViewMore     models.AttachActionName = "view more"
	ViewLess     models.AttachActionName = "view less"
	ViewProgress models.AttachActionName = "view progress"
	Closeout     models.AttachActionName = "closeout"
	MoreOptions  models.AttachActionName = "more_options"
	HomeOptions  models.AttachActionName = "home_options"
	No           models.AttachActionName = "no"
	Cancel       models.AttachActionName = "cancel"
	Enable       models.AttachActionName = "enable"
	Confirm      models.AttachActionName = "confirm"
)
var	ViewOpenObjectives = common2.ViewOpenObjectives
const(
	Individual          = "Individual Objective"
	CapabilityObjective = "Capability Objective"
	StrategyInitiative  = "Initiative"
	FinancialObjective  = "Financial Objective"
	CustomerObjective   = "Customer Objective"
)

func IDOCoaches(userID string, teamID models.TeamID, conn daosCommon.DynamoDBConnection) []models.KvPair {
	var filteredCoaches []models.KvPair
	coaches := objectives.IDOCoaches(userID, teamID,
		communityUsersTable, string(adaptiveCommunityUser.PlatformIDCommunityIDIndex), 
		utilsUser.UserIDsToDisplayNamesConnUnsafe(conn))
	for _, each := range coaches {
		if each.Value != userID {
			// Filtering out the self-user
			filteredCoaches = append(filteredCoaches, each)
		}
	}
	return filteredCoaches
}

// ---------------------------------------------------------------------------------------------------------------------
// COMMON METHODS

// Listing all the objectives for a user
func ListObjectives(userID, channelID string, objectives []models.UserObjective, fn UserObjectiveAttachments, threadTs string) {
	// Display each objective
	for _, each := range objectives {
		year, month := core.CurrentYearMonth()
		mc := models.MessageCallback{Module: "objectives", Source: userID, Target: each.ID, Topic: "init", Action: "ask",
			Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
		publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Attachments: fn(mc, each), ThreadTs: threadTs})
	}
}

// TimeStamp extracts timestamp from the original message
// When the original message is from a thread, we need to post to the same thread
// Below logic checks if the incoming message is from a thread
func TimeStamp(request slack.InteractionCallback) string {
	ts := request.OriginalMessage.ThreadTimestamp
	if ts == "" {
		ts = request.MessageTs
	}
	return ts
}

func createObjectiveNow(conn daosCommon.DynamoDBConnection, message slack.InteractionCallback, userID string,
	actionValue string, mc *models.MessageCallback,
) {
	teamID := models.ParseTeamID(conn.PlatformID)
	id := actionValue // callbackId
	mc.Set("Target", "")
	pValues := platformValues(teamID)
	logger.Infof("Retrieved values for %s platform: %v", teamID, pValues)
	initsObjsValues := append(InitsAndObjectives(userID, teamID),
		pValues...)
	survey := utils.AttachmentSurvey(string(ObjectiveAddAnotherDialogTitle),
		objectives.EditObjectiveSurveyElems2(nil, IDOCoaches(userID, teamID, conn),
			objectives.DevelopmentObjectiveDates(namespace, ""), initsObjsValues))
	api := getSlackClient(message)
	// Open a survey associated with the engagement
	marshaledSurvey, _ := json.Marshal(survey)
	logger.Infof("Survey for ObjectiveCreate with %s id: %s", id, string(marshaledSurvey))
	err2 := dialogFromSurvey(api, message, survey, id, false, "")
	core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not open dialog from %s survey", id+":"+message.CallbackID))
}

func PartnerSelectingUserEngagement(conn daosCommon.DynamoDBConnection, 
	mc models.MessageCallback, userID string, text ui.RichText, users []string) {
	var userOpts []ebm.MenuOption

	teamID := models.ParseTeamID(conn.PlatformID)
	for _, each := range users {
		var users1 []models.User
		var err2 error
		users1, err2 = daosUser.ReadOrEmpty(each)(conn)
		core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not query for %s user", userID))
		for _, user := range users1 {
			userOpts = append(userOpts, ebm.Option(each, ui.PlainText(user.DisplayName)))
		}
	}
	utils.AddChatEngagement(mc, "",
		string(text), string(PartnerSelectingUserEngagementFallbackText),
		userID, []ebm.AttachmentAction{*models.SelectAttachAction(mc, models.Select,
			string(PickAUserMenuPrompt), userOpts,
			models.EmptyActionMenuOptionGroups())},
		[]ebm.AttachmentField{}, teamID, true, engagementTable, d, namespace,
		time.Now().Unix(), common2.EngagementEmptyCheck)
}

// ---------------------------------------------------------------------------------------------------------------------

func detailedViewFields(conn daosCommon.DynamoDBConnection, u *models.UserObjective, typ, alignment string) []ebm.AttachmentField {
	fields := compactViewFields(u, typ, alignment)
	// For ViewMore action, we only need the latest comment
	ops := LatestProgressUpdateByObjectiveID(u.ID)

	var comments []ui.RichText
	for _, each := range ops {
		comments = append(comments, ui.Sprintf("%s (%s percent, %s status)", each.Comments, each.PercentTimeLapsed, each.StatusColor))
	}

	var status ui.PlainText
	if u.Cancelled == 1 {
		status = StatusCancelled
	} else if u.Completed == 0 {
		status = StatusPending
	} else if u.Completed == 1 && u.PartnerVerifiedCompletion {
		status = StatusCompletedAndPartnerVerifiedCompletion
	} else if u.Completed == 1 && !u.PartnerVerifiedCompletion {
		status = StatusCompletedAndNotPartnerVerifiedCompletion
	}
	user, err2 := daosUser.Read(u.AccountabilityPartner)(conn)
	core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not read %s user", u.AccountabilityPartner))

	fields = append(fields, models.AttachmentFields([]models.KvPair{
		{
			Key:   string(AccountabilityPartnerLabel),
			Value: user.DisplayName,
		},
		{
			Key:   string(StatusLabel),
			Value: string(status),
		},
		{
			Key:   string(LastReportedProgressLabel),
			Value: string(JoinRichText(comments, "\n")),
		},
	})...)
	return fields
}

func showProgressField(u *models.UserObjective) (af []ebm.AttachmentField) {
	ops, err := userObjectiveProgressByID(u.ID, -1)

	var comments []ui.RichText
	for _, each := range ops {
		comments = append(comments, ui.Sprintf("%s (%s)", each.Comments, each.CreatedOn))
	}

	if err == nil {
		commentsJoined := ui.ListItems(comments...)
		progressTitle := ProgressTitle(*u)
		if commentsJoined == "" {
			progressTitle = ProgressAbsentTitle(*u)
		}
		af = []ebm.AttachmentField{
			{
				Title: string(progressTitle),
				Value: string(commentsJoined),
				Short: true,
			},
		}
	} else {
		logger.WithField("error", err).Errorf("Could not get progress updates for %s objective", u.ID)
	}
	return
}

func renderObjectiveViewDate(userObj models.UserObjective) ui.PlainText {
	return ui.PlainText(common.ObjectiveDate{
		CreatedDate:     userObj.CreatedDate,
		ExpectedEndDate: userObj.ExpectedEndDate,
	}.Render(AdaptiveDateFormat, AdaptiveReadableDateFormat, namespace))
}

func compactViewFields(userObj *models.UserObjective, typ, alignment string) []ebm.AttachmentField {
	kvs := []models.KvPair{
		{
			Key:   NameLabel,
			Value: userObj.Name,
		},
		{
			Key:   DescriptionLabel,
			Value: userObj.Description,
		},
		{
			Key:   TimelineLabel,
			Value: string(renderObjectiveViewDate(*userObj)),
		},
	}
	if typ != core.EmptyString {
		kvs = append(kvs, models.KvPair{
			Key:   "Type",
			Value: typ,
		})
	}
	if alignment != core.EmptyString {
		kvs = append(kvs, models.KvPair{
			Key:   string(StrategyAssociationFieldLabel),
			Value: alignment,
		})
	}
	return models.AttachmentFields(kvs)
}

func progressFields(comments ui.PlainText, status models.ObjectiveStatusColor, obj models.UserObjective) []ebm.AttachmentField {
	today := time.Now().Format(DateFormat)
	timeProgressLabel := ObjectiveProgressText(obj, today)
	return models.AttachmentFields([]models.KvPair{
		{
			Key:   NameLabel,
			Value: string(timeProgressLabel),
		},
		{
			Key:   DescriptionLabel,
			Value: obj.Description,
		},
		// {
		//	Key:   "Strategy Association(s)",
		//	Value: objectiveType(obj),
		// },
		{
			Key:   string(ProgressStatusLabel),
			Value: string(models.ObjectiveStatusColorLabels[status]),
		},
		{
			Key:   string(CommentsLabel),
			Value: string(comments),
		},
	})
}

func parseGwSNSRequest(np models.NamespacePayload) (slackevents.EventsAPIEvent, error) {
	ueRequest, err := url.QueryUnescape(np.Payload)
	requestPayload := strings.Replace(ueRequest, "payload=", "", -1)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not un-escape the request body"))

	return utils.ParseApiRequest(requestPayload)
}

func publish(msg models.PlatformSimpleNotification) {
	_, err := sns.Publish(msg, platformNotificationTopic)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not pusblish message to %s topic", platformNotificationTopic))
}

func progressCommentSurveyElements(objName ui.PlainText, startDate string) []ebm.AttachmentActionTextElement {
	nameConstrained := ObjectiveCommentsTitle(objName)
	elapsedDays := common.DurationDays(startDate, TodayISOString(), AdaptiveDateFormat, namespace)
	return []ebm.AttachmentActionTextElement{
		{
			Label:    string(ObjectiveStatusLabel(elapsedDays, startDate)),
			Name:     ObjectiveStatusColor,
			ElemType: models.MenuSelectType,
			Options:  utils.AttachActionElementOptions(models.ObjectiveStatusColorKeyValues),
			// Value:    string(models.ObjectiveStatusColorLabels[statusValue]), // it's necessary to fill afterwards
		},
		ebm.NewTextArea(ObjectiveProgressComments, nameConstrained, ObjectiveProgressCommentsPlaceholder, ""),
	}
}

func updateAttachmentActions(mc models.MessageCallback, teamID models.TeamID,
	obj *models.UserObjective,
	showMore, showProgress, strategyFlag bool) []ebm.AttachmentAction {
	var actions []ebm.AttachmentAction
	typeLabel := objectiveTypeLabel(*obj)
	if showMore {
		actions = append(actions,
			*models.SimpleAttachAction(*mc.WithTarget(obj.ID), ViewLess, PrevPageOfOptionsActionLabel))
	} else if !strategyFlag {
		// Show `Show details` button only for individual objectives
		actions = append(actions,
			*models.SimpleAttachAction(*mc.WithTarget(obj.ID), ViewMore, NextPageOfOptionsActionLabel))

	}
	if !strategyFlag {
		if obj.Cancelled == 1 {
			actions = append(actions,
				*models.GenAttachAction(*mc.WithTarget(obj.ID), Enable,
					string(CancelledObjectiveActivateActionLabel),
					utils.AttachmentConfirm(
						string(CancelledObjectiveActivateConfirmationDialogTitle),
						string(cancelledObjectiveActivateConfirmationDialogText(typeLabel))), true))
		} else if obj.Completed == 0 {
			// Show below options only if objective is still in progress
			actions = append(actions,
				*models.DialogAttachAction(*mc.WithTarget(obj.ID),
					models.Update,
					ObjectiveModifyActionLabel),
				*models.SimpleAttachAction(*mc.WithTarget(obj.ID).WithTopic(CommentsName), models.Now,
					ObjectiveAddProgressInfoActionLabel),
				*models.GenAttachAction(*mc.WithTarget(obj.ID).WithTopic("init"), Cancel,
					string(ObjectiveCancelActionLabel),
					utils.AttachmentConfirm(
						string(ObjectiveCancellationConfirmationDialogTitle),
						string(objectiveCancellationConfirmationDialogText(typeLabel))), true),
				*models.SimpleAttachAction(*mc.WithTarget(obj.ID).WithTopic("init").WithAction("ask"), MoreOptions,
					ObjectiveMoreOptionsActionLabel))
		}
	} else {
		actions = append(actions,
			*models.SimpleAttachAction(*mc.WithTarget(obj.ID).WithTopic(CommentsName), models.Now,
				ObjectiveAddProgressInfoActionLabel))
		// Do not show entire progress under progress view context
		if !showProgress {
			actions = append(actions, *models.SimpleAttachAction(*mc.WithTarget(obj.ID).WithTopic("init"), ViewProgress,
				ObjectiveShowEntireProgressActionLabel))
		}
		actions = append(actions,
			*models.GenAttachAction(mc, Closeout,
				string(ObjectiveCloseoutActionLabel),
				utils.AttachmentConfirm(
					string(ObjectiveCloseoutConfirmationDialogTitle),
					string(objectiveCloseoutConfirmationDialogText(typeLabel))), true))
		// Show details only under progress view context
		if showProgress {
			actions = append(actions,
				*models.SimpleAttachAction(*mc.WithTopic("init"), HomeOptions,
					ObjectiveDetailsActionLabel))
		}
	}
	return actions
}

func moreOptionsAttachmentActions(mc models.MessageCallback, obj *models.UserObjective, teamID models.TeamID) []ebm.AttachmentAction {
	typeLabel := objectiveTypeLabel(*obj)
	var actions = []ebm.AttachmentAction{
		*models.SimpleAttachAction(*mc.WithTarget(obj.ID), ViewProgress, ObjectiveShowEntireProgressActionLabel),
	}
	if obj.Completed == 0 {
		actions = append(actions,
			*models.GenAttachAction(mc, Closeout,
				string(ObjectiveCloseoutActionLabel),
				utils.AttachmentConfirm(
					string(ObjectiveCloseoutConfirmationDialogTitle),
					string(objectiveCloseoutConfirmationDialogText(typeLabel))), true))
	}
	actions = append(actions,
		*models.DialogAttachAction(mc, models.Now, ObjectiveAddAnotherActionLabel),
		*models.SimpleAttachAction(mc, HomeOptions, ObjectiveLessOptionsActionLabel))
	return actions
}

func updateProgressAttachmentActions(mc models.MessageCallback, actionName models.AttachActionName, objName ui.PlainText, start, end string) []ebm.AttachmentAction {
	return []ebm.AttachmentAction{*models.GenAttachAction(mc, actionName,
		string(ObjectiveProgressChangeCommentsActionLabel), models.EmptyActionConfirm(), true)}
}

// An attachment that displays information about objectives
func updateObjAttachment(conn daosCommon.DynamoDBConnection, 
	mc models.MessageCallback, title, text, fallback ui.PlainText, 
	uObj *models.UserObjective, showMore, showProgress, strategyFlag bool) []ebm.Attachment {
	teamID := models.ParseTeamID(conn.PlatformID)
	var fields []ebm.AttachmentField
	typeLabel, alignment := objectiveType(teamID)(*uObj)
	fmt.Println(fmt.Sprintf("Objective type for %s: %s", uObj.Name, typeLabel))
	if uObj.Cancelled == 1 {
		// Objective is cancelled
	}
	if showProgress {
		fields = showProgressField(uObj)
	} else {
		fields = compactViewFields(uObj, typeLabel, alignment)
		if showMore {
			fields = detailedViewFields(conn, uObj, typeLabel, alignment)
		}
	}
	attach := utils.ChatAttachment(string(title), string(text), string(fallback), mc.ToCallbackID(),
		updateAttachmentActions(mc, teamID, uObj, showMore, showProgress, strategyFlag), fields, time.Now().Unix())
	return []ebm.Attachment{*attach}
}

func moreOptionsAttachment(mc models.MessageCallback, title, fallback ui.PlainText, userObj *models.UserObjective,
	teamID models.TeamID) []ebm.Attachment {
	attach := utils.ChatAttachment(string(title), core.EmptyString, string(fallback), mc.ToCallbackID(),
		moreOptionsAttachmentActions(mc, userObj, teamID), compactViewFields(userObj, core.EmptyString, core.EmptyString),
		time.Now().Unix())
	return []ebm.Attachment{*attach}
}

func viewProgressAttachment(mc models.MessageCallback, title, fallback ui.PlainText, comments ui.PlainText,
	status models.ObjectiveStatusColor, obj models.UserObjective, actionName models.AttachActionName) []ebm.Attachment {
	attach := utils.ChatAttachment(string(title), core.EmptyString, string(fallback), mc.ToCallbackID(),
		updateProgressAttachmentActions(mc, actionName, ui.PlainText(obj.Name),
			obj.CreatedDate, obj.ExpectedEndDate),
		progressFields(comments, status, obj), time.Now().Unix())
	return []ebm.Attachment{*attach}
}

type MsgState struct {
	Ts          string `json:"ts"`
	ThreadTs    string `json:"thread_ts"`
	Update      bool   `json:"update"`
	ObjectiveId string `json:"objective_id"`
}

func GetMsgStateUnsafe(request slack.InteractionCallback) (msgState MsgState) {
	err := json.Unmarshal([]byte(request.State), &msgState)
	core.ErrorHandler(err, namespace, "Couldn't unmarshal MsgState")
	return
}

func dialogFromSurvey(api *slack.Client, message slack.InteractionCallback, survey ebm.AttachmentActionSurvey,
	id string, update bool, objectiveId string) error {
	survState := func() string {
		// When the original message is from a thread, we need to post to the same thread
		// Below logic checks if the incoming message is from a thread
		var ts string
		if message.OriginalMessage.ThreadTimestamp == "" {
			ts = message.MessageTs
		} else {
			ts = message.OriginalMessage.ThreadTimestamp
		}
		msgStateBytes, err := json.Marshal(MsgState{Ts: message.MessageTs, ThreadTs: ts, Update: update, ObjectiveId: objectiveId})
		core.ErrorHandler(err, namespace, "Could not marshal MsgState")
		return string(msgStateBytes)
	}

	return utils.SlackSurvey(api, message, survey, id, survState)
}

type UserObjectiveAttachments func(mc models.MessageCallback, objective models.UserObjective) []ebm.Attachment

func fillCommentsSurveyValues(sur ebm.AttachmentActionSurvey, comments string, status models.ObjectiveStatusColor) ebm.AttachmentActionSurvey {
	return models.FillSurvey(sur, map[string]string{
		ObjectiveProgressComments: comments,
		ObjectiveStatusColor:      string(status),
	})
}

func DeleteOriginalEng(userID, channel, ts string) {
	utils.DeleteOriginalEng(userID, channel, ts, func(notification models.PlatformSimpleNotification) {
		publish(notification)
	})
}

type ObjectivePredicate func(models.UserObjective) bool

func filterObjectives(objs []models.UserObjective, predicate ObjectivePredicate) (filtered []models.UserObjective) {
	for _, each := range objs {
		if predicate(each) {
			filtered = append(filtered, each)
		}
	}
	return
}
func ListObjectivesWithEvaluation(userID, channelID string, fn ObjectivePredicate,
	uofn UserObjectiveAttachments, typ models.DevelopmentObjectiveType, threadTs string) {
	// Times in AWS are in UTC
	allObjs := objectives.AllUserObjectives(userID, userObjectivesTable, string(userObjective.UserIDTypeIndex), typ, 0)
	openObjs := filterObjectives(allObjs, fn)
	if len(openObjs) > 0 {
		ListObjectives(userID, channelID, allObjs, uofn, threadTs)
	}
}

func callback(userID, topic, action string) models.MessageCallback {
	// We are writing month rather than year in engagement because quarter can always be inferred from month
	year, month := core.CurrentYearMonth()
	mc := models.MessageCallback{Module: "objectives", Source: userID, Topic: topic, Action: action, Target: "", Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
	return mc
}

// Get a string field value from an interface
func GetFieldString(i interface{}, field string) string {
	// Create a value for the slice.
	v := reflect.ValueOf(i)
	// Get the field of the slice element that we want to set.
	f := v.FieldByName(field)
	// Get value
	return f.String()
}

func renderStrategyAssociations(prefix, field string, entities ...interface{}) string {
	var op string
	if len(entities) > 0 {
		var acc ui.RichText
		for _, entity := range entities {
			acc += ui.Sprintf("%s %s \n", BlueDiamondEmoji, GetFieldString(entity, field))
		}
		op = fmt.Sprintf("*%s* \n%s", prefix, acc)
	}
	return op
}

func objectiveType(teamID models.TeamID) func(uObj models.UserObjective) (typ string, alignment string) {
	return func(uObj models.UserObjective) (typ string, alignment string) {
		typ = "Not aligned with strategy"
		alignment = "Not aligned with any strategy"
		switch utilsIssues.DetectIssueType(uObj) {
		case issues.IDO:
			typ = Individual
			switch uObj.StrategyAlignmentEntityType {
			case models.ObjectiveStrategyObjectiveAlignment:
				capObj := strategy.StrategyObjectiveByID(teamID, uObj.StrategyAlignmentEntityID, strategyObjectivesTableName)
				alignment = renderStrategyAssociations("Objective", "Name", capObj)
			case models.ObjectiveStrategyInitiativeAlignment:
				initiative := strategy.StrategyInitiativeByID(teamID, uObj.StrategyAlignmentEntityID, strategyInitiativesTableName)
				alignment = renderStrategyAssociations("Initiative", "Name", initiative)
			case models.ObjectiveCompetencyAlignment:
				valueID := uObj.StrategyAlignmentEntityID
				values, err2 := valueDao.ReadOrEmpty(valueID)
				core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not read adaptive value %s ", valueID))
				for _, value := range values {
					alignment = renderStrategyAssociations("Competency", "Name", value)
				}
			default:
				log.Printf("WARN Unsupported uObj.StrategyAlignmentEntityType=%s\n", uObj.StrategyAlignmentEntityType)
			}
		case issues.SObjective:
			typ = CapabilityObjective
			splits := strings.Split(uObj.ID, "_")
			if len(splits) == 2 {
				so := strategy.StrategyObjectiveByID(teamID, splits[0], strategyObjectivesTableName)
				capComm := strategy.CapabilityCommunityByID(teamID, splits[1], capabilityCommunitiesTableName)
				alignment = fmt.Sprintf("%s%s",
					renderStrategyAssociations("Capability Communities", "Name", capComm),
					renderStrategyAssociations("Capability Objectives", "Name", so))
			} else {
				so := strategy.StrategyObjectiveByID(teamID, uObj.ID, strategyObjectivesTableName)
				alignment = fmt.Sprintf("`%s Objective` : `%s`\n", so.ObjectiveType, so.Name)
			}
		case issues.Initiative:
			typ = StrategyInitiative
			si := strategy.StrategyInitiativeByID(teamID, uObj.ID, strategyInitiativesTableName)
			initCommID := si.InitiativeCommunityID
			capObjID := si.CapabilityObjective
			initComm := strategy.InitiativeCommunityByID(teamID, initCommID, strategyInitiativeCommunitiesTable)
			capObj := strategy.StrategyObjectiveByID(teamID, capObjID, strategyObjectivesTableName)
			alignment = fmt.Sprintf("%s%s",
				renderStrategyAssociations("Initiative Communities", "Name", initComm),
				renderStrategyAssociations("Objectives", "Name", capObj))
		}
		return
	}
}

func CommentsSurvey(title ui.PlainText, elemLabel ui.PlainText, elemName string, value ui.PlainText) ebm.AttachmentActionSurvey {
	return utils.AttachmentSurvey(string(title), []ebm.AttachmentActionTextElement{
		ebm.NewTextArea(elemName, elemLabel, CommentsSurveyPlaceholder, value),
	})
}

func HandleRequest(ctx context.Context, e events.SNSEvent) (err error) {
	logger = logger.WithLambdaContext(ctx)
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("Error in user-objectives-lambda %v", err2)
		}
	}()
	for _, record := range e.Records {
		fmt.Println(record)
		message := record.SNS.Message

		if message == "warmup" {
			log.Println("Warmed up...")
		} else {
			var np models.NamespacePayload4
			err = json.Unmarshal([]byte(message), &np)
			core.ErrorHandler(err, namespace, fmt.Sprintf("Could not unmarshal sns record to NamespacePayload4"))
			teamID := np.TeamID
			conn := connGen.ForPlatformID(teamID.ToPlatformID())
			if strings.HasPrefix(np.InteractionCallback.CallbackID, "/") {
				err = wfRoutes.InvokeWorkflow(np, conn)
			} else {
				if np.Namespace == "objectives" {
					switch np.SlackRequest.Type {
					case models.InteractionSlackRequestType:
						fmt.Printf("### interactive_message event: %v\n", np)

						message := np.SlackRequest.InteractionCallback
						request := message
						action := message.ActionCallback.AttachmentActions[0]

						// TODO: Check if the request is coming from a user community or one of the strategy related communities

						// 'menu_list' is for the options that are presented to the user
						if action.Name == "menu_list" {
							err = onMenuList(np)
						} else if strings.HasPrefix(action.Name, "confirm") {
							onConfirm(request, teamID)
						} else if strings.HasPrefix(action.Name, "ask") {
							onAsk(request, teamID)
						} else if strings.HasPrefix(action.Name, string(Closeout)) {
							onCloseoutRequest(request, teamID)
						} else if strings.HasPrefix(action.Name, user.StaleIDOsForMe) {
							onViewStaleIDOs(request, teamID)
						} else if strings.HasPrefix(action.Name, ViewOpenObjectives) {
							onViewOpenObjectives(request, teamID)
						} else if strings.HasPrefix(action.Name, user.CapabilityObjectiveUpdateDueWithinWeek) ||
							strings.HasPrefix(action.Name, user.CapabilityObjectiveUpdateDueWithinMonth) ||
							strings.HasPrefix(action.Name, user.CapabilityObjectiveUpdateDueWithinQuarter) ||
							strings.HasPrefix(action.Name, user.CapabilityObjectiveUpdateDueWithinYear) ||
							strings.HasPrefix(action.Name, user.InitiativeUpdateDueWithinWeek) ||
							strings.HasPrefix(action.Name, user.InitiativeUpdateDueWithinMonth) ||
							strings.HasPrefix(action.Name, user.InitiativeUpdateDueWithinQuarter) ||
							strings.HasPrefix(action.Name, user.InitiativeUpdateDueWithinYear) ||
							strings.HasPrefix(action.Name, user.IDOUpdateDueWithinWeek) ||
							strings.HasPrefix(action.Name, user.IDOUpdateDueWithinMonth) ||
							strings.HasPrefix(action.Name, user.IDOUpdateDueWithinQuarter) ||
							strings.HasPrefix(action.Name, user.IDOUpdateDueWithinYear) {
							onCapabilityObjectiveUpdateDueWithinWeek(request, teamID)
						}
					case models.DialogSubmissionSlackRequestType:
						fmt.Printf("### dialog_submission event: %v\n", np)
						// Handling dialog submission for each answer
						dialog := np.SlackRequest.InteractionCallback
						onDialogSubmission(dialog, teamID)
					}
				}
			}
		}
	}
	if err != nil {
		log. // WithField("Error", err).
			Printf("user-community HandleRequest error %+v", err)
	}
	return err
}

func onMenuList(np models.NamespacePayload4) (err error) {
	request := np.SlackRequest.InteractionCallback
	action := request.ActionCallback.AttachmentActions[0]
	
	selected := action.SelectedOptions[0]
	if clientID == "" {
		err = errors.New("onMenuList: clientID == ''")
		return
	}
	conn := connGen.ForPlatformID(np.TeamID.ToPlatformID())
	switch selected.Value {
	// case objectives.CreateIDO:
	// 	// err = invokeWorkflow(np)
	// 	// if err != nil {// temporary fallback to the old mechanism TODO: remove
	// 	// Create an objective
	// 	mc := models.MessageCallback{Module: "objectives", Source: userID, Topic: "init", Action: "ask",
	// 		Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
	// 	// Posting engagement to user prompting to create an objective
	// 	coaches := IDOCoaches(userID, teamID)
	// 	initObjs := InitsAndObjectives(userID, teamID)
	// 	objectives.CreateObjectiveEng(engagementTable, mc, coaches, objectives.DevelopmentObjectiveDates(namespace, ""),
	// 		initObjs, "Would you like to add a development objective?",
	// 		"Development objectives", true, dns, common2.IDOCreateCheck, teamID)
	// 	DeleteOriginalEng(userID, channelID, message.MessageTs)
	// 	// }
	case objectives.CreateIDO, objectives.CreateIDONow:
		err = wfRoutes.EnterWorkflow(wfRoutes.IssuesWorkflow, np, conn, "") //onCreateIDONow(np)
		// if err != nil { // temporary fallback to the old mechanism TODO: remove
		// 	fmt.Printf("Handling %s event", objectives.CreateIDONow)
		// 	mc := models.MessageCallback{Module: "objectives", Source: userID, Topic: "init", Action: "ask",
		// 		Month: strconv.Itoa(int(month)), Year: strconv.Itoa(year)}
		// 	createObjectiveNow(message, userID, teamID, mc.ToCallbackID(), &mc)
		// }
	case strategy.CreateStrategyObjective, 
		strategy.CreateFinancialObjective, 
		strategy.CreateCustomerObjective:
		err = wfRoutes.EnterWorkflow(wfRoutes.IssuesWorkflow, np, conn, issues.CreateIssueByTypeEvent(issues.SObjective))
	case strategy.ViewStrategyObjectives:
		err = wfRoutes.EnterWorkflow(wfRoutes.IssuesWorkflow, np, conn, issues.ViewListOfIssuesByTypeEvent(issues.SObjective))
	// case strategy.ViewCapabilityCommunityObjectives:
	// 	err = wfRoutes.EnterWorkflow(wfRoutes.IssuesWorkflow, np, issues.ViewListOfIssuesByTypeEvent(issues.SObjective))
	// case strategy.ViewAdvocacyObjectives: - already handled
	// 	err = wfRoutes.EnterWorkflow(wfRoutes.IssuesWorkflow, np, issues.ViewMyListOfIssuesByTypeEvent(issues.SObjective))
	case user.ViewObjectives:
		// List all objectives
		// onViewIDOs(np)
		err = wfRoutes.EnterWorkflow(wfRoutes.IssuesWorkflow, np, conn, issues.ViewListOfIssuesByTypeEvent(issues.IDO)) // "view-idos")
		// userObjs := objectives.AllUserObjectives(userID, userObjectivesTable,
		// 	string(userObjective.UserIDTypeIndex), models.IndividualDevelopmentObjective, 0)
		// if len(userObjs) > 0 {
		// 	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
		// 		Message: "You can find the list of your Individual Development Objectives in the thread."})
		// 	ListObjectives(userID, channelID, userObjs,
		// 		func(mc models.MessageCallback, objective models.UserObjective) []ebm.Attachment {
		// 			strategyFlag := core.IfThenElse(objective.Type == models.StrategyDevelopmentObjective,
		// 				true, false).(bool)
		// 			return updateObjAttachment(mc, teamID, "", "", "", &objective,
		// 				false, false, strategyFlag)
		// 		}, TimeStamp(message))
		// } else {
		// 	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
		// 		Message: "You do not have any Individual Development Objectives yet."})
		// }
	case user.StaleObjectivesForMe:
		err = wfRoutes.EnterWorkflow(wfRoutes.IssuesWorkflow, np, conn, issues.ViewListOfStaleIssuesByTypeEvent(issues.SObjective))
		// userObjs := strategy.UserCapabilityObjectivesWithNoProgressInAMonth(userID, Today(),
		// 	userObjectivesTable, string(userObjective.UserIDTypeIndex), userObjectivesProgressTable, 0)
		// if len(userObjs) > 0 {
		// 	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
		// 		Message: "You can find the list of your Capability Objectives that haven't been updated in the last 30 days in the thread below"})
		// 	ListObjectives(userID, channelID, userObjs,
		// 		func(mc models.MessageCallback, objective models.UserObjective) []ebm.Attachment {
		// 			return updateObjAttachment(mc, teamID, "", "", "", &objective,
		// 				false, false, true)
		// 		}, TimeStamp(message))
		// } else {
		// 	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
		// 		Message: "All your objectives have been updated"})
		// }
		// DeleteOriginalEng(userID, channelID, message.MessageTs)
	case strategy.CreateInitiative:
		err = wfRoutes.EnterWorkflow(wfRoutes.IssuesWorkflow, np, conn, issues.CreateIssueByTypeEvent(issues.Initiative))
	case user.StaleInitiativesForMe:
		err = wfRoutes.EnterWorkflow(wfRoutes.IssuesWorkflow, np, conn, issues.ViewListOfStaleIssuesByTypeEvent(issues.Initiative))
		// userObjs := strategy.UserInitiativesWithNoProgressInAWeek(userID, Today(),
		// 	userObjectivesTable, string(userObjective.UserIDTypeIndex), userObjectivesProgressTable, 0)
		// if len(userObjs) > 0 {
		// 	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
		// 		Message: "You can find the list of your Initiatives that haven't been updated in the last 7 days in the thread below."})
		// 	ListObjectives(userID, channelID, userObjs,
		// 		func(mc models.MessageCallback, objective models.UserObjective) []ebm.Attachment {
		// 			return updateObjAttachment(mc, teamID, "", "", "", &objective,
		// 				false, false, true)
		// 		}, TimeStamp(message))
		// } else {
		// 	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
		// 		Message: "All of your initiatives have been updated"})
		// }
	case user.StaleIDOsForMe:
		err = wfRoutes.EnterWorkflow(wfRoutes.IssuesWorkflow, np, conn, issues.ViewListOfStaleIssuesByTypeEvent(issues.IDO))
		// fmt.Printf("UserIDOsWithNoProgressInAWeek(%s, %s, %s, %s, %s)",
		// 	userID, Today(),
		// 	userObjectivesTable, string(userObjective.UserIDTypeIndex), userObjectivesProgressTable)

		// userObjs := objectives.UserIDOsWithNoProgressInAWeek(userID, Today(),
		// 	userObjectivesTable, string(userObjective.UserIDTypeIndex), userObjectivesProgressTable)
		// if len(userObjs) > 0 {
		// 	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
		// 		Message: "You can find  a list of your stale Individual Development Objectives in the thread"})
		// 	ListObjectives(userID, channelID, userObjs,
		// 		func(mc models.MessageCallback, objective models.UserObjective) []ebm.Attachment {
		// 			return updateObjAttachment(mc, teamID, "", "", "", &objective,
		// 				false, false, true)
		// 		}, TimeStamp(message))
		// } else {
		// 	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
		// 		Message: "All of your Individual Development Objectives have been updated"})
		// }
	case
		strategy.ViewInitiativeCommunityInitiatives,
		strategy.ViewCapabilityCommunityInitiatives:
		err = wfRoutes.EnterWorkflow(wfRoutes.IssuesWorkflow, np, conn, issues.ViewListOfIssuesByTypeEvent(issues.Initiative))
	case strategy.ViewAdvocacyObjectives:
		err = wfRoutes.EnterWorkflow(wfRoutes.IssuesWorkflow, np, conn, issues.ViewMyListOfIssuesByTypeEvent(issues.SObjective))

		// userObjs := strategy.UserAdvocacyObjectives(userID, userObjectivesTable, userObjectivesTypeIndex, 0)
		// if len(userObjs) > 0 {
		// 	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
		// 		Message: "You can find the Capability Objectives assigned to you in the thread below"})
		// 	ListObjectives(userID, channelID, userObjs,
		// 		func(mc models.MessageCallback, objective models.UserObjective) []ebm.Attachment {
		// 			strategyFlag := core.IfThenElse(objective.Type == models.StrategyDevelopmentObjective,
		// 				true, false).(bool)
		// 			return updateObjAttachment(mc, teamID, "", "", "", &objective,
		// 				false, false, strategyFlag)
		// 		}, TimeStamp(message))
		// } else {
		// 	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
		// 		Message: "There are no Capability Objectives assigned to you"})
		// }
	case strategy.ViewAdvocacyInitiatives:
		err = wfRoutes.EnterWorkflow(wfRoutes.IssuesWorkflow, np, conn, issues.ViewListOfAdvocacyIssuesByTypeEvent(issues.Initiative))
		// fmt.Printf("UserAdvocacyInitiatives(%s, %s, %s, %d)",
		// 	userID,
		// 	userObjectivesTable,
		// 	userObjectivesTypeIndex, 0)
		// inits := strategy.UserAdvocacyInitiatives(userID, userObjectivesTable, userObjectivesTypeIndex, 0)
		// if len(inits) > 0 {
		// 	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
		// 		Message: "You can find the Initiatives assigned to you in the thread below"})
		// 	ListObjectives(userID, channelID, inits,
		// 		func(mc models.MessageCallback, objective models.UserObjective) []ebm.Attachment {
		// 			strategyFlag := core.IfThenElse(objective.Type == models.StrategyDevelopmentObjective,
		// 				true, false).(bool)
		// 			return updateObjAttachment(mc, teamID, "", "", "", &objective,
		// 				false, false, strategyFlag)
		// 		}, TimeStamp(message))
		// } else {
		// 	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
		// 		Message: "There are no Initiatives assigned to you"})
		// }
	case strategy.ViewCommunityAdvocateObjectives:
		err = wfRoutes.EnterWorkflow(wfRoutes.IssuesWorkflow, np, conn, issues.ViewListOfAdvocacyIssuesByTypeEvent(issues.SObjective))

		// // List objectives for which the user is an advocate for
		// stratObjs := objectives.AllUserObjectives(userID, userObjectivesTable, string(userObjective.UserIDTypeIndex),
		// 	models.StrategyDevelopmentObjective, 0)
		// if len(stratObjs) > 0 {
		// 	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
		// 		Message: "Below thread contains Strategy Objectives that you are advocate for"})
		// 	ListObjectives(userID, channelID, stratObjs,
		// 		func(mc models.MessageCallback, objective models.UserObjective) []ebm.Attachment {
		// 			strategyFlag := core.IfThenElse(objective.Type == models.StrategyDevelopmentObjective,
		// 				true, false).(bool)
		// 			return updateObjAttachment(mc, teamID, "", "", "", &objective,
		// 				false, false, strategyFlag)
		// 		}, TimeStamp(message))
		// } else {
		// 	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
		// 		Message: "There are no Strategy Objectives assigned to you"})
	}
	return
}

func onConfirm(request slack.InteractionCallback, teamID models.TeamID) {
	// userID := request.User.ID
	// channelID := request.Channel.ID
	action := request.ActionCallback.AttachmentActions[0]
	message := request
	// Parse callback Id to messageCallback
	mc, err := utils.ParseToCallback(message.CallbackID)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse to callback"))
	act := strings.TrimPrefix(action.Name, "confirm_")
	if mc.Topic == CommentsName {
		switch act {
		case string(models.Update):
			id := action.Value
			var comments string
			var status models.ObjectiveStatusColor
			objsByID := LatestProgressUpdateByObjectiveID(mc.Target)
			uObj := userObjectiveByID(mc.Target)
			if len(objsByID) > 0 {
				comments = objsByID[0].Comments
				status = objsByID[0].StatusColor
			}
			today := time.Now().Format(DateFormat)
			label := ObjectiveProgressText2(uObj, today)
			survey := utils.AttachmentSurvey(string(label),
				progressCommentSurveyElements(ui.PlainText(uObj.Name), uObj.CreatedDate))
			val := fillCommentsSurveyValues(survey, comments, status)
			api := getSlackClient(message)
			// Open a survey associated with the engagement
			err = dialogFromSurvey(api, message, val, id, false, mc.Target)
			core.ErrorHandler(err, namespace, fmt.Sprintf("Could not open dialog from %s survey", id+":"+message.CallbackID))
		}
	}

}

func onAsk(request slack.InteractionCallback, teamID models.TeamID) {
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
	userID := request.User.ID
	channelID := request.Channel.ID
	action := request.ActionCallback.AttachmentActions[0]
	message := request
	// Parse callback Id to messageCallback
	mc, err := utils.ParseToCallback(action.Value)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse to callback"))
	act := strings.TrimPrefix(action.Name, "ask_")

	if mc.Topic == CommentsName {
		switch act {
		case CommentsName, string(models.Now):
			// Ask for comments
			id := action.Value // callbackId
			api := getSlackClient(message)

			comments := ""
			var status models.ObjectiveStatusColor
			objsByID := LatestProgressUpdateByObjectiveID(mc.Target)
			uObj := userObjectiveByID(mc.Target)
			for _, each := range objsByID {
				if each.CreatedOn == TodayISOString() {
					comments = each.Comments
					status = each.StatusColor
				}
			}
			label := ObjectiveProgressText2(uObj, TodayISOString())
			survey := utils.AttachmentSurvey(string(label),
				progressCommentSurveyElements(ui.PlainText(uObj.Name), uObj.CreatedDate))
			val := fillCommentsSurveyValues(survey, comments, status)
			// Open a survey associated with the engagement
			// When the original message is from a thread, we need to post to the same thread
			// Below logic checks if the incoming message is from a thread
			var ts string
			if message.OriginalMessage.ThreadTimestamp == "" {
				ts = message.MessageTs
			} else {
				ts = message.OriginalMessage.ThreadTimestamp
			}

			survState := func() string {
				msgStateBytes, err := json.Marshal(MsgState{Ts: message.MessageTs, ThreadTs: ts, Update: false, ObjectiveId: mc.Target})
				core.ErrorHandler(err, namespace, "Could not marshal MsgState")
				return string(msgStateBytes)
			}
			err = utils.SlackSurvey(api, message, val, id, survState)
			core.ErrorHandler(err, namespace, fmt.Sprintf("Could not open dialog from %s survey", id+":"+message.CallbackID))
		case string(models.Ignore):
			utils.UpdateEngAsIgnored(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
			DeleteOriginalEng(message.User.ID, message.Channel.ID, message.MessageTs)
		}
	} else if mc.Topic == "init" {
		switch act {
		case string(models.Now):
			createObjectiveNow(conn, message, userID, action.Value, mc)
		case string(models.Update):
			// We should write this an engagement for update
			id := message.CallbackID // callbackId
			target := mc.Target
			mc.Set("Target", "")
			initsObjsValues := append(InitsAndObjectives(userID, teamID), platformValues(models.TeamID(teamID))...)
			uObj := userObjectiveByID(target)
			val := utils.AttachmentSurvey(string(ObjectiveModifyDialogTitle), 
				objectives.EditObjectiveSurveyElems2(&uObj, IDOCoaches(userID, teamID, conn),
				objectives.DevelopmentObjectiveDates(namespace, ""), initsObjsValues))

			// Is the AttachmentActionSurvey non-empty
			api := getSlackClient(message)

			// Open a survey associated with the engagement
			err = dialogFromSurvey(api, message, val, id, true, target)
			core.ErrorHandler(err, namespace, fmt.Sprintf("Could not open dialog from %s survey", id+":"+message.CallbackID))
		case string(ViewMore):
			uObj := userObjectiveByID(mc.Target)
			// Set target for message callback as the id for the engagement
			mc.Set("Target", mc.Target)
			strategyFlag := core.IfThenElse(uObj.ObjectiveType == models.StrategyDevelopmentObjective, true, false).(bool)
			attachs := updateObjAttachment(conn, *mc, "", "", "", &uObj, true, false, strategyFlag)
			publish(models.PlatformSimpleNotification{UserId: message.User.ID, Channel: message.Channel.ID, Attachments: attachs, Ts: message.MessageTs})
		case string(ViewLess):
			uObj := userObjectiveByID(mc.Target)
			// Set target for message callback as the id for the engagement
			mc.Set("Target", mc.Target)
			strategyFlag := core.IfThenElse(uObj.ObjectiveType == models.StrategyDevelopmentObjective, true, false).(bool)
			attachs := updateObjAttachment(conn, *mc, "", "", "", &uObj, false, false, strategyFlag)
			publish(models.PlatformSimpleNotification{UserId: message.User.ID, Channel: message.Channel.ID, Message: "", Attachments: attachs, Ts: message.MessageTs})
		case string(ViewProgress):
			uObj := userObjectiveByID(mc.Target)
			// Set target for message callback as the id for the engagement
			mc.Set("Target", mc.Target)
			strategyFlag := core.IfThenElse(uObj.ObjectiveType == models.StrategyDevelopmentObjective, true, false).(bool)
			attachs := updateObjAttachment(conn, *mc, "", "", "", &uObj, false, true, strategyFlag)
			publish(models.PlatformSimpleNotification{UserId: message.User.ID, Channel: message.Channel.ID, Message: "", Attachments: attachs, Ts: message.MessageTs})
		case string(models.Ok):
			DeleteOriginalEng(message.User.ID, message.Channel.ID, message.MessageTs)
		case string(models.Ignore):
			utils.UpdateEngAsIgnored(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
			DeleteOriginalEng(message.User.ID, message.Channel.ID, message.MessageTs)
		case string(Cancel):
			uObj := userObjectiveByID(mc.Target)
			SetObjectiveField(uObj, "cancelled", 1)
			// publish the message to the user
			publish(models.PlatformSimpleNotification{UserId: message.User.ID, Channel: message.Channel.ID,
				Message: core.TextWrap(fmt.Sprintf("Ok, cancelled the following objective: `%s`", uObj.Name),
					core.Underscore)})
			if uObj.Accepted == 1 {
				// post only if the objective has a coach
				publish(models.PlatformSimpleNotification{UserId: uObj.AccountabilityPartner,
					Message: core.TextWrap(fmt.Sprintf("<@%s> canceled the following objective: `%s`", uObj.UserID, uObj.Name),
						core.Underscore)})
			}
		case string(Enable):
			uObj := userObjectiveByID(mc.Target)
			SetObjectiveField(uObj, "cancelled", 0)
			publish(models.PlatformSimpleNotification{UserId: message.User.ID, Channel: message.Channel.ID,
				Message: core.TextWrap(fmt.Sprintf("Made the following objective active again: `%s`", uObj.Name),
					core.Underscore)})
			if uObj.Accepted == 1 {
				// post only if the objective has a coach
				publish(models.PlatformSimpleNotification{UserId: uObj.AccountabilityPartner,
					Message: core.TextWrap(fmt.Sprintf("<@%s> made the following objective active again: `%s`", uObj.UserID, uObj.Name),
						core.Underscore)})
			}
		case CommentsName:
			// name := utils.SlackFieldValue(message.OriginalMessage.Attachments[0], NameLabel)

			params := map[string]*dynamodb.AttributeValue{
				"id":      dynString(mc.Target),
				"user_id": dynString(message.User.ID),
			}
			var uObj models.UserObjective
			found, err2 := d.GetItemOrEmptyFromTable(userObjectivesTable, params, &uObj)
			if !found {
				uObj = models.UserObjective{}
			}
			core.ErrorHandler(err2, namespace, fmt.Sprintf("Could not query %s table for default values", userObjectivesTable))

			mc.Set("Topic", CommentsName)
			mc.Set("Target", uObj.ID)
			objectives.CommentsEng(engagementTable, *mc,
				"Nice! A little progress every day adds up to big results. Add some to the objective below?",
				"", &uObj, true, dns, common2.EngagementEmptyCheck)
			// Remove original message
			publish(models.PlatformSimpleNotification{UserId: message.User.ID, Channel: message.Channel.ID, Ts: message.MessageTs})
		case string(Closeout):
			logger.Infof("WARN: Old Closeout mechanism has been triggered")
			// uObj := userObjectiveByID(mc.Target)
			// typLabel := objectiveTypeLabel(uObj)
			// // If there is no partner assigned, send a message to the user that issue can't be closed-out until there is a coach
			// if uObj.AccountabilityPartner == "requested" || uObj.AccountabilityPartner == "none" {
			// 	publish(models.PlatformSimpleNotification{UserId: message.User.ID, Channel: message.Channel.ID, AsUser: true, Message: core.TextWrap(fmt.Sprintf(
			// 		"You do not have a coach for the objective: `%s`. Please get a coach before attemping to close out.", uObj.Name), core.Underscore, "*")})
			// } else {
			// 	// send a notification to the partner
			// 	objectives.ObjectiveCloseoutEng(engagementTable, *mc, uObj.AccountabilityPartner,
			// 		fmt.Sprintf("<@%s> wants to close the following %s. Are you ok with that?", message.User.ID, typLabel),
			// 		fmt.Sprintf("*%s*: %s \n *%s*: %s", NameLabel, uObj.Name, DescriptionLabel, uObj.Description),
			// 		string(closeoutLabel(uObj.ID)), objectiveCloseoutPath, false, dns, common2.EngagementEmptyCheck, teamID)
			// 	// Mark objective as closed
			// 	SetObjectiveField(uObj, "completed", 1)
			// 	// send a notification to the coachee that partner has been notified
			// 	publish(models.PlatformSimpleNotification{UserId: message.User.ID, Channel: message.Channel.ID,
			// 		Message: core.TextWrap(fmt.Sprintf("Awesome! I’ll schedule time with <@%s> to close out the %s: `%s`",
			// 			uObj.AccountabilityPartner, typLabel, uObj.Name), core.Underscore)})
			// }
		case string(HomeOptions):
			// Take to original options page
			// Go to ViewLess case
			uObj := userObjectiveByID(mc.Target)
			// Set target for message callback as the id for the engagement
			mc.Set("Target", mc.Target)
			strategyFlag := core.IfThenElse(uObj.ObjectiveType == models.StrategyDevelopmentObjective, true, false).(bool)
			attachs := updateObjAttachment(conn, *mc.WithTopic("init"), "", "", "", &uObj, false, false, strategyFlag)
			publish(models.PlatformSimpleNotification{UserId: message.User.ID, Channel: message.Channel.ID, Message: "", Attachments: attachs, Ts: message.MessageTs})
		case string(MoreOptions):
			// Show more options
			uObj := userObjectiveByID(mc.Target)

			attachs := moreOptionsAttachment(*mc, "", "", &uObj, teamID)
			publish(models.PlatformSimpleNotification{UserId: message.User.ID, Channel: message.Channel.ID, Message: "", Attachments: attachs, Ts: message.MessageTs})
		}
	} else if mc.Topic == "coaching" {
		// coaching related handlers
		switch act {
		case string(Confirm):
			onCoachConfirmAction(userID, channelID, message.MessageTs, *mc)
		case string(models.No), string(models.Update):
			id := mc.ToCallbackID()
			api := getSlackClient(message)
			// Open a survey associated with the engagement
			comments := utils.SlackFieldValue(message.OriginalMessage.Attachments[0], CommentsName)
			survey := CommentsSurvey(CoachingLabel, CoacheeRejectionReasonLabel, CommentsName, ui.PlainText(comments))
			err = dialogFromSurvey(api, message, survey, id, false, mc.Target)
			core.ErrorHandler(err, namespace, fmt.Sprintf("Could not open dialog from %s survey", id+":"+message.CallbackID))
		}
	}
}

func onCloseoutRequest(request slack.InteractionCallback, teamID models.TeamID) {
	// userID := request.User.ID
	// channelID := request.Channel.ID
	action := request.ActionCallback.AttachmentActions[0]
	message := request
	mc, err := utils.ParseToCallback(action.Value)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse to callback"))
	act := strings.TrimPrefix(action.Name, fmt.Sprintf("%s%s", Closeout, core.Underscore))
	switch act {
	case string(models.Now):
		// Mark the objective as closed for the coachee
		uObj := userObjectiveByID(mc.Target)
		// Optional comments for closeout
		id := message.CallbackID
		val := ebm.AttachmentActionSurvey{
			Title:       string(closeoutLabel(uObj.ID)),
			SubmitLabel: models.SubmitLabel,
			Elements: []ebm.AttachmentActionTextElement{
				ebm.NewTextArea(ObjectiveCloseoutComment, "Closeout Comments", ebm.EmptyPlaceholder, ""),
			},
		}
		api := getSlackClient(message)

		// Open a survey associated with the engagement
		err = dialogFromSurvey(api, message, val, id, true, uObj.ID)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not open dialog from %s survey", id+":"+message.CallbackID))
	case string(models.Update):
		// This is to update the closeout comments
		comments := utils.SlackFieldValue(message.OriginalMessage.Attachments[0], CommentsName)
		// Optional comments for closeout
		id := message.CallbackID
		val := ebm.AttachmentActionSurvey{
			Title:       string(closeoutLabel(mc.Target)),
			SubmitLabel: models.SubmitLabel,
			Elements: []ebm.AttachmentActionTextElement{
				ebm.NewTextArea(ObjectiveCloseoutComment, "Closeout Comments", ebm.EmptyPlaceholder, ui.PlainText(comments)),
			},
		}
		api := getSlackClient(message)

		// Open a survey associated with the engagement
		err = dialogFromSurvey(api, message, val, id, true, mc.Target)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not open dialog from %s survey", id+":"+message.CallbackID))
	case string(No):
		// Mark the objective as closed for the coachee
		uObj := userObjectiveByID(mc.Target)
		value := utils.SlackFieldValue(message.OriginalMessage.Attachments[0], CommentsName)
		// Optional comments for closeout
		id := message.CallbackID
		val := ebm.AttachmentActionSurvey{
			Title:       string(closeoutLabel(mc.Target)),
			SubmitLabel: models.SubmitLabel,
			Elements: []ebm.AttachmentActionTextElement{
				ebm.NewTextArea(ObjectiveNoCloseoutComment, objectives.ObjectiveCloseoutWhyDisagreeSurveyLabel, ebm.EmptyPlaceholder, ui.PlainText(value)),
			},
		}
		api := getSlackClient(message)

		// Open a survey associated with the engagement
		err = dialogFromSurvey(api, message, val, id, true, uObj.ID)
		core.ErrorHandler(err, namespace, fmt.Sprintf("Could not open dialog from %s survey", id+":"+message.CallbackID))
	}

}

func onViewStaleIDOs(request slack.InteractionCallback, teamID models.TeamID) {
	userID := request.User.ID
	channelID := request.Channel.ID
	action := request.ActionCallback.AttachmentActions[0]
	message := request
	mc, err := utils.ParseToCallback(action.Value)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse to callback"))
	act := strings.TrimPrefix(action.Name, fmt.Sprintf("%s%s", user.StaleIDOsForMe, core.Underscore))	
	switch act {
	case string(models.Now):
		// List the objectives with no progress
		// TODO: Check what this is
		// ListObjectivesWithNoProgress(userID, channelID, teamID, models.IndividualDevelopmentObjective)
		// Update engagement as answered and remove the original engagement
		utils.UpdateEngAsAnswered(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
		DeleteOriginalEng(userID, channelID, message.OriginalMessage.Timestamp)
	case string(models.Ignore):
		// Update engagement as ignored and remove the original engagement
		utils.UpdateEngAsIgnored(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
		DeleteOriginalEng(userID, channelID, message.OriginalMessage.Timestamp)
	}
}

func onViewOpenObjectives(request slack.InteractionCallback, teamID models.TeamID) {
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
	userID := request.User.ID
	channelID := request.Channel.ID
	action := request.ActionCallback.AttachmentActions[0]
	message := request
	mc, err := utils.ParseToCallback(action.Value)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse to callback"))
	act := strings.TrimPrefix(action.Name, fmt.Sprintf("%s%s", ViewOpenObjectives, core.Underscore))
	switch act {
	case string(models.Now):
		// List the objectives with no progress
		ListObjectivesWithEvaluation(userID, channelID, func(objective models.UserObjective) bool {
			return objective.Completed == 0
		}, func(mc models.MessageCallback, objective models.UserObjective) []ebm.Attachment {
			strategyFlag := core.IfThenElse(objective.ObjectiveType == models.StrategyDevelopmentObjective,
				true, false).(bool)
			return updateObjAttachment(conn, mc, "", "", "", &objective,
				false, false, strategyFlag)
		}, models.IndividualDevelopmentObjective, TimeStamp(message))
		// Update engagement as answered and remove the original engagement
		utils.UpdateEngAsAnswered(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
		DeleteOriginalEng(userID, channelID, message.OriginalMessage.Timestamp)
	case string(models.Ignore):
		// Update engagement as ignored and remove the original engagement
		utils.UpdateEngAsIgnored(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
		DeleteOriginalEng(userID, channelID, message.OriginalMessage.Timestamp)
	}
}

func objectiveProgressOnDate(objID string, date string) (uop models.UserObjectiveProgress, found bool, err error) {
	params := map[string]*dynamodb.AttributeValue{
		"id":         dynString(objID),
		"created_on": dynString(date),
	}

	found, err = d.GetItemOrEmptyFromTable(userObjectivesProgressTable, params, &uop)
	return
}

func onCapabilityObjectiveUpdateDueWithinWeek(request slack.InteractionCallback, teamID models.TeamID) {
	conn := connGen.ForPlatformID(teamID.ToPlatformID())

	userID := request.User.ID
	channelID := request.Channel.ID
	action := request.ActionCallback.AttachmentActions[0]
	message := request
	mc, err := utils.ParseToCallback(action.Value)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse to callback"))
	if strings.HasSuffix(action.Name, string(models.Now)) {
		act := strings.TrimSuffix(action.Name, fmt.Sprintf("_%s", string(models.Now)))
		var objs []models.UserObjective
		var text string
		var emptyListMessage string
		switch act {
		case user.IDOUpdateDueWithinWeek:
			objs = objectives.UserIDOsWithNoProgressInAWeek(userID, Today(), userObjectivesTable,
				string(userObjective.UserIDTypeIndex), userObjectivesProgressTable)
			text = "You can find the list of stale Individual Development Objectives in the thread"
			emptyListMessage = "All of your Individual Development Objectives has been updated"
		case user.CapabilityObjectiveUpdateDueWithinMonth:
			objs = strategy.UserCapabilityObjectivesWithNoProgressInAMonth(userID, Today(), userObjectivesTable,
				string(userObjective.UserIDTypeIndex), userObjectivesProgressTable, 0)
			text = "You can find the list of stale Capability Objectives in the thread"
			emptyListMessage = "All of your capability objectives has been updated"
		case user.InitiativeUpdateDueWithinWeek:
			objs = strategy.UserInitiativesWithNoProgressInAWeek(userID, Today(), userObjectivesTable,
				string(userObjective.UserIDTypeIndex), userObjectivesProgressTable, 0)
			text = "You can find the list of stale Initiatives in the thread"
			emptyListMessage = "All of your Initiatives has been updated"
		}
		if len(objs) > 0 {
			publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
				Message: text})
			ListObjectives(userID, channelID, objs,
				func(mc models.MessageCallback, objective models.UserObjective) []ebm.Attachment {
					return updateObjAttachment(conn, mc, "", "", "", &objective,
						false, false, true)
				}, TimeStamp(message))
		} else {
			publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Ts: message.MessageTs,
				Message: emptyListMessage})
		}
	} else if strings.HasSuffix(action.Name, string(models.Ignore)) {
		// Mark as ignored and delete the engagement
		utils.UpdateEngAsIgnored(userID, mc.ToCallbackID(), engagementTable, d, namespace)
		DeleteOriginalEng(userID, channelID, message.MessageTs)
	}
}

func onDialogSubmission(request slack.InteractionCallback, teamID models.TeamID) {
	dialog := request
	// Parse callback Id to messageCallback
	mc, err := utils.ParseToCallback(dialog.CallbackID)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not parse to callback"))
	fmt.Println("### callback in submission: " + dialog.CallbackID)
	// conn := connGen.ForPlatformID(teamID.ToPlatformID())

	if mc.Topic == "init" && mc.Action == "ask" {
		fmt.Println("ERROR: Invocation of the old onUserObjectiveSubmitted implementation which has been removed")
	} else if mc.Topic == "init" && mc.Action == string(Closeout) {
		fmt.Println("ERROR: Invocation of the old Closeout implementation which has been removed")
	} else if mc.Topic == "coaching" {
		switch mc.Action {
		case "ask":
			onCoachingRequestRejected(dialog, *mc)
		}
	}
}

// Get the alignment type for the aligned objective
func getAlignedStrategyTypeFromAlignmentID(strategyEntityID string) (alignment models.AlignedStrategyType, alignmentID string) {
	alignment = models.ObjectiveNoStrategyAlignment
	// strategy entity id is of the form 'initiative:<id>' or 'capability:<id>'
	splits := strings.Split(strategyEntityID, ":")
	if len(splits) == 2 {
		alignmentID = splits[1]
		switch splits[0] {
		case string(community.Capability):
			alignment = models.ObjectiveStrategyObjectiveAlignment
		case string(community.Initiative):
			alignment = models.ObjectiveStrategyInitiativeAlignment
		case string(community.Competency):
			alignment = models.ObjectiveCompetencyAlignment
		}
	}
	return
}

func convertDialogSubmissionToUserObjective(
	objectiveID string,
	userID string,
	teamID models.TeamID,
	form map[string]string,
) (userObj models.UserObjective) {
	objName := form[objectives.ObjectiveName]
	objDescription := form[objectives.ObjectiveDescription]
	partner := form[objectives.ObjectiveAccountabilityPartner]
	endDate := form[objectives.ObjectiveEndDate]
	strategyEntityID := form[objectives.ObjectiveStrategyAlignment]
	// Get the alignment type for the aligned objective
	alignment, alignmentID := getAlignedStrategyTypeFromAlignmentID(strategyEntityID)
	year, quarter := core.CurrentYearQuarter()

	userObj = models.UserObjective{
		UserID:                      userID,
		Name:                        objName,
		ID:                          objectiveID,
		Description:                 objDescription,
		AccountabilityPartner:       partner,
		ObjectiveType:               models.IndividualDevelopmentObjective,
		StrategyAlignmentEntityID:   alignmentID,
		StrategyAlignmentEntityType: alignment,
		Quarter:                     quarter,
		Year:                        year,
		CreatedDate:                 core.ISODateLayout.Format(time.Now()),
		ExpectedEndDate:             endDate,
		PlatformID:                  teamID.ToPlatformID(),
	}

	return
}

// This is to handle the case when a partner not accepts coaching
func onCoachingRequestRejected(request slack.InteractionCallback, mc models.MessageCallback) {
	dialog := request
	msgState := GetMsgStateUnsafe(request)
	userID := request.User.ID
	rejectionComments := dialog.Submission[CommentsName]

	// apr := models.AccountabilityPartnerShipRejection{ObjectiveID: mc.Target,
	// 	CreatedOn: core.CurrentRFCTimestamp(), UserID: mc.Source,
	// 	AccountabilityPartnerID: dialog.User.ID, Comments: rejectionComments}
	// err := d.PutTableEntry(apr, partnershipRejectionsTable)
	// if err == nil {
		notes, coachAttachs := coachingRejectionRequestNotifications(mc, userID, rejectionComments, msgState.ThreadTs)
		platformInstance.PublishAll(notes)

		utils.ECAnalysis(rejectionComments, IDOCoachingRejectionContext, "Non-acceptance comments",
			dialogTableName, mc.ToCallbackID(), dialog.User.ID, dialog.Channel.ID, msgState.Ts,
			msgState.ThreadTs, coachAttachs, s, platformNotificationTopic, namespace)
	// } else {
	// 	logger.WithField("error", err).Errorf("Could not write partner rejection entry in %s table", partnershipRejectionsTable)
	// }
}

func coachingRejectionRequestNotifications(mc models.MessageCallback, coachID, comments string, coachMsgThreadTs string) (
	notes []models.PlatformSimpleNotification, coachAttachs []ebm.Attachment) {
	// coach message
	coachAttachs = viewCommentsWithUpdateAttachment(mc, string(CoachingRequestRejectionReasonTitleToCoach), comments)
	notes = append(notes, models.PlatformSimpleNotification{UserId: coachID, Attachments: coachAttachs, ThreadTs: coachMsgThreadTs})

	// coachee message
	ido := userObjectiveByID(mc.Target)
	coacheeAttachs :=
		viewCommentsAttachment(mc, fmt.Sprintf("%s did not accept to coach you for the objective: %s",
			common.TaggedUser(coachID), ido.Name), comments)
	notes = append(notes, models.PlatformSimpleNotification{UserId: ido.UserID, Attachments: coacheeAttachs})
	return
}

func statusLabel(status models.ObjectiveStatusColor) ui.PlainText {
	return models.ObjectiveStatusColorLabels[status]
}

func splitObjectiveWithDateUnsafe(str string) (objID string, date string) {
	splits := strings.Split(str, "_")
	// this occurs for objectives, which has the form <objective_id>_<community_id>_<date>
	if len(splits) == 3 {
		objID = fmt.Sprintf("%s_%s", splits[0], splits[1])
		date = splits[2]
	} else if len(splits) == 2 {
		objID = splits[0]
		date = splits[1]
	} else {
		logger.Errorf("Expected 3 or 2 elements elements after split")
	}
	return
}

func viewCommentsProgressAttachment(mc models.MessageCallback, title, comments string, status models.ObjectiveStatusColor, notes string) []ebm.Attachment {
	attach := utils.ChatAttachment(title, "", "", mc.ToCallbackID(), updateCommentsAttachmentActions(mc),
		objCommentsProgressFields(comments, status, notes), time.Now().Unix())
	return []ebm.Attachment{*attach}
}

func commentsField(comments string) ebm.AttachmentField {
	return ebm.AttachmentField{
		Title: string(CommentsLabel),
		Value: comments,
		Short: true,
	}
}

func viewCommentsWithUpdateAttachment(mc models.MessageCallback, title, comments string) []ebm.Attachment {
	attach := utils.ChatAttachment(title, "", "", mc.ToCallbackID(), updateCommentsAttachmentActions(mc),
		[]ebm.AttachmentField{commentsField(comments)}, time.Now().Unix())
	return []ebm.Attachment{*attach}
}

func viewCommentsAttachment(mc models.MessageCallback, title, comments string) []ebm.Attachment {
	attach := utils.ChatAttachment(title, "", "", mc.ToCallbackID(), []ebm.AttachmentAction{},
		[]ebm.AttachmentField{commentsField(comments)}, time.Now().Unix())
	return []ebm.Attachment{*attach}
}

func updateCommentsAttachmentActions(mc models.MessageCallback) []ebm.AttachmentAction {
	return []ebm.AttachmentAction{
		*models.SimpleAttachAction(mc, models.Update, "I would like to change this"),
	}
}

func CreateUserObjective(userObj models.UserObjective, mc *models.MessageCallback, channelID string,
	teamID models.TeamID, threadTs string, update bool) {
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
	userID := userObj.UserID
	editStatus := core.IfThenElse(update, "updated", "created").(string)

	// It's an update, meaning we won't check for it's existence. Directly, update the entry
	err := d.PutTableEntry(userObj, userObjectivesTable)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not add entry to %s table", userObjectivesTable))
	strategyFlag := core.IfThenElse(userObj.ObjectiveType == models.StrategyDevelopmentObjective, true, false).(bool)
	attachs := updateObjAttachment(conn, *mc,
		ui.PlainText(ui.Sprintf("Here is the objective that you %s", editStatus)),
		"", "", &userObj, false, false, strategyFlag)
	publish(models.PlatformSimpleNotification{UserId: userID, Channel: channelID, Message: "", Attachments: attachs, Ts: threadTs})
	mc.Set("Target", "")
	utils.UpdateEngAsAnswered(mc.Source, mc.ToCallbackID(), engagementTable, d, namespace)
	// Do analysis on objective description
	utils.ECAnalysis(userObj.Description, IDODescriptionContext, "Development objective", dialogTableName, mc.ToCallbackID(),
		userID, channelID, threadTs, threadTs, attachs, s, platformNotificationTopic, namespace)

}

// COMMON FUNCTIONS //
var confirm = ebm.AttachmentActionConfirm{
	OkText:      models.YesLabel,
	DismissText: models.CancelLabel,
}

// This engagement is for a user (included in callback) asking 'target' for a partnership
// This engagement will be sent to the target since source requests it
func AskForPartnershipEngagement(teamID models.TeamID, mc models.MessageCallback, partner, title, text,
	fallback, learnLink string, urgent bool) {
	actions := models.AppendOptionalAction(
		[]ebm.AttachmentAction{
			*models.ConfirmAttachAction(mc, Confirm, "Yes", confirm),
			*models.SimpleAttachAction(mc, No, "No"), // Danger? models.EmptyActionSurvey(), models.EmptyActionConfirm(), true),
		},
		models.LearnMoreAction(models.ConcatPrefixOpt("docs/general/", learnLink)),
	)
	utils.AddChatEngagement(mc, title, text, fallback, partner, actions, []ebm.AttachmentField{}, teamID, urgent,
		engagementTable, d, namespace, time.Now().Unix(), common2.EngagementEmptyCheck)
}

func ObjectiveProgressAskEngagement(teamID models.TeamID, mc models.MessageCallback, userID, text string) {
	actions := []ebm.AttachmentAction{
		*models.SimpleAttachAction(mc, models.Now, "Yes"),
		*models.SimpleAttachAction(mc, models.Back, "I intended a different coachee"),
		*models.SimpleAttachAction(mc, models.Ignore, "Skip this"),
	}
	utils.AddChatEngagement(mc, "", text, "Adaptive at your service", userID, actions, []ebm.AttachmentField{},
		teamID, true, engagementTable, d, namespace, time.Now().Unix(), common2.EngagementEmptyCheck)
}

func CoachCoacheeProgressReviewAskEngagement(teamID models.TeamID, mc models.MessageCallback, userID, title,
	objName, objDesc, objComments string, percentDone string, urgent bool, emptyComments bool) {
	var actions []ebm.AttachmentAction
	if !emptyComments {
		actions = append(actions,
			*models.SimpleAttachAction(mc, models.Now, "Add my response"),
		)
	}
	actions = append(actions,
		*models.SimpleAttachAction(mc, models.Ignore, "Skip this"))
	utils.AddChatEngagement(mc, title, fmt.Sprintf("%s - %s", objName, objDesc), "Adaptive at your service",
		userID, actions, objCommentsPercentDoneProgressFields(objComments, percentDone, core.EmptyString), teamID,
		urgent, engagementTable, d, namespace, time.Now().Unix(), common2.EngagementEmptyCheck)
}

func objCommentsProgressFields(comments string, status models.ObjectiveStatusColor, notes string) []ebm.AttachmentField {
	return models.AttachmentFields([]models.KvPair{
		{Key: string(CommentsLabel), Value: comments},
		{Key: string(ProgressStatusLabel), Value: string(models.ObjectiveStatusColorLabels[status])},
		{Key: "Notes", Value: notes},
	})
}

func objCommentsPercentDoneProgressFields(comments string, percentDone string, notes string) []ebm.AttachmentField {
	return models.AttachmentFields([]models.KvPair{
		{Key: string(CommentsLabel), Value: comments},
		{Key: string(PercentDoneLabel), Value: percentDone},
		{Key: "Notes", Value: notes},
	})
}

func CommentsProgressSurvey(title, statusLabel ui.PlainText, statusName string,
	commentsLabel ui.PlainText, commentsName string, statusColor models.ObjectiveStatusColor) ebm.AttachmentActionSurvey {

	surveyElems := []ebm.AttachmentActionTextElement{
		{
			Label:    string(statusLabel),
			Name:     statusName,
			ElemType: models.MenuSelectType,
			Options:  utils.AttachActionElementOptions(models.ObjectiveStatusColorKeyValues),
			Value:    string(statusColor),
		},
		ebm.NewTextArea(commentsName, commentsLabel, ebm.EmptyPlaceholder, ""),
	}
	return utils.AttachmentSurvey(string(title), surveyElems)
}

// COMMON FUNCTIONS //

func surveyState(message slack.InteractionCallback, target string) string {
	// When the original message is from a thread, we need to post to the same thread
	// Below logic checks if the incoming message is from a thread
	var ts string
	if message.OriginalMessage.ThreadTimestamp == "" {
		ts = message.MessageTs
	} else {
		ts = message.OriginalMessage.ThreadTimestamp
	}
	fmt.Println("### dialog callback id: " + message.CallbackID)
	msgStateBytes, err := json.Marshal(MsgState{ThreadTs: ts, ObjectiveId: target})
	core.ErrorHandler(err, namespace, "Could not marshal MsgState")
	return string(msgStateBytes)
}

func percentTimeLapsed(today, start, end string) int {
	d1 := common.DurationDays(start, today, AdaptiveDateFormat,
		namespace)
	d2 := common.DurationDays(start, end, AdaptiveDateFormat, namespace)
	return int(float32(d1) / float32(d2) * float32(100))
}

func clipString(str string, prefixLength int) string {
	if len(str) < prefixLength {
		return str
	}
	return fmt.Sprintf("%s...", str[0:prefixLength-3])
}

// initiativesGroup formats one option group with initiatives
func initiativesGroup(userID string) (res []ebm.AttachmentActionElementOptionGroup, opInits []models.StrategyInitiative) {
	opInits = strategy.UserInitiativeCommunityInitiatives(userID, strategyInitiativesTableName,
		string(strategyInitiative.InitiativeCommunityIDIndex), communityUsersTable, string(adaptiveCommunityUser.UserIDIndex))

	if len(opInits) != 0 {
		grp := ebm.AttachmentActionElementOptionGroup{}
		options := grp.Options
		for _, each := range opInits {
			options = append(options,
				ebm.AttachmentActionElementOption{
					Label: clipString(each.Name, 30), // get first
					Value: fmt.Sprintf("%s:%s", community.Initiative, each.ID),
				})
		}
		sort.Sort(MenuOptionLabelSorter(options))
		grp.Options = options
		grp.Label = ui.PlainText("Initiatives")
		res = append(res, grp)
	}
	return
}

type MenuOptionLabelSorter []ebm.AttachmentActionElementOption

func (a MenuOptionLabelSorter) Len() int           { return len(a) }
func (a MenuOptionLabelSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MenuOptionLabelSorter) Less(i, j int) bool { return a[i].Label < a[j].Label }

// objectives formats one option group with objectives
func objectivesGroup(userID string, teamID models.TeamID, initiatives []models.StrategyInitiative) (res []ebm.AttachmentActionElementOptionGroup) {
	conn := platform.GetConnectionForUserFromEnvUnsafe(userID)

	capabilityObjectives := strategy.UserStrategyObjectives(userID, strategyObjectivesTableName,
		string(strategyObjective.PlatformIDIndex), userObjectivesTable,
		communityUsersTable, string(adaptiveCommunityUser.UserIDCommunityIDIndex), 
		conn,
	)

	var initiativeRelatedCapabilityObjectiveIDs []string
	// Getting the last of related Capability Objective for each of the Initiatives
	for _, each := range initiatives {
		initiativeRelatedCapabilityObjectiveIDs = append(initiativeRelatedCapabilityObjectiveIDs, each.CapabilityObjective)
	}

	var capabilityObjectiveIDs []string
	var groupOptions []ebm.AttachmentActionElementOption

	if len(capabilityObjectives) > 0 {
		// Add each Capability Objective from retrieved Capability Objectives in options
		for _, each := range capabilityObjectives {
			groupOptions = append(groupOptions,
				ebm.AttachmentActionElementOption{
					Label: clipString(each.Name, 30),
					Value: fmt.Sprintf("%s:%s", community.Capability, each.ID),
				})
			capabilityObjectiveIDs = append(capabilityObjectiveIDs, each.ID)
		}
	}
	objectivesIDsFromInitiativesNotInOptions := core.InBButNotA(capabilityObjectiveIDs, initiativeRelatedCapabilityObjectiveIDs)
	fmt.Println(capabilityObjectiveIDs)
	fmt.Println(initiativeRelatedCapabilityObjectiveIDs)
	fmt.Println(fmt.Sprintf("### objectivesIDsFromInitiativesNotInOptions: %v", objectivesIDsFromInitiativesNotInOptions))
	// Add related Capability Objectives from Initiatives, that are not already in options
	for _, each := range objectivesIDsFromInitiativesNotInOptions {
		obj := strategy.StrategyObjectiveByID(teamID, each, strategyObjectivesTableName)
		groupOptions = append(groupOptions,
			ebm.AttachmentActionElementOption{
				Label: clipString(obj.Name, 30),
				Value: fmt.Sprintf("%s:%s", community.Capability, obj.ID),
			})
	}
	// adding options to group only when they exist
	// reference error: Element 2 field `options` must have at least one option
	if len(groupOptions) > 0 {
		grp := ebm.AttachmentActionElementOptionGroup{
			Label:   "Objectives",
			Options: groupOptions,
		}
		res = append(res, grp)
	}
	return
}

// InitsAndObjectives returns initiatives and objectives
func InitsAndObjectives(userID string, teamID models.TeamID) (res []ebm.AttachmentActionElementOptionGroup) {
	i, inits := initiativesGroup(userID)
	res = append(res, i...)
	o := objectivesGroup(userID, teamID, inits)
	res = append(res, o...)
	return
}

func platformValues(teamID models.TeamID) (res []ebm.AttachmentActionElementOptionGroup) {
	vs := values.PlatformValues(teamID)
	if len(vs) != 0 {
		grp := ebm.AttachmentActionElementOptionGroup{}
		options := grp.Options
		for _, each := range vs {
			options = append(options,
				ebm.AttachmentActionElementOption{
					Label: each.Name,
					Value: fmt.Sprintf("%s:%s", community.Competency, each.ID),
				})
		}
		grp.Options = options
		grp.Label = "Competencies" // ui.PlainText(strings.Title(string(community.Strategy))) //
		res = append(res, grp)
	}
	return
}

func progressLabel(userObjID string) ui.PlainText {
	// suffix := "Progress"
	// userObj := userObjectiveByID(userObjID)
	// var prefix = objectiveTypeLabel(userObj)
	return ui.PlainText("Responsibility Progress")
}

func closeoutLabel(userObjID string) ui.PlainText {
	// suffix := "Closeout"
	// userObj := userObjectiveByID(userObjID)
	// var prefix = objectiveTypeLabel(userObj)
	// return ui.PlainText(strings.Join([]string{prefix, suffix}, " "))
	return ui.PlainText("Responsibility Closeout")
}

func closeoutAgreementContext(userObj models.UserObjective) (context string) {
	typeLabel := objectiveTypeLabel(userObj)
	switch typeLabel {
	case Individual:
		context = IDOCloseoutAgreementContext
	case CapabilityObjective:
		context = CapabilityObjectiveCloseoutAgreementContext
	case StrategyInitiative:
		context = InitiativeCloseoutAgreementContext
	}
	return
}

func closeoutDisagreementContext(userObj models.UserObjective) (context string) {
	issueType := utilsIssues.DetectIssueType(userObj)
	context = issueType.FoldString(
		IDOCloseoutDisagreementContext,
		CapabilityObjectiveCloseoutDisagreementContext,
		InitiativeCloseoutDisagreementContext,
	)
	return
}

func objectiveTypeLabel(userObj models.UserObjective) (prefix string) {
	return string(utilsIssues.ObjectiveTypeLabel(userObj))
}

func progressUpdateContext(userObj models.UserObjective) (context string) {
	return utilsIssues.DetectIssueType(userObj).
		FoldString(
			IDOProgressUpdateContext,
			CapabilityObjectiveProgressUpdateContext,
			InitiativeProgressUpdateContext,
		)
}

func responseUpdateContext(userObj models.UserObjective) (context string) {
	return utilsIssues.DetectIssueType(userObj).
		FoldString(
			IDOResponseObjectiveUpdateContext,
			CapabilityObjectiveUpdateResponseContext,
			InitiativeUpdateResponseContext,
		)
}

// onCoachConfirmAction handles the action when a coach is attempting to confirm a coaching request
func onCoachConfirmAction(coachID, channelID, ts string, mc models.MessageCallback) {
	coacheeObjective := userObjectiveByID(mc.Target)
	// From the time when the coaching request engagement has been posted to requested to coach,
	// before the acceptance, the user could have changed the coach to a different person
	// Here, we are checking if the initial requested coach is still the current requested coach
	if coacheeObjective.AccountabilityPartner == coachID {
		SetObjectiveField(coacheeObjective, "accepted", 1)
		// Send notification to coachee
		publish(models.PlatformSimpleNotification{UserId: coacheeObjective.UserID,
			Message: fmt.Sprintf("Your requested coach, <@%s>, has agreed to coach you for your development objective: `%s`.",
				coacheeObjective.AccountabilityPartner, coacheeObjective.Name)})
		// Send notification to partner
		publish(models.PlatformSimpleNotification{UserId: coacheeObjective.AccountabilityPartner,
			Message: fmt.Sprintf("Awesome! You will be coaching <@%s> for the development objective: `%s`.",
				mc.Source, coacheeObjective.Name)})
	} else {
		publish(models.PlatformSimpleNotification{UserId: coachID,
			Message: fmt.Sprintf("<@%s> has requested a different coach for the development objective: `%s`.",
				mc.Source, coacheeObjective.Name)})
	}
	// Update engagement as answered
	utils.UpdateEngAsAnswered(coachID, mc.ToCallbackID(), engagementTable, d, namespace)
	DeleteOriginalEng(coachID, channelID, ts)
}
