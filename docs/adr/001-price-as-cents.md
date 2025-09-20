
# ADR-001 — Armazenar valores monetários em **centavos** (inteiros)

**Data:** 2025-09-20  
**Status:** Aceita

## Contexto
Precisamos representar valores monetários (preços de produtos/serviços, totais de orçamentos). Tipos de ponto flutuante (`float/double`) introduzem imprecisão binária e erros de arredondamento cumulativos (ex.: `0.1 + 0.2 != 0.3`), o que é crítico em finanças.

## Decisão
Armazenar **todos os valores monetários em centavos (inteiros)** no banco e nas trocas internas entre serviços. A camada de apresentação/DTOs pode exibir em unidade humana (ex.: “R$ 12,34”).

- Exemplos:
    - `R$ 10,00` → `1000`
    - `R$ 25,75` → `2575`

Criar/usar um tipo utilitário `Money` para conversões, soma, subtração, multiplicação por inteiros, arredondamento e formatação.

## Alternativas consideradas
1. **`DECIMAL(precision, scale)` no banco e `string/decimal` na aplicação**
    - ✅ Precisão exata; ❌ Complexidade maior e risco de conversões inconsistentes entre camadas/linguagens.
2. **`float/double`**
    - ✅ Simples; ❌ Impreciso para finanças.

## Consequências

**Positivas**
- Precisão determinística (operações inteiras).
- Comparações e agregações simples.
- Padronização entre serviços.

**Negativas**
- Conversão para exibição (centavos ↔ unidades) obrigatória.
- Atenção ao arredondamento em operações não-inteiras (ex.: rateio, impostos).

## Regras de uso
- Toda persistência interna em **centavos**.
- Toda matemática de domínio em **centavos**.
- Arredondamentos centralizados no tipo `Money`.
- Não misturar `float` em cálculos; apenas na **formatação** para exibição.
