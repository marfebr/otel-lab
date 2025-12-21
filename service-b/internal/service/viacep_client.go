package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// ViaCEPClient cliente para a API ViaCEP
type ViaCEPClient struct {
	baseURL string
	client  *http.Client
	tracer  trace.Tracer
}

// NewViaCEPClient cria uma nova instância do cliente ViaCEP
func NewViaCEPClient(baseURL string, tracer trace.Tracer) *ViaCEPClient {
	return &ViaCEPClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		tracer: tracer,
	}
}

// GetAddressByCEP busca endereço por CEP na API ViaCEP
func (c *ViaCEPClient) GetAddressByCEP(ctx context.Context, cep string) (*ViaCEPResponse, error) {
	ctx, span := c.tracer.Start(ctx, "viacep-request")
	defer span.End()

	// Criar URL da requisição
	url := fmt.Sprintf("%s/%s/json", c.baseURL, cep)

	// Criar requisição HTTP
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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

	// Decodificar resposta
	var viaCEPResp ViaCEPResponse
	if err := json.Unmarshal(body, &viaCEPResp); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Verificar se o CEP foi encontrado
	if viaCEPResp.Erro {
		span.RecordError(ErrCEPNotFound)
		return nil, ErrCEPNotFound
	}

	// Verificar se a localidade está vazia (CEP inválido)
	if viaCEPResp.Localidade == "" {
		span.RecordError(ErrCEPNotFound)
		return nil, ErrCEPNotFound
	}

	return &viaCEPResp, nil
}

type ViaCEPRequest struct {
	Cep         string `json:"cep"`
	Longradouro string `json:"longradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type AddressResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func BuscaViaCepApi(cep string) (AddressResponse, error) {
	return BuscaViaCepApiComURL(cep, "https://viacep.com.br/ws")
}

func BuscaViaCepApiComURL(cep string, baseURL string) (AddressResponse, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: tr}

	url := fmt.Sprintf("%s/%s/json/", baseURL, cep)
	log.Printf("ViaCEP URL: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return AddressResponse{}, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return AddressResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return AddressResponse{}, fmt.Errorf("viacep status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AddressResponse{}, err
	}
	log.Printf("ViaCEP response body: %s", string(body))
	var viaCEP ViaCEPRequest
	if err := json.Unmarshal(body, &viaCEP); err != nil {
		return AddressResponse{}, err
	}
	return AddressResponse{
		Cep:          viaCEP.Cep,
		State:        viaCEP.Uf,
		City:         viaCEP.Localidade,
		Neighborhood: viaCEP.Bairro,
		Street:       viaCEP.Longradouro,
		Service:      "ViaCEP",
	}, nil
}
