package a

type T struct {
	I int
}

func funcA() (*T, error) {
	return nil, nil
}

func funcB() (*T, error) {
	var err error
	if err != nil {
		return nil, nil // want "return err"
	}
	return nil, nil
}

func funcC() (*T, error) {
	t, err := funcB()
	if err != nil {
		return nil, nil // want "return err"
	}
	return t, nil
}

func funcD() (*T, error) {
	if _, err := funcB(); err != nil {
		return nil, nil // want "return err"
	}
	return nil, nil
}

func funcE() (*T, error) {
	if x := 1 + 2; x == 3 {
		return &T{I: x}, nil
	}
	return nil, nil
}

func funcF() (*T, string, error) {
	if t, err := funcB(); err != nil {
		return t, "", nil // want "return err"
	}
	return nil, "", nil
}

// return err
func funcG() error {
	if _, err := funcB(); err != nil {
		return nil // want "return err"
	}
	return nil
}
