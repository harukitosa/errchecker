package sample_test1

import "errors"

func sample1() error { // want "It returns nil in all the places where it should return error"
	if false {
		s := 0
		s++
		if s >= 1 {
			s++
		} else {
			s++
			for {
				return nil
			}
		}
	}
	return nil
}

func sample2() error {
	if false {
		s := 0
		s++
		if s >= 1 {
			s++
		} else {
			s++
			for {
				return errors.New("error")
			}
		}
	}
	return errors.New("error")
}
