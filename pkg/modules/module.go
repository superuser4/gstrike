package modules

type Module interface {
	Name() string
	Execute(beaconID string, args []string) (string, error)
}
