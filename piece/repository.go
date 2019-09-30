package piece

import (
	"sync"
)

type Repository interface {
	Save(piece uint32, data []byte) error
	Get(piece uint32) ([]byte, error)
}

type inMemoryRepo struct {
	lock   sync.RWMutex
	pieces [][]byte
}

func NewRepo(n uint32) Repository {
	return &inMemoryRepo{pieces: make([][]byte, n)}
}

func (r *inMemoryRepo) Save(piece uint32, data []byte) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.pieces[piece] = data
	return nil
}

func (r *inMemoryRepo) Get(piece uint32) ([]byte, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.pieces[piece], nil
}
