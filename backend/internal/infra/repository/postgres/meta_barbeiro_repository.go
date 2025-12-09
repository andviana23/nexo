// Package postgres implementa os repositórios usando PostgreSQL e sqlc.
package postgres

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
)

// MetaBarbeiroRepository implementa port.MetaBarbeiroRepository usando sqlc.
type MetaBarbeiroRepository struct {
	queries *db.Queries
}

// NewMetaBarbeiroRepository cria uma nova instância do repositório.
func NewMetaBarbeiroRepository(queries *db.Queries) *MetaBarbeiroRepository {
	return &MetaBarbeiroRepository{
		queries: queries,
	}
}

// Create persiste uma nova meta de barbeiro.
func (r *MetaBarbeiroRepository) Create(ctx context.Context, meta *entity.MetaBarbeiro) error {
	tenantUUID := entityUUIDToPgtype(meta.TenantID)
	barbeiroUUID := uuidStringToPgtype(meta.BarbeiroID)

	params := db.CreateMetaBarbeiroParams{
		TenantID:           tenantUUID,
		BarbeiroID:         barbeiroUUID,
		MesAno:             meta.MesAno.String(),
		MetaServicosGerais: moneyToNumeric(meta.MetaServicosGerais),
		MetaServicosExtras: moneyToNumeric(meta.MetaServicosExtras),
		MetaProdutos:       moneyToNumeric(meta.MetaProdutos),
	}

	result, err := r.queries.CreateMetaBarbeiro(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar meta de barbeiro: %w", err)
	}

	meta.ID = pgUUIDToString(result.ID)
	meta.CriadoEm = timestamptzToTime(result.CriadoEm)
	meta.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)

	return nil
}

// FindByID busca uma meta por ID.
func (r *MetaBarbeiroRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.MetaBarbeiro, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	result, err := r.queries.GetMetaBarbeiroByID(ctx, db.GetMetaBarbeiroByIDParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar meta: %w", err)
	}

	return r.toDomain(&result)
}

// FindByBarbeiroMesAno busca meta de um barbeiro em um mês.
func (r *MetaBarbeiroRepository) FindByBarbeiroMesAno(ctx context.Context, tenantID, barbeiroID string, mesAno valueobject.MesAno) (*entity.MetaBarbeiro, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	barbeiroUUID := uuidStringToPgtype(barbeiroID)

	result, err := r.queries.GetMetaBarbeiroByMesAno(ctx, db.GetMetaBarbeiroByMesAnoParams{
		TenantID:   tenantUUID,
		BarbeiroID: barbeiroUUID,
		MesAno:     mesAno.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar meta por barbeiro e mês: %w", err)
	}

	return r.toDomain(&result)
}

// Update atualiza uma meta existente.
func (r *MetaBarbeiroRepository) Update(ctx context.Context, meta *entity.MetaBarbeiro) error {
	tenantUUID := entityUUIDToPgtype(meta.TenantID)
	idUUID := uuidStringToPgtype(meta.ID)

	params := db.UpdateMetaBarbeiroParams{
		ID:                 idUUID,
		TenantID:           tenantUUID,
		MetaServicosGerais: moneyToNumeric(meta.MetaServicosGerais),
		MetaServicosExtras: moneyToNumeric(meta.MetaServicosExtras),
		MetaProdutos:       moneyToNumeric(meta.MetaProdutos),
	}

	result, err := r.queries.UpdateMetaBarbeiro(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar meta: %w", err)
	}

	meta.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)
	return nil
}

// Delete remove uma meta.
func (r *MetaBarbeiroRepository) Delete(ctx context.Context, tenantID, id string) error {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	err := r.queries.DeleteMetaBarbeiro(ctx, db.DeleteMetaBarbeiroParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	})
	if err != nil {
		return fmt.Errorf("erro ao deletar meta: %w", err)
	}

	return nil
}

// ListByBarbeiro lista todas as metas de um barbeiro.
func (r *MetaBarbeiroRepository) ListByBarbeiro(ctx context.Context, tenantID, barbeiroID string) ([]*entity.MetaBarbeiro, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	barbeiroUUID := uuidStringToPgtype(barbeiroID)

	results, err := r.queries.ListMetasBarbeiroByBarbeiro(ctx, db.ListMetasBarbeiroByBarbeiroParams{
		TenantID:   tenantUUID,
		BarbeiroID: barbeiroUUID,
		Limit:      1000,
		Offset:     0,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar metas do barbeiro: %w", err)
	}

	return r.toDomainList(results)
}

// ListByMesAno lista todas as metas de barbeiros de um mês.
func (r *MetaBarbeiroRepository) ListByMesAno(ctx context.Context, tenantID string, mesAno valueobject.MesAno) ([]*entity.MetaBarbeiro, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	results, err := r.queries.ListMetasBarbeiroByMesAno(ctx, db.ListMetasBarbeiroByMesAnoParams{
		TenantID: tenantUUID,
		MesAno:   mesAno.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar metas por mês: %w", err)
	}

	return r.toDomainList(results)
}

// toDomain converte modelo sqlc para entidade de domínio.
func (r *MetaBarbeiroRepository) toDomain(model *db.MetasBarbeiro) (*entity.MetaBarbeiro, error) {
	mesAno, err := valueobject.NewMesAno(model.MesAno)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear mes_ano: %w", err)
	}

	meta := &entity.MetaBarbeiro{
		ID:                       pgUUIDToString(model.ID),
		TenantID: pgtypeToEntityUUID(model.TenantID),
		BarbeiroID:               pgUUIDToString(model.BarbeiroID),
		MesAno:                   mesAno,
		MetaServicosGerais:       numericToMoney(model.MetaServicosGerais),
		MetaServicosExtras:       numericToMoney(model.MetaServicosExtras),
		MetaProdutos:             numericToMoney(model.MetaProdutos),
		RealizadoServicosGerais:  valueobject.Zero(),
		RealizadoServicosExtras:  valueobject.Zero(),
		RealizadoProdutos:        valueobject.Zero(),
		PercentualServicosGerais: valueobject.ZeroPercent(),
		PercentualServicosExtras: valueobject.ZeroPercent(),
		PercentualProdutos:       valueobject.ZeroPercent(),
		CriadoEm:                 timestamptzToTime(model.CriadoEm),
		AtualizadoEm:             timestamptzToTime(model.AtualizadoEm),
	}

	return meta, nil
}

// toDomainList converte lista de modelos para lista de entidades.
func (r *MetaBarbeiroRepository) toDomainList(models []db.MetasBarbeiro) ([]*entity.MetaBarbeiro, error) {
	result := make([]*entity.MetaBarbeiro, 0, len(models))
	for i := range models {
		meta, err := r.toDomain(&models[i])
		if err != nil {
			return nil, err
		}
		result = append(result, meta)
	}
	return result, nil
}
