package core

import (
	"ha2bct/core/models"
	"strconv"
	"testing"
)

func TestInitialize(t *testing.T) {
	t.Run("Initialize should update symbolsTable for initializing label", func(t *testing.T) {
		initializer := NewHackAssemblerElement16BitInitializer(New16BitsSymbolsTable())
		labels := map[string]int{
			"firstLabel":  1,
			"secondLabel": 3,
			"thirdLabel":  5,
		}
		countElementsIntSymbolsTableBeforeInitialize := len(initializer.symbolsTable.data)
		countInitializedLabels := 0
		for label, address := range labels {
			err := initializer.Initialize(address, models.HackAssemblerLabel{Value: label})
			if err != nil {
				t.Errorf("Expected no error but got err: %s", err)
			}
			if savedAddress, ok := initializer.symbolsTable.Get(label); ok {
				expectedAddress := address - countInitializedLabels
				if expectedAddress != savedAddress {
					t.Errorf("Expected saved address %v but got %v", expectedAddress, savedAddress)
				}
			} else {
				t.Errorf("Expected saved address %v but got not ok for Get address from symbols table%v", address, savedAddress)
			}
			countInitializedLabels++
		}
		countElementsIntSymbolsTableAfterInitialize := len(initializer.symbolsTable.data)
		countAddedElements := countElementsIntSymbolsTableAfterInitialize - countElementsIntSymbolsTableBeforeInitialize
		if countAddedElements != len(labels) {
			t.Errorf("Expected added %d elements, but added %d", len(labels), countAddedElements)
		}

	})

	t.Run("Initialize should skip not labels", func(t *testing.T) {
		initializer := NewHackAssemblerElement16BitInitializer(New16BitsSymbolsTable())
		countElementsIntSymbolsTableBeforeInitialize := len(initializer.symbolsTable.data)

		err1 := initializer.Initialize(1, models.HackAssemblerCInstruction{})
		err2 := initializer.Initialize(2, models.HackAssemblerAlphabeticAInstruction{})
		err3 := initializer.Initialize(2, models.HackAssemblerNumericAInstruction{})

		countElementsIntSymbolsTableAfterInitialize := len(initializer.symbolsTable.data)
		if countElementsIntSymbolsTableBeforeInitialize != countElementsIntSymbolsTableAfterInitialize {
			t.Errorf("Expected skip not labels but has affect on symbolsTable, beforeInitializeSize: '%d' and afterInitializeSize: '%d'",
				countElementsIntSymbolsTableBeforeInitialize, countElementsIntSymbolsTableAfterInitialize)
		}
		if err1 != nil || err2 != nil || err3 != nil {
			t.Errorf("Expected no error but got err1: %s, err2: %s, err3: %s", err1, err2, err3)
		}
	})

	t.Run("Initialize should return error if label has incorrect address", func(t *testing.T) {
		initializer := NewHackAssemblerElement16BitInitializer(New16BitsSymbolsTable())
		for _, address := range []int{-1, Max15bitValue + 1} {
			err := initializer.Initialize(address, models.HackAssemblerLabel{Value: "some" + strconv.Itoa(address)})
			if err == nil {
				t.Errorf("Expected error but got nil for address: %d", address)
			}
		}
	})
}
