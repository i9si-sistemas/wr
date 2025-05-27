package coff

import (
	"fmt"
	"os"
	"reflect"

	"github.com/i9si-sistemas/wr/bin"
)

func (c *Coff) WriteFile(filename string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	w := bin.Writer{W: out}

	bin.Walk(c, func(v reflect.Value, path string) error {
		if bin.Plain(v.Kind()) {
			w.WriteLE(v.Interface())
			return nil
		}
		vv, ok := v.Interface().(bin.SizedReader)
		if ok {
			w.WriteFromSized(vv)
			return bin.ErrWalkSkip
		}
		return nil
	})

	if w.Err != nil {
		return fmt.Errorf("error writing output file: %s", w.Err)
	}

	return nil
}
