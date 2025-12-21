package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/marfebr/otel-lab/service-a/internal/service"
	"go.opentelemetry.io/otel/trace"
)

// CEPHandler handler para endpoints relacionados a CEP
type CEPHandler struct {
	cepValidator   service.CEPValidator
	weatherService service.WeatherService
	tracer         trace.Tracer
}

// NewCEPHandler cria uma nova instância do handler de CEP
func NewCEPHandler(cepValidator service.CEPValidator, weatherService service.WeatherService, tracer trace.Tracer) *CEPHandler {
	return &CEPHandler{
		cepValidator:   cepValidator,
		weatherService: weatherService,
		tracer:         tracer,
	}
}

// HandleCEPValidation processa a validação de CEP
func (h *CEPHandler) HandleCEPValidation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Criar span para tracing
	ctx, span := h.tracer.Start(ctx, "cep-validation")
	defer span.End()

	// Verificar método HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar request
	var req CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validar CEP
	if err := h.cepValidator.ValidateCEP(req.CEP); err != nil {
		span.RecordError(err)
		h.sendErrorResponse(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Buscar dados de clima no Serviço B
	weatherResp, err := h.weatherService.GetWeatherByCEP(ctx, req.CEP)
	if err != nil {
		span.RecordError(err)
		log.Printf("Erro retornado por GetWeatherByCEP: %v", err)
		// Verificar se é erro do Serviço B e propagar status code
		if err.Error() == "service B error: invalid zipcode" {
			h.sendErrorResponse(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}
		if err.Error() == "service B error: can not find zipcode" || strings.Contains(err.Error(), "can not find zipcode") {
			h.sendErrorResponse(w, "can not find zipcode", http.StatusNotFound)
			return
		}
		// Erro interno do servidor
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Retornar dados de clima
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weatherResp)
}

// sendErrorResponse envia resposta de erro padronizada
func (h *CEPHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	})
}
