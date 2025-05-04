package computes

import (
	"maps"
	"testing"
)

var expectedComputeByString = map[string]Compute{
	"0":   Zero,
	"1":   One,
	"-1":  NegativeOne,
	"D":   D,
	"A":   A,
	"!D":  NotD,
	"!A":  NotA,
	"-D":  NegativeD,
	"-A":  NegativeA,
	"D+1": DPlusOne,
	"A+1": APlusOne,
	"D-1": DMinusOne,
	"A-1": AMinusOne,
	"D+A": DPlusA,
	"D-A": DMinusA,
	"A-D": AMinusD,
	"D&A": DAndA,
	"D|A": DOrA,
	"M":   M,
	"!M":  NotM,
	"-M":  NegativeM,
	"M+1": MPlusOne,
	"M-1": MMinusOne,
	"D+M": DPlusM,
	"D-M": DMinusM,
	"M-D": MMinusD,
	"D&M": DAndM,
	"D|M": DOrM,
}

func TestNewCompute(t *testing.T) {

	t.Run("NewCompute from string when string is valid should return ok and compute instance", func(t *testing.T) {
		for computeString := range maps.Keys(expectedComputeByString) {
			expectedCompute := expectedComputeByString[computeString]
			compute, ok := NewCompute(computeString)
			if !ok {
				t.Errorf("NewCompute returned wrong ok, expected 'true' but return 'false' for %q", computeString)
			}
			if compute != expectedCompute {
				t.Errorf("NewCompute returned wrong compute instance, expected: '%d' but got '%d'", expectedCompute, compute)
			}
		}
	})

	t.Run("NewCompute from string when string is valid should return not ok and -1", func(t *testing.T) {
		unknownComputeString := "unknown-compute"
		compute, ok := NewCompute(unknownComputeString)
		if ok {
			t.Errorf("NewCompute returned wrong ok, expected 'false' but return 'true' for %q", unknownComputeString)
		}
		if compute != -1 {
			t.Errorf("NewCompute returned wrong compute instance, expected: '-1' but got '%d'", compute)
		}
	})
}
