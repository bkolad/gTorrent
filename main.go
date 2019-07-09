package main

import (
	"fmt"
	"io/ioutil"

	i "github.com/bkolad/gTorrent/init"

	"github.com/bkolad/gTorrent/torrent"
	"github.com/bkolad/gTorrent/tracker"
)

func main() {
	conf := i.NewConf()
	initState := i.NewInitState()
	data, err := ioutil.ReadFile(conf.TorrentPath)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	dec := torrent.NewTorrentDecoder(string(data))
	info, err := dec.Decode()

	if err != nil {
		fmt.Println(err)
		return
	}

	tracker, _ := tracker.NewTracker(info, initState, conf)

	peers, err := tracker.Peers()
	if err != nil {
		fmt.Println(err)
		//	return
	}

	for _, p := range peers {
		fmt.Println(p.IP, p.Port)
	}
}
