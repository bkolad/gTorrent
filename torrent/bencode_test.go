package torrent

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const PATH = "../testData/"

func TestTorrentFiles(t *testing.T) {
	torrents, err := torrents()
	require.NoError(t, err)

	for _, torrent := range torrents {
		data, err := ioutil.ReadFile(torrent)
		require.NoError(t, err)

		ben, err := benData(string(data))
		require.NoError(t, err)

		ben2, err := benData(ben.String())
		require.NoError(t, err)
		require.NotNil(t, ben2)

		require.Equal(t, ben, ben2)
	}
}

func benData(data string) (Bencode, error) {
	p := NewParser(string(data))
	return p.Parse()
}

func torrents() ([]string, error) {
	matches, err := filepath.Glob(PATH + "*.torrent")
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, match := range matches {
		paths = append(paths, match)
	}
	return paths, nil
}
