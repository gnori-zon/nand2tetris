package core

import (
	"fmt"
	"ha2bct/core/models/computes"
	"ha2bct/core/models/dests"
	"ha2bct/core/models/jumps"
)

type BinaryCodeProvider interface {
	EncodeAInstruction(address int) (string, error)
	EncodeCInstruction(dest dests.Dest, compute computes.Compute, jump jumps.Jump) string
}

type BinaryCode16BitsProvider struct {
}

func NewBinaryCode16BitsProvider() BinaryCodeProvider {
	return &BinaryCode16BitsProvider{}
}

var addressTo15Bits = "%015b"

func (p *BinaryCode16BitsProvider) EncodeAInstruction(address int) (string, error) {
	if err := Validate15BitAddress(address); err != nil {
		return "", err
	}
	return "0" + fmt.Sprintf(addressTo15Bits, address), nil
}

func (p *BinaryCode16BitsProvider) EncodeCInstruction(dest dests.Dest, compute computes.Compute, jump jumps.Jump) string {
	destBinaryCode := encodeDest(dest)
	computeBinaryCode := encodeCompute(compute)
	jumpBinaryCode := encodeJump(jump)
	return "111" + computeBinaryCode + destBinaryCode + jumpBinaryCode
}

// region encodeDest
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

func encodeDest(dest dests.Dest) string {
	return destBinaryCodes[dest]
}

// endregion

// region encodeCompute
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

func encodeCompute(compute computes.Compute) string {
	if needLoadComputeBinaryCode, ok := needLoadComputesBinaryCodes[compute]; ok {
		return "1" + needLoadComputeBinaryCode
	} else {
		return "0" + notNeedLoadComputesBinaryCodes[compute]
	}
}

// endregion

// region encodeJump
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

func encodeJump(jump jumps.Jump) string {
	return jumpsBinaryCodes[jump]
}

// endregion
