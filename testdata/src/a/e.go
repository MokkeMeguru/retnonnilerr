package a

import "log"

func funcEA() bool {
	if _, err := funcB(); err != nil {
		return false // want "return err"
	}
	return true
}

func funcEB() error {
	if _, err := funcB(); err != nil {
		log.Printf("error is happend: %w", err)
		return nil // want "return err"
	}
	return nil
}
