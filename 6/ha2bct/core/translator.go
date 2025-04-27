package core

import (
	"errors"
	"fmt"
	"ha2bct/core/computes"
	"ha2bct/core/dests"
	"ha2bct/core/jumps"
	"regexp"
	"strconv"
	"strings"
)

type HackAssemblerToBinaryCodeTranslator interface {
	Translate(rows []string) ([]string, error)
}

type HackAssemblerTo16BinaryCodeTranslator struct {
	symbolsTable SymbolsTable
}

var max15bitValue = 32768

func New16bitTranslator() HackAssemblerToBinaryCodeTranslator {
	return &HackAssemblerTo16BinaryCodeTranslator{
		symbolsTable: *NewSymbolsTable(16, max15bitValue),
	}
}

//TIP Translate
// + 1 - need clean comments for each rows (replace by expression '//.*' on ”)
// + 2 - skip blank rows
// + 3 - need remove spaces
// + 4 - need read all and add all labels by (value = index - count_labels) addr (check is 15bit value)
// 			need store count processed labels
// 5 - translate (skip labels)
// 5.1 - if A-command
// if alphabetic => check SymbolsTable (add addr if needed from [16...N])
// else check is 15bit and generate 15bit binary code and add first bit 0
// 5.2 - if C-command
// split by pattern (dests = comp; jumps) '=' and ';' is optional
// parts is [dests =]?, [comp]?, [;jmp]?
// every is part is enum and may be optional
// translate every part to binary code and concat

func (t *HackAssemblerTo16BinaryCodeTranslator) Translate(rows []string) ([]string, error) {
	cleanedRows := t.cleaned(rows)
	cleanedWithoutLabelsRows, err := t.processAndRemoveLabels(cleanedRows)
	if err != nil {
		return nil, err
	}
	return t.translateInstructions(cleanedWithoutLabelsRows)
}

// region cleaned

func (t *HackAssemblerTo16BinaryCodeTranslator) cleaned(rows []string) []string {
	cleanedRows := make([]string, 0, len(rows))
	for _, row := range rows {
		withoutCommentsRow := withoutComments(row)
		if isNotBlank(withoutCommentsRow) {
			cleanedRow := withoutSpaces(withoutCommentsRow)
			cleanedRows = append(cleanedRows, cleanedRow)
		}
	}
	return cleanedRows
}

var commentsRegex = regexp.MustCompile(`//.*`)

func withoutComments(value string) string {
	return commentsRegex.ReplaceAllString(value, "")
}

var spacesRegex = regexp.MustCompile(`\s+`)

func withoutSpaces(value string) string {
	return spacesRegex.ReplaceAllString(value, "")
}

// endregion

// region processAndRemoveLabels

func (t *HackAssemblerTo16BinaryCodeTranslator) processAndRemoveLabels(rows []string) ([]string, error) {
	withoutLabelsRows := make([]string, 0, len(rows))
	countProcessedLabels := 0
	for i, row := range rows {
		if label, ok := extractLabel(row); ok {
			address := i - countProcessedLabels
			countProcessedLabels++
			if err := t.symbolsTable.Add(label, address); err != nil {
				return nil, err
			}
		} else {
			withoutLabelsRows = append(withoutLabelsRows, row)
		}
	}
	return withoutLabelsRows, nil
}

var labelRegex = regexp.MustCompile(`^\(.*\)$`)

func extractLabel(row string) (string, bool) {
	if !labelRegex.MatchString(row) {
		return "", false
	}
	withoutBracketsRow := withoutBrackets(row)
	if isNotBlank(withoutBracketsRow) {
		return withoutBracketsRow, true
	}
	return "", false
}

func withoutBrackets(value string) string {
	runes := []rune(value)
	if len(runes) < 2 {
		return ""
	}
	return string(runes[1 : len(runes)-1])
}

// endregion

// region translateInstructions

func (t *HackAssemblerTo16BinaryCodeTranslator) translateInstructions(rows []string) ([]string, error) {
	translatedRows := make([]string, 0, len(rows))
	for _, row := range rows {
		translated, err := t.translateRow(row)
		if err != nil {
			return nil, err
		} else {
			translatedRows = append(translatedRows, translated)
		}
	}
	return translatedRows, nil
}

func (t *HackAssemblerTo16BinaryCodeTranslator) translateRow(row string) (string, error) {
	if parsedA, ok := extractAInstruction(row); ok {
		return t.processAInstruction(parsedA)
	} else if parsedC, ok := extractCInstruction(row); ok {
		return t.processCInstruction(parsedC)
	} else {
		return "", errors.New("bad row: " + row)
	}
}

var aInstructionRegex = regexp.MustCompile(`^@[A-Za-z0-9]+$`)

func extractAInstruction(row string) (string, bool) {
	if !aInstructionRegex.MatchString(row) {
		return "", false
	}
	runes := []rune(row)
	parsedAInstruction := string(runes[1:])
	return parsedAInstruction, true
}

func (t *HackAssemblerTo16BinaryCodeTranslator) processAInstruction(instruction string) (string, error) {
	address, err := t.resolveAddress(instruction)
	if err != nil {
		return "", err
	}
	binaryAddress := addressToBinaryCode(address)
	return "0" + binaryAddress, nil
}

