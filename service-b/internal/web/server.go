package web

import (
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/marfebr/otel-lab/service-b/internal/handler"
	"github.com/marfebr/otel-lab/service-b/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/trace"
)

// Server representa o servidor web
type Server struct {
	router         *chi.Mux
	weatherHandler *handler.WeatherHandler
	tracer         trace.Tracer
}

// NewServer cria uma nova instância do servidor
func NewServer(
	tracer trace.Tracer,
) *Server {
	// Criar handler de clima
	weatherHandler := handler.NewWeatherHandler(service.NewWeatherOrchestrator(tracer), tracer)

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
	router.Post("/weather", weatherHandler.HandleWeatherRequest)

	return &Server{
		router:         router,
		weatherHandler: weatherHandler,
		tracer:         tracer,
	}
}

// GetRouter retorna o router configurado
func (s *Server) GetRouter() *chi.Mux {
	return s.router
}
