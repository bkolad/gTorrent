package torrent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const torrentContent = "d8:announce39:http://torrent.ubuntu.com:6969/announce" +
	"13:announce-listll39:http://torrent.ubuntu.com:6969/announcee" +
	"l44:http://ipv6.torrent.ubuntu.com:6969/announceee" +
	"7:comment29:Ubuntu CD releases.ubuntu.com13:creation datei1445507299e" +
	"4:info" +
	"d6:lengthi1e" +
	"4:name30:ubuntu-15.10-desktop-amd64.iso" +
	"12:piece lengthi524288e" +
	"6:pieces20:aaaaaaaaaaaaaaaaaaaaee"

func TestDecode(t *testing.T) {
	dec := NewTorrentDecoder(torrentContent)
	decodedInfo, err := dec.Decode()
	require.NoError(t, err)
	info := info()
	require.Equal(t, decodedInfo, info)
}

func info() *Info {
	info := new(Info)
	info.Announce = "http://torrent.ubuntu.com:6969/announce"
	info.AnnounceList = [][]string{
		{"http://torrent.ubuntu.com:6969/announce"},
		{"http://ipv6.torrent.ubuntu.com:6969/announce"},
	}
	info.files = nil
	info.Length = 1
	info.InfoHash =
		[]byte{98, 194, 202, 18, 139, 80, 209, 76, 165,
			195, 230, 13, 19, 178, 186, 49, 28, 102, 203, 88}
	info.name = "ubuntu-15.10-desktop-amd64.iso"
	var pieceHashe []byte
	for i := 0; i < 20; i++ {
		pieceHashe = append(pieceHashe, 97)
	}
	info.PieceHashes = [][]byte{pieceHashe}
	info.PieceSize = 524288
	return info
}
