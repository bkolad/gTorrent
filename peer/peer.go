package peer

type Peer interface {
	doHandshake()
}

type peer struct {
	msgs chan MSG
	net  Network
}

//TODO add listener

func (p peer) doHandshake() {
	p.net.SendHandshake()
}
