package service

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/sayanli/calculator/internal/entity"
)

func TestCalculateInstructions(t *testing.T) {
	buf := &bytes.Buffer{}
	log := slog.New(slog.NewTextHandler(buf, nil))
	svc := NewCalculationService(log, 100)

	tests := []struct {
		name         string
		instructions []entity.Instruction
		expected     []entity.Result
		expectErr    bool
	}{
		{
			name: "OK",
			instructions: []entity.Instruction{
				{Type: "calc", Op: "+", Var: "x", Left: int64(2), Right: int64(3)},
				{Type: "calc", Op: "*", Var: "y", Left: "x", Right: int64(4)},
				{Type: "calc", Op: "-", Var: "z", Left: "x", Right: "y"},
				{Type: "print", Var: "x"},
				{Type: "print", Var: "y"},
				{Type: "print", Var: "z"},
			},
			expected: []entity.Result{
				{Var: "x", Value: 5},
				{Var: "y", Value: 20},
				{Var: "z", Value: -15},
			},
			expectErr: false,
		},
		{
			name: "Missing Variable",
			instructions: []entity.Instruction{
				{Type: "print", Var: "z"},
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "Invalid Operation",
			instructions: []entity.Instruction{
				{Type: "calc", Op: "invalid", Var: "x", Left: int64(2), Right: int64(3)},
				{Type: "calc", Op: "+", Var: "y", Left: string("x"), Right: int64(2)},
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "float values",
			instructions: []entity.Instruction{
				{Type: "calc", Op: "+", Var: "x", Left: float64(2.0), Right: float64(3.0)},
				{Type: "calc", Op: "*", Var: "y", Left: "x", Right: float64(4.0)},
				{Type: "calc", Op: "-", Var: "z", Left: "x", Right: "y"},
				{Type: "print", Var: "x"},
				{Type: "print", Var: "y"},
				{Type: "print", Var: "z"},
			},
			expected: []entity.Result{
				{Var: "x", Value: 5.0},
				{Var: "y", Value: 20.0},
				{Var: "z", Value: -15.0},
			},
			expectErr: false,
		},
		{
			name: "invalid left type",
			instructions: []entity.Instruction{
				{Type: "calc", Op: "+", Var: "x", Left: true, Right: int64(3)},
				{Type: "print", Var: "x"},
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "invalid right type",
			instructions: []entity.Instruction{
				{Type: "calc", Op: "+", Var: "x", Left: int64(3), Right: true},
				{Type: "print", Var: "x"},
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "unknown variable",
			instructions: []entity.Instruction{
				{Type: "calc", Op: "+", Var: "x", Left: "unknown", Right: int64(4)},
				{Type: "print", Var: "x"},
			},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := svc.CalculateInstructions(tt.instructions)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if len(results) != len(tt.expected) {
				t.Fatalf("expected %d results, got %d", len(tt.expected), len(results))
			}

			for i, res := range results {
				if res.Var != tt.expected[i].Var || res.Value != tt.expected[i].Value {
					t.Fatalf("expected %s = %d, got %s = %d", tt.expected[i].Var, tt.expected[i].Value, res.Var, res.Value)
				}
			}
		})
	}
}

func BenchmarkCalculateInstructions(b *testing.B) {
	buf := &bytes.Buffer{}
	log := slog.New(slog.NewTextHandler(buf, nil))
	svc := NewCalculationService(log, 100)
	instructions := []entity.Instruction{
		{Type: "calc", Op: "+", Var: "x", Left: int64(10), Right: int64(2)},
		{Type: "calc", Op: "*", Var: "y", Left: "x", Right: int64(5)},
		{Type: "calc", Op: "-", Var: "q", Left: "y", Right: int64(20)},
		{Type: "calc", Op: "+", Var: "unusedA", Left: "y", Right: int64(100)},
		{Type: "calc", Op: "*", Var: "unusedB", Left: "unusedA", Right: int64(2)},
		{Type: "print", Var: "q"},
		{Type: "calc", Op: "-", Var: "z", Left: "x", Right: int64(15)},
		{Type: "print", Var: "z"},
		{Type: "calc", Op: "+", Var: "ignoreC", Left: "z", Right: "y"},
		{Type: "print", Var: "x"},
	}

	for i := 0; i < b.N; i++ {
		_, err := svc.CalculateInstructions(instructions)
		if err != nil {
			b.Fatal(err)
		}

	}
}
