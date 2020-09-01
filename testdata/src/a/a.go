package a

import "errors"

// testcase 1
func sample1() (string, error) {
	return "helloworld", nil // want "It returns nil in all the places where it should return error"
}

func sample2() (int, error) {
	return 0, nil // want "It returns nil in all the places where it should return error"
}

// testcase 2
func sample3() (int, error) {
	return 0, errors.New("error")
}
