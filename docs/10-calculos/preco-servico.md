# Cálculo do Preço de Serviço (Financeiro Completo)

## 1. Finalidade
- Definir preço de venda do serviço considerando insumos, comissão, impostos e taxas, visando margem desejada.

## 2. Fórmula Matemática Exata
$$\text{Preco Final} = \frac{\text{Custo de Insumos} + \text{Comissao} + \text{Impostos} + \text{Taxas}}{1 - \text{Margem Desejada}}$$

## 3. Definição de cada variável
- Custo de Insumos — decimal (R$); origem: consumo de insumos do serviço (ver custo-insumo-servico.md).
- Comissao — decimal (R$); origem: política de comissão do barbeiro/profissional.
- Impostos — decimal (R$); origem: tributação do serviço.
- Taxas — decimal (R$); origem: adquirência/gateway/outros custos diretos.
- Margem Desejada — decimal (0–1); origem: configuração de pricing.

## 4. Regras de Arredondamento
- Cálculo interno em 4 casas; preço exibido com 2 casas (round half-up).

## 5. Regras de Exceção e Validação
- Margem Desejada >= 1 → inválido.
- Valores negativos → rejeitar.
- Se Comissão ou Taxas forem % do preço, é necessário iteração/convergência (não especificado). **INFORMACAO AUSENTE — CONFIRMACAO NECESSARIA**

## 6. Exemplo Numérico Real (Passo a Passo)
- Insumos = 10; Comissao = 20; Impostos = 3; Taxas = 2; Margem = 0,35.  
  Preco Final = (10 + 20 + 3 + 2) / (1 – 0,35) = 35 / 0,65 = 53,85.

## 7. Onde essa fórmula é usada no sistema
- Tabela de serviços e orçamentos; dashboards de rentabilidade de serviços.

## 8. Notas para Desenvolvedores (Dev Notes)
- Definir claramente se comissão/taxas são baseadas no PV; se sim, usar equação com PV em ambos lados ou método iterativo.
- Garantir consumo de estoque integrado (para custo de insumos).
