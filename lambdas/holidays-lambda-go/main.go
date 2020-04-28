package lambda

import (
	"github.com/adaptiveteam/adaptive/daos/adHocHoliday"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	eholidays "github.com/adaptiveteam/adaptive/adaptive-engagements/holidays"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	plat "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	eb "github.com/adaptiveteam/adaptive/engagement-builder"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/nlopes/slack"
)

const (
	HolidaysNamespace          = "holidays"
	HolidaysListMenuItem       = HolidaysNamespace + ":list"
	HolidaysSimpleListMenuItem = HolidaysNamespace + ":simple_list"
	HolidaysCreateNewMenuItem  = HolidaysNamespace + ":create_new"
	HolidayIDQueryField        = "holiday_id"
	dateFormat                 = "2006-01-02" // Only these constants are valid (https://stackoverflow.com/questions/14106541/go-parsing-date-time-strings-which-are-not-standard-formats)
	isInteractiveDebugEnabled  = false
)

var (
	HolidaysListAction = models.NewActionPath(models.NewPath(HolidaysNamespace, "list"))
)

func newHolidayTemplate() models.AdHocHoliday {
	return models.AdHocHoliday{
		ID:               core.Uuid(),
		Name:             "",
		Description:      "",
		Date:             "",
		ScopeCommunities: "",
	}
}

const (
	SubmitNewAdHocHolidayAction     = "submit-new-ad-hoc-holiday"
	DeleteAdHocHolidayAction        = "delete-ad-hoc-holiday"
	EditAdHocHolidayAction          = "edit-ad-hoc-holiday"
	CreateAdHocHolidayAction        = "create-ad-hoc-holiday"
	ModifyAdHocHolidayAction        = "modify-ad-hoc-holiday"
	SubmitUpdatedAdHocHolidayAction = "submit-updated-ad-hoc-holiday"
	ModifyListAdHocHolidayAction    = "modify-list-ad-hoc-holiday"
)

var (
	MenuListRequestHandlers = utils.RequestHandlers{
		HolidaysListMenuItem:       detailedListMenuItemFunc, // uses the menu message to attach the thread
		HolidaysSimpleListMenuItem: runAlso(listMenuItemFunc, clearOriginalMessage),
		HolidaysCreateNewMenuItem:  runAlso(createNewHolidayMenuItemFunc, clearOriginalMessage),
	}
	RequestHandlers = utils.RequestHandlers{
		"menu_list":              MenuListRequestHandlers.DispatchByRule(utils.SelectedOptionRule),
		DeleteAdHocHolidayAction: deleteHolidayButtonFunc,
		ModifyAdHocHolidayAction: editHolidayButtonFunc,
		EditAdHocHolidayAction:   editHolidayButtonFunc,
		// HolidaysCreateNewMenuItem:    createNewHolidayMenuItemFuncFromMainMenu, // I suspect that this line isn't working
		ModifyListAdHocHolidayAction: detailedListMenuItemFunc,
		CreateAdHocHolidayAction:     createNewHolidayMenuItemFunc,
	}
	DialogSubmissionHandlers = utils.DialogSubmissionHandlers{
		SubmitNewAdHocHolidayAction:     createNewHolidayDialogSubmissionHandler,
		SubmitUpdatedAdHocHolidayAction: updateHolidayDialogSubmissionHandler,
	}
	LambdaRouting = utils.LambdaHandler{
		Namespace:                             namespace,
		DispatchSlackInteractionCallback:      platform.DispatchInteractionCallback(RequestHandlers),
		DispatchSlackDialogSubmissionCallback: DispatchDialogSubmissionByRule(platform, DialogSubmissionHandlers, callbackIDRule),
	}
)

func runAlso(r1, r2 utils.RequestHandler) utils.RequestHandler {
	return func(request slack.InteractionCallback, conn daosCommon.DynamoDBConnection) (resp []models.PlatformSimpleNotification, err error){
		resp, err = r1(request, conn)
		if err == nil {
			var notes2 []models.PlatformSimpleNotification
			notes2, err = r2(request, conn)
			resp = append(resp, notes2...)
		}
		return
	}
}

func clearOriginalMessage(request slack.InteractionCallback, conn daosCommon.DynamoDBConnection) (resp []models.PlatformSimpleNotification, err error){
	resp = responses(utils.InteractionCallbackOverrideRequestMessage(request, ""))
	return 
}

func main() {
	LambdaRouting.StartHandler()
}


// func globalConnectionGen() daosCommon.DynamoDBConnectionGen {
// 	return daosCommon.CreateConnectionGenFromEnv()
// }

func readUser(userID string) (models.User, error) {
	conn := daosCommon.CreateConnectionGenFromEnv()
	dao := user.DAOFromConnectionGen(conn)
	// dao := user
	return dao.Read(userID)
}

