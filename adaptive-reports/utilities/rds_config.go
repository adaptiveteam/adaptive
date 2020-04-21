package utilities

import (
	"strings"
	"database/sql"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
)

// RDSConfig configuration that allows to create an sql db connection
type RDSConfig struct {
	Driver           string
	ConnectionString string
}

// ReadRDSConfigFromEnv read config from env
func ReadRDSConfigFromEnv() RDSConfig {
	rdsEndpoint := utils.NonEmptyEnv("RDS_ENDPOINT")
	parts := strings.Split(rdsEndpoint, ":")

	rdsHost := parts[0]
	rdsPort := "3306"
	if len(parts) > 1 {
		rdsPort = parts[1]
	}
	GlobalRDSConfig := RDSConfig{
		Driver: "mysql", 
		ConnectionString: ConnectionString(
			rdsHost,
			utils.NonEmptyEnv("RDS_USER"),
			utils.NonEmptyEnv("RDS_PASSWORD"),
			rdsPort,
			utils.NonEmptyEnv("RDS_DB_NAME"),
		)}
	return GlobalRDSConfig
}

// SQLOpenUnsafe - opens a connection to an SQL database
func (RDSConfig RDSConfig)SQLOpenUnsafe() (conn *sql.DB) {
	conn = SQLOpenUnsafe(RDSConfig.Driver, RDSConfig.ConnectionString)
	// defer utilities.CloseUnsafe(db)
	return
}
