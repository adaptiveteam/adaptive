package communities


import (
	// "github.com/pkg/errors"
	"strings"
	"github.com/adaptiveteam/adaptive/daos/community"
	"github.com/adaptiveteam/adaptive/daos/capabilityCommunity"
	"github.com/adaptiveteam/adaptive/daos/strategyInitiativeCommunity"
	"github.com/adaptiveteam/adaptive/daos/strategyCommunity"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunity"
	"github.com/adaptiveteam/adaptive/daos/migration"
	"github.com/adaptiveteam/adaptive/daos/common"
	"fmt"
)

// func findAdaptiveCommunityByID(acs []adaptiveCommunity.AdaptiveCommunity, id string) (found []adaptiveCommunity.AdaptiveCommunity) {
// 	for _, ac := range acs {
// 		idParts := strings.Split(ac.ID, ":")
// 		if ac.ID == id || (len(idParts) > 1 && idParts[1] == id) {
// 			found = append(found, ac)
// 		}
// 	}
// 	return
// }

func groupAdaptiveCommunityByID(acs []adaptiveCommunity.AdaptiveCommunity) (m map[string]adaptiveCommunity.AdaptiveCommunity) {
	m = make(map[string]adaptiveCommunity.AdaptiveCommunity, len(acs))
	for _, ac := range acs {
		var id string
		idParts := strings.Split(ac.ID, ":")
		if len(idParts) > 1 {
			id = idParts[1]
		} else {
			id = ac.ID
		}
		m[id] = ac
	}
	return
}

func indexStrategyCommunityByID(scs []strategyCommunity.StrategyCommunity) (m map[string]strategyCommunity.StrategyCommunity) {
	m = make(map[string]strategyCommunity.StrategyCommunity, len(scs))
	for _, sc := range scs {
		m[sc.ID] = sc
	}
	return
}

// MigrateStrategyCommunities - copies communities to Community table.
func MigrateStrategyCommunities(conn common.DynamoDBConnection, m *migration.Migration) (err error) {
	var scs []strategyCommunity.StrategyCommunity
	err = conn.Dynamo.ScanTable(strategyCommunity.TableName(conn.ClientID), &scs)
	if err != nil {
		return
	}

	total := 0
	if err == nil {
		for _, sc := range scs {
			if sc.PlatformID == conn.PlatformID || conn.PlatformID == "-all" {
				total ++
				c := convertStrategyCommunity(sc)
				 
				err = community.CreateOrUpdate(c)(conn)
				if err == nil {
					m.SuccessCount ++
				} else {
					fmt.Printf("FAILED: %+v\n", err)
					m.FailuresCount ++
				}
			}
			if err != nil {
				break
			}
		}
	}
	fmt.Printf("StrategyCommunity: success %d + failures %d of total %d.\n", m.SuccessCount, m.FailuresCount, total)
	return
}

func convertStrategyCommunity(sc strategyCommunity.StrategyCommunity) (c community.Community) {
	c = community.Community{
		ID: sc.ID,
		Advocate: sc.Advocate,
		AccountabilityPartner: sc.AccountabilityPartner,
		ChannelID: sc.ChannelID,
		CreatedAt: sc.CreatedAt,
		ModifiedAt: sc.ModifiedAt,
		CreatedBy: sc.AccountabilityPartner,
		ModifiedBy: sc.AccountabilityPartner,
		PlatformID: sc.PlatformID,
		Name: "",
		Description: "",
		CommunityKind: common.StrategyCommunity,
		DeactivatedAt: "", // active
		ParentCommunityID: string(sc.ParentCommunity),
	}
	fmt.Printf("community %s, channel %s: ", sc.ID, sc.ChannelID)
	switch sc.Community { // capability/initiative
	case common.Capability:
		c.CommunityKind = common.ObjectiveCommunity
	case common.Initiative:
		c.CommunityKind = common.InitiativeCommunity
	default:
		c.CommunityKind = common.ObjectiveCommunity
		fmt.Printf("Invalid community kind: %s", sc.Community)
	}
	if sc.ChannelID == "none" || sc.ChannelCreated == 0 {
		c.ChannelID = ""
		fmt.Printf("cleared channel id because it's not created")
	}
	fmt.Println()
	return
}

