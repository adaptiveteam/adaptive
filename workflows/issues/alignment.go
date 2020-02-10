package issues

import (
	"log"
	"fmt"
	"strings"

	community "github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	userObjective "github.com/adaptiveteam/adaptive/daos/userObjective"
	ui "github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

// Get the alignment type for the aligned objective
func getAlignedStrategyTypeFromStrategyEntityID(strategyEntityID string) (alignment userObjective.AlignedStrategyType, alignmentID string) {
	alignment = userObjective.ObjectiveNoStrategyAlignment
	// strategy entity id is of the form 'initiative:<id>' or 'capability:<id>'
	splits := strings.Split(strategyEntityID, ":")
	if len(splits) == 2 {
		alignmentID = splits[1]
		switch splits[0] {
		case string(community.Capability):
			alignment = userObjective.ObjectiveStrategyObjectiveAlignment
		case string(community.Initiative):
			alignment = userObjective.ObjectiveStrategyInitiativeAlignment
		case string(community.Competency):
			alignment = userObjective.ObjectiveCompetencyAlignment
		}
	}
	return
}

func (IDOImpl) GetAlignment(issue Issue) (alignment string) {
	switch issue.StrategyAlignmentEntityType {
	case userObjective.ObjectiveStrategyObjectiveAlignment:
		alignment = renderStrategyAssociations("Capability Objective", "Name", issue.AlignedCapabilityObjective)
	case userObjective.ObjectiveStrategyInitiativeAlignment:
		alignment = renderStrategyAssociations("Initiative", "Name", issue.AlignedCapabilityInitiative)
	case userObjective.ObjectiveCompetencyAlignment:
		alignment = fmt.Sprintf("Competency: `%s`", issue.AlignedCompetency.Name)
		log.Printf("[alignment.go:40] issue (uo.id=%s).AlignedCompetency (id=%s, name=%s)",
			issue.UserObjective.ID,
			issue.AlignedCompetency.ID,
			issue.AlignedCompetency.Name)
	}
	return
}

func (SObjectiveImpl) GetAlignment(issue Issue) (alignment string) {
	splits := strings.Split(issue.UserObjective.ID, "_")
	if len(splits) == 2 {
		alignment = fmt.Sprintf("%s%s",
			renderStrategyAssociations("Capability Communities", "Name", issue.PrefetchedData.AlignedCapabilityCommunity),
			renderStrategyAssociations("Capability Objectives", "Name", issue.StrategyObjective))
	} else {
		alignment = fmt.Sprintf("`%s Objective` : `%s`\n", issue.StrategyObjective.ObjectiveType, issue.StrategyObjective.Name)
	}
	return
}
func (InitiativeImpl) GetAlignment(issue Issue) (alignment string) {
	alignment = fmt.Sprintf("%s - %s",
		renderStrategyAssociations("Initiative Communities", "Name", issue.AlignedInitiativeCommunity),
		renderStrategyAssociations("Capability Objectives", "Name", issue.AlignedCapabilityObjective),
	)
	return
}

// func objectiveType(platformID models.PlatformID) func(uObj models.UserObjective) (typ string, alignment string) {
// 	return func(uObj models.UserObjective) (typ string, alignment string) {
// 		typ = "Not aligned with strategy"
// 		if uObj.Type == models.IndividualDevelopmentObjective {
// 			typ = Individual
// 			switch uObj.StrategyAlignmentEntityType {
// 			case models.ObjectiveStrategyObjectiveAlignment:
// 				capObj := strategy.StrategyObjectiveByID(platformID, uObj.StrategyAlignmentEntityID, strategyObjectivesTableName)
// 				alignment = renderStrategyAssociations("Capability Objective", "Name", capObj)
// 			case models.ObjectiveStrategyInitiativeAlignment:
// 				initiative := strategy.StrategyInitiativeByID(platformID, uObj.StrategyAlignmentEntityID, strategyInitiativesTableName)
// 				alignment = renderStrategyAssociations("Initiative", "Name", initiative)
// 			case models.ObjectiveCompetencyAlignment:
// 				valueID := uObj.StrategyAlignmentEntityID
// 				dns := common.DeprecatedGetGlobalDns()
// 				valueDao := values.DAOImpl{DNS: &dns}
// 				valueDao.Name = valuesTableName
// 				valueDao.PlatformIDIndex = valuesPlatformIdIndex
// 				value, err2 := valueDao.Read(valueID)
// 				core.ErrorHandler(err, namespace, fmt.Sprintf("Could not read value from %s table", valuesTableName))
// 				alignment = fmt.Sprintf("Value: `%s`", value.Name)
// 			}
// 		} else if uObj.Type == models.StrategyDevelopmentObjective {
// 			switch uObj.StrategyAlignmentEntityType {
// 			case models.ObjectiveStrategyObjectiveAlignment:
// 				typ = CapabilityObjective
// 				splits := strings.Split(uObj.ID, "_")
// 				if len(splits) == 2 {
// 					so := strategy.StrategyObjectiveByID(platformID, splits[0], strategyObjectivesTableName)
// 					capComm := strategy.CapabilityCommunityByID(platformID, splits[1], capabilityCommunitiesTableName)
// 					alignment = fmt.Sprintf("%s%s",
// 						renderStrategyAssociations("Capability Communities", "Name", capComm),
// 						renderStrategyAssociations("Capability Objectives", "Name", so))
// 				} else {
// 					so := strategy.StrategyObjectiveByID(platformID, uObj.ID, strategyObjectivesTableName)
// 					alignment = fmt.Sprintf("`%s Objective` : `%s`\n", so.Type, so.Name)
// 				}
// 			case models.ObjectiveStrategyInitiativeAlignment:
// 				typ = StrategyInitiative
// 				si := strategy.StrategyInitiativeByID(platformID, uObj.ID, strategyInitiativesTableName)
// 				initCommID := si.InitiativeCommunityID
// 				capObjID := si.CapabilityObjective
// 				initComm := strategy.InitiativeCommunityByID(platformID, initCommID, strategyInitiativeCommunitiesTable)
// 				capObj := strategy.StrategyObjectiveByID(platformID, capObjID, strategyObjectivesTableName)
// 				alignment = fmt.Sprintf("%s%s",
// 					renderStrategyAssociations("Initiative Communities", "Name", initComm),
// 					renderStrategyAssociations("Capability Objectives", "Name", capObj))
// 			case models.ObjectiveNoStrategyAlignment:
// 				alignment = "Not aligned with any strategy"
// 			}
// 		}
// 		return
// 	}
// }
func renderStrategyAssociations(prefix, field string, entities ...interface{}) string {
	var op string
	if len(entities) > 0 {
		var acc ui.RichText
		for _, entity := range entities {
			acc += ui.Sprintf("%s %s \n", BlueDiamondEmoji, GetFieldString(entity, field))
		}
		op = fmt.Sprintf("*%s* \n%s", prefix, acc)
	}
	return op
}
