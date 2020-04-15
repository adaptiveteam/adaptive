package values

import (
	"sort"
	"github.com/adaptiveteam/adaptive/daos/adaptiveValue"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	// "github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	// utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
)

// ReadAndSortAllAdaptiveValues returns active competencies
func ReadAndSortAllAdaptiveValues(conn daosCommon.DynamoDBConnection) (adaptiveValues []adaptiveValue.AdaptiveValue) {
	adaptiveValues = adaptiveValue.ReadByPlatformIDUnsafe(conn.PlatformID)(conn)
	adaptiveValues = adaptiveValue.AdaptiveValueFilterActive(adaptiveValues)
	sort.Slice(adaptiveValues, func(i, j int) bool {
		return adaptiveValues[i].Name < adaptiveValues[j].Name
	})
	return
}

// PlatformValues returns values for platform id
// Deprecated. Use ReadAndSortAllAdaptiveValues
func PlatformValues(teamID models.TeamID) []adaptiveValue.AdaptiveValue {
	connGen := daosCommon.CreateConnectionGenFromEnv()
	conn := connGen.ForPlatformID(teamID.ToPlatformID())
	return ReadAndSortAllAdaptiveValues(conn)
}
