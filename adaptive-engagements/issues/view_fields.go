package issues

import (
	"log"
	// "fmt"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	// community "github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	// objectives "github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	userObjectiveProgress "github.com/adaptiveteam/adaptive/daos/userObjectiveProgress"
	strategy "github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	// wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	models "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	ui "github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

func attachmentFieldNewOld(label ui.PlainText, prop func(Issue) ui.PlainText, newAndOldIssues NewAndOldIssues) ebm.AttachmentField {
	return ebm.AttachmentField{
		Title: string(label),
		Value: strategy.NewAndOld(string(prop(newAndOldIssues.NewIssue)), string(prop(newAndOldIssues.OldIssue))),
	}
}

func attachmentField(label ui.PlainText, value ui.PlainText) ebm.AttachmentField {
	return ebm.AttachmentField{
		Title: string(label),
		Value: string(value),
	}
}

func objectiveToFieldsNameDescriptionType(newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField) {
	fields = []ebm.AttachmentField{
		attachmentFieldNewOld(NameLabel, getName, newAndOldIssues),
		attachmentFieldNewOld(DescriptionLabel, getDescription, newAndOldIssues),
		// {Title: string("Type"), Value: "Individual"},
	}
	return
}

func ObjectiveProgressLabelTemplate(itype IssueType) ui.PlainText {
	return itype.Template() + " Progress"
}

var (
	getName                  = func(issue Issue) ui.PlainText { return ui.PlainText(issue.UserObjective.Name) }
	getDescription           = func(issue Issue) ui.PlainText { return ui.PlainText(issue.UserObjective.Description) }
	getInitiativeCommunity   = func(issue Issue) ui.PlainText { return ui.PlainText(issue.AlignedInitiativeCommunity.Name) }
	getRelatedObjective      = func(issue Issue) ui.PlainText { return ui.PlainText(issue.AlignedCapabilityObjective.Name) }
	getDefinitionOfVictory   = func(issue Issue) ui.PlainText { return ui.PlainText(issue.DefinitionOfVictory) }
	getBudget                = func(issue Issue) ui.PlainText { return ui.PlainText(issue.Budget) }
	// getAccountabilityPartner = func(issue Issue) ui.PlainText { return ui.PlainText("<@" + issue.UserObjective.AccountabilityPartner + ">") }
	getAccountabilityPartner = func(issue Issue) (res ui.PlainText) { 
		if issue.UserObjective.AccountabilityPartner == "none" || issue.UserObjective.AccountabilityPartner == "" {
			res = "None"
		} else {
			res = ui.PlainText("<@" + issue.UserObjective.AccountabilityPartner + ">") 
		}
		return 
	}
	
)


func renderObjectiveViewDate(issue Issue) ui.PlainText {
	defer core.RecoverAsLogError("renderObjectiveViewDate")
	if issue.UserObjective.CreatedDate == "" || issue.UserObjective.ExpectedEndDate == "" {
		return ui.PlainText("[" + issue.UserObjective.CreatedDate + "," + issue.UserObjective.ExpectedEndDate + "]")
	}
	objectiveDate := common.ObjectiveDate{
		CreatedDate:     issue.UserObjective.CreatedDate,
		ExpectedEndDate: issue.UserObjective.ExpectedEndDate,
	}
	return ui.PlainText(objectiveDate.Render(core.ISODateLayout, core.USDateLayout, "IDO issues workflow"))
}

func getStatus(issue Issue) (status ui.PlainText) {
	if issue.Cancelled == 1 {
		status = StatusCancelled
	} else if issue.Completed == 0 {
		status = StatusPending
	} else if issue.Completed == 1 && issue.PartnerVerifiedCompletion {
		status = StatusCompletedAndPartnerVerifiedCompletion
	} else if issue.Completed == 1 && !issue.PartnerVerifiedCompletion {
		status = StatusCompletedAndNotPartnerVerifiedCompletion
	}
	return
}

func getLatestComments(progress []userObjectiveProgress.UserObjectiveProgress) (status ui.PlainText) {
	comments := getCommentsFromProgress(progress)
	return ui.PlainText(ui.Join(comments, "\n"))
}


func getCommentsFromProgress(objectiveProgress []userObjectiveProgress.UserObjectiveProgress) (comments []ui.RichText) {
	for _, each := range objectiveProgress {
		comments = append(comments, ui.Sprintf("%s (%s percent, [%s] status)", each.Comments, each.PercentTimeLapsed, models.ObjectiveStatusColorLabels[each.StatusColor]))
	}
	return
}

// func getAccountabilityPartner(issue Issue) ui.PlainText {
// 	return readUserDisplayName(issue.AccountabilityPartner)
// }

func getObjectiveProgressComment(op userObjectiveProgress.UserObjectiveProgress) ui.RichText {
	res := ui.Sprintf("[%s] %s (%s)", models.ObjectiveStatusColorLabels[op.StatusColor], op.Comments, op.CreatedOn)

	if op.PartnerComments != "" {
		partnerStatusLabel := ui.RichText("")
		partnerStatus := models.ObjectiveStatusColor(op.PartnerReportedProgress)
		if partnerStatus != op.StatusColor {
			partnerStatusLabel = ui.Sprintf("[`%s`] ", models.ObjectiveStatusColorLabels[partnerStatus])
		}
		res = res + ui.Sprintf("\nPartner: %s%s", partnerStatusLabel, op.PartnerComments)
	}
	return res
}

func userObjectiveProgressField(progress []userObjectiveProgress.UserObjectiveProgress) (field ebm.AttachmentField) {
	// ops, err2 := userObjectiveProgressByID(item.ID, -1)
	comments := mapObjectiveProgressToRichText(progress, getObjectiveProgressComment)
	progressBody := ui.ListItems(comments...)
	
	if progressBody == "" {
		progressBody = "No progress"
	}
	return ebm.AttachmentField{
		Title: string("Progress"),
		Value: string(progressBody),
		Short: true,
	}
}

func mapObjectiveProgressToRichText(ops []userObjectiveProgress.UserObjectiveProgress, f func(userObjectiveProgress.UserObjectiveProgress) ui.RichText) (texts []ui.RichText) {
	for _, each := range ops {
		texts = append(texts, f(each))
	}
	return
}

func formatDate(date string, ipLayout, opLayout core.AdaptiveDateLayout) (res string) {
	defer core.RecoverAsLogError("formatDate")
	log.Printf("formatDate(%s, ipLayout=%s, opLayout=%s)", date, ipLayout, opLayout)
	if date == "" {
		res = "\"\""
	} else {
		var err error
		res, err = common.FormatDateWithIndefiniteOption(date, ipLayout, opLayout, "issue workflow formatDate")
		if err != nil {
			log.Printf("Could not parse string '%s' to date: %+v", date, err)
			res = date
		}
		
	}
	return
}

func renderStrategyAssociations(prefix string, fieldValue string) ui.PlainText {
	return ui.PlainText(ui.Sprintf("*%s* \n%s %s \n", prefix, BlueDiamondEmoji, fieldValue))
}

