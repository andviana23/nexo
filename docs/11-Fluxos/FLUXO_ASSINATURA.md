# Fluxo de Assinatura — NEXO

```
[Início]
   ↓
[Usuário clica em "Assinaturas"]
   ↓
[Clica em "Nova Assinatura"]
   ↓
[Seleciona Cliente]
   ↓
Cliente existe no sistema?
   → Não → [Criar cliente no sistema] → continua
   → Sim  → continua
   ↓
Cliente existe no Asaas?
   → Não → [Criar cliente no Asaas]
   → Sim  → [Recuperar customerId]
   ↓
[Selecionar Plano de Assinatura]
   ↓
[Definir Dia de Vencimento]
   ↓
[Criar Assinatura via API Asaas]
   ↓
[Receber subscriptionId]
   ↓
[Gerar Link de Pagamento]
   ↓
[Exibir link para Recepção copiar]
   ↓
[Salvar assinatura no banco local]
   ↓
Status inicial = “Pendente de Pagamento”
   ↓
[Fim]
```
