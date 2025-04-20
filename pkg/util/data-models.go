package util

import (
	"sync"
	"time"
)

var (
	agents  = make(map[string]Agent)
	tasks   = make(map[string][]Task)
	results = make([]Result, 0)
	mutex   sync.Mutex
)

type Agent struct {
	ID       string    `json:"id"`
	Hostname string    `json:"hostname"`
	IP       string    `json:"ip"`
	LastSeen time.Time `json:"last_seen"`
}

type Task struct {
	ID      string `json:"id"`
	Command string `json:"command"`
}

type Result struct {
	AgentID string `json:"agent_id"`
	TaskID  string `json:"task_id"`
	Output  string `json:"output"`
}
