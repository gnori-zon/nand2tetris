package dests

import (
	"maps"
	"testing"
)

var expectedDestByString = map[string]Dest{
	"":    Empty,
	"M":   M,
	"D":   D,
	"MD":  MD,
	"A":   A,
	"AM":  AM,
	"AD":  AD,
	"AMD": AMD,
}

func TestNewDest(t *testing.T) {

	t.Run("NewDest from string when string is valid should return ok and dest instance", func(t *testing.T) {
		for computeDest := range maps.Keys(expectedDestByString) {
			expectedCompute := expectedDestByString[computeDest]
			compute, ok := NewDest(computeDest)
			if !ok {
				t.Errorf("NewDest returned wrong ok, expected 'true' but return 'false' for %q", computeDest)
			}
			if compute != expectedCompute {
				t.Errorf("NewDest returned wrong dest instance, expected: '%d' but got '%d'", expectedCompute, compute)
			}
		}
	})

	t.Run("NewDest from string when string is valid should return not ok and -1", func(t *testing.T) {
		unknownDestString := "unknown-dest"
		dest, ok := NewDest(unknownDestString)
		if ok {
			t.Errorf("NewDest returned wrong ok, expected 'false' but return 'true' for %q", unknownDestString)
		}
		if dest != -1 {
			t.Errorf("NewDest returned wrong dest instance, expected: '-1' but got '%d'", dest)
		}
	})
}
