package service

import "github.com/sayanli/calculator/internal/entity"

type Calculation interface {
	CompleteInstructions(instructions []entity.Instruction) ([]entity.Result, error)
}

type Services struct {
	Calculation Calculation
}

func NewServices() *Services {
	return &Services{
		Calculation: NewCalculationService(),
	}
}
