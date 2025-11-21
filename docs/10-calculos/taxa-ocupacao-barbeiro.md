# Taxa de Ocupação por Barbeiro

## 1. Finalidade
- Medir produtividade individual de cada barbeiro em relação à agenda disponível.
- Apoiar gestão de equipe e comissionamento variável.

## 2. Fórmula Matemática Exata
$$\text{Ocupacao do Barbeiro} = \frac{\text{Horas Trabalhadas}}{\text{Horas de Agenda Disponiveis}}$$

## 3. Definição de cada variável
- Horas Trabalhadas — decimal (horas); origem: soma das duracoes de servicos concluídos pelo barbeiro no período.
- Horas de Agenda Disponiveis — decimal (horas); origem: slots de agenda configurados para o barbeiro no período.

## 4. Regras de Arredondamento
- Percentual com 2 casas decimais.

## 5. Regras de Exceção e Validação
- Disponiveis = 0 → retornar “sem agenda configurada”.
- Serviços cancelados ou no-show não entram nas horas trabalhadas.

## 6. Exemplo Numérico Real (Passo a Passo)
- Trabalhadas = 60h; Disponiveis = 80h.
- Ocupação = 60 / 80 = 0,75 (75%).

## 7. Onde essa fórmula é usada no sistema
- Relatórios por colaborador; dashboards de equipe; cálculo de performance.

## 8. Notas para Desenvolvedores (Dev Notes)
- Usar mesma granularidade temporal da agenda (dia/semana/mes).
- Efetuar soma convertendo minutos para horas com precisão decimal.
