package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/i9si-sistemas/assert"
	"github.com/i9si-sistemas/wr/command"
)

const name = "wr_windows_amd64.syso"

func TestBuildSucceeds(t *testing.T) {
	tests := []struct {
		comment string
		args    []string
	}{{
		comment: "icon",
		args:    []string{"-ico", "i9si.ico"},
	}, {
		comment: "manifest",
		args:    []string{"-manifest", "manifest.xml"},
	}, {
		comment: "manifest & icon",
		args:    []string{"-manifest", "manifest.xml", "-ico", "i9si.ico"},
	}}
	for _, tt := range tests {
		t.Run(tt.comment, func(t *testing.T) {
			dir, err := os.Getwd()
			assert.NoError(t, err)
			dir = filepath.Join(dir, "build")

			os.Stdout.Write([]byte("-- compiling resource(s)...\n"))
			defer os.Remove(filepath.Join(dir, name))
			cmd :=	command.New().
				Execute("go", "run", "../wr.go", "-arch", "amd64").
				AppendArgs(tt.args...).WithDir(dir)
			err = cmd.Run()
			assert.NoError(t, err)
			_, err = os.Stat(filepath.Join(dir, name))
			assert.NoError(t, err)

			defer os.Setenv("GOOS", os.Getenv("GOOS"))
			defer os.Setenv("GOARCH", os.Getenv("GOARCH"))
			os.Setenv("GOOS", "windows")
			os.Setenv("GOARCH", "amd64")
			_, err = os.Stdout.Write([]byte("-- compiling app...\n"))
			assert.NoError(t, err)
			cmd = command.New().Execute("go", "build", "../").WithDir(dir)
			err = cmd.Run()
			assert.NoError(t, err)

			if runtime.GOOS == "windows" && runtime.GOARCH == "amd64" {
				_, err = os.Stdout.Write([]byte("-- running app...\n"))
				assert.NoError(t, err)
				cmd = command.New().WithPath("wr.exe").WithDir(dir)
				out, err := cmd.CombinedOutput()
				output := string(out)
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "exit status 1")
				assert.Equal(t, output, "USAGE:\n\nwr.exe [-manifest FILE.exe.manifest] [-ico FILE.ico[,FILE2.ico...]] [OPTIONS...]\n  Generates a .syso file with specified resources embedded in .wr section,\n  aimed for consumption by Go linker when building Win32 excecutables.\n\nThe generated *.syso files should get automatically recognized by 'go build'\ncommand and linked into an executable/library, as long as there are any *.go\nfiles in the same directory.\n\nOPTIONS:\n  -arch string\n    \tarchitecture of output file - one of: 386, amd64, [EXPERIMENTAL: arm, arm64] (default \"amd64\")\n  -ico string\n    \tcomma-separated list of paths to .ico files to embed\n  -manifest string\n    \tpath to a Windows manifest file to embed\n  -o string\n    \tname of output COFF (.res or .syso) file; if set to empty, will default to 'wr_windows_{arch}.syso'\n")
			}
		})
	}
}
