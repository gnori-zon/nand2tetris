package core

import (
	"strconv"
	"testing"
)

func TestSymbolsTable_AddGenerated(t *testing.T) {
	symbolsTable := New16BitsSymbolsTable()

	t.Run("AddGenerated should add symbols and return incremented address from 16 to max15BitValue", func(t *testing.T) {
		for i := 16; i <= Max15bitValue; i++ {
			address, addGeneratedErr := symbolsTable.AddGenerated("some-key" + strconv.Itoa(i))
			assertNilErr(t, addGeneratedErr)
			if address != i {
				t.Errorf("Expected address is %d but got %d", i, address)
			}
		}
	})

	t.Run("AddGenerated should return error if address is exceed than max15BitValue", func(t *testing.T) {
		address, exceedErr := symbolsTable.AddGenerated("exceed-key")
		assertNotNilErr(t, exceedErr)
		if address != -1 {
			t.Errorf("Expected address is %d but got %d", -1, address)
		}
	})
}

func TestSymbolsTable_Add(t *testing.T) {
	symbolsTable := New16BitsSymbolsTable()
	firstSymbol := "some"
	t.Run("Add should add symbols when address is valid and symbols table not has key", func(t *testing.T) {
		minAddressErr := symbolsTable.Add(firstSymbol, 0)
		assertNilErr(t, minAddressErr)
		maxAddressErr := symbolsTable.Add(firstSymbol+"-1", Max15bitValue)
		assertNilErr(t, maxAddressErr)
	})

	t.Run("Add should return error when address is valid but symbols table already has key", func(t *testing.T) {
		err := symbolsTable.Add(firstSymbol, 1)
		assertNotNilErr(t, err)
	})

	t.Run("Add should return error when address is not valid", func(t *testing.T) {
		negativeAddressErr := symbolsTable.Add("new-some-key", -1)
		assertNotNilErr(t, negativeAddressErr)
		exceedAddressErr := symbolsTable.Add("new-some-key", Max15bitValue+1)
		assertNotNilErr(t, exceedAddressErr)
	})
}

func TestNew16BitsSymbolsTable(t *testing.T) {
	symbolsTable := New16BitsSymbolsTable()

	for i := 0; i <= 15; i++ {
		assertGetReturnExpectedAddress(symbolsTable, t, "R"+strconv.Itoa(i), i)
	}
	expectedScreenAddress := 16384
	assertGetReturnExpectedAddress(symbolsTable, t, "SCREEN", expectedScreenAddress)

	expectedKbdAddress := 24576
	assertGetReturnExpectedAddress(symbolsTable, t, "KBD", expectedKbdAddress)

	expectedSpAddress := 0
	assertGetReturnExpectedAddress(symbolsTable, t, "SP", expectedSpAddress)

	expectedLclAddress := 1
	assertGetReturnExpectedAddress(symbolsTable, t, "LCL", expectedLclAddress)

	expectedArgAddress := 2
	assertGetReturnExpectedAddress(symbolsTable, t, "ARG", expectedArgAddress)

	expectedThisAddress := 3
	assertGetReturnExpectedAddress(symbolsTable, t, "THIS", expectedThisAddress)

	expectedThatAddress := 4
	assertGetReturnExpectedAddress(symbolsTable, t, "THAT", expectedThatAddress)
}

func assertGetReturnExpectedAddress(symbolsTable *SymbolsTable, t *testing.T, key string, expectedAddress int) {
	address, ok := symbolsTable.Get(key)
	if !ok {
		t.Errorf("%q not found but expected", key)
	}
	if address != expectedAddress {
		t.Errorf("%q address got %d but expected %d", key, address, expectedAddress)
	}
}

func assertNotNilErr(t *testing.T, err error) {
	if err == nil {
		t.Error("Expected error but not return")
	}
}

func assertNilErr(t *testing.T, err error) {
	if err != nil {
		t.Error("Expected no error but got", err)
	}
}
