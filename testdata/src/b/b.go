package b

import (
	"errors"
	"log"
)

// 返り値が存在しない関数
func noreturn() {
	log.Print("hello")
}

func sample2() (int, error) { // want "It returns nil in all the places where it should return error"
	return 0, nil
}

// 無名関数の場合も指摘
func voidfunc() error {
	s := func(src string) error { return errors.New("error") }
	return s("hogehoge")
}

func voidfunc2() error {
	s := func(src string) error { return nil } // want "It returns nil in all the places where it should return error"
	return s("hogehoge")
}
