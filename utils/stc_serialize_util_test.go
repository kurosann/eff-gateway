package utils

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	s := struct {
		Name string
		Id   string
	}{
		Name: "test",
		Id:   "id1",
	}
	encode, err := Encode(s)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(s)
	fmt.Println(reflect.TypeOf(s).String())

	var value struct {
		Name string
		Id   string
	}
	err = Decode(encode, &value)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(value)
	fmt.Println(reflect.TypeOf(value).String())
}
