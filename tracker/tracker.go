package tracker

type PeerInfo struct {
}

type Tracker interface {
	Peers() []PeerInfo
}

type tracker struct {
}

func (t *tracker) Peers() []PeerInfo {
	return nil
}