// MigrateSimpleCommunities - copies communities to Community table.
func MigrateSimpleCommunities(conn common.DynamoDBConnection, m *migration.Migration) (err error) {
	var acs []adaptiveCommunity.AdaptiveCommunity
	err = conn.Dynamo.ScanTable(adaptiveCommunity.TableName(conn.ClientID), &acs)
	if err != nil {
		return
	}


	total := 0
	if err == nil {
		for _, ac := range acs {
			if ac.PlatformID == conn.PlatformID || conn.PlatformID == "-all" {
				cs := convertSimpleCommunity(ac)
				 
				for _, c := range cs {
					total ++
					err = community.CreateOrUpdate(c)(conn)
					if err == nil {
						m.SuccessCount ++
					} else {
						fmt.Printf("FAILED: %+v\n", err)
						m.FailuresCount ++
					}
				}
			}
			if err != nil {
				break
			}
		}
	}
	fmt.Printf("AdaptiveCommunity (simple): success %d + failures %d of total %d.\n", m.SuccessCount, m.FailuresCount, total)
	return
}

func convertSimpleCommunity(ac adaptiveCommunity.AdaptiveCommunity) (c []community.Community) {
	var id common.CommunityKind 
	id = ""
	switch ac.ID {
	case "admin": id = common.AdminCommunity
	case "hr": id = common.HRCommunity
	case "coaching": id = common.CoachingCommunity
	case "competency": id = common.CompetencyCommunity
	case "strategy": id = common.StrategyCommunity
	case "user": id = common.UserCommunity
	// case "values": id = common.CompetencyCommunity
	}
	if id != "" {
		c = append(c, community.Community{
			ID: string(id),
			ChannelID: ac.ChannelID,
			CreatedAt: ac.CreatedAt,
			ModifiedAt: ac.ModifiedAt,
			PlatformID: ac.PlatformID,
			Name: strings.Title(ac.ID),
			Description: "",
			CommunityKind: id,
			DeactivatedAt: ac.DeactivatedAt,
			ParentCommunityID: "",
		})
		fmt.Printf("community %s, channel %s\n", c[0].ID, c[0].ChannelID)
	} else {
		fmt.Printf("Skipped id=%s", ac.ID)
	}
	return
}
// MigrateCapabilityCommunities - copies communities to Community table.
func MigrateCapabilityCommunities(conn common.DynamoDBConnection, m *migration.Migration) (err error) {
	var scs []strategyCommunity.StrategyCommunity
	err = conn.Dynamo.ScanTable(strategyCommunity.TableName(conn.ClientID), &scs)
	if err != nil {
		return
	}
	scsMap := indexStrategyCommunityByID(scs)

	var ccs []capabilityCommunity.CapabilityCommunity
	err = conn.Dynamo.ScanTable(capabilityCommunity.TableName(conn.ClientID), &ccs)
	if err != nil {
		return
	}

	total := 0
	if err == nil {
		for _, cc := range ccs {
			if cc.PlatformID == conn.PlatformID || conn.PlatformID == "-all" {
				total ++
				sc, ok := scsMap[cc.ID]
				if ok {
					c := convertObjectiveCommunity(cc, sc)
					err = community.CreateOrUpdate(c)(conn)
					if err == nil {
						m.SuccessCount ++
					} else {
						fmt.Printf("FAILED: %+v\n", err)
						m.FailuresCount ++
					}
				} else {
					fmt.Println("FAILED: strategy community not found")
					m.FailuresCount ++
				}
			}
			if err != nil {
				break
			}
		}
	}
	fmt.Printf("CapabilityCommunity: success %d + failures %d of total %d.\n", m.SuccessCount, m.FailuresCount, total)
	return
}

