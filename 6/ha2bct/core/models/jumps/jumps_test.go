package jumps

import (
	"maps"
	"testing"
)

var expectedJumpByString = map[string]Jump{
	"":    Empty,
	"JMP": JMP,
	"JEQ": JEQ,
	"JNE": JNE,
	"JGT": JGT,
	"JGE": JGE,
	"JLT": JLT,
	"JLE": JLE,
}

func TestNewJump(t *testing.T) {

	t.Run("NewJump from string when string is valid should return ok and jump instance", func(t *testing.T) {
		for computeDest := range maps.Keys(expectedJumpByString) {
			expectedCompute := expectedJumpByString[computeDest]
			compute, ok := NewJump(computeDest)
			if !ok {
				t.Errorf("NewJump returned wrong ok, expected 'true' but return 'false' for %q", computeDest)
			}
			if compute != expectedCompute {
				t.Errorf("NewJump returned wrong jump instance, expected: '%d' but got '%d'", expectedCompute, compute)
			}
		}
	})

	t.Run("NewJump from string when string is valid should return not ok and -1", func(t *testing.T) {
		unknownJumpString := "unknown-jump"
		jump, ok := NewJump(unknownJumpString)
		if ok {
			t.Errorf("NewJump returned wrong ok, expected 'false' but return 'true' for %q", unknownJumpString)
		}
		if jump != -1 {
			t.Errorf("NewJump returned wrong jump instance, expected: '-1' but got '%d'", jump)
		}
	})
}
