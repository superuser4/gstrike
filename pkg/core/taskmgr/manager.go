package taskmgr

import (
	"fmt"
	"gstrike/pkg/util"
	"time"
)

// Task creation + results parsing and tasks updating

// Use for both DB and in memory tasking? You can also write a mock Tasker for testing
type Tasker interface {
	CreateTask(agentID string, cmd string) (Task, error)
	GetPendingTasks(agentID string) ([]Task, error)
	MarkTaskComplete(taskID string) error
}

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

func NewTask(cmd string, beacon string) {
	id, _ := util.RandomString(10)
	t := Task{
		TaskID:    id,
		BeaconID:  beacon,
		Command:   cmd,
		CreatedAt: time.Now(),
		Status:    "pending",
	}
	Tasks = append(Tasks, t)
	fmt.Printf("%s Queued new command '%s' for beacon: %s\n", util.PrintGood, cmd, t.BeaconID)
}

func UpdateTask(res Task) {
	for i := 0; i < len(Tasks); i++ {
		if Tasks[i].TaskID == res.TaskID {
			Tasks[i] = res
			break
		}
	}
}

func NextTask(beaconId string) Task {
	var next Task
	for i := 0; i < len(Tasks); i++ {
		if Tasks[i].BeaconID == beaconId && Tasks[i].Status == "pending" {
			next = Tasks[i]
			break
		}
	}
	return next
}

func PrintTasks() {
	util.ListDisplay([]string{"Task ID", "Beacon ID", "Command", "Status", "Created At", "Finished At", "Output"})

	for i := 0; i < len(Tasks); i++ {
		t := Tasks[i]
		fmt.Printf("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\n", t.TaskID, t.BeaconID, t.Command, t.Status, t.CreatedAt, t.FinishedAt, t.Output)
	}
}

func PrintTask(beacon *string) {
	util.ListDisplay([]string{"Task ID", "Beacon ID", "Command", "Status", "Created At", "Finished At", "Output"})
	var exist bool

	for i := 0; i < len(Tasks); i++ {
		t := Tasks[i]
		if t.BeaconID == *beacon {
			exist = true
			fmt.Printf("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\n", t.TaskID, t.BeaconID, t.Command, t.Status, t.CreatedAt, t.FinishedAt, t.Output)
		}
	}
	if !exist {
		fmt.Printf("%s No such beacon ID found...\n", util.PrintBad)
	}
}
