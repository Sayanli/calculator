package entity

type InstructionType string

type Instruction struct {
	Type  string      `json:"type"`
	Op    string      `json:"op,omitempty"`
	Var   string      `json:"var"`
	Left  interface{} `json:"left,omitempty"`
	Right interface{} `json:"right,omitempty"`
}

type Result struct {
	Var   string `json:"var"`
	Value int64  `json:"value"`
}
