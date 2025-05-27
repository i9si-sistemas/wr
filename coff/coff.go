package coff

import (
	"debug/pe"
	"encoding/binary"
	"errors"
	"reflect"
	"regexp"
	"sort"

	"github.com/i9si-sistemas/wr/bin"
)

type DataEntry struct {
	OffsetToData uint32
	Size1        uint32
	CodePage     uint32
	Reserved     uint32
}

type RelocationEntry struct {
	RVA         uint32
	SymbolIndex uint32
	Type        uint16
}

type Auxiliary [18]byte

type Symbol struct {
	Name           [8]byte
	Value          uint32
	SectionNumber  uint16
	Type           uint16
	StorageClass   uint8
	AuxiliaryCount uint8
	Auxiliaries    []Auxiliary
}

type StringsHeader struct {
	Length uint32
}

const (
	RT_ICON           = 3
	RT_GROUP_ICON     = 3 + 11
	RT_MANIFEST       = 24
	MASK_SUBDIRECTORY = 1 << 31
)

const (
	DT_PTR  = 1
	T_UCHAR = 12
)

var (
	STRING_WR  = [8]byte{'.', 'r', 's', 'r', 'c', 0, 0, 0}
	LANG_ENTRY = DirEntry{NameOrId: 0x0409}
)

type Sizer interface {
	Size() int64
}

type Coff struct {
	pe.FileHeader
	pe.SectionHeader32

	*Dir
	DataEntries []DataEntry
	Data        []PaddedData

	Relocations []RelocationEntry
	Symbols     []Symbol
	StringsHeader
	Strings []Sizer
}

func (coff *Coff) Arch(arch string) error {
	switch arch {
	case "386":
		coff.Machine = pe.IMAGE_FILE_MACHINE_I386
	case "amd64":
		coff.Machine = pe.IMAGE_FILE_MACHINE_AMD64
	case "arm":
		coff.Machine = pe.IMAGE_FILE_MACHINE_ARMNT
	case "arm64":
		coff.Machine = pe.IMAGE_FILE_MACHINE_ARM64
	default:
		return errors.New("coff: unknown architecture: " + arch)
	}
	return nil
}

const MEM_READ = 0x40000040

func NewWR() *Coff {
	coff := &Coff{
		pe.FileHeader{
			Machine:              pe.IMAGE_FILE_MACHINE_I386,
			NumberOfSections:     1,
			TimeDateStamp:        0,
			NumberOfSymbols:      1,
			SizeOfOptionalHeader: 0,
			Characteristics:      0x0104,
		},
		pe.SectionHeader32{
			Name:            STRING_WR,
			Characteristics: MEM_READ,
		},
		&Dir{},
		[]DataEntry{},
		[]PaddedData{},
		[]RelocationEntry{},
		[]Symbol{{
			Name:          STRING_WR,
			SectionNumber: 1,
			StorageClass:  3,
		}},
		StringsHeader{
			Length: uint32(binary.Size(StringsHeader{})),
		},
		[]Sizer{},
	}
	return coff
}

func (coff *Coff) AddResource(kind uint32, id uint16, data Sizer) {
	re := RelocationEntry{}
	switch coff.Machine {
	case pe.IMAGE_FILE_MACHINE_I386:
		re.Type = _IMAGE_REL_I386_DIR32NB
	case pe.IMAGE_FILE_MACHINE_AMD64:
		re.Type = _IMAGE_REL_AMD64_ADDR32NB
	case pe.IMAGE_FILE_MACHINE_ARMNT:
		re.Type = _IMAGE_REL_ARM_ADDR32NB
	case pe.IMAGE_FILE_MACHINE_ARM64:
		re.Type = _IMAGE_REL_ARM64_ADDR32NB
	}
	coff.Relocations = append(coff.Relocations, re)
	coff.SectionHeader32.NumberOfRelocations++
	entries0 := coff.Dir.DirEntries
	dirs0 := coff.Dir.Dirs
	i0 := sort.Search(len(entries0), func(i int) bool {
		return entries0[i].NameOrId >= kind
	})
	if i0 >= len(entries0) || entries0[i0].NameOrId != kind {
		entries0 = append(entries0[:i0], append([]DirEntry{{NameOrId: kind}}, entries0[i0:]...)...)
		dirs0 = append(dirs0[:i0], append([]Dir{{}}, dirs0[i0:]...)...)
		coff.Dir.NumberOfIdEntries++
	}
	coff.Dir.DirEntries = entries0
	coff.Dir.Dirs = dirs0

	dirs0[i0].DirEntries = append(dirs0[i0].DirEntries, DirEntry{NameOrId: uint32(id)})
	dirs0[i0].Dirs = append(dirs0[i0].Dirs, Dir{
		NumberOfIdEntries: 1,
		DirEntries:        DirEntries{LANG_ENTRY},
	})
	dirs0[i0].NumberOfIdEntries++

	n := 0
	for _, dir0 := range dirs0[:i0+1] {
		n += len(dir0.DirEntries)
	}
	n--

	coff.DataEntries = append(coff.DataEntries[:n], append([]DataEntry{{Size1: uint32(data.Size())}}, coff.DataEntries[n:]...)...)
	coff.Data = append(coff.Data[:n], append([]PaddedData{pad(data)}, coff.Data[n:]...)...)
}

// Freeze fills in some important offsets in resulting file.
func (coff *Coff) Freeze() {
	switch coff.SectionHeader32.Name {
	case STRING_WR:
		coff.freezeRSRC()
	}
}

func (coff *Coff) freezeCommon1(path string, offset, diroff uint32) (newdiroff uint32) {
	switch path {
	case "/Dir":
		coff.SectionHeader32.PointerToRawData = offset
		diroff = offset
	case "/Relocations":
		coff.SectionHeader32.PointerToRelocations = offset
		coff.SectionHeader32.SizeOfRawData = offset - diroff
	case "/Symbols":
		coff.FileHeader.PointerToSymbolTable = offset
	}
	return diroff
}

func (coff *Coff) freezeRSRC() {
	leafwalker := make(chan *DirEntry)
	go func() {
		for _, dir1 := range coff.Dir.Dirs {
			for _, dir2 := range dir1.Dirs {
				for i := range dir2.DirEntries {
					leafwalker <- &dir2.DirEntries[i]
				}
			}
		}
	}()

	var offset, diroff uint32
	bin.Walk(coff, func(v reflect.Value, path string) error {
		diroff = coff.freezeCommon1(path, offset, diroff)

		RE := regexp.MustCompile
		const N = `\[(\d+)\]`
		m := matcher{}
		switch {
		case m.Find(path, RE("^/Dir/Dirs"+N+"$")):
			coff.Dir.DirEntries[m[0]].OffsetToData = MASK_SUBDIRECTORY | (offset - diroff)
		case m.Find(path, RE("^/Dir/Dirs"+N+"/Dirs"+N+"$")):
			coff.Dir.Dirs[m[0]].DirEntries[m[1]].OffsetToData = MASK_SUBDIRECTORY | (offset - diroff)
		case m.Find(path, RE("^/DataEntries"+N+"$")):
			direntry := <-leafwalker
			direntry.OffsetToData = offset - diroff
		case m.Find(path, RE("^/DataEntries"+N+"/OffsetToData$")):
			coff.Relocations[m[0]].RVA = offset - diroff
		case m.Find(path, RE("^/Data"+N+"$")):
			coff.DataEntries[m[0]].OffsetToData = offset - diroff
		}

		return freezeCommon(v, &offset)
	})
}
