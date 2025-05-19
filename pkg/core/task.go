package core

import (
	"fmt"
	"gstrike/pkg/util"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	TaskID     string    `json:"id"`
	BeaconID   string    `json:"beaconID"`
	Command    string    `json:"command"` // raw command or API name
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at"`
	Status     string    `json:"status"` // "success", "failed", "pending"
	Output     string    `json:"output"`
}

var Tasks []Task

func NewTask(cmd string) {
	t := Task{
		TaskID:    uuid.NewString(),
		BeaconID:  SelectedBeaconId,
		Command:   cmd,
		CreatedAt: time.Now(),
		Status:    "pending",
	}
	Tasks = append(Tasks, t)
	fmt.Printf("%s Queued new command '%s' for beacon: %s\n", util.PrintGood, cmd, t.BeaconID)
}
