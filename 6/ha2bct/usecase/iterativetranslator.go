package usecase

import (
	"ha2bct/core"
)

type HackAssemblerTo16BitIterativeTranslator struct {
	elementParser      core.HackAssemblerElement16BitParser
	elementInitializer core.HackAssemblerElement16BitInitializer
	elementTranslator  core.HackAssemblerElementTo16BitTranslator
}

func NewNewHackAssemblerTo16BitIterativeTranslator() *HackAssemblerTo16BitIterativeTranslator {
	ptrSymbolsTable := core.New16BitsSymbolsTable()
	return &HackAssemblerTo16BitIterativeTranslator{
		elementParser:      *core.NewHackAssembler16BitParser(),
		elementInitializer: *core.NewHackAssemblerElement16BitInitializer(ptrSymbolsTable),
		elementTranslator:  *core.NewHackAssemblerElementTo16BitTranslator(ptrSymbolsTable),
	}
}

type RowReader func() (row string, hasNext bool)
type TranslatedElementCollector func(translated string)

func (translator *HackAssemblerTo16BitIterativeTranslator) InitializeAll(reader RowReader) error {
	numberElement := 0
	row, okRead := reader()
	for okRead {
		if element, ok := translator.elementParser.Parse(row); ok {
			if err := translator.elementInitializer.Initialize(numberElement, element); err != nil {
				return err
			}
			numberElement++
		}
		row, okRead = reader()
	}
	return nil
}

func (translator *HackAssemblerTo16BitIterativeTranslator) TranslateAll(
	reader RowReader,
	collector TranslatedElementCollector,
) error {
	row, okRead := reader()
	for okRead {
		if element, ok := translator.elementParser.Parse(row); ok {
			if translator.elementTranslator.NeedTranslate(element) {
				translated, err := translator.elementTranslator.Translate(element)
				if err != nil {
					return err
				}
				collector(translated)
			}
		}
		row, okRead = reader()
	}
	return nil
}
