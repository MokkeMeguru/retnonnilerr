package a

import "errors"

func funcH() error {
	if _, err := funcB(); err != nil {
		return errors.New("error is %w")
	}
	return nil
}
