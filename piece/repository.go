package piece

import (
	"sync"
)

type Repository interface {
	Save(piece uint32, data []byte)
	Get(piece, offset, size uint32) []byte
}

type inMemoryRepo struct {
	lock   sync.RWMutex
	pieces [][]byte
}

func NewRepo(n uint32) Repository {
	return &inMemoryRepo{pieces: make([][]byte, n)}
}

func (r *inMemoryRepo) Save(piece uint32, data []byte) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.pieces[piece] = data
}

func (r *inMemoryRepo) Get(piece, offset, size uint32) []byte {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.pieces[piece][offset : offset+size]
}
