package init

// Configuration of the local peer.
type Configuration struct {
	Port        int
	PeerID      string
	TorrentPath string
}

// NewConf creates default configuration.
func NewConf() Configuration {
	c := Configuration{}
	c.Port = 6881
	c.PeerID = peerID()
	c.TorrentPath = "testData/ubuntu-19.04-desktop-amd64.iso.torrent"
	return c
}

func peerID() string {
	return "-GT001-" + "d0p22uiake0bd"
}
