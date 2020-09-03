package b

import (
	"errors"
	"fmt"
)

// 無名関数の場合も指摘
func anonfunc() (int, error) {
	s := func(src string) (int, error) { return 2, errors.New("error") }
	return s("hogehoge")
}

func anonfunc2() error {
	s := func(src string) error { return nil } // want "It returns nil in all the places where it should return error"
	return s("hogehoge")
}

func anonfunc3() error { // want "It returns nil in all the places where it should return error"
	s := func(src string) error { return nil } // want "It returns nil in all the places where it should return error"
	fmt.Println(s("hogehoge"))
	return nil
}
