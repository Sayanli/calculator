package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/sayanli/calculator/internal/entity"
	"golang.org/x/sync/errgroup"
)

type CalculationService struct {
	log     *slog.Logger
	golimit int
}

type InstructionsState struct {
	intermediateValues map[string]int64
	mask               map[string]struct{}
	pr                 []string
	mu                 sync.Mutex
	cond               *sync.Cond
}

func NewInstructionsState() *InstructionsState {
	state := &InstructionsState{
		intermediateValues: make(map[string]int64),
		mask:               make(map[string]struct{}),
		pr:                 make([]string, 0),
		mu:                 sync.Mutex{},
	}
	state.cond = sync.NewCond(&state.mu)
	return state
}

func NewCalculationService(log *slog.Logger, golimit int) *CalculationService {
	return &CalculationService{
		log:     log,
		golimit: golimit,
	}
}

func (s *CalculationService) CalculateInstructions(instructions []entity.Instruction) ([]entity.Result, error) {
	const op = "calculationService.CalculateInstructions"
	log := s.log.With(slog.String("op", op))

	state := NewInstructionsState()

	tmp := make([]string, 0, len(instructions))

	for _, instruction := range instructions {
		if instruction.Type == "calc" {
			state.mask[instruction.Var] = struct{}{}
		}
	}

	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(s.golimit)

	for _, instruction := range instructions {
		if instruction.Type == "calc" {
			g.Go(func() error {
				return s.calculate(ctx, state, instruction)
			})
		} else if instruction.Type == "print" {
			tmp = append(tmp, instruction.Var)
		}
	}

	if err := g.Wait(); err != nil {
		log.Error("error waiting for goroutines", "error", err)
		return nil, err
	}
	results := make([]entity.Result, 0, len(tmp))
	for _, varName := range tmp {
		if _, ok := state.intermediateValues[varName]; !ok {
			log.Error("variable not found", "variable", varName)
			return nil, fmt.Errorf("variable %s not found", varName)
		}
		results = append(results, entity.Result{
			Var:   varName,
			Value: state.intermediateValues[varName],
		})
	}
	return results, nil
}

func (s *CalculationService) calculate(ctx context.Context, state *InstructionsState, instruction entity.Instruction) error {
	const op = "calculationService.calculate"
	log := s.log.With("op", op)

	leftValue, err := s.getValue(ctx, state, instruction.Left)
	if err != nil {
		log.Error("error getting left value", "error", err)
		state.cond.Broadcast()
		return err
	}
	rightValue, err := s.getValue(ctx, state, instruction.Right)
	if err != nil {
		log.Error("error getting right value", "error", err)
		state.cond.Broadcast()
		return err
	}

	var variable int64
	switch instruction.Op {
	case "+":
		variable = leftValue + rightValue
	case "-":
		variable = leftValue - rightValue
	case "*":
		variable = leftValue * rightValue
	default:
		log.Error("unknown operator")
		state.cond.Broadcast()
		return fmt.Errorf("unknown operator")
	}

	state.mu.Lock()
	state.intermediateValues[instruction.Var] = variable
	state.cond.Broadcast()
	state.mu.Unlock()
	return nil
}

func (s *CalculationService) getValue(ctx context.Context, state *InstructionsState, input interface{}) (int64, error) {
	switch v := input.(type) {
	case string:
		state.mu.Lock()
		defer state.mu.Unlock()
		for {
			if _, ok := state.mask[v]; !ok {
				return 0, fmt.Errorf("unknown variable %s", v)
			}
			if value, ok := state.intermediateValues[v]; ok {
				return value, nil
			}
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			default:

			}
			state.cond.Wait()
		}
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("unknown type: %T", input)
	}
}
