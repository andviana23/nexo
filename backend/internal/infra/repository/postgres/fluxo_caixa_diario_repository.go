package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/andviana23/barber-analytics-backend/internal/domain/entity"
	"github.com/andviana23/barber-analytics-backend/internal/domain/valueobject"
	db "github.com/andviana23/barber-analytics-backend/internal/infra/db/sqlc"
)

// FluxoCaixaDiarioRepository implementa port.FluxoCaixaDiarioRepository com sqlc.
type FluxoCaixaDiarioRepository struct {
	queries *db.Queries
}

// NewFluxoCaixaDiarioRepository cria uma instância do repositório.
func NewFluxoCaixaDiarioRepository(queries *db.Queries) *FluxoCaixaDiarioRepository {
	return &FluxoCaixaDiarioRepository{queries: queries}
}

// Create insere um novo fluxo diário.
func (r *FluxoCaixaDiarioRepository) Create(ctx context.Context, fluxo *entity.FluxoCaixaDiario) error {
	tenantID := entityUUIDToPgtype(fluxo.TenantID)

	result, err := r.queries.CreateFluxoCaixaDiario(ctx, db.CreateFluxoCaixaDiarioParams{
		TenantID:            tenantID,
		Data:                dateToDate(fluxo.Data),
		SaldoInicial:        moneyToNumeric(fluxo.SaldoInicial),
		SaldoFinal:          moneyToNumeric(fluxo.SaldoFinal),
		EntradasConfirmadas: moneyToNumeric(fluxo.EntradasConfirmadas),
		EntradasPrevistas:   moneyToNumeric(fluxo.EntradasPrevistas),
		SaidasPagas:         moneyToNumeric(fluxo.SaidasPagas),
		SaidasPrevistas:     moneyToNumeric(fluxo.SaidasPrevistas),
		ProcessadoEm:        timestampToTimestamptz(fluxo.ProcessadoEm),
	})
	if err != nil {
		return fmt.Errorf("erro ao criar fluxo diário: %w", err)
	}

	fluxo.ID = pgUUIDToString(result.ID)
	fluxo.TenantID = pgtypeToEntityUUID(result.TenantID)
	fluxo.Data = dateToTime(result.Data)
	fluxo.ProcessadoEm = timestamptzToTime(result.ProcessadoEm)
	fluxo.CriadoEm = timestamptzToTime(result.CriadoEm)
	fluxo.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)
	return nil
}

