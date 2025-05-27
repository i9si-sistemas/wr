package ico

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/i9si-sistemas/wr/bin"
	"github.com/i9si-sistemas/wr/coff"
)

type _GRPICONDIR struct {
	Dir
	Entries []_GRPICONDIRENTRY
}

func (group _GRPICONDIR) Size() int64 {
	return int64(binary.Size(group.Dir) + len(group.Entries)*binary.Size(group.Entries[0]))
}

type _GRPICONDIRENTRY struct {
	DirEntryCommon
	Id uint16
}

func Embed(
	outputFilePath,
	arch,
	inputFilePath,
	iconFilePath string,
) error {
	lastid := uint16(0)
	newid := func() uint16 {
		lastid++
		return lastid
	}

	out := coff.NewWR()
	err := out.Arch(arch)
	if err != nil {
		return err
	}

	if inputFilePath != "" {
		manifest, err := bin.SizedOpen(inputFilePath)
		if err != nil {
			return fmt.Errorf("wr: error opening manifest file '%s': %s", inputFilePath, err)
		}
		defer manifest.Close()

		id := newid()
		out.AddResource(coff.RT_MANIFEST, id, manifest)
	}
	if iconFilePath != "" {
		for iconFileNameSingle := range strings.SplitSeq(iconFilePath, ",") {
			f, err := addIcon(out, iconFileNameSingle, newid)
			if err != nil {
				return err
			}
			defer f.Close()
		}
	}

	out.Freeze()

	return out.WriteFile(outputFilePath)
}

func addIcon(out *coff.Coff, fname string, newid func() uint16) (io.Closer, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	icons, err := DecodeHeaders(f)
	if err != nil {
		f.Close()
		return nil, err
	}

	if len(icons) > 0 {
		group := _GRPICONDIR{Dir: Dir{
			Type:  1,
			Count: uint16(len(icons)),
		}}
		gid := newid()
		for _, icon := range icons {
			id := newid()
			r := io.NewSectionReader(f, int64(icon.ImageOffset), int64(icon.BytesInRes))
			out.AddResource(coff.RT_ICON, id, r)
			group.Entries = append(group.Entries, _GRPICONDIRENTRY{icon.DirEntryCommon, id})
		}
		out.AddResource(coff.RT_GROUP_ICON, gid, group)
	}

	return f, nil
}
