package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sayanli/calculator/internal/entity"
	"github.com/sayanli/calculator/internal/service"
)

type ErrorResponse struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type CalculationRouter struct {
	calculationService service.Calculation
}

func NewCalculationRouter(calculationService service.Calculation) *CalculationRouter {
	return &CalculationRouter{
		calculationService: calculationService,
	}
}

// Calculate godoc
// @Summary      Calculate instructions
// @Description  Calculate instructions
// @Accept       json
// @Produce      json
// @Success      200  {object}  []entity.Result
// @Failure      400  {object}  ErrorResponse
// @Router       /calculate [post]
func (c *CalculationRouter) Calculate(w http.ResponseWriter, r *http.Request) {
	var instructions []entity.Instruction

	if err := json.NewDecoder(r.Body).Decode(&instructions); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if len(instructions) == 0 {
		writeError(w, http.StatusUnprocessableEntity, "Empty instructions list")
		return
	}
	result, err := c.calculationService.CalculateInstructions(instructions)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := ErrorResponse{
		Code:      statusCode,
		Message:   message,
		Timestamp: time.Now().UTC(),
	}

	json.NewEncoder(w).Encode(response)
}
