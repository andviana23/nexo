-- Migration: Alterar tamanho da coluna CPF para suportar CNPJ
-- Descrição: Aumenta o tamanho da coluna cpf de 11 para 14 caracteres
--             para permitir armazenar CPF (11 dígitos) ou CNPJ (14 dígitos)

-- Alterar tipo da coluna cpf na tabela profissionais
ALTER TABLE profissionais 
  ALTER COLUMN cpf TYPE VARCHAR(14);

-- Adicionar comentário explicativo
COMMENT ON COLUMN profissionais.cpf IS 'CPF (11 dígitos) ou CNPJ (14 dígitos) sem formatação';
