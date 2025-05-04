package models

import "testing"

func TestType(t *testing.T) {
	assertType(t, HackAssemblerLabel{}, Label)
	assertType(t, HackAssemblerAlphabeticAInstruction{}, AInstruction)
	assertType(t, HackAssemblerNumericAInstruction{}, AInstruction)
	assertType(t, HackAssemblerCInstruction{}, CInstruction)
}

func TestAInstructionSubType(t *testing.T) {
	assertAInstructionSubType(t, HackAssemblerAlphabeticAInstruction{}, AlphabeticAInstruction)
	assertAInstructionSubType(t, HackAssemblerNumericAInstruction{}, NumericAInstruction)
}

func assertType(t *testing.T, element HackAssemblerElement, expectedType HackAssemblerElementType) {
	if element.Type() != expectedType {
		t.Errorf("Expected type '%d' but got '%d'", expectedType, element.Type())
	}
}

func assertAInstructionSubType(t *testing.T, aInstruction HackAssemblerAInstruction, expectedSubType HackAssemblerAInstructionType) {
	if aInstruction.SubType() != expectedSubType {
		t.Errorf("Expected subType '%d' but got '%d'", expectedSubType, aInstruction.SubType())
	}
}
