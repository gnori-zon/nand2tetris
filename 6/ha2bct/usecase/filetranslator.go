package usecase

import (
	"bufio"
	"os"
	"path/filepath"
)

type HackAssemblerTo16BitFileTranslator struct {
	iterativeTranslator HackAssemblerTo16BitIterativeTranslator
}

func NewHackAssemblerTo16BitFileTranslator() *HackAssemblerTo16BitFileTranslator {
	return &HackAssemblerTo16BitFileTranslator{iterativeTranslator: *NewNewHackAssemblerTo16BitIterativeTranslator()}
}

func (translator HackAssemblerTo16BitFileTranslator) TranslateAll(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	rowReaderOnInitialize := buildFileRowReader(file)
	if err := translator.iterativeTranslator.InitializeAll(rowReaderOnInitialize); err != nil {
		return "", err
	}

	if err := resetRowReader(file); err != nil {
		return "", err
	}
	outputFile, err := createOutputFile(filename)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	translatedRowsCollector := func(translated string) {
		writer.WriteString(translated + "\n")
	}
	rowReaderOnTranslate := buildFileRowReader(file)
	if err := translator.iterativeTranslator.TranslateAll(rowReaderOnTranslate, translatedRowsCollector); err != nil {
		return "", err
	}
	if err := writer.Flush(); err != nil {
		return "", err
	}
	return outputFile.Name(), nil
}

func resetRowReader(file *os.File) error {
	_, err := file.Seek(0, 0)
	return err
}

func buildFileRowReader(file *os.File) func() (string, bool) {
	scanner := bufio.NewScanner(file)
	return func() (string, bool) {
		okRead := scanner.Scan()
		if okRead {
			return scanner.Text(), true
		}
		return "", false
	}
}

func createOutputFile(sourceFilename string) (*os.File, error) {
	dir := filepath.Dir(sourceFilename)
	newName := withoutExtension(sourceFilename) + ".hack"
	return os.Create(filepath.Join(dir, newName))
}

func withoutExtension(filename string) string {
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	return base[0 : len(base)-len(ext)]
}
