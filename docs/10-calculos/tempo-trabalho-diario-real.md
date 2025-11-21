# Tempo de Trabalho Diário Real

## 1. Finalidade
- Quantificar o esforço diário entregue em serviços realizados.
- Suportar produtividade e timesheets.

## 2. Fórmula Matemática Exata
$$\text{Tempo Real} = \sum (\text{Duracao de Cada Servico Comandado})$$

## 3. Definição de cada variável
- Duracao de Cada Servico — decimal (minutos ou horas); origem: duração registrada na agenda/OS para cada serviço concluído no dia.

## 4. Regras de Arredondamento
- Somar em minutos; apresentar em horas com 2 casas (round half-up).

## 5. Regras de Exceção e Validação
- Serviços sem duração → considerar 0 ou rejeitar conforme regra de negócio.
- Converter unidades (min → h) antes de exibir.

## 6. Exemplo Numérico Real (Passo a Passo)
- Serviços: 30min, 45min, 60min.
- Soma = 135min = 2,25h.

## 7. Onde essa fórmula é usada no sistema
- Dashboard diário de operação; relatórios de time tracking.

## 8. Notas para Desenvolvedores (Dev Notes)
- Garantir timezone do dia de referência.
- Excluir serviços cancelados ou não concluídos.