func convertObjectiveCommunity(
	cc capabilityCommunity.CapabilityCommunity,
	sc strategyCommunity.StrategyCommunity,
) (c community.Community) {
	c = community.Community{
		ID: cc.ID,

		Advocate: cc.Advocate,
		AccountabilityPartner: cc.CreatedBy,
		ChannelID: sc.ChannelID,
		CreatedAt: cc.CreatedAt,
		ModifiedAt: cc.ModifiedAt,
		CreatedBy: cc.CreatedBy,
		ModifiedBy: cc.CreatedBy,
		PlatformID: cc.PlatformID,
		Name: cc.Name,
		Description: cc.Description,
		CommunityKind: common.ObjectiveCommunity,
		DeactivatedAt: "",
		ParentCommunityID: sc.ParentCommunityChannelID,
	}
	fmt.Printf("community %s, channel %s\n", c.ID, c.ChannelID)
	return
}

// MigrateInitiativeCommunities - copies communities to Community table.
func MigrateInitiativeCommunities(conn common.DynamoDBConnection, m *migration.Migration) (err error) {
	var scs []strategyCommunity.StrategyCommunity
	err = conn.Dynamo.ScanTable(strategyCommunity.TableName(conn.ClientID), &scs)
	if err != nil {
		return
	}
	scsMap := indexStrategyCommunityByID(scs)
	// var acs []adaptiveCommunity.AdaptiveCommunity
	// err = conn.Dynamo.ScanTable(adaptiveCommunity.TableName(conn.ClientID), &acs)
	// if err != nil {
	// 	return
	// }
	// acsMap := groupAdaptiveCommunityByID(acs)
	var sics []strategyInitiativeCommunity.StrategyInitiativeCommunity
	err = conn.Dynamo.ScanTable(strategyInitiativeCommunity.TableName(conn.ClientID), &sics)
	if err != nil {
		return
	}

	total := 0
	if err == nil {
		for _, sic := range sics {
			if sic.PlatformID == conn.PlatformID || conn.PlatformID == "-all" {
				total ++
				sc, ok := scsMap[sic.ID]
				if ok {
					c := convertStrategyInitiativeCommunity(sic, sc)
					err = community.CreateOrUpdate(c)(conn)
					if err == nil {
						m.SuccessCount ++
					} else {
						fmt.Printf("FAILED: %+v\n", err)
						m.FailuresCount ++
					}
				} else {
					fmt.Println("FAILED: strategy community not found")
					m.FailuresCount ++
				}
				
			}
			if err != nil {
				break
			}
		}
	}
	fmt.Printf("StrategyCommunity: success %d + failures %d of total %d.\n", m.SuccessCount, m.FailuresCount, total)
	return
}

func convertStrategyInitiativeCommunity(
	sic strategyInitiativeCommunity.StrategyInitiativeCommunity,
	sc strategyCommunity.StrategyCommunity,
) (c community.Community) {
	c = community.Community{
		ID: sic.ID,
		Advocate: sic.Advocate,
		AccountabilityPartner: sic.CreatedBy,
		ChannelID: sc.ChannelID,
		CreatedAt: sic.CreatedAt,
		ModifiedAt: sic.ModifiedAt,
		CreatedBy: sic.CreatedBy,
		ModifiedBy: sic.CreatedBy,
		PlatformID: sic.PlatformID,
		Name: sic.Name,
		Description: sic.Description,
		CommunityKind: common.InitiativeCommunity,
		DeactivatedAt: "",
		ParentCommunityID: sic.CapabilityCommunityID,
	}

	fmt.Printf("community %s, channel %s", c.ID, c.ChannelID)
	// if sc.ParentCommunityChannelID != sic.CapabilityCommunityID {
	// 	fmt.Printf("sc.ParentCommunityChannelID=%s != sic.CapabilityCommunityID=%s", sc.ParentCommunityChannelID, sic.CapabilityCommunityID)
	// }
	fmt.Println()
	return
}
