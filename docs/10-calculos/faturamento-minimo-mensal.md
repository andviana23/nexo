# Faturamento Mínimo Mensal

## 1. Finalidade
- Definir a meta mínima de receita mensal para cobrir despesas e comissões previstas.

## 2. Fórmula Matemática Exata
$$\text{Faturamento Minimo} = \text{Despesas Fixas} + \text{Despesas Variaveis} + \text{Comissoes Previstas}$$

## 3. Definição de cada variável
- Despesas Fixas — decimal (R$); origem: despesas marcadas como fixas.
- Despesas Variaveis — decimal (R$); origem: despesas variáveis/COGS.
- Comissoes Previstas — decimal (R$); origem: projeção de comissões no período.

## 4. Regras de Arredondamento
- 2 casas decimais.

## 5. Regras de Exceção e Validação
- Valores nulos → tratar como 0.
- Comissao como % do faturamento → requer iteração; **INFORMACAO AUSENTE — CONFIRMACAO NECESSARIA**.

## 6. Exemplo Numérico Real (Passo a Passo)
- Fixas = 20.000; Variaveis = 8.000; Comissoes = 5.000 → Minimo = 33.000.

## 7. Onde essa fórmula é usada no sistema
- Planejamento mensal; alertas de atingimento de meta mínima.

## 8. Notas para Desenvolvedores (Dev Notes)
- Definir se projeção de comissão depende do faturamento; se sim, resolver equação simultânea.
- Usar período mensal fechado e timezone consistente.
