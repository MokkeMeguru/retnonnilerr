package a

import "errors"

func funcBA() error {
	if _, err := funcB(); err != nil {
		return errors.New("error is %w")
	}
	return nil
}
