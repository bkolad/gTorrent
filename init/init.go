package init

// Global state of the local peer.
type State struct {
	Downloaded int
	Uploaded   int
	Left       int
}

// NewInitState retrieves initial state from the DB.
func NewInitState() State {
	return State{0, 0, 10000}
}
