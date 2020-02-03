package lambda

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

const (
	SObjectiveName        = "s_objective_name"
	SObjectiveDescription = "s_objective_description"
	SObjectiveMeasures    = "s_objective_measures"
	SObjectiveTargets     = "s_objective_targets"
	SObjectiveType        = "s_objective_type"
	SObjectiveAdvocate    = "s_objective_advocate"
	SObjectiveEndDate     = "s_objective_end_Date"

	// labels
	SObjectiveNameLabel        ui.PlainText = "Name"
	SObjectiveDescriptionLabel ui.PlainText = "Description"
	SObjectiveMeasuresLabel    ui.PlainText = "Measures"
	SObjectiveTargetsLabel     ui.PlainText = "Targets"
	SObjectiveTypeLabel                     = "Type"
	SObjectiveAdvocateLabel                 = "Advocate"
	SObjectiveEndDateLabel                  = "Time to work on this"
)

func EditSObjectiveSurveyElems(obj *models.StrategyObjective, types, advocates, dates []models.KvPair) []ebm.AttachmentActionTextElement {
	var op []ebm.AttachmentActionTextElement
	// Survey box should consist of a menu option to select rating and a text area for a user to enter the feedback
	if obj == nil {
		op = []ebm.AttachmentActionTextElement{
			{
				Label:    SObjectiveTypeLabel,
				Name:     SObjectiveType,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(types),
				Value: ObjectiveTypeDefaultValue,
			},
			ebm.NewTextBox(SObjectiveName, SObjectiveNameLabel, ebm.EmptyPlaceholder, ""),
			ebm.NewTextArea(SObjectiveDescription, SObjectiveDescriptionLabel, ebm.EmptyPlaceholder, ""),
			ebm.NewTextArea(SObjectiveMeasures, SObjectiveMeasuresLabel, ebm.EmptyPlaceholder, ""),
			ebm.NewTextArea(SObjectiveTargets, SObjectiveTargetsLabel, ebm.EmptyPlaceholder, ""),
			{
				Label:    SObjectiveAdvocateLabel,
				Name:     SObjectiveAdvocate,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(advocates),
			},
			{
				Label:    SObjectiveEndDateLabel,
				Name:     SObjectiveEndDate,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(dates),
			},
		}
	} else {
		op = []ebm.AttachmentActionTextElement{
			{
				Label:    SObjectiveTypeLabel,
				Name:     SObjectiveType,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(types),
				Value:    string(obj.ObjectiveType),
			},
			ebm.NewTextBox(SObjectiveName, SObjectiveNameLabel, ebm.EmptyPlaceholder, ui.PlainText(obj.Name)),
			ebm.NewTextArea(SObjectiveDescription, SObjectiveDescriptionLabel, ebm.EmptyPlaceholder, ui.PlainText(obj.Description)),
			ebm.NewTextArea(SObjectiveMeasures, SObjectiveMeasuresLabel, ebm.EmptyPlaceholder, ui.PlainText(obj.AsMeasuredBy)),
			ebm.NewTextArea(SObjectiveTargets, SObjectiveTargetsLabel, ebm.EmptyPlaceholder, ui.PlainText(obj.Targets)),
			{
				Label:    SObjectiveAdvocateLabel,
				Name:     SObjectiveAdvocate,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(advocates),
				Value:    obj.Advocate,
			},
			{
				Label:    SObjectiveEndDateLabel,
				Name:     SObjectiveEndDate,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(dates),
				Value:    obj.ExpectedEndDate,
			},
		}
	}
	return op
}

const ObjectiveTypeDefaultValue = "No Type"
func ObjectiveTypes() []models.KvPair {
	return []models.KvPair{
		// {Key: "Customer Strategy Objective", Value: string(strategy.CustomerStrategyObjective)},
		// {Key: "Financial Strategy Objective", Value: string(strategy.FinancialStrategyObjective)},
		// {Key: "Capability Strategy Objective", Value: string(strategy.CapabilityStrategyObjective)},
		{Key: "No Type", Value: "No Type"},
		{Key: "Financial Performance", Value: "Financial Performance"},
		{Key: "Effective Resource Use", Value: "Effective Resource Use"},
		{Key: "Customer Value", Value: "Customer Value"},
		{Key: "Customer Satisfaction", Value: "Customer Satisfaction"},
		{Key: "Customer Retention", Value: "Customer Retention"},
		{Key: "Efficiency", Value: "Efficiency"},
		{Key: "Quality", Value: "Quality"},
		{Key: "People Development", Value: "People Development"},
		{Key: "Infrastructure", Value: "Infrastructure"},
		{Key: "Technology", Value: "Technology"},
		{Key: "Culture", Value: "Culture"},
	}
}