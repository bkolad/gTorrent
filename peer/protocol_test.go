package peer

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytesToBits(t *testing.T) {
	tests := []uint32{232, 24111, 3453463, 0, 255, 9999}
	for _, testValue := range tests {
		buf := &bytes.Buffer{}
		err := binary.Write(buf, binary.LittleEndian, testValue)
		require.NoError(t, err)
		bytes := buf.Bytes()
		require.Equal(t, bytes, bitsToBytes(bytesToBits(bytes)))
	}
}

func TestHaveToIndex(t *testing.T) {
}
