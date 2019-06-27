package piece

type Manager interface {
	Next() int
	SetDone(int, []byte)
	SetInProgress(int)
	Done() []int
	InProgress() []int
	PieceLength() int
	LastPieceLength() int
	PieceHash(int) []byte
}
