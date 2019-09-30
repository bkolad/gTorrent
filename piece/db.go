package piece

import (
	"github.com/bkolad/gTorrent/db"
	"github.com/bkolad/gTorrent/torrent"
	"github.com/gchaincl/dotsql"
)

type repoDb struct {
	dot  *dotsql.DotSql
	db   *db.DB
	info *torrent.Info
}

func NewRepoDb(db *db.DB, info *torrent.Info) (*repoDb, error) {
	dot, err := dotsql.LoadFromFile("piece/queries.sql")
	if err != nil {
		return nil, err
	}
	return &repoDb{dot, db, info}, nil
}

func (r *repoDb) Save(pieceIndex uint32, data []byte) error {
	_, err := r.dot.Exec(r.db, "save-piece", r.info.InfoHash, pieceIndex, data)
	return err
}

func (r *repoDb) Get(pieceIndex uint32) ([]byte, error) {
	row, err := r.dot.QueryRow(r.db, "get-piece", r.info.InfoHash, pieceIndex)
	if err != nil {
		return nil, err
	}

	var data []byte
	err = row.Scan(&data)
	return data, nil
}
