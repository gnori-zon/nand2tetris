package computes

type Compute int

const (
	Zero Compute = iota
	One
	NegativeOne
	D
	A
	NotD
	NotA
	NegativeD
	NegativeA
	DPlusOne
	APlusOne
	DMinusOne
	AMinusOne
	DPlusA
	DMinusA
	AMinusD
	DAndA
	DOrA
	M
	NotM
	NegativeM
	MPlusOne
	MMinusOne
	DPlusM
	DMinusM
	MMinusD
	DAndM
	DOrM
)

var computeByString = map[string]Compute{
	"0":   Zero,
	"1":   One,
	"-1":  NegativeOne,
	"D":   D,
	"A":   A,
	"!D":  NotD,
	"!A":  NotA,
	"-D":  NegativeD,
	"-A":  NegativeA,
	"D+1": DPlusOne,
	"A+1": APlusOne,
	"D-1": DMinusOne,
	"A-1": AMinusOne,
	"D+A": DPlusA,
	"D-A": DMinusA,
	"A-D": AMinusD,
	"D&A": DAndA,
	"D|A": DOrA,
	"M":   M,
	"!M":  NotM,
	"-M":  NegativeM,
	"M+1": MPlusOne,
	"M-1": MMinusOne,
	"D+M": DPlusM,
	"D-M": DMinusM,
	"M-D": MMinusD,
	"D&M": DAndM,
	"D|M": DOrM,
}

func NewCompute(value string) (Compute, bool) {
	if compute, ok := computeByString[value]; ok {
		return compute, true
	}
	return -1, false
}
