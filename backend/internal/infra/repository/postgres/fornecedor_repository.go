package postgres

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/port"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/google/uuid"
)

// FornecedorRepositoryPG implementa FornecedorRepository usando PostgreSQL
type FornecedorRepositoryPG struct {
	queries *db.Queries
}

// NewFornecedorRepositoryPG cria nova instancia do repositorio
func NewFornecedorRepositoryPG(queries *db.Queries) port.FornecedorRepository {
	return &FornecedorRepositoryPG{queries: queries}
}

// Create persiste um novo fornecedor
func (r *FornecedorRepositoryPG) Create(ctx context.Context, fornecedor *entity.Fornecedor) error {
	params := db.CreateFornecedorParams{
		TenantID:            uuidToPgUUID(fornecedor.TenantID),
		RazaoSocial:         fornecedor.RazaoSocial,
		NomeFantasia:        strPtrToPgText(fornecedor.NomeFantasia),
		Cnpj:                fornecedor.CNPJ,
		Email:               strPtrToPgText(fornecedor.Email),
		Telefone:            strPtrToPgText(fornecedor.Telefone),
		Celular:             strPtrToPgText(fornecedor.Celular),
		EnderecoLogradouro:  strPtrToPgText(fornecedor.EnderecoLogradouro),
		EnderecoNumero:      strPtrToPgText(fornecedor.EnderecoNumero),
		EnderecoComplemento: strPtrToPgText(fornecedor.EnderecoComplemento),
		EnderecoBairro:      strPtrToPgText(fornecedor.EnderecoBairro),
		EnderecoCidade:      strPtrToPgText(fornecedor.EnderecoCidade),
		EnderecoEstado:      strPtrToPgText(fornecedor.EnderecoEstado),
		EnderecoCep:         strPtrToPgText(fornecedor.EnderecoCEP),
		Banco:               strPtrToPgText(fornecedor.Banco),
		Agencia:             strPtrToPgText(fornecedor.Agencia),
		Conta:               strPtrToPgText(fornecedor.Conta),
		Observacoes:         strPtrToPgText(fornecedor.Observacoes),
		Ativo:               fornecedor.Ativo,
	}

	created, err := r.queries.CreateFornecedor(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar fornecedor: %w", err)
	}

	fornecedor.ID = pgUUIDToUUID(created.ID)
	fornecedor.CreatedAt = created.CriadoEm.Time
	fornecedor.UpdatedAt = created.AtualizadoEm.Time

	return nil
}

// FindByID busca fornecedor por ID
func (r *FornecedorRepositoryPG) FindByID(ctx context.Context, tenantID, fornecedorID uuid.UUID) (*entity.Fornecedor, error) {
	params := db.GetFornecedorByIDParams{
		ID:       uuidToPgUUID(fornecedorID),
		TenantID: uuidToPgUUID(tenantID),
	}

	result, err := r.queries.GetFornecedorByID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar fornecedor: %w", err)
	}

	return r.toDomain(&result), nil
}

// FindByCNPJ busca fornecedor por CNPJ (nota: query precisa incluir tenant_id)
func (r *FornecedorRepositoryPG) FindByCNPJ(ctx context.Context, cnpj string) (*entity.Fornecedor, error) {
	// TODO: Atualizar query para aceitar tenantID como parametro
	// Por ora, busca sem filtro de tenant (problema de seguranca)
	return nil, fmt.Errorf("FindByCNPJ nao implementado - requer tenant_id")
}

