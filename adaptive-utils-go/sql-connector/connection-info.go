package sqlconnector

import (
	"fmt"
	"strings"
	// postgres driver
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
)

var (
	// NonEmptyEnv is the environment that is checked for emptiness
	NonEmptyEnv = core.NonEmptyEnv
)

// ConnectionInfo describes connection information
type ConnectionInfo struct {
	Driver string
	// Endpoint string
	Host string
	Port string // might be empty
	UserName string
	Password string
	DatabaseName string
}

// SplitEndpoint extracts host:port from endpoint
func SplitEndpoint(endpoint string) (host, port string) {
	rdsEndpoint := NonEmptyEnv("RDS_ENDPOINT")
	parts := strings.Split(rdsEndpoint, ":")

	host = parts[0]
	if len(parts) > 1 {
		port = parts[1]
	}
	return
}
// ReadConnectionInfoFromEnv reads connection info from environment
func ReadConnectionInfoFromEnv() (connectionInfo ConnectionInfo) {
	// Endpoint string
	rdsEndpoint := NonEmptyEnv("RDS_ENDPOINT")
	rdsHost, rdsPort := SplitEndpoint(rdsEndpoint)
	connectionInfo = ConnectionInfo{
		Driver: ReadDriverFromEnv(),
		Host: rdsHost,
		Port: rdsPort, // might be empty
		UserName: NonEmptyEnv("RDS_USER"),
		Password: NonEmptyEnv("RDS_PASSWORD"),
		DatabaseName: NonEmptyEnv("RDS_DB_NAME"),
	}
	return
}

// ReadDriverFromEnv is the single place where we specify which driver we use.
func ReadDriverFromEnv() string {
	return "mysql" // "postgres"
}

func prependUnlessEmpty(prefix, port string) (res string) {
	if port == "" {
		res = ""
	} else {
		res = prefix + port
	}
	return
}

// ConnectionString constructs connection string
func (ci ConnectionInfo)ConnectionString() (cs string) {
	switch ci.Driver {
	case "mysql":
		cs = ci.UserName + ":" + ci.Password + "@tcp(" + ci.Host + prependUnlessEmpty(":", ci.Port) + ")/" + ci.DatabaseName + "?charset=utf8&parseTime=True&loc=Local"
	case "postgres":
		cs = fmt.Sprintf("host=%s dbname=%s user=%s password=%s", ci.Host, ci.DatabaseName , ci.UserName, ci.Password) +
			prependUnlessEmpty(" port=", ci.Port)
	}
	return
}
// ToRDSConfig converts to RDSConfig
func (ci ConnectionInfo)ToRDSConfig() RDSConfig {
	return RDSConfig{
		Driver: ci.Driver,
		ConnectionString: ci.ConnectionString(),
	}
}
