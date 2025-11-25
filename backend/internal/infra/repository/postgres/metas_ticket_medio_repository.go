// Package postgres implementa os repositórios usando PostgreSQL e sqlc.
package postgres

import (
	"context"
	"fmt"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// MetasTicketMedioRepository implementa port.MetasTicketMedioRepository usando sqlc.
type MetasTicketMedioRepository struct {
	queries *db.Queries
}

// NewMetasTicketMedioRepository cria uma nova instância do repositório.
func NewMetasTicketMedioRepository(queries *db.Queries) *MetasTicketMedioRepository {
	return &MetasTicketMedioRepository{
		queries: queries,
	}
}

// Create persiste uma nova meta de ticket médio.
func (r *MetasTicketMedioRepository) Create(ctx context.Context, meta *entity.MetaTicketMedio) error {
	tenantUUID := uuidStringToPgtype(meta.TenantID)

	barbeiroUUID := pgtype.UUID{Valid: false}
	if meta.BarbeiroID != nil {
		barbeiroUUID = uuidStringToPgtype(*meta.BarbeiroID)
	}

	tipoStr := string(meta.Tipo)

	params := db.CreateMetaTicketMedioParams{
		TenantID:   tenantUUID,
		Tipo:       &tipoStr,
		BarbeiroID: barbeiroUUID,
		MesAno:     meta.MesAno.String(),
		MetaValor:  moneyToRawDecimal(meta.MetaValor),
	}

	result, err := r.queries.CreateMetaTicketMedio(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao criar meta ticket médio: %w", err)
	}

	meta.ID = pgUUIDToString(result.ID)
	meta.CriadoEm = timestamptzToTime(result.CriadoEm)
	meta.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)

	return nil
}

// FindByID busca uma meta por ID.
func (r *MetasTicketMedioRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.MetaTicketMedio, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	result, err := r.queries.GetMetaTicketMedioByID(ctx, db.GetMetaTicketMedioByIDParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar meta ticket médio: %w", err)
	}

	return r.toDomain(&result)
}

// FindGeralByMesAno busca meta geral de um mês.
func (r *MetasTicketMedioRepository) FindGeralByMesAno(ctx context.Context, tenantID string, mesAno valueobject.MesAno) (*entity.MetaTicketMedio, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	result, err := r.queries.GetMetaTicketMedioGeralByMesAno(ctx, db.GetMetaTicketMedioGeralByMesAnoParams{
		TenantID: tenantUUID,
		MesAno:   mesAno.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar meta geral: %w", err)
	}

	return r.toDomain(&result)
}

// FindBarbeiroByMesAno busca meta de um barbeiro em um mês.
func (r *MetasTicketMedioRepository) FindBarbeiroByMesAno(ctx context.Context, tenantID, barbeiroID string, mesAno valueobject.MesAno) (*entity.MetaTicketMedio, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	barbeiroUUID := uuidStringToPgtype(barbeiroID)

	result, err := r.queries.GetMetaTicketMedioBarbeiroByMesAno(ctx, db.GetMetaTicketMedioBarbeiroByMesAnoParams{
		TenantID:   tenantUUID,
		BarbeiroID: barbeiroUUID,
		MesAno:     mesAno.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar meta de barbeiro: %w", err)
	}

	return r.toDomain(&result)
}

// Update atualiza uma meta existente.
func (r *MetasTicketMedioRepository) Update(ctx context.Context, meta *entity.MetaTicketMedio) error {
	tenantUUID := uuidStringToPgtype(meta.TenantID)
	idUUID := uuidStringToPgtype(meta.ID)

	params := db.UpdateMetaTicketMedioParams{
		ID:        idUUID,
		TenantID:  tenantUUID,
		MetaValor: moneyToRawDecimal(meta.MetaValor),
	}

	result, err := r.queries.UpdateMetaTicketMedio(ctx, params)
	if err != nil {
		return fmt.Errorf("erro ao atualizar meta ticket médio: %w", err)
	}

	meta.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)
	return nil
}

// Delete remove uma meta.
func (r *MetasTicketMedioRepository) Delete(ctx context.Context, tenantID, id string) error {
	tenantUUID := uuidStringToPgtype(tenantID)
	idUUID := uuidStringToPgtype(id)

	err := r.queries.DeleteMetaTicketMedio(ctx, db.DeleteMetaTicketMedioParams{
		ID:       idUUID,
		TenantID: tenantUUID,
	})
	if err != nil {
		return fmt.Errorf("erro ao deletar meta ticket médio: %w", err)
	}

	return nil
}

// ListByMesAno lista todas as metas de um mês (gerais e por barbeiro).
func (r *MetasTicketMedioRepository) ListByMesAno(ctx context.Context, tenantID string, mesAno valueobject.MesAno) ([]*entity.MetaTicketMedio, error) {
	tenantUUID := uuidStringToPgtype(tenantID)

	results, err := r.queries.ListMetasTicketMedioByMesAno(ctx, db.ListMetasTicketMedioByMesAnoParams{
		TenantID: tenantUUID,
		MesAno:   mesAno.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar metas do mês: %w", err)
	}

	return r.toDomainList(results)
}

// ListByBarbeiro lista metas de ticket médio de um barbeiro específico.
func (r *MetasTicketMedioRepository) ListByBarbeiro(ctx context.Context, tenantID, barbeiroID string) ([]*entity.MetaTicketMedio, error) {
	tenantUUID := uuidStringToPgtype(tenantID)
	barbeiroUUID := uuidStringToPgtype(barbeiroID)

	results, err := r.queries.ListMetasTicketMedioByBarbeiro(ctx, db.ListMetasTicketMedioByBarbeiroParams{
		TenantID:   tenantUUID,
		BarbeiroID: barbeiroUUID,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar metas do barbeiro: %w", err)
	}

	return r.toDomainList(results)
}

// toDomain converte modelo sqlc para entidade de domínio.
func (r *MetasTicketMedioRepository) toDomain(model *db.MetasTicketMedio) (*entity.MetaTicketMedio, error) {
	mesAno, err := valueobject.NewMesAno(model.MesAno)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear mes_ano: %w", err)
	}

	var barbeiroID *string
	if model.BarbeiroID.Valid {
		bid := pgUUIDToString(model.BarbeiroID)
		barbeiroID = &bid
	}

	var tipo valueobject.TipoMetaTicket
	if model.Tipo != nil {
		tipo = valueobject.TipoMetaTicket(*model.Tipo)
	}

	meta := &entity.MetaTicketMedio{
		ID:                   pgUUIDToString(model.ID),
		TenantID:             pgUUIDToString(model.TenantID),
		Tipo:                 tipo,
		BarbeiroID:           barbeiroID,
		MesAno:               mesAno,
		MetaValor:            rawDecimalToMoney(model.MetaValor),
		TicketMedioRealizado: valueobject.Zero(),
		Percentual:           valueobject.ZeroPercent(),
		CriadoEm:             timestamptzToTime(model.CriadoEm),
		AtualizadoEm:         timestamptzToTime(model.AtualizadoEm),
	}

	return meta, nil
}

// toDomainList converte lista de modelos para lista de entidades.
func (r *MetasTicketMedioRepository) toDomainList(models []db.MetasTicketMedio) ([]*entity.MetaTicketMedio, error) {
	result := make([]*entity.MetaTicketMedio, 0, len(models))
	for i := range models {
		meta, err := r.toDomain(&models[i])
		if err != nil {
			return nil, err
		}
		result = append(result, meta)
	}
	return result, nil
}
