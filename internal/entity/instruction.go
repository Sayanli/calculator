package entity

type InstructionType string

const (
	CalcType  InstructionType = "calc"
	PrintType InstructionType = "print"
)

type OperationType string

const (
	AddOperation OperationType = "+"
	SubOperation OperationType = "-"
	MulOperation OperationType = "*"
)

type Instruction struct {
	Type  InstructionType `json:"type"`
	Op    OperationType   `json:"op,omitempty"`
	Var   string          `json:"var"`
	Left  interface{}     `json:"left,omitempty"`
	Right interface{}     `json:"right,omitempty"`
}

type Result struct {
	Var   string `json:"var"`
	Value int    `json:"value"`
}
