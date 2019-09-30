package main

import (
	"github.com/bkolad/gTorrent/db"
	i "github.com/bkolad/gTorrent/init"

	log "github.com/bkolad/gTorrent/logger"
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

	galaxy, err := newGalaxy(db, initState, conf)
	if err != nil {
		panic(err)
	}

	err = galaxy.saveTorrent()
	if err != nil {
		panic(err)
	}

	go func() {
		galaxy.retrievePeersFromTracker()
	}()

	go func() {
		galaxy.connectToPeers()
	}()
	select {}
}
