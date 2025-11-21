# CAC (Custo de Aquisição de Cliente)

## 1. Finalidade
- Medir o custo médio para adquirir cada novo cliente em um período.
- Comparar com LTV para avaliar eficiência de marketing e payback.

## 2. Fórmula Matemática Exata
$$CAC = \frac{\text{Investimento em Marketing}}{\text{Numero de Novos Clientes}}$$

## 3. Definição de cada variável
- Investimento em Marketing — decimal (R$); origem: despesas de marketing no período.
- Numero de Novos Clientes — inteiro; origem: contagem de clientes criados/ativados no período.

## 4. Regras de Arredondamento
- Resultado: round half-up, 2 casas (R$).

## 5. Regras de Exceção e Validação
- Numero de Novos Clientes = 0 → retorno “nao definido” ou erro 422.
- Investimento negativo → invalidar.

## 6. Exemplo Numérico Real (Passo a Passo)
- Investimento = 5.000; Novos Clientes = 50.
- CAC = 5.000 / 50 = 100,00.

## 7. Onde essa fórmula é usada no sistema
- Dashboard de marketing/financeiro; relatórios de aquisição.
- Comparação CAC vs LTV.

## 8. Notas para Desenvolvedores (Dev Notes)
- Filtrar despesas de marketing e clientes pelo mesmo período e timezone.
- Considerar clientes únicos; evitar reativações contadas como “novos” se não for regra de negócio.