func (t *HackAssemblerTo16BinaryCodeTranslator) resolveAddress(instruction string) (int, error) {
	if isAlphabetic(instruction) {
		address, err := t.symbolsTable.AddGenerated(instruction)
		if err != nil {
			return -1, err
		}
		return address, nil
	} else {
		address, err := strconv.Atoi(instruction)
		if err != nil {
			return -1, err
		}
		if address >= max15bitValue {
			return -1, errors.New(fmt.Sprintf("address: %d exceed max address%d", address, max15bitValue))
		}
		return address, nil
	}
}

func addressToBinaryCode(address int) string {
	return fmt.Sprintf("%015b", address)
}

var alphabeticRegex = regexp.MustCompile(`^[A-Za-z]+$`)

func isAlphabetic(value string) bool {
	return alphabeticRegex.MatchString(value)
}

type RawCInstruction struct {
	dest    string
	compute string
	jump    string
	origin  string
}

func extractCInstruction(row string) (RawCInstruction, bool) {
	destAndRightPart := strings.Split(row, "=")
	countParts := len(destAndRightPart)
	if countParts > 2 {
		return RawCInstruction{}, false
	}

	var dest string
	var rightPart string
	if countParts == 2 {
		dest = destAndRightPart[0]
		rightPart = destAndRightPart[1]
	} else {
		dest = ""
		rightPart = destAndRightPart[0]
	}

	computeAndJump := strings.Split(rightPart, ";")
	countParts = len(computeAndJump)
	if countParts > 2 {
		return RawCInstruction{}, false
	}
	var compute string
	var jump string
	if countParts == 2 {
		compute = computeAndJump[0]
		jump = computeAndJump[1]
	} else {
		compute = computeAndJump[0]
		jump = ""
	}
	return RawCInstruction{dest: dest, compute: compute, jump: jump, origin: row}, true
}

func (t *HackAssemblerTo16BinaryCodeTranslator) processCInstruction(rawInstruction RawCInstruction) (string, error) {
	dest, okDest := dests.NewDest(rawInstruction.dest)
	if !okDest {
		return "", errors.New("bad dest part: " + rawInstruction.origin)
	}
	compute, okCompute := computes.NewCompute(rawInstruction.compute)
	if !okCompute {
		return "", errors.New("bad compute part: " + rawInstruction.origin)
	}
	jump, okJump := jumps.NewJump(rawInstruction.jump)
	if !okJump {
		return "", errors.New("bad jump part: " + rawInstruction.origin)
	}

	destBinaryCode := destToBinaryCode(dest)
	computeBinaryCode := computeToBinaryCode(compute)
	jumpBinaryCode := jumpToBinaryCode(jump)

	return "111" + computeBinaryCode + destBinaryCode + jumpBinaryCode, nil
}

var jumpsBinaryCodes = map[jumps.Jump]string{
	jumps.Empty: "000",
	jumps.JGT:   "001",
	jumps.JEQ:   "010",
	jumps.JGE:   "011",
	jumps.JLT:   "100",
	jumps.JNE:   "101",
	jumps.JLE:   "110",
	jumps.JMP:   "111",
}

func jumpToBinaryCode(jump jumps.Jump) string {
	return jumpsBinaryCodes[jump]
}

var needLoadComputesBinaryCodes = map[computes.Compute]string{
	computes.M:         "110000",
	computes.NotM:      "110001",
	computes.NegativeM: "110011",
	computes.MPlusOne:  "110111",
	computes.MMinusOne: "110010",
	computes.DPlusM:    "000010",
	computes.DMinusM:   "010011",
	computes.MMinusD:   "000111",
	computes.DAndM:     "000000",
	computes.DOrM:      "010101",
}

var notNeedLoadComputesBinaryCodes = map[computes.Compute]string{
	computes.Zero:        "101010",
	computes.One:         "111111",
	computes.NegativeOne: "111010",
	computes.D:           "001100",
	computes.A:           "110000",
	computes.NotD:        "001101",
	computes.NotA:        "110001",
	computes.NegativeD:   "001111",
	computes.NegativeA:   "110011",
	computes.DPlusOne:    "011111",
	computes.APlusOne:    "110111",
	computes.DMinusOne:   "001110",
	computes.AMinusOne:   "110010",
	computes.DPlusA:      "000010",
	computes.DMinusA:     "010011",
	computes.AMinusD:     "000111",
	computes.DAndA:       "000000",
	computes.DOrA:        "010101",
}

func computeToBinaryCode(compute computes.Compute) string {
	if needLoadComputeBinaryCode, ok := needLoadComputesBinaryCodes[compute]; ok {
		return "1" + needLoadComputeBinaryCode
	} else {
		return "0" + notNeedLoadComputesBinaryCodes[compute]
	}
}

var destBinaryCodes = map[dests.Dest]string{
	dests.Empty: "000",
	dests.M:     "001",
	dests.D:     "010",
	dests.MD:    "011",
	dests.A:     "100",
	dests.AM:    "101",
	dests.AD:    "110",
	dests.AMD:   "111",
}

func destToBinaryCode(dest dests.Dest) string {
	return destBinaryCodes[dest]
}

// endregion

// region Helpers

func isNotBlank(value string) bool {
	return len(strings.TrimSpace(value)) > 0
}

// endregion
