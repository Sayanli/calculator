package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/sayanli/calculator/internal/entity"
	"golang.org/x/sync/errgroup"
)

type CalculationService struct {
	sync.Mutex
	cond *sync.Cond
}

func NewCalculationService() *CalculationService {
	service := &CalculationService{}
	service.cond = sync.NewCond(&service.Mutex)
	return service
}

func (s *CalculationService) CompleteInstructions(instructions []entity.Instruction) ([]entity.Result, error) {
	intermediateValues := make(map[string]int)
	results := make([]entity.Result, 0, len(instructions))
	tmp := make([]string, 0, len(instructions))
	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(1000)
	for _, instruction := range instructions {
		if instruction.Type == entity.CalcType {
			g.Go(func() error {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				return s.calculate(intermediateValues, instruction)
			})
		} else if instruction.Type == entity.PrintType {
			tmp = append(tmp, instruction.Var)
		}
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	for _, varName := range tmp {
		if _, ok := intermediateValues[varName]; !ok {
			return nil, fmt.Errorf("variable %s not found", varName)
		}
		results = append(results, entity.Result{
			Var:   varName,
			Value: intermediateValues[varName],
		})
	}
	return results, nil
}

func (s *CalculationService) calculate(intermediateValues map[string]int, instruction entity.Instruction) error {
	leftValue, err := s.getValue(intermediateValues, instruction.Left)
	if err != nil {
		return err
	}
	rightValue, err := s.getValue(intermediateValues, instruction.Right)
	if err != nil {
		return err
	}
	var variable int
	switch instruction.Op {
	case entity.AddOperation:
		variable = leftValue + rightValue
	case entity.SubOperation:
		variable = leftValue - rightValue
	case entity.MulOperation:
		variable = leftValue * rightValue
	default:
		return fmt.Errorf("Неизвестный оператор")
	}
	s.Lock()
	intermediateValues[instruction.Var] = variable
	s.cond.Broadcast()
	s.Unlock()
	return nil
}

func (s *CalculationService) getValue(intermediateValues map[string]int, input interface{}) (int, error) {
	switch v := input.(type) {
	case string:
		s.Lock()
		defer s.Unlock()
		for {
			if value, ok := intermediateValues[v]; ok {
				return value, nil
			}
			s.cond.Wait()
		}
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("Неизвестный тип данных: %T", input)
	}
}
