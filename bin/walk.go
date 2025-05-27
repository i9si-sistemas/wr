package bin

import (
	"errors"
	"fmt"
	"path"
	"reflect"
)

var ErrWalkSkip error = errors.New("walk skipped")

type Walker func(v reflect.Value, path string) error

func Walk(value any, walker Walker) error {
	err := walk(reflect.ValueOf(value), "/", walker)
	if err == ErrWalkSkip {
		err = nil
	}
	return err
}

func stopping(err error) bool {
	return err != nil && err != ErrWalkSkip
}

func walk(v reflect.Value, spath string, walker Walker) error {
	err := walker(v, spath)
	if err != nil {
		return err
	}
	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := range v.Len() {
			err = walk(v.Index(i), spath+fmt.Sprintf("[%d]", i), walker)
			if stopping(err) {
				return err
			}
		}
	case reflect.Interface:
		err = walk(v.Elem(), spath, walker)
		if stopping(err) {
			return err
		}
	case reflect.Struct:
		for i := range v.NumField() {
			vv := v.Field(i)
			err = walk(vv, path.Join(spath, v.Type().Field(i).Name), walker)
			if stopping(err) {
				return err
			}
		}
	default:
		return nil
	}
	return nil
}
