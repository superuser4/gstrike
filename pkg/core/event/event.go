package event

import "time"

type EventLog struct {
	EventID     string    `json:"event_id"`
	Timestamp   time.Time `json:"timestamp"`
	Source      string    `json:"source"` // e.g. "agent", "operator", "system"
	AgentID     string    `json:"agent_id,omitempty"`
	OperatorID  string    `json:"operator_id,omitempty"`
	EventType   string    `json:"event_type"` // e.g. "login", "task_sent", "error"
	Description string    `json:"description"`
}
