package service

import (
	"log/slog"

	"github.com/sayanli/calculator/internal/entity"
)

type Calculation interface {
	CalculateInstructions(instructions []entity.Instruction) ([]entity.Result, error)
}

type Services struct {
	log         *slog.Logger
	Calculation Calculation
}

func NewServices(log *slog.Logger, golimit int) *Services {
	return &Services{
		Calculation: NewCalculationService(log, golimit),
	}
}
