package core

import "time"

type Task struct {
	TaskID     string    `json:"id"`
	BeaconID   string    `json:"beaconID"`
	Command    string    `json:"command"` // raw command or API name
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at"`
	Status     string    `json:"status"` // "success", "failed", "pending"
	Output     string    `json:"output"`
}

var Tasks = make(map[string][]Task)

func NewTask() Task {
	t := Task{}
	return t
}
