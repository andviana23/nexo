// Package postgres implementa os repositórios usando PostgreSQL e sqlc.
package postgres

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
)

// MetaMensalRepository implementa port.MetaMensalRepository usando sqlc.
type MetaMensalRepository struct {
	queries *db.Queries
}

// NewMetaMensalRepository cria uma nova instância do repositório.
func NewMetaMensalRepository(queries *db.Queries) *MetaMensalRepository {
	return &MetaMensalRepository{
		queries: queries,
	}
}

// Create persiste uma nova meta mensal.
func (r *MetaMensalRepository) Create(ctx context.Context, meta *entity.MetaMensal) error {
	tenantUUID := uuidStringToPgtype(meta.TenantID)
	// CriadoPor vem do DB, não da entidade - usar UUID vazio ou extrair de contexto
	criadoPorUUID := uuidStringToPgtype("00000000-0000-0000-0000-000000000000")
	origemStr := meta.Origem.String()
	statusStr := meta.Status
	params := db.CreateMetaMensalParams{
		TenantID:        tenantUUID,
		MesAno:          meta.MesAno.String(),
		MetaFaturamento: moneyToRawDecimal(meta.MetaFaturamento),
		Origem:          &origemStr,
		Status:          &statusStr,
		CriadoPor:       criadoPorUUID,
	}
	result, err := r.queries.CreateMetaMensal(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar meta mensal: %w", err)
	}
	meta.ID = pgUUIDToString(result.ID)
	meta.CriadoEm = timestamptzToTime(result.CriadoEm)
	meta.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)
	return nil
}

// FindByID busca uma meta por ID.
func (r *MetaMensalRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.MetaMensal, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)
	result, err := r.queries.GetMetaMensalByID(ctx, db.GetMetaMensalByIDParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar meta: %w", err)
	}
	return r.toDomain(&result)
}

// FindByMesAno busca meta de um mês específico.
func (r *MetaMensalRepository) FindByMesAno(ctx context.Context, tenantID string, mesAno valueobject.MesAno) (*entity.MetaMensal, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	result, err := r.queries.GetMetaMensalByMesAno(ctx, db.GetMetaMensalByMesAnoParams{
		TenantID: tenantUUID,
		MesAno:   mesAno.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar meta por mês: %w", err)
	}
	return r.toDomain(&result)
}

// Update atualiza uma meta existente.
func (r *MetaMensalRepository) Update(ctx context.Context, meta *entity.MetaMensal) error {
	tenantUUID := uuidStringToPgtype(meta.TenantID)
	idUUID := uuidStringToPgtype(meta.ID)
	params := db.UpdateMetaMensalParams{
		ID:              idUUID,
		TenantID:        tenantUUID,
		MetaFaturamento: moneyToRawDecimal(meta.MetaFaturamento),
		Status:          &meta.Status,
	}
	result, err := r.queries.UpdateMetaMensal(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar meta: %w", err)
	}
	meta.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)
	return nil
}

// Delete remove uma meta.
func (r *MetaMensalRepository) Delete(ctx context.Context, tenantID, id string) error {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)
	err := r.queries.DeleteMetaMensal(ctx, db.DeleteMetaMensalParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	})
	if err != nil {
		return fmt.Errorf("erro ao deletar meta: %w", err)
	}
	return nil
}

// ListAtivas lista metas ativas.
func (r *MetaMensalRepository) ListAtivas(ctx context.Context, tenantID string) ([]*entity.MetaMensal, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	results, err := r.queries.ListMetasMensaisByTenant(ctx, db.ListMetasMensaisByTenantParams{
		TenantID: tenantUUID,
		Limit:    1000,
		Offset:   0,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar metas: %w", err)
	}
	return r.toDomainList(results)
}

// ListByPeriod lista metas em um período.
func (r *MetaMensalRepository) ListByPeriod(ctx context.Context, tenantID string, inicio, fim valueobject.MesAno) ([]*entity.MetaMensal, error) {
	// Usar ListAtivas e filtrar em memória (ou adicionar query específica no futuro)
	metas, err := r.ListAtivas(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	var filtered []*entity.MetaMensal
	for _, meta := range metas {
		if !meta.MesAno.Before(inicio) && !meta.MesAno.After(fim) {
			filtered = append(filtered, meta)
		}
	}
	return filtered, nil
}

// toDomain converte modelo sqlc para entidade de domínio.
func (r *MetaMensalRepository) toDomain(model *db.MetasMensai) (*entity.MetaMensal, error) {
	// Parse string "2025-01" diretamente com NewMesAno
	mesAno, err := valueobject.NewMesAno(model.MesAno)
	if err != nil {
		return nil, fmt.Errorf("criar MesAno: %w", err)
	}
	var origem valueobject.OrigemMeta
	if model.Origem != nil {
		origem = valueobject.OrigemMeta(*model.Origem)
	} else {
		origem = valueobject.OrigemMetaManual
	}
	var status string
	if model.Status != nil {
		status = *model.Status
	} else {
		status = "ATIVA"
	}
	meta := &entity.MetaMensal{
		ID:              pgUUIDToString(model.ID),
		TenantID:        pgUUIDToString(model.TenantID),
		MesAno:          mesAno,
		MetaFaturamento: rawDecimalToMoney(model.MetaFaturamento),
		Origem:          origem,
		Status:          status,
		CriadoEm:        timestamptzToTime(model.CriadoEm),
		AtualizadoEm:    timestamptzToTime(model.AtualizadoEm),
	}
	return meta, nil
}

// toDomainList converte lista de modelos para lista de entidades.
func (r *MetaMensalRepository) toDomainList(models []db.MetasMensai) ([]*entity.MetaMensal, error) {
	result := make([]*entity.MetaMensal, 0, len(models))
	for i := range models {
		meta, err := r.toDomain(&models[i])
		if err != nil {
			return nil, err
		}
		result = append(result, meta)
	}
	return result, nil
}
