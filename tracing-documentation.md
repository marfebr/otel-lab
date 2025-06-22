# Documentação de Tracing Distribuído

## Visão Geral
O sistema implementa tracing distribuído usando OpenTelemetry (OTEL) com visualização no Zipkin. Todos os serviços propagam contexto OTEL para manter a continuidade dos traces.

## Hierarquia de Spans

### Serviço A (service-a)
```
cep-validation
├── service-b-weather-request (chamada para Serviço B)
```

### Serviço B (service-b)
```
weather-request (handler HTTP)
├── weather-orchestration (orquestrador)
    ├── viacep-request (API ViaCEP)
    └── openweather-request (API OpenWeatherMap)
```

## Detalhamento dos Spans

### Serviço A - Spans

#### `cep-validation`
- **Localização**: `service-a/internal/handler/cep_handler.go`
- **Propósito**: Mede o tempo total de processamento da requisição no Serviço A
- **Inclui**: Validação de CEP, chamada para Serviço B, formatação de resposta
- **Erros registrados**: CEP inválido, erros de JSON, erros do Serviço B

#### `service-b-weather-request`
- **Localização**: `service-a/internal/service/service_b_client.go`
- **Propósito**: Mede o tempo da chamada HTTP para o Serviço B
- **Inclui**: Serialização de request, chamada HTTP, deserialização de response
- **Propagação**: Contexto OTEL injetado nos headers HTTP
- **Erros registrados**: Falhas de rede, erros de serialização, status codes de erro

### Serviço B - Spans

#### `weather-request`
- **Localização**: `service-b/internal/handler/weather_handler.go`
- **Propósito**: Mede o tempo de processamento da requisição HTTP no Serviço B
- **Inclui**: Deserialização de request, orquestração, formatação de resposta
- **Erros registrados**: CEP inválido, CEP não encontrado, erros internos

#### `weather-orchestration`
- **Localização**: `service-b/internal/service/weather_orchestrator.go`
- **Propósito**: Mede o tempo total da orquestração de serviços
- **Inclui**: Validação, busca no ViaCEP, busca no OpenWeatherMap, conversão de temperaturas
- **Erros registrados**: Erros de validação, falhas nas APIs externas

#### `viacep-request`
- **Localização**: `service-b/internal/service/viacep_client.go`
- **Propósito**: Mede o tempo da chamada para a API ViaCEP
- **Inclui**: Construção de URL, chamada HTTP, processamento de resposta
- **Erros registrados**: CEP não encontrado, falhas de rede, erros de parsing

#### `openweather-request`
- **Localização**: `service-b/internal/service/weather_client.go`
- **Propósito**: Mede o tempo da chamada para a API OpenWeatherMap
- **Inclui**: Construção de URL com parâmetros, chamada HTTP, processamento de resposta
- **Erros registrados**: Falhas de rede, status codes de erro, erros de parsing

## Propagação de Contexto

### Entre Serviços
- **Serviço A → Serviço B**: Contexto OTEL propagado via headers HTTP
- **Método**: `otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))`
- **Resultado**: Mesmo trace ID mantido entre serviços

### Dentro dos Serviços
- **Contexto**: Passado através de todas as funções
- **Spans**: Criados como filhos do span pai
- **Resultado**: Hierarquia completa de spans visível no Zipkin

## Visualização no Zipkin

### Acesso
- **URL**: http://localhost:9411
- **Método**: Após subir `docker-compose up`

### Informações Visíveis
- **Trace ID**: Identificador único do trace completo
- **Span ID**: Identificador de cada operação
- **Tempo**: Duração de cada span
- **Hierarquia**: Relacionamento pai-filho entre spans
- **Erros**: Spans marcados como erro quando falham
- **Tags**: Informações adicionais (CEP, cidade, etc.)

## Exemplo de Trace Completo

```
1. Cliente faz POST /cep com CEP "29902555"
2. Serviço A: span "cep-validation" inicia
3. Serviço A: span "service-b-weather-request" inicia
4. Serviço B: span "weather-request" inicia (mesmo trace ID)
5. Serviço B: span "weather-orchestration" inicia
6. Serviço B: span "viacep-request" inicia
7. ViaCEP retorna dados de Vitória/ES
8. Serviço B: span "openweather-request" inicia
9. OpenWeatherMap retorna temperatura em Celsius
10. Serviço B: converte temperaturas
11. Todos os spans finalizam em ordem reversa
12. Cliente recebe resposta com dados de clima
```

## Benefícios do Tracing

1. **Observabilidade**: Visibilidade completa do fluxo de dados
2. **Debugging**: Identificação rápida de gargalos e falhas
3. **Performance**: Medição de tempo de cada componente
4. **Distribuído**: Rastreamento através de múltiplos serviços
5. **Erros**: Identificação precisa de onde falhas ocorrem 