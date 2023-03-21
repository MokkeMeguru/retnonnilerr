package a

func funcDA() error {
	x := 1
	if x == 1 {
		if _, err := funcB(); err != nil {
			return nil // want "return err"
		}
	}
	return nil
}
