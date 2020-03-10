package sqlconnector

import "github.com/jinzhu/gorm"

func NewMySqlConnection() (db *gorm.DB, err error) {
	connector := mySqlConnector{}
	return connector.GetConnection()
}
