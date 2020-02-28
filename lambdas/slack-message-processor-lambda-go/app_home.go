package lambda

import (
	"encoding/json"
	// "github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

type AppHomeOpenedEvent = slackevents.AppHomeOpenedEvent

const (
	AppHomeOpened = slackevents.AppHomeOpened
)

type SlackAppHomeEvent struct {
	Token     string             `json:"token"`
	TeamID    string             `json:"team_id"`
	ApiAppID  string             `json:"api_app_id"`
	Event     AppHomeOpenedEvent `json:"event"`
	Type      string             `json:"type"`
	EventID   string             `json:"event_id"`
	EventTime json.Number        `json:"event_time"`
}
