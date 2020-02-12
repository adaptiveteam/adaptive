package issues

import (
	"log"
	"strconv"
	"time"

	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	"github.com/adaptiveteam/adaptive/workflows/exchange"
	"github.com/pkg/errors"

	common "github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	engIssues "github.com/adaptiveteam/adaptive/adaptive-engagements/issues"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	issuesUtils "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	models "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	platform "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	userObjective "github.com/adaptiveteam/adaptive/daos/userObjective"
	userObjectiveProgress "github.com/adaptiveteam/adaptive/daos/userObjectiveProgress"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	ui "github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

func (w workflowImpl) OnProgressCancel(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	issueID := ctx.Data[issueIDKey]
	log := w.AdaptiveLogger.WithField("issueID", issueID)
	log.Info("OnProgressCancel")
	var newAndOldIssues NewAndOldIssues
	newAndOldIssues, err = w.WorkflowContext.GetNewAndOldIssues(ctx)
	err = errors.Wrapf(err, "OnProgressCancel getNewAndOldIssues")
	if err != nil {
		return
	}
	log.Infof("Found issue id=%s, cancelled=%v", newAndOldIssues.NewIssue.UserObjective.ID, newAndOldIssues.NewIssue.UserObjective.Cancelled)
	err = issuesUtils.SetCancelled(issueID)(w.DynamoDBConnection)
	err = errors.Wrapf(err, "OnProgressCancel SetCancelled")
	if err != nil {
		return
	}
	if err == nil {
		out, err = w.standardView(ctx)
	}
	itype := newAndOldIssues.NewIssue.GetIssueType()
	tc := getTypeClass(itype)

	issueTypeNameText := tc.IssueTypeName()
	out = ctx.Reply(ui.Sprintf("Ok, cancelled the following %s: `%s`",
		issueTypeNameText,
		newAndOldIssues.NewIssue.UserObjective.Name))

	if newAndOldIssues.NewIssue.UserObjective.Accepted == 1 { // post only if the objective has a coach
		out.Responses = append(out.Responses,
			platform.Post(platform.ConversationID(newAndOldIssues.NewIssue.UserObjective.AccountabilityPartner),
				platform.MessageContent{
					Message: ui.Sprintf("<@%s> has cancelled the following %s: `%s`",
						newAndOldIssues.NewIssue.UserObjective.UserID,
						issueTypeNameText,
						newAndOldIssues.NewIssue.UserObjective.Name),
				},
			),
		)
	}
	return
}

func (w workflowImpl) OnProgressIntermediate(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	issueID := ctx.Data[issueIDKey]
	w.AdaptiveLogger.WithField("issueID", issueID).Info("OnProgressIntermediate")
	var newAndOldIssues NewAndOldIssues
	newAndOldIssues, err = w.WorkflowContext.GetNewAndOldIssues(ctx)
	if err != nil {
		return
	}
	err = w.prefetch(ctx, &newAndOldIssues.NewIssue)
	if err != nil {
		return
	}
	comments := ui.PlainText("")
	var status models.ObjectiveStatusColor
	objectiveProgress := newAndOldIssues.NewIssue.Progress
	if len(objectiveProgress) > 0 {
		comments = ui.PlainText(objectiveProgress[0].Comments)
		status = models.ObjectiveStatusColor(objectiveProgress[0].StatusColor)
	}

	today := TodayISOString()
	// item := userObjectiveByID(itemID)
	label := ObjectiveProgressText2(newAndOldIssues.NewIssue.UserObjective, today)

	survey := utils.AttachmentSurvey(string(label),
		progressCommentSurveyElements(ui.PlainText(newAndOldIssues.NewIssue.UserObjective.Name),
			newAndOldIssues.NewIssue.UserObjective.CreatedDate))
	surveyWithValues := fillCommentsSurveyValues(survey, comments, status)
	out = out.WithSurvey(surveyWithValues)
	out.KeepOriginal = true
	return
}
func (w workflowImpl) OnProgressCloseout(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	issueID := ctx.Data[issueIDKey]
	w.AdaptiveLogger.WithField("issueID", issueID).Info("OnProgressCloseout")
	var newAndOldIssues NewAndOldIssues
	newAndOldIssues, err = w.WorkflowContext.GetNewAndOldIssues(ctx)
	uo := newAndOldIssues.NewIssue.UserObjective
	itype := getIssueTypeFromContext(ctx)
	tc := getTypeClass(itype)
	typLabel := tc.IssueTypeName()
	// If there is no partner assigned, send a message to the user that issue can't be closed-out until there is a coach
	if utilsUser.IsSpecialOrEmptyUserID(uo.AccountabilityPartner) {
		out.Interaction.Messages = wf.InteractiveMessages(wf.InteractiveMessage{
			PassiveMessage: wf.PassiveMessage{
				Text: ui.Sprintf("You do not have a coach for the %s: `%s`. Please get a coach before attemping to close out.",
					typLabel, uo.Name),
			},
		})
	} else {
		issue := newAndOldIssues.NewIssue
		issue.Completed = 1
		issue.CompletedDate = core.ISODateLayout.Format(time.Now())
		err = issuesUtils.Save(issue)(w.DynamoDBConnection)
		if err == nil {
			// send a notification to the coachee that partner has been notified
			out = out.WithInteractiveMessage(wf.InteractiveMessage{
				PassiveMessage: wf.PassiveMessage{
					Text: ui.Sprintf("Awesome! I’ll schedule time with <@%s> to close out the %s: `%s`",
						uo.AccountabilityPartner, typLabel, uo.Name),
				},
			})
			out = out.WithPostponedEvent(
				exchange.RequestCloseoutForIssue(newAndOldIssues.NewIssue),
			)
		}
	}
	return
}
func (w workflowImpl) OnProgressFormSubmitted(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	issueID := ctx.Data[issueIDKey]
	w.AdaptiveLogger.WithField("issueID", issueID).Info("OnProgressFormSubmitted")
	var newAndOldIssues NewAndOldIssues
	ctx.SetFlag(isShowingProgressKey, true) // enable show progress. This will make sure that progress is prefetched
	newAndOldIssues, err = w.WorkflowContext.GetNewAndOldIssues(ctx)
	ctx.RuntimeData = runtimeData(newAndOldIssues)
	uo := newAndOldIssues.NewIssue.UserObjective
	var progress userObjectiveProgress.UserObjectiveProgress
	progress, err = extractObjectiveProgressFromContext(ctx, uo)
	if err != nil {
		return
	}
	isProgressAvailableForToday := false
	if len(newAndOldIssues.NewIssue.Progress) > 0 {
		isProgressAvailableForToday = newAndOldIssues.NewIssue.Progress[0].CreatedOn == progress.CreatedOn
	}
	err = UserObjectiveProgressSave(progress)(w.DynamoDBConnection)
	if err != nil {
		return
	}
	err = w.prefetch(ctx, &newAndOldIssues.NewIssue)
	if err != nil {
		return
	}
	ctx.RuntimeData = runtimeData(newAndOldIssues)
	// attachs := viewProgressAttachment(mc,
	// 	ui.PlainText(Sprintf("This is your reported progress for the below %s", typLabel)),
	// 	"",
	// 	comments,
	// 	statusColor, item, models.Update)
	// publish(models.PlatformSimpleNotification{UserId: dialog.User.ID, Channel: dialog.Channel.ID, Ts: msgState.ThreadTs, Attachments: attachs})
	itype := newAndOldIssues.NewIssue.GetIssueType()
	tc := getTypeClass(itype)
	if err == nil {
		ctx.SetFlag(isShowingProgressKey, true)

		eventDescription := ObjectiveProgressCreatedUpdatedStatusTemplate(isProgressAvailableForToday, ctx.Request.User.ID)
		out, err = w.onNewOrUpdatedItemAvailable(ctx, tc, newAndOldIssues, ProgressUpdateContext,
			eventDescription, true)
		// if err == nil {
		// 	out.ImmediateEvent = "ProgressFormShown"
		// }
	} else {
		w.AdaptiveLogger.WithError(err).Error("OnProgressFormSubmitted error")
		out.Interaction = wf.SimpleResponses(
			platform.Post(platform.ConversationID(ctx.Request.User.ID),
				platform.MessageContent{Message: ui.Sprintf("Couldn't save the entered %s progress", tc.IssueTypeName())},
			),
		)
		err = nil // we want to show error interaction
	}
	out.RuntimeData = runtimeData(newAndOldIssues)
	return
}

func extractObjectiveProgressFromContext(ctx wf.EventHandlingContext, item userObjective.UserObjective) (progress userObjectiveProgress.UserObjectiveProgress, err error) {
	defer recoverToErrorVar("extractObjectiveProgressFromContext", &err)
	form := ctx.Request.DialogSubmissionCallback.Submission

	comments := form[ObjectiveProgressComments]
	statusColor := form[ObjectiveStatusColor]
	today := TodayISOString()

	progress = userObjectiveProgress.UserObjectiveProgress{
		ID:                item.ID,
		CreatedOn:         today,
		UserID:            ctx.Request.User.ID,
		Comments:          comments,
		PlatformID:        item.PlatformID,
		PartnerID:         item.AccountabilityPartner,
		PercentTimeLapsed: strconv.Itoa(percentTimeLapsed(today, item.CreatedDate, item.ExpectedEndDate)),
		StatusColor:       models.ObjectiveStatusColor(statusColor)}
	return
}

// func ObjectiveProgressText(objective models.UserObjective, today string) ui.RichText {
// 	timeUsed := fmt.Sprintf("%d days elapsed since %s",
// 		common.DurationDays(objective.CreatedDate, today, AdaptiveDateFormat, namespace), objective.CreatedDate)
// 	fmt.Printf("Time used for %s objective: %s", objective.Name, timeUsed)
// 	if objective.ExpectedEndDate == common.StrategyIndefiniteDateValue {
// 		return ui.Sprintf("%s", objective.Name)
// 	} else {
// 		return ui.Sprintf("%s", objective.Name)
// 	}
// }

// // JoinRichText concatenates elements of `a` placing `sep` in between.
// func JoinRichText(a []ui.RichText, sep ui.RichText) ui.RichText {
// 	s := make([]string, len(a))
// 	for i := 0; i < len(a); i++ {
// 		s[i] = string(a[i])
// 	}
// 	return ui.RichText(strings.Join(s, string(sep)))
// }

// // IntToString converts int to string
// func IntToString(i int) string {
// 	return fmt.Sprintf("%d", i)
// }

func ObjectiveProgressText2(objective userObjective.UserObjective, today string) ui.PlainText {
	var labelText ui.PlainText
	if objective.ExpectedEndDate == common.StrategyIndefiniteDateValue {
		labelText = "Progress"
	} else {
		percentElapsed := percentTimeLapsed(today, objective.CreatedDate, objective.ExpectedEndDate)
		labelText = ui.PlainText(ui.Sprintf("Time used - %d %%", percentElapsed))
	}
	return labelText
}

func percentTimeLapsed(today, start, end string) (percent int) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("percentTimeLapsed: %+v", e)
			percent = 0
		}
	}()
	d1 := common.DurationDays(start, today, core.ISODateLayout,
		"percentTimeLapsed")
	if end == common.StrategyIndefiniteDateValue {
		percent = 0
	} else {
		d2 := common.DurationDays(start, end, core.ISODateLayout, "percentTimeLapsed")
		percent = int(float32(d1) / float32(d2) * float32(100))
	}
	return
}

