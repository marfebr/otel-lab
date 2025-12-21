package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/marfebr/otel-lab/service-b/internal/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// WeatherHandler handler para endpoints relacionados a clima
type WeatherHandler struct {
	weatherOrchestrator *service.WeatherOrchestrator
	tracer              trace.Tracer
}

// NewWeatherHandler cria uma nova instância do handler de clima
func NewWeatherHandler(weatherOrchestrator *service.WeatherOrchestrator, tracer trace.Tracer) *WeatherHandler {
	return &WeatherHandler{
		weatherOrchestrator: weatherOrchestrator,
		tracer:              tracer,
	}
}

// HandleWeatherRequest processa a requisição de dados de clima
func (h *WeatherHandler) HandleWeatherRequest(w http.ResponseWriter, r *http.Request) {
	// Extrair contexto OTEL do header HTTP
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	// Criar span para tracing
	ctx, span := h.tracer.Start(ctx, "weather-request")
	defer span.End()

	// Verificar método HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar request
	var req service.CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("CEP recebido no handler: %s", req.CEP)

	// Buscar dados de clima
	weatherResp, err := h.weatherOrchestrator.GetWeatherByCEP(ctx, req.CEP)
	if err != nil {
		span.RecordError(err)
		log.Printf("Erro retornado por GetWeatherByCEP: %v", err)
		// Verificar tipo de erro e retornar status code apropriado
		switch err {
		case service.ErrInvalidCEP:
			h.sendErrorResponse(w, err.Error(), http.StatusUnprocessableEntity)
			return
		case service.ErrCEPNotFound:
			h.sendErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Retornar dados de clima
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weatherResp)
}

// sendErrorResponse envia resposta de erro padronizada
func (h *WeatherHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(service.ErrorResponse{
		Error: message,
	})
}
