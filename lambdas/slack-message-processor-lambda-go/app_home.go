package lambda

import "encoding/json"

type AppHomeOpenedEvent struct {
	Type           string      `json:"type"`
	User           string      `json:"user"`
	Channel        string      `json:"channel"`
	EventTimeStamp json.Number `json:"event_ts"`
}

const (
	AppHomeOpened = "app_home_opened"
)

type SlackAppHomeEvent struct {
	Token     string             `json:"token"`
	TeamId    string             `json:"team_id"`
	ApiAppId  string             `json:"api_app_id"`
	Event     AppHomeOpenedEvent `json:"event"`
	Type      string             `json:"type"`
	EventId   string             `json:"event_id"`
	EventTime json.Number        `json:"event_time"`
}
