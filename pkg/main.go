package main

import (
	"crypto/tls"
	"gstrike/pkg/middleware"
	"gstrike/pkg/relay"
	"gstrike/pkg/router"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("static/dist/assets"))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/dist/index.html")
	})
	r.HandleFunc("/tasks", router.PostTaskHandler).Methods("POST")

	// Agent API (secured by HMAC middleware)
	agentAPI := r.PathPrefix("/").Subrouter()
	agentAPI.Use(middleware.AgentHMACAuth)
	agentAPI.HandleFunc("/register", router.RegisterAgentHandler).Methods("POST")
	agentAPI.HandleFunc("/tasks/{agentID}", router.GetTasksHandler).Methods("GET")
	agentAPI.HandleFunc("/results", router.PostResultHandler).Methods("POST")

	r.HandleFunc("/ws", relay.WSHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":443",
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	log.Println("[*] Gstrike C2 server listening on https://localhost:443")
	err := server.ListenAndServeTLS("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}
}
