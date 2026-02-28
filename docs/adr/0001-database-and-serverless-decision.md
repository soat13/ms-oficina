# ADR 0001 — Escolha do Banco de Dados Relacional e Impacto da Adoção de Serverless

## Status
Accepted

## Date
2026-02-28

## Authors
Marcos Souza  
Lucas Prioli  
Samuel Gonçalves

---

## 1. Contexto

O sistema **Oficina API** foi concebido para suportar o fluxo operacional completo de uma oficina mecânica, contemplando:

- Cadastro de clientes e veículos
- Criação e gestão de Ordens de Serviço (OS)
- Geração e aprovação de orçamentos
- Controle de itens (produtos e serviços)
- Controle de acesso por perfis

O domínio apresenta forte interdependência entre entidades, com relacionamentos bem definidos e invariantes de negócio que exigem consistência transacional.

Inicialmente, adotou-se o **PostgreSQL** como banco de dados relacional. Posteriormente, com a inclusão da execução em ambiente **serverless (AWS Lambda)**, tornou-se necessário reavaliar a decisão arquitetural quanto ao modelo de persistência, considerando impactos operacionais e de escalabilidade.

---

## 2. Problema Arquitetural

A adoção de funções serverless altera o modelo de execução da aplicação, introduzindo:

- Escalabilidade automática e concorrência variável
- Conexões efêmeras ao banco de dados
- Possibilidade de picos de criação de conexões
- Retries automáticos e execução idempotente
- Cold start e variabilidade de latência

Diante desse novo cenário, surge a questão:

> A arquitetura de persistência relacional permanece adequada sob o modelo serverless?

---

## 3. Decisão

Foi decidido:

1. Manter o **PostgreSQL** como banco principal.
2. Manter o **modelo relacional existente**, conforme documentado no DER.
3. Endereçar as preocupações introduzidas pelo ambiente serverless no nível operacional, e não estrutural.

---

## 4. Justificativa Técnica

### 4.1 Adequação ao Modelo de Domínio

O domínio da aplicação apresenta características típicas de sistemas transacionais fortemente relacionais:

- Dependência de integridade referencial entre agregados.
- Necessidade de atomicidade em operações compostas.
- Regras de consistência envolvendo múltiplas entidades.
- Uso extensivo de constraints (FK, UNIQUE, CHECK).

Exemplos críticos incluem:

- Aprovação de orçamento e consequente alteração de estado da Ordem de Serviço.
- Controle de estoque de produtos.
- Garantia de unicidade de documentos e e-mails.
- Manutenção de relacionamentos cliente → veículo → OS → orçamento → itens.

A utilização de um banco relacional com suporte a **ACID**, constraints declarativas e isolamento transacional reduz significativamente o risco de inconsistência e simplifica a implementação das invariantes de negócio.

A adoção de solução não relacional implicaria:

- Desnormalização de dados.
- Gestão manual de integridade referencial.
- Complexidade adicional em transações distribuídas.
- Aumento da responsabilidade da aplicação na manutenção da consistência.

Dado o perfil do domínio, tais trade-offs não se mostraram vantajosos.

---

### 4.2 Impacto da Adoção de Serverless

A arquitetura serverless impacta principalmente aspectos operacionais:

- Gerenciamento de conexões com o banco.
- Controle de concorrência sob auto-scaling.
- Necessidade de handlers idempotentes.
- Estratégias de retry e timeout.

Entretanto, esses desafios são mitigáveis por meio de:

- Estratégias de pooling ou uso de proxy de conexão.
- Controle de concorrência via configuração de Lambda.
- Design idempotente dos casos de uso.
- Monitoramento e observabilidade adequados.

Importante destacar que tais aspectos não alteram a natureza do modelo de dados nem invalidam os benefícios do modelo relacional para o domínio em questão.

---

### 4.3 Manutenção do Mesmo Banco

Optou-se por manter o mesmo banco (PostgreSQL) após a adoção de serverless com base nos seguintes critérios:

- Coerência arquitetural.
- Redução de complexidade cognitiva.
- Preservação das garantias transacionais.
- Evitar introdução de múltiplos paradigmas de persistência sem justificativa clara.
- Minimização de risco técnico e operacional.

A mudança do modelo de execução (container/Kubernetes → serverless) não implica, por si só, necessidade de mudança no modelo de persistência.

---

## 5. Consequências

### 5.1 Consequências Positivas

- Forte integridade de dados garantida pelo banco.
- Clareza estrutural do domínio.
- Simplificação das regras de negócio.
- Redução da complexidade arquitetural.

### 5.2 Consequências Negativas / Trade-offs

- Necessidade de atenção ao gerenciamento de conexões em ambiente serverless.
- Potencial necessidade de proxy ou limitação de concorrência.
- Dependência de banco relacional para operações críticas.

---

## 6. Alternativas Consideradas

### 6.1 Migração para Banco NoSQL

Avaliada a possibilidade de adoção de modelo NoSQL orientado a documentos.  
Foi descartada devido à:

- Alta interdependência entre entidades.
- Necessidade de integridade referencial.
- Complexidade transacional.
- Ausência de requisito de escala global massiva.

### 6.2 Separação de Banco Exclusivo para Funções Serverless

Considerada a hipótese de utilizar banco distinto para funções Lambda.  
Descartada por introduzir:

- Complexidade adicional.
- Problemas de sincronização.
- Risco de inconsistência.
- Sobrecarga operacional desnecessária.

---

## 7. Conclusão

A decisão de manter PostgreSQL e o modelo relacional existente permanece tecnicamente adequada, mesmo após a adoção de execução serverless.

O ambiente de execução impacta aspectos operacionais, mas não altera os fundamentos estruturais do domínio, que continua sendo fortemente relacional e transacional.

A decisão está alinhada com os princípios de:

- Coerência arquitetural
- Minimização de complexidade desnecessária
- Adequação da tecnologia ao domínio
- Engenharia orientada a trade-offs conscientes

---

## Referências

- DER: `docs/der.md`
- Arquitetura: `docs/architecture.md`