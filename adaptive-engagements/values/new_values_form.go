package values

import (
	// "github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"time"
)

// CreateAdaptiveValuesEng This is an engagement prompting a user to create 
// a new adaptive value
func CreateAdaptiveValuesEng(mc models.MessageCallback,
	fallback, learnTrailPath, namespace string, check models.UserEngagementCheckWithValue,
	platformID models.PlatformID) models.UserEngagement {
	title := "Create adaptive competency"
	actions := models.AppendOptionalAction(
		[]ebm.AttachmentAction{
			*models.DialogAttachAction(mc, models.Now, "Create"),
		},
		models.LearnMoreAction(learnTrailPath))
	fields := []ebm.AttachmentField{}
	urgent := true
	return utils.MakeUserEngagement(mc, title, core.EmptyString, fallback,
		mc.Source, actions, fields, urgent, namespace, time.Now().Unix(), check, platformID)
}

const (
	NameLabel              ui.PlainText = "Name"
	NamePlaceholder        ui.PlainText = "Skills, Contribution or some other competency"
	DescriptionLabel       ui.PlainText = "Question to ask during 360 review" //
	DescriptionPlaceholder ui.PlainText = "Please ask a question to reveal the competency estimate. " +
		"For instance, \"How satisfied were you with this person's " +
		"skills this quarter?\"."
	DateLabel       ui.PlainText = "Date of the Holiday (YYYY-MM-DD)"
	DatePlaceholder ui.PlainText = ebm.EmptyPlaceholder

	ValueTypeLabel       ui.PlainText = "Type"
	ValueTypePlaceholder ui.PlainText = models.ValueTypeEnumPerformance + " or " + models.ValueTypeEnumRelationship
)

func option(name string, label ui.PlainText) ebm.AttachmentActionElementPlainTextOption {
	return ebm.AttachmentActionElementPlainTextOption{
		Label: label,
		Value: name,
	}
}

func options(o ...ebm.AttachmentActionElementPlainTextOption) []ebm.AttachmentActionElementPlainTextOption {
	return o
}

// EditAdaptiveValueForm creates a form for modifying/constructing adaptive value
func EditAdaptiveValueForm(obj *models.AdaptiveValue) []ebm.AttachmentActionTextElement {
	if obj == nil {
		obj = &models.AdaptiveValue{}
	}
	return []ebm.AttachmentActionTextElement{
		ebm.NewTextBox("Name", NameLabel, NamePlaceholder, ui.PlainText(obj.Name)),
		ebm.NewSimpleOptionsSelect("ValueType", ValueTypeLabel, ValueTypePlaceholder, obj.ValueType, options(
			option(models.ValueTypeEnumPerformance, "Performance"),
			option(models.ValueTypeEnumRelationship, "Relationship"),
		)...),
		ebm.NewTextArea("Description", DescriptionLabel, DescriptionPlaceholder, ui.PlainText(obj.Description)),
	}

}
