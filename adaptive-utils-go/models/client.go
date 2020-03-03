package models

import (
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/clientPlatformToken"
)

type ClientContact struct {
	ContactFirstName string `json:"contact_first_name"`
	ContactLastName  string `json:"contact_last_name"`
	ContactMail      string `json:"contact_mail"`
}

type ClientPlatform struct {
	PlatformName  common.PlatformName `json:"platform_name"` // should be slack or ms-teams
	PlatformToken string       `json:"platform_token"`
}

type ClientPlatformRequest struct {
	// team_id is either the AppID (api_app_id) or TeamID (team_id) from Slack
	TeamID TeamID `json:"team_id"`
	Org string `json:"platform_org"`
}

type ClientPlatformToken = clientPlatformToken.ClientPlatformToken

// ClientPlatformTokenTableSchema is the schema of _adaptive_client_config table
type ClientPlatformTokenTableSchema struct {
	Name string
	// platform_id is the hash key 
}
// ClientPlatformTokenTableSchemaForClientID creates table schema given client id
func ClientPlatformTokenTableSchemaForClientID(clientID string) ClientPlatformTokenTableSchema {
	return ClientPlatformTokenTableSchema{
		Name: clientID + "_adaptive_client_config",
	}
}
