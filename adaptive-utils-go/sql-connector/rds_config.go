package sqlconnector

import (
	"github.com/jinzhu/gorm"
	"log"
	"database/sql"
)

// RDSConfig configuration that allows to create an sql db connection
type RDSConfig struct {
	Driver           string
	ConnectionString string
}

// ReadRDSConfigFromEnv read config from env
func ReadRDSConfigFromEnv() RDSConfig {
	return ReadConnectionInfoFromEnv().ToRDSConfig()
}

// SQLOpenUnsafe - opens a connection to an SQL database
func (RDSConfig RDSConfig) SQLOpenUnsafe() (conn *sql.DB) {
	db, err2 := RDSConfig.SQLOpen()
	if err2 != nil {
		log.Panicf("Error creating database: %+v", err2)
	}
	return db
}

// SQLOpen - opens a connection to an SQL database
func (RDSConfig RDSConfig) SQLOpen() (conn *sql.DB, err error) {
	conn, err = sql.Open(RDSConfig.Driver, RDSConfig.ConnectionString)
	return
}

// GormOpen - opens a connection to an SQL database
func (RDSConfig RDSConfig) GormOpen() (conn *gorm.DB, err error) {
	conn, err = gorm.Open(RDSConfig.Driver, RDSConfig.ConnectionString)
	return
}
