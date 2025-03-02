package httpserver

import (
	"github.com/go-chi/chi"
	_ "github.com/sayanli/calculator/docs"
	"github.com/sayanli/calculator/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

type HttpServer struct {
	calculationService service.Calculation
}

func NewHttpServer(calculationService service.Calculation) *HttpServer {
	return &HttpServer{
		calculationService: calculationService,
	}
}

// @title           Calculator Service
// @version         1.0
// @description     This is a service for performing basic arithmetic operations.
// @contact.name   Ilya Veselov
// @contact.email  notmydev@gmail.com

// @host      localhost:8080
// @BasePath  /
func (h *HttpServer) Router() *chi.Mux {
	r := chi.NewRouter()
	CalculationRouter := NewCalculationRouter(h.calculationService)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	r.Route("/v1", func(r chi.Router) {
		r.Post("/calculate", CalculationRouter.Calculate)
	})
	return r
}
