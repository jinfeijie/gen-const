package exp

type Num int

//go:generate gen-const -type Num
const (
	numStart = Num(iota) // start
	numA                 // A
	numB                 // B
	numC                 // C
	numD                 // D
	numE                 // E
)
