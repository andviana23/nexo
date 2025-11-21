# Custo de Insumo por Serviço

## 1. Finalidade
- Calcular o custo total dos insumos consumidos em uma execução de serviço.
- Base para precificação e baixa de estoque.

## 2. Fórmula Matemática Exata
$$\text{Custo de Insumo} = \sum (\text{Quantidade Consumida} \times \text{Custo Unitario})$$

## 3. Definição de cada variável
- Quantidade Consumida — decimal (unidade do insumo); origem: ficha técnica/consumo real.
- Custo Unitario — decimal (R$/un); origem: custo médio/última compra do insumo.

## 4. Regras de Arredondamento
- Custo unitário em 4 casas; resultado em 2 casas (round half-up).

## 5. Regras de Exceção e Validação
- Quantidade ou custo ausentes → tratar como 0 ou bloquear, conforme política.
- Quantidade negativa → inválido.

## 6. Exemplo Numérico Real (Passo a Passo)
- Pomada: 0,2 un × 50 = 10; Pó: 0,1 un × 30 = 3 → Total = 13,00.

## 7. Onde essa fórmula é usada no sistema
- Precificação de serviços; baixa de estoque; relatórios de rentabilidade por serviço.

## 8. Notas para Desenvolvedores (Dev Notes)
- Integrar com estoque para registrar consumo (manual ou automático).
- Garantir unidade de medida consistente na ficha técnica e no estoque.
