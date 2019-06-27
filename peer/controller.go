package peer

type Controller interface {
	Peer() (bool, Peer)
}
