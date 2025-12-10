-- Adicionar vínculo entre Sangria e Contas a Pagar
ALTER TABLE operacoes_caixa 
ADD COLUMN conta_pagar_id UUID REFERENCES contas_a_pagar(id);

CREATE INDEX idx_operacoes_caixa_conta_pagar ON operacoes_caixa(conta_pagar_id);

-- Renomear tabela de assinaturas da plataforma para evitar confusão com subscriptions (clientes)
ALTER TABLE assinaturas RENAME TO nexo_assinaturas;

-- Garantir índices do DRE para updates rápidos
CREATE INDEX IF NOT EXISTS idx_dre_mensal_tenant_mes_ano ON dre_mensal(tenant_id, mes_ano);
