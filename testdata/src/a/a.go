package a

func sample1() (string, error) {
	return "helloworld", nil // want "It returns nil in all the places where it should return error"
}

func sample2() (int, error) {
	return 0, nil // want "It returns nil in all the places where it should return error"
}
