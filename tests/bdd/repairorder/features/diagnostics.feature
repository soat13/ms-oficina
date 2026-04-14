# language: pt
Funcionalidade: Diagnóstico da Ordem de Reparo
  Como mecânico da oficina
  Quero iniciar e finalizar o diagnóstico de um veículo
  Para que um orçamento seja gerado automaticamente para o cliente aprovar

  Contexto:
    Dado que existe um produto "Filtro de Óleo" com estoque 10
    E existe um produto "Pastilha de Freio" com estoque 5
    E existe um serviço "Troca de Óleo"
    E existe um serviço "Revisão de Freios"

  Cenário: Mecânico finaliza o diagnóstico com sucesso e um orçamento é gerado
    Dado que existe uma Ordem de Reparo no status "received"
    Quando eu inicio o diagnóstico
    Então a OS deve estar no status "in_diagnostics"
    Quando eu finalizo o diagnóstico com:
      | tipo    | nome              | quantidade |
      | produto | Filtro de Óleo    | 2          |
      | produto | Pastilha de Freio | 1          |
      | servico | Troca de Óleo     | 1          |
      | servico | Revisão de Freios | 1          |
    Então a OS deve estar no status "awaiting_approval"
    E um orçamento deve ter sido criado para a OS no status "awaiting_approval"
    E o orçamento deve conter 4 itens

  Cenário: Não é possível iniciar diagnóstico fora do status received
    Dado que existe uma Ordem de Reparo no status "in_execution"
    Quando eu inicio o diagnóstico
    Então a resposta HTTP deve ser 409
    E a OS deve permanecer no status "in_execution"

  Cenário: Finalização falha por estoque insuficiente
    Dado que existe uma Ordem de Reparo no status "in_diagnostics"
    Quando eu finalizo o diagnóstico com:
      | tipo    | nome              | quantidade |
      | produto | Pastilha de Freio | 99         |
      | servico | Revisão de Freios | 1          |
    Então a resposta HTTP deve ser 422
    E a OS deve permanecer no status "in_diagnostics"
