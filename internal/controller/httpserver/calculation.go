package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/sayanli/calculator/internal/entity"
	"github.com/sayanli/calculator/internal/service"
)

type CalculationRouter struct {
	calculationService service.Calculation
}

func NewCalculationRouter(calculationService service.Calculation) *CalculationRouter {
	return &CalculationRouter{
		calculationService: calculationService,
	}
}

// @Summary Calculate instructions
// @Description Calculate instructions
// @Tags calculation
// @Accept json
// @Produce json
// @Success 200 {object} entity.Result
// @Failure 400 {object} errors.HTTPError
// @Failure 500 {object} errors.HTTPError
// @Router /v1/calculate [post]
func (c *CalculationRouter) Calculate(w http.ResponseWriter, r *http.Request) {
	var instructions []entity.Instruction

	err := json.NewDecoder(r.Body).Decode(&instructions)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result, err := c.calculationService.CompleteInstructions(instructions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