func retrieveUserToken(request slack.InteractionCallback) (string, error) {
	return plat.GetTokenForUser(d, clientID, request.User.ID)
}

func teamID(request slack.InteractionCallback) models.TeamID {
	user, err2 := readUser(request.User.ID)
	core.ErrorHandler(err2, "teamID", "readUser failed") 
	return models.ParseTeamID(user.PlatformID)
}

func responses(notifications ...models.PlatformSimpleNotification) []models.PlatformSimpleNotification {
	return notifications
}


func callbackID(request slack.InteractionCallback, action string) models.ActionPath {
	return models.NewActionPath(
		models.NewPath(HolidaysNamespace, action),
		models.CurrentQuarterHashPair(),
		models.P("user_id", request.User.ID),
	)
}

func callbackIDRule(request slack.InteractionCallback) string {
	actionPath := models.ParseActionPath(request.CallbackID)
	return actionPath.Tail().Path.Head()
}

// DispatchDialogSubmissionByRule dispatches request using provided routing table
func DispatchDialogSubmissionByRule(p utils.Platform, r utils.DialogSubmissionHandlers, rule utils.RequestRoutingRule) func(slack.InteractionCallback, slack.DialogSubmissionCallback, daosCommon.DynamoDBConnection) {
	return func(request slack.InteractionCallback, dialog slack.DialogSubmissionCallback, conn daosCommon.DynamoDBConnection) {
		defer p.RecoverGracefully(request)
		notes, err2 := r.DispatchByRule(rule)(request, dialog, conn)
		core.ErrorHandler(err2, "DispatchDialogSubmissionByRule", "DispatchByRule")
		p.PublishAll(notes)
	}
}

// listMenuItemFunc renders the list of holidays directly in chat.
func listMenuItemFunc(request slack.InteractionCallback, conn daosCommon.DynamoDBConnection) (resp []models.PlatformSimpleNotification, err error){
	platform.Debug(request, "Listing holidays...")
	holidays := platformDAO(teamID(request)).AllUnsafe()
	sort.Slice(holidays, func(i, j int) bool {
		return holidays[i].Date < holidays[j].Date
	})
	holidayItems := mapAdHocHolidayString(holidays, RenderAdHocHolidayItem)
	listOfItems := strings.Join(holidayItems, "\n")
	mc := callbackID(request, ModifyListAdHocHolidayAction)

	actions := []ebm.AttachmentAction{}
	if DoesUserHaveWriteAccessToHolidays(request.User.ID) {
		actions = []ebm.AttachmentAction{
			modifyListAdHocHolidayAction(request),
			addAnotherAdHocHolidayAction(request),
		}
	}
	attachment := eb.NewAttachmentBuilder().
		Title(HolidaysTemplate(HolidaysListTitleSubject)).
		Text(listOfItems).
		Actions(actions).
		CallbackId(mc.Encode()).
		ToAttachment()
	p := utils.InteractionCallbackSimpleResponse(request, "") // RenderAdHocHolidayItem(h)
	p.Attachments = []ebm.Attachment{attachment}

	resp = responses(p)
	return
}

var (
	deleteAdHocHolidayAction = func(request slack.InteractionCallback, h models.AdHocHoliday) ebm.AttachmentAction {
		return eb.NewButtonDanger(DeleteAdHocHolidayAction, h.ID,
			ui.PlainText(HolidaysTemplate(DeleteButtonSubject)),
			ui.PlainText(HolidaysTemplate(CancelDeleteSubject)),
		)
	}

	editAdHocHolidayAction = func(request slack.InteractionCallback, h models.AdHocHoliday) ebm.AttachmentAction {
		return eb.NewButton(EditAdHocHolidayAction, h.ID,
			ui.PlainText(HolidaysTemplate(EditButtonSubject)))
	}

	addAnotherAdHocHolidayAction = func(request slack.InteractionCallback) ebm.AttachmentAction {
		return eb.NewButton(CreateAdHocHolidayAction, "",
			ui.PlainText(HolidaysTemplate(CreateButtonSubject)))
	}

	modifyListAdHocHolidayAction = func(request slack.InteractionCallback) ebm.AttachmentAction {
		return eb.NewButton(ModifyListAdHocHolidayAction, "",
			ui.PlainText(HolidaysTemplate(ModifyListButtonSubject)))
	}
)

