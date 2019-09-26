package main

import (
	"github.com/bkolad/gTorrent/db"
	"github.com/bkolad/gTorrent/torrent"
	"github.com/gchaincl/dotsql"
)

func saveTorrent(db *db.DB, info *torrent.Info) error {
	dot, err := dotsql.LoadFromFile("queries.sql")
	if err != nil {
		return err
	}

	_, err = dot.Exec(db, "create-torrent", info.Name, "pending", info.InfoHash)

	if err != nil {
		return err
	}
	return nil
}
