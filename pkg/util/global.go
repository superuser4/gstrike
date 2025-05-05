package util

import (
	"sync"
	"time"
)

const BEACON_CALLBACK string = "beacon_callback"
const BEACON_REGISTER string = "beacon_register"

var (
	Agents       = make(map[string]Agent)
	Tasks        = make(map[string][]Task)
	Results      = make([]Result, 0)
	Mutex        sync.Mutex
	SharedSecret = "378432999013382759857861340953603067"
)

type Agent struct {
	Type       string    `json:"type"`
	ID         string    `json:"id"`
	Hostname   string    `json:"hostname"`
	ExternalIP string    `json:"external_ip"`
	LastSeen   time.Time `json:"last_seen"`
}

type Task struct {
	ID      string `json:"id"`
	Command string `json:"command"`
}

type Result struct {
	Type    string `json:"type"`
	AgentID string `json:"agent_id"`
	TaskID  string `json:"task_id"`
	Output  string `json:"output"`
}
