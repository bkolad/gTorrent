package torrent

import (
	"encoding/binary"
	"net"
)

// RSP is response from the tracker containing informations about remote peers.
type RSP struct {
	PeersInfo  []PeerInfo
	complete   int
	incomplete int
	interval   int
}

// PeerInfo holds information about remote peer.
type PeerInfo struct {
	IP   string
	Port uint16
}

// Decoder for bencoded tracker response.
type TrackerRSPDecoder interface {
	Decode() (*RSP, error)
}

type trackerDec struct {
	str string
}

// NewTrackerRspDecoder
func NewTrackerRspDecoder(str string) TrackerRSPDecoder {
	return &trackerDec{str}
}

func (dec *trackerDec) Decode() (*RSP, error) {
	p := NewParser(dec.str)
	ben, err := p.Parse()

	if err != nil {
		return nil, err
	}

	dict, ok := ben.(*bDict)
	if !ok {
		return nil, wrongTypeError("Torrent content ", "dictionary")
	}

	complete, _, err := intValue(dict, "complete")
	if err != nil {
		return nil, err
	}

	incomplete, _, err := intValue(dict, "incomplete")
	if err != nil {
		return nil, err
	}

	interval, _, err := intValue(dict, "interval")
	if err != nil {
		return nil, err
	}

	peers, err := fromDict(dict, "peers")
	if err != nil {
		return nil, err
	}

	rsp := RSP{}
	rsp.PeersInfo = parseIPAndPort(peers.PrettyString())
	rsp.complete = complete
	rsp.incomplete = incomplete
	rsp.interval = interval
	return &rsp, nil
}

func parseIPAndPort(peers string) []PeerInfo {
	var peerInfos []PeerInfo
	for i := 0; i <= len(peers)-6; i = i + 6 {
		addr := peers[i : i+6]
		ip := net.IPv4(addr[0], addr[1], addr[2], addr[3])
		port := binary.BigEndian.Uint16([]byte(addr[4:6]))
		peerInfo := PeerInfo{ip.String(), port}
		peerInfos = append(peerInfos, peerInfo)
	}
	return peerInfos
}
