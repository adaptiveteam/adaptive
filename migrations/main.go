package main

import (
	"github.com/adaptiveteam/adaptive/migrations/communities"
	"github.com/adaptiveteam/adaptive/daos/migration"
	"os"
	"github.com/adaptiveteam/adaptive/daos/common"
	"fmt"
)

type MigrationFunction struct {
	ID  string
	PerformMigration func(conn common.DynamoDBConnection, m *migration.Migration) (err error)
}

var migrations = []MigrationFunction {
	{ID: "010 Community - migrate users", PerformMigration: communities.MigrateUsers},
	{ID: "011 Community - migrate simple communities", PerformMigration: communities.MigrateSimpleCommunities},
	{ID: "012 Community - migrate initiative communities", PerformMigration: communities.MigrateInitiativeCommunities},
	{ID: "013 Community - migrate objective communities", PerformMigration: communities.MigrateCapabilityCommunities},
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	} else {
		connGen := common.CreateConnectionGenFromEnv()
		platformID := common.PlatformID(os.Args[1])//adaptive_utils_go.NonEmptyEnv("")
		conn := connGen.ForPlatformID(platformID)
		err2 := runMigrations(conn)
		
		if err2 != nil {
			fmt.Printf("ERROR:\n%+v\n", err2)
		}
	}
}

func usage() {
	fmt.Println("Usage: migrations <platform-id>")
	fmt.Println("       or")
	fmt.Println("       migrations -all")
}

func runMigrations(conn common.DynamoDBConnection) (err error) {
	for _, mf := range migrations {
		err = runMigration(mf, conn)
		if err != nil {
			break
		}
	}
	return
}

func runMigration(mf MigrationFunction, conn common.DynamoDBConnection) (err error) {
	var ms [] migration.Migration
	ms, err = migration.ReadOrEmpty(conn.PlatformID, mf.ID)(conn)
	if err != nil {
		return
	}
	m := migration.Migration{}
	if len(ms) > 0 {
		m = ms[0]
	} else {
		err = migration.CreateOrUpdate(m)(conn)
	}
	if err != nil {
		return
	}
	err = mf.PerformMigration(conn, &m)
	if err != nil {
		return
	}
	err = migration.CreateOrUpdate(m)(conn)
	return
}
