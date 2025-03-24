package grpcserver

import (
	"context"
	"log/slog"

	"github.com/sayanli/calculator/internal/entity"
	calculator "github.com/sayanli/calculator/protos/gen/go/calculator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Calculator interface {
	CalculateInstructions(instructions []entity.Instruction) ([]entity.Result, error)
}

type server struct {
	calculator.UnimplementedCalculatorServer
	log         *slog.Logger
	calcService Calculator
}

func RegisterServer(s *grpc.Server, log *slog.Logger, calcService Calculator) {
	calculator.RegisterCalculatorServer(s, &server{
		log:         log,
		calcService: calcService,
	})
}

func (s *server) Calculate(ctx context.Context, req *calculator.OperationsRequest) (*calculator.OperationsResponse, error) {
	const op = "grpcserver.Calculate"
	log := s.log.With(slog.String("op", op))

	instructions := make([]entity.Instruction, len(req.Operations))
	instructions = parseInstructions(req.Operations)
	results, err := s.calcService.CalculateInstructions(instructions)
	if err != nil {
		log.Error("failed to calculate instructions", "error", err)

		return nil, status.Errorf(codes.InvalidArgument, "invalid instructions: %v", err)
	}
	res := make([]*calculator.Result, len(results))
	for i, result := range results {
		res[i] = &calculator.Result{
			Var:   result.Var,
			Value: result.Value,
		}
	}
	return &calculator.OperationsResponse{
		Results: res,
	}, nil
}

func parseInstructions(instructions []*calculator.Operation) []entity.Instruction {
	var parsedInstructions []entity.Instruction
	for _, op := range instructions {
		var left interface{}
		var right interface{}
		if op.Left.GetStringValue() == "" {
			left = op.Left.GetNumberValue()
		} else {
			left = op.Left.GetStringValue()
		}
		if op.Right.GetStringValue() == "" {
			right = op.Right.GetNumberValue()
		} else {
			right = op.Right.GetStringValue()
		}
		instruction := entity.Instruction{
			Type:  op.Type,
			Var:   op.Var,
			Op:    op.Op,
			Left:  left,
			Right: right,
		}
		parsedInstructions = append(parsedInstructions, instruction)
	}
	return parsedInstructions
}
