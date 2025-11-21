# Margem de Lucro

## 1. Finalidade
- Medir rentabilidade percentual de um item/serviço.
- Usada em relatórios de lucratividade e formação de preço.

## 2. Fórmula Matemática Exata
$$\text{Margem} = \frac{\text{Lucro}}{\text{Preco de Venda}}$$

## 3. Definição de cada variável
- Lucro — decimal (R$); origem: Preco de Venda – Custo Total (insumos/taxas/comissão).
- Preco de Venda — decimal (R$); origem: tabela de preços do item/serviço.

## 4. Regras de Arredondamento
- Percentual com 2 casas; Lucro em 2 casas.

## 5. Regras de Exceção e Validação
- Preco de Venda = 0 → erro.
- Lucro negativo → margem negativa (mostrar alerta).

## 6. Exemplo Numérico Real (Passo a Passo)
- PV = 100,00; Custo Total = 70,00 → Lucro = 30,00 → Margem = 0,30 (30%).

## 7. Onde essa fórmula é usada no sistema
- Relatórios financeiros e comparativo de itens; suporte a KPIs de rentabilidade.

## 8. Notas para Desenvolvedores (Dev Notes)
- Garantir custo total consistente com configurações de preço (ver preco-produto.md e preco-servico.md).
- Evitar confundir margem (%) com markup (fator).
