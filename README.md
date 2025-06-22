# Lab OTEL

## Objetivo: Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin) juntamente com a cidade. Esse sistema deverá implementar OTEL(Open Telemetry) e Zipkin, com propagação de trace distribuído entre os serviços.

1. Baseado no cenário conhecido "Sistema de temperatura por CEP" denominado Serviço B, será incluso um novo projeto, denominado Serviço A.

2. Requisitos - Serviço A (responsável pelo input):

    O sistema deve receber um input de 8 dígitos via POST, através do schema:  { "cep": "29902555" }
    O sistema deve validar se o input é valido (contem 8 dígitos) e é uma STRING
        Caso seja válido, buscar a cidade correspondente no ViaCEP
        Caso não seja válido, deve retornar:
            Código HTTP: 422
            Mensagem: invalid zipcode
    Se o CEP for encontrado, encaminhar para o Serviço B via HTTP o nome da cidade:
        { "city": "Vitória" }
    Se o CEP não for encontrado:
        Código HTTP: 404
        Mensagem: can not find zipcode

3. Requisitos - Serviço B (responsável pela orquestração):

    O sistema deve receber o nome da cidade
    O sistema deve realizar a pesquisa do clima no OpenWeatherMap e retornar as temperaturas formatadas em: Celsius, Fahrenheit, Kelvin juntamente com o nome da cidade.
    O sistema deve responder adequadamente nos seguintes cenários:
        Em caso de sucesso:
            Código HTTP: 200
            Response Body: { "city": "Vitória", "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.65 }
        Em caso de falha, caso a cidade não seja encontrada:
            Código HTTP: 404
            Mensagem: can not find city

4. Após a implementação dos serviços, adicione a implementação do OTEL + Zipkin:

    Implementar tracing distribuído entre Serviço A - Serviço B
    Utilizar span para medir o tempo de resposta do serviço de busca de CEP e busca de temperatura

Dicas:

    use o codigo do diretorio "base/comunicacao-ms" como base para a funcionalidade

## Propagação de Trace OTEL

O sistema implementa tracing distribuído usando OpenTelemetry (OTEL) e Zipkin. O trace é propagado automaticamente do **service-a** para o **service-b** e de volta, permitindo rastrear toda a jornada da requisição, desde o recebimento do CEP até a resposta final com o clima.

- O **service-a** injeta o contexto OTEL nos headers HTTP ao chamar o service-b.
- O **service-b** extrai o contexto OTEL dos headers HTTP no início do handler, garantindo a continuidade do trace.
- Todos os spans (validação, requisições externas, orquestração) são encadeados e visualizáveis no Zipkin.

### Exemplo de visualização de trace no Zipkin
1. Acesse http://localhost:9411/zipkin/
2. Clique em "Find traces" para ver as requisições recentes.
3. Clique em um trace para ver a hierarquia de spans, por exemplo:
   - `POST /cep` (service-a)
     - `service-b-weather-request` (service-a → service-b)
       - `weather-request` (service-b handler)
         - `viacep-request` (consulta ViaCEP)
         - `weather-orchestration` (orquestração do clima)
         - `weatherapi-request` (consulta WeatherAPI)

Assim, é possível acompanhar toda a cadeia de chamadas e identificar gargalos ou falhas.

## Como Executar

### Pré-requisitos
- Docker e Docker Compose instalados
- API Key do OpenWeatherMap (gratuita em https://openweathermap.org/api)

### Configuração
1. Configure a variável de ambiente para a API do OpenWeatherMap:
```bash
export OPENWEATHER_API_KEY=sua_api_key_aqui
```

### Execução
1. Suba todos os serviços:
```bash
docker-compose up --build
```

2. Aguarde todos os serviços estarem prontos (pode levar alguns minutos na primeira execução)

### Testando o Sistema

#### Exemplo de requisição válida para o Serviço A:
```bash
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "70636240"}'
```
Resposta:
```json
{
  "city": "Brasília",
  "temp_C": 20.2,
  "temp_F": 68.4,
  "temp_K": 293.35
}
```

#### Exemplo de requisição com CEP inválido:
```bash
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "123"}'
```
Resposta:
```json
HTTP/1.1 422 Unprocessable Entity
{"error":"invalid zipcode"}
```

#### Exemplo de requisição com CEP não encontrado:
```bash
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "00000000"}'
```
Resposta:
```json
HTTP/1.1 404 Not Found
{"error":"can not find zipcode"}
```

#### Exemplo de requisição direta para o Serviço B:
```bash
curl -X POST http://localhost:8181/weather \
  -H "Content-Type: application/json" \
  -d '{"cep": "70636240"}'
```
Resposta:
```json
{
  "city": "Brasília",
  "temp_C": 20.2,
  "temp_F": 68.4,
  "temp_K": 293.35
}
```

## Observabilidade: Visualizando Traces

Após subir o ambiente com `docker-compose up`, acesse o Zipkin para visualizar os traces distribuídos:

- URL: http://localhost:9411/zipkin/

No Zipkin, você poderá acompanhar o tracing distribuído entre os serviços, com todos os spans encadeados.

### Como usar o Zipkin:
1. Acesse http://localhost:9411/zipkin/
2. Clique em "Find traces" para ver todos os traces
3. Use os filtros para buscar por serviço específico
4. Clique em um trace para ver a hierarquia de spans
5. Analise o tempo de cada operação e identifique gargalos

### Endpoints disponíveis:
- **Serviço A**: http://localhost:8080/cep (POST, recebe CEP)
- **Serviço B**: http://localhost:8181/weather (POST, recebe CEP)
- **Métricas**: http://localhost:9090 (Prometheus)
- **Zipkin**: http://localhost:9411/zipkin/ (Traces)

### Estrutura do Projeto:
```
Otel-lab/
├── service-a/          # Serviço A - Validação de CEP e busca de cidade
├── service-b/          # Serviço B - Orquestração e clima
├── base/              # Código base original
├── docker-compose.yaml # Orquestração completa
└── README.md          # Este arquivo
```

## Resultados dos Testes Finais

- **CEP válido:**
  - Status: 200
  - Resposta: cidade e clima retornados corretamente.

- **CEP inválido (menos de 8 dígitos):**
  - Status: 422
  - Resposta: {"error":"invalid zipcode"}

- **CEP válido mas inexistente:**
  - Status: 404
  - Resposta: {"error":"can not find zipcode"}

O sistema está robusto, com tracing distribuído OTEL funcionando ponta a ponta, e responde corretamente para todos os cenários principais e de erro.
