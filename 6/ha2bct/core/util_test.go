package core

import "testing"

func TestValidate15BitAddress(t *testing.T) {
	t.Run("Validate15BitAddress should return nil error if address is valid", func(t *testing.T) {
		for address := 0; address <= Max15bitValue; address++ {
			err := Validate15BitAddress(address)
			if err != nil {
				t.Errorf("Validate15BitAddress expected returned nil error, but return %v", err)
			}
		}
	})

	t.Run("Validate15BitAddress should return error if address is invalid", func(t *testing.T) {
		for _, address := range []int{-1, Max15bitValue + 1} {
			err := Validate15BitAddress(address)
			if err == nil {
				t.Errorf("Validate15BitAddress expected returned error, but return nil")
			}
		}
	})
}
