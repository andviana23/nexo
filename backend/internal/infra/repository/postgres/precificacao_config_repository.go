// Package postgres implementa os repositórios usando PostgreSQL e sqlc.
package postgres

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
)

// PrecificacaoConfigRepository implementa port.PrecificacaoConfigRepository usando sqlc.
type PrecificacaoConfigRepository struct {
	queries *db.Queries
}

// NewPrecificacaoConfigRepository cria uma nova instância do repositório.
func NewPrecificacaoConfigRepository(queries *db.Queries) *PrecificacaoConfigRepository {
	return &PrecificacaoConfigRepository{
		queries: queries,
	}
}

// Create persiste uma nova configuração de precificação.
func (r *PrecificacaoConfigRepository) Create(ctx context.Context, config *entity.PrecificacaoConfig) error {
	tenantUUID := uuidStringToPgtype(config.TenantID)

	params := db.CreatePrecificacaoConfigParams{
		TenantID:                  tenantUUID,
		MargemDesejada:            percentageToNumeric(config.MargemDesejada),
		MarkupAlvo:                decimalToNumeric(config.MarkupAlvo),
		ImpostoPercentual:         percentageToNumeric(config.ImpostoPercentual),
		ComissaoPercentualDefault: percentageToNumeric(config.ComissaoPercentualDefault),
	}

	result, err := r.queries.CreatePrecificacaoConfig(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar configuração de precificação: %w", err)
	}

	config.ID = pgUUIDToString(result.ID)
	config.CriadoEm = timestamptzToTime(result.CriadoEm)
	config.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)

	return nil
}

// FindByID busca uma configuração por ID.
func (r *PrecificacaoConfigRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.PrecificacaoConfig, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	params := db.GetPrecificacaoConfigByIDParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	}

	model, err := r.queries.GetPrecificacaoConfigByID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar configuração: %w", err)
	}

	return r.toDomain(&model)
}

// FindByTenantID busca a configuração de um tenant.
func (r *PrecificacaoConfigRepository) FindByTenantID(ctx context.Context, tenantID string) (*entity.PrecificacaoConfig, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	model, err := r.queries.GetPrecificacaoConfigByTenant(ctx, tenantUUID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar configuração do tenant: %w", err)
	}

	return r.toDomain(&model)
}

// Update atualiza uma configuração.
func (r *PrecificacaoConfigRepository) Update(ctx context.Context, config *entity.PrecificacaoConfig) error {
	tenantUUID := uuidStringToPgtype(config.TenantID)
	idUUID := uuidStringToPgtype(config.ID)

	params := db.UpdatePrecificacaoConfigParams{
		ID:                        idUUID,
		TenantID:                  tenantUUID,
		MargemDesejada:            percentageToNumeric(config.MargemDesejada),
		MarkupAlvo:                decimalToNumeric(config.MarkupAlvo),
		ImpostoPercentual:         percentageToNumeric(config.ImpostoPercentual),
		ComissaoPercentualDefault: percentageToNumeric(config.ComissaoPercentualDefault),
	}

	result, err := r.queries.UpdatePrecificacaoConfig(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar configuração: %w", err)
	}

	config.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)
	return nil
}

// Delete remove uma configuração (por tenant - uma única config por tenant).
func (r *PrecificacaoConfigRepository) Delete(ctx context.Context, tenantID string) error {
	tenantUUID := uuidStringToPgtype(tenantID)

	// Primeiro buscar a config para obter o ID
	config, err := r.FindByTenantID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("erro ao buscar configuração para deletar: %w", err)
	}

	idUUID := uuidStringToPgtype(config.ID)

	params := db.DeletePrecificacaoConfigParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	}

	err = r.queries.DeletePrecificacaoConfig(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao deletar configuração: %w", err)
	}

	return nil
}

// toDomain converte modelo sqlc para entidade de domínio.
func (r *PrecificacaoConfigRepository) toDomain(model *db.PrecificacaoConfig) (*entity.PrecificacaoConfig, error) {
	margemDesejada, err := numericToPercentage(model.MargemDesejada)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter margem_desejada: %w", err)
	}

	impostoPercentual, err := numericToPercentage(model.ImpostoPercentual)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter imposto_percentual: %w", err)
	}

	comissaoDefault, err := numericToPercentage(model.ComissaoPercentualDefault)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter comissao_percentual_default: %w", err)
	}

	config := &entity.PrecificacaoConfig{
		ID:                        pgUUIDToString(model.ID),
		TenantID:                  pgUUIDToString(model.TenantID),
		MargemDesejada:            margemDesejada,
		MarkupAlvo:                numericToDecimal(model.MarkupAlvo),
		ImpostoPercentual:         impostoPercentual,
		ComissaoPercentualDefault: comissaoDefault,
		CriadoEm:                  timestamptzToTime(model.CriadoEm),
		AtualizadoEm:              timestamptzToTime(model.AtualizadoEm),
	}

	return config, nil
}
