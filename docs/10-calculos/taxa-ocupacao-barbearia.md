# Taxa de Ocupação da Barbearia

## 1. Finalidade
- Medir a utilização global da capacidade de atendimento (horas) da barbearia.
- Indicar ociosidade ou saturação da agenda.

## 2. Fórmula Matemática Exata
$$\text{Taxa de Ocupacao} = \frac{\text{Horas Totais de Servicos Realizados}}{\text{Horas Totais Disponiveis}}$$

## 3. Definição de cada variável
- Horas Totais de Servicos Realizados — decimal (horas); origem: soma das duracoes de servicos concluídos no período.
- Horas Totais Disponiveis — decimal (horas); origem: capacidade da agenda (slots configurados) no período.

## 4. Regras de Arredondamento
- Percentual com 2 casas decimais.

## 5. Regras de Exceção e Validação
- Horas Totais Disponiveis = 0 → retornar “agenda nao configurada” ou erro.
- Duracoes faltantes → tratar como 0.

## 6. Exemplo Numérico Real (Passo a Passo)
- Realizados = 320h; Disponiveis = 400h.
- Taxa = 320 / 400 = 0,80 (80%).

## 7. Onde essa fórmula é usada no sistema
- Dashboard operacional; relatório de produtividade global.

## 8. Notas para Desenvolvedores (Dev Notes)
- Mesmo timezone e janela (dia/semana/mes) para ambos numerador e denominador.
- Excluir no-shows se não houver execução real.
