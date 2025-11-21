# Ticket Médio (TM)

## 1. Finalidade
- Medir a receita média por cliente atendido no período.
- Insumo direto para o cálculo de LTV.

## 2. Fórmula Matemática Exata
$$TM = \frac{\text{Faturamento}}{\text{Numero de Clientes}}$$

## 3. Definição de cada variável
- Faturamento — decimal (R$); origem: receitas no período.
- Numero de Clientes — inteiro; origem: clientes únicos atendidos no período.

## 4. Regras de Arredondamento
- 2 casas decimais (round half-up).

## 5. Regras de Exceção e Validação
- Numero de Clientes = 0 → erro/“sem clientes”.
- Faturamento negativo → inválido.

## 6. Exemplo Numérico Real (Passo a Passo)
- Faturamento = 50.000; Clientes = 400 → TM = 125,00.

## 7. Onde essa fórmula é usada no sistema
- Dash de receita; insumo para LTV.

## 8. Notas para Desenvolvedores (Dev Notes)
- Clientes devem ser únicos no período (distintos).
- Usar timezone e mesmo período para numerador/denominador.
