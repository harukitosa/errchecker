package b

import "log"

// 返り値が存在しない関数
func noreturn() {
	log.Print("hello")
}

func sample2() (int, error) { // want "It returns nil in all the places where it should return error"
	return 0, nil
}

func voidfunc() {}
