package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/coaching"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/holidays"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/user"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/values"
)

var bindings = map[string]string{
	"FetchEngagementsForMe":      user.AskForEngagements,
	"StaleObjectivesExistForMe":  user.StaleObjectivesForMe,
	"StaleInitiativesExistForMe": user.StaleInitiativesForMe,
	"AllIDOsForCoachees":         "all_idos_for_coachees",
	"AllObjectivesForMe":         strategy.ViewAdvocacyObjectives,
	"ProvideFeedback":            coaching.GiveFeedback,
	"RequestFeedback":            coaching.RequestFeedback,
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
	"CreateVision":               strategy.CreateVision,
	"CreateCapabilityCommunity":  strategy.CreateCapabilityCommunity,
	"CreateInitiativeCommunity":  strategy.CreateInitiativeCommunity,
	"CreateValues":               values.AdaptiveValuesCreateNewMenuItem,
	"CreateHolidays":             holidays.HolidaysCreateNewMenuItem,
	"AssignCapabilityObjective":  strategy.AssociateStrategyObjectiveToCapabilityCommunity,
	"UserSettings":               user.UpdateSettings,
	"StrategyPerformanceReport":  "StrategyPerformanceReport",
	"IDOPerformanceReport":       "IDOPerformanceReport",
}
