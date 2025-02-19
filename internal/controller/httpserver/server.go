package httpserver

import (
	"github.com/go-chi/chi"
	"github.com/sayanli/calculator/internal/service"
)

type HttpServer struct {
	calculationService service.Calculation
}

func NewHttpServer(calculationService service.Calculation) *HttpServer {
	return &HttpServer{
		calculationService: calculationService,
	}
}

func (h *HttpServer) Router() *chi.Mux {
	r := chi.NewRouter()
	CalculationRouter := NewCalculationRouter(h.calculationService)
	r.Route("/v1", func(r chi.Router) {
		r.Post("/calculate", CalculationRouter.Calculate)
	})
	return r
}
