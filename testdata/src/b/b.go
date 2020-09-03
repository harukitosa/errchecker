package b

import (
	"errors"
)

// 無名関数の場合も指摘
func voidfunc() (int, error) {
	s := func(src string) (int, error) { return 2, errors.New("error") }
	return s("hogehoge")
}

func voidfunc2() error {
	s := func(src string) error { return nil } // want "It returns nil in all the places where it should return error"
	return s("hogehoge")
}
