package service

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// WeatherService interface para serviços de clima
type WeatherService interface {
	GetWeatherByCEP(ctx context.Context, cep string) (*WeatherResponse, error)
}

// weatherService implementação do serviço de clima
type weatherService struct {
	serviceBClient *ServiceBClient
	tracer         trace.Tracer
}

// NewWeatherService cria uma nova instância do serviço de clima
func NewWeatherService(serviceBClient *ServiceBClient, tracer trace.Tracer) WeatherService {
	return &weatherService{
		serviceBClient: serviceBClient,
		tracer:         tracer,
	}
}

// GetWeatherByCEP busca dados de clima por CEP
func (s *weatherService) GetWeatherByCEP(ctx context.Context, cep string) (*WeatherResponse, error) {
	ctx, span := s.tracer.Start(ctx, "request-weather-by-cep")
	defer span.End()

	// Apenas encaminhar o CEP para o Service B
	return s.serviceBClient.GetWeatherByCEP(ctx, cep)
}
