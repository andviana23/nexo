-- 059 - Ajuste FK de compensacoes_bancarias.receita_id
-- Motivo: novos fluxos usam contas_a_receber como receita canônica;
-- a FK antiga para receitas (legado) bloqueia criação automática de compensações D+.

ALTER TABLE compensacoes_bancarias
    DROP CONSTRAINT IF EXISTS compensacoes_bancarias_receita_id_fkey;

