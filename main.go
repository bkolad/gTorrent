package main

import (
	"fmt"
	"io/ioutil"

	i "github.com/bkolad/gTorrent/init"
	"github.com/bkolad/gTorrent/peer"

	log "github.com/bkolad/gTorrent/logger"
	"github.com/bkolad/gTorrent/torrent"
	"github.com/bkolad/gTorrent/tracker"
)

func main() {

	log.Info("Starting gTorrent..")
	conf := i.NewConf()
	initState := i.NewInitState()
	log.Debug("Local peer ID: " + conf.PeerID)

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
		return
	}

	peerInfoChan := make(chan torrent.PeerInfo, 100)
	go func() {
		for _, p := range peers {
			peerInfoChan <- p
		}
	}()
	h := peer.NewHandshake(conf, info)
	peerManager := peer.NewManager(peerInfoChan, h)
	go peerManager.ConnectToPeers()
	//	h := p.NewHandshake(conf, info)

	//	net := network.NewNetwork(peers[0], h)
	//	net.Send()
	select {}
}
