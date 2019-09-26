package main

import (
	"fmt"
	"io/ioutil"

	"github.com/bkolad/gTorrent/db"
	i "github.com/bkolad/gTorrent/init"
	"github.com/bkolad/gTorrent/peer"
	"github.com/bkolad/gTorrent/piece"
	"github.com/bkolad/gTorrent/tracker"

	log "github.com/bkolad/gTorrent/logger"
	"github.com/bkolad/gTorrent/torrent"
)

func main() {

	db, err := db.InitDB()
	if err != nil {
		panic(err)
	}

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

	err = saveTorrent(db, info)
	if err != nil {
		panic(err)
	}

	pieceRepo, err := piece.NewRepoDb(db, info)
	if err != nil {
		panic(err)
	}

	err = pieceRepo.Save(22, []byte("some data"))
	if err != nil {
		panic(err)
	}
	data, err = pieceRepo.Get(22)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(data))

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

	pieceManager := piece.NewManager(*info, nil)
	repo := piece.NewRepo(pieceManager.PieceCount())

	peerManager := peer.NewManager(peerInfoChan, handshake, pieceManager, repo)
	go peerManager.ConnectToPeers()
	select {}
}
