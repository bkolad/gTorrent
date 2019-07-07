package init

type State struct {
	Downloaded int
	Uploaded   int
	Left       int
}

func NewInitState() State {
	return State{0, 0, 10000}
}
