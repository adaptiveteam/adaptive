package main

import (
	"github.com/adaptiveteam/adaptive/daos/userObjective"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	// "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"os"
	// "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/daos/common"
	"fmt"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	} else {
		connGen := common.CreateConnectionGenFromEnv()
		platformID := common.PlatformID(os.Args[1])//adaptive_utils_go.NonEmptyEnv("")
		conn := connGen.ForPlatformID(platformID)
		var uos []userObjective.UserObjective
		err2 := conn.Dynamo.ScanTable(userObjective.TableName(conn.ClientID), &uos)
		count := 0
		errorCount := 0
		total := 0
		if err2 == nil {
			for _, uo := range uos {
				if uo.PlatformID == platformID {
					total ++
					newType := issues.DetectIssueType(uo).GetObjectiveType()
					fmt.Printf("%s .ObjectiveType == %s: ", uo.ID, uo.ObjectiveType)
					if uo.ObjectiveType == newType {
						fmt.Println("leaving unchanged")
					} else {
						fields, ok := uo.CollectEmptyFields()
						if ok {
							fmt.Println("-> ", newType)
							uo.ObjectiveType = newType
							err2 = userObjective.CreateOrUpdate(uo)(conn)
							count ++
						} else {
							fmt.Printf("FAILED - empty fields: %v\n", fields)
							errorCount ++
						}
					}
				}
				if err2 != nil {
					break
				}
			}
		}
		if err2 != nil {
			fmt.Printf("ERROR:\n%+v\n", err2)
		}
		fmt.Printf("Overall updated %d user objectives of total %d. Failed to update %d\n", count, total, errorCount)
	}
}

func usage() {
	fmt.Println("Usage: update-issues <platform-id>")
}
