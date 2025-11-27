-- Migration: Reverter alteração do tamanho da coluna CPF
-- Descrição: Volta o tamanho da coluna cpf para 11 caracteres (apenas CPF)

-- Reverter tipo da coluna cpf na tabela profissionais
ALTER TABLE profissionais 
  ALTER COLUMN cpf TYPE VARCHAR(11);

-- Restaurar comentário original
COMMENT ON COLUMN profissionais.cpf IS 'CPF sem formatação (11 dígitos)';