func TodayISOString() string {
	return core.ISODateLayout.Format(time.Now())
}

// func Percentages() (progressPercentValues []models.KvPair) {
// 	// show progress from 0 to 100 in increments of 10
// 	for i := 0; i <= 10; i++ {
// 		progressPercentValues = append(progressPercentValues, models.KvPair{
// 			// %% is required to use `%`. https://github.com/golang/go/commit/29499858bfa616b19c5108510d3cc6c9fa937bcc
// 			Key:   string(ui.Sprintf("%d %%", i*10)),
// 			Value: strconv.Itoa(i * 10),
// 		})
// 	}
// 	return
// }

// func objectiveCloseoutConfirmationDialogText(typ string) ui.PlainText {
// 	return ui.PlainText(fmt.Sprintf("Congratulations! Good job closing out this %s. I’m going to ask your partner if they agree. If they do, I’ll close this out for you.", typ))
// }

// func objectiveCancellationConfirmationDialogText(typ string) ui.PlainText {
// 	return ui.PlainText(fmt.Sprintf("You are attempting to cancel the %s", typ))
// }

// func cancelledObjectiveActivateConfirmationDialogText(typ string) ui.PlainText {
// 	return ui.PlainText(fmt.Sprintf("You are attempting re-activate a cancelled %s", typ))
// }

