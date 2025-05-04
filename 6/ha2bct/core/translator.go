package core

import (
	"errors"
	"ha2bct/core/models"
)

type HackAssemblerElementToBinaryCodeTranslator interface {
	Translate(element models.HackAssemblerElement) (string, error)
	NeedTranslate(element models.HackAssemblerElement) bool
}

type HackAssemblerElementTo16BitTranslator struct {
	symbolsTable       *SymbolsTable
	binaryCodeProvider BinaryCodeProvider
	typedTranslators   map[models.HackAssemblerElementType]TypedHackAssemblerElement16bitTranslator
}

func NewHackAssemblerElementTo16BitTranslator(symbolsTable *SymbolsTable) *HackAssemblerElementTo16BitTranslator {
	binaryCodeProvider := NewBinaryCode16BitsProvider()
	typedTranslators := map[models.HackAssemblerElementType]TypedHackAssemblerElement16bitTranslator{
		models.AInstruction: HackAssemblerAInstructionTo16BitTranslator{symbolsTable, binaryCodeProvider},
		models.CInstruction: HackAssemblerCInstructionTo16BitTranslator{binaryCodeProvider},
	}
	return &HackAssemblerElementTo16BitTranslator{
		symbolsTable:       symbolsTable,
		binaryCodeProvider: NewBinaryCode16BitsProvider(),
		typedTranslators:   typedTranslators,
	}
}

func (translator *HackAssemblerElementTo16BitTranslator) Translate(element models.HackAssemblerElement) (string, error) {
	if typedTranslator, ok := translator.typedTranslators[element.Type()]; ok {
		return typedTranslator.Translate(element)
	}
	return "", errors.New("unsupported element type" + element.Type().ToString())
}

func (translator *HackAssemblerElementTo16BitTranslator) NeedTranslate(element models.HackAssemblerElement) bool {
	return element.Type() != models.Label
}

type TypedHackAssemblerElement16bitTranslator interface {
	Translate(element models.HackAssemblerElement) (string, error)
	SupportsType() models.HackAssemblerElementType
}

// region AInstructionTranslator

type HackAssemblerAInstructionTo16BitTranslator struct {
	symbolsTable       *SymbolsTable
	binaryCodeProvider BinaryCodeProvider
}

func (translator HackAssemblerAInstructionTo16BitTranslator) Translate(element models.HackAssemblerElement) (string, error) {
	if aInstruction, ok := element.(models.HackAssemblerAInstruction); ok {
		address, err := translator.resolveAddress(aInstruction)
		if err != nil {
			return "", err
		}
		if errValidate := Validate15BitAddress(address); errValidate != nil {
			return "", errValidate
		}
		return translator.binaryCodeProvider.EncodeAInstruction(address)
	}
	return "", errors.New("element is not a-instruction, type: " + element.Type().ToString())
}

func (translator HackAssemblerAInstructionTo16BitTranslator) SupportsType() models.HackAssemblerElementType {
	return models.AInstruction
}

func (translator HackAssemblerAInstructionTo16BitTranslator) resolveAddress(instruction models.HackAssemblerAInstruction) (int, error) {
	if instruction.SubType() == models.NumericAInstruction {
		if numericAInstruction, ok := instruction.(models.HackAssemblerNumericAInstruction); ok {
			return numericAInstruction.Value, nil
		}
	}
	if instruction.SubType() == models.AlphabeticAInstruction {
		if alphabeticAInstruction, ok := instruction.(models.HackAssemblerAlphabeticAInstruction); ok {
			if existAddress, ok := translator.symbolsTable.Get(alphabeticAInstruction.Value); ok {
				return existAddress, nil
			}
			return translator.symbolsTable.AddGenerated(alphabeticAInstruction.Value)
		}
	}
	return -1, errors.New("element is not supported a-instruction, type: " + instruction.Type().ToString())
}

// endregion

// region CInstructionTranslator

type HackAssemblerCInstructionTo16BitTranslator struct {
	binaryCodeProvider BinaryCodeProvider
}

func (translator HackAssemblerCInstructionTo16BitTranslator) Translate(element models.HackAssemblerElement) (string, error) {
	if cInstruction, ok := element.(models.HackAssemblerCInstruction); ok {
		return translator.binaryCodeProvider.EncodeCInstruction(cInstruction.Dest, cInstruction.Compute, cInstruction.Jump), nil
	}
	return "", errors.New("element is not c-instruction, type: " + element.Type().ToString())
}

func (translator HackAssemblerCInstructionTo16BitTranslator) SupportsType() models.HackAssemblerElementType {
	return models.CInstruction
}

// endregion
