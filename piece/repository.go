package piece

type repository interface {
	Save(piece uint32, data []byte)
	Get(piece, offset, size uint32) []byte
}
