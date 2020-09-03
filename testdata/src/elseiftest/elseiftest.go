package elseiftest

import "errors"

func elseiffunc() error {
	if false {
		return nil
	} else if true {
		return nil
	} else if true {
		return nil
	} else {
		return errors.New("error")
	}
}