// adHocHolidayInlineView creates an inline representation of the holiday with edit/delete buttons
func adHocHolidayInlineView(request slack.InteractionCallback, h models.AdHocHoliday) models.PlatformSimpleNotification {
	mc := callbackID(request, ModifyAdHocHolidayAction)

	actions := []ebm.AttachmentAction{}
	if DoesUserHaveWriteAccessToHolidays(request.User.ID) {
		actions = []ebm.AttachmentAction{
			editAdHocHolidayAction(request, h),
			deleteAdHocHolidayAction(request, h),
			addAnotherAdHocHolidayAction(request),
		}
	}
	b := eb.NewAttachmentBuilder()
	attachment :=
		b.
			Fields(ShowAdHocHolidayItemAsFields(h)).
			Actions(actions).
			CallbackId(mc.Encode()).
			ToAttachment()
	p := utils.InteractionCallbackSimpleResponse(request, "")
	p.Attachments = []ebm.Attachment{attachment}
	return p
}

func adHocHolidayEditChatMessage(request slack.InteractionCallback) func(models.AdHocHoliday) models.PlatformSimpleNotification {
	return func(h models.AdHocHoliday) models.PlatformSimpleNotification {
		view := adHocHolidayInlineView(request, h)
		view.ThreadTs = timeStamp(request)
		return view
	}
}

// detailedListMenuItemFunc sends the list of holidays to thread
func detailedListMenuItemFunc(request slack.InteractionCallback, conn daosCommon.DynamoDBConnection) (notes []models.PlatformSimpleNotification, err error) {
	platform.Debug(request, "Listing holidays...")
	holidays := platformDAO(teamID(request)).AllUnsafe()
	sort.Slice(holidays, func(i, j int) bool {
		return holidays[i].Date < holidays[j].Date
	})
	holidayItems := mapAdHocHolidayPlatformSimpleNotification(holidays,
		adHocHolidayEditChatMessage(request))
	replacementMessage := utils.InteractionCallbackOverrideRequestMessage(request,
		HolidaysTemplate(InviteToListOfHolidaysSubject))
	ts := time.Now()
	for _, item := range holidayItems {
		ts = ts.Add(time.Millisecond)
		item.Ts = core.TimestampLayout.Format(ts) // This is actually ignored in lambda that communicates with Slack
	}

	notes = append(responses(
		replacementMessage,
		//	title,
	),
		holidayItems...)
	return
}

func validateDate(request slack.InteractionCallback, dateStr string) {
	_, err := time.Parse(models.AdHocHolidayDateFormat, dateStr)
	platform.ErrorHandler(request, "Date: Couldn't parse "+dateStr, err)
}

func convertFormToAdHocHoliday(request slack.InteractionCallback, form map[string]string) models.AdHocHoliday {
	name := form["Name"]
	description := form["Description"]
	dateStr := form["Date"]
	user, err2 := readUser(request.User.ID)
	core.ErrorHandler(err2, "teamID", "readUser failed") 

	locations := user.Timezone // We don't have locations in the form anymore // form"Locations"]
	validateDate(request, dateStr)
	holiday := models.AdHocHoliday{
		ID:               core.Uuid(),
		Name:             name,
		Description:      description,
		Date:             dateStr,
		ScopeCommunities: locations,
		PlatformID:       user.PlatformID,
	}
	return holiday
}

func createNewHolidayDialogSubmissionHandler(request slack.InteractionCallback, dialog slack.DialogSubmissionCallback, conn daosCommon.DynamoDBConnection) (resp []models.PlatformSimpleNotification, err error){
	resp, ok := ensureUserHasWriteAccessToHolidays(request, conn)
	if ok {
		platform.Debug(request, "Creating holiday "+dialog.Submission["Name"])
		holiday := convertFormToAdHocHoliday(request, dialog.Submission)
		holiday.PlatformID = conn.PlatformID
		adHocHoliday.CreateUnsafe(holiday)(conn)
		platform.Debug(request, "Created holiday "+holiday.Name)
		view := adHocHolidayInlineView(request, holiday)
		//view.Ts = dialog.State
		resp = responses(view)
	}
	return
}

func createNewHolidayMenuItemFunc(request slack.InteractionCallback, conn daosCommon.DynamoDBConnection) (resp []models.PlatformSimpleNotification, err error){
	resp, ok := ensureUserHasWriteAccessToHolidays(request, conn)
	if ok {
		showCreateDialog(request)
		resp = responses()
	}

	return
}