func progressCommentSurveyElements(objName ui.PlainText, startDate string) []ebm.AttachmentActionTextElement {
	nameConstrained := engIssues.ObjectiveCommentsTitle(objName)
	elapsedDays := common.DurationDays(startDate, TodayISOString(), core.ISODateLayout, "progressCommentSurveyElements")
	return []ebm.AttachmentActionTextElement{
		{
			Label:    string(engIssues.ObjectiveStatusLabel(elapsedDays, startDate)),
			Name:     ObjectiveStatusColor,
			ElemType: models.MenuSelectType,
			Options:  utils.AttachActionElementOptions(models.ObjectiveStatusColorKeyValues),
		},
		ebm.NewTextArea(ObjectiveProgressComments, nameConstrained, ObjectiveProgressCommentsPlaceholder, ""),
	}
}

func fillCommentsSurveyValues(sur ebm.AttachmentActionSurvey, comments ui.PlainText, status models.ObjectiveStatusColor) ebm.AttachmentActionSurvey {
	return models.FillSurvey(sur, map[string]string{
		ObjectiveProgressComments: string(comments),
		ObjectiveStatusColor:      string(status),
	})
}

const (
	ObjectiveProgressComments                         = "objective_progress"
	ObjectiveProgressCommentsPlaceholder ui.PlainText = ebm.EmptyPlaceholder
)

const (
	SlackLabelLimit = 48
)

const (
	ObjectiveStatusColor       = "objective_status_color"
	ObjectiveCloseoutComment   = "objective_closeout_comment"
	ObjectiveNoCloseoutComment = "objective_no_closeout_comment"
	ReviewUserProgressSelect   = "review_user_progress_select"
	UberCoach                  = "uber_coach"
)
