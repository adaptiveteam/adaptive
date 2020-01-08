package test_checks

import (
	acfn "github.com/adaptiveteam/adaptive/adaptive-check-function-names"
	"github.com/adaptiveteam/adaptive/checks"
)

var AllTrueTestProfile = checks.CheckFunctionMap {
	acfn.FeedbackGivenThisQuarter                  :checks.ReturnsFalse,
	acfn.InLastMonthOfQuarter                      :checks.ReturnsFalse,

	// Due date tests
	acfn.IDOsDueWithinTheWeek                      :checks.ReturnsFalse,
	acfn.IDOsDueWithinTheMonth                     :checks.ReturnsFalse,
	acfn.IDOsDueWithinTheQuarter                   :checks.ReturnsFalse,

	acfn.InitiativesDueWithinTheWeek               :checks.ReturnsFalse,
	acfn.InitiativesDueWithinTheMonth              :checks.ReturnsFalse,
	acfn.InitiativesDueWithinTheQuarter            :checks.ReturnsFalse,

	acfn.ObjectivesDueWithinTheWeek                 :checks.ReturnsFalse,
	acfn.ObjectivesDueWithinTheMonth                :checks.ReturnsFalse,
	acfn.ObjectivesDueWithinTheQuarter              :checks.ReturnsFalse,

	// Community membership tests
	acfn.InCapabilityCommunity                     :checks.ReturnsFalse,
	acfn.InValuesCommunity                         :checks.ReturnsFalse,
	acfn.InHRCommunity                             :checks.ReturnsFalse,
	acfn.InStrategyCommunity                       :checks.ReturnsFalse,
	acfn.InInitiativeCommunity                     :checks.ReturnsFalse,

	// Component existence tests

	// Miscellaneous
	acfn.UserSettingsExist                         :checks.ReturnsFalse,
	acfn.HolidaysExist                             :checks.ReturnsFalse,
	acfn.CollaborationReportExists                 :checks.ReturnsFalse,
	acfn.UndeliveredEngagementsExistForMe          :checks.ReturnsFalse,

	// Strategy component existence tests independent of the user
	acfn.TeamValuesExist                           :checks.ReturnsFalse,
	acfn.CompanyVisionExists                       :checks.ReturnsFalse,
	acfn.ObjectivesExist                           :checks.ReturnsFalse,
	acfn.InitiativesExist                          :checks.ReturnsFalse,

	// Strategy component existence tests for a given user
	acfn.IDOsExistForMe                            :checks.ReturnsTrue,
	acfn.ObjectivesExistForMe                      :checks.ReturnsTrue,
	acfn.InitiativesExistForMe                     :checks.ReturnsTrue,

	// Stale components exist for a specfc individual
	acfn.StaleIDOsExistForMe                       :checks.ReturnsTrue,
	acfn.StaleInitiativesExistForMe                :checks.ReturnsTrue,
	acfn.StaleObjectivesExistForMe                 :checks.ReturnsTrue,

	// Community existence tests
	acfn.CapabilityCommunityExists                 :checks.ReturnsFalse,
	acfn.MultipleCapabilityCommunitiesExists       :checks.ReturnsFalse,
	acfn.InitiativeCommunityExists                 :checks.ReturnsFalse,
	acfn.MultipleInitiativeCommunitiesExists       :checks.ReturnsFalse,

	// State of community tests
	acfn.ObjectivesExistInMyCapabilityCommunities  :checks.ReturnsFalse,
	acfn.InitiativesExistInMyCapabilityCommunities :checks.ReturnsFalse,
	acfn.InitiativesExistInMyInitiativeCommunities :checks.ReturnsFalse,
}