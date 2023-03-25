package a

func funcCA() (*T, error) {
	var err error
	if err != nil {
		return nil, nil // want "return err"
	}
	return nil, nil
}

func funcCB() (*T, error) {
	t, err := funcB()
	if err != nil {
		//lint:ignore retnonnilerr ignore
		return nil, nil
	}
	return t, nil
}

func funcCC() (*T, error) {
	t, err := funcB()
	if err != nil {
		return nil, err
	}
	t, err = funcB()
	if err != nil {
		//lint:ignore retnonnilerr ignore
		return nil, nil
	}
	return t, nil
}
