package core

import (
	"encoding/json"
	"fmt"
	"gstrike/pkg/util"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// beacon handler

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
var SelectedBeaconId = ""

func NewBeacon(w http.ResponseWriter, r *http.Request) (Beacon, error) {
	var beacon Beacon
	if err := json.NewDecoder(r.Body).Decode(&beacon); err != nil {
		fmt.Printf("%s Json decode error: %v\n", util.PrintBad, err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return beacon, err
	}

	beacon.ID = uuid.New().String()
	beacon.FirstSeen = time.Now()
	beacon.LastSeen = time.Now()
	beacon.ExternalIP = r.RemoteAddr

	Beacons = append(Beacons, beacon)
	return Beacons[len(Beacons)-1], nil
}
func UpdateBeacon() {}

func ListBeacons(list *string) {
	var lol = []string{"ID", "External IP", "Internal IP", "Hostname", "Username", "Domain", "OS", "Arch", "PID", "LastSeen", "FirstSeen"}
	util.ListDisplay(lol)

	for i := 0; i < len(Beacons); i++ {
		c := Beacons[i]

		if *list != "" && c.ID != *list {
			continue
		}
		fmt.Printf("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t", c.ID, c.ExternalIP, c.InternalIp, c.Hostname, c.Username, c.Domain, c.OS, c.Arch, strconv.Itoa(c.PID), c.LastSeen, c.FirstSeen)
	}
	fmt.Printf("\n")
}

func SelectBeacon(beacon *string) {
	if *beacon != "" {
		var exist bool
		for i := 0; i < len(Beacons); i++ {
			if Beacons[i].ID == *beacon {
				exist = true
				break
			}
		}
		if exist {
			fmt.Printf("%s Beacon selected: %s\n", util.PrintStatus, *beacon)
			SelectedBeaconId = *beacon
		}
	}
}

func ResetBeacon() {
	SelectedBeaconId = ""
}
