package entity

type InstructionType string

const (
	CalcType  InstructionType = "calc"
	PrintType InstructionType = "print"
)

type Instruction struct {
	Type  InstructionType `json:"type"`
	Op    string          `json:"op,omitempty"`
	Var   string          `json:"var"`
	Left  interface{}     `json:"left,omitempty"`
	Right interface{}     `json:"right,omitempty"`
}
