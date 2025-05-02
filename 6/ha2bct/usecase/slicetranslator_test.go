package usecase

import (
	"fmt"
	"testing"
)

var assemblerCode = []string{
	" // Computes R2 = max(R0, R1)  (R0,R1,R2 refer to RAM[0],RAM[1],RAM[2])   	",
	" // Usage: Before executing, put two values in R0 and R1.                 	",
	"                                                                          	",
	" // D = R0 - R1														   	",
	" @R0																		",
	" D=	M																	",
	" @R1																		",
	" D=D-M																		",
	" // If (D > 0) goto ITSR0													",
	" @ITSR0																	",
	" D;JGT																		",
	" // Its R1																	",
	" @R1																		",
	" D=M																		",
	" @OUTPUT_D																	",
	" 0;JMP																		",
	" (ITSR0)																	",
	" @R0																		",
	" D=M																		",
	" (OUTPUT_D		)															",
	" @R2																		",
	" M=D																		",
	" (END)																		",
	" @END																		",
	" 0;JMP																		",
}

var expectedBinaryCode = []string{
	"0000000000000000",
	"1111110000010000",
	"0000000000000001",
	"1111010011010000",
	"0000000000001010",
	"1110001100000001",
	"0000000000000001",
	"1111110000010000",
	"0000000000001100",
	"1110101010000111",
	"0000000000000000",
	"1111110000010000",
	"0000000000000010",
	"1110001100001000",
	"0000000000001110",
	"1110101010000111",
}

func TestTranslateAll(t *testing.T) {
	translator := NewHackAssemblerTo16BitSliceTranslator()
	binaryCode, translateErr := translator.TranslateAll(assemblerCode)

	needCheckEach := true
	errs := make([]string, 0)
	if translateErr != nil {
		errs = append(errs, fmt.Sprintf("Error is not expected, but has %d", translateErr))
	}
	if len(expectedBinaryCode) != len(binaryCode) {
		needCheckEach = false
		errs = append(errs, fmt.Sprintf("Expected binary code length: %d But actual: %d", len(expectedBinaryCode), len(binaryCode)))
	}
	for i := 0; needCheckEach && i < len(binaryCode); i++ {
		if expectedBinaryCode[i] != binaryCode[i] {
			errs = append(errs, fmt.Sprintf("Expected [%d] binary code: %s But actual: %s", i, expectedBinaryCode[i], binaryCode[i]))
		}
	}
	if len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}
