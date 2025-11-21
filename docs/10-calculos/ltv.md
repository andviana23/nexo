# LTV (Lifetime Value)

## 1. Finalidade
- Estimar o valor total gerado por um cliente ao longo do ciclo de vida.
- Suportar decisões de investimento (CAC vs LTV) e segmentação.

## 2. Fórmula Matemática Exata
$$LTV = \text{Ticket Medio} \times \text{Frequencia Mensal} \times \text{Tempo de Retencao (meses)}$$

## 3. Definição de cada variável
- Ticket Medio — decimal (R$); origem: faturamento/numero de clientes no período (ver ticket-medio.md).
- Frequencia Mensal — decimal (transacoes por cliente por mes); origem: contagem de visitas/servicos por cliente / meses.
- Tempo de Retencao (meses) — decimal (meses); origem: dados históricos ou churn.

## 4. Regras de Arredondamento
- Ticket e Frequencia: 2 casas; Retencao: 2 casas; LTV final: 2 casas.

## 5. Regras de Exceção e Validação
- Qualquer variavel nula ou zero → LTV = 0.
- Retencao negativa → erro.

## 6. Exemplo Numérico Real (Passo a Passo)
- Ticket Medio = 120; Frequencia = 2,5; Retencao = 12.
- LTV = 120 × 2,5 × 12 = 3.600,00.

## 7. Onde essa fórmula é usada no sistema
- Dashboard financeiro; relatórios de coorte.
- Comparação com CAC para payback.

## 8. Notas para Desenvolvedores (Dev Notes)
- Garantir que Ticket Medio e Frequencia sejam calculados no mesmo período-base.
- Usar timezone consistente; evitar double binário para moeda.
