package lambda

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"strconv"
	"strings"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	"time"
)

const (
	SlackLabelLimit = 48
)
const (
	ObjectiveProgressComments                         = "objective_progress"
	ObjectiveProgressCommentsPlaceholder ui.PlainText = ebm.EmptyPlaceholder

	CommentsName                  = "Comments"
	CommentsLabel    ui.PlainText = "Comments"
	PercentDoneLabel ui.PlainText = "Percent Done"

	CommentsSurveyPlaceholder ui.PlainText = ebm.EmptyPlaceholder
	CommentsPlaceholder       ui.PlainText = ebm.EmptyPlaceholder

	CoachingName                             = "coaching"
	CoachingLabel               ui.PlainText = "Coaching"
	CoachRejectionReasonLabel   ui.PlainText = "Reason for not accepting the coach"
	CoacheeRejectionReasonLabel ui.PlainText = "Reason for not coaching"

	CoachingRequestRejectionReasonTitleToCoach = ui.RichText("You provided the following information for not accepting the coachee")

	NameLabel                            = "Name"
	DescriptionLabel                     = "Description"
	TimelineLabel                        = "Timeline"
	ProgressCommentsLabel   ui.PlainText = "Comments on Progress"
	ProgressStatusLabel     ui.PlainText = "Current Status"
	ObjectiveProgressLabel  ui.PlainText = "Objective Progress"
	PerceptionOfStatusLabel ui.PlainText = "Your perception of status"
	PerceptionOfStatusName = "perception_of_status"
	// CoachNotNeededOption ui.PlainText = "Coach not needed"
	// RequestACoachOption  ui.PlainText = "Request a coach"

	ListObjectivesCaption                                        ui.RichText = "I got it, here are your personal improvement objectives"
	ListObjectivesEmptyListOfObjectivesNeededProgressUpdatesText ui.RichText = "Good job! You added this week's progress for all your objectives"

	PartnerSelectingUserEngagementFallbackText ui.RichText  = "Adaptive at your service"
	PickAUserMenuPrompt                        ui.PlainText = "Pick a user..."

	CoacheeProgressSelectionPrompt ui.RichText = "Hello! I can fetch your associated users' progress. Whose are you looking for?"

	StatusCancelled                                ui.PlainText = "Cancelled"
	StatusPending                                  ui.PlainText = "Pending"
	StatusCompletedAndPartnerVerifiedCompletion    ui.PlainText = "Completed by you and closeout approved by your partner"
	StatusCompletedAndNotPartnerVerifiedCompletion ui.PlainText = "Completed by you and pending closeout approval from your partner"

	AccountabilityPartnerLabel ui.PlainText = "Accountability Partner"
	StatusLabel                ui.PlainText = "Status"
	LastReportedProgressLabel  ui.PlainText = "Last reported progress"
)

func ProgressTitle(objective models.UserObjective) ui.RichText {
	return ui.Sprintf("Entire progress for `%s`", objective.Name)
}

func ProgressAbsentTitle(objective models.UserObjective) ui.RichText {
	return ui.Sprintf("No progress reported for `%s`", objective.Name)
}

var (
	StrategyAssociationFieldLabel ui.PlainText = "Strategic Alignment"
)

func limitPlainText(text ui.PlainText, maxLength int) ui.PlainText {
	if len(text) < maxLength {
		return text
	}
	return text[:maxLength-3] + "..."
}

func ObjectiveCommentsTitle(objName ui.PlainText) ui.PlainText {
	nameConstrained := limitPlainText(ui.PlainText("Comments on "+objName), SlackLabelLimit)
	return nameConstrained
}

func ObjectiveStatusLabel(elapsedDays int, startDate string) ui.PlainText {
	return ui.PlainText(ui.Sprintf("Status (%d days since %s)", elapsedDays, startDate))
}

func ObjectiveProgressText(objective models.UserObjective, today string) ui.RichText {
	timeUsed := fmt.Sprintf("%d days elapsed since %s",
		common.DurationDays(objective.CreatedDate, today, AdaptiveDateFormat, namespace), objective.CreatedDate)
	fmt.Printf("Time used for %s objective: %s", objective.Name, timeUsed)
	if objective.ExpectedEndDate == common.StrategyIndefiniteDateValue {
		return ui.Sprintf("%s", objective.Name)
	} else {
		return ui.Sprintf("%s", objective.Name)
	}
}

