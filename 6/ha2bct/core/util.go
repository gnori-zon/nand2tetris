package core

import (
	"errors"
	"fmt"
)

var Max15bitValue = 32767

func Validate15BitAddress(address int) error {
	if address < 0 {
		return errors.New("address must be positive")
	}
	if address > Max15bitValue {
		return errors.New(fmt.Sprintf("address: %d exceed max address: %d", address, Max15bitValue))
	}
	return nil
}
