# Ponto de Equilíbrio (PE)

## 1. Finalidade
- Determinar o faturamento mínimo para cobrir todas as despesas fixas, dado o mix de custos variáveis e margem de contribuição.
- Apoiar metas mensais e alertas de viabilidade financeira.

## 2. Fórmula Matemática Exata
- Fórmula principal:  
  $$PE = \frac{\text{Despesas Fixas}}{\text{Margem de Contribuicao}}$$
- Auxiliar (margem):  
  $$\text{Margem de Contribuicao} = 1 - \frac{\text{Custos Variaveis}}{\text{Faturamento}}$$

## 3. Definição de cada variável
- Despesas Fixas — decimal (R$); origem: despesas marcadas como fixas na tabela de despesas ou input manual.
- Custos Variaveis — decimal (R$); origem: despesas variáveis/COGS do período.
- Faturamento — decimal (R$); origem: receitas brutas do período.
- Margem de Contribuicao — decimal (0–1); derivada das variáveis acima.

## 4. Regras de Arredondamento
- Margem: 4 casas decimais.
- PE: round half-up com 2 casas (R$).

## 5. Regras de Exceção e Validação
- Faturamento = 0 → erro “divisao por zero”.
- Custos Variaveis > Faturamento → margem negativa; marcar como “margem negativa/impraticavel”.
- Despesas Fixas ausentes → considerar 0.

## 6. Exemplo Numérico Real (Passo a Passo)
- Despesas Fixas = 20.000; Faturamento = 50.000; Custos Variaveis = 20.000.
- Margem = 1 – (20.000 / 50.000) = 0,60.
- PE = 20.000 / 0,60 = 33.333,33.

## 7. Onde essa fórmula é usada no sistema
- Módulo Financeiro → Dashboard de saúde financeira e relatórios de viabilidade.
- Endpoint/resumo financeiro mensal.

## 8. Notas para Desenvolvedores (Dev Notes)
- Calcular sempre com dados do mesmo período e timezone.
- Garantir tipagem monetária (decimal/currency) e não double binário.
- Armazenar PE calculado ou em visão/cached para dashboards; recalcular após novas despesas/receitas.
