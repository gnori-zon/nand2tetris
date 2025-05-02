package dests

type Dest int

const (
	Empty Dest = iota
	M
	D
	MD
	A
	AM
	AD
	AMD
)

var destByString = map[string]Dest{
	"":    Empty,
	"M":   M,
	"D":   D,
	"MD":  MD,
	"A":   A,
	"AM":  AM,
	"AD":  AD,
	"AMD": AMD,
}

func NewDest(value string) (Dest, bool) {
	dest, ok := destByString[value]
	return dest, ok
}
