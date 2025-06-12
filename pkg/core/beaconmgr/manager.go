package beaconmgr

import (
	"encoding/json"
	"gstrike/pkg/util"
	"net/http"
	"time"
)

type Beacon struct {
	ID         string    `json:"id"` // unique agent/beacon ID
	Hostname   string    `json:"hostname"`
	Username   string    `json:"username"` // current user context
	Domain     string    `json:"domain"`   // Active Directory domain, if any
	OS         string    `json:"os"`       // e.g. "windows", "linux"
	Arch       string    `json:"arch"`     // .eg. "amd64", "arm64"
	PID        int       `json:"pid"`      // process ID of the implant
	InternalIp string    `json:"internal_ip"`
	ExternalIP string    `json:"external_ip"`
	LastSeen   time.Time `json:"last_seen"`  // RFC3339 timestamp
	FirstSeen  time.Time `json:"first_seen"` // when the implant first checked in
}

var Beacons []Beacon

func NewBeacon(w http.ResponseWriter, r *http.Request) (Beacon, error) {
	var beacon Beacon
	if err := json.NewDecoder(r.Body).Decode(&beacon); err != nil {
		return beacon, err
	}

	id, _ := util.RandomString(12)
	beacon.ID = id
	beacon.FirstSeen = time.Now()
	beacon.LastSeen = time.Now()
	beacon.ExternalIP = r.RemoteAddr

	Beacons = append(Beacons, beacon)
	return Beacons[len(Beacons)-1], nil
}

