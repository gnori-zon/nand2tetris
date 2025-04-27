package jumps

type Jump int

const (
	Empty Jump = iota
	JMP
	JEQ
	JNE
	JGT
	JGE
	JLT
	JLE
)

var jumpByString = map[string]Jump{
	"":    Empty,
	"JMP": JMP,
	"JEQ": JEQ,
	"JNE": JNE,
	"JGT": JGT,
	"JGE": JGE,
	"JLT": JLT,
	"JLE": JLE,
}

func NewJump(value string) (Jump, bool) {
	jump, ok := jumpByString[value]
	return jump, ok
}
