package taskmgr

import (
	"gstrike/pkg/util"
	"time"
)

type Tasker interface {
	CreateTask(agentID string, cmd string) (Task, error)
	GetPendingTasks(agentID string) ([]Task, error)
	MarkTaskComplete(taskID string) error
}
type status int

const (
	success status = iota
	failed
	pending
)

type Task struct {
	TaskID     string    `json:"id"`
	BeaconID   string    `json:"beaconID"`
	Command    string    `json:"command"` // raw command or API name
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at"`
	Status     status    `json:"status"` // "success", "failed", "pending"
	Output     string    `json:"output"`
}

var Tasks []Task

func NewTask(cmd string, beacon string) error {
	id, err := util.RandomString(10)

	if err != nil {
		return err
	}

	t := Task{
		TaskID:    id,
		BeaconID:  beacon,
		Command:   cmd,
		CreatedAt: time.Now(),
		Status:    pending,
	}
	Tasks = append(Tasks, t)
	return nil
}

func UpdateTask(res Task) {
	for i := 0; i < len(Tasks); i++ {
		if Tasks[i].TaskID == res.TaskID {
			Tasks[i] = res
			break
		}
	}
}

func NextTasks(beaconId string) []Task {
	var tasks []Task
	for i := 0; i < len(Tasks); i++ {
		if Tasks[i].BeaconID == beaconId {
			tasks = Tasks[:i]
		}
	}
	return tasks
}
