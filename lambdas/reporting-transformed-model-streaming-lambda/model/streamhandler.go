package model

type StreamEventType string

var (
	StreamEventAdd    StreamEventType = "add"
	StreamEventEdit   StreamEventType = "edit"
	StreamEventDelete StreamEventType = "delete"
)

type StreamEntity struct {
	TableName string          `json:"table_name"`
	OldEntity interface{}     `json:"old_entity"`
	NewEntity interface{}     `json:"new_entity"`
	EventType StreamEventType `json:"event_type"`
}