// FindByID busca um fluxo diário pelo ID.
func (r *FluxoCaixaDiarioRepository) FindByID(ctx context.Context, tenantID, id string) (*entity.FluxoCaixaDiario, error) {
	tenantPg := uuidStringToPgtype(tenantID)
	idPg := uuidStringToPgtype(id)

	result, err := r.queries.GetFluxoCaixaDiarioByID(ctx, db.GetFluxoCaixaDiarioByIDParams{
		ID:       idPg,
		TenantID: tenantPg,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar fluxo diário por id: %w", err)
	}

	return r.modelToFluxo(&result)
}

// FindByData busca um fluxo diário por data.
func (r *FluxoCaixaDiarioRepository) FindByData(ctx context.Context, tenantID string, data time.Time) (*entity.FluxoCaixaDiario, error) {
	tenantPg := uuidStringToPgtype(tenantID)

	result, err := r.queries.GetFluxoCaixaDiarioByData(ctx, db.GetFluxoCaixaDiarioByDataParams{
		TenantID: tenantPg,
		Data:     dateToDate(data),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar fluxo diário por data: %w", err)
	}

	return r.modelToFluxo(&result)
}

// Update atualiza um fluxo diário.
func (r *FluxoCaixaDiarioRepository) Update(ctx context.Context, fluxo *entity.FluxoCaixaDiario) error {
	tenantPg := entityUUIDToPgtype(fluxo.TenantID)
	idPg := uuidStringToPgtype(fluxo.ID)

	result, err := r.queries.UpdateFluxoCaixaDiario(ctx, db.UpdateFluxoCaixaDiarioParams{
		ID:                  idPg,
		TenantID:            tenantPg,
		SaldoInicial:        moneyToNumeric(fluxo.SaldoInicial),
		SaldoFinal:          moneyToNumeric(fluxo.SaldoFinal),
		EntradasConfirmadas: moneyToNumeric(fluxo.EntradasConfirmadas),
		EntradasPrevistas:   moneyToNumeric(fluxo.EntradasPrevistas),
		SaidasPagas:         moneyToNumeric(fluxo.SaidasPagas),
		SaidasPrevistas:     moneyToNumeric(fluxo.SaidasPrevistas),
		ProcessadoEm:        timestampToTimestamptz(fluxo.ProcessadoEm),
	})
	if err != nil {
		return fmt.Errorf("erro ao atualizar fluxo diário: %w", err)
	}

	fluxo.AtualizadoEm = timestamptzToTime(result.AtualizadoEm)
	return nil
}

// Delete remove um fluxo diário.
func (r *FluxoCaixaDiarioRepository) Delete(ctx context.Context, tenantID, id string) error {
	tenantPg := uuidStringToPgtype(tenantID)
	idPg := uuidStringToPgtype(id)

	if err := r.queries.DeleteFluxoCaixaDiario(ctx, db.DeleteFluxoCaixaDiarioParams{
		ID:       idPg,
		TenantID: tenantPg,
	}); err != nil {
		return fmt.Errorf("erro ao deletar fluxo diário: %w", err)
	}
	return nil
}

// ListByDateRange lista fluxos por período.
func (r *FluxoCaixaDiarioRepository) ListByDateRange(ctx context.Context, tenantID string, inicio, fim time.Time) ([]*entity.FluxoCaixaDiario, error) {
	tenantPg := uuidStringToPgtype(tenantID)

	results, err := r.queries.ListFluxoCaixaDiarioByPeriod(ctx, db.ListFluxoCaixaDiarioByPeriodParams{
		TenantID: tenantPg,
		Data:     dateToDate(inicio),
		Data_2:   dateToDate(fim),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar fluxo diário por período: %w", err)
	}

	fluxos := make([]*entity.FluxoCaixaDiario, 0, len(results))
	for i := range results {
		f, err := r.modelToFluxo(&results[i])
		if err != nil {
			return nil, fmt.Errorf("erro ao converter fluxo diário %d: %w", i, err)
		}
		fluxos = append(fluxos, f)
	}
	return fluxos, nil
}

// List retorna fluxos paginados (usado internamente se necessário).
func (r *FluxoCaixaDiarioRepository) List(ctx context.Context, tenantID string, page, pageSize int) ([]*entity.FluxoCaixaDiario, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	tenantPg := uuidStringToPgtype(tenantID)

	results, err := r.queries.ListFluxoCaixaDiarioByTenant(ctx, db.ListFluxoCaixaDiarioByTenantParams{
		TenantID: tenantPg,
		Limit:    int32(pageSize),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar fluxo diário: %w", err)
	}

	fluxos := make([]*entity.FluxoCaixaDiario, 0, len(results))
	for i := range results {
		f, err := r.modelToFluxo(&results[i])
		if err != nil {
			return nil, fmt.Errorf("erro ao converter fluxo diário %d: %w", i, err)
		}
		fluxos = append(fluxos, f)
	}
	return fluxos, nil
}

// SumEntradas soma entradas em um período.
func (r *FluxoCaixaDiarioRepository) SumEntradas(ctx context.Context, tenantID string, inicio, fim time.Time) (valueobject.Money, error) {
	tenantPg := uuidStringToPgtype(tenantID)

	result, err := r.queries.SumEntradasByPeriod(ctx, db.SumEntradasByPeriodParams{
		TenantID: tenantPg,
		Data:     dateToDate(inicio),
		Data_2:   dateToDate(fim),
	})
	if err != nil {
		return valueobject.Money{}, fmt.Errorf("erro ao somar entradas: %w", err)
	}

	dec, err := aggregateNumericToDecimal(result)
	if err != nil {
		return valueobject.Money{}, fmt.Errorf("erro ao converter soma de entradas: %w", err)
	}
	return valueobject.NewMoneyFromDecimal(dec), nil
}

// SumSaidas soma saídas em um período.
func (r *FluxoCaixaDiarioRepository) SumSaidas(ctx context.Context, tenantID string, inicio, fim time.Time) (valueobject.Money, error) {
	tenantPg := uuidStringToPgtype(tenantID)

	result, err := r.queries.SumSaidasByPeriod(ctx, db.SumSaidasByPeriodParams{
		TenantID: tenantPg,
		Data:     dateToDate(inicio),
		Data_2:   dateToDate(fim),
	})
	if err != nil {
		return valueobject.Money{}, fmt.Errorf("erro ao somar saídas: %w", err)
	}

	dec, err := aggregateNumericToDecimal(result)
	if err != nil {
		return valueobject.Money{}, fmt.Errorf("erro ao converter soma de saídas: %w", err)
	}
	return valueobject.NewMoneyFromDecimal(dec), nil
}

// modelToFluxo converte model sqlc para entidade.
func (r *FluxoCaixaDiarioRepository) modelToFluxo(m *db.FluxoCaixaDiario) (*entity.FluxoCaixaDiario, error) {
	return &entity.FluxoCaixaDiario{
		ID:                  pgUUIDToString(m.ID),
		TenantID:            pgtypeToEntityUUID(m.TenantID),
		Data:                dateToTime(m.Data),
		SaldoInicial:        numericToMoney(m.SaldoInicial),
		SaldoFinal:          numericToMoney(m.SaldoFinal),
		EntradasConfirmadas: numericToMoney(m.EntradasConfirmadas),
		EntradasPrevistas:   numericToMoney(m.EntradasPrevistas),
		SaidasPagas:         numericToMoney(m.SaidasPagas),
		SaidasPrevistas:     numericToMoney(m.SaidasPrevistas),
		ProcessadoEm:        timestamptzToTime(m.ProcessadoEm),
		CriadoEm:            timestamptzToTime(m.CriadoEm),
		AtualizadoEm:        timestamptzToTime(m.AtualizadoEm),
	}, nil
}