// Update atualiza fornecedor existente
func (r *FornecedorRepositoryPG) Update(ctx context.Context, fornecedor *entity.Fornecedor) error {
	params := db.UpdateFornecedorParams{
		ID:                  uuidToPgUUID(fornecedor.ID),
		TenantID:            uuidToPgUUID(fornecedor.TenantID),
		RazaoSocial:         fornecedor.RazaoSocial,
		NomeFantasia:        strPtrToPgText(fornecedor.NomeFantasia),
		Email:               strPtrToPgText(fornecedor.Email),
		Telefone:            strPtrToPgText(fornecedor.Telefone),
		Celular:             strPtrToPgText(fornecedor.Celular),
		EnderecoLogradouro:  strPtrToPgText(fornecedor.EnderecoLogradouro),
		EnderecoNumero:      strPtrToPgText(fornecedor.EnderecoNumero),
		EnderecoComplemento: strPtrToPgText(fornecedor.EnderecoComplemento),
		EnderecoBairro:      strPtrToPgText(fornecedor.EnderecoBairro),
		EnderecoCidade:      strPtrToPgText(fornecedor.EnderecoCidade),
		EnderecoEstado:      strPtrToPgText(fornecedor.EnderecoEstado),
		EnderecoCep:         strPtrToPgText(fornecedor.EnderecoCEP),
		Banco:               strPtrToPgText(fornecedor.Banco),
		Agencia:             strPtrToPgText(fornecedor.Agencia),
		Conta:               strPtrToPgText(fornecedor.Conta),
		Observacoes:         strPtrToPgText(fornecedor.Observacoes),
	}

	updated, err := r.queries.UpdateFornecedor(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar fornecedor: %w", err)
	}

	fornecedor.UpdatedAt = updated.AtualizadoEm.Time
	return nil
}

// Delete realiza soft delete do fornecedor
func (r *FornecedorRepositoryPG) Delete(ctx context.Context, tenantID, fornecedorID uuid.UUID) error {
	params := db.DeleteFornecedorParams{
		ID:       uuidToPgUUID(fornecedorID),
		TenantID: uuidToPgUUID(tenantID),
	}

	if err := r.queries.DeleteFornecedor(ctx, params); err != nil {
		return fmt.Errorf("erro ao deletar fornecedor: %w", err)
	}

	return nil
}

// ListAll lista todos os fornecedores do tenant
func (r *FornecedorRepositoryPG) ListAll(ctx context.Context, tenantID uuid.UUID) ([]*entity.Fornecedor, error) {
	results, err := r.queries.ListFornecedores(ctx, uuidToPgUUID(tenantID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar fornecedores: %w", err)
	}

	fornecedores := make([]*entity.Fornecedor, len(results))
	for i, result := range results {
		fornecedores[i] = r.toDomain(&result)
	}

	return fornecedores, nil
}

// ListAtivos lista apenas fornecedores ativos do tenant
func (r *FornecedorRepositoryPG) ListAtivos(ctx context.Context, tenantID uuid.UUID) ([]*entity.Fornecedor, error) {
	results, err := r.queries.ListFornecedoresAtivos(ctx, uuidToPgUUID(tenantID))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar fornecedores ativos: %w", err)
	}

	fornecedores := make([]*entity.Fornecedor, len(results))
	for i, result := range results {
		fornecedores[i] = r.toDomain(&result)
	}

	return fornecedores, nil
}

// toDomain converte modelo do sqlc para entidade de dominio
func (r *FornecedorRepositoryPG) toDomain(f *db.Fornecedore) *entity.Fornecedor {
	return &entity.Fornecedor{
		ID:                  pgUUIDToUUID(f.ID),
		TenantID:            pgUUIDToUUID(f.TenantID),
		RazaoSocial:         f.RazaoSocial,
		NomeFantasia:        pgTextToStr(f.NomeFantasia),
		CNPJ:                f.Cnpj,
		Email:               pgTextToStr(f.Email),
		Telefone:            pgTextToStr(f.Telefone),
		Celular:             pgTextToStr(f.Celular),
		EnderecoLogradouro:  pgTextToStr(f.EnderecoLogradouro),
		EnderecoNumero:      pgTextToStr(f.EnderecoNumero),
		EnderecoComplemento: pgTextToStr(f.EnderecoComplemento),
		EnderecoBairro:      pgTextToStr(f.EnderecoBairro),
		EnderecoCidade:      pgTextToStr(f.EnderecoCidade),
		EnderecoEstado:      pgTextToStr(f.EnderecoEstado),
		EnderecoCEP:         pgTextToStr(f.EnderecoCep),
		Banco:               pgTextToStr(f.Banco),
		Agencia:             pgTextToStr(f.Agencia),
		Conta:               pgTextToStr(f.Conta),
		Observacoes:         pgTextToStr(f.Observacoes),
		Ativo:               f.Ativo,
		CreatedAt:           f.CriadoEm.Time,
		UpdatedAt:           f.AtualizadoEm.Time,
	}
}
