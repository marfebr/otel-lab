package web

import (
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/marfebr/otel-lab/service-a/internal/handler"
	"github.com/marfebr/otel-lab/service-a/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/trace"
)

// Server representa o servidor web
type Server struct {
	router     *chi.Mux
	cepHandler *handler.CEPHandler
	tracer     trace.Tracer
}

// NewServer cria uma nova instância do servidor
func NewServer(tracer trace.Tracer, serviceBURL string) *Server {
	// Criar validador de CEP
	cepValidator := service.NewCEPValidator()

	// Criar cliente do Serviço B
	serviceBClient := service.NewServiceBClient(serviceBURL, tracer)

	// Criar serviço de clima
	weatherService := service.NewWeatherService(serviceBClient, tracer)

	// Criar handler de CEP
	cepHandler := handler.NewCEPHandler(cepValidator, weatherService, tracer)

	// Criar router
	router := chi.NewRouter()

	// Configurar middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))

	// Configurar endpoints
	router.Handle("/metrics", promhttp.Handler())
	router.Post("/cep", cepHandler.HandleCEPValidation)

	return &Server{
		router:     router,
		cepHandler: cepHandler,
		tracer:     tracer,
	}
}

// GetRouter retorna o router configurado
func (s *Server) GetRouter() *chi.Mux {
	return s.router
}
