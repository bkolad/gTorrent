package peer

type Peer interface {
	doHandshake()
}

type peer struct {
	msgs chan MSG
	net  Network
}

func (p peer) doHandshake() {
	err := p.net.SendHandshake()
	if err != nil {
		p.msgs <- handshakeError{}
	}
}
