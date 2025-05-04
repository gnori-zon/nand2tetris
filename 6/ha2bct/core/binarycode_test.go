package core

import (
	"ha2bct/core/models/computes"
	"ha2bct/core/models/dests"
	"ha2bct/core/models/jumps"
	"maps"
	"testing"
)

var expectedBinaryCodesForValues = map[int]string{
	0:             "000000000000000",
	200:           "000000011001000",
	Max15bitValue: "111111111111111",
}

func TestEncodeAInstruction(t *testing.T) {
	binaryCodeProvider := NewBinaryCode16BitsProvider()
	t.Run("EncodeAInstruction when address is valid should return binary code representation", func(t *testing.T) {
		for address := range maps.Keys(expectedBinaryCodesForValues) {
			binaryCode, err := binaryCodeProvider.EncodeAInstruction(address)
			expectedBinaryCode := "0" + expectedBinaryCodesForValues[address]
			if expectedBinaryCode != binaryCode {
				t.Errorf("Expected binary code to be %s, but got %s", expectedBinaryCode, binaryCode)
			}
			if err != nil {
				t.Errorf("Expected error to be nil, but got %s", err)
			}
		}
	})

	t.Run("EncodeAInstruction when address is not valid should return error and empty binary code", func(t *testing.T) {
		for _, address := range []int{-1, Max15bitValue + 1} {
			binaryCode, err := binaryCodeProvider.EncodeAInstruction(address)
			if binaryCode != "" {
				t.Errorf("Expected empty binary code, but got %s", binaryCode)
			}
			if err == nil {
				t.Errorf("Expected error to not be nil, but got nil")
			}
		}
	})
}

func TestEncodeCInstruction(t *testing.T) {
	binaryCodeProvider := NewBinaryCode16BitsProvider()

	expectedFirstBinaryCode := "1110000111000010"
	firstBinaryCode := binaryCodeProvider.EncodeCInstruction(dests.Empty, computes.AMinusD, jumps.JEQ)
	if firstBinaryCode != expectedFirstBinaryCode {
		t.Errorf("Expected binary code %q, but got %q", expectedFirstBinaryCode, firstBinaryCode)
	}

	expectedSecondBinaryCode := "1111010011101000"
	secondBinaryCode := binaryCodeProvider.EncodeCInstruction(dests.AM, computes.DMinusM, jumps.Empty)
	if secondBinaryCode != expectedSecondBinaryCode {
		t.Errorf("Expected binary code %q, but got %q", expectedSecondBinaryCode, secondBinaryCode)
	}
}
