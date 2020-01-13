package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/holidays"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/values"
)

var bindings = map[string]string{
	"FetchEngagementsForMe":      user.AskForEngagements,
	"StaleIDOsExistForMe":        user.StaleIDOsForMe,
	"StaleObjectivesExistForMe":  user.StaleObjectivesForMe,
	"StaleInitiativesExistForMe": user.StaleInitiativesForMe,
	"AllIDOsForMe":               user.ViewObjectives,
	"AllIDOsForCoachees":         "all_idos_for_coachees",
	"AllObjectivesForMe":         strategy.ViewAdvocacyObjectives,
	"AllInitiativesForMe":        strategy.ViewAdvocacyInitiatives,
	"ProvideFeedback":            coaching.GiveFeedback,
	"RequestFeedback":            coaching.RequestFeedback,
	"ViewEditObjectives":         strategy.ViewStrategyObjectives,
	"ViewCommunityObjectives":    strategy.ViewCapabilityCommunityObjectives,
	"ViewEditInitiatives":        strategy.ViewCapabilityCommunityInitiatives,
	"ViewCommunityInitiatives":   strategy.ViewInitiativeCommunityInitiatives,
	// values
	"ViewValues":     values.AdaptiveValuesListMenuItem,
	"ViewEditValues": values.AdaptiveValuesListMenuItem,
	//
	"ViewCollaborationReport":    user.FetchReport,
	"ViewHolidays":               holidays.HolidaysListMenuItem,
	"ViewEditHolidays":           holidays.HolidaysListMenuItem,
	"ViewScheduleNextQuarter":    user.NextQuarterSchedule,
	"ViewScheduleCurrentQuarter": user.CurrentQuarterSchedule,
	"ViewVision":                 strategy.ViewVision,
	"ViewEditVision":             strategy.ViewEditVision,
	"ViewCoachees":               coaching.ViewCoachees,
	"ViewAdvocates":              coaching.ViewAdvocates,
	"CreateIDO":                  objectives.CreateIDONow,
	"CreateVision":               strategy.CreateVision,
	"CreateFinancialObjectives":  strategy.CreateFinancialObjective,
	"CreateCustomerObjectives":   strategy.CreateCustomerObjective,
	"CreateCapabilityObjectives": strategy.CreateStrategyObjective,
	"CreateCapabilityCommunity":  strategy.CreateCapabilityCommunity,
	"CreateInitiatives":          strategy.CreateInitiative,
	"CreateInitiativeCommunity":  strategy.CreateInitiativeCommunity,
	"CreateValues":               values.AdaptiveValuesCreateNewMenuItem,
	"CreateHolidays":             holidays.HolidaysCreateNewMenuItem,
	"AssignCapabilityObjective":  strategy.AssociateStrategyObjectiveToCapabilityCommunity,
	"UserSettings":               user.UpdateSettings,
	"StrategyPerformanceReport":  "StrategyPerformanceReport",
	"IDOPerformanceReport":       "IDOPerformanceReport",
}
