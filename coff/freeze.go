package coff

import (
	"encoding/binary"
	"reflect"

	"github.com/i9si-sistemas/wr/bin"
)

func freezeCommon(v reflect.Value, offset *uint32) error {
	if bin.Plain(v.Kind()) {
		*offset += uint32(binary.Size(v.Interface()))
		return nil
	}
	vv, ok := v.Interface().(Sizer)
	if ok {
		*offset += uint32(vv.Size())
		return bin.ErrWalkSkip
	}
	return nil
}
