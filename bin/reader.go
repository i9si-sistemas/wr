package bin

import (
	"encoding/binary"
	"io"
)

var Reader = reader{}

type reader struct{}

func (reader) LittleEndian(r io.Reader, data any) error {
	return binary.Read(r, binary.LittleEndian, data)
}
