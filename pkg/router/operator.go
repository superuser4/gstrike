package router

import (
	"encoding/json"
	"gstrike/pkg/util"
	"net/http"

	"github.com/google/uuid"
)

// for the operators to queue tasks to agents
func PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task struct {
		AgentID string `json:"agent_id"`
		Command string `json:"command"`
	}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid task", http.StatusBadRequest)
		return
	}

	util.Mutex.Lock()
	defer util.Mutex.Unlock()

	newTask := util.Task{
		ID:      uuid.New().String(),
		Command: task.Command,
	}
	util.Tasks[task.AgentID] = append(util.Tasks[task.AgentID], newTask)

	w.WriteHeader(http.StatusCreated)
}
