package issues

import (
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ui "github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

const (
	SlackLabelLimit = 48
)

const (
	InitiativeCommunityLabel ui.PlainText = "Initiative Community"
	RelatedObjectiveLabel    ui.PlainText = "Related Objective"
	DefinitionOfVictoryLabel ui.PlainText = "Definition of Victory"
	BudgetLabel              ui.PlainText = "Budget"

	InitiativeNameLabel        ui.PlainText = "Initiative Name"
	InitiativeDescriptionLabel ui.PlainText = "Initiative Description"
	InitiativeVictoryLabel     ui.PlainText = "Definition of Victory"

	InitiativeAdvocateLabel            = "Advocate"
	InitiativeBudgetLabel              = "Budget($) in the following format: 1234.56"
	InitiateEndDateLabel               = "Time to work on this"
	InitiativeCapabilityObjectiveLabel = "Related Capability Objective"

	InitiativeName                    = "initiative_name"
	InitiativeDescriptionName         = "initiative_description"
	InitiativeVictoryName             = "definition_of_victory"
	InitiativeAdvocateName            = "advocate"
	InitiativeBudgetName              = "initiative_budget_name"
	InitiateEndDateName               = "time_to_work_on_this"
	InitiativeCapabilityObjectiveName = "initiative_capability_objective"
)


const (
	NameLabel                            = "Name"
	DescriptionLabel                     = "Description"
	TimelineLabel                        = "Timeline"
	ProgressCommentsLabel   ui.PlainText = "Comments on Progress"
	ProgressStatusLabel     ui.PlainText = "Current Status"
	PerceptionOfStatusLabel ui.PlainText = "Your perception of status"
	PerceptionOfStatusName               = "perception_of_status"

	StrategyAssociationFieldLabel ui.PlainText = "Strategic Alignment"
)

const (
	AccountabilityPartnerLabel ui.PlainText = "Accountability Partner"
	StatusLabel                ui.PlainText = "Status"
	LastReportedProgressLabel  ui.PlainText = "Last reported progress"
)

const (
	StatusCancelled                                ui.PlainText = "Cancelled"
	StatusPending                                  ui.PlainText = "Pending"
	StatusCompletedAndPartnerVerifiedCompletion    ui.PlainText = "Completed by you and closeout approved by your partner"
	StatusCompletedAndNotPartnerVerifiedCompletion ui.PlainText = "Completed by you and pending closeout approval from your partner"
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

const ObjectiveTypeDefaultValue = "No Type"

const (
	BlueDiamondEmoji                         = ":small_blue_diamond:"
)

func ObjectiveCommentsTitle(objName ui.PlainText) ui.PlainText {
	return ui.PlainText(core.ClipString("Comments on "+string(objName), SlackLabelLimit, "…"))
}

func ObjectiveStatusLabel(elapsedDays int, startDate string) ui.PlainText {
	return ui.PlainText(ui.Sprintf("Status (%d days since %s)", elapsedDays, startDate))
}
