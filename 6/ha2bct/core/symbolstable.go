package core

import (
	"errors"
	"fmt"
)

type SymbolsTable struct {
	data                     map[string]int
	currentAddressToGenerate int
	maxAddressToGenerate     int
}

func New16BitsSymbolsTable() *SymbolsTable {
	symbolsTable := &SymbolsTable{
		data:                     make(map[string]int),
		currentAddressToGenerate: 16,
		maxAddressToGenerate:     Max15bitValue,
	}
	symbolsTable.initDefaultSymbols()
	return symbolsTable
}

func (t *SymbolsTable) Get(name string) (int, bool) {
	address, ok := t.data[name]
	return address, ok
}

func (t *SymbolsTable) Add(name string, address int) error {
	if existAddress, ok := t.Get(name); ok {
		return errors.New(fmt.Sprintf("name already has address: %d", existAddress))
	}
	if address < 0 {
		return errors.New("address must be positive")
	}
	if address >= t.maxAddressToGenerate {
		return errors.New("new generated address exceed max address")
	}
	t.data[name] = address
	return nil
}

func (t *SymbolsTable) AddGenerated(name string) (int, error) {
	if existAddress, ok := t.data[name]; ok {
		return -1, errors.New(fmt.Sprintf("name already has address: %d", existAddress))
	}
	if t.currentAddressToGenerate >= t.maxAddressToGenerate {
		return -1, errors.New("new generated address exceed max address")
	}
	t.data[name] = t.currentAddressToGenerate
	defer t.incrementAddressToGenerate()
	return t.currentAddressToGenerate, nil
}

func (t *SymbolsTable) incrementAddressToGenerate() {
	t.currentAddressToGenerate++
}

func (t *SymbolsTable) initDefaultSymbols() {
	for i := 0; i <= 15; i++ {
		name := fmt.Sprintf("R%d", i)
		t.data[name] = i
	}
	t.data["SCREEN"] = 16384
	t.data["KBD"] = 24576
	t.data["SP"] = 0
	t.data["LCL"] = 1
	t.data["ARG"] = 2
	t.data["THIS"] = 3
	t.data["THAT"] = 4
}
