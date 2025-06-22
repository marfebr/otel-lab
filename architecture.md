# Arquitetura do Sistema - Temperatura por CEP com OTEL

## Análise Comparativa: Atual vs. Requisitos

### Arquitetura Atual (Código Base)
```
Cliente → goapp (porta 8080) → goapp2 (porta 8181) → goapp3 (porta 8282)
   ↑                                                              ↓
   └─────────────── Resposta HTML ←──────────────────────────────┘
```

### Arquitetura Requerida
```
Cliente → Serviço A (validação) → Serviço B (orquestração) → APIs Externas
   ↑                                                              ↓
   └─────────────── Resposta JSON com clima ←────────────────────┘
```

## Definição da Arquitetura

### Serviço A - API Gateway/Input Handler
**Baseado no `goapp` atual**

**Responsabilidades:**
- Receber requisições POST com CEP
- Validar formato do CEP (8 dígitos, string)
- Retornar erro 422 para CEPs inválidos
- Encaminhar CEPs válidos para Serviço B
- Propagar contexto OTEL para tracing distribuído

**Endpoints:**
- `POST /cep` - Receber e validar CEP

**Estrutura de Request:**
```json
{
  "cep": "29902555"
}
```

**Estrutura de Response (sucesso):**
```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

**Estrutura de Response (erro):**
```json
{
  "error": "invalid zipcode"
}
```

### Serviço B - Orquestrador
**Baseado no `goapp2` atual**

**Responsabilidades:**
- Receber CEP válido do Serviço A
- Buscar endereço via API ViaCEP
- Extrair nome da cidade
- Buscar dados de clima via API OpenWeatherMap
- Converter temperaturas (Celsius, Fahrenheit, Kelvin)
- Retornar dados formatados
- Tratamento de erros específicos

**Endpoints:**
- `POST /weather` - Processar CEP e retornar clima

**Integrações Externas:**
- **ViaCEP API**: `https://viacep.com.br/ws/{cep}/json/`
- **OpenWeatherMap API**: `https://api.openweathermap.org/data/2.5/weather`

**Códigos de Erro:**
- `422` - CEP inválido (formato incorreto)
- `404` - CEP não encontrado

## Fluxo de Dados

```
1. Cliente → POST /cep → Serviço A
2. Serviço A → Valida CEP
   ├─ CEP inválido → Retorna 422
   └─ CEP válido → Continua
3. Serviço A → POST /weather → Serviço B
4. Serviço B → GET ViaCEP API
   ├─ CEP não encontrado → Retorna 404
   └─ CEP encontrado → Continua
5. Serviço B → GET OpenWeatherMap API
6. Serviço B → Converte temperaturas
7. Serviço B → Retorna dados → Serviço A → Cliente
```

## Observabilidade (OTEL + Zipkin)

### Tracing Distribuído
- **Serviço A**: Span para validação de CEP
- **Serviço B**: Span para busca de CEP + Span para busca de clima
- **Propagação**: Contexto OTEL entre serviços
- **Medição**: Tempo de resposta das APIs externas

### Configuração
- **OTEL Collector**: Processamento centralizado
- **Zipkin**: Visualização de traces (substitui Jaeger)
- **Prometheus**: Métricas dos serviços
- **Grafana**: Dashboards (opcional)

## Estrutura de Dados

### Request/Response Padrões
```go
// Request CEP
type CEPRequest struct {
    CEP string `json:"cep"`
}

// Response Clima
type WeatherResponse struct {
    City   string  `json:"city"`
    TempC  float64 `json:"temp_C"`
    TempF  float64 `json:"temp_F"`
    TempK  float64 `json:"temp_K"`
}

// Response Erro
type ErrorResponse struct {
    Error string `json:"error"`
}
```

### Conversão de Temperaturas
- **Celsius → Fahrenheit**: `°F = (°C × 9/5) + 32`
- **Celsius → Kelvin**: `K = °C + 273.15`

## Configuração de Ambiente

### Variáveis de Ambiente
```bash
# Serviço A
HTTP_PORT=:8080
SERVICE_B_URL=http://service-b:8181
OTEL_SERVICE_NAME=service-a
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317

# Serviço B
HTTP_PORT=:8181
VIACEP_BASE_URL=https://viacep.com.br/ws
OPENWEATHER_API_KEY=your_api_key
OPENWEATHER_BASE_URL=https://api.openweathermap.org/data/2.5
OTEL_SERVICE_NAME=service-b
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317
```

## Adaptações Necessárias do Código Base

### Componentes Reutilizáveis
- ✅ Estrutura OTEL (provider, tracer, propagator)
- ✅ Chi Router com middleware
- ✅ Configuração Viper
- ✅ Docker setup
- ✅ Graceful shutdown

### Modificações Necessárias
- 🔄 Remover template HTML (não necessário)
- 🔄 Adicionar endpoints específicos
- 🔄 Implementar validação de CEP
- 🔄 Integrar APIs externas
- 🔄 Implementar conversão de temperaturas
- 🔄 Substituir Jaeger por Zipkin 