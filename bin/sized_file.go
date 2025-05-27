package bin

import (
	"io"
	"os"
)

func SizedOpen(filename string) (*SizedFile, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	section := SizedSection{
		file: f,
	}
	section.SectionReader = io.NewSectionReader(f, 0, info.Size())
	return &SizedFile{section}, nil
}

type SizedReader interface {
	io.Reader
	Size() int64
}

type SizedFile struct {
	SizedSection
}

type SizedSection struct {
	file *os.File
	*io.SectionReader
}

func (r SizedSection) ReadSection(p []byte) (int, error) {
	return r.Read(p)
}

func (r SizedSection) CloseFile() error {
	return r.file.Close()
}

func (r SizedSection) SectionSize() int64 {
	return r.SectionReader.Size()
}

func (r *SizedFile) Read(p []byte) (int, error) {
	return r.ReadSection(p)
}
func (r *SizedFile) Size() int64 {
	return r.SectionSize()
}
func (r *SizedFile) Close() error {
	return r.CloseFile()
}