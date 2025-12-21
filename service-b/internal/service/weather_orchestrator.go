package service

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
)

// WeatherOrchestrator orquestra a busca de dados de clima por cidade
type WeatherOrchestrator struct {
	tracer trace.Tracer
}

// NewWeatherOrchestrator cria uma nova instância do orquestrador
func NewWeatherOrchestrator(tracer trace.Tracer) *WeatherOrchestrator {
	return &WeatherOrchestrator{
		tracer: tracer,
	}
}

// GetWeatherByCity orquestra o processo de busca de clima por cidade
func (o *WeatherOrchestrator) GetWeatherByCity(ctx context.Context, city string) (*WeatherResponse, error) {
	ctx, span := o.tracer.Start(ctx, "weather-orchestration")
	defer span.End()

	// Buscar dados de clima no WeatherAPI (padrão cloud-run)
	temps, err := GetWeatherAPICall(city)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Montar resposta final
	response := &WeatherResponse{
		City:  city,
		TempC: temps.TempC,
		TempF: temps.TempF,
		TempK: temps.TempK,
	}

	return response, nil
}

// GetWeatherByCEP orquestra o processo de busca de clima por CEP
func (o *WeatherOrchestrator) GetWeatherByCEP(ctx context.Context, cep string) (*WeatherResponse, error) {
	ctx, span := o.tracer.Start(ctx, "weather-orchestration")
	defer span.End()

	// Buscar cidade no ViaCEP
	address, err := BuscaViaCepApi(cep)
	if err != nil {
		span.RecordError(err)
		if err.Error() == "can not find zipcode" {
			return nil, ErrCEPNotFound
		}
		return nil, err
	}
	if address.City == "" {
		log.Printf("[DEBUG] Nome da cidade vazio para o CEP: %s", cep)
		span.RecordError(ErrCEPNotFound)
		return nil, ErrCEPNotFound
	}
	log.Printf("[DEBUG] Nome da cidade retornado pelo ViaCEP: %s", address.City)

	// Buscar dados de clima na WeatherAPI
	useMock := viper.GetString("WEATHER_API") == ""
	temps, err := GetWeatherAPICall(address.City)
	if err != nil {
		span.RecordError(err)
		if useMock {
			// Se mock está ativo, sempre retorna dados fixos
			return &WeatherResponse{
				City:  address.City,
				TempC: 25.0,
				TempF: 77.0,
				TempK: 298.15,
			}, nil
		}
		return nil, fmt.Errorf("can not find city")
	}

	// Montar resposta final
	response := &WeatherResponse{
		City:  address.City,
		TempC: temps.TempC,
		TempF: temps.TempF,
		TempK: temps.TempK,
	}

	return response, nil
}
