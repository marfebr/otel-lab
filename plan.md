# Plano de Desenvolvimento - Sistema de Temperatura por CEP com OTEL

## ✅ Fase 1: Análise e Preparação
1. **✅ Analisar o código base existente**
   - ✅ Examinar o código em `base/comunicacao-ms/` para entender a estrutura atual
   - ✅ Identificar componentes reutilizáveis
   - ✅ Mapear dependências e configurações existentes

2. **✅ Definir arquitetura do sistema**
   - ✅ Serviço A: API Gateway/Input Handler
   - ✅ Serviço B: Orquestrador (baseado no código existente)
   - ✅ Comunicação HTTP entre serviços
   - ✅ Integração com OTEL e Zipkin
   - [[./architecture.md]] 

## 🔄 Fase 2: Desenvolvimento do Serviço A
1. **✅ Criar estrutura do projeto**
   - ✅ Configurar `go.mod` e dependências
   - ✅ Estrutura de pastas: `cmd/`, `internal/`, `pkg/`
   - ✅ Dockerfile e docker-compose

2. **✅ Implementar validação de CEP**
   - ✅ Endpoint POST para receber CEP
   - ✅ Validação de formato (8 dígitos, string)
   - ✅ Retorno de erro 422 para CEPs inválidos

3. **✅ Implementar comunicação com Serviço B**
   - ✅ Cliente HTTP para chamar Serviço B
   - ✅ Tratamento de respostas e erros
   - ✅ Propagação de status codes

## 🔄 Fase 3: Adaptação do Serviço B
1. **✅ Analisar código existente**
   - ✅ Identificar funcionalidades de busca de CEP
   - ✅ Identificar integração com serviço de clima
   - ✅ Mapear endpoints e estruturas de resposta
   - [[./service-b-analysis.md]]

2. **✅ Implementar novos requisitos**
   - ✅ Validação de CEP (formato 8 dígitos)
   - ✅ Tratamento de erros específicos (422, 404)
   - ✅ Formatação de resposta padronizada
   - ✅ Conversão de temperaturas (Celsius, Fahrenheit, Kelvin)
   - ✅ Integração ViaCEP
   - ✅ Integração OpenWeatherMap
   - ✅ Orquestração e endpoint POST /weather

## 🔄 Fase 4: Implementação do OTEL + Zipkin
1. **✅ Configurar OpenTelemetry**
   - ✅ Adicionar serviço Zipkin ao docker-compose
   - ✅ Configurar exportador Zipkin no OTEL Collector
   - ✅ Documentar instrução de acesso ao Zipkin
   - ✅ Garantir integração dos serviços com OTEL Collector

2. **✅ Implementar tracing distribuído**
   - ✅ Instrumentar Serviço A com spans
   - ✅ Instrumentar Serviço B com spans
   - ✅ Configurar propagação de contexto entre serviços
   - ✅ Medir tempo de resposta das APIs externas
   - ✅ Documentar hierarquia de spans implementados

3. **✅ Configurar Zipkin**
   - ✅ Adicionar Zipkin ao docker-compose
   - ✅ Configurar endpoint de exportação
   - ✅ Testar visualização de traces
   - ✅ Documentar instruções de acesso e uso

## 🔄 Fase 5: Integração e Testes
1. **✅ Configurar comunicação entre serviços**
   - ✅ Atualizar docker-compose para incluir ambos os serviços
   - ✅ Configurar networking entre containers
   - ✅ Definir variáveis de ambiente

2. **✅ Implementar testes**
   - ✅ Testes unitários para validação de CEP
   - ✅ Testes de integração entre serviços
   - ✅ Testes de cenários de erro
   - ✅ Testes de tracing

3. **✅ Testar propagação de contexto OTEL (tracing distribuído)**
   - ✅ Gerar requisição completa Service A → Service B
   - ✅ Validar no Zipkin a presença de spans de ambos os serviços
   - ✅ Conferir hierarquia e detalhes dos spans

4. **✅ Documentação e validação**
   - ✅ Documentar endpoints, exemplos de uso e troubleshooting
   - ✅ Validar todos os requisitos do README

5. **✅ Refinamento e Otimização**
   - ✅ Revisar implementação
   - ✅ Verificar conformidade com requisitos
   - ✅ Otimizar performance
   - ✅ Melhorar tratamento de erros
   - ✅ Organizar estrutura final do projeto
   - ✅ Preparar instruções de execução
   - ✅ Teste final de ponta a ponta

---


```