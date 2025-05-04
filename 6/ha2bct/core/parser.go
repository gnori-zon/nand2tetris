package core

import (
	"ha2bct/core/models"
	"ha2bct/core/models/computes"
	"ha2bct/core/models/dests"
	"ha2bct/core/models/jumps"
	"regexp"
	"strconv"
	"strings"
)

type HackAssemblerParser interface {
	Parse(row string) (models.HackAssemblerElement, bool)
}

type HackAssemblerElement16BitParser struct {
}

func NewHackAssembler16BitParser() *HackAssemblerElement16BitParser {
	return &HackAssemblerElement16BitParser{}
}

func (parser *HackAssemblerElement16BitParser) Parse(row string) (models.HackAssemblerElement, bool) {
	if cleanedRow, ok := cleaned(row); ok {
		if parsedElement, okParsed := parseHackAssemblerElement(cleanedRow); okParsed {
			return parsedElement, true
		}
	}
	return nil, false
}

// region cleaned
func cleaned(row string) (string, bool) {
	withoutCommentsRow := withoutComments(row)
	if isNotBlank(withoutCommentsRow) {
		return withoutSpaces(withoutCommentsRow), true
	}
	return "", false
}

var commentsRegex = regexp.MustCompile(`//.*`)

func withoutComments(value string) string {
	return commentsRegex.ReplaceAllString(value, "")
}

func isNotBlank(value string) bool {
	return len(strings.TrimSpace(value)) > 0
}

var spacesRegex = regexp.MustCompile(`\s+`)

func withoutSpaces(value string) string {
	return spacesRegex.ReplaceAllString(value, "")
}

// endregion

func parseHackAssemblerElement(row string) (models.HackAssemblerElement, bool) {
	if label, ok := extractLabel(row); ok {
		return label, true
	}
	if aInstruction, ok := extractAInstruction(row); ok {
		return aInstruction, true
	}
	if cInstruction, ok := extractCInstruction(row); ok {
		return cInstruction, true
	}
	return nil, false
}

// region extractLabel

func extractLabel(row string) (models.HackAssemblerLabel, bool) {
	if !labelRegex.MatchString(row) {
		return models.HackAssemblerLabel{}, false
	}
	withoutBracketsRow := withoutBrackets(row)
	return models.HackAssemblerLabel{Value: withoutBracketsRow}, true
}

var labelRegex = regexp.MustCompile(`^\(.+\)$`)

func withoutBrackets(value string) string {
	runes := []rune(value)
	if len(runes) < 2 {
		return ""
	}
	return string(runes[1 : len(runes)-1])
}

// endregion

// region extractAInstruction

func extractAInstruction(row string) (models.HackAssemblerAInstruction, bool) {
	if !aInstructionRegex.MatchString(row) {
		return nil, false
	}
	runes := []rune(row)
	instructionValue := string(runes[1:])
	isNumeric := numericRegex.MatchString(instructionValue)
	if isNumeric {
		if intInstructionValue, err := strconv.Atoi(instructionValue); err == nil {
			return models.HackAssemblerNumericAInstruction{Value: intInstructionValue}, true
		}
	} else {
		return models.HackAssemblerAlphabeticAInstruction{Value: instructionValue}, true
	}
	return nil, false
}

var aInstructionRegex = regexp.MustCompile(`^@[A-Za-z_$.\-0-9]+$`)

var numericRegex = regexp.MustCompile(`^[0-9]+$`)

// endregion

// region extractCInstruction

// examples C-instruction:
//
//	M=D-1;JMP
//	D-1;JMP
//	;JMP
func extractCInstruction(row string) (models.HackAssemblerCInstruction, bool) {
	dest, rightPart, okSplitOnDestAndRightPart := splitOnDestAndRightPart(row)
	if !okSplitOnDestAndRightPart {
		return models.HackAssemblerCInstruction{}, false
	}
	compute, jump, okSplitOnComputeAndJump := splitOnComputeAndJump(rightPart)
	if !okSplitOnComputeAndJump {
		return models.HackAssemblerCInstruction{}, false
	}
	return models.HackAssemblerCInstruction{Dest: dest, Compute: compute, Jump: jump}, true
}

func splitOnDestAndRightPart(row string) (dests.Dest, string, bool) {
	destAndRightPart := strings.Split(row, "=")
	countParts := len(destAndRightPart)
	if countParts > 2 {
		return -1, "", false
	}
	var rawDest string
	var rightPart string
	if countParts == 2 {
		rawDest = destAndRightPart[0]
		rightPart = destAndRightPart[1]
	} else {
		rawDest = ""
		rightPart = destAndRightPart[0]
	}
	if dest, ok := dests.NewDest(rawDest); ok {
		return dest, rightPart, true
	}
	return -1, "", false
}

func splitOnComputeAndJump(rightPart string) (computes.Compute, jumps.Jump, bool) {
	computeAndJump := strings.Split(rightPart, ";")
	countParts := len(computeAndJump)
	if countParts > 2 {
		return -1, -1, false
	}
	var rawCompute string
	var rawJump string
	if countParts == 2 {
		rawCompute = computeAndJump[0]
		rawJump = computeAndJump[1]
	} else {
		rawCompute = computeAndJump[0]
		rawJump = ""
	}
	compute, okCompute := computes.NewCompute(rawCompute)
	jump, okJump := jumps.NewJump(rawJump)
	if !okCompute || !okJump {
		return -1, -1, false
	}
	return compute, jump, true
}

// endregion
