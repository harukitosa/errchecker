package a

import "errors"

// testcase 1
// return nil
func sample1() (string, error) { // want "It returns nil in all the places where it should return error"
	return "helloworld", nil
}
func sample2() (int, error) { // want "It returns nil in all the places where it should return error"
	return 0, nil
}

// testcase 2
// return error
func sample3() (int, error) {
	return 0, errors.New("error")
}

// testcase3
// nest return
func sample4() (int, error) { // want "It returns nil in all the places where it should return error"
	if false {
		return 0, nil
	}
	return 0, nil
}

func sample5() (int, error) {
	if false {
		return 0, errors.New("error")
	}
	return 0, nil
}

func sample6() (int, error) { // want "It returns nil in all the places where it should return error"
	if false {
		if false {
			return 0, nil
		}
	}
	return 0, nil
}

func sample7() (int, error) {
	if false {
		if false {
			return 0, errors.New("error")
		}
	}
	return 0, nil
}

// test case 8
// for statement
func sample8() (int, error) { // want "It returns nil in all the places where it should return error"
	for {
		if false {
			return 0, nil
		}
	}
}

func sample9() (string, error) {
	s := "hoge"
	for {
		if false {
			return s, errors.New("error")
		}
	}
}

func sample10() (int, error) {
	if false {
		return 0, nil
	}

	a := 0
	if false {
		a = 1
	} else {
		return 0, errors.New("error")
	}

	return a, nil
}
