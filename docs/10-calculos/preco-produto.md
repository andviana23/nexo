# Cálculo do Preço de Produto (Financeiro Completo)

## 1. Finalidade
- Definir preço de venda de produto considerando custo, impostos e taxa de adquirência, respeitando margem desejada.

## 2. Fórmula Matemática Exata
$$\text{Preco Final} = \frac{\text{Custo do Produto} + \text{Impostos} + \text{Taxa da Maquininha}}{1 - \text{Margem Desejada}}$$

## 3. Definição de cada variável
- Custo do Produto — decimal (R$); origem: custo unitário médio/última compra.
- Impostos — decimal (R$); origem: tributação aplicável ao item.
- Taxa da Maquininha — decimal (R$); origem: % de adquirência * preço ou valor nominal. **INFORMACAO AUSENTE — CONFIRMACAO NECESSARIA** sobre se é percentual do PV (iterativo) ou valor fixo.
- Margem Desejada — decimal (0–1); origem: configuração de pricing do negócio.

## 4. Regras de Arredondamento
- Cálculo interno em 4 casas; preço exibido com 2 casas (round half-up).

## 5. Regras de Exceção e Validação
- Margem Desejada >= 1 → inválido (divisão por zero/negativo).
- Valores negativos → rejeitar.
- Se Taxa da Maquininha for % do PV, iterar até convergência (não especificado).

## 6. Exemplo Numérico Real (Passo a Passo)
- Custo = 30; Impostos = 5; Taxa = 2; Margem = 0,40.  
- Preco Final = (30 + 5 + 2) / (1 – 0,40) = 37 / 0,60 = 61,67.

## 7. Onde essa fórmula é usada no sistema
- Catálogo de produtos; definição de preço no POS.
- Relatórios de rentabilidade por produto.

## 8. Notas para Desenvolvedores (Dev Notes)
- Documentar claramente se Taxa da Maquininha é % do preço ou valor fixo; caso % → usar iteração ou equação fechada.
- Usar decimal monetário para evitar erro de ponto flutuante.
