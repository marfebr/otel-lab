package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// ServiceBClient cliente para comunicação com Serviço B
type ServiceBClient struct {
	baseURL string
	client  *http.Client
	tracer  trace.Tracer
}

// NewServiceBClient cria uma nova instância do cliente do Serviço B
func NewServiceBClient(baseURL string, tracer trace.Tracer) *ServiceBClient {
	return &ServiceBClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		tracer: tracer,
	}
}

// GetWeatherByCEP envia o CEP ao Service B e retorna o clima
func (c *ServiceBClient) GetWeatherByCEP(ctx context.Context, cep string) (*WeatherResponse, error) {
	ctx, span := c.tracer.Start(ctx, "service-b-weather-request")
	defer span.End()

	log.Printf("Enviando CEP para Service B: %s", cep)
	// Preparar request
	requestBody := CEPRequest{CEP: cep}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Criar requisição HTTP
	url := fmt.Sprintf("%s/weather", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Configurar headers
	req.Header.Set("Content-Type", "application/json")

	// Propagar contexto OTEL
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	// Executar requisição
	resp, err := c.client.Do(req)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Ler resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Verificar status code
	if resp.StatusCode != http.StatusOK {
		span.RecordError(fmt.Errorf("service B returned status %d", resp.StatusCode))

		// Tentar decodificar erro
		var errorResp ErrorResponse
		if json.Unmarshal(body, &errorResp) == nil {
			return nil, fmt.Errorf("service B error: %s", errorResp.Error)
		}

		return nil, fmt.Errorf("service B returned status %d: %s", resp.StatusCode, string(body))
	}

	// Decodificar resposta de sucesso
	var weatherResp WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &weatherResp, nil
}
