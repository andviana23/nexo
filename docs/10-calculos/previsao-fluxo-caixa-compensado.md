# Previsão de Fluxo de Caixa Compensado

## 1. Finalidade
- Estimar a data de entrada de valores considerando o prazo de compensação (cartão/boleto) para projeção de caixa.

## 2. Fórmula Matemática Exata
$$\text{Data de Entrada} = \text{Data da Venda} + \text{Prazo de Compensacao}$$

## 3. Definição de cada variável
- Data da Venda — data; origem: registro da transação.
- Prazo de Compensacao — inteiro (dias); origem: regra do meio de pagamento (cartão/boleto/PIX). **INFORMACAO AUSENTE — CONFIRMACAO NECESSARIA** para prazos padrões.

## 4. Regras de Arredondamento
- Não se aplica a valores monetários; somar dias em calendário.

## 5. Regras de Exceção e Validação
- Prazo ausente → usar padrão por meio de pagamento (não especificado).
- Ajustes por feriados/banco → não especificados.

## 6. Exemplo Numérico Real (Passo a Passo)
- Venda em 10/03; Prazo = 30 dias.
- Data de Entrada = 09/04.

## 7. Onde essa fórmula é usada no sistema
- Projeção de caixa; conciliação de recebíveis.

## 8. Notas para Desenvolvedores (Dev Notes)
- Parametrizar prazos por bandeira/meio; considerar D+N e se conta fins de semana.
- Datas em timezone do negócio.
