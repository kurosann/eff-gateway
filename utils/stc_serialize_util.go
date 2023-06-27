package utils

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func Encode(value interface{}) ([]byte, error) {
	kind := reflect.TypeOf(value).Kind()
	v := reflect.ValueOf(value)
	if kind == reflect.Pointer {
		v = v.Elem()
		kind = v.Kind()
	}
	if kind == reflect.String {
		return []byte(value.(string)), nil
	}
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(v.Uint(), 10)), nil
	}
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(value)
	if err != nil {
		return []byte{}, err
	}
	return buff.Bytes(), nil
}
func Decode(bs []byte, value interface{}) error {
	kind := reflect.TypeOf(value).Kind()
	if kind != reflect.Pointer {
		return errors.New("value not a pointer")
	}
	kind = reflect.ValueOf(value).Elem().Kind()
	v := reflect.ValueOf(value).Elem()
	var err error
	switch kind {
	case reflect.String:
		v.SetString(string(bs))
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i int64
		i, err = strconv.ParseInt(string(bs), 10, 64)
		v.SetInt(i)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var i uint64
		i, err = strconv.ParseUint(string(bs), 10, 64)
		v.SetUint(i)
		return nil
	}

	if err != nil {
		return errors.New(fmt.Sprintf("value not a %s", kind.String()))
	}
	dec := gob.NewDecoder(bytes.NewBuffer(bs))
	return dec.Decode(value)
}
