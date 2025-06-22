# Análise do Código Base - Serviço B

## Estrutura Atual dos Microserviços

### Arquitetura Existente
```
goapp (porta 8080) → goapp2 (porta 8181) → goapp3 (porta 8282)
```

### Decisão: Adaptar goapp2 como Serviço B
- **Porta**: 8181 (compatível com arquitetura definida)
- **Posição**: Meio da cadeia (orquestrador)
- **Configuração**: Já possui chamadas externas

## Funcionalidades Existentes

### ✅ Componentes Reutilizáveis
1. **Estrutura OTEL Completa**
   - Provider, tracer, propagator configurados
   - Integração com OTEL Collector
   - Spans para medir tempo de resposta
   - Propagação de contexto entre serviços

2. **Chi Router com Middleware**
   - RequestID, RealIP, Recoverer, Logger
   - Timeout de 60 segundos
   - Endpoint `/metrics` para Prometheus

3. **Configuração Viper**
   - Variáveis de ambiente configuráveis
   - Valores padrão definidos
   - Configuração automática

4. **Graceful Shutdown**
   - Tratamento de sinais (SIGINT)
   - Timeout de shutdown
   - Limpeza de recursos

5. **Docker Setup**
   - Multi-stage build
   - Configuração de ambiente
   - Integração com OTEL Collector

### 🔄 Funcionalidades a Adaptar
1. **Chamadas Externas**
   - Atual: Chamada simples para goapp3
   - Necessário: Integração com ViaCEP + OpenWeatherMap

2. **Template HTML**
   - Atual: Renderização de template
   - Necessário: Endpoints JSON REST

3. **Estruturas de Dados**
   - Atual: TemplateData para HTML
   - Necessário: Request/Response JSON

## Lacunas Identificadas

### ❌ Funcionalidades Ausentes
1. **Integração com ViaCEP**
   - Não existe busca de CEP
   - Não há estruturas para endereços

2. **Integração com OpenWeatherMap**
   - Não existe busca de clima
   - Não há estruturas para temperaturas

3. **Validação de CEP**
   - Não existe validação de formato
   - Não há tratamento de erros específicos

4. **Conversão de Temperaturas**
   - Não existe conversão Celsius/Fahrenheit/Kelvin

5. **Endpoints REST**
   - Atual: Apenas GET / (HTML)
   - Necessário: POST /weather (JSON)

## Configurações Existentes

### Variáveis de Ambiente (goapp2)
```bash
TITLE=Microservice Demo 2
BACKGROUND_COLOR=blue
EXTERNAL_CALL_URL=http://goapp3:8282
EXTERNAL_CALL_METHOD=GET
RESPONSE_TIME=2000
REQUEST_NAME_OTEL=microservice-demo2-request
OTEL_SERVICE_NAME=microservice-demo2
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317
HTTP_PORT=:8181
```

### Dependências (go.mod)
- ✅ Chi Router (v5.2.1)
- ✅ Viper (v1.20.1)
- ✅ OpenTelemetry (v1.36.0)
- ✅ Prometheus (v1.22.0)
- ✅ gRPC (v1.73.0)

## Plano de Adaptação

### 1. Manter (Reutilizar)
- ✅ Estrutura OTEL completa
- ✅ Chi Router e middleware
- ✅ Configuração Viper
- ✅ Graceful shutdown
- ✅ Docker setup

### 2. Remover
- 🔄 Template HTML (não necessário)
- 🔄 Chamada para goapp3
- 🔄 TemplateData struct

### 3. Adicionar
- 🔄 Endpoint POST /weather
- 🔄 Integração ViaCEP
- 🔄 Integração OpenWeatherMap
- 🔄 Validação de CEP
- 🔄 Conversão de temperaturas
- 🔄 Estruturas JSON

### 4. Modificar
- 🔄 Handler principal
- 🔄 Variáveis de ambiente
- 🔄 Estruturas de dados

## Próximos Passos

1. **Criar estrutura do service-b** baseada no código base
2. **Implementar integração ViaCEP**
3. **Implementar integração OpenWeatherMap**
4. **Implementar conversão de temperaturas**
5. **Criar endpoint POST /weather**
6. **Configurar variáveis de ambiente**
7. **Testar integração completa** 