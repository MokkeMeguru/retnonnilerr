package a

import (
	"errors"
	"fmt"
)

type CustomErr struct {
	Err error
}

func (e *CustomErr) Error() string {
	return e.Err.Error()
}

func NewCustomErr(e error) *CustomErr {
	return &CustomErr{Err: e}
}

func funcBA() error {
	if _, err := funcB(); err != nil {
		return errors.New(fmt.Sprintf("error is %w", err))
	}
	return nil
}

func funcBB() error {
	if _, err := funcB(); err != nil {
		return &CustomErr{Err: err}
	}
	return nil
}

func funcBC() error {
	if _, err := funcB(); err != nil {
		return NewCustomErr(err)
	}
	return nil
}
