# Curva ABC

## 1. Finalidade
- Classificar itens (produtos/serviços/clientes) por representatividade no valor total, para priorização de gestão e estoque.

## 2. Fórmula Matemática Exata
$$\%\text{Representatividade} = \frac{\text{Valor Total do Item}}{\text{Valor Total Geral}}$$

## 3. Definição de cada variável
- Valor Total do Item — decimal (R$); origem: faturamento ou valor movimentado pelo item.
- Valor Total Geral — decimal (R$); origem: soma de todos os itens do conjunto analisado.

## 4. Regras de Arredondamento
- Percentual com 2 casas decimais.

## 5. Regras de Exceção e Validação
- Total Geral = 0 → erro.
- Valores negativos → inválidos.

## 6. Exemplo Numérico Real (Passo a Passo)
- Item = 5.000; Total = 50.000 → 10%.

## 7. Onde essa fórmula é usada no sistema
- Relatórios de estoque/vendas; segmentação A/B/C (cortes típicos: A até ~80%, B até ~95%, C restante).

## 8. Notas para Desenvolvedores (Dev Notes)
- Ordenar itens desc, acumular % e marcar faixas A/B/C; cortes exatos devem ser configuráveis se necessário.
- Usar período e base de valor consistentes para todos os itens.
