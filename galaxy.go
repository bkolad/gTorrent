package main

import (
	"io/ioutil"

	"github.com/bkolad/gTorrent/db"
	i "github.com/bkolad/gTorrent/init"
	"github.com/bkolad/gTorrent/peer"
	"github.com/bkolad/gTorrent/piece"
	"github.com/bkolad/gTorrent/torrent"
	"github.com/bkolad/gTorrent/tracker"
	"github.com/gchaincl/dotsql"
)

type galaxy struct {
	db          *db.DB
	tracker     tracker.Tracker
	peerManager peer.Manager
	info        *torrent.Info
	peerInfos   chan torrent.PeerInfo
}

func newGalaxy(
	db *db.DB,
	initState i.State,
	conf i.Configuration) (*galaxy, error) {
	data, err := ioutil.ReadFile(conf.TorrentPath)
	if err != nil {
		return nil, err
	}

	dec := torrent.NewTorrentDecoder(string(data))
	info, err := dec.Decode()

	if err != nil {
		return nil, err
	}

	handshake := peer.NewHandshake(conf, info)
	tracker, _ := tracker.NewTracker(info, initState, conf)

	pieceRepo, err := piece.NewRepoDb(db, info)
	if err != nil {
		return nil, err
	}
	pieceManager := piece.NewManager(*info, pieceRepo)
	peerInfos := make(chan torrent.PeerInfo, 100)
	peerManager := peer.NewManager(peerInfos, handshake, pieceManager)

	return &galaxy{
		db:          db,
		tracker:     tracker,
		peerManager: peerManager,
		info:        info,
		peerInfos:   peerInfos,
	}, nil
}

func (g *galaxy) saveTorrent() error {
	dot, err := dotsql.LoadFromFile("queries.sql")
	if err != nil {
		return err
	}
	_, err = dot.Exec(g.db, "create-torrent", g.info.Name, "pending", g.info.InfoHash)

	if err != nil {
		return err
	}
	return nil
}

func (g *galaxy) retrievePeersFromTracker() error {
	peers, err := g.tracker.Peers()
	if err != nil {
		return err
	}
	for _, peer := range peers {
		g.peerInfos <- peer
	}
	close(g.peerInfos)
	return nil
}

func (g *galaxy) connectToPeers() {
	g.peerManager.ConnectToPeers()
}
