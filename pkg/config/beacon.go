package config

type BeaconConfig struct {
	Sleep     int
	Jitter    float64
	C2URL     string
	Cert      string
	UserAgent string
}
