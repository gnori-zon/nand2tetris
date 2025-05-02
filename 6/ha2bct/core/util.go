package core

import (
	"errors"
	"fmt"
)

var Max15bitValue = 32768

func Validate15BitAddress(address int) error {
	if address >= Max15bitValue {
		return errors.New(fmt.Sprintf("address: %d exceed max address%d", address, Max15bitValue))
	}
	return nil
}
