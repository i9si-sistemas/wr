package coff

type Dir struct {
	Characteristics      uint32
	TimeDateStamp        uint32
	MajorVersion         uint16
	MinorVersion         uint16
	NumberOfNamedEntries uint16
	NumberOfIdEntries    uint16
	DirEntries
	Dirs
}

type DirEntry struct {
	NameOrId     uint32
	OffsetToData uint32
}

type DirEntries []DirEntry

type Dirs []Dir
