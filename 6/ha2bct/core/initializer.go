package core

import (
	"ha2bct/core/models"
)

type HackAssemblerElementInitializer interface {
	Initialize(number int, element models.HackAssemblerElement) error
}

type HackAssemblerElement16BitInitializer struct {
	symbolsTable         *SymbolsTable
	countProcessedLabels int
}

func NewHackAssemblerElement16BitInitializer(symbolsTable *SymbolsTable) *HackAssemblerElement16BitInitializer {
	return &HackAssemblerElement16BitInitializer{symbolsTable: symbolsTable, countProcessedLabels: 0}
}

func (initializer *HackAssemblerElement16BitInitializer) Initialize(
	number int,
	element models.HackAssemblerElement,
) error {
	if element.Type() == models.Label {
		if label, ok := element.(models.HackAssemblerLabel); ok {
			labelAddress := number - initializer.countProcessedLabels
			if err := Validate15BitAddress(labelAddress); err != nil {
				return err
			}
			if err := initializer.symbolsTable.Add(label.Value, labelAddress); err != nil {
				return err
			}
			initializer.countProcessedLabels++
		}
	}
	return nil
}
