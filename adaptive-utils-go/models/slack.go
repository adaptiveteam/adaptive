package models

type GroupLeftEvent struct {
	Type    string `json:"type"` // group_left
	Channel string `json:"channel"`
	ActorId string `json:"actor_id"`
	EventTs string `json:"event_ts"`
}
