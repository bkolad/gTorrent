package torrent

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseIPAndPort(t *testing.T) {
	peerInfos := []PeerInfo{
		{"198.48.140.173", 21025},
		{"94.21.25.78", 52111},
		{"81.182.215.168", 65534},
	}

	buf := new(bytes.Buffer)
	for _, pi := range peerInfos {
		binary.Write(buf, binary.BigEndian, net.ParseIP(pi.IP).To4())
		binary.Write(buf, binary.BigEndian, pi.Port)
	}

	parsedInfos := parseIPAndPort(string(buf.Bytes()))
	for i := range peerInfos {
		require.Equal(t, *parsedInfos[i], peerInfos[i])
	}
}

type testCase struct {
	path       string
	isSuccess  bool
	complete   int
	incomplete int
	peersInfo  []PeerInfo
}

func TestTrackerRSPDecode(t *testing.T) {
	path := "../testData/"

	tests := []testCase{
		{path + "announceSuccess.rsp",
			true,
			2480,
			45,
			[]PeerInfo{{"198.48.140.173", 21025}, {"2.95.177.152", 51413}},
		},
		{path + "announceFailure.rsp",
			false,
			0,
			0,
			nil,
		},
	}

	for _, test := range tests {
		if test.isSuccess {
			testSuccess(t, test)
		} else {
			testFailure(t, test)
		}
	}
}

func testSuccess(t *testing.T, test testCase) {
	data, err := ioutil.ReadFile(test.path)
	require.NoError(t, err)

	dec := NewTrackerRspDecoder(string(data))
	rsp, err := dec.Decode()
	require.NoError(t, err)

	require.Equal(t, rsp.complete, test.complete)
	require.Equal(t, rsp.incomplete, test.incomplete)

	for _, p := range test.peersInfo {
		require.True(t, contains(rsp.PeersInfo, p))
	}
}

func testFailure(t *testing.T, test testCase) {
	data, err := ioutil.ReadFile(test.path)
	require.NoError(t, err)

	dec := NewTrackerRspDecoder(string(data))
	_, err = dec.Decode()
	require.Error(t, err)
}

func contains(allPeerInfo []*PeerInfo, peerInfo PeerInfo) bool {
	for _, pi := range allPeerInfo {
		if *pi == peerInfo {
			return true
		}
	}
	return false
}