func showCreateDialog(request slack.InteractionCallback) {
	mc := callbackID(request, SubmitNewAdHocHolidayAction)
	callbackID := mc.Encode()
	Ts := core.TimestampLayout.Format(time.Now())
	h := newHolidayTemplate()
	survey := utils.AttachmentSurvey("Holidays", eholidays.EditAdHocHolidayForm(&h))

	survey2 := ebm.AttachmentActionSurvey2{
		AttachmentActionSurvey: survey,
		CallbackID: callbackID,
		State: Ts,
		TriggerID: request.TriggerID,
	}
	// dialog, err := utils.ConvertSurveyToSlackDialog(survey, request.TriggerID, callbackID, Ts) // timeStamp(request))
	// platform.ErrorHandler(request,
	// 	fmt.Sprintf("Could not convert dialog to survey"),
	// 	err,
	// )
	// Open a survey associated with the engagement
	err2 := slackAPI(teamID(request)).ShowDialog(survey2)
	platform.Debug(request, "createNewHolidayMenuItemFunc: OpenDialog - done")
	platform.ErrorHandler(request,
		fmt.Sprintf("Could not open dialog from %s survey", request.CallbackID),
		err2,
	)
}
func deleteHolidayButtonFunc(request slack.InteractionCallback, conn daosCommon.DynamoDBConnection) (resp []models.PlatformSimpleNotification, err error){
	resp, ok := ensureUserHasWriteAccessToHolidays(request, conn)
	if ok {
		action := request.ActionCallback.AttachmentActions[0]
		id := action.Value
		err := adHocHoliday.Deactivate(id)(conn)
		platform.ErrorHandler(request,
			fmt.Sprintf("Could not delete %s", id),
			err,
		)
		view := utils.InteractionCallbackOverrideRequestMessage(request, HolidaysTemplate(DeletedHolidayNoticeSubject))
		resp = responses(view)
	}
	return
}

func editHolidayButtonFunc(request slack.InteractionCallback, conn daosCommon.DynamoDBConnection) (resp []models.PlatformSimpleNotification, err error){
	resp, ok := ensureUserHasWriteAccessToHolidays(request, conn)
	if ok {
		action := request.ActionCallback.AttachmentActions[0]
		id := action.Value
		// ut := retrieveUserToken(request)
		holiday := adHocHoliday.ReadUnsafe(id)(conn)
		platform.Debug(request, "Found holiday: "+holiday.Name)
		mc := callbackID(request, SubmitUpdatedAdHocHolidayAction)
		mc.Values.Set(HolidayIDQueryField, id)
		callbackID := mc.Encode()

		survey := utils.AttachmentSurvey("Holidays", eholidays.EditAdHocHolidayForm(&holiday))

		survey2 := ebm.AttachmentActionSurvey2{
			TriggerID: request.TriggerID,
			AttachmentActionSurvey: survey,
			State: request.MessageTs,
			CallbackID: callbackID,
		}
		// dialog, err := utils.ConvertSurveyToSlackDialog(survey, request.TriggerID, callbackID, request.MessageTs)
		// platform.ErrorHandler(request,
		// 	fmt.Sprintf("Could not convert dialog to survey"),
		// 	err,
		// )
		conn := daosCommon.DynamoDBConnection{
			Dynamo: d,
			ClientID: clientID,
			PlatformID: teamID(request).ToPlatformID(),
		}
		// Open a survey associated with the engagement
		err2 := mapper.SlackAdapterForTeamID(conn).ShowDialog(survey2)
		platform.Debug(request, "createNewHolidayMenuItemFunc: OpenDialog - done")
		platform.ErrorHandler(request,
			fmt.Sprintf("Could not open dialog from %s survey", request.CallbackID),
			err2,
		)
		// utils.InteractionCallbackOverrideRequestMessage(request, HolidaysTemplate(EditingHolidayNoticeSubject))
		// When just pressing Edit button we don't want to update anything.
		resp = responses()
	}
	return
}

func updateHolidayDialogSubmissionHandler(request slack.InteractionCallback, dialog slack.DialogSubmissionCallback, conn daosCommon.DynamoDBConnection) (resp []models.PlatformSimpleNotification, err error){
	resp, ok := ensureUserHasWriteAccessToHolidays(request, conn)
	if ok {
		platform.Debug(request, "Updating holiday "+dialog.Submission["Name"])
		mc := models.ParseActionPath(request.CallbackID)
		holiday := convertFormToAdHocHoliday(request, dialog.Submission)
		holiday.ID = mc.Values.Get(HolidayIDQueryField)
		err2 := adHocHoliday.CreateOrUpdate(holiday)(conn)
		platform.ErrorHandler(request, "Updating holiday", err2)
		platform.Debug(request, "Updated holiday "+dialog.Submission["Name"])
		view := adHocHolidayInlineView(request, holiday)
		view.Ts = dialog.State
		resp = responses(view)
	}
	return
}

func ensureUserHasWriteAccessToHolidays(request slack.InteractionCallback, conn daosCommon.DynamoDBConnection) (resp []models.PlatformSimpleNotification, ok bool) {
	ok = DoesUserHaveWriteAccessToHolidays(request.User.ID)
	if ok {
		log.Printf("Write access allowed %s", request.User.ID)
	} else {
		log.Printf("Write access denied %s", request.User.ID)
		resp = responses(
			utils.InteractionCallbackOverrideRequestMessage(request, "It looks like you don't have write access at the moment"),
		)
	}
	return
}