// JoinRichText concatenates elements of `a` placing `sep` in between.
func JoinRichText(a []ui.RichText, sep ui.RichText) ui.RichText {
	s := make([]string, len(a))
	for i := 0; i < len(a); i++ {
		s[i] = string(a[i])
	}
	return ui.RichText(strings.Join(s, string(sep)))
}

// IntToString converts int to string
func IntToString(i int) string {
	return fmt.Sprintf("%d", i)
}

func ObjectiveProgressText2(objective models.UserObjective, today string) ui.PlainText {
	var labelText ui.PlainText
	if objective.ExpectedEndDate == common.StrategyIndefiniteDateValue {
		labelText = "Progress"
	} else {
		percentElapsed := percentTimeLapsed(today, objective.CreatedDate, objective.ExpectedEndDate)
		labelText = ui.PlainText(ui.Sprintf("Time used - %d %%", percentElapsed))
	}
	return labelText
}

func TodayISOString() string {
	return core.ISODateLayout.Format(time.Now())
}

// Today returns the current date
func Today() business_time.Date {
	return business_time.Today(time.UTC)
}

func Percentages() (progressPercentValues []models.KvPair) {
	// show progress from 0 to 100 in increments of 10
	for i := 0; i <= 10; i++ {
		progressPercentValues = append(progressPercentValues, models.KvPair{
			// %% is required to use `%`. https://github.com/golang/go/commit/29499858bfa616b19c5108510d3cc6c9fa937bcc
			Key:   string(ui.Sprintf("%d %%", i*10)),
			Value: strconv.Itoa(i * 10),
		})
	}
	return
}

func objectiveCloseoutConfirmationDialogText(typ string) ui.PlainText {
	return ui.PlainText(fmt.Sprintf("Congratulations! Good job closing out this %s. I’m going to ask your partner if they agree. If they do, I’ll close this out for you.", typ))
}

func objectiveCancellationConfirmationDialogText(typ string) ui.PlainText {
	return ui.PlainText(fmt.Sprintf("You are attempting to cancel the %s", typ))
}

func cancelledObjectiveActivateConfirmationDialogText(typ string) ui.PlainText {
	return ui.PlainText(fmt.Sprintf("You are attempting re-activate a cancelled %s", typ))
}

var (
	NextPageOfOptionsActionLabel ui.PlainText = "Show details"
	PrevPageOfOptionsActionLabel ui.PlainText = "Show less"

	DefaultConfirmationDialogTitle ui.PlainText = "Are you sure?"

	CancelledObjectiveActivateActionLabel             ui.PlainText = "Make active"
	CancelledObjectiveActivateConfirmationDialogTitle ui.PlainText = DefaultConfirmationDialogTitle

	ObjectiveModifyActionLabel                   ui.PlainText = "I want to modify it"
	ObjectiveModifyDialogTitle                   ui.PlainText = "Individual Objectives"
	ObjectiveAddProgressInfoActionLabel          ui.PlainText = "Made some progress" // WARN: It's used in two places!
	ObjectiveCancelActionLabel                   ui.PlainText = "Cancel this"
	ObjectiveCancellationConfirmationDialogTitle ui.PlainText = DefaultConfirmationDialogTitle

	ObjectiveMoreOptionsActionLabel        ui.PlainText = "More options"
	ObjectiveLessOptionsActionLabel        ui.PlainText = "Original options"
	ObjectiveShowEntireProgressActionLabel ui.PlainText = "Show entire progress"

	// WARN: these closeout labels are used more than once
	ObjectiveCloseoutActionLabel             ui.PlainText = "Closeout"
	ObjectiveCloseoutConfirmationDialogTitle ui.PlainText = DefaultConfirmationDialogTitle

	ObjectiveDetailsActionLabel ui.PlainText = "Details"

	ObjectiveAddAnotherActionLabel ui.PlainText = "How about another?"
	ObjectiveAddAnotherDialogTitle ui.PlainText = "Individual Objectives"

	ObjectiveProgressChangeCommentsActionLabel ui.PlainText = "Change my comments"

	ObjectiveProgressChangeCommentsDialogTitle ui.PlainText = "Individual Objectives" // NB! this title might be irrelevant

	ObjectivePartnerSelectActionLabel        ui.PlainText = "Partner on this"
	ObjectivePartnerSelectionSkipActionLabel ui.PlainText = "Skip this"
)
