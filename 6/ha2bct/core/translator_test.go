package core

import (
	"ha2bct/core/models"
	"ha2bct/core/models/computes"
	"ha2bct/core/models/dests"
	"ha2bct/core/models/jumps"
	"testing"
)

func TestTranslate(t *testing.T) {
	translator := NewHackAssemblerElementTo16BitTranslator(New16BitsSymbolsTable())

	t.Run("Translate when model is (c/a)-instruction should return translated", func(t *testing.T) {
		expectedTranslatedByElement := map[models.HackAssemblerElement]string{
			models.HackAssemblerAlphabeticAInstruction{Value: "some_a"}:                           "0000000000010000",
			models.HackAssemblerNumericAInstruction{Value: 12}:                                    "0000000000001100",
			models.HackAssemblerCInstruction{Dest: dests.D, Compute: computes.A, Jump: jumps.JLE}: "1110110000010110",
		}
		for element, expectedTranslated := range expectedTranslatedByElement {
			translated, err := translator.Translate(element)
			if err != nil {
				t.Errorf("Translate expected return nil error, but got: %v", err)
			}
			if expectedTranslated != translated {
				t.Errorf("Translate expected return %q but got: %q", expectedTranslated, translated)
			}
		}
	})

	t.Run("Translate when model is label should return error", func(t *testing.T) {
		translated, err := translator.Translate(models.HackAssemblerLabel{Value: "some_label"})
		if err == nil {
			t.Error("Translate expected return error, but got nil")
		}
		if translated != "" {
			t.Errorf("Translate expected return empty but got: %q", translated)
		}
	})

	t.Run("Translate when model is invalid a-instruction should return error", func(t *testing.T) {
		invalidAInstruction := []models.HackAssemblerAInstruction{
			models.HackAssemblerNumericAInstruction{Value: -1},
			models.HackAssemblerNumericAInstruction{Value: Max15bitValue + 1},
		}
		for _, element := range invalidAInstruction {
			translated, err := translator.Translate(element)
			if err == nil {
				t.Error("Translate expected return error, but got nil")
			}
			if translated != "" {
				t.Errorf("Translate expected return empty but got: %q", translated)
			}
		}
	})
}

func TestNeedTranslate(t *testing.T) {
	translator := NewHackAssemblerElementTo16BitTranslator(New16BitsSymbolsTable())

	t.Run("NeedTranslate should return true if element is (c/a)-instruction", func(t *testing.T) {
		needTranslateElements := []models.HackAssemblerElement{models.HackAssemblerAlphabeticAInstruction{}, models.HackAssemblerNumericAInstruction{}, models.HackAssemblerCInstruction{}}
		for _, needTranslateElement := range needTranslateElements {
			isNeedTranslate := translator.NeedTranslate(needTranslateElement)
			if !isNeedTranslate {
				t.Errorf("NeedTranslate expected return 'true' but return 'false'")
			}
		}
	})

	t.Run("NeedTranslate should return false if element is label", func(t *testing.T) {
		notNeedTranslateElements := []models.HackAssemblerElement{models.HackAssemblerLabel{}}
		for _, notNeedTranslateElement := range notNeedTranslateElements {
			isNeedTranslate := translator.NeedTranslate(notNeedTranslateElement)
			if isNeedTranslate {
				t.Errorf("NeedTranslate expected return 'false' but return 'true'")
			}
		}
	})
}
