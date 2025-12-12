-- 059 - Rollback FK de compensacoes_bancarias.receita_id para receitas (legado)

ALTER TABLE compensacoes_bancarias
    ADD CONSTRAINT compensacoes_bancarias_receita_id_fkey
        FOREIGN KEY (receita_id) REFERENCES receitas(id) ON DELETE CASCADE;

