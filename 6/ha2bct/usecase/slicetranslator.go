package usecase

type HackAssemblerTo16BitSliceTranslator struct {
	iterativeTranslator HackAssemblerTo16BitIterativeTranslator
}

func NewHackAssemblerTo16BitSliceTranslator() *HackAssemblerTo16BitSliceTranslator {
	return &HackAssemblerTo16BitSliceTranslator{iterativeTranslator: *NewNewHackAssemblerTo16BitIterativeTranslator()}
}

func (translator HackAssemblerTo16BitSliceTranslator) TranslateAll(rows []string) ([]string, error) {
	rowReaderOnInitialize := buildSliceRowReader(rows)
	if err := translator.iterativeTranslator.InitializeAll(rowReaderOnInitialize); err != nil {
		return nil, err
	}

	rowReaderOnTranslate := buildSliceRowReader(rows)
	translatedRows := make([]string, 0, len(rows))
	translatedRowsCollector := func(translated string) {
		translatedRows = append(translatedRows, translated)
	}
	if err := translator.iterativeTranslator.TranslateAll(rowReaderOnTranslate, translatedRowsCollector); err != nil {
		return nil, err
	}
	return translatedRows, nil
}

func buildSliceRowReader(rows []string) func() (string, bool) {
	translateCounter := 0
	return func() (string, bool) {
		if translateCounter < len(rows) {
			currentIndex := translateCounter
			translateCounter++
			return rows[currentIndex], true
		}
		return "", false
	}
}
