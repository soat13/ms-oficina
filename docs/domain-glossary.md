## De-Para: Termos em Português → Inglês

### Contextos (Bounded Contexts)

| Português         | Inglês       | Pasta                    | Explicação                                                                     |
|-------------------|--------------|--------------------------|--------------------------------------------------------------------------------|
| Cliente           | Customer     | `internal/customer/`     | Representa o dono do veículo e titular do contrato                             |
| Usuário           | User         | `internal/user/`         | Usuários do sistema (Atendente, Mecânico, Gerente) com credenciais de acesso   |
| Veículo           | Vehicle      | `internal/vehicle/`      | Automóvel que será atendido na oficina                                         |
| Serviço           | Service      | `internal/service/`      | Atividade técnica do catálogo (ex: "Alinhamento", "Troca de Óleo")             |
| Produto           | Product      | `internal/product/`      | Peça ou insumo do catálogo (ex: "Filtro de Cabine", "Óleo 5W30")               |
| Ordem de Serviço  | Repair Order | `internal/repairorder/`  | OS - Autorização formal para executar serviços no veículo                      |
| Orçamento         | Estimate     | `internal/estimate/`     | Proposta de preços com itens de serviço e produtos                             |
| Autenticação      | Auth         | `internal/auth/`         | Sistema de login e geração de tokens JWT                                       |

### Atores (Roles)

| Português  | Inglês (código) | Constante   | Explicação                                                      |
|------------|-----------------|-------------|-----------------------------------------------------------------|
| Cliente    | Customer        | —           | Dono do veículo que solicita serviços e aprova orçamentos       |
| Atendente  | Attendant       | `attendant` | Cadastra clientes, cria OS, gerencia o fluxo de atendimento     |
| Mecânico   | Mechanic        | `mechanic`  | Executa diagnósticos e serviços técnicos                        |
| Gerente    | Manager         | `manager`   | Supervisiona operações e aprova ações administrativas           |

### Estados da Ordem de Serviço (Repair Order Status)

| Português              | Inglês (código)      | Constante              | Explicação                                                     |
|------------------------|----------------------|----------------------  |----------------------------------------------------------------|
| Recebida               | Received             | `received`             | Veículo chegou na oficina                                      |
| Em Diagnóstico         | In Diagnostics       | `in_diagnostics`       | Mecânico está avaliando o problema                             |
| Diagnóstico finalizado | Diagnostics Finished | `diagnostics_finished` | Diagnóstico do mecânico finalizado                             |
| Aguardando Aprovação   | Awaiting Approval    | `awaiting_approval`    | Cliente precisa aprovar o orçamento                            |
| Aprovada               | Approved             | `approved`             | Cliente aprovou, estoque atualizado, pode iniciar execução     |
| Em Execução            | In Execution         | `in_execution`         | Serviços sendo executados                                      |
| Finalizada             | Finished             | `finished`             | Serviços concluídos, aguardando retirada                       |
| Liberada               | Released             | `released`             | Veículo entregue ao cliente                                    |
| Cancelada              | Canceled             | `canceled`             | OS foi cancelada                                               |

### Estados do Orçamento (Estimate Status)

| Português            | Inglês (código)   | Constante           | Explicação                                              |
|----------------------|-------------------|---------------------|---------------------------------------------------------|
| Aguardando Aprovação | Awaiting Approval | `awaiting_approval` | Orçamento criado, aguardando decisão do cliente         |
| Aprovado             | Approved          | `approved`          | Totalmente aprovado e com estoque confirmado            |
| Rejeitado            | Rejected          | `rejected`          | Cliente recusou o orçamento                             |
| Cancelado            | Canceled          | `canceled`          | Orçamento foi cancelado (ex: OS cancelada)              |
