package main

import (
	"fmt"
	"io/ioutil"

	i "github.com/bkolad/gTorrent/init"
	"github.com/bkolad/gTorrent/peer"
	"github.com/bkolad/gTorrent/piece"

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
	fmt.Println("Pieces Length:", info.Length)
	fmt.Println("Pieces Size:", info.PieceSize)
	tracker, _ := tracker.NewTracker(info, initState, conf)

	peers, err := tracker.Peers()
	if err != nil {
		fmt.Println(err)
		return
	}

	peerInfoChan := make(chan torrent.PeerInfo, 100)
	go func() {
		for _, peer := range peers {
			peerInfoChan <- peer
		}
		close(peerInfoChan)
	}()

	handshake := peer.NewHandshake(conf, info)
	pieceManager := piece.NewManager(*info)

	peerManager := peer.NewManager(peerInfoChan, handshake, pieceManager)
	go peerManager.ConnectToPeers()
	select {}
}
