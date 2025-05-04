package core

import (
	"ha2bct/core/models"
	"ha2bct/core/models/computes"
	"ha2bct/core/models/dests"
	"ha2bct/core/models/jumps"
	"strconv"
	"testing"
)

func TestParse(t *testing.T) {

	parser := NewHackAssembler16BitParser()

	t.Run("Parse when row is valid should return correct element and true for ok", func(t *testing.T) {

		expectedAlphabeticAInstructionValue := "asd$asd.while"
		alphabeticAInstruction, okParsedAlphabeticAInstruction := parser.Parse("@" + expectedAlphabeticAInstructionValue)
		assertTrueOk(t, okParsedAlphabeticAInstruction)
		assertElementType(t, alphabeticAInstruction, models.AInstruction)
		assertAInstructionSubType(t, alphabeticAInstruction.(models.HackAssemblerAInstruction), models.AlphabeticAInstruction)
		alphabeticAInstructionValue := alphabeticAInstruction.(models.HackAssemblerAlphabeticAInstruction).Value
		if expectedAlphabeticAInstructionValue != alphabeticAInstructionValue {
			t.Errorf("Incorrect alphabetic a-instruction value expected: %s but got: %s", expectedAlphabeticAInstructionValue, alphabeticAInstructionValue)
		}

		expectedNumericAInstructionValue := 123123
		numericAInstruction, okParsedNumericAInstruction := parser.Parse("@" + strconv.Itoa(expectedNumericAInstructionValue))
		assertTrueOk(t, okParsedNumericAInstruction)
		assertElementType(t, numericAInstruction, models.AInstruction)
		assertAInstructionSubType(t, numericAInstruction.(models.HackAssemblerAInstruction), models.NumericAInstruction)
		numericAInstructionValue := numericAInstruction.(models.HackAssemblerNumericAInstruction).Value
		if expectedNumericAInstructionValue != numericAInstructionValue {
			t.Errorf("Incorrect numeric a-instruction value expected: %d but got: %d", expectedNumericAInstructionValue, numericAInstructionValue)
		}

		expectedLabelValue := "asddasd.dasd$_asd-Asd"
		label, okParsedLabel := parser.Parse("(" + expectedLabelValue + ")")
		assertTrueOk(t, okParsedLabel)
		assertElementType(t, label, models.Label)
		labelValue := label.(models.HackAssemblerLabel).Value
		if expectedLabelValue != labelValue {
			t.Errorf("Incorrect label value expected: %s but got: %s", expectedLabelValue, labelValue)
		}

		expectedCInstruction := models.HackAssemblerCInstruction{Dest: dests.M, Compute: computes.APlusOne, Jump: jumps.JLE}
		cInstruction, okParsedCInstruction := parser.Parse("M=A+1;JLE")
		assertTrueOk(t, okParsedCInstruction)
		assertElementType(t, cInstruction, models.CInstruction)
		if expectedCInstruction != cInstruction {
			t.Errorf("Incorrect cInstruction value expected: %v but got: %v", expectedCInstruction, cInstruction)
		}
	})

	t.Run("Parse when row is comment or blank should return nil element and false for ok", func(t *testing.T) {
		rows := []string{
			"",
			"\t\t\t",
			"    \t   ",
			"// comment",
			"\t   // some-comment",
		}
		for _, row := range rows {
			element, ok := parser.Parse(row)
			assertFalseOk(t, ok)
			assertNilElement(t, element)
		}
	})

	t.Run("Parse when row is invalid should return nil element and false for ok", func(t *testing.T) {
		invalidRows := []string{
			"1+123",
			"@ASd_Asd  asdl;m,asd ;asdlasd",
			"(as.dasd",
		}
		for _, invalidRow := range invalidRows {
			element, ok := parser.Parse(invalidRow)
			assertFalseOk(t, ok)
			assertNilElement(t, element)
		}
	})
}

func assertTrueOk(t *testing.T, ok bool) {
	if !ok {
		t.Error("Expected ok to be true")
	}
}

func assertFalseOk(t *testing.T, ok bool) {
	if ok {
		t.Error("Expected ok to be false")
	}
}

func assertNilElement(t *testing.T, element models.HackAssemblerElement) {
	if element != nil {
		t.Errorf("Incorrect element value expected: %v but got: %v", nil, element)
	}
}

func assertElementType(t *testing.T, element models.HackAssemblerElement, expectedType models.HackAssemblerElementType) {
	if element.Type() != expectedType {
		t.Errorf("Expected element type to be %d, but got %d", expectedType, element.Type())
	}
}

func assertAInstructionSubType(t *testing.T, aInstruction models.HackAssemblerAInstruction, expectedSubType models.HackAssemblerAInstructionType) {
	if aInstruction.SubType() != expectedSubType {
		t.Errorf("Expected a-instruction subType to be %d, but got %d", expectedSubType, aInstruction.SubType())
	}
}
