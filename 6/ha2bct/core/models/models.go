package models

import (
	"ha2bct/core/models/computes"
	"ha2bct/core/models/dests"
	"ha2bct/core/models/jumps"
)

type HackAssemblerElement interface {
	Type() HackAssemblerElementType
}

type HackAssemblerElementType int

var elementTypeStrings = map[HackAssemblerElementType]string{
	Label:        "Label",
	AInstruction: "A-Instruction",
	CInstruction: "C-Instruction",
}

func (elementType HackAssemblerElementType) ToString() string {
	return elementTypeStrings[elementType]
}

const (
	Label HackAssemblerElementType = iota
	AInstruction
	CInstruction
)

// region Label

type HackAssemblerLabel struct {
	Value string
}

func (label HackAssemblerLabel) Type() HackAssemblerElementType {
	return Label
}

// endregion

// region AInstruction

type HackAssemblerAInstruction interface {
	Type() HackAssemblerElementType
	SubType() HackAssemblerAInstructionType
}

type HackAssemblerAInstructionType int

const (
	AlphabeticAInstruction HackAssemblerAInstructionType = iota
	NumericAInstruction
)

type HackAssemblerAlphabeticAInstruction struct {
	Value string
}

func (aInstruction HackAssemblerAlphabeticAInstruction) Type() HackAssemblerElementType {
	return AInstruction
}

func (aInstruction HackAssemblerAlphabeticAInstruction) SubType() HackAssemblerAInstructionType {
	return AlphabeticAInstruction
}

type HackAssemblerNumericAInstruction struct {
	Value int
}

func (aInstruction HackAssemblerNumericAInstruction) Type() HackAssemblerElementType {
	return AInstruction
}

func (aInstruction HackAssemblerNumericAInstruction) SubType() HackAssemblerAInstructionType {
	return NumericAInstruction
}

// endregion

// region CInstruction

type HackAssemblerCInstruction struct {
	Dest    dests.Dest
	Compute computes.Compute
	Jump    jumps.Jump
}

func (cInstruction HackAssemblerCInstruction) Type() HackAssemblerElementType {
	return CInstruction
}

// endregion
