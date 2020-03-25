package utilities

import (
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
	rdsHost := utils.NonEmptyEnv("RDS_HOST")
	GlobalRDSConfig := RDSConfig{
		Driver: "mysql", 
		ConnectionString: ConnectionString(
			rdsHost,
			utils.NonEmptyEnv("RDS_USER"),
			utils.NonEmptyEnv("RDS_PASSWORD"),
			utils.NonEmptyEnv("RDS_PORT"),
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
