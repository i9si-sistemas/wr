package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/i9si-sistemas/wr/ico"
)

var usage = `USAGE:

%s [-manifest FILE.exe.manifest] [-ico FILE.ico[,FILE2.ico...]] [OPTIONS...]
  Generates a .syso file with specified resources embedded in .wr section,
  aimed for consumption by Go linker when building Win32 excecutables.

The generated *.syso files should get automatically recognized by 'go build'
command and linked into an executable/library, as long as there are any *.go
files in the same directory.

OPTIONS:
`

func main() {
	var inputFile, iconFile, outputFile, Arch string
	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.StringVar(&inputFile, "manifest", "", "path to a Windows manifest file to embed")
	flags.StringVar(&iconFile, "ico", "", "comma-separated list of paths to .ico files to embed")
	flags.StringVar(&outputFile, "o", "", "name of output COFF (.res or .syso) file; if set to empty, will default to 'wr_windows_{arch}.syso'")
	flags.StringVar(&Arch, "arch", "amd64", "architecture of output file - one of: 386, amd64, [EXPERIMENTAL: arm, arm64]")
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flags.PrintDefaults()
	}
	_ = flags.Parse(os.Args[1:])
	if inputFile == "" && iconFile == "" {
		flags.Usage()
		os.Exit(1)
	}
	if outputFile == "" {
		outputFile = "wr_windows_" + Arch + ".syso"
	}

	if err := ico.Embed(outputFile, Arch, inputFile, iconFile); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
