package service

import (
	"fmt"

	"github.com/sayanli/calculator/internal/entity"
)

type CalculationService struct {
}

func NewCalculationService() *CalculationService {
	return &CalculationService{}
}

func (s *CalculationService) CompleteInstructions(instructions []entity.Instruction) ([]entity.Result, error) {
	result := make([]entity.Result, 0, len(instructions))
	intermediateValues := make(map[string]int)
	for _, instruction := range instructions {
		if instruction.Type == entity.CalcType {
			if err := Calculate(intermediateValues, instruction); err != nil {
				return nil, err
			}
		} else if instruction.Type == entity.PrintType {
			result = append(result, entity.Result{
				Var:   instruction.Var,
				Value: intermediateValues[instruction.Var],
			})
		}
	}
	return result, nil
}

func Calculate(variables map[string]int, instruction entity.Instruction) error {
	leftValue, err := getValue(variables, instruction.Left)
	if err != nil {
		return err
	}
	rightValue, err := getValue(variables, instruction.Right)
	if err != nil {
		return err
	}
	switch instruction.Op {
	case entity.AddOperation:
		variables[instruction.Var] = leftValue + rightValue
	case entity.SubOperation:
		variables[instruction.Var] = leftValue - rightValue
	case entity.MulOperation:
		variables[instruction.Var] = leftValue * rightValue
	default:
		return fmt.Errorf("Неизвестный оператор")
	}
	return nil
}

func getValue(variables map[string]int, input interface{}) (int, error) {
	switch v := input.(type) {
	case string:
		if value, ok := variables[v]; ok {
			return value, nil
		} else {
			return 0, fmt.Errorf("Переменная %s не найдена", v)
		}
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("Неизвестный тип данных: %T", input)
	}
}
