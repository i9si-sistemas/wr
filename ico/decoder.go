package ico

import (
	"fmt"
	"io"

	"github.com/i9si-sistemas/wr/bin"
)

func DecodeHeaders(r io.Reader) ([]DirEntry, error) {
	var hdr Dir
	err := bin.Reader.LittleEndian(r, &hdr)
	if err != nil {
		return nil, err
	}
	if hdr.Reserved != 0 || hdr.Type != 1 {
		return nil, fmt.Errorf("bad magic number")
	}

	entries := make([]DirEntry, hdr.Count)
	for i := range len(entries) {
		err = bin.Reader.LittleEndian(r, &entries[i])
		if err != nil {
			return nil, err
		}
	}
	return entries, nil
}
