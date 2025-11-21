# Markup

## 1. Finalidade
- Fator multiplicador sobre o custo para chegar ao preço de venda.
- Útil em cenários de precificação rápida quando custos estão consolidados.

## 2. Fórmula Matemática Exata
$$Markup = \frac{\text{Preco de Venda}}{\text{Custo}}$$

## 3. Definição de cada variável
- Preco de Venda — decimal (R$); origem: tabela de preços.
- Custo — decimal (R$); origem: custo total do item/serviço.

## 4. Regras de Arredondamento
- Markup em 4 casas (cálculo); exibição em 2 casas.

## 5. Regras de Exceção e Validação
- Custo = 0 → erro.
- Preço negativo → inválido.

## 6. Exemplo Numérico Real (Passo a Passo)
- PV = 150; Custo = 75 → Markup = 2,00.

## 7. Onde essa fórmula é usada no sistema
- Precificação de catálogo; relatórios de preço.

## 8. Notas para Desenvolvedores (Dev Notes)
- Não confundir markup (fator) com margem (%).
- Usar decimal exato para evitar erro de ponto flutuante.
