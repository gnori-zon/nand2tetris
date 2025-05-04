package usecase

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

var fileAssemblerCode = []string{
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

var fileExpectedBinaryCode = []string{
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

func TestHackAssemblerTo16BitFileTranslator_TranslateAllTranslateAll(t *testing.T) {
	translator := NewHackAssemblerTo16BitFileTranslator()
	file, err := createTempFileWithAssembler()
	if err != nil {
		t.Errorf("Bad create temp file: %s", err)
	}
	defer file.Close()
	outputFileName, translateErr := translator.TranslateAll(file.Name())

	errs := make([]string, 0)
	if translateErr != nil {
		errs = append(errs, fmt.Sprintf("Error is not expected but has %d", translateErr))
	}
	fileWithBinaryCodes, err := os.Open(outputFileName)
	defer func() {
		fileWithBinaryCodes.Close()
		os.Remove(outputFileName)
	}()
	if err != nil {
		t.Errorf("Bad open output file: %s", err)
	}
	scanner := bufio.NewScanner(fileWithBinaryCodes)
	for i := 0; i < len(fileExpectedBinaryCode); i++ {
		if !scanner.Scan() {
			t.Errorf("Bad read at position: %d in file", i)
		}
		translatedByteCode := scanner.Text()
		if fileExpectedBinaryCode[i] != translatedByteCode {
			errs = append(errs, fmt.Sprintf("Expected [%d] binary code: %s But actual: %s", i, fileExpectedBinaryCode[i], translatedByteCode))
		}
	}
	for scanner.Scan() {
		nextRow := scanner.Text()
		if nextRow != "" {
			t.Errorf("Expected empty text but got %s", nextRow)
		}
	}
	if len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func createTempFileWithAssembler() (*os.File, error) {
	file, err := os.CreateTemp("", "")
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(file)
	for _, assemblerRow := range fileAssemblerCode {
		_, err := writer.WriteString(assemblerRow + "\n")
		if err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return file, nil
}
